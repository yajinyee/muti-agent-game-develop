// lucky_mirror_split_handler.go — 幸運鏡像分裂魚系統（DAY-250）
// 業界原創「鏡像分裂+雙重目標」機制
//
// 設計：擊破 T208 後，場上隨機 4 個目標被「鏡像分裂」：
//   - 每個目標在其鏡像位置（X 軸對稱）生成一個「鏡像副本」
//   - 鏡像副本 HP = 原目標 50%，倍率 = 原目標 × 0.6（個人獎勵）
//   - 鏡像副本存活 15 秒，玩家擊破獲得個人獎勵
//   - 15 秒後所有未擊破的鏡像副本「鏡像消融」：每個消融給全服 ×0.3 倍率共享獎勵
//   - 個人冷卻 22 秒；全服冷卻 35 秒
//
// 設計差異：
//   - 與分身魚（T200，射擊產生分身子彈）不同，鏡像分裂是「目標本身分裂成兩個」
//     讓玩家有「突然多了一倍的目標可以打」的驚喜感
//   - 「鏡像副本 HP 只有 50%」讓玩家有「這些比較好打」的成就感
//   - 「消融獎勵」確保即使沒打完也有收益，降低挫敗感
//   - 「X 軸鏡像位置」讓玩家有「要同時注意兩個位置」的空間感
//   - 全服廣播分裂位置讓所有玩家都看到鏡像副本，製造「全服一起打」的社交感
package game

import (
	"fmt"
	"log"
	"math/rand"
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
	LuckyMirrorSplitPersonalCD  = 22 * time.Second // 個人冷卻
	LuckyMirrorSplitGlobalCD    = 35 * time.Second // 全服冷卻
	LuckyMirrorSplitDuration    = 15 * time.Second // 鏡像副本存活時間
	LuckyMirrorSplitCount       = 4                // 分裂目標數量
	LuckyMirrorSplitHPRatio     = 0.5              // 鏡像副本 HP 比例
	LuckyMirrorSplitKillMult    = 0.6              // 鏡像副本擊破倍率（個人）
	LuckyMirrorSplitFadeMult    = 0.3              // 消融倍率（全服共享）
	LuckyMirrorSplitCenterX     = 500.0            // 鏡像軸 X（場景中央）
)

// mirrorSplitEntry 鏡像副本記錄
type mirrorSplitEntry struct {
	mirrorInstanceID string    // 鏡像副本 InstanceID
	origInstanceID   string    // 原目標 InstanceID
	origDefID        string    // 原目標 DefID
	origMult         float64   // 原目標倍率
	expiresAt        time.Time // 消融時間
}

// luckyMirrorSplitManager 幸運鏡像分裂魚管理器
type luckyMirrorSplitManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 活躍的鏡像副本（mirrorInstanceID → entry）
	activeMirrors map[string]*mirrorSplitEntry
}

func newLuckyMirrorSplitManager() *luckyMirrorSplitManager {
	return &luckyMirrorSplitManager{
		personalCooldowns: make(map[string]time.Time),
		activeMirrors:     make(map[string]*mirrorSplitEntry),
	}
}

// isLuckyMirrorSplitFish 判斷是否為幸運鏡像分裂魚
func isLuckyMirrorSplitFish(defID string) bool {
	return defID == "T208"
}

// isMirrorSplitTarget 判斷是否為鏡像副本（供 handleKill 使用）
func (m *luckyMirrorSplitManager) isMirrorSplitTarget(instanceID string) (bool, *mirrorSplitEntry) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if entry, ok := m.activeMirrors[instanceID]; ok {
		return true, entry
	}
	return false, nil
}

// removeMirrorEntry 移除鏡像副本記錄
func (m *luckyMirrorSplitManager) removeMirrorEntry(instanceID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.activeMirrors, instanceID)
}

// tryLuckyMirrorSplitFish 擊破 T208 後觸發鏡像分裂
func (g *Game) tryLuckyMirrorSplitFish(p *player.Player) {
	m := g.LuckyMirrorSplit
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
	m.personalCooldowns[p.ID] = now.Add(LuckyMirrorSplitPersonalCD)
	m.globalCooldownUntil = now.Add(LuckyMirrorSplitGlobalCD)
	m.mu.Unlock()

	// 選取場上隨機目標進行分裂
	g.mu.Lock()
	var candidates []*target.Target
	for _, t := range g.Targets {
		if !t.IsAlive {
			continue
		}
		// 排除 BOSS、Bonus 目標、鏡像副本本身
		if def, ok := data.Targets[t.DefID]; ok {
			if def.Type == data.TargetTypeBoss || def.Type == data.TargetTypeBonus {
				continue
			}
		}
		// 排除已是鏡像副本的目標
		if _, isMirror := m.activeMirrors[t.InstanceID]; isMirror {
			continue
		}
		candidates = append(candidates, t)
	}

	// 隨機選取最多 LuckyMirrorSplitCount 個
	shuffleTargets(candidates)
	if len(candidates) > LuckyMirrorSplitCount {
		candidates = candidates[:LuckyMirrorSplitCount]
	}

	if len(candidates) == 0 {
		g.mu.Unlock()
		return
	}

	// 建立鏡像副本
	type mirrorInfo struct {
		mirrorID   string
		origID     string
		defID      string
		origMult   float64
		mirrorX    float64
		mirrorY    float64
		mirrorHP   int
		mirrorMult float64
	}
	var mirrors []mirrorInfo
	expiresAt := now.Add(LuckyMirrorSplitDuration)

	for _, orig := range candidates {
		// 鏡像 X 座標（以場景中央 X=500 為軸）
		mirrorX := 2*LuckyMirrorSplitCenterX - orig.X
		// 邊界限制
		if mirrorX < 50 {
			mirrorX = 50
		}
		if mirrorX > 950 {
			mirrorX = 950
		}
		mirrorY := orig.Y

		mirrorHP := int(float64(orig.HP) * LuckyMirrorSplitHPRatio)
		if mirrorHP < 1 {
			mirrorHP = 1
		}
		mirrorMult := orig.Multiplier * LuckyMirrorSplitKillMult

		mirrorID := uuid.New().String()

		// 建立鏡像副本 Target
		mirrorTarget := &target.Target{
			InstanceID: mirrorID,
			DefID:      orig.DefID,
			Def:        orig.Def,
			HP:         mirrorHP,
			MaxHP:      mirrorHP,
			Multiplier: mirrorMult,
			X:          mirrorX,
			Y:          mirrorY,
			SpawnedAt:  now,
			IsAlive:    true,
			IsEcho:     false,
		}
		g.Targets[mirrorID] = mirrorTarget

		// 記錄鏡像副本
		entry := &mirrorSplitEntry{
			mirrorInstanceID: mirrorID,
			origInstanceID:   orig.InstanceID,
			origDefID:        orig.DefID,
			origMult:         orig.Multiplier,
			expiresAt:        expiresAt,
		}
		m.mu.Lock()
		m.activeMirrors[mirrorID] = entry
		m.mu.Unlock()

		mirrors = append(mirrors, mirrorInfo{
			mirrorID:   mirrorID,
			origID:     orig.InstanceID,
			defID:      orig.DefID,
			origMult:   orig.Multiplier,
			mirrorX:    mirrorX,
			mirrorY:    mirrorY,
			mirrorHP:   mirrorHP,
			mirrorMult: mirrorMult,
		})
	}
	g.mu.Unlock()

	splitCount := len(mirrors)
	log.Printf("[MirrorSplit] player=%s 鏡像分裂！生成 %d 個鏡像副本", p.ID, splitCount)

	// 建立廣播用的鏡像資訊
	type mirrorSpawnInfo struct {
		MirrorID   string  `json:"mirror_id"`
		OrigID     string  `json:"orig_id"`
		DefID      string  `json:"def_id"`
		MirrorX    float64 `json:"mirror_x"`
		MirrorY    float64 `json:"mirror_y"`
		MirrorHP   int     `json:"mirror_hp"`
		MirrorMult float64 `json:"mirror_mult"`
	}
	mirrorList := make([]mirrorSpawnInfo, 0, len(mirrors))
	for _, mi := range mirrors {
		mirrorList = append(mirrorList, mirrorSpawnInfo{
			MirrorID:   mi.mirrorID,
			OrigID:     mi.origID,
			DefID:      mi.defID,
			MirrorX:    mi.mirrorX,
			MirrorY:    mi.mirrorY,
			MirrorHP:   mi.mirrorHP,
			MirrorMult: mi.mirrorMult,
		})
	}

	// 個人訊息：分裂啟動
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyMirrorSplit,
		Payload: ws.LuckyMirrorSplitPayload{
			Event:       "mirror_split_start",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			SplitCount:  splitCount,
			DurationSec: int(LuckyMirrorSplitDuration.Seconds()),
			KillMult:    LuckyMirrorSplitKillMult,
			FadeMult:    LuckyMirrorSplitFadeMult,
			Mirrors:     mirrorList,
		},
	})

	// 全服廣播
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyMirrorSplit,
		Payload: ws.LuckyMirrorSplitPayload{
			Event:      "mirror_split_broadcast",
			PlayerName: p.DisplayName,
			SplitCount: splitCount,
			Mirrors:    mirrorList,
		},
	})

	// 全服公告
	g.Announce.Create(announce.EventLuckyMirrorSplit, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🪞 %s 觸發鏡像分裂！%d 個目標分裂成鏡像副本！", p.DisplayName, splitCount),
		"color":   "#8E44AD",
	})

	// 啟動消融計時 goroutine
	go g.runMirrorSplitFade(p, mirrors, expiresAt)
}

// notifyMirrorSplitKill 鏡像副本被擊破時的處理（由 handleKill 呼叫）
func (g *Game) notifyMirrorSplitKill(p *player.Player, mirrorInstanceID string, entry *mirrorSplitEntry) float64 {
	g.LuckyMirrorSplit.removeMirrorEntry(mirrorInstanceID)

	reward := int(float64(data.GetBetDef(p.BetLevel).BetCost) * entry.origMult * LuckyMirrorSplitKillMult)
	if reward < 1 {
		reward = 1
	}
	p.AddCoins(reward)

	log.Printf("[MirrorSplit] player=%s 擊破鏡像副本 %s，獎勵 %d", p.ID, mirrorInstanceID, reward)

	// 廣播擊破事件
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyMirrorSplit,
		Payload: ws.LuckyMirrorSplitPayload{
			Event:      "mirror_split_kill",
			PlayerName: p.DisplayName,
			InstanceID: mirrorInstanceID,
			Reward:     reward,
			KillMult:   LuckyMirrorSplitKillMult,
		},
	})

	return LuckyMirrorSplitKillMult
}

// runMirrorSplitFade 鏡像消融計時 goroutine
func (g *Game) runMirrorSplitFade(p *player.Player, mirrors interface{}, expiresAt time.Time) {
	remaining := time.Until(expiresAt)
	if remaining <= 0 {
		return
	}

	timer := time.NewTimer(remaining)
	defer timer.Stop()

	select {
	case <-timer.C:
		g.doMirrorSplitFade(p)
	case <-g.stopCh:
		return
	}
}

// doMirrorSplitFade 執行鏡像消融
func (g *Game) doMirrorSplitFade(p *player.Player) {
	m := g.LuckyMirrorSplit
	m.mu.Lock()

	// 收集所有仍存活的鏡像副本
	var fadeIDs []string
	for mirrorID := range m.activeMirrors {
		fadeIDs = append(fadeIDs, mirrorID)
	}
	// 清除所有記錄
	for _, id := range fadeIDs {
		delete(m.activeMirrors, id)
	}
	m.mu.Unlock()

	if len(fadeIDs) == 0 {
		return
	}

	// 消滅場上的鏡像副本
	g.mu.Lock()
	for _, id := range fadeIDs {
		delete(g.Targets, id)
	}
	g.mu.Unlock()

	// 計算全服共享消融獎勵
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

	totalReward := int(float64(avgBet) * LuckyMirrorSplitFadeMult * float64(len(fadeIDs)))
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

	log.Printf("[MirrorSplit] 鏡像消融！消融 %d 個副本，全服獎勵 %d", len(fadeIDs), totalReward)

	// 全服廣播消融
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyMirrorSplit,
		Payload: ws.LuckyMirrorSplitPayload{
			Event:       "mirror_split_fade",
			FadeCount:   len(fadeIDs),
			TotalReward: totalReward,
			FadeMult:    LuckyMirrorSplitFadeMult,
		},
	})

	// 個人結束通知
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyMirrorSplit,
		Payload: ws.LuckyMirrorSplitPayload{
			Event:       "mirror_split_end",
			FadeCount:   len(fadeIDs),
			TotalReward: totalReward,
		},
	})

	if len(fadeIDs) >= 2 {
		g.Announce.Create(announce.EventLuckyMirrorSplit, p.DisplayName, totalReward, map[string]string{
			"message": fmt.Sprintf("🪞 %s 鏡像消融！%d 個副本消融，全服獲得 %d 籌碼！",
				p.DisplayName, len(fadeIDs), totalReward),
			"color": "#6C3483",
		})
	}
}

// shuffleTargets 隨機打亂目標切片（Fisher-Yates）
func shuffleTargets(targets []*target.Target) {
	for i := len(targets) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		targets[i], targets[j] = targets[j], targets[i]
	}
}
