// festival_handler.go — 賽季節日活動系統 handler（DAY-109）
package game

import (
	"log"

	"digital-twin/server/internal/game/achievement"
	"digital-twin/server/internal/game/festival"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// sendFestivalState 發送節日狀態給玩家
func (g *Game) sendFestivalState(p *player.Player) {
	if g.Festival == nil {
		return
	}
	snap := g.Festival.GetSnapshot(p.ID)
	if err := g.Hub.Send(p.ID, &ws.Message{
		Type:    ws.MsgFestivalUpdate,
		Payload: snap,
	}); err != nil {
		log.Printf("[Festival] send state error: %v", err)
	}
}

// handleGetFestival 處理查詢節日狀態請求
func (g *Game) handleGetFestival(p *player.Player) {
	g.sendFestivalState(p)
}

// handleClaimFestivalTask 處理領取節日任務獎勵
func (g *Game) handleClaimFestivalTask(p *player.Player, taskID string) {
	if g.Festival == nil {
		return
	}
	coins := g.Festival.ClaimTaskReward(p.ID, taskID)
	if coins <= 0 {
		if err := g.Hub.Send(p.ID, &ws.Message{
			Type:    ws.MsgFestivalError,
			Payload: ws.FestivalErrorPayload{Message: "任務未完成或已領取"},
		}); err != nil {
			log.Printf("[Festival] send error: %v", err)
		}
		return
	}

	// 發放金幣
	p.AddCoins(coins)
	log.Printf("[Festival] player=%s claimed task=%s reward=%d coins", p.ID, taskID, coins)

	// 通知玩家
	if err := g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgFestivalTaskClaimed,
		Payload: ws.FestivalTaskClaimedPayload{
			TaskID:      taskID,
			RewardCoins: coins,
		},
	}); err != nil {
		log.Printf("[Festival] send task claimed error: %v", err)
	}

	// 更新玩家狀態
	g.sendPlayerUpdate(p)

	// 檢查是否可以領取稱號
	g.checkAndGrantFestivalTitle(p)

	// 更新節日狀態
	g.sendFestivalState(p)
}

// checkAndGrantFestivalTitle 檢查並發放節日稱號
func (g *Game) checkAndGrantFestivalTitle(p *player.Player) {
	if g.Festival == nil {
		return
	}
	titleID, titleName, titleColor, ok := g.Festival.ClaimTitle(p.ID)
	if !ok {
		return
	}

	log.Printf("[Festival] player=%s earned title=%s", p.ID, titleID)

	// 解鎖稱號（透過 achievement 系統）
	if p.Titles != nil {
		p.Titles.TryUnlockByID(achievement.TitleID(titleID))
	}

	// 通知玩家
	if err := g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgFestivalTitleEarned,
		Payload: ws.FestivalTitleEarnedPayload{
			TitleID:    titleID,
			TitleName:  titleName,
			TitleColor: titleColor,
		},
	}); err != nil {
		log.Printf("[Festival] send title earned error: %v", err)
	}
}

// notifyFestivalKill 在擊破目標後更新節日任務進度
func (g *Game) notifyFestivalKill(p *player.Player, targetID string) {
	if g.Festival == nil || !g.Festival.IsActive() {
		return
	}
	updated, completedTaskID := g.Festival.RecordKill(p.ID, targetID)
	if !updated {
		return
	}
	if completedTaskID != "" {
		log.Printf("[Festival] player=%s completed task=%s", p.ID, completedTaskID)
		// 通知任務完成（可領取獎勵）
		if err := g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgFestivalTaskReady,
			Payload: ws.FestivalTaskReadyPayload{
				TaskID: completedTaskID,
			},
		}); err != nil {
			log.Printf("[Festival] send task ready error: %v", err)
		}
	}
	// 更新節日狀態
	g.sendFestivalState(p)
}

// notifyFestivalBonus 在完成 Bonus 後更新節日任務進度
func (g *Game) notifyFestivalBonus(p *player.Player) {
	if g.Festival == nil || !g.Festival.IsActive() {
		return
	}
	updated, completedTaskID := g.Festival.RecordBonus(p.ID)
	if !updated {
		return
	}
	if completedTaskID != "" {
		if err := g.Hub.Send(p.ID, &ws.Message{
			Type:    ws.MsgFestivalTaskReady,
			Payload: ws.FestivalTaskReadyPayload{TaskID: completedTaskID},
		}); err != nil {
			log.Printf("[Festival] send task ready error: %v", err)
		}
	}
	g.sendFestivalState(p)
}

// notifyFestivalStreak 在達成連擊後更新節日任務進度
func (g *Game) notifyFestivalStreak(p *player.Player, streak int) {
	if g.Festival == nil || !g.Festival.IsActive() {
		return
	}
	updated, completedTaskID := g.Festival.RecordStreak(p.ID, streak)
	if !updated {
		return
	}
	if completedTaskID != "" {
		if err := g.Hub.Send(p.ID, &ws.Message{
			Type:    ws.MsgFestivalTaskReady,
			Payload: ws.FestivalTaskReadyPayload{TaskID: completedTaskID},
		}); err != nil {
			log.Printf("[Festival] send task ready error: %v", err)
		}
	}
	g.sendFestivalState(p)
}

// notifyFestivalJackpot 在觸發 Jackpot 後更新節日任務進度
func (g *Game) notifyFestivalJackpot(p *player.Player) {
	if g.Festival == nil || !g.Festival.IsActive() {
		return
	}
	updated, completedTaskID := g.Festival.RecordJackpot(p.ID)
	if !updated {
		return
	}
	if completedTaskID != "" {
		if err := g.Hub.Send(p.ID, &ws.Message{
			Type:    ws.MsgFestivalTaskReady,
			Payload: ws.FestivalTaskReadyPayload{TaskID: completedTaskID},
		}); err != nil {
			log.Printf("[Festival] send task ready error: %v", err)
		}
	}
	g.sendFestivalState(p)
}

// notifyFestivalChain 在觸發連鎖爆炸後更新節日任務進度
func (g *Game) notifyFestivalChain(p *player.Player) {
	if g.Festival == nil || !g.Festival.IsActive() {
		return
	}
	updated, completedTaskID := g.Festival.RecordChain(p.ID)
	if !updated {
		return
	}
	if completedTaskID != "" {
		if err := g.Hub.Send(p.ID, &ws.Message{
			Type:    ws.MsgFestivalTaskReady,
			Payload: ws.FestivalTaskReadyPayload{TaskID: completedTaskID},
		}); err != nil {
			log.Printf("[Festival] send task ready error: %v", err)
		}
	}
	g.sendFestivalState(p)
}

// getFestivalRewardMult 取得節日獎勵倍率加成
func (g *Game) getFestivalRewardMult() float64 {
	if g.Festival == nil || !g.Festival.IsActive() {
		return 1.0
	}
	return g.Festival.GetRewardMult()
}

// getFestivalJackpotMult 取得節日 Jackpot 倍率加成
func (g *Game) getFestivalJackpotMult() float64 {
	if g.Festival == nil || !g.Festival.IsActive() {
		return 1.0
	}
	return g.Festival.GetJackpotMult()
}

// getFestivalBonusChanceAdd 取得節日 Bonus 觸發率加成
func (g *Game) getFestivalBonusChanceAdd() float64 {
	if g.Festival == nil || !g.Festival.IsActive() {
		return 0.0
	}
	return g.Festival.GetBonusChanceAdd()
}

// broadcastFestivalUpdate 廣播節日狀態給所有玩家
func (g *Game) broadcastFestivalUpdate() {
	if g.Festival == nil {
		return
	}
	g.mu.RLock()
	players := make([]*player.Player, 0, len(g.Players))
	for _, p := range g.Players {
		players = append(players, p)
	}
	g.mu.RUnlock()

	for _, p := range players {
		g.sendFestivalState(p)
	}
}

// getFestivalSpecialTargetMult 取得節日限定目標物的倍率（若不是節日目標則回傳 0）
func (g *Game) getFestivalSpecialTargetMult(targetID string) float64 {
	if g.Festival == nil || !g.Festival.IsActive() {
		return 0
	}
	for _, st := range g.Festival.GetSpecialTargets() {
		if st.ID == targetID {
			return st.Multiplier
		}
	}
	return 0
}

// shouldSpawnFestivalTarget 判斷是否應該生成節日限定目標物
// 回傳 (targetID, shouldSpawn)
func (g *Game) shouldSpawnFestivalTarget() (string, bool) {
	if g.Festival == nil || !g.Festival.IsActive() {
		return "", false
	}
	targets := g.Festival.GetSpecialTargets()
	if len(targets) == 0 {
		return "", false
	}

	// 依生成機率隨機選擇
	r := festival.RandFloat64()
	cumulative := 0.0
	for _, st := range targets {
		cumulative += st.SpawnRate
		if r < cumulative {
			return st.ID, true
		}
	}
	return "", false
}

// festivalTargetName 取得節日目標物名稱
func (g *Game) festivalTargetName(targetID string) string {
	if g.Festival == nil {
		return targetID
	}
	for _, st := range g.Festival.GetSpecialTargets() {
		if st.ID == targetID {
			return st.Name
		}
	}
	return targetID
}
