package hospital

import (
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/olekukonko/tablewriter"

	_ "embed"

	_ "github.com/go-sql-driver/mysql"
)

var (
	datasource = os.ExpandEnv("${MYSQL_USER}:${MYSQL_PASSWORD}@tcp(${MYSQL_HOST}:3306)/geo01?charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=true")
	dialect    = "mysql"
	Output     = os.Stdout
)

type HospitalResult struct {
	Id       int64
	Name     string
	Location string
	Distance float64
}

func NewHospitalResultsByNeighborhood(db *sql.DB, lat float32, lng float32, distance float32) (*HospitalResults, error) {

	// MySQLの場合、パラメータは"?"記号で指定し、変数は、同じ順序で引数として追加する必要がある
	// sql libraryは、GoからSQLへの型変換をdriverに基づいて行う
	rows, err := db.Query("call neighborhood(?, ?, ?)", lat, lng, distance)
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}
	defer rows.Close()

	rs := HospitalResults{}

	for rows.Next() {
		hr := HospitalResult{}

		if err := rows.Scan(&hr.Id, &hr.Name, &hr.Location, &hr.Distance); err != nil {
			return nil, fmt.Errorf("could not scan row: %w", err)
		}

		rs = append(rs, hr)
	}

	return &rs, nil
}

type HospitalResults []HospitalResult

func toStrings(i interface{}) (name *[]string, value *[]string) {

	n := make([]string, 0)
	s := make([]string, 0)

	v := reflect.ValueOf(i)
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := v.Type().Field(i)
		value := v.Field(i)

		// 必要な型だけ文字列に変換
		n = append(n, field.Name)
		switch value := value.Interface().(type) {
		case int:
			s = append(s, strconv.Itoa(value))
			break
		case int16:
			s = append(s, strconv.FormatInt(int64(value), 10))
			break
		case int32:
			s = append(s, strconv.FormatInt(int64(value), 10))
			break
		case int64:
			s = append(s, strconv.FormatInt(value, 10))
			break
		case string:
			s = append(s, value)
			break
		case float64:
			s = append(s, fmt.Sprintf("%f", value))
		}
	}

	return &n, &s
}

func (s *HospitalResults) PrintTable() {
	n := make([]string, 0)
	d := make([][]string, 0)

	for _, v := range *s {
		i, j := toStrings(v)
		if len(n) == 0 {
			n = *i
		}
		d = append(d, *j)
	}

	printTable(n, d)
}

func printTable(header []string, data [][]string) {
	table := tablewriter.NewWriter(Output)
	table.SetHeader(header)

	for _, v := range data {
		table.Append(v)
	}
	table.Render() // Send output
}
