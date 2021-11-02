package db

import (
	"time"
)

type Envelope struct{
	Envelope_id 	uint64  `gorm:"envelope_id"`
	User_id     	uint64  `gorm:"user_id"`
	Open_stat   	bool    `gorm:"open_stat"`
	Value       	int     `gorm:"value"`
	Snatched_time   time.Time `gorm:"snatched_time"`
}

