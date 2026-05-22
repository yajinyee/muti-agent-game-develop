// lucky_egg_fish_handler.go — 幸運彩蛋魚系統 handler（DAY-172）
// 業界依據：JILI Mega Fishing 2026「Giant Prize Fish lets you easily win great prizes,
// with the chance for 5x multipliers」+ Ocean King 2026「Egg Fish drops golden eggs
// containing random rewards — coins, multiplier boost, or weapon charge」
// 擊破 T130 後掉落 1-5 個彩蛋，每個彩蛋隨機包含：
//   - 金幣獎勵（50%）：betLevel × 5-20x
//   - 倍率加成 2x 持續 5 秒（30%）：觸發玩家 5 秒內所有擊破 ×2
//   - 特殊武器充能（20%）：隨機充能一種特殊武器
// 設計差異：與冰釣輪盤（玩家選擇停止）不同，彩蛋是「自動掉落+隨機開啟」，
// 製造「每個彩蛋都是驚喜」的期待感；與幸運星魚（固定 ×2）不同，
// 彩蛋的倍率加成是「短暫但疊加」的，多個彩蛋可以連續觸發
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	// LuckyEggMinCount 最少掉落彩蛋數
	LuckyEggMinCount = 1
	// LuckyEggMaxCount 最多掉落彩蛋數
	LuckyEggMaxCount = 5
	// LuckyEggDropIntervalMs 彩蛋掉落間隔（ms）
	LuckyEggDropIntervalMs = 200
	// LuckyEggMultBoost 彩蛋倍率加成
	LuckyEggMultBoost = 2.0
	// LuckyEggMultDurationSec 彩蛋倍率加成持續時間（秒）
	LuckyEggMultDurationSec = 5
	// LuckyEggCooldownSec 個人冷卻時間（秒）
	LuckyEggCooldownSec = 30
	// LuckyEggAnnounceThreshold 全服公告門檻（彩蛋數）
	LuckyEggAnnounceThreshold = 4
)

// luckyEggRewardType 彩蛋獎勵類型
type luckyEggRewardType string

const (
	luckyEggRewardCoins  luckyEggRewardType = "coins"   // 金幣獎勵（50%）
	luckyEggRewardMult   luckyEggRewardType = "mult"    // 倍率加成（30%）
	luckyEggRewardWeapon luckyEggRewardType = "weapon"  // 武器充能（20%）
)

// luckyEggSession 玩家彩蛋倍率加成 session
type luckyEggSession struct {
	active    bool
	expiresAt time.Time
}

// luckyEggManager 幸運彩蛋魚管理器
type luckyEggManager struct {
	mu       sync.Mutex
	sessions map[string]*luckyEggSession // playerID → session
	cooldown map[string]time.Time        // playerID → 冷卻結束時間
}

// newLuckyEggManager 建立幸運彩蛋魚管理器
func newLuckyEggManager() *luckyEggManager {
	return &luckyEggManager{
		sessions: make(map[string]*luckyEggSession),
		cooldown: make(map[string]time.Time),
	}
}

// isLuckyEggFish 判斷是否為幸運彩蛋魚（T130）
func isLuckyEggFish(defID string) bool {
	return defID == "T130"
}

// getLuckyEggMult 取得玩家當前彩蛋倍率加成（供 handleKill 使用）
// 回傳 2.0（有加成）或 1.0（無加成）
func (g *Game) getLuckyEggMult(playerID string) float64 {
	g.LuckyEgg.mu.Lock()
	defer g.LuckyEgg.mu.Unlock()

	sess, ok := g.LuckyEgg.sessions[playerID]
	if !ok || !sess.active {
		return 1.0
	}
	if time.Now().After(sess.expiresAt) {
		sess.active = false
		return 1.0
	}
	return LuckyEggMultBoost
}

// activateLuckyEggMult 激活玩家彩蛋倍率加成
func (g *Game) activateLuckyEggMult(playerID string) {
	g.LuckyEgg.mu.Lock()
	defer g.LuckyEgg.mu.Unlock()

	sess, ok := g.LuckyEgg.sessions[playerID]
	if !ok {
		sess = &luckyEggSession{}
		g.LuckyEgg.sessions[playerID] = sess
	}
	// 如果已有加成，延長時間；否則重新計時
	if sess.active && time.Now().Before(sess.expiresAt) {
		sess.expiresAt = sess.expiresAt.Add(time.Duration(LuckyEggMultDurationSec) * time.Second)
	} else {
		sess.active = true
		sess.expiresAt = time.Now().Add(time.Duration(LuckyEggMultDurationSec) * time.Second)
	}
}

// deactivateLuckyEggMult 停用玩家彩蛋倍率加成（倍率結束後廣播）
func (g *Game) deactivateLuckyEggMult(playerID string, playerName string) {
	g.LuckyEgg.mu.Lock()
	sess, ok := g.LuckyEgg.sessions[playerID]
	if ok {
		sess.active = false
	}
	g.LuckyEgg.mu.Unlock()

	// 廣播倍率結束
	g.Hub.Send(playerID, &ws.Message{
		Type: ws.MsgLuckyEggFish,
		Payload: ws.LuckyEggFishPayload{
			Phase:      "mult_end",
			PlayerID:   playerID,
			PlayerName: playerName,
		},
	})
}

// isLuckyEggOnCooldown 檢查玩家是否在冷卻中
func (m *luckyEggManager) isOnCooldown(playerID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	cd, ok := m.cooldown[playerID]
	if !ok {
		return false
	}
	return time.Now().Before(cd)
}

// setLuckyEggCooldown 設定玩家冷卻
func (m *luckyEggManager) setCooldown(playerID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cooldown[playerID] = time.Now().Add(time.Duration(LuckyEggCooldownSec) * time.Second)
}

// pickEggCount 加權隨機決定彩蛋數量（1-5個，低數量高機率）
func pickEggCount(rng *rand.Rand) int {
	// 1個:40%, 2個:30%, 3個:15%, 4個:10%, 5個:5%
	weights := []int{40, 30, 15, 10, 5}
	total := 0
	for _, w := range weights {
		total += w
	}
	r := rng.Intn(total)
	cumulative := 0
	for i, w := range weights {
		cumulative += w
		if r < cumulative {
			return i + 1
		}
	}
	return 1
}

// pickEggRewardType 加權隨機決定彩蛋獎勵類型
func pickEggRewardType(rng *rand.Rand) luckyEggRewardType {
	// 金幣:50%, 倍率:30%, 武器:20%
	r := rng.Intn(100)
	if r < 50 {
		return luckyEggRewardCoins
	} else if r < 80 {
		return luckyEggRewardMult
	}
	return luckyEggRewardWeapon
}

// pickCoinReward 計算金幣彩蛋獎勵（betLevel × 5-20x）
func pickCoinReward(rng *rand.Rand, betLevel int) int {
	mult := 5 + rng.Intn(16) // 5-20
	reward := betLevel * mult
	if reward < 1 {
		reward = 1
	}
	return reward
}

// tryLuckyEggFish 擊破 T130 後觸發彩蛋掉落（DAY-172）
// 由 handleKill 呼叫（在 goroutine 中執行）
func (g *Game) tryLuckyEggFish(p *player.Player, triggerID string, triggerX, triggerY float64) {
	// 個人冷卻檢查
	if g.LuckyEgg.isOnCooldown(p.ID) {
		return
	}
	g.LuckyEgg.setCooldown(p.ID)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	eggCount := pickEggCount(rng)

	log.Printf("[LuckyEgg] player=%s triggered, dropping %d eggs", p.ID, eggCount)

	// 廣播彩蛋掉落開始（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyEggFish,
		Payload: ws.LuckyEggFishPayload{
			Phase:      "egg_start",
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			EggCount:   eggCount,
			TriggerX:   triggerX,
			TriggerY:   triggerY,
		},
	})

	// 逐一掉落彩蛋（每 200ms 一個）
	var eggResults []ws.LuckyEggResult
	totalCoins := 0
	multCount := 0
	weaponCount := 0

	for i := 0; i < eggCount; i++ {
		if i > 0 {
			time.Sleep(time.Duration(LuckyEggDropIntervalMs) * time.Millisecond)
		}

		rewardType := pickEggRewardType(rng)
		result := ws.LuckyEggResult{
			EggIndex:   i,
			RewardType: string(rewardType),
		}

		switch rewardType {
		case luckyEggRewardCoins:
			coins := pickCoinReward(rng, p.BetLevel)
			p.AddReward(coins)
			totalCoins += coins
			result.CoinsReward = coins
			result.Label = fmt.Sprintf("+%d 金幣", coins)
			result.Color = "#FFD700" // 金色

		case luckyEggRewardMult:
			g.activateLuckyEggMult(p.ID)
			multCount++
			result.MultBoost = LuckyEggMultBoost
			result.DurationSec = LuckyEggMultDurationSec
			result.Label = fmt.Sprintf("×%.0f 加成 %ds", LuckyEggMultBoost, LuckyEggMultDurationSec)
			result.Color = "#FF69B4" // 粉紅色

			// 5 秒後廣播倍率結束
			go func(pid, pname string) {
				time.Sleep(time.Duration(LuckyEggMultDurationSec) * time.Second)
				g.deactivateLuckyEggMult(pid, pname)
			}(p.ID, p.DisplayName)

		case luckyEggRewardWeapon:
			// 隨機充能一種特殊武器（通知 specialweapon 系統）
			go g.notifyLuckyEggWeaponCharge(p)
			weaponCount++
			result.Label = "武器充能 ×1"
			result.Color = "#00BFFF" // 天藍色
		}

		eggResults = append(eggResults, result)

		// 廣播單個彩蛋開啟（個人）
		g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgLuckyEggFish,
			Payload: ws.LuckyEggFishPayload{
				Phase:      "egg_open",
				PlayerID:   p.ID,
				PlayerName: p.DisplayName,
				EggIndex:   i,
				EggResult:  result,
			},
		})

		log.Printf("[LuckyEgg] egg[%d] type=%s coins=%d mult=%v weapon=%v",
			i, rewardType, result.CoinsReward, result.MultBoost > 0, weaponCount > 0)
	}

	// 廣播彩蛋結果（個人）
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyEggFish,
		Payload: ws.LuckyEggFishPayload{
			Phase:       "egg_result",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			EggCount:    eggCount,
			EggResults:  eggResults,
			TotalCoins:  totalCoins,
			MultCount:   multCount,
			WeaponCount: weaponCount,
		},
	})

	// 全服廣播（≥4 個彩蛋）
	if eggCount >= LuckyEggAnnounceThreshold {
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyEggFish,
			Payload: ws.LuckyEggFishPayload{
				Phase:      "egg_broadcast",
				PlayerID:   p.ID,
				PlayerName: p.DisplayName,
				EggCount:   eggCount,
				TotalCoins: totalCoins,
				MultCount:  multCount,
			},
		})
		g.announceLuckyEggFish(p.DisplayName, eggCount, totalCoins, multCount)
	}

	log.Printf("[LuckyEgg] player=%s eggs=%d totalCoins=%d multCount=%d weaponCount=%d",
		p.ID, eggCount, totalCoins, multCount, weaponCount)
}

// notifyLuckyEggWeaponCharge 彩蛋武器充能（隨機充能一種特殊武器）
func (g *Game) notifyLuckyEggWeaponCharge(p *player.Player) {
	if g.SpecialWeapon == nil {
		return
	}
	// 使用現有的 notifySpecialWeaponCharge 機制，傳入高倍率觸發充能
	// 彩蛋武器充能 = 相當於擊破一個 30x 目標（足以觸發大多數武器的充能進度）
	g.notifySpecialWeaponCharge(p, 30.0)
}

// announceLuckyEggFish 全服公告幸運彩蛋魚（DAY-172）
func (g *Game) announceLuckyEggFish(playerName string, eggCount int, totalCoins int, multCount int) {
	msg := fmt.Sprintf("🥚 %s 的幸運彩蛋魚掉落 %d 個彩蛋！", playerName, eggCount)
	if multCount > 0 {
		msg += fmt.Sprintf(" 獲得 %d 次倍率加成！", multCount)
	}
	if totalCoins > 0 {
		msg += fmt.Sprintf(" 共 %d 金幣！", totalCoins)
	}

	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAnnouncement,
		Payload: map[string]interface{}{
			"event_type": "lucky_egg_fish",
			"message":    msg,
			"color":      "#FFD700",
			"duration":   4.0,
			"priority":   2,
		},
	})
}
