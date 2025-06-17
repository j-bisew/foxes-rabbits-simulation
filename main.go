package main

import (
	"flag"

	"github.com/j-bisew/foxes-rabbits-simulation/gui"
	"github.com/j-bisew/foxes-rabbits-simulation/world"
)

func main() {
    width := flag.Int("w", 400, "width of board")
    height := flag.Int("h", 200, "height of board")

    flag.Parse()

    world := world.NewWorld(*width, *height)
    
    gui := gui.NewGUI(world)
    gui.Run()
}