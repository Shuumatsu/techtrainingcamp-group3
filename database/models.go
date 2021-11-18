package database

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"techtrainingcamp-group3/pkg/models"
	"time"
)

type User struct {
	UserId       uint64    `gorm:"column:uid; PRIMARY_KEY; uniqueIndex" json:"uid"`
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

func ParseEnvelopeList(envelopeList string) ([]uint64, error) {
	ret := []uint64{}
	for _, token := range strings.Split(envelopeList, ",") {
		if len(token) == 0 {
			continue
		}
		eid, err := strconv.ParseUint(token, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invaild: the envelope id can not change to number")
		}
		ret = append(ret, eid)
	}
	return ret, nil
}

type Envelope struct {
	EnvelopeId uint64 `gorm:"column:envelope_id; PRIMARY_KEY; uniqueIndex" json:"envelope_id"`
	UserId     uint64 `gorm:"column:uid; index:uid" json:"uid"`
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
