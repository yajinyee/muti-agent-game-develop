// lucky_freeze_world_handler.go — 幸運冰凍世界魚系統（DAY-237）
// 業界原創「全場冰凍+冰裂爆發」機制
//
// 設計：擊破 T195 後觸發「冰凍世界」（8 秒）：
//   - 場上所有目標物移動速度降低 80%（幾乎靜止）
//   - 冰凍世界期間擊破任何目標獲得 ×2.0 倍率加成（乘法）
//   - 8 秒後「冰裂爆發」：所有目標 HP -50%（冰裂傷害），同時移動速度恢復
//   - 個人冷卻 20 秒；全服冷卻 30 秒
//
// 設計差異：
//   - 與時間凍結魚（DAY-212，全場完全靜止）不同，冰凍世界是「速度降低 80%」，
//     目標仍在緩慢移動，讓玩家有「要趁目標慢的時候趕快打」的緊迫感
//   - 「冰裂爆發 HP -50%」比鏡面崩潰（-35%）更強，讓玩家有更大的爆發感
//   - ×2.0 倍率加成（全場有效）讓玩家有「趁冰凍期間多打」的動機
//   - 全服廣播冰凍狀態讓所有玩家都知道目標變慢了，製造「全服一起趁機打」的社交感
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
	LuckyFreezeWorldPersonalCD  = 20 * time.Second // 個人冷卻
	LuckyFreezeWorldGlobalCD    = 30 * time.Second // 全服冷卻
	LuckyFreezeWorldDuration    = 8 * time.Second  // 冰凍世界持續時間
	LuckyFreezeWorldKillBoost   = 2.0              // 冰凍世界期間擊破倍率加成（乘法）
	LuckyFreezeWorldSpeedFactor = 0.20             // 冰凍期間速度係數（降低 80%）
	LuckyFreezeWorldCrackHP     = 0.50             // 冰裂爆發 HP 削減比例
)

// luckyFreezeWorldManager 幸運冰凍世界魚管理器
type luckyFreezeWorldManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 冰凍世界狀態
	active      bool
	activeUntil time.Time
	instanceID  string
}

func newLuckyFreezeWorldManager() *luckyFreezeWorldManager {
	return &luckyFreezeWorldManager{
		personalCooldowns: make(map[string]time.Time),
	}
}

// isLuckyFreezeWorldFish 判斷是否為幸運冰凍世界魚
func isLuckyFreezeWorldFish(defID string) bool {
	return defID == "T195"
}

// isFreezeWorldActive 判斷冰凍世界是否啟動中（供 target 移動邏輯使用）
func (g *Game) isFreezeWorldActive() bool {
	mgr := g.LuckyFreezeWorld
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

// getLuckyFreezeWorldBoost 取得冰凍世界倍率加成（供 handleKill 使用）
func (g *Game) getLuckyFreezeWorldBoost() float64 {
	if g.isFreezeWorldActive() {
		return LuckyFreezeWorldKillBoost
	}
	return 1.0
}

// tryLuckyFreezeWorldFish 擊破 T195 後觸發冰凍世界（供 handleKill 使用）
func (g *Game) tryLuckyFreezeWorldFish(p *player.Player) {
	mgr := g.LuckyFreezeWorld
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
	// 已有冰凍世界啟動中
	if mgr.active && time.Now().Before(mgr.activeUntil) {
		mgr.mu.Unlock()
		return
	}

	// 設定冷卻
	mgr.personalCooldowns[p.ID] = time.Now().Add(LuckyFreezeWorldPersonalCD)
	mgr.globalCooldownUntil = time.Now().Add(LuckyFreezeWorldGlobalCD)

	// 啟動冰凍世界
	mgr.active = true
	mgr.activeUntil = time.Now().Add(LuckyFreezeWorldDuration)
	instanceID := fmt.Sprintf("freeze_%d", time.Now().UnixNano())
	mgr.instanceID = instanceID
	mgr.mu.Unlock()

	log.Printf("[LuckyFreezeWorld] player=%s activated freeze world for %v", p.ID, LuckyFreezeWorldDuration)

	// 套用冰凍速度降低（所有目標速度 × 0.20）
	frozenCount := g.applyFreezeWorldSpeed(LuckyFreezeWorldSpeedFactor)

	// 全服廣播：冰凍世界開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyFreezeWorld,
		Payload: ws.LuckyFreezeWorldPayload{
			Event:        "freeze_start",
			PlayerID:     p.ID,
			PlayerName:   p.DisplayName,
			DurationSec:  int(LuckyFreezeWorldDuration.Seconds()),
			KillBoost:    LuckyFreezeWorldKillBoost,
			SpeedFactor:  LuckyFreezeWorldSpeedFactor,
			FrozenCount:  frozenCount,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyFreezeWorld, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("❄️ %s 觸發冰凍世界！%d 個目標速度降低 80%%，×%.1f 倍率加成！",
			p.DisplayName, frozenCount, LuckyFreezeWorldKillBoost),
		"color": "#5DADE2",
	})
	g.broadcastAnnouncement(ann)

	// 啟動冰凍世界計時器
	go g.runLuckyFreezeWorld(p, instanceID)
}

// applyFreezeWorldSpeed 套用冰凍速度（廣播給 Client 端降速，Server 端記錄狀態）
func (g *Game) applyFreezeWorldSpeed(factor float64) int {
	g.mu.RLock()
	defer g.mu.RUnlock()

	count := 0
	for _, t := range g.Targets {
		if t.HP <= 0 {
			continue
		}
		count++
	}
	return count
}

// restoreFreezeWorldSpeed 恢復冰凍速度（廣播給 Client 端恢復速度）
func (g *Game) restoreFreezeWorldSpeed(factor float64) {
	// 速度恢復由 freeze_end 廣播通知 Client 端處理
	// Server 端不直接修改 Target 速度（速度由 Def.Speed 決定）
}

// runLuckyFreezeWorld 冰凍世界計時器（goroutine）
func (g *Game) runLuckyFreezeWorld(p *player.Player, instanceID string) {
	timer := time.NewTimer(LuckyFreezeWorldDuration)
	defer timer.Stop()

	<-timer.C

	// 確認 instanceID 仍有效
	g.LuckyFreezeWorld.mu.Lock()
	if g.LuckyFreezeWorld.instanceID != instanceID {
		g.LuckyFreezeWorld.mu.Unlock()
		return
	}
	g.LuckyFreezeWorld.active = false
	g.LuckyFreezeWorld.mu.Unlock()

	log.Printf("[LuckyFreezeWorld] freeze world ended, triggering ice crack")

	// 恢復速度
	g.restoreFreezeWorldSpeed(LuckyFreezeWorldSpeedFactor)

	// 觸發冰裂爆發
	g.doFreezeWorldCrack(p, instanceID)
}

// doFreezeWorldCrack 冰裂爆發（所有目標 HP -50%）
func (g *Game) doFreezeWorldCrack(p *player.Player, instanceID string) {
	g.mu.Lock()

	crackedCount := 0
	for _, t := range g.Targets {
		if t.HP <= 0 {
			continue
		}
		// HP -50%，保留最少 1
		damage := int(float64(t.HP) * LuckyFreezeWorldCrackHP)
		if damage < 1 {
			damage = 1
		}
		t.HP -= damage
		if t.HP < 1 {
			t.HP = 1
		}
		crackedCount++
	}
	g.mu.Unlock()

	log.Printf("[LuckyFreezeWorld] ice crack: affected %d targets", crackedCount)

	// 廣播冰裂爆發
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyFreezeWorld,
		Payload: ws.LuckyFreezeWorldPayload{
			Event:        "freeze_crack",
			PlayerID:     p.ID,
			CrackedCount: crackedCount,
		},
	})

	// 廣播冰凍世界結束
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyFreezeWorld,
		Payload: ws.LuckyFreezeWorldPayload{
			Event: "freeze_end",
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyFreezeWorld, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("❄️ 冰裂爆發！%d 個目標 HP -50%%！速度恢復！",
			crackedCount),
		"color": "#2E86C1",
	})
	g.broadcastAnnouncement(ann)
}
