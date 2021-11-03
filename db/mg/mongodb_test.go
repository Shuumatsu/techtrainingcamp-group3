package mg

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"techtrainingcamp-group3/models"
	"testing"
	"time"
)

type userTest struct {
	UID      int64 `bson:"uid"`
	CurCount int64 `bson:"cur_count"`
	models.WalletListData
}

func TestEnvelope(t *testing.T) {
	collection := MG.Collection("userTest")
	user1 := userTest{
		UID:      1,
		CurCount: 2,
		WalletListData: models.WalletListData{
			Amount: 10086,
			EnvelopeList: []models.Envelope{
				{
					EnvelopeId: 1,
					Opened:     false,
					Value:      6546,
					SnatchTime: time.Now().Unix(),
				},
				{
					EnvelopeId: 2,
					Opened:     false,
					Value:      56161,
					SnatchTime: time.Now().Unix(),
				},
			},
		},
	}
	user2 := userTest{
		UID:      2,
		CurCount: 1,
		WalletListData: models.WalletListData{
			Amount: 15446,
			EnvelopeList: []models.Envelope{
				{
					EnvelopeId: 3,
					Opened:     false,
					Value:      654456,
					SnatchTime: time.Now().Unix(),
				},
			},
		},
	}
	user3 := userTest{
		UID:      3,
		CurCount: 0,
		WalletListData: models.WalletListData{
			Amount:       10086,
			EnvelopeList: nil,
		},
	}
	users := []interface{}{user1, user2, user3}
	insertResult, err := collection.InsertMany(context.TODO(), users)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Inserted userTest: ", insertResult.InsertedIDs)
	newEnvelopeList := []models.Envelope{
		{
			EnvelopeId: 4,
			Opened:     false,
			Value:      754456,
			SnatchTime: time.Now().Unix(),
		},
		{
			EnvelopeId: 5,
			Opened:     false,
			Value:      653356,
			SnatchTime: time.Now().Unix(),
		},
	}
	filter := bson.D{{"uid", 3}}
	update := bson.D{
		{"$inc", bson.D{
			{"cur_count", 2},
		}},
		{"$set", bson.D{
			{"WalletListData.EnvelopeList", newEnvelopeList},
		}},
	}
	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Matched %v documents and updated %v documents.\n",
		updateResult.MatchedCount,
		updateResult.ModifiedCount)
	// Here's an array in which you can store the decoded documents
	var results []*userTest
	// Passing bson.D{{}} as the filter matches all documents in the collection
	cur, err := collection.Find(context.TODO(), bson.D{{}}, options.Find())
	if err != nil {
		t.Fatal(err)
	}
	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var elem userTest
		err := cur.Decode(&elem)
		if err != nil {
			t.Fatal(err)
		}
		results = append(results, &elem)
	}
	if err := cur.Err(); err != nil {
		t.Fatal(err)
	}
	// Close the cursor once finished
	cur.Close(context.TODO())
	fmt.Printf("Found multiple documents (array of pointers): %+v\n", results)
	for _, u := range results {
		fmt.Println(*u)
	}
	deleteResult, err := collection.DeleteMany(context.TODO(), bson.D{{}})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)
}

type trainer struct {
	Name string
	Age  int
	City string
}

func TestClient(t *testing.T) {
	collection := MG.Collection("trainers")
	ash := trainer{"Ash", 10, "Pallet Town"}
	misty := trainer{"Misty", 10, "Cerulean City"}
	brock := trainer{"Brock", 15, "Pewter City"}
	// test insert one
	insertResult, err := collection.InsertOne(context.TODO(), ash)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	// test insert many
	trainers := []interface{}{misty, brock}
	insertManyResult, err := collection.InsertMany(context.TODO(), trainers)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)
	// test update
	filter := bson.D{{"name", "Ash"}}
	update := bson.D{
		{"$inc", bson.D{
			{"age", 1},
		}},
	}
	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Matched %v documents and updated %v documents.\n",
		updateResult.MatchedCount,
		updateResult.ModifiedCount)
	updateResult, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Matched %v documents and updated %v documents.\n",
		updateResult.MatchedCount,
		updateResult.ModifiedCount)
	// test find one
	// create a value into which the result can be decoded
	var result trainer
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Found a single document: %+v\n", result)
	// test find many
	// Pass these options to the Find method
	findOptions := options.Find()
	findOptions.SetLimit(2)
	// Here's an array in which you can store the decoded documents
	var results []*trainer
	// Passing bson.D{{}} as the filter matches all documents in the collection
	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		t.Fatal(err)
	}
	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var elem trainer
		err := cur.Decode(&elem)
		if err != nil {
			t.Fatal(err)
		}
		results = append(results, &elem)
	}
	if err := cur.Err(); err != nil {
		t.Fatal(err)
	}
	// Close the cursor once finished
	cur.Close(context.TODO())
	fmt.Printf("Found multiple documents (array of pointers): %+v\n", results)
	// test delete
	deleteResult, err := collection.DeleteMany(context.TODO(), bson.D{{}})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)
}
