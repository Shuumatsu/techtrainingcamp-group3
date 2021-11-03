package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"math/rand"
	"techtrainingcamp-group3/models"
	"testing"
	"time"
)

var users []interface{}

func setup() {
	insertTestDate()
}
func teardown() {
	deleteTestData()
}
func TestMain(t *testing.M) {
	setup()
	t.Run()
	teardown()
}
func TestFindUserByEID(t *testing.T) {
	showTestData()
	user, err := FindUserByEID(1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(user)
}

func TestFindUserByUID(t *testing.T) {
	showTestData()
	user, err := FindUserByUID(1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(user)
}
func TestAddEnvelopeToUserByUID(t *testing.T) {
	showTestData()
	err := AddEnvelopeToUserByUID(1, models.Envelope{
		EnvelopeId: 646,
		Opened:     false,
		Value:      646461,
		SnatchTime: time.Now().Unix() + 10000,
	})
	if err != nil {
		t.Fatal(err)
	}
	showTestData()
}

func insertTestDate() {
	rand.Seed(time.Now().Unix())
	users = make([]interface{}, 0)
	users = append(users, models.User{
		Uid: 5,
		Wallet: models.WalletListData{
			Amount:       0,
			EnvelopeList: nil,
		},
	})
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
		users = append(users, user)
	}
	collection := MG.Collection(models.User{}.CollectionName())
	insertManyResult, err := collection.InsertMany(context.TODO(), users)
	if err != nil {
		log.Fatal("insert test data error:", err)
	}
	fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)
}
func showTestData() {
	collection := MG.Collection(models.User{}.CollectionName())
	// Here's an array in which you can store the decoded documents
	var results []*models.User
	// Passing bson.D{{}} as the filter matches all documents in the collection
	cur, err := collection.Find(context.TODO(), bson.D{{}}, options.Find())
	if err != nil {
		log.Fatal(err)
	}
	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var elem models.User
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, &elem)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	// Close the cursor once finished
	cur.Close(context.TODO())
	fmt.Printf("Found multiple documents (array of pointers): %+v\n", results)
	for _, u := range results {
		fmt.Println(*u)
	}
}
func deleteTestData() {
	collection := MG.Collection(models.User{}.CollectionName())
	deleteResult, err := collection.DeleteMany(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatal("delete test data error:", err)
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)
}
