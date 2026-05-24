// lucky_star_burst_handler.go — 幸運星爆魚系統（DAY-282）
// 業界依據：業界原創「星爆連鎖+全場星雨+倍率爆炸」機制
// 結合「多點爆炸+累積倍率+共鳴爆發」三個元素
//
// 設計：擊破 T240 後，觸發「星爆」：
//   - Server 隨機生成 5-8 個「星爆點」，依序在 3 秒內爆炸
//   - 每個星爆點爆炸時，場上所有目標 HP -35%（全場 AOE）
//   - 每個爆炸給觸發玩家 ×1.3 累積倍率（最高 ×6.0）
//   - 若有 2 個以上星爆點在 0.5 秒內同時爆炸 → 觸發「星爆共鳴」：全服 ×2.0 加成 5 秒
//   - 全服廣播星爆位置和爆炸結果
//   - 個人冷卻 24 秒；全服冷卻 40 秒
//
// 設計差異：
//   - 與黃金颶風（T234，螺旋掃場 HP-30%）不同，星爆是「多點同時爆炸」
//     讓玩家有「5-8 個星爆點同時炸，全場魚都受傷」的爽感
//   - 「累積倍率 ×1.3/次，最高 ×6.0」讓玩家有「爆炸越多倍率越高」的期待感
//   - 「星爆共鳴（2個以上同時爆炸）」讓玩家有「要是多個星爆點同時炸就觸發共鳴」的驚喜感
//   - 「全服 ×2.0 加成 5 秒」讓所有玩家都受益，製造「全服一起爽」的社交感
//   - 「全服廣播星爆結果」讓所有玩家看到「有幾個星爆點炸了，命中幾條魚」
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
	LuckyStarBurstPersonalCD    = 24 * time.Second // 個人冷卻
	LuckyStarBurstGlobalCD      = 40 * time.Second // 全服冷卻
	LuckyStarBurstMinPoints     = 5                // 最少星爆點數
	LuckyStarBurstMaxPoints     = 8                // 最多星爆點數
	LuckyStarBurstInterval      = 400 * time.Millisecond // 每個星爆點間隔
	LuckyStarBurstHPReduction   = 0.35             // HP 降低比例（-35%）
	LuckyStarBurstAccumMult     = 1.3              // 每次爆炸累積倍率
	LuckyStarBurstMaxAccumMult  = 6.0              // 最高累積倍率
	LuckyStarBurstResonanceMult = 2.0              // 星爆共鳴全服倍率
	LuckyStarBurstResonanceDur  = 5 * time.Second  // 星爆共鳴持續時間
	LuckyStarBurstResonanceWin  = 500 * time.Millisecond // 共鳴判定視窗
)

// starBurstResonanceBoost 星爆共鳴全服加成
type starBurstResonanceBoost struct {
	mult      float64
	expiresAt time.Time
}

// luckyStarBurstManager 幸運星爆魚管理器
type luckyStarBurstManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 星爆共鳴全服加成
	resonanceBoost *starBurstResonanceBoost
}

func newLuckyStarBurstManager() *luckyStarBurstManager {
	return &luckyStarBurstManager{
		personalCooldowns: make(map[string]time.Time),
	}
}

// isLuckyStarBurstFish 判斷是否為幸運星爆魚
func isLuckyStarBurstFish(defID string) bool {
	return defID == "T240"
}

// getStarBurstResonanceMult 取得星爆共鳴全服倍率（供 handleKill 使用）
func (m *luckyStarBurstManager) getStarBurstResonanceMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.resonanceBoost != nil && time.Now().Before(m.resonanceBoost.expiresAt) {
		return m.resonanceBoost.mult
	}
	return 1.0
}

// tryLuckyStarBurstFish 擊破 T240 後觸發星爆（供 handleKill 使用）
func (g *Game) tryLuckyStarBurstFish(p *player.Player) {
	mgr := g.LuckyStarBurst
	mgr.mu.Lock()

	now := time.Now()

	// 全服冷卻檢查
	if now.Before(mgr.globalCooldownUntil) {
		mgr.mu.Unlock()
		return
	}
	// 個人冷卻檢查
	if cd, ok := mgr.personalCooldowns[p.ID]; ok && now.Before(cd) {
		mgr.mu.Unlock()
		return
	}

	// 設定冷卻
	mgr.personalCooldowns[p.ID] = now.Add(LuckyStarBurstPersonalCD)
	mgr.globalCooldownUntil = now.Add(LuckyStarBurstGlobalCD)
	mgr.mu.Unlock()

	// 決定星爆點數量
	count := LuckyStarBurstMinPoints + rand.Intn(LuckyStarBurstMaxPoints-LuckyStarBurstMinPoints+1)

	log.Printf("[StarBurst] player=%s triggered star burst, points=%d", p.ID, count)

	// 全服廣播：星爆觸發
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyStarBurst,
		Payload: ws.LuckyStarBurstPayload{
			Event:      "burst_start",
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			BurstCount: count,
			AccumMult:  1.0,
			Duration:   int(time.Duration(count) * LuckyStarBurstInterval / time.Second),
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyStarBurst, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("⭐ %s 觸發星爆！%d 個星爆點即將爆炸！全場 HP -35%%！",
			p.DisplayName, count),
		"color": "#FFD700",
	})
	g.broadcastAnnouncement(ann)

	// 啟動星爆序列
	go g.runStarBurstSequence(p, count)
}

// runStarBurstSequence 執行星爆序列（goroutine）
func (g *Game) runStarBurstSequence(p *player.Player, count int) {
	accumMult := 1.0
	totalHits := 0
	totalReward := 0
	lastExplodeAt := time.Time{}
	resonanceTriggered := false

	for i := 0; i < count; i++ {
		select {
		case <-time.After(LuckyStarBurstInterval):
		case <-g.stopCh:
			return
		}

		// 對場上所有目標施加 HP -35%
		g.mu.Lock()
		hitCount := 0
		for _, t := range g.Targets {
			if t.HP > 0 && !isLuckyStarBurstFish(t.DefID) {
				reduction := int(float64(t.HP) * LuckyStarBurstHPReduction)
				if reduction < 1 {
					reduction = 1
				}
				t.HP -= reduction
				if t.HP < 1 {
					t.HP = 1
				}
				hitCount++
			}
		}
		g.mu.Unlock()

		totalHits += hitCount

		// 累積倍率
		if accumMult < LuckyStarBurstMaxAccumMult {
			accumMult = accumMult * LuckyStarBurstAccumMult
			if accumMult > LuckyStarBurstMaxAccumMult {
				accumMult = LuckyStarBurstMaxAccumMult
			}
		}

		// 計算本次爆炸獎勵（給觸發玩家）
		snap := p.Snapshot()
		betCost := snap.BetCost
		burstReward := int(float64(betCost) * accumMult * float64(hitCount) * 0.1)
		totalReward += burstReward
		if burstReward > 0 {
			p.AddCoins(burstReward)
		}

		now := time.Now()

		// 星爆共鳴判定：與上一個爆炸時間差 ≤ 0.5 秒
		if !lastExplodeAt.IsZero() && !resonanceTriggered {
			if now.Sub(lastExplodeAt) <= LuckyStarBurstResonanceWin {
				resonanceTriggered = true
				go g.doStarBurstResonance(p)
			}
		}
		lastExplodeAt = now

		log.Printf("[StarBurst] burst %d/%d: hitCount=%d accumMult=%.2f reward=%d",
			i+1, count, hitCount, accumMult, burstReward)

		// 廣播本次爆炸
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyStarBurst,
			Payload: ws.LuckyStarBurstPayload{
				Event:      "burst_explode",
				BurstIndex: i + 1,
				HitCount:   hitCount,
				AccumMult:  accumMult,
			},
		})
	}

	// 結算
	log.Printf("[StarBurst] settle: player=%s totalBursts=%d totalHits=%d finalMult=%.2f totalReward=%d",
		p.ID, count, totalHits, accumMult, totalReward)

	// 全服廣播：星爆結算
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyStarBurst,
		Payload: ws.LuckyStarBurstPayload{
			Event:       "burst_end",
			PlayerName:  p.DisplayName,
			TotalBursts: count,
			TotalHits:   totalHits,
			FinalMult:   accumMult,
			TotalReward: totalReward,
		},
	})

	// 高倍率結算公告
	if accumMult >= 4.0 {
		ann := g.Announce.Create(announce.EventLuckyStarBurst, p.DisplayName, 0, map[string]string{
			"message": fmt.Sprintf("⭐ %s 星爆結算！%d 次爆炸，命中 %d 個目標，最終倍率 ×%.1f！",
				p.DisplayName, count, totalHits, accumMult),
			"color": "#FFD700",
		})
		g.broadcastAnnouncement(ann)
	}
}

// doStarBurstResonance 觸發星爆共鳴（全服 ×2.0 加成 5 秒）
func (g *Game) doStarBurstResonance(p *player.Player) {
	mgr := g.LuckyStarBurst
	mgr.mu.Lock()
	mgr.resonanceBoost = &starBurstResonanceBoost{
		mult:      LuckyStarBurstResonanceMult,
		expiresAt: time.Now().Add(LuckyStarBurstResonanceDur),
	}
	mgr.mu.Unlock()

	log.Printf("[StarBurst] RESONANCE triggered by player=%s, global x%.1f for %ds",
		p.ID, LuckyStarBurstResonanceMult, int(LuckyStarBurstResonanceDur.Seconds()))

	// 全服廣播：星爆共鳴
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyStarBurst,
		Payload: ws.LuckyStarBurstPayload{
			Event:          "burst_resonance",
			PlayerName:     p.DisplayName,
			GlobalMult:     LuckyStarBurstResonanceMult,
			GlobalDuration: int(LuckyStarBurstResonanceDur.Seconds()),
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyStarBurst, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("⭐✨ 星爆共鳴！全服 ×%.1f 加成 %d 秒！",
			LuckyStarBurstResonanceMult, int(LuckyStarBurstResonanceDur.Seconds())),
		"color": "#FFD700",
	})
	g.broadcastAnnouncement(ann)

	// 計時器：共鳴結束後清除
	select {
	case <-time.After(LuckyStarBurstResonanceDur):
	case <-g.stopCh:
		return
	}

	mgr.mu.Lock()
	mgr.resonanceBoost = nil
	mgr.mu.Unlock()
}
