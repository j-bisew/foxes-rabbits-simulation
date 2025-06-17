package gui

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type SetupPage struct {
	container *fyne.Container
	
	grassEntry *widget.Entry
	rabbitEntry *widget.Entry
	foxEntry *widget.Entry
	startBtn *widget.Button
	
	onStartCallback func(grassCount, rabbitCount, foxCount int, spawnMode string)
}

type SetupConfig struct {
	GrassCount  int
	RabbitCount int
	FoxCount    int
	SpawnMode   string
}

func NewSetupPage(onStart func(grassCount, rabbitCount, foxCount int, spawnMode string)) *SetupPage {
	setup := &SetupPage{
		onStartCallback: onStart,
	}
	
	setup.buildUI()
	return setup
}

func (s *SetupPage) buildUI() {
	title := widget.NewCard("Foxes & Rabbits Ecosystem", "Simulation Setup", 
		widget.NewRichTextFromMarkdown(`Configure your ecosystem parameters before starting the simulation.

**Instructions:**
- Set initial grass coverage (1% to 60%)
- Set the initial number of rabbits and foxes
- All animals spawn randomly across the board
- Grass grows randomly during simulation`))

	s.grassEntry = widget.NewEntry()
	s.grassEntry.SetText("30")
	s.grassEntry.Validator = func(text string) error {
		if val, err := strconv.ParseFloat(text, 64); err != nil || val < 1.0 || val > 60.0 {
			return fmt.Errorf("must be between 1%% and 60%%")
		}
		return nil
	}
	grassForm := container.NewBorder(nil, nil, 
		widget.NewRichTextFromMarkdown("**Initial grass coverage (%):**"), nil, s.grassEntry)

	s.rabbitEntry = widget.NewEntry()
	s.rabbitEntry.SetText("20")
	s.rabbitEntry.Validator = func(text string) error {
		if val, err := strconv.Atoi(text); err != nil || val < 1 {
			return fmt.Errorf("must be a number ≥ 1")
		}
		return nil
	}
	rabbitForm := container.NewBorder(nil, nil, 
		widget.NewRichTextFromMarkdown("**Rabbits:**"), nil, s.rabbitEntry)

	s.foxEntry = widget.NewEntry()
	s.foxEntry.SetText("5")
	s.foxEntry.Validator = func(text string) error {
		if val, err := strconv.Atoi(text); err != nil || val < 1 {
			return fmt.Errorf("must be a number ≥ 1")
		}
		return nil
	}
	foxForm := container.NewBorder(nil, nil, 
		widget.NewRichTextFromMarkdown("**Foxes:**"), nil, s.foxEntry)

	spawnInfo := widget.NewCard("Spawn Info", "", 
		widget.NewRichTextFromMarkdown(`**Animal Spawning:** All rabbits and foxes spawn randomly across the board

**Grass Growth:** Grass spawns initially at the specified percentage, then grows randomly during simulation with 0.2% chance per empty cell each tick

**Maximum Grass:** Grass coverage is automatically capped at 70% of the board to prevent overcrowding`))

	s.startBtn = widget.NewButton("Start Simulation", s.onStartClicked)
	s.startBtn.Importance = widget.HighImportance

	formContainer := container.NewVBox(
		grassForm,
		widget.NewSeparator(),
		rabbitForm,
		widget.NewSeparator(), 
		foxForm,
		widget.NewSeparator(),
		s.startBtn,
	)

	leftColumn := container.NewVBox(
		title,
		widget.NewCard("Parameters", "", formContainer),
	)

	rightColumn := container.NewVBox(
		spawnInfo,
	)

	mainContent := container.NewHSplit(leftColumn, rightColumn)
	mainContent.SetOffset(0.6)

	s.container = container.NewCenter(mainContent)
}

func (s *SetupPage) onStartClicked() {
	if err := s.grassEntry.Validate(); err != nil {
		return
	}
	if err := s.rabbitEntry.Validate(); err != nil {
		return
	}
	if err := s.foxEntry.Validate(); err != nil {
		return
	}

	grassPercentage, _ := strconv.ParseFloat(s.grassEntry.Text, 64)
	rabbitCount, _ := strconv.Atoi(s.rabbitEntry.Text)
	foxCount, _ := strconv.Atoi(s.foxEntry.Text)

	if s.onStartCallback != nil {
		s.onStartCallback(int(grassPercentage*100), rabbitCount, foxCount, "Random")
	}
}

func (s *SetupPage) GetContainer() *fyne.Container {
	return s.container
}

func (s *SetupPage) GetConfig() SetupConfig {
	grassPercentage, _ := strconv.ParseFloat(s.grassEntry.Text, 64)
	rabbitCount, _ := strconv.Atoi(s.rabbitEntry.Text)
	foxCount, _ := strconv.Atoi(s.foxEntry.Text)
	
	return SetupConfig{
		GrassCount:  int(grassPercentage * 100),
		RabbitCount: rabbitCount,
		FoxCount:    foxCount,
		SpawnMode:   "Random",
	}
}