package mg

import (
	"context"
	"techtrainingcamp-group3/logger"
	"techtrainingcamp-group3/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindUserByEID(eid models.EID) (*models.User, error) {
	var user *models.User
	collection := MG.Collection(models.User{}.CollectionName())
	filter := bson.D{{"wallet.envelope_list.envelope_id", eid}}
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		logger.Sugar.Errorw("FindUserByEID", "collection find one", err.Error())
		return nil, err
	}
	return user, nil
}

func FindUserByUID(uid models.UID) (*models.User, error) {
	var user *models.User
	collection := MG.Collection(models.User{}.CollectionName())
	filter := bson.D{{"uid", uid}}
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		logger.Sugar.Errorw("FindUserByUID", "collection find one", err.Error())
		return nil, err
	}
	return user, nil
}

func AddEnvelopeToUserByUID(uid models.UID, envelope models.Envelope) error {
	collection := MG.Collection(models.User{}.CollectionName())
	filter := bson.D{{"uid", uid}}
	update := bson.D{
		{"$push", bson.D{
			{"wallet.envelope_list", envelope},
		}},
	}
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		logger.Sugar.Errorw("AddEnvelopeToUserByUID", "update", err.Error())
		return err
	}
	return nil
}

func SetDefaultUserByUID(uid models.UID) (*models.User, error) {
	user, err := FindUserByUID(uid)
	if err == nil {
		return user, nil
	} else if err == mongo.ErrNoDocuments {
		collection := MG.Collection(models.User{}.CollectionName())
		newUser := models.User{
			Uid: uid,
			Wallet: models.WalletListData{
				Amount: 0,
				EnvelopeList: []models.Envelope{},
			},
		}
		_, ierr := collection.InsertOne(context.TODO(), newUser)
		if ierr != nil {
			return nil, ierr
		}
		return &newUser, nil
	} else {
		return nil, err
	}
}
