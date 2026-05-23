// lucky_black_hole_handler.go — 幸運黑洞魚系統（DAY-221）
// 業界原創「重力黑洞」機制
//
// 設計：擊破 T179 後在場上建立「重力黑洞」（持續 10 秒）：
//   - 黑洞建立在場景中央附近（隨機偏移），半徑 350px
//   - 黑洞範圍內所有目標每 1 秒被「吸引」（HP -10%，模擬重力傷害）
//   - 黑洞範圍內目標被擊破：獎勵 ×2.0 倍率加成（乘法）
//   - 10 秒後「奇點爆炸」：黑洞範圍內所有目標 85% 擊破機率（0.70x 倍率，全服共享）
//   - 個人冷卻 22 秒；全服冷卻 35 秒
//
// 設計差異：
//   - 與幸運熱區魚（DAY-210，空間限定 ×2.0，脈衝 HP -15%）不同，
//     黑洞魚是「重力吸引 + 奇點爆炸」，視覺上更震撼
//   - 「重力傷害」讓玩家感受到「黑洞在幫我削血」的輔助感
//   - 「奇點爆炸」是「等待→爆發」的高潮設計，85% 擊破機率是最高的
//   - 黑洞範圍 350px（比熱區 280px 更大），覆蓋更多目標
//   - 全服廣播黑洞位置讓所有玩家都往同一個地方打，製造「全服聚焦」的社交感
//   - 「奇點爆炸」的視覺效果（黑洞收縮→爆炸）是業界最震撼的特效之一
package game

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"

	"digital-twin/server/internal/data"
	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	LuckyBlackHolePersonalCD  = 22 * time.Second // 個人冷卻
	LuckyBlackHoleGlobalCD    = 35 * time.Second // 全服冷卻
	LuckyBlackHoleDuration    = 10 * time.Second // 黑洞持續時間
	LuckyBlackHoleRadius      = 350.0            // 黑洞半徑（px）
	LuckyBlackHolePulseDmg    = 0.10             // 每秒重力傷害（HP -10%）
	LuckyBlackHoleKillMult    = 2.0              // 黑洞範圍內擊破倍率加成（乘法）
	LuckyBlackHoleBlastChance = 0.85             // 奇點爆炸擊破機率
	LuckyBlackHoleBlastMult   = 0.70             // 奇點爆炸倍率
)

// luckyBlackHoleManager 幸運黑洞魚管理器
type luckyBlackHoleManager struct {
	mu sync.Mutex

	// 個人冷卻
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 當前黑洞狀態
	active       bool
	blackHoleX   float64
	blackHoleY   float64
	expiresAt    time.Time
	instanceID   string // 用於識別當前黑洞實例
}

func newLuckyBlackHoleManager() *luckyBlackHoleManager {
	return &luckyBlackHoleManager{
		personalCooldowns: make(map[string]time.Time),
	}
}

// isLuckyBlackHoleFish 判斷是否為幸運黑洞魚
func isLuckyBlackHoleFish(defID string) bool {
	return defID == "T179"
}

// isInBlackHole 判斷目標是否在黑洞範圍內（供 handleKill 使用）
func (g *Game) isInBlackHole(x, y float64) bool {
	mgr := g.LuckyBlackHole
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	if !mgr.active || time.Now().After(mgr.expiresAt) {
		return false
	}
	dx := x - mgr.blackHoleX
	dy := y - mgr.blackHoleY
	return math.Sqrt(dx*dx+dy*dy) <= LuckyBlackHoleRadius
}

// getLuckyBlackHoleMultiplier 取得黑洞倍率加成（供 handleKill 使用）
// 若目標在黑洞範圍內，回傳 2.0；否則回傳 1.0
func (g *Game) getLuckyBlackHoleMultiplier(x, y float64) float64 {
	if g.isInBlackHole(x, y) {
		return LuckyBlackHoleKillMult
	}
	return 1.0
}

// tryLuckyBlackHoleFish 擊破 T179 後觸發黑洞（供 handleKill 使用）
func (g *Game) tryLuckyBlackHoleFish(p *player.Player) {
	mgr := g.LuckyBlackHole
	mgr.mu.Lock()

	// 全服冷卻檢查
	if time.Now().Before(mgr.globalCooldownUntil) {
		mgr.mu.Unlock()
		return
	}

	// 個人冷卻檢查
	if cd, ok := mgr.personalCooldowns[p.ID]; ok {
		if time.Now().Before(cd) {
			mgr.mu.Unlock()
			return
		}
	}

	// 已有黑洞在運作中
	if mgr.active {
		mgr.mu.Unlock()
		return
	}

	// 設定冷卻
	mgr.personalCooldowns[p.ID] = time.Now().Add(LuckyBlackHolePersonalCD)
	mgr.globalCooldownUntil = time.Now().Add(LuckyBlackHoleGlobalCD)

	// 決定黑洞位置（場景中央偏移，確保在畫面內）
	// 假設場景寬 1280，高 720，黑洞中心在中央區域隨機
	bhX := 400.0 + rand.Float64()*480.0 // 400~880
	bhY := 180.0 + rand.Float64()*360.0 // 180~540

	// 啟動黑洞
	mgr.active = true
	mgr.blackHoleX = bhX
	mgr.blackHoleY = bhY
	mgr.expiresAt = time.Now().Add(LuckyBlackHoleDuration)
	instanceID := fmt.Sprintf("bh_%d", time.Now().UnixNano())
	mgr.instanceID = instanceID
	mgr.mu.Unlock()

	// 全服廣播：黑洞建立
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyBlackHole,
		Payload: ws.LuckyBlackHolePayload{
			Event:       "blackhole_start",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			X:           bhX,
			Y:           bhY,
			Radius:      LuckyBlackHoleRadius,
			DurationSec: int(LuckyBlackHoleDuration.Seconds()),
			InstanceID:  instanceID,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyBlackHole, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🌑 %s 召喚黑洞！10 秒後奇點爆炸！", p.DisplayName),
		"color":   "#8B00FF",
	})
	g.broadcastAnnouncement(ann)

	log.Printf("[LuckyBlackHole] player=%s created black hole at (%.0f,%.0f) instance=%s",
		p.ID, bhX, bhY, instanceID)

	// 啟動黑洞 goroutine（每秒脈衝 + 10 秒後奇點爆炸）
	go g.runLuckyBlackHole(p, bhX, bhY, instanceID)
}

// runLuckyBlackHole 黑洞運作 goroutine
func (g *Game) runLuckyBlackHole(p *player.Player, bhX, bhY float64, instanceID string) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	deadline := time.NewTimer(LuckyBlackHoleDuration)
	defer deadline.Stop()

	pulseCount := 0

	for {
		select {
		case <-ticker.C:
			pulseCount++
			g.doBlackHolePulse(bhX, bhY, instanceID, pulseCount)

		case <-deadline.C:
			// 奇點爆炸
			g.doBlackHoleSingularityBlast(p, bhX, bhY, instanceID)
			return
		}
	}
}

// doBlackHolePulse 黑洞每秒重力脈衝（HP -10%）
func (g *Game) doBlackHolePulse(bhX, bhY float64, instanceID string, pulseNum int) {
	// 確認黑洞仍有效
	mgr := g.LuckyBlackHole
	mgr.mu.Lock()
	if !mgr.active || mgr.instanceID != instanceID {
		mgr.mu.Unlock()
		return
	}
	mgr.mu.Unlock()

	g.mu.Lock()
	affectedCount := 0
	for _, t := range g.Targets {
		if t.HP <= 0 {
			continue
		}
		dx := t.X - bhX
		dy := t.Y - bhY
		if math.Sqrt(dx*dx+dy*dy) <= LuckyBlackHoleRadius {
			// 重力傷害：HP -10%（最少保留 1）
			dmg := int(float64(t.HP) * LuckyBlackHolePulseDmg)
			if dmg < 1 {
				dmg = 1
			}
			t.HP -= dmg
			if t.HP < 1 {
				t.HP = 1
			}
			affectedCount++
		}
	}
	g.mu.Unlock()

	if affectedCount > 0 {
		// 廣播脈衝
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyBlackHole,
			Payload: ws.LuckyBlackHolePayload{
				Event:         "blackhole_pulse",
				InstanceID:    instanceID,
				PulseNum:      pulseNum,
				AffectedCount: affectedCount,
				X:             bhX,
				Y:             bhY,
			},
		})
		log.Printf("[LuckyBlackHole] pulse#%d: affected=%d targets", pulseNum, affectedCount)
	}
}

// doBlackHoleSingularityBlast 奇點爆炸（黑洞結束時觸發）
func (g *Game) doBlackHoleSingularityBlast(p *player.Player, bhX, bhY float64, instanceID string) {
	// 清除黑洞狀態
	mgr := g.LuckyBlackHole
	mgr.mu.Lock()
	if mgr.instanceID != instanceID {
		mgr.mu.Unlock()
		return
	}
	mgr.active = false
	mgr.mu.Unlock()

	// 廣播奇點爆炸開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyBlackHole,
		Payload: ws.LuckyBlackHolePayload{
			Event:      "singularity_blast",
			InstanceID: instanceID,
			X:          bhX,
			Y:          bhY,
			Radius:     LuckyBlackHoleRadius,
		},
	})

	// 稍等 300ms 讓 Client 播放爆炸動畫
	time.Sleep(300 * time.Millisecond)

	// 對黑洞範圍內所有目標執行奇點爆炸
	g.mu.Lock()
	type blastTarget struct {
		id         string
		multiplier float64
		x, y       float64
	}
	var targets []blastTarget
	for id, t := range g.Targets {
		if t.HP <= 0 {
			continue
		}
		dx := t.X - bhX
		dy := t.Y - bhY
		if math.Sqrt(dx*dx+dy*dy) <= LuckyBlackHoleRadius {
			targets = append(targets, blastTarget{
				id:         id,
				multiplier: t.Multiplier,
				x:          t.X,
				y:          t.Y,
			})
		}
	}
	g.mu.Unlock()

	// 計算平均 betCost（全服共享獎勵）
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
	g.mu.RUnlock()

	killedCount := 0
	totalReward := 0

	for _, bt := range targets {
		if rand.Float64() >= LuckyBlackHoleBlastChance {
			continue
		}

		g.mu.Lock()
		t := g.Targets[bt.id]
		if t == nil || t.HP <= 0 {
			g.mu.Unlock()
			continue
		}
		t.HP = 0
		delete(g.Targets, bt.id)
		g.mu.Unlock()

		reward := int(bt.multiplier * float64(avgBet) * LuckyBlackHoleBlastMult)
		totalReward += reward
		killedCount++

		// 廣播單個目標爆炸
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyBlackHole,
			Payload: ws.LuckyBlackHolePayload{
				Event:    "singularity_hit",
				TargetID: bt.id,
				Reward:   reward,
				X:        bt.x,
				Y:        bt.y,
			},
		})
	}

	// 全服共享獎勵
	if totalReward > 0 {
		g.mu.RLock()
		playerCount := len(g.Players)
		players := make([]*player.Player, 0, playerCount)
		for _, pl := range g.Players {
			players = append(players, pl)
		}
		g.mu.RUnlock()

		if playerCount > 0 {
			share := totalReward / playerCount
			if share < 1 {
				share = 1
			}
			for _, pl := range players {
				pl.AddCoins(share)
			}
		}
	}

	// 廣播奇點爆炸結算
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyBlackHole,
		Payload: ws.LuckyBlackHolePayload{
			Event:       "singularity_result",
			InstanceID:  instanceID,
			KilledCount: killedCount,
			TotalReward: totalReward,
		},
	})

	log.Printf("[LuckyBlackHole] singularity blast: killed=%d totalReward=%d", killedCount, totalReward)

	// 全服公告（≥4 個擊破才公告）
	if killedCount >= 4 {
		color := "#8B00FF"
		if killedCount >= 8 {
			color = "#FF00FF"
		}
		ann := g.Announce.Create(announce.EventLuckyBlackHole, p.DisplayName, killedCount, map[string]string{
			"message": fmt.Sprintf("🌑 奇點爆炸！消滅 %d 個目標！獎勵 %d 金幣！", killedCount, totalReward),
			"color":   color,
		})
		g.broadcastAnnouncement(ann)
	}
}
