package service

import (
	"znews/app/dao"
)

type msgRec struct {
	account string
	message string
	crtDte  string
}

func GetMsgRecord(account string) ([]msgRec, error) {
	var msgArr []msgRec

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
		var msg msgRec
		if err := rows.Scan(&msg.account, &msg.message, &msg.crtDte); err != nil {
			return nil, err
		}
		msgArr = append(msgArr, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return msgArr, nil

	// if err := dao.GormSession.Select("*").Where("account_from=? or account_to=?", account, account).Order("created_at ASC").Find(&msg).Error; err != nil {
	// 	return nil, err
	// } else {
	// 	return msg, nil
	// }
}
