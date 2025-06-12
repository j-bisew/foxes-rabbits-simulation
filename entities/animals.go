package entities

import "kroliki/geom"

type Entity interface {
	Update(*World)
	GetPosition() figures.Point
	GetSpecies() string
	GetFoodType() string
	IsAlive() bool
}

type Animal struct {
	Pos figures.Point
	Energy float64
	ReproduceCooldown int
	Species string
	Food string
	Alive bool
} 

func (a *Animal) GetPosition() figures.Point { return a.Pos }
func (a *Animal) GetSpecies() string { return a.Species }
func (a *Animal) GetFoodType() string { return a.Food }
func (a *Animal) IsAlive() bool { return a.Alive }

func (a *Animal) Move(dx, dy float64) {
	a.Pos.X += dx
	a.Pos.Y += dy
}

type Rabbit struct {
	Animal
} 

func NewRabbit(x, y float64) *Rabbit {
	return &Rabbit{
		Animal: Animal{
			Pos: figures.Point{X: x, Y: y},
			Energy: 10,
			ReproduceCooldown: 0,
			Species: "rabbit",
			Food: "grass",
			Alive: true,
		},
	}
}

func (r *Rabbit) GetPosition() figures.Point { return r.Animal.GetPosition() }
func (r *Rabbit) GetSpecies() string { return r.Animal.GetSpecies() }
func (r *Rabbit) GetFoodType() string { return r.Animal.GetFoodType() }
func (r *Rabbit) IsAlive() bool { return r.Animal.IsAlive() }

type Fox struct {
	Animal
} 

func NewFox(x, y float64) *Fox {
	return &Fox{
		Animal: Animal{
			Pos: figures.Point{X: x, Y: y},
			Energy: 10,
			ReproduceCooldown: 0,
			Species: "fox",
			Food: "rabbit",
			Alive: true,
		},
	}
}

func (f *Fox) GetPosition() figures.Point { return f.Animal.GetPosition() }
func (f *Fox) GetSpecies() string { return f.Animal.GetSpecies() }
func (f *Fox) GetFoodType() string { return f.Animal.GetFoodType() }
func (f *Fox) IsAlive() bool { return f.Animal.IsAlive() }