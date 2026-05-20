// Package wheel 幸運轉盤系統（DAY-084）
// 擊殺特殊目標（T103流星/T104金草/B001BOSS）有機率觸發轉盤
// 轉盤有 8 個格子，倍率從 2x 到 100x
package wheel

import (
	"math/rand"
	"sync"
	"time"
)

// Slot 轉盤格子
type Slot struct {
	Multiplier float64 `json:"multiplier"`
	Label      string  `json:"label"`
	Color      string  `json:"color"`
	Weight     int     `json:"weight"`
}

// Slots 轉盤格子定義（8格）
var Slots = []Slot{
	{Multiplier: 2, Label: "2x", Color: "#4CAF50", Weight: 30},
	{Multiplier: 5, Label: "5x", Color: "#8BC34A", Weight: 25},
	{Multiplier: 10, Label: "10x", Color: "#FFC107", Weight: 20},
	{Multiplier: 20, Label: "20x", Color: "#FF9800", Weight: 12},
	{Multiplier: 30, Label: "30x", Color: "#FF5722", Weight: 7},
	{Multiplier: 50, Label: "50x", Color: "#E91E63", Weight: 4},
	{Multiplier: 80, Label: "80x", Color: "#9C27B0", Weight: 1},
	{Multiplier: 100, Label: "100x", Color: "#FFD700", Weight: 1},
}

// TriggerChance 各目標觸發轉盤的機率
var TriggerChance = map[string]float64{
	"T103": 0.15,
	"T104": 0.20,
	"B001": 0.50,
}

var totalWeight int

func init() {
	for _, s := range Slots {
		totalWeight += s.Weight
	}
}

// Manager 轉盤管理器
type Manager struct {
	mu  sync.Mutex
	rng *rand.Rand
}

// NewManager 建立新轉盤管理器
func NewManager() *Manager {
	return &Manager{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// ShouldTrigger 判斷是否觸發轉盤
func (m *Manager) ShouldTrigger(defID string) bool {
	chance, ok := TriggerChance[defID]
	if !ok {
		return false
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.rng.Float64() < chance
}

// Spin 執行轉盤，回傳中獎格子索引和格子資訊
func (m *Manager) Spin() (slotIndex int, slot Slot) {
	m.mu.Lock()
	defer m.mu.Unlock()
	r := m.rng.Intn(totalWeight)
	cumulative := 0
	for i, s := range Slots {
		cumulative += s.Weight
		if r < cumulative {
			return i, s
		}
	}
	return 0, Slots[0]
}

// SpinResult 轉盤結果
type SpinResult struct {
	SlotIndex   int     `json:"slot_index"`
	Multiplier  float64 `json:"multiplier"`
	Label       string  `json:"label"`
	Color       string  `json:"color"`
	BaseReward  int     `json:"base_reward"`
	FinalReward int     `json:"final_reward"`
}
