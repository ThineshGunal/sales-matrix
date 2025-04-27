package sales

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"salesmatrix/common"
	"salesmatrix/dbconfig"
	"salesmatrix/model"
	"strconv"
	"strings"
	"time"

	"github.com/gocarina/gocsv"
)

func RefreshDataAPI(resp http.ResponseWriter, req *http.Request) {

	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Set("Access-Control-Allow-Credentials", "true")
	resp.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, credentials")

	log.Println("RefreshDataAPI(+)")

	var lRefreshResp model.RefreshDataResp

	lRefreshResp.Status = common.SuccessStatus
	lRefreshResp.StatusCode = http.StatusText(http.StatusOK)
	lRefreshResp.Message = "process successfully completed"

	if strings.EqualFold(req.Method, "POST") {

		lErr := RefreshData(common.Manual)
		if lErr != nil {
			log.Println("RDA01", lErr)
			lRefreshResp.Status = common.ErrorStatus
			lRefreshResp.StatusCode = "RDA01"
			lRefreshResp.Message = lErr.Error()
			goto marshal

		}

	marshal:
		{
			lResponse, lErr := json.Marshal(lRefreshResp)
			if lErr != nil {
				log.Println("error occured at marshal", lErr)
				fmt.Fprint(resp, string(lResponse))
				return

			}
			fmt.Fprint(resp, string(lResponse))
			return
		}

	}
	lRefreshResp.Status = common.ErrorStatus
	lRefreshResp.StatusCode = http.StatusText(http.StatusMethodNotAllowed)
	lRefreshResp.Message = req.Method + " not allowed"

}
func RefreshData(pType string) error {
	log.Println("RefreshData(+)")

	lErr := InsertSalesRefreshLog(common.Pending, pType)
	if lErr != nil {
		log.Println("RD01 ", lErr)
		return lErr
	}

	lFile, lErr := os.Open("./csv/sample.csv")
	if lErr != nil {
		log.Println("RD02 ", lErr)
		return lErr
	}

	defer lFile.Close()

	var lSalesRequest []*model.SalesData

	lErr = CSVFileReader(lFile, &lSalesRequest)
	if lErr != nil {
		log.Println("RD03", lErr)
		return lErr
	}

	for _, lSalesData := range lSalesRequest {

		var lSalesProcess model.SalesProcessData

		lErr = ProcessCustomerData(&lSalesProcess, lSalesData)
		if lErr != nil {
			log.Println("RD04", lErr)
			return lErr
		}
		lErr = ProcessProductData(&lSalesProcess, lSalesData)
		if lErr != nil {
			log.Println("RD05", lErr)
			return lErr
		}
		lErr = ProcessOrdersData(&lSalesProcess, lSalesData)
		if lErr != nil {
			log.Println("RD06", lErr)
			return lErr
		}
		lErr = ProcessOrderItemsData(&lSalesProcess, lSalesData)
		if lErr != nil {
			log.Println("RD07", lErr)
			return lErr
		}

	}

	lErr = UpdateRefreshLog(common.Completed)
	if lErr != nil {
		log.Println("RD07", lErr)
		return lErr
	}

	log.Println("RefreshData(-)")

	return nil
}
func ProcessSalesData(pSalesData <-chan *model.SalesData, pErrChan chan<- error) {

}

func UpdateRefreshLog(pSatus string) error {
	log.Println("UpdateRefreshLog (+)")
	var latestLogEntry model.SalesRefreshLog
	if lErr := dbconfig.GDB.Order("created_date desc").First(&latestLogEntry).Error; lErr != nil {
		log.Println("UR01 ", lErr)
		return lErr
	}

	latestLogEntry.Status = pSatus
	latestLogEntry.UpdatedDate = time.Now()

	if lErr := dbconfig.GDB.Model(&latestLogEntry).Updates(model.SalesRefreshLog{
		Status:      latestLogEntry.Status,
		UpdatedDate: latestLogEntry.UpdatedDate,
	}).Error; lErr != nil {
		log.Println("UR02 ", lErr)
		return lErr
	}

	log.Println("UpdateRefreshLog (-)")

	return nil
}
func CSVFileReader(file *os.File, Struct interface{}) error {
	log.Println("CSVFileReader(+)")
	lErr := gocsv.UnmarshalFile(file, Struct)

	if lErr != nil {
		log.Println("CFR01", lErr)
		return lErr
	}

	log.Println("CSVFileReader(-)")

	return nil
}

func ProcessCustomerData(lSalesProcess *model.SalesProcessData, lSalesData *model.SalesData) error {
	log.Println("ProcessCustomerData(+)")
	lSalesProcess.Customer.CustomerId = lSalesData.CustomerID
	lSalesProcess.Customer.Name = lSalesData.CustomerName
	lSalesProcess.Customer.Email = lSalesData.CustomerEmail
	lSalesProcess.Customer.Address = lSalesData.CustomerAddress

	lErr := dbconfig.GDB.Where("customer_id=?", lSalesProcess.Customer.CustomerId).FirstOrCreate(&lSalesProcess.Customer).Error
	if lErr != nil {
		log.Println("PCD01", lErr)
		return lErr
	}
	log.Println("ProcessCustomerData(-)")

	return nil
}

func ProcessProductData(lSalesProcess *model.SalesProcessData, lSalesData *model.SalesData) error {
	log.Println("ProcessProductData(+)")

	lSalesProcess.Product.ProductId = lSalesData.ProductID
	lSalesProcess.Product.ProductName = lSalesData.ProductName
	lSalesProcess.Product.Category = lSalesData.Category

	lErr := dbconfig.GDB.Where("product_id=?", lSalesProcess.Product.ProductId).FirstOrCreate(&lSalesProcess.Product).Error
	if lErr != nil {
		log.Println("PPD01", lErr)
		return lErr
	}
	log.Println("ProcessProductData(-)")

	return nil
}

func ProcessOrdersData(lSalesProcess *model.SalesProcessData, lSalesData *model.SalesData) error {
	log.Println("ProcessOrdersData(+)")

	var lErr error
	lSalesProcess.Orders.OrderId = lSalesData.OrderID
	lSalesProcess.Orders.CustomerId = lSalesData.CustomerID

	lSalesProcess.Orders.DateofSale, lErr = time.Parse("2006-01-02", lSalesData.DateOfSale)
	if lErr != nil {
		log.Println("POD01", lErr.Error())
		return lErr
	}

	lSalesProcess.Orders.PaymentMethod = lSalesData.PaymentMethod
	lSalesProcess.Orders.Region = lSalesData.Region

	lErr = dbconfig.GDB.Where("order_id=?", lSalesProcess.Orders.OrderId).FirstOrCreate(&lSalesProcess.Orders).Error

	if lErr != nil {
		log.Println("POD02", lErr.Error())
		return lErr
	}
	log.Println("ProcessOrdersData(-)")

	return nil
}

func ProcessOrderItemsData(lSalesProcess *model.SalesProcessData, lSalesData *model.SalesData) error {
	log.Println("ProcessOrderItemsData(+)")

	var lErr error

	lSalesProcess.OrderItems.OrderId = lSalesData.OrderID
	lSalesProcess.OrderItems.ProductId = lSalesData.ProductID

	lSalesProcess.OrderItems.QuantitySold, lErr = strconv.Atoi(lSalesData.QuantitySold)
	if lErr != nil {
		log.Println("POID01", lErr.Error())
		return lErr
	}

	lSalesProcess.OrderItems.ShippingCost, lErr = strconv.ParseFloat(lSalesData.ShippingCost, 64)
	if lErr != nil {
		log.Println("POID02", lErr.Error())
		return lErr
	}

	lSalesProcess.OrderItems.UnitPrice, lErr = strconv.ParseFloat(lSalesData.UnitPrice, 64)
	if lErr != nil {
		log.Println("POID03", lErr.Error())
		return lErr
	}

	lSalesProcess.OrderItems.Discount, lErr = strconv.ParseFloat(lSalesData.Discount, 64)
	if lErr != nil {
		log.Println("POID04", lErr.Error())
		return lErr
	}

	lErr = dbconfig.GDB.Where("order_id=? and product_id=?", lSalesProcess.Orders.OrderId, lSalesProcess.Product.ProductId).FirstOrCreate(&lSalesProcess.OrderItems).Error
	if lErr != nil {
		log.Println("POID05", lErr)
		return lErr
	}
	log.Println("ProcessOrderItemsData(-)")

	return nil
}

func InsertSalesRefreshLog(pStatus, pType string) error {
	log.Println("InsertSalesRefreshLog(+)")

	lLogEntry := model.SalesRefreshLog{
		Status:      pStatus,
		Type:        pType,
		CreatedDate: time.Now(),
		UpdatedDate: time.Now(),
	}

	if lErr := dbconfig.GDB.Create(&lLogEntry).Error; lErr != nil {
		log.Println("ISR01", lErr)
		return lErr
	}
	log.Println("InsertSalesRefreshLog(-)")

	return nil
}
