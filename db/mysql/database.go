package mysql

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"techtrainingcamp-group3/config"
)

type Database struct {
	*gorm.DB
}

var sql *gorm.DB

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
		panic(err)
	}
	sql = db
}