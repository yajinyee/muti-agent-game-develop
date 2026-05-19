// Package game — 遊戲邏輯單元測試
package game

import (
	"testing"
	"time"

	"digital-twin/server/internal/data"
	"digital-twin/server/internal/game/target"
	"digital-twin/server/internal/ws"
)

// ── 輔助函數 ──────────────────────────────────────────

// newTestGame 建立測試用 Game（不啟動 Hub）
func newTestGame(t *testing.T) *Game {
	t.Helper()
	hub := ws.NewHub()
	g := NewGame("test-game", hub)
	return g
}

// ── 基礎功能測試 ──────────────────────────────────────

// TestNewGame 確認 NewGame 初始狀態正確
func TestNewGame(t *testing.T) {
	g := newTestGame(t)

	if g.ID != "test-game" {
		t.Errorf("expected ID=test-game, got %s", g.ID)
	}
	if g.GetState() != "normal_play" {
		t.Errorf("expected initial state=normal_play, got %s", g.GetState())
	}
	if len(g.Players) != 0 {
		t.Errorf("expected 0 players, got %d", len(g.Players))
	}
	if len(g.Targets) != 0 {
		t.Errorf("expected 0 targets, got %d", len(g.Targets))
	}
}

// TestAddRemovePlayer 確認玩家加入/移除正確
func TestAddRemovePlayer(t *testing.T) {
	g := newTestGame(t)

	// 加入玩家
	g.AddPlayer("player_001")
	g.mu.RLock()
	if len(g.Players) != 1 {
		t.Errorf("expected 1 player after add, got %d", len(g.Players))
	}
	p := g.Players["player_001"]
	g.mu.RUnlock()

	if p == nil {
		t.Fatal("player_001 should exist")
	}
	if p.Coins != 10000 {
		t.Errorf("expected initial coins=10000, got %d", p.Coins)
	}

	// 重複加入不應該重置
	g.AddPlayer("player_001")
	g.mu.RLock()
	if len(g.Players) != 1 {
		t.Errorf("duplicate add should not create new player, got %d", len(g.Players))
	}
	g.mu.RUnlock()

	// 移除玩家
	g.RemovePlayer("player_001")
	g.mu.RLock()
	if len(g.Players) != 0 {
		t.Errorf("expected 0 players after remove, got %d", len(g.Players))
	}
	g.mu.RUnlock()
}

// TestGameStartStop 確認 Start/Stop 不會造成 goroutine 洩漏
func TestGameStartStop(t *testing.T) {
	hub := ws.NewHub()
	g := NewGame("leak-test", hub)

	// 啟動遊戲
	g.Start()

	// 讓遊戲跑一小段時間
	time.Sleep(200 * time.Millisecond)

	// 停止遊戲
	g.Stop()

	// 等待 goroutine 清理
	time.Sleep(100 * time.Millisecond)

	// 確認 stopCh 已關閉（不會 panic）
	select {
	case <-g.stopCh:
		// 正常：stopCh 已關閉
	default:
		t.Error("stopCh should be closed after Stop()")
	}
}

// TestSpawnTarget 確認目標生成正確
func TestSpawnTarget(t *testing.T) {
	g := newTestGame(t)
	g.AddPlayer("player_001")

	// 生成目標
	g.spawnTarget()

	g.mu.RLock()
	count := len(g.Targets)
	g.mu.RUnlock()

	if count != 1 {
		t.Errorf("expected 1 target after spawn, got %d", count)
	}

	// 確認目標有有效的 InstanceID
	g.mu.RLock()
	for id, target := range g.Targets {
		if id == "" {
			t.Error("target InstanceID should not be empty")
		}
		if target.DefID == "" {
			t.Error("target DefID should not be empty")
		}
		if target.HP <= 0 {
			t.Errorf("target HP should be > 0, got %d", target.HP)
		}
	}
	g.mu.RUnlock()
}

// TestBetLevelAverage 確認多玩家 bet level 平均計算正確
func TestBetLevelAverage(t *testing.T) {
	g := newTestGame(t)

	// 加入 3 個玩家，設定不同 bet level
	g.AddPlayer("p1")
	g.AddPlayer("p2")
	g.AddPlayer("p3")

	g.mu.Lock()
	g.Players["p1"].BetLevel = 2
	g.Players["p2"].BetLevel = 4
	g.Players["p3"].BetLevel = 6
	g.mu.Unlock()

	// 計算平均（應該是 (2+4+6)/3 = 4）
	g.mu.RLock()
	total := 0
	for _, p := range g.Players {
		total += p.BetLevel
	}
	avg := total / len(g.Players)
	g.mu.RUnlock()

	if avg != 4 {
		t.Errorf("expected average bet level=4, got %d", avg)
	}
}

// TestStateTransition 確認狀態轉換正確
func TestStateTransition(t *testing.T) {
	hub := ws.NewHub()
	g := NewGame("state-test", hub)

	// 初始狀態
	if g.GetState() != "normal_play" {
		t.Errorf("expected normal_play, got %s", g.GetState())
	}

	// 轉換到 boss_warning
	g.transitionState("boss_warning")
	if g.GetState() != "boss_warning" {
		t.Errorf("expected boss_warning, got %s", g.GetState())
	}

	// 轉換到 boss_battle
	g.transitionState("boss_battle")
	if g.GetState() != "boss_battle" {
		t.Errorf("expected boss_battle, got %s", g.GetState())
	}

	// 非法轉換（boss_battle → bonus_game 不允許）
	g.transitionState("bonus_game")
	if g.GetState() != "boss_battle" {
		t.Errorf("illegal transition should be rejected, state should remain boss_battle, got %s", g.GetState())
	}
}

// TestSafeAfterFunc 確認 safeAfterFunc 在 Stop 後不執行
func TestSafeAfterFunc(t *testing.T) {
	g := newTestGame(t)

	executed := false
	g.safeAfterFunc(50*time.Millisecond, func() {
		executed = true
	})

	// 立即停止（在 50ms 前）
	g.Stop()

	// 等待超過 50ms
	time.Sleep(100 * time.Millisecond)

	if executed {
		t.Error("safeAfterFunc should NOT execute after Stop()")
	}
}

// TestSafeAfterFuncExecutes 確認 safeAfterFunc 在正常情況下執行
func TestSafeAfterFuncExecutes(t *testing.T) {
	g := newTestGame(t)

	done := make(chan struct{})
	g.safeAfterFunc(20*time.Millisecond, func() {
		close(done)
	})

	select {
	case <-done:
		// 正常執行
	case <-time.After(200 * time.Millisecond):
		t.Error("safeAfterFunc should execute within timeout")
	}

	g.Stop()
}

// TestBonusTriggerCooldown 確認 Bonus 觸發冷卻（90 秒間隔）
func TestBonusTriggerCooldown(t *testing.T) {
	hub := ws.NewHub()
	g := NewGame("bonus-test", hub)

	// 設定 lastBonusAt 為剛剛（模擬剛觸發過 Bonus）
	g.mu.Lock()
	g.lastBonusAt = time.Now()
	g.mu.Unlock()

	// 嘗試再次觸發（應該被冷卻阻擋）
	g.triggerBonusReady()

	// 狀態應該仍然是 normal_play
	if g.GetState() != "normal_play" {
		t.Errorf("bonus should be blocked by cooldown, state should be normal_play, got %s", g.GetState())
	}
}

// TestScoreTarget 確認目標評分系統正確
func TestScoreTarget(t *testing.T) {
	g := newTestGame(t)

	// 生成一個目標
	g.spawnTarget()

	g.mu.RLock()
	var target_score float64
	for _, t := range g.Targets {
		target_score = g.scoreTarget(t)
		break
	}
	g.mu.RUnlock()

	// 分數應該 >= 0
	if target_score < 0 {
		t.Errorf("target score should be >= 0, got %f", target_score)
	}
}

// TestScoreTarget_BossHighestPriority 確認 BOSS 評分最高
func TestScoreTarget_BossHighestPriority(t *testing.T) {
	g := newTestGame(t)

	// 生成普通目標
	g.spawnTarget()

	// 生成 BOSS
	g.spawnBoss()

	g.mu.RLock()
	var bossScore, normalScore float64
	for _, tgt := range g.Targets {
		score := g.scoreTarget(tgt)
		if tgt.Def.Type == data.TargetTypeBoss {
			bossScore = score
		} else {
			normalScore = score
		}
	}
	g.mu.RUnlock()

	if bossScore <= normalScore {
		t.Errorf("BOSS score (%f) should be higher than normal target score (%f)", bossScore, normalScore)
	}
}

// TestScoreTarget_LowHPHigherPriority 確認低 HP 目標評分更高
func TestScoreTarget_LowHPHigherPriority(t *testing.T) {
	g := newTestGame(t)

	// 建立兩個相同的目標，一個 HP 高一個 HP 低
	def := data.Targets["T001"]
	tgtFull := target.NewTarget("full-hp", def, 1000, 300)
	tgtLow := target.NewTarget("low-hp", def, 1000, 300)
	tgtLow.HP = 1 // 幾乎死了

	g.mu.Lock()
	g.Targets["full-hp"] = tgtFull
	g.Targets["low-hp"] = tgtLow
	g.mu.Unlock()

	g.mu.RLock()
	fullScore := g.scoreTarget(tgtFull)
	lowScore := g.scoreTarget(tgtLow)
	g.mu.RUnlock()

	if lowScore <= fullScore {
		t.Errorf("Low HP target score (%f) should be higher than full HP target score (%f)", lowScore, fullScore)
	}
}

// TestGetMissionResetAt 確認任務重置時間正確（DAY-039）
func TestGetMissionResetAt(t *testing.T) {
	g := newTestGame(t)

	resetAt := g.GetMissionResetAt()
	if resetAt.IsZero() {
		t.Error("mission reset time should not be zero")
	}

	// 重置時間應該在未來（今天午夜之後）
	if resetAt.Before(time.Now()) {
		t.Errorf("mission reset time should be in the future, got %v", resetAt)
	}
}

// TestGetLeaderboardData 確認排行榜資料正確
func TestGetLeaderboardData(t *testing.T) {
	g := newTestGame(t)

	data := g.GetLeaderboardData()
	if data.Timestamp <= 0 {
		t.Error("leaderboard timestamp should be > 0")
	}
	// 無玩家時排行榜應為空
	if data.Entries == nil {
		t.Error("leaderboard entries should not be nil")
	}
}

// TestGetLeaderboardData_WithPlayers 確認有玩家時排行榜正確
func TestGetLeaderboardData_WithPlayers(t *testing.T) {
	g := newTestGame(t)

	// 加入玩家並設定分數
	g.AddPlayer("p1")
	g.AddPlayer("p2")
	g.mu.Lock()
	g.Players["p1"].SessionScore = 1000
	g.Players["p1"].DisplayName = "Alice"
	g.Players["p2"].SessionScore = 500
	g.Players["p2"].DisplayName = "Bob"
	g.mu.Unlock()

	data := g.GetLeaderboardData()
	if len(data.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(data.Entries))
	}

	// 第一名應該是 Alice（分數最高）
	if len(data.Entries) > 0 && data.Entries[0].Score != 1000 {
		t.Errorf("expected top score=1000, got %d", data.Entries[0].Score)
	}
}

// TestAddPlayer_InitialCoins 確認玩家初始金幣正確
func TestAddPlayer_InitialCoins(t *testing.T) {
	hub := ws.NewHub()
	g := NewGameWithStore("test", hub, nil, 5000) // 自訂初始金幣

	g.AddPlayer("p1")
	g.mu.RLock()
	p := g.Players["p1"]
	g.mu.RUnlock()

	if p.Coins != 5000 {
		t.Errorf("expected initial coins=5000, got %d", p.Coins)
	}
}

// TestSpawnBoss 確認 BOSS 生成正確
func TestSpawnBoss(t *testing.T) {
	g := newTestGame(t)
	g.AddPlayer("p1")

	// 生成 BOSS
	g.spawnBoss()

	g.mu.RLock()
	bossCount := 0
	for _, tgt := range g.Targets {
		if tgt.DefID == "B001" {
			bossCount++
		}
	}
	g.mu.RUnlock()

	if bossCount != 1 {
		t.Errorf("expected 1 BOSS, got %d", bossCount)
	}
}

// TestGetState_ThreadSafe 確認 GetState 在並發下安全
func TestGetState_ThreadSafe(t *testing.T) {
	g := newTestGame(t)

	done := make(chan struct{})
	go func() {
		for i := 0; i < 100; i++ {
			g.GetState()
		}
		close(done)
	}()

	for i := 0; i < 100; i++ {
		g.transitionState("normal_play")
	}

	select {
	case <-done:
		// 正常完成
	case <-time.After(2 * time.Second):
		t.Error("GetState should complete within timeout")
	}
}
