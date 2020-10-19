package main

import (
	"database/sql"
	"fmt"
	_ "github.com/ClickHouse/clickhouse-go"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"reflect"
)

const my = "sdtd:OdeMucbPBkaVcyST@tcp(192.168.11.10:3306)/sdtd_td?timeout=5s&readTimeout=1m&writeTimeout=1m"
const clickhouse = "tcp://192.168.11.37:9000?username=default&password=&database=sdtd_td&timeout=5s&read_timeout=1m&write_timeout=1m"

func main() {
	table := "dsp_test"
	dbc, _ := sql.Open("clickhouse", clickhouse)

	exist := false
	query := dbc.QueryRow(fmt.Sprintf("exists table `%s`", table))
	err := query.Scan(&exist)
	switch {
	case err != nil:
		return
	case !exist:
		err = fmt.Errorf("数据表[%s]不存在", table)
		panic(err)
	}

	rows, err := dbc.Query(fmt.Sprintf("describe `%s`", table))
	if err != nil {
		panic(err)
		return
	}
	defer dbc.Close()

	type Row struct {
		Name              string
		Type              string
		DefaultType       string
		DefaultExpression string
		Comment           string
		CodecExpression   string
		TtlExpression     string
	}
	rowv := Row{}
	rowp := make([]interface{}, 0)

	rv := reflect.ValueOf(&rowv).Elem()
	for i := 0; i < rv.NumField(); i++ {
		rowp = append(rowp, rv.Field(i).Addr().Interface())
	}

	var all = make(map[string]string, 0)
	for rows.Next() {
		err := rows.Scan(rowp...)
		if err != nil {
			panic(err)
		}

		all[rowv.Name] = rowv.Type
	}

	log.Printf("%s", all)
}
