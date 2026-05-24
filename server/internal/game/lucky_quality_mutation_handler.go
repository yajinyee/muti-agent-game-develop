// lucky_quality_mutation_handler.go — 幸運品質突變魚系統（DAY-272）
// 業界依據：Fishing Frenzy Chapter 3（2026-05-11）Quality Roll 系統
//           Fisch Roblox 的 Mutation 機制（150+ 種突變，0.1x 到 17x 倍率）
//           2026 年最熱門「品質突變+稀有度分層」機制，讓玩家有「這條魚有突變！品質越高獎勵越多！」的驚喜感
//
// 設計：擊破 T230 後，觸發「品質突變」：
//   - Server 為觸發玩家的下一次擊破「品質突變」（5個品質等級）
//   - Normal（40%）×1.0 / Rare（30%）×1.8 / Epic（18%）×3.5 / Legendary（9%）×6.0 / Mythic（3%）×10.0
//   - 品質效果持續到下一次擊破（一次性）
//   - Mythic 品質全服廣播 + 全服公告
//   - 個人冷卻 18 秒；全服冷卻 30 秒
//
// 設計差異：
//   - 與重擲（T229，動態累積最高值）不同，品質突變是「一次性稀有度抽取」，讓玩家有「這次突變是什麼品質？」的期待感
//   - 「5個品質等級」讓玩家有「Normal 到 Mythic 的稀有度層次感」
//   - 「Mythic ×10.0（3% 機率）」讓玩家有「要是突變到 Mythic 就賺大了」的期待感
//   - 「品質顏色視覺（灰/藍/紫/橙/彩虹）」讓玩家一眼看出突變品質
//   - 「Mythic 全服廣播」讓所有玩家看到「有人突變到 Mythic！」，製造羨慕感
//   - 「品質突變視覺特效（目標物閃光+品質顏色光暈）」讓玩家感受到「這條魚有突變！」
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

// 品質等級常數
const (
	QualityNormal    = "normal"    // 普通（灰色）
	QualityRare      = "rare"      // 稀有（藍色）
	QualityEpic      = "epic"      // 史詩（紫色）
	QualityLegendary = "legendary" // 傳說（橙色）
	QualityMythic    = "mythic"    // 神話（彩虹）
)

// 品質突變冷卻常數
const (
	LuckyQualityMutationPersonalCD = 18 * time.Second // 個人冷卻
	LuckyQualityMutationGlobalCD   = 30 * time.Second // 全服冷卻
	LuckyQualityMutationSessionTTL = 25 * time.Second // session 存活時間（等待下一次擊破）
)

// qualityTier 品質等級定義
type qualityTier struct {
	Name   string  // 品質名稱
	Mult   float64 // 倍率
	Weight int     // 抽取權重（總和 100）
	Color  string  // 顯示顏色
	Emoji  string  // 顯示 emoji
}

// qualityTiers 品質等級列表（按稀有度排列）
var qualityTiers = []qualityTier{
	{Name: QualityNormal, Mult: 1.0, Weight: 40, Color: "#AAAAAA", Emoji: "⬜"},
	{Name: QualityRare, Mult: 1.8, Weight: 30, Color: "#4A90D9", Emoji: "🔵"},
	{Name: QualityEpic, Mult: 3.5, Weight: 18, Color: "#9B59B6", Emoji: "🟣"},
	{Name: QualityLegendary, Mult: 6.0, Weight: 9, Color: "#FF8C00", Emoji: "🟠"},
	{Name: QualityMythic, Mult: 10.0, Weight: 3, Color: "#FF69B4", Emoji: "🌈"},
}

// qualityMutationSession 品質突變 session（等待下一次擊破）
type qualityMutationSession struct {
	playerID   string
	playerName string
	quality    qualityTier
	expiresAt  time.Time
	used       bool // 是否已用於一次擊破
}

// luckyQualityMutationManager 幸運品質突變魚管理器
type luckyQualityMutationManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 活躍 session（playerID → session）
	activeSessions map[string]*qualityMutationSession
}

func newLuckyQualityMutationManager() *luckyQualityMutationManager {
	return &luckyQualityMutationManager{
		personalCooldowns: make(map[string]time.Time),
		activeSessions:    make(map[string]*qualityMutationSession),
	}
}

// isLuckyQualityMutationFish 判斷是否為幸運品質突變魚
func isLuckyQualityMutationFish(defID string) bool {
	return defID == "T230"
}

// rollQuality 隨機抽取品質等級
func rollQuality(rng *rand.Rand) qualityTier {
	total := 0
	for _, t := range qualityTiers {
		total += t.Weight
	}
	roll := rng.Intn(total)
	cumulative := 0
	for _, t := range qualityTiers {
		cumulative += t.Weight
		if roll < cumulative {
			return t
		}
	}
	return qualityTiers[0] // fallback: Normal
}

// getQualityMutationMult 取得品質突變倍率（供 handleKill 使用）
// 回傳 (倍率, 是否有品質突變加成, 品質名稱)
func (m *luckyQualityMutationManager) getQualityMutationMult(playerID string) (float64, bool, string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, ok := m.activeSessions[playerID]
	if !ok || session.used || time.Now().After(session.expiresAt) {
		if ok {
			delete(m.activeSessions, playerID)
		}
		return 1.0, false, ""
	}

	// 標記為已使用
	session.used = true
	return session.quality.Mult, true, session.quality.Name
}

// consumeQualityMutationSession 消耗 session（擊破後呼叫，取得完整品質資訊）
func (m *luckyQualityMutationManager) consumeQualityMutationSession(playerID string) (*qualityMutationSession, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, ok := m.activeSessions[playerID]
	if !ok {
		return nil, false
	}
	delete(m.activeSessions, playerID)
	return session, true
}

// tryLuckyQualityMutationFish 擊破 T230 後觸發品質突變
func (g *Game) tryLuckyQualityMutationFish(p *player.Player) {
	m := g.LuckyQualityMutation

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
	m.personalCooldowns[p.ID] = now.Add(LuckyQualityMutationPersonalCD)
	m.globalCooldownUntil = now.Add(LuckyQualityMutationGlobalCD)
	m.mu.Unlock()

	// 抽取品質等級
	rng := rand.New(rand.NewSource(now.UnixNano()))
	quality := rollQuality(rng)

	log.Printf("[QualityMutation] player=%s 觸發品質突變！品質=%s 倍率=×%.1f",
		p.ID, quality.Name, quality.Mult)

	// 建立 session
	session := &qualityMutationSession{
		playerID:   p.ID,
		playerName: p.DisplayName,
		quality:    quality,
		expiresAt:  now.Add(LuckyQualityMutationSessionTTL),
	}

	m.mu.Lock()
	m.activeSessions[p.ID] = session
	m.mu.Unlock()

	// 個人通知：品質突變開始
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyQualityMutation,
		Payload: ws.LuckyQualityMutationPayload{
			Event:       "mutation_start",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			Quality:     quality.Name,
			QualityMult: quality.Mult,
			QualityColor: quality.Color,
			QualityEmoji: quality.Emoji,
		},
	})

	// 全服廣播（Legendary 以上才廣播）
	if quality.Name == QualityLegendary || quality.Name == QualityMythic {
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyQualityMutation,
			Payload: ws.LuckyQualityMutationPayload{
				Event:       "mutation_broadcast",
				PlayerName:  p.DisplayName,
				Quality:     quality.Name,
				QualityMult: quality.Mult,
				QualityColor: quality.Color,
				QualityEmoji: quality.Emoji,
			},
		})
	}

	// 全服公告（Mythic 才公告）
	if quality.Name == QualityMythic {
		g.Announce.Create(announce.EventLuckyQualityMutation, p.DisplayName, 0, map[string]string{
			"message": fmt.Sprintf("🌈 %s 觸發神話突變！×%.1f 等待下一擊！",
				p.DisplayName, quality.Mult),
			"color": "#FF69B4",
		})
	} else if quality.Name == QualityLegendary {
		g.Announce.Create(announce.EventLuckyQualityMutation, p.DisplayName, 0, map[string]string{
			"message": fmt.Sprintf("🟠 %s 觸發傳說突變！×%.1f 等待下一擊！",
				p.DisplayName, quality.Mult),
			"color": "#FF8C00",
		})
	}

	// session 超時清理
	go func() {
		time.Sleep(LuckyQualityMutationSessionTTL)
		m.mu.Lock()
		if s, ok := m.activeSessions[p.ID]; ok && s == session {
			delete(m.activeSessions, p.ID)
			log.Printf("[QualityMutation] player=%s session 超時，未使用", p.ID)
		}
		m.mu.Unlock()

		// 通知超時
		_ = g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgLuckyQualityMutation,
			Payload: ws.LuckyQualityMutationPayload{
				Event:    "mutation_expire",
				PlayerID: p.ID,
			},
		})
	}()
}

// notifyQualityMutationKill 品質突變加成被使用時呼叫（由 handleKill 在套用倍率後呼叫）
func (g *Game) notifyQualityMutationKill(p *player.Player, targetName string, reward int, quality qualityTier) {
	log.Printf("[QualityMutation] player=%s 使用品質突變！%s ×%.1f，目標=%s，獎勵=%d",
		p.ID, quality.Name, quality.Mult, targetName, reward)

	// 個人通知：品質突變結算
	_ = g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyQualityMutation,
		Payload: ws.LuckyQualityMutationPayload{
			Event:        "mutation_used",
			PlayerID:     p.ID,
			TargetName:   targetName,
			Reward:       reward,
			Quality:      quality.Name,
			QualityMult:  quality.Mult,
			QualityColor: quality.Color,
			QualityEmoji: quality.Emoji,
		},
	})

	// 全服廣播（Mythic 才廣播結果）
	if quality.Name == QualityMythic {
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyQualityMutation,
			Payload: ws.LuckyQualityMutationPayload{
				Event:        "mutation_result_broadcast",
				PlayerName:   p.DisplayName,
				Quality:      quality.Name,
				QualityMult:  quality.Mult,
				QualityColor: quality.Color,
				QualityEmoji: quality.Emoji,
				Reward:       reward,
			},
		})

		g.Announce.Create(announce.EventLuckyQualityMutation, p.DisplayName, 0, map[string]string{
			"message": fmt.Sprintf("🌈 %s 神話突變命中！×%.1f 大獎 %d！",
				p.DisplayName, quality.Mult, reward),
			"color": "#FF69B4",
		})
	} else if quality.Name == QualityLegendary {
		g.Announce.Create(announce.EventLuckyQualityMutation, p.DisplayName, 0, map[string]string{
			"message": fmt.Sprintf("🟠 %s 傳說突變命中！×%.1f 獎勵 %d！",
				p.DisplayName, quality.Mult, reward),
			"color": "#FF8C00",
		})
	}
}
