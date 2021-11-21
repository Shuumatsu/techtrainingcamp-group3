package bloomfilter

import (
	"github.com/go-redis/redis/v8"
	"techtrainingcamp-group3/config"
	"techtrainingcamp-group3/pkg/db/dbmodels"
	"techtrainingcamp-group3/pkg/db/rds"
	"techtrainingcamp-group3/pkg/logger"
)

const LargestPrime32 = (1<<31) -1

type HashF func(uint64) int64
type RedisBloom struct {
	Cli     *redis.Client
	Key      string
	HashFuncs []HashF
}

var badUserFilter      bool
var badEnvelopeFilter  bool
var userRedisFilter *RedisBloom
var envelopeRedisFilter *RedisBloom
var userLimiterRedisFilter *RedisBloom

func NewFuncArrUID() []HashF{
	m := make([]HashF,0)
	var f HashF
	f = LargestPrime32Module
	m = append(m,f)
	return m
}

func NewFuncArrSnowID() []HashF{
	m := make([]HashF,0)
	var f HashF
	f = LargestPrime32Module
	m = append(m,f)
	f = LargestPrime32Divide
	m = append(m,f)
	return m
}

//num % LargestPrime32
func LargestPrime32Module(num uint64)  int64{
	return int64(num % LargestPrime32)
}

//num / LargestPrime32
func LargestPrime32Divide(num uint64) int64{
	return int64(num / LargestPrime32)
}

func NewRedisBloom(cli *redis.Client,key string,f func() []HashF) *RedisBloom{
	return &RedisBloom{Cli: cli,Key:key,HashFuncs: f()}
}

//Add the key to the bloom filter,key need to be uint64 type
func (b *RedisBloom)Add(id uint64) error{
	var err error
	for _,f := range b.HashFuncs{
		offset := f(id)
		err = b.Cli.SetBit(rds.Ctx,b.Key,offset,1).Err()
		if err != nil{
			logger.Sugar.Errorw("bloomAddError","id",id)
			return err
		}
	}
	return err
}


//Test if id exists on the bloom filter, if a redis error happened, return true
func (b *RedisBloom) Exist(id uint64) (bool,error){
	for _,f := range b.HashFuncs {
		offset := f(id)
		cmd := b.Cli.GetBit(rds.Ctx,b.Key,offset)
		 v,err :=cmd.Result();
		 if err != nil{
			return true,err
		}
		if v != 1 {
			return false,nil
		}
	}
	return true,nil
}

func init() {

	if config.Env.Bloomfilter == "true" {
		userRedisFilter = NewRedisBloom(rds.DB, "userBloomFilter",NewFuncArrUID)
		envelopeRedisFilter = NewRedisBloom(rds.DB, "envelopeBloomFilter",NewFuncArrSnowID)
		userLimiterRedisFilter = NewRedisBloom(rds.DB,"userLimiterBloomFilter",NewFuncArrUID)
		badUserFilter = false
		badEnvelopeFilter = false
	}
}

func RedisAddUser(uid dbmodels.UID) {
	if config.Env.Bloomfilter != "true" || badUserFilter{
		return
	}
	//Add user failed
	if err := userRedisFilter.Add(uint64(uid));err != nil{
		badUserFilter = true
		logger.Sugar.Errorw("badUserFilter")
	}
}
func RedisTestUser(uid dbmodels.UID) bool {
	if config.Env.Bloomfilter != "true" || badUserFilter {
		return true
	}
	ret,err := userRedisFilter.Exist(uint64(uid))
	if err != nil{
		badUserFilter = true
		logger.Sugar.Errorw("badUserFilter")
	}
	return ret
}

func RedisAddEnvelope(eid dbmodels.EID) {
	if config.Env.Bloomfilter != "true" || badEnvelopeFilter{
		return
	}
	if err := envelopeRedisFilter.Add(uint64(eid));err != nil{
		badEnvelopeFilter = true
		logger.Sugar.Errorw("badEnvelopeFilter")
	}
}
func RedisTestEnvelope(eid dbmodels.EID) bool {
	if config.Env.Bloomfilter != "true" || badEnvelopeFilter{
		return true
	}
	ret,err := envelopeRedisFilter.Exist(uint64(eid))
	if err != nil{
		badEnvelopeFilter = true
		logger.Sugar.Errorw("badEnvelopeFilter")
	}
	return ret
}

func RedisAddLimitUser(uid dbmodels.UID) {
	if config.Env.Bloomfilter != "true" || badUserFilter{
		return
	}
	//Add user failed
	if err := userLimiterRedisFilter.Add(uint64(uid));err != nil{
		badUserFilter = true
		logger.Sugar.Errorw("badUserLimiterFilter")
	}
}

func RedisTestLimitUser(uid dbmodels.UID) bool {
	if config.Env.Bloomfilter != "true" || badUserFilter {
		return true
	}
	ret,err := userLimiterRedisFilter.Exist(uint64(uid))
	if err != nil{
		badUserFilter = true
		logger.Sugar.Errorw("badUserFilter")
	}
	return ret
}