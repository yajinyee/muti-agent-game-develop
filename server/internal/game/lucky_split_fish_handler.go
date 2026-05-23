// lucky_split_fish_handler.go — 幸運分裂魚系統（DAY-224）
// 業界原創「一魚分三」機制
//
// 設計：擊破 T182 後觸發「分裂爆炸」：
//   - T182 分裂成 3 個「分裂碎片」（HP = 原 HP × 30%，倍率 ×1.8）
//   - 分裂碎片在場上存活 8 秒，被擊破獲得 ×1.8 倍率加成（乘法）
//   - 8 秒後所有未被擊破的分裂碎片「二次爆炸」（65% 擊破機率，0.60x 倍率）
//   - 個人冷卻 18 秒；全服廣播分裂事件
//
// 設計差異：
//   - 與幸運鏡像魚（DAY-215，複製分身 HP 50%，×1.5）不同，
//     分裂魚是「一魚分三」，讓玩家有「打一個得三個機會」的爽感
//   - 分裂碎片 HP 只有 30%，更容易擊破，讓玩家有「快速連殺」的節奏感
//   - 分裂碎片倍率 ×1.8，比鏡像魚（×1.5）更高，讓玩家有「分裂碎片比本體更值錢」的感覺
//   - 「二次爆炸」讓玩家有「等待→爆發」的高潮設計
//   - 全服廣播讓所有玩家都看到分裂碎片，製造「全服競爭搶打碎片」的社交感
package game

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"

	"digital-twin/server/internal/data"
	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/game/target"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	LuckySplitPersonalCD    = 18 * time.Second // 個人冷卻
	LuckySplitDuration      = 8 * time.Second  // 分裂碎片存活時間
	LuckySplitCount         = 3                // 分裂碎片數量
	LuckySplitHPRatio       = 0.30             // 分裂碎片 HP 比例
	LuckySplitKillMult      = 1.8              // 分裂碎片擊破倍率加成
	LuckySplitBlastChance   = 0.65             // 二次爆炸擊破機率
	LuckySplitBlastMult     = 0.60             // 二次爆炸倍率
	LuckySplitSpreadRadius  = 120.0            // 分裂碎片散佈半徑
)

// splitFragmentEntry 分裂碎片記錄
type splitFragmentEntry struct {
	instanceID string
	playerID   string // 觸發者
}

// luckySplitFishManager 幸運分裂魚管理器
type luckySplitFishManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldown map[string]time.Time

	// 活躍的分裂碎片（fragmentInstanceID → entry）
	activeFragments map[string]*splitFragmentEntry
}

func newLuckySplitFishManager() *luckySplitFishManager {
	return &luckySplitFishManager{
		personalCooldown: make(map[string]time.Time),
		activeFragments:  make(map[string]*splitFragmentEntry),
	}
}

// isLuckySplitFish 判斷是否為幸運分裂魚
func isLuckySplitFish(defID string) bool {
	return defID == "T182"
}

// getLuckySplitFragmentMult 取得分裂碎片倍率加成（供 handleKill 使用）
func (g *Game) getLuckySplitFragmentMult(instanceID string) float64 {
	mgr := g.LuckySplitFish
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	if _, ok := mgr.activeFragments[instanceID]; ok {
		return LuckySplitKillMult
	}
	return 1.0
}

// removeLuckySplitFragment 分裂碎片被擊破後移除（供 handleKill 使用）
func (g *Game) removeLuckySplitFragment(instanceID string) {
	mgr := g.LuckySplitFish
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	delete(mgr.activeFragments, instanceID)
}

// tryLuckySplitFish 擊破 T182 後觸發分裂（供 handleKill 使用）
func (g *Game) tryLuckySplitFish(p *player.Player, originX, originY float64, originHP int) {
	mgr := g.LuckySplitFish
	mgr.mu.Lock()

	// 個人冷卻檢查
	if until, ok := mgr.personalCooldown[p.ID]; ok && time.Now().Before(until) {
		mgr.mu.Unlock()
		return
	}

	// 設定個人冷卻
	mgr.personalCooldown[p.ID] = time.Now().Add(LuckySplitPersonalCD)
	mgr.mu.Unlock()

	// 取得 T182 的定義（用於分裂碎片）
	def, ok := data.Targets["T182"]
	if !ok {
		return
	}

	// 計算分裂碎片 HP
	fragmentHP := int(float64(originHP) * LuckySplitHPRatio)
	if fragmentHP < 1 {
		fragmentHP = 1
	}

	// 生成 3 個分裂碎片
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	instanceID := fmt.Sprintf("split_%d", time.Now().UnixNano())

	type fragmentInfo struct {
		instanceID string
		x, y       float64
	}
	var fragments []fragmentInfo

	for i := 0; i < LuckySplitCount; i++ {
		angle := float64(i) * (360.0 / float64(LuckySplitCount)) * (math.Pi / 180.0)
		angle += rng.Float64() * 0.5 // 隨機偏移
		fx := originX + LuckySplitSpreadRadius*float64(i+1)/float64(LuckySplitCount)*0.8*math.Cos(angle)
		fy := originY + LuckySplitSpreadRadius*float64(i+1)/float64(LuckySplitCount)*0.5*math.Sin(angle)

		// 邊界限制
		if fx < 50 {
			fx = 50
		}
		if fx > 950 {
			fx = 950
		}
		if fy < 50 {
			fy = 50
		}
		if fy > 550 {
			fy = 550
		}

		fragID := fmt.Sprintf("%s_f%d", instanceID, i)

		// 在遊戲中生成分裂碎片目標
		g.mu.Lock()
		fragTarget := target.NewTarget(fragID, def, fx, fy)
		fragTarget.HP = fragmentHP
		fragTarget.Multiplier = def.MultiplierMin * LuckySplitKillMult
		g.Targets[fragID] = fragTarget
		g.mu.Unlock()

		// 記錄分裂碎片
		mgr.mu.Lock()
		mgr.activeFragments[fragID] = &splitFragmentEntry{
			instanceID: instanceID,
			playerID:   p.ID,
		}
		mgr.mu.Unlock()

		fragments = append(fragments, fragmentInfo{instanceID: fragID, x: fx, y: fy})
	}

	// 全服廣播：分裂爆炸
	type fragPayload struct {
		InstanceID string  `json:"instance_id"`
		X          float64 `json:"x"`
		Y          float64 `json:"y"`
		HP         int     `json:"hp"`
		Mult       float64 `json:"mult"`
	}
	var fragPayloads []fragPayload
	for _, f := range fragments {
		fragPayloads = append(fragPayloads, fragPayload{
			InstanceID: f.instanceID,
			X:          f.x,
			Y:          f.y,
			HP:         fragmentHP,
			Mult:       def.MultiplierMin * LuckySplitKillMult,
		})
	}

	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckySplitFish,
		Payload: ws.LuckySplitFishPayload{
			Event:      "split_start",
			InstanceID: instanceID,
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			OriginX:    originX,
			OriginY:    originY,
			Fragments:  fragPayloads,
			DurationSec: int(LuckySplitDuration.Seconds()),
			KillMult:   LuckySplitKillMult,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckySplitFish, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("💥 %s 觸發分裂爆炸！一魚分三！×%.1f 倍率加成！",
			p.DisplayName, LuckySplitKillMult),
		"color": "#FF6B35",
	})
	g.broadcastAnnouncement(ann)

	log.Printf("[LuckySplit] player=%s triggered split instance=%s fragments=%d",
		p.ID, instanceID, LuckySplitCount)

	// 8 秒後觸發二次爆炸
	go g.runLuckySplitBlast(instanceID, p)
}

// runLuckySplitBlast 8 秒後觸發二次爆炸（goroutine）
func (g *Game) runLuckySplitBlast(instanceID string, p *player.Player) {
	time.Sleep(LuckySplitDuration)

	mgr := g.LuckySplitFish
	mgr.mu.Lock()

	// 找出仍存活的分裂碎片
	var survivingFrags []string
	for fragID, entry := range mgr.activeFragments {
		if entry.instanceID == instanceID {
			survivingFrags = append(survivingFrags, fragID)
		}
	}

	// 清除所有此次分裂的碎片記錄
	for _, fragID := range survivingFrags {
		delete(mgr.activeFragments, fragID)
	}
	mgr.mu.Unlock()

	if len(survivingFrags) == 0 {
		// 所有碎片都被玩家打掉了，廣播結束
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckySplitFish,
			Payload: ws.LuckySplitFishPayload{
				Event:      "split_end",
				InstanceID: instanceID,
				BlastCount: 0,
			},
		})
		return
	}

	// 二次爆炸
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	blastCount := 0
	totalReward := 0

	g.mu.Lock()
	for _, fragID := range survivingFrags {
		t, ok := g.Targets[fragID]
		if !ok || t.HP <= 0 {
			continue
		}
		if rng.Float64() < LuckySplitBlastChance {
			// 擊破
			reward := int(float64(t.Multiplier) * LuckySplitBlastMult)
			if reward < 1 {
				reward = 1
			}
			t.HP = 0
			totalReward += reward
			blastCount++
			delete(g.Targets, fragID)
		}
	}
	g.mu.Unlock()

	// 發放獎勵給觸發者
	if totalReward > 0 {
		g.mu.RLock()
		pl, ok := g.Players[p.ID]
		g.mu.RUnlock()
		if ok {
			pl.AddCoins(totalReward)
		}
	}

	log.Printf("[LuckySplit] blast: instance=%s surviving=%d blasted=%d reward=%d",
		instanceID, len(survivingFrags), blastCount, totalReward)

	// 廣播二次爆炸結算
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckySplitFish,
		Payload: ws.LuckySplitFishPayload{
			Event:       "split_blast",
			InstanceID:  instanceID,
			BlastCount:  blastCount,
			TotalReward: totalReward,
		},
	})

	// 全服公告（≥2 個爆炸時）
	if blastCount >= 2 {
		color := "#FF6B35"
		if blastCount >= 3 {
			color = "#FF4500"
		}
		ann := g.Announce.Create(announce.EventLuckySplitFish, p.DisplayName, blastCount, map[string]string{
			"message": fmt.Sprintf("💥 分裂二次爆炸！%d 個碎片爆炸！獎勵 %d 金幣！",
				blastCount, totalReward),
			"color": color,
		})
		g.broadcastAnnouncement(ann)
	}
}
