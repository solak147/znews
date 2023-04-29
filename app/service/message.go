package service

import (
	"znews/app/dao"
	"znews/app/model"
)

func GetMsgRecord(account string) ([]model.MsgRec, error) {
	var msgArr []model.MsgRec

	// 执行 SQL 查询语句
	query := `select account, 
					(SELECT MAX(created_at) FROM msg_records WHERE (account_from = ? and account_to = a.account) or (account_from = a.account and account_to = ? )) crtDte, 
					(SELECT message FROM msg_records WHERE (account_from = ? and account_to = a.account) or (account_from = a.account and account_to = ? ) ORDER BY created_at DESC LIMIT 1) message,
					(SELECT COUNT(*) FROM msg_records WHERE account_from = ? and account_to = a.account AND is_read = '0') notReadCnt  
				FROM (
					SELECT account_to account
					FROM msg_records WHERE account_from = ? GROUP BY account_to 
					UNION 
					SELECT account_from account
					FROM msg_records WHERE account_to = ? GROUP BY account_from
				)  a
				group by account`
	rows, err := dao.DbSession.Query(query, account, account, account, account, account, account, account)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var msg model.MsgRec
		if err := rows.Scan(&msg.Account, &msg.CrtDte, &msg.Message, &msg.NotReadCnt); err != nil {
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
