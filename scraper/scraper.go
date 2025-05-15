package scraper

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Status struct {
	Code    string
	Message string
}

type Formatted struct {
	CreatedAt string `json:"createdat"`
}

type TrackingInfo struct {
	Description string      `json:"description"`
	Timestamp   string      `json:"timestamp"`
	Formatted   []Formatted `json:"formatted"`
}

type AWBData struct {
	ReceivedBy string         `json:"receivedBy"`
	Histories  []TrackingInfo `json:"histories"`
}

func ScrapeAWB(awb string) (*AWBData, error) {
	url := fmt.Sprintf("https://be-otten-test.netlify.app/%s", awb)

	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return nil, Response("408", "request timed out")
		}
		if opErr, ok := err.(*net.OpError); ok {
			return nil, Response("NETWORK_ERROR_CODE", fmt.Sprintf("network error: %v", opErr.Err))
		}
		return nil, Response("400", fmt.Sprintf("failed to perform HTTP request: %v", err))
	}
	defer resp.Body.Close()

	// Validasi status HTTP
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, Response("404", "AWB number not found")
		} else {
			return nil, Response(strconv.Itoa(resp.StatusCode), fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
		}
	}

	// Parse dokumen HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, Response("400", "failed to parse response HTML")
	}

	var history []TrackingInfo

	doc.Find(".history-details tr").Each(func(i int, s *goquery.Selection) {
		tds := s.Find("td")
		if tds.Length() >= 2 {
			desc := strings.TrimSpace(tds.Eq(1).Text())
			timestamp := strings.TrimSpace(tds.Eq(0).Text())
			if desc != "" && timestamp != "" {
				createdat := timestamp
				// Parsing dari string ke time.Time
				inputTimeStamp := timestamp
				t, err := time.Parse("02-01-2006 15:04", inputTimeStamp)
				if err == nil {
					// Format kembali jadi string
					createdat = t.Format("02 January 2006, 15:04 WIB")
				}

				history = append(history, TrackingInfo{
					Description: desc,
					Timestamp:   timestamp,
					Formatted: []Formatted{
						{
							CreatedAt: createdat,
						},
					},
				})

			}
		}
	})

	trConsignee := doc.Find(".shipper-details tr").Eq(1)
	consignee := trConsignee.Find("td").Eq(1).Text()

	if len(history) == 0 {
		return nil, Response("404", "AWB number not found")
	}

	return &AWBData{
		ReceivedBy: consignee,
		Histories:  history,
	}, nil
}

func (e *Status) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Fungsi bantu untuk membuat error
func Response(code, message string) *Status {
	return &Status{
		Code:    code,
		Message: message,
	}
}
