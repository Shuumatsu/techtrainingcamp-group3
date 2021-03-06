package tools

import (
	"fmt"
	"techtrainingcamp-group3/config"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	fmt.Println("start")
	p := NewPool()
	p.Start()
	time.Sleep(time.Second)
	sum := 0
	for i := 0; i < config.TotalAmount; i++ {
		re := p.Snatch()
		fmt.Printf("Get red envelope %d, it is %d yuan\n", re.Eid, re.Money)
		sum += re.Money
	}
	fmt.Printf("Total money %d\n", sum)
}
