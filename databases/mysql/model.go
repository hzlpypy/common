package mysql

import (
	"errors"
	"gorm.io/gorm/logger"
	"time"
)

type Op string

const (
	Op_CREATE Op = "create"
	Op_UPDATE Op = "update"
	Op_DELETE Op = "delete"
)

func (o Op) Validate(opInter *OpInter) error {
	if o == Op_CREATE {
		return nil
	}
	if len(opInter.Where) == 0 {
		return errors.New("update where is nil")
	}
	return nil
}

type OpInter struct {
	Op     Op
	Object interface{}
	Where  string
}

type Config struct {
	Username string
	Password string
	DBName string
	Host     string
	Port     int
	DbName   string
	Charset  string
	NETWORK  string
	Debug bool
	// conn pool
	ConnMaxLifetime                          time.Duration
	MaxIdleConns                             int
	MaxOpenConns                             int
	LogLevel                                 logger.LogLevel
	DisableForeignKeyConstraintWhenMigrating bool
}