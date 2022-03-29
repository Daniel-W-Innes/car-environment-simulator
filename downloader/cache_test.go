package downloader

import (
	"image"
	"math"
	"reflect"
	"testing"
)

func comparePoint(point, point2 *Point) bool {
	return math.Abs(point.distance-point2.distance) < 0.001 && reflect.DeepEqual(point.images, point2.images)
}

func TestCache_add(t *testing.T) {
	type fields struct {
		pointCache map[Location]*Point
	}
	type args struct {
		request DownloadRequest
		img     image.Image
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		result fields
	}{
		{"add_to_empty", fields{pointCache: map[Location]*Point{}}, args{DownloadRequest{Location{1, 1}, 1}, image.Black},
			fields{pointCache: map[Location]*Point{
				Location{1, 1}: {images: map[int]image.Image{1: image.Black}, distance: math.MaxFloat64}}}},
		{"add_to_existing", fields{pointCache: map[Location]*Point{Location{1, 1}: {images: map[int]image.Image{1: image.Black}, distance: math.MaxFloat64}}}, args{DownloadRequest{Location{2, 2}, 2}, image.White},
			fields{pointCache: map[Location]*Point{
				Location{1, 1}: {images: map[int]image.Image{1: image.Black}, distance: math.MaxFloat64},
				Location{2, 2}: {images: map[int]image.Image{2: image.White}, distance: math.MaxFloat64}}}},
		{"add_to_not_empty", fields{pointCache: map[Location]*Point{Location{1, 1}: {images: map[int]image.Image{1: image.Black}, distance: math.MaxFloat64}}}, args{DownloadRequest{Location{1, 1}, 2}, image.White},
			fields{pointCache: map[Location]*Point{
				Location{1, 1}: {images: map[int]image.Image{1: image.Black, 2: image.White}, distance: math.MaxFloat64}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cache{
				pointCache: tt.fields.pointCache,
			}
			c.add(tt.args.request, tt.args.img)
			for location, point := range c.pointCache {
				if point2, ok := tt.result.pointCache[location]; ok {
					if comparePoint(point, point2) {
						continue
					}
				}
				t.FailNow()
			}
		})
	}
}

func TestCache_getAndClean(t *testing.T) {
	type fields struct {
		pointCache map[Location]*Point
	}
	type args struct {
		request DownloadRequest
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   image.Image
		result fields
	}{
		{"get_empty", fields{pointCache: map[Location]*Point{}}, args{DownloadRequest{Location{1, 1}, 1}}, nil,
			fields{pointCache: map[Location]*Point{}}},
		{"get_same", fields{pointCache: map[Location]*Point{Location{1, 1}: {images: map[int]image.Image{1: image.Black}, distance: math.MaxFloat64}}}, args{DownloadRequest{Location{1, 1}, 1}}, image.Black,
			fields{pointCache: map[Location]*Point{Location{1, 1}: {images: map[int]image.Image{1: image.Black}, distance: 0}}}},
		{"get_near", fields{pointCache: map[Location]*Point{Location{1, 1}: {images: map[int]image.Image{1: image.Black}, distance: math.MaxFloat64}}}, args{DownloadRequest{Location{2, 2}, 1}}, image.Black,
			fields{pointCache: map[Location]*Point{Location{1, 1}: {images: map[int]image.Image{1: image.Black}, distance: 157225.43203}}}},
		{"get_angle", fields{pointCache: map[Location]*Point{Location{1, 1}: {images: map[int]image.Image{1: image.Black, 2: image.White}, distance: math.MaxFloat64}}}, args{DownloadRequest{Location{1, 1}, 2}}, image.White,
			fields{pointCache: map[Location]*Point{Location{1, 1}: {images: map[int]image.Image{1: image.Black, 2: image.White}, distance: 0}}}},
		{"get_clean", fields{pointCache: map[Location]*Point{
			Location{1, 1}: {images: map[int]image.Image{1: image.Black}, distance: 0},
			Location{2, 2}: {images: map[int]image.Image{1: image.White}, distance: math.MaxFloat64}}},
			args{DownloadRequest{Location{2, 2}, 1}}, image.White,
			fields{pointCache: map[Location]*Point{Location{2, 2}: {images: map[int]image.Image{1: image.White}, distance: 0}}}},
		{"get_clean_top", fields{pointCache: map[Location]*Point{Location{1, 1}: {images: map[int]image.Image{1: image.White}, distance: 0}}}, args{DownloadRequest{Location{2, 2}, 1}}, image.White,
			fields{pointCache: map[Location]*Point{Location{1, 1}: {images: map[int]image.Image{1: image.White}, distance: 157225.43203}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cache{
				pointCache: tt.fields.pointCache,
			}
			if got := c.getAndClean(tt.args.request); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getAndClean() = %v, want %v", got, tt.want)
			}
			for location, point := range c.pointCache {
				if point2, ok := tt.result.pointCache[location]; ok {
					if comparePoint(point, point2) {
						continue
					}
				}
				t.FailNow()
			}
		})
	}
}

func TestCache_has(t *testing.T) {
	type fields struct {
		pointCache map[Location]*Point
	}
	type args struct {
		request DownloadRequest
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"has_empty", fields{pointCache: map[Location]*Point{}}, args{DownloadRequest{Location{1, 1}, 1}}, false},
		{"has_missing", fields{pointCache: map[Location]*Point{Location{1, 1}: {images: map[int]image.Image{1: image.Black}, distance: math.MaxFloat64}}}, args{DownloadRequest{Location{2, 2}, 1}}, false},
		{"has_missing_angle", fields{pointCache: map[Location]*Point{Location{1, 1}: {images: map[int]image.Image{2: image.Black}, distance: math.MaxFloat64}}}, args{DownloadRequest{Location{1, 1}, 1}}, false},
		{"has", fields{pointCache: map[Location]*Point{Location{1, 1}: {images: map[int]image.Image{1: image.Black}, distance: math.MaxFloat64}}}, args{DownloadRequest{Location{1, 1}, 1}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cache{
				pointCache: tt.fields.pointCache,
			}
			if got := c.has(tt.args.request); got != tt.want {
				t.Errorf("has() = %v, want %v", got, tt.want)
			}
		})
	}
}
