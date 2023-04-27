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
					(SELECT message FROM msg_records WHERE (account_from = ? and account_to = a.account) or (account_from = a.account and account_to = ? ) ORDER BY created_at DESC LIMIT 1) message 
				FROM (
					SELECT account_to account
					FROM msg_records WHERE account_from = ? GROUP BY account_to 
					UNION 
					SELECT account_from account
					FROM msg_records WHERE account_to = ? GROUP BY account_from
				)  a
				group by account`
	rows, err := dao.DbSession.Query(query, account, account, account, account, account, account)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var msg model.MsgRec
		if err := rows.Scan(&msg.Account, &msg.CrtDte, &msg.Message); err != nil {
			return nil, err
		}
		msgArr = append(msgArr, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return msgArr, nil

}
