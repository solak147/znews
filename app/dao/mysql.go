package dao

import (
	"database/sql"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	GormSession *gorm.DB
	DbSession   *sql.DB
)

func GormInit(dbConfig string) (*gorm.DB, error) {
	var err error
	GormSession, err = gorm.Open(mysql.Open(dbConfig), &gorm.Config{})
	return GormSession, err
}

func DbInit(dbConfig string) (*sql.DB, error) {
	var err error
	DbSession, err = sql.Open("mysql", dbConfig)
	return DbSession, err
}
