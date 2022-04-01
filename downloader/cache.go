package downloader

import (
	"image"
	"log"
	"math"
	"sync"
)

type Cache struct {
	mux        sync.RWMutex
	pointCache map[Location]*Point
}

type Point struct {
	mux      sync.RWMutex
	distance float64
	images   map[int]image.Image
}

func (c *Cache) add(request DownloadRequest, img image.Image) {
	log.Printf("adding to cache %s, %d\n", request.Location, request.Angle)
	c.mux.RLock()
	if point, ok := c.pointCache[request.Location]; ok {
		defer c.mux.RUnlock()
		point.mux.Lock()
		defer point.mux.Unlock()
		point.images[request.Angle] = img
	} else {
		c.mux.RUnlock()
		p := Point{mux: sync.RWMutex{}, distance: math.MaxFloat64, images: map[int]image.Image{}}
		p.images[request.Angle] = img
		c.mux.Lock()
		defer c.mux.Unlock()
		c.pointCache[request.Location] = &p
	}
}

func (c *Cache) has(request DownloadRequest) bool {
	c.mux.RLock()
	defer c.mux.RUnlock()
	if point, ok := c.pointCache[request.Location]; ok {
		point.mux.RLock()
		defer point.mux.RUnlock()
		_, ok = point.images[request.Angle]
		return ok
	}
	return false
}

func (p *Point) update(l1, l2 Location, angle int, minDistance float64) (float64, image.Image, bool, bool) {
	p.mux.Lock()
	defer p.mux.Unlock()
	distance := l1.distance(l2)
	next := distance < minDistance
	remove := distance > p.distance
	p.distance = distance
	return distance, p.images[angle], remove, next
}

func (c *Cache) removeInLoop(location Location) {
	c.mux.RUnlock()
	c.mux.Lock()
	defer c.mux.RLock()
	defer c.mux.Unlock()
	delete(c.pointCache, location)
}

func (c *Cache) getAndClean(request DownloadRequest) image.Image {
	c.mux.RLock()
	defer c.mux.RUnlock()
	minDistance := math.MaxFloat64
	var next image.Image
	var toRemove []Location
	var nextLoc Location
	for l, point := range c.pointCache {
		newDistance, img, remove, newNext := point.update(l, request.Location, request.Angle, minDistance)
		if newNext {
			minDistance = newDistance
			next = img
			nextLoc = l
		}
		if remove {
			toRemove = append(toRemove, l)
		}
	}
	for _, location := range toRemove {
		if location != nextLoc {
			c.removeInLoop(location)
		}
	}
	return next
}

func (c *Cache) exporter(input <-chan DownloadRequest, output chan<- image.Image) {
	for {
		downloadRequest, ok := <-input
		if !ok {
			close(output)
			return
		}
		output <- c.getAndClean(downloadRequest)
	}
}
