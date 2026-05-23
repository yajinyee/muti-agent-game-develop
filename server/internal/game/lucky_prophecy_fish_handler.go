// lucky_prophecy_fish_handler.go — 幸運預言魚系統（DAY-243）
// 業界原創「預言指定目標」機制
//
// 設計：擊破 T201 後，Server 隨機「預言」場上 1 個目標（標記持續 12 秒）：
//   - 玩家在 12 秒內擊破預言目標 → 獲得 ×3.5 倍率加成（「預言成真」）
//   - 若預言目標在 12 秒內自然消失 → 自動「預言轉移」到下一個目標（最多轉移 2 次）
//   - 若 12 秒後仍未擊破 → 「預言失敗」，全場 HP -20%（安慰獎）
//   - 個人冷卻 20 秒；全服廣播讓所有玩家看到「有人的預言目標是哪條魚」
//
// 設計差異：
//   - 與時間炸彈魚（T193，倒數計時+提前引爆）不同，預言魚是「指定目標」，讓玩家有「要集中火力打那條魚」的聚焦感
//   - 「預言轉移」讓玩家不會因為目標消失而完全失去機會，降低挫敗感
//   - 「預言失敗 HP -20%」讓玩家有「要趕快打，不然全場魚都受傷」的緊迫感
//   - 全服廣播讓其他玩家也知道「有人在追那條魚」，製造社交感
//   - ×3.5 倍率是目前個人指定目標類最高倍率，讓玩家有「值得集中火力」的動機
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
	LuckyProphecyFishPersonalCD  = 20 * time.Second // 個人冷卻
	LuckyProphecyFishDuration    = 12 * time.Second // 預言標記持續時間
	LuckyProphecyFishKillMult    = 3.5              // 預言成真倍率
	LuckyProphecyFishMaxTransfer = 2                // 最多轉移次數
	LuckyProphecyFishFailHPLoss  = 0.20             // 預言失敗 HP 損失比例
)

// prophecyEntry 單個玩家的預言狀態
type prophecyEntry struct {
	playerID      string
	targetID      string // 當前預言目標的 instanceID
	targetDefID   string // 目標 defID（用於廣播）
	targetX       float64
	targetY       float64
	transferCount int       // 已轉移次數
	expiresAt     time.Time // 預言到期時間
}

// luckyProphecyFishManager 幸運預言魚管理器
type luckyProphecyFishManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 當前預言狀態（playerID → prophecyEntry）
	activeProphecies map[string]*prophecyEntry
}

func newLuckyProphecyFishManager() *luckyProphecyFishManager {
	return &luckyProphecyFishManager{
		personalCooldowns: make(map[string]time.Time),
		activeProphecies:  make(map[string]*prophecyEntry),
	}
}

// isLuckyProphecyFish 判斷是否為幸運預言魚
func isLuckyProphecyFish(defID string) bool {
	return defID == "T201"
}

// isProphecyTarget 判斷某個目標是否為某玩家的預言目標（供 handleKill 使用）
func (g *Game) isProphecyTarget(playerID, instanceID string) bool {
	mgr := g.LuckyProphecyFish
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	entry, ok := mgr.activeProphecies[playerID]
	if !ok {
		return false
	}
	if time.Now().After(entry.expiresAt) {
		delete(mgr.activeProphecies, playerID)
		return false
	}
	return entry.targetID == instanceID
}

// getLuckyProphecyMult 取得預言成真倍率（供 handleKill 使用）
// 若玩家有預言且命中預言目標，回傳 ×3.5；否則回傳 1.0
func (g *Game) getLuckyProphecyMult(playerID, instanceID string) float64 {
	if g.isProphecyTarget(playerID, instanceID) {
		return LuckyProphecyFishKillMult
	}
	return 1.0
}

// removeProphecyEntry 移除預言（預言成真後呼叫）
func (g *Game) removeProphecyEntry(playerID string) {
	mgr := g.LuckyProphecyFish
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	delete(mgr.activeProphecies, playerID)
}

// notifyProphecyKill 玩家擊破預言目標時呼叫（廣播預言成真）
func (g *Game) notifyProphecyKill(p *player.Player, instanceID string, reward int) {
	mgr := g.LuckyProphecyFish
	mgr.mu.Lock()
	entry, ok := mgr.activeProphecies[p.ID]
	if !ok || entry.targetID != instanceID {
		mgr.mu.Unlock()
		return
	}
	delete(mgr.activeProphecies, p.ID)
	mgr.mu.Unlock()

	log.Printf("[LuckyProphecy] player=%s prophecy fulfilled! target=%s reward=%d", p.ID, instanceID, reward)

	// 個人訊息：預言成真
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyProphecyFish,
		Payload: ws.LuckyProphecyFishPayload{
			Event:      "prophecy_fulfilled",
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			TargetID:   instanceID,
			Reward:     reward,
			KillMult:   LuckyProphecyFishKillMult,
		},
	})

	// 全服廣播：預言成真
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyProphecyFish,
		Payload: ws.LuckyProphecyFishPayload{
			Event:      "prophecy_broadcast_fulfilled",
			PlayerName: p.DisplayName,
			TargetID:   instanceID,
			Reward:     reward,
			KillMult:   LuckyProphecyFishKillMult,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyProphecyFish, p.DisplayName, reward, map[string]string{
		"message": fmt.Sprintf("🔮 %s 預言成真！×%.1f 倍率！獲得 %d 金幣！",
			p.DisplayName, LuckyProphecyFishKillMult, reward),
		"color": "#9B59B6",
	})
	g.broadcastAnnouncement(ann)
}

// notifyProphecyTargetGone 預言目標消失時呼叫（嘗試轉移）
// 由 updateNormalPlay 的目標清理邏輯呼叫
func (g *Game) notifyProphecyTargetGone(instanceID string) {
	mgr := g.LuckyProphecyFish
	mgr.mu.Lock()

	// 找到所有以此目標為預言目標的玩家
	var affectedPlayers []*prophecyEntry
	for _, entry := range mgr.activeProphecies {
		if entry.targetID == instanceID && time.Now().Before(entry.expiresAt) {
			affectedPlayers = append(affectedPlayers, entry)
		}
	}
	mgr.mu.Unlock()

	for _, entry := range affectedPlayers {
		go g.tryTransferProphecy(entry.playerID)
	}
}

// tryTransferProphecy 嘗試將預言轉移到新目標
func (g *Game) tryTransferProphecy(playerID string) {
	mgr := g.LuckyProphecyFish
	mgr.mu.Lock()

	entry, ok := mgr.activeProphecies[playerID]
	if !ok || time.Now().After(entry.expiresAt) {
		delete(mgr.activeProphecies, playerID)
		mgr.mu.Unlock()
		return
	}

	// 檢查轉移次數
	if entry.transferCount >= LuckyProphecyFishMaxTransfer {
		// 已達最大轉移次數，預言失敗
		delete(mgr.activeProphecies, playerID)
		mgr.mu.Unlock()
		go g.doProphecyFail(playerID)
		return
	}
	mgr.mu.Unlock()

	// 選擇新的預言目標
	g.mu.RLock()
	var candidates []string
	for id, t := range g.Targets {
		if t.HP > 0 && t.DefID != "B001" && t.DefID != "T201" {
			candidates = append(candidates, id)
		}
	}
	g.mu.RUnlock()

	if len(candidates) == 0 {
		// 場上沒有目標，預言失敗
		mgr.mu.Lock()
		delete(mgr.activeProphecies, playerID)
		mgr.mu.Unlock()
		go g.doProphecyFail(playerID)
		return
	}

	// 隨機選一個新目標
	newTargetID := candidates[rand.Intn(len(candidates))]

	g.mu.RLock()
	newTarget := g.Targets[newTargetID]
	var newDefID string
	var newX, newY float64
	if newTarget != nil {
		newDefID = newTarget.DefID
		newX = newTarget.X
		newY = newTarget.Y
	}
	g.mu.RUnlock()

	if newTarget == nil {
		mgr.mu.Lock()
		delete(mgr.activeProphecies, playerID)
		mgr.mu.Unlock()
		go g.doProphecyFail(playerID)
		return
	}

	// 更新預言目標
	mgr.mu.Lock()
	entry, ok = mgr.activeProphecies[playerID]
	if !ok {
		mgr.mu.Unlock()
		return
	}
	oldTargetID := entry.targetID
	entry.targetID = newTargetID
	entry.targetDefID = newDefID
	entry.targetX = newX
	entry.targetY = newY
	entry.transferCount++
	transferCount := entry.transferCount
	mgr.mu.Unlock()

	log.Printf("[LuckyProphecy] player=%s prophecy transferred from=%s to=%s (transfer #%d)",
		playerID, oldTargetID, newTargetID, transferCount)

	// 取得玩家資訊
	g.mu.RLock()
	p := g.Players[playerID]
	g.mu.RUnlock()
	if p == nil {
		return
	}

	// 個人訊息：預言轉移
	g.Hub.Send(playerID, &ws.Message{
		Type: ws.MsgLuckyProphecyFish,
		Payload: ws.LuckyProphecyFishPayload{
			Event:         "prophecy_transfer",
			PlayerID:      playerID,
			PlayerName:    p.DisplayName,
			TargetID:      newTargetID,
			TargetDefID:   newDefID,
			X:             newX,
			Y:             newY,
			TransferCount: transferCount,
			KillMult:      LuckyProphecyFishKillMult,
		},
	})

	// 全服廣播：預言轉移
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyProphecyFish,
		Payload: ws.LuckyProphecyFishPayload{
			Event:         "prophecy_broadcast_transfer",
			PlayerName:    p.DisplayName,
			TargetID:      newTargetID,
			TargetDefID:   newDefID,
			X:             newX,
			Y:             newY,
			TransferCount: transferCount,
		},
	})
}

// doProphecyFail 預言失敗：全場 HP -20%
func (g *Game) doProphecyFail(playerID string) {
	g.mu.RLock()
	p := g.Players[playerID]
	g.mu.RUnlock()

	playerName := "某玩家"
	if p != nil {
		playerName = p.DisplayName
	}

	// 全場 HP -20%（保留最少 1）
	g.mu.Lock()
	affectedCount := 0
	for _, t := range g.Targets {
		if t.HP > 0 && t.DefID != "B001" {
			loss := int(float64(t.MaxHP) * LuckyProphecyFishFailHPLoss)
			if loss < 1 {
				loss = 1
			}
			t.HP -= loss
			if t.HP < 1 {
				t.HP = 1
			}
			affectedCount++
		}
	}
	g.mu.Unlock()

	log.Printf("[LuckyProphecy] player=%s prophecy failed, HP-20%% applied to %d targets", playerID, affectedCount)

	// 廣播預言失敗
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyProphecyFish,
		Payload: ws.LuckyProphecyFishPayload{
			Event:         "prophecy_fail",
			PlayerName:    playerName,
			AffectedCount: affectedCount,
			HPLossPct:     LuckyProphecyFishFailHPLoss,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyProphecyFish, playerName, 0, map[string]string{
		"message": fmt.Sprintf("🔮 %s 預言失敗！全場目標 HP -20%%！",
			playerName),
		"color": "#7F8C8D",
	})
	g.broadcastAnnouncement(ann)
}

// tryLuckyProphecyFish 擊破 T201 後觸發預言
func (g *Game) tryLuckyProphecyFish(p *player.Player) {
	mgr := g.LuckyProphecyFish
	mgr.mu.Lock()

	// 個人冷卻檢查
	if cd, ok := mgr.personalCooldowns[p.ID]; ok && time.Now().Before(cd) {
		mgr.mu.Unlock()
		return
	}

	// 設定個人冷卻
	mgr.personalCooldowns[p.ID] = time.Now().Add(LuckyProphecyFishPersonalCD)
	mgr.mu.Unlock()

	// 選擇預言目標（場上隨機一個非 BOSS 目標）
	g.mu.RLock()
	var candidates []string
	for id, t := range g.Targets {
		if t.HP > 0 && t.DefID != "B001" && t.DefID != "T201" {
			candidates = append(candidates, id)
		}
	}
	g.mu.RUnlock()

	if len(candidates) == 0 {
		log.Printf("[LuckyProphecy] player=%s no valid targets for prophecy", p.ID)
		return
	}

	// 隨機選一個目標
	targetID := candidates[rand.Intn(len(candidates))]

	g.mu.RLock()
	target := g.Targets[targetID]
	var targetDefID string
	var targetX, targetY float64
	if target != nil {
		targetDefID = target.DefID
		targetX = target.X
		targetY = target.Y
	}
	g.mu.RUnlock()

	if target == nil {
		return
	}

	// 建立預言
	expiresAt := time.Now().Add(LuckyProphecyFishDuration)
	mgr.mu.Lock()
	mgr.activeProphecies[p.ID] = &prophecyEntry{
		playerID:      p.ID,
		targetID:      targetID,
		targetDefID:   targetDefID,
		targetX:       targetX,
		targetY:       targetY,
		transferCount: 0,
		expiresAt:     expiresAt,
	}
	mgr.mu.Unlock()

	log.Printf("[LuckyProphecy] player=%s prophecy set on target=%s (%s) for %v",
		p.ID, targetID, targetDefID, LuckyProphecyFishDuration)

	// 個人訊息：預言開始
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyProphecyFish,
		Payload: ws.LuckyProphecyFishPayload{
			Event:       "prophecy_start",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			TargetID:    targetID,
			TargetDefID: targetDefID,
			X:           targetX,
			Y:           targetY,
			DurationSec: int(LuckyProphecyFishDuration.Seconds()),
			KillMult:    LuckyProphecyFishKillMult,
		},
	})

	// 全服廣播：預言開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyProphecyFish,
		Payload: ws.LuckyProphecyFishPayload{
			Event:       "prophecy_broadcast",
			PlayerName:  p.DisplayName,
			TargetID:    targetID,
			TargetDefID: targetDefID,
			X:           targetX,
			Y:           targetY,
			DurationSec: int(LuckyProphecyFishDuration.Seconds()),
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyProphecyFish, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🔮 %s 觸發預言！%d 秒內擊破指定目標獲得 ×%.1f 倍率！",
			p.DisplayName, int(LuckyProphecyFishDuration.Seconds()), LuckyProphecyFishKillMult),
		"color": "#9B59B6",
	})
	g.broadcastAnnouncement(ann)

	// 啟動計時器，到期後觸發預言失敗
	go func() {
		time.Sleep(LuckyProphecyFishDuration)

		mgr.mu.Lock()
		entry, ok := mgr.activeProphecies[p.ID]
		if !ok {
			// 預言已被消耗（成真）
			mgr.mu.Unlock()
			return
		}
		// 預言仍存在，表示未成真 → 失敗
		delete(mgr.activeProphecies, p.ID)
		mgr.mu.Unlock()

		log.Printf("[LuckyProphecy] player=%s prophecy expired (target=%s), triggering fail", p.ID, entry.targetID)
		g.doProphecyFail(p.ID)
	}()
}
