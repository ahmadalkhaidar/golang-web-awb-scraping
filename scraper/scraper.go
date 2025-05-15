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

type TrackingInfo struct {
	Description string `json:"description"`
	Timestamp   string `json:"timestamp"`
}

type AWBData struct {
	WaybillNumber string         `json:"waybill_number"`
	History       []TrackingInfo `json:"history"`
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
			desc := strings.TrimSpace(tds.Eq(0).Text())
			timestamp := strings.TrimSpace(tds.Eq(1).Text())

			if desc != "" && timestamp != "" {
				history = append(history, TrackingInfo{
					Description: desc,
					Timestamp:   timestamp,
				})
			}
		}
	})

	if len(history) == 0 {
		return nil, Response("404", "AWB number not found")
	}

	return &AWBData{
		WaybillNumber: awb,
		History:       history,
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
