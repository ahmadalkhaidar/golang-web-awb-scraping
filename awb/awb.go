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
		if appErr, ok := err.(*scraper.Status); ok {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status": map[string]interface{}{
					"code":    appErr.Code,
					"message": appErr.Message,
				},
				"data": nil,
			})
			return
		} else {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status": map[string]interface{}{
					"code":    "500",
					"message": err.Error(),
				},
				"data": nil,
			})
			return
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": map[string]interface{}{
			"code":    "060101",
			"message": "Delivery tracking detail fetched successfully",
		},
		"data": result,
	})
}
