// lucky_mirror_time_handler.go — 幸運鏡像時空魚系統（DAY-227）
// 業界原創「時間倒流」機制
//
// 設計：擊破 T185 後觸發「時間倒流」（8 秒）：
//   - 場上所有目標物的 HP 恢復到「滿血狀態」（HP = MaxHP）
//   - 但倍率提升 ×2.0（讓玩家有「HP 回滿但更值錢」的高報酬感）
//   - 玩家在這 8 秒內擊破任何目標都獲得 ×2.0 倍率加成（乘法）
//   - 8 秒後「時間崩潰」：所有目標 HP -40%（補償玩家）
//   - 個人冷卻 25 秒；全服廣播時間倒流/崩潰
//
// 設計差異：
//   - 與時間凍結魚（DAY-212，全場靜止）不同，時間倒流是「HP 回滿但倍率翻倍」
//   - 與傳送魚（DAY-223，全場瞬間移動）不同，時間倒流是「HP 狀態回溯」
//   - 「HP 回滿但倍率 ×2.0」讓玩家有「要趁 HP 回滿前趕快打」的緊迫感
//   - 「時間崩潰 HP -40%」讓玩家有「等待→爆發」的高潮設計
//   - 全服廣播讓所有玩家都知道「現在是時間倒流，所有魚 HP 回滿但更值錢」
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
	LuckyMirrorTimePersonalCD  = 25 * time.Second // 個人冷卻
	LuckyMirrorTimeDuration    = 8 * time.Second  // 時間倒流持續時間
	LuckyMirrorTimeBoostMult   = 2.0              // 時間倒流期間倍率加成（乘法）
	LuckyMirrorTimeCollapseHP  = 0.40             // 時間崩潰 HP 削減比例
)

// luckyMirrorTimeManager 幸運鏡像時空魚管理器
type luckyMirrorTimeManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldown map[string]time.Time

	// 時間倒流狀態
	active      bool
	activeUntil time.Time
	instanceID  string
}

func newLuckyMirrorTimeManager() *luckyMirrorTimeManager {
	return &luckyMirrorTimeManager{
		personalCooldown: make(map[string]time.Time),
	}
}

// isLuckyMirrorTimeFish 判斷是否為幸運鏡像時空魚
func isLuckyMirrorTimeFish(defID string) bool {
	return defID == "T185"
}

// isLuckyMirrorTimeActive 判斷時間倒流是否啟動（供 handleKill 使用）
func (g *Game) isLuckyMirrorTimeActive() bool {
	mgr := g.LuckyMirrorTime
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	return mgr.active && time.Now().Before(mgr.activeUntil)
}

// getLuckyMirrorTimeBoost 取得時間倒流倍率加成（供 handleKill 使用）
func (g *Game) getLuckyMirrorTimeBoost() float64 {
	if g.isLuckyMirrorTimeActive() {
		return LuckyMirrorTimeBoostMult
	}
	return 1.0
}

// tryLuckyMirrorTimeFish 擊破 T185 後觸發時間倒流（供 handleKill 使用）
func (g *Game) tryLuckyMirrorTimeFish(p *player.Player) {
	mgr := g.LuckyMirrorTime
	mgr.mu.Lock()

	// 個人冷卻檢查
	if until, ok := mgr.personalCooldown[p.ID]; ok && time.Now().Before(until) {
		mgr.mu.Unlock()
		return
	}

	// 設定個人冷卻
	mgr.personalCooldown[p.ID] = time.Now().Add(LuckyMirrorTimePersonalCD)

	// 建立 instance ID
	instanceID := fmt.Sprintf("mtime_%d", time.Now().UnixNano())
	mgr.active = true
	mgr.activeUntil = time.Now().Add(LuckyMirrorTimeDuration)
	mgr.instanceID = instanceID
	mgr.mu.Unlock()

	log.Printf("[LuckyMirrorTime] player=%s triggered time rewind instance=%s", p.ID, instanceID)

	// 執行時間倒流：所有目標 HP 回滿
	g.mu.Lock()
	rewindCount := 0
	for _, t := range g.Targets {
		if t.IsAlive && t.Def.Type != "boss" {
			t.HP = t.MaxHP
			rewindCount++
		}
	}
	g.mu.Unlock()

	log.Printf("[LuckyMirrorTime] rewound %d targets to full HP", rewindCount)

	// 全服廣播時間倒流開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyMirrorTime,
		Payload: ws.LuckyMirrorTimePayload{
			Event:       "time_rewind_start",
			InstanceID:  instanceID,
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			BoostMult:   LuckyMirrorTimeBoostMult,
			DurationSec: int(LuckyMirrorTimeDuration.Seconds()),
			RewindCount: rewindCount,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyMirrorTimeFish, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("⏪ %s 觸發時間倒流！%d 個目標 HP 回滿，倍率 ×%.1f！",
			p.DisplayName, rewindCount, LuckyMirrorTimeBoostMult),
		"color": "#00BFFF",
	})
	g.broadcastAnnouncement(ann)

	// 8 秒後觸發時間崩潰
	go func() {
		time.Sleep(LuckyMirrorTimeDuration)

		// 清除時間倒流狀態
		mgr.mu.Lock()
		if mgr.instanceID == instanceID {
			mgr.active = false
		}
		mgr.mu.Unlock()

		// 時間崩潰：所有目標 HP -40%
		g.mu.Lock()
		collapseCount := 0
		for _, t := range g.Targets {
			if t.IsAlive && t.Def.Type != "boss" {
				newHP := int(float64(t.HP) * (1.0 - LuckyMirrorTimeCollapseHP))
				if newHP < 1 {
					newHP = 1
				}
				t.HP = newHP
				collapseCount++
			}
		}
		g.mu.Unlock()

		log.Printf("[LuckyMirrorTime] time collapse! %d targets HP -40%%", collapseCount)

		// 全服廣播時間崩潰
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyMirrorTime,
			Payload: ws.LuckyMirrorTimePayload{
				Event:         "time_collapse",
				InstanceID:    instanceID,
				CollapseCount: collapseCount,
				CollapseRatio: LuckyMirrorTimeCollapseHP,
			},
		})

		// 全服公告
		ann2 := g.Announce.Create(announce.EventLuckyMirrorTimeFish, p.DisplayName, 0, map[string]string{
			"message": fmt.Sprintf("💥 時間崩潰！%d 個目標 HP -40%%！",
				collapseCount),
			"color": "#FF6B35",
		})
		g.broadcastAnnouncement(ann2)
	}()
}
