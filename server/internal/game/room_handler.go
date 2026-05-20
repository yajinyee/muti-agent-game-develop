// Package game — 房間難度系統 handler（DAY-091）
// 業界依據：Ocean King 系列多難度房間是 2026 年捕魚機標配
// 4 個難度：初級/中級/高級/VIP，不同獎勵倍率和 Jackpot 累積速度
package game

import (
	"fmt"
	"log"

	"digital-twin/server/internal/player"
	"digital-twin/server/internal/room"
	"digital-twin/server/internal/ws"
)

// getRoomDifficulty 取得玩家當前房間難度定義
func (g *Game) getRoomDifficulty(p *player.Player) *room.DifficultyDef {
	diffID := room.Difficulty(p.GetRoomDifficulty())
	return room.GetDifficulty(diffID)
}

// getRoomRewardMult 取得玩家當前房間的獎勵倍率
func (g *Game) getRoomRewardMult(p *player.Player) float64 {
	return g.getRoomDifficulty(p).RewardMult
}

// getRoomJackpotMult 取得玩家當前房間的 Jackpot 貢獻倍率
func (g *Game) getRoomJackpotMult(p *player.Player) float64 {
	return g.getRoomDifficulty(p).JackpotMult
}

// handleGetRoomList 處理查詢房間列表請求（Client → Server）
func (g *Game) handleGetRoomList(playerID string) {
	g.mu.RLock()
	p, ok := g.Players[playerID]
	g.mu.RUnlock()
	if !ok {
		return
	}

	currentDiff := p.GetRoomDifficulty()
	playerCoins := p.GetCoins()

	// 統計各難度房間的玩家數
	diffCounts := make(map[string]int)
	g.mu.RLock()
	for _, pl := range g.Players {
		d := pl.GetRoomDifficulty()
		diffCounts[d]++
	}
	g.mu.RUnlock()

	// 建立房間列表
	defs := room.AllDifficulties()
	rooms := make([]ws.RoomDifficultyInfo, 0, len(defs))
	for _, def := range defs {
		count := diffCounts[string(def.ID)]
		isAvailable := count < def.MaxPlayers && playerCoins >= def.EntryFee
		rooms = append(rooms, ws.RoomDifficultyInfo{
			ID:          string(def.ID),
			Name:        def.Name,
			Icon:        def.Icon,
			Color:       def.Color,
			MinBetCost:  def.MinBetCost,
			MaxBetCost:  def.MaxBetCost,
			MaxPlayers:  def.MaxPlayers,
			PlayerCount: count,
			RewardMult:  def.RewardMult,
			JackpotMult: def.JackpotMult,
			EntryFee:    def.EntryFee,
			Description: def.Description,
			IsAvailable: isAvailable,
			IsCurrent:   string(def.ID) == currentDiff,
		})
	}

	g.Hub.Send(playerID, &ws.Message{
		Type: ws.MsgRoomList,
		Payload: ws.RoomListPayload{
			Rooms:       rooms,
			CurrentRoom: currentDiff,
		},
	})
}

// handleSwitchRoom 處理切換房間請求（Client → Server）
func (g *Game) handleSwitchRoom(playerID string, payload ws.SwitchRoomPayload) {
	g.mu.RLock()
	p, ok := g.Players[playerID]
	g.mu.RUnlock()
	if !ok {
		return
	}

	targetDiff := room.Difficulty(payload.RoomID)

	// 驗證難度 ID
	def, valid := room.Difficulties[targetDiff]
	if !valid {
		g.Hub.Send(playerID, &ws.Message{
			Type: ws.MsgRoomError,
			Payload: ws.RoomErrorPayload{
				Code:    "room_not_found",
				Message: "找不到指定房間",
			},
		})
		return
	}

	// 檢查是否已在該房間
	currentDiff := p.GetRoomDifficulty()
	playerCoins := p.GetCoins()

	if currentDiff == string(targetDiff) {
		// 已在該房間，直接回傳當前狀態
		g.Hub.Send(playerID, &ws.Message{
			Type: ws.MsgRoomSwitched,
			Payload: ws.RoomSwitchedPayload{
				RoomID:      string(targetDiff),
				RoomName:    def.Name,
				RoomIcon:    def.Icon,
				RoomColor:   def.Color,
				RewardMult:  def.RewardMult,
				JackpotMult: def.JackpotMult,
				EntryFee:    0,
				NewBalance:  playerCoins,
			},
		})
		return
	}

	// 檢查房間人數
	diffCounts := make(map[string]int)
	g.mu.RLock()
	for _, pl := range g.Players {
		d := pl.GetRoomDifficulty()
		diffCounts[d]++
	}
	g.mu.RUnlock()

	if diffCounts[string(targetDiff)] >= def.MaxPlayers {
		g.Hub.Send(playerID, &ws.Message{
			Type: ws.MsgRoomError,
			Payload: ws.RoomErrorPayload{
				Code:    "room_full",
				Message: def.Name + " 房間已滿（" + def.Icon + "），請稍後再試",
			},
		})
		return
	}

	// 檢查進場費
	entryFee := def.EntryFee
	newBalance := playerCoins
	if entryFee > 0 {
		var ok bool
		newBalance, ok = p.DeductCoins(entryFee)
		if !ok {
			g.Hub.Send(playerID, &ws.Message{
				Type: ws.MsgRoomError,
				Payload: ws.RoomErrorPayload{
					Code:    "insufficient_coins",
					Message: fmt.Sprintf("金幣不足，進入 %s 需要 %s 金幣", def.Name, formatCoins(entryFee)),
				},
			})
			return
		}
	}

	// 切換房間
	p.SetRoomDifficulty(string(targetDiff))

	log.Printf("[Room] Player %s switched to %s (entry_fee=%d, balance=%d)",
		playerID, def.Name, entryFee, newBalance)

	// 通知玩家切換成功
	g.Hub.Send(playerID, &ws.Message{
		Type: ws.MsgRoomSwitched,
		Payload: ws.RoomSwitchedPayload{
			RoomID:      string(targetDiff),
			RoomName:    def.Name,
			RoomIcon:    def.Icon,
			RoomColor:   def.Color,
			RewardMult:  def.RewardMult,
			JackpotMult: def.JackpotMult,
			EntryFee:    entryFee,
			NewBalance:  newBalance,
		},
	})

	// 更新玩家狀態（讓 Client 看到新餘額）
	g.sendPlayerUpdate(p)
}

// formatCoins 格式化金幣數量（用於錯誤訊息）
func formatCoins(n int) string {
	if n >= 10000 {
		return fmt.Sprintf("%dk", n/1000)
	}
	return fmt.Sprintf("%d", n)
}
