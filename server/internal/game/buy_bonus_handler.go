// buy_bonus_handler.go — Buy Bonus 系統（DAY-114）
// 參考 BGaming Fishing Club 2（2026-04）的 Buy Bonus 機制
// 玩家可以花費金幣直接觸發 Bonus，不需要等待勞動值累積
// 標準 Bonus：BetCost × 100（期望回報 ×60）
// TNT Bonus：BetCost × 150（期望回報 ×100，倍率加成 1.5x）
package game

import (
	"log"
	"sync"
	"time"

	"digital-twin/server/internal/data"
	"digital-twin/server/internal/game/state"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	// BuyBonusDailyLimit 每日購買上限
	BuyBonusDailyLimit = 3

	// BuyBonusStandardMultiplier 標準 Bonus 費用倍率（BetCost × 100）
	BuyBonusStandardMultiplier = 100

	// BuyBonusTNTMultiplier TNT Bonus 費用倍率（BetCost × 150）
	BuyBonusTNTMultiplier = 150

	// BuyBonusTNTBonusMult TNT Bonus 的額外倍率加成
	BuyBonusTNTBonusMult = 1.5
)

// buyBonusRecord 玩家的 Buy Bonus 使用記錄
type buyBonusRecord struct {
	mu        sync.Mutex
	dailyUsed int
	lastDate  string // "2026-05-21" 格式
	tntMult   float64 // 下次 Bonus 的 TNT 倍率加成（0=無加成）
}

// buyBonusRecords 所有玩家的記錄（playerID → record）
var (
	buyBonusRecordsMu sync.RWMutex
	buyBonusRecords   = make(map[string]*buyBonusRecord)
)

// getBuyBonusRecord 取得或建立玩家的 Buy Bonus 記錄
func getBuyBonusRecord(playerID string) *buyBonusRecord {
	buyBonusRecordsMu.RLock()
	rec, ok := buyBonusRecords[playerID]
	buyBonusRecordsMu.RUnlock()
	if ok {
		return rec
	}
	buyBonusRecordsMu.Lock()
	defer buyBonusRecordsMu.Unlock()
	rec = &buyBonusRecord{}
	buyBonusRecords[playerID] = rec
	return rec
}

// todayDate 取得今日日期字串（UTC+8）
func todayDate() string {
	loc, _ := time.LoadLocation("Asia/Taipei")
	return time.Now().In(loc).Format("2006-01-02")
}

// getDailyUsed 取得今日已使用次數（跨日自動重置）
func (r *buyBonusRecord) getDailyUsed() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	today := todayDate()
	if r.lastDate != today {
		r.dailyUsed = 0
		r.lastDate = today
	}
	return r.dailyUsed
}

// incrementDailyUsed 增加今日使用次數
func (r *buyBonusRecord) incrementDailyUsed() {
	r.mu.Lock()
	defer r.mu.Unlock()
	today := todayDate()
	if r.lastDate != today {
		r.dailyUsed = 0
		r.lastDate = today
	}
	r.dailyUsed++
}

// setTNTMult 設定下次 Bonus 的 TNT 倍率加成
func (r *buyBonusRecord) setTNTMult(mult float64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tntMult = mult
}

// consumeTNTMult 消費 TNT 倍率加成（取出並清零）
func (r *buyBonusRecord) consumeTNTMult() float64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	mult := r.tntMult
	r.tntMult = 0
	return mult
}

// handleBuyBonus 處理購買 Bonus 請求（Client → Server）
func (g *Game) handleBuyBonus(p *player.Player, msg *ws.Message) {
	var payload ws.BuyBonusPayload
	if err := remarshal(msg.Payload, &payload); err != nil {
		g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgBuyBonusError,
			Payload: ws.BuyBonusErrorPayload{
				Reason:  "invalid_payload",
				Message: "請求格式錯誤",
			},
		})
		return
	}

	bonusType := payload.BonusType
	if bonusType != "standard" && bonusType != "tnt" {
		bonusType = "standard"
	}

	// 檢查遊戲狀態（只有正常遊戲中才能購買）
	g.mu.RLock()
	currentState := g.State
	g.mu.RUnlock()

	if currentState != state.StateNormalPlay {
		g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgBuyBonusError,
			Payload: ws.BuyBonusErrorPayload{
				Reason:  "game_busy",
				Message: "遊戲進行中，請稍後再試",
			},
		})
		return
	}

	// 計算費用
	betDef := data.GetBetDef(p.BetLevel)
	var cost int
	var multBonus float64 = 1.0
	if bonusType == "tnt" {
		cost = betDef.BetCost * BuyBonusTNTMultiplier
		multBonus = BuyBonusTNTBonusMult
	} else {
		cost = betDef.BetCost * BuyBonusStandardMultiplier
	}

	// 檢查每日限制
	rec := getBuyBonusRecord(p.ID)
	dailyUsed := rec.getDailyUsed()
	if dailyUsed >= BuyBonusDailyLimit {
		g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgBuyBonusError,
			Payload: ws.BuyBonusErrorPayload{
				Reason:  "daily_limit",
				Message: "今日購買次數已達上限（3次）",
				Cost:    cost,
				Balance: p.GetCoins(),
			},
		})
		return
	}

	// 檢查金幣是否足夠
	if p.GetCoins() < cost {
		g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgBuyBonusError,
			Payload: ws.BuyBonusErrorPayload{
				Reason:  "insufficient_coins",
				Message: "金幣不足，無法購買",
				Cost:    cost,
				Balance: p.GetCoins(),
			},
		})
		return
	}

	// 扣除金幣
	p.AddCoins(-cost)
	rec.incrementDailyUsed()
	dailyLeft := BuyBonusDailyLimit - rec.getDailyUsed()

	// 如果是 TNT Bonus，設定倍率加成（在 endBonusGame 中消費）
	if bonusType == "tnt" {
		rec.setTNTMult(multBonus)
	}

	log.Printf("[BuyBonus] player=%s type=%s cost=%d balance=%d dailyLeft=%d",
		p.ID, bonusType, cost, p.GetCoins(), dailyLeft)

	// 發送購買成功通知
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgBuyBonusSuccess,
		Payload: ws.BuyBonusSuccessPayload{
			BonusType:  bonusType,
			Cost:       cost,
			NewBalance: p.GetCoins(),
			DailyLeft:  dailyLeft,
			MultBonus:  multBonus,
		},
	})

	// 觸發 Bonus（3秒後開始，讓玩家看到購買成功動畫）
	g.safeAfterFunc(1*time.Second, func() {
		g.triggerBonusReady()
	})
}

// handleGetBuyBonusStatus 查詢今日購買狀態（Client → Server）
func (g *Game) handleGetBuyBonusStatus(p *player.Player) {
	rec := getBuyBonusRecord(p.ID)
	dailyUsed := rec.getDailyUsed()
	dailyLeft := BuyBonusDailyLimit - dailyUsed

	betDef := data.GetBetDef(p.BetLevel)
	standardCost := betDef.BetCost * BuyBonusStandardMultiplier
	tntCost := betDef.BetCost * BuyBonusTNTMultiplier

	g.mu.RLock()
	currentState := g.State
	g.mu.RUnlock()

	canBuy := dailyLeft > 0 && currentState == state.StateNormalPlay

	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgBuyBonusStatus,
		Payload: ws.BuyBonusStatusPayload{
			DailyLimit:   BuyBonusDailyLimit,
			DailyUsed:    dailyUsed,
			DailyLeft:    dailyLeft,
			StandardCost: standardCost,
			TNTCost:      tntCost,
			CanBuy:       canBuy,
		},
	})
}

// getBuyBonusTNTMult 取得玩家的 TNT Bonus 倍率加成（由 endBonusGame 呼叫）
// 如果玩家購買了 TNT Bonus，回傳倍率加成並清零
func (g *Game) getBuyBonusTNTMult(playerID string) float64 {
	rec := getBuyBonusRecord(playerID)
	return rec.consumeTNTMult()
}
