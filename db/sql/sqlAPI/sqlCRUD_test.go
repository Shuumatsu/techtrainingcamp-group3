package sqlAPI

import (
	"fmt"
	"math/rand"
	"techtrainingcamp-group3/db/dbmodels"
	"techtrainingcamp-group3/db/sql"
	"techtrainingcamp-group3/logger"
	"testing"
	"time"
)

func ShowMySqlTestData() {
	var users []dbmodels.User
	if err := sql.DB.Table(
		dbmodels.User{}.TableName()).Scan(
		&users).Error; err != nil {
		logger.Sugar.Fatalw("ShowMySqlTestData", "error", err)
	}
	fmt.Printf("There are %v user in the table %v\n", len(users), dbmodels.User{}.TableName())
	for i := 0; i < len(users); i++ {
		fmt.Println(users[i])
	}
	var envelopes []dbmodels.Envelope
	if err := sql.DB.Table(
		dbmodels.Envelope{}.TableName()).Scan(
		&envelopes).Error; err != nil {
		logger.Sugar.Fatalw("ShowMySqlTestData", "error", err)
	}
	fmt.Printf("There are %v envelope in the table %v\n", len(envelopes), dbmodels.Envelope{}.TableName())
	for i := 0; i < len(envelopes); i++ {
		fmt.Println(envelopes[i])
	}
}

func DeleteSqlTestData() {
	err := sql.DB.Table(dbmodels.User{}.TableName()).Delete(dbmodels.User{}, "uid >= ?", 0).Error
	if err != nil {
		logger.Sugar.Fatalw("DeleteSqlTestData", "user error", err)
	}
	err = sql.DB.Table(dbmodels.Envelope{}.TableName()).Delete(dbmodels.Envelope{}, "envelope_id >= ?", 0).Error
	if err != nil {
		logger.Sugar.Fatalw("DeleteSqlTestData", "envelope error", err)
	}
}

func TestFindOrCreateUserByUID(t *testing.T) {
	user, err := FindOrCreateUserByUID(1, dbmodels.User{
		Uid:          1,
		Amount:       5880,
		EnvelopeList: ",1",
	})
	if err != nil {
		t.Fatal("TestFindOrCreateUserByUID", "error", err)
	}
	fmt.Printf("user: %v\n", user)
}

func TestFindUserByUID(t *testing.T) {
	var uid dbmodels.UID = 1
	user, err := FindUserByUID(uid)
	if err != nil && err != Error.NotFound {
		t.Fatal("TestFindOrCreateUserByUID", "error", err)
	}
	if err == Error.NotFound {
		fmt.Println("not found user", uid)
	}
	fmt.Printf("user: %v\n", user)
}

func TestAddEnvelopeToUserByUID(t *testing.T) {
	for i := 1; i <= 10; i++ {
		_, err := FindOrCreateUserByUID(dbmodels.UID(i), dbmodels.User{Uid: dbmodels.UID(i)})
		if err != nil {
			t.Fatal("TestFindOrCreateUserByUID", "error", err)
		}
		for j := 1; j <= 10; j++ {
			err := AddEnvelopeToUserByUID(dbmodels.UID(i), dbmodels.Envelope{
				EnvelopeId: dbmodels.EID((i-1)*10 + j),
				Uid:        dbmodels.UID(i),
				Opened:     false,
				Value:      uint64(rand.Int31n(1008688)),
				SnatchTime: time.Now().Unix(),
			})
			if err != nil {
				t.Fatal("TestAddEnvelopeToUserByUID", "error", err)
			}
		}
	}
	ShowMySqlTestData()
	DeleteSqlTestData()
}
