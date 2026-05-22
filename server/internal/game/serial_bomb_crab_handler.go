// serial_bomb_crab_handler.go — 連環炸彈蟹系統（DAY-201）
// 業界依據：Royal Fishing JILI「Serial Bomb Crab — orange crab with panda face and skull bomb designs.
// Triggers large-scale multiple explosions across screen, capturing fish within each explosion range.
// Worth 70x, this explosive crustacean triggers multiple large-scale detonations.
// Each bomb creates expanding capture zones for massive multi-target eliminations.」
//
// 設計：擊破 T159 後觸發「連環爆炸」：
//   1. 3-5 顆炸彈依序爆炸（每顆間隔 600ms）
//   2. 每顆炸彈：250px 半徑，75% 擊破機率，0.65x 倍率
//   3. 炸彈位置：隨機分散在場上（製造「全場覆蓋」感）
//   4. 全服廣播每顆炸彈的爆炸位置和結果
//
// 設計差異：
//   - 與連鎖爆炸魚（BFS 傳播，從擊破點擴散）不同，連環炸彈蟹是「預設位置的多顆炸彈」，
//     讓玩家看到「炸彈在場上各處爆炸」的壯觀感
//   - 與鑽頭龍蝦（穿透移動）不同，連環炸彈蟹是「靜態爆炸」，但多顆同時覆蓋全場
//   - 炸彈位置隨機分散，讓每次觸發都有不同的「覆蓋模式」，增加驚喜感
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

// 連環炸彈蟹常數
const (
	SerialBombCrabCooldownSec = 35    // 全服冷卻 35 秒
	SerialBombCrabMinBombs    = 3     // 最少炸彈數
	SerialBombCrabMaxBombs    = 5     // 最多炸彈數
	SerialBombCrabRadius      = 250.0 // 爆炸半徑 250px
	SerialBombCrabKillChance  = 0.75  // 擊破機率 75%
	SerialBombCrabMult        = 0.65  // 獎勵倍率 0.65x
	SerialBombCrabInterval    = 600   // 炸彈間隔 600ms
)

// serialBombCrabManager 連環炸彈蟹管理器（全服共享）
type serialBombCrabManager struct {
	mu          sync.Mutex
	isActive    bool
	cooldownEnd time.Time
}

func newSerialBombCrabManager() *serialBombCrabManager {
	return &serialBombCrabManager{}
}

// isSerialBombCrab 判斷是否為連環炸彈蟹（T159）
func isSerialBombCrab(defID string) bool {
	return defID == "T159"
}

// trySerialBombCrabExplosion 擊破 T159 後觸發連環爆炸
func (g *Game) trySerialBombCrabExplosion(p *player.Player, killX, killY float64) {
	mgr := g.SerialBombCrab
	mgr.mu.Lock()

	// 全服冷卻檢查
	if mgr.isActive || time.Now().Before(mgr.cooldownEnd) {
		mgr.mu.Unlock()
		return
	}
	mgr.isActive = true
	mgr.mu.Unlock()

	// 決定炸彈數量（3-5 顆）
	bombCount := SerialBombCrabMinBombs + rand.Intn(SerialBombCrabMaxBombs-SerialBombCrabMinBombs+1)

	log.Printf("[SerialBombCrab] player=%s triggered %d bombs", p.ID, bombCount)

	// 生成炸彈位置（隨機分散在場上）
	bombs := make([]serialBombPos, bombCount)
	for i := 0; i < bombCount; i++ {
		bombs[i] = serialBombPos{
			x: 150 + rand.Float64()*980, // 150-1130px
			y: 80 + rand.Float64()*560,  // 80-640px
		}
	}

	// 全服廣播：連環爆炸開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgSerialBombCrab,
		Payload: ws.SerialBombCrabPayload{
			Event:      "bomb_start",
			KillerName: p.DisplayName,
			BombCount:  bombCount,
			KillX:      killX,
			KillY:      killY,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventMegaWin, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("💣 %s 觸發連環炸彈蟹！%d 顆炸彈即將爆炸！", p.DisplayName, bombCount),
	})
	g.broadcastAnnouncement(ann)

	// 依序引爆每顆炸彈
	go g.runSerialBombCrabExplosions(p, bombs)
}

// runSerialBombCrabExplosions 依序引爆所有炸彈
type serialBombPos struct {
	x, y float64
}

func (g *Game) runSerialBombCrabExplosions(p *player.Player, bombs []serialBombPos) {
	totalKills := 0
	totalReward := 0

	for i, bomb := range bombs {
		time.Sleep(SerialBombCrabInterval * time.Millisecond)

		// 找出爆炸範圍內的目標
		g.mu.RLock()
		type hitTarget struct {
			id   string
			mult float64
		}
		var inRange []hitTarget
		for _, t := range g.Targets {
			if t.HP <= 0 || t.DefID == "B001" || isGhostFishClone(t.DefID) {
				continue
			}
			dx := t.X - bomb.x
			dy := t.Y - bomb.y
			dist := dx*dx + dy*dy
			if dist <= SerialBombCrabRadius*SerialBombCrabRadius {
				inRange = append(inRange, hitTarget{id: t.InstanceID, mult: t.Multiplier})
			}
		}
		g.mu.RUnlock()

		// 對範圍內目標執行爆炸
		bombKills := 0
		bombReward := 0
		for _, ht := range inRange {
			if rand.Float64() < SerialBombCrabKillChance {
				g.mu.Lock()
				t, ok := g.Targets[ht.id]
				if ok && t.HP > 0 {
					r := int(ht.mult * float64(p.BetLevel) * SerialBombCrabMult)
					if r < 1 {
						r = 1
					}
					delete(g.Targets, ht.id)
					bombKills++
					bombReward += r
					totalKills++
					totalReward += r
					// 給觸發者獎勵
					if pp, ok2 := g.Players[p.ID]; ok2 {
						pp.Coins += r
					}
				}
				g.mu.Unlock()
			}
		}

		// 廣播單顆炸彈結果
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgSerialBombCrab,
			Payload: ws.SerialBombCrabPayload{
				Event:      "bomb_explode",
				BombIndex:  i + 1,
				BombCount:  len(bombs),
				BombX:      bomb.x,
				BombY:      bomb.y,
				BombKills:  bombKills,
				BombReward: bombReward,
			},
		})
	}

	// 最終結算廣播
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgSerialBombCrab,
		Payload: ws.SerialBombCrabPayload{
			Event:       "bomb_result",
			KillerName:  p.DisplayName,
			TotalKills:  totalKills,
			TotalReward: totalReward,
		},
	})

	// 全服公告（≥4 個擊破才公告）
	if totalKills >= 4 {
		ann := g.Announce.Create(announce.EventMegaWin, p.DisplayName, totalReward, map[string]string{
			"message": fmt.Sprintf("💣💥 連環炸彈蟹！擊破 %d 個目標！獎勵 %d 金幣！", totalKills, totalReward),
		})
		g.broadcastAnnouncement(ann)
	}

	// 重置管理器
	mgr := g.SerialBombCrab
	mgr.mu.Lock()
	mgr.isActive = false
	mgr.cooldownEnd = time.Now().Add(SerialBombCrabCooldownSec * time.Second)
	mgr.mu.Unlock()

	log.Printf("[SerialBombCrab] complete: bombs=%d kills=%d reward=%d", len(bombs), totalKills, totalReward)
}
