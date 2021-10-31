package db

import (
	"fmt"
	"log"
	"techtrainingcamp-group3/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

var DB *gorm.DB

// Opening a database and save the reference to `Database` struct.
func init() {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := fmt.Sprintf("user:pass@tcp(%s:%s)/dbname?charset=utf8mb4&parseTime=True&loc=Local", config.Env.DBHost, config.Env.DBPort)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	DB = db
}
