package model

import (
	"time"
)

type Product struct {
	Name    string    `gorm:"primary_key;size:10;not null;" json:"name"`
	Price   string    `gorm:"size:4;not null;" json:"price"`
	ChiName    string    `gorm:"size:10;not null;" json:"chiName"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type ProductForm struct {
	Card string `form:"card"  example:"yearCard"`
}
