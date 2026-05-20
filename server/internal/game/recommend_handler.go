// recommend_handler.go — 智慧推薦系統 handler（DAY-110）
package game

import (
	"log"
	"time"

	"digital-twin/server/internal/game/recommend"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// handleGetRecommendations 處理查詢推薦請求
func (g *Game) handleGetRecommendations(p *player.Player) {
	recs := g.buildRecommendations(p)
	if err := g.Hub.Send(p.ID, &ws.Message{
		Type:    ws.MsgRecommendations,
		Payload: ws.RecommendationsPayload{
			Recommendations: recs,
			GeneratedAt:     time.Now().UnixMilli(),
		},
	}); err != nil {
		log.Printf("[Recommend] send error: %v", err)
	}
}

// buildRecommendations 根據玩家狀態建立推薦
func (g *Game) buildRecommendations(p *player.Player) []ws.RecommendationPayload {
	behavior := buildPlayerBehavior(p)
	engine := recommend.New()
	recs := engine.Analyze(behavior)

	var result []ws.RecommendationPayload
	for _, r := range recs {
		result = append(result, ws.RecommendationPayload{
			Type:        string(r.Type),
			Title:       r.Title,
			Description: r.Description,
			Icon:        r.Icon,
			Priority:    r.Priority,
			TargetBetLv: r.TargetBetLv,
			Confidence:  r.Confidence,
		})
	}
	return result
}

// buildPlayerBehavior 從玩家狀態建立行為摘要
func buildPlayerBehavior(p *player.Player) recommend.PlayerBehavior {
	b := recommend.PlayerBehavior{
		CurrentBetLv: p.BetLevel,
		CurrentCoins: p.Coins,
	}

	// 登入連續天數
	b.LoginStreak = p.LoginStreak

	// 統計資料
	if p.Stats != nil {
		snap := p.Stats.Snapshot()
		b.TotalShots = snap.TotalShots
		b.TotalKills = snap.TotalKills
		b.TotalBet = snap.TotalBet
		b.TotalReward = snap.TotalReward
		b.TotalBonuses = snap.TotalBonuses
		b.TotalBossKills = snap.TotalBossKills
		b.BestStreak = snap.BestStreak
		b.BestMultiplier = snap.BestMultiplier
		b.JackpotWins = snap.JackpotWins
	}

	return b
}
