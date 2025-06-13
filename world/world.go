package world

type World struct {
	
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