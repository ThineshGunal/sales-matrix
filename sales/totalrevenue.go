package sales

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"salesmatrix/common"
	"salesmatrix/dbconfig"
	"salesmatrix/model"

	"strings"

	"github.com/gorilla/mux"
)

func GetTotalRevenueAPI(resp http.ResponseWriter, req *http.Request) {
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
		lType := "total"

		lErr := ValidateRequest(lFromDate, lToDate, lType, "")
		if lErr != nil {
			log.Println("GR01", lErr)
			lRevenueResp.Status = common.ErrorStatus
			lRevenueResp.StatusCode = "GR01"
			lRevenueResp.Message = lErr.Error()
			goto marshal
		}

		lRevenueResp.RevenueType = lType

		lRevenueResp.Revenue, lErr = GetTotalRevenue(lFromDate, lToDate)
		if lErr != nil {
			log.Println("GR02", lErr)
			lRevenueResp.Status = common.ErrorStatus
			lRevenueResp.StatusCode = "GR02"
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
func GetTotalRevenue(pFromDate, pToDate string) (float64, error) {
	log.Println("GetTotalRevenue(+)")

	var lTotalRevenue sql.NullFloat64
	query := dbconfig.GDB.Table("order_items")

	if pFromDate != "" && pToDate != "" {
		query = query.Joins("JOIN orders ON orders.order_id = order_items.order_id").
			Where("orders.date_of_sale BETWEEN ? AND ?", pFromDate, pToDate)

		lRow := query.Select("SUM((unit_price - discount) * quantity_sold)").Row()
		if lErr := lRow.Scan(&lTotalRevenue); lErr != nil {
			log.Println("GTR01", lErr)
			return lTotalRevenue.Float64, lErr
		}
	}
	if !lTotalRevenue.Valid {
		log.Println("Total revenue is NULL")
		return lTotalRevenue.Float64, nil
	}

	log.Println("GetTotalRevenue(-)")
	return lTotalRevenue.Float64, nil

}
func GetRevenueByProduct(pFromDate, pToDate, pProductId string) (float64, error) {
	log.Println("GetRevenueByProduct(+)")

	var lRevenueProduct sql.NullFloat64
	lQuery := dbconfig.GDB.Table("order_items").Where("product_id = ?", pProductId)

	if pFromDate != "" && pToDate != "" {
		lQuery = lQuery.Joins("JOIN orders ON orders.order_id = order_items.order_id").
			Where("orders.date_of_sale BETWEEN ? AND ?", pFromDate, pToDate)

		lRow := lQuery.Select("SUM((unit_price - discount) * quantity_sold)").Row()
		if lErr := lRow.Scan(&lRevenueProduct); lErr != nil {
			log.Println("GRP01", lErr)
			return lRevenueProduct.Float64, lErr
		}
	}
	if !lRevenueProduct.Valid {
		log.Println("revenue by product is NULL")
		return lRevenueProduct.Float64, nil
	}

	log.Println("GetRevenueByProduct(-)")
	return lRevenueProduct.Float64, nil

}
func ValidateRequest(pFromDate, pToDate, pType, pID string) error {
	log.Println("ValidateRequest(+)")

	log.Println("pFromDate", pFromDate)
	log.Println("pToDate", pToDate)
	log.Println("pType", pType)

	if pFromDate == "" {
		return errors.New("please provide the fromdate")
	}
	if pToDate == "" {
		return errors.New("please provide the todate")
	}

	if pType == "" {
		return errors.New("please provide the process type")
	}
	if strings.EqualFold(pType, "product") && pID == "" {
		return errors.New("please provide the prduct id")

	}

	var lIsMatch bool
	lDatePattern := `^\d{4}-\d{2}-\d{2}$`

	re := regexp.MustCompile(lDatePattern)

	lIsMatch = re.MatchString(pFromDate)
	if !lIsMatch {
		return errors.New("please provide the valid fromdate in (yyyy-mm-dd) format")
	}
	lIsMatch = re.MatchString(pToDate)
	if !lIsMatch {
		return errors.New("please provide the valid todate in (yyyy-mm-dd) format")
	}

	log.Println("ValidateRequest(-)")
	return nil
}
