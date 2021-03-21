package simple

import (
	"log"
	"os"
)

type App struct {
	currentScene Scene
	widgets      Widget
	sr           *SimpleRunner
	scene        chan Scene
	quit         chan bool
}

func NewApp(scene Scene) *App {
	sr := NewSimpleRunner(false)
	a := &App{
		currentScene: scene,
		sr:           sr,
		scene:        make(chan Scene),
		quit:         make(chan bool, 1),
	}
	sr.Start()
	return a
}

func (a *App) run() {
	a.runScene(a.currentScene)
	for {
		log.Println("selecting")
		select {
		case scene := <-a.scene:
			log.Printf("scene received: %T\n", scene)
			a.currentScene = scene
			a.runScene(a.currentScene)
		case event := <-a.sr.Events:
			log.Println("event received:", event.Type().String())
			a.handleEvent(event)
		case <-a.quit:
			log.Println("quit received, quitting")
			a.sr.Stop()
			os.Exit(0) // maybe return instead
		case <-a.sr.Exited:
			// user input causes the simple process to exit
			// render updates afterwards
			log.Println("exited received, rerendering")
			a.runScene(a.currentScene)
		}
	}
}

type BoundEventHandler func(app *App) error

func (a *App) handleEvent(event Event) {
	var eventHandlers []BoundEventHandler

	handlers, err := a.widgets.Update(event)
	if err != nil {
		log.Printf("while handling %s event: %s", event.Type().String(), err.Error())
	}
	eventHandlers = append(eventHandlers, handlers...)

	for _, evHandler := range eventHandlers {
		evHandler(a)
	}
}

func (a *App) SetScene(scene Scene) {
	a.currentScene = scene
	a.scene <- scene
}

func (a *App) RunForever() {
	a.run()
}

func (a *App) Quit() {
	log.Println("app: quit sent")
	a.quit <- true
}

func (a *App) RunInBackground() {
	go a.run()
}

func (a *App) ScreenWidth() int {
	return 1404
}

func (a *App) ScreenHeight() int {
	return 1872
}
