// chainlong_king_handler.go — 長龍王雙環輪盤系統（DAY-194）
// 業界依據：Royal Fishing JILI「ChainLong King — dual-ring roulette activates when captured.
// You control when the pointer stops, multiplying inner and outer ring values together.
// Maximum combination delivers 350X, whilst the ChainLong King itself can award up to 1000X mega wins.」
// 設計：擊破 T152 後觸發「雙環輪盤」互動（個人）：
//   - 內環：5x/10x/20x/50x（玩家點擊停止）
//   - 外環：1x/2x/3x/5x/7x（玩家點擊停止）
//   - 最終獎勵 = 內環 × 外環 × betLevel（最高 350x）
//   - 特殊：1% 機率觸發「千倍大獎」（1000x），跳過輪盤直接給獎
// 設計差異：與 DAY-113 雙層倍率輪盤（全服共享，自動停止）不同，
//   長龍王是「個人互動輪盤」（玩家主動點擊停止），製造「我控制命運」的掌控感；
//   與 DAY-139 雙環輪盤（全服事件）不同，長龍王是「擊破特定目標觸發」，更有目標感；
//   千倍大獎（1%）讓玩家每次擊破長龍王都有「說不定這次就是千倍」的期待感
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// 長龍王雙環輪盤常數
const (
	ChainLongKingCooldownSec = 30 // 個人冷卻 30 秒
	ChainLongKingTimeoutSec  = 15 // 輪盤互動超時 15 秒（超時自動停止）
	ChainLongKingMegaChance  = 1  // 千倍大獎機率（1%）
)

// 內環倍率（5x/10x/20x/50x，加權）
var chainLongInnerRing = []struct {
	Mult   int
	Weight int
}{
	{5, 45},
	{10, 30},
	{20, 18},
	{50, 7},
}

// 外環乘數（1x/2x/3x/5x/7x，加權）
var chainLongOuterRing = []struct {
	Mult   int
	Weight int
}{
	{1, 40},
	{2, 28},
	{3, 18},
	{5, 10},
	{7, 4},
}

// chainLongKingSession 個人輪盤 session
type chainLongKingSession struct {
	PlayerID    string
	InstanceID  string
	Phase       string    // "inner_spin" → "outer_spin" → "result"
	InnerResult int       // 內環停止結果（0 = 未停止）
	OuterResult int       // 外環停止結果（0 = 未停止）
	IsMega      bool      // 是否千倍大獎
	StartAt     time.Time
	CooldownEnd time.Time
}

// chainLongKingManager 長龍王雙環輪盤管理器
type chainLongKingManager struct {
	mu       sync.Mutex
	sessions map[string]*chainLongKingSession // playerID → session
}

func newChainLongKingManager() *chainLongKingManager {
	return &chainLongKingManager{
		sessions: make(map[string]*chainLongKingSession),
	}
}

// isChainLongKingFish 判斷是否為長龍王（T152，DAY-194）
func isChainLongKingFish(defID string) bool {
	return defID == "T152"
}

// tryChainLongKingRoulette 擊破 T152 後觸發雙環輪盤
func (g *Game) tryChainLongKingRoulette(p *player.Player, instanceID string) {
	mgr := g.ChainLongKing
	mgr.mu.Lock()

	// 檢查個人冷卻
	if sess, ok := mgr.sessions[p.ID]; ok {
		if time.Now().Before(sess.CooldownEnd) {
			mgr.mu.Unlock()
			return
		}
	}

	// 1% 機率千倍大獎（跳過輪盤）
	isMega := rand.Intn(100) < ChainLongKingMegaChance

	sess := &chainLongKingSession{
		PlayerID:   p.ID,
		InstanceID: instanceID,
		Phase:      "inner_spin",
		IsMega:     isMega,
		StartAt:    time.Now(),
	}
	mgr.sessions[p.ID] = sess
	mgr.mu.Unlock()

	if isMega {
		// 千倍大獎：直接給獎，不需要輪盤互動
		g.resolveChainLongKingMega(p, sess)
		return
	}

	// 廣播輪盤開始（個人）
	if err := g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgChainLongKing,
		Payload: ws.ChainLongKingPayload{
			Phase:      "roulette_start",
			InstanceID: instanceID,
			InnerRing:  buildChainLongRingDef(chainLongInnerRing),
			OuterRing:  buildChainLongRingDef(chainLongOuterRing),
		},
	}); err != nil {
		log.Printf("[ChainLongKing] send roulette_start error: %v", err)
	}

	log.Printf("[ChainLongKing] player=%s triggered dual-ring roulette (instanceID=%s)", p.ID, instanceID)

	// 啟動超時 goroutine（15 秒後自動停止）
	go func() {
		time.Sleep(ChainLongKingTimeoutSec * time.Second)
		mgr.mu.Lock()
		s, ok := mgr.sessions[p.ID]
		if !ok || s.InstanceID != instanceID {
			mgr.mu.Unlock()
			return
		}
		// 超時自動停止：隨機選結果
		if s.Phase == "inner_spin" {
			s.InnerResult = pickChainLongRing(chainLongInnerRing)
			s.Phase = "outer_spin"
			mgr.mu.Unlock()
			g.Hub.Send(p.ID, &ws.Message{
				Type: ws.MsgChainLongKing,
				Payload: ws.ChainLongKingPayload{
					Phase:       "inner_stop",
					InstanceID:  instanceID,
					InnerResult: s.InnerResult,
					IsTimeout:   true,
				},
			})
			time.Sleep(1 * time.Second)
			mgr.mu.Lock()
		}
		if s.Phase == "outer_spin" {
			s.OuterResult = pickChainLongRing(chainLongOuterRing)
			s.Phase = "result"
			mgr.mu.Unlock()
			g.Hub.Send(p.ID, &ws.Message{
				Type: ws.MsgChainLongKing,
				Payload: ws.ChainLongKingPayload{
					Phase:       "outer_stop",
					InstanceID:  instanceID,
					OuterResult: s.OuterResult,
					IsTimeout:   true,
				},
			})
			time.Sleep(500 * time.Millisecond)
			g.resolveChainLongKingResult(p, s)
			return
		}
		mgr.mu.Unlock()
	}()
}

// handleChainLongKingStop 處理玩家點擊停止輪盤（由 HandleMessage 呼叫）
func (g *Game) handleChainLongKingStop(p *player.Player, payload ws.ChainLongKingStopPayload) {
	mgr := g.ChainLongKing
	mgr.mu.Lock()
	sess, ok := mgr.sessions[p.ID]
	if !ok || sess.Phase == "result" {
		mgr.mu.Unlock()
		return
	}

	switch sess.Phase {
	case "inner_spin":
		// 停止內環：隨機選結果（玩家點擊時機決定「感覺」，實際是隨機）
		sess.InnerResult = pickChainLongRing(chainLongInnerRing)
		sess.Phase = "outer_spin"
		mgr.mu.Unlock()

		if err := g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgChainLongKing,
			Payload: ws.ChainLongKingPayload{
				Phase:       "inner_stop",
				InstanceID:  sess.InstanceID,
				InnerResult: sess.InnerResult,
				IsTimeout:   false,
			},
		}); err != nil {
			log.Printf("[ChainLongKing] send inner_stop error: %v", err)
		}

	case "outer_spin":
		// 停止外環：隨機選結果
		sess.OuterResult = pickChainLongRing(chainLongOuterRing)
		sess.Phase = "result"
		mgr.mu.Unlock()

		if err := g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgChainLongKing,
			Payload: ws.ChainLongKingPayload{
				Phase:       "outer_stop",
				InstanceID:  sess.InstanceID,
				OuterResult: sess.OuterResult,
				IsTimeout:   false,
			},
		}); err != nil {
			log.Printf("[ChainLongKing] send outer_stop error: %v", err)
		}

		// 短暫延遲後結算
		go func() {
			time.Sleep(600 * time.Millisecond)
			g.resolveChainLongKingResult(p, sess)
		}()

	default:
		mgr.mu.Unlock()
	}
}

// resolveChainLongKingResult 結算雙環輪盤獎勵
func (g *Game) resolveChainLongKingResult(p *player.Player, sess *chainLongKingSession) {
	// 計算獎勵：內環 × 外環 × betLevel
	totalMult := sess.InnerResult * sess.OuterResult
	reward := totalMult * p.BetLevel

	// 給予獎勵
	g.mu.Lock()
	p.Coins += reward
	g.mu.Unlock()

	// 設定冷卻
	mgr := g.ChainLongKing
	mgr.mu.Lock()
	if s, ok := mgr.sessions[p.ID]; ok && s.InstanceID == sess.InstanceID {
		s.CooldownEnd = time.Now().Add(ChainLongKingCooldownSec * time.Second)
	}
	mgr.mu.Unlock()

	// 廣播結算結果（個人）
	isBigWin := totalMult >= 100
	if err := g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgChainLongKing,
		Payload: ws.ChainLongKingPayload{
			Phase:       "result",
			InstanceID:  sess.InstanceID,
			InnerResult: sess.InnerResult,
			OuterResult: sess.OuterResult,
			TotalMult:   totalMult,
			Reward:      reward,
			IsBigWin:    isBigWin,
		},
	}); err != nil {
		log.Printf("[ChainLongKing] send result error: %v", err)
	}

	// 全服廣播（≥100x 才廣播）
	if isBigWin {
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgChainLongKing,
			Payload: ws.ChainLongKingPayload{
				Phase:       "broadcast",
				PlayerName:  p.DisplayName,
				TotalMult:   totalMult,
				Reward:      reward,
			},
		})
		// 全服公告
		ann := g.Announce.Create(announce.EventMegaWin, p.DisplayName, reward, map[string]string{
			"message": formatChainLongKingAnnounce(p.DisplayName, totalMult),
			"color":   chainLongKingColor(totalMult),
		})
		g.broadcastAnnouncement(ann)
	}

	log.Printf("[ChainLongKing] player=%s result: inner=%dx outer=%dx total=%dx reward=%d",
		p.ID, sess.InnerResult, sess.OuterResult, totalMult, reward)
}

// resolveChainLongKingMega 千倍大獎直接結算
func (g *Game) resolveChainLongKingMega(p *player.Player, sess *chainLongKingSession) {
	reward := 1000 * p.BetLevel

	// 給予獎勵
	g.mu.Lock()
	p.Coins += reward
	g.mu.Unlock()

	// 設定冷卻
	mgr := g.ChainLongKing
	mgr.mu.Lock()
	if s, ok := mgr.sessions[p.ID]; ok && s.InstanceID == sess.InstanceID {
		s.CooldownEnd = time.Now().Add(ChainLongKingCooldownSec * time.Second)
	}
	mgr.mu.Unlock()

	// 廣播千倍大獎（個人）
	if err := g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgChainLongKing,
		Payload: ws.ChainLongKingPayload{
			Phase:      "mega_win",
			InstanceID: sess.InstanceID,
			TotalMult:  1000,
			Reward:     reward,
			IsMega:     true,
		},
	}); err != nil {
		log.Printf("[ChainLongKing] send mega_win error: %v", err)
	}

	// 全服廣播千倍大獎
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgChainLongKing,
		Payload: ws.ChainLongKingPayload{
			Phase:      "mega_broadcast",
			PlayerName: p.DisplayName,
			TotalMult:  1000,
			Reward:     reward,
			IsMega:     true,
		},
	})

	// 全服公告（金色，最高優先級）
	ann := g.Announce.Create(announce.EventMegaWin, p.DisplayName, reward, map[string]string{
		"message": "🐉 " + p.DisplayName + " 觸發長龍王千倍大獎！獲得 1000x！",
		"color":   "#FFD700",
	})
	g.broadcastAnnouncement(ann)

	log.Printf("[ChainLongKing] MEGA WIN! player=%s reward=%d (1000x)", p.ID, reward)
}

// pickChainLongRing 加權隨機選取輪盤結果
func pickChainLongRing(rings []struct {
	Mult   int
	Weight int
}) int {
	total := 0
	for _, r := range rings {
		total += r.Weight
	}
	n := rand.Intn(total)
	for _, r := range rings {
		n -= r.Weight
		if n < 0 {
			return r.Mult
		}
	}
	return rings[0].Mult
}

// buildChainLongRingDef 建立輪盤定義（供 Client 顯示）
func buildChainLongRingDef(rings []struct {
	Mult   int
	Weight int
}) []int {
	result := make([]int, len(rings))
	for i, r := range rings {
		result[i] = r.Mult
	}
	return result
}

// formatChainLongKingAnnounce 格式化全服公告文字
func formatChainLongKingAnnounce(playerName string, totalMult int) string {
	switch {
	case totalMult >= 350:
		return fmt.Sprintf("🐉 %s 觸發長龍王最高倍率 %dx！", playerName, totalMult)
	case totalMult >= 200:
		return fmt.Sprintf("🐉 %s 長龍王雙環輪盤獲得 %dx 大獎！", playerName, totalMult)
	default:
		return fmt.Sprintf("🐉 %s 長龍王輪盤獲得 %dx！", playerName, totalMult)
	}
}

// chainLongKingColor 依倍率決定公告顏色
func chainLongKingColor(totalMult int) string {
	switch {
	case totalMult >= 350:
		return "#FF4500" // 橙紅色（最高倍率）
	case totalMult >= 150:
		return "#FFD700" // 金色
	default:
		return "#FFA500" // 橙色
	}
}
