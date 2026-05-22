// jackpot_dragon_handler.go — 獎池龍 Jackpot 抽獎系統（DAY-205）
// 業界依據：JILI Jackpot Fishing「special targets like the Jackpot Fish and Jackpot Dragon
// offering chances at substantial prizes. With the potential for high payouts up to 1000 times the bet.」
//
// 設計：擊破 T163 後觸發「獎池抽獎」（個人）：
//   - 加權隨機選擇 Jackpot 等級：Mini(70%) / Minor(20%) / Major(8%) / Grand(2%)
//   - 直接觸發對應等級的 Jackpot（ForceWin），立即給予獎勵
//   - 個人冷卻 60 秒（防止刷獎）
//   - Grand/Major 全服廣播 + 全服公告
//
// 設計差異（與普通 Jackpot 累積的區別）：
//   - 普通 Jackpot：每次射擊累積 0.5%，達到門檻後機率觸發（被動）
//   - 獎池龍：擊破後直接抽獎（主動），讓玩家有「我要去打那條龍」的目標感
//   - Grand 2% 機率讓玩家每次擊破都有「說不定這次就是 Grand」的期待感
//   - 個人冷卻確保不會被單一玩家壟斷
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/game/jackpot"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// 獎池龍常數（DAY-205）
const (
	JackpotDragonCooldownSec = 60 // 個人冷卻 60 秒
)

// 獎池龍抽獎等級權重
var jackpotDragonWeights = []struct {
	Level  jackpot.Level
	Weight int
}{
	{jackpot.LevelMini, 70},  // 70% Mini
	{jackpot.LevelMinor, 20}, // 20% Minor
	{jackpot.LevelMajor, 8},  // 8% Major
	{jackpot.LevelGrand, 2},  // 2% Grand
}

// jackpotDragonManager 獎池龍管理器（個人冷卻）
type jackpotDragonManager struct {
	mu        sync.Mutex
	cooldowns map[string]time.Time // playerID → 冷卻結束時間
}

func newJackpotDragonManager() *jackpotDragonManager {
	return &jackpotDragonManager{
		cooldowns: make(map[string]time.Time),
	}
}

// isJackpotDragon 判斷是否為獎池龍（T163，DAY-205）
func isJackpotDragon(defID string) bool {
	return defID == "T163"
}

// tryJackpotDragonDraw 擊破 T163 後觸發獎池抽獎
func (g *Game) tryJackpotDragonDraw(p *player.Player) {
	mgr := g.JackpotDragon
	mgr.mu.Lock()

	// 個人冷卻檢查
	if cd, ok := mgr.cooldowns[p.ID]; ok && time.Now().Before(cd) {
		mgr.mu.Unlock()
		return
	}
	mgr.cooldowns[p.ID] = time.Now().Add(JackpotDragonCooldownSec * time.Second)
	mgr.mu.Unlock()

	// 加權隨機選擇 Jackpot 等級
	level := pickJackpotDragonLevel()

	log.Printf("[JackpotDragon] player=%s drawing jackpot level=%s", p.ID, level)

	// 直接觸發對應等級的 Jackpot
	win := g.jackpotMgr.ForceWin(level, p.ID)
	if win == nil {
		log.Printf("[JackpotDragon] ForceWin returned nil for level=%s", level)
		return
	}

	// 使用現有的 handleJackpotWin 處理獎勵和廣播
	g.handleJackpotWin(p, win)

	// 額外廣播：獎池龍觸發通知（讓 Client 顯示龍的特效）
	levelName, levelColor, levelIcon := jackpot.GetLevelInfo(level)
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgJackpotDragon,
		Payload: ws.JackpotDragonPayload{
			Event:      "dragon_draw",
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			Level:      string(level),
			LevelName:  levelName,
			LevelColor: levelColor,
			LevelIcon:  levelIcon,
			Amount:     win.Amount,
			IsGrand:    level == jackpot.LevelGrand,
			IsMajor:    level == jackpot.LevelMajor,
		},
	})

	// Grand/Major 額外全服公告
	if level == jackpot.LevelGrand || level == jackpot.LevelMajor {
		ann := g.Announce.Create(announce.EventGrandJackpot, p.DisplayName, win.Amount, map[string]string{
			"message": formatJackpotDragonAnnounce(p.DisplayName, levelName, levelIcon, win.Amount),
			"color":   levelColor,
		})
		g.broadcastAnnouncement(ann)
	}

	log.Printf("[JackpotDragon] player=%s won %s jackpot: %d coins", p.ID, level, win.Amount)
}

// pickJackpotDragonLevel 加權隨機選擇 Jackpot 等級
func pickJackpotDragonLevel() jackpot.Level {
	total := 0
	for _, w := range jackpotDragonWeights {
		total += w.Weight
	}
	n := rand.Intn(total)
	for _, w := range jackpotDragonWeights {
		n -= w.Weight
		if n < 0 {
			return w.Level
		}
	}
	return jackpot.LevelMini
}

// formatJackpotDragonAnnounce 格式化全服公告文字
func formatJackpotDragonAnnounce(playerName, levelName, levelIcon string, amount int) string {
	return fmt.Sprintf("%s %s 擊破獎池龍觸發 %s JACKPOT！獲得 %d 金幣！",
		levelIcon, playerName, levelName, amount)
}
