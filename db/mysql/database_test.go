package mysql

import (
	"fmt"
	"techtrainingcamp-group3/config"
	"testing"
)

func TestDB(t *testing.T) {
	if DB == nil {
		t.Errorf("db is nil")
	}
	fmt.Println(config.Env)
}
