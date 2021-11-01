package db

import (
	"log"
	"testing"
	"time"
)
type envelope struct {
	EnvelopeId uint64    `gorm:"eid"`
	Value      uint64    `gorm:"value"`
	Opened     bool      `gorm:"opened"`
	SnatchTime time.Time `gorm:"snatched_time"`
}

func (envelope) TableName() string {
	return "envelope_list"
}
func TestDB(t *testing.T) {
	if DB == nil {
		t.Errorf("db is nil")
	}
	eid := 0
	conditons := map[string]interface{}{
		"eid": eid,
	}
	var envelopes []*envelope
	if err := DB.Table(envelope{}.TableName()).
		Where(conditons).
		Find(&envelopes).Error; err != nil {
			t.Errorf(err.Error())
	}
	if envelopes == nil {
		t.Errorf("the eid %v is not exist", eid)
	}
	for _, evlp := range envelopes {
		log.Println(*evlp)
	}
}
