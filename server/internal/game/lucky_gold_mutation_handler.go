// lucky_gold_mutation_handler.go — 幸運黃金突變魚系統（DAY-281）
// 業界依據：Fisch Roblox「Lucky Gold Mutation（6.14×）」+ Fishing Legend 2025「高倍率目標」機制
// 業界原創「黃金突變+全場感染+突變連鎖」機制
//
// 設計：擊破 T239 後，觸發「黃金突變」：
//   - Server 隨機選場上 2-4 個目標，讓它們「突變為黃金版本」
//     （HP 降低 50%，擊破倍率 ×3.0）
//   - 突變目標有 30% 機率在被擊破時「感染」相鄰目標（也突變，HP 降低 50%，擊破倍率 ×2.0）
//   - 突變持續 15 秒（超時後突變消失，目標恢復正常）
//   - 全服廣播突變目標位置和數量
//   - 個人冷卻 20 秒；全服冷卻 35 秒
//
// 設計差異：
//   - 與品質突變（T230，個人下一次擊破品質抽取）不同，黃金突變是「場上特定目標突變」
//     讓玩家有「快去打那幾條黃金魚！」的緊迫感
//   - 「感染相鄰目標（30%）」讓玩家有「打一條黃金魚，旁邊的也可能突變」的驚喜感
//   - 「突變目標 HP -50%」讓黃金魚更容易打，製造「黃金魚比普通魚更值得打」的策略感
//   - 「突變持續 15 秒」讓玩家有「要趁 15 秒內打完所有黃金魚」的緊迫感
//   - 「全服廣播突變目標位置」讓所有玩家看到「哪幾條魚突變了」，製造全服搶打感
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
	LuckyGoldMutationPersonalCD    = 20 * time.Second // 個人冷卻
	LuckyGoldMutationGlobalCD      = 35 * time.Second // 全服冷卻
	LuckyGoldMutationDuration      = 15 * time.Second // 突變持續時間
	LuckyGoldMutationMinTargets    = 2                // 最少突變目標數
	LuckyGoldMutationMaxTargets    = 4                // 最多突變目標數
	LuckyGoldMutationKillMult      = 3.0              // 突變目標擊破倍率
	LuckyGoldMutationInfectMult    = 2.0              // 感染目標擊破倍率
	LuckyGoldMutationInfectChance  = 0.30             // 感染機率（30%）
	LuckyGoldMutationHPReduction   = 0.50             // HP 降低比例（-50%）
)

// goldMutatedTarget 黃金突變目標
type goldMutatedTarget struct {
	instanceID string
	defID      string
	name       string
	mult       float64   // 擊破倍率（3.0 或 2.0）
	isInfected bool      // 是否為感染突變（非原始突變）
	expiresAt  time.Time
}

// luckyGoldMutationManager 幸運黃金突變魚管理器
type luckyGoldMutationManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 當前突變目標（instanceID → goldMutatedTarget）
	mutatedTargets map[string]*goldMutatedTarget
}

func newLuckyGoldMutationManager() *luckyGoldMutationManager {
	return &luckyGoldMutationManager{
		personalCooldowns: make(map[string]time.Time),
		mutatedTargets:    make(map[string]*goldMutatedTarget),
	}
}

// isLuckyGoldMutationFish 判斷是否為幸運黃金突變魚
func isLuckyGoldMutationFish(defID string) bool {
	return defID == "T239"
}

// isGoldMutated 判斷目標是否為黃金突變目標
func (m *luckyGoldMutationManager) isGoldMutated(instanceID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if mt, ok := m.mutatedTargets[instanceID]; ok {
		return time.Now().Before(mt.expiresAt)
	}
	return false
}

// getGoldMutationMult 取得黃金突變倍率（供 handleKill 使用）
func (m *luckyGoldMutationManager) getGoldMutationMult(instanceID string) float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if mt, ok := m.mutatedTargets[instanceID]; ok {
		if time.Now().Before(mt.expiresAt) {
			return mt.mult
		}
	}
	return 1.0
}

// removeGoldMutation 移除突變記錄（目標被擊破後）
func (m *luckyGoldMutationManager) removeGoldMutation(instanceID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.mutatedTargets, instanceID)
}

// tryLuckyGoldMutationFish 擊破 T239 後觸發黃金突變（供 handleKill 使用）
func (g *Game) tryLuckyGoldMutationFish(p *player.Player) {
	mgr := g.LuckyGoldMutation
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
	mgr.personalCooldowns[p.ID] = now.Add(LuckyGoldMutationPersonalCD)
	mgr.globalCooldownUntil = now.Add(LuckyGoldMutationGlobalCD)
	mgr.mu.Unlock()

	// 選取 2-4 個目標
	g.mu.RLock()
	var candidates []string
	var candidateNames []string
	for iid, t := range g.Targets {
		if t.HP > 0 && !isLuckyGoldMutationFish(t.DefID) {
			candidates = append(candidates, iid)
			candidateNames = append(candidateNames, t.Def.Name)
		}
	}
	g.mu.RUnlock()

	if len(candidates) < 1 {
		log.Printf("[GoldMutation] not enough targets (%d)", len(candidates))
		return
	}

	// 隨機選取
	count := LuckyGoldMutationMinTargets + rand.Intn(LuckyGoldMutationMaxTargets-LuckyGoldMutationMinTargets+1)
	if count > len(candidates) {
		count = len(candidates)
	}
	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
		candidateNames[i], candidateNames[j] = candidateNames[j], candidateNames[i]
	})
	selectedIIDs := candidates[:count]
	selectedNames := candidateNames[:count]

	expiresAt := now.Add(LuckyGoldMutationDuration)

	// 對選中目標施加突變（HP -50%）
	mgr.mu.Lock()
	for i, iid := range selectedIIDs {
		mgr.mutatedTargets[iid] = &goldMutatedTarget{
			instanceID: iid,
			defID:      "",
			name:       selectedNames[i],
			mult:       LuckyGoldMutationKillMult,
			isInfected: false,
			expiresAt:  expiresAt,
		}
	}
	mgr.mu.Unlock()

	// 對選中目標施加 HP -50%
	g.mu.Lock()
	for _, iid := range selectedIIDs {
		if t, ok := g.Targets[iid]; ok && t.HP > 0 {
			reduction := int(float64(t.HP) * LuckyGoldMutationHPReduction)
			if reduction < 1 {
				reduction = 1
			}
			t.HP -= reduction
			if t.HP < 1 {
				t.HP = 1
			}
		}
	}
	g.mu.Unlock()

	log.Printf("[GoldMutation] player=%s triggered gold mutation, targets=%v", p.ID, selectedIIDs)

	// 全服廣播：黃金突變觸發
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyGoldMutation,
		Payload: ws.LuckyGoldMutationPayload{
			Event:       "mutation_start",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			TargetIIDs:  selectedIIDs,
			TargetNames: selectedNames,
			KillMult:    LuckyGoldMutationKillMult,
			Duration:    int(LuckyGoldMutationDuration.Seconds()),
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyGoldMutation, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("✨ %s 觸發黃金突變！%d 個目標變成黃金魚！HP -50%%，擊破得 ×%.1f！15 秒內快打！",
			p.DisplayName, count, LuckyGoldMutationKillMult),
		"color": "#FFD700",
	})
	g.broadcastAnnouncement(ann)

	// 啟動突變計時器
	go g.runGoldMutationTimer(selectedIIDs, expiresAt)
}

// notifyGoldMutationKill 黃金突變目標被擊破時的感染處理（供 handleKill 使用）
func (g *Game) notifyGoldMutationKill(p *player.Player, killedIID string, killMult float64) {
	mgr := g.LuckyGoldMutation

	// 移除突變記錄
	mgr.removeGoldMutation(killedIID)

	log.Printf("[GoldMutation] target %s killed by player=%s mult=x%.1f", killedIID, p.ID, killMult)

	// 廣播擊破通知
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyGoldMutation,
		Payload: ws.LuckyGoldMutationPayload{
			Event:      "mutation_kill",
			PlayerName: p.DisplayName,
			KilledIID:  killedIID,
			KillMult:   killMult,
		},
	})

	// 30% 機率感染相鄰目標
	if rand.Float64() < LuckyGoldMutationInfectChance {
		go g.tryGoldMutationInfect(p, killedIID)
	}
}

// tryGoldMutationInfect 嘗試感染相鄰目標
func (g *Game) tryGoldMutationInfect(p *player.Player, sourceIID string) {
	// 找一個未突變的目標
	g.mu.RLock()
	var candidates []string
	var candidateNames []string
	for iid, t := range g.Targets {
		if t.HP > 0 && !isLuckyGoldMutationFish(t.DefID) && !g.LuckyGoldMutation.isGoldMutated(iid) {
			candidates = append(candidates, iid)
			candidateNames = append(candidateNames, t.Def.Name)
		}
	}
	g.mu.RUnlock()

	if len(candidates) == 0 {
		return
	}

	// 隨機選一個感染
	idx := rand.Intn(len(candidates))
	infectIID := candidates[idx]
	infectName := candidateNames[idx]
	expiresAt := time.Now().Add(LuckyGoldMutationDuration / 2) // 感染突變持續 7.5 秒

	g.LuckyGoldMutation.mu.Lock()
	g.LuckyGoldMutation.mutatedTargets[infectIID] = &goldMutatedTarget{
		instanceID: infectIID,
		name:       infectName,
		mult:       LuckyGoldMutationInfectMult,
		isInfected: true,
		expiresAt:  expiresAt,
	}
	g.LuckyGoldMutation.mu.Unlock()

	// HP -50%
	g.mu.Lock()
	if t, ok := g.Targets[infectIID]; ok && t.HP > 0 {
		reduction := int(float64(t.HP) * LuckyGoldMutationHPReduction)
		if reduction < 1 {
			reduction = 1
		}
		t.HP -= reduction
		if t.HP < 1 {
			t.HP = 1
		}
	}
	g.mu.Unlock()

	log.Printf("[GoldMutation] INFECT! source=%s infected=%s mult=x%.1f", sourceIID, infectIID, LuckyGoldMutationInfectMult)

	// 全服廣播：感染突變
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyGoldMutation,
		Payload: ws.LuckyGoldMutationPayload{
			Event:      "mutation_infect",
			PlayerName: p.DisplayName,
			TargetIIDs: []string{infectIID},
			TargetNames: []string{infectName},
			KillMult:   LuckyGoldMutationInfectMult,
		},
	})
}

// runGoldMutationTimer 黃金突變計時器（goroutine）
func (g *Game) runGoldMutationTimer(targetIIDs []string, expiresAt time.Time) {
	remaining := time.Until(expiresAt)
	if remaining <= 0 {
		return
	}

	select {
	case <-time.After(remaining):
	case <-g.stopCh:
		return
	}

	// 清除所有突變記錄
	g.LuckyGoldMutation.mu.Lock()
	for _, iid := range targetIIDs {
		delete(g.LuckyGoldMutation.mutatedTargets, iid)
	}
	g.LuckyGoldMutation.mu.Unlock()

	log.Printf("[GoldMutation] mutation expired for %d targets", len(targetIIDs))

	// 全服廣播：突變消失
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyGoldMutation,
		Payload: ws.LuckyGoldMutationPayload{
			Event: "mutation_expire",
		},
	})
}
