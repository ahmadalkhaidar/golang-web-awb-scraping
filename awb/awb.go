package awb

import (
	"encoding/json"
	"golang-web-scraping/scraper"
	"net/http"

	"github.com/gorilla/mux"
)

func GetAWBInfo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	awbNumber := params["awb_number"]

	result, err := scraper.ScrapeAWB(awbNumber)
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": map[string]interface{}{
				"code":    "40003",
				"message": "AWB number not found",
			},
			"data": nil,
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": map[string]interface{}{
			"code":    "00000",
			"message": "Delivery tracking detail fetched successfully",
		},
		"data": result,
	})
}
