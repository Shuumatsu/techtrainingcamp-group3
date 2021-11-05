package rds

import (
	"techtrainingcamp-group3/logger"
	"testing"
)

func TestClient(t *testing.T) {
	logger.Sugar.Debugw("test", "redis config:", DB)
	DB.Set("1", 100, 0)
	DB.Set("2", 156451, 0)
	val1, err := DB.Get("1").Int()
	if err != nil {
		t.Error(err)
	}
	val2, err := DB.Get("2").Int()
	if err != nil {
		t.Error(err)
	}
	if val1 != 100 || val2 != 156451 {
		t.Errorf("the value is error")
	}
}
