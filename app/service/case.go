package service

import (
	"fmt"
	"sync"
	"time"
	"znews/app/dao"
	"znews/app/model"
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

	caseId := fmt.Sprintf("%d%02d%03d", year, month, serial.No)
	return caseId, nil
}
