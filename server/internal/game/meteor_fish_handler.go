// meteor_fish_handler.go — 隕石魚隕石雨系統 handler（DAY-184）
// 業界依據：Royal Fishing JILI「Dragon Wrath — unleash a massive meteorite attack across
// the centre screen, simultaneously targeting multiple fish including Immortal Bosses」
// 擊破 T142 後觸發「隕石雨」：
//   1. 5-10 顆隕石從天而降，每顆命中隨機目標
//   2. 每顆隕石有 70% 擊破機率，獎勵 0.60x 倍率
//   3. 每顆隕石間隔 300ms，製造「連續轟炸」的視覺爽感
//   4. 全服廣播每顆隕石落點，讓所有玩家看到「隕石在轟炸全場」
// 設計差異：
//   - 與閃電魚（時間驅動/8秒/16次）不同，隕石魚是「數量驅動」（5-10顆），
//     每顆隕石都是獨立的「天降神兵」，視覺上更有衝擊感
//   - 與漩渦魚（吸引同類）不同，隕石魚是「隨機轟炸」（任意目標），
//     讓玩家感受到「天降神兵，無差別攻擊」的爽感
//   - 隕石可以命中 BOSS（但機率降低到 30%），讓玩家有「隕石打 BOSS」的驚喜感
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
	// MeteorFishMinMeteors 最少隕石數
	MeteorFishMinMeteors = 5
	// MeteorFishMaxMeteors 最多隕石數
	MeteorFishMaxMeteors = 10
	// MeteorFishIntervalMs 每顆隕石間隔（ms）
	MeteorFishIntervalMs = 300
	// MeteorFishNormalKillChance 普通目標擊破機率（70%）
	MeteorFishNormalKillChance = 0.70
	// MeteorFishBossKillChance BOSS 擊破機率（30%）
	MeteorFishBossKillChance = 0.30
	// MeteorFishRewardMult 擊破獎勵倍率
	MeteorFishRewardMult = 0.60
	// MeteorFishCooldownSec 全服冷卻時間（秒）
	MeteorFishCooldownSec = 35
	// MeteorFishAnnounceMinKills 全服公告最低擊破數
	MeteorFishAnnounceMinKills = 4
)

// meteorFishManager 隕石魚管理器（全服共享）
type meteorFishManager struct {
	mu          sync.Mutex
	isActive    bool
	cooldownEnd time.Time
}

// newMeteorFishManager 建立隕石魚管理器
func newMeteorFishManager() *meteorFishManager {
	return &meteorFishManager{}
}

// isMeteorFish 判斷是否為隕石魚（T142）
func isMeteorFish(defID string) bool {
	return defID == "T142"
}

// isOnCooldownMeteor 檢查是否在全服冷卻中
func (m *meteorFishManager) isOnCooldownMeteor() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return time.Now().Before(m.cooldownEnd)
}

// activateMeteor 激活隕石雨
func (m *meteorFishManager) activateMeteor() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.isActive {
		return false
	}
	m.isActive = true
	m.cooldownEnd = time.Now().Add(time.Duration(MeteorFishCooldownSec) * time.Second)
	return true
}

// deactivateMeteor 結束隕石雨
func (m *meteorFishManager) deactivateMeteor() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.isActive = false
}

// tryMeteorFishShower 擊破 T142 後觸發隕石雨（DAY-184）
// 由 handleKill 呼叫（在 goroutine 中執行）
func (g *Game) tryMeteorFishShower(p *player.Player, triggerID string, triggerX, triggerY float64) {
	// 全服冷卻檢查
	if g.MeteorFish.isOnCooldownMeteor() {
		return
	}
	if !g.MeteorFish.activateMeteor() {
		return // 已有其他隕石雨在進行
	}
	defer g.MeteorFish.deactivateMeteor()

	// 決定隕石數量（5-10 顆，加權：5-7 顆機率高）
	meteorCount := pickMeteorCount()

	log.Printf("[MeteorFish] player=%s triggered meteor shower, meteors=%d", p.ID, meteorCount)

	// 廣播隕石雨開始（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgMeteorFish,
		Payload: ws.MeteorFishPayload{
			Phase:       "meteor_start",
			TriggerID:   triggerID,
			TriggerX:    triggerX,
			TriggerY:    triggerY,
			KillerID:    p.ID,
			KillerName:  p.DisplayName,
			MeteorCount: meteorCount,
		},
	})

	totalReward := 0
	totalKills := 0

	for i := 0; i < meteorCount; i++ {
		time.Sleep(time.Duration(MeteorFishIntervalMs) * time.Millisecond)

		// 從全場隨機選一個目標（包含 BOSS，但機率不同）
		g.mu.RLock()
		type candidate struct {
			instanceID string
			defID      string
			x, y       float64
			multiplier float64
			isBoss     bool
		}
		var normalCandidates []candidate
		var bossCandidates []candidate
		for id, t := range g.Targets {
			if id == triggerID || t.HP <= 0 {
				continue
			}
			c := candidate{t.InstanceID, t.DefID, t.X, t.Y, t.Multiplier, t.DefID == "B001"}
			if c.isBoss {
				bossCandidates = append(bossCandidates, c)
			} else {
				normalCandidates = append(normalCandidates, c)
			}
		}
		g.mu.RUnlock()

		// 優先選普通目標，10% 機率選 BOSS（如果有的話）
		var dt *candidate
		if len(bossCandidates) > 0 && rand.Float64() < 0.10 {
			c := bossCandidates[rand.Intn(len(bossCandidates))]
			dt = &c
		} else if len(normalCandidates) > 0 {
			c := normalCandidates[rand.Intn(len(normalCandidates))]
			dt = &c
		} else if len(bossCandidates) > 0 {
			c := bossCandidates[rand.Intn(len(bossCandidates))]
			dt = &c
		}

		if dt == nil {
			break // 沒有目標，結束
		}

		// 廣播隕石落點（全服）
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgMeteorFish,
			Payload: ws.MeteorFishPayload{
				Phase:     fmt.Sprintf("meteor_%d", i+1),
				TargetID:  dt.instanceID,
				TargetX:   dt.x,
				TargetY:   dt.y,
				MeteorNum: i + 1,
				KillerID:  p.ID,
				IsBoss:    dt.isBoss,
			},
		})

		// 判斷擊破機率
		killChance := MeteorFishNormalKillChance
		if dt.isBoss {
			killChance = MeteorFishBossKillChance
		}
		if rand.Float64() >= killChance {
			continue // 未擊破
		}

		// 擊破目標
		g.mu.Lock()
		t, ok := g.Targets[dt.instanceID]
		if !ok || t.HP <= 0 {
			g.mu.Unlock()
			continue
		}
		reward := int(float64(p.BetLevel) * dt.multiplier * MeteorFishRewardMult)
		if reward < 1 {
			reward = 1
		}
		t.HP = 0
		delete(g.Targets, dt.instanceID)
		g.mu.Unlock()

		totalReward += reward
		totalKills++

		// 廣播目標擊破
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgTargetKill,
			Payload: ws.TargetKillPayload{
				InstanceID: dt.instanceID,
				KillerID:   p.ID,
				Reward:     reward,
				Multiplier: dt.multiplier,
			},
		})

		log.Printf("[MeteorFish] meteor[%d] target=%s mult=%.0f reward=%d boss=%v",
			i+1, dt.instanceID, dt.multiplier, reward, dt.isBoss)
	}

	if totalReward > 0 {
		p.AddReward(totalReward)
	}

	// 廣播隕石雨結果（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgMeteorFish,
		Payload: ws.MeteorFishPayload{
			Phase:       "meteor_result",
			TriggerID:   triggerID,
			MeteorCount: meteorCount,
			TotalKills:  totalKills,
			TotalReward: totalReward,
			KillerID:    p.ID,
			KillerName:  p.DisplayName,
		},
	})

	// 個人結果通知
	if totalReward > 0 {
		g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgReward,
			Payload: ws.RewardPayload{
				Source:     "meteor_fish",
				Amount:     totalReward,
				Multiplier: float64(totalKills),
				NewBalance: p.Coins,
			},
		})
	}

	// 全服公告：擊破 ≥ 4 個
	if totalKills >= MeteorFishAnnounceMinKills {
		g.announceMeteorFish(p.DisplayName, meteorCount, totalKills, totalReward)
	}

	log.Printf("[MeteorFish] player=%s meteors=%d kills=%d total_reward=%d",
		p.ID, meteorCount, totalKills, totalReward)
}

// pickMeteorCount 加權隨機選擇隕石數量（5-10 顆）
// 5顆:30%, 6顆:25%, 7顆:20%, 8顆:15%, 9顆:7%, 10顆:3%
func pickMeteorCount() int {
	r := rand.Float64()
	switch {
	case r < 0.30:
		return 5
	case r < 0.55:
		return 6
	case r < 0.75:
		return 7
	case r < 0.90:
		return 8
	case r < 0.97:
		return 9
	default:
		return 10
	}
}

// announceMeteorFish 全服公告隕石魚隕石雨（DAY-184）
func (g *Game) announceMeteorFish(playerName string, meteors, kills, reward int) {
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAnnouncement,
		Payload: map[string]interface{}{
			"event_type": "meteor_fish",
			"message":    fmt.Sprintf("☄️ %s 的隕石魚觸發隕石雨！%d 顆隕石擊破 %d 個目標！獲得 %d 金幣！", playerName, meteors, kills, reward),
			"color":      "#FF6600", // 橙紅色（隕石感）
			"duration":   5.0,
			"priority":   3,
		},
	})
}
