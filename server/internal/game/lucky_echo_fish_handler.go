// lucky_echo_fish_handler.go — 幸運回聲魚系統（DAY-233）
// 業界原創「回聲分身+層疊倍率」機制
//
// 設計：擊破 T191 後，玩家的下一次擊破會產生「回聲分身」：
//   - 在原位置複製一個相同目標（HP = 原 HP × 50%，倍率 ×1.5）
//   - 分身被擊破時再次產生回聲（最多 3 層）
//   - 每層倍率遞增：第1層 ×1.5 → 第2層 ×2.0 → 第3層 ×2.5
//   - 個人冷卻 18 秒
//
// 設計差異：
//   - 與分裂魚（DAY-224，一魚分三，同時生成）不同，回聲魚是「連鎖回聲」，
//     讓玩家有「打一個，再打一個，再打一個」的連鎖爽感
//   - 「層疊倍率遞增」讓玩家有「越打越值錢」的期待感
//   - 「最多 3 層」確保 RTP 平衡，不會無限連鎖
//   - 「HP 50%」讓分身更容易擊破，讓玩家有「快速連殺」的節奏感
//   - 個人機制（不是全服），讓每個玩家都有自己的回聲節奏
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/game/target"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	LuckyEchoPersonalCD  = 18 * time.Second // 個人冷卻
	LuckyEchoMaxLayers   = 3                // 最多回聲層數
	LuckyEchoHPRatio     = 0.50             // 分身 HP 比例
	LuckyEchoLayer1Mult  = 1.5              // 第1層倍率加成
	LuckyEchoLayer2Mult  = 2.0              // 第2層倍率加成
	LuckyEchoLayer3Mult  = 2.5              // 第3層倍率加成
	LuckyEchoSpreadRange = 80.0             // 分身散佈範圍（px）
	LuckyEchoLifetime    = 10.0             // 分身存活時間（秒）
)

// echoSession 回聲模式 session（每個玩家一個）
type echoSession struct {
	playerID   string
	expiresAt  time.Time // session 過期時間（觸發後 15 秒內有效）
}

// echoEntry 回聲分身記錄
type echoEntry struct {
	instanceID string  // 分身目標的 instanceID
	layer      int     // 回聲層數（1-3）
	ownerID    string  // 觸發玩家 ID
}

// luckyEchoFishManager 幸運回聲魚管理器
type luckyEchoFishManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 回聲模式 session（playerID → echoSession）
	// 擊破 T191 後進入回聲模式，下一次擊破觸發回聲
	activeSessions map[string]*echoSession

	// 回聲分身記錄（instanceID → echoEntry）
	echoTargets map[string]*echoEntry
}

func newLuckyEchoFishManager() *luckyEchoFishManager {
	return &luckyEchoFishManager{
		personalCooldowns: make(map[string]time.Time),
		activeSessions:    make(map[string]*echoSession),
		echoTargets:       make(map[string]*echoEntry),
	}
}

// isLuckyEchoFish 判斷是否為幸運回聲魚
func isLuckyEchoFish(defID string) bool {
	return defID == "T191"
}

// isEchoTarget 判斷是否為回聲分身目標
func (g *Game) isEchoTarget(instanceID string) (bool, *echoEntry) {
	mgr := g.LuckyEchoFish
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	entry, ok := mgr.echoTargets[instanceID]
	return ok, entry
}

// removeEchoEntry 移除回聲分身記錄（分身被擊破後呼叫）
func (g *Game) removeEchoEntry(instanceID string) {
	mgr := g.LuckyEchoFish
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	delete(mgr.echoTargets, instanceID)
}

// getLuckyEchoKillMult 取得回聲分身擊破倍率加成（供 handleKill 使用）
// 回傳 (倍率加成, 是否為回聲分身, 層數)
func (g *Game) getLuckyEchoKillMult(instanceID string) (float64, bool, int) {
	isEcho, entry := g.isEchoTarget(instanceID)
	if !isEcho {
		return 1.0, false, 0
	}
	switch entry.layer {
	case 1:
		return LuckyEchoLayer1Mult, true, 1
	case 2:
		return LuckyEchoLayer2Mult, true, 2
	case 3:
		return LuckyEchoLayer3Mult, true, 3
	default:
		return 1.0, true, entry.layer
	}
}

// isEchoModeActive 判斷玩家是否在回聲模式中（供 handleKill 使用）
func (g *Game) isEchoModeActive(playerID string) bool {
	mgr := g.LuckyEchoFish
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	sess, ok := mgr.activeSessions[playerID]
	if !ok {
		return false
	}
	if time.Now().After(sess.expiresAt) {
		delete(mgr.activeSessions, playerID)
		return false
	}
	return true
}

// consumeEchoSession 消耗回聲 session（觸發回聲後移除）
func (g *Game) consumeEchoSession(playerID string) bool {
	mgr := g.LuckyEchoFish
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	sess, ok := mgr.activeSessions[playerID]
	if !ok {
		return false
	}
	if time.Now().After(sess.expiresAt) {
		delete(mgr.activeSessions, playerID)
		return false
	}
	delete(mgr.activeSessions, playerID)
	return true
}

// tryLuckyEchoFish 擊破 T191 後觸發回聲模式（供 handleKill 使用）
func (g *Game) tryLuckyEchoFish(p *player.Player) {
	mgr := g.LuckyEchoFish
	mgr.mu.Lock()

	// 個人冷卻檢查
	if cd, ok := mgr.personalCooldowns[p.ID]; ok && time.Now().Before(cd) {
		mgr.mu.Unlock()
		return
	}

	// 設定個人冷卻
	mgr.personalCooldowns[p.ID] = time.Now().Add(LuckyEchoPersonalCD)

	// 建立回聲 session（15 秒內有效）
	mgr.activeSessions[p.ID] = &echoSession{
		playerID:  p.ID,
		expiresAt: time.Now().Add(15 * time.Second),
	}
	mgr.mu.Unlock()

	log.Printf("[LuckyEcho] player=%s entered echo mode (next kill will spawn echo)", p.ID)

	// 個人訊息：回聲模式啟動
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyEchoFish,
		Payload: ws.LuckyEchoFishPayload{
			Event:      "echo_ready",
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			Layer:      0,
		},
	})

	// 全服廣播（小橫幅）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyEchoFish,
		Payload: ws.LuckyEchoFishPayload{
			Event:      "echo_broadcast",
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyEchoFish, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🔊 %s 觸發回聲模式！下次擊破將產生回聲分身！", p.DisplayName),
		"color":   "#9B59B6",
	})
	g.broadcastAnnouncement(ann)
}

// notifyEchoKill 擊破目標後觸發回聲分身（供 handleKill 使用）
// 當玩家在回聲模式中擊破任何目標時呼叫
// 回傳是否成功觸發回聲
func (g *Game) notifyEchoKill(p *player.Player, killedTarget *target.Target, killedInstanceID string) bool {
	// 消耗回聲 session
	if !g.consumeEchoSession(p.ID) {
		return false
	}

	// 生成回聲分身（第1層）
	return g.spawnEchoTarget(p, killedTarget, killedInstanceID, 1)
}

// notifyEchoTargetKill 擊破回聲分身後觸發下一層回聲（供 handleKill 使用）
// 回傳是否成功觸發下一層回聲
func (g *Game) notifyEchoTargetKill(p *player.Player, killedTarget *target.Target, killedInstanceID string, layer int) bool {
	if layer >= LuckyEchoMaxLayers {
		// 已達最大層數，不再產生回聲
		log.Printf("[LuckyEcho] player=%s killed echo layer=%d (max reached)", p.ID, layer)
		return false
	}

	// 生成下一層回聲分身
	return g.spawnEchoTarget(p, killedTarget, killedInstanceID, layer+1)
}

// spawnEchoTarget 生成回聲分身目標
func (g *Game) spawnEchoTarget(p *player.Player, originalTarget *target.Target, originalInstanceID string, layer int) bool {
	g.mu.Lock()

	// 計算分身位置（在原位置附近隨機散佈）
	offsetX := (rand.Float64()*2 - 1) * LuckyEchoSpreadRange
	offsetY := (rand.Float64()*2 - 1) * LuckyEchoSpreadRange
	echoX := originalTarget.X + offsetX
	echoY := originalTarget.Y + offsetY

	// 邊界限制
	if echoX < 80 {
		echoX = 80
	}
	if echoX > 920 {
		echoX = 920
	}
	if echoY < 60 {
		echoY = 60
	}
	if echoY > 540 {
		echoY = 540
	}

	// 計算分身 HP（原 HP × 50%，最少 1）
	echoHP := int(float64(originalTarget.MaxHP) * LuckyEchoHPRatio)
	if echoHP < 1 {
		echoHP = 1
	}

	// 計算分身倍率（原倍率 × 層數加成）
	var layerMult float64
	switch layer {
	case 1:
		layerMult = LuckyEchoLayer1Mult
	case 2:
		layerMult = LuckyEchoLayer2Mult
	case 3:
		layerMult = LuckyEchoLayer3Mult
	default:
		layerMult = LuckyEchoLayer1Mult
	}
	echoMult := originalTarget.Multiplier * layerMult

	// 生成新的 instanceID
	echoInstanceID := fmt.Sprintf("echo_%s_L%d_%d", originalInstanceID, layer, time.Now().UnixNano())

	// 建立回聲分身目標（使用原目標的 DefID）
	echoTarget := &target.Target{
		InstanceID:  echoInstanceID,
		DefID:       originalTarget.DefID,
		Def:         originalTarget.Def,
		HP:          echoHP,
		MaxHP:       echoHP,
		Multiplier:  echoMult,
		X:           echoX,
		Y:           echoY,
		SpawnedAt:   time.Now(),
		IsAlive:     true,
		IsEcho:      true,
		EchoLayer:   layer,
	}

	g.Targets[echoInstanceID] = echoTarget
	g.mu.Unlock()

	// 記錄回聲分身
	mgr := g.LuckyEchoFish
	mgr.mu.Lock()
	mgr.echoTargets[echoInstanceID] = &echoEntry{
		instanceID: echoInstanceID,
		layer:      layer,
		ownerID:    p.ID,
	}
	mgr.mu.Unlock()

	log.Printf("[LuckyEcho] player=%s spawned echo layer=%d instanceID=%s mult=%.2f hp=%d",
		p.ID, layer, echoInstanceID, echoMult, echoHP)

	// 廣播回聲分身生成
	var multLabel string
	switch layer {
	case 1:
		multLabel = fmt.Sprintf("×%.1f", LuckyEchoLayer1Mult)
	case 2:
		multLabel = fmt.Sprintf("×%.1f", LuckyEchoLayer2Mult)
	case 3:
		multLabel = fmt.Sprintf("×%.1f", LuckyEchoLayer3Mult)
	}

	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyEchoFish,
		Payload: ws.LuckyEchoFishPayload{
			Event:          "echo_spawn",
			PlayerID:       p.ID,
			PlayerName:     p.DisplayName,
			Layer:          layer,
			EchoInstanceID: echoInstanceID,
			OriginalID:     originalInstanceID,
			EchoX:          echoX,
			EchoY:          echoY,
			EchoHP:         echoHP,
			EchoMult:       echoMult,
			MultLabel:      multLabel,
		},
	})

	// 個人訊息：回聲分身生成提示
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyEchoFish,
		Payload: ws.LuckyEchoFishPayload{
			Event:     "echo_spawn_personal",
			Layer:     layer,
			MultLabel: multLabel,
		},
	})

	return true
}

// notifyEchoTargetExpire 回聲分身逃跑/消失時清除記錄（供 gameLoop 使用）
func (g *Game) notifyEchoTargetExpire(instanceID string) {
	mgr := g.LuckyEchoFish
	mgr.mu.Lock()
	_, wasEcho := mgr.echoTargets[instanceID]
	delete(mgr.echoTargets, instanceID)
	mgr.mu.Unlock()

	if wasEcho {
		log.Printf("[LuckyEcho] echo target expired: instanceID=%s", instanceID)
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyEchoFish,
			Payload: ws.LuckyEchoFishPayload{
				Event:          "echo_expire",
				EchoInstanceID: instanceID,
			},
		})
	}
}
