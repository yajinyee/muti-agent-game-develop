// vortex_fish_handler.go — 漩渦魚群吸引系統（DAY-169）
// 業界依據：Ocean King（Google Play 2026）「Vortex Fish — catching a Vortex Fish will suck
// all fish of the same species in the area into a whirlpool, capturing them all at once.」
// 設計：擊破 T127 漩渦魚後，場上所有相同「目標類型」的目標被吸入漩渦，全部擊破
// 與黑洞（吸引所有目標）不同：漩渦魚是「同類型吸引」，讓玩家有策略性選擇
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"digital-twin/server/internal/data"
	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/game/target"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// vortexFishManager 漩渦魚群管理器
type vortexFishManager struct {
	mu        sync.Mutex
	cooldown  time.Time // 全服冷卻（防止頻繁觸發）
	CooldownSec int
}

func newVortexFishManager() *vortexFishManager {
	return &vortexFishManager{
		CooldownSec: 20, // 20 秒全服冷卻
	}
}

// isVortexFish 判斷是否為漩渦魚
func isVortexFish(defID string) bool {
	return defID == "T127"
}

// getVortexTargetGroup 取得漩渦魚的目標群組（同類型目標）
// 回傳：目標群組名稱（用於廣播顯示）和篩選函數
func getVortexTargetGroup(defID string) (groupName string, filter func(*target.Target) bool) {
	// T127 漩渦魚：吸引所有基礎目標（T001-T006）
	// 設計理由：基礎目標是最常見的，讓玩家感受到「一網打盡」的爽感
	// 同時不影響特殊目標（保持遊戲平衡）
	return "基礎目標群", func(t *target.Target) bool {
		return t.Def != nil && t.Def.Type == data.TargetTypeBasic
	}
}

// tryVortexFishSuck 漩渦魚擊破後觸發漩渦吸引
func (g *Game) tryVortexFishSuck(p *player.Player, vortexInstanceID string, vx, vy float64) {
	g.VortexFish.mu.Lock()
	if time.Now().Before(g.VortexFish.cooldown) {
		g.VortexFish.mu.Unlock()
		return
	}
	g.VortexFish.cooldown = time.Now().Add(time.Duration(g.VortexFish.CooldownSec) * time.Second)
	g.VortexFish.mu.Unlock()

	// 取得目標群組
	groupName, filter := getVortexTargetGroup("T127")

	// 收集場上符合條件的目標
	g.mu.RLock()
	var targets []*target.Target
	for _, t := range g.Targets {
		if t.InstanceID == vortexInstanceID {
			continue // 跳過漩渦魚本身（已被擊破）
		}
		if filter(t) {
			targets = append(targets, t)
		}
	}
	g.mu.RUnlock()

	if len(targets) == 0 {
		return
	}

	// 廣播漩渦開始（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgVortexFish,
		Payload: ws.VortexFishPayload{
			Phase:       "vortex_start",
			TriggerID:   p.ID,
			TriggerName: p.DisplayName,
			VortexX:     vx,
			VortexY:     vy,
			GroupName:   groupName,
			TargetCount: len(targets),
		},
	})

	// 逐步吸入並擊破目標（每 150ms 一個，製造「漩渦吸入」的視覺效果）
	go func() {
		totalReward := 0
		killedCount := 0
		var killedEntries []ws.VortexKillEntry

		for i, t := range targets {
			time.Sleep(150 * time.Millisecond)

			// 從場上移除目標
			g.mu.Lock()
			if _, exists := g.Targets[t.InstanceID]; !exists {
				g.mu.Unlock()
				continue // 目標已被其他玩家擊破
			}
			delete(g.Targets, t.InstanceID)
			g.mu.Unlock()

			// 計算獎勵（漩渦吸引的獎勵 = 目標倍率 × betLevel × 0.55，比直接擊破低）
			betCost := 0
			if bd := data.GetBetDef(p.BetLevel); bd != nil {
				betCost = bd.BetCost
			}
			mult := t.Def.MultiplierMin
			if t.Def.MultiplierMax > mult {
				// 取中間值
				mult = (t.Def.MultiplierMin + t.Def.MultiplierMax) / 2
			}
			reward := int(mult * float64(betCost) * 0.55)
			if reward < 1 {
				reward = 1
			}

			totalReward += reward
			killedCount++

			killedEntries = append(killedEntries, ws.VortexKillEntry{
				InstanceID: t.InstanceID,
				DefID:      t.DefID,
				Multiplier: mult,
				Reward:     reward,
				X:          t.X,
				Y:          t.Y,
			})

			// 廣播每個目標被吸入（讓 Client 做吸入動畫）
			g.Hub.Broadcast(&ws.Message{
				Type: ws.MsgVortexFish,
				Payload: ws.VortexFishPayload{
					Phase:       "vortex_suck",
					TriggerID:   p.ID,
					TriggerName: p.DisplayName,
					VortexX:     vx,
					VortexY:     vy,
					GroupName:   groupName,
					TargetCount: len(targets),
					SuckIndex:   i,
					SuckEntry: &ws.VortexKillEntry{
						InstanceID: t.InstanceID,
						DefID:      t.DefID,
						Multiplier: mult,
						Reward:     reward,
						X:          t.X,
						Y:          t.Y,
					},
				},
			})
		}

		if killedCount == 0 {
			return
		}

		// 發放總獎勵給觸發玩家
		p.AddReward(totalReward)
		g.sendPlayerUpdate(p)

		// 廣播漩渦結束（全服）
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgVortexFish,
			Payload: ws.VortexFishPayload{
				Phase:        "vortex_end",
				TriggerID:    p.ID,
				TriggerName:  p.DisplayName,
				VortexX:      vx,
				VortexY:      vy,
				GroupName:    groupName,
				TargetCount:  len(targets),
				KilledCount:  killedCount,
				TotalReward:  totalReward,
				KilledEntries: killedEntries,
			},
		})

		// 全服公告：≥5 個目標被吸入時廣播
		if killedCount >= 5 {
			g.announceVortexFish(p.DisplayName, killedCount, totalReward)
		}

		log.Printf("[VortexFish] player=%s triggered vortex: killed=%d totalReward=%d",
			p.ID, killedCount, totalReward)
	}()
}

// announceVortexFish 全服公告漩渦魚大豐收
func (g *Game) announceVortexFish(playerName string, killedCount, totalReward int) {
	ann := g.Announce.Create(announce.EventVortexFish, playerName, totalReward, map[string]string{
		"killed": fmt.Sprintf("%d", killedCount),
		"reward": fmt.Sprintf("%d", totalReward),
	})
	g.broadcastAnnouncement(ann)
}
