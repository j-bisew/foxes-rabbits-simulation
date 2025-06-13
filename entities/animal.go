package entities

import (
	"math"
	"math/rand"

	"github.com/j-bisew/foxes-rabbits-simulation/geom"
	"github.com/j-bisew/foxes-rabbits-simulation/world"
)

type Entity interface {
	Update(*world.World)
	GetPosition() geom.Point
	GetSpecies() string
	GetFoodType() string
	IsAlive() bool
	GetEnergy() float64
}

type Animal struct {
	Pos geom.Point
	Energy, MaxEnergy float64
	HungerTreshold, ReproductionTreshold float64
	WellFedCooldown, ReproduceCooldown int
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


// Hunger & Reproduction
func (a *Animal) IsHungry() bool { return a.Energy < a.HungerTreshold }
func (a *Animal) IsWellFed() bool { return a.Energy >= a.ReproductionTreshold && a.WellFedCooldown > 0 }
func (a *Animal) UpdateWellFed() { a.WellFedCooldown -= 1 }

func (a *Animal) Eat(entity Entity, world *world.World) {
	a.Energy += entity.GetEnergy()
	if a.Energy > a.MaxEnergy { a.Energy = a.MaxEnergy }
	a.WellFedCooldown = 15
	a.ReproduceCooldown = 0

	world.Consume(entity)
}

func (a *Animal) CanReproduce() bool { return a.IsWellFed() && a.ReproduceCooldown == 0 }
func (a *Animal) UpdateReproduce() { a.ReproduceCooldown -= 1 }

func (a *Animal) Reproduce(mate *Animal, world *world.World) {
	newAnimal := world.CreateOffspring(a, mate)

	a.ReproduceCooldown = 20
	mate.ReproduceCooldown = 20
}

func (a *Animal) ConsumeEnergy() {
	a.Energy -= 1.0
	if a.Energy <= 0 { a.Alive = false }
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

func (a *Animal) search(world *world.World, searchType string) {
	searchRect := geom.Rectangle{
		X: a.Pos.X - a.SearchRadius,
		Y: a.Pos.Y - a.SearchRadius,
		Width: a.SearchRadius*2,
		Height: a.SearchRadius*2,
	}

	var nearbyEntities []Entity
	world.QuadTree.Query(searchRect, &nearbyEntities)

	var targets []Entity

	if searchType == "food" {
		for _, entity := range nearbyEntities {
			if entity.GetSpecies() == a.FoodType && entity.IsAlive() {
				targets = append(targets, entity)
			}
		}
	} else if searchType == "mate" {
		for _, entity := range nearbyEntities {
			animal, ok := entity.(*Animal)
			if ok &&
				entity.GetSpecies() == a.Species &&
				entity.IsAlive() &&
				animal.CanReproduce() &&
				entity != a {
					targets = append(targets, entity)
				}
		}
	}
	
	if len(targets) > 0 {
		closest := a.findClosest(targets)
		a.MoveTowards(closest.GetPosition())

		_, _, distance := a.DistanceTo(closest.GetPosition())
		if distance < 5.0 {
			if searchType == "food" {
				a.Eat(closest, world)
			} else {
				a.Reproduce(closest.(*Animal), world)
			}
		}
	} else {
		a.MoveRandomly()
	}
}

// Main Behavior
func (a *Animal) Update(world *world.World) {
	a.ConsumeEnergy()
	if !a.IsAlive() {
		return
	}
	a.UpdateReproduce()
	a.UpdateWellFed()

	if a.IsHungry() {
		a.search(world, "food")
	} else if a.CanReproduce() {
		a.search(world, "mate")
	} else {
		a.MoveRandomly()
	}
}
