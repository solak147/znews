package service

import (
	"errors"
	"fmt"
	"znews/app/dao"
	"znews/app/model"
)

func GetMsgRecord(account string) ([]model.MsgRec, error) {
	var msgArr []model.MsgRec

	// 执行 SQL 查询语句
	query := `select account, 
					(SELECT MAX(created_at) FROM msg_records WHERE (account_from = ? and account_to = a.account) or (account_from = a.account and account_to = ? )) crtDte, 
					(SELECT message FROM msg_records WHERE (account_from = ? and account_to = a.account) or (account_from = a.account and account_to = ? ) ORDER BY created_at DESC LIMIT 1) message,
					(SELECT is_system FROM msg_records WHERE (account_from = ? and account_to = a.account) or (account_from = a.account and account_to = ? ) ORDER BY created_at DESC LIMIT 1) isSystem, 
					(SELECT COUNT(*) FROM msg_records WHERE account_from = a.account and account_to = ? AND is_read = '0') notReadCnt   
				FROM (
					SELECT account_to account
					FROM msg_records WHERE account_from = ? GROUP BY account_to 
					UNION 
					SELECT account_from account
					FROM msg_records WHERE account_to = ? GROUP BY account_from
				)  a
				group by account`
	rows, err := dao.DbSession.Query(query, account, account, account, account, account, account, account, account, account)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var msg model.MsgRec
		if err := rows.Scan(&msg.Account, &msg.CrtDte, &msg.Message, &msg.IsSystem, &msg.NotReadCnt); err != nil {
			return nil, err
		}
		msgArr = append(msgArr, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return msgArr, nil

}

func GetMsgRecordDetail(from string, to string) ([]model.MsgRecord, error) {
	var msg []model.MsgRecord
	err := dao.GormSession.Where("account_from = ? and account_to = ? or account_from = ? and account_to = ?", from, to, to, from).Order("created_at asc").Find(&msg).Error
	if err != nil {
		return nil, err
	} else {
		return msg, nil
	}
}

func SendMsg(account string, m model.MsgSend) error {
	msg := model.MsgRecord{
		AccountFrom: account,
		AccountTo:   m.AccountTo,
		Message:     m.Message,
	}

	err := dao.GormSession.Model(&model.MsgRecord{}).Create(&msg).Error
	if err != nil {
		return err
	} else {
		return nil
	}
}

func ChkNoRead(account string) (int64, error) {
	var cnt int64
	if err := dao.GormSession.Model(&model.MsgRecord{}).Where("account_to = ? and is_read = '0'", account).Select("count(*)").Count(&cnt).Error; err != nil {
		return 0, err
	} else {
		return cnt, nil
	}
}

func UpdateRead(to string, from string) error {

	msg := model.MsgRecord{
		IsRead: "1",
	}

	err := dao.GormSession.Model(&model.MsgRecord{}).Where("account_from = ? and account_to = ?", from, to).Updates(msg).Error
	if err != nil {
		return err
	} else {
		return nil
	}
}

func MsgDeal(account string, m model.MsgDeal) error {

	casem := model.Casem{}
	//檢查案件是否已成交
	if err := dao.GormSession.Select("status").Where("case_id=?", m.CaseId).First(&casem).Error; err != nil {
		return err
	}

	if casem.Status == "1" {
		return errors.New("案件已有成交紀錄")
	}

	tx := dao.GormSession.Begin()

	msg := model.MsgRecord{
		AccountFrom: account,
		AccountTo:   m.Quoter,
		Message:     fmt.Sprintf("%s-=%s-=%d-=%d", m.CaseId, m.Title, m.PriceS, m.PriceE),
		IsSystem:    "2",
	}
	if err := tx.Model(&model.MsgRecord{}).Create(&msg).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 更新報價紀錄 1:已成交
	quote := model.Quote{
		Deal: 1,
	}
	if err := tx.Model(&model.Quote{}).Where("case_id = ? and account = ?", m.CaseId, m.Quoter).Updates(quote).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 更新案件狀態：1已成交
	c := model.Casem{
		Status: "1",
	}
	if err := tx.Model(&model.Casem{}).Where("case_id = ?", m.CaseId).Updates(c).Error; err != nil {
		tx.Rollback()
		return err
	} else {
		tx.Commit()
		return nil
	}
}
