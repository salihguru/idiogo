package db

import (
	"context"
	"fmt"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	Debug    bool
}

func NewPostgres(ctx context.Context, cnf PostgresConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cnf.Host, cnf.Port, cnf.User, cnf.Password, cnf.DBName, cnf.SSLMode)
	conf := &gorm.Config{}
	if cnf.Debug {
		conf.Logger = logger.Default.LogMode(logger.Info)
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), conf)
	if err != nil {
		return nil, err
	}
	return gormDB, nil
}
