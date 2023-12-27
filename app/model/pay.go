package model

import (
	"time"
)

type Product struct {
	Name      string    `gorm:"primary_key;size:10;not null;" json:"name"`
	Price     string    `gorm:"size:4;not null;" json:"price"`
	ChiName   string    `gorm:"size:10;not null;" json:"chiName"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type ProductForm struct {
	Card string `form:"card"  example:"yearCard"`
}

type Order struct {
	MerchantID           string `gorm:"primary_key;size:10;not null; COMMENT:'特店編號';" json:"MerchantID"`
	MerchantTradeNo      string `gorm:"primary_key;size:20;not null; COMMENT:'訂單產生時傳送給綠界的特店交易編號';" json:"MerchantTradeNo"`
	RtnCode              string `gorm:"size:5; not null; COMMENT:'若回傳值為1時，為付款成功,其餘代碼皆為交易異常';" json:"RtnCode"`
	RtnMsg               string `gorm:"size:200; not null; COMMENT:'付款成功回傳:RtnMsg=交易成功, 排程回傳:RtnMsg=paid, OrderResultURL(client端)回傳：RtnMsg=Succeeded';" json:"RtnMsg"`
	TradeNo              string `gorm:"size:20; not null; COMMENT:'綠界的交易編號';" json:"TradeNo"`
	TradeAmt             string `gorm:"size:5; not null; COMMENT:'交易金額';" json:"TradeAmt"`
	PaymentDate          string `gorm:"size:20; not null; COMMENT:'付款時間 yyyy/MM/dd HH:mm:ss';" json:"PaymentDate"`
	PaymentType          string `gorm:"size:20;not null; COMMENT:'特店選擇的付款方式';" json:"PaymentType"`
	PaymentTypeChargeFee string `gorm:"size:3; not null; COMMENT:'交易手續費金額';" json:"PaymentTypeChargeFee"`
	TradeDate            string `gorm:"size:20; not null; COMMENT:'訂單成立時間 yyyy/MM/dd HH:mm:ss';" json:"TradeDate"`
	SimulatePaid         string `gorm:"size:1; not null; COMMENT:'0：代表此交易非模擬付款。1：代表此交易為模擬付款';" json:"SimulatePaid"`
	CustomField1         string `gorm:"size:10; not null; COMMENT:'產品名稱';" json:"CustomField1"`
	CustomField2         string `gorm:"size:50; not null; COMMENT:'帳號';" json:"CustomField2"`
	CheckMacValue        string `gorm:"not null; COMMENT:'檢查碼';" json:"CheckMacValue"`
	IsCheck              string `gorm:"size:1; not null; COMMENT:'檢查碼比對正確';" json:"IsCheck"`

	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}
