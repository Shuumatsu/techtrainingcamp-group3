package sql

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strconv"
	"techtrainingcamp-group3/config"
	"techtrainingcamp-group3/db/dbmodels"
	"techtrainingcamp-group3/logger"
	"time"
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
	if !DB.Migrator().HasTable(&dbmodels.User{}) {
		err = DB.AutoMigrate(&dbmodels.User{})
		if err != nil {
			panic(err)
		}
	}
	if !DB.Migrator().HasTable(&dbmodels.Envelope{}) {
		err = DB.AutoMigrate(&dbmodels.Envelope{})
		if err != nil {
			panic(err)
		}
	}
	maxIdleConns, err := strconv.Atoi(config.Env.DBMaxIdleConns)
	if err != nil {
		logger.Sugar.Fatalw("mysql init", "maxIdleConns error", err)
	}
	maxOpenConns, err := strconv.Atoi(config.Env.DBMaxOpenConns)
	if err != nil {
		logger.Sugar.Fatalw("mysql init", "maxOpenConns error", err)
	}
	// 获取通用数据库对象 sql.DB ，然后使用其提供的功能
	sqlDB, err := DB.DB()
	// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	sqlDB.SetMaxIdleConns(maxIdleConns)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(maxOpenConns)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
}
