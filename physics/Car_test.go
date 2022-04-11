package physics

import (
	"github.com/Daniel-W-Innes/car-environment-simulator/downloader"
	"reflect"
	"testing"
)

func TestCar_GetPosition(t *testing.T) {
	type fields struct {
		easting    float64
		northing   float64
		zoneNumber int
		angle      int
		zoneLetter string
	}
	tests := []struct {
		name   string
		fields fields
		want   downloader.DownloadRequest
	}{
		{"start_pos", fields{
			easting:    443714.419729,
			northing:   5019240.130888,
			zoneNumber: 18,
			zoneLetter: "N",
			angle:      160,
		}, downloader.DownloadRequest{
			Location: downloader.Location{
				Latitude:  45.32441,
				Longitude: -75.71821},
			Angle: 160,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Car{
				easting:    tt.fields.easting,
				northing:   tt.fields.northing,
				zoneNumber: tt.fields.zoneNumber,
				angle:      tt.fields.angle,
				zoneLetter: tt.fields.zoneLetter,
			}
			if got := c.GetPosition(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPosition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCar_Run(t *testing.T) {
	type fields struct {
		Input chan Command
	}
	type args struct {
		lat    float64
		lng    float64
		north  bool
		angle  int
		output chan downloader.DownloadRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		input   []Command
		wantErr bool
	}{
		{"no_input", fields{Input: make(chan Command)}, args{45.32441, -75.71821, true, 160, make(chan downloader.DownloadRequest)}, []Command{}, false},
		{"forward", fields{Input: make(chan Command)}, args{45.32441, -75.71821, true, 160, make(chan downloader.DownloadRequest)}, []Command{Forward}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Car{
				Input: tt.fields.Input,
			}
			if err := c.Run(tt.args.lat, tt.args.lng, tt.args.north, tt.args.angle, tt.args.output); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
			for _, command := range tt.input {
				tt.fields.Input <- command
			}
		})
	}
}

func TestCar_update(t *testing.T) {
	type fields struct {
		easting     float64
		northing    float64
		v           float64
		a           float64
		zoneNumber  int
		angle       int
		zoneLetter  string
		lastUpdated int64
		j           int8
	}
	type args struct {
		nextTime int64
		output   chan<- downloader.DownloadRequest
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"v_10", fields{j: 0, a: 0, v: 10, easting: 443706.25, northing: 5019264.20, angle: 160, zoneLetter: "T", zoneNumber: 18, lastUpdated: 0}, args{nextTime: 2000000000, output: make(chan downloader.DownloadRequest)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Car{
				easting:     tt.fields.easting,
				northing:    tt.fields.northing,
				v:           tt.fields.v,
				a:           tt.fields.a,
				zoneNumber:  tt.fields.zoneNumber,
				angle:       tt.fields.angle,
				zoneLetter:  tt.fields.zoneLetter,
				lastUpdated: tt.fields.lastUpdated,
				j:           tt.fields.j,
			}
			c.update(tt.args.nextTime, tt.args.output)
		})
	}
}
