// ghost_fish_handler.go — 幽靈魚分身系統（DAY-198）
// 原創設計，靈感來自 Fisch Phantom Mutation（4x 倍率）和業界「分身」概念
// 設計：T156 幽靈魚出現時，同時在場上生成 2-3 個「幻影分身」（外觀相同但 HP=1）：
//   1. 玩家需要找出「真身」（HP 正常）才能獲得完整獎勵
//   2. 擊破幻影分身：給 1x betLevel 小獎勵（安慰獎），廣播「幻影消散」
//   3. 擊破真身：觸發「幽靈爆發」— 所有幻影分身同時爆炸（50% 擊破機率，0.50x 倍率）
//   4. 全服廣播幽靈魚出現（但不告訴玩家哪個是真身）
//   5. 幽靈魚離開時，所有幻影分身一起消失
// 設計差異：
//   - 與普通目標（直接擊破）不同，幽靈魚製造「哪條是真的？」的懸疑感
//   - 與連鎖爆炸魚（靜態爆炸）不同，幽靈魚的爆炸是「找到真身才觸發」，
//     讓玩家有「我找到了！」的成就感
//   - 幻影分身的安慰獎讓玩家不會完全空手，但真身的爆炸獎勵更豐厚
//   - 全服廣播讓所有玩家都在「找真身」，製造競爭感
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"

	"digital-twin/server/internal/data"
	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/game/target"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// 幽靈魚常數
const (
	GhostFishCooldownSec    = 35    // 全服冷卻 35 秒
	GhostFishCloneMin       = 2     // 最少幻影分身數
	GhostFishCloneMax       = 3     // 最多幻影分身數
	GhostFishPhantomReward  = 1     // 幻影分身安慰獎（1x betLevel）
	GhostFishExplodeChance  = 0.50  // 幻影爆炸擊破機率 50%
	GhostFishExplodeMult    = 0.50  // 幻影爆炸獎勵倍率 0.50x
)

// ghostFishSession 幽靈魚會話（追蹤真身和幻影分身）
type ghostFishSession struct {
	realID     string   // 真身 InstanceID
	cloneIDs   []string // 幻影分身 InstanceID 列表
	multiplier float64  // 真身倍率
	killerID   string   // 擊破真身的玩家 ID（空=未被擊破）
}

// ghostFishManager 幽靈魚管理器（全服共享）
type ghostFishManager struct {
	mu          sync.Mutex
	session     *ghostFishSession
	cooldownEnd time.Time
}

func newGhostFishManager() *ghostFishManager {
	return &ghostFishManager{}
}

// isGhostFish 判斷是否為幽靈魚真身（T156）
func isGhostFish(defID string) bool {
	return defID == "T156"
}

// isGhostFishClone 判斷是否為幽靈魚幻影分身（T156C）
func isGhostFishClone(defID string) bool {
	return defID == "T156C"
}

// notifyGhostFishSpawn T156 生成時觸發幻影分身
func (g *Game) notifyGhostFishSpawn(instanceID string, x, y float64, multiplier float64) {
	mgr := g.GhostFish
	mgr.mu.Lock()

	// 全服冷卻檢查
	if mgr.session != nil || time.Now().Before(mgr.cooldownEnd) {
		mgr.mu.Unlock()
		return
	}

	// 決定幻影分身數量（2-3 個）
	cloneCount := GhostFishCloneMin + rand.Intn(GhostFishCloneMax-GhostFishCloneMin+1)

	session := &ghostFishSession{
		realID:     instanceID,
		multiplier: multiplier,
		cloneIDs:   make([]string, 0, cloneCount),
	}
	mgr.session = session
	mgr.mu.Unlock()

	// 生成幻影分身（在真身附近隨機位置）
	cloneDef := data.Targets["T156C"]
	if cloneDef == nil {
		// T156C 不存在時降級處理
		log.Printf("[GhostFish] T156C not found, skipping clones")
		return
	}

	cloneIDs := make([]string, 0, cloneCount)
	for i := 0; i < cloneCount; i++ {
		cloneID := uuid.New().String()
		// 在真身附近 ±200px 隨機位置
		cloneX := x + (rand.Float64()*400 - 200)
		cloneY := y + (rand.Float64()*300 - 150)
		// 確保在畫面範圍內
		if cloneX < 100 {
			cloneX = 100
		}
		if cloneX > 1180 {
			cloneX = 1180
		}
		if cloneY < 50 {
			cloneY = 50
		}
		if cloneY > 650 {
			cloneY = 650
		}

		clone := target.NewTarget(cloneID, cloneDef, cloneX, cloneY)
		g.mu.Lock()
		g.Targets[cloneID] = clone
		g.mu.Unlock()

		cloneIDs = append(cloneIDs, cloneID)

		// 廣播幻影分身生成
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgTargetSpawn,
			Payload: ws.TargetSpawnPayload{
				InstanceID: cloneID,
				DefID:      cloneDef.ID,
				Name:       cloneDef.Name,
				Type:       string(cloneDef.Type),
				X:          cloneX,
				Y:          cloneY,
				HP:         cloneDef.HP,
				MaxHP:      cloneDef.HP,
				Multiplier: (cloneDef.MultiplierMin + cloneDef.MultiplierMax) / 2,
				Speed:      cloneDef.Speed,
				Lifetime:   cloneDef.Lifetime,
			},
		})
	}

	// 更新 session 的 cloneIDs
	mgr.mu.Lock()
	if mgr.session != nil && mgr.session.realID == instanceID {
		mgr.session.cloneIDs = cloneIDs
	}
	mgr.mu.Unlock()

	// 廣播幽靈魚出現（全服，不告訴玩家哪個是真身）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgGhostFish,
		Payload: ws.GhostFishPayload{
			Phase:      "ghost_appear",
			RealID:     instanceID,
			CloneIDs:   cloneIDs,
			CloneCount: cloneCount,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventMegaWin, "幽靈魚", 0, map[string]string{
		"message": fmt.Sprintf("👻 幽靈魚出現！場上有 %d 個分身，哪個是真身？快找！", cloneCount+1),
	})
	g.broadcastAnnouncement(ann)

	log.Printf("[GhostFish] real=%s clones=%v spawned", instanceID, cloneIDs)
}

// notifyGhostFishCloneKill 玩家擊破幻影分身
func (g *Game) notifyGhostFishCloneKill(p *player.Player, cloneID string) {
	mgr := g.GhostFish
	mgr.mu.Lock()
	if mgr.session == nil {
		mgr.mu.Unlock()
		return
	}
	// 確認是本次 session 的幻影分身
	isMyClone := false
	for _, id := range mgr.session.cloneIDs {
		if id == cloneID {
			isMyClone = true
			break
		}
	}
	mgr.mu.Unlock()

	if !isMyClone {
		return
	}

	// 給安慰獎
	reward := GhostFishPhantomReward * p.BetLevel
	g.mu.Lock()
	p.Coins += reward
	g.mu.Unlock()

	// 廣播幻影消散
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgGhostFish,
		Payload: ws.GhostFishPayload{
			Phase:    "phantom_vanish",
			CloneID:  cloneID,
			KillerID: p.ID,
			KillerName: p.DisplayName,
			Reward:   reward,
		},
	})

	log.Printf("[GhostFish] player=%s hit phantom clone=%s reward=%d", p.ID, cloneID, reward)
}

// notifyGhostFishRealKill 玩家擊破幽靈魚真身，觸發幽靈爆發
func (g *Game) notifyGhostFishRealKill(p *player.Player, instanceID string, baseMult float64) {
	mgr := g.GhostFish
	mgr.mu.Lock()
	if mgr.session == nil || mgr.session.realID != instanceID {
		mgr.mu.Unlock()
		return
	}
	cloneIDs := make([]string, len(mgr.session.cloneIDs))
	copy(cloneIDs, mgr.session.cloneIDs)
	mgr.session.killerID = p.ID
	mgr.mu.Unlock()

	// 廣播找到真身（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgGhostFish,
		Payload: ws.GhostFishPayload{
			Phase:      "real_found",
			RealID:     instanceID,
			KillerID:   p.ID,
			KillerName: p.DisplayName,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventMegaWin, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("👻💥 %s 找到幽靈魚真身！幻影爆炸開始！", p.DisplayName),
	})
	g.broadcastAnnouncement(ann)

	// 幻影爆炸：所有幻影分身同時爆炸
	time.Sleep(300 * time.Millisecond)
	explodeKills := 0
	explodeReward := 0

	for _, cloneID := range cloneIDs {
		g.mu.RLock()
		t, ok := g.Targets[cloneID]
		g.mu.RUnlock()
		if !ok {
			continue // 幻影分身已被其他玩家擊破
		}

		if rand.Float64() < GhostFishExplodeChance {
			r := int(t.Multiplier * float64(p.BetLevel) * GhostFishExplodeMult)
			g.mu.Lock()
			if _, stillOk := g.Targets[cloneID]; stillOk {
				delete(g.Targets, cloneID)
				p.Coins += r
				explodeKills++
				explodeReward += r
			}
			g.mu.Unlock()
		}
	}

	// 廣播幻影爆炸結果（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgGhostFish,
		Payload: ws.GhostFishPayload{
			Phase:         "ghost_explode",
			RealID:        instanceID,
			KillerID:      p.ID,
			KillerName:    p.DisplayName,
			ExplodeKills:  explodeKills,
			ExplodeReward: explodeReward,
		},
	})

	// 清理 session
	mgr.mu.Lock()
	mgr.session = nil
	mgr.cooldownEnd = time.Now().Add(GhostFishCooldownSec * time.Second)
	mgr.mu.Unlock()

	// 全服公告（≥1 個幻影爆炸）
	if explodeKills >= 1 {
		ann2 := g.Announce.Create(announce.EventMegaWin, p.DisplayName, explodeReward, map[string]string{
			"message": fmt.Sprintf("👻✨ %s 幽靈爆炸！擊破 %d 個幻影！獎勵 %d 金幣！",
				p.DisplayName, explodeKills, explodeReward),
		})
		g.broadcastAnnouncement(ann2)
	}

	log.Printf("[GhostFish] player=%s found real, explode: kills=%d reward=%d",
		p.ID, explodeKills, explodeReward)
}

// onGhostFishLeave 幽靈魚真身離開畫面（未被擊破）
func (g *Game) onGhostFishLeave(instanceID string) {
	mgr := g.GhostFish
	mgr.mu.Lock()
	if mgr.session == nil || mgr.session.realID != instanceID {
		mgr.mu.Unlock()
		return
	}
	cloneIDs := make([]string, len(mgr.session.cloneIDs))
	copy(cloneIDs, mgr.session.cloneIDs)
	mgr.session = nil
	mgr.cooldownEnd = time.Now().Add(GhostFishCooldownSec * time.Second)
	mgr.mu.Unlock()

	// 移除所有幻影分身
	for _, cloneID := range cloneIDs {
		g.mu.Lock()
		delete(g.Targets, cloneID)
		g.mu.Unlock()
	}

	// 廣播幽靈魚逃跑（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgGhostFish,
		Payload: ws.GhostFishPayload{
			Phase:    "ghost_escape",
			RealID:   instanceID,
			CloneIDs: cloneIDs,
		},
	})

	log.Printf("[GhostFish] real=%s escaped, removed %d clones", instanceID, len(cloneIDs))
}
