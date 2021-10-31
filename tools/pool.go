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

const WorkerNUM = 10

type UnopenedRedEnvelope struct {
	money int
	eid   string
}

type Pool struct {
	lanes []chan UnopenedRedEnvelope
	set   []reflect.SelectCase
}

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
	for i := 0; i < WorkerNUM; i++ {
		if i != WorkerNUM-1 {
			go p.Work(i, config.TotalMoney/WorkerNUM, config.TotalAmount/WorkerNUM)
		} else {
			go p.Work(i, config.TotalMoney-config.TotalMoney/WorkerNUM*i, config.TotalAmount-config.TotalAmount/WorkerNUM*i)
		}
	}
}

func NewPool() Pool {
	rand.Seed(time.Now().UnixNano())
	channels := []chan UnopenedRedEnvelope{}
	for i := 0; i < WorkerNUM; i++ {
		channels = append(channels, make(chan UnopenedRedEnvelope, 2))
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
