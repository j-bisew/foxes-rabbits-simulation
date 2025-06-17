package interfaces

import "github.com/j-bisew/foxes-rabbits-simulation/geom"

type Entity interface {
	Update(WorldInterface)
	GetPosition() geom.Point
	GetSpecies() string
	GetFoodType() string
	IsAlive() bool
	GetEnergy() float64
	UpdateEnergy(float64)
	Kill()
}