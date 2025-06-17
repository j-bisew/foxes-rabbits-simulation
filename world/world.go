package world

import (
	"math/rand"

	"github.com/j-bisew/foxes-rabbits-simulation/interfaces"
	"github.com/j-bisew/foxes-rabbits-simulation/geom"
	"github.com/j-bisew/foxes-rabbits-simulation/quadtree"
	"github.com/j-bisew/foxes-rabbits-simulation/entities"
)

type Entity = interfaces.Entity

type World struct {
	Width, Height int
	QuadTree      *quadtree.QuadTree
	Entities      []Entity
	GrassSpawnRate    float64
	MaxGrassCount     int
}

func NewWorld(width, height int) *World {
	boundary := geom.Rectangle{
		X: 0, Y: 0,
		Width: float64(width),
		Height: float64(height),
	}
	
	world := &World{
		Width:         width,
		Height:        height,
		QuadTree:      quadtree.NewQuadTree(10, boundary),
		Entities:      make([]Entity, 0),
		GrassSpawnRate: 0.5,
		MaxGrassCount: int(float64(width * height) * 0.70),
	}
	
	return world
}

// Grass
func (w *World) SpawnInitialGrassRandom(count int) {
	for i := 0; i < count; i++ {
		x := rand.Float64() * float64(w.Width)
		y := rand.Float64() * float64(w.Height)
		grass := entities.NewGrass(x, y)
		w.AddEntity(grass)
	}
}

func (w *World) spawnGrass() {
	grassCount := w.CountGrass()
	if grassCount >= w.MaxGrassCount {
		return
	}
	
	if rand.Float64() < w.GrassSpawnRate {
		x := rand.Float64() * float64(w.Width)
		y := rand.Float64() * float64(w.Height)
		grass := entities.NewGrass(x, y)
		w.AddEntity(grass)
	}
}

// Adders
func (w *World) AddEntity(entity Entity) {
	w.Entities = append(w.Entities, entity)
}

func (w *World) AddRabbit(x, y float64) {
	rabbit := entities.NewRabbit(x, y)
	w.AddEntity(rabbit)
}

func (w *World) AddFox(x, y float64) {
	fox := entities.NewFox(x, y)
	w.AddEntity(fox)
}

// Cleaner
func (w *World) removeDeadEntities() {
	alive := make([]Entity, 0, len(w.Entities))
	for _, entity := range w.Entities {
		if entity.IsAlive() {
			alive = append(alive, entity)
		}
	}
	w.Entities = alive
}

func (w *World) ClearEntities() {
	w.Entities = w.Entities[:0]
	w.QuadTree.Clear()
}

// Reproduction
func (w *World) CreateOffspring(parent1, parent2 Entity) Entity {
	newX := (parent1.GetPosition().X + parent2.GetPosition().X) / 2
	newY := (parent1.GetPosition().Y + parent2.GetPosition().Y) / 2
	
	newX += (rand.Float64() - 0.5) * 10.0
	newY += (rand.Float64() - 0.5) * 10.0
	
	if newX < 0 { newX = 0 }
	if newX >= float64(w.Width) { newX = float64(w.Width - 1) }
	if newY < 0 { newY = 0 }
	if newY >= float64(w.Height) { newY = float64(w.Height - 1) }
	
	var offspring Entity
	switch parent1.GetSpecies() {
	case "rabbit":
		offspring = entities.NewRabbit(newX, newY)
	case "fox":
		offspring = entities.NewFox(newX, newY)
	default:
		return nil
	}

	w.AddEntity(offspring)
	return offspring
}

// Methods for Entities
func (w *World) FindNearbyEntities(pos geom.Point, radius float64, species string) []Entity {
	searchRect := geom.Rectangle{
		X: pos.X - radius,
		Y: pos.Y - radius,
		Width: radius * 2,
		Height: radius * 2,
	}
	
	var nearbyEntities []Entity
	w.QuadTree.Query(searchRect, &nearbyEntities)
	
	var filtered []Entity
	for _, entity := range nearbyEntities {
		if entity.GetSpecies() == species && entity.IsAlive() && entity.GetEnergy() > 0{
			filtered = append(filtered, entity)
		}
	}
	
	return filtered
}

func (w *World) IsValidPosition(x, y float64) bool {
	return x >= 0 && x < float64(w.Width) && y >= 0 && y < float64(w.Height)
}

func (w *World) ConsumeFood(food Entity, eater Entity) float64 {
	switch food.GetSpecies() {
	case "grass":
		grass := food.(*entities.Grass)
		wantedAmount := 20.0 + rand.Float64()*20.0
		if grass.GetEnergy() < wantedAmount {
			return grass.Consume(grass.GetEnergy())
		}
		return grass.Consume(wantedAmount)
	case "rabbit":
		energyGain := 80.0 + (food.GetEnergy() * 0.3)
		food.Kill()
		return energyGain
	}
	return 0.0
}

// Main Method
func (w *World) Update() {
	w.QuadTree.Clear()
	
	for _, entity := range w.Entities {
		if entity.IsAlive() {
			entity.Update(w)
			w.QuadTree.Insert(entity.GetPosition(), entity)
		}
	}
	
	w.removeDeadEntities()
	w.spawnGrass()
}

// Getters for stats
func (w *World) CountRabbits() int {
	count := 0
	for _, entity := range w.Entities {
		if entity.GetSpecies() == "rabbit" && entity.IsAlive() {
			count++
		}
	}
	return count
}

func (w *World) CountFoxes() int {
	count := 0
	for _, entity := range w.Entities {
		if entity.GetSpecies() == "fox" && entity.IsAlive() {
			count++
		}
	}
	return count
}

func (w *World) CountGrass() int {
	count := 0
	for _, entity := range w.Entities {
		if entity.GetSpecies() == "grass" && entity.IsAlive() {
			count++
		}
	}
	return count
}