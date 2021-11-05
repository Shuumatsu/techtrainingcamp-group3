package dbmodels

import (
	"encoding/json"
	"techtrainingcamp-group3/models"
	"time"
)

type UID uint64

func (u UID) String() string {
	return int2str(uint64(u))
}

type EID uint64

func (e EID) String() string {
	return int2str(uint64(e))
}

type User struct {
	Uid          UID       `gorm:"uid" json:"uid"`
	Amount       uint64    `gorm:"amount" json:"amount"`
	EnvelopeList string    `gorm:"envelope_list" json:"envelope_list"`
	UpdatedAt    time.Time `gorm:"updated_at" json:"updated_at"`
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
	EnvelopeId EID    `gorm:"envelope_id" json:"envelope_id"`
	Uid        UID    `gorm:"uid" json:"uid"`
	Opened     bool   `gorm:"opened" json:"opened"`
	Value      uint64 `gorm:"value" json:"value"`
	SnatchTime int64  `gorm:"snatch_time" json:"snatch_time"`
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
func (e Envelope) ToResqModel() models.Envelope {
	var envelope models.Envelope
	envelope.EnvelopeId = uint64(e.EnvelopeId)
	envelope.Opened = e.Opened
	if e.Opened != false {
		envelope.Value = e.Value
	}
	envelope.SnatchTime = e.SnatchTime
	return envelope
}

func int2str(num uint64) string {
	if num == 0 {
		return "0"
	}
	var ret []byte
	for num != 0 {
		ret = append(ret, byte(num%10)+'0')
		num /= 10
	}
	for i, j := 0, len(ret)-1; i < j; i, j = i+1, j-1 {
		ret[i], ret[j] = ret[j], ret[i]
	}
	return string(ret)
}
