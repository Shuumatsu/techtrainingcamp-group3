package database

import (
	"gorm.io/gorm"
)

func GetUser(uid uint64, txs ...*gorm.DB) (user *User, err error) {
	if len(txs) > 0 {
		txs[0].Where(&User{UserId: uid}).Take(&user)
	} else {
		Client.Where(&User{UserId: uid}).Take(&user)
	}

	return
}
func SetUser(user *User, txs ...*gorm.DB) (err error) {
	if len(txs) > 0 {
		txs[0].Save(user)
	} else {
		Client.Save(user)
	}

	return
}
func UpdateUser(uid uint64, attributes map[string]interface{}, txs ...*gorm.DB) (err error) {
	if len(txs) > 0 {
		txs[0].Model(&User{UserId: uid}).Updates(attributes)
	} else {
		Client.Model(&User{UserId: uid}).Updates(attributes)
	}

	return
}

func GetEnvelope(eid uint64, txs ...*gorm.DB) (envelope *Envelope, err error) {
	if len(txs) > 0 {
		txs[0].Where(&Envelope{EnvelopeId: eid}).Take(&envelope)
	} else {
		Client.Where(&Envelope{EnvelopeId: eid}).Take(&envelope)
	}

	return
}
func SetEnvelope(envelope *Envelope, txs ...*gorm.DB) (err error) {
	if len(txs) > 0 {
		txs[0].Save(envelope)
	} else {
		Client.Save(envelope)
	}

	return
}
func UpdateEnvelope(eid uint64, attributes map[string]interface{}, txs ...*gorm.DB) (err error) {
	if len(txs) > 0 {
		txs[0].Model(&Envelope{EnvelopeId: eid}).Updates(attributes)
	} else {
		Client.Model(&Envelope{EnvelopeId: eid}).Updates(attributes)
	}

	return
}
