package geojson

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var (
	Verbose bool      = false
	Output  io.Writer = os.Stdout
)

type FeatureCollection struct {
	Type     string   `json:"type"`
	Name     string   `json:"name"`
	Crs      Crs      `json:"crs"`
	Features Features `json:"features"`
}

type Crs struct {
	Type       string        `json:"type"`
	Properties CrsProperties `json:"properties"`
}

type CrsProperties struct {
	Name string `json:"name"`
}

type Properties struct {
	// 医療機関分類（P04_001）
	P04001 int `json:"P04_001"`
	// 施設名称（P04_002）
	P04002 string `json:"P04_002"`
	// 所在地（P04_003）
	P04003 string `json:"P04_003"`
	// 診療科目１（P04_004）
	P04004 string `json:"P04_004"`
	// 診療科目２（P04_005）
	P04005 string `json:"P04_005"`
	// 診療科目３（P04_006）
	P04006 string `json:"P04_006"`
	// 開設者分類（P04_007）
	P04007 int `json:"P04_007"`
	// 病床数（P04_008）
	P04008 int `json:"P04_008"`
	// 救急告示病院（P04_009）
	P04009 int `json:"P04_009"`
	// 災害拠点病院（P04_010）
	P04010 int `json:"P04_010"`
}

//
type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

//
type Feature struct {
	Type       string     `json:"type"`
	Properties Properties `json:"properties"`
	Geometry   Geometry   `json:"geometry"`
}

type Features []Feature

func (fs Features) Dump() {
	// https://developers.google.com/maps/documentation/urls/get-started#search-action
	for i, f := range fs {
		fmt.Fprintf(Output, "%v, %v, %v, %v, https://www.google.com/maps/search/?api=1&query=%v,%v",
			i,
			f.Properties.P04002,
			f.Geometry.Coordinates[0], f.Geometry.Coordinates[1],
			f.Geometry.Coordinates[1], f.Geometry.Coordinates[0],
		)
	}
}

type SqlOption int

const (
	Transaction SqlOption = 0
	TableLock             = 1
	AutoCommit            = 2
	MultiValue            = 3
)

type PrePostSQL struct {
	Pre  string
	Post string
}

var PrePostSQLs = map[SqlOption]PrePostSQL{
	Transaction: {
		// local env, 1m45.372s
		Pre:  "start transaction;\n",
		Post: "commit;\n",
	},
	TableLock: {
		// local env, more than 43m ...
		Pre: `
        /*!50503 SET NAMES utf8mb4 */;
        /*!40103 SET TIME_ZONE='+00:00' */;
        /*!40014 SET UNIQUE_CHECKS=0 */;
        /*!40014 SET FOREIGN_KEY_CHECKS=0 */;
        /*!40101 SET SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
        LOCK TABLES hospital WRITE;
        /*!40000 ALTER TABLE hospital DISABLE KEYS */;
        `,
		Post: `
        /*!40000 ALTER TABLE hospital ENABLE KEYS */;
        UNLOCK TABLES;
        `,
	},
	AutoCommit: {
		// local env,  44m29.460s
		Pre:  "",
		Post: "",
	},
	MultiValue: {
		// local env, 1m9.798s
		Pre:  "start transaction;\n",
		Post: "commit;\n",
	},
}

func (fs Features) printSQL(sqlOption SqlOption) error {
	fmt.Fprint(Output, PrePostSQLs[sqlOption].Pre)
	insert := "insert into hospital (name, location) values"
	eol := ";\n"
	for i, f := range fs {
		switch i {
		case 0:
			if sqlOption == MultiValue {
				eol = ",\n"
			}
		case 1:
			if sqlOption == MultiValue {
				insert = ""
			}
		case len(fs) - 1:
			eol = ";\n"
		}

		fmt.Fprintf(Output, "%v('%v', st_geomfromtext('point(%v %v)', 4326))%v",
			insert,
			strings.ReplaceAll(f.Properties.P04002, "'", "\\'"),
			f.Geometry.Coordinates[1], f.Geometry.Coordinates[0],
			eol)
	}
	fmt.Fprint(Output, PrePostSQLs[sqlOption].Post)

	return nil
}

func (fs Features) printTsv(sep string) error {
	for _, f := range fs {
		s := []string{
			strings.ReplaceAll(f.Properties.P04002, "'", "\\'"),
			strconv.FormatFloat(f.Geometry.Coordinates[1], 'f', -1, 32),
			strconv.FormatFloat(f.Geometry.Coordinates[0], 'f', -1, 32),
		}
		fmt.Fprintln(Output, strings.Join(s, sep))
	}
	return nil
}

func (fs Features) Print(sqlOption SqlOption, args ...interface{}) (err error) {
	switch sqlOption {
	case MultiValue:
		sep := "\t"
		if args != nil {
			sep = args[0].(string)
		}
		err = fs.printTsv(sep)
	default:
		err = fs.printSQL(sqlOption)
	}
	return err
}

func NewFeatures(jsonFile string) (*FeatureCollection, error) {

	file, err := os.Open(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("os.Open: %w", err)
	}
	defer file.Close()

	if Verbose {
		fmt.Fprintf(Output, "Successfully Opened %v\n", jsonFile)
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var fc FeatureCollection

	if err := json.Unmarshal(bytes, &fc); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}

	if Verbose {
		fc.Features.Dump()
	}

	return &fc, nil
}
