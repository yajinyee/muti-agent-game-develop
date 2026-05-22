// icefishing_handler.go — 冰釣幸運輪盤系統（DAY-171）
// 業界依據：Cozy Fishing Life（2026-05-10）「Winter Wheel — 8 segments x2-x10 multipliers
// + bonus mode triggers」+ Ice Fishing Live（Evolution Gaming）「wheel triggers bonus fishing rounds」
// 設計：擊破 T129 冰釣魚後觸發「冰釣幸運輪盤」（8格：2x-10x 倍率加成）
// 玩家在 5 秒內點擊停止，結果預先決定（公平性保證）
// 觸發後玩家在 8 秒內所有擊破獎勵套用輪盤倍率
// 設計差異：與巨型章魚輪盤（950x 大獎）不同，冰釣輪盤是「倍率加成型」（2x-10x）
// 讓玩家在短時間內所有擊破都有倍率加成，製造「黃金 8 秒」的爽感
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// iceFishingWheelSlot 冰釣輪盤格子定義
type iceFishingWheelSlot struct {
	Multiplier float64 // 倍率加成（2x-10x）
	Label      string  // 顯示文字
	Color      string  // 格子顏色
	Weight     int     // 抽取權重
}

// iceFishingWheelSlots 冰釣輪盤格子（8格）
var iceFishingWheelSlots = []iceFishingWheelSlot{
	{Multiplier: 2.0, Label: "×2", Color: "#64B5F6", Weight: 30},  // 淺藍（最常見）
	{Multiplier: 3.0, Label: "×3", Color: "#42A5F5", Weight: 25},  // 藍
	{Multiplier: 4.0, Label: "×4", Color: "#2196F3", Weight: 18},  // 中藍
	{Multiplier: 5.0, Label: "×5", Color: "#1E88E5", Weight: 12},  // 深藍
	{Multiplier: 6.0, Label: "×6", Color: "#1565C0", Weight: 8},   // 更深藍
	{Multiplier: 7.0, Label: "×7", Color: "#0D47A1", Weight: 4},   // 深海藍
	{Multiplier: 8.0, Label: "×8", Color: "#00BCD4", Weight: 2},   // 青色（稀有）
	{Multiplier: 10.0, Label: "×10", Color: "#00E5FF", Weight: 1}, // 冰藍（最稀有）
}

// iceFishingSession 冰釣輪盤 session（per-player）
type iceFishingSession struct {
	wheelResult  int       // 預先決定的輪盤結果（格子索引）
	multiplier   float64   // 輪盤倍率
	activatedAt  time.Time // 倍率激活時間
	expiresAt    time.Time // 倍率過期時間
	killCount    int       // 倍率期間擊破數
	totalBonus   int       // 倍率期間額外獎勵
	isWheelSpun  bool      // 是否已停止輪盤
}

// iceFishingManager 冰釣幸運輪盤管理器
type iceFishingManager struct {
	mu       sync.Mutex
	sessions map[string]*iceFishingSession // playerID -> session
	cooldowns map[string]time.Time          // playerID -> 冷卻結束時間
	CooldownSec int
	MultDurationSec int
}

func newIceFishingManager() *iceFishingManager {
	return &iceFishingManager{
		sessions:        make(map[string]*iceFishingSession),
		cooldowns:       make(map[string]time.Time),
		CooldownSec:     45, // 45 秒個人冷卻
		MultDurationSec: 8,  // 倍率持續 8 秒
	}
}

// isIceFish 判斷是否為冰釣魚
func isIceFish(defID string) bool {
	return defID == "T129"
}

// getIceFishingMult 取得冰釣輪盤倍率（供 handleKill 使用）
func (g *Game) getIceFishingMult(playerID string) float64 {
	g.IceFishing.mu.Lock()
	defer g.IceFishing.mu.Unlock()

	sess, ok := g.IceFishing.sessions[playerID]
	if !ok || !sess.isWheelSpun {
		return 1.0
	}
	if time.Now().After(sess.expiresAt) {
		delete(g.IceFishing.sessions, playerID)
		return 1.0
	}
	return sess.multiplier
}

// recordIceFishingKill 記錄冰釣倍率期間的擊破
func (g *Game) recordIceFishingKill(playerID string, bonusReward int) {
	g.IceFishing.mu.Lock()
	defer g.IceFishing.mu.Unlock()

	sess, ok := g.IceFishing.sessions[playerID]
	if !ok || !sess.isWheelSpun {
		return
	}
	if time.Now().After(sess.expiresAt) {
		return
	}
	sess.killCount++
	sess.totalBonus += bonusReward
}

// tryIceFishingWheel 冰釣魚擊破後觸發輪盤
func (g *Game) tryIceFishingWheel(p *player.Player, instanceID string, fx, fy float64) {
	g.IceFishing.mu.Lock()
	// 檢查冷卻
	if cd, ok := g.IceFishing.cooldowns[p.ID]; ok && time.Now().Before(cd) {
		g.IceFishing.mu.Unlock()
		return
	}
	// 預先決定輪盤結果（公平性保證）
	totalWeight := 0
	for _, slot := range iceFishingWheelSlots {
		totalWeight += slot.Weight
	}
	r := rand.Intn(totalWeight)
	resultIdx := 0
	cumWeight := 0
	for i, slot := range iceFishingWheelSlots {
		cumWeight += slot.Weight
		if r < cumWeight {
			resultIdx = i
			break
		}
	}
	resultSlot := iceFishingWheelSlots[resultIdx]

	// 建立 session
	sess := &iceFishingSession{
		wheelResult: resultIdx,
		multiplier:  resultSlot.Multiplier,
		isWheelSpun: false,
	}
	g.IceFishing.sessions[p.ID] = sess
	g.IceFishing.cooldowns[p.ID] = time.Now().Add(time.Duration(g.IceFishing.CooldownSec) * time.Second)
	g.IceFishing.mu.Unlock()

	// 廣播輪盤開始（個人）
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgIceFishingWheel,
		Payload: ws.IceFishingWheelPayload{
			Phase:       "wheel_start",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			WheelResult: resultIdx,
			Multiplier:  resultSlot.Multiplier,
			Label:       resultSlot.Label,
			Color:       resultSlot.Color,
			SpinSec:     5, // 5 秒旋轉時間
		},
	})

	// 全服廣播（讓其他玩家看到）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgIceFishingWheel,
		Payload: ws.IceFishingWheelPayload{
			Phase:      "wheel_broadcast",
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
		},
	})

	log.Printf("[IceFishing] player=%s triggered wheel, result=×%.0f (slot=%d)",
		p.ID, resultSlot.Multiplier, resultIdx)

	// 5 秒後自動停止（如果玩家沒有手動停止）
	go func() {
		time.Sleep(5 * time.Second)
		g.handleIceFishingWheelStop(p)
	}()
}

// handleIceFishingWheelStop 處理玩家停止輪盤（或自動停止）
func (g *Game) handleIceFishingWheelStop(p *player.Player) {
	g.IceFishing.mu.Lock()
	sess, ok := g.IceFishing.sessions[p.ID]
	if !ok || sess.isWheelSpun {
		g.IceFishing.mu.Unlock()
		return
	}
	sess.isWheelSpun = true
	sess.activatedAt = time.Now()
	sess.expiresAt = time.Now().Add(time.Duration(g.IceFishing.MultDurationSec) * time.Second)
	resultIdx := sess.wheelResult
	multiplier := sess.multiplier
	g.IceFishing.mu.Unlock()

	resultSlot := iceFishingWheelSlots[resultIdx]

	// 廣播輪盤結果（個人）
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgIceFishingWheel,
		Payload: ws.IceFishingWheelPayload{
			Phase:       "wheel_result",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			WheelResult: resultIdx,
			Multiplier:  multiplier,
			Label:       resultSlot.Label,
			Color:       resultSlot.Color,
			DurationSec: g.IceFishing.MultDurationSec,
		},
	})

	// 全服廣播結果（≥5x 才廣播）
	if multiplier >= 5.0 {
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgIceFishingWheel,
			Payload: ws.IceFishingWheelPayload{
				Phase:      "wheel_result_broadcast",
				PlayerID:   p.ID,
				PlayerName: p.DisplayName,
				Multiplier: multiplier,
				Label:      resultSlot.Label,
			},
		})
	}

	// 全服公告（≥7x 才公告）
	if multiplier >= 7.0 {
		g.announceIceFishing(p.DisplayName, multiplier)
	}

	// 等待倍率結束，廣播結束通知
	go func() {
		time.Sleep(time.Duration(g.IceFishing.MultDurationSec) * time.Second)

		g.IceFishing.mu.Lock()
		sess2, ok2 := g.IceFishing.sessions[p.ID]
		var killCount, totalBonus int
		if ok2 {
			killCount = sess2.killCount
			totalBonus = sess2.totalBonus
			delete(g.IceFishing.sessions, p.ID)
		}
		g.IceFishing.mu.Unlock()

		// 廣播倍率結束（個人）
		g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgIceFishingWheel,
			Payload: ws.IceFishingWheelPayload{
				Phase:       "mult_end",
				PlayerID:    p.ID,
				PlayerName:  p.DisplayName,
				Multiplier:  multiplier,
				KillCount:   killCount,
				TotalBonus:  totalBonus,
				DurationSec: g.IceFishing.MultDurationSec,
			},
		})

		// 全服公告：≥7x 且擊破 ≥3 個時廣播結果
		if multiplier >= 7.0 && killCount >= 3 {
			ann := g.Announce.Create(announce.EventIceFishingResult, p.DisplayName, totalBonus, map[string]string{
				"mult":   fmt.Sprintf("%.0f", multiplier),
				"kills":  fmt.Sprintf("%d", killCount),
				"reward": fmt.Sprintf("%d", totalBonus),
			})
			g.broadcastAnnouncement(ann)
		}

		log.Printf("[IceFishing] player=%s mult=×%.0f ended: kills=%d bonus=%d",
			p.ID, multiplier, killCount, totalBonus)
	}()
}

// announceIceFishing 全服公告冰釣幸運輪盤高倍率
func (g *Game) announceIceFishing(playerName string, multiplier float64) {
	ann := g.Announce.Create(announce.EventIceFishing, playerName, int(multiplier), map[string]string{
		"mult": fmt.Sprintf("%.0f", multiplier),
	})
	g.broadcastAnnouncement(ann)
}
