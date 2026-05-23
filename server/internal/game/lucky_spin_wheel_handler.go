// lucky_spin_wheel_handler.go — 幸運輪盤魚系統（DAY-269）
// 業界依據：Royal Fishing ChainLong King Wheel 個人輪盤版本
//
// 設計：擊破 T227 後，觸發「個人幸運輪盤」：
//   - 5 個扇區：×2.0 / ×3.0 / ×5.0 / ×8.0 / ×0.5（機率：35%/30%/20%/10%/5%）
//   - Server 隨機決定結果，廣播給觸發玩家
//   - 結果倍率套用到觸發玩家接下來 8 秒的所有擊破（個人）
//   - 個人冷卻 25 秒；全服冷卻 40 秒
//
// 設計差異：
//   - 與競速賽（T223，排名競爭）不同，輪盤是「個人運氣」，讓玩家有「轉輪盤，看看這次是幾倍」的期待感
//   - 「5 個扇區不同倍率」讓每次觸發都有不同結果，增加多樣性
//   - 「×8.0 最高倍率（10% 機率）」讓玩家有「要是轉到 8 倍就賺大了」的期待感
//   - 「×0.5 懲罰扇區（5% 機率）」讓輪盤有風險感，不是每次都穩賺
//   - 「8 秒加成期間」讓玩家有「要趁 8 秒內多打幾條魚」的緊迫感
//   - 「全服廣播輪盤結果」讓所有玩家看到「有人轉到 8 倍了！」，製造羨慕感
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
	LuckySpinWheelPersonalCD = 25 * time.Second // 個人冷卻
	LuckySpinWheelGlobalCD   = 40 * time.Second // 全服冷卻
	LuckySpinWheelDuration   = 8 * time.Second  // 加成持續時間
)

// spinWheelSector 輪盤扇區定義
type spinWheelSector struct {
	mult   float64
	weight int // 權重（總和 100）
}

// 輪盤扇區定義（5 個扇區）
var spinWheelSectors = []spinWheelSector{
	{mult: 2.0, weight: 35}, // 35% 機率
	{mult: 3.0, weight: 30}, // 30% 機率
	{mult: 5.0, weight: 20}, // 20% 機率
	{mult: 8.0, weight: 10}, // 10% 機率
	{mult: 0.5, weight: 5},  // 5% 機率（懲罰）
}

// spinWheelBoost 輪盤加成 session
type spinWheelBoost struct {
	playerID    string
	mult        float64
	expiresAt   time.Time
	killCount   int
	totalReward int
}

// luckySpinWheelManager 幸運輪盤魚管理器
type luckySpinWheelManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 活躍加成（playerID → boost）
	activeBoosts map[string]*spinWheelBoost
}

func newLuckySpinWheelManager() *luckySpinWheelManager {
	return &luckySpinWheelManager{
		personalCooldowns: make(map[string]time.Time),
		activeBoosts:      make(map[string]*spinWheelBoost),
	}
}

// isLuckySpinWheelFish 判斷是否為幸運輪盤魚
func isLuckySpinWheelFish(defID string) bool {
	return defID == "T227"
}

// getLuckySpinWheelMult 取得輪盤加成倍率（供 handleKill 使用）
func (m *luckySpinWheelManager) getLuckySpinWheelMult(playerID string) float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	boost, ok := m.activeBoosts[playerID]
	if !ok {
		return 1.0
	}
	if time.Now().After(boost.expiresAt) {
		delete(m.activeBoosts, playerID)
		return 1.0
	}
	return boost.mult
}

// spinWheel 隨機抽取輪盤結果
func spinWheel(rng *rand.Rand) (int, float64) {
	total := 0
	for _, s := range spinWheelSectors {
		total += s.weight
	}
	r := rng.Intn(total)
	cumulative := 0
	for i, s := range spinWheelSectors {
		cumulative += s.weight
		if r < cumulative {
			return i, s.mult
		}
	}
	return 0, spinWheelSectors[0].mult
}

// tryLuckySpinWheelFish 擊破 T227 後觸發輪盤
func (g *Game) tryLuckySpinWheelFish(p *player.Player) {
	m := g.LuckySpinWheel

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
	m.personalCooldowns[p.ID] = now.Add(LuckySpinWheelPersonalCD)
	m.globalCooldownUntil = now.Add(LuckySpinWheelGlobalCD)
	m.mu.Unlock()

	// 抽取輪盤結果
	rng := rand.New(rand.NewSource(now.UnixNano()))
	sectorIdx, resultMult := spinWheel(rng)

	log.Printf("[SpinWheel] player=%s 觸發輪盤！結果 ×%.1f（扇區 %d）",
		p.ID, resultMult, sectorIdx)

	// 建立加成 session
	boost := &spinWheelBoost{
		playerID:  p.ID,
		mult:      resultMult,
		expiresAt: now.Add(LuckySpinWheelDuration),
	}
	m.mu.Lock()
	m.activeBoosts[p.ID] = boost
	m.mu.Unlock()

	// 收集扇區倍率列表
	sectors := make([]float64, len(spinWheelSectors))
	for i, s := range spinWheelSectors {
		sectors[i] = s.mult
	}

	// 個人訊息：觸發者（輪盤開始）
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckySpinWheel,
		Payload: ws.LuckySpinWheelPayload{
			Event:      "spin_start",
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			Sectors:    sectors,
		},
	})

	// 短暫延遲模擬輪盤旋轉動畫（0.8 秒）
	time.Sleep(800 * time.Millisecond)

	// 個人訊息：輪盤結果
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckySpinWheel,
		Payload: ws.LuckySpinWheelPayload{
			Event:       "spin_result",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			ResultMult:  resultMult,
			SectorIndex: sectorIdx,
			Duration:    LuckySpinWheelDuration.Seconds(),
		},
	})

	// 全服廣播
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckySpinWheel,
		Payload: ws.LuckySpinWheelPayload{
			Event:      "spin_broadcast",
			PlayerName: p.DisplayName,
			ResultMult: resultMult,
		},
	})

	// 全服公告（×5.0 以上才公告）
	if resultMult >= 5.0 {
		g.Announce.Create(announce.EventLuckySpinWheel, p.DisplayName, 0, map[string]string{
			"message": fmt.Sprintf("🎡 %s 幸運輪盤轉到 ×%.1f！8 秒黃金時間！",
				p.DisplayName, resultMult),
			"color": "#FFD700",
		})
	} else {
		g.Announce.Create(announce.EventLuckySpinWheel, p.DisplayName, 0, map[string]string{
			"message": fmt.Sprintf("🎡 %s 觸發幸運輪盤！結果 ×%.1f！",
				p.DisplayName, resultMult),
			"color": "#FF69B4",
		})
	}

	// 啟動超時 goroutine
	go g.runSpinWheelTimeout(p, boost)
}

// notifySpinWheelBoostKill 輪盤加成期間擊破目標時呼叫（由 handleKill 在套用倍率後呼叫）
func (g *Game) notifySpinWheelBoostKill(p *player.Player, targetName string, reward int) {
	m := g.LuckySpinWheel

	m.mu.Lock()
	boost, ok := m.activeBoosts[p.ID]
	if !ok || time.Now().After(boost.expiresAt) {
		m.mu.Unlock()
		return
	}
	boost.killCount++
	boost.totalReward += reward
	currentMult := boost.mult
	m.mu.Unlock()

	// 個人通知
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckySpinWheel,
		Payload: ws.LuckySpinWheelPayload{
			Event:       "spin_boost_kill",
			PlayerID:    p.ID,
			TargetName:  targetName,
			Reward:      reward,
			CurrentMult: currentMult,
		},
	})
}

// runSpinWheelTimeout 加成超時後結算
func (g *Game) runSpinWheelTimeout(p *player.Player, boost *spinWheelBoost) {
	timer := time.NewTimer(LuckySpinWheelDuration)
	defer timer.Stop()

	<-timer.C

	m := g.LuckySpinWheel
	m.mu.Lock()
	currentBoost, ok := m.activeBoosts[p.ID]
	if !ok || currentBoost != boost {
		m.mu.Unlock()
		return
	}
	killCount := boost.killCount
	totalReward := boost.totalReward
	delete(m.activeBoosts, p.ID)
	m.mu.Unlock()

	log.Printf("[SpinWheel] player=%s 加成結束！×%.1f，擊破 %d 個，總獎勵 %d",
		p.ID, boost.mult, killCount, totalReward)

	// 個人結算通知
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckySpinWheel,
		Payload: ws.LuckySpinWheelPayload{
			Event:       "spin_expire",
			PlayerID:    p.ID,
			TotalReward: totalReward,
			KillCount:   killCount,
			CurrentMult: boost.mult,
		},
	})
}
