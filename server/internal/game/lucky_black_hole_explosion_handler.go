// lucky_black_hole_explosion_handler.go — 幸運黑洞爆炸魚系統（DAY-249）
// 業界原創「黑洞吸收+能量爆炸」機制
//
// 設計：擊破 T207 後，場景中央生成「黑洞」（持續 10 秒）：
//   - 黑洞每 1.5 秒「吸收」場上距離最近的目標（直接消滅，×1.2 倍率，個人獎勵）
//   - 黑洞最多吸收 6 個目標，每吸收一個「能量充能 +1」
//   - 10 秒後「黑洞爆炸」：能量值 × 場上目標數 × 0.8 倍率（全服共享）
//   - 個人冷卻 20 秒；全服冷卻 30 秒
//
// 設計差異：
//   - 與龍捲風（螺旋移動+爆發）不同，黑洞是「直接吸收消滅+能量累積爆炸」
//     讓玩家看到「魚一條一條被黑洞吸走消失」的爽感
//   - 「能量充能」讓玩家有「吸越多，爆炸越強」的期待感
//   - 「全服共享爆炸」讓所有玩家都受益，製造社交感
//   - 「最多 6 個」讓玩家有「要趁黑洞還在時多打周圍的魚」的策略感
//   - 黑洞爆炸倍率 = 能量值 × 場上目標數 × 0.8，讓玩家有「場上魚越多，爆炸越強」的動機
package game

import (
	"fmt"
	"log"
	"math"
	"sync"
	"time"

	"digital-twin/server/internal/data"
	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	LuckyBlackHoleExplosionPersonalCD    = 20 * time.Second            // 個人冷卻
	LuckyBlackHoleExplosionGlobalCD      = 30 * time.Second            // 全服冷卻
	LuckyBlackHoleExplosionDuration      = 10 * time.Second            // 黑洞持續時間
	LuckyBlackHoleExplosionAbsorbInterval = 1500 * time.Millisecond    // 吸收間隔
	LuckyBlackHoleExplosionMaxAbsorb     = 6                           // 最多吸收數量
	LuckyBlackHoleExplosionAbsorbMult    = 1.2                         // 吸收倍率（個人）
	LuckyBlackHoleExplosionBlastBaseMult = 0.8                         // 爆炸基礎倍率係數
	LuckyBlackHoleExplosionCenterX       = 500.0                       // 黑洞中心 X
	LuckyBlackHoleExplosionCenterY       = 300.0                       // 黑洞中心 Y
	LuckyBlackHoleExplosionAbsorbRadius  = 600.0                       // 吸收搜尋範圍（px）
)

// luckyBlackHoleExplosionManager 幸運黑洞爆炸魚管理器
type luckyBlackHoleExplosionManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 黑洞是否啟動
	blackHoleActive bool

	// 觸發者 playerID
	triggerPlayerID string

	// 已吸收數量（能量值）
	absorbedCount int
}

func newLuckyBlackHoleExplosionManager() *luckyBlackHoleExplosionManager {
	return &luckyBlackHoleExplosionManager{
		personalCooldowns: make(map[string]time.Time),
	}
}

// isLuckyBlackHoleExplosionFish 判斷是否為幸運黑洞爆炸魚
func isLuckyBlackHoleExplosionFish(defID string) bool {
	return defID == "T207"
}

// isBlackHoleExplosionActive 判斷黑洞是否啟動（供外部使用）
func (m *luckyBlackHoleExplosionManager) isBlackHoleExplosionActive() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.blackHoleActive
}

// tryLuckyBlackHoleExplosionFish 擊破 T207 後觸發黑洞
func (g *Game) tryLuckyBlackHoleExplosionFish(p *player.Player) {
	m := g.LuckyBlackHoleExplosion
	m.mu.Lock()

	now := time.Now()

	// 全服冷卻檢查
	if now.Before(m.globalCooldownUntil) {
		m.mu.Unlock()
		return
	}
	// 個人冷卻檢查
	if cd, ok := m.personalCooldowns[p.ID]; ok && now.Before(cd) {
		m.mu.Unlock()
		return
	}
	// 黑洞已啟動
	if m.blackHoleActive {
		m.mu.Unlock()
		return
	}

	// 設定冷卻和狀態
	m.personalCooldowns[p.ID] = now.Add(LuckyBlackHoleExplosionPersonalCD)
	m.globalCooldownUntil = now.Add(LuckyBlackHoleExplosionGlobalCD)
	m.blackHoleActive = true
	m.triggerPlayerID = p.ID
	m.absorbedCount = 0
	m.mu.Unlock()

	log.Printf("[BlackHoleExplosion] player=%s 黑洞生成！持續 %v，最多吸收 %d 個目標",
		p.ID, LuckyBlackHoleExplosionDuration, LuckyBlackHoleExplosionMaxAbsorb)

	// 個人訊息：黑洞啟動
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyBlackHoleExplosion,
		Payload: ws.LuckyBlackHoleExplosionPayload{
			Event:       "blackhole_start",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			CenterX:     LuckyBlackHoleExplosionCenterX,
			CenterY:     LuckyBlackHoleExplosionCenterY,
			DurationSec: int(LuckyBlackHoleExplosionDuration.Seconds()),
			MaxAbsorb:   LuckyBlackHoleExplosionMaxAbsorb,
			AbsorbMult:  LuckyBlackHoleExplosionAbsorbMult,
		},
	})

	// 全服廣播
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyBlackHoleExplosion,
		Payload: ws.LuckyBlackHoleExplosionPayload{
			Event:      "blackhole_broadcast",
			PlayerName: p.DisplayName,
			CenterX:    LuckyBlackHoleExplosionCenterX,
			CenterY:    LuckyBlackHoleExplosionCenterY,
		},
	})

	// 全服公告
	g.Announce.Create(announce.EventLuckyBlackHoleExplosion, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🕳️ %s 觸發黑洞！吸收目標累積能量，爆炸全服共享！", p.DisplayName),
		"color":   "#2C3E50",
	})

	// 啟動黑洞 goroutine
	go g.runLuckyBlackHoleExplosion(p)
}

// runLuckyBlackHoleExplosion 黑洞主 goroutine
func (g *Game) runLuckyBlackHoleExplosion(p *player.Player) {
	ticker := time.NewTicker(LuckyBlackHoleExplosionAbsorbInterval)
	defer ticker.Stop()

	endTimer := time.NewTimer(LuckyBlackHoleExplosionDuration)
	defer endTimer.Stop()

	for {
		select {
		case <-ticker.C:
			// 吸收最近目標
			absorbed := g.doBlackHoleAbsorb(p)
			if absorbed {
				// 檢查是否達到最大吸收數
				g.LuckyBlackHoleExplosion.mu.Lock()
				count := g.LuckyBlackHoleExplosion.absorbedCount
				g.LuckyBlackHoleExplosion.mu.Unlock()
				if count >= LuckyBlackHoleExplosionMaxAbsorb {
					// 提前爆炸
					g.doBlackHoleExplosion(p)
					return
				}
			}

		case <-endTimer.C:
			g.doBlackHoleExplosion(p)
			return

		case <-g.stopCh:
			return
		}
	}
}

// doBlackHoleAbsorb 吸收最近的目標
// 回傳是否成功吸收
func (g *Game) doBlackHoleAbsorb(p *player.Player) bool {
	// 找距離黑洞中心最近的存活目標
	g.mu.Lock()

	var nearestID string
	nearestDist := math.MaxFloat64

	for _, t := range g.Targets {
		if !t.IsAlive {
			continue
		}
		dx := t.X - LuckyBlackHoleExplosionCenterX
		dy := t.Y - LuckyBlackHoleExplosionCenterY
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist < nearestDist && dist <= LuckyBlackHoleExplosionAbsorbRadius {
			nearestDist = dist
			nearestID = t.InstanceID
		}
	}

	if nearestID == "" {
		g.mu.Unlock()
		return false
	}

	// 取得目標資訊並消滅
	t, exists := g.Targets[nearestID]
	if !exists || !t.IsAlive {
		g.mu.Unlock()
		return false
	}

	defID := t.DefID
	targetName := defID
	if def, ok := data.Targets[defID]; ok {
		targetName = def.Name
	}
	targetX := t.X
	targetY := t.Y

	delete(g.Targets, nearestID)
	g.mu.Unlock()

	// 計算獎勵（個人）
	betDef := data.GetBetDef(p.BetLevel)
	reward := int(float64(betDef.BetCost) * LuckyBlackHoleExplosionAbsorbMult)
	if reward < 1 {
		reward = 1
	}
	p.AddCoins(reward)

	// 更新能量計數
	g.LuckyBlackHoleExplosion.mu.Lock()
	g.LuckyBlackHoleExplosion.absorbedCount++
	absorbedCount := g.LuckyBlackHoleExplosion.absorbedCount
	g.LuckyBlackHoleExplosion.mu.Unlock()

	log.Printf("[BlackHoleExplosion] player=%s 吸收 #%d：%s（距離 %.0fpx），獎勵 %d",
		p.ID, absorbedCount, targetName, nearestDist, reward)

	// 廣播吸收事件（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyBlackHoleExplosion,
		Payload: ws.LuckyBlackHoleExplosionPayload{
			Event:         "blackhole_absorb",
			AbsorbCount:   absorbedCount,
			MaxAbsorb:     LuckyBlackHoleExplosionMaxAbsorb,
			TargetName:    targetName,
			TargetX:       targetX,
			TargetY:       targetY,
			InstanceID:    nearestID,
			Reward:        reward,
		},
	})

	return true
}

// doBlackHoleExplosion 黑洞爆炸（計時結束或達到最大吸收數時觸發）
func (g *Game) doBlackHoleExplosion(p *player.Player) {
	m := g.LuckyBlackHoleExplosion

	// 取得能量值並重置狀態
	m.mu.Lock()
	energy := m.absorbedCount
	m.blackHoleActive = false
	m.mu.Unlock()

	// 計算場上目標數
	g.mu.RLock()
	targetCount := 0
	for _, t := range g.Targets {
		if t.IsAlive {
			targetCount++
		}
	}
	g.mu.RUnlock()

	// 爆炸倍率 = 能量值 × 場上目標數 × 0.8（最少 1.0）
	blastMult := float64(energy) * float64(targetCount) * LuckyBlackHoleExplosionBlastBaseMult
	if blastMult < 1.0 {
		blastMult = 1.0
	}
	// 上限 50.0 防止爆炸過強
	if blastMult > 50.0 {
		blastMult = 50.0
	}

	// 計算全服共享獎勵
	g.mu.RLock()
	totalBet := 0
	playerCount := 0
	for _, pl := range g.Players {
		betDef := data.GetBetDef(pl.BetLevel)
		totalBet += betDef.BetCost
		playerCount++
	}
	g.mu.RUnlock()

	avgBet := 1
	if playerCount > 0 {
		avgBet = totalBet / playerCount
	}
	if avgBet < 1 {
		avgBet = 1
	}

	totalReward := int(float64(avgBet) * blastMult)
	if totalReward < 1 {
		totalReward = 1
	}

	// 全服共享獎勵
	g.mu.RLock()
	players := make([]*player.Player, 0, len(g.Players))
	for _, pl := range g.Players {
		players = append(players, pl)
	}
	g.mu.RUnlock()

	if len(players) > 0 {
		share := totalReward / len(players)
		if share < 1 {
			share = 1
		}
		for _, pl := range players {
			pl.AddCoins(share)
		}
	}

	log.Printf("[BlackHoleExplosion] player=%s 黑洞爆炸！能量=%d，場上目標=%d，倍率=%.1f，總獎勵=%d",
		p.ID, energy, targetCount, blastMult, totalReward)

	// 全服廣播爆炸
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyBlackHoleExplosion,
		Payload: ws.LuckyBlackHoleExplosionPayload{
			Event:       "blackhole_explosion",
			Energy:      energy,
			TargetCount: targetCount,
			BlastMult:   blastMult,
			TotalReward: totalReward,
		},
	})

	// 結束通知（個人）
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyBlackHoleExplosion,
		Payload: ws.LuckyBlackHoleExplosionPayload{
			Event:       "blackhole_end",
			Energy:      energy,
			BlastMult:   blastMult,
			TotalReward: totalReward,
		},
	})

	// 全服公告（能量 ≥ 3 才公告）
	if energy >= 3 {
		g.Announce.Create(announce.EventLuckyBlackHoleExplosion, p.DisplayName, totalReward, map[string]string{
			"message": fmt.Sprintf("🕳️ %s 黑洞爆炸！吸收 %d 個目標，×%.1f 倍率，全服獲得 %d 籌碼！",
				p.DisplayName, energy, blastMult, totalReward),
			"color": "#1A252F",
		})
	}
}
