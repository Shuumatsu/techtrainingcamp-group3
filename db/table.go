package db

type User struct {
	Uid          uint64 `gorm:"uid"`
	EnvelopeList string `gorm:"envelope_list"`
	Amount       uint64 `gorm:"amount"`
}

func (User) TableName() string {
	return "user"
}

type Envelope struct {
	EnvelopeId uint64 `gorm:"envelope_id" json:"envelope_id"`
	Opened     bool   `gorm:"opened" json:"opened"`
	Value      uint64 `gorm:"value" json:"value,omitempty"`
	SnatchTime uint64 `gorm:"snatch_time" json:"snatch_time"`
}

func (Envelope) TableName() string {
	return "envelope"
}
