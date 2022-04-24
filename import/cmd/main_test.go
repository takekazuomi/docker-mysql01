package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/takekazuomi/docker-mysql01/import/geojson"

	flag "github.com/spf13/pflag"
)

func BenchmarkMain(b *testing.B) {
	type args []string
	tests := []struct {
		name string
		args args
	}{
		//		{name: "P04-20.geojson", args: []string{"test", "-s", "0", "-j", "../geojson/testdata/one.json"}},
		//		{name: "P04-20.geojson", args: []string{"test", "-s", "1", "-j", "../geojson/testdata/one.json"}},
		{name: "P04-20.geojson", args: []string{"test", "-s", "0", "-j", "../../data/P04-20.geojson"}},
		{name: "P04-20.geojson", args: []string{"test", "-s", "1", "-j", "../../data/P04-20.geojson"}},
		{name: "P04-20.geojson", args: []string{"test", "-s", "2", "-j", "../../data/P04-20.geojson"}},
		{name: "P04-20.geojson", args: []string{"test", "-s", "3", "-j", "../../data/P04-20.geojson"}},
	}
	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			os.Args = tt.args
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			geojson.Output = ioutil.Discard
			main()
		})
	}
}
