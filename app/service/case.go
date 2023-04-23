package service

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"znews/app/dao"
	"znews/app/model"

	"github.com/gin-gonic/gin"
)

var (
	mu sync.Mutex
)

func CreateCase(form model.CreateCase) error {
	if err := regexpRigister(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, form.Account); err != nil {
		return err
	}

	if err := regexpRigister(`^.{3,20}$`, form.Title); err != nil {
		return err
	}

	if err := regexpRigister(`^.{4}$`, form.Type); err != nil {
		return err
	}

	if err := regexpRigister(`^[o,i]$`, form.Kind); err != nil {
		return err
	}

	// 0-30 or yyyy/mm/dd
	if err := regexpRigister(`^([0-9]{4}/(0[1-9]|1[012])/(0[1-9]|[12][0-9]|3[01])|([12][0-9]|30|[0-9])|)$`, form.ExpectDate); err != nil {
		return err
	}

	if err := regexpRigister(`^[1,2,3]$`, form.ExpectDateChk); err != nil {
		return err
	}

	if err := regexpRigister(`^.{1,200}$`, form.WorkContent); err != nil {
		return err
	}

	if err := regexpRigister(`^\S.{0,13}\S?$`, form.Name); err != nil {
		return err
	}

	if err := regexpRigister(`^\d{1,15}$`, form.Phone); err != nil {
		return err
	}

	if err := regexpRigister(`^\d{0,4}$`, form.CityTalk); err != nil {
		return err
	}

	if err := regexpRigister(`^\d{0,10}$`, form.CityTalk2); err != nil {
		return err
	}

	if err := regexpRigister(`^\d{0,5}$`, form.Extension); err != nil {
		return err
	}

	if err := regexpRigister(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, form.Email); err != nil {
		return err
	}

	if err := regexpRigister(`^[a-zA-Z0-9_-]*$`, form.Line); err != nil {
		return err
	}

	caseId, genErr := genCaseId()
	if genErr != nil {
		return genErr
	}

	casem := model.Casem{
		CaseId:        caseId,
		Account:       form.Account,
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

	err := dao.SqlSession.Model(&model.Casem{}).Create(&casem).Error
	if err != nil {
		return err
	}

	for _, name := range form.FilesName {
		file := model.CaseFile{
			CaseId:   caseId,
			FileName: name,
		}

		err := dao.SqlSession.Model(&model.CaseFile{}).Create(&file).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func genCaseId() (string, error) {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now()
	year, month, _ := now.Date()
	monthFmt := fmt.Sprintf("%02d", month)

	serial := model.SerialNo{}

	err := dao.SqlSession.Select("*").Where("year=? and month=?", year, monthFmt).First(&serial).Error
	if err != nil {
		return "", err
	}

	uptNo := model.SerialNo{
		No: serial.No + 1,
	}

	insErr := dao.SqlSession.Model(&model.SerialNo{}).Where("year=? and month=?", year, monthFmt).Updates(uptNo).Error
	if insErr != nil {
		return "", insErr
	}

	caseId := fmt.Sprintf("%d%02d%03d", year, month, serial.No+1)
	return caseId, nil
}

func GetCase(c *gin.Context) ([]interface{}, error, int64) {

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
		orStr = fmt.Sprintf(`,"should": [
									{"term": {"work_area_chk": "1"}},
									{"match_phrase": {"work_area": "%s"}}
								]`, city)

	}

	matchSub := ""
	if search != "" {
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

	var cnt int64
	if err := dao.SqlSession.Model(&model.Casem{}).Count(&cnt).Error; err != nil {
		return nil, err, 0
	}

	return data, nil, cnt
}

func GetCaseDetail(caseId string, isAuth bool) (*model.Casem, []model.CaseFile, error) {

	fields := []string{"title", "type", "kind", "expect_date", "expect_date_chk", "expect_money", "work_area", "work_area_chk", "work_content", "updated_at"}
	casem := &model.Casem{}

	var err error
	if isAuth {
		err = dao.SqlSession.Select("*").Where("case_id=?", caseId).First(&casem).Error
	} else {
		err = dao.SqlSession.Select(fields).Where("case_id=?", caseId).First(&casem).Error
	}

	if err != nil {
		return nil, nil, err
	}

	if true {
		if len(casem.Name) > 3 {
			casem.Name = casem.Name[:1] + "**" + casem.Name[4:]
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

		if len(casem.Line) > 2 {
			casem.Line = casem.Line[:1] + "****" + casem.Line[3:]
		}

		emailArr := strings.Split(casem.Email, "@")
		casem.Email = "*****@" + emailArr[1]
	}

	var files []model.CaseFile
	filesErr := dao.SqlSession.Select("*").Where("case_id=?", caseId).Find(&files).Error
	if filesErr != nil {
		return nil, nil, err
	} else {
		return casem, files, nil
	}
}
