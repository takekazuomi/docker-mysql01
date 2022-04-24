package main

import (
	"fmt"

	. "github.com/takekazuomi/docker-mysql01/import/geojson"

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
	flag.Int32VarP(&sqlOption, "sql", "s", 0, "sql option, 0 is transaction, 1 is table lock, 2 is auto commit, 3 is multi value insert, 4 is tsv")
	flag.Parse()

	if help {
		flag.PrintDefaults()
		return
	}

	Verbose = verbose

	fc, err := NewFeatures(jsonFile)
	if err != nil {
		fmt.Printf("g.NewFeatures: %v", err)
		return
	}

	if verbose {
		fc.Features.Dump()
	}

	fc.Features.Print(SqlOption(sqlOption))

}
