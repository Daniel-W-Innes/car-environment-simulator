package downloader

import (
	"os"
	"testing"
)

func BenchmarkInput(b *testing.B) {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		b.FailNow()
	}
	downloader := New()
	downloader.Run(apiKey, true)
	for i := 0; i < b.N; i++ {
		downloader.Input <- DownloadRequest{Location: Location{45.3219512062345, -75.71679090749016}, Angle: 160}
	}
}
