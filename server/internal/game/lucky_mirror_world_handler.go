// lucky_mirror_world_handler.go — 幸運鏡面世界魚系統（DAY-236）
// 業界原創「全場鏡像反轉+鏡面崩潰」機制
//
// 設計：擊破 T194 後觸發「鏡面世界」（10 秒）：
//   - 場上所有目標物的 X 座標鏡像反轉（以場景中央 X=500 為軸）
//   - 鏡面世界期間擊破任何目標獲得 ×2.3 倍率加成（乘法）
//   - 10 秒後「鏡面崩潰」：所有目標 HP -35%（保留最少 1）
//   - 個人冷卻 22 秒；全服冷卻 35 秒
//
// 設計差異：
//   - 與傳送魚（DAY-223，全場瞬間移動到隨機位置）不同，鏡面魚是「對稱翻轉」，
//     讓玩家有「要重新瞄準鏡像位置」的空間感
//   - 「鏡像反轉」讓玩家感受到「世界顛倒了」的視覺衝擊
//   - ×2.3 倍率加成（全場有效）讓玩家有「趕快在 10 秒內多打」的緊迫感
//   - 「鏡面崩潰 HP -35%」讓玩家有「等待→爆發」的高潮設計
//   - 全服廣播鏡像位置讓所有玩家都看到目標的新位置，製造「全服一起重新瞄準」的社交感
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
	LuckyMirrorWorldPersonalCD  = 22 * time.Second // 個人冷卻
	LuckyMirrorWorldGlobalCD    = 35 * time.Second // 全服冷卻
	LuckyMirrorWorldDuration    = 10 * time.Second // 鏡面世界持續時間
	LuckyMirrorWorldKillBoost   = 2.3              // 鏡面世界期間擊破倍率加成（乘法）
	LuckyMirrorWorldCollapseHP  = 0.35             // 鏡面崩潰 HP 削減比例
	LuckyMirrorWorldCenterX     = 500.0            // 鏡像軸 X（場景中央）
)

// luckyMirrorWorldManager 幸運鏡面世界魚管理器
type luckyMirrorWorldManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 鏡面世界狀態
	active      bool
	activeUntil time.Time
	instanceID  string
}

func newLuckyMirrorWorldManager() *luckyMirrorWorldManager {
	return &luckyMirrorWorldManager{
		personalCooldowns: make(map[string]time.Time),
	}
}

// isLuckyMirrorWorldFish 判斷是否為幸運鏡面世界魚
func isLuckyMirrorWorldFish(defID string) bool {
	return defID == "T194"
}

// isMirrorWorldActive 判斷鏡面世界是否啟動中（供 handleKill 使用）
func (g *Game) isMirrorWorldActive() bool {
	mgr := g.LuckyMirrorWorld
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

// getLuckyMirrorWorldBoost 取得鏡面世界倍率加成（供 handleKill 使用）
func (g *Game) getLuckyMirrorWorldBoost() float64 {
	if g.isMirrorWorldActive() {
		return LuckyMirrorWorldKillBoost
	}
	return 1.0
}

// tryLuckyMirrorWorldFish 擊破 T194 後觸發鏡面世界（供 handleKill 使用）
func (g *Game) tryLuckyMirrorWorldFish(p *player.Player) {
	mgr := g.LuckyMirrorWorld
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
	// 已有鏡面世界啟動中
	if mgr.active && time.Now().Before(mgr.activeUntil) {
		mgr.mu.Unlock()
		return
	}

	// 設定冷卻
	mgr.personalCooldowns[p.ID] = time.Now().Add(LuckyMirrorWorldPersonalCD)
	mgr.globalCooldownUntil = time.Now().Add(LuckyMirrorWorldGlobalCD)

	// 啟動鏡面世界
	mgr.active = true
	mgr.activeUntil = time.Now().Add(LuckyMirrorWorldDuration)
	instanceID := fmt.Sprintf("mirror_%d", time.Now().UnixNano())
	mgr.instanceID = instanceID
	mgr.mu.Unlock()

	log.Printf("[LuckyMirrorWorld] player=%s activated mirror world for %v", p.ID, LuckyMirrorWorldDuration)

	// 執行鏡像反轉，取得所有目標的新位置
	mirroredPositions := g.doMirrorFlip(instanceID)

	// 全服廣播：鏡面世界開始（含所有目標的鏡像位置）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyMirrorWorld,
		Payload: ws.LuckyMirrorWorldPayload{
			Event:       "mirror_start",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			DurationSec: int(LuckyMirrorWorldDuration.Seconds()),
			KillBoost:   LuckyMirrorWorldKillBoost,
			Positions:   mirroredPositions,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyMirrorWorld, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🪞 %s 觸發鏡面世界！目標左右翻轉，×%.1f 倍率加成！",
			p.DisplayName, LuckyMirrorWorldKillBoost),
		"color": "#8E44AD",
	})
	g.broadcastAnnouncement(ann)

	// 啟動鏡面世界計時器
	go g.runLuckyMirrorWorld(p, instanceID)
}

// doMirrorFlip 執行鏡像反轉（所有目標 X 座標以 CenterX 為軸翻轉）
func (g *Game) doMirrorFlip(instanceID string) interface{} {
	g.mu.Lock()
	defer g.mu.Unlock()

	type mirroredPos struct {
		ID string  `json:"id"`
		X  float64 `json:"x"`
		Y  float64 `json:"y"`
	}

	var positions []mirroredPos

	for id, t := range g.Targets {
		if t.HP <= 0 {
			continue
		}
		// 以 CenterX 為軸鏡像反轉 X 座標
		newX := 2*LuckyMirrorWorldCenterX - t.X
		// 邊界限制
		if newX < 50 {
			newX = 50
		}
		if newX > 950 {
			newX = 950
		}
		t.X = newX
		positions = append(positions, mirroredPos{ID: id, X: t.X, Y: t.Y})
	}

	log.Printf("[LuckyMirrorWorld] mirrored %d targets", len(positions))
	return positions
}

// runLuckyMirrorWorld 鏡面世界計時器（goroutine）
func (g *Game) runLuckyMirrorWorld(p *player.Player, instanceID string) {
	timer := time.NewTimer(LuckyMirrorWorldDuration)
	defer timer.Stop()

	<-timer.C

	// 確認 instanceID 仍有效
	g.LuckyMirrorWorld.mu.Lock()
	if g.LuckyMirrorWorld.instanceID != instanceID {
		g.LuckyMirrorWorld.mu.Unlock()
		return
	}
	g.LuckyMirrorWorld.active = false
	g.LuckyMirrorWorld.mu.Unlock()

	log.Printf("[LuckyMirrorWorld] mirror world ended, triggering collapse")
	g.doMirrorCollapse(p, instanceID)
}

// doMirrorCollapse 鏡面崩潰（所有目標 HP -35%）
func (g *Game) doMirrorCollapse(p *player.Player, instanceID string) {
	g.mu.Lock()

	collapsedCount := 0
	for _, t := range g.Targets {
		if t.HP <= 0 {
			continue
		}
		// HP -35%，保留最少 1
		damage := int(float64(t.HP) * LuckyMirrorWorldCollapseHP)
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

	log.Printf("[LuckyMirrorWorld] collapse: affected %d targets", collapsedCount)

	// 廣播鏡面崩潰
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyMirrorWorld,
		Payload: ws.LuckyMirrorWorldPayload{
			Event:          "mirror_collapse",
			PlayerID:       p.ID,
			CollapsedCount: collapsedCount,
		},
	})

	// 廣播鏡面世界結束
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyMirrorWorld,
		Payload: ws.LuckyMirrorWorldPayload{
			Event: "mirror_end",
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyMirrorWorld, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🪞 鏡面崩潰！%d 個目標 HP -35%%！",
			collapsedCount),
		"color": "#6C3483",
	})
	g.broadcastAnnouncement(ann)
}
