// lucky_time_rift_handler.go — 幸運時間裂縫魚系統（DAY-278）
// 業界原創「時間裂縫+最高倍率重現+裂縫複製體」機制
//
// 設計：擊破 T236 後，觸發「時間裂縫」：
//   - Server 查找玩家過去 30 秒內「最高倍率的那次擊破記錄」
//   - 立即給予觸發玩家 ×2.5 加成（個人，基於最高倍率目標的 bet × mult × 2.5）
//   - 同時在場上生成一個「裂縫複製體」（同種類目標，HP 只有 30%，擊破給 ×3.0 大獎）
//   - 若過去 30 秒無擊破記錄 → 給予 ×1.5 保底獎勵 + 生成隨機裂縫複製體
//   - 全服廣播「時間裂縫重現了什麼目標/倍率」
//   - 個人冷卻 20 秒；全服冷卻 32 秒
//
// 設計差異：
//   - 與時光倒流（T205，重播最近 5 個目標 ×1.6）不同，時間裂縫是「找最高倍率那次」
//     讓玩家有「我最好的那次擊破被時間裂縫記住了」的成就感
//   - 「裂縫複製體 HP 30%」讓玩家有「這條很容易打，要趕快打」的緊迫感
//   - 「裂縫複製體擊破 ×3.0」讓玩家有「打裂縫複製體比打普通魚更值」的策略感
//   - 「全服廣播裂縫複製體出現位置」讓所有玩家看到「有人觸發了時間裂縫，裂縫複製體在哪裡」
//   - 「30 秒歷史視窗」比時光倒流的 10 秒更長，讓玩家有更多機會觸發高倍率重現
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
	LuckyTimeRiftPersonalCD      = 20 * time.Second // 個人冷卻
	LuckyTimeRiftGlobalCD        = 32 * time.Second // 全服冷卻
	LuckyTimeRiftHistoryWindow   = 30 * time.Second // 歷史記錄視窗（30 秒）
	LuckyTimeRiftReplayMult      = 2.5              // 最高倍率重現加成
	LuckyTimeRiftFallbackMult    = 1.5              // 無記錄時的保底倍率
	LuckyTimeRiftCloneHPRatio    = 0.30             // 裂縫複製體 HP 比例（30%）
	LuckyTimeRiftCloneKillMult   = 3.0              // 裂縫複製體擊破倍率
	LuckyTimeRiftCloneLifetime   = 12               // 裂縫複製體存活時間（秒）
	LuckyTimeRiftCloneSpeedMult  = 0.7              // 裂縫複製體速度倍率（比原版慢）
)

// riftKillRecord 時間裂縫用的擊破記錄
type riftKillRecord struct {
	instanceID string
	defID      string
	name       string
	mult       int // 目標倍率（MultiplierMin）
	killedAt   time.Time
}

// luckyTimeRiftManager 幸運時間裂縫魚管理器
type luckyTimeRiftManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 玩家擊破歷史（playerID → []riftKillRecord，最近 30 秒）
	killHistory map[string][]riftKillRecord

	// 裂縫複製體標記（instanceID → true，用於識別裂縫複製體）
	riftClones map[string]bool
}

func newLuckyTimeRiftManager() *luckyTimeRiftManager {
	return &luckyTimeRiftManager{
		personalCooldowns: make(map[string]time.Time),
		killHistory:       make(map[string][]riftKillRecord),
		riftClones:        make(map[string]bool),
	}
}

// isLuckyTimeRiftFish 判斷是否為幸運時間裂縫魚
func isLuckyTimeRiftFish(defID string) bool {
	return defID == "T236"
}

// isRiftClone 判斷是否為裂縫複製體
func (m *luckyTimeRiftManager) isRiftClone(instanceID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.riftClones[instanceID]
}

// recordRiftKillHistory 記錄玩家擊破歷史（由 handleKill 呼叫）
func (m *luckyTimeRiftManager) recordRiftKillHistory(playerID, instanceID, defID, name string, mult int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-LuckyTimeRiftHistoryWindow)

	// 清理過期記錄
	history := m.killHistory[playerID]
	valid := history[:0]
	for _, r := range history {
		if r.killedAt.After(cutoff) {
			valid = append(valid, r)
		}
	}

	// 追加新記錄
	valid = append(valid, riftKillRecord{
		instanceID: instanceID,
		defID:      defID,
		name:       name,
		mult:       mult,
		killedAt:   now,
	})
	m.killHistory[playerID] = valid
}

// getBestKill 取得玩家過去 30 秒內最高倍率的擊破記錄
func (m *luckyTimeRiftManager) getBestKill(playerID string) (riftKillRecord, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-LuckyTimeRiftHistoryWindow)

	history := m.killHistory[playerID]
	var valid []riftKillRecord
	for _, r := range history {
		if r.killedAt.After(cutoff) {
			valid = append(valid, r)
		}
	}
	m.killHistory[playerID] = valid

	if len(valid) == 0 {
		return riftKillRecord{}, false
	}

	// 找最高倍率
	best := valid[0]
	for _, r := range valid[1:] {
		if r.mult > best.mult {
			best = r
		}
	}
	return best, true
}

// tryLuckyTimeRiftFish 擊破 T236 後觸發時間裂縫（供 handleKill 使用）
func (g *Game) tryLuckyTimeRiftFish(p *player.Player) {
	mgr := g.LuckyTimeRift
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
	mgr.personalCooldowns[p.ID] = now.Add(LuckyTimeRiftPersonalCD)
	mgr.globalCooldownUntil = now.Add(LuckyTimeRiftGlobalCD)
	mgr.mu.Unlock()

	log.Printf("[TimeRift] player=%s triggered time rift", p.ID)

	// 取得 betCost
	betDef := data.GetBetDef(p.BetLevel)
	betCost := betDef.BetCost
	if betCost < 1 {
		betCost = 1
	}

	// 查找最高倍率擊破記錄
	bestKill, hasBest := mgr.getBestKill(p.ID)

	var replayMult float64
	var replayDefID, replayName string
	var replayMultVal int

	if hasBest {
		replayMult = LuckyTimeRiftReplayMult
		replayDefID = bestKill.defID
		replayName = bestKill.name
		replayMultVal = bestKill.mult
	} else {
		// 無記錄：保底獎勵
		replayMult = LuckyTimeRiftFallbackMult
		replayDefID = "T001" // 預設用最基礎目標
		replayName = "裂縫幻影"
		replayMultVal = 2
	}

	// 計算即時獎勵（bet × mult × replayMult）
	immediateReward := int(float64(betCost) * float64(replayMultVal) * replayMult)
	if immediateReward < betCost {
		immediateReward = betCost
	}
	p.AddCoins(immediateReward)

	log.Printf("[TimeRift] player=%s bestKill=%s mult=%d replayMult=%.1f immediateReward=%d",
		p.ID, replayDefID, replayMultVal, replayMult, immediateReward)

	// 生成裂縫複製體
	cloneIID := g.spawnRiftClone(replayDefID, replayMultVal)

	// 個人訊息：時間裂縫觸發
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyTimeRift,
		Payload: ws.LuckyTimeRiftPayload{
			Event:           "rift_start",
			PlayerID:        p.ID,
			PlayerName:      p.DisplayName,
			HasBestKill:     hasBest,
			ReplayDefID:     replayDefID,
			ReplayName:      replayName,
			ReplayMult:      replayMultVal,
			RiftMult:        replayMult,
			ImmediateReward: immediateReward,
			CloneInstanceID: cloneIID,
			CloneKillMult:   LuckyTimeRiftCloneKillMult,
		},
	})

	// 全服廣播：時間裂縫觸發
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyTimeRift,
		Payload: ws.LuckyTimeRiftPayload{
			Event:           "rift_broadcast",
			PlayerName:      p.DisplayName,
			ReplayName:      replayName,
			ReplayMult:      replayMultVal,
			RiftMult:        replayMult,
			ImmediateReward: immediateReward,
			CloneInstanceID: cloneIID,
		},
	})

	// 全服公告
	var annMsg, annColor string
	if hasBest {
		annMsg = fmt.Sprintf("🌀 %s 觸發時間裂縫！最高倍率 %s（×%d）重現！即時獎勵 %d 籌碼！裂縫複製體出現！",
			p.DisplayName, replayName, replayMultVal, immediateReward)
		annColor = "#9B59B6" // 紫色
	} else {
		annMsg = fmt.Sprintf("🌀 %s 觸發時間裂縫！保底獎勵 %d 籌碼！裂縫複製體出現！",
			p.DisplayName, immediateReward)
		annColor = "#3498DB" // 藍色
	}
	ann := g.Announce.Create(announce.EventLuckyTimeRift, p.DisplayName, immediateReward, map[string]string{
		"message": annMsg,
		"color":   annColor,
	})
	g.broadcastAnnouncement(ann)
}

// spawnRiftClone 生成裂縫複製體（HP 30%，擊破給 ×3.0）
// 回傳裂縫複製體的 instanceID
func (g *Game) spawnRiftClone(defID string, origMult int) string {
	// 取得目標定義
	def, ok := data.Targets[defID]
	if !ok {
		// 找不到定義時用 T001
		def = data.Targets["T001"]
	}

	// 生成裂縫複製體（HP 30%，速度 70%）
	cloneHP := int(float64(def.HP) * LuckyTimeRiftCloneHPRatio)
	if cloneHP < 1 {
		cloneHP = 1
	}
	cloneSpeed := int(float64(def.Speed) * LuckyTimeRiftCloneSpeedMult)
	if cloneSpeed < 1 {
		cloneSpeed = 1
	}

	// 隨機生成位置（右側進入）
	spawnX := float64(1280 + 50)
	spawnY := float64(100 + rand.Intn(500))

	// 建立裂縫複製體目標
	cloneIID := uuid.New().String()
	clone := target.NewTarget(cloneIID, def, spawnX, spawnY)
	clone.HP = cloneHP
	clone.MaxHP = cloneHP

	// 加入遊戲
	g.mu.Lock()
	g.Targets[cloneIID] = clone
	g.mu.Unlock()

	// 標記為裂縫複製體
	g.LuckyTimeRift.mu.Lock()
	g.LuckyTimeRift.riftClones[cloneIID] = true
	g.LuckyTimeRift.mu.Unlock()

	// 廣播裂縫複製體生成
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgTargetSpawn,
		Payload: ws.TargetSpawnPayload{
			InstanceID: cloneIID,
			DefID:      def.ID,
			Name:       def.Name + "（裂縫）",
			Type:       string(def.Type),
			X:          spawnX,
			Y:          spawnY,
			HP:         cloneHP,
			MaxHP:      cloneHP,
			Speed:      float64(cloneSpeed),
			Lifetime:   float64(LuckyTimeRiftCloneLifetime),
			Behavior:   "linear",
			Multiplier: float64(origMult),
			Quality:    "legendary",
			QualityColor: "#9B59B6", // 紫色裂縫光暈
		},
	})

	log.Printf("[TimeRift] spawned rift clone: iid=%s defID=%s hp=%d speed=%d",
		cloneIID, defID, cloneHP, cloneSpeed)

	return cloneIID
}

// notifyRiftCloneKill 裂縫複製體被擊破時的結算（由 handleKill 呼叫）
func (g *Game) notifyRiftCloneKill(p *player.Player, t *target.Target) {
	// 計算裂縫複製體擊破獎勵（bet × mult × 3.0）
	betDef := data.GetBetDef(p.BetLevel)
	betCost := betDef.BetCost
	if betCost < 1 {
		betCost = 1
	}

	cloneReward := int(float64(betCost) * float64(t.Def.MultiplierMin) * LuckyTimeRiftCloneKillMult)
	if cloneReward < betCost {
		cloneReward = betCost
	}
	p.AddCoins(cloneReward)

	// 清除裂縫複製體標記
	g.LuckyTimeRift.mu.Lock()
	delete(g.LuckyTimeRift.riftClones, t.InstanceID)
	g.LuckyTimeRift.mu.Unlock()

	log.Printf("[TimeRift] rift clone killed: player=%s iid=%s reward=%d",
		p.ID, t.InstanceID, cloneReward)

	// 個人通知
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyTimeRift,
		Payload: ws.LuckyTimeRiftPayload{
			Event:       "rift_clone_kill",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			ReplayName:  t.Def.Name + "（裂縫）",
			CloneKillMult: LuckyTimeRiftCloneKillMult,
			CloneReward: cloneReward,
		},
	})

	// 全服廣播
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyTimeRift,
		Payload: ws.LuckyTimeRiftPayload{
			Event:       "rift_clone_broadcast",
			PlayerName:  p.DisplayName,
			ReplayName:  t.Def.Name + "（裂縫）",
			CloneReward: cloneReward,
		},
	})

	// 全服公告（高獎勵才公告）
	if cloneReward >= betCost*5 {
		ann := g.Announce.Create(announce.EventLuckyTimeRift, p.DisplayName, cloneReward, map[string]string{
			"message": fmt.Sprintf("🌀 %s 擊破裂縫複製體！獲得 %d 籌碼！×%.1f 裂縫大獎！",
				p.DisplayName, cloneReward, LuckyTimeRiftCloneKillMult),
			"color": "#FFD700",
		})
		g.broadcastAnnouncement(ann)
	}
}
