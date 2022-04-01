package downloader

import (
	"image"
	"sync"
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
		cache:            &Cache{mux: sync.RWMutex{}, pointCache: map[Location]*Point{}},
	}
}

func (d *Downloader) Run(key string, preventDownload bool) {
	go preload(d.Input, d.downloadRequests, key)
	go download(d.downloadRequests, d.cache, key, preventDownload)
	go d.cache.exporter(d.LocationUpdater, d.Output)
}
