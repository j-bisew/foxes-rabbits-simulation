package entities

import (
	"github.com/j-bisew/foxes-rabbits-simulation/geom"
	"github.com/j-bisew/foxes-rabbits-simulation/world"
)
type Rabbit struct {
	Animal
} 

func NewRabbit(x, y float64) *Rabbit {
	return &Rabbit{
		Animal: Animal{
			Pos: geom.Point{X: x, Y: y},
			Energy: 150.0,
			MaxEnergy: 250.0,
			HungerTreshold: 30.0,
			ReproductionTreshold: 80.0,
			WellFedCooldown: 0,
			ReproduceCooldown: 0,
			SearchRadius: 10.0,
			Species: "rabbit",
			FoodType: "grass",
			MovementSpeed: 3.5,
			Alive: true,
		},
	}
}

func (r *Rabbit) GetPosition() geom.Point { return r.Animal.GetPosition() }
func (r *Rabbit) GetEnergy() float64 { return r.Animal.GetEnergy() }
func (r *Rabbit) GetSpecies() string { return r.Animal.GetSpecies() }
func (r *Rabbit) GetFoodType() string { return r.Animal.GetFoodType() }
func (r *Rabbit) IsAlive() bool { return r.Animal.IsAlive() }

func (r *Rabbit) IsHungry() bool { return r.Animal.IsHungry() }
func (r *Rabbit) IsWellFed() bool { return r.Animal.IsWellFed() }
func (r *Rabbit) UpdateWellFed() { r.Animal.UpdateWellFed() }
func (r *Rabbit) Eat(entity Entity, world *world.World) { r.Animal.Eat(entity, world) }

func (r *Rabbit) CanReproduce() bool { return r.Animal.CanReproduce() }
func (r *Rabbit) UpdateReproduce() { r.Animal.UpdateReproduce() }
func (r *Rabbit) Reproduce(mate *Animal, world *world.World) { r.Animal.Reproduce(mate, world) }

func (r *Rabbit) ConsumeEnergy() { r.Animal.ConsumeEnergy() }

func (r *Rabbit) Move(dx, dy float64) { r.Animal.Move(dx, dy) }
func (r *Rabbit) MoveRandomly() { r.Animal.MoveRandomly() }
func (r *Rabbit) DistanceTo(target geom.Point) (float64, float64, float64) { return r.Animal.DistanceTo(target) }
func (r *Rabbit) MoveTowards(target geom.Point) { r.Animal.MoveTowards(target) }

func (r *Rabbit) findClosest(entities []Entity) Entity { return r.Animal.findClosest(entities) }

func (r *Rabbit) search(world *world.World, searchType string) { r.Animal.search(world, searchType) }

func (r *Rabbit) Update(world *world.World) { r.Animal.Update(world) }