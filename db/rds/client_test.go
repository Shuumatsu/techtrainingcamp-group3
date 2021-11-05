package rds

import (
	"techtrainingcamp-group3/logger"
	"testing"
)

func TestClient(t *testing.T) {
	logger.Sugar.Debugw("test", "redis config:", RD)
	RD.Set("1", 100, 0)
	RD.Set("2", 156451, 0)
	val1, err := RD.Get("1").Int()
	if err != nil {
		t.Error(err)
	}
	val2, err := RD.Get("2").Int()
	if err != nil {
		t.Error(err)
	}
	if val1 != 100 || val2 != 156451 {
		t.Errorf("the value is error")
	}
}
