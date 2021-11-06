package tools

import (
	"github.com/godruoyi/go-snowflake"
	"math"
	"math/rand"
	"reflect"
	"techtrainingcamp-group3/config"
	"time"
)

type UnopenedRedEnvelope struct {
	Money int
	Eid   uint64
}

type Pool struct {
	lanes []chan UnopenedRedEnvelope
	set   []reflect.SelectCase
}

var REPool Pool

func (p *Pool) Work(id int, money int, amount int) {
	snowflake.SetMachineID(uint16(id))
	mean := float64(money) / float64(amount)
	StdDev := math.Min(float64(config.MaxMoney)-mean, mean-float64(config.MinMoney)) / 3
	restMoney := money
	for i := 0; i < amount; i++ {
		rightRange := restMoney - (amount-i-1)*config.MinMoney
		if rightRange > config.MaxMoney {
			rightRange = config.MaxMoney
		}
		now := int(rand.NormFloat64()*StdDev + mean)
		if now > rightRange {
			now = rightRange
		}
		if now < config.MinMoney {
			now = config.MinMoney
		}
		restMoney -= now
		eid := snowflake.ID()
		p.lanes[id] <- UnopenedRedEnvelope{now, eid}
	}
}

func (p *Pool) Snatch() UnopenedRedEnvelope {
	_, value, _ := reflect.Select(p.set)
	return value.Interface().(UnopenedRedEnvelope)
}

func (p *Pool) Start() {
	for i := 0; i < config.PoolWorkerNUM; i++ {
		if i != config.PoolWorkerNUM-1 {
			go p.Work(i, config.TotalMoney/config.PoolWorkerNUM, config.TotalAmount/config.PoolWorkerNUM)
		} else {
			go p.Work(i, config.TotalMoney-config.TotalMoney/config.PoolWorkerNUM*i, config.TotalAmount-config.TotalAmount/config.PoolWorkerNUM*i)
		}
	}
}

func NewPool() Pool {
	rand.Seed(time.Now().UnixNano())
	channels := make([]chan UnopenedRedEnvelope, 0)
	for i := 0; i < config.PoolWorkerNUM; i++ {
		channels = append(channels, make(chan UnopenedRedEnvelope, config.PoolCapacity))
	}
	set := make([]reflect.SelectCase, 0)
	for _, ch := range channels {
		set = append(set, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		})
	}
	return Pool{channels, set}
}

func init() {
	REPool = NewPool()
	REPool.Start()
}
