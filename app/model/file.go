package model

import (
	"time"
)

type CaseFile struct {
	CaseId    string    `gorm:"size:10;not null;" json:"caseId"`
	FileName  string    `gorm:"size:20;not null;" json:"filename"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type SerialNo struct {
	Year  string `gorm:"size:4;not null;" json:"year"`
	Month string `gorm:"size:2;not null;" json:"month"`
	No    int64  `gorm:"size:4;not null;" json:"no"`
}

// 個人作品
type SohoWork struct {
	Account   string    `gorm:"primary_key;size:50;not null;" json:"account"`
	FileName  string    `gorm:"primary_key;size:20;not null;" json:"filename"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}
