package model

import (
	"time"
)

type MsgRecord struct {
	AccountFrom string    `gorm:"size:50;not null;" json:"accountFrom"`
	AccountTo   string    `gorm:"size:50;not null;" json:"accountTo"`
	Message     string    `gorm:"size:100;not null;" json:"message"`
	IsRead      string    `gorm:"size:1;default:0;" json:"isRead"`
	IsSystem    string    `gorm:"size:1;default:0;" json:"isSystem"`
	CreatedAt   time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
}

type MsgRec struct {
	Account    string `json:"account"`
	Message    string `json:"message"`
	CrtDte     string `json:"crtDte"`
	NotReadCnt string `json:"notReadCnt"`
}

type MsgSend struct {
	Message   string `form:"message"  example:"test"`
	AccountTo string `form:"accountTo"  example:"Mike"`
}

type MsgUpdateRead struct {
	AccountFrom string `form:"accountFrom"  example:"Mike"`
}
