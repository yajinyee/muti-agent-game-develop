// lucky_zodiac_fate_handler.go — 幸運星座命運魚系統（DAY-259）
// 業界原創「星座命運+星座祝福+星座標記」機制
//
// 設計：擊破 T217 後，Server 隨機抽取「今日星座」（12 星座之一）：
//   - 對應星座的玩家獲得「星座祝福」（×3.0 倍率加成，10 秒）
//   - 非對應星座玩家獲得「星座庇護」（×1.5 倍率加成，5 秒）
//   - 同時場上隨機 3 個目標被「星座標記」（持續 15 秒）
//     擊破標記目標獲得 ×2.0 倍率（全服共享）
//   - 個人冷卻 28 秒；全服冷卻 45 秒
//
// 設計差異：
//   - 與龍王降臨（T212，龍息攻擊+護盾+爆發）不同，星座命運是「命運分配」，
//     讓玩家有「今天是我的幸運星座嗎？」的期待感
//   - 「星座祝福 ×3.0」讓對應星座玩家有「今天是我的幸運日！」的爽感
//   - 「星座庇護 ×1.5」確保所有玩家都有收益，不會讓非對應星座玩家感到被排除
//   - 「星座標記 ×2.0 全服共享」讓所有玩家都有「要趕快打標記目標」的動機
//   - 「12 星座隨機抽取」讓每次觸發都有不同的星座，增加多樣性和話題性
//   - 全服廣播「今日星座」讓所有玩家都知道「這次是什麼星座」，製造社交討論感
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

const (
	LuckyZodiacFatePersonalCD    = 28 * time.Second // 個人冷卻
	LuckyZodiacFateGlobalCD      = 45 * time.Second // 全服冷卻
	LuckyZodiacFateMarkDuration  = 15 * time.Second // 星座標記持續時間
	LuckyZodiacFateBlessDuration = 10 * time.Second // 星座祝福持續時間（對應星座）
	LuckyZodiacFateShieldDuration = 5 * time.Second // 星座庇護持續時間（非對應星座）
	LuckyZodiacFateBlessMult     = 3.0              // 星座祝福倍率
	LuckyZodiacFateShieldMult    = 1.5              // 星座庇護倍率
	LuckyZodiacFateMarkMult      = 2.0              // 星座標記擊破倍率
	LuckyZodiacFateMarkCount     = 3                // 星座標記目標數
)

// zodiacSigns 12 星座
var zodiacSigns = []string{
	"牡羊座", "金牛座", "雙子座", "巨蟹座",
	"獅子座", "處女座", "天秤座", "天蠍座",
	"射手座", "摩羯座", "水瓶座", "雙魚座",
}

// zodiacEmojis 星座對應 emoji
var zodiacEmojis = map[string]string{
	"牡羊座": "♈", "金牛座": "♉", "雙子座": "♊", "巨蟹座": "♋",
	"獅子座": "♌", "處女座": "♍", "天秤座": "♎", "天蠍座": "♏",
	"射手座": "♐", "摩羯座": "♑", "水瓶座": "♒", "雙魚座": "♓",
}

// zodiacColors 星座對應顏色
var zodiacColors = map[string]string{
	"牡羊座": "#FF4500", "金牛座": "#228B22", "雙子座": "#FFD700", "巨蟹座": "#87CEEB",
	"獅子座": "#FFA500", "處女座": "#9370DB", "天秤座": "#FF69B4", "天蠍座": "#8B0000",
	"射手座": "#4169E1", "摩羯座": "#696969", "水瓶座": "#00CED1", "雙魚座": "#7B68EE",
}

// zodiacMarkEntry 星座標記目標
type zodiacMarkEntry struct {
	instanceID string
	defID      string
	expiresAt  time.Time
}

// luckyZodiacFateManager 幸運星座命運魚管理器
type luckyZodiacFateManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 當前活躍的星座標記（instanceID → entry）
	activeMarks map[string]*zodiacMarkEntry

	// 當前星座祝福（playerID → expiresAt）
	blessBoosts map[string]time.Time

	// 當前星座庇護（playerID → expiresAt）
	shieldBoosts map[string]time.Time
}

func newLuckyZodiacFateManager() *luckyZodiacFateManager {
	return &luckyZodiacFateManager{
		personalCooldowns: make(map[string]time.Time),
		activeMarks:       make(map[string]*zodiacMarkEntry),
		blessBoosts:       make(map[string]time.Time),
		shieldBoosts:      make(map[string]time.Time),
	}
}

// isLuckyZodiacFateFish 判斷是否為幸運星座命運魚
func isLuckyZodiacFateFish(defID string) bool {
	return defID == "T217"
}

// isZodiacMarkTarget 判斷是否為星座標記目標
func (m *luckyZodiacFateManager) isZodiacMarkTarget(instanceID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	entry, ok := m.activeMarks[instanceID]
	if !ok {
		return false
	}
	return time.Now().Before(entry.expiresAt)
}

// removeZodiacMark 移除星座標記
func (m *luckyZodiacFateManager) removeZodiacMark(instanceID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.activeMarks, instanceID)
}

// getLuckyZodiacFateMult 取得星座倍率加成（供 handleKill 使用）
func (m *luckyZodiacFateManager) getLuckyZodiacFateMult(playerID string) float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now()
	if exp, ok := m.blessBoosts[playerID]; ok && now.Before(exp) {
		return LuckyZodiacFateBlessMult
	}
	if exp, ok := m.shieldBoosts[playerID]; ok && now.Before(exp) {
		return LuckyZodiacFateShieldMult
	}
	return 1.0
}

// tryLuckyZodiacFateFish 擊破 T217 後觸發星座命運
func (g *Game) tryLuckyZodiacFateFish(p *player.Player) {
	m := g.LuckyZodiacFate

	m.mu.Lock()
	now := time.Now()

	// 全服冷卻檢查
	if now.Before(m.globalCooldownUntil) {
		m.mu.Unlock()
		return
	}
	// 個人冷卻檢查
	if cd, ok := m.personalCooldowns[p.ID]; ok && now.Before(cd) {
		m.mu.Unlock()
		return
	}

	// 設定冷卻
	m.personalCooldowns[p.ID] = now.Add(LuckyZodiacFatePersonalCD)
	m.globalCooldownUntil = now.Add(LuckyZodiacFateGlobalCD)
	m.mu.Unlock()

	// 隨機抽取今日星座
	zodiac := zodiacSigns[rand.Intn(len(zodiacSigns))]
	emoji := zodiacEmojis[zodiac]
	color := zodiacColors[zodiac]

	log.Printf("[ZodiacFate] player=%s 觸發星座命運！今日星座：%s %s",
		p.ID, emoji, zodiac)

	// 選取星座標記目標（隨機 3 個）
	g.mu.Lock()
	var alive []*target.Target
	for _, t := range g.Targets {
		if t.IsAlive && !isLuckyZodiacFateFish(t.DefID) {
			alive = append(alive, t)
		}
	}
	g.mu.Unlock()

	// Fisher-Yates 打亂
	for i := len(alive) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		alive[i], alive[j] = alive[j], alive[i]
	}

	markCount := LuckyZodiacFateMarkCount
	if markCount > len(alive) {
		markCount = len(alive)
	}

	markedTargets := alive[:markCount]
	markExpiry := now.Add(LuckyZodiacFateMarkDuration)

	m.mu.Lock()
	for _, t := range markedTargets {
		m.activeMarks[t.InstanceID] = &zodiacMarkEntry{
			instanceID: t.InstanceID,
			defID:      t.DefID,
			expiresAt:  markExpiry,
		}
	}
	m.mu.Unlock()

	// 對所有在線玩家分配星座祝福/庇護
	g.mu.RLock()
	players := make([]*player.Player, 0, len(g.Players))
	for _, pl := range g.Players {
		players = append(players, pl)
	}
	g.mu.RUnlock()

	blessExpiry := now.Add(LuckyZodiacFateBlessDuration)
	shieldExpiry := now.Add(LuckyZodiacFateShieldDuration)

	blessedPlayers := []string{}
	shieldedPlayers := []string{}

	m.mu.Lock()
	for _, pl := range players {
		// 依玩家 ID 的 hash 決定星座（確保同一玩家每次觸發星座一致）
		playerZodiacIdx := 0
		for _, c := range pl.ID {
			playerZodiacIdx += int(c)
		}
		playerZodiac := zodiacSigns[playerZodiacIdx%len(zodiacSigns)]

		if playerZodiac == zodiac {
			m.blessBoosts[pl.ID] = blessExpiry
			blessedPlayers = append(blessedPlayers, pl.DisplayName)
		} else {
			m.shieldBoosts[pl.ID] = shieldExpiry
			shieldedPlayers = append(shieldedPlayers, pl.ID)
		}
	}
	m.mu.Unlock()

	// 建立標記目標資訊
	type markInfo struct {
		InstanceID string  `json:"instance_id"`
		DefID      string  `json:"def_id"`
		X          float64 `json:"x"`
		Y          float64 `json:"y"`
	}
	marks := make([]markInfo, 0, len(markedTargets))
	for _, t := range markedTargets {
		marks = append(marks, markInfo{
			InstanceID: t.InstanceID,
			DefID:      t.DefID,
			X:          t.X,
			Y:          t.Y,
		})
	}

	// 對每個玩家發送個人化訊息
	for _, pl := range players {
		playerZodiacIdx := 0
		for _, c := range pl.ID {
			playerZodiacIdx += int(c)
		}
		playerZodiac := zodiacSigns[playerZodiacIdx%len(zodiacSigns)]
		isBlessed := playerZodiac == zodiac

		var boostMult float64
		var boostDuration int
		var boostType string
		if isBlessed {
			boostMult = LuckyZodiacFateBlessMult
			boostDuration = int(LuckyZodiacFateBlessDuration.Seconds())
			boostType = "bless"
		} else {
			boostMult = LuckyZodiacFateShieldMult
			boostDuration = int(LuckyZodiacFateShieldDuration.Seconds())
			boostType = "shield"
		}

		_ = g.Hub.Send(pl.ID, &ws.Message{
			Type: ws.MsgLuckyZodiacFate,
			Payload: ws.LuckyZodiacFatePayload{
				Event:         "zodiac_start",
				PlayerID:      p.ID,
				PlayerName:    p.DisplayName,
				Zodiac:        zodiac,
				ZodiacEmoji:   emoji,
				ZodiacColor:   color,
				BoostType:     boostType,
				BoostMult:     boostMult,
				BoostDuration: boostDuration,
				MarkCount:     len(markedTargets),
				MarkMult:      LuckyZodiacFateMarkMult,
				MarkDuration:  int(LuckyZodiacFateMarkDuration.Seconds()),
			},
		})
	}

	// 全服廣播
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyZodiacFate,
		Payload: ws.LuckyZodiacFatePayload{
			Event:        "zodiac_broadcast",
			PlayerName:   p.DisplayName,
			Zodiac:       zodiac,
			ZodiacEmoji:  emoji,
			ZodiacColor:  color,
			BlessedCount: len(blessedPlayers),
			MarkCount:    len(markedTargets),
			MarkMult:     LuckyZodiacFateMarkMult,
		},
	})

	// 全服公告
	g.Announce.Create(announce.EventLuckyZodiacFate, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("%s %s 今日星座：%s！%d 人獲得 ×%.1f 星座祝福！%d 個目標被星座標記 ×%.1f！",
			emoji, p.DisplayName, zodiac,
			len(blessedPlayers), LuckyZodiacFateBlessMult,
			len(markedTargets), LuckyZodiacFateMarkMult),
		"color": color,
	})

	// 啟動星座標記超時清理
	go g.runZodiacMarkTimeout(markedTargets, zodiac, emoji, p.DisplayName)
}

// notifyZodiacMarkKill 星座標記目標被擊破
func (g *Game) notifyZodiacMarkKill(p *player.Player, t *target.Target) {
	m := g.LuckyZodiacFate
	m.removeZodiacMark(t.InstanceID)

	avgBet := g.getAvgBetCost()
	reward := int(float64(avgBet) * t.Multiplier * LuckyZodiacFateMarkMult)
	if reward < 1 {
		reward = 1
	}
	g.distributeRewardToAll(reward)

	log.Printf("[ZodiacFate] player=%s 擊破星座標記目標 %s！全服獎勵 %d（×%.1f）",
		p.ID, t.Def.Name, reward, LuckyZodiacFateMarkMult)

	// 廣播標記擊破
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyZodiacFate,
		Payload: ws.LuckyZodiacFatePayload{
			Event:       "zodiac_mark_kill",
			PlayerName:  p.DisplayName,
			TargetName:  t.Def.Name,
			MarkMult:    LuckyZodiacFateMarkMult,
			TotalReward: reward,
		},
	})
}

// runZodiacMarkTimeout 星座標記超時清理
func (g *Game) runZodiacMarkTimeout(targets []*target.Target, zodiac, emoji, triggerName string) {
	timer := time.NewTimer(LuckyZodiacFateMarkDuration)
	defer timer.Stop()

	select {
	case <-timer.C:
		m := g.LuckyZodiacFate
		// 清理剩餘標記
		remaining := 0
		m.mu.Lock()
		for _, t := range targets {
			if _, ok := m.activeMarks[t.InstanceID]; ok {
				delete(m.activeMarks, t.InstanceID)
				remaining++
			}
		}
		m.mu.Unlock()

		log.Printf("[ZodiacFate] 星座標記超時，清理 %d 個剩餘標記", remaining)

		// 廣播標記結束
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyZodiacFate,
			Payload: ws.LuckyZodiacFatePayload{
				Event:      "zodiac_end",
				PlayerName: triggerName,
				Zodiac:     zodiac,
				ZodiacEmoji: emoji,
			},
		})

	case <-g.stopCh:
		return
	}
}
