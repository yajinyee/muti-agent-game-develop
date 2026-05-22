// lucky_trident_handler.go — 幸運三叉魚互動三轉盤系統（DAY-211）
// 業界依據：TaDa Gaming TriLuck™ Series 2026
// 「Within the TriLuck™ Series, you can trigger three different feature
//  specifications, ranging from win multipliers, jackpot bonuses,
//  collecting all rewards, and more unique features.」
//
// 設計：擊破 T169 後觸發「三叉幸運儀式」（個人互動）：
//   三個獨立轉盤同時旋轉，玩家依序點擊停止：
//   - 轉盤 A（金幣）：即時金幣獎勵 10x/20x/30x/50x/100x × betLevel
//   - 轉盤 B（倍率）：下一次擊破倍率加成 ×1.5/×2.0/×2.5/×3.0/×5.0（持續 15 秒）
//   - 轉盤 C（特效）：隨機特效（HP削減/免費射擊/全服廣播/小型清場）
//   三個結果同時生效，製造「三重爽感」
//   個人冷卻 25 秒；超時 12 秒自動停止（依序自動選擇）
//
// 設計差異（與其他互動系統的區別）：
//   - 巨型章魚輪盤（DAY-144）：單一轉盤，950x 大獎
//   - 長龍王雙環輪盤（DAY-194）：雙環，最高 350x
//   - 幸運三叉魚（DAY-211）：三個獨立轉盤，每個決定不同類型獎勵，三重疊加
//   - 「三個轉盤」讓玩家有「每個轉盤都是一次期待」的連續爽感
//   - 轉盤 C 的特效多樣性讓每次觸發都有驚喜感
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

// 幸運三叉魚常數（DAY-211）
const (
	LuckyTridentCooldownSec  = 25   // 個人冷卻 25 秒
	LuckyTridentTimeoutSec   = 12   // 超時自動停止 12 秒
	LuckyTridentMultDuration = 15   // 倍率加成持續 15 秒
)

// tridentWheelA 轉盤 A：金幣獎勵
type tridentWheelA struct {
	Multiplier int
	Weight     int
	Label      string
}

// tridentWheelB 轉盤 B：倍率加成
type tridentWheelB struct {
	Mult   float64
	Weight int
	Label  string
}

// tridentWheelC 轉盤 C：特效
type tridentWheelC struct {
	Effect string
	Weight int
	Label  string
	Desc   string
}

var tridentWheelATable = []tridentWheelA{
	{Multiplier: 10, Weight: 40, Label: "💰 ×10"},
	{Multiplier: 20, Weight: 30, Label: "💰 ×20"},
	{Multiplier: 30, Weight: 18, Label: "💰 ×30"},
	{Multiplier: 50, Weight: 9, Label: "💰 ×50"},
	{Multiplier: 100, Weight: 3, Label: "💰 ×100"},
}

var tridentWheelBTable = []tridentWheelB{
	{Mult: 1.5, Weight: 40, Label: "⚡ ×1.5"},
	{Mult: 2.0, Weight: 30, Label: "⚡ ×2.0"},
	{Mult: 2.5, Weight: 18, Label: "⚡ ×2.5"},
	{Mult: 3.0, Weight: 9, Label: "⚡ ×3.0"},
	{Mult: 5.0, Weight: 3, Label: "⚡ ×5.0"},
}

var tridentWheelCTable = []tridentWheelC{
	{Effect: "hp_drain", Weight: 35, Label: "🩸 HP削減", Desc: "全場目標 HP -30%"},
	{Effect: "free_shot", Weight: 30, Label: "🎯 免費射擊", Desc: "5 秒免費射擊"},
	{Effect: "broadcast", Weight: 20, Label: "📢 全服廣播", Desc: "全服廣播你的三叉結果"},
	{Effect: "mini_blast", Weight: 15, Label: "💥 小型清場", Desc: "隨機 3 個目標 80% 擊破"},
}

// luckyTridentSession 個人三叉儀式 session
type luckyTridentSession struct {
	playerID   string
	betLevel   int
	wheelAIdx  int     // -1 = 未停止
	wheelBIdx  int
	wheelCIdx  int
	startAt    time.Time
	done       bool
}

// luckyTridentManager 幸運三叉魚管理器
type luckyTridentManager struct {
	mu        sync.Mutex
	cooldowns map[string]time.Time          // playerID -> 冷卻結束時間
	sessions  map[string]*luckyTridentSession // playerID -> 進行中 session
	// 倍率加成狀態（個人）
	multBoosts  map[string]float64   // playerID -> 倍率加成
	multEnds    map[string]time.Time // playerID -> 加成結束時間
	freeShotEnds map[string]time.Time // playerID -> 免費射擊結束時間
}

func newLuckyTridentManager() *luckyTridentManager {
	return &luckyTridentManager{
		cooldowns:    make(map[string]time.Time),
		sessions:     make(map[string]*luckyTridentSession),
		multBoosts:   make(map[string]float64),
		multEnds:     make(map[string]time.Time),
		freeShotEnds: make(map[string]time.Time),
	}
}

// isLuckyTridentFish 判斷是否為幸運三叉魚（T169，DAY-211）
func isLuckyTridentFish(defID string) bool {
	return defID == "T169"
}

// getLuckyTridentMultBoost 取得個人倍率加成（供 handleKill 使用）
func (g *Game) getLuckyTridentMultBoost(playerID string) float64 {
	mgr := g.LuckyTrident
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	if end, ok := mgr.multEnds[playerID]; ok && time.Now().Before(end) {
		if mult, ok2 := mgr.multBoosts[playerID]; ok2 {
			return mult
		}
	}
	return 1.0
}

// tryLuckyTrident 擊破 T169 後觸發三叉幸運儀式
func (g *Game) tryLuckyTrident(p *player.Player, t *target.Target) {
	mgr := g.LuckyTrident
	mgr.mu.Lock()

	// 個人冷卻檢查
	if cd, ok := mgr.cooldowns[p.ID]; ok && time.Now().Before(cd) {
		mgr.mu.Unlock()
		return
	}
	// 已有進行中 session
	if _, ok := mgr.sessions[p.ID]; ok {
		mgr.mu.Unlock()
		return
	}

	mgr.cooldowns[p.ID] = time.Now().Add(LuckyTridentCooldownSec * time.Second)

	// 預先決定三個轉盤結果（公平性保證）
	aIdx := pickTridentWheelA()
	bIdx := pickTridentWheelB()
	cIdx := pickTridentWheelC()

	sess := &luckyTridentSession{
		playerID:  p.ID,
		betLevel:  p.BetLevel,
		wheelAIdx: aIdx,
		wheelBIdx: bIdx,
		wheelCIdx: cIdx,
		startAt:   time.Now(),
	}
	mgr.sessions[p.ID] = sess
	mgr.mu.Unlock()

	log.Printf("[LuckyTrident] player=%s triggered: A=%s B=%s C=%s",
		p.ID,
		tridentWheelATable[aIdx].Label,
		tridentWheelBTable[bIdx].Label,
		tridentWheelCTable[cIdx].Label,
	)

	// 廣播三叉儀式開始（個人）
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyTrident,
		Payload: ws.LuckyTridentPayload{
			Event:      "trident_start",
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			TimeoutSec: LuckyTridentTimeoutSec,
		},
	})

	// 超時自動停止
	go func() {
		time.Sleep(LuckyTridentTimeoutSec * time.Second)
		g.resolveLuckyTrident(p, true)
	}()
}

// handleLuckyTridentStop 玩家點擊停止（Client → Server）
func (g *Game) handleLuckyTridentStop(p *player.Player, msg *ws.Message) {
	var payload ws.LuckyTridentStopPayload
	if err := remarshal(msg.Payload, &payload); err != nil {
		return
	}
	g.resolveLuckyTrident(p, false)
}

// resolveLuckyTrident 結算三叉儀式
func (g *Game) resolveLuckyTrident(p *player.Player, isTimeout bool) {
	mgr := g.LuckyTrident
	mgr.mu.Lock()
	sess, ok := mgr.sessions[p.ID]
	if !ok || sess.done {
		mgr.mu.Unlock()
		return
	}
	sess.done = true
	delete(mgr.sessions, p.ID)
	mgr.mu.Unlock()

	aResult := tridentWheelATable[sess.wheelAIdx]
	bResult := tridentWheelBTable[sess.wheelBIdx]
	cResult := tridentWheelCTable[sess.wheelCIdx]

	// 轉盤 A：即時金幣獎勵
	coinReward := aResult.Multiplier * sess.betLevel
	if coinReward < 1 {
		coinReward = 1
	}
	g.mu.Lock()
	if pp, ok := g.Players[p.ID]; ok {
		pp.Coins += coinReward
	}
	g.mu.Unlock()

	// 轉盤 B：倍率加成（持續 15 秒）
	mgr.mu.Lock()
	mgr.multBoosts[p.ID] = bResult.Mult
	mgr.multEnds[p.ID] = time.Now().Add(LuckyTridentMultDuration * time.Second)
	mgr.mu.Unlock()

	log.Printf("[LuckyTrident] player=%s resolved: coin=%d mult=×%.1f effect=%s timeout=%v",
		p.ID, coinReward, bResult.Mult, cResult.Effect, isTimeout)

	// 廣播結算結果（個人）
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyTrident,
		Payload: ws.LuckyTridentPayload{
			Event:       "trident_result",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			WheelALabel: aResult.Label,
			WheelBLabel: bResult.Label,
			WheelCLabel: cResult.Label,
			CoinReward:  coinReward,
			MultBoost:   bResult.Mult,
			MultSec:     LuckyTridentMultDuration,
			Effect:      cResult.Effect,
			EffectDesc:  cResult.Desc,
			IsTimeout:   isTimeout,
		},
	})

	// 執行轉盤 C 特效
	go g.executeTridentEffect(p, cResult)

	// 大獎公告（A ≥ 50x 或 B ≥ 3.0x）
	if aResult.Multiplier >= 50 || bResult.Mult >= 3.0 {
		color := "#9B59B6"
		if aResult.Multiplier >= 100 || bResult.Mult >= 5.0 {
			color = "#FF00FF"
		}
		msg := fmt.Sprintf("🔱 %s 觸發三叉幸運！%s + %s + %s！",
			p.DisplayName, aResult.Label, bResult.Label, cResult.Label)
		ann := g.Announce.Create(announce.EventBigWin, p.DisplayName, coinReward, map[string]string{
			"message": msg,
			"color":   color,
		})
		g.broadcastAnnouncement(ann)
	}

	// 倍率加成結束後清除
	go func() {
		time.Sleep(LuckyTridentMultDuration * time.Second)
		mgr.mu.Lock()
		delete(mgr.multBoosts, p.ID)
		delete(mgr.multEnds, p.ID)
		mgr.mu.Unlock()
		g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgLuckyTrident,
			Payload: ws.LuckyTridentPayload{
				Event:    "trident_mult_end",
				PlayerID: p.ID,
			},
		})
	}()
}

// executeTridentEffect 執行轉盤 C 特效
func (g *Game) executeTridentEffect(p *player.Player, c tridentWheelC) {
	switch c.Effect {
	case "hp_drain":
		// 全場目標 HP -30%
		count := 0
		g.mu.Lock()
		for _, t := range g.Targets {
			if !t.IsAlive || isLuckyTridentFish(t.DefID) {
				continue
			}
			reduction := int(float64(t.HP) * 0.30)
			if reduction < 1 {
				reduction = 1
			}
			t.HP -= reduction
			if t.HP < 1 {
				t.HP = 1
			}
			count++
		}
		g.mu.Unlock()
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyTrident,
			Payload: ws.LuckyTridentPayload{
				Event:         "trident_effect",
				Effect:        "hp_drain",
				AffectedCount: count,
			},
		})
		log.Printf("[LuckyTrident] hp_drain: affected=%d", count)

	case "free_shot":
		// 5 秒免費射擊（不扣費）
		g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgLuckyTrident,
			Payload: ws.LuckyTridentPayload{
				Event:    "trident_effect",
				Effect:   "free_shot",
				PlayerID: p.ID,
				FreeSec:  5,
			},
		})
		// Server 端：5 秒內標記玩家免費射擊（使用 LuckyTrident manager 的 freeShotEnd map）
		mgr := g.LuckyTrident
		mgr.mu.Lock()
		mgr.freeShotEnds[p.ID] = time.Now().Add(5 * time.Second)
		mgr.mu.Unlock()
		time.Sleep(5 * time.Second)
		mgr.mu.Lock()
		delete(mgr.freeShotEnds, p.ID)
		mgr.mu.Unlock()
		g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgLuckyTrident,
			Payload: ws.LuckyTridentPayload{
				Event:  "trident_effect_end",
				Effect: "free_shot",
			},
		})

	case "broadcast":
		// 全服廣播三叉結果
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyTrident,
			Payload: ws.LuckyTridentPayload{
				Event:      "trident_broadcast",
				PlayerID:   p.ID,
				PlayerName: p.DisplayName,
				Effect:     "broadcast",
			},
		})

	case "mini_blast":
		// 隨機 3 個目標 80% 擊破機率
		type blastTarget struct {
			t      *target.Target
			reward int
		}
		var candidates []blastTarget
		g.mu.RLock()
		for _, t := range g.Targets {
			if !t.IsAlive || isLuckyTridentFish(t.DefID) {
				continue
			}
			candidates = append(candidates, blastTarget{t: t})
		}
		g.mu.RUnlock()

		// 隨機選 3 個
		rand.Shuffle(len(candidates), func(i, j int) {
			candidates[i], candidates[j] = candidates[j], candidates[i]
		})
		if len(candidates) > 3 {
			candidates = candidates[:3]
		}

		killed := 0
		totalReward := 0
		g.mu.Lock()
		for _, bt := range candidates {
			if !bt.t.IsAlive {
				continue
			}
			if rand.Float64() < 0.80 {
				bt.t.IsAlive = false
				bt.t.HP = 0
				def := bt.t.Def
				reward := int((def.MultiplierMin+def.MultiplierMax)/2.0 * 0.60)
				if reward < 1 {
					reward = 1
				}
				totalReward += reward
				killed++
			}
		}
		if pp, ok := g.Players[p.ID]; ok {
			pp.Coins += totalReward
		}
		g.mu.Unlock()

		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyTrident,
			Payload: ws.LuckyTridentPayload{
				Event:         "trident_effect",
				Effect:        "mini_blast",
				AffectedCount: killed,
				CoinReward:    totalReward,
			},
		})
		log.Printf("[LuckyTrident] mini_blast: killed=%d reward=%d", killed, totalReward)
	}
}

// pickTridentWheelA 加權隨機選轉盤 A
func pickTridentWheelA() int {
	total := 0
	for _, t := range tridentWheelATable {
		total += t.Weight
	}
	r := rand.Intn(total)
	cum := 0
	for i, t := range tridentWheelATable {
		cum += t.Weight
		if r < cum {
			return i
		}
	}
	return 0
}

// pickTridentWheelB 加權隨機選轉盤 B
func pickTridentWheelB() int {
	total := 0
	for _, t := range tridentWheelBTable {
		total += t.Weight
	}
	r := rand.Intn(total)
	cum := 0
	for i, t := range tridentWheelBTable {
		cum += t.Weight
		if r < cum {
			return i
		}
	}
	return 0
}

// pickTridentWheelC 加權隨機選轉盤 C
func pickTridentWheelC() int {
	total := 0
	for _, t := range tridentWheelCTable {
		total += t.Weight
	}
	r := rand.Intn(total)
	cum := 0
	for i, t := range tridentWheelCTable {
		cum += t.Weight
		if r < cum {
			return i
		}
	}
	return 0
}
