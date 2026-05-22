// lucky_hot_zone_handler.go — 幸運熱區魚空間策略系統（DAY-210）
// 業界依據：King of Ocean 2026（Galaxsys）
// 「The electric jellyfish chains current between adjacent targets,
//  paying multipliers from every link in the chain.」
// + Ocean King 4 Brand New World 2025
// 「Golden Zone — a glowing area appears on screen, all fish inside
//  receive a 2x multiplier bonus. Zone lasts 8 seconds then explodes,
//  capturing all remaining fish within the zone.」
//
// 設計：擊破 T168 後在場上建立「幸運熱區」（半徑 280px，持續 8 秒）：
//   1. 熱區內所有目標獲得 ×2.0 倍率加成（乘法，最強加成）
//   2. 每 1 秒「熱區脈衝」：熱區內目標 HP 降低 15%（讓玩家更容易擊破）
//   3. 8 秒後「熱區爆炸」：熱區內所有目標 75% 擊破機率（0.65x 倍率，全服共享）
//   4. 個人冷卻 20 秒；全服冷卻 30 秒（防止熱區疊加）
//   5. 全服廣播熱區位置（讓所有玩家都能往熱區打）
//
// 設計差異（與其他倍率加成系統的區別）：
//   - 幸運星魚（DAY-160）：個人 ×2，10秒，全場
//   - 黃金波浪魚（DAY-207）：全服 ×2.0，8秒，全場
//   - 幸運熱區魚（DAY-210）：全服 ×2.0，8秒，「空間限定」（只在熱區內）
//   - 「空間限定」讓玩家有「要把魚引到熱區裡打」的策略感
//   - 熱區脈衝讓玩家感受到「熱區在幫我削血」的輔助感
//   - 熱區爆炸是「等待→爆發」的高潮設計
//   - 全服廣播熱區位置讓所有玩家都往同一個地方打，製造「全服聚焦」的社交感
package game

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/game/target"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// 幸運熱區常數（DAY-210）
const (
	LuckyHotZoneRadius       = 280.0  // 熱區半徑（px）
	LuckyHotZoneDuration     = 8      // 熱區持續時間（秒）
	LuckyHotZoneMultiplier   = 2.0    // 熱區倍率加成（乘法）
	LuckyHotZonePulseSec     = 1      // 脈衝間隔（秒）
	LuckyHotZonePulseHPRatio = 0.15   // 脈衝 HP 降低比例（15%）
	LuckyHotZoneBlastProb    = 0.75   // 爆炸擊破機率（75%）
	LuckyHotZoneBlastMult    = 0.65   // 爆炸獎勵倍率（0.65x）
	LuckyHotZonePersonalCD   = 20     // 個人冷卻（秒）
	LuckyHotZoneGlobalCD     = 30     // 全服冷卻（秒）
)

// luckyHotZoneManager 幸運熱區管理器
type luckyHotZoneManager struct {
	mu          sync.Mutex
	cooldowns   map[string]time.Time // playerID -> 個人冷卻結束時間
	globalCDEnd time.Time            // 全服冷卻結束時間
	zoneActive  bool                 // 熱區是否正在進行
	zoneX       float64              // 熱區中心 X
	zoneY       float64              // 熱區中心 Y
	zoneEnd     time.Time            // 熱區結束時間
}

func newLuckyHotZoneManager() *luckyHotZoneManager {
	return &luckyHotZoneManager{
		cooldowns: make(map[string]time.Time),
	}
}

// isLuckyHotZoneFish 判斷是否為幸運熱區魚（T168，DAY-210）
func isLuckyHotZoneFish(defID string) bool {
	return defID == "T168"
}

// getLuckyHotZoneMultiplier 取得熱區倍率加成（供 handleKill 使用）
// 若目標在熱區內且熱區正在進行，回傳 2.0；否則回傳 1.0
func (g *Game) getLuckyHotZoneMultiplier(tx, ty float64) float64 {
	mgr := g.LuckyHotZone
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	if !mgr.zoneActive || time.Now().After(mgr.zoneEnd) {
		return 1.0
	}
	dist := math.Sqrt(math.Pow(tx-mgr.zoneX, 2) + math.Pow(ty-mgr.zoneY, 2))
	if dist <= LuckyHotZoneRadius {
		return LuckyHotZoneMultiplier
	}
	return 1.0
}

// tryLuckyHotZone 擊破 T168 後觸發幸運熱區
func (g *Game) tryLuckyHotZone(p *player.Player, t *target.Target) {
	mgr := g.LuckyHotZone
	mgr.mu.Lock()

	// 全服冷卻檢查
	if time.Now().Before(mgr.globalCDEnd) {
		mgr.mu.Unlock()
		return
	}
	// 個人冷卻檢查
	if cd, ok := mgr.cooldowns[p.ID]; ok && time.Now().Before(cd) {
		mgr.mu.Unlock()
		return
	}
	// 熱區已在進行中
	if mgr.zoneActive && time.Now().Before(mgr.zoneEnd) {
		mgr.mu.Unlock()
		return
	}

	// 設定冷卻
	mgr.cooldowns[p.ID] = time.Now().Add(LuckyHotZonePersonalCD * time.Second)
	mgr.globalCDEnd = time.Now().Add(LuckyHotZoneGlobalCD * time.Second)

	// 決定熱區位置（以擊破點為中心，確保在畫面內）
	zoneX := clampFloat(t.X, LuckyHotZoneRadius, 1280-LuckyHotZoneRadius)
	zoneY := clampFloat(t.Y, LuckyHotZoneRadius, 720-LuckyHotZoneRadius)

	mgr.zoneActive = true
	mgr.zoneX = zoneX
	mgr.zoneY = zoneY
	mgr.zoneEnd = time.Now().Add(LuckyHotZoneDuration * time.Second)
	mgr.mu.Unlock()

	log.Printf("[LuckyHotZone] player=%s triggered zone at (%.0f,%.0f) for %ds",
		p.ID, zoneX, zoneY, LuckyHotZoneDuration)

	// 廣播熱區開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyHotZone,
		Payload: ws.LuckyHotZonePayload{
			Event:      "zone_start",
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			ZoneX:      zoneX,
			ZoneY:      zoneY,
			Radius:     LuckyHotZoneRadius,
			DurationSec: LuckyHotZoneDuration,
			Multiplier: LuckyHotZoneMultiplier,
		},
	})

	// 全服公告
	msg := fmt.Sprintf("🔥 %s 觸發幸運熱區！熱區內目標 ×%.0f 倍率！快往熱區打！",
		p.DisplayName, LuckyHotZoneMultiplier)
	ann := g.Announce.Create(announce.EventBigWin, p.DisplayName, 0, map[string]string{
		"message": msg,
		"color":   "#FF6600",
	})
	g.broadcastAnnouncement(ann)

	// 啟動熱區 goroutine（脈衝 + 爆炸）
	go g.runLuckyHotZone(p, zoneX, zoneY)
}

// runLuckyHotZone 熱區主循環（脈衝 + 最終爆炸）
func (g *Game) runLuckyHotZone(p *player.Player, zoneX, zoneY float64) {
	ticker := time.NewTicker(LuckyHotZonePulseSec * time.Second)
	defer ticker.Stop()

	endTimer := time.NewTimer(LuckyHotZoneDuration * time.Second)
	defer endTimer.Stop()

	pulseCount := 0

	for {
		select {
		case <-ticker.C:
			pulseCount++
			g.doLuckyHotZonePulse(zoneX, zoneY, pulseCount)

		case <-endTimer.C:
			// 熱區結束 → 爆炸清場
			g.doLuckyHotZoneBlast(p, zoneX, zoneY)
			return
		}
	}
}

// doLuckyHotZonePulse 熱區脈衝（每秒降低熱區內目標 HP 15%）
func (g *Game) doLuckyHotZonePulse(zoneX, zoneY float64, pulseNum int) {
	affectedCount := 0
	g.mu.Lock()
	for _, t := range g.Targets {
		if !t.IsAlive || isLuckyHotZoneFish(t.DefID) {
			continue
		}
		dist := math.Sqrt(math.Pow(t.X-zoneX, 2) + math.Pow(t.Y-zoneY, 2))
		if dist > LuckyHotZoneRadius {
			continue
		}
		// 降低 HP 15%
		reduction := int(float64(t.HP) * LuckyHotZonePulseHPRatio)
		if reduction < 1 {
			reduction = 1
		}
		t.HP -= reduction
		if t.HP < 1 {
			t.HP = 1
		}
		affectedCount++
	}
	g.mu.Unlock()

	// 廣播脈衝（每 2 次脈衝廣播一次，減少網路流量）
	if pulseNum%2 == 0 || pulseNum == 1 {
		remaining := LuckyHotZoneDuration - pulseNum
		if remaining < 0 {
			remaining = 0
		}
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyHotZone,
			Payload: ws.LuckyHotZonePayload{
				Event:         "zone_pulse",
				ZoneX:         zoneX,
				ZoneY:         zoneY,
				Radius:        LuckyHotZoneRadius,
				PulseNum:      pulseNum,
				AffectedCount: affectedCount,
				RemainingSec:  remaining,
			},
		})
	}

	log.Printf("[LuckyHotZone] pulse #%d: affected=%d targets", pulseNum, affectedCount)
}

// doLuckyHotZoneBlast 熱區爆炸（清場）
func (g *Game) doLuckyHotZoneBlast(p *player.Player, zoneX, zoneY float64) {
	mgr := g.LuckyHotZone
	mgr.mu.Lock()
	mgr.zoneActive = false
	mgr.mu.Unlock()

	// 收集熱區內所有存活目標
	type blastTarget struct {
		t      *target.Target
		reward int
	}
	var blastTargets []blastTarget

	g.mu.Lock()
	for _, t := range g.Targets {
		if !t.IsAlive || isLuckyHotZoneFish(t.DefID) {
			continue
		}
		dist := math.Sqrt(math.Pow(t.X-zoneX, 2) + math.Pow(t.Y-zoneY, 2))
		if dist > LuckyHotZoneRadius {
			continue
		}
		// 75% 擊破機率
		if rand.Float64() < LuckyHotZoneBlastProb {
			def := t.Def
			mult := (def.MultiplierMin + def.MultiplierMax) / 2.0
			reward := int(mult * LuckyHotZoneBlastMult)
			if reward < 1 {
				reward = 1
			}
			blastTargets = append(blastTargets, blastTarget{t: t, reward: reward})
		}
	}
	g.mu.Unlock()

	// 執行擊破並分配獎勵給所有玩家
	totalReward := 0
	killedCount := 0
	g.mu.Lock()
	playerCount := len(g.Players)
	if playerCount < 1 {
		playerCount = 1
	}
	for _, bt := range blastTargets {
		if !bt.t.IsAlive {
			continue
		}
		bt.t.IsAlive = false
		bt.t.HP = 0
		killedCount++
		totalReward += bt.reward
	}
	// 按玩家數平均分配獎勵
	rewardPerPlayer := totalReward / playerCount
	if rewardPerPlayer < 1 && totalReward > 0 {
		rewardPerPlayer = 1
	}
	for _, pp := range g.Players {
		pp.Coins += rewardPerPlayer
	}
	g.mu.Unlock()

	log.Printf("[LuckyHotZone] blast: killed=%d totalReward=%d perPlayer=%d",
		killedCount, totalReward, rewardPerPlayer)

	// 廣播爆炸結果
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyHotZone,
		Payload: ws.LuckyHotZonePayload{
			Event:         "zone_blast",
			PlayerID:      p.ID,
			PlayerName:    p.DisplayName,
			ZoneX:         zoneX,
			ZoneY:         zoneY,
			Radius:        LuckyHotZoneRadius,
			KilledCount:   killedCount,
			TotalReward:   totalReward,
			RewardPerPlayer: rewardPerPlayer,
		},
	})

	// 全服公告（≥5 個擊破）
	if killedCount >= 5 {
		color := "#FF6600"
		if killedCount >= 10 {
			color = "#FF4500"
		}
		msg := fmt.Sprintf("💥 幸運熱區爆炸！擊破 %d 個目標！每位玩家獲得 %d 金幣！",
			killedCount, rewardPerPlayer)
		ann := g.Announce.Create(announce.EventBigWin, p.DisplayName, totalReward, map[string]string{
			"message": msg,
			"color":   color,
		})
		g.broadcastAnnouncement(ann)
	}
}

// clampFloat 限制浮點數在 [min, max] 範圍內
func clampFloat(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
