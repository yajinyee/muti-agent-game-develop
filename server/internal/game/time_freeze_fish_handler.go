// time_freeze_fish_handler.go — 時間凍結魚系統（DAY-212）
// 業界依據：Evolution Ice Fishing Live 2026「frozen paradise」概念
// + 業界原創設計：「時間停止」讓玩家在凍結期間免費射擊所有靜止目標
//
// 設計：擊破 T170 後觸發「時間凍結」（5秒）：
//   1. 全場所有目標物靜止不動（Speed=0，Server 端停止移動更新）
//   2. 凍結期間玩家射擊命中率 +30%（kill_chance 乘以 1.3，上限 1.0）
//   3. 凍結結束後「解凍爆炸」：所有被命中過的目標 HP 降低 50%
//   4. 個人冷卻 20 秒；全服冷卻 35 秒（防止凍結疊加）
//   5. 全服廣播凍結開始/結束（讓所有玩家知道可以趁機打靜止目標）
//
// 設計差異（與其他時間/速度相關系統的區別）：
//   - 黃金海龜時間停止（DAY-159）：個人 5 秒，只影響個人射擊
//   - 冰凍炸彈魚（DAY-170）：凍結特定目標，不是全場
//   - 時間凍結魚（DAY-212）：全場靜止 + 命中率 +30% + 解凍爆炸
//   - 「全場靜止」讓玩家有「趕快打靜止的魚」的緊迫感
//   - 「命中率 +30%」讓玩家感受到「凍結期間打魚更容易」
//   - 「解凍爆炸」讓被打過的目標在解凍時有額外傷害，製造「凍結期間打越多，解凍後爆炸越多」的策略感
//   - 全服廣播讓所有玩家都知道「現在是凍結期間，趕快打」
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/game/target"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// 時間凍結常數（DAY-212）
const (
	TimeFreezeDuration      = 5     // 凍結持續時間（秒）
	TimeFreezeHitBonus      = 0.30  // 命中率加成（+30%）
	TimeFreezeThawHPRatio   = 0.50  // 解凍爆炸 HP 降低比例（50%）
	TimeFreezeThawKillProb  = 0.60  // 解凍爆炸擊破機率（60%）
	TimeFreezeThawMult      = 0.55  // 解凍爆炸獎勵倍率（0.55x）
	TimeFreezePersonalCD    = 20    // 個人冷卻（秒）
	TimeFreezeGlobalCD      = 35    // 全服冷卻（秒）
)

// timeFreezeManager 時間凍結管理器
type timeFreezeManager struct {
	mu           sync.Mutex
	cooldowns    map[string]time.Time // playerID -> 個人冷卻結束時間
	globalCDEnd  time.Time            // 全服冷卻結束時間
	freezeActive bool                 // 凍結是否正在進行
	freezeEnd    time.Time            // 凍結結束時間
	hitTargets   map[string]bool      // 凍結期間被命中的目標 instanceID
}

func newTimeFreezeManager() *timeFreezeManager {
	return &timeFreezeManager{
		cooldowns:  make(map[string]time.Time),
		hitTargets: make(map[string]bool),
	}
}

// isTimeFreezeFish 判斷是否為時間凍結魚（T170，DAY-212）
func isTimeFreezeFish(defID string) bool {
	return defID == "T170"
}

// isTimeFreezeActive 判斷時間凍結是否正在進行（供 handleAttack 使用）
func (g *Game) isTimeFreezeActive() bool {
	mgr := g.TimeFreeze
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	return mgr.freezeActive && time.Now().Before(mgr.freezeEnd)
}

// getTimeFreezeHitBonus 取得命中率加成（供 combat 使用）
// 凍結期間回傳 0.30（+30%），否則回傳 0.0
func (g *Game) getTimeFreezeHitBonus() float64 {
	if g.isTimeFreezeActive() {
		return TimeFreezeHitBonus
	}
	return 0.0
}

// recordTimeFreezeHit 記錄凍結期間被命中的目標（供 handleAttack 使用）
func (g *Game) recordTimeFreezeHit(instanceID string) {
	mgr := g.TimeFreeze
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	if mgr.freezeActive && time.Now().Before(mgr.freezeEnd) {
		mgr.hitTargets[instanceID] = true
	}
}

// tryTimeFreezeFish 擊破 T170 後觸發時間凍結
func (g *Game) tryTimeFreezeFish(p *player.Player) {
	mgr := g.TimeFreeze
	mgr.mu.Lock()

	// 全服冷卻檢查
	if time.Now().Before(mgr.globalCDEnd) {
		mgr.mu.Unlock()
		return
	}
	// 個人冷卻檢查
	if cd, ok := mgr.cooldowns[p.ID]; ok && time.Now().Before(cd) {
		mgr.mu.Unlock()
		return
	}
	// 凍結已在進行中
	if mgr.freezeActive && time.Now().Before(mgr.freezeEnd) {
		mgr.mu.Unlock()
		return
	}

	// 設定冷卻
	mgr.cooldowns[p.ID] = time.Now().Add(TimeFreezePersonalCD * time.Second)
	mgr.globalCDEnd = time.Now().Add(TimeFreezeGlobalCD * time.Second)

	// 啟動凍結
	mgr.freezeActive = true
	mgr.freezeEnd = time.Now().Add(TimeFreezeDuration * time.Second)
	mgr.hitTargets = make(map[string]bool) // 清空上次記錄
	mgr.mu.Unlock()

	log.Printf("[TimeFreeze] player=%s triggered time freeze for %ds", p.ID, TimeFreezeDuration)

	// 廣播凍結開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgTimeFreezeFish,
		Payload: ws.TimeFreezeFishPayload{
			Event:       "freeze_start",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			DurationSec: TimeFreezeDuration,
			HitBonus:    TimeFreezeHitBonus,
		},
	})

	// 全服公告
	msg := fmt.Sprintf("❄️ %s 觸發時間凍結！全場靜止 %d 秒！命中率 +%.0f%%！趕快打！",
		p.DisplayName, TimeFreezeDuration, TimeFreezeHitBonus*100)
	ann := g.Announce.Create(announce.EventBigWin, p.DisplayName, 0, map[string]string{
		"message": msg,
		"color":   "#00BFFF",
	})
	g.broadcastAnnouncement(ann)

	// 啟動凍結計時器
	go g.runTimeFreezeTimer(p)
}

// runTimeFreezeTimer 凍結計時器（等待結束後觸發解凍爆炸）
func (g *Game) runTimeFreezeTimer(p *player.Player) {
	timer := time.NewTimer(TimeFreezeDuration * time.Second)
	defer timer.Stop()

	<-timer.C

	// 凍結結束
	mgr := g.TimeFreeze
	mgr.mu.Lock()
	mgr.freezeActive = false
	hitTargetIDs := make([]string, 0, len(mgr.hitTargets))
	for id := range mgr.hitTargets {
		hitTargetIDs = append(hitTargetIDs, id)
	}
	mgr.mu.Unlock()

	log.Printf("[TimeFreeze] freeze ended, %d targets were hit during freeze", len(hitTargetIDs))

	// 廣播凍結結束
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgTimeFreezeFish,
		Payload: ws.TimeFreezeFishPayload{
			Event:      "freeze_end",
			HitCount:   len(hitTargetIDs),
		},
	})

	// 解凍爆炸：對被命中的目標執行 HP -50% + 60% 擊破機率
	if len(hitTargetIDs) > 0 {
		time.Sleep(300 * time.Millisecond) // 短暫延遲，讓 Client 播放解凍動畫
		g.doTimeFreezeThaw(p, hitTargetIDs)
	}
}

// doTimeFreezeThaw 解凍爆炸（對凍結期間被命中的目標執行額外傷害）
func (g *Game) doTimeFreezeThaw(p *player.Player, hitTargetIDs []string) {
	type thawResult struct {
		instanceID string
		x, y       float64
		reward     int
		killed     bool
	}

	var results []thawResult
	killedCount := 0
	totalReward := 0

	g.mu.Lock()
	playerCount := len(g.Players)
	if playerCount < 1 {
		playerCount = 1
	}

	for _, id := range hitTargetIDs {
		t, ok := g.Targets[id]
		if !ok || !t.IsAlive {
			continue
		}

		// HP -50%
		reduction := int(float64(t.HP) * TimeFreezeThawHPRatio)
		if reduction < 1 {
			reduction = 1
		}
		t.HP -= reduction

		// 60% 擊破機率
		killed := false
		reward := 0
		if t.HP <= 0 || rand.Float64() < TimeFreezeThawKillProb {
			t.IsAlive = false
			t.HP = 0
			killed = true
			def := t.Def
			mult := (def.MultiplierMin + def.MultiplierMax) / 2.0
			reward = int(mult * TimeFreezeThawMult)
			if reward < 1 {
				reward = 1
			}
			killedCount++
			totalReward += reward
		}

		results = append(results, thawResult{
			instanceID: id,
			x:          t.X,
			y:          t.Y,
			reward:     reward,
			killed:     killed,
		})
	}

	// 按玩家數平均分配獎勵
	rewardPerPlayer := totalReward / playerCount
	if rewardPerPlayer < 1 && totalReward > 0 {
		rewardPerPlayer = 1
	}
	for _, pp := range g.Players {
		pp.Coins += rewardPerPlayer
	}
	g.mu.Unlock()

	log.Printf("[TimeFreeze] thaw: hit=%d killed=%d totalReward=%d perPlayer=%d",
		len(hitTargetIDs), killedCount, totalReward, rewardPerPlayer)

	// 廣播解凍爆炸結果
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgTimeFreezeFish,
		Payload: ws.TimeFreezeFishPayload{
			Event:           "thaw_blast",
			PlayerID:        p.ID,
			PlayerName:      p.DisplayName,
			HitCount:        len(hitTargetIDs),
			KilledCount:     killedCount,
			TotalReward:     totalReward,
			RewardPerPlayer: rewardPerPlayer,
		},
	})

	// 全服公告（≥4 個擊破）
	if killedCount >= 4 {
		color := "#00BFFF"
		if killedCount >= 8 {
			color = "#1E90FF"
		}
		msg := fmt.Sprintf("💥 解凍爆炸！擊破 %d 個目標！每位玩家獲得 %d 金幣！",
			killedCount, rewardPerPlayer)
		ann := g.Announce.Create(announce.EventBigWin, p.DisplayName, totalReward, map[string]string{
			"message": msg,
			"color":   color,
		})
		g.broadcastAnnouncement(ann)
	}

	// 廣播每個被擊破的目標（讓 Client 播放死亡動畫）
	for _, r := range results {
		if r.killed {
			g.Hub.Broadcast(&ws.Message{
				Type: ws.MsgTargetKill,
				Payload: ws.TargetKillPayload{
					InstanceID: r.instanceID,
					DefID:      "T170_thaw",
					Multiplier: TimeFreezeThawMult,
					Reward:     r.reward,
					KillerID:   p.ID,
				},
			})
		}
	}
}

// notifyTimeFreezeTargetHit 通知凍結期間目標被命中（供 handleAttack 使用）
// 在 handleAttack 中，命中目標後呼叫此函數記錄
func (g *Game) notifyTimeFreezeTargetHit(t *target.Target) {
	if t == nil {
		return
	}
	g.recordTimeFreezeHit(t.InstanceID)
}
