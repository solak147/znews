package model

import (
	"time"
)

type Casem struct {
	CaseId        string `gorm:"primary_key;size:10;" json:"caseId"`
	Account       string `gorm:"size:50;not null;" json:"account"`
	Title         string `gorm:"size:20;not null;" json:"title"`
	Type          string `gorm:"size:4;not null;" json:"type"`
	Kind          string `gorm:"size:1;not null;" json:"kind"`
	ExpectDate    string `gorm:"size:10;" json:"expectDate"`
	ExpectDateChk string `gorm:"size:1;not null;" json:"expectDateChk"`
	ExpectMoney   string `gorm:"size:10;not null;" json:"expectMoney"`
	WorkArea      string `gorm:"size:7;not null;" json:"workArea"`
	WorkAreaChk   string `gorm:"size:1;not null;" json:"workAreaChk"`
	WorkContent   string `gorm:"size:200;not null;" json:"workContent"`

	Name        string `gorm:"size:20;not null;" json:"name"`
	Phone       string `gorm:"size:20;not null;" json:"phone"`
	CityTalk    string `gorm:"size:4;" json:"CityTalk"`
	CityTalk2   string `gorm:"size:10;" json:"CityTalk2"`
	Extension   string `gorm:"size:5;" json:"Extension"`
	ContactTime string `gorm:"size:7;not null;" json:"contactTime"`
	Email       string `gorm:"size:50;not null" json:"email"`
	Line        string `gorm:"size:20;" json:"line"`

	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type CreateCase struct {
	Account       string `form:"account"  example:"kevin@gmail.com"`
	Title         string `form:"title"  example:"電傷平台架設"`
	Type          string `form:"type"  example:"程式開發"`
	Kind          string `form:"kind"  example:"o:一般案件 i:急件"`
	ExpectDate    string `form:"expectDate"  example:"2022/02/02"`
	ExpectDateChk string `form:"expectDateChk"  example:"1"`
	ExpectMoney   string `form:"expectMoney"  example:"5000"`
	WorkArea      string `form:"workArea"  example:"台北市 信義區"`
	WorkAreaChk   string `form:"workAreaChk"  example:"1"`
	WorkContent   string `form:"workContent"  example:"電傷平台架設，伺服器維護..."`

	Name        string   `form:"name"  example:"kevin"`
	Phone       string   `form:"phone"  example:"0999999999"`
	CityTalk    string   `form:"cityTalk"  example:"04"`
	CityTalk2   string   `form:"cityTalk2"  example:"26232356"`
	Extension   string   `form:"extension"  example:"0000"`
	ContactTime string   `form:"contactTime"  example:"m,a"`
	Email       string   `form:"email"  example:"kevin@gmail.com"`
	Line        string   `form:"line"  example:"imlineid"`
	FilesName   []string `form:"filesName"  example:"[a.jpg]"`
}
