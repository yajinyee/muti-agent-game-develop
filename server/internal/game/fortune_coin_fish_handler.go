// fortune_coin_fish_handler.go — 幸運金幣魚即時獎勵系統（DAY-209）
// 業界依據：Galaxsys King of Ocean 2026
// 「Free Spin Fish, Captain Fish, and Money Fish trigger bonus rounds,
//  extra multipliers, and instant payouts.」
//
// 設計：擊破 T167 後立即觸發「金幣爆發」：
//   1. 加權隨機選擇即時獎勵：5x(50%)/10x(30%)/20x(15%)/50x(5%) × betLevel
//   2. 3% 機率觸發「黃金爆發」：全場所有目標 HP 降低 80%（持續 5 秒）
//   3. 個人冷卻 15 秒（快節奏設計，讓玩家頻繁觸發）
//   4. 全服廣播（讓其他玩家看到有人觸發了金幣爆發）
//
// 設計差異（與其他即時獎勵系統的區別）：
//   - 獎池龍（DAY-205）：抽 Jackpot 等級（Mini/Minor/Major/Grand），個人冷卻 60 秒
//   - 幸運金幣魚（DAY-209）：即時小獎（5-50x），個人冷卻 15 秒，節奏更快
//   - 「黃金爆發」（3%）讓玩家有「說不定這次全場清場」的期待感
//   - 15 秒個人冷卻讓玩家可以頻繁觸發，製造「快速節奏」的遊戲感
//   - 個人獎勵（不是全服共享），讓玩家有「我打到了！」的個人爽感
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/game/target"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// 幸運金幣魚常數（DAY-209）
const (
	FortuneCoinCooldownSec     = 15    // 個人冷卻 15 秒
	FortuneCoinGoldenBurstProb = 0.03  // 黃金爆發機率 3%
	FortuneCoinHPReduction     = 0.80  // 黃金爆發 HP 降低 80%
	FortuneCoinBurstDuration   = 5     // 黃金爆發持續 5 秒
)

// fortuneCoinRewardTier 即時獎勵等級
type fortuneCoinRewardTier struct {
	Multiplier int
	Weight     int
	Label      string
}

// fortuneCoinRewardTable 即時獎勵表（加權隨機）
var fortuneCoinRewardTable = []fortuneCoinRewardTier{
	{Multiplier: 5, Weight: 50, Label: "💰 ×5"},
	{Multiplier: 10, Weight: 30, Label: "💰 ×10"},
	{Multiplier: 20, Weight: 15, Label: "💰 ×20"},
	{Multiplier: 50, Weight: 5, Label: "💰 ×50"},
}

// fortuneCoinFishManager 幸運金幣魚管理器（個人冷卻）
type fortuneCoinFishManager struct {
	mu          sync.Mutex
	cooldowns   map[string]time.Time // playerID -> 冷卻結束時間
	burstActive bool                 // 黃金爆發是否正在進行
	burstEnd    time.Time            // 黃金爆發結束時間
}

func newFortuneCoinFishManager() *fortuneCoinFishManager {
	return &fortuneCoinFishManager{
		cooldowns: make(map[string]time.Time),
	}
}

// isFortuneCoinFish 判斷是否為幸運金幣魚（T167，DAY-209）
func isFortuneCoinFish(defID string) bool {
	return defID == "T167"
}

// isFortuneCoinBurstActive 黃金爆發是否正在進行（供 combat 查詢）
func (g *Game) isFortuneCoinBurstActive() bool {
	mgr := g.FortuneCoinFish
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	return mgr.burstActive && time.Now().Before(mgr.burstEnd)
}

// tryFortuneCoinFish 擊破 T167 後觸發金幣爆發
func (g *Game) tryFortuneCoinFish(p *player.Player, t *target.Target) {
	mgr := g.FortuneCoinFish
	mgr.mu.Lock()

	// 個人冷卻檢查
	if cd, ok := mgr.cooldowns[p.ID]; ok && time.Now().Before(cd) {
		mgr.mu.Unlock()
		return
	}
	mgr.cooldowns[p.ID] = time.Now().Add(FortuneCoinCooldownSec * time.Second)
	mgr.mu.Unlock()

	// 選擇即時獎勵等級（加權隨機）
	tier := pickFortuneCoinTier()
	reward := tier.Multiplier * p.BetLevel
	if reward < 1 {
		reward = 1
	}

	// 給予個人獎勵
	g.mu.Lock()
	if pp, ok := g.Players[p.ID]; ok {
		pp.Coins += reward
	}
	g.mu.Unlock()

	log.Printf("[FortuneCoin] player=%s tier=%s reward=%d", p.ID, tier.Label, reward)

	// 廣播個人金幣爆發
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgFortuneCoinFish,
		Payload: ws.FortuneCoinFishPayload{
			Event:      "coin_burst",
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			Multiplier: tier.Multiplier,
			Reward:     reward,
			Label:      tier.Label,
		},
	})

	// 全服廣播（讓其他玩家看到）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgFortuneCoinFish,
		Payload: ws.FortuneCoinFishPayload{
			Event:      "coin_broadcast",
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			Multiplier: tier.Multiplier,
			Reward:     reward,
			Label:      tier.Label,
		},
	})

	// 大獎公告（≥20x）
	if tier.Multiplier >= 20 {
		color := "#FFD700"
		if tier.Multiplier >= 50 {
			color = "#FF8C00"
		}
		msg := fmt.Sprintf("💰 %s 觸發幸運金幣魚！獲得 %s 即時獎勵！",
			p.DisplayName, tier.Label)
		ann := g.Announce.Create(announce.EventBigWin, p.DisplayName, reward, map[string]string{
			"message": msg,
			"color":   color,
		})
		g.broadcastAnnouncement(ann)
	}

	// 3% 機率觸發黃金爆發
	if rand.Float64() < FortuneCoinGoldenBurstProb {
		go g.executeFortuneCoinGoldenBurst(p)
	}
}

// executeFortuneCoinGoldenBurst 執行黃金爆發（全場 HP 降低 80%，持續 5 秒）
func (g *Game) executeFortuneCoinGoldenBurst(p *player.Player) {
	mgr := g.FortuneCoinFish
	mgr.mu.Lock()
	if mgr.burstActive {
		mgr.mu.Unlock()
		return
	}
	mgr.burstActive = true
	mgr.burstEnd = time.Now().Add(FortuneCoinBurstDuration * time.Second)
	mgr.mu.Unlock()

	log.Printf("[FortuneCoin] Golden Burst triggered by %s! HP -80%% for %ds", p.ID, FortuneCoinBurstDuration)

	// 降低全場所有目標 HP 80%
	affectedCount := 0
	g.mu.Lock()
	for _, t := range g.Targets {
		if !t.IsAlive || isFortuneCoinFish(t.DefID) {
			continue
		}
		reduction := int(float64(t.HP) * FortuneCoinHPReduction)
		if reduction < 1 {
			reduction = 1
		}
		t.HP -= reduction
		if t.HP < 1 {
			t.HP = 1
		}
		affectedCount++
	}
	g.mu.Unlock()

	log.Printf("[FortuneCoin] Golden Burst affected %d targets", affectedCount)

	// 廣播黃金爆發開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgFortuneCoinFish,
		Payload: ws.FortuneCoinFishPayload{
			Event:         "golden_burst_start",
			PlayerID:      p.ID,
			PlayerName:    p.DisplayName,
			AffectedCount: affectedCount,
			BurstSec:      FortuneCoinBurstDuration,
		},
	})

	// 全服公告
	msg := fmt.Sprintf("💥 %s 觸發黃金爆發！全場 %d 個目標 HP 降低 80%%！快打！",
		p.DisplayName, affectedCount)
	ann := g.Announce.Create(announce.EventMegaWin, p.DisplayName, 0, map[string]string{
		"message": msg,
		"color":   "#FFD700",
	})
	g.broadcastAnnouncement(ann)

	// 等待黃金爆發結束
	time.Sleep(FortuneCoinBurstDuration * time.Second)

	mgr.mu.Lock()
	mgr.burstActive = false
	mgr.mu.Unlock()

	// 廣播黃金爆發結束
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgFortuneCoinFish,
		Payload: ws.FortuneCoinFishPayload{
			Event: "golden_burst_end",
		},
	})
	log.Printf("[FortuneCoin] Golden Burst ended")
}

// pickFortuneCoinTier 加權隨機選擇獎勵等級
func pickFortuneCoinTier() fortuneCoinRewardTier {
	totalWeight := 0
	for _, t := range fortuneCoinRewardTable {
		totalWeight += t.Weight
	}
	r := rand.Intn(totalWeight)
	cumulative := 0
	for _, t := range fortuneCoinRewardTable {
		cumulative += t.Weight
		if r < cumulative {
			return t
		}
	}
	return fortuneCoinRewardTable[0]
}
