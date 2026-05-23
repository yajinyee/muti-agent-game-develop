// lucky_teleport_fish_handler.go — 幸運傳送魚系統（DAY-223）
// 業界原創「傳送混亂」機制
//
// 設計：擊破 T181 後觸發「傳送漩渦」（10 秒）：
//   - 場上所有目標物立即隨機傳送到新位置（瞬間移動）
//   - 傳送後 3 秒內擊破任何目標：獎勵 ×2.5 倍率加成（「傳送混亂」加成）
//   - 每 3 秒再次傳送（最多 3 次，共 4 次傳送）
//   - 個人冷卻 20 秒；全服廣播傳送事件
//
// 設計差異：
//   - 與時間凍結魚（DAY-212，全場靜止）不同，傳送魚是「全場瞬間移動」，
//     讓玩家有「趕快在傳送後 3 秒內打」的緊迫感
//   - 每次傳送都是新的機會，讓玩家保持高度專注
//   - 視覺上所有魚瞬間移動，製造「混亂爽感」
//   - 「傳送混亂加成」讓玩家有「傳送後立刻打」的強烈動機
//   - 全服廣播傳送位置讓所有玩家都看到魚的新位置，製造「全服一起搶打」的社交感
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

const (
	LuckyTeleportPersonalCD    = 20 * time.Second // 個人冷卻
	LuckyTeleportDuration      = 10 * time.Second // 傳送漩渦持續時間
	LuckyTeleportInterval      = 3 * time.Second  // 每次傳送間隔
	LuckyTeleportMaxWaves      = 4                // 最多傳送次數（含初始）
	LuckyTeleportBonusDuration = 3 * time.Second  // 傳送混亂加成持續時間
	LuckyTeleportBonusMult     = 2.5              // 傳送混亂倍率加成
	LuckyTeleportMinX          = 100.0            // 傳送目標 X 最小值
	LuckyTeleportMaxX          = 900.0            // 傳送目標 X 最大值
	LuckyTeleportMinY          = 80.0             // 傳送目標 Y 最小值
	LuckyTeleportMaxY          = 520.0            // 傳送目標 Y 最大值
)

// luckyTeleportFishManager 幸運傳送魚管理器
type luckyTeleportFishManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldown map[string]time.Time

	// 傳送混亂加成狀態（playerID → bonusUntil）
	bonusUntil map[string]time.Time

	// 當前傳送漩渦狀態
	active     bool
	instanceID string
	waveCount  int
}

func newLuckyTeleportFishManager() *luckyTeleportFishManager {
	return &luckyTeleportFishManager{
		personalCooldown: make(map[string]time.Time),
		bonusUntil:       make(map[string]time.Time),
	}
}

// isLuckyTeleportFish 判斷是否為幸運傳送魚
func isLuckyTeleportFish(defID string) bool {
	return defID == "T181"
}

// getLuckyTeleportBonus 取得傳送混亂倍率加成（供 handleKill 使用）
func (g *Game) getLuckyTeleportBonus(p *player.Player) float64 {
	mgr := g.LuckyTeleportFish
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	until, ok := mgr.bonusUntil[p.ID]
	if !ok || time.Now().After(until) {
		return 1.0
	}
	return LuckyTeleportBonusMult
}

// tryLuckyTeleportFish 擊破 T181 後觸發傳送漩渦（供 handleKill 使用）
func (g *Game) tryLuckyTeleportFish(p *player.Player) {
	mgr := g.LuckyTeleportFish
	mgr.mu.Lock()

	// 個人冷卻檢查
	if until, ok := mgr.personalCooldown[p.ID]; ok && time.Now().Before(until) {
		mgr.mu.Unlock()
		return
	}

	// 已有傳送漩渦在運作中（不允許疊加）
	if mgr.active {
		mgr.mu.Unlock()
		return
	}

	// 設定個人冷卻
	mgr.personalCooldown[p.ID] = time.Now().Add(LuckyTeleportPersonalCD)

	// 啟動傳送漩渦
	mgr.active = true
	instanceID := fmt.Sprintf("tp_%d", time.Now().UnixNano())
	mgr.instanceID = instanceID
	mgr.waveCount = 0
	mgr.mu.Unlock()

	log.Printf("[LuckyTeleport] player=%s triggered teleport vortex instance=%s", p.ID, instanceID)

	// 全服廣播：傳送漩渦開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyTeleportFish,
		Payload: ws.LuckyTeleportFishPayload{
			Event:       "teleport_start",
			InstanceID:  instanceID,
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			DurationSec: int(LuckyTeleportDuration.Seconds()),
			MaxWaves:    LuckyTeleportMaxWaves,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyTeleportFish, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🌀 %s 觸發傳送漩渦！所有目標瞬間移動！傳送後 3 秒內擊破獲得 ×%.1f 倍率！",
			p.DisplayName, LuckyTeleportBonusMult),
		"color": "#9B59B6",
	})
	g.broadcastAnnouncement(ann)

	// 立即執行第一次傳送
	go g.runLuckyTeleportWaves(instanceID)
}

// runLuckyTeleportWaves 執行傳送波次（goroutine）
func (g *Game) runLuckyTeleportWaves(instanceID string) {
	for wave := 0; wave < LuckyTeleportMaxWaves; wave++ {
		// 第一波立即執行，後續每 3 秒一次
		if wave > 0 {
			time.Sleep(LuckyTeleportInterval)
		}

		mgr := g.LuckyTeleportFish
		mgr.mu.Lock()
		if !mgr.active || mgr.instanceID != instanceID {
			mgr.mu.Unlock()
			return
		}
		mgr.waveCount = wave + 1
		mgr.mu.Unlock()

		// 執行傳送
		teleportedTargets := g.doLuckyTeleport(instanceID, wave+1)

		if len(teleportedTargets) == 0 {
			continue
		}

		// 設定傳送混亂加成（所有在線玩家）
		bonusEnd := time.Now().Add(LuckyTeleportBonusDuration)
		mgr.mu.Lock()
		g.mu.RLock()
		for playerID := range g.Players {
			mgr.bonusUntil[playerID] = bonusEnd
		}
		g.mu.RUnlock()
		mgr.mu.Unlock()

		// 廣播傳送波次
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyTeleportFish,
			Payload: ws.LuckyTeleportFishPayload{
				Event:      "teleport_wave",
				InstanceID: instanceID,
				Wave:       wave + 1,
				MaxWaves:   LuckyTeleportMaxWaves,
				Targets:    teleportedTargets,
				BonusSec:   int(LuckyTeleportBonusDuration.Seconds()),
				BonusMult:  LuckyTeleportBonusMult,
			},
		})

		log.Printf("[LuckyTeleport] wave=%d teleported=%d targets", wave+1, len(teleportedTargets))
	}

	// 傳送漩渦結束
	mgr := g.LuckyTeleportFish
	mgr.mu.Lock()
	if mgr.active && mgr.instanceID == instanceID {
		mgr.active = false
	}
	mgr.mu.Unlock()

	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyTeleportFish,
		Payload: ws.LuckyTeleportFishPayload{
			Event:      "teleport_end",
			InstanceID: instanceID,
		},
	})

	log.Printf("[LuckyTeleport] vortex ended instance=%s", instanceID)
}

// doLuckyTeleport 執行一次傳送，回傳傳送的目標列表
func (g *Game) doLuckyTeleport(instanceID string, wave int) []ws.TeleportTargetInfo {
	g.mu.Lock()
	defer g.mu.Unlock()

	var teleported []ws.TeleportTargetInfo
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for _, t := range g.Targets {
		if t.HP <= 0 {
			continue
		}
		// 傳送到隨機新位置
		newX := LuckyTeleportMinX + rng.Float64()*(LuckyTeleportMaxX-LuckyTeleportMinX)
		newY := LuckyTeleportMinY + rng.Float64()*(LuckyTeleportMaxY-LuckyTeleportMinY)

		t.X = newX
		t.Y = newY

		teleported = append(teleported, ws.TeleportTargetInfo{
			TargetID: t.InstanceID,
			NewX:     newX,
			NewY:     newY,
		})
	}

	return teleported
}
