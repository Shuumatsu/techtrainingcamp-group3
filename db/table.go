package db

import "time"

type Envelope struct {
	Envelope_id   uint64    `gorm:"envelope_id;primaryKey"`
	Uid           uint64    `gorm:"uid"`
	Open_stat     bool      `gorm:"open_stat"`
	Value         int       `gorm:"value"`
	Snatched_time time.Time `gorm:"snatched_time"`
}

func (Envelope) TableName() string {
	return "envelope"
}

type User struct {
	Uid          uint64 `gorm:"uid"`
	EnvelopeList string `gorm:"envelope_list"`
	Amount       uint64 `gorm:"amount"`
}

func (User) TableName() string {
	return "user"
}
