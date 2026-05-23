// lucky_quantum_entangle_handler.go — 幸運量子糾纏魚系統（DAY-251）
// 業界原創「量子糾纏+同步爆炸+量子共鳴」機制
//
// 設計：擊破 T209 後，場上隨機 2 個目標被「量子糾纏」（持續 20 秒）：
//   - 任何玩家擊破其中一個 → 另一個立刻「同步爆炸」（×1.8 倍率，全服共享）
//   - 若兩個在 1.5 秒內被不同玩家擊破 → 觸發「量子共鳴」：全服 ×3.5 倍率大獎
//   - 20 秒後未擊破 → 「量子衰變」：兩個目標 HP -60%（安慰獎）
//   - 個人冷卻 25 秒；全服冷卻 40 秒
//
// 設計差異：
//   - 與鏡像分裂（T208，目標分裂成副本）不同，量子糾纏是「兩個真實目標互相連結」
//     讓玩家有「打一個，另一個也爆」的驚喜感
//   - 「量子共鳴」鼓勵多玩家協作，1.5 秒內同時擊破兩個目標，製造「全服合力」的社交感
//   - ×3.5 倍率是目前全服合力類最高倍率，讓玩家有「要趕快找隊友一起打」的動機
//   - 「量子衰變 HP -60%」確保即使沒人打也有安慰獎，降低挫敗感
//   - 全服廣播糾纏位置讓所有玩家都看到「這兩條魚是連結的」，製造「全服一起盯著」的緊張感
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"

	"digital-twin/server/internal/data"
	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/game/target"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	LuckyQuantumEntanglePersonalCD  = 25 * time.Second  // 個人冷卻
	LuckyQuantumEntangleGlobalCD    = 40 * time.Second  // 全服冷卻
	LuckyQuantumEntangleDuration    = 20 * time.Second  // 糾纏持續時間
	LuckyQuantumEntangleSyncMult    = 1.8               // 同步爆炸倍率（全服共享）
	LuckyQuantumEntangleResonMult   = 3.5               // 量子共鳴倍率（全服共享）
	LuckyQuantumEntangleDecayHP     = 0.6               // 量子衰變 HP 扣除比例
	LuckyQuantumEntangleResonWindow = 1500 * time.Millisecond // 量子共鳴時間窗口
)

// quantumEntangleSession 量子糾纏會話
type quantumEntangleSession struct {
	sessionID    string
	triggerID    string    // 觸發玩家 ID
	triggerName  string    // 觸發玩家名稱
	targetA      string    // 糾纏目標 A InstanceID
	targetB      string    // 糾纏目標 B InstanceID
	defIDA       string    // 目標 A DefID
	defIDB       string    // 目標 B DefID
	multA        float64   // 目標 A 倍率
	multB        float64   // 目標 B 倍率
	expiresAt    time.Time // 糾纏到期時間
	firstKillAt  time.Time // 第一個被擊破的時間（零值表示尚未擊破）
	firstKillerID string   // 第一個擊破者 ID
	firstKillerName string // 第一個擊破者名稱
	firstKilledID string   // 第一個被擊破的目標 ID
	mu           sync.Mutex
}

// luckyQuantumEntangleManager 幸運量子糾纏魚管理器
type luckyQuantumEntangleManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 活躍的糾纏會話（sessionID → session）
	activeSessions map[string]*quantumEntangleSession

	// 目標 → 會話映射（instanceID → sessionID）
	targetToSession map[string]string
}

func newLuckyQuantumEntangleManager() *luckyQuantumEntangleManager {
	return &luckyQuantumEntangleManager{
		personalCooldowns: make(map[string]time.Time),
		activeSessions:    make(map[string]*quantumEntangleSession),
		targetToSession:   make(map[string]string),
	}
}

// isLuckyQuantumEntangleFish 判斷是否為幸運量子糾纏魚
func isLuckyQuantumEntangleFish(defID string) bool {
	return defID == "T209"
}

// isQuantumEntangleTarget 判斷是否為量子糾纏目標（供 handleKill 使用）
// 回傳 (是否為糾纏目標, sessionID)
func (m *luckyQuantumEntangleManager) isQuantumEntangleTarget(instanceID string) (bool, string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if sid, ok := m.targetToSession[instanceID]; ok {
		return true, sid
	}
	return false, ""
}

// getSession 取得會話（需在外部加鎖）
func (m *luckyQuantumEntangleManager) getSession(sessionID string) *quantumEntangleSession {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.activeSessions[sessionID]
}

// removeSession 移除會話及其目標映射
func (m *luckyQuantumEntangleManager) removeSession(sessionID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if sess, ok := m.activeSessions[sessionID]; ok {
		delete(m.targetToSession, sess.targetA)
		delete(m.targetToSession, sess.targetB)
		delete(m.activeSessions, sessionID)
	}
}

// tryLuckyQuantumEntangleFish 擊破 T209 後觸發量子糾纏
func (g *Game) tryLuckyQuantumEntangleFish(p *player.Player) {
	m := g.LuckyQuantumEntangle
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

	// 設定冷卻
	m.personalCooldowns[p.ID] = now.Add(LuckyQuantumEntanglePersonalCD)
	m.globalCooldownUntil = now.Add(LuckyQuantumEntangleGlobalCD)
	m.mu.Unlock()

	// 選取場上 2 個隨機目標進行糾纏
	g.mu.Lock()
	var candidates []*target.Target
	for _, t := range g.Targets {
		if !t.IsAlive {
			continue
		}
		// 排除 BOSS、Bonus 目標
		if def, ok := data.Targets[t.DefID]; ok {
			if def.Type == data.TargetTypeBoss || def.Type == data.TargetTypeBonus {
				continue
			}
		}
		// 排除已被糾纏的目標
		if _, alreadyEntangled := m.targetToSession[t.InstanceID]; alreadyEntangled {
			continue
		}
		candidates = append(candidates, t)
	}

	shuffleTargets(candidates)
	if len(candidates) < 2 {
		g.mu.Unlock()
		return
	}
	tA := candidates[0]
	tB := candidates[1]

	sessionID := uuid.New().String()
	expiresAt := now.Add(LuckyQuantumEntangleDuration)

	sess := &quantumEntangleSession{
		sessionID:   sessionID,
		triggerID:   p.ID,
		triggerName: p.DisplayName,
		targetA:     tA.InstanceID,
		targetB:     tB.InstanceID,
		defIDA:      tA.DefID,
		defIDB:      tB.DefID,
		multA:       tA.Multiplier,
		multB:       tB.Multiplier,
		expiresAt:   expiresAt,
	}

	m.mu.Lock()
	m.activeSessions[sessionID] = sess
	m.targetToSession[tA.InstanceID] = sessionID
	m.targetToSession[tB.InstanceID] = sessionID
	m.mu.Unlock()

	// 廣播用資訊
	type entangleInfo struct {
		InstanceID string  `json:"instance_id"`
		DefID      string  `json:"def_id"`
		X          float64 `json:"x"`
		Y          float64 `json:"y"`
		Mult       float64 `json:"mult"`
	}
	targets := []entangleInfo{
		{InstanceID: tA.InstanceID, DefID: tA.DefID, X: tA.X, Y: tA.Y, Mult: tA.Multiplier},
		{InstanceID: tB.InstanceID, DefID: tB.DefID, X: tB.X, Y: tB.Y, Mult: tB.Multiplier},
	}
	g.mu.Unlock()

	log.Printf("[QuantumEntangle] player=%s 量子糾纏！目標 A=%s B=%s", p.ID, tA.InstanceID, tB.InstanceID)

	// 個人訊息：糾纏啟動
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyQuantumEntangle,
		Payload: ws.LuckyQuantumEntanglePayload{
			Event:       "entangle_start",
			SessionID:   sessionID,
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			Targets:     targets,
			DurationSec: int(LuckyQuantumEntangleDuration.Seconds()),
			SyncMult:    LuckyQuantumEntangleSyncMult,
			ResonMult:   LuckyQuantumEntangleResonMult,
		},
	})

	// 全服廣播
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyQuantumEntangle,
		Payload: ws.LuckyQuantumEntanglePayload{
			Event:      "entangle_broadcast",
			SessionID:  sessionID,
			PlayerName: p.DisplayName,
			Targets:    targets,
			SyncMult:   LuckyQuantumEntangleSyncMult,
			ResonMult:  LuckyQuantumEntangleResonMult,
		},
	})

	// 全服公告
	g.Announce.Create(announce.EventLuckyQuantumEntangle, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("⚛️ %s 觸發量子糾纏！2 個目標被量子連結！同時擊破可觸發 ×%.1f 量子共鳴！",
			p.DisplayName, LuckyQuantumEntangleResonMult),
		"color": "#1A5276",
	})

	// 啟動衰變計時 goroutine
	go g.runQuantumEntangleDecay(sessionID, expiresAt)
}

// notifyQuantumEntangleKill 量子糾纏目標被擊破時的處理（由 handleKill 呼叫）
func (g *Game) notifyQuantumEntangleKill(p *player.Player, killedInstanceID string, sessionID string) {
	m := g.LuckyQuantumEntangle
	sess := m.getSession(sessionID)
	if sess == nil {
		return
	}

	sess.mu.Lock()
	now := time.Now()

	// 判斷是第一個還是第二個被擊破
	isFirstKill := sess.firstKillAt.IsZero()

	if isFirstKill {
		// 第一個被擊破
		sess.firstKillAt = now
		sess.firstKillerID = p.ID
		sess.firstKillerName = p.DisplayName
		sess.firstKilledID = killedInstanceID
		sess.mu.Unlock()

		// 找出另一個目標
		var otherID string
		if killedInstanceID == sess.targetA {
			otherID = sess.targetB
		} else {
			otherID = sess.targetA
		}

		log.Printf("[QuantumEntangle] player=%s 擊破第一個糾纏目標 %s，同步爆炸 %s", p.ID, killedInstanceID, otherID)

		// 同步爆炸另一個目標
		g.doQuantumSyncExplosion(sess, otherID, p)

	} else {
		// 第二個被擊破
		timeDiff := now.Sub(sess.firstKillAt)
		firstKillerName := sess.firstKillerName
		sess.mu.Unlock()

		// 移除會話
		m.removeSession(sessionID)

		if timeDiff <= LuckyQuantumEntangleResonWindow {
			// 量子共鳴！兩個在 1.5 秒內被擊破
			log.Printf("[QuantumEntangle] 量子共鳴！timeDiff=%.2fs", timeDiff.Seconds())
			g.doQuantumResonance(sess, p, firstKillerName, timeDiff)
		} else {
			// 普通第二次擊破（同步爆炸已在第一次觸發，這裡只廣播結束）
			g.Hub.Broadcast(&ws.Message{
				Type: ws.MsgLuckyQuantumEntangle,
				Payload: ws.LuckyQuantumEntanglePayload{
					Event:     "entangle_end",
					SessionID: sessionID,
				},
			})
		}
	}
}

// doQuantumSyncExplosion 同步爆炸另一個糾纏目標（全服共享獎勵）
func (g *Game) doQuantumSyncExplosion(sess *quantumEntangleSession, otherInstanceID string, triggerPlayer *player.Player) {
	// 計算全服共享獎勵
	g.mu.Lock()
	var otherMult float64
	if t, ok := g.Targets[otherInstanceID]; ok && t.IsAlive {
		otherMult = t.Multiplier
		// 消滅目標
		t.IsAlive = false
		delete(g.Targets, otherInstanceID)
	}
	g.mu.Unlock()

	// 移除糾纏映射（只移除另一個目標的映射，保留第一個以便第二次擊破判斷）
	g.LuckyQuantumEntangle.mu.Lock()
	delete(g.LuckyQuantumEntangle.targetToSession, otherInstanceID)
	g.LuckyQuantumEntangle.mu.Unlock()

	if otherMult <= 0 {
		otherMult = 1.0
	}

	// 計算全服共享獎勵
	g.mu.RLock()
	players := make([]*player.Player, 0, len(g.Players))
	for _, pl := range g.Players {
		players = append(players, pl)
	}
	avgBet := 1
	if len(players) > 0 {
		totalBet := 0
		for _, pl := range players {
			totalBet += data.GetBetDef(pl.BetLevel).BetCost
		}
		avgBet = totalBet / len(players)
	}
	g.mu.RUnlock()

	if avgBet < 1 {
		avgBet = 1
	}

	totalReward := int(float64(avgBet) * otherMult * LuckyQuantumEntangleSyncMult)
	if totalReward < 1 {
		totalReward = 1
	}

	// 全服共享獎勵
	if len(players) > 0 {
		share := totalReward / len(players)
		if share < 1 {
			share = 1
		}
		for _, pl := range players {
			pl.AddCoins(share)
		}
	}

	log.Printf("[QuantumEntangle] 同步爆炸！otherID=%s mult=%.1f 全服獎勵=%d", otherInstanceID, otherMult, totalReward)

	// 全服廣播同步爆炸
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyQuantumEntangle,
		Payload: ws.LuckyQuantumEntanglePayload{
			Event:       "entangle_sync",
			SessionID:   sess.sessionID,
			PlayerName:  triggerPlayer.DisplayName,
			InstanceID:  otherInstanceID,
			SyncMult:    LuckyQuantumEntangleSyncMult,
			TotalReward: totalReward,
		},
	})

	// 全服公告（同步爆炸）
	g.Announce.Create(announce.EventLuckyQuantumEntangle, triggerPlayer.DisplayName, totalReward, map[string]string{
		"message": fmt.Sprintf("⚛️ %s 擊破糾纏目標！量子同步爆炸！全服獲得 %d 籌碼！",
			triggerPlayer.DisplayName, totalReward),
		"color": "#2471A3",
	})
}

// doQuantumResonance 量子共鳴（兩個目標在 1.5 秒內被擊破）
func (g *Game) doQuantumResonance(sess *quantumEntangleSession, secondKiller *player.Player, firstKillerName string, timeDiff time.Duration) {
	// 計算量子共鳴全服大獎
	g.mu.RLock()
	players := make([]*player.Player, 0, len(g.Players))
	for _, pl := range g.Players {
		players = append(players, pl)
	}
	avgBet := 1
	if len(players) > 0 {
		totalBet := 0
		for _, pl := range players {
			totalBet += data.GetBetDef(pl.BetLevel).BetCost
		}
		avgBet = totalBet / len(players)
	}
	g.mu.RUnlock()

	if avgBet < 1 {
		avgBet = 1
	}

	// 量子共鳴獎勵 = avgBet × ResonMult × 2（兩個目標的平均倍率）
	avgMult := (sess.multA + sess.multB) / 2.0
	totalReward := int(float64(avgBet) * avgMult * LuckyQuantumEntangleResonMult)
	if totalReward < 1 {
		totalReward = 1
	}

	// 全服共享大獎
	if len(players) > 0 {
		share := totalReward / len(players)
		if share < 1 {
			share = 1
		}
		for _, pl := range players {
			pl.AddCoins(share)
		}
	}

	log.Printf("[QuantumEntangle] 量子共鳴！timeDiff=%.2fs 全服大獎=%d", timeDiff.Seconds(), totalReward)

	// 全服廣播量子共鳴
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyQuantumEntangle,
		Payload: ws.LuckyQuantumEntanglePayload{
			Event:           "entangle_resonance",
			SessionID:       sess.sessionID,
			PlayerName:      secondKiller.DisplayName,
			FirstKillerName: firstKillerName,
			ResonMult:       LuckyQuantumEntangleResonMult,
			TotalReward:     totalReward,
			TimeDiffMs:      int(timeDiff.Milliseconds()),
		},
	})

	// 全服公告（量子共鳴大獎）
	g.Announce.Create(announce.EventLuckyQuantumEntangle, secondKiller.DisplayName, totalReward, map[string]string{
		"message": fmt.Sprintf("⚛️ 量子共鳴！%s + %s 同時擊破糾纏目標！全服獲得 %d 籌碼大獎！",
			firstKillerName, secondKiller.DisplayName, totalReward),
		"color": "#1A5276",
	})
}

// runQuantumEntangleDecay 量子衰變計時 goroutine
func (g *Game) runQuantumEntangleDecay(sessionID string, expiresAt time.Time) {
	remaining := time.Until(expiresAt)
	if remaining <= 0 {
		return
	}

	timer := time.NewTimer(remaining)
	defer timer.Stop()

	select {
	case <-timer.C:
		g.doQuantumEntangleDecay(sessionID)
	case <-g.stopCh:
		return
	}
}

// doQuantumEntangleDecay 執行量子衰變（20 秒後未擊破）
func (g *Game) doQuantumEntangleDecay(sessionID string) {
	m := g.LuckyQuantumEntangle
	sess := m.getSession(sessionID)
	if sess == nil {
		// 已被擊破，不需要衰變
		return
	}

	// 移除會話
	m.removeSession(sessionID)

	// 對兩個目標執行 HP -60%
	g.mu.Lock()
	decayCount := 0
	for _, id := range []string{sess.targetA, sess.targetB} {
		if t, ok := g.Targets[id]; ok && t.IsAlive {
			reduction := int(float64(t.HP) * LuckyQuantumEntangleDecayHP)
			t.HP -= reduction
			if t.HP < 1 {
				t.HP = 1
			}
			decayCount++
		}
	}
	g.mu.Unlock()

	if decayCount == 0 {
		return
	}

	log.Printf("[QuantumEntangle] 量子衰變！sessionID=%s 衰變 %d 個目標 HP-60%%", sessionID, decayCount)

	// 全服廣播衰變
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyQuantumEntangle,
		Payload: ws.LuckyQuantumEntanglePayload{
			Event:      "entangle_decay",
			SessionID:  sessionID,
			PlayerName: sess.triggerName,
			DecayCount: decayCount,
			DecayHP:    int(LuckyQuantumEntangleDecayHP * 100),
		},
	})

	// 全服公告（衰變）
	g.Announce.Create(announce.EventLuckyQuantumEntangle, sess.triggerName, 0, map[string]string{
		"message": fmt.Sprintf("⚛️ 量子衰變！%s 的糾纏目標衰變，HP -60%%！趕快擊破！",
			sess.triggerName),
		"color": "#7F8C8D",
	})
}
