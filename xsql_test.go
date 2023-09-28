package xsql_test

import (
	"fmt"
	"testing"
	"xsql"
)

func TestZeroGroup(t *testing.T) {
	dbConfig := &xsql.Config{
		Host:     "127.0.0.1",
		Port:     3306,
		UserName: "root",
		Password: "",
		DBName:   "test",
		Charset:  "utf8",
	}

	db := xsql.NewXSQL(dbConfig)

	//db.Table("bag_0000").Primary("uid").Fields("uid", "item", "expire", "itime").Where("uid in (6,8)")
	//rawsql := db.GenRawSQL()
	//fmt.Println("rawsql:::::::", rawsql)

	//data, err := db.Table("bag_0000").Primary("uid").Fields("uid", "item", "expire", "itime").Where("uid in (6,8)").Query()
	data, err := db.Table("bag_0000").Primary("uid").Fields("item", "expire", "itime").Where("uid in (6)").QueryRow()
	fmt.Println("err:::::::", err)
	fmt.Println("data:::::::", data)

}
