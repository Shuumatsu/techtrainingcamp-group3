package sql

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"runtime"
	"strconv"
	"sync"
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
	// create table `user` and `envelope`
	if !DB.Migrator().HasTable(&dbmodels.User{}) {
		err = DB.AutoMigrate(&dbmodels.User{})
		if err != nil {
			panic(err)
		}
		DB.Logger = DB.Logger.LogMode(0)
		err = RegisterDefaultUser(config.UserAmount)
		DB.Logger = DB.Logger.LogMode(2)
		if err != nil {
			logger.Sugar.Debugw(`mysql table "user" create`, "error", err)
			panic(err)
		}
		logger.Sugar.Debugw(`mysql table "user" create success`)
	}
	if !DB.Migrator().HasTable(&dbmodels.Envelope{}) {
		err = DB.AutoMigrate(&dbmodels.Envelope{})
		if err != nil {
			panic(err)
		}
		logger.Sugar.Debugw(`mysql table "envelope" create success`)
	}
}

func RegisterDefaultUser(n uint64) error {
	const step uint64 = 16383 // m * n < 65535, 此为user一次性插入的最大数目。
	bar := progressbar.Default(int64(n), "register user")
	doRegister := func(lo, hi uint64) {
		users := make([]dbmodels.User, hi-lo+1)
		for lo <= hi {
			users[hi-lo].Uid = dbmodels.UID(lo)
			lo++
			bar.Add(1)
		}
		err := DB.Table(dbmodels.User{}.TableName()).Create(users).Error
		if err != nil {
			logger.Sugar.Errorw("register user", "error", err)
		}
	}
	var i uint64 = 1
	if n > step {
		var wg sync.WaitGroup
		ch := make(chan struct{}, runtime.NumCPU())
		for j := 0; j < cap(ch); j++ {
			ch <- struct{}{}
		}
		logger.Sugar.Debugw("register user start", "cpuNum", cap(ch))
		for i = 1; i <= n-step; i += step {
			wg.Add(1)
			go func(i uint64, ch chan struct{}) {
				<-ch
				doRegister(i, i+step-1)
				wg.Done()
				ch <- struct{}{}
			}(i, ch)
		}
		wg.Wait()
	}
	doRegister(i, n)
	logger.Sugar.Debugw("register default user success", "user num", n)
	return nil
}
