// immortalboss_handler.go — 不死 BOSS 連勝系統 handler（DAY-129）
// 業界依據：JILI Royal Fishing 2026 Immortal Boss
// Golden Toad（金蟾蜍）和 Ancient Crocodile（古鱷魚）隨機出現，
// 每次命中都給獎勵（50x-150x），直到牠們自己離開畫面。
// 製造「連續獲勝序列」的爽感，是 2026 年捕魚機最熱門機制之一。
//
// 設計：不死 BOSS 出現時，玩家每次射擊有 25-35% 機率「順帶命中」不死 BOSS，
// 給予額外獎勵。這讓玩家在正常遊戲中也能感受到不死 BOSS 的存在，
// 不需要特別瞄準，符合捕魚機「廣撒網」的玩法精神。
package game

import (
	"fmt"
	"log"
	"math/rand"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"

	"github.com/google/uuid"
)

// trySpawnImmortalBoss 在 spawnTarget 時嘗試觸發不死 BOSS
// 由 spawnTarget 呼叫（goroutine 安全）
func (g *Game) trySpawnImmortalBoss() {
	if g.ImmortalBoss == nil {
		return
	}

	ok, def := g.ImmortalBoss.ShouldTrigger()
	if !ok || def == nil {
		return
	}

	instanceID := uuid.New().String()
	g.ImmortalBoss.StartSession(instanceID, def)

	log.Printf("[ImmortalBoss] spawned: type=%s, instance=%s, mult=%.0f-%.0fx, duration=%.0fs",
		def.ID, instanceID, def.MinMult, def.MaxMult, def.Duration)

	// 全服廣播：不死 BOSS 出現
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgImmortalBossSpawn,
		Payload: ws.ImmortalBossSpawnPayload{
			InstanceID:      instanceID,
			BossType:        string(def.ID),
			BossName:        def.Name,
			BossIcon:        def.Icon,
			BossColor:       def.Color,
			MinMult:         def.MinMult,
			MaxMult:         def.MaxMult,
			DurationSeconds: def.Duration,
			Message:         fmt.Sprintf("%s %s 出現了！每次命中獲得 %.0f-%.0fx 獎勵！", def.Icon, def.Name, def.MinMult, def.MaxMult),
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventBossWarning, def.Name, 0, map[string]string{
		"message": fmt.Sprintf("⚡ 不死 BOSS %s %s 出現！命中即得 %.0f-%.0fx！", def.Icon, def.Name, def.MinMult, def.MaxMult),
	})
	g.broadcastAnnouncement(ann)
}

// tryImmortalBossHit 嘗試命中不死 BOSS（每次射擊時呼叫）
// 命中機率：25%（LV1-6）/ 35%（LV7+）
func (g *Game) tryImmortalBossHit(p *player.Player) {
	if g.ImmortalBoss == nil || !g.ImmortalBoss.IsActive() {
		return
	}

	// 命中機率依投注等級
	hitChance := 0.25
	if p.BetLevel >= 7 {
		hitChance = 0.35
	}
	if rand.Float64() >= hitChance {
		return
	}

	instanceID := g.ImmortalBoss.GetActiveInstanceID()
	if instanceID == "" {
		return
	}

	g.notifyImmortalBossHit(p, instanceID)
}

// notifyImmortalBossHit 處理命中不死 BOSS，發放獎勵並廣播
func (g *Game) notifyImmortalBossHit(p *player.Player, instanceID string) {
	if g.ImmortalBoss == nil {
		return
	}

	betDef := p.GetBetDef()
	mult, reward, ok := g.ImmortalBoss.RecordHit(instanceID, p.ID, p.DisplayName, betDef.BetCost)
	if !ok {
		return
	}

	// 發放獎勵
	p.AddCoins(reward)

	snap := g.ImmortalBoss.GetSnapshot()
	isHighMult := mult >= 100.0

	log.Printf("[ImmortalBoss] hit: player=%s, mult=%.0fx, reward=%d, total_hits=%d",
		p.ID, mult, reward, snap.HitCount)

	// 廣播命中事件（全服可見，讓其他玩家也想來打）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgImmortalBossHit,
		Payload: ws.ImmortalBossHitPayload{
			InstanceID:  instanceID,
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			Multiplier:  mult,
			Reward:      reward,
			NewBalance:  p.GetCoins(),
			HitCount:    snap.HitCount,
			TotalReward: snap.TotalReward,
			IsHighMult:  isHighMult,
		},
	})

	// 高倍率（≥100x）時全服公告
	if isHighMult {
		ann := g.Announce.Create(announce.EventMegaWin, p.DisplayName, reward, map[string]string{
			"message": fmt.Sprintf("🌟 %s 命中 %s %s 獲得 %.0fx 大獎！", p.DisplayName, snap.BossIcon, snap.BossName, mult),
		})
		g.broadcastAnnouncement(ann)
	}

	// 動態牆：高倍率命中（≥100x）
	if isHighMult {
		go g.notifyFeedMegaWin(p, mult, reward)
	}
}

// tickImmortalBoss 定期檢查不死 BOSS 是否過期（由 gameLoop 每次 update 呼叫）
func (g *Game) tickImmortalBoss() {
	if g.ImmortalBoss == nil {
		return
	}

	expired := g.ImmortalBoss.CheckExpiry()
	if expired == nil {
		return
	}

	log.Printf("[ImmortalBoss] left: type=%s, hits=%d, total_reward=%d",
		expired.Def.ID, expired.HitCount, expired.TotalReward)

	// 全服廣播：不死 BOSS 離開
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgImmortalBossLeave,
		Payload: ws.ImmortalBossLeavePayload{
			InstanceID:  expired.InstanceID,
			BossName:    expired.Def.Name,
			BossIcon:    expired.Def.Icon,
			HitCount:    expired.HitCount,
			TotalReward: expired.TotalReward,
			Message:     fmt.Sprintf("%s %s 離開了！共被命中 %d 次，總獎勵 %d 金幣", expired.Def.Icon, expired.Def.Name, expired.HitCount, expired.TotalReward),
		},
	})
}

// sendImmortalBossStatus 發送不死 BOSS 狀態給玩家（登入時呼叫）
func (g *Game) sendImmortalBossStatus(p *player.Player) {
	if g.ImmortalBoss == nil {
		return
	}

	snap := g.ImmortalBoss.GetSnapshot()
	if err := g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgImmortalBossStatus,
		Payload: ws.ImmortalBossStatusPayload{
			Active:           snap.Active,
			InstanceID:       snap.InstanceID,
			BossType:         string(snap.BossType),
			BossName:         snap.BossName,
			BossIcon:         snap.BossIcon,
			BossColor:        snap.BossColor,
			MinMult:          snap.MinMult,
			MaxMult:          snap.MaxMult,
			HitCount:         snap.HitCount,
			TotalReward:      snap.TotalReward,
			RemainingSeconds: snap.RemainingSeconds,
		},
	}); err != nil {
		log.Printf("[ImmortalBoss] send status error: %v", err)
	}
}
