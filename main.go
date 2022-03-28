package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"github.com/Daniel-W-Innes/car-environment-simulator/physics"
	"github.com/Daniel-W-Innes/street-view-image-manager"
	"image"
	"log"
	"os"
	"time"
)

func main() {
	cse := app.New()
	w := cse.NewWindow("car environment simulator")

	car := physics.Car{Input: make(chan physics.Command)}

	downloader := manager.Downloader{
		Input:           make(chan manager.DownloadRequest),
		LocationUpdater: make(chan manager.DownloadRequest),
		Output:          make(chan image.Image),
	}
	downloader.Run(os.Getenv("API_KEY"))

	img := canvas.NewImageFromFile("cash/45.3219512062345,-75.71679090749016,70.jpg")

	ticker := time.NewTicker(1 / 60 * time.Second)
	go func() {
		for range ticker.C {
			downloader.LocationUpdater <- car.GetPosition()
			newImg, ok := <-downloader.Output
			if !ok {
				return
			}
			img = canvas.NewImageFromImage(newImg)
		}
	}()

	w.SetContent(img)
	w.Canvas().SetOnTypedKey(func(event *fyne.KeyEvent) {
		switch event.Name {
		case fyne.KeyW:
			car.Input <- physics.Forward
		case fyne.KeyS:
			car.Input <- physics.Backward
		case fyne.KeyA:
			car.Input <- physics.Left
		case fyne.KeyD:
			car.Input <- physics.Right
		case fyne.KeySpace:
			car.Input <- physics.CruiseControl
		case fyne.KeyBackspace:
			car.Input <- physics.Stop
		}
		log.Println(car.ToString())
	})

	err := car.Run(45.3219512062345, -75.71679090749016, true, downloader.Input)
	if err != nil {
		log.Fatalln(err)
	}

	w.SetOnClosed(func() {
		car.Input <- physics.Exit
		ticker.Stop()
		close(downloader.LocationUpdater)
	})

	w.ShowAndRun()
}
