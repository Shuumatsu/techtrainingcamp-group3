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

func CreateSqlTestData() {
	for i := 1; i <= 10; i++ {
		_, err := FindOrCreateUserByUID(dbmodels.User{Uid: dbmodels.UID(i)})
		if err != nil {
			logger.Sugar.Fatalw("TestFindOrCreateUserByUID", "error", err)
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
				logger.Sugar.Fatalw("TestAddEnvelopeToUserByUID", "error", err)
			}
		}
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

// func TestMain(t *testing.M) {
// 	rand.Seed(time.Now().Unix())
// 	DeleteSqlTestData()
// 	t.Run()
// }
//
// func TestFindOrCreateUserByUID(t *testing.T) {
// 	user, err := FindOrCreateUserByUID(dbmodels.User{
// 		Uid:          1,
// 		Amount:       5880,
// 		EnvelopeList: "",
// 	})
// 	if err != nil {
// 		t.Fatal("TestFindOrCreateUserByUID", "error", err)
// 	}
// 	fmt.Printf("user: %v\n", user)
// 	DeleteSqlTestData()
// }
//
// func TestFindUserByUID(t *testing.T) {
// 	CreateSqlTestData()
// 	var uid dbmodels.UID = 1
// 	user, err := FindUserByUID(uid)
// 	if err != nil && err != Error.NotFound {
// 		t.Fatal("TestFindOrCreateUserByUID", "error", err)
// 	}
// 	if err == Error.NotFound {
// 		fmt.Println("not found user", uid)
// 	}
// 	fmt.Printf("user: %v\n", user)
// 	DeleteSqlTestData()
// }
//
// func TestFindEnvelopesByUID(t *testing.T) {
// 	CreateSqlTestData()
// 	envelopes, err := FindEnvelopesByUID(1)
// 	if err != nil {
// 		t.Fatal("TestFindEnvelopesByUID", "error", err)
// 	}
// 	for _, envelope := range envelopes {
// 		fmt.Println(envelope)
// 	}
// 	ShowMySqlTestData()
// 	DeleteSqlTestData()
// }
//
// func TestAddEnvelopeToUserByUID(t *testing.T) {
// 	for i := 1; i <= 10; i++ {
// 		_, err := FindOrCreateUserByUID(dbmodels.User{Uid: dbmodels.UID(i)})
// 		if err != nil {
// 			t.Fatal("TestFindOrCreateUserByUID", "error", err)
// 		}
// 		for j := 1; j <= 10; j++ {
// 			err := AddEnvelopeToUserByUID(dbmodels.UID(i), dbmodels.Envelope{
// 				EnvelopeId: dbmodels.EID((i-1)*10 + j),
// 				Uid:        dbmodels.UID(i),
// 				Opened:     false,
// 				Value:      uint64(rand.Int31n(1008688)),
// 				SnatchTime: time.Now().Unix(),
// 			})
// 			if err != nil {
// 				t.Fatal("TestAddEnvelopeToUserByUID", "error", err)
// 			}
// 		}
// 	}
// 	ShowMySqlTestData()
// 	DeleteSqlTestData()
// }
//
// func TestFindEnvelopeByEID(t *testing.T) {
// 	CreateSqlTestData()
// 	// exist
// 	envelope, err := FindEnvelopeByEID(dbmodels.EID(rand.Intn(100)))
// 	if err != nil {
// 		if err == Error.ErrorParam {
// 			fmt.Println(Error.ErrorParam)
// 		}
// 		t.Fatal(err)
// 	}
// 	fmt.Println(envelope)
// 	// not exist
// 	envelope, err = FindEnvelopeByEID(dbmodels.EID(rand.Intn(100) + 101))
// 	if err != Error.NotFound {
// 		t.Fatal(err)
// 	}
// 	DeleteSqlTestData()
// }
//
// func TestOpenEnvelopeByEID(t *testing.T) {
// 	CreateSqlTestData()
// 	// exist
// 	_,err := OpenEnvelope(75, 8)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// not exist
// 	_,err = OpenEnvelope(414, 4)
// 	if err != Error.NotFound {
// 		t.Fatal(err)
// 	}
// 	// error owner
// 	_,err = OpenEnvelope(40, 1)
// 	if err != dbmodels.Error.ErrorEnvelopeOwner {
// 		t.Fatal(err)
// 	}
// 	// already open
// 	_,err = OpenEnvelope(75, 8)
// 	if err != dbmodels.Error.EnvelopeAlreadyOpen {
// 		t.Fatal(err)
// 	}
// 	ShowMySqlTestData()
// 	DeleteSqlTestData()
// }
func showUser() {
	var users []dbmodels.User
	if err := sql.DB.Table(
		dbmodels.User{}.TableName()).Scan(
		&users).Error; err != nil {
		logger.Sugar.Fatalw("ShowMySqlTestData", "error", err)
	}
	for i := 0; i < len(users); i++ {
		fmt.Println(users[i])
	}
	fmt.Printf("There are %v user in the table %v\n", len(users), dbmodels.User{}.TableName())
}
func showUserAmount() {
	var users []dbmodels.User
	if err := sql.DB.Table(
		dbmodels.User{}.TableName()).Scan(
		&users).Error; err != nil {
		logger.Sugar.Fatalw("ShowMySqlTestData", "error", err)
	}
	fmt.Printf("There are %v user in the table %v\n", len(users), dbmodels.User{}.TableName())
}
func showEnvelope() {
	var envelopes []dbmodels.Envelope
	if err := sql.DB.Table(
		dbmodels.Envelope{}.TableName()).Scan(
		&envelopes).Error; err != nil {
		logger.Sugar.Fatalw("ShowMySqlTestData", "error", err)
	}
	for i := 0; i < len(envelopes); i++ {
		fmt.Println(envelopes[i])
	}
	fmt.Printf("There are %v envelope in the table %v\n", len(envelopes), dbmodels.Envelope{}.TableName())
}
func showEnvelopeAmount() {
	var envelopes []dbmodels.Envelope
	if err := sql.DB.Table(
		dbmodels.Envelope{}.TableName()).Scan(
		&envelopes).Error; err != nil {
		logger.Sugar.Fatalw("ShowMySqlTestData", "error", err)
	}
	fmt.Printf("There are %v envelope in the table %v\n", len(envelopes), dbmodels.Envelope{}.TableName())
}
func TestMysql(t *testing.T) {
	showEnvelopeAmount()
}
