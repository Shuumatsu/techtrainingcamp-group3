package mysql

import (
	"fmt"
	"runtime"
	"sync"
	"techtrainingcamp-group3/registry/models"
	"time"

	"github.com/schollz/progressbar/v3"

	"gorm.io/driver/mysql"

	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

var DB *gorm.DB

func NewDatabase(user, password, host, port, name string, maxIdleConns, maxOpenConns int, userAmount uint64) *Database {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// 获取通用数据库对象 sql.DB ，然后使用其提供的功能
	mysql, err := DB.DB()
	if err != nil {
		panic(err)
	}
	// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	mysql.SetMaxIdleConns(maxIdleConns)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	mysql.SetMaxOpenConns(maxOpenConns)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	mysql.SetConnMaxLifetime(time.Hour)

	if !DB.Migrator().HasTable(&models.User{}) {
		if err := DB.AutoMigrate(&models.User{}); err != nil {
			panic(err)
		}
		DB.Logger = DB.Logger.LogMode(0)
		if err := RegisterDefaultUser(userAmount); err != nil {
			panic(err)
		}
		DB.Logger = DB.Logger.LogMode(2)
	}
	if !DB.Migrator().HasTable(&models.Envelope{}) {
		err = DB.AutoMigrate(&models.Envelope{})
		if err != nil {
			panic(err)
		}
	}

	return &Database{db}
}

func RegisterDefaultUser(n uint64) error {
	const step uint64 = 16383 // m * n < 65535, 此为user一次性插入的最大数目。
	bar := progressbar.Default(int64(n), "register user")
	doRegister := func(lo, hi uint64) {
		users := make([]models.User, hi-lo+1)
		for lo <= hi {
			users[hi-lo].Uid = models.UID(lo)
			lo++
			bar.Add(1)
		}
		err := DB.Table(models.User{}.TableName()).Create(users).Error
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
