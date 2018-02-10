package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type block struct {
	sprite *sdl.Texture
	w      int32
	h      int32
	x      int32
	y      int32
}

func newBlock(r *sdl.Renderer, x, y, w, h int32) (*block, error) {
	t, err := img.LoadTexture(r, "res/block.png")
	if err != nil {
		return nil, fmt.Errorf("unable to load block sprite: %v", err)
	}

	return &block{
		sprite: t,
		w:      w,
		h:      h,
		x:      x,
		y:      y,
	}, nil
}

func (b *block) paint(r *sdl.Renderer) error {
	rect := &sdl.Rect{
		X: b.x,
		Y: b.y,
		W: b.w,
		H: b.h,
	}

	if err := r.Copy(b.sprite, nil, rect); err != nil {
		return fmt.Errorf("failed to paint block to surface: %v", err)
	}

	return nil
}
