// Copyright 2022 The wangkai. ALL rights reserved.

/*
Package models
*/
package models

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var DB *sql.DB

// InitDBConnectionPool 初始化DB连接池
func InitDBConnectionPool(dataSourceName string) (*sql.DB, error) {
	// 连接 db
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(100)                    //  设置连接数总数, 需要根据实际业务来测算, 应小于 mysql.max_connection (应该远远小于), 后续根据指标进行调整
	db.SetMaxIdleConns(100)                    //  设置最大空闲连接数, 该数值应该小于等于 SetMaxOpenConns 设置的值,  需要根据实际业务来测算, 后续根据指标进行调整
	db.SetConnMaxLifetime(0)                   // 设置连接最大生命周期, 默认为 0(不限制), 我不建议设置该值, 只有当 mysql 服务器出现问题, 会导致连接报错, 恢复后可以自动恢复正常, 而我们配置了时间也不能卡住出问题的时间, 配置小还不如使用 SetConnMaxIdleTime 来解决
	db.SetConnMaxIdleTime(86400 * time.Second) // 设置空闲状态最大生命周期, 该值应小于 mysql.wait_timeout 的值, 以避免被服务端断开连接, 产生报错影响业务。

	// 创建连接池
	DB = db

	return DB, nil
}

// SendPhoneNumberList SendPhoneNumber table send_phone_number_list
type SendPhoneNumberList struct {
	RequestId string `db:"request_id"`
	ProjectId string `db:"project_id"`
	AreaCode  string `db:"area_code"`
	Number    string `db:"number"`
	Status    int    `db:"status"`
	CancelAt  string `db:"cancel_at"`
	SmsCode   string `db:"sms_code"`
}

// Insert 插入数据
func (s *SendPhoneNumberList) Insert() (int64, error) {
	// 该死必须全部匹配
	stmt, _ := DB.Prepare("INSERT INTO send_phone_number_list ( `request_id`, `area_code`, `number`, `status`, `project_id` )" +
		"VALUES (?, ?, ?, ?, ?) ")
	dateTime := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println(dateTime)
	res, err := stmt.Exec(
		s.RequestId,
		s.Number,
		s.AreaCode,
		s.Status,
		s.ProjectId,
	)

	if err == nil {
		insertId, _ := res.LastInsertId()
		return insertId, nil
	} else {
		return 0, err
	}
}

// UpdateSmsSendSuccessStatus 修改成功短信状态
func (s *SendPhoneNumberList) UpdateSmsSendSuccessStatus() (rowsAffected int64) {
	stmt, _ := DB.Prepare("UPDATE `send_phone_number_list` SET sms_code = ?, status = 1 WHERE request_id = ?")
	res, err := stmt.Exec(
		s.SmsCode,
		s.RequestId,
	)

	if err == nil {
		rowsAffected, _ = res.RowsAffected()
	} else {
		panic(err)
	}

	return rowsAffected
}

// CancelSmsSend 取消短信发送
func (s *SendPhoneNumberList) CancelSmsSend(requestId string) (rowsAffected int64) {
	stmt, _ := DB.Prepare("UPDATE `send_phone_number_list` SET status = 2, cancel_at = ? WHERE request_id = ?")
	res, err := stmt.Exec(
		time.Now().Format("2006-01-02 15:04:05"),
		requestId,
	)

	if err == nil {
		rowsAffected, _ = res.RowsAffected()
	} else {
		panic(err)
	}

	return rowsAffected
}

// GetLastActivePhoneNumber 获取最新一条可用手机号码
func GetLastActivePhoneNumber(projectId int) (SendPhoneNumberList, error) {
	var queryData SendPhoneNumberList

	fmt.Println(projectId)
	err := DB.QueryRow("SELECT request_id RequestId, project_id ProjectId, area_code AreaCode, number Number, status Status, cancel_at CancelAt,  "+
		"sms_code SmsCode FROM `send_phone_number_list` "+
		"WHERE `status` = 0 AND project_id = ? ORDER BY id DESC LIMIT 1", projectId).
		Scan(&queryData.RequestId, &queryData.ProjectId, &queryData.AreaCode, &queryData.Number, &queryData.Status, &queryData.CancelAt, &queryData.SmsCode)

	if err != nil {
		if err == sql.ErrNoRows {
			return queryData, nil
		}
		return queryData, err
	}

	return queryData, nil
}
