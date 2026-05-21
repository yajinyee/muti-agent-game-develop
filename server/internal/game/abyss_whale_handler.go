// abyss_whale_handler.go — 深淵巨鯨全服 Boss 挑戰系統（DAY-164）
// 業界依據：Fishing Frenzy Chapter 3 2026「Boss Fish as endgame content for higher-level players」
// + Ocean King 2026「Abyss Whale — massive HP boss requiring full server cooperation」
// 深淵巨鯨擁有超高 HP（500），需要全服玩家合力攻擊才能擊破
// 擊破後觸發「深淵寶藏」，按傷害貢獻比例分配獎勵（最高 500x betLevel）
// 設計：合作型終局內容，與船長魚（競技型）形成對比，製造「全服合力打 Boss」的緊張爽感
package game

import (
	"fmt"
	"log"
	"sort"
	"sync"
	"time"

	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// abyssWhaleManager 深淵巨鯨全服 Boss 挑戰管理器
type abyssWhaleManager struct {
	mu           sync.Mutex
	isActive     bool          // 是否有深淵巨鯨在場
	instanceID   string        // 當前巨鯨的 InstanceID
	spawnAt      time.Time     // 出現時間
	totalHP      int           // 總 HP（固定 500）
	currentHP    int           // 當前 HP
	contributions map[string]*whaleContribution // playerID -> 貢獻
	cooldownEnd  time.Time     // 冷卻結束時間
	lastBroadcast time.Time    // 上次廣播時間（節流）
}

type whaleContribution struct {
	PlayerID    string
	PlayerName  string
	Damage      int // 累積傷害
	BetLevel    int
}

func newAbyssWhaleManager() *abyssWhaleManager {
	return &abyssWhaleManager{
		contributions: make(map[string]*whaleContribution),
	}
}

// isAbyssWhale 判斷是否為深淵巨鯨
func isAbyssWhale(defID string) bool {
	return defID == "T124"
}

// notifyAbyssWhaleSpawn 深淵巨鯨出現時呼叫（由 spawnTarget 呼叫）
func (g *Game) notifyAbyssWhaleSpawn(instanceID string, x, y float64) {
	if g.AbyssWhale == nil {
		return
	}

	const totalHP = 500
	const cooldown = 180.0 // 180 秒冷卻（Boss 級別，冷卻更長）

	g.AbyssWhale.mu.Lock()
	if g.AbyssWhale.isActive {
		g.AbyssWhale.mu.Unlock()
		return // 已有巨鯨在場
	}
	if time.Now().Before(g.AbyssWhale.cooldownEnd) {
		g.AbyssWhale.mu.Unlock()
		return // 冷卻中
	}

	g.AbyssWhale.isActive = true
	g.AbyssWhale.instanceID = instanceID
	g.AbyssWhale.spawnAt = time.Now()
	g.AbyssWhale.totalHP = totalHP
	g.AbyssWhale.currentHP = totalHP
	g.AbyssWhale.contributions = make(map[string]*whaleContribution)
	g.AbyssWhale.lastBroadcast = time.Now()
	g.AbyssWhale.mu.Unlock()

	log.Printf("[AbyssWhale] spawned instanceID=%s totalHP=%d", instanceID, totalHP)

	// 廣播深淵巨鯨出現（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAbyssWhale,
		Payload: ws.AbyssWhalePayload{
			Phase:      "whale_spawn",
			InstanceID: instanceID,
			X:          x,
			Y:          y,
			TotalHP:    totalHP,
			CurrentHP:  totalHP,
			HPPercent:  1.0,
		},
	})

	// 全服公告
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAnnouncement,
		Payload: map[string]interface{}{
			"event_type": "abyss_whale_spawn",
			"message":    "🐋 深淵巨鯨出現！全服合力擊破，按貢獻分配深淵寶藏！",
			"color":      "#0066CC",
			"duration":   6.0,
			"priority":   6,
		},
	})
}

// notifyAbyssWhaleHit 深淵巨鯨被命中時呼叫（由 handleAttack 呼叫）
func (g *Game) notifyAbyssWhaleHit(p *player.Player, instanceID string, damage int) {
	if g.AbyssWhale == nil {
		return
	}

	g.AbyssWhale.mu.Lock()
	if !g.AbyssWhale.isActive || g.AbyssWhale.instanceID != instanceID {
		g.AbyssWhale.mu.Unlock()
		return
	}

	// 記錄貢獻
	contrib, ok := g.AbyssWhale.contributions[p.ID]
	if !ok {
		contrib = &whaleContribution{
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			BetLevel:   p.BetLevel,
		}
		g.AbyssWhale.contributions[p.ID] = contrib
	}
	contrib.Damage += damage

	// 扣除 HP
	g.AbyssWhale.currentHP -= damage
	if g.AbyssWhale.currentHP < 0 {
		g.AbyssWhale.currentHP = 0
	}
	currentHP := g.AbyssWhale.currentHP
	totalHP := g.AbyssWhale.totalHP
	hpPercent := float64(currentHP) / float64(totalHP)

	// 節流廣播（每 0.5 秒最多廣播一次 HP 更新）
	shouldBroadcast := time.Since(g.AbyssWhale.lastBroadcast) >= 500*time.Millisecond
	if shouldBroadcast {
		g.AbyssWhale.lastBroadcast = time.Now()
	}
	g.AbyssWhale.mu.Unlock()

	// 廣播 HP 更新
	if shouldBroadcast {
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgAbyssWhale,
			Payload: ws.AbyssWhalePayload{
				Phase:      "whale_hp_update",
				InstanceID: instanceID,
				TotalHP:    totalHP,
				CurrentHP:  currentHP,
				HPPercent:  hpPercent,
				AttackerID: p.ID,
			},
		})
	}
}

// notifyAbyssWhaleKill 深淵巨鯨被擊破時呼叫（由 handleKill 呼叫）
func (g *Game) notifyAbyssWhaleKill(p *player.Player, instanceID string, x, y float64) {
	if g.AbyssWhale == nil {
		return
	}

	g.AbyssWhale.mu.Lock()
	if !g.AbyssWhale.isActive || g.AbyssWhale.instanceID != instanceID {
		g.AbyssWhale.mu.Unlock()
		return
	}

	g.AbyssWhale.isActive = false
	g.AbyssWhale.cooldownEnd = time.Now().Add(180 * time.Second)

	// 收集貢獻列表
	contribs := make([]*whaleContribution, 0, len(g.AbyssWhale.contributions))
	totalDamage := 0
	for _, c := range g.AbyssWhale.contributions {
		contribs = append(contribs, c)
		totalDamage += c.Damage
	}
	// 確保擊破者有貢獻記錄
	found := false
	for _, c := range contribs {
		if c.PlayerID == p.ID {
			found = true
			break
		}
	}
	if !found {
		contribs = append(contribs, &whaleContribution{
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			Damage:     1,
			BetLevel:   p.BetLevel,
		})
		totalDamage++
	}
	g.AbyssWhale.mu.Unlock()

	log.Printf("[AbyssWhale] killed by player=%s, %d contributors, totalDamage=%d",
		p.ID, len(contribs), totalDamage)

	// 按貢獻比例分配獎勵
	entries := g.distributeAbyssWhaleRewards(contribs, totalDamage)

	// 廣播深淵巨鯨擊破（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAbyssWhale,
		Payload: ws.AbyssWhalePayload{
			Phase:       "whale_killed",
			InstanceID:  instanceID,
			X:           x,
			Y:           y,
			KillerID:    p.ID,
			KillerName:  p.DisplayName,
			Entries:     entries,
			TotalDamage: totalDamage,
		},
	})

	// 全服公告
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAnnouncement,
		Payload: map[string]interface{}{
			"event_type": "abyss_whale_killed",
			"message":    fmt.Sprintf("🏆 %s 擊破深淵巨鯨！%d 名玩家獲得深淵寶藏！", p.DisplayName, len(entries)),
			"color":      "#00AAFF",
			"duration":   6.0,
			"priority":   6,
		},
	})
}

// distributeAbyssWhaleRewards 按貢獻比例分配深淵寶藏獎勵
func (g *Game) distributeAbyssWhaleRewards(contribs []*whaleContribution, totalDamage int) []ws.AbyssWhaleEntry {
	// 按貢獻排序（降序）
	sort.Slice(contribs, func(i, j int) bool {
		return contribs[i].Damage > contribs[j].Damage
	})

	g.mu.RLock()
	players := g.Players
	g.mu.RUnlock()

	entries := make([]ws.AbyssWhaleEntry, 0, len(contribs))
	for i, c := range contribs {
		if totalDamage == 0 {
			break
		}
		// 貢獻比例（0.0 ~ 1.0）
		ratio := float64(c.Damage) / float64(totalDamage)

		// 獎勵計算：最高 500x betLevel（第一名全貢獻），按比例縮放
		// 最低保底：1x betLevel（確保每個參與者都有獎勵）
		p, ok := players[c.PlayerID]
		if !ok {
			continue
		}
		maxReward := p.BetLevel * 500
		bonus := int(float64(maxReward) * ratio)
		if bonus < p.BetLevel {
			bonus = p.BetLevel // 保底 1x betLevel
		}

		p.AddReward(bonus)

		log.Printf("[AbyssWhale] rank=%d player=%s damage=%d ratio=%.1f%% bonus=%d",
			i+1, c.PlayerID, c.Damage, ratio*100, bonus)

		entry := ws.AbyssWhaleEntry{
			Rank:       i + 1,
			PlayerID:   c.PlayerID,
			PlayerName: c.PlayerName,
			Damage:     c.Damage,
			Ratio:      ratio,
			Bonus:      bonus,
		}
		entries = append(entries, entry)

		// 個人通知
		_ = g.Hub.Send(c.PlayerID, &ws.Message{
			Type: ws.MsgAbyssWhale,
			Payload: ws.AbyssWhalePayload{
				Phase:      "whale_reward",
				InstanceID: g.AbyssWhale.instanceID,
				KillerID:   contribs[0].PlayerID,
				KillerName: contribs[0].PlayerName,
				Entries:    entries,
				MyRank:     i + 1,
				MyBonus:    bonus,
				MyDamage:   c.Damage,
				MyRatio:    ratio,
			},
		})

		// 第一名（最高貢獻）全服公告
		if i == 0 && bonus >= p.BetLevel*100 {
			g.Hub.Broadcast(&ws.Message{
				Type: ws.MsgAnnouncement,
				Payload: map[string]interface{}{
					"event_type": "abyss_whale_top",
					"message":    fmt.Sprintf("🐋 %s 貢獻最多傷害（%.0f%%），獲得深淵寶藏 %d 金幣！", c.PlayerName, ratio*100, bonus),
					"color":      "#00CCFF",
					"duration":   5.0,
					"priority":   5,
				},
			})
		}
	}
	return entries
}
