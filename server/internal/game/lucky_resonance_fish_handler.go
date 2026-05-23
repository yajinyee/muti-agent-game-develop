// lucky_resonance_fish_handler.go — 幸運共鳴魚系統（DAY-222）
// 業界原創「全服共鳴」機制
//
// 設計：擊破 T180 後觸發「共鳴模式」（15 秒）：
//   - 全服所有玩家的每次射擊都累積「共鳴能量」（+1/shot）
//   - 共鳴能量達到 30 點 → 觸發「共鳴爆發」：
//     全場所有目標 HP -50% + 全服 ×1.8 倍率加成 6 秒
//   - 共鳴爆發獎勵按「貢獻比例」分配（射擊越多，分到越多）
//   - 15 秒內未達到 30 點 → 觸發「小型共鳴」：
//     全場所有目標 HP -25% + 全服 ×1.3 倍率加成 3 秒
//   - 全服冷卻 40 秒
//
// 設計差異：
//   - 與深海龍王（DAY-208，全服合力蓄力，20 點，12 秒）不同，
//     共鳴魚是「貢獻比例分配獎勵」，讓玩家有「我射擊越多，我分到越多」的個人動機
//   - 「共鳴能量進度條」讓全服玩家看到「還差幾槍就爆發」的期待感
//   - 「貢獻比例分配」讓積極射擊的玩家獲得更多獎勵，製造「競爭合作」的社交感
//   - 「共鳴爆發」的 ×1.8 倍率加成讓所有玩家在爆發後 6 秒內都受益
//   - 全服廣播每 5 點進度讓玩家感受到「大家一起在累積」的社群感
package game

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"digital-twin/server/internal/data"
	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	LuckyResonanceGlobalCD      = 40 * time.Second // 全服冷卻
	LuckyResonanceDuration      = 15 * time.Second // 共鳴模式持續時間
	LuckyResonanceTarget        = 30               // 共鳴能量目標
	LuckyResonanceBoostDuration = 6 * time.Second  // 共鳴爆發倍率加成持續時間
	LuckyResonanceSmallDuration = 3 * time.Second  // 小型共鳴倍率加成持續時間
	LuckyResonanceFullBoost     = 1.8              // 共鳴爆發倍率加成
	LuckyResonanceSmallBoost    = 1.3              // 小型共鳴倍率加成
	LuckyResonanceFullHPDrain   = 0.50             // 共鳴爆發 HP 削減
	LuckyResonanceSmallHPDrain  = 0.25             // 小型共鳴 HP 削減
)

// luckyResonanceFishManager 幸運共鳴魚管理器
type luckyResonanceFishManager struct {
	mu sync.Mutex

	// 全服冷卻
	globalCooldownUntil time.Time

	// 共鳴模式狀態
	active     bool
	instanceID string

	// 共鳴能量（atomic，並發安全）
	resonanceCount int64

	// 玩家貢獻記錄（playerID → 射擊次數）
	contributions map[string]int

	// 倍率加成狀態
	boostActive    bool
	boostUntil     time.Time
	boostMultiplier float64
}

func newLuckyResonanceFishManager() *luckyResonanceFishManager {
	return &luckyResonanceFishManager{
		contributions: make(map[string]int),
	}
}

// isLuckyResonanceFish 判斷是否為幸運共鳴魚
func isLuckyResonanceFish(defID string) bool {
	return defID == "T180"
}

// getLuckyResonanceBoost 取得共鳴倍率加成（供 handleKill 使用）
func (g *Game) getLuckyResonanceBoost() float64 {
	mgr := g.LuckyResonanceFish
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	if !mgr.boostActive || time.Now().After(mgr.boostUntil) {
		mgr.boostActive = false
		return 1.0
	}
	return mgr.boostMultiplier
}

// notifyResonanceShot 每次射擊累積共鳴能量（供 handleAttack 使用）
func (g *Game) notifyResonanceShot(p *player.Player) {
	mgr := g.LuckyResonanceFish
	mgr.mu.Lock()
	if !mgr.active {
		mgr.mu.Unlock()
		return
	}
	instanceID := mgr.instanceID
	mgr.mu.Unlock()

	// 累積共鳴能量
	newCount := atomic.AddInt64(&mgr.resonanceCount, 1)

	// 記錄玩家貢獻
	mgr.mu.Lock()
	mgr.contributions[p.ID]++
	mgr.mu.Unlock()

	// 每 5 點廣播進度
	if newCount%5 == 0 {
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyResonanceFish,
			Payload: ws.LuckyResonanceFishPayload{
				Event:      "resonance_progress",
				InstanceID: instanceID,
				Count:      int(newCount),
				Target:     LuckyResonanceTarget,
				PlayerID:   p.ID,
				PlayerName: p.DisplayName,
			},
		})
	}

	// 達到目標 → 觸發共鳴爆發
	if newCount == int64(LuckyResonanceTarget) {
		go g.triggerResonanceBurst(instanceID, true)
	}
}

// tryLuckyResonanceFish 擊破 T180 後觸發共鳴模式（供 handleKill 使用）
func (g *Game) tryLuckyResonanceFish(p *player.Player) {
	mgr := g.LuckyResonanceFish
	mgr.mu.Lock()

	// 全服冷卻檢查
	if time.Now().Before(mgr.globalCooldownUntil) {
		mgr.mu.Unlock()
		return
	}

	// 已有共鳴模式在運作中
	if mgr.active {
		mgr.mu.Unlock()
		return
	}

	// 設定全服冷卻
	mgr.globalCooldownUntil = time.Now().Add(LuckyResonanceGlobalCD)

	// 啟動共鳴模式
	mgr.active = true
	instanceID := fmt.Sprintf("res_%d", time.Now().UnixNano())
	mgr.instanceID = instanceID
	atomic.StoreInt64(&mgr.resonanceCount, 0)
	mgr.contributions = make(map[string]int)
	mgr.mu.Unlock()

	// 全服廣播：共鳴模式開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyResonanceFish,
		Payload: ws.LuckyResonanceFishPayload{
			Event:       "resonance_start",
			InstanceID:  instanceID,
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			Target:      LuckyResonanceTarget,
			DurationSec: int(LuckyResonanceDuration.Seconds()),
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyResonanceFish, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🎵 %s 觸發共鳴模式！全服合力射擊 %d 次觸發共鳴爆發！", p.DisplayName, LuckyResonanceTarget),
		"color":   "#00BFFF",
	})
	g.broadcastAnnouncement(ann)

	log.Printf("[LuckyResonance] player=%s activated resonance mode instance=%s", p.ID, instanceID)

	// 啟動計時器，15 秒後若未達到目標則觸發小型共鳴
	go func() {
		time.Sleep(LuckyResonanceDuration)

		mgr.mu.Lock()
		if !mgr.active || mgr.instanceID != instanceID {
			mgr.mu.Unlock()
			return
		}
		currentCount := atomic.LoadInt64(&mgr.resonanceCount)
		if currentCount < int64(LuckyResonanceTarget) {
			mgr.mu.Unlock()
			// 未達目標，觸發小型共鳴
			go g.triggerResonanceBurst(instanceID, false)
		} else {
			mgr.mu.Unlock()
		}
	}()
}

// triggerResonanceBurst 觸發共鳴爆發（isFull=true 完整爆發，false 小型共鳴）
func (g *Game) triggerResonanceBurst(instanceID string, isFull bool) {
	mgr := g.LuckyResonanceFish
	mgr.mu.Lock()
	if !mgr.active || mgr.instanceID != instanceID {
		mgr.mu.Unlock()
		return
	}

	// 取得貢獻記錄
	contributions := make(map[string]int, len(mgr.contributions))
	for k, v := range mgr.contributions {
		contributions[k] = v
	}
	totalShots := int(atomic.LoadInt64(&mgr.resonanceCount))

	// 清除共鳴模式
	mgr.active = false

	// 設定倍率加成
	var boostMult float64
	var boostDur time.Duration
	var hpDrain float64
	if isFull {
		boostMult = LuckyResonanceFullBoost
		boostDur = LuckyResonanceBoostDuration
		hpDrain = LuckyResonanceFullHPDrain
	} else {
		boostMult = LuckyResonanceSmallBoost
		boostDur = LuckyResonanceSmallDuration
		hpDrain = LuckyResonanceSmallHPDrain
	}
	mgr.boostActive = true
	mgr.boostUntil = time.Now().Add(boostDur)
	mgr.boostMultiplier = boostMult
	mgr.mu.Unlock()

	// 廣播爆發開始
	eventName := "resonance_burst"
	if !isFull {
		eventName = "resonance_small_burst"
	}
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyResonanceFish,
		Payload: ws.LuckyResonanceFishPayload{
			Event:       eventName,
			InstanceID:  instanceID,
			TotalShots:  totalShots,
			BoostMult:   boostMult,
			BoostSec:    int(boostDur.Seconds()),
		},
	})

	// 全場 HP 削減
	g.mu.Lock()
	affectedCount := 0
	for _, t := range g.Targets {
		if t.HP <= 0 {
			continue
		}
		dmg := int(float64(t.HP) * hpDrain)
		if dmg < 1 {
			dmg = 1
		}
		t.HP -= dmg
		if t.HP < 1 {
			t.HP = 1
		}
		affectedCount++
	}
	g.mu.Unlock()

	// 計算貢獻比例獎勵（全服共享獎勵池 = avgBet × totalShots × 0.3）
	g.mu.RLock()
	avgBet := 1
	if len(g.Players) > 0 {
		total := 0
		for _, pl := range g.Players {
			betDef := data.GetBetDef(pl.BetLevel)
			if betDef != nil {
				total += betDef.BetCost
			}
		}
		avgBet = total / len(g.Players)
	}
	players := make(map[string]*player.Player, len(g.Players))
	for id, pl := range g.Players {
		players[id] = pl
	}
	g.mu.RUnlock()

	// 獎勵池
	rewardPool := avgBet * totalShots / 3
	if rewardPool < 1 {
		rewardPool = 1
	}

	// 按貢獻比例分配
	type playerReward struct {
		playerID string
		reward   int
	}
	var rewards []playerReward

	if totalShots > 0 {
		for playerID, shots := range contributions {
			share := rewardPool * shots / totalShots
			if share < 1 {
				share = 1
			}
			rewards = append(rewards, playerReward{playerID: playerID, reward: share})
		}
	}

	// 發放獎勵
	for _, pr := range rewards {
		if pl, ok := players[pr.playerID]; ok {
			pl.AddCoins(pr.reward)
		}
	}

	log.Printf("[LuckyResonance] burst: isFull=%v totalShots=%d affected=%d rewardPool=%d",
		isFull, totalShots, affectedCount, rewardPool)

	// 廣播結算
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyResonanceFish,
		Payload: ws.LuckyResonanceFishPayload{
			Event:         "resonance_result",
			InstanceID:    instanceID,
			AffectedCount: affectedCount,
			RewardPool:    rewardPool,
			TotalShots:    totalShots,
		},
	})

	// 全服公告
	if isFull || affectedCount >= 5 {
		color := "#00BFFF"
		if isFull {
			color = "#00FFFF"
		}
		ann := g.Announce.Create(announce.EventLuckyResonanceFish, "", affectedCount, map[string]string{
			"message": fmt.Sprintf("🎵 共鳴爆發！全服合力 %d 槍！HP -%.0f%%！×%.1f 倍率加成 %d 秒！",
				totalShots, hpDrain*100, boostMult, int(boostDur.Seconds())),
			"color": color,
		})
		g.broadcastAnnouncement(ann)
	}

	// 倍率加成結束後廣播
	go func() {
		time.Sleep(boostDur)
		mgr.mu.Lock()
		mgr.boostActive = false
		mgr.mu.Unlock()

		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyResonanceFish,
			Payload: ws.LuckyResonanceFishPayload{
				Event:      "resonance_boost_end",
				InstanceID: instanceID,
			},
		})
	}()
}
