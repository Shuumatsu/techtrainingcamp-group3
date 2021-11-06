package sql

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"techtrainingcamp-group3/config"
	"techtrainingcamp-group3/db/dbmodels"
	"techtrainingcamp-group3/logger"
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
		panic(err)
	}
	DB = db
	logger.Sugar.Debugw("mysql init", "mysql config", dsn)
	if !db.Migrator().HasTable(&dbmodels.User{}) {
		err = DB.AutoMigrate(&dbmodels.User{})
		if err != nil {
			panic(err)
		}
	}
	if !db.Migrator().HasTable(&dbmodels.Envelope{}) {
		err = DB.AutoMigrate(&dbmodels.Envelope{})
		if err != nil {
			panic(err)
		}
	}
}
