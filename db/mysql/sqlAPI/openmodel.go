package sqlAPI

import (
	"time"
)

type Envelope struct {
	Envelope_id   uint64    `gorm:"envelope_id;primaryKey"`
	User_id       uint64    `gorm:"user_id"`
	Open_stat     bool      `gorm:"open_stat"`
	Value         int       `gorm:"value"`
	Snatched_time time.Time `gorm:"snatched_time"`
}

func (Envelope) TableName() string {
	return "envelope"
}

type OpenUser struct {
	User_id uint64 `gorm:"user_id;primaryKey"`
	Amount  uint64 `gorm:"amount"`
}

func (OpenUser) TableName() string {
	return "user"
}
