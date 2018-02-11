package main

import (
	"fmt"
	"sync"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

var gravity = -0.11
var jumpConstant = 1.2
var jumpIncrease = 0.25
var initialJumpVelocity = float64(1)
var ground = int32(360)

type hero struct {
	mu           sync.RWMutex
	sprites      *sprites
	w            int32
	h            int32
	x            int32
	y            int32
	speed        float32
	time         int
	direction    int
	jumpVelocity float64
	onGround     bool
	stepping     int
	startJump    bool
	stepSound    *mix.Music
	jumpSound    *mix.Music
}

type sprites struct {
	idle     *sdl.Texture
	running  []*sdl.Texture
	jumping  *sdl.Texture
	landing  *sdl.Texture
	inAir    *sdl.Texture
	grabbing *sdl.Texture
}

func newHero(r *sdl.Renderer) (*hero, error) {
	s := sprites{}

	var err error

	s.idle, err = img.LoadTexture(r, "res/idle.gif")
	if err != nil {
		return nil, fmt.Errorf("failed to load texture: %v", err)
	}

	running := make([]*sdl.Texture, 0)
	for i := 1; i < 9; i++ {
		t, err := img.LoadTexture(r, fmt.Sprintf("res/run-%d.png", i))
		if err != nil {
			return nil, fmt.Errorf("failed to load texture: %v", err)
		}
		running = append(running, t)
	}
	s.running = running

	s.jumping, err = img.LoadTexture(r, "res/jump.png")
	if err != nil {
		return nil, fmt.Errorf("failed to load texture: %v", err)
	}

	s.landing, err = img.LoadTexture(r, "res/landing.png")
	if err != nil {
		return nil, fmt.Errorf("failed to load texture: %v", err)
	}

	s.inAir, err = img.LoadTexture(r, "res/mid-air.gif")
	if err != nil {
		return nil, fmt.Errorf("failed to load texture: %v", err)
	}

	s.grabbing, err = img.LoadTexture(r, "res/ledge-grab.gif")
	if err != nil {
		return nil, fmt.Errorf("failed to load texture: %v", err)
	}

	step, err := mix.LoadMUS("res/footstep.mp3")
	if err != nil {
		return nil, fmt.Errorf("failed to load foostep sound: %v", err)
	}

	jump, err := mix.LoadMUS("res/jump.mp3")
	if err != nil {
		return nil, fmt.Errorf("failed to load jump sound: %v", err)
	}

	return &hero{
		sprites:      &s,
		w:            38,
		h:            68,
		x:            100,
		y:            ground + 68,
		speed:        0,
		time:         0,
		onGround:     true,
		jumpVelocity: 0.0,
		stepping:     0,
		startJump:    false,
		stepSound:    step,
		jumpSound:    jump,
	}, nil
}

func (h *hero) update(factor int32, b *block) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.x += int32(h.speed) * factor
	if h.x+h.w > 800 {
		h.x = 800 - h.w
	}
	if h.x < 0 {
		h.x = 0
	}

	if h.speed > 0 && (h.x+h.w) > b.x && (h.x+h.w) < b.x+b.w && (h.y+h.h) > b.y {
		h.x = b.x - h.w
	}

	if h.speed < 0 && h.x < (b.x+b.w) && h.x > b.x && h.y+h.h > b.y {
		h.x = b.x + b.w
	}

	if h.jumpVelocity > 0 {
		y := gravity*h.jumpVelocity*h.jumpVelocity + jumpConstant*h.jumpVelocity
		h.jumpVelocity += jumpIncrease
		h.y -= int32(y)
		h.onGround = false
	}

	if (h.y+h.h) > b.y && (h.x+h.w) > b.x && h.x < (b.x+b.w) {
		h.y = b.y - h.h
		h.onGround = true
		h.jumpVelocity = 0
	}

	if h.y+h.h < ground && h.onGround && (h.x+h.w < b.x || h.x > b.x+b.w) {
		h.jumpVelocity = 12.5
		h.onGround = false
	}

	if h.y+h.h > ground {
		h.y = ground - h.h
		h.jumpVelocity = 0
		h.onGround = true
	}

	if h.speed > 0 {
		h.direction = 1
	} else if h.speed < 0 {
		h.direction = -1
	}
	if h.time == int(^uint(0)>>1) {
		h.time = -1
	}
	h.time++
	return nil
}

func (h *hero) paint(r *sdl.Renderer) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	rect := &sdl.Rect{
		X: h.x,
		Y: h.y,
		W: h.w,
		H: h.h,
	}

	var t *sdl.Texture
	flip := sdl.FLIP_NONE

	if h.speed == 0 {
		t = h.sprites.idle
		if h.direction == -1 {
			flip = sdl.FLIP_HORIZONTAL
		}
	} else if h.jumpVelocity > 0 {
		t = h.sprites.jumping
		if h.direction == -1 {
			flip = sdl.FLIP_HORIZONTAL
		}
	} else if h.jumpVelocity < 0 {
		t = h.sprites.landing
		if h.direction == -1 {
			flip = sdl.FLIP_HORIZONTAL
		}
	} else {
		i := h.time / 10 % len(h.sprites.running)
		t = h.sprites.running[i]
		if h.speed < 0 {
			flip = sdl.FLIP_HORIZONTAL
		}
		if i%4 == 0 {
			h.stepSound.Play(1)
			h.stepping = i
		}
	}

	if h.startJump {
		h.jumpSound.Play(1)
		h.startJump = false
	}

	if err := r.CopyEx(t, nil, rect, 0, nil, flip); err != nil {
		return fmt.Errorf("failed to paint sprite to surface: %v", err)
	}

	return nil
}

func (h *hero) destroy() {
	h.sprites.grabbing.Destroy()
	h.sprites.idle.Destroy()
	h.sprites.jumping.Destroy()
	h.sprites.landing.Destroy()
	h.sprites.inAir.Destroy()
	for _, t := range h.sprites.running {
		t.Destroy()
	}
}

func (h *hero) goRight() {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.speed = 1
}

func (h *hero) goLeft() {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.speed = -1
}

func (h *hero) stop() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.speed = 0
}

func (h *hero) jump() {
	h.mu.Lock()
	defer h.mu.Unlock()

	if !h.onGround || h.jumpVelocity > 0.0 {
		return
	}

	h.jumpVelocity = initialJumpVelocity
	h.startJump = true
}
