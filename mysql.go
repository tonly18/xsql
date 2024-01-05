package xsql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cast"
	"runtime"
	"slices"
	"strings"
	"sync"
	"time"
)

var once sync.Once
var dbConn *sql.DB

type XSQL struct {
	ctx context.Context
	db  *sql.DB

	table     string   //表名
	fields    []string //字段
	values    []any    //字段-值
	where     []string //条件
	order     string   //排序
	group     string   //分组
	have      string   //分组条件
	leftJoin  string   //左关联
	rightJoin string   //右关联
	on        string   //on条件
	sql       string   //sql
}

// NewXSQL
func NewXSQL(ctx context.Context, config *Config) *XSQL {
	xsql := &XSQL{
		ctx:    ctx,
		fields: make([]string, 0, 20),
		values: make([]any, 0, 20),
		where:  make([]string, 0, 5),
	}
	once.Do(func() {
		if err := xsql.connect(config); err != nil {
			panic(fmt.Errorf(`[new xsql] once.do error: %w`, err))
		}
	})
	if dbConn == nil {
		if err := xsql.connect(config); err != nil {
			panic(fmt.Errorf(`[new xsql] dbConn is nil, error: %w`, err))
		}
	}
	xsql.db = dbConn

	//return
	return xsql
}

// connect
func (d *XSQL) connect(config *Config) error {
	if config.Charset == "" {
		config.Charset = "utf8"
	}
	if config.MaxLifetime == 0 {
		config.MaxLifetime = time.Second * maxLifeTime
	}
	if config.MaxIdleTime == 0 {
		config.MaxIdleTime = time.Second * maxIdleTime
	}
	if config.MaxOpenConns == 0 {
		config.MaxOpenConns = maxOpenConns
	}
	if config.MaxIdleConns == 0 {
		config.MaxIdleConns = maxIdleConns
	}

	var err error
	dbConn, err = sql.Open("mysql", fmt.Sprintf(`%s:%s@tcp(%s:%d)/%s?charset=%s`, config.UserName, config.Password, config.Host, config.Port, config.DBName, config.Charset))
	if err != nil {
		return fmt.Errorf(`[connect] sql.Open error: %w`, err)
	}
	if dbConn == nil {
		return fmt.Errorf(`%w`, errors.New("[connect] db conn is nil"))
	}
	if err = dbConn.Ping(); err != nil {
		return fmt.Errorf(`[connect] ping error: %w`, err)
	}

	//设置连接池里的连接最大存活时长(通常比mysql服务器wait_timeout小)
	dbConn.SetConnMaxLifetime(config.MaxLifetime)
	//设置连接池里的连接最大空闲时长(连接每次被使用后,持续空闲时长会被重置,从0开始从新计算)
	dbConn.SetConnMaxIdleTime(config.MaxIdleTime)
	//设置连接池最多同时打开的连接数,如果n<=0(默认值为0,无限制)
	dbConn.SetMaxOpenConns(config.MaxOpenConns)
	//设置连接池里最大空闲连接数,如果n<=0(则不保留任何空闲连接), 必须要比maxOpenConns小
	dbConn.SetMaxIdleConns(config.MaxIdleConns)

	//Finalizer
	runtime.SetFinalizer(dbConn, func(conn *sql.DB) {
		conn.Close()
	})

	//return
	return nil
}

// Table 字段
func (d *XSQL) Table(table string) *XSQL {
	d.table = table

	//return
	return d
}

// Fields 字段
func (d *XSQL) Fields(fields ...string) *XSQL {
	if len(fields) > 0 {
		d.fields = append(d.fields, fields...)
	}

	//return
	return d
}

// Where 条件
func (d *XSQL) Where(condition string) *XSQL {
	if condition != "" {
		if len(d.where) == 0 {
			d.where = append(d.where, condition)
		} else {
			d.where = append(d.where, " AND ", condition)
		}
	}

	//return
	return d
}

// ORWhere 条件
func (d *XSQL) ORWhere(condition string) *XSQL {
	if condition != "" {
		if len(d.where) == 0 {
			d.where = append(d.where, condition)
		} else {
			d.where = append(d.where, " OR ", condition)
		}
	}

	//return
	return d
}

// GroupBy 分组
func (d *XSQL) GroupBy(group string) *XSQL {
	d.group = group

	//return
	return d
}

// Having 分组条件
func (d *XSQL) Having(having string) *XSQL {
	d.have = having

	//return
	return d
}

// LeftJoin 关联
func (d *XSQL) LeftJoin(join string) *XSQL {
	d.leftJoin = join

	//return
	return d
}

// RightJoin 关联
func (d *XSQL) RightJoin(join string) *XSQL {
	d.rightJoin = join

	//return
	return d
}

// ON 关联
func (d *XSQL) ON(on string) *XSQL {
	d.on = on

	//return
	return d
}

// OrderBy 排序
func (d *XSQL) OrderBy(order string) *XSQL {
	d.order = order

	//return
	return d
}

// QueryRow 查询单条数据
func (d *XSQL) QueryRow() (map[string]any, error) {
	data, err := d.Query()
	if err != nil {
		return nil, err
	}

	return data[0], nil
}

// Query 查询数据
func (d *XSQL) Query() ([]map[string]any, error) {
	defer d.RestSQL()

	//生成SQL
	rawsql := d.GenRawSQL()

	//QUERY
	rows, err := d.db.Query(rawsql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if len(d.fields) == 0 {
		d.fields, _ = rows.Columns()
	}

	//字段 - 数据
	data := make([]map[string]any, 0, 20)
	entity := genEntity(len(d.fields))
	for rows.Next() {
		if err := rows.Scan(entity...); err != nil {
			return nil, err
		}
		data = append(data, genRecord(entity, d.fields))
	}
	if len(data) == 0 {
		return nil, sql.ErrNoRows
	}

	//return
	return data, nil
}

// QueryMap 查询数据
func (d *XSQL) QueryMap(field string) (map[int]map[string]any, error) {
	defer d.RestSQL()

	//field
	if len(d.fields) > 0 {
		if false == slices.Contains(d.fields, field) {
			d.fields = append(d.fields, field)
		}
	}

	//generate sql
	rawsql := d.GenRawSQL()

	//query
	rows, err := d.db.Query(rawsql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//field
	if len(d.fields) == 0 {
		d.fields, _ = rows.Columns()
	}

	//字段 - 数据
	data := make(map[int]map[string]any, 50)
	entity := genEntity(len(d.fields))
	for rows.Next() {
		if err := rows.Scan(entity...); err != nil {
			return nil, err
		}
		record := genRecord(entity, d.fields)
		data[cast.ToInt(record[field])] = record
	}
	if len(data) == 0 {
		return nil, sql.ErrNoRows
	}

	//return
	return data, nil
}

// GenRawSQL 生成查询SQL
func (d *XSQL) GenRawSQL() string {
	//modify|insert|delete sql
	if len(d.sql) > 0 {
		return d.sql
	}

	//query sql
	var rawsql strings.Builder
	if len(d.fields) == 0 {
		rawsql.WriteString(fmt.Sprintf(`SELECT * FROM %v`, d.table))
	} else {
		rawsql.WriteString(fmt.Sprintf(`SELECT %v FROM %v`, strings.Join(d.fields, ","), d.table))
	}
	if len(d.where) > 0 {
		rawsql.WriteString(fmt.Sprintf(` WHERE %v`, strings.Join(d.where, "")))
	}
	if d.group != "" {
		rawsql.WriteString(fmt.Sprintf(` GROUP BY %v`, d.order))
	}
	if d.have != "" {
		rawsql.WriteString(fmt.Sprintf(` HAVING %v`, d.have))
	}
	if d.leftJoin != "" {
		rawsql.WriteString(fmt.Sprintf(` LEFT JOIN %v`, d.leftJoin))
	}
	if d.rightJoin != "" {
		rawsql.WriteString(fmt.Sprintf(` RIGHT JOIN %v`, d.rightJoin))
	}
	if d.on != "" {
		rawsql.WriteString(fmt.Sprintf(` ON %v`, d.on))
	}
	if d.order != "" {
		rawsql.WriteString(fmt.Sprintf(` ORDER BY %v`, d.order))
	}

	//return
	return rawsql.String()
}

// Insert 插入数据
func (d *XSQL) Insert(params map[string]any) *XSQL {
	for k, v := range params {
		d.fields = append(d.fields, k)
		d.values = append(d.values, v)
	}
	if len(params) > 0 {
		d.sql = fmt.Sprintf("INSERT INTO %v(%v) VALUES (%v)", d.table, strings.Join(d.fields, ","), strings.Repeat(",?", len(d.fields))[1:])
	}

	//return
	return d
}

// Modify 修改数据
func (d *XSQL) Modify(params map[string]any) *XSQL {
	for k, v := range params {
		d.fields = append(d.fields, fmt.Sprintf(`%v=?`, k))
		d.values = append(d.values, v)
	}
	if len(params) > 0 {
		d.sql = fmt.Sprintf(`UPDATE %v SET %v`, d.table, strings.Join(d.fields, ","))
	}
	if len(d.where) > 0 {
		d.sql = fmt.Sprintf(`%v WHERE %v`, d.sql, strings.Join(d.where, ""))
	}

	//return
	return d
}

// Delete 删除数据
func (d *XSQL) Delete() *XSQL {
	d.sql = fmt.Sprintf(`DELETE FROM %v`, d.table)
	if len(d.where) > 0 {
		d.sql = fmt.Sprintf(`%v WHERE %v`, d.sql, strings.Join(d.where, ""))
	}

	//return
	return d
}

// Exec 执行SQL
func (d *XSQL) Exec() (sql.Result, error) {
	defer d.RestSQL()

	stmt, err := d.db.Prepare(d.sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(d.values...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// RestSQL
func (d *XSQL) RestSQL() {
	d.sql = ""

	if d.table != "" {
		d.table = ""
	}
	if len(d.fields) > 0 {
		d.fields = make([]string, 0, 20)
	}
	if len(d.values) > 0 {
		d.values = make([]any, 0, 20)
	}
	if len(d.where) > 0 {
		d.where = make([]string, 0, 5)
	}
	if d.group != "" {
		d.group = ""
	}
	if d.have != "" {
		d.have = ""
	}
	if d.order != "" {
		d.order = ""
	}
	if d.on != "" {
		d.on = ""
	}
	if d.leftJoin != "" {
		d.leftJoin = ""
	}
	if d.rightJoin != "" {
		d.rightJoin = ""
	}
}

// Transaction
func (d *XSQL) Begin() (*sql.Tx, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return nil, err
	}
	return tx, nil
}
