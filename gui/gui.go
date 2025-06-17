package gui

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	
	"github.com/j-bisew/foxes-rabbits-simulation/world"
)

type GUI struct {
	app fyne.App
	window fyne.Window
	world *world.World

	setupPage *SetupPage
	simContainer *fyne.Container

	gameCanvas *canvas.Raster
	chart *canvas.Raster
	statsLabel *widget.Label
	startBtn *widget.Button
	stopBtn *widget.Button
	backBtn *widget.Button

	running bool
	ticker *time.Ticker

	rabbitHistory []int
	foxHistory []int
	grassHistory []int
	maxHistory int
}

func NewGUI(w *world.World) *GUI {
	myApp := app.New()
	window := myApp.NewWindow("Foxes & Rabbits Ecosystem Simulation")
	window.Resize(fyne.NewSize(1200, 800))

	gui := &GUI{
		app: myApp,
		window: window,
		world: w,
		running: false,
		maxHistory: 200,
		rabbitHistory: make([]int, 0, 200),
		foxHistory: make([]int, 0, 200),
		grassHistory: make([]int, 0, 200),
	}

	gui.setupUI()
	gui.showSetupPage()
	return gui
}

func (g *GUI) setupUI() {
	g.setupPage = NewSetupPage(g.onConfigurationComplete)
	g.setupSimulationPage()
}

func (g *GUI) setupSimulationPage() {
	g.gameCanvas = canvas.NewRaster(g.drawGame)
	g.gameCanvas.Resize(fyne.NewSize(600, 300))

	g.chart = canvas.NewRaster(g.drawChart)
	g.chart.Resize(fyne.NewSize(600, 200))

	g.statsLabel = widget.NewLabel("Rabbits: 0, Foxes: 0, Grass: 0")
	g.statsLabel.TextStyle = fyne.TextStyle{Bold: true}

	g.startBtn = widget.NewButton("Start", g.startSimulation)
	g.stopBtn = widget.NewButton("Stop", g.stopSimulation)
	g.backBtn = widget.NewButton("Back to Setup", g.showSetupPage)

	g.stopBtn.Disable()

	controlsContainer := container.NewHBox(
		g.startBtn,
		g.stopBtn,
		widget.NewSeparator(),
		g.backBtn,
		widget.NewSeparator(),
		g.statsLabel,
	)

	gameContainer := container.NewBorder(
		widget.NewLabel("Game Board"),
		nil, nil, nil,
		g.gameCanvas,
	)
	
	chartContainer := container.NewBorder(
		widget.NewLabel("Population Chart"),
		nil, nil, nil,
		g.chart,
	)

	g.simContainer = container.NewBorder(
		nil,
		controlsContainer,
		nil,
		nil,
		container.NewVSplit(
			gameContainer,
			chartContainer,
		),
	)
}

func (g *GUI) showSetupPage() {
	if g.running {
		g.stopSimulation()
		time.Sleep(100 * time.Millisecond)
	}
	
	g.world.ClearEntities()
	
	g.rabbitHistory = g.rabbitHistory[:0]
	g.foxHistory = g.foxHistory[:0]
	g.grassHistory = g.grassHistory[:0]
	
	totalCells := g.world.Width * g.world.Height
	g.world.MaxGrassCount = int(float64(totalCells) * 0.70)
	g.world.GrassSpawnRate = 0.002
	
	g.window.SetContent(g.setupPage.GetContainer())
}

func (g *GUI) showSimulationPage() {
	g.window.SetContent(g.simContainer)
	g.updateStats()
}

func (g *GUI) onConfigurationComplete(grassPercentageBasisPoints, rabbitCount, foxCount int, spawnMode string) {
	g.initializeWorld(grassPercentageBasisPoints, rabbitCount, foxCount, spawnMode)
	g.showSimulationPage()
}

func (g *GUI) initializeWorld(grassPercentageBasisPoints, rabbitCount, foxCount int, spawnMode string) {
	g.world.ClearEntities()
	
	totalCells := g.world.Width * g.world.Height
	g.world.MaxGrassCount = int(float64(totalCells) * 0.70)
	
	g.rabbitHistory = g.rabbitHistory[:0]
	g.foxHistory = g.foxHistory[:0]
	g.grassHistory = g.grassHistory[:0]

	maxGrass := g.world.MaxGrassCount
	requestedGrass := int(float64(totalCells) * float64(grassPercentageBasisPoints) / 10000.0)
	
	grassCount := requestedGrass
	if grassCount > maxGrass {
		grassCount = maxGrass
	}

	g.world.SpawnInitialGrassRandom(grassCount)
	
	for i := 0; i < rabbitCount; i++ {
		x := rand.Float64() * float64(g.world.Width)
		y := rand.Float64() * float64(g.world.Height)
		g.world.AddRabbit(x, y)
	}

	for i := 0; i < foxCount; i++ {
		x := rand.Float64() * float64(g.world.Width)
		y := rand.Float64() * float64(g.world.Height)
		g.world.AddFox(x, y)
	}
}

func (g *GUI) drawGame(w,h int) image.Image {
	img := image.NewRGBA(image.Rect(0,0,w,h))

	for y := 0; y<h; y++ {
		for x := 0; x<w; x++ {
			img.Set(x,y, color.RGBA{0,0,0,255})
		}
	}

	for _,entity := range g.world.Entities {
		if !entity.IsAlive() { continue }

		pos := entity.GetPosition()
		x := int(pos.X * float64(w) / float64(g.world.Width))
		y := int(pos.Y * float64(h) / float64(g.world.Height))

		if x >= 0 && x < w && y >= 0 && y < h {
			var c color.RGBA

			switch entity.GetSpecies() {
			case "rabbit":
				c = color.RGBA{150,150,150,255}
			case "fox":
				c = color.RGBA{255,100,100,255}
			case "grass":
				intensity := uint8(entity.GetEnergy() * 2.5)
				if intensity > 255 {
					intensity = 255
				}
				c = color.RGBA{0, intensity, 0, 255}
			}

			img.Set(x, y, c)
		}
	}
	return img
}

func (g *GUI) drawChart(w, h int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, color.RGBA{240, 240, 240, 255})
		}
	}
	
	if len(g.rabbitHistory) < 2 {
		return img
	}
	
	maxPop := 1
	for _, count := range g.rabbitHistory {
		if count > maxPop { maxPop = count }
	}
	for _, count := range g.foxHistory {
		if count > maxPop { maxPop = count }
	}
	for _, count := range g.grassHistory {
		if count/10 > maxPop { maxPop = count/10 }
	}
	
	if maxPop == 0 { maxPop = 1 }
	
	g.drawLine(img, g.rabbitHistory, color.RGBA{150, 150, 150, 255}, w, h, maxPop)
	g.drawLine(img, g.foxHistory, color.RGBA{255, 0, 0, 255}, w, h, maxPop)
	
	grassScaled := make([]int, len(g.grassHistory))
	for i, count := range g.grassHistory {
		grassScaled[i] = count / 10
	}
	g.drawLine(img, grassScaled, color.RGBA{0, 150, 0, 255}, w, h, maxPop)
	
	img.Set(10, 10, color.RGBA{150, 150, 150, 255})
	img.Set(10, 20, color.RGBA{255, 0, 0, 255})
	img.Set(10, 30, color.RGBA{0, 150, 0, 255})
	
	return img
}

func (g *GUI) drawLine(img *image.RGBA, data []int, col color.RGBA, w, h, maxVal int) {
	if len(data) < 2 {
		return
	}
	
	for i := 1; i < len(data); i++ {
		x1 := (i-1) * w / len(data)
		y1 := h - (data[i-1] * h / maxVal)
		x2 := i * w / len(data)
		y2 := h - (data[i] * h / maxVal)
		
		g.drawLineSegment(img, x1, y1, x2, y2, col, w, h)
	}
}

func (g *GUI) drawLineSegment(img *image.RGBA, x1, y1, x2, y2 int, col color.RGBA, w, h int) {
	steps := abs(x2-x1) + abs(y2-y1)
	if steps == 0 {
		steps = 1
	}
	
	for i := 0; i <= steps; i++ {
		x := x1 + (x2-x1)*i/steps
		y := y1 + (y2-y1)*i/steps
		
		if x >= 0 && x < w && y >= 0 && y < h {
			img.Set(x, y, col)
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (g *GUI) updateStats() {
	rabbits := g.world.CountRabbits()
	foxes := g.world.CountFoxes()
	grass := g.world.CountGrass()

	g.rabbitHistory = append(g.rabbitHistory, rabbits)
	g.foxHistory = append(g.foxHistory, foxes)
	g.grassHistory = append(g.grassHistory, grass)

	if len(g.rabbitHistory) > g.maxHistory {
		g.rabbitHistory = g.rabbitHistory[1:]
		g.foxHistory = g.foxHistory[1:]
		g.grassHistory = g.grassHistory[1:]
	}

	fyne.Do(func() {
		maxGrass := g.world.MaxGrassCount
		grassPercent := float64(grass) / float64(g.world.Width * g.world.Height) * 100
		totalEntities := len(g.world.Entities)
		
		g.statsLabel.SetText(fmt.Sprintf("Rabbits: %d, Foxes: %d, Grass: %d/%d (%.1f%%) | Total Entities: %d", 
			rabbits, foxes, grass, maxGrass, grassPercent, totalEntities))
		g.gameCanvas.Refresh()
		g.chart.Refresh()
	})
}

func (g *GUI) startSimulation() {
	if g.running {
		return
	}

	if g.world == nil || len(g.world.Entities) == 0 {
		return
	}

	g.running = true
	g.startBtn.Disable()
	g.stopBtn.Enable()

	g.ticker = time.NewTicker(200 * time.Millisecond)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Simulation panic recovered: %v\n", r)
				g.running = false
			}
		}()
		
		for range g.ticker.C {
			if !g.running { 
				return 
			}
			
			if g.world == nil || g.world.QuadTree == nil {
				g.running = false
				return
			}
			
			g.world.Update()
			g.updateStats()
		}
	}()
}

func (g *GUI) stopSimulation() {
	if !g.running {
		return
	}
	
	g.running = false

	if g.ticker != nil { 
		g.ticker.Stop()
		g.ticker = nil
	}

	fyne.Do(func() {
		g.startBtn.Enable()
		g.stopBtn.Disable()
	})
}

func (g *GUI) Run() {
	g.window.ShowAndRun()
}