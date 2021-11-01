package db

import (
	"fmt"
	"log"
	"strconv"
	"strings"
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
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Env.DBUser,
		config.Env.DBPasswd,
		config.Env.DBHost,
		config.Env.DBPort,
		config.Env.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	DB = db
}

func ParseEnvelopeList(envelopeList string) ([]uint64, error) {
	envelopesID := make([]uint64, 0)
	for _, token := range strings.Split(envelopeList, ",") {
		if len(token) == 0 {
			continue
		}
		eid, err := strconv.Atoi(token)
		if err != nil {
			return nil, err
		}
		envelopesID = append(envelopesID, uint64(eid))
	}
	return envelopesID, nil
}
