package downloader

import (
	"image"
)

type Downloader struct {
	Input            chan DownloadRequest
	LocationUpdater  chan DownloadRequest
	Output           chan image.Image
	downloadRequests chan DownloadRequest
	cache            *Cache
}

type DownloadRequest struct {
	Location Location
	Angle    int
}

func New() Downloader {
	return Downloader{
		Input:            make(chan DownloadRequest),
		LocationUpdater:  make(chan DownloadRequest),
		Output:           make(chan image.Image),
		downloadRequests: make(chan DownloadRequest),
		cache:            &Cache{},
	}
}

func exporter(cache *Cache, input <-chan DownloadRequest, output chan<- image.Image) {
	for {
		downloadRequest, ok := <-input
		if !ok {
			close(output)
			return
		}
		output <- cache.getAndClean(downloadRequest)
	}
}

func (d *Downloader) Run(key string) {
	go preload(d.Input, d.downloadRequests, key)
	go download(d.downloadRequests, d.cache, key)
	go exporter(d.cache, d.LocationUpdater, d.Output)
}
