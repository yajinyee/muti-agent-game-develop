// triple_lucky_fish_handler.go — 三重幸運魚系統 handler（DAY-190）
// 業界靈感：TaDa Gaming TriLuck™ 2026「trigger three different feature specifications simultaneously,
// ranging from win multipliers, jackpot bonuses, collecting all rewards, and more unique features」
// 設計：擊破 T148 後同時觸發三重效果：
//   1. 金幣雨：立即發放 betLevel × 20-50x 金幣（隨機）
//   2. 倍率加成：+50% 加成持續 12 秒（加法，與其他加成疊加）
//   3. 武器充能：龍怒/魚雷/軌道炮隨機充能一發
// 三個效果同時生效，製造「三重爽感」；全服廣播讓其他玩家看到「有人觸發了三重幸運」
// 設計差異：與幸運草魚（全服+50%）不同，三重幸運魚是「個人三重效果」，更有個人成就感；
// 與幸運彩蛋魚（隨機開彩蛋）不同，三重幸運魚是「三個效果保證觸發」，沒有隨機性，更有確定感
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"digital-twin/server/internal/game/specialweapon"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	TripleLuckyFishDefID       = "T148"
	TripleLuckyMultBonus       = 0.50  // +50% 倍率加成（加法）
	TripleLuckyMultDuration    = 12.0  // 倍率加成持續秒數
	TripleLuckyCoinMinMult     = 20.0  // 金幣雨最小倍率
	TripleLuckyCoinMaxMult     = 50.0  // 金幣雨最大倍率
	TripleLuckyCooldownSecs    = 30    // 個人冷卻秒數
)

// tripleLuckySession 三重幸運魚個人 session
type tripleLuckySession struct {
	PlayerID    string
	MultEnd     time.Time    // 倍率加成結束時間
	CooldownEnd time.Time    // 冷卻結束時間
}

// tripleLuckyFishManager 三重幸運魚管理器
type tripleLuckyFishManager struct {
	mu       sync.Mutex
	sessions map[string]*tripleLuckySession // playerID -> session
}

func newTripleLuckyFishManager() *tripleLuckyFishManager {
	return &tripleLuckyFishManager{
		sessions: make(map[string]*tripleLuckySession),
	}
}

// isTripleLuckyFish 判斷是否為三重幸運魚
func isTripleLuckyFish(defID string) bool {
	return defID == TripleLuckyFishDefID
}

// getTripleLuckyMultBonus 取得三重幸運魚倍率加成（供 handleKill 使用）
// 若玩家有活躍 session，回傳 0.5（+50%），否則回傳 0.0
func (g *Game) getTripleLuckyMultBonus(playerID string) float64 {
	if g.TripleLucky == nil {
		return 0.0
	}
	g.TripleLucky.mu.Lock()
	defer g.TripleLucky.mu.Unlock()

	sess, ok := g.TripleLucky.sessions[playerID]
	if !ok {
		return 0.0
	}
	if time.Now().After(sess.MultEnd) {
		// 倍率已過期，清除 session
		delete(g.TripleLucky.sessions, playerID)
		return 0.0
	}
	return TripleLuckyMultBonus
}

// tryTripleLuckyFish 擊破三重幸運魚後觸發三重效果（由 handleKill 呼叫）
func (g *Game) tryTripleLuckyFish(p *player.Player, killedInstanceID string) {
	if g.TripleLucky == nil {
		return
	}

	g.TripleLucky.mu.Lock()
	// 檢查冷卻
	if sess, ok := g.TripleLucky.sessions[p.ID]; ok {
		if time.Now().Before(sess.CooldownEnd) {
			g.TripleLucky.mu.Unlock()
			log.Printf("[TripleLucky] player=%s on cooldown", p.ID)
			return
		}
	}

	// 建立新 session
	now := time.Now()
	sess := &tripleLuckySession{
		PlayerID:    p.ID,
		MultEnd:     now.Add(time.Duration(TripleLuckyMultDuration * float64(time.Second))),
		CooldownEnd: now.Add(time.Duration(TripleLuckyCooldownSecs) * time.Second),
	}
	g.TripleLucky.sessions[p.ID] = sess
	g.TripleLucky.mu.Unlock()

	log.Printf("[TripleLucky] player=%s triggered triple lucky fish!", p.ID)

	// ===== 效果一：金幣雨 =====
	coinMult := TripleLuckyCoinMinMult + rand.Float64()*(TripleLuckyCoinMaxMult-TripleLuckyCoinMinMult)
	coinReward := int(float64(p.BetLevel) * coinMult)
	if coinReward < 1 {
		coinReward = 1
	}
	p.AddReward(coinReward)

	// ===== 效果二：倍率加成（已在 session 中設定，handleKill 會自動套用）=====
	// 廣播倍率加成開始
	multEndUnix := sess.MultEnd.Unix()

	// ===== 效果三：武器充能 =====
	weaponCharged := g.chargeRandomWeapon(p)

	// 廣播三重幸運觸發（個人詳細 + 全服廣播）
	g.Hub.Send(p.ID, &ws.Message{ //nolint:errcheck
		Type: ws.MsgTripleLuckyFish,
		Payload: ws.TripleLuckyFishPayload{
			Phase:         "triple_start",
			PlayerID:      p.ID,
			PlayerName:    p.DisplayName,
			CoinReward:    coinReward,
			CoinMult:      coinMult,
			MultBonus:     TripleLuckyMultBonus,
			MultDuration:  TripleLuckyMultDuration,
			MultEndUnix:   multEndUnix,
			WeaponCharged: weaponCharged,
			NewBalance:    p.Coins,
			Message:       fmt.Sprintf("🍀 三重幸運觸發！金幣+%d / 倍率+50%% / %s充能！", coinReward, weaponCharged),
		},
	})

	// 全服廣播（讓其他玩家看到）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgTripleLuckyFish,
		Payload: ws.TripleLuckyFishPayload{
			Phase:      "triple_broadcast",
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			CoinReward: coinReward,
			CoinMult:   coinMult,
			Message:    fmt.Sprintf("🍀 %s 觸發三重幸運！金幣+%d / 倍率+50%% / %s充能！", p.DisplayName, coinReward, weaponCharged),
		},
	})

	// 全服公告
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAnnouncement,
		Payload: map[string]interface{}{
			"event_type": "triple_lucky_fish",
			"message":    fmt.Sprintf("🍀 %s 觸發三重幸運魚！金幣+%d / 倍率+50%% / %s充能！", p.DisplayName, coinReward, weaponCharged),
			"color":      "#FFD700",
			"duration":   5.0,
			"priority":   3,
		},
	})

	// 12 秒後廣播倍率加成結束
	go func() {
		time.Sleep(time.Duration(TripleLuckyMultDuration) * time.Second)
		g.Hub.Send(p.ID, &ws.Message{ //nolint:errcheck
			Type: ws.MsgTripleLuckyFish,
			Payload: ws.TripleLuckyFishPayload{
				Phase:    "mult_end",
				PlayerID: p.ID,
				Message:  "🍀 三重幸運倍率加成結束",
			},
		})
		log.Printf("[TripleLucky] player=%s mult bonus ended", p.ID)
	}()

	log.Printf("[TripleLucky] player=%s coin=%d(%.0fx) mult=+50%% weapon=%s",
		p.ID, coinReward, coinMult, weaponCharged)
}

// chargeRandomWeapon 隨機充能一個特殊武器（龍怒/魚雷/軌道炮）
// 回傳充能的武器名稱（用於廣播）
func (g *Game) chargeRandomWeapon(p *player.Player) string {
	if g.SpecialWeapon == nil {
		return "龍怒"
	}

	// 隨機選擇武器類型
	weapons := []specialweapon.WeaponType{
		specialweapon.WeaponDragonWrath,
		specialweapon.WeaponTorpedo,
		specialweapon.WeaponRailgun,
	}
	weaponNames := map[specialweapon.WeaponType]string{
		specialweapon.WeaponDragonWrath: "龍怒",
		specialweapon.WeaponTorpedo:     "魚雷",
		specialweapon.WeaponRailgun:     "軌道炮",
	}

	chosen := weapons[rand.Intn(len(weapons))]
	g.SpecialWeapon.AddCharge(p.ID, chosen)
	g.sendSpecialWeaponUpdate(p, false)

	return weaponNames[chosen]
}
