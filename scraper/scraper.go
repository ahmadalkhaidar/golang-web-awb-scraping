package scraper

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

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
		// Validasi timeout dan error jaringan lainnya
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return nil, errors.New("request timed out")
		}
		if opErr, ok := err.(*net.OpError); ok {
			return nil, fmt.Errorf("network error: %v", opErr.Err)
		}
		return nil, fmt.Errorf("failed to perform HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Validasi status HTTP
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse dokumen HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, errors.New("failed to parse response HTML")
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
		return nil, errors.New("no tracking information found")
	}

	return &AWBData{
		WaybillNumber: awb,
		History:       history,
	}, nil
}
