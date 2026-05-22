// lion_dance_handler.go — 獅子舞大獎爆發系統（DAY-168）
// 業界依據：Fortune King Jackpot（TaDa Gaming 2026）「Lion Dance bonus — triggered by special fish,
// delivers burst multiplier payouts with festive visual effects」
// 擊破 T126 獅子舞魚後，觸發「獅子舞爆發」：全場隨機 3-5 個目標被「獅子舞光環」標記，
// 玩家在 15 秒內擊破標記目標獲得 3x-10x 額外倍率加成。
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

// lionDanceSession 獅子舞爆發 session（per-player）
type lionDanceSession struct {
	markedTargets map[string]float64 // instanceID -> 倍率加成（3x-10x）
	expiresAt     time.Time
	bonusMult     float64 // 本次爆發的倍率（3-10x）
}

// lionDanceManager 獅子舞爆發管理器
type lionDanceManager struct {
	mu                  sync.Mutex
	sessions            map[string]*lionDanceSession // playerID -> session
	globalCooldownUntil time.Time                   // 全服冷卻（30 秒，防止多人同時觸發）
}

func newLionDanceManager() *lionDanceManager {
	return &lionDanceManager{
		sessions: make(map[string]*lionDanceSession),
	}
}

// isLionDance 判斷是否為獅子舞魚
func isLionDance(defID string) bool {
	return defID == "T126"
}

// tryLionDanceBurst 擊破 T126 後觸發獅子舞爆發
func (g *Game) tryLionDanceBurst(p *player.Player, instanceID string, x, y float64) {
	g.LionDance.mu.Lock()
	// 全服冷卻檢查（30 秒）
	if time.Now().Before(g.LionDance.globalCooldownUntil) {
		g.LionDance.mu.Unlock()
		return
	}
	// 設定全服冷卻
	g.LionDance.globalCooldownUntil = time.Now().Add(30 * time.Second)
	g.LionDance.mu.Unlock()

	// 取得所有存活目標（排除觸發者自己剛擊破的那個）
	g.mu.RLock()
	type targetInfo struct {
		id string
		x  float64
		y  float64
	}
	targets := make([]targetInfo, 0, len(g.Targets))
	for id, t := range g.Targets {
		if t.HP > 0 && id != instanceID {
			targets = append(targets, targetInfo{id, t.X, t.Y})
		}
	}
	g.mu.RUnlock()

	// 決定本次爆發倍率（3x-10x，加權：低倍率高機率）
	multOptions := []float64{3, 4, 5, 6, 7, 8, 10}
	multWeights := []int{30, 25, 20, 12, 7, 4, 2}
	burstMult := pickWeightedFloat(multOptions, multWeights)

	// 隨機選 3-5 個目標標記
	markCount := 3 + rand.Intn(3) // 3, 4, or 5
	if len(targets) < markCount {
		markCount = len(targets)
	}
	if markCount == 0 {
		log.Printf("[LionDance] player=%s no targets to mark, skip", p.ID)
		return
	}

	// 隨機打亂目標列表
	rand.Shuffle(len(targets), func(i, j int) {
		targets[i], targets[j] = targets[j], targets[i]
	})
	marked := targets[:markCount]

	// 建立 session
	session := &lionDanceSession{
		markedTargets: make(map[string]float64),
		expiresAt:     time.Now().Add(15 * time.Second),
		bonusMult:     burstMult,
	}
	for _, t := range marked {
		session.markedTargets[t.id] = burstMult
	}

	g.LionDance.mu.Lock()
	g.LionDance.sessions[p.ID] = session
	g.LionDance.mu.Unlock()

	// 建立標記目標列表（供廣播）
	markedList := make([]ws.LionDanceMarkedTarget, 0, markCount)
	for _, t := range marked {
		markedList = append(markedList, ws.LionDanceMarkedTarget{
			InstanceID: t.id,
			X:          t.x,
			Y:          t.y,
		})
	}

	// 全服廣播：獅子舞爆發開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLionDanceBurst,
		Payload: ws.LionDanceBurstPayload{
			Phase:         "burst_start",
			TriggerPlayer: p.ID,
			TriggerName:   p.DisplayName,
			BurstMult:     burstMult,
			MarkedTargets: markedList,
			DurationSec:   15,
		},
	})

	log.Printf("[LionDance] player=%s triggered burst mult=%.0fx marks=%d",
		p.ID, burstMult, markCount)

	// 全服公告（≥7x 才公告）
	if burstMult >= 7 {
		ann := g.Announce.Create(announce.EventLionDance, p.DisplayName, int(burstMult), map[string]string{
			"mult": fmt.Sprintf("%.0f", burstMult),
		})
		g.broadcastAnnouncement(ann)
	}

	// 15 秒後結束 session
	go func() {
		time.Sleep(15 * time.Second)
		g.LionDance.mu.Lock()
		sess, ok := g.LionDance.sessions[p.ID]
		remaining := 0
		if ok {
			remaining = len(sess.markedTargets)
			delete(g.LionDance.sessions, p.ID)
		}
		g.LionDance.mu.Unlock()

		// 廣播結束
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLionDanceBurst,
			Payload: ws.LionDanceBurstPayload{
				Phase:            "burst_end",
				TriggerPlayer:    p.ID,
				TriggerName:      p.DisplayName,
				RemainingTargets: remaining,
			},
		})
	}()
}

// getLionDanceMult 取得獅子舞爆發倍率（供 handleKill 使用）
// 若目標在玩家的標記列表中，回傳倍率加成；否則回傳 1.0
func (g *Game) getLionDanceMult(playerID, targetInstanceID string) float64 {
	g.LionDance.mu.Lock()
	defer g.LionDance.mu.Unlock()

	sess, ok := g.LionDance.sessions[playerID]
	if !ok {
		return 1.0
	}
	if time.Now().After(sess.expiresAt) {
		delete(g.LionDance.sessions, playerID)
		return 1.0
	}
	mult, marked := sess.markedTargets[targetInstanceID]
	if !marked {
		return 1.0
	}
	// 移除已擊破的標記目標
	delete(sess.markedTargets, targetInstanceID)
	return mult
}

// pickWeightedFloat 加權隨機選擇 float64
func pickWeightedFloat(options []float64, weights []int) float64 {
	total := 0
	for _, w := range weights {
		total += w
	}
	r := rand.Intn(total)
	cumulative := 0
	for i, w := range weights {
		cumulative += w
		if r < cumulative {
			return options[i]
		}
	}
	return options[0]
}
