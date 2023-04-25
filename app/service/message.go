package service

import (
	"znews/app/dao"
	"znews/app/model"
)

func GetMsgRecord(account string) ([]model.MsgRecord, error) {
	var msg []model.MsgRecord
	if err := dao.SqlSession.Select("*").Where("account_from=? or account_to=?", account, account).Order("created_at ASC").Find(&msg).Error; err != nil {
		return nil, err
	} else {
		return msg, nil
	}
}
