package dbconfig

import (
	"database/sql"
	"fmt"
	"log"
	"salesmatrix/common"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var GDB *gorm.DB

func DBConnect() (*gorm.DB, error) {

	conf := common.GlobalConfig
	var lSqlDB *sql.DB

	var lDialector gorm.Dialector
	lConnectDB := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.DB.Username,
		conf.DB.Password,
		conf.DB.Server,
		conf.DB.Port,
		conf.DB.Database,
	)
	lDialector = mysql.Open(lConnectDB)

	log.Println("lConnectDB =>", lConnectDB)

	lGormDb, lErr := gorm.Open(lDialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if lErr != nil {
		log.Println("DB01", lErr)
		return lGormDb, lErr
	}

	lSqlDB, lErr = lGormDb.DB()
	if lErr != nil {
		log.Println("DB02", lErr)
		return lGormDb, lErr

	}

	lSqlDB.SetMaxIdleConns(conf.DBConfig.MaxIdleConnection)

	lSqlDB.SetMaxOpenConns(conf.DBConfig.MaxConnection)

	lSqlDB.SetConnMaxLifetime(60 * time.Minute)

	return lGormDb, nil
}
