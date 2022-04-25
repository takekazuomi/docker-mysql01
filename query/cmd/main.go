package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "embed"

	. "query/hospital"

	_ "github.com/go-sql-driver/mysql"

	flag "github.com/spf13/pflag"
)

var (
	datasource = os.ExpandEnv("${MYSQL_USER}:${MYSQL_PASSWORD}@tcp(${MYSQL_HOST}:3306)/geo?charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=true")
	dialect    = "mysql"
)

func Open() (*sql.DB, error) {
	// `sql.Open` function は、新しい `*sql.DB` インスタンスを作成する. driver名と、データーベースのURIを指定する。
	db, err := sql.Open(dialect, datasource)
	if err != nil {
		return db, fmt.Errorf("driver: %v, datasource:%v: %w", dialect, datasource, err)
	}
	return db, nil
}

var (
	lat      float32
	lng      float32
	distance float32
)

func init() {
	flag.Float32VarP(&lat, "lat", "a", 35.6884204226699, "lat")
	flag.Float32VarP(&lng, "lng", "n", 139.72515649841105, "lng")
	flag.Float32VarP(&distance, "distance", "d", 1, "distance")
}

func main() {
	db, err := Open()
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}
	defer db.Close()

	// database instanceのコネクションを確認ために、`Ping` メソッドを呼ぶ。
	// もしエラーが帰ってこなければ接続が成功している。
	if err := db.Ping(); err != nil {
		log.Fatalf("unable to reach database: %v", err)
	}

	r, err := NewHospitalResultsByNeighborhood(db, lat, lng, distance)
	if err == nil {
		r.PrintTable()
	} else {
		log.Fatalf("query error: %v", err)
	}
}
