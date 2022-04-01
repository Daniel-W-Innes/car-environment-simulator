package downloader

import (
	"sync"
	"testing"
)

func BenchmarkDownloader(b *testing.B) {
	input := make(chan DownloadRequest)
	cache := &Cache{mux: sync.RWMutex{}, pointCache: map[Location]*Point{}}
	go download(input, cache, "", true)
	for i := 0; i < b.N; i++ {
		input <- DownloadRequest{Location: Location{45.3219512062345, -75.71679090749016}, Angle: 160}
	}
}
