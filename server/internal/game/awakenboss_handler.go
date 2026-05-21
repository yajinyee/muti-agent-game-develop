// awakenboss_handler.go — 覺醒 BOSS 系統 handler（DAY-130）
// 業界依據：JILI Royal Fishing 2026 Awaken Boss
// 覺醒 BOSS（覺醒龍 90-180x / 冰鳳凰 120-300x）比不死 BOSS 更強，
// 並有 Power Up 機制：每 5-8 次命中觸發一次 Power Up（6x-10x 加成），
// 製造「蓄力爆發」的爽感，是 2026 年捕魚機最高倍率機制之一。
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

// trySpawnAwakenBoss 在 spawnTarget 時嘗試觸發覺醒 BOSS
func (g *Game) trySpawnAwakenBoss() {
	if g.AwakenBoss == nil {
		return
	}

	ok, def := g.AwakenBoss.ShouldTrigger()
	if !ok || def == nil {
		return
	}

	instanceID := uuid.New().String()
	g.AwakenBoss.StartSession(instanceID, def)

	log.Printf("[AwakenBoss] spawned: type=%s, instance=%s, mult=%.0f-%.0fx, powerup=%.0f-%.0fx, threshold=%d",
		def.ID, instanceID, def.MinMult, def.MaxMult, def.PowerUpMinMult, def.PowerUpMaxMult, def.PowerUpThreshold)

	// 全服廣播：覺醒 BOSS 出現
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAwakenBossSpawn,
		Payload: ws.AwakenBossSpawnPayload{
			InstanceID:       instanceID,
			BossType:         string(def.ID),
			BossName:         def.Name,
			BossIcon:         def.Icon,
			BossColor:        def.Color,
			MinMult:          def.MinMult,
			MaxMult:          def.MaxMult,
			PowerUpMinMult:   def.PowerUpMinMult,
			PowerUpMaxMult:   def.PowerUpMaxMult,
			PowerUpThreshold: def.PowerUpThreshold,
			DurationSeconds:  def.Duration,
			Message:          fmt.Sprintf("%s %s 覺醒！每 %d 次命中觸發 Power Up（%.0f-%.0fx 加成）！", def.Icon, def.Name, def.PowerUpThreshold, def.PowerUpMinMult, def.PowerUpMaxMult),
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventBossWarning, def.Name, 0, map[string]string{
		"message": fmt.Sprintf("⚡ 覺醒 BOSS %s %s 降臨！命中 %d 次觸發 Power Up！", def.Icon, def.Name, def.PowerUpThreshold),
	})
	g.broadcastAnnouncement(ann)
}

// tryAwakenBossHit 嘗試命中覺醒 BOSS（每次射擊時呼叫）
// 命中機率：15%（LV1-6）/ 25%（LV7+）
// 比不死 BOSS 低（覺醒 BOSS 更稀有，命中更難）
func (g *Game) tryAwakenBossHit(p *player.Player) {
	if g.AwakenBoss == nil || !g.AwakenBoss.IsActive() {
		return
	}

	hitChance := 0.15
	if p.BetLevel >= 7 {
		hitChance = 0.25
	}
	if rand.Float64() >= hitChance {
		return
	}

	instanceID := g.AwakenBoss.GetActiveInstanceID()
	if instanceID == "" {
		return
	}

	g.notifyAwakenBossHit(p, instanceID)
}

// notifyAwakenBossHit 處理命中覺醒 BOSS，發放獎勵並廣播
func (g *Game) notifyAwakenBossHit(p *player.Player, instanceID string) {
	if g.AwakenBoss == nil {
		return
	}

	betDef := p.GetBetDef()
	mult, reward, isPowerUp, ok := g.AwakenBoss.RecordHit(instanceID, p.ID, p.DisplayName, betDef.BetCost)
	if !ok {
		return
	}

	// 發放獎勵
	p.AddCoins(reward)

	snap := g.AwakenBoss.GetSnapshot()

	log.Printf("[AwakenBoss] hit: player=%s, mult=%.0fx, reward=%d, isPowerUp=%v, hits=%d",
		p.ID, mult, reward, isPowerUp, snap.HitCount)

	if isPowerUp {
		// Power Up 觸發：特殊廣播
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgAwakenBossPowerUp,
			Payload: ws.AwakenBossPowerUpPayload{
				InstanceID:   instanceID,
				BossName:     snap.BossName,
				BossIcon:     snap.BossIcon,
				PlayerID:     p.ID,
				PlayerName:   p.DisplayName,
				Multiplier:   mult,
				Reward:       reward,
				NewBalance:   p.GetCoins(),
				PowerUpCount: snap.PowerUpCount,
				Message:      fmt.Sprintf("⚡ %s 觸發 %s %s Power Up！獲得 %.0fx 大獎！", p.DisplayName, snap.BossIcon, snap.BossName, mult),
			},
		})

		// Power Up 全服公告
		ann := g.Announce.Create(announce.EventMegaWin, p.DisplayName, reward, map[string]string{
			"message": fmt.Sprintf("⚡ %s 觸發 %s Power Up！%.0fx！", p.DisplayName, snap.BossName, mult),
		})
		g.broadcastAnnouncement(ann)

		// 動態牆：Power Up 大獎
		go g.notifyFeedMegaWin(p, mult, reward)
	} else {
		// 普通命中廣播
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgAwakenBossHit,
			Payload: ws.AwakenBossHitPayload{
				InstanceID:      instanceID,
				PlayerID:        p.ID,
				PlayerName:      p.DisplayName,
				Multiplier:      mult,
				Reward:          reward,
				NewBalance:      p.GetCoins(),
				HitCount:        snap.HitCount,
				PowerUpProgress: snap.PowerUpProgress,
				TotalReward:     snap.TotalReward,
			},
		})
	}
}

// tickAwakenBoss 定期檢查覺醒 BOSS 是否過期
func (g *Game) tickAwakenBoss() {
	if g.AwakenBoss == nil {
		return
	}

	expired := g.AwakenBoss.CheckExpiry()
	if expired == nil {
		return
	}

	log.Printf("[AwakenBoss] left: type=%s, hits=%d, powerups=%d, total_reward=%d",
		expired.Def.ID, expired.HitCount, expired.PowerUpCount, expired.TotalReward)

	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAwakenBossLeave,
		Payload: ws.AwakenBossLeavePayload{
			InstanceID:   expired.InstanceID,
			BossName:     expired.Def.Name,
			BossIcon:     expired.Def.Icon,
			HitCount:     expired.HitCount,
			PowerUpCount: expired.PowerUpCount,
			TotalReward:  expired.TotalReward,
			Message:      fmt.Sprintf("%s %s 離去！共命中 %d 次，觸發 %d 次 Power Up，總獎勵 %d 金幣", expired.Def.Icon, expired.Def.Name, expired.HitCount, expired.PowerUpCount, expired.TotalReward),
		},
	})
}

// sendAwakenBossStatus 發送覺醒 BOSS 狀態給玩家（登入時呼叫）
func (g *Game) sendAwakenBossStatus(p *player.Player) {
	if g.AwakenBoss == nil {
		return
	}

	snap := g.AwakenBoss.GetSnapshot()
	if err := g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgAwakenBossStatus,
		Payload: ws.AwakenBossStatusPayload{
			Active:           snap.Active,
			InstanceID:       snap.InstanceID,
			BossType:         string(snap.BossType),
			BossName:         snap.BossName,
			BossIcon:         snap.BossIcon,
			BossColor:        snap.BossColor,
			MinMult:          snap.MinMult,
			MaxMult:          snap.MaxMult,
			PowerUpMinMult:   snap.PowerUpMinMult,
			PowerUpMaxMult:   snap.PowerUpMaxMult,
			PowerUpThreshold: snap.PowerUpThreshold,
			HitCount:         snap.HitCount,
			PowerUpCount:     snap.PowerUpCount,
			PowerUpProgress:  snap.PowerUpProgress,
			TotalReward:      snap.TotalReward,
			RemainingSeconds: snap.RemainingSeconds,
		},
	}); err != nil {
		log.Printf("[AwakenBoss] send status error: %v", err)
	}
}
