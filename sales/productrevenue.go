package sales

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"salesmatrix/common"
	"salesmatrix/model"
	"strings"

	"github.com/gorilla/mux"
)

func GetProductRevenueAPI(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Set("Access-Control-Allow-Credentials", "true")
	resp.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, credentials")

	log.Println("GetRevenue(+)")
	var lRevenueResp model.CommonRevenueResp
	lRevenueResp.Status = common.SuccessStatus
	lRevenueResp.StatusCode = http.StatusText(http.StatusOK)
	lRevenueResp.Message = "process successfully completed"

	if strings.EqualFold(req.Method, http.MethodGet) {

		lReq := mux.Vars(req)
		lFromDate := lReq["from-date"]
		lToDate := lReq["to-date"]
		lId := lReq["id"]
		lType := "product"

		lErr := ValidateRequest(lFromDate, lToDate, lType, lId)
		if lErr != nil {
			log.Println("GPR01", lErr)
			lRevenueResp.Status = common.ErrorStatus
			lRevenueResp.StatusCode = "GPR01"
			lRevenueResp.Message = lErr.Error()
			goto marshal
		}

		lRevenueResp.RevenueType = lType

		lRevenueResp.Revenue, lErr = GetRevenueByProduct(lFromDate, lToDate, lId)
		if lErr != nil {
			log.Println("GPR02", lErr)
			lRevenueResp.Status = common.ErrorStatus
			lRevenueResp.StatusCode = "GPR02"
			lRevenueResp.Message = lErr.Error()
			goto marshal
		}

	marshal:
		{
			lResponse, lErr := json.Marshal(lRevenueResp)
			if lErr != nil {
				log.Println("error occured at marshal", lErr)
				fmt.Fprint(resp, string(lResponse))
				return

			}
			fmt.Fprint(resp, string(lResponse))
			return
		}
	}
	lRevenueResp.Status = common.ErrorStatus
	lRevenueResp.StatusCode = http.StatusText(http.StatusMethodNotAllowed)
	lRevenueResp.Message = req.Method + " not allowed"

	log.Println("GetRevenue(-)")

}
