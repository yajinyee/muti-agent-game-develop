// lucky_mirror_fish_handler.go — T121 幸運鏡像魚系統
// server-event-agent 負責維護
// 業界依據：Fishing Fortune 2026「Mirror Fish — duplicates your next 3 shots, each shot fires twice」
// 設計：擊破 T121 後，觸發「鏡像模式」：玩家下 3 次攻擊每次自動複製一次（等同打兩發）
// 複製攻擊命中同一目標，傷害 ×1.0（不額外加成）
// 若 3 次複製攻擊全部命中 → 「完美鏡像」：全服 ×1.8 加成 5 秒
// 個人冷卻 16 秒；全服冷卻 28 秒
package game

import (
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyMirrorFishManager struct {
	mu              sync.Mutex
	playerCooldowns map[string]time.Time
	globalCooldown  time.Time
	// 活躍的鏡像會話：playerID -> session
	activeSessions map[string]*mirrorFishSession
	// 完美鏡像全服加成
	perfectBoost *mirrorPerfectBoost
}

type mirrorFishSession struct {
	playerID    string
	playerName  string
	shotsLeft   int // 剩餘鏡像次數（最多 3）
	hitCount    int // 成功命中次數
	totalReward int
	expiresAt   time.Time
	settled     bool
}

type mirrorPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyMirrorFishManager() *luckyMirrorFishManager {
	return &luckyMirrorFishManager{
		playerCooldowns: make(map[string]time.Time),
		activeSessions:  make(map[string]*mirrorFishSession),
	}
}

func isLuckyMirrorFish(defID string) bool {
	return defID == "T121"
}

func (m *luckyMirrorFishManager) getMirrorPerfectMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckyMirrorFishManager) isMirrorActive(playerID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	s, ok := m.activeSessions[playerID]
	if !ok || s.settled {
		return false
	}
	if time.Now().After(s.expiresAt) {
		return false
	}
	return s.shotsLeft > 0
}

// notifyMirrorShot 玩家攻擊時，若在鏡像模式中，觸發複製攻擊
// 回傳是否觸發了鏡像（供 handleAttack 使用）
func (g *Game) notifyMirrorShot(playerID string, targetID string) bool {
	m := g.luckyMirrorFish
	m.mu.Lock()
	s, ok := m.activeSessions[playerID]
	if !ok || s.settled || s.shotsLeft <= 0 || time.Now().After(s.expiresAt) {
		m.mu.Unlock()
		return false
	}
	s.shotsLeft--
	playerName := s.playerName
	m.mu.Unlock()

	// 複製攻擊：對同一目標再打一次
	go g.doMirrorShot(playerID, playerName, targetID)
	return true
}

func (g *Game) doMirrorShot(playerID string, playerName string, targetID string) {
	time.Sleep(80 * time.Millisecond) // 稍微延遲，讓原始攻擊先處理

	g.mu.Lock()
	p, ok := g.players[playerID]
	if !ok {
		g.mu.Unlock()
		return
	}
	t, tOk := g.targets[targetID]
	if !tOk || t.HP <= 0 {
		g.mu.Unlock()
		return
	}

	betCost := p.GetBetDef().BetCost
	damage := p.GetBetDef().AttackPower
	t.HP -= damage
	hitCount := 0
	reward := 0

	if t.HP <= 0 {
		t.HP = 0
		reward = int(float64(betCost) * t.Def.Multiplier)
		p.AddCoins(reward)
		g.sendPlayerUpdate(playerID)
		hitCount = 1
	}

	g.hub.Broadcast(protocol.MsgTargetUpdate, protocol.TargetUpdatePayload{
		InstanceID: t.InstanceID,
		HP:         t.HP,
		MaxHP:      t.MaxHP,
		X:          t.X,
		Y:          t.Y,
	})
	g.mu.Unlock()

	// 更新鏡像會話
	m := g.luckyMirrorFish
	m.mu.Lock()
	s, ok2 := m.activeSessions[playerID]
	if ok2 && !s.settled {
		s.hitCount += hitCount
		s.totalReward += reward
		shotsLeft := s.shotsLeft
		hitCnt := s.hitCount
		totalReward := s.totalReward

		// 廣播鏡像命中
		g.hub.Broadcast(protocol.MsgLuckyMirrorFish, protocol.LuckyMirrorFishPayload{
			Event:       "mirror_hit",
			TriggerID:   playerID,
			TriggerName: playerName,
			ShotsLeft:   shotsLeft,
			HitCount:    hitCnt,
			TotalReward: totalReward,
		})

		// 3 次全部用完 → 結算
		if shotsLeft == 0 {
			s.settled = true
			isPerfect := hitCnt >= 3
			m.mu.Unlock()
			g.settleMirrorFish(playerID, playerName, hitCnt, totalReward, isPerfect)
			return
		}
	}
	m.mu.Unlock()
}

func (g *Game) settleMirrorFish(playerID string, playerName string, hitCount int, totalReward int, isPerfect bool) {
	if isPerfect {
		// 完美鏡像：全服 ×1.8 加成 5 秒
		m := g.luckyMirrorFish
		m.mu.Lock()
		m.perfectBoost = &mirrorPerfectBoost{
			mult:      1.8,
			expiresAt: time.Now().Add(5 * time.Second),
		}
		m.mu.Unlock()

		g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
			Message:  "🪞✨ 完美鏡像！" + playerName + " 全服 ×1.8 加成 5 秒！",
			Priority: "high",
			Color:    "#E0AAFF",
		})
		g.hub.Broadcast(protocol.MsgLuckyMirrorFish, protocol.LuckyMirrorFishPayload{
			Event:       "perfect_mirror",
			TriggerID:   playerID,
			TriggerName: playerName,
			HitCount:    hitCount,
			TotalReward: totalReward,
		})

		// 5 秒後廣播加成結束
		go func() {
			time.Sleep(5 * time.Second)
			g.hub.Broadcast(protocol.MsgLuckyMirrorFish, protocol.LuckyMirrorFishPayload{
				Event:       "perfect_end",
				TriggerID:   playerID,
				TriggerName: playerName,
			})
		}()
	}

	g.hub.Broadcast(protocol.MsgLuckyMirrorFish, protocol.LuckyMirrorFishPayload{
		Event:       "settle",
		TriggerID:   playerID,
		TriggerName: playerName,
		HitCount:    hitCount,
		TotalReward: totalReward,
	})

	log.Printf("[MirrorFish] Player %s: hits=%d, reward=%d, perfect=%v",
		playerID, hitCount, totalReward, isPerfect)
}

func (g *Game) tryLuckyMirrorFish(playerID string, killerName string) {
	m := g.luckyMirrorFish
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCooldown) {
		m.mu.Unlock()
		return
	}
	if cd, ok := m.playerCooldowns[playerID]; ok && now.Before(cd) {
		m.mu.Unlock()
		return
	}
	m.playerCooldowns[playerID] = now.Add(16 * time.Second)
	m.globalCooldown = now.Add(28 * time.Second)

	m.activeSessions[playerID] = &mirrorFishSession{
		playerID:   playerID,
		playerName: killerName,
		shotsLeft:  3,
		expiresAt:  now.Add(20 * time.Second), // 20 秒內用完
	}
	m.mu.Unlock()

	g.hub.Broadcast(protocol.MsgLuckyMirrorFish, protocol.LuckyMirrorFishPayload{
		Event:       "trigger",
		TriggerID:   playerID,
		TriggerName: killerName,
		ShotsLeft:   3,
	})
	g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
		Message:  "🪞 " + killerName + " 觸發鏡像魚！下 3 次攻擊自動複製！",
		Priority: "high",
		Color:    "#E0AAFF",
	})

	// 20 秒後超時結算
	go func() {
		time.Sleep(20 * time.Second)
		m.mu.Lock()
		s, ok := m.activeSessions[playerID]
		if !ok || s.settled {
			m.mu.Unlock()
			return
		}
		s.settled = true
		hitCount := s.hitCount
		totalReward := s.totalReward
		m.mu.Unlock()

		g.hub.Broadcast(protocol.MsgLuckyMirrorFish, protocol.LuckyMirrorFishPayload{
			Event:       "timeout",
			TriggerID:   playerID,
			TriggerName: killerName,
			HitCount:    hitCount,
			TotalReward: totalReward,
		})
	}()
}
