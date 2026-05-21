// weathersurge_handler.go — 天氣湧現事件 handler（DAY-127）
// 業界依據：Fisch（Roblox）2026-05-21 Sovereign Surge — 特殊天氣事件讓稀有目標群湧出現
// 設計：特定天氣（暴風雨/豔陽/暴雪）觸發 30-45 秒的稀有目標群湧，全服廣播
package game

import (
	"log"
	"time"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/game/weather"
	"digital-twin/server/internal/ws"
)

// WeatherSurgeDef 天氣湧現事件定義
type WeatherSurgeDef struct {
	WeatherType   weather.WeatherType
	SurgeName     string  // 湧現名稱
	SurgeIcon     string  // 圖示
	SurgeMessage  string  // 廣播訊息
	Duration      time.Duration
	RareBonus     float64 // 稀有目標加成（加到 SpecialRatio）
	GoldBonus     float64 // 金幣魚加成（加到 HighRatio）
	TriggerChance float64 // 天氣切換時觸發機率（0.0-1.0）
	Color         string  // 橫幅顏色
}

// weatherSurgeDefs 各天氣的湧現定義
var weatherSurgeDefs = map[weather.WeatherType]*WeatherSurgeDef{
	weather.WeatherStorm: {
		WeatherType:   weather.WeatherStorm,
		SurgeName:     "暴風湧現",
		SurgeIcon:     "⛈️🌊",
		SurgeMessage:  "暴風雨帶來神秘生物！稀有目標大量湧現！",
		Duration:      40 * time.Second,
		RareBonus:     0.15, // 稀有目標出現率 +15%
		GoldBonus:     0.05, // 金幣魚出現率 +5%
		TriggerChance: 0.60, // 60% 機率觸發
		Color:         "#4A90D9",
	},
	weather.WeatherSunshine: {
		WeatherType:   weather.WeatherSunshine,
		SurgeName:     "黃金湧現",
		SurgeIcon:     "🌞✨",
		SurgeMessage:  "豔陽照耀！金幣魚和稀有目標大量出現！",
		Duration:      35 * time.Second,
		RareBonus:     0.08,
		GoldBonus:     0.20, // 金幣魚出現率 +20%
		TriggerChance: 0.70, // 70% 機率觸發
		Color:         "#FFD700",
	},
	weather.WeatherBlizzard: {
		WeatherType:   weather.WeatherBlizzard,
		SurgeName:     "冰封湧現",
		SurgeIcon:     "❄️👾",
		SurgeMessage:  "冰封海域！神秘生物從冰層中湧現！",
		Duration:      45 * time.Second,
		RareBonus:     0.20, // 稀有目標出現率 +20%（最強）
		GoldBonus:     0.10,
		TriggerChance: 0.80, // 80% 機率觸發（稀有天氣，高觸發率）
		Color:         "#A8D8EA",
	},
	weather.WeatherFog: {
		WeatherType:   weather.WeatherFog,
		SurgeName:     "迷霧湧現",
		SurgeIcon:     "🌫️👁️",
		SurgeMessage:  "濃霧中隱藏著神秘生物！稀有目標若隱若現！",
		Duration:      30 * time.Second,
		RareBonus:     0.12,
		GoldBonus:     0.03,
		TriggerChance: 0.50,
		Color:         "#B0BEC5",
	},
}

// tryTriggerWeatherSurge 天氣切換時嘗試觸發湧現事件（由 tickAndBroadcastWeather 呼叫）
func (g *Game) tryTriggerWeatherSurge(weatherType weather.WeatherType) {
	def, ok := weatherSurgeDefs[weatherType]
	if !ok {
		return // 晴天/下雨不觸發湧現
	}

	// 機率判斷
	if randFloat() > def.TriggerChance {
		return
	}

	// 如果已有湧現事件，不重複觸發
	g.mu.Lock()
	if g.weatherSurgeActive {
		g.mu.Unlock()
		return
	}
	g.weatherSurgeActive = true
	g.weatherSurgeEndAt = time.Now().Add(def.Duration)
	g.weatherSurgeRareBonus = def.RareBonus
	g.weatherSurgeGoldBonus = def.GoldBonus
	g.weatherSurgeName = def.SurgeName
	g.mu.Unlock()

	log.Printf("[WeatherSurge] Triggered: %s (rare+%.0f%%, gold+%.0f%%, %ds)",
		def.SurgeName, def.RareBonus*100, def.GoldBonus*100, int(def.Duration.Seconds()))

	// 廣播湧現開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgWeatherSurgeStart,
		Payload: ws.WeatherSurgeStartPayload{
			SurgeName:    def.SurgeName,
			SurgeIcon:    def.SurgeIcon,
			SurgeMessage: def.SurgeMessage,
			Duration:     int(def.Duration.Seconds()),
			RareBonus:    def.RareBonus,
			GoldBonus:    def.GoldBonus,
			Color:        def.Color,
		},
	})

	// 全服公告
	g.announceWeatherSurge(def.SurgeName, def.SurgeIcon)

	// 立即生成 3 個稀有目標（湧現感）
	go func() {
		for i := 0; i < 3; i++ {
			time.Sleep(time.Duration(i*500) * time.Millisecond)
			g.triggerSpecialEvent()
		}
	}()
}

// tickWeatherSurge 檢查湧現事件是否過期（由 gameLoop 每秒呼叫）
func (g *Game) tickWeatherSurge() {
	g.mu.Lock()
	if !g.weatherSurgeActive {
		g.mu.Unlock()
		return
	}
	if time.Now().Before(g.weatherSurgeEndAt) {
		g.mu.Unlock()
		return
	}
	// 湧現結束
	surgeName := g.weatherSurgeName
	g.weatherSurgeActive = false
	g.weatherSurgeRareBonus = 0
	g.weatherSurgeGoldBonus = 0
	g.weatherSurgeName = ""
	g.mu.Unlock()

	log.Printf("[WeatherSurge] Ended: %s", surgeName)

	// 廣播湧現結束
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgWeatherSurgeEnd,
		Payload: ws.WeatherSurgeEndPayload{
			SurgeName: surgeName,
			Message:   surgeName + " 結束了",
		},
	})
}

// announceWeatherSurge 全服公告天氣湧現（整合 announce 系統）
func (g *Game) announceWeatherSurge(surgeName, icon string) {
	if g.Announce == nil {
		return
	}
	ann := g.Announce.Create(announce.EventWeatherSurge, "", 0, map[string]string{
		"surge_name": surgeName,
		"surge_icon": icon,
	})
	g.broadcastAnnouncement(ann)
}

// randFloat 生成 0.0-1.0 的隨機浮點數（避免 import math/rand 衝突）
func randFloat() float64 {
	return float64(time.Now().UnixNano()%1000) / 1000.0
}
