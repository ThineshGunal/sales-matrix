package main

import (
	"log"
	"net/http"
	"os"
	"salesmatrix/common"
	"salesmatrix/dbconfig"
	"salesmatrix/sales"

	"github.com/gorilla/mux"

	"time"
)

func main() {
	lFile, lErr := os.OpenFile("./log/logfile"+time.Now().Format("02012006.15.04.05.000000000")+".txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if lErr != nil {
		log.Panic("error occured in log file create", lErr.Error())
	}
	defer lFile.Close()
	log.SetOutput(lFile)

	log.Println("server started in " + common.AppName + " !!... ")

	dbconfig.GDB, lErr = dbconfig.DBConnect()
	if lErr != nil {
		log.Fatal("db connection failed", lErr)
	}

	router := mux.NewRouter()
	srv := &http.Server{
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 600 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      router,
		Addr:         ":29093",
	}
	//Automatic refresh sales data
	go AutoRefresh()

	//refresh sales data manaul
	router.HandleFunc("/refreshdata", sales.RefreshDataAPI).Methods("POST")

	//this api is used to fetch the total revenue
	router.HandleFunc("/revenue/total/{from-date}/{to-date}", sales.GetTotalRevenueAPI).Methods("GET")

	//this api is used to fetch the revenue by product
	router.HandleFunc("/revenue/product/{id}/{from-date}/{to-date}", sales.GetProductRevenueAPI).Methods("GET")

	log.Println("listening on port 29093...")

	if lErr = srv.ListenAndServe(); lErr != nil {
		log.Fatal("server failed to start: ", lErr)
	}

}

func init() {
	common.ReadTomlFile("./toml/config.toml")

}
func AutoRefresh() {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {

		log.Println("executing daily basis task...")

		lErr := sales.RefreshData(common.Auto)
		if lErr != nil {
			log.Println("error occured while auto refreshing...", lErr.Error())
		} else {
			log.Println("data refresh completed successfully")
		}
	}

}
