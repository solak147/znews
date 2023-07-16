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
	Phone       string `gorm:"size:15;not null;" json:"phone"`
	CityTalk    string `gorm:"size:4;" json:"cityTalk"`
	CityTalk2   string `gorm:"size:10;" json:"cityTalk2"`
	Extension   string `gorm:"size:5;" json:"extension"`
	ContactTime string `gorm:"size:7;not null;" json:"contactTime"`
	Email       string `gorm:"size:50;not null" json:"email"`
	Line        string `gorm:"size:20;" json:"line"`

	QuoteTotal  int    `gorm:"size:2;default:0;COMMENT:'報價人數';" json:"quoteTotal"`
	Status      string `gorm:"size:2;default:0;COMMENT:'-1:下架,0:已發案,1:已成交,2:已完成,3:案主評價,4:接案評價';" json:"status"`
	BossStar    string `gorm:"size:1;" json:"bossStar"`
	BossComment string `gorm:"size:30;" json:"bossComment"`
	SohoStar    string `gorm:"size:1;" json:"sohoStar"`
	SohoComment string `gorm:"size:30;" json:"sohoComment"`

	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type CaseByIdForm struct {
	CaseId string `form:"caseId"  example:"202402001"`
}

type CaseForm struct {
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

type CaseDetailRtn struct {
	Casem        *Casem
	CaseFile     []CaseFile
	IsVip        bool
	IsCollection bool
	Error        error
}

// 未成交的報價紀錄案件
type QuoteCaseRec struct {
	CaseId      string    `form:"caseId"  example:"202402001"`
	Title       string    `form:"title"  example:"電傷平台架設"`
	ExpectMoney string    `form:"expectMoney"  example:"5000"`
	WorkArea    string    `form:"workArea"  example:"台北市 信義區"`
	WorkAreaChk string    `form:"workAreaChk"  example:"1"`
	WorkContent string    `form:"workContent"  example:"電傷平台架設，伺服器維護..."`
	PriceS      int       `form:"priceS"  example:"1000"`
	PriceE      int       `form:"priceE"  example:"2000"`
	QuoteTotal  string    `form:"quoteTotal"  example:"1"`
	Status      string    `form:"status"  example:"1"`
	UpdatedAt   time.Time `form:"updatedAt"  example:"2023-04-21 04:16:50"`
}

// 報價人列表
type QuotePerRtn struct {
	Account   string
	Name      string
	Email     string
	Phone     string
	PriceS    int
	PriceE    int
	Day       int
	UpdatedAt time.Time
}

// 報價
type Quote struct {
	CaseId    string    `gorm:"primary_key;size:10;" json:"caseId"`
	Account   string    `gorm:"primary_key;size:50;not null;" json:"account"`
	PriceS    int       `gorm:"type:int;not null;" json:"priceS"`
	PriceE    int       `gorm:"type:int;not null;" json:"priceE"`
	Day       int       `gorm:"type:int;not null;" json:"day"`
	Deal      int       `gorm:"size:1;default:0" json:"deal"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type QuoteForm struct {
	Title   string `form:"title"  example:"app繪圖軟體製作"`
	CaseId  string `form:"caseId"  example:"202304005"`
	Account string `form:"account"  example:"kevin@gmail.com"`
	PriceS  int    `form:"priceS"  example:"1000"`
	PriceE  int    `form:"priceE"  example:"2000"`
	Day     int    `form:"day"  example:"10"`
}

// 案件流程
type CaseFlow struct {
	CaseId    string    `gorm:"primary_key;size:10;not null;" json:"caseId"`
	Status    string    `gorm:"primary_key;size:1;not null;COMMENT:'1:已成交,2:已結案,3：案主評價,4:接案評價';" json:"status"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
}

type Flow struct {
	CaseId string `form:"caseId"  example:"202304005"`
	Status string `form:"status"  example:"2"`

	BossStar    string `form:"bossStar"  example:"5"`
	BossComment string `form:"bossComment"  example:"very good"`
	SohoStar    string `form:"sohoStar"  example:"5"`
	SohoComment string `form:"sohoComment"  example:"very good"`
}

// 收藏案件
type CaseCollect struct {
	CaseId    string    `gorm:"primary_key;size:10;" json:"caseId"`
	Account   string    `gorm:"primary_key;size:50;not null;" json:"account"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
}

type CaseCollectForm struct {
	CaseId string `form:"caseId"  example:"202304005"`
	IsLike string `form:"isLike"  example:"1"`
}

// 收藏案件紀錄
type CaseCollectRec struct {
	CaseId      string    `form:"caseId"  example:"202402001"`
	Title       string    `form:"title"  example:"電傷平台架設"`
	ExpectMoney string    `form:"expectMoney"  example:"5000"`
	WorkArea    string    `form:"workArea"  example:"台北市 信義區"`
	WorkAreaChk string    `form:"workAreaChk"  example:"1"`
	WorkContent string    `form:"workContent"  example:"電傷平台架設，伺服器維護..."`
	QuoteTotal  string    `form:"quoteTotal"  example:"1"`
	Status      string    `form:"status"  example:"1"`
	UpdatedAt   time.Time `form:"updatedAt"  example:"2023-04-21 04:16:50"`
}

// 發布案件紀錄
type CaseReleaseRec struct {
	CaseId      string    `form:"caseId"  example:"202402001"`
	Title       string    `form:"title"  example:"電傷平台架設"`
	ExpectMoney string    `form:"expectMoney"  example:"5000"`
	WorkArea    string    `form:"workArea"  example:"台北市 信義區"`
	WorkAreaChk string    `form:"workAreaChk"  example:"1"`
	WorkContent string    `form:"workContent"  example:"電傷平台架設，伺服器維護..."`
	QuoteTotal  string    `form:"quoteTotal"  example:"1"`
	Status      string    `form:"status"  example:"1"`
	UpdatedAt   time.Time `form:"updatedAt"  example:"2023-04-21 04:16:50"`
}

// 下架案件紀錄
type CaseCloseRec struct {
	CaseId      string    `form:"caseId"  example:"202402001"`
	Title       string    `form:"title"  example:"電傷平台架設"`
	ExpectMoney string    `form:"expectMoney"  example:"5000"`
	WorkArea    string    `form:"workArea"  example:"台北市 信義區"`
	WorkAreaChk string    `form:"workAreaChk"  example:"1"`
	WorkContent string    `form:"workContent"  example:"電傷平台架設，伺服器維護..."`
	QuoteTotal  string    `form:"quoteTotal"  example:"1"`
	Status      string    `form:"status"  example:"1"`
	UpdatedAt   time.Time `form:"updatedAt"  example:"2023-04-21 04:16:50"`
}
