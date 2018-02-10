package main

import (
	"fmt"
	"log"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type scene struct {
	background *sdl.Texture
	hero       *hero
	block      *block
}

func newScene(r *sdl.Renderer) (*scene, error) {
	t, err := img.LoadTexture(r, "res/background.png")
	if err != nil {
		return nil, fmt.Errorf("unable to load texture: %v", err)
	}

	h, err := newHero(r)
	if err != nil {
		return nil, err
	}

	b, err := newBlock(r, 400, ground-50, 50, 50)
	if err != nil {
		return nil, err
	}

	return &scene{background: t, hero: h, block: b}, nil
}

func (s *scene) run(r *sdl.Renderer, events chan sdl.Event) chan error {
	errc := make(chan error)

	go func() {
		updateTick := time.Tick(10 * time.Millisecond)
		renderTick := time.Tick(32 * time.Millisecond)

		lastUpdate := time.Now()

		for {
			select {
			case e := <-events:
				if quit := s.handleEvent(e); quit {
					return
				}
			case <-updateTick:
				now := time.Now()
				duration := now.Sub(lastUpdate)
				factor := int32(duration.Nanoseconds()/int64(1000000)) / 3
				s.update(factor)
				lastUpdate = now
			case <-renderTick:
				if err := s.paint(r); err != nil {
					errc <- err
				}
			}

		}
	}()

	return errc
}

func (s *scene) handleEvent(event sdl.Event) bool {
	switch event.(type) {
	case *sdl.QuitEvent:
		return true
	case *sdl.KeyboardEvent:
		ke := event.(*sdl.KeyboardEvent)
		//log.Printf("keyboard event, key:%v, up/down: %v repeat: %v", ke.Keysym.Scancode, ke.State, ke.Repeat)
		switch ke.Keysym.Scancode {
		case 79:
			if ke.State == 1 {
				s.hero.goRight()
			} else {
				s.hero.stop()
			}
		case 80:
			if ke.State == 1 {
				s.hero.goLeft()
			} else {
				s.hero.stop()
			}
		case 82:
			if ke.State == 1 {
				s.hero.jump()
			}
		case 20: // Q
			if ke.State == 1 {
				gravity -= 0.01
			}
		case 4: // A
			if ke.State == 1 {
				gravity += 0.01
			}
		case 26: // W
			if ke.State == 1 {
				jumpConstant += 0.1
			}
		case 22: // S
			if ke.State == 1 {
				jumpConstant -= 0.1
			}
		case 8: // E
			if ke.State == 1 {
				jumpIncrease += 0.1
			}
		case 7: // D
			if ke.State == 1 {
				jumpIncrease -= 0.1
			}
		case 21: // R
			if ke.State == 1 {
				initialJumpVelocity += 0.1
			}
		case 9: // F
			if ke.State == 1 {
				initialJumpVelocity -= 0.1
			}
		case 41: // ESC
			if ke.State == 1 {
				gravity = -0.11
				jumpConstant = 1.2
				jumpIncrease = 0.25
				initialJumpVelocity = float64(1)
				s.hero.x = 100
				s.hero.y = 290
			}
		default:
			log.Printf("unknown key %v", ke.Keysym.Scancode)
		}
	case *sdl.MouseMotionEvent, *sdl.WindowEvent, *sdl.TouchFingerEvent, *sdl.CommonEvent:
	default:
		log.Printf("unknown event %T", event)
	}
	return false
}

func (s *scene) update(f int32) error {
	return s.hero.update(f, s.block)
}

func (s *scene) paint(r *sdl.Renderer) error {
	r.Clear()

	if err := r.Copy(s.background, nil, nil); err != nil {
		return fmt.Errorf("unable to paint background: %v", err)
	}

	s.hero.paint(r)

	s.block.paint(r)

	if err := drawText(r, 20, 20, fmt.Sprintf("Gravity: %f", gravity)); err != nil {
		return err
	}
	if err := drawText(r, 20, 40, fmt.Sprintf("Jump inc: %f", jumpIncrease)); err != nil {
		return nil
	}
	if err := drawText(r, 20, 60, fmt.Sprintf("Jump const: %f", jumpConstant)); err != nil {
		return nil
	}
	if err := drawText(r, 20, 80, fmt.Sprintf("Jump vel: %f", initialJumpVelocity)); err != nil {
		return nil
	}

	r.Present()
	return nil
}

func (s *scene) destroy() {
	s.hero.destroy()
	s.background.Destroy()
}

func drawText(r *sdl.Renderer, x, y int32, text string) error {
	f, err := ttf.OpenFont("font/PressStart2P.ttf", 10)
	if err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	defer f.Close()

	c := sdl.Color{R: 255, G: 255, B: 255, A: 255}
	s, err := f.RenderUTF8Solid(text, c)
	if err != nil {
		return fmt.Errorf("could not render title: %v", err)
	}
	defer s.Free()

	t, err := r.CreateTextureFromSurface(s)
	if err != nil {
		return fmt.Errorf("could not create texture: %v", err)
	}
	defer t.Destroy()

	dest := &sdl.Rect{X: x, Y: y, W: 200, H: 20}

	if err := r.Copy(t, nil, dest); err != nil {
		return fmt.Errorf("could not copy texture: %v", err)
	}

	return nil
}
