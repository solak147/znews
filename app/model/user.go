package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        int64     `gorm:"primary_key;auto_increment" json:"id"`
	Account   string    `gorm:"size:30;not null;unique" json:"account" form:"account"`
	Password  string    `gorm:"size:30;not null;" json:"password" form:"password"`
	Email     string    `gorm:"size:30;" json:"email"`
	Phone     string    `gorm:"size:20;" json:"phone"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type Login struct {
	Account  string `form:"account"  example:"kevin"`
	Password string `form:"password"  example:"123456"`
}
