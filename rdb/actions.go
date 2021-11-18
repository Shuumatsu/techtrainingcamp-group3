package rdb

import (
	"fmt"
	"strconv"
	"techtrainingcamp-group3/database"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

type Namespace string

const (
	UserNamespace     Namespace = "user"
	EnvelopeNamespace Namespace = "envelope"
)

func GetUserKey(uid uint64) string {
	return fmt.Sprintf("%s:%s", UserNamespace, strconv.FormatUint(uid, 10))
}

func GetEnvelopeKey(eid uint64) string {
	return fmt.Sprintf("%s:%s", EnvelopeNamespace, strconv.FormatUint(eid, 10))
}

func SetUser(ctx context.Context, user *database.User, expiration time.Duration, txs ...*redis.Tx) (err error) {
	if len(txs) > 0 {
		_, err = txs[0].Set(ctx, GetUserKey(user.UserId), user, expiration).Result()
	} else {
		_, err = Client.Set(ctx, GetUserKey(user.UserId), user, expiration).Result()
	}

	return
}
func GetUser(ctx context.Context, uid uint64, txs ...*redis.Tx) (user *database.User, err error) {
	var resp string
	if len(txs) > 0 {
		resp, err = Client.Get(ctx, GetUserKey(uid)).Result()
	} else {
		resp, err = Client.Get(ctx, GetUserKey(uid)).Result()
	}

	if err == nil {
		user.UnmarshalBinary([]byte(resp))
	}
	return
}

func SetEnvelope(ctx context.Context, envelope *database.Envelope, expiration time.Duration, txs ...*redis.Tx) (err error) {
	if len(txs) > 0 {
		_, err = txs[0].Set(ctx, GetEnvelopeKey(envelope.EnvelopeId), envelope, expiration).Result()
	} else {
		_, err = Client.Set(ctx, GetEnvelopeKey(envelope.EnvelopeId), envelope, expiration).Result()
	}
	return
}
func GetEnvelope(ctx context.Context, eid uint64, txs ...*redis.Tx) (envelope *database.Envelope, err error) {
	var resp string
	if len(txs) > 0 {
		resp, err = txs[0].Get(ctx, GetEnvelopeKey(eid)).Result()
	} else {
		resp, err = Client.Get(ctx, GetEnvelopeKey(eid)).Result()
	}

	if err == nil {
		envelope.UnmarshalBinary([]byte(resp))
	}
	return
}
