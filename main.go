package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"github.com/Daniel-W-Innes/car-environment-simulator/physics"
	"log"
)

func main() {
	cse := app.New()
	w := cse.NewWindow("car environment simulator")

	car := physics.Car{Input: make(chan physics.Command)}

	image := canvas.NewImageFromFile("cash/45.3219512062345,-75.71679090749016,70.jpg")

	w.SetContent(image)
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

	err := car.Run(45.3219512062345, -75.71679090749016, true)
	if err != nil {
		log.Fatalln(err)
	}

	w.SetOnClosed(func() {
		car.Input <- physics.Exit
	})

	w.ShowAndRun()
}
