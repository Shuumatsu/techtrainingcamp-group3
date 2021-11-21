package rds

import (
	"techtrainingcamp-group3/pkg/logger"
	"testing"
)

func TestClient(t *testing.T) {
	logger.Sugar.Debugw("test", "redis config:", DB)
	DB.Set(Ctx, "1", 100, 0)
	DB.Set(Ctx, "2", 156451, 0)
	val1, err := DB.Get(Ctx, "1").Int()
	if err != nil {
		t.Error(err)
	}
	val2, err := DB.Get(Ctx, "2").Int()
	if err != nil {
		t.Error(err)
	}
	if val1 != 100 || val2 != 156451 {
		t.Errorf("the value is error")
	}
}
