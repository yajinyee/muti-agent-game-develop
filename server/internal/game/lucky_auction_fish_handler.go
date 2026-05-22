// lucky_auction_fish_handler.go — 幸運拍賣魚系統（DAY-217）
// 業界原創「全服競標」機制
//
// 設計：T175 幸運拍賣魚出現後，開啟「全服競標」（持續 8 秒）：
//   - 任何玩家可以「出價」（消耗 betLevel × 5 籌碼），出價最高者獲得「大獎控制權」
//   - 競標結束後，最高出價者獲得「大獎控制權」（5 秒內自動射擊最高價值目標，0.85x 倍率）
//   - 競標失敗者退還 50% 出價籌碼
//   - 若無人競標，T175 自動逃跑
//   - 個人冷卻 20 秒；全服廣播競標開始/出價/結算
//
// 設計差異：
//   - 與幸運三叉魚（個人互動三轉盤）不同，幸運拍賣魚是「全服競標」，製造「全服競爭搶標」的社交感
//   - 「出價越高越有機會贏」讓玩家有「要不要加碼」的策略決策
//   - 「失敗者退還 50%」降低競標風險，讓更多玩家願意參與
//   - 「大獎控制權」讓贏家有「我掌控全場」的英雄感
//   - 全服廣播讓所有玩家看到競標進度，製造「緊張刺激的競標氛圍」
package game

import (
	"fmt"
	"log"
	"sort"
	"sync"
	"time"

	"digital-twin/server/internal/data"
	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/game/target"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	LuckyAuctionDuration    = 8 * time.Second  // 競標持續時間
	LuckyAuctionControlSec  = 5                // 大獎控制權持續秒數
	LuckyAuctionControlMult = 0.85             // 大獎控制權射擊倍率
	LuckyAuctionRefundRate  = 0.5              // 競標失敗退還比例
	LuckyAuctionBidBase     = 5                // 出價基礎倍數（× betLevel BetCost）
	LuckyAuctionPersonalCD  = 20 * time.Second // 個人冷卻
	LuckyAuctionShotInterval = 600 * time.Millisecond // 控制權射擊間隔
)

// auctionBid 競標出價記錄
type auctionBid struct {
	playerID    string
	playerName  string
	betLevel    int
	bidAmount   int
	bidTime     time.Time
}

// luckyAuctionFishManager 幸運拍賣魚管理器
type luckyAuctionFishManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldown map[string]time.Time

	// 當前競標狀態
	auctionActive   bool
	auctionEndsAt   time.Time
	auctionInstanceID string // 拍賣魚的 instanceID

	// 競標記錄（playerID → auctionBid）
	bids map[string]*auctionBid

	// 控制權狀態
	controlActive    bool
	controlPlayerID  string
	controlPlayerName string
	controlEndsAt    time.Time
}

func newLuckyAuctionFishManager() *luckyAuctionFishManager {
	return &luckyAuctionFishManager{
		personalCooldown: make(map[string]time.Time),
		bids:             make(map[string]*auctionBid),
	}
}

// isLuckyAuctionFish 判斷是否為幸運拍賣魚
func isLuckyAuctionFish(defID string) bool {
	return defID == "T175"
}

// getLuckyAuctionControlMult 取得大獎控制權射擊倍率（供 handleKill 使用）
// 若控制權持有者正在控制期間，回傳 0.85x 乘法加成
func (g *Game) getLuckyAuctionControlMult(playerID string) float64 {
	mgr := g.LuckyAuctionFish
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	if mgr.controlActive && mgr.controlPlayerID == playerID && time.Now().Before(mgr.controlEndsAt) {
		return LuckyAuctionControlMult
	}
	return 1.0
}

// notifyLuckyAuctionFishSpawn 幸運拍賣魚生成時開啟競標
// 由 spawnTarget 呼叫
func (g *Game) notifyLuckyAuctionFishSpawn(t *target.Target) {
	mgr := g.LuckyAuctionFish
	mgr.mu.Lock()

	// 若已有競標進行中，不重複觸發
	if mgr.auctionActive {
		mgr.mu.Unlock()
		return
	}

	mgr.auctionActive = true
	mgr.auctionEndsAt = time.Now().Add(LuckyAuctionDuration)
	mgr.auctionInstanceID = t.InstanceID
	// 清空舊競標記錄
	mgr.bids = make(map[string]*auctionBid)
	mgr.mu.Unlock()

	log.Printf("[LuckyAuctionFish] auction started: instanceID=%s, ends in %v", t.InstanceID, LuckyAuctionDuration)

	// 全服廣播：競標開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyAuctionFish,
		Payload: ws.LuckyAuctionFishPayload{
			Event:      "auction_start",
			InstanceID: t.InstanceID,
			DurationSec: int(LuckyAuctionDuration.Seconds()),
			BidBase:    LuckyAuctionBidBase,
			ControlSec: LuckyAuctionControlSec,
			ControlMult: LuckyAuctionControlMult,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyAuctionFish, "", 0, map[string]string{
		"message": fmt.Sprintf("🏆 幸運拍賣魚出現！全服競標開始！出價最高者獲得 %d 秒大獎控制權（×%.2f 倍率）！",
			LuckyAuctionControlSec, LuckyAuctionControlMult),
		"color": "#FFD700",
	})
	g.broadcastAnnouncement(ann)

	// 啟動競標計時器
	go g.runLuckyAuctionTimer(t.InstanceID)
}

// handleLuckyAuctionBid 玩家出價（由 HandleMessage 呼叫）
func (g *Game) handleLuckyAuctionBid(p *player.Player) {
	mgr := g.LuckyAuctionFish
	mgr.mu.Lock()

	// 競標是否進行中
	if !mgr.auctionActive || time.Now().After(mgr.auctionEndsAt) {
		mgr.mu.Unlock()
		g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgError,
			Payload: ws.ErrorPayload{Message: "競標已結束"},
		})
		return
	}

	// 個人冷卻檢查
	if cd, ok := mgr.personalCooldown[p.ID]; ok && time.Now().Before(cd) {
		mgr.mu.Unlock()
		g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgError,
			Payload: ws.ErrorPayload{Message: "競標冷卻中"},
		})
		return
	}

	// 計算出價金額
	betDef := data.GetBetDef(p.BetLevel)
	if betDef == nil {
		mgr.mu.Unlock()
		return
	}
	bidAmount := LuckyAuctionBidBase * betDef.BetCost

	// 檢查餘額
	if p.Coins < bidAmount {
		mgr.mu.Unlock()
		g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgError,
			Payload: ws.ErrorPayload{Message: "餘額不足"},
		})
		return
	}

	// 若玩家已出價，累加（不退還舊出價）
	existingBid, alreadyBid := mgr.bids[p.ID]
	if alreadyBid {
		existingBid.bidAmount += bidAmount
		existingBid.bidTime = time.Now()
	} else {
		mgr.bids[p.ID] = &auctionBid{
			playerID:   p.ID,
			playerName: p.DisplayName,
			betLevel:   p.BetLevel,
			bidAmount:  bidAmount,
			bidTime:    time.Now(),
		}
	}

	// 設定個人冷卻（防止瘋狂出價）
	mgr.personalCooldown[p.ID] = time.Now().Add(LuckyAuctionPersonalCD)

	// 取得當前最高出價
	topBid := mgr.getTopBidLocked()
	totalBids := len(mgr.bids)
	mgr.mu.Unlock()

	// 扣除出價金額
	p.Coins -= bidAmount

	log.Printf("[LuckyAuctionFish] player=%s bid=%d (total=%d), top=%s(%d)",
		p.ID, bidAmount, totalBids, topBid.playerName, topBid.bidAmount)

	// 全服廣播：新出價
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyAuctionFish,
		Payload: ws.LuckyAuctionFishPayload{
			Event:       "auction_bid",
			PlayerName:  p.DisplayName,
			BidAmount:   bidAmount,
			TopBidder:   topBid.playerName,
			TopBidAmount: topBid.bidAmount,
			TotalBidders: totalBids,
		},
	})
}

// getTopBidLocked 取得當前最高出價（需在 mu.Lock() 內呼叫）
func (mgr *luckyAuctionFishManager) getTopBidLocked() *auctionBid {
	var top *auctionBid
	for _, bid := range mgr.bids {
		if top == nil || bid.bidAmount > top.bidAmount ||
			(bid.bidAmount == top.bidAmount && bid.bidTime.Before(top.bidTime)) {
			top = bid
		}
	}
	return top
}

// runLuckyAuctionTimer 競標計時器 goroutine
func (g *Game) runLuckyAuctionTimer(instanceID string) {
	time.Sleep(LuckyAuctionDuration)

	mgr := g.LuckyAuctionFish
	mgr.mu.Lock()

	// 確認是同一場競標
	if !mgr.auctionActive || mgr.auctionInstanceID != instanceID {
		mgr.mu.Unlock()
		return
	}

	mgr.auctionActive = false

	// 取得競標結果
	if len(mgr.bids) == 0 {
		mgr.mu.Unlock()
		log.Printf("[LuckyAuctionFish] auction ended with no bids")
		// 廣播：無人競標
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyAuctionFish,
			Payload: ws.LuckyAuctionFishPayload{
				Event: "auction_no_bid",
			},
		})
		return
	}

	// 找出最高出價者
	winner := mgr.getTopBidLocked()

	// 收集所有出價（用於退款）
	allBids := make([]*auctionBid, 0, len(mgr.bids))
	for _, bid := range mgr.bids {
		allBids = append(allBids, bid)
	}

	// 設定控制權
	mgr.controlActive = true
	mgr.controlPlayerID = winner.playerID
	mgr.controlPlayerName = winner.playerName
	mgr.controlEndsAt = time.Now().Add(time.Duration(LuckyAuctionControlSec) * time.Second)
	mgr.mu.Unlock()

	log.Printf("[LuckyAuctionFish] auction winner: player=%s bid=%d", winner.playerID, winner.bidAmount)

	// 退還失敗者 50% 出價
	g.mu.RLock()
	players := make(map[string]*player.Player, len(g.Players))
	for id, p := range g.Players {
		players[id] = p
	}
	g.mu.RUnlock()

	// 排序出價（由高到低，用於廣播）
	sort.Slice(allBids, func(i, j int) bool {
		return allBids[i].bidAmount > allBids[j].bidAmount
	})

	refundList := make([]ws.AuctionBidResult, 0, len(allBids))
	for _, bid := range allBids {
		isWinner := bid.playerID == winner.playerID
		refundList = append(refundList, ws.AuctionBidResult{
			PlayerName: bid.playerName,
			BidAmount:  bid.bidAmount,
			IsWinner:   isWinner,
		})
		if !isWinner {
			// 退還 50%
			refund := int(float64(bid.bidAmount) * LuckyAuctionRefundRate)
			if p, ok := players[bid.playerID]; ok {
				p.Coins += refund
				g.Hub.Send(bid.playerID, &ws.Message{
					Type: ws.MsgReward,
					Payload: ws.RewardPayload{
						Source:     "auction_refund",
						Amount:     refund,
						Multiplier: LuckyAuctionRefundRate,
						NewBalance: p.Coins,
					},
				})
			}
		}
	}

	// 全服廣播：競標結算
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyAuctionFish,
		Payload: ws.LuckyAuctionFishPayload{
			Event:        "auction_result",
			WinnerName:   winner.playerName,
			WinnerBid:    winner.bidAmount,
			ControlSec:   LuckyAuctionControlSec,
			ControlMult:  LuckyAuctionControlMult,
			BidResults:   refundList,
			RefundRate:   LuckyAuctionRefundRate,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyAuctionFish, winner.playerName, winner.bidAmount, map[string]string{
		"message": fmt.Sprintf("🏆 %s 以 %d 金幣贏得競標！獲得 %d 秒大獎控制權（×%.2f 倍率）！",
			winner.playerName, winner.bidAmount, LuckyAuctionControlSec, LuckyAuctionControlMult),
		"color": "#FFD700",
	})
	g.broadcastAnnouncement(ann)

	// 啟動大獎控制權自動射擊
	go g.runLuckyAuctionControl(winner.playerID, winner.playerName)
}

// runLuckyAuctionControl 大獎控制權自動射擊 goroutine
func (g *Game) runLuckyAuctionControl(playerID, playerName string) {
	ticker := time.NewTicker(LuckyAuctionShotInterval)
	defer ticker.Stop()

	endTimer := time.NewTimer(time.Duration(LuckyAuctionControlSec) * time.Second)
	defer endTimer.Stop()

	shotCount := 0
	totalReward := 0

	g.mu.RLock()
	p, ok := g.Players[playerID]
	g.mu.RUnlock()
	if !ok {
		return
	}

	for {
		select {
		case <-endTimer.C:
			// 控制權結束
			mgr := g.LuckyAuctionFish
			mgr.mu.Lock()
			mgr.controlActive = false
			mgr.mu.Unlock()

			log.Printf("[LuckyAuctionFish] control ended: player=%s shots=%d reward=%d",
				playerID, shotCount, totalReward)

			// 廣播：控制權結束
			g.Hub.Broadcast(&ws.Message{
				Type: ws.MsgLuckyAuctionFish,
				Payload: ws.LuckyAuctionFishPayload{
					Event:       "control_end",
					PlayerName:  playerName,
					ShotCount:   shotCount,
					TotalReward: totalReward,
				},
			})

			// 全服公告（≥3 次擊破才公告）
			if shotCount >= 3 {
				ann := g.Announce.Create(announce.EventLuckyAuctionFish, playerName, totalReward, map[string]string{
					"message": fmt.Sprintf("🎯 %s 大獎控制權結束！共擊破 %d 個目標，獲得 %d 金幣！",
						playerName, shotCount, totalReward),
					"color": "#FFA500",
				})
				g.broadcastAnnouncement(ann)
			}
			return

		case <-ticker.C:
			// 自動射擊最高價值目標
			reward := g.doLuckyAuctionShot(p)
			if reward > 0 {
				shotCount++
				totalReward += reward

				// 廣播：控制權射擊結果
				g.Hub.Broadcast(&ws.Message{
					Type: ws.MsgLuckyAuctionFish,
					Payload: ws.LuckyAuctionFishPayload{
						Event:      "control_shot",
						PlayerName: playerName,
						ShotReward: reward,
						ShotCount:  shotCount,
					},
				})
			}
		}
	}
}

// doLuckyAuctionShot 執行一次大獎控制權自動射擊
// 選最高價值目標，80% 擊破機率，0.85x 倍率
func (g *Game) doLuckyAuctionShot(p *player.Player) int {
	g.mu.Lock()
	defer g.mu.Unlock()

	// 找最高價值目標（排除 BOSS）
	var best *target.Target
	bestScore := 0.0
	for _, t := range g.Targets {
		if t.HP <= 0 || t.DefID == "B001" {
			continue
		}
		def, ok := data.Targets[t.DefID]
		if !ok {
			continue
		}
		score := def.MultiplierMax * (1.0 - float64(t.HP)/float64(def.HP))
		if best == nil || score > bestScore {
			best = t
			bestScore = score
		}
	}

	if best == nil {
		return 0
	}

	// 80% 擊破機率
	if randFloat() >= 0.80 {
		return 0
	}

	// 擊破目標
	def, ok := data.Targets[best.DefID]
	if !ok {
		return 0
	}

	best.HP = 0

	// 計算獎勵（0.85x 倍率）
	betDef := data.GetBetDef(p.BetLevel)
	if betDef == nil {
		return 0
	}
	mult := (def.MultiplierMin + def.MultiplierMax) / 2.0
	reward := int(float64(betDef.BetCost) * mult * LuckyAuctionControlMult)
	p.Coins += reward

	log.Printf("[LuckyAuctionFish] control shot: player=%s target=%s reward=%d",
		p.ID, best.DefID, reward)

	// 廣播目標被擊破
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgTargetKill,
		Payload: ws.TargetKillPayload{
			InstanceID: best.InstanceID,
			DefID:      best.DefID,
			KillerID:   p.ID,
			Reward:     reward,
		},
	})

	// 從場上移除
	delete(g.Targets, best.InstanceID)

	return reward
}

// notifyLuckyAuctionFishKill 玩家擊破幸運拍賣魚本身
// 若競標進行中，立即結算（提前結束競標）
func (g *Game) notifyLuckyAuctionFishKill(p *player.Player) {
	mgr := g.LuckyAuctionFish
	mgr.mu.Lock()

	if !mgr.auctionActive {
		mgr.mu.Unlock()
		return
	}

	// 提前結束競標
	mgr.auctionActive = false
	mgr.auctionEndsAt = time.Now() // 讓 timer goroutine 知道已結束
	mgr.mu.Unlock()

	log.Printf("[LuckyAuctionFish] auction fish killed by player=%s, auction ended early", p.ID)

	// 廣播：拍賣魚被擊破，競標提前結束
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyAuctionFish,
		Payload: ws.LuckyAuctionFishPayload{
			Event:      "auction_fish_killed",
			PlayerName: p.DisplayName,
		},
	})
}
