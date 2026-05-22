// rainbow_prism_handler.go — 彩虹稜鏡魚系統（DAY-213）
// 業界依據：Dive Down 2026「Rainbow is the strongest mutation with 3.0x multiplier」
// + 業界原創「稜鏡折射染色」機制
//
// 設計：擊破 T171 後觸發「稜鏡折射」：
//   - 隨機選場上最多 5 個目標，分別染成 5 種顏色
//   - 紅(×1.5) / 橙(×2.0) / 黃(×2.5) / 綠(×3.0) / 藍(×5.0)
//   - 染色持續 10 秒；染色期間擊破對應顏色目標獲得對應倍率加成（乘法）
//   - 10 秒後「彩虹爆炸」：所有仍存活的染色目標同時爆炸（70% 擊破機率，0.65x 倍率）
//   - 個人冷卻 25 秒；全服廣播染色開始/爆炸結算
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

// prismColor 稜鏡顏色定義
type prismColor struct {
	Name      string  // 顏色名稱
	ColorHex  string  // 顏色 hex（供 Client 顯示）
	MultBonus float64 // 倍率加成（乘法）
}

// prismColors 5 種稜鏡顏色（由低到高）
var prismColors = []prismColor{
	{Name: "red", ColorHex: "#FF4444", MultBonus: 1.5},
	{Name: "orange", ColorHex: "#FF8C00", MultBonus: 2.0},
	{Name: "yellow", ColorHex: "#FFD700", MultBonus: 2.5},
	{Name: "green", ColorHex: "#00CC44", MultBonus: 3.0},
	{Name: "blue", ColorHex: "#0088FF", MultBonus: 5.0},
}

// prismColoredTarget 被染色的目標
type prismColoredTarget struct {
	TargetID  string
	ColorIdx  int
	ColorName string
	MultBonus float64
}

// rainbowPrismManager 彩虹稜鏡魚管理器
type rainbowPrismManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → 下次可觸發時間）
	personalCooldown map[string]time.Time

	// 當前染色目標（instanceID → prismColoredTarget）
	coloredTargets map[string]*prismColoredTarget

	// 是否正在稜鏡模式中
	active bool
}

func newRainbowPrismManager() *rainbowPrismManager {
	return &rainbowPrismManager{
		personalCooldown: make(map[string]time.Time),
		coloredTargets:   make(map[string]*prismColoredTarget),
	}
}

// isRainbowPrismFish 判斷是否為彩虹稜鏡魚
func isRainbowPrismFish(defID string) bool {
	return defID == "T171"
}

// getRainbowPrismMultiplier 取得稜鏡倍率加成（供 handleKill 使用）
// 若目標被染色，回傳對應倍率加成；否則回傳 1.0
func (g *Game) getRainbowPrismMultiplier(instanceID string) float64 {
	g.RainbowPrism.mu.Lock()
	defer g.RainbowPrism.mu.Unlock()

	if ct, ok := g.RainbowPrism.coloredTargets[instanceID]; ok {
		return ct.MultBonus
	}
	return 1.0
}

// removeRainbowPrismColor 移除目標的染色（目標被擊破後呼叫）
func (g *Game) removeRainbowPrismColor(instanceID string) {
	g.RainbowPrism.mu.Lock()
	defer g.RainbowPrism.mu.Unlock()
	delete(g.RainbowPrism.coloredTargets, instanceID)
}

// tryRainbowPrismFish 擊破 T171 後觸發稜鏡折射
func (g *Game) tryRainbowPrismFish(p *player.Player) {
	mgr := g.RainbowPrism
	mgr.mu.Lock()

	// 個人冷卻檢查（25 秒）
	if t, ok := mgr.personalCooldown[p.ID]; ok && time.Now().Before(t) {
		mgr.mu.Unlock()
		return
	}
	// 若已有稜鏡模式進行中，不重複觸發
	if mgr.active {
		mgr.mu.Unlock()
		return
	}

	// 設定個人冷卻
	mgr.personalCooldown[p.ID] = time.Now().Add(25 * time.Second)
	mgr.active = true
	mgr.mu.Unlock()

	// 選取場上最多 5 個目標進行染色
	g.mu.RLock()
	instanceIDs := make([]string, 0, len(g.Targets))
	for id, t := range g.Targets {
		if t.HP > 0 && !isRainbowPrismFish(t.DefID) {
			instanceIDs = append(instanceIDs, id)
		}
	}
	g.mu.RUnlock()

	// 隨機打亂，取最多 5 個
	rand.Shuffle(len(instanceIDs), func(i, j int) {
		instanceIDs[i], instanceIDs[j] = instanceIDs[j], instanceIDs[i]
	})
	count := len(instanceIDs)
	if count > 5 {
		count = 5
	}
	selectedIDs := instanceIDs[:count]

	// 建立染色映射（每個目標對應一種顏色）
	mgr.mu.Lock()
	coloredList := make([]ws.PrismColoredTargetInfo, 0, count)
	for i, tid := range selectedIDs {
		colorIdx := i % len(prismColors)
		ct := &prismColoredTarget{
			TargetID:  tid,
			ColorIdx:  colorIdx,
			ColorName: prismColors[colorIdx].Name,
			MultBonus: prismColors[colorIdx].MultBonus,
		}
		mgr.coloredTargets[tid] = ct
		coloredList = append(coloredList, ws.PrismColoredTargetInfo{
			TargetID:  tid,
			ColorName: prismColors[colorIdx].Name,
			ColorHex:  prismColors[colorIdx].ColorHex,
			MultBonus: prismColors[colorIdx].MultBonus,
		})
	}
	mgr.mu.Unlock()

	log.Printf("[RainbowPrism] player=%s triggered prism refraction, colored %d targets", p.ID, count)

	// 全服廣播：稜鏡折射開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgRainbowPrism,
		Payload: ws.RainbowPrismPayload{
			Event:          "prism_start",
			TriggerPlayer:  p.DisplayName,
			ColoredTargets: coloredList,
			Duration:       10,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventRainbowPrism, p.DisplayName, count, map[string]string{
		"message": fmt.Sprintf("🌈 %s 觸發彩虹稜鏡！場上 %d 個目標被染色！快打高倍率顏色！", p.DisplayName, count),
		"color":   "#FF69B4",
	})
	g.broadcastAnnouncement(ann)

	// 10 秒後觸發彩虹爆炸
	go g.runRainbowPrismBlast(p, 10*time.Second)
}

// runRainbowPrismBlast 等待 10 秒後對所有仍存活的染色目標觸發彩虹爆炸
func (g *Game) runRainbowPrismBlast(p *player.Player, delay time.Duration) {
	time.Sleep(delay)

	mgr := g.RainbowPrism
	mgr.mu.Lock()
	// 取出所有仍在染色列表中的目標
	remainingTargets := make([]string, 0, len(mgr.coloredTargets))
	for tid := range mgr.coloredTargets {
		remainingTargets = append(remainingTargets, tid)
	}
	// 清空染色列表
	mgr.coloredTargets = make(map[string]*prismColoredTarget)
	mgr.active = false
	mgr.mu.Unlock()

	if len(remainingTargets) == 0 {
		// 所有染色目標已被玩家擊破，廣播結算
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgRainbowPrism,
			Payload: ws.RainbowPrismPayload{
				Event:       "prism_blast",
				BlastKills:  0,
				BlastReward: 0,
			},
		})
		return
	}

	// 對仍存活的染色目標觸發彩虹爆炸（70% 擊破機率，0.65x 倍率）
	totalKills := 0
	totalReward := 0

	g.mu.Lock()
	for _, tid := range remainingTargets {
		t, ok := g.Targets[tid]
		if !ok || t.HP <= 0 {
			continue
		}
		if rand.Float64() < 0.70 {
			// 擊破目標
			reward := int(float64(t.Multiplier) * float64(p.BetLevel) * 0.65)
			t.HP = 0
			delete(g.Targets, tid)
			totalKills++
			totalReward += reward
			// 給觸發者獎勵
			p.Coins += reward
		}
	}
	g.mu.Unlock()

	log.Printf("[RainbowPrism] blast: kills=%d, reward=%d", totalKills, totalReward)

	// 全服廣播：彩虹爆炸結算
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgRainbowPrism,
		Payload: ws.RainbowPrismPayload{
			Event:       "prism_blast",
			BlastKills:  totalKills,
			BlastReward: totalReward,
		},
	})

	// 全服公告（≥3 個擊破才公告）
	if totalKills >= 3 {
		color := "#FF69B4"
		if totalKills >= 5 {
			color = "#FF00FF" // 紫紅（5 個以上）
		}
		ann := g.Announce.Create(announce.EventRainbowPrism, p.DisplayName, totalReward, map[string]string{
			"message": fmt.Sprintf("🌈💥 %s 彩虹爆炸！擊破 %d 個目標！獲得 %d 金幣！", p.DisplayName, totalKills, totalReward),
			"color":   color,
		})
		g.broadcastAnnouncement(ann)
	}
}
