package xsql_test

import (
	"context"
	"fmt"
	"github.com/tonly18/xsql"
	"testing"
)

func TestXSQL(t *testing.T) {
	dbConfig := &xsql.Config{
		Host:     "127.0.0.1",
		Port:     3306,
		UserName: "root",
		Password: "",
		DBName:   "test",
		Charset:  "utf8",
	}
	db := xsql.NewXSQL(context.Background(), dbConfig)

	//select
	db.Table("bag_0000").Fields("uid", "item", "expire", "itime").Where("uid in (6,8)")
	rawsql := db.GenRawSQL()
	fmt.Println("rawsql:::::::", rawsql)

	//data, err := db.Table("bag_0000").Fields("uid", "item2", "expire", "itime").Where("uid in (6,8)").Query()
	//data, err := db.Table("bag_0000").Fields("item", "expire", "itime").Where("uid in (6)").QueryRow()
	//data, err := db.Table("bag_0000").Fields("item").Where("uid in (6)").QueryMap("uid")
	//data, err := db.Table("bag_0000").Where("uid in (6)").Query()
	//sql := db.Table("bag_0000").Fields("uid,item").Where("uid in (6)").GenRawSQL()
	//fmt.Println("err:::::::", err)
	//fmt.Println("data:::::::", data)
	//fmt.Println("sql:::::::", sql)

	//if errors.Is(err, sql.ErrNoRows) {
	//	fmt.Println("sql.ErrNoRows::::::", sql.ErrNoRows)
	//}

	rows, err := db.RawQuery("SELECT `item`, count(item) as num FROM `bag_0000` group by `item`")
	fmt.Println("err:::::", err)
	fmt.Println("rows:::::", rows)
	for rows.Next() {
		//var item, expire string
		//var uid, itime int
		var item string
		var num int
		if err := rows.Scan(&item, &num); err != nil {
			fmt.Println("err::::", err)
			continue
		}
		fmt.Println("uid, item, expire:::", item, num)
	}

	//Insert
	//result, err := db.Table("bag_0000").Insert(map[string]any{
	//	"uid":    18,
	//	"item":   "item-18",
	//	"expire": "expire-18",
	//	"itime":  1988120018,
	//}).Exec()
	//fmt.Println("err::::::::", err)
	//count, err := result.RowsAffected()
	//fmt.Println("result-count,err::::::::", count, err)
	//newId, err := result.LastInsertId()
	//fmt.Println("result-newId,err::::::::", newId, err)

	result, err := db.RawExec("insert into bag_0000(uid,item,expire,itime) values(101, \"item-101\", 123456, 789)")
	fmt.Println("err::::", err)
	n, _ := result.RowsAffected()
	fmt.Println("result::::", n)

	//modify
	//result, err := db.Table("bag_0001").Where("uid=17").Modify(map[string]any{
	//	"item":   "item-17-m",
	//	"expire": "expire-17-m",
	//}).Exec()
	//fmt.Println("err::::::::", err)
	//count, err := result.RowsAffected()
	//fmt.Println("result-count,err::::::::", count, err)
	//newId, err := result.LastInsertId()
	//fmt.Println("result-newId,err::::::::", newId, err)

	//delete
	//result, err := db.Table("bag_0000").Where("uid=4").Delete().Exec()
	//fmt.Println("err::::::::", err)
	//count, err := result.RowsAffected()
	//fmt.Println("result-count,err::::::::", count, err)
	//newId, err := result.LastInsertId()
	//fmt.Println("result-newId,err::::::::", newId, err)

	//Transaction
	//tx, err := db.Begin()
	//fmt.Println("err::::::::", err)
	//fmt.Println("tx:::::::::", tx)
}
