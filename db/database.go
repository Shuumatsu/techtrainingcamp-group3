package db

import (
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

var DB *gorm.DB

// Opening a database and save the reference to `Database` struct.
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := fmt.Sprintf("user:pass@tcp(%s:%s)/dbname?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	DB = db
}
