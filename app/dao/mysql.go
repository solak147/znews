package dao

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	SqlSession *gorm.DB
)

func Initialize(dbConfig string) (*gorm.DB, error) {
	var err error
	SqlSession, err = gorm.Open(mysql.Open(dbConfig), &gorm.Config{})
	return SqlSession, err
}
