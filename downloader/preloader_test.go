package downloader

import (
	"os"
	"testing"
)

func BenchmarkPreload(b *testing.B) {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		b.FailNow()
	}
	input := make(chan DownloadRequest)
	output := make(chan DownloadRequest)
	go preload(input, output, apiKey)
	for i := 0; i < b.N; i++ {
		input <- DownloadRequest{Location: Location{45.3219512062345, -75.71679090749016}, Angle: 160}
		for j := 0; j < 10; j++ {
			<-output
		}
	}
	close(input)
}
