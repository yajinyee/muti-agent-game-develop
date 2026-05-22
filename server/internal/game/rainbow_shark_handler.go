// rainbow_shark_handler.go — 彩虹鯊魚爆發系統 handler（DAY-180）
// 業界依據：JILI 2026 新特性「Rainbow Shark — triggers a rainbow burst that randomly assigns
// 1.5x-3x multiplier bonuses to all targets on screen for 10 seconds」
// 擊破 T138 後觸發「彩虹爆發」：
//   1. 場上所有存活目標隨機獲得 1.5x/2.0x/2.5x/3.0x 倍率加成標記（加權隨機）
//   2. 持續 10 秒，玩家擊破標記目標可獲得額外倍率加成
//   3. 全服廣播「彩虹鯊魚爆發」，讓所有玩家看到每個目標的倍率標記
// 設計差異：
//   - 與黃金鯊魚（全服固定 ×1.5）不同，彩虹鯊魚是「每個目標倍率不同」（1.5x-3x），
//     製造「哪個目標倍率最高？快去打！」的策略感
//   - 與幸運星魚（個人 ×2，10秒）不同，彩虹鯊魚是「全服共享，但每個目標倍率不同」，
//     讓玩家需要快速判斷哪個目標最值得打
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	// RainbowSharkDurationSec 彩虹爆發持續時間（秒）
	RainbowSharkDurationSec = 10
	// RainbowSharkCooldownSec 全服冷卻時間（秒）
	RainbowSharkCooldownSec = 40
	// RainbowSharkAnnounceMinMult 全服公告最低倍率門檻（有 3x 目標才公告）
	RainbowSharkAnnounceMinMult = 3.0
)

// rainbowSharkMultWeights 彩虹爆發倍率加權（低倍率高機率）
var rainbowSharkMultWeights = []struct {
	Mult   float64
	Weight int
}{
	{1.5, 40}, // 40% 機率 1.5x
	{2.0, 30}, // 30% 機率 2.0x
	{2.5, 20}, // 20% 機率 2.5x
	{3.0, 10}, // 10% 機率 3.0x
}

// RainbowSharkMarkedTarget 彩虹爆發標記目標
type RainbowSharkMarkedTarget struct {
	InstanceID string  `json:"instance_id"`
	DefID      string  `json:"def_id"`
	X          float64 `json:"x"`
	Y          float64 `json:"y"`
	BurstMult  float64 `json:"burst_mult"` // 1.5/2.0/2.5/3.0
}

// rainbowSharkManager 彩虹鯊魚爆發管理器（全服共享）
type rainbowSharkManager struct {
	mu          sync.Mutex
	isActive    bool
	expiresAt   time.Time
	cooldownEnd time.Time
	// 標記目標的倍率（instanceID → burstMult）
	markedMults map[string]float64
}

// newRainbowSharkManager 建立彩虹鯊魚爆發管理器
func newRainbowSharkManager() *rainbowSharkManager {
	return &rainbowSharkManager{
		markedMults: make(map[string]float64),
	}
}

// isRainbowShark 判斷是否為彩虹鯊魚（T138）
func isRainbowShark(defID string) bool {
	return defID == "T138"
}

// isOnCooldown 檢查是否在全服冷卻中
func (m *rainbowSharkManager) isOnCooldown() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return time.Now().Before(m.cooldownEnd)
}

// activate 激活彩虹爆發，設定標記目標
func (m *rainbowSharkManager) activate(marks []RainbowSharkMarkedTarget) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.isActive = true
	m.expiresAt = time.Now().Add(time.Duration(RainbowSharkDurationSec) * time.Second)
	m.cooldownEnd = time.Now().Add(time.Duration(RainbowSharkCooldownSec) * time.Second)
	m.markedMults = make(map[string]float64, len(marks))
	for _, mk := range marks {
		m.markedMults[mk.InstanceID] = mk.BurstMult
	}
}

// getRainbowSharkMult 取得目標的彩虹爆發倍率（供 handleKill 使用）
// 回傳 1.0 表示無加成
func (m *rainbowSharkManager) getRainbowSharkMult(instanceID string) float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.isActive {
		return 1.0
	}
	if time.Now().After(m.expiresAt) {
		m.isActive = false
		m.markedMults = make(map[string]float64)
		return 1.0
	}
	mult, ok := m.markedMults[instanceID]
	if !ok {
		return 1.0
	}
	return mult
}

// removeMarked 移除已擊破的標記目標
func (m *rainbowSharkManager) removeMarked(instanceID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.markedMults, instanceID)
}

// deactivate 結束彩虹爆發
func (m *rainbowSharkManager) deactivate() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.isActive = false
	m.markedMults = make(map[string]float64)
}

// pickRainbowSharkMult 加權隨機選擇彩虹爆發倍率
func pickRainbowSharkMult() float64 {
	total := 0
	for _, w := range rainbowSharkMultWeights {
		total += w.Weight
	}
	r := rand.Intn(total)
	cumulative := 0
	for _, w := range rainbowSharkMultWeights {
		cumulative += w.Weight
		if r < cumulative {
			return w.Mult
		}
	}
	return 1.5
}

// tryRainbowSharkBurst 擊破 T138 後觸發彩虹爆發（DAY-180）
// 由 handleKill 呼叫（在 goroutine 中執行）
func (g *Game) tryRainbowSharkBurst(p *player.Player, triggerID string) {
	// 全服冷卻檢查
	if g.RainbowShark.isOnCooldown() {
		return
	}

	// 收集場上所有存活目標（排除 BOSS 和觸發者）
	g.mu.RLock()
	var marks []RainbowSharkMarkedTarget
	maxMult := 0.0
	for id, t := range g.Targets {
		if id == triggerID || t.HP <= 0 || t.DefID == "B001" {
			continue
		}
		mult := pickRainbowSharkMult()
		marks = append(marks, RainbowSharkMarkedTarget{
			InstanceID: t.InstanceID,
			DefID:      t.DefID,
			X:          t.X,
			Y:          t.Y,
			BurstMult:  mult,
		})
		if mult > maxMult {
			maxMult = mult
		}
	}
	g.mu.RUnlock()

	if len(marks) == 0 {
		return
	}

	// 激活彩虹爆發
	g.RainbowShark.activate(marks)

	log.Printf("[RainbowShark] player=%s triggered burst: %d targets marked, maxMult=%.1f",
		p.ID, len(marks), maxMult)

	// 廣播彩虹爆發開始（全服）
	// 轉換為 ws payload 格式
	wsMarks := make([]ws.RainbowSharkMarkedTargetPayload, len(marks))
	for i, mk := range marks {
		wsMarks[i] = ws.RainbowSharkMarkedTargetPayload{
			InstanceID: mk.InstanceID,
			DefID:      mk.DefID,
			X:          mk.X,
			Y:          mk.Y,
			BurstMult:  mk.BurstMult,
		}
	}
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgRainbowSharkBurst,
		Payload: ws.RainbowSharkBurstPayload{
			Phase:           "burst_start",
			TriggerPlayerID: p.ID,
			TriggerName:     p.DisplayName,
			MarkedTargets:   wsMarks,
			DurationSec:     RainbowSharkDurationSec,
		},
	})

	// 全服公告（有 3x 目標才公告）
	if maxMult >= RainbowSharkAnnounceMinMult {
		g.announceRainbowSharkBurst(p.DisplayName, len(marks), maxMult)
	}

	// 等待爆發結束
	time.Sleep(time.Duration(RainbowSharkDurationSec) * time.Second)

	// 結束彩虹爆發
	g.RainbowShark.deactivate()

	// 廣播彩虹爆發結束（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgRainbowSharkBurst,
		Payload: ws.RainbowSharkBurstPayload{
			Phase: "burst_end",
		},
	})

	log.Printf("[RainbowShark] burst ended for player=%s", p.ID)
}

// getRainbowSharkMult 取得彩虹鯊魚爆發倍率（供 handleKill 使用）
func (g *Game) getRainbowSharkMult(instanceID string) float64 {
	return g.RainbowShark.getRainbowSharkMult(instanceID)
}

// removeRainbowSharkMark 移除已擊破的彩虹標記（供 handleKill 使用）
func (g *Game) removeRainbowSharkMark(instanceID string) {
	g.RainbowShark.removeMarked(instanceID)
}

// announceRainbowSharkBurst 全服公告彩虹鯊魚爆發（DAY-180）
func (g *Game) announceRainbowSharkBurst(playerName string, targetCount int, maxMult float64) {
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAnnouncement,
		Payload: map[string]interface{}{
			"event_type": "rainbow_shark_burst",
			"message":    fmt.Sprintf("🌈 %s 擊破彩虹鯊魚！%d 個目標獲得彩虹加成！最高 %.0fx！", playerName, targetCount, maxMult),
			"color":      "#FF69B4", // 熱粉紅（彩虹感）
			"duration":   5.0,
			"priority":   3,
		},
	})
}
