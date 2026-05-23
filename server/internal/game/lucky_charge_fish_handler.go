// lucky_charge_fish_handler.go — 幸運充能魚系統（DAY-225）
// 業界原創「射擊充能→爆發」機制
//
// 設計：擊破 T183 後觸發「充能模式」（12 秒）：
//   - 玩家的每次射擊都累積「充能值」（+1/shot）
//   - 充能值達到 10 → 自動觸發「充能爆發」：
//     下一次擊破任何目標獲得 ×5.0 倍率加成（一次性）
//   - 充能爆發後重置，可再次累積（12 秒內可觸發多次）
//   - 個人冷卻 22 秒；個人機制（不是全服）
//
// 設計差異：
//   - 與幸運共鳴魚（DAY-222，全服合力射擊，貢獻比例分配）不同，
//     充能魚是「個人射擊累積→個人爆發」，讓玩家有「我自己充能，我自己爆發」的個人英雄感
//   - ×5.0 一次性爆發是目前最高的個人倍率，讓玩家有「要選最高價值目標觸發爆發」的策略感
//   - 12 秒內可多次觸發，讓積極射擊的玩家有更多爆發機會
//   - 充能進度條讓玩家清楚看到「還差幾槍就爆發」的期待感
//   - 個人機制讓每個玩家都有自己的充能節奏，不受其他玩家影響
package game

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	LuckyChargePersonalCD    = 22 * time.Second // 個人冷卻
	LuckyChargeDuration      = 12 * time.Second // 充能模式持續時間
	LuckyChargeTarget        = 10               // 充能目標（射擊次數）
	LuckyChargeBurstMult     = 5.0              // 充能爆發倍率加成（一次性）
)

// chargeSession 充能 session（每個玩家獨立）
type chargeSession struct {
	instanceID string
	count      int64 // atomic
	burstReady bool  // 充能爆發就緒（下一次擊破觸發）
	activeUntil time.Time
}

// luckyChargeFishManager 幸運充能魚管理器
type luckyChargeFishManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldown map[string]time.Time

	// 充能 session（playerID → session）
	sessions map[string]*chargeSession
}

func newLuckyChargeFishManager() *luckyChargeFishManager {
	return &luckyChargeFishManager{
		personalCooldown: make(map[string]time.Time),
		sessions:         make(map[string]*chargeSession),
	}
}

// isLuckyChargeFish 判斷是否為幸運充能魚
func isLuckyChargeFish(defID string) bool {
	return defID == "T183"
}

// getLuckyChargeBurst 取得充能爆發倍率（供 handleKill 使用，一次性消耗）
func (g *Game) getLuckyChargeBurst(p *player.Player) float64 {
	mgr := g.LuckyChargeFish
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	sess, ok := mgr.sessions[p.ID]
	if !ok || !sess.burstReady || time.Now().After(sess.activeUntil) {
		return 1.0
	}

	// 消耗爆發，重置充能
	sess.burstReady = false
	atomic.StoreInt64(&sess.count, 0)
	return LuckyChargeBurstMult
}

// notifyChargeShot 每次射擊累積充能值（供 handleAttack 使用）
func (g *Game) notifyChargeShot(p *player.Player) {
	mgr := g.LuckyChargeFish
	mgr.mu.Lock()

	sess, ok := mgr.sessions[p.ID]
	if !ok || time.Now().After(sess.activeUntil) {
		mgr.mu.Unlock()
		return
	}

	// 已有爆發就緒，不再累積
	if sess.burstReady {
		mgr.mu.Unlock()
		return
	}

	instanceID := sess.instanceID
	mgr.mu.Unlock()

	// 累積充能值
	newCount := atomic.AddInt64(&sess.count, 1)

	// 廣播充能進度（每 3 點一次）
	if newCount%3 == 0 || newCount == int64(LuckyChargeTarget) {
		g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgLuckyChargeFish,
			Payload: ws.LuckyChargeFishPayload{
				Event:      "charge_progress",
				InstanceID: instanceID,
				Count:      int(newCount),
				Target:     LuckyChargeTarget,
			},
		})
	}

	// 達到目標 → 充能爆發就緒
	if newCount >= int64(LuckyChargeTarget) {
		mgr.mu.Lock()
		if sess2, ok2 := mgr.sessions[p.ID]; ok2 && sess2.instanceID == instanceID {
			sess2.burstReady = true
		}
		mgr.mu.Unlock()

		// 通知玩家充能爆發就緒
		g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgLuckyChargeFish,
			Payload: ws.LuckyChargeFishPayload{
				Event:      "charge_ready",
				InstanceID: instanceID,
				BurstMult:  LuckyChargeBurstMult,
			},
		})

		log.Printf("[LuckyCharge] player=%s charge ready! burst=×%.1f", p.ID, LuckyChargeBurstMult)
	}
}

// tryLuckyChargeFish 擊破 T183 後觸發充能模式（供 handleKill 使用）
func (g *Game) tryLuckyChargeFish(p *player.Player) {
	mgr := g.LuckyChargeFish
	mgr.mu.Lock()

	// 個人冷卻檢查
	if until, ok := mgr.personalCooldown[p.ID]; ok && time.Now().Before(until) {
		mgr.mu.Unlock()
		return
	}

	// 設定個人冷卻
	mgr.personalCooldown[p.ID] = time.Now().Add(LuckyChargePersonalCD)

	// 建立充能 session
	instanceID := fmt.Sprintf("chg_%d", time.Now().UnixNano())
	sess := &chargeSession{
		instanceID:  instanceID,
		activeUntil: time.Now().Add(LuckyChargeDuration),
	}
	atomic.StoreInt64(&sess.count, 0)
	mgr.sessions[p.ID] = sess
	mgr.mu.Unlock()

	log.Printf("[LuckyCharge] player=%s activated charge mode instance=%s", p.ID, instanceID)

	// 通知玩家充能模式開始（個人訊息）
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyChargeFish,
		Payload: ws.LuckyChargeFishPayload{
			Event:       "charge_start",
			InstanceID:  instanceID,
			Target:      LuckyChargeTarget,
			DurationSec: int(LuckyChargeDuration.Seconds()),
			BurstMult:   LuckyChargeBurstMult,
		},
	})

	// 全服廣播（讓其他玩家知道有人觸發充能）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyChargeFish,
		Payload: ws.LuckyChargeFishPayload{
			Event:      "charge_broadcast",
			InstanceID: instanceID,
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyChargeFish, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("⚡ %s 觸發充能模式！射擊 %d 次觸發 ×%.1f 爆發！",
			p.DisplayName, LuckyChargeTarget, LuckyChargeBurstMult),
		"color": "#F39C12",
	})
	g.broadcastAnnouncement(ann)

	// 12 秒後清除 session
	go func() {
		time.Sleep(LuckyChargeDuration)
		mgr.mu.Lock()
		if s, ok := mgr.sessions[p.ID]; ok && s.instanceID == instanceID {
			delete(mgr.sessions, p.ID)
		}
		mgr.mu.Unlock()

		// 通知玩家充能模式結束
		g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgLuckyChargeFish,
			Payload: ws.LuckyChargeFishPayload{
				Event:      "charge_end",
				InstanceID: instanceID,
			},
		})
	}()
}

// notifyLuckyChargeBurstUsed 充能爆發被使用後廣播（供 handleKill 使用）
func (g *Game) notifyLuckyChargeBurstUsed(p *player.Player, instanceID string, reward int) {
	// 通知玩家爆發已觸發
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyChargeFish,
		Payload: ws.LuckyChargeFishPayload{
			Event:      "charge_burst",
			InstanceID: instanceID,
			BurstMult:  LuckyChargeBurstMult,
			Reward:     reward,
		},
	})

	// 全服廣播（讓其他玩家看到爆發）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyChargeFish,
		Payload: ws.LuckyChargeFishPayload{
			Event:      "charge_burst_broadcast",
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			BurstMult:  LuckyChargeBurstMult,
		},
	})

	log.Printf("[LuckyCharge] player=%s burst used! reward=%d", p.ID, reward)
}

// getChargeInstanceID 取得玩家當前充能 session ID（供 handleKill 使用）
func (g *Game) getChargeInstanceID(playerID string) string {
	mgr := g.LuckyChargeFish
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	if sess, ok := mgr.sessions[playerID]; ok {
		return sess.instanceID
	}
	return ""
}
