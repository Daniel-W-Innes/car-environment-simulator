package downloader

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const SIZE = "1280x960"

func getImageFromGoogle(request DownloadRequest, key string, prevent bool) (image.Image, error) {
	if os.Getenv("USE_GOOGLE") != "y" || prevent {
		return nil, errors.New("tried to download image from google")
	}
	log.Printf("getting image from google %s, %d\n", request.Location.String(), request.Angle)
	path := fmt.Sprintf("https://maps.googleapis.com/maps/api/streetview?size=%s&location=%s&heading=%d&key=%s", SIZE, request.Location.String(), request.Angle, key)
	response, err := http.Get(path)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error bad status from googleapis %d", response.StatusCode)
	}
	defer func() {
		err := response.Body.Close()
		if err != nil {
			log.Fatalf(err.Error())
		}
	}()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	decode, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		log.Fatalln(err)
	}
	return decode, nil
}

func getImage(request DownloadRequest, key string, preventDownload bool) (image.Image, error) {
	path := fmt.Sprintf("/home/daniel/.cache/car-environment-simulator/%s,%d.jpg", request.Location.String(), request.Angle)
	f, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		img, err := getImageFromGoogle(request, key, preventDownload)
		if err != nil {
			return img, err
		}
		out, err := os.Create(path)
		if err != nil {
			return img, err
		}
		opt := jpeg.Options{Quality: 100}
		err = jpeg.Encode(out, img, &opt)
		return img, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	return img, err
}

func download(input <-chan DownloadRequest, cache *Cache, key string, preventDownload bool) {
	for {
		downloadRequest, ok := <-input
		if !ok {
			return
		}
		if !cache.has(downloadRequest) {
			img, err := getImage(downloadRequest, key, preventDownload)
			if err != nil {
				log.Fatalln(err)
				return
			}
			cache.add(downloadRequest, img)
		}
	}
}
