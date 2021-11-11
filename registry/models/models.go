package models

import (
	"encoding/json"
	"time"
)

type UID uint64

func (u UID) Key() string {
	return "user:" + u.String()
}
func (u UID) String() string {
	return int2str(uint64(u))
}

type EID uint64

func (e EID) Key() string {
	return "envelope:" + e.String()
}
func (e EID) String() string {
	return int2str(uint64(e))
}

type User struct {
	Uid          UID       `gorm:"column:uid; PRIMARY_KEY; uniqueIndex" json:"uid"`
	Amount       uint64    `gorm:"column:amount" json:"amount"`
	EnvelopeList string    `gorm:"column:envelope_list" json:"envelope_list"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (User) TableName() string {
	return "user"
}
func (u *User) MarshalBinary() (data []byte, err error) {
	return json.Marshal(u)
}
func (u *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}

type Envelope struct {
	EnvelopeId EID    `gorm:"column:envelope_id; PRIMARY_KEY; uniqueIndex" json:"envelope_id"`
	Uid        UID    `gorm:"column:uid; index:uid" json:"uid"`
	Opened     bool   `gorm:"column:opened" json:"opened"`
	Value      uint64 `gorm:"column:value" json:"value"`
	SnatchTime int64  `gorm:"column:snatch_time" json:"snatch_time"`
}

func (Envelope) TableName() string {
	return "envelope"
}
func (e *Envelope) MarshalBinary() (data []byte, err error) {
	return json.Marshal(e)
}
func (e *Envelope) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, e)
}
