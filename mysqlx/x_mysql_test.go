package mysqlx

import (
	"log"
	"testing"
)

var testSetting = &MysqlSetting{
	RunMode: "develop",

	Host:     "localhost:3306",
	User:     "xiaoqucloud",
	Password: "123456",

	MySqlMaxIdleConn: 10,
	MySqlMaxOpenConn: 50,
	MySqlConnMaxLife: 60,
}

const testDbName string = "xiaoqucloud_platform"

func TestMysql(t *testing.T) {
	Setup(testSetting)
	db, err := GetDb(testDbName)
	if err != nil {
		t.Error(err)
	}
	var testResult map[string]interface{}
	db.Table("users").Find(&testResult)
	log.Println(testResult)
}
