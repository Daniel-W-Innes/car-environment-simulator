package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"github.com/Daniel-W-Innes/car-environment-simulator/downloader"
	"github.com/Daniel-W-Innes/car-environment-simulator/physics"
	"log"
	"os"
	"time"
)

func main() {
	cse := app.New()
	w := cse.NewWindow("car environment simulator")

	car := physics.Car{Input: make(chan physics.Command)}

	backend := downloader.New()
	backend.Run(os.Getenv("API_KEY"), false)

	img := canvas.NewImageFromFile("res/45.3219512062345,-75.71679090749016,70.jpg")

	ticker := time.NewTicker(17 * time.Millisecond)
	go func() {
		for range ticker.C {
			backend.LocationUpdater <- car.GetPosition()
			newImg, ok := <-backend.Output
			if !ok {
				return
			}
			img = canvas.NewImageFromImage(newImg)
			w.SetContent(img)
		}
	}()

	w.SetContent(img)
	w.Canvas().SetOnTypedKey(func(event *fyne.KeyEvent) {
		log.Println(event.Name)
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
		case fyne.KeyP:
			log.Println(car.String())
		}
	})

	err := car.Run(45.32441, -75.71821, true, 160, backend.Input)
	if err != nil {
		log.Fatalln(err)
	}

	w.SetOnClosed(func() {
		car.Input <- physics.Exit
		ticker.Stop()
		close(backend.LocationUpdater)
	})

	w.Resize(fyne.Size{Width: 1020, Height: 1020})
	w.ShowAndRun()
}
