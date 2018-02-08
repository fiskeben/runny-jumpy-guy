package main

type mover interface {
	Update() error
}

type collider interface {
	CheckCollision(c *collider) error
}
