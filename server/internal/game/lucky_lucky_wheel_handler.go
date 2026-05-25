// lucky_lucky_wheel_handler.go — T125 幸運大轉盤魚系統
// server-event-agent 負責維護
// 業界依據：Jili Fishing「Lucky Wheel — spin the wheel for random multipliers, jackpots, or special effects」
// 設計：擊破 T125 後，觸發「幸運大轉盤」：轉盤有 8 個格子，隨機停在一格
// 格子內容：×2（30%）/ ×5（20%）/ ×10（15%）/ ×20（12%）/ ×50（8%）/ ×100（6%）/ 全場 HP-50%（5%）/ 大獎（4%）
// 大獎 = 當前大獎池 × 50%（最低 5000）
// 個人冷卻 18 秒；全服冷卻 30 秒
package game

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyLuckyWheelManager struct {
	mu              sync.Mutex
	playerCooldowns map[string]time.Time
	globalCooldown  time.Time
	jackpotPool     int // 大獎池
}

type wheelSlot struct {
	Name    string
	Type    string  // "mult" | "aoe" | "jackpot"
	Mult    float64 // 倍率（type=mult 時使用）
	Weight  int
}

var luckyWheelSlots = []wheelSlot{
	{Name: "×2", Type: "mult", Mult: 2, Weight: 30},
	{Name: "×5", Type: "mult", Mult: 5, Weight: 20},
	{Name: "×10", Type: "mult", Mult: 10, Weight: 15},
	{Name: "×20", Type: "mult", Mult: 20, Weight: 12},
	{Name: "×50", Type: "mult", Mult: 50, Weight: 8},
	{Name: "×100", Type: "mult", Mult: 100, Weight: 6},
	{Name: "全場 HP-50%", Type: "aoe", Mult: 1, Weight: 5},
	{Name: "大獎", Type: "jackpot", Mult: 1, Weight: 4},
}

func newLuckyLuckyWheelManager() *luckyLuckyWheelManager {
	return &luckyLuckyWheelManager{
		playerCooldowns: make(map[string]time.Time),
		jackpotPool:     20000, // 初始大獎池
	}
}

func isLuckyLuckyWheelFish(defID string) bool {
	return defID == "T125"
}

func (m *luckyLuckyWheelManager) contributeToPool(amount int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.jackpotPool += amount
}

func (m *luckyLuckyWheelManager) rollWheel() wheelSlot {
	totalWeight := 0
	for _, s := range luckyWheelSlots {
		totalWeight += s.Weight
	}
	r := rand.Intn(totalWeight)
	cumulative := 0
	for _, s := range luckyWheelSlots {
		cumulative += s.Weight
		if r < cumulative {
			return s
		}
	}
	return luckyWheelSlots[0]
}

func (g *Game) tryLuckyLuckyWheel(playerID string, killerName string) {
	m := g.luckyLuckyWheel
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCooldown) {
		m.mu.Unlock()
		return
	}
	if cd, ok := m.playerCooldowns[playerID]; ok && now.Before(cd) {
		m.mu.Unlock()
		return
	}
	m.playerCooldowns[playerID] = now.Add(18 * time.Second)
	m.globalCooldown = now.Add(30 * time.Second)
	poolSize := m.jackpotPool
	m.mu.Unlock()

	g.hub.Broadcast(protocol.MsgLuckyLuckyWheel, protocol.LuckyLuckyWheelPayload{
		Event:       "trigger",
		TriggerID:   playerID,
		TriggerName: killerName,
		PoolSize:    poolSize,
	})
	g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
		Message:  "🎡 " + killerName + " 觸發幸運大轉盤！",
		Priority: "high",
		Color:    "#FF69B4",
	})

	// 轉盤動畫延遲（2 秒）
	go func() {
		time.Sleep(2 * time.Second)

		m.mu.Lock()
		slot := m.rollWheel()
		m.mu.Unlock()

		g.mu.Lock()
		p, ok := g.players[playerID]
		if !ok {
			g.mu.Unlock()
			return
		}
		betCost := p.GetBetDef().BetCost
		reward := 0

		switch slot.Type {
		case "mult":
			reward = int(float64(betCost) * slot.Mult)
			p.AddCoins(reward)
			g.sendPlayerUpdate(playerID)

		case "aoe":
			// 全場 HP -50%
			for _, t := range g.targets {
				if t.Def.Type == "boss" {
					continue
				}
				damage := t.MaxHP / 2
				t.HP -= damage
				if t.HP < 1 {
					t.HP = 1
				}
				g.hub.Broadcast(protocol.MsgTargetUpdate, protocol.TargetUpdatePayload{
					InstanceID: t.InstanceID,
					HP:         t.HP,
					MaxHP:      t.MaxHP,
					X:          t.X,
					Y:          t.Y,
				})
			}
			reward = betCost * 5

		case "jackpot":
			m.mu.Lock()
			jackpotReward := m.jackpotPool / 2
			if jackpotReward < 5000 {
				jackpotReward = 5000
			}
			m.jackpotPool -= jackpotReward
			if m.jackpotPool < 5000 {
				m.jackpotPool = 5000
			}
			m.mu.Unlock()
			reward = jackpotReward
			p.AddCoins(reward)
			g.sendPlayerUpdate(playerID)
		}
		g.mu.Unlock()

		g.hub.Broadcast(protocol.MsgLuckyLuckyWheel, protocol.LuckyLuckyWheelPayload{
			Event:       "spin_result",
			TriggerID:   playerID,
			TriggerName: killerName,
			SlotName:    slot.Name,
			SlotType:    slot.Type,
			SlotMult:    slot.Mult,
			Reward:      reward,
			PoolSize:    g.luckyLuckyWheel.jackpotPool,
		})

		if slot.Type == "jackpot" {
			g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
				Message:  "🎡🏆 " + killerName + " 中大獎！獲得 " + string(rune('0'+reward/1000)) + "K！",
				Priority: "critical",
				Color:    "#FFD700",
			})
		} else {
			g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
				Message:  "🎡 " + killerName + " 轉到 " + slot.Name + "！獎勵 " + string(rune('0'+reward)) + "！",
				Priority: "normal",
				Color:    "#FF69B4",
			})
		}

		log.Printf("[LuckyWheel] Player %s: slot=%s, reward=%d", playerID, slot.Name, reward)
	}()
}
