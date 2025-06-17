package entities

import (
	"github.com/j-bisew/foxes-rabbits-simulation/geom"
)

type Fox struct {
	Animal
} 

func NewFox(x, y float64) *Fox {
	return &Fox{
		Animal: Animal{
			Pos: geom.Point{X: x, Y: y},
			Energy: 200.0,
			MaxEnergy: 300.0,
			EnergyLoss: 3.0,
			CriticalHungerThreshold: 170.0,
			ReproduceCooldown: 0,
			SearchRadius: 60.0,
			Species: "fox",
			FoodType: "rabbit",
			MovementSpeed: 2.0,
			Alive: true,
		},
	}
}