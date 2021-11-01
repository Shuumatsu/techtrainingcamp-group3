package tools

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"techtrainingcamp-group3/config"
	"time"

	"github.com/go-basic/uuid"
)

type UnopenedRedEnvelope struct {
	Money int
	Eid   string
}

type Pool struct {
	lanes []chan UnopenedRedEnvelope
	set   []reflect.SelectCase
}

var REPool Pool

func (p *Pool) Work(id int, money int, amount int) {
	fmt.Printf("Worker %d, money %d, amount %d\n", id, money, amount)
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
		p.lanes[id] <- UnopenedRedEnvelope{now, uuid.New()}
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
	channels := []chan UnopenedRedEnvelope{}
	for i := 0; i < config.PoolWorkerNUM; i++ {
		channels = append(channels, make(chan UnopenedRedEnvelope, config.PoolCapacity))
	}
	set := []reflect.SelectCase{}
	for _, ch := range channels {
		set = append(set, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		})
	}
	return Pool{channels, set}
}

func PoolInit() {
	REPool = NewPool()
	REPool.Start()
}
