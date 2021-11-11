package models

import (
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestInt2Str(t *testing.T) {
	rand.Seed(time.Now().Unix())
	for i := 0; i < 1e7; i++ {
		num := rand.Int31n(1e8)
		int2strAns := int2str(uint64(num))
		realAns := strconv.Itoa(int(num))
		if int2strAns != realAns {
			t.Fatal("int2str is error")
		}
	}
}
