package interfaces

import "github.com/j-bisew/foxes-rabbits-simulation/geom"

type WorldInterface interface {
	FindNearbyEntities(pos geom.Point, radius float64, species string) []Entity
	CreateOffspring(parent1, parent2 Entity) Entity
	IsValidPosition(x, y float64) bool
	ConsumeFood(entity Entity, eater Entity) float64
}