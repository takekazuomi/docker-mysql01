package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	g "github.com/takekazuomi/docker-mysql01/import/geojson"

	flag "github.com/spf13/pflag"
)

var (
	jsonFile  string
	help      bool
	verbose   bool
	sqlOption int32
)

func main() {
	flag.BoolVarP(&help, "help", "h", false, "show help message")
	flag.BoolVarP(&verbose, "verbose", "v", false, "show verbose message")
	flag.StringVarP(&jsonFile, "json", "j", "_data/P04-20_11_GML/small.geojson", "source geo json file")
	flag.Int32VarP(&sqlOption, "sql", "s", 0, "sql option, 0 is transaction, 1 is table lock, 2 is auto commit, 3 is multi value insert")
	flag.Parse()

	if help {
		flag.PrintDefaults()
		return
	}

	file, err := os.Open(jsonFile)

	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	if verbose {
		fmt.Printf("Successfully Opened %v\n", jsonFile)
	}

	bytes, _ := ioutil.ReadAll(file)

	var fc g.FeatureCollection

	json.Unmarshal(bytes, &fc)

	if verbose {
		fc.Features.Dump()
	}

	fc.Features.PrintSQL(g.SqlOption(sqlOption))

}
