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

type DayStats struct {
	TotalDrunk int     `json:"total_drunk"`
	DailyGoal  int     `json:"daily_goal"`
	Missing    int     `json:"missing"`
	Events     []Event `json:"events"`
}

type WaterFile struct {
	LastUpdated string              `json:"last_updated"`
	Events      map[string]DayStats `json:"events"`
}

var currentStats DayStats
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

func loadStats() DayStats {
	today := time.Now().Format("2006-01-02")

	file, err := os.ReadFile(configPath)
	if err != nil {
		return DayStats{DailyGoal: 4250, Missing: 4250, Events: []Event{}}
	}

	var wf WaterFile
	_ = json.Unmarshal(file, &wf)

	if wf.Events == nil {
		return DayStats{DailyGoal: 4250, Missing: 4250, Events: []Event{}}
	}

	day, ok := wf.Events[today]
	if !ok {
		dailyGoal := 4250
		var lastDate string
		for date, d := range wf.Events {
			if date > lastDate {
				lastDate = date
				dailyGoal = d.DailyGoal
			}
		}
		return DayStats{DailyGoal: dailyGoal, Missing: dailyGoal, Events: []Event{}}
	}

	day.Missing = max(day.DailyGoal-day.TotalDrunk, 0)
	return day
}

func saveStats() {
	today := time.Now().Format("2006-01-02")

	var wf WaterFile
	file, err := os.ReadFile(configPath)
	if err == nil {
		_ = json.Unmarshal(file, &wf)
	}
	if wf.Events == nil {
		wf.Events = make(map[string]DayStats)
	}

	currentStats.Missing = max(currentStats.DailyGoal-currentStats.TotalDrunk, 0)

	wf.LastUpdated = time.Now().Format(time.RFC3339)
	wf.Events[today] = currentStats

	data, _ := json.MarshalIndent(wf, "", "  ")
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
