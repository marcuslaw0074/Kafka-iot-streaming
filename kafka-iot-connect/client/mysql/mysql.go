package mysql

import (
	"sync"

	"gorm.io/gorm"
)

type MysqlDB struct {
	*sync.RWMutex
	DB     *gorm.DB
	DbType string
	Name   string
	Path   string
	Tables map[string]string
}
