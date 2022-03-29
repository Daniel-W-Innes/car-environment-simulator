package physics

import (
	"math"
	"testing"
)

func Test_xComponent(t *testing.T) {
	type args struct {
		d     float64
		angle int
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"30-60-90 a=30,d=2", args{2, 30}, 1},
		{"30-60-90 a=150,d=2", args{2, 150}, 1},
		{"30-60-90 a=210,d=2", args{2, 210}, -1},
		{"30-60-90 a=330,d=2", args{2, 330}, -1},

		{"45-45-90 a=45,d=2", args{math.Sqrt(2), 45}, 1},
		{"45-45-90 a=135,d=2", args{math.Sqrt(2), 135}, 1},
		{"45-45-90 a=225,d=2", args{math.Sqrt(2), 225}, -1},
		{"45-45-90 a=315,d=2", args{math.Sqrt(2), 315}, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := xComponent(tt.args.d, tt.args.angle); math.Abs(got-tt.want) > 0.001 {
				t.Errorf("xComponent() = %v, want %v, diff %f", got, tt.want, math.Abs(got-tt.want))
			}
		})
	}
}

func Test_yComponent(t *testing.T) {
	type args struct {
		d     float64
		angle int
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"30-60-90 a=60,d=2", args{2, 60}, 1},
		{"30-60-90 a=120,d=2", args{2, 120}, -1},
		{"30-60-90 a=240,d=2", args{2, 240}, -1},
		{"30-60-90 a=300,d=2", args{2, 300}, 1},

		{"45-45-90 a=45,d=2", args{math.Sqrt(2), 45}, 1},
		{"45-45-90 a=135,d=2", args{math.Sqrt(2), 135}, -1},
		{"45-45-90 a=225,d=2", args{math.Sqrt(2), 225}, -1},
		{"45-45-90 a=315,d=2", args{math.Sqrt(2), 315}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := yComponent(tt.args.d, tt.args.angle); math.Abs(got-tt.want) > 0.001 {
				t.Errorf("yComponent() = %v, want %v, diff %f", got, tt.want, math.Abs(got-tt.want))
			}
		})
	}
}
