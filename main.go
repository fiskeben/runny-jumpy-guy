package main

import (
	"fmt"
	"runtime"

	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("Error running the game: %v\n", err)
	}
}

func run() error {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return err
	}
	defer sdl.Quit()

	if err := ttf.Init(); err != nil {
		return fmt.Errorf("could not initialize TTF: %v", err)
	}
	defer ttf.Quit()

	w, r, err := sdl.CreateWindowAndRenderer(800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		return err
	}
	defer w.Destroy()

	w.SetTitle("Runny Jumpy Guy")

	if err = mix.OpenAudio(mix.DEFAULT_FREQUENCY, mix.INIT_MP3, mix.DEFAULT_CHANNELS, mix.DEFAULT_CHUNKSIZE); err != nil {
		return fmt.Errorf("failed to initialize audio: %v", err)
	}

	s, err := newScene(r)
	if err != nil {
		return err
	}
	defer s.destroy()

	events := make(chan sdl.Event)
	errc := s.run(r, events)
	runtime.LockOSThread()

	for {
		select {
		case events <- sdl.WaitEvent():
		case err := <-errc:
			return err
		}
	}
}
