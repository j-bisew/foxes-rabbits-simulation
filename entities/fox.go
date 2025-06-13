package entities

import (
	"github.com/j-bisew/foxes-rabbits-simulation/geom"
	"github.com/j-bisew/foxes-rabbits-simulation/world"
)

type Fox struct {
	Animal
} 

func NewFox(x, y float64) *Fox {
	return &Fox{
		Animal: Animal{
			Pos: geom.Point{X: x, Y: y},
			Energy: 150.0,
			MaxEnergy: 300.0,
			HungerTreshold: 50.0,
			ReproductionTreshold: 120.0,
			WellFedCooldown: 0,
			ReproduceCooldown: 0,
			SearchRadius: 20.0,
			Species: "fox",
			FoodType: "rabbit",
			MovementSpeed: 2.0,
			Alive: true,
		},
	}
}


func (f *Fox) GetPosition() geom.Point { return f.Animal.GetPosition() }
func (f *Fox) GetEnergy() float64 { return f.Animal.GetEnergy() }
func (f *Fox) GetSpecies() string { return f.Animal.GetSpecies() }
func (f *Fox) GetFoodType() string { return f.Animal.GetFoodType() }
func (f *Fox) IsAlive() bool { return f.Animal.IsAlive() }

func (f *Fox) IsHungry() bool { return f.Animal.IsHungry() }
func (f *Fox) IsWellFed() bool { return f.Animal.IsWellFed() }
func (f *Fox) UpdateWellFed() { f.Animal.UpdateWellFed() }
func (f *Fox) Eat(entity Entity, world *world.World) { f.Animal.Eat(entity, world) }

func (f *Fox) CanReproduce() bool { return f.Animal.CanReproduce() }
func (f *Fox) UpdateReproduce() { f.Animal.UpdateReproduce() }
func (f *Fox) Reproduce(mate *Animal, world *world.World) { f.Animal.Reproduce(mate, world) }

func (f *Fox) ConsumeEnergy() { f.Animal.ConsumeEnergy() }

func (f *Fox) Move(dx, dy float64) { f.Animal.Move(dx, dy) }
func (f *Fox) MoveRandomly() { f.Animal.MoveRandomly() }
func (f *Fox) DistanceTo(target geom.Point) (float64, float64, float64) { return f.Animal.DistanceTo(target) }
func (f *Fox) MoveTowards(target geom.Point) { f.Animal.MoveTowards(target) }

func (f *Fox) findClosest(entities []Entity) Entity { return f.Animal.findClosest(entities) }

func (f *Fox) search(world *world.World, searchType string) { f.Animal.search(world, searchType) }

func (f *Fox) Update(world *world.World) { f.Animal.Update(world) }