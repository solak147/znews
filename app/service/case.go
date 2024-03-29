package service

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"znews/app/dao"
	"znews/app/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	mu sync.Mutex
)

func CreateCase(c *gin.Context) error {
	account, _ := c.Get("account")

	var form model.CaseForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "表單綁定失敗 : " + err.Error(),
		})
		return err
	}

	if err := chkCaseForm(form); err != nil {
		return err
	}

	tx := dao.GormSession.Begin()

	caseId, genErr := genCaseId(tx)
	if genErr != nil {
		tx.Rollback()
		return genErr
	}

	casem := model.Casem{
		CaseId:        caseId,
		Account:       fmt.Sprintf("%v", account),
		Title:         form.Title,
		Type:          form.Type,
		Kind:          form.Kind,
		ExpectDate:    form.ExpectDate,
		ExpectDateChk: form.ExpectDateChk,
		ExpectMoney:   form.ExpectMoney,
		WorkArea:      form.WorkArea,
		WorkAreaChk:   form.WorkAreaChk,
		WorkContent:   form.WorkContent,

		Name:        form.Name,
		Phone:       form.Phone,
		CityTalk:    form.CityTalk,
		CityTalk2:   form.CityTalk2,
		Extension:   form.Extension,
		ContactTime: form.ContactTime,
		Email:       form.Email,
		Line:        form.Line,
	}

	err := tx.Model(&model.Casem{}).Create(&casem).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, name := range form.FilesName {
		file := model.CaseFile{
			CaseId:   caseId,
			FileName: name,
		}

		err := tx.Model(&model.CaseFile{}).Create(&file).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := Uploads(c, caseId); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func genCaseId(tx *gorm.DB) (string, error) {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now()
	year, month, _ := now.Date()
	monthFmt := fmt.Sprintf("%02d", month)

	serial := model.SerialNo{}

	err := tx.Select("*").Where("year=? and month=?", year, monthFmt).First(&serial).Error
	if err != nil {
		return "", err
	}

	uptNo := model.SerialNo{
		No: serial.No + 1,
	}

	insErr := tx.Model(&model.SerialNo{}).Where("year=? and month=?", year, monthFmt).Updates(uptNo).Error
	if insErr != nil {
		return "", insErr
	}

	caseId := fmt.Sprintf("%d%02d%03d", year, month, serial.No+1)
	return caseId, nil
}

func GetCase(c *gin.Context) ([]interface{}, error, int) {

	search := c.Query("search")
	city := c.Query("city")
	typ := c.Query("type")
	price := c.Query("price")
	from := c.Query("from")
	size := os.Getenv("Sel_Case_Size")

	fromInt, _ := strconv.Atoi(from)
	sizeInt, _ := strconv.Atoi(size)

	orStr := ""
	if city != "" {
		// 工作地點不限 or 指定定點
		// match_phrase 可以指定以片語來搜尋。片語必須要完全符合，也就是不會被拆開成單詞。
		orStr = fmt.Sprintf(`,"should": [
									{"term": {"work_area_chk": "1"}},
									{"match_phrase": {"work_area": "%s"}}
								]`, city)

	}

	matchSub := `{"term": {"status": "0"}}`
	if search != "" {
		matchSub += ","
		// 在 title, work_content 欄位中找包含search輸入字串的資料
		matchSub += fmt.Sprintf(`{"multi_match": {"query": "%s", "fields": ["title", "work_content"]}}`, search)
	}

	if typ != "" {
		if matchSub != "" {
			matchSub += ","
		}
		matchSub += fmt.Sprintf(`{"match_phrase": {"type": "%s"}}`, typ)
	}

	if price != "" {
		if matchSub != "" {
			matchSub += ","
		}
		matchSub += fmt.Sprintf(`{"match_phrase": {"expect_money": "%s"}}`, price)
	}

	// bool query 用於將多個條件組合在一起，而他主要由三個部份組成 :
	// must : 所有條件都必須完全匹配，等於 AND。
	// should : 至少一個條件要匹配，等於 OR。
	// must_not : 所有條件都不能匹配，等於 NOT。
	match := ""
	if matchSub != "" || orStr != "" {
		match = fmt.Sprintf(`{

			"bool": {
				"must": [
					%s
				]
				%s
			}
			
		  }`, matchSub, orStr)

	} else {
		match = `{"match_all": {}}`
	}

	query := fmt.Sprintf(`{
        "query": %s,
        "sort": [
            {
              "updated_at": {
                "order": "desc"
              }
            }
          ],

        "from": %d, 
        "size": %d
    }`, match, fromInt, sizeInt)

	data, eserr := dao.GetCase(query)
	if eserr != nil {
		return nil, eserr, 0
	}

	cnt, cnterr := dao.GetCaseCount()
	if cnterr != nil {
		return nil, cnterr, 0
	}

	// var cnt int64
	// if err := dao.GormSession.Model(&model.Casem{}).Count(&cnt).Error; err != nil {
	// 	return nil, err, 0
	// }

	return data, nil, cnt
}

func UpdateCase(c *gin.Context) error {

	caseId := c.Params.ByName("caseId")

	var form model.CaseForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "表單綁定失敗 : " + err.Error(),
		})
		return err
	}

	if err := chkCaseForm(form); err != nil {
		return err
	}

	tx := dao.GormSession.Begin()

	casem := model.Casem{
		Title:         form.Title,
		Type:          form.Type,
		Kind:          form.Kind,
		ExpectDate:    form.ExpectDate,
		ExpectDateChk: form.ExpectDateChk,
		ExpectMoney:   form.ExpectMoney,
		WorkArea:      form.WorkArea,
		WorkAreaChk:   form.WorkAreaChk,
		WorkContent:   form.WorkContent,

		Name:        form.Name,
		Phone:       form.Phone,
		CityTalk:    form.CityTalk,
		CityTalk2:   form.CityTalk2,
		Extension:   form.Extension,
		ContactTime: form.ContactTime,
		Email:       form.Email,
		Line:        form.Line,
		UpdatedAt:   time.Now(),
	}

	if err := tx.Model(&model.Casem{}).Where("case_id = ?", caseId).Updates(casem).Error; err != nil {
		tx.Rollback()
		return err
	}

	file := model.CaseFile{}
	if err := tx.Where("case_id = ?", caseId).Delete(file).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, name := range form.FilesName {
		file := model.CaseFile{
			CaseId:   caseId,
			FileName: name,
		}

		err := tx.Model(&model.CaseFile{}).Create(&file).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := Uploads(c, caseId); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func GetCaseDetailOri(caseId string, account string) model.CaseDetailRtn {
	res := model.CaseDetailRtn{}
	casem := model.Casem{}

	if err := dao.GormSession.Select("*").Where("case_id=?", caseId).First(&casem).Error; err != nil {
		res.Error = err
		return res
	}

	if account != casem.Account {
		res.Error = errors.New("非你所擁有的案件")
	}

	var files []model.CaseFile
	if err := dao.GormSession.Select("*").Where("case_id=?", caseId).Find(&files).Error; err != nil {
		res.Error = err
		return res
	} else {
		res.Casem = &casem
		res.CaseFile = files

		return res
	}
}

func GetCaseDetail(caseId string, account string) model.CaseDetailRtn {

	res := model.CaseDetailRtn{}

	fields := []string{"case_id", "account", "title", "type", "kind", "expect_date", "expect_date_chk", "expect_money", "work_area", "work_area_chk", "work_content", "updated_at", "status"}
	casem := model.Casem{}

	//是否已登入
	var err error
	if account != "<nil>" {
		err = dao.GormSession.Select("*").Where("case_id=?", caseId).First(&casem).Error
	} else {
		err = dao.GormSession.Select(fields).Where("case_id=?", caseId).First(&casem).Error
	}

	if err != nil {
		res.Error = err
		return res
	}

	if account != "<nil>" {

		user := model.User{}
		if err := dao.GormSession.Select("vip_date").Where("account=?", account).First(&user).Error; err != nil {
			res.Error = err
			return res
		}

		// 计算时间戳与当前时间之间的时间差
		// duration := time.Since(user.VipDate)
		// 计算时间差对应的月数
		// months := int(duration.Hours() / 24 / 30)

		now := time.Now()

		if now.After(user.VipDate) {

			if len(casem.Name) > 3 {
				casem.Name = casem.Name[:1] + "****" + casem.Name[4:]
			} else {
				casem.Name = casem.Name[:1] + "**"
			}

			if len(casem.Phone) > 10 {
				casem.Phone = casem.Phone[:3] + "*******" + casem.Phone[11:]
			} else {
				casem.Phone = casem.Phone[:3] + "*******"
			}

			if len(casem.CityTalk) > 1 {
				casem.CityTalk = "**"
			}

			if len(casem.CityTalk2) >= 4 {
				casem.CityTalk2 = casem.CityTalk2[:4] + "****"
			}

			if len(casem.Extension) >= 2 {
				casem.Extension = casem.Extension[:2] + "**"
			}

			if len(casem.Line) > 3 {
				casem.Line = casem.Line[:1] + "****" + casem.Line[4:]
			} else if len(casem.Line) > 0 {
				casem.Line = casem.Line[:1] + "***"
			}

			emailArr := strings.Split(casem.Email, "@")
			casem.Email = "*****@" + emailArr[1]
		} else {
			// vip是否到期
			res.IsVip = true
		}

		var cnt int64
		if err := dao.GormSession.Model(&model.CaseCollect{}).Where("case_id= ? and account = ?", caseId, account).Count(&cnt).Error; err != nil {
			res.Error = err
			return res
		} else {
			if cnt > 0 {
				// 已收藏
				res.IsCollection = true
			}
		}

	}

	var files []model.CaseFile
	filesErr := dao.GormSession.Select("*").Where("case_id=?", caseId).Find(&files).Error
	if filesErr != nil {
		res.Error = filesErr
		return res
	} else {
		res.Casem = &casem
		res.CaseFile = files

		return res
	}
}

func Quote(account string, m model.QuoteForm) error {

	if err := regexpRigister(`^\d{1,8}$`, strconv.Itoa(m.PriceS)); err != nil {
		return err
	}
	if err := regexpRigister(`^\d{1,8}$`, strconv.Itoa(m.PriceE)); err != nil {
		return err
	}
	if err := regexpRigister(`^\d{1,4}$`, strconv.Itoa(m.Day)); err != nil {
		return err
	}

	//檢查是否有購買vip
	user := model.User{}
	if err := dao.GormSession.Where("account = ?", account).Select("account, vip_date").First(&user).Error; err != nil {
		return err
	}

	vipDate := user.VipDate
	//使用 time.Time 类型的零值来表示“无效”时间值
	if vipDate.IsZero() {
		return errors.New("請先購買vip才能使用此功能")
	}

	now := time.Now()

	if now.After(vipDate) {
		return errors.New("vip已過期")
	}

	tx := dao.GormSession.Begin()

	quote := model.Quote{
		CaseId:  m.CaseId,
		Account: account, //報價者
		PriceS:  m.PriceS,
		PriceE:  m.PriceE,
		Day:     m.Day,
	}

	err := tx.Model(&model.Quote{}).Create(&quote).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	casem := model.Casem{}
	if err := tx.Where("case_id = ?", m.CaseId).Select("account").First(&casem).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 總報價人數+1，elasticsearch撈資料用
	uc := model.Casem{
		QuoteTotal: casem.QuoteTotal + 1,
	}
	if err := tx.Model(&model.Casem{}).Where("case_id = ?", m.CaseId).Updates(uc).Error; err != nil {
		return err
	}

	// 流程狀態紀錄
	flow := model.CaseFlow{
		CaseId: m.CaseId,
		Status: "1",
	}
	if err := tx.Model(&model.CaseFlow{}).Create(&flow).Error; err != nil {
		tx.Rollback()
		return err
	}

	msg := model.MsgRecord{
		AccountFrom: account,
		AccountTo:   casem.Account,
		Message:     fmt.Sprintf("%s-=%s-=%d-=%d-=%d", m.CaseId, m.Title, m.PriceS, m.PriceE, m.Day),
		IsSystem:    "1",
	}

	msgErr := tx.Model(&model.MsgRecord{}).Create(&msg).Error
	if msgErr != nil {
		tx.Rollback()
		return msgErr
	} else {
		tx.Commit()
		return nil
	}
}

// 案主成交 & 人才成交 ＆ 人才報價紀錄 共用
func QuoteRecord(c *gin.Context) ([]model.QuoteCaseRec, error) {

	accountAny, _ := c.Get("account")
	account := fmt.Sprintf("%v", accountAny)

	deal := c.Params.ByName("deal")
	userType := c.Query("userType")
	status := c.Query("status")
	finish := c.Query("finish")

	caseArr := []model.QuoteCaseRec{}

	var (
		query     string
		statusSql string
		finishSql string
		rows      *sql.Rows
		err       error
	)

	if status != "" {
		statusSql = fmt.Sprintf(`AND status = %s`, status)
	}

	if finish == "false" {
		finishSql = `AND status != 4`
	}

	if userType == "boss" {
		query = fmt.Sprintf(`SELECT case_id, title, expect_money, work_area, work_area_chk, work_content, quote_total, status, updated_at ,
					(SELECT price_s FROM quotes WHERE case_id =  a.case_id AND deal = '1') price_s,
					(SELECT price_e FROM quotes WHERE case_id =  a.case_id AND deal = '1') price_e
				FROM casems a 
				WHERE account = ? AND status > 0 %s %s
				ORDER BY updated_at DESC`, statusSql, finishSql)
		rows, err = dao.DbSession.Query(query, account)
	} else {
		query = fmt.Sprintf(`SELECT case_id, title, expect_money, work_area, work_area_chk, work_content, quote_total, status, updated_at ,
					(SELECT price_s FROM quotes WHERE account = ?  AND case_id =  a.case_id) price_s,
					(SELECT price_e FROM quotes WHERE account = ?  AND case_id =  a.case_id) price_e
				FROM casems a 
				WHERE case_id IN (
						SELECT case_id FROM quotes
						WHERE account = ? AND deal = ?) %s %s
				ORDER BY updated_at DESC`, statusSql, finishSql)
		rows, err = dao.DbSession.Query(query, account, account, account, deal)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c model.QuoteCaseRec
		if err := rows.Scan(&c.CaseId, &c.Title, &c.ExpectMoney, &c.WorkArea, &c.WorkAreaChk, &c.WorkContent, &c.QuoteTotal, &c.Status, &c.UpdatedAt, &c.PriceS, &c.PriceE); err != nil {
			return nil, err
		}
		caseArr = append(caseArr, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return caseArr, nil
}

// 報價人列表
func QuotePerLst(caseId string) ([]model.QuotePerRtn, error) {

	caseArr := []model.QuotePerRtn{}

	var (
		query string
		rows  *sql.Rows
		err   error
	)

	query = `SELECT (SELECT name from users WHERE account = a.account) name,
			(SELECT email from users WHERE account = a.account) email,
			(SELECT phone from users WHERE account = a.account) phone,
			account, price_s, price_e, day, updated_at
			FROM quotes a  
			WHERE case_id = ?`
	rows, err = dao.DbSession.Query(query, caseId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c model.QuotePerRtn
		if err := rows.Scan(&c.Name, &c.Email, &c.Phone, &c.Account, &c.PriceS, &c.PriceE, &c.Day, &c.UpdatedAt); err != nil {
			return nil, err
		}
		caseArr = append(caseArr, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return caseArr, nil
}

// 案主 報價紀錄
func QuoteBossRecord(account string) ([]model.QuoteCaseRec, error) {

	caseArr := []model.QuoteCaseRec{}

	var (
		query string
		rows  *sql.Rows
		err   error
	)

	query = `SELECT case_id, title, expect_money, work_area, work_area_chk, work_content, quote_total, status, updated_at
			FROM casems a
			WHERE account = ? AND status = 0
			AND EXISTS (SELECT 1 FROM quotes WHERE case_id = a.case_id)
			ORDER BY updated_at DESC`
	rows, err = dao.DbSession.Query(query, account)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c model.QuoteCaseRec
		if err := rows.Scan(&c.CaseId, &c.Title, &c.ExpectMoney, &c.WorkArea, &c.WorkAreaChk, &c.WorkContent, &c.QuoteTotal, &c.Status, &c.UpdatedAt); err != nil {
			return nil, err
		}
		caseArr = append(caseArr, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return caseArr, nil
}

func GetFlow(caseId string) (*model.Casem, []model.CaseFlow, error) {
	var flow []model.CaseFlow

	if err := dao.GormSession.Select("*").Where("case_id=?", caseId).Find(&flow).Error; err != nil {
		return nil, nil, err
	}

	var casem model.Casem
	if err := dao.GormSession.Select("account, title, soho_star, boss_star, soho_comment, boss_comment").Where("case_id=?", caseId).Order("created_at asc").Find(&casem).Error; err != nil {
		return nil, nil, err
	} else {
		return &casem, flow, nil
	}

}

func Flow(form model.Flow) error {

	tx := dao.GormSession.Begin()
	casem := model.Casem{
		Status: form.Status,
	}

	if form.BossStar != "" {
		casem.BossStar = form.BossStar
	}

	if form.BossComment != "" {
		casem.BossComment = form.BossComment
	}

	if form.SohoStar != "" {
		casem.SohoStar = form.SohoStar
	}

	if form.SohoComment != "" {
		casem.SohoComment = form.SohoComment
	}

	if err := tx.Model(&model.Casem{}).Where("case_id = ?", form.CaseId).Updates(casem).Error; err != nil {
		tx.Rollback()
		return err
	}

	//結案
	flow := model.CaseFlow{
		CaseId: form.CaseId,
		Status: form.Status,
	}
	if err := tx.Model(&model.CaseFlow{}).Create(&flow).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func UpdateCollect(account string, form model.CaseCollectForm) error {
	if form.IsLike == "1" {
		collect := model.CaseCollect{}
		if err := dao.GormSession.Where("account = ? AND case_id = ?", account, form.CaseId).Delete(collect).Error; err != nil {
			return err
		}
	} else {
		collect := model.CaseCollect{
			Account: account,
			CaseId:  form.CaseId,
		}
		if err := dao.GormSession.Model(&model.CaseCollect{}).Create(collect).Error; err != nil {
			return err
		}

	}
	return nil
}

func GetCollect(account string) ([]model.CaseCollectRec, error) {

	var caseArr []model.CaseCollectRec

	query := `SELECT case_id, title, expect_money, work_area, work_area_chk, work_content, quote_total, status, updated_at 
			FROM casems  
			WHERE case_id IN (SELECT case_id FROM case_collects WHERE account = ? )
			ORDER BY updated_at DESC`

	rows, err := dao.DbSession.Query(query, account)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c model.CaseCollectRec
		if err := rows.Scan(&c.CaseId, &c.Title, &c.ExpectMoney, &c.WorkArea, &c.WorkAreaChk, &c.WorkContent, &c.QuoteTotal, &c.Status, &c.UpdatedAt); err != nil {
			return nil, err
		}
		caseArr = append(caseArr, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return caseArr, nil
}

func GetCaseRelease(account string) ([]model.CaseReleaseRec, error) {

	var caseArr []model.CaseReleaseRec

	query := `SELECT case_id, title, expect_money, work_area, work_area_chk, work_content, quote_total, status, updated_at 
			FROM casems  
			WHERE status = '0' AND account = ?
			ORDER BY updated_at DESC`

	rows, err := dao.DbSession.Query(query, account)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c model.CaseReleaseRec
		if err := rows.Scan(&c.CaseId, &c.Title, &c.ExpectMoney, &c.WorkArea, &c.WorkAreaChk, &c.WorkContent, &c.QuoteTotal, &c.Status, &c.UpdatedAt); err != nil {
			return nil, err
		}
		caseArr = append(caseArr, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return caseArr, nil
}

// 下架案件
func CloseCase(caseId string) error {

	casem := model.Casem{
		Status:    "-1",
		UpdatedAt: time.Now(),
	}
	if err := dao.GormSession.Model(&model.Casem{}).Where("case_id=?", caseId).Updates(casem).Error; err != nil {
		return err
	}

	return nil
}

// 取得下架案件
func GetCloseCase(account string) ([]model.CaseCloseRec, error) {

	var caseArr []model.CaseCloseRec

	query := `SELECT case_id, title, expect_money, work_area, work_area_chk, work_content, quote_total, status, updated_at 
			FROM casems  
			WHERE status = '-1' AND account = ?
			ORDER BY updated_at DESC`

	rows, err := dao.DbSession.Query(query, account)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c model.CaseCloseRec
		if err := rows.Scan(&c.CaseId, &c.Title, &c.ExpectMoney, &c.WorkArea, &c.WorkAreaChk, &c.WorkContent, &c.QuoteTotal, &c.Status, &c.UpdatedAt); err != nil {
			return nil, err
		}
		caseArr = append(caseArr, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return caseArr, nil
}

// 重新上架
func RePublish(caseId string) error {

	casem := model.Casem{
		Status:    "0",
		UpdatedAt: time.Now(),
	}
	if err := dao.GormSession.Model(&model.Casem{}).Where("case_id=?", caseId).Updates(casem).Error; err != nil {
		return err
	}

	return nil
}

// 刪除已下架案件
func DeleteCloseCase(caseId string) error {

	tx := dao.GormSession.Begin()

	if err := tx.Where("case_id = ?", caseId).Delete(model.Casem{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("case_id = ?", caseId).Delete(model.CaseFile{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	path := os.Getenv("Case_FILE_PATH")
	if err := os.RemoveAll(path + "/" + caseId); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
