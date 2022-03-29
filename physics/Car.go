package physics

import (
	"fmt"
	"github.com/Daniel-W-Innes/car-environment-simulator/downloader"
	"github.com/im7mortal/UTM"
	"math"
	"time"
)

type Command uint8

const (
	Forward Command = iota
	Backward
	Left
	Right
	CruiseControl
	Stop
	Exit
)

const (
	maxJ         = 1
	maxA         = 12
	timeDilation = 2
)

type Car struct {
	Input                   chan Command
	easting, northing, v, a float64
	zoneNumber, angle       int
	zoneLetter              string
	ticker                  *time.Ticker
	lastUpdated             int64
	j                       int8
}

func (c Car) ToString() string {
	return fmt.Sprintf("x %f, y %f", c.easting, c.northing)
}

func (c Car) GetPosition() downloader.DownloadRequest {
	latitude, longitude, err := UTM.ToLatLon(c.easting, c.northing, c.zoneNumber, c.zoneLetter)
	if err != nil {
		return downloader.DownloadRequest{}
	}
	return downloader.DownloadRequest{Location: downloader.Location{Latitude: latitude, Longitude: longitude}, Angle: c.angle}
}

func (c *Car) Run(lat, lng float64, north bool, output chan<- downloader.DownloadRequest) error {
	easting, northing, zoneNumber, zoneLetter, err := UTM.FromLatLon(lat, lng, north)
	if err != nil {
		return err
	}
	c.easting = easting
	c.northing = northing
	c.zoneNumber = zoneNumber
	c.zoneLetter = zoneLetter

	go func(car *Car) {
		for {
			switch <-c.Input {
			case Exit:
				close(output)
			case Forward:
				c.a += maxJ
				if c.a > maxA {
					c.a = maxA
				}
			case Backward:
				c.a -= maxJ
				if c.a < -maxA {
					c.a = -maxA
				}
			case Left:
				if c.angle == 0 {
					c.angle = 359
				} else {
					c.angle -= 1
				}
			case Right:
				if c.angle == 359 {
					c.angle = 0
				} else {
					c.angle += 1
				}
			case CruiseControl:
				c.j = 0
				c.a = 0
				c.v = 10
			case Stop:
				c.j = 0
				c.a = 0
				c.v = 0
			}
		}
	}(c)

	c.ticker = time.NewTicker(1 * time.Millisecond)

	go func(c *Car) {
		for range c.ticker.C {
			next := time.Now().UnixNano()
			dt := (float64(next-c.lastUpdated) * math.Pow(10, -9)) / timeDilation
			if c.a > 0 {
				c.j = -maxJ
			} else {
				c.j = maxJ
			}
			d := deltaD(dt, c.v, c.a, float64(c.j))
			c.easting += xComponent(d, c.angle)
			c.northing += yComponent(d, c.angle)
			c.v += deltaV(dt, c.a, float64(c.j))
			c.a += deltaA(dt, c.a, float64(c.j))
			c.lastUpdated = next
			for i := 0.5; i < 2; i += 0.5 {
				d = deltaD(i, c.v, c.a, float64(c.j))
				latitude, longitude, err := UTM.ToLatLon(c.easting+xComponent(d, c.angle), c.northing+yComponent(d, c.angle), c.zoneNumber, c.zoneLetter)
				if err != nil {
					return
				}
				output <- downloader.DownloadRequest{Location: downloader.Location{Latitude: latitude, Longitude: longitude}, Angle: c.angle}
			}
		}
	}(c)

	return nil
}
