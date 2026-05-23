// lucky_gravity_flip_handler.go — 幸運重力反轉魚系統（DAY-238）
// 業界原創「重力反轉+上下顛倒移動+重力崩潰」機制
//
// 設計：擊破 T196 後觸發「重力反轉」（10 秒）：
//   - 場上所有目標物的 Y 座標以場景中央（Y=300）為軸翻轉
//   - 重力反轉期間擊破任何目標獲得 ×2.1 倍率加成（乘法）
//   - 10 秒後「重力崩潰」：所有目標 HP -45%（保留最少 1），Y 座標恢復
//   - 個人冷卻 22 秒；全服冷卻 32 秒
//
// 設計差異：
//   - 與鏡面世界魚（DAY-236，X 座標翻轉）不同，重力反轉是「Y 座標翻轉」，
//     讓玩家有「要重新瞄準上下顛倒的位置」的空間感
//   - 「上下顛倒」讓玩家感受到「世界倒過來了」的視覺衝擊，比左右鏡像更有衝擊感
//   - ×2.1 倍率加成（全場有效）讓玩家有「趕快在 10 秒內多打」的緊迫感
//   - 「重力崩潰 HP -45%」比鏡面崩潰（-35%）更強，讓玩家有更大的爆發感
//   - 全服廣播翻轉位置讓所有玩家都看到目標的新位置，製造「全服一起重新瞄準」的社交感
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	LuckyGravityFlipPersonalCD = 22 * time.Second // 個人冷卻
	LuckyGravityFlipGlobalCD   = 32 * time.Second // 全服冷卻
	LuckyGravityFlipDuration   = 10 * time.Second // 重力反轉持續時間
	LuckyGravityFlipKillBoost  = 2.1              // 重力反轉期間擊破倍率加成（乘法）
	LuckyGravityFlipCollapseHP = 0.45             // 重力崩潰 HP 削減比例
	LuckyGravityFlipCenterY    = 300.0            // 翻轉軸 Y（場景中央）
)

// luckyGravityFlipManager 幸運重力反轉魚管理器
type luckyGravityFlipManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 重力反轉狀態
	active      bool
	activeUntil time.Time
	instanceID  string
}

func newLuckyGravityFlipManager() *luckyGravityFlipManager {
	return &luckyGravityFlipManager{
		personalCooldowns: make(map[string]time.Time),
	}
}

// isLuckyGravityFlipFish 判斷是否為幸運重力反轉魚
func isLuckyGravityFlipFish(defID string) bool {
	return defID == "T196"
}

// isGravityFlipActive 判斷重力反轉是否啟動中（供 handleKill 使用）
func (g *Game) isGravityFlipActive() bool {
	mgr := g.LuckyGravityFlip
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	if !mgr.active {
		return false
	}
	if time.Now().After(mgr.activeUntil) {
		mgr.active = false
		return false
	}
	return true
}

// getLuckyGravityFlipBoost 取得重力反轉倍率加成（供 handleKill 使用）
func (g *Game) getLuckyGravityFlipBoost() float64 {
	if g.isGravityFlipActive() {
		return LuckyGravityFlipKillBoost
	}
	return 1.0
}

// tryLuckyGravityFlipFish 擊破 T196 後觸發重力反轉（供 handleKill 使用）
func (g *Game) tryLuckyGravityFlipFish(p *player.Player) {
	mgr := g.LuckyGravityFlip
	mgr.mu.Lock()

	// 個人冷卻檢查
	if cd, ok := mgr.personalCooldowns[p.ID]; ok && time.Now().Before(cd) {
		mgr.mu.Unlock()
		return
	}
	// 全服冷卻檢查
	if time.Now().Before(mgr.globalCooldownUntil) {
		mgr.mu.Unlock()
		return
	}
	// 已有重力反轉啟動中
	if mgr.active && time.Now().Before(mgr.activeUntil) {
		mgr.mu.Unlock()
		return
	}

	// 設定冷卻
	mgr.personalCooldowns[p.ID] = time.Now().Add(LuckyGravityFlipPersonalCD)
	mgr.globalCooldownUntil = time.Now().Add(LuckyGravityFlipGlobalCD)

	// 啟動重力反轉
	mgr.active = true
	mgr.activeUntil = time.Now().Add(LuckyGravityFlipDuration)
	instanceID := fmt.Sprintf("gravity_%d", time.Now().UnixNano())
	mgr.instanceID = instanceID
	mgr.mu.Unlock()

	log.Printf("[LuckyGravityFlip] player=%s activated gravity flip for %v", p.ID, LuckyGravityFlipDuration)

	// 執行 Y 座標翻轉，取得所有目標的新位置
	flippedPositions := g.doGravityFlip(instanceID)

	// 全服廣播：重力反轉開始（含所有目標的翻轉後位置）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyGravityFlip,
		Payload: ws.LuckyGravityFlipPayload{
			Event:       "gravity_start",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			DurationSec: int(LuckyGravityFlipDuration.Seconds()),
			KillBoost:   LuckyGravityFlipKillBoost,
			Positions:   flippedPositions,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyGravityFlip, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🔄 %s 觸發重力反轉！目標上下翻轉，×%.1f 倍率加成！",
			p.DisplayName, LuckyGravityFlipKillBoost),
		"color": "#E67E22",
	})
	g.broadcastAnnouncement(ann)

	// 啟動重力反轉計時器
	go g.runLuckyGravityFlip(p, instanceID)
}

// doGravityFlip 執行重力翻轉（所有目標 Y 座標以 CenterY 為軸翻轉）
func (g *Game) doGravityFlip(instanceID string) interface{} {
	g.mu.Lock()
	defer g.mu.Unlock()

	type flippedPos struct {
		ID string  `json:"id"`
		X  float64 `json:"x"`
		Y  float64 `json:"y"`
	}

	var positions []flippedPos

	for id, t := range g.Targets {
		if t.HP <= 0 {
			continue
		}
		// 以 CenterY 為軸翻轉 Y 座標
		newY := 2*LuckyGravityFlipCenterY - t.Y
		// 邊界限制（場景 Y 範圍 60-540）
		if newY < 60 {
			newY = 60
		}
		if newY > 540 {
			newY = 540
		}
		t.Y = newY
		positions = append(positions, flippedPos{ID: id, X: t.X, Y: t.Y})
	}

	log.Printf("[LuckyGravityFlip] flipped %d targets", len(positions))
	return positions
}

// runLuckyGravityFlip 重力反轉計時器（goroutine）
func (g *Game) runLuckyGravityFlip(p *player.Player, instanceID string) {
	timer := time.NewTimer(LuckyGravityFlipDuration)
	defer timer.Stop()

	<-timer.C

	// 確認 instanceID 仍有效
	g.LuckyGravityFlip.mu.Lock()
	if g.LuckyGravityFlip.instanceID != instanceID {
		g.LuckyGravityFlip.mu.Unlock()
		return
	}
	g.LuckyGravityFlip.active = false
	g.LuckyGravityFlip.mu.Unlock()

	log.Printf("[LuckyGravityFlip] gravity flip ended, triggering collapse")
	g.doGravityCollapse(p, instanceID)
}

// doGravityCollapse 重力崩潰（所有目標 HP -45%）
func (g *Game) doGravityCollapse(p *player.Player, instanceID string) {
	g.mu.Lock()

	collapsedCount := 0
	for _, t := range g.Targets {
		if t.HP <= 0 {
			continue
		}
		// HP -45%，保留最少 1
		damage := int(float64(t.HP) * LuckyGravityFlipCollapseHP)
		if damage < 1 {
			damage = 1
		}
		t.HP -= damage
		if t.HP < 1 {
			t.HP = 1
		}
		collapsedCount++
	}
	g.mu.Unlock()

	log.Printf("[LuckyGravityFlip] collapse: affected %d targets", collapsedCount)

	// 廣播重力崩潰
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyGravityFlip,
		Payload: ws.LuckyGravityFlipPayload{
			Event:          "gravity_collapse",
			PlayerID:       p.ID,
			CollapsedCount: collapsedCount,
		},
	})

	// 廣播重力反轉結束
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyGravityFlip,
		Payload: ws.LuckyGravityFlipPayload{
			Event: "gravity_end",
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyGravityFlip, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🔄 重力崩潰！%d 個目標 HP -45%%！重力恢復！",
			collapsedCount),
		"color": "#D35400",
	})
	g.broadcastAnnouncement(ann)
}
