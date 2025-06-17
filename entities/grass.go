package entities

import (
	"math/rand"

	"github.com/j-bisew/foxes-rabbits-simulation/geom"
	"github.com/j-bisew/foxes-rabbits-simulation/interfaces"
)

type Grass struct {
	Pos geom.Point
	Amount, MaxAmount, GrowthRate float64
	Alive bool
}

func NewGrass(x,y float64) *Grass {
	return &Grass{
		Pos: geom.Point{X: x, Y: y},
		Amount: 0.0,
		MaxAmount: 50.0 + rand.Float64()*50.0,
		GrowthRate: 0.5 + rand.Float64(),
		Alive: true,
	}
}

func (g *Grass) GetPosition() geom.Point { return g.Pos }
func (g *Grass) GetSpecies() string      { return "grass" }
func (g *Grass) GetFoodType() string     { return "" }
func (g *Grass) IsAlive() bool           { return g.Alive }
func (g *Grass) GetEnergy() float64      { return g.Amount }
func (g *Grass) UpdateEnergy(amount float64)  {
	g.Amount += amount
}
func (g *Grass) Kill() { g.Alive = false }

func (g *Grass) Update(world interfaces.WorldInterface) {
	g.Amount += g.GrowthRate
	if g.Amount > g.MaxAmount { g.Amount = g.MaxAmount }

	if g.Amount <= 0 { g.Alive = false }
}

func (g *Grass) Consume(wantedEnergy float64) float64 {
	if g.Amount <= 0 { return 0.0 }

	g.Amount -= wantedEnergy
	if g.Amount <= 0.0 {
		g.Kill()
	}

	return wantedEnergy
}