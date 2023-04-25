package model

import (
	"time"
)

type MsgRecord struct {
	AccountFrom string    `gorm:"size:50;not null;" json:"accountFrom"`
	AccountTo   string    `gorm:"size:50;not null;" json:"accountTo"`
	Message     string    `gorm:"size:100;not null;" json:"message"`
	CreatedAt   time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
}

type Msg struct {
	Account string `gorm:"size:50;not null;" json:"account"`
}
