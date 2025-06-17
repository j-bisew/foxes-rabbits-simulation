package entities

import (
	"github.com/j-bisew/foxes-rabbits-simulation/geom"
)
type Rabbit struct {
	Animal
} 

func NewRabbit(x, y float64) *Rabbit {
	return &Rabbit{
		Animal: Animal{
			Pos: geom.Point{X: x, Y: y},
			Energy: 125.0,
			MaxEnergy: 200.0,
			EnergyLoss: 1.5,
			CriticalHungerThreshold: 100.0,
			ReproduceCooldown: 0,
			SearchRadius: 40.0,
			Species: "rabbit",
			FoodType: "grass",
			MovementSpeed: 3.5,
			Alive: true,
		},
	}
}