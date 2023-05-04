package model

import (
	"time"
)

type User struct {
	//gorm.Model
	//ID           int64     `gorm:"primary_key;auto_increment" json:"id"`
	Account      string    `gorm:"primary_key;size:50;not null;unique" json:"account"`
	Password     string    `gorm:"size:30;not null;" json:"password"`
	Name         string    `gorm:"size:20;not null;"  json:"name" `
	Email        string    `gorm:"size:50;" json:"email"`
	Phone        string    `gorm:"size:15;" json:"phone"`
	Zipcode      string    `gorm:"size:3;"  json:"zipcode"`
	VipLevel     string    `gorm:"size:1;"  json:"viplevel"`
	VipDate      time.Time `gorm:"type:timestamp;" json:"vipdate"`
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

type ProfileSave struct {
	Account      string `form:"account"  example:"kevin@gmail.com"`
	OldPassword  string `form:"oldPassword"  example:"123456"`
	Password     string `form:"password"  example:"123456"`
	Name         string `form:"name"  example:"桐谷和人"`
	Phone        string `form:"phone"  example:"0999999999"`
	Zipcode      string `form:"zipcode"  example:"200"`
	Introduction string `form:"introduction"  example:"我有 8000 名部下"`
	PwdSwitch    bool   `form:"pwdSwitch"  example:"true"`
}

// soho設定(案主可公開查詢)
type SohoSetting struct {
	Account     string `gorm:"primary_key;size:50;not null;unique" json:"account"`
	Open        string `gorm:"size:1;not null;" json:"open"`
	Name        string `gorm:"size:20;not null;"  json:"name" `
	Role        string `gorm:"size:1;not null;COMMENT:'1:個人兼職,2:專職SOHO,3:工作室,4:兼職上班族,5:公司,6:學生';" json:"role"`
	Phone       string `gorm:"size:15;not null;" json:"phone"`
	CityTalk    string `gorm:"size:4;" json:"cityTalk"`
	CityTalk2   string `gorm:"size:10;" json:"cityTalk2"`
	Extension   string `gorm:"size:5;" json:"extension"`
	Email       string `gorm:"size:50;not null;" json:"email"`
	Zipcode     string `gorm:"size:3;not null;"  json:"zipcode"`
	Type        string `gorm:"size:20;not null;" json:"type"`
	Exp         string `gorm:"size:1;not null;COMMENT:'0:無接案經驗';" json:"exp"`
	Description string `gorm:"size:200;not null;" json:"description"`
}

type SohoSettingForm struct {
	Open        string `form:"open"  example:"1"`
	Name        string `form:"name"  example:"kevin"`
	Role        string `form:"role"  example:"0"`
	Phone       string `form:"phone"  example:"0999999999"`
	CityTalk    string `form:"cityTalk"  example:"04"`
	CityTalk2   string `form:"cityTalk2"  example:"29992222"`
	Extension   string `form:"extension"  example:"9527"`
	Email       string `form:"email"  example:"kevin@gmail.com"`
	Zipcode     string `form:"zipcode"  example:"200"`
	Type        string `form:"type"  example:"1,2,3"`
	ExpVal      string `form:"expVal"  example:"無接案經驗"`
	Description string `form:"description"  example:"i have ten years experience"`
}
