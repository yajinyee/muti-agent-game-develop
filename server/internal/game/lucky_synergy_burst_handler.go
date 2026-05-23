// lucky_synergy_burst_handler.go — 幸運共鳴爆發魚系統（DAY-239）
// 業界原創「多效疊加共鳴爆發」機制
//
// 設計：擊破 T197 後，偵測場上當前同時啟動的幸運效果數量：
//   - ≥2 個效果同時啟動 → 「共鳴爆發」：所有現有效果倍率加成額外 ×1.5（疊加乘法），持續 6 秒
//   - 1 個效果啟動 → 「小型共鳴」：該效果倍率 ×1.3，持續 4 秒
//   - 0 個效果 → 「基礎爆發」：全場 HP -30%，個人 ×1.8 倍率加成 5 秒
//   - 個人冷卻 25 秒；全服冷卻 40 秒
//
// 設計差異：
//   - 與幸運共鳴魚（DAY-222，全服合力射擊累積能量）不同，共鳴爆發魚是「偵測現有效果疊加」，
//     讓玩家有「要先觸發多個幸運魚，再打共鳴爆發魚」的策略深度
//   - 「效果越多，爆發越強」讓玩家有「組合技」的成就感
//   - 「基礎爆發」確保即使沒有其他效果也有獎勵，降低挫敗感
//   - 全服廣播讓所有玩家都看到共鳴爆發的效果數量，製造「哇，同時有這麼多效果」的驚嘆感
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	LuckySynergyBurstPersonalCD    = 25 * time.Second // 個人冷卻
	LuckySynergyBurstGlobalCD      = 40 * time.Second // 全服冷卻
	LuckySynergyBurstFullDuration  = 6 * time.Second  // 共鳴爆發持續時間（≥2 效果）
	LuckySynergyBurstSmallDuration = 4 * time.Second  // 小型共鳴持續時間（1 效果）
	LuckySynergyBurstBaseDuration  = 5 * time.Second  // 基礎爆發持續時間（0 效果）
	LuckySynergyBurstFullMult      = 1.5              // 共鳴爆發額外倍率乘數（疊加到現有效果）
	LuckySynergyBurstSmallMult     = 1.3              // 小型共鳴額外倍率乘數
	LuckySynergyBurstBaseMult      = 1.8              // 基礎爆發個人倍率加成
	LuckySynergyBurstBaseHPDmg     = 0.30             // 基礎爆發全場 HP 削減比例
)

// luckySynergyBurstManager 幸運共鳴爆發魚管理器
type luckySynergyBurstManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 共鳴爆發狀態（個人）
	activePlayers map[string]synergyBurstSession
}

// synergyBurstSession 個人共鳴爆發 session
type synergyBurstSession struct {
	burstType   string    // "full" / "small" / "base"
	activeUntil time.Time
	extraMult   float64   // 額外倍率乘數
	effectCount int       // 觸發時的效果數量
}

func newLuckySynergyBurstManager() *luckySynergyBurstManager {
	return &luckySynergyBurstManager{
		personalCooldowns: make(map[string]time.Time),
		activePlayers:     make(map[string]synergyBurstSession),
	}
}

// isLuckySynergyBurstFish 判斷是否為幸運共鳴爆發魚
func isLuckySynergyBurstFish(defID string) bool {
	return defID == "T197"
}

// getLuckySynergyBurstMult 取得共鳴爆發額外倍率（供 handleKill 使用）
func (g *Game) getLuckySynergyBurstMult(playerID string) float64 {
	mgr := g.LuckySynergyBurst
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	sess, ok := mgr.activePlayers[playerID]
	if !ok {
		return 1.0
	}
	if time.Now().After(sess.activeUntil) {
		delete(mgr.activePlayers, playerID)
		return 1.0
	}
	return sess.extraMult
}

// countActiveEffects 計算當前場上同時啟動的幸運效果數量
func (g *Game) countActiveEffects() (count int, names []string) {
	if g.isMirrorWorldActive() {
		count++
		names = append(names, "鏡面世界")
	}
	if g.isFreezeWorldActive() {
		count++
		names = append(names, "冰凍世界")
	}
	if g.isGravityFlipActive() {
		count++
		names = append(names, "重力反轉")
	}
	if g.isMagnetActive() {
		count++
		names = append(names, "磁力場")
	}
	if g.isLuckyMirrorTimeActive() {
		count++
		names = append(names, "時間倒流")
	}
	if g.isVortexActive() {
		count++
		names = append(names, "漩渦")
	}
	if g.isTimeFreezeActive() {
		count++
		names = append(names, "時間凍結")
	}
	return
}

// tryLuckySynergyBurstFish 擊破 T197 後觸發共鳴爆發（供 handleKill 使用）
func (g *Game) tryLuckySynergyBurstFish(p *player.Player) {
	mgr := g.LuckySynergyBurst
	mgr.mu.Lock()

	// 個人冷卻檢查
	if cd, ok := mgr.personalCooldowns[p.ID]; ok && time.Now().Before(cd) {
		mgr.mu.Unlock()
		return
	}
	// 全服冷卻檢查
	if time.Now().Before(mgr.globalCooldownUntil) {
		mgr.mu.Unlock()
		return
	}

	// 設定冷卻
	mgr.personalCooldowns[p.ID] = time.Now().Add(LuckySynergyBurstPersonalCD)
	mgr.globalCooldownUntil = time.Now().Add(LuckySynergyBurstGlobalCD)
	mgr.mu.Unlock()

	// 計算當前啟動的效果數量
	effectCount, effectNames := g.countActiveEffects()

	log.Printf("[LuckySynergyBurst] player=%s triggered synergy burst, active effects=%d %v",
		p.ID, effectCount, effectNames)

	switch {
	case effectCount >= 2:
		g.doSynergyFullBurst(p, effectCount, effectNames)
	case effectCount == 1:
		g.doSynergySmallBurst(p, effectNames)
	default:
		g.doSynergyBaseBurst(p)
	}
}

// doSynergyFullBurst 共鳴爆發（≥2 個效果）
func (g *Game) doSynergyFullBurst(p *player.Player, effectCount int, effectNames []string) {
	mgr := g.LuckySynergyBurst
	mgr.mu.Lock()
	mgr.activePlayers[p.ID] = synergyBurstSession{
		burstType:   "full",
		activeUntil: time.Now().Add(LuckySynergyBurstFullDuration),
		extraMult:   LuckySynergyBurstFullMult,
		effectCount: effectCount,
	}
	mgr.mu.Unlock()

	namesStr := ""
	for i, n := range effectNames {
		if i > 0 {
			namesStr += "+"
		}
		namesStr += n
	}

	// 全服廣播
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckySynergyBurst,
		Payload: ws.LuckySynergyBurstPayload{
			Event:       "synergy_full",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			EffectCount: effectCount,
			EffectNames: effectNames,
			ExtraMult:   LuckySynergyBurstFullMult,
			DurationSec: int(LuckySynergyBurstFullDuration.Seconds()),
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckySynergyBurst, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("✨ %s 觸發共鳴爆發！%d 個效果疊加（%s），所有倍率 ×%.1f！",
			p.DisplayName, effectCount, namesStr, LuckySynergyBurstFullMult),
		"color": "#FF6B9D",
	})
	g.broadcastAnnouncement(ann)

	// 6 秒後清除
	go func() {
		time.Sleep(LuckySynergyBurstFullDuration)
		mgr.mu.Lock()
		if sess, ok := mgr.activePlayers[p.ID]; ok && sess.burstType == "full" {
			delete(mgr.activePlayers, p.ID)
		}
		mgr.mu.Unlock()
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckySynergyBurst,
			Payload: ws.LuckySynergyBurstPayload{
				Event:    "synergy_end",
				PlayerID: p.ID,
			},
		})
	}()
}

// doSynergySmallBurst 小型共鳴（1 個效果）
func (g *Game) doSynergySmallBurst(p *player.Player, effectNames []string) {
	mgr := g.LuckySynergyBurst
	mgr.mu.Lock()
	mgr.activePlayers[p.ID] = synergyBurstSession{
		burstType:   "small",
		activeUntil: time.Now().Add(LuckySynergyBurstSmallDuration),
		extraMult:   LuckySynergyBurstSmallMult,
		effectCount: 1,
	}
	mgr.mu.Unlock()

	effectName := ""
	if len(effectNames) > 0 {
		effectName = effectNames[0]
	}

	// 全服廣播
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckySynergyBurst,
		Payload: ws.LuckySynergyBurstPayload{
			Event:       "synergy_small",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			EffectCount: 1,
			EffectNames: effectNames,
			ExtraMult:   LuckySynergyBurstSmallMult,
			DurationSec: int(LuckySynergyBurstSmallDuration.Seconds()),
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckySynergyBurst, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("✨ %s 觸發小型共鳴！%s 效果倍率 ×%.1f！",
			p.DisplayName, effectName, LuckySynergyBurstSmallMult),
		"color": "#C39BD3",
	})
	g.broadcastAnnouncement(ann)

	// 4 秒後清除
	go func() {
		time.Sleep(LuckySynergyBurstSmallDuration)
		mgr.mu.Lock()
		if sess, ok := mgr.activePlayers[p.ID]; ok && sess.burstType == "small" {
			delete(mgr.activePlayers, p.ID)
		}
		mgr.mu.Unlock()
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckySynergyBurst,
			Payload: ws.LuckySynergyBurstPayload{
				Event:    "synergy_end",
				PlayerID: p.ID,
			},
		})
	}()
}

// doSynergyBaseBurst 基礎爆發（0 個效果）
func (g *Game) doSynergyBaseBurst(p *player.Player) {
	mgr := g.LuckySynergyBurst
	mgr.mu.Lock()
	mgr.activePlayers[p.ID] = synergyBurstSession{
		burstType:   "base",
		activeUntil: time.Now().Add(LuckySynergyBurstBaseDuration),
		extraMult:   LuckySynergyBurstBaseMult,
		effectCount: 0,
	}
	mgr.mu.Unlock()

	// 全場 HP -30%
	g.mu.Lock()
	damagedCount := 0
	for _, t := range g.Targets {
		if t.HP <= 0 {
			continue
		}
		damage := int(float64(t.HP) * LuckySynergyBurstBaseHPDmg)
		if damage < 1 {
			damage = 1
		}
		t.HP -= damage
		if t.HP < 1 {
			t.HP = 1
		}
		damagedCount++
	}
	g.mu.Unlock()

	log.Printf("[LuckySynergyBurst] base burst: HP -30%% on %d targets, player=%s", damagedCount, p.ID)

	// 全服廣播
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckySynergyBurst,
		Payload: ws.LuckySynergyBurstPayload{
			Event:        "synergy_base",
			PlayerID:     p.ID,
			PlayerName:   p.DisplayName,
			EffectCount:  0,
			ExtraMult:    LuckySynergyBurstBaseMult,
			DurationSec:  int(LuckySynergyBurstBaseDuration.Seconds()),
			DamagedCount: damagedCount,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckySynergyBurst, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("✨ %s 觸發基礎爆發！全場 HP -30%%，個人 ×%.1f 倍率加成！",
			p.DisplayName, LuckySynergyBurstBaseMult),
		"color": "#F39C12",
	})
	g.broadcastAnnouncement(ann)

	// 5 秒後清除
	go func() {
		time.Sleep(LuckySynergyBurstBaseDuration)
		mgr.mu.Lock()
		if sess, ok := mgr.activePlayers[p.ID]; ok && sess.burstType == "base" {
			delete(mgr.activePlayers, p.ID)
		}
		mgr.mu.Unlock()
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckySynergyBurst,
			Payload: ws.LuckySynergyBurstPayload{
				Event:    "synergy_end",
				PlayerID: p.ID,
			},
		})
	}()
}
