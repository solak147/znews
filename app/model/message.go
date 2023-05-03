package model

import (
	"time"
)

type MsgRecord struct {
	AccountFrom string    `gorm:"size:50;not null;" json:"accountFrom"`
	AccountTo   string    `gorm:"size:50;not null;" json:"accountTo"`
	Message     string    `gorm:"size:100;not null;" json:"message"`
	IsRead      string    `gorm:"size:1;default:0;" json:"isRead"`
	IsSystem    string    `gorm:"size:1;default:0;COMMENT:'1:報價訊息,2:成交訊息';" json:"isSystem"`
	CreatedAt   time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
}

type MsgRec struct {
	Account    string `json:"account"`
	Message    string `json:"message"`
	IsSystem   string `json:"isSystem"`
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

type MsgDeal struct {
	Quoter string `form:"quoter"  example:"test@gmail.com"`
	CaseId string `form:"caseId"  example:"202305001"`
	Title  string `form:"title"  example:"APP繪圖軟體"`
	PriceS int    `form:"priceS"  example:"1000"`
	PriceE int    `form:"priceE"  example:"2000"`
}
