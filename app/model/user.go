package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID           int64     `gorm:"primary_key;auto_increment" json:"id"`
	Account      string    `gorm:"size:20;not null;unique" json:"account"`
	Password     string    `gorm:"size:30;not null;" json:"password"`
	Name         string    `gorm:"size:20;not null;"  json:"name" `
	Email        string    `gorm:"size:20;" json:"email"`
	Phone        string    `gorm:"size:20;" json:"phone"`
	Zipcode      string    `gorm:"size:3;"  json:"zipcode"`
	Introduction string    `gorm:"size:200;"  json:"introduction"`
	CreatedAt    time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type Login struct {
	Account  string `form:"account"  example:"kevin@gmail.com"`
	Password string `form:"password"  example:"123456"`
}

type RegisterStep1 struct {
	Email string `form:"email"  example:"kevin@gmail.com"`
}

type RegisterStep3 struct {
	Account      string `form:"account"  example:"kevin@gmail.com"`
	Password     string `form:"password"  example:"123456"`
	Name         string `form:"name"  example:"桐谷和人"`
	Phone        string `form:"phone"  example:"0999999999"`
	Email        string `form:"email"  example:"kevin@gmail.com"`
	Zipcode      string `form:"zipcode"  example:"200"`
	Introduction string `form:"introduction"  example:"我有 8000 名部下"`
}
