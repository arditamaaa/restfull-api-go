package database

import (
	"fmt"
	"log"
	"simple-api-go/config"
	util "simple-api-go/util"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect(dbHost, dbName string) *gorm.DB {
	dbUserPass := config.DBUsername
	if config.DBPassword != "" {
		dbUserPass = dbUserPass + ":" + config.DBPassword
	}
	dsn := fmt.Sprintf("%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUserPass, dbHost, config.DBPort, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// Logger:                 logger.Default.LogMode(logger.Info),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		TranslateError:         true,
	})

	if err != nil {
		util.Log.Errorf("Failed to connect to database: %+v", err)
	}

	logPrint := fmt.Sprintf("connect MySQL dbName %s", dbName)
	log.Println(logPrint)

	return db
}
