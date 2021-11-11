package mysql

import (
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

var DB *gorm.DB
