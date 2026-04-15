package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
)

//go:embed assets/icon.png
var iconData []byte

type Event struct {
	Amount    int    `json:"amount"`
	Timestamp string `json:"timestamp"`
}

type Stats struct {
	TotalDrunk  int     `json:"total_drunk"`
	DailyGoal   int     `json:"daily_goal"`
	Missing     int     `json:"missing"` // Novo campo: quanto falta para a meta
	LastUpdated string  `json:"last_updated"`
	Events      []Event `json:"events"`
}

var currentStats Stats
var configPath string

func init() {
	home, _ := os.UserHomeDir()
	configDir := filepath.Join(home, ".config", "rdw")
	configPath = filepath.Join(configDir, "water.json")

	_ = os.MkdirAll(configDir, 0755)
	currentStats = loadStats()
}

func main() {
	systray.Run(onReady, onExit)
}

func loadStats() Stats {
	file, err := os.ReadFile(configPath)
	if err != nil {
		// Meta de 4250ml [cite: 2026-03-06]
		return Stats{
			TotalDrunk:  0,
			DailyGoal:   4250,
			Missing:     4250,
			LastUpdated: time.Now().Format("2006-01-02"),
			Events:      []Event{},
		}
	}
	var s Stats
	_ = json.Unmarshal(file, &s)
	// Recalcula o missing no load para garantir consistência
	s.Missing = s.DailyGoal - s.TotalDrunk
	if s.Missing < 0 {
		s.Missing = 0
	}
	return s
}

func saveStats() {
	currentStats.LastUpdated = time.Now().Format("2006-01-02 15:04:05")

	// Atualiza o missing antes de salvar
	currentStats.Missing = currentStats.DailyGoal - currentStats.TotalDrunk
	if currentStats.Missing < 0 {
		currentStats.Missing = 0
	}

	data, _ := json.MarshalIndent(currentStats, "", "  ")
	_ = os.WriteFile(configPath, data, 0644)
}

func onReady() {
	systray.SetIcon(iconData)
	systray.SetTooltip(fmt.Sprintf("RDW - Meta: %dml", currentStats.DailyGoal))

	mStatus := systray.AddMenuItem(fmt.Sprintf("Progresso: %d/%dml", currentStats.TotalDrunk, currentStats.DailyGoal), "")
	mStatus.Disable()
	systray.AddSeparator()

	m80 := systray.AddMenuItem("Golão (80ml)", "Benchmark 80ml")
	m1000 := systray.AddMenuItem("Garrafa (1000ml)", "Tupppaware 1L")
	m1500 := systray.AddMenuItem("Garrafona (1500ml)", "Normalzona 1.5L")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Sair", "")

	go func() {
		ticker := time.NewTicker(15 * time.Minute)
		for range ticker.C {
			if currentStats.TotalDrunk < currentStats.DailyGoal {
				beeep.Alert("Hidratação!", fmt.Sprintf("Faltam %dml para bater sua meta!", currentStats.Missing), "")
			}
		}
	}()

	go func() {
		for {
			select {
			case <-m80.ClickedCh:
				addWater(80, mStatus)
			case <-m1000.ClickedCh:
				addWater(1000, mStatus)
			case <-m1500.ClickedCh:
				addWater(1500, mStatus)
			case <-mQuit.ClickedCh:
				systray.Quit()
			}
		}
	}()
}

func addWater(amount int, item *systray.MenuItem) {
	event := Event{
		Amount:    amount,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	currentStats.Events = append(currentStats.Events, event)
	currentStats.TotalDrunk += amount
	saveStats() // saveStats já recalcula o missing

	percentage := (float64(currentStats.TotalDrunk) / float64(currentStats.DailyGoal)) * 100
	item.SetTitle(fmt.Sprintf("Progresso: %d/%dml (%.1f%%)", currentStats.TotalDrunk, currentStats.DailyGoal, percentage))

	beeep.Notify("RDW", fmt.Sprintf("Faltam %dml. Total: %dml", currentStats.Missing, currentStats.TotalDrunk), "")
}

func onExit() {}
