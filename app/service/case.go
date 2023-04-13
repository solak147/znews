package service

import (
	"znews/app/dao"
	"znews/app/model"
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

	if err := regexpRigister(`^[o,i]$`, form.ExpectDate); err != nil {
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

	casem := model.Case{
		Account:       form.Account,
		Title:         form.Title,
		Type:          form.Type,
		Kind:          form.Kind,
		ExpectDate:    form.ExpectDate,
		ExpectDateChk: form.ExpectDateChk,
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

	insertErr := dao.SqlSession.Model(&model.Case{}).Create(&casem).Error
	return insertErr
}
