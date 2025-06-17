package entities

import (
	"math"
	"math/rand"

	"github.com/j-bisew/foxes-rabbits-simulation/geom"
	"github.com/j-bisew/foxes-rabbits-simulation/interfaces"
)
type Entity = interfaces.Entity
type WorldInterface = interfaces.WorldInterface

type Animal struct {
	Pos geom.Point
	Energy, MaxEnergy, EnergyLoss float64
	CriticalHungerThreshold float64
	ReproduceCooldown int
	SearchRadius float64
	Species string
	FoodType string
	MovementSpeed float64
	Alive bool
} 

// Getters
func (a *Animal) GetPosition() geom.Point { return a.Pos }
func (a *Animal) GetEnergy() float64 { return a.Energy }
func (a *Animal) GetSpecies() string { return a.Species }
func (a *Animal) GetFoodType() string { return a.FoodType }
func (a *Animal) IsAlive() bool { return a.Alive }
func (a *Animal) Kill() { a.Alive = false }


// Hunger & Reproduction
func (a *Animal) ConsumeEnergy() { a.Energy -= a.EnergyLoss }

func (a *Animal) CanReproduce() bool { return a.ReproduceCooldown == 0 }
func (a *Animal) UpdateReproduce() { 
    if a.ReproduceCooldown > 0 {
        a.ReproduceCooldown--
    }
}
func (a *Animal) UpdateEnergy(amount float64) { 
    if a.Energy > 0 {
        a.Energy += amount
    }
}

// Movement & Search
func (a *Animal) Move(dx, dy float64) {
	a.Pos.X += dx
	a.Pos.Y += dy
}

func (a *Animal) MoveRandomly() {
	angle := rand.Float64() * 2 * math.Pi
	dx := math.Cos(angle) * a.MovementSpeed
	dy := math.Sin(angle) * a.MovementSpeed
	a.Move(dx, dy)
}

func (a *Animal) DistanceTo(target geom.Point) (float64, float64, float64) {
	dx := target.X - a.Pos.X
	dy := target.Y - a.Pos.Y
	
	return dx, dy, math.Sqrt(dx*dx + dy*dy)
}

func (a *Animal) MoveTowards(target geom.Point) {
	dx, dy, distance := a.DistanceTo(target)
	if distance > 0 {
		dx = (dx / distance) * a.MovementSpeed
		dy = (dy / distance) * a.MovementSpeed
		a.Move(dx, dy)
	}
}

func (a *Animal) findClosest(entities []Entity) Entity {
	if len(entities) == 0 {
		return nil
	}

	closest := entities[0]
	_,_,minDistance := a.DistanceTo(closest.GetPosition())

	for _, entity := range entities[1:] {
		_,_,distance := a.DistanceTo(entity.GetPosition())
		if distance < minDistance {
			minDistance = distance
			closest = entity
		}
	}

	return closest
}

func (a *Animal) search(world WorldInterface, searchType string) {
    var targets []Entity
    if searchType == "food" {
        targets = world.FindNearbyEntities(a.Pos, a.SearchRadius, a.FoodType)
	} else if searchType == "mate" {
		nearby := world.FindNearbyEntities(a.Pos, a.SearchRadius, a.Species)
		
		for _, entity := range nearby {
			if entity != Entity(a) && entity.IsAlive() {
				if rabbit, ok := entity.(*Rabbit); ok && rabbit.ReproduceCooldown == 0 {
					targets = append(targets, entity)
				} else if fox, ok := entity.(*Fox); ok && fox.ReproduceCooldown == 0 {
					targets = append(targets, entity)
				}
			}
		}
	}
	
	if len(targets) > 0 {
		closest := a.findClosest(targets)
		a.MoveTowards(closest.GetPosition())

		_, _, distance := a.DistanceTo(closest.GetPosition())
		if distance < 5.0 {
			if searchType == "food" {
				energyGained := world.ConsumeFood(closest, Entity(a))
				a.Energy += energyGained
				if a.Energy > a.MaxEnergy { a.Energy = a.MaxEnergy }
			} else if searchType == "mate" {
				world.CreateOffspring(Entity(a), closest)
				a.ReproduceCooldown = 40
				a.UpdateEnergy(-10.0)
				closest.UpdateEnergy(-10.0)
			}
		}
	} else {
		a.MoveRandomly()
	}
}

// Main Behavior
func (a *Animal) Update(world WorldInterface) {
    a.ConsumeEnergy()
	
    if a.Energy <= 0 {
        a.Alive = false
        return
    }

	a.UpdateReproduce()
    
    if a.Energy < a.CriticalHungerThreshold {
        a.search(world, "food")
    } else if a.ReproduceCooldown == 0 {
        a.search(world, "mate")
    } else {
        a.search(world, "food")
    }
}
