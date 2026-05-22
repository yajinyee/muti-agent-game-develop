// golden_treasure_fish_handler.go — 黃金寶藏魚系統 handler（DAY-177）
// 業界依據：Ocean King 3 Plus「Golden Treasure feature — catching triggers treasure chests
// to appear, players open them for random rewards」
// + JILI Giant Anglerfish「shoot electricity to open treasure chests」
// 擊破 T135 後觸發「黃金寶藏」：場上出現 3 個寶藏箱，玩家在 12 秒內點擊開啟，
// 每個寶藏箱隨機包含：金幣獎勵（50%）、倍率加成 ×3 持續 8 秒（30%）、特殊武器充能（20%）
// 設計差異：與幸運彩蛋魚（自動掉落+隨機開啟）不同，黃金寶藏是「玩家主動點擊開啟」，
// 製造「我選擇開哪個箱子」的互動感；寶藏箱在場上可見，讓其他玩家也能看到「有人觸發了寶藏」
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	// GoldenTreasureCooldownSec 個人冷卻時間（秒）
	GoldenTreasureCooldownSec = 35
	// GoldenTreasureTimeoutSec 寶藏箱超時時間（秒）
	GoldenTreasureTimeoutSec = 12
	// GoldenTreasureChestCount 寶藏箱數量
	GoldenTreasureChestCount = 3
	// GoldenTreasureMultBonus 倍率加成值
	GoldenTreasureMultBonus = 3.0
	// GoldenTreasureMultDurationSec 倍率加成持續時間（秒）
	GoldenTreasureMultDurationSec = 8
	// GoldenTreasureAnnounceThreshold 全服公告門檻（倍率加成時公告）
	GoldenTreasureAnnounceThreshold = 2 // 開出 ≥2 個倍率加成時公告
)

// TreasureRewardType 寶藏獎勵類型
type TreasureRewardType string

const (
	TreasureRewardCoins  TreasureRewardType = "coins"   // 金幣獎勵（50%）
	TreasureRewardMult   TreasureRewardType = "mult"    // 倍率加成（30%）
	TreasureRewardWeapon TreasureRewardType = "weapon"  // 特殊武器充能（20%）
)

// TreasureChest 單個寶藏箱
type TreasureChest struct {
	ID         int                // 箱子編號（0-2）
	RewardType TreasureRewardType // 預先決定的獎勵類型
	CoinReward int                // 金幣獎勵數量（RewardType=coins 時有效）
	Opened     bool               // 是否已開啟
}

// goldenTreasureSession 黃金寶藏 session（per-player）
type goldenTreasureSession struct {
	PlayerID  string
	Chests    [GoldenTreasureChestCount]*TreasureChest
	StartedAt time.Time
	MultCount int // 本次 session 開出的倍率加成數量
}

// goldenTreasureManager 黃金寶藏魚管理器
type goldenTreasureManager struct {
	mu       sync.Mutex
	sessions map[string]*goldenTreasureSession // playerID → session
	cooldown map[string]time.Time              // playerID → 冷卻結束時間
	// 倍率加成狀態（per-player）
	multActive    map[string]bool
	multExpiresAt map[string]time.Time
}

// newGoldenTreasureManager 建立黃金寶藏魚管理器
func newGoldenTreasureManager() *goldenTreasureManager {
	return &goldenTreasureManager{
		sessions:      make(map[string]*goldenTreasureSession),
		cooldown:      make(map[string]time.Time),
		multActive:    make(map[string]bool),
		multExpiresAt: make(map[string]time.Time),
	}
}

// isGoldenTreasureFish 判斷是否為黃金寶藏魚（T135）
func isGoldenTreasureFish(defID string) bool {
	return defID == "T135"
}

// isOnCooldown 檢查玩家是否在冷卻中
func (m *goldenTreasureManager) isOnCooldown(playerID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	cd, ok := m.cooldown[playerID]
	if !ok {
		return false
	}
	return time.Now().Before(cd)
}

// hasActiveSession 檢查玩家是否有活躍 session
func (m *goldenTreasureManager) hasActiveSession(playerID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	sess, ok := m.sessions[playerID]
	if !ok {
		return false
	}
	return time.Now().Before(sess.StartedAt.Add(time.Duration(GoldenTreasureTimeoutSec) * time.Second))
}

// startSession 開始新 session
func (m *goldenTreasureManager) startSession(playerID string, rng *rand.Rand, betLevel int) *goldenTreasureSession {
	m.mu.Lock()
	defer m.mu.Unlock()

	chests := [GoldenTreasureChestCount]*TreasureChest{}
	for i := 0; i < GoldenTreasureChestCount; i++ {
		rewardType := pickTreasureRewardType(rng)
		coinReward := 0
		if rewardType == TreasureRewardCoins {
			// 金幣獎勵：betLevel × 8-25x
			coinReward = betLevel * (8 + rng.Intn(18))
		}
		chests[i] = &TreasureChest{
			ID:         i,
			RewardType: rewardType,
			CoinReward: coinReward,
			Opened:     false,
		}
	}

	sess := &goldenTreasureSession{
		PlayerID:  playerID,
		Chests:    chests,
		StartedAt: time.Now(),
	}
	m.sessions[playerID] = sess
	return sess
}

// openChest 開啟寶藏箱（玩家點擊）
// 回傳 (chest, ok)
func (m *goldenTreasureManager) openChest(playerID string, chestID int) (*TreasureChest, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	sess, ok := m.sessions[playerID]
	if !ok {
		return nil, false
	}
	// 檢查 session 是否超時
	if time.Now().After(sess.StartedAt.Add(time.Duration(GoldenTreasureTimeoutSec) * time.Second)) {
		return nil, false
	}
	if chestID < 0 || chestID >= GoldenTreasureChestCount {
		return nil, false
	}
	chest := sess.Chests[chestID]
	if chest.Opened {
		return nil, false
	}
	chest.Opened = true
	if chest.RewardType == TreasureRewardMult {
		sess.MultCount++
	}
	return chest, true
}

// endSession 結束 session，設定冷卻
func (m *goldenTreasureManager) endSession(playerID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.sessions, playerID)
	m.cooldown[playerID] = time.Now().Add(time.Duration(GoldenTreasureCooldownSec) * time.Second)
}

// activateMult 激活倍率加成
func (m *goldenTreasureManager) activateMult(playerID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.multActive[playerID] = true
	m.multExpiresAt[playerID] = time.Now().Add(time.Duration(GoldenTreasureMultDurationSec) * time.Second)
}

// getGoldenTreasureMult 取得倍率加成（供 handleKill 使用）
func (m *goldenTreasureManager) getGoldenTreasureMult(playerID string) float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.multActive[playerID] {
		return 1.0
	}
	if time.Now().After(m.multExpiresAt[playerID]) {
		m.multActive[playerID] = false
		return 1.0
	}
	return GoldenTreasureMultBonus
}

// pickTreasureRewardType 加權隨機選擇獎勵類型
// 金幣 50%、倍率 30%、武器 20%
func pickTreasureRewardType(rng *rand.Rand) TreasureRewardType {
	r := rng.Intn(100)
	switch {
	case r < 50:
		return TreasureRewardCoins
	case r < 80:
		return TreasureRewardMult
	default:
		return TreasureRewardWeapon
	}
}

// tryGoldenTreasureFish 擊破 T135 後觸發黃金寶藏（DAY-177）
// 由 handleKill 呼叫（在 goroutine 中執行）
func (g *Game) tryGoldenTreasureFish(p *player.Player) {
	// 個人冷卻 + 活躍 session 檢查
	if g.GoldenTreasure.isOnCooldown(p.ID) {
		return
	}
	if g.GoldenTreasure.hasActiveSession(p.ID) {
		return
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 建立 session（預先決定所有箱子獎勵）
	sess := g.GoldenTreasure.startSession(p.ID, rng, p.BetLevel)

	log.Printf("[GoldenTreasure] player=%s session started, chests: %v/%v/%v",
		p.ID, sess.Chests[0].RewardType, sess.Chests[1].RewardType, sess.Chests[2].RewardType)

	// 廣播寶藏箱出現（個人）
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgGoldenTreasureFish,
		Payload: ws.GoldenTreasureFishPayload{
			Phase:       "treasure_start",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			ChestCount:  GoldenTreasureChestCount,
			TimeoutSec:  GoldenTreasureTimeoutSec,
		},
	})

	// 全服廣播：有人觸發寶藏
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgGoldenTreasureFish,
		Payload: ws.GoldenTreasureFishPayload{
			Phase:      "treasure_broadcast",
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
		},
	})

	// 等待玩家開箱或超時
	timeout := time.After(time.Duration(GoldenTreasureTimeoutSec) * time.Second)
	<-timeout

	// 超時後自動開啟未開的箱子（給玩家最後的獎勵）
	g.autoOpenRemainingChests(p, sess)

	// 結束 session
	g.GoldenTreasure.endSession(p.ID)

	log.Printf("[GoldenTreasure] player=%s session ended", p.ID)
}

// autoOpenRemainingChests 超時後自動開啟未開的箱子
func (g *Game) autoOpenRemainingChests(p *player.Player, sess *goldenTreasureSession) {
	for i := 0; i < GoldenTreasureChestCount; i++ {
		chest := sess.Chests[i]
		if chest.Opened {
			continue
		}
		// 自動開啟，給予獎勵（金幣獎勵減半，倍率/武器不給）
		if chest.RewardType == TreasureRewardCoins && chest.CoinReward > 0 {
			halfReward := chest.CoinReward / 2
			p.AddReward(halfReward)
			g.Hub.Send(p.ID, &ws.Message{
				Type: ws.MsgGoldenTreasureFish,
				Payload: ws.GoldenTreasureFishPayload{
					Phase:      "treasure_auto_open",
					ChestID:    i,
					RewardType: string(chest.RewardType),
					Reward:     halfReward,
					IsAuto:     true,
				},
			})
		}
	}

	// 廣播寶藏結束
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgGoldenTreasureFish,
		Payload: ws.GoldenTreasureFishPayload{
			Phase:    "treasure_end",
			PlayerID: p.ID,
		},
	})
}

// handleGoldenTreasureOpen 處理玩家點擊開箱（由 HandleMessage 呼叫）
func (g *Game) handleGoldenTreasureOpen(p *player.Player, chestID int) {
	chest, ok := g.GoldenTreasure.openChest(p.ID, chestID)
	if !ok {
		return
	}

	reward := 0
	multActivated := false

	switch chest.RewardType {
	case TreasureRewardCoins:
		reward = chest.CoinReward
		p.AddReward(reward)

	case TreasureRewardMult:
		// 激活倍率加成
		g.GoldenTreasure.activateMult(p.ID)
		multActivated = true
		// 廣播倍率激活
		g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgGoldenTreasureFish,
			Payload: ws.GoldenTreasureFishPayload{
				Phase:         "treasure_mult_start",
				MultBonus:     GoldenTreasureMultBonus,
				MultDurationSec: GoldenTreasureMultDurationSec,
			},
		})
		// 8 秒後廣播倍率結束
		go func() {
			time.Sleep(time.Duration(GoldenTreasureMultDurationSec) * time.Second)
			g.Hub.Send(p.ID, &ws.Message{
				Type: ws.MsgGoldenTreasureFish,
				Payload: ws.GoldenTreasureFishPayload{
					Phase: "treasure_mult_end",
				},
			})
		}()

	case TreasureRewardWeapon:
		// 武器充能（相當於擊破 30x 目標）
		g.notifyGoldenTreasureWeaponCharge(p)
	}

	// 廣播開箱結果（個人）
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgGoldenTreasureFish,
		Payload: ws.GoldenTreasureFishPayload{
			Phase:         "treasure_open",
			ChestID:       chestID,
			RewardType:    string(chest.RewardType),
			Reward:        reward,
			MultActivated: multActivated,
		},
	})

	// 全服公告：開出倍率加成時廣播
	if chest.RewardType == TreasureRewardMult {
		g.announceGoldenTreasureMult(p.DisplayName)
	}

	log.Printf("[GoldenTreasure] player=%s opened chest=%d type=%s reward=%d",
		p.ID, chestID, chest.RewardType, reward)
}

// notifyGoldenTreasureWeaponCharge 武器充能通知
func (g *Game) notifyGoldenTreasureWeaponCharge(p *player.Player) {
	// 使用現有的 notifySpecialWeaponCharge 機制，傳入高倍率觸發充能
	// 黃金寶藏武器充能 = 相當於擊破一個 30x 目標
	g.notifySpecialWeaponCharge(p, 30.0)
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgGoldenTreasureFish,
		Payload: ws.GoldenTreasureFishPayload{
			Phase:        "treasure_weapon_charge",
			WeaponCharge: 30,
		},
	})
}

// getGoldenTreasureMult 取得黃金寶藏倍率加成（供 handleKill 使用）
func (g *Game) getGoldenTreasureMult(playerID string) float64 {
	return g.GoldenTreasure.getGoldenTreasureMult(playerID)
}

// announceGoldenTreasureMult 全服公告黃金寶藏倍率加成（DAY-177）
func (g *Game) announceGoldenTreasureMult(playerName string) {
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAnnouncement,
		Payload: map[string]interface{}{
			"event_type": "golden_treasure_fish",
			"message":    fmt.Sprintf("💰 %s 開啟黃金寶藏！獲得 ×%.0f 倍率加成！", playerName, GoldenTreasureMultBonus),
			"color":      "#FFD700",
			"duration":   4.0,
			"priority":   3,
		},
	})
}
