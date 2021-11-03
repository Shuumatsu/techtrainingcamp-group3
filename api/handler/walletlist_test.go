package handler

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"techtrainingcamp-group3/db/rds"
	"techtrainingcamp-group3/models"
	"testing"
	"time"
)

func InsertTestDataInRedis() {
	rand.Seed(time.Now().Unix())
	rds.RD.Set(strconv.Itoa(5), models.User{
		Uid: 5,
		Wallet: models.WalletListData{
			Amount:       7646,
			EnvelopeList: nil,
		},
	}, 0)
	e := 0
	for i := 0; i < 5; i++ {
		envelopes := make([]models.Envelope, 0)
		for j := 0; j < rand.Intn(6); j++ {
			envelopes = append(envelopes,
				models.Envelope{
					EnvelopeId: models.EID(e),
					Opened:     j%2 == 1,
					Value:      uint64(rand.Intn(1008688)),
					SnatchTime: time.Now().Unix() + int64(rand.Intn(1000)),
				})
			e++
		}
		user := models.User{
			Uid: models.UID(i),
			Wallet: models.WalletListData{
				Amount:       uint64(rand.Intn(1008688)),
				EnvelopeList: envelopes,
			},
		}
		_, err := rds.RD.Set(user.Uid.String(), &user, 0).Result()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func TestRedisGetUserByUID(t *testing.T) {
	InsertTestDataInRedis()
	user, err := GetUserByUID(1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(user)
}
