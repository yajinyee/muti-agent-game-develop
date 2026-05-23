// lucky_tornado_handler.go — 幸運龍捲風魚系統（DAY-248）
// 業界原創「龍捲風吸引+螺旋爆發」機制
//
// 設計：擊破 T206 後，場景中央生成「龍捲風」（持續 12 秒）：
//   - 龍捲風每 2 秒「吸引」場上所有目標向中央螺旋移動（每次移動 80px，帶旋轉角度）
//   - 龍捲風期間擊破任何目標獲得 ×2.2 倍率加成（乘法）
//   - 12 秒後「龍捲風爆發」：中央 250px 範圍內所有目標 85% 擊破機率（×1.5 倍率，全服共享）
//   - 個人冷卻 22 秒；全服冷卻 35 秒
//
// 設計差異：
//   - 與磁力魚（DAY-232，直線向中央移動）不同，龍捲風是「螺旋移動」
//     目標繞著中央旋轉靠近，讓玩家看到「魚群被龍捲風捲起來」的視覺爽感
//   - 與漩渦魚（DAY-234，磁力聚集）不同，龍捲風有「旋轉角度」
//     讓目標移動路徑更有視覺感，不是單純的直線聚集
//   - 「螺旋移動」讓玩家有「要趁目標螺旋靠近時趕快打」的緊迫感
//   - 「龍捲風爆發」讓玩家有「等待→爆發」的高潮設計
//   - 全服廣播龍捲風位置讓所有玩家都往中央打，製造「全服聚焦」的社交感
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
	LuckyTornadoPersonalCD   = 22 * time.Second // 個人冷卻
	LuckyTornadoGlobalCD     = 35 * time.Second // 全服冷卻
	LuckyTornadoDuration     = 12 * time.Second // 龍捲風持續時間
	LuckyTornadoPullInterval = 2 * time.Second  // 螺旋吸引間隔
	LuckyTornadoPullDist     = 80.0             // 每次移動距離（px）
	LuckyTornadoSpiralAngle  = 30.0             // 螺旋旋轉角度（度）
	LuckyTornadoBlastRadius  = 250.0            // 爆發範圍（px）
	LuckyTornadoBlastChance  = 0.85             // 爆發擊破機率
	LuckyTornadoKillMult     = 2.2              // 龍捲風期間擊破倍率
	LuckyTornadoBlastMult    = 1.5              // 爆發倍率
	LuckyTornadoCenterX      = 500.0            // 場景中央 X
	LuckyTornadoCenterY      = 300.0            // 場景中央 Y
)

// luckyTornadoManager 幸運龍捲風魚管理器
type luckyTornadoManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 龍捲風是否啟動
	tornadoActive bool

	// 觸發者 playerID
	triggerPlayerID string

	// 龍捲風啟動時間
	tornadoStartAt time.Time
}

func newLuckyTornadoManager() *luckyTornadoManager {
	return &luckyTornadoManager{
		personalCooldowns: make(map[string]time.Time),
	}
}

// isLuckyTornadoFish 判斷是否為幸運龍捲風魚
func isLuckyTornadoFish(defID string) bool {
	return defID == "T206"
}

// isTornadoActive 判斷龍捲風是否啟動（供 handleKill 使用）
func (m *luckyTornadoManager) isTornadoActive() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.tornadoActive
}

// getLuckyTornadoBoost 取得龍捲風期間的倍率加成（供 handleKill 使用）
func (m *luckyTornadoManager) getLuckyTornadoBoost() float64 {
	if m.isTornadoActive() {
		return LuckyTornadoKillMult
	}
	return 1.0
}

// tryLuckyTornadoFish 擊破 T206 後觸發龍捲風
func (g *Game) tryLuckyTornadoFish(p *player.Player) {
	m := g.LuckyTornado
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
	// 龍捲風已啟動
	if m.tornadoActive {
		m.mu.Unlock()
		return
	}

	// 設定冷卻和狀態
	m.personalCooldowns[p.ID] = now.Add(LuckyTornadoPersonalCD)
	m.globalCooldownUntil = now.Add(LuckyTornadoGlobalCD)
	m.tornadoActive = true
	m.triggerPlayerID = p.ID
	m.tornadoStartAt = now
	m.mu.Unlock()

	log.Printf("[Tornado] player=%s 龍捲風啟動！持續 %v", p.ID, LuckyTornadoDuration)

	// 個人訊息：龍捲風啟動
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyTornado,
		Payload: ws.LuckyTornadoPayload{
			Event:       "tornado_start",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			CenterX:     LuckyTornadoCenterX,
			CenterY:     LuckyTornadoCenterY,
			DurationSec: int(LuckyTornadoDuration.Seconds()),
			KillMult:    LuckyTornadoKillMult,
			BlastMult:   LuckyTornadoBlastMult,
		},
	})

	// 全服廣播
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyTornado,
		Payload: ws.LuckyTornadoPayload{
			Event:      "tornado_broadcast",
			PlayerName: p.DisplayName,
			CenterX:    LuckyTornadoCenterX,
			CenterY:    LuckyTornadoCenterY,
		},
	})

	// 全服公告
	g.Announce.Create(announce.EventLuckyTornado, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🌪️ %s 觸發龍捲風！螺旋吸引所有目標，×%.1f 倍率！", p.DisplayName, LuckyTornadoKillMult),
		"color":   "#1ABC9C",
	})

	// 啟動龍捲風 goroutine
	go g.runLuckyTornado(p)
}

// runLuckyTornado 龍捲風主 goroutine
func (g *Game) runLuckyTornado(p *player.Player) {
	ticker := time.NewTicker(LuckyTornadoPullInterval)
	defer ticker.Stop()

	endTimer := time.NewTimer(LuckyTornadoDuration)
	defer endTimer.Stop()

	pullCount := 0

	for {
		select {
		case <-ticker.C:
			pullCount++
			g.doTornadoSpiral(pullCount)

		case <-endTimer.C:
			g.doTornadoBlast(p)
			return

		case <-g.stopCh:
			return
		}
	}
}

// doTornadoSpiral 螺旋吸引：所有目標向中央螺旋移動
func (g *Game) doTornadoSpiral(pullCount int) {
	g.mu.Lock()

	type movedTarget struct {
		instanceID string
		newX       float64
		newY       float64
	}
	var moved []movedTarget

	for _, t := range g.Targets {
		if !t.IsAlive {
			continue
		}

		// 計算目標到中央的向量
		dx := LuckyTornadoCenterX - t.X
		dy := LuckyTornadoCenterY - t.Y
		dist := math.Sqrt(dx*dx + dy*dy)

		if dist < 30.0 {
			// 已在中央附近，不再移動
			continue
		}

		// 螺旋移動：先旋轉角度，再向中央移動
		angleRad := LuckyTornadoSpiralAngle * math.Pi / 180.0
		// 旋轉向量（順時針）
		rotX := dx*math.Cos(angleRad) + dy*math.Sin(angleRad)
		rotY := -dx*math.Sin(angleRad) + dy*math.Cos(angleRad)
		rotDist := math.Sqrt(rotX*rotX + rotY*rotY)

		// 移動距離：min(pullDist, dist)
		moveDist := LuckyTornadoPullDist
		if moveDist > dist {
			moveDist = dist
		}

		// 新位置 = 當前位置 + 旋轉後方向 × 移動距離
		newX := t.X + (rotX/rotDist)*moveDist
		newY := t.Y + (rotY/rotDist)*moveDist

		// 邊界限制
		if newX < 50 {
			newX = 50
		}
		if newX > 950 {
			newX = 950
		}
		if newY < 50 {
			newY = 50
		}
		if newY > 550 {
			newY = 550
		}

		t.X = newX
		t.Y = newY
		moved = append(moved, movedTarget{instanceID: t.InstanceID, newX: newX, newY: newY})
	}
	g.mu.Unlock()

	if len(moved) == 0 {
		return
	}

	// 廣播螺旋移動
	type tornadoPos struct {
		InstanceID string  `json:"instance_id"`
		X          float64 `json:"x"`
		Y          float64 `json:"y"`
	}
	positions := make([]tornadoPos, 0, len(moved))
	for _, m := range moved {
		positions = append(positions, tornadoPos{
			InstanceID: m.instanceID,
			X:          m.newX,
			Y:          m.newY,
		})
	}

	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyTornado,
		Payload: ws.LuckyTornadoPayload{
			Event:       "tornado_spiral",
			PullCount:   pullCount,
			MovedCount:  len(moved),
			Positions:   positions,
		},
	})

	log.Printf("[Tornado] 螺旋吸引 #%d：移動 %d 個目標", pullCount, len(moved))
}

// doTornadoBlast 龍捲風爆發（計時結束時觸發）
func (g *Game) doTornadoBlast(p *player.Player) {
	m := g.LuckyTornado

	// 重置狀態
	m.mu.Lock()
	m.tornadoActive = false
	m.mu.Unlock()

	// 計算 betCost（全服共享，用平均 bet）
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

	// 收集中央範圍內的目標
	g.mu.Lock()
	var blastTargets []string
	for _, t := range g.Targets {
		if !t.IsAlive {
			continue
		}
		dx := t.X - LuckyTornadoCenterX
		dy := t.Y - LuckyTornadoCenterY
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist <= LuckyTornadoBlastRadius {
			blastTargets = append(blastTargets, t.InstanceID)
		}
	}
	g.mu.Unlock()

	blastCount := 0
	totalReward := 0

	for _, instanceID := range blastTargets {
		// 85% 擊破機率
		if randFloat() > LuckyTornadoBlastChance {
			continue
		}

		g.mu.Lock()
		t, exists := g.Targets[instanceID]
		if !exists || !t.IsAlive {
			g.mu.Unlock()
			continue
		}
		delete(g.Targets, instanceID)
		g.mu.Unlock()

		reward := int(float64(avgBet) * LuckyTornadoBlastMult)
		totalReward += reward
		blastCount++

		// 全服共享獎勵
		g.mu.RLock()
		players := make([]*player.Player, 0, len(g.Players))
		for _, pl := range g.Players {
			players = append(players, pl)
		}
		g.mu.RUnlock()

		if len(players) > 0 {
			share := reward / len(players)
			if share < 1 {
				share = 1
			}
			for _, pl := range players {
				pl.AddCoins(share)
			}
		}
	}

	log.Printf("[Tornado] player=%s 龍捲風爆發！爆炸 %d 個目標，總獎勵 %d", p.ID, blastCount, totalReward)

	// 全服廣播爆發
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyTornado,
		Payload: ws.LuckyTornadoPayload{
			Event:       "tornado_blast",
			BlastCount:  blastCount,
			TotalReward: totalReward,
			BlastMult:   LuckyTornadoBlastMult,
		},
	})

	// 結束通知（個人）
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyTornado,
		Payload: ws.LuckyTornadoPayload{
			Event:       "tornado_end",
			BlastCount:  blastCount,
			TotalReward: totalReward,
		},
	})

	if blastCount >= 3 {
		g.Announce.Create(announce.EventLuckyTornado, p.DisplayName, totalReward, map[string]string{
			"message": fmt.Sprintf("🌪️ %s 龍捲風爆發！捲走 %d 個目標，全服獲得 %d 籌碼！",
				p.DisplayName, blastCount, totalReward),
			"color": "#16A085",
		})
	}
}
