// Package recommend 智慧推薦系統（DAY-110）
// 根據玩家行為模式，推薦最適合的投注等級和遊戲策略
package recommend

import (
	"fmt"
	"math"
)

// PlayerBehavior 玩家行為摘要（輸入）
type PlayerBehavior struct {
	TotalShots     int     // 總射擊次數
	TotalKills     int     // 總擊破次數
	TotalBet       int     // 總投注
	TotalReward    int     // 總獲得
	TotalBonuses   int     // Bonus 觸發次數
	TotalBossKills int     // BOSS 擊殺次數
	BestStreak     int     // 最高連擊
	BestMultiplier float64 // 最高倍率
	CurrentBetLv   int     // 當前投注等級（1-10）
	CurrentCoins   int     // 當前金幣
	VIPLevel       int     // VIP 等級（0-5）
	LoginStreak    int     // 登入連續天數
	JackpotWins    int     // Jackpot 中獎次數
}

// RecommendationType 推薦類型
type RecommendationType string

const (
	RecommendBetUp      RecommendationType = "bet_up"      // 建議提升投注
	RecommendBetDown    RecommendationType = "bet_down"    // 建議降低投注
	RecommendBetStay    RecommendationType = "bet_stay"    // 建議維持投注
	RecommendAutoMode   RecommendationType = "auto_mode"   // 建議開啟自動模式
	RecommendLockTarget RecommendationType = "lock_target" // 建議鎖定高倍目標
	RecommendBonus      RecommendationType = "bonus_focus" // 建議專注累積 Bonus
	RecommendBoss       RecommendationType = "boss_focus"  // 建議集中火力打 BOSS
	RecommendJackpot    RecommendationType = "jackpot"     // 建議高難度房間衝 Jackpot
)

// Recommendation 推薦結果
type Recommendation struct {
	Type        RecommendationType `json:"type"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Icon        string             `json:"icon"`
	Priority    int                `json:"priority"` // 1=最高
	TargetBetLv int                `json:"target_bet_lv,omitempty"` // 建議的投注等級
	Confidence  float64            `json:"confidence"` // 信心度 0.0-1.0
}

// Engine 推薦引擎
type Engine struct{}

// New 建立推薦引擎
func New() *Engine {
	return &Engine{}
}

// Analyze 分析玩家行為，回傳最多 3 條推薦
func (e *Engine) Analyze(b PlayerBehavior) []Recommendation {
	var recs []Recommendation

	// 需要足夠的資料才能推薦
	if b.TotalShots < 20 {
		return []Recommendation{
			{
				Type:        RecommendBetStay,
				Title:       "繼續探索",
				Description: "再多玩幾局，我就能給你更精準的建議！",
				Icon:        "🎯",
				Priority:    1,
				Confidence:  0.5,
			},
		}
	}

	rtp := 0.0
	if b.TotalBet > 0 {
		rtp = float64(b.TotalReward) / float64(b.TotalBet)
	}

	hitRate := 0.0
	if b.TotalShots > 0 {
		hitRate = float64(b.TotalKills) / float64(b.TotalShots)
	}

	bonusRate := 0.0
	if b.TotalKills > 0 {
		bonusRate = float64(b.TotalBonuses) / float64(b.TotalKills)
	}

	// 1. 投注等級推薦
	betRec := e.analyzeBetLevel(b, rtp)
	if betRec != nil {
		recs = append(recs, *betRec)
	}

	// 2. 遊戲策略推薦
	stratRec := e.analyzeStrategy(b, rtp, hitRate, bonusRate)
	if stratRec != nil {
		recs = append(recs, *stratRec)
	}

	// 3. 特殊機會推薦
	specialRec := e.analyzeSpecialOpportunity(b, rtp)
	if specialRec != nil {
		recs = append(recs, *specialRec)
	}

	// 最多回傳 3 條
	if len(recs) > 3 {
		recs = recs[:3]
	}
	return recs
}

// analyzeBetLevel 分析投注等級是否合適
func (e *Engine) analyzeBetLevel(b PlayerBehavior, rtp float64) *Recommendation {
	currentLv := b.CurrentBetLv
	if currentLv < 1 {
		currentLv = 1
	}

	// 金幣不足以支撐當前投注等級（少於 100 次投注）
	betCost := betLevelCost(currentLv)
	if b.CurrentCoins < betCost*100 && currentLv > 1 {
		suggestLv := currentLv - 1
		for suggestLv > 1 && b.CurrentCoins < betLevelCost(suggestLv)*100 {
			suggestLv--
		}
		return &Recommendation{
			Type:        RecommendBetDown,
			Title:       "建議降低投注",
			Description: fmt.Sprintf("金幣不足，建議降到 LV%d 延長遊戲時間", suggestLv),
			Icon:        "⬇️",
			Priority:    1,
			TargetBetLv: suggestLv,
			Confidence:  0.9,
		}
	}

	// RTP 很高（> 1.2）且金幣充足，建議提升投注
	if rtp > 1.2 && b.TotalShots >= 50 && currentLv < 10 {
		suggestLv := min(currentLv+1, 10)
		// 確保金幣足夠
		if b.CurrentCoins >= betLevelCost(suggestLv)*200 {
			confidence := math.Min(0.95, 0.6+(rtp-1.2)*0.5)
			return &Recommendation{
				Type:        RecommendBetUp,
				Title:       "手氣正旺！提升投注",
				Description: fmt.Sprintf("你的 RTP 達 %.0f%%，升到 LV%d 獎勵更豐厚", rtp*100, suggestLv),
				Icon:        "⬆️",
				Priority:    2,
				TargetBetLv: suggestLv,
				Confidence:  confidence,
			}
		}
	}

	// RTP 很低（< 0.7）且已有足夠樣本，建議降低投注
	if rtp < 0.7 && b.TotalShots >= 100 && currentLv > 1 {
		suggestLv := max(currentLv-1, 1)
		return &Recommendation{
			Type:        RecommendBetDown,
			Title:       "調整策略",
			Description: fmt.Sprintf("當前 RTP %.0f%%，降到 LV%d 減少損失", rtp*100, suggestLv),
			Icon:        "🔄",
			Priority:    2,
			TargetBetLv: suggestLv,
			Confidence:  0.75,
		}
	}

	return nil
}

// analyzeStrategy 分析遊戲策略
func (e *Engine) analyzeStrategy(b PlayerBehavior, rtp, hitRate, bonusRate float64) *Recommendation {
	// 命中率很低，建議開啟自動模式（優先級最高）
	if hitRate < 0.05 && b.TotalShots >= 50 {
		return &Recommendation{
			Type:        RecommendAutoMode,
			Title:       "開啟自動模式",
			Description: "自動模式會智慧選擇目標，提升命中效率",
			Icon:        "🤖",
			Priority:    2,
			Confidence:  0.8,
		}
	}

	// 連擊很高但倍率偏低，建議鎖定高倍目標（優先於 Bonus 和 BOSS）
	if b.BestStreak >= 10 && b.BestMultiplier < 20 {
		return &Recommendation{
			Type:        RecommendLockTarget,
			Title:       "鎖定高倍目標",
			Description: "你的連擊能力強！試著鎖定流星或寶箱怪，倍率更高",
			Icon:        "🎯",
			Priority:    2,
			Confidence:  0.75,
		}
	}

	// BOSS 擊殺次數多，建議集中火力策略
	if b.TotalBossKills >= 3 {
		return &Recommendation{
			Type:        RecommendBoss,
			Title:       "BOSS 獵人策略",
			Description: "你擅長打 BOSS！BOSS 出現時集中火力，獎勵最高 500x",
			Icon:        "⚔️",
			Priority:    3,
			Confidence:  0.8,
		}
	}

	// Bonus 觸發率很低，建議專注累積勞動值
	if bonusRate < 0.01 && b.TotalKills >= 30 {
		return &Recommendation{
			Type:        RecommendBonus,
			Title:       "累積勞動值",
			Description: "多打小目標快速累積勞動值，觸發 Bonus 獎勵更豐厚",
			Icon:        "🌾",
			Priority:    3,
			Confidence:  0.7,
		}
	}

	return nil
}

// analyzeSpecialOpportunity 分析特殊機會
func (e *Engine) analyzeSpecialOpportunity(b PlayerBehavior, rtp float64) *Recommendation {
	// 從未中過 Jackpot 且金幣充足，建議高難度房間
	if b.JackpotWins == 0 && b.CurrentCoins >= 50000 && b.CurrentBetLv >= 5 {
		return &Recommendation{
			Type:        RecommendJackpot,
			Title:       "挑戰 Jackpot",
			Description: "你還沒中過 Jackpot！高難度房間 Jackpot 貢獻更快，試試看",
			Icon:        "🎰",
			Priority:    3,
			Confidence:  0.65,
		}
	}

	// VIP 等級高但投注等級低，建議提升
	if b.VIPLevel >= 3 && b.CurrentBetLv <= 3 && b.CurrentCoins >= 100000 {
		return &Recommendation{
			Type:        RecommendBetUp,
			Title:       "VIP 玩家應有的體驗",
			Description: fmt.Sprintf("你是 VIP%d，試試高投注等級，獎勵倍率更高", b.VIPLevel),
			Icon:        "👑",
			Priority:    2,
			TargetBetLv: 5,
			Confidence:  0.7,
		}
	}

	// 登入連續天數長，給予鼓勵
	if b.LoginStreak >= 7 && rtp >= 0.9 {
		return &Recommendation{
			Type:        RecommendBetStay,
			Title:       "連續登入獎勵加成中",
			Description: fmt.Sprintf("連續登入 %d 天！保持節奏，今天的 RTP 表現不錯", b.LoginStreak),
			Icon:        "🔥",
			Priority:    3,
			Confidence:  0.6,
		}
	}

	return nil
}

// betLevelCost 取得投注等級的每次費用
func betLevelCost(lv int) int {
	costs := []int{0, 1, 2, 3, 5, 10, 20, 30, 50, 80, 100}
	if lv < 1 || lv > 10 {
		return 10
	}
	return costs[lv]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
