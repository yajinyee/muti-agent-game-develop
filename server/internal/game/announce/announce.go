// Package announce 全服公告系統（DAY-097）
// 當重大事件發生時，廣播全服通知，增加社交感和緊張感
package announce

import (
	"fmt"
	"sync"
	"time"
)

// EventType 公告事件類型
type EventType string

const (
	// 玩家事件
	EventJackpotWin    EventType = "jackpot_win"    // Jackpot 中獎
	EventBigWin        EventType = "big_win"         // 大獎（≥50x）
	EventMegaWin       EventType = "mega_win"        // 超大獎（≥100x）
	EventBossKill      EventType = "boss_kill"       // BOSS 擊殺
	EventStreakRecord  EventType = "streak_record"   // 連擊新記錄（≥20）
	EventPlayerJoin    EventType = "player_join"     // 玩家加入
	EventPlayerLeave   EventType = "player_leave"    // 玩家離開

	// 系統事件
	EventWeatherChange EventType = "weather_change"  // 天氣變化
	EventEventStart    EventType = "event_start"     // 限時活動開始
	EventDailyReset    EventType = "daily_reset"     // 每日重置
	EventBossWarning   EventType = "boss_warning"    // BOSS 即將出現
	EventGrandJackpot  EventType = "grand_jackpot"   // Grand Jackpot 中獎（最高優先）
	EventGoldenTime    EventType = "golden_time"     // 黃金時間開始（DAY-125）
	EventWeatherSurge  EventType = "weather_surge"   // 天氣湧現事件（DAY-127）
	EventLightningChain EventType = "lightning_chain" // 閃電鰻連鎖擊破（DAY-132）
	EventFeverMode      EventType = "fever_mode"      // 狂熱模式觸發（DAY-133）
	EventUnluckyBonus   EventType = "unlucky_bonus"   // 失敗補償觸發（DAY-135）
	EventSpeedRace      EventType = "speed_race"      // 競速獵殺開始（DAY-136）
	EventSpeedRaceWin   EventType = "speed_race_win"  // 競速獵殺第一名（DAY-136）
	EventBountyPosted   EventType = "bounty_posted"   // 懸賞發布（DAY-137）
	EventBountyClaimed  EventType = "bounty_claimed"  // 懸賞被領取（DAY-137）
	EventMultStorm      EventType = "mult_storm"      // 倍率風暴觸發（DAY-138）
	EventDualRoulette   EventType = "dual_roulette"   // 雙環輪盤高倍率（DAY-139）
	EventMegaCatch      EventType = "mega_catch"      // Mega Catch 事件觸發（DAY-140）
	EventThunderboltLobster       EventType = "thunderbolt_lobster"        // 雷霆龍蝦免費射擊觸發（DAY-150）
	EventThunderboltLobsterResult EventType = "thunderbolt_lobster_result" // 雷霆龍蝦免費射擊結果（DAY-150）
	EventRainbowPhoenix           EventType = "rainbow_phoenix"            // 彩虹鳳凰 Power Up 觸發（DAY-151）
	EventRainbowPhoenixResult     EventType = "rainbow_phoenix_result"     // 彩虹鳳凰 Power Up 結果（DAY-151）
	EventVampireBloodMoon         EventType = "vampire_blood_moon"         // 吸血鬼血月模式觸發（DAY-152）
	EventVampireKill              EventType = "vampire_kill"               // 吸血鬼血月模式被擊破（DAY-152）
	EventCrystalDragon            EventType = "crystal_dragon"             // 水晶龍地獄龍大獎觸發（DAY-153）
	EventLionDance                EventType = "lion_dance"                 // 獅子舞大獎爆發觸發（DAY-168）
	EventVortexFish               EventType = "vortex_fish"                // 漩渦魚群吸引觸發（DAY-169）
	EventFreezeBomb               EventType = "freeze_bomb"               // 冰凍炸彈魚觸發（DAY-170）
	EventIceFishing               EventType = "ice_fishing"               // 冰釣幸運輪盤觸發（DAY-171）
	EventIceFishingResult         EventType = "ice_fishing_result"        // 冰釣幸運輪盤結果（DAY-171）
	EventRainbowPrism             EventType = "rainbow_prism"             // 彩虹稜鏡魚觸發（DAY-213）
	EventGoldenAccumulator        EventType = "golden_accumulator"        // 黃金累積魚觸發（DAY-214）
	EventLuckyMirrorFish          EventType = "lucky_mirror_fish"         // 幸運鏡像魚觸發（DAY-215）
	EventCursedPoisonFish         EventType = "cursed_poison_fish"        // 詛咒毒魚觸發（DAY-216）
	EventLuckyAuctionFish         EventType = "lucky_auction_fish"        // 幸運拍賣魚觸發（DAY-217）
	EventLuckyEvolutionFish       EventType = "lucky_evolution_fish"      // 幸運進化魚觸發（DAY-218）
	EventLuckyInfectionFish       EventType = "lucky_infection_fish"      // 幸運連鎖感染魚觸發（DAY-219）
	EventLuckyRicochetFish        EventType = "lucky_ricochet_fish"       // 幸運反彈魚觸發（DAY-220）
	EventLuckyBlackHole           EventType = "lucky_black_hole"          // 幸運黑洞魚觸發（DAY-221）
	EventLuckyResonanceFish       EventType = "lucky_resonance_fish"      // 幸運共鳴魚觸發（DAY-222）
	EventLuckyTeleportFish        EventType = "lucky_teleport_fish"       // 幸運傳送魚觸發（DAY-223）
	EventLuckySplitFish           EventType = "lucky_split_fish"          // 幸運分裂魚觸發（DAY-224）
	EventLuckyChargeFish          EventType = "lucky_charge_fish"         // 幸運充能魚觸發（DAY-225）
	EventLuckyChainBombFish       EventType = "lucky_chain_bomb_fish"     // 幸運鏈鎖爆炸魚觸發（DAY-226）
	EventLuckyMirrorTimeFish      EventType = "lucky_mirror_time_fish"    // 幸運鏡像時空魚觸發（DAY-227）
	EventLuckyQuantumFish         EventType = "lucky_quantum_fish"        // 幸運量子魚觸發（DAY-228）
	EventLuckyParasiteFish        EventType = "lucky_parasite_fish"       // 幸運寄生魚觸發（DAY-229）
	EventLuckyStormFish           EventType = "lucky_storm_fish"          // 幸運風暴魚觸發（DAY-230）
	EventLuckyBoomerangFish       EventType = "lucky_boomerang_fish"      // 幸運迴旋鏢魚觸發（DAY-231）
	EventLuckyMagnetFish          EventType = "lucky_magnet_fish"         // 幸運磁力魚觸發（DAY-232）
	EventLuckyEchoFish            EventType = "lucky_echo_fish"           // 幸運回聲魚觸發（DAY-233）
	EventLuckyVortexFish          EventType = "lucky_vortex_fish"         // 幸運漩渦魚觸發（DAY-234）
	EventLuckyTimeBombFish        EventType = "lucky_time_bomb_fish"      // 幸運時間炸彈魚觸發（DAY-235）
	EventLuckyMirrorWorld         EventType = "lucky_mirror_world"        // 幸運鏡面世界魚觸發（DAY-236）
)

// Priority 公告優先級
type Priority int

const (
	PriorityLow    Priority = 1 // 低優先（玩家加入/離開）
	PriorityNormal Priority = 2 // 普通（大獎/天氣）
	PriorityHigh   Priority = 3 // 高優先（BOSS/Jackpot）
	PriorityCritical Priority = 4 // 最高（Grand Jackpot）
)

// Announcement 單筆公告
type Announcement struct {
	ID         string    `json:"id"`
	EventType  EventType `json:"event_type"`
	Priority   Priority  `json:"priority"`
	Title      string    `json:"title"`       // 公告標題
	Message    string    `json:"message"`     // 公告內容
	PlayerName string    `json:"player_name"` // 相關玩家名稱（可空）
	Amount     int       `json:"amount"`      // 相關金額（可空）
	Icon       string    `json:"icon"`        // 顯示圖示
	Color      string    `json:"color"`       // 顯示顏色
	Duration   int       `json:"duration"`    // 顯示時長（毫秒）
	CreatedAt  time.Time `json:"created_at"`
	CreatedAtMs int64    `json:"created_at_ms"`
}

// Manager 公告管理器
type Manager struct {
	mu      sync.RWMutex
	history []Announcement
	maxSize int
	counter int
}

// NewManager 建立公告管理器
func NewManager() *Manager {
	return &Manager{
		history: make([]Announcement, 0, 50),
		maxSize: 50,
	}
}

// Create 建立一筆公告
func (m *Manager) Create(eventType EventType, playerName string, amount int, extra map[string]string) Announcement {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.counter++
	now := time.Now()

	title, message, icon, color, priority, duration := m.buildContent(eventType, playerName, amount, extra)

	ann := Announcement{
		ID:          fmt.Sprintf("ann_%d_%d", now.UnixMilli(), m.counter),
		EventType:   eventType,
		Priority:    priority,
		Title:       title,
		Message:     message,
		PlayerName:  playerName,
		Amount:      amount,
		Icon:        icon,
		Color:       color,
		Duration:    duration,
		CreatedAt:   now,
		CreatedAtMs: now.UnixMilli(),
	}

	// 加入歷史
	m.history = append([]Announcement{ann}, m.history...)
	if len(m.history) > m.maxSize {
		m.history = m.history[:m.maxSize]
	}

	return ann
}

// buildContent 根據事件類型建立公告內容
func (m *Manager) buildContent(eventType EventType, playerName string, amount int, extra map[string]string) (title, message, icon, color string, priority Priority, duration int) {
	name := playerName
	if name == "" {
		name = "玩家"
	}

	switch eventType {
	case EventGrandJackpot:
		title = "👑 GRAND JACKPOT！"
		message = fmt.Sprintf("%s 中了 Grand Jackpot！獲得 🪙%d！", name, amount)
		icon = "👑"
		color = "#FF0080"
		priority = PriorityCritical
		duration = 8000

	case EventJackpotWin:
		level := "JACKPOT"
		if extra != nil {
			if l, ok := extra["level_name"]; ok {
				level = l
			}
		}
		title = fmt.Sprintf("🎰 %s 中獎！", level)
		message = fmt.Sprintf("%s 中了 %s！獲得 🪙%d！", name, level, amount)
		icon = "🎰"
		color = "#FFD700"
		priority = PriorityHigh
		duration = 5000

	case EventMegaWin:
		mult := ""
		if extra != nil {
			if m, ok := extra["multiplier"]; ok {
				mult = m + "x "
			}
		}
		title = "🌟 MEGA WIN！"
		message = fmt.Sprintf("%s 獲得 %s超大獎！🪙%d！", name, mult, amount)
		icon = "🌟"
		color = "#FF6B35"
		priority = PriorityHigh
		duration = 5000

	case EventBigWin:
		title = "✨ BIG WIN！"
		message = fmt.Sprintf("%s 獲得大獎！🪙%d！", name, amount)
		icon = "✨"
		color = "#FFD700"
		priority = PriorityNormal
		duration = 3500

	case EventBossKill:
		bossName := "BOSS"
		if extra != nil {
			if b, ok := extra["boss_name"]; ok {
				bossName = b
			}
		}
		title = "⚔️ BOSS 擊殺！"
		message = fmt.Sprintf("%s 擊敗了 %s！獲得 🪙%d！", name, bossName, amount)
		icon = "⚔️"
		color = "#FF4444"
		priority = PriorityHigh
		duration = 5000

	case EventStreakRecord:
		title = "🔥 連擊記錄！"
		message = fmt.Sprintf("%s 達成 %d 連擊！", name, amount)
		icon = "🔥"
		color = "#FF8C00"
		priority = PriorityNormal
		duration = 3000

	case EventPlayerJoin:
		title = "👋 玩家加入"
		message = fmt.Sprintf("%s 加入了遊戲！", name)
		icon = "👋"
		color = "#4CAF50"
		priority = PriorityLow
		duration = 2500

	case EventPlayerLeave:
		title = "👋 玩家離開"
		message = fmt.Sprintf("%s 離開了遊戲。", name)
		icon = "🚪"
		color = "#9E9E9E"
		priority = PriorityLow
		duration = 2000

	case EventWeatherChange:
		weatherName := "天氣變化"
		if extra != nil {
			if w, ok := extra["weather_name"]; ok {
				weatherName = w
			}
		}
		title = "🌤️ 天氣變化"
		message = fmt.Sprintf("天氣變為「%s」！", weatherName)
		icon = "🌤️"
		color = "#64B5F6"
		priority = PriorityNormal
		duration = 3000

	case EventEventStart:
		eventName := "限時活動"
		if extra != nil {
			if e, ok := extra["event_name"]; ok {
				eventName = e
			}
		}
		title = "⚡ 限時活動開始！"
		message = fmt.Sprintf("「%s」開始了！", eventName)
		icon = "⚡"
		color = "#FFC107"
		priority = PriorityNormal
		duration = 4000

	case EventBossWarning:
		title = "⚠️ BOSS 即將出現！"
		message = "強大的 BOSS 即將現身！準備好了嗎？"
		icon = "⚠️"
		color = "#FF5722"
		priority = PriorityHigh
		duration = 4000

	case EventDailyReset:
		title = "🌅 每日重置"
		message = "新的一天開始了！任務和獎勵已重置。"
		icon = "🌅"
		color = "#81C784"
		priority = PriorityNormal
		duration = 3000

	case EventGoldenTime:
		tierName := "✨ 黃金時間"
		multStr := "2.0"
		if extra != nil {
			if t, ok := extra["tier_name"]; ok {
				tierName = t
			}
			if ms, ok := extra["mult"]; ok {
				multStr = ms
			}
		}
		title = tierName + " 開始！"
		message = fmt.Sprintf("全場獎勵 ×%s！把握機會！", multStr)
		icon = "✨"
		color = "#FFD700"
		priority = PriorityHigh
		duration = 5000

	case EventWeatherSurge:
		surgeName := "天氣湧現"
		surgeIcon := "🌊"
		if extra != nil {
			if sn, ok := extra["surge_name"]; ok {
				surgeName = sn
			}
			if si, ok := extra["surge_icon"]; ok {
				surgeIcon = si
			}
		}
		title = surgeIcon + " " + surgeName + "！"
		message = "稀有目標大量湧現！快來獵殺！"
		icon = surgeIcon
		color = "#4A90D9"
		priority = PriorityHigh
		duration = 6000

	case EventLightningChain:
		kills := ""
		reward := ""
		if extra != nil {
			if k, ok := extra["kills"]; ok {
				kills = k
			}
			if r, ok := extra["reward"]; ok {
				reward = r
			}
		}
		title = "⚡ 閃電連鎖！"
		message = fmt.Sprintf("%s 閃電連鎖擊破 %s 個目標！獲得 🪙%s！", name, kills, reward)
		icon = "⚡"
		color = "#FFE066"
		priority = PriorityNormal
		duration = 4000

	case EventFeverMode:
		mult := "1.5"
		if extra != nil {
			if m, ok := extra["mult"]; ok {
				mult = m
			}
		}
		title = "🔥 狂熱模式！"
		message = fmt.Sprintf("%s 進入狂熱模式！獎勵 ×%s！", name, mult)
		icon = "🔥"
		color = "#FF4500"
		priority = PriorityNormal
		duration = 4000

	case EventUnluckyBonus:
		title = "🍀 運氣補償！"
		message = fmt.Sprintf("%s 獲得運氣補償 🪙%d！繼續加油！", name, amount)
		icon = "🍀"
		color = "#4CAF50"
		priority = PriorityNormal
		duration = 3500

	case EventSpeedRace:
		targetName := "目標"
		multStr := "?"
		if extra != nil {
			if tn, ok := extra["target_name"]; ok {
				targetName = tn
			}
			if ms, ok := extra["mult"]; ok {
				multStr = ms
			}
		}
		title = "🏆 競速獵殺！"
		message = fmt.Sprintf("搶先擊破【%s】(×%s) 獲得 3x 獎勵！", targetName, multStr)
		icon = "🏆"
		color = "#FFD700"
		priority = PriorityHigh
		duration = 4000

	case EventSpeedRaceWin:
		targetName := "目標"
		multStr := "3.0"
		if extra != nil {
			if tn, ok := extra["target_name"]; ok {
				targetName = tn
			}
			if ms, ok := extra["bonus_mult"]; ok {
				multStr = ms
			}
		}
		title = "🥇 競速第一！"
		message = fmt.Sprintf("%s 搶先擊破【%s】！獲得 ×%s 獎勵！", name, targetName, multStr)
		icon = "🥇"
		color = "#FFD700"
		priority = PriorityHigh
		duration = 5000

	case EventBountyPosted:
		targetName := "目標"
		multStr := "?"
		if extra != nil {
			if tn, ok := extra["target_name"]; ok {
				targetName = tn
			}
			if ms, ok := extra["mult"]; ok {
				multStr = ms
			}
		}
		title = "💰 懸賞發布！"
		message = fmt.Sprintf("%s 對【%s】(×%s) 懸賞 🪙%d！", name, targetName, multStr, amount)
		icon = "💰"
		color = "#FFC107"
		priority = PriorityNormal
		duration = 3500

	case EventBountyClaimed:
		targetName := "目標"
		if extra != nil {
			if tn, ok := extra["target_name"]; ok {
				targetName = tn
			}
		}
		title = "💰 懸賞領取！"
		message = fmt.Sprintf("%s 擊破懸賞目標【%s】！獲得 🪙%d！", name, targetName, amount)
		icon = "💰"
		color = "#4CAF50"
		priority = PriorityNormal
		duration = 4000

	case EventMultStorm:
		tierName := "⚡ 倍率風暴"
		tierIcon := "⚡"
		multStr := "2"
		durStr := "20"
		if extra != nil {
			if tn, ok := extra["tier_name"]; ok {
				tierName = tn
			}
			if ti, ok := extra["tier_icon"]; ok {
				tierIcon = ti
			}
			if ms, ok := extra["mult_boost"]; ok {
				multStr = ms
			}
			if ds, ok := extra["duration"]; ok {
				durStr = ds
			}
		}
		title = tierIcon + " " + tierName + "！"
		message = fmt.Sprintf("全場倍率 ×%s！持續 %s 秒！快速擊破！", multStr, durStr)
		icon = tierIcon
		color = "#FF69B4"
		priority = PriorityHigh
		duration = 6000
		_ = tierName // suppress unused warning

	case EventDualRoulette:
		playerName := name
		combinedStr := "50"
		bonusStr := "0"
		if extra != nil {
			if cs, ok := extra["combined"]; ok {
				combinedStr = cs
			}
			if bs, ok := extra["bonus_reward"]; ok {
				bonusStr = bs
			}
		}
		title = "🎡 雙環輪盤大獎！"
		message = fmt.Sprintf("%s 觸發 %sx 雙環輪盤，獲得 %s 金幣！", playerName, combinedStr, bonusStr)
		icon = "🎡"
		color = "#FFD700"
		priority = PriorityHigh
		duration = 5000

	case EventMegaCatch:
		tierName := "🎣 大豐收"
		tierIcon := "🎣"
		rewardStr := "1.5"
		durStr := "12"
		if extra != nil {
			if tn, ok := extra["tier_name"]; ok {
				tierName = tn
			}
			if ti, ok := extra["tier_icon"]; ok {
				tierIcon = ti
			}
			if rs, ok := extra["reward_boost"]; ok {
				rewardStr = rs
			}
			if ds, ok := extra["duration"]; ok {
				durStr = ds
			}
		}
		title = tierIcon + " " + tierName + "！"
		message = fmt.Sprintf("全場獎勵 ×%s！高倍率目標大量湧現！持續 %s 秒！", rewardStr, durStr)
		icon = tierIcon
		color = "#66CCFF"
		priority = PriorityHigh
		duration = 6000
		_ = tierName // suppress unused warning

	case EventThunderboltLobster:
		title = "⚡ 雷霆龍蝦！"
		message = fmt.Sprintf("⚡ %s 觸發了雷霆龍蝦！免費射擊 15 秒！", playerName)
		icon = "⚡"
		color = "#FF6600"
		priority = PriorityHigh
		duration = 5000

	case EventThunderboltLobsterResult:
		title = "⚡ 雷霆龍蝦結果！"
		message = fmt.Sprintf("⚡ %s 的雷霆龍蝦免費射擊擊破了 %d 個目標！獲得 %d 金幣！", playerName, amount, 0)
		if extra != nil {
			if reward, ok := extra["reward"]; ok {
				message = fmt.Sprintf("⚡ %s 的雷霆龍蝦免費射擊擊破了 %d 個目標！獲得 %s 金幣！", playerName, amount, reward)
			}
		}
		icon = "⚡"
		color = "#FF9900"
		priority = PriorityNormal
		duration = 4000

	case EventRainbowPhoenix:
		multStr := fmt.Sprintf("%.0f", float64(amount))
		title = "🌈 彩虹鳳凰 Power Up！"
		message = fmt.Sprintf("🌈 %s 觸發彩虹鳳凰！Power Up %sx！持續 8 秒！", playerName, multStr)
		icon = "🌈"
		color = "#FF66FF"
		priority = PriorityHigh
		duration = 5000

	case EventRainbowPhoenixResult:
		killsStr := "0"
		multStr := "6"
		if extra != nil {
			if k, ok := extra["kills"]; ok {
				killsStr = k
			}
			if m, ok := extra["mult"]; ok {
				multStr = m
			}
		}
		title = "🌈 彩虹鳳凰結果！"
		message = fmt.Sprintf("🌈 %s 的彩虹鳳凰 Power Up %sx 擊破了 %s 個目標！獲得 %d 金幣！", playerName, multStr, killsStr, amount)
		icon = "🌈"
		color = "#CC44FF"
		priority = PriorityNormal
		duration = 4000

	case EventVampireBloodMoon:
		multStr := "5"
		if extra != nil {
			if m, ok := extra["mult"]; ok {
				multStr = m
			}
		}
		title = "🩸 吸血鬼血月模式！"
		message = fmt.Sprintf("🩸 吸血鬼進入血月模式！倍率 ×%s！誰能擊破它？", multStr)
		icon = "🩸"
		color = "#CC0000"
		priority = PriorityHigh
		duration = 5000

	case EventVampireKill:
		hitsStr := "0"
		multStr := "5"
		if extra != nil {
			if h, ok := extra["hits"]; ok {
				hitsStr = h
			}
			if m, ok := extra["mult"]; ok {
				multStr = m
			}
		}
		title = "🦇 吸血鬼血月擊破！"
		message = fmt.Sprintf("🦇 %s 擊破了血月吸血鬼！命中 %s 次，×%s 倍率！獲得 %d 金幣！", playerName, hitsStr, multStr, amount)
		icon = "🦇"
		color = "#880000"
		priority = PriorityNormal
		duration = 4000

	case EventCrystalDragon:
		totalStr := ""
		if extra != nil {
			if t, ok := extra["total"]; ok {
				totalStr = t
			}
		}
		title = "🐉 地獄龍大獎！"
		if totalStr != "" {
			message = fmt.Sprintf("🐉 %s 貢獻最多水晶！地獄龍大獎！全服共獲得 %s 金幣！", playerName, totalStr)
		} else {
			message = fmt.Sprintf("🐉 %s 貢獻最多水晶！地獄龍大獎觸發！", playerName)
		}
		icon = "🐉"
		color = "#7B2FBE"
		priority = PriorityHigh
		duration = 5000

	case EventLionDance:
		multStr := ""
		if extra != nil {
			if m, ok := extra["mult"]; ok {
				multStr = m
			}
		}
		title = "🦁 獅子舞爆發！"
		if multStr != "" {
			message = fmt.Sprintf("🦁 %s 觸發獅子舞爆發！全場標記目標獲得 %sx 額外倍率！", playerName, multStr)
		} else {
			message = fmt.Sprintf("🦁 %s 觸發獅子舞爆發！快去擊破標記目標！", playerName)
		}
		icon = "🦁"
		color = "#FF6B00"
		priority = PriorityHigh
		duration = 5000

	case EventVortexFish:
		killedStr := "0"
		rewardStr := "0"
		if extra != nil {
			if k, ok := extra["killed"]; ok {
				killedStr = k
			}
			if r, ok := extra["reward"]; ok {
				rewardStr = r
			}
		}
		title = "🌀 漩渦魚大豐收！"
		message = fmt.Sprintf("🌀 %s 觸發漩渦魚！一口氣吸入 %s 個目標！獲得 %s 金幣！", playerName, killedStr, rewardStr)
		icon = "🌀"
		color = "#00BFFF"
		priority = PriorityHigh
		duration = 5000

	case EventFreezeBomb:
		title = "❄️ 冰凍炸彈魚！"
		message = fmt.Sprintf("❄️ %s 觸發冰凍炸彈魚！%d 個特殊目標被冰凍 6 秒！快去擊破！", playerName, amount)
		icon = "❄️"
		color = "#00CFFF"
		priority = PriorityHigh
		duration = 5000

	case EventIceFishing:
		multStr := fmt.Sprintf("%d", amount)
		if extra != nil {
			if m, ok := extra["mult"]; ok {
				multStr = m
			}
		}
		title = "🎣 冰釣幸運輪盤！"
		message = fmt.Sprintf("🎣 %s 觸發冰釣輪盤！獲得 ×%s 倍率加成！8 秒黃金時間！", playerName, multStr)
		icon = "🎣"
		color = "#00E5FF"
		priority = PriorityHigh
		duration = 5000

	case EventIceFishingResult:
		multStr := "?"
		killsStr := "0"
		rewardStr := "0"
		if extra != nil {
			if m, ok := extra["mult"]; ok {
				multStr = m
			}
			if k, ok := extra["kills"]; ok {
				killsStr = k
			}
			if r, ok := extra["reward"]; ok {
				rewardStr = r
			}
		}
		title = "🎣 冰釣輪盤結果！"
		message = fmt.Sprintf("🎣 %s 的 ×%s 冰釣輪盤擊破 %s 個目標！獲得 %s 金幣！", playerName, multStr, killsStr, rewardStr)
		icon = "🎣"
		color = "#80DFFF"
		priority = PriorityNormal
		duration = 4000

	case EventRainbowPrism:
		msg := fmt.Sprintf("🌈 %s 觸發彩虹稜鏡！", playerName)
		if extra != nil {
			if m, ok := extra["message"]; ok {
				msg = m
			}
		}
		c := "#FF69B4"
		if extra != nil {
			if cv, ok := extra["color"]; ok {
				c = cv
			}
		}
		title = "🌈 彩虹稜鏡！"
		message = msg
		icon = "🌈"
		color = c
		priority = PriorityNormal
		duration = 4000

	case EventGoldenAccumulator:
		msg := fmt.Sprintf("🌟 黃金累積魚觸發！", )
		if extra != nil {
			if m, ok := extra["message"]; ok {
				msg = m
			}
		}
		c := "#FFD700"
		if extra != nil {
			if cv, ok := extra["color"]; ok {
				c = cv
			}
		}
		title = "🌟 黃金累積魚！"
		message = msg
		icon = "🌟"
		color = c
		priority = PriorityHigh
		duration = 5000

	case EventLuckyMirrorFish:
		msg := "🪞 幸運鏡像魚觸發！"
		if extra != nil {
			if m, ok := extra["message"]; ok {
				msg = m
			}
		}
		c := "#00FFFF"
		if extra != nil {
			if cv, ok := extra["color"]; ok {
				c = cv
			}
		}
		title = "🪞 幸運鏡像魚！"
		message = msg
		icon = "🪞"
		color = c
		priority = PriorityHigh
		duration = 5000

	case EventCursedPoisonFish:
		msg := "☠️ 詛咒毒魚觸發！"
		if extra != nil {
			if m, ok := extra["message"]; ok {
				msg = m
			}
		}
		c := "#9B59B6"
		if extra != nil {
			if cv, ok := extra["color"]; ok {
				c = cv
			}
		}
		title = "☠️ 詛咒毒魚！"
		message = msg
		icon = "☠️"
		color = c
		priority = PriorityHigh
		duration = 5000

	case EventLuckyAuctionFish:
		msg := "🏆 幸運拍賣魚競標！"
		if extra != nil {
			if m, ok := extra["message"]; ok {
				msg = m
			}
		}
		c := "#FFD700"
		if extra != nil {
			if cv, ok := extra["color"]; ok {
				c = cv
			}
		}
		title = "🏆 幸運拍賣魚！"
		message = msg
		icon = "🏆"
		color = c
		priority = PriorityHigh
		duration = 5000

	case EventLuckyEvolutionFish:
		msg := "🌟 幸運進化魚觸發！"
		if extra != nil {
			if m, ok := extra["message"]; ok {
				msg = m
			}
		}
		c := "#00FF88"
		if extra != nil {
			if cv, ok := extra["color"]; ok {
				c = cv
			}
		}
		title = "🌟 幸運進化魚！"
		message = msg
		icon = "🌟"
		color = c
		priority = PriorityHigh
		duration = 5000

	case EventLuckyInfectionFish:
		msg := "🦠 感染蔓延！"
		if extra != nil {
			if m, ok := extra["message"]; ok {
				msg = m
			}
		}
		c := "#00FF88"
		if extra != nil {
			if cv, ok := extra["color"]; ok {
				c = cv
			}
		}
		title = "🦠 幸運感染魚！"
		message = msg
		icon = "🦠"
		color = c
		priority = PriorityNormal
		duration = 4000

	case EventLuckyRicochetFish:
		msg := "🎯 反彈模式觸發！"
		if extra != nil {
			if m, ok := extra["message"]; ok {
				msg = m
			}
		}
		c := "#FF8C00"
		if extra != nil {
			if cv, ok := extra["color"]; ok {
				c = cv
			}
		}
		title = "🎯 幸運反彈魚！"
		message = msg
		icon = "🎯"
		color = c
		priority = PriorityNormal
		duration = 4000

	case EventLuckyBlackHole:
		msg := "🌑 黑洞召喚！"
		if extra != nil {
			if m, ok := extra["message"]; ok {
				msg = m
			}
		}
		c := "#8B00FF"
		if extra != nil {
			if cv, ok := extra["color"]; ok {
				c = cv
			}
		}
		title = "🌑 幸運黑洞魚！"
		message = msg
		icon = "🌑"
		color = c
		priority = PriorityHigh
		duration = 5000

	case EventLuckyResonanceFish:
		msg := "🎵 共鳴模式觸發！"
		if extra != nil {
			if m, ok := extra["message"]; ok {
				msg = m
			}
		}
		c := "#00BFFF"
		if extra != nil {
			if cv, ok := extra["color"]; ok {
				c = cv
			}
		}
		title = "🎵 幸運共鳴魚！"
		message = msg
		icon = "🎵"
		color = c
		priority = PriorityHigh
		duration = 5000

	case EventLuckyTeleportFish:
		msg := "🌀 傳送漩渦觸發！"
		if extra != nil {
			if m, ok := extra["message"]; ok {
				msg = m
			}
		}
		c := "#9B59B6"
		if extra != nil {
			if cv, ok := extra["color"]; ok {
				c = cv
			}
		}
		title = "🌀 幸運傳送魚！"
		message = msg
		icon = "🌀"
		color = c
		priority = PriorityHigh
		duration = 5000

	case EventLuckySplitFish:
		msg := "💥 分裂爆炸觸發！"
		if extra != nil {
			if m, ok := extra["message"]; ok {
				msg = m
			}
		}
		c := "#FF6B35"
		if extra != nil {
			if cv, ok := extra["color"]; ok {
				c = cv
			}
		}
		title = "💥 幸運分裂魚！"
		message = msg
		icon = "💥"
		color = c
		priority = PriorityHigh
		duration = 5000

	case EventLuckyChargeFish:
		msg := "⚡ 充能模式觸發！"
		if extra != nil {
			if m, ok := extra["message"]; ok {
				msg = m
			}
		}
		c := "#F39C12"
		if extra != nil {
			if cv, ok := extra["color"]; ok {
				c = cv
			}
		}
		title = "⚡ 幸運充能魚！"
		message = msg
		icon = "⚡"
		color = c
		priority = PriorityHigh
		duration = 4000

	case EventLuckyChainBombFish:
		msg := "💣 鏈鎖爆炸觸發！"
		if extra != nil {
			if m, ok := extra["message"]; ok {
				msg = m
			}
		}
		c := "#FF4500"
		if extra != nil {
			if cv, ok := extra["color"]; ok {
				c = cv
			}
		}
		title = "💣 幸運鏈鎖爆炸魚！"
		message = msg
		icon = "💣"
		color = c
		priority = PriorityHigh
		duration = 4000

	case EventLuckyMirrorTimeFish:
		msg := "⏪ 時間倒流觸發！"
		if extra != nil {
			if m, ok := extra["message"]; ok {
				msg = m
			}
		}
		c := "#00BFFF"
		if extra != nil {
			if cv, ok := extra["color"]; ok {
				c = cv
			}
		}
		title = "⏪ 幸運鏡像時空魚！"
		message = msg
		icon = "⏪"
		color = c
		priority = PriorityHigh
		duration = 4000

	case EventLuckyQuantumFish:
		msg := "⚛️ 量子疊加觸發！"
		if extra != nil {
			if m, ok := extra["message"]; ok {
				msg = m
			}
		}
		c := "#9B59B6"
		if extra != nil {
			if cv, ok := extra["color"]; ok {
				c = cv
			}
		}
		title = "⚛️ 幸運量子魚！"
		message = msg
		icon = "⚛️"
		color = c
		priority = PriorityHigh
		duration = 4000

	case EventLuckyParasiteFish:
		msg := "🦠 寄生釋放觸發！"
		if extra != nil {
			if m, ok := extra["message"]; ok {
				msg = m
			}
		}
		c := "#27AE60"
		if extra != nil {
			if cv, ok := extra["color"]; ok {
				c = cv
			}
		}
		title = "🦠 幸運寄生魚！"
		message = msg
		icon = "🦠"
		color = c
		priority = PriorityHigh
		duration = 4000

	case EventLuckyStormFish:
		msg := "🌪️ 風暴觸發！"
		if extra != nil {
			if m, ok := extra["message"]; ok {
				msg = m
			}
		}
		c := "#1ABC9C"
		if extra != nil {
			if cv, ok := extra["color"]; ok {
				c = cv
			}
		}
		title = "🌪️ 幸運風暴魚！"
		message = msg
		icon = "🌪️"
		color = c
		priority = PriorityHigh
		duration = 4000

	case EventLuckyBoomerangFish:
		msg := "🪃 迴旋鏢模式觸發！"
		if extra != nil {
			if m, ok := extra["message"]; ok {
				msg = m
			}
		}
		c := "#E67E22"
		if extra != nil {
			if cv, ok := extra["color"]; ok {
				c = cv
			}
		}
		title = "🪃 幸運迴旋鏢魚！"
		message = msg
		icon = "🪃"
		color = c
		priority = PriorityHigh
		duration = 4000

	case EventLuckyMagnetFish:
		msg := "🧲 磁力場觸發！"
		if extra != nil {
			if m, ok := extra["message"]; ok {
				msg = m
			}
		}
		c := "#3498DB"
		if extra != nil {
			if cv, ok := extra["color"]; ok {
				c = cv
			}
		}
		title = "🧲 幸運磁力魚！"
		message = msg
		icon = "🧲"
		color = c
		priority = PriorityHigh
		duration = 4000

	case EventLuckyEchoFish:
		msg := "🔊 回聲模式觸發！"
		if extra != nil {
			if m, ok := extra["message"]; ok {
				msg = m
			}
		}
		c := "#9B59B6"
		if extra != nil {
			if cv, ok := extra["color"]; ok {
				c = cv
			}
		}
		title = "🔊 幸運回聲魚！"
		message = msg
		icon = "🔊"
		color = c
		priority = PriorityHigh
		duration = 3500

	case EventLuckyVortexFish:
		msg := "🌀 漩渦觸發！"
		if extra != nil {
			if m, ok := extra["message"]; ok {
				msg = m
			}
		}
		c := "#16A085"
		if extra != nil {
			if cv, ok := extra["color"]; ok {
				c = cv
			}
		}
		title = "🌀 幸運漩渦魚！"
		message = msg
		icon = "🌀"
		color = c
		priority = PriorityHigh
		duration = 4000

	case EventLuckyTimeBombFish:
		msg := "💣 時間炸彈！"
		if extra != nil {
			if m, ok := extra["message"]; ok {
				msg = m
			}
		}
		c := "#E74C3C"
		if extra != nil {
			if cv, ok := extra["color"]; ok {
				c = cv
			}
		}
		title = "💣 幸運時間炸彈魚！"
		message = msg
		icon = "💣"
		color = c
		priority = PriorityHigh
		duration = 4000

	case EventLuckyMirrorWorld:
		msg := "🪞 鏡面世界！"
		if extra != nil {
			if m, ok := extra["message"]; ok {
				msg = m
			}
		}
		c := "#8E44AD"
		if extra != nil {
			if cv, ok := extra["color"]; ok {
				c = cv
			}
		}
		title = "🪞 幸運鏡面世界魚！"
		message = msg
		icon = "🪞"
		color = c
		priority = PriorityHigh
		duration = 4000

	default:
		title = "📢 公告"
		message = "系統公告"
		icon = "📢"
		color = "#FFFFFF"
		priority = PriorityLow
		duration = 2500
	}

	return
}

// GetRecent 取得最近 n 筆公告
func (m *Manager) GetRecent(n int) []Announcement {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if n > len(m.history) {
		n = len(m.history)
	}
	result := make([]Announcement, n)
	copy(result, m.history[:n])
	return result
}

// Count 取得公告總數
func (m *Manager) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.history)
}
