package geojson_test

import (
	"reflect"
	"testing"

	. "github.com/takekazuomi/docker-mysql01/import/geojson"
)

var (
	fcOne *FeatureCollection = &FeatureCollection{
		Type: "FeatureCollection",
		Name: "P04-20",
		Crs: Crs{
			Type: "name",
			Properties: CrsProperties{
				Name: "urn:ogc:def:crs:EPSG::6668",
			},
		},
		Features: Features{
			Feature{
				Type: "Feature",
				Properties: Properties{
					P04001: 1,
					P04002: "医療法人樹恵会石田病院",
					P04003: "標津郡中標津町りんどう町5番地6",
					P04004: "内\u3000リハ",
					P04005: "",
					P04006: "",
					P04007: 4,
					P04008: 60,
					P04009: 9,
					P04010: 9}, Geometry: Geometry{
					Type:        "Point",
					Coordinates: []float64{144.9270714, 43.54427951},
				},
			},
		},
	}
)

func TestNewFeatures(t *testing.T) {
	type args struct {
		jsonFile string
	}
	tests := []struct {
		name    string
		args    args
		want    *FeatureCollection
		wantErr bool
	}{
		{"file not exists", args{"testdata/nofile.json"}, nil, true},
		{"broken.json", args{"testdata/broken.jsonb"}, nil, true},
		{"one.json", args{"testdata/one.json"}, fcOne, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFeatures(tt.args.jsonFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFeatures() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				t.Log(err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFeatures() = %v, want %v", got, tt.want)
			}
		})
	}
}
