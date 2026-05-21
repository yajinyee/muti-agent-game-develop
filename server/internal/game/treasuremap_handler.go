// treasuremap_handler.go — 寶藏地圖系統 handler（DAY-122）
// 業界依據：bsu.edu（2026）確認「Hidden Treasure Unlocks」是 2026 年捕魚機最新趨勢
// 玩家擊破特定目標物收集地圖格子，集滿一行/列/對角線觸發寶藏獎勵（類似賓果）
package game

import (
	"fmt"
	"log"

	"digital-twin/server/internal/data"
	"digital-twin/server/internal/game/treasuremap"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// notifyTreasureMapKill 在擊破目標後更新寶藏地圖（由 handleKill 呼叫）
func (g *Game) notifyTreasureMapKill(p *player.Player, defID string) {
	if g.TreasureMap == nil {
		return
	}

	filled, newLines, fullDone := g.TreasureMap.RecordKill(p.ID, defID)
	if !filled {
		return // 此目標物不在地圖上，或已填滿
	}

	betDef := data.GetBetDef(p.BetLevel)
	if betDef == nil {
		return
	}

	// 發送地圖狀態更新
	g.sendTreasureMapUpdate(p)

	// 處理完成的行/列/對角線
	for _, line := range newLines {
		reward := treasuremap.CalcLineReward(betDef.BetCost)
		p.AddCoins(reward)

		lineMsg := lineTypeToMessage(line.Type)
		log.Printf("[TreasureMap] player=%s completed %s reward=%d", p.ID, line.Type, reward)

		g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgTreasureMapLine,
			Payload: ws.TreasureMapLinePayload{
				LineType:   line.Type,
				Reward:     reward,
				NewBalance: p.Coins,
				Message:    lineMsg,
			},
		})

		// 動態牆：完成行/列/對角線（uncommon 稀有度）
		go g.notifyFeedAchievement(p, "🗺️ 寶藏地圖", lineMsg, "uncommon")
	}

	// 處理完成整張地圖
	if fullDone {
		reward := treasuremap.CalcFullReward(betDef.BetCost)
		p.AddCoins(reward)

		log.Printf("[TreasureMap] player=%s completed FULL MAP reward=%d", p.ID, reward)

		g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgTreasureMapFull,
			Payload: ws.TreasureMapFullPayload{
				Reward:     reward,
				NewBalance: p.Coins,
				Message:    fmt.Sprintf("🏆 傳說寶藏！獲得 %d 金幣！", reward),
			},
		})

		// 全服公告
		g.announceBigWin(p.DisplayName, float64(reward)/float64(betDef.BetCost), reward)
		// 動態牆：完成整張地圖（legendary 稀有度）
		go g.notifyFeedAchievement(p, "🏆 傳說寶藏", "集滿寶藏地圖！", "legendary")
	}
}

// handleGetTreasureMap 處理查詢寶藏地圖請求（Client→Server）
func (g *Game) handleGetTreasureMap(p *player.Player) {
	g.sendTreasureMapUpdate(p)
}

// sendTreasureMapUpdate 發送寶藏地圖狀態給玩家
func (g *Game) sendTreasureMapUpdate(p *player.Player) {
	if g.TreasureMap == nil {
		return
	}

	snap := g.TreasureMap.GetSnapshot(p.ID)
	cells := buildTreasureMapCells(snap)

	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgTreasureMapUpdate,
		Payload: ws.TreasureMapUpdatePayload{
			Cells:       cells,
			FilledCount: snap.FilledCount,
			LinesCount:  len(snap.Lines),
			FullDone:    snap.FullDone,
			Date:        snap.Date,
		},
	})
}

// buildTreasureMapCells 建立地圖格子 Payload 列表
func buildTreasureMapCells(snap *treasuremap.PlayerMap) []ws.TreasureMapCellPayload {
	cells := make([]ws.TreasureMapCellPayload, 0, treasuremap.GridSize*treasuremap.GridSize)
	for r := 0; r < treasuremap.GridSize; r++ {
		for c := 0; c < treasuremap.GridSize; c++ {
			def := treasuremap.GetCellDef(r, c)
			if def == nil {
				continue
			}
			cells = append(cells, ws.TreasureMapCellPayload{
				Row:    r,
				Col:    c,
				DefID:  def.DefID,
				Name:   def.Name,
				Icon:   def.Icon,
				Filled: snap.Cells[r][c],
			})
		}
	}
	return cells
}

// lineTypeToMessage 將行/列/對角線類型轉換為顯示訊息
func lineTypeToMessage(lineType string) string {
	switch lineType {
	case "row0":
		return "🗺️ 完成第一行！"
	case "row1":
		return "🗺️ 完成第二行！"
	case "row2":
		return "🗺️ 完成第三行！"
	case "col0":
		return "🗺️ 完成第一列！"
	case "col1":
		return "🗺️ 完成第二列！"
	case "col2":
		return "🗺️ 完成第三列！"
	case "diag0":
		return "🗺️ 完成對角線！"
	case "diag1":
		return "🗺️ 完成反對角線！"
	default:
		return "🗺️ 完成一條線！"
	}
}
