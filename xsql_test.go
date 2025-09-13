package xsql_test

import (
	"context"
	"fmt"
	"github.com/tonly18/xsql"
	"math/rand"
	"testing"
	"time"
)

func TestXSQL(t *testing.T) {
	dbConfig := &xsql.Config{
		Host:     "127.0.0.1",
		Port:     3306,
		UserName: "root",
		Password: "123456",
		DBName:   "test",
		Charset:  "utf8",
	}
	db := xsql.NewXSQL(context.Background(), dbConfig)

	//select
	//db.Table("bag_0000").Fields("uid", "item", "expire", "itime").Where("uid in (6,8)")
	//rawsql := db.GenRawSQL()
	//fmt.Println("rawsql:::::::", rawsql)

	data1, err := db.Table("bag_0000").Fields("uid", "item", "expire", "itime").Where("uid in (2)").Query()
	fmt.Println("data1::::::", data1)

	data2, err := db.Table("bag_0000").Fields("uid", "item", "expire", "itime").Where("uid in (6, 8, 100, 2)").QueryMap("uid")
	fmt.Println("data2::::::", data2)
	for k, v := range data2 {
		fmt.Println("\nk-v::::::", k)
		for key, val := range v {
			fmt.Println("	key-val::::::", key, string(val))
		}
	}

	fmt.Println("over！")

	result, err := db.RawExec("insert into bag_0000(uid,item,expire,itime) values(3, \"item-3\", 123456, 333)")
	fmt.Println("err::::", err)
	if result != nil {
		n, _ := result.RowsAffected()
		fmt.Println("result::::", n)
	}

	return

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

	rows, err := db.RawQuery("SELECT * FROM `bag_0000` where uid=?", 100)
	fmt.Println("err:::::", err)
	defer rows.Close()
	for rows.Next() {
		var item, expire string
		var uid, itime int
		if err := rows.Scan(&uid, &item, &expire, &itime); err != nil {
			fmt.Println("err::::", err)
			continue
		}
		fmt.Println("uid, item, expire:::", uid, item, expire, itime)
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

// 并发安全测试
func TestConcurrencySafety(t *testing.T) {
	db := xsql.NewXSQL(context.Background(), &xsql.Config{
		Host:     "127.0.0.1",
		Port:     3306,
		UserName: "root",
		Password: "123456",
		DBName:   "test",
		Charset:  "utf8",
	})

	for i := 0; i < 100; i++ {
		go func(x int) {
			id := rand.Intn(6) + 1
			where := fmt.Sprintf(`id=%d`, id)

			rawsql := db.Table("employees").Fields("id", "name", "age").Where(where).GenRawSQL()
			fmt.Println("rawsql:", rawsql)

			data, err := db.Table("employees").Fields("id", "name", "age").Where(where).Query()
			if err != nil {
				panic(err)
			}
			for _, v := range data {
				fmt.Println("x:", x, ", id:", string(v["id"]), " name:", string(v["name"]), " age:", string(v["age"]))
			}
		}(i)
	}

	time.Sleep(time.Second * 3)
	fmt.Println("main-over")
}
