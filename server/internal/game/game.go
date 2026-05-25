// Package game — 遊戲主循環與狀態機
// server-core-agent 負責維護
package game

import (
	"encoding/json"
	"log"
	"math/rand"
	"sync"
	"time"

	"chiikawa-game/internal/data"
	"chiikawa-game/internal/protocol"
	"chiikawa-game/internal/ws"
)

// GameState 遊戲狀態
type GameState string

const (
	StateNormalPlay  GameState = "normal_play"
	StateBossWarning GameState = "boss_warning"
	StateBossBattle  GameState = "boss_battle"
	StateBossResult  GameState = "boss_result"
	StateBonusReady  GameState = "bonus_ready"
	StateBonusGame   GameState = "bonus_game"
	StateBonusResult GameState = "bonus_result"
)

// Game 遊戲實例
type Game struct {
	mu      sync.RWMutex
	hub     *ws.Hub
	state   GameState
	players map[string]*Player
	targets map[string]*Target
	boss    *Target

	// 計時器
	spawnTimer    float64
	bossTimer     float64 // BOSS 戰剩餘時間
	bossWarningAt time.Time
	nextBossIn    float64 // 下次 BOSS 觸發倒數（秒）
	bonusPlayer   string  // 觸發 Bonus 的玩家 ID
	bonusTimer    float64
	bonusScore    map[string]int // playerID -> score
	bonusWeedHP   map[string]int // weedID -> hp（BG002 需要連點2次）

	// 幸運特殊魚系統
	luckyChainLightning  *luckyChainLightningManager
	luckyCrabTorpedo     *luckyCrabTorpedoManager
	luckyVortex          *luckyVortexManager
	luckyGoldenDragon    *luckyGoldenDragonManager
	luckyThunderLobster  *luckyThunderLobsterManager
	luckyAwakenedPhoenix *luckyAwakenedPhoenixManager
	luckyShockwaveBomb   *luckyShockwaveBombManager
	luckyDrillTorpedo    *luckyDrillTorpedoManager
	luckyTimeFreeze      *luckyTimeFreezeManager
	luckyChainExplosion  *luckyChainExplosionManager
	// DAY-295 新增
	luckyChainLongKing  *luckyChainLongKingManager
	luckyDragonShotgun  *luckyDragonShotgunManager
	luckyRocketCannon   *luckyRocketCannonManager
	luckyDeepWhirlpool  *luckyDeepWhirlpoolManager
	luckyVampireMult    *luckyVampireMultManager
	// DAY-296 新增
	luckyMirrorFish   *luckyMirrorFishManager
	luckyGoldenRain   *luckyGoldenRainManager
	luckyFreezeBomb   *luckyFreezeBombManager
	luckyThunderStorm *luckyThunderStormManager
	luckyLuckyWheel   *luckyLuckyWheelManager

	lastTick time.Time
}

func NewGame(hub *ws.Hub) *Game {
	g := &Game{
		hub:        hub,
		state:      StateNormalPlay,
		players:    make(map[string]*Player),
		targets:    make(map[string]*Target),
		bonusScore: make(map[string]int),
		bonusWeedHP: make(map[string]int),
		lastTick:   time.Now(),

		// 初始化幸運特殊魚系統
		luckyChainLightning:  newLuckyChainLightningManager(),
		luckyCrabTorpedo:     newLuckyCrabTorpedoManager(),
		luckyVortex:          newLuckyVortexManager(),
		luckyGoldenDragon:    newLuckyGoldenDragonManager(),
		luckyThunderLobster:  newLuckyThunderLobsterManager(),
		luckyAwakenedPhoenix: newLuckyAwakenedPhoenixManager(),
		luckyShockwaveBomb:   newLuckyShockwaveBombManager(),
		luckyDrillTorpedo:    newLuckyDrillTorpedoManager(),
		luckyTimeFreeze:      newLuckyTimeFreezeManager(),
		luckyChainExplosion:  newLuckyChainExplosionManager(),
		// DAY-295 新增
		luckyChainLongKing:  newLuckyChainLongKingManager(),
		luckyDragonShotgun:  newLuckyDragonShotgunManager(),
		luckyRocketCannon:   newLuckyRocketCannonManager(),
		luckyDeepWhirlpool:  newLuckyDeepWhirlpoolManager(),
		luckyVampireMult:    newLuckyVampireMultManager(),
		// DAY-296 新增
		luckyMirrorFish:   newLuckyMirrorFishManager(),
		luckyGoldenRain:   newLuckyGoldenRainManager(),
		luckyFreezeBomb:   newLuckyFreezeBombManager(),
		luckyThunderStorm: newLuckyThunderStormManager(),
		luckyLuckyWheel:   newLuckyLuckyWheelManager(),
	}
	g.nextBossIn = 180 + rand.Float64()*120 // 3-5 分鐘
	return g
}

// Start 啟動遊戲循環
func (g *Game) Start() {
	go g.loop()
}

// AddPlayer 玩家加入
func (g *Game) AddPlayer(id string) {
	g.mu.Lock()
	p := NewPlayer(id)
	g.players[id] = p
	g.mu.Unlock()

	log.Printf("[Game] Player joined: %s", id)
	g.sendPlayerUpdate(id)
	g.hub.Send(id, protocol.MsgGameState, protocol.GameStatePayload{
		State:     string(g.state),
		Timestamp: time.Now().UnixMilli(),
	})
	// 把現有目標物同步給新玩家
	g.mu.RLock()
	for _, t := range g.targets {
		g.hub.Send(id, protocol.MsgTargetSpawn, g.targetSpawnPayload(t))
	}
	if g.boss != nil {
		g.hub.Send(id, protocol.MsgTargetSpawn, g.targetSpawnPayload(g.boss))
	}
	g.mu.RUnlock()
}

// RemovePlayer 玩家離開
func (g *Game) RemovePlayer(id string) {
	g.mu.Lock()
	delete(g.players, id)
	g.mu.Unlock()
	log.Printf("[Game] Player left: %s", id)
}

// HandleMessage 處理玩家訊息
func (g *Game) HandleMessage(clientID string, msgType string, payload json.RawMessage) {
	switch msgType {
	case protocol.MsgAttack:
		var req protocol.AttackRequest
		if err := json.Unmarshal(payload, &req); err != nil {
			return
		}
		g.handleAttack(clientID, req)
	case protocol.MsgAutoToggle:
		g.handleAutoToggle(clientID)
	case protocol.MsgBetChange:
		var req protocol.BetChangeRequest
		if err := json.Unmarshal(payload, &req); err != nil {
			return
		}
		g.handleBetChange(clientID, req.BetLevel)
	case protocol.MsgLock:
		var req protocol.LockRequest
		if err := json.Unmarshal(payload, &req); err != nil {
			return
		}
		g.handleLock(clientID, req.TargetID)
	case protocol.MsgBonusClick:
		var req protocol.BonusClickRequest
		if err := json.Unmarshal(payload, &req); err != nil {
			return
		}
		g.handleBonusClick(clientID, req)
	case protocol.MsgPing:
		g.hub.Send(clientID, protocol.MsgPong, map[string]interface{}{})
	case protocol.MsgTriggerBoss:
		g.triggerBoss()
	case protocol.MsgTriggerBonus:
		g.triggerBonus(clientID)
	case protocol.MsgCollectGoldenCoin:
		var req protocol.CollectGoldenCoinRequest
		if err := json.Unmarshal(payload, &req); err != nil {
			return
		}
		g.collectGoldenCoin(clientID, req.CoinID)
	}
}

// ── 主循環 ────────────────────────────────────────────────────

func (g *Game) loop() {
	ticker := time.NewTicker(50 * time.Millisecond) // 20 ticks/sec
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		delta := now.Sub(g.lastTick).Seconds()
		g.lastTick = now
		g.tick(delta)
	}
}

func (g *Game) tick(delta float64) {
	g.mu.Lock()
	defer g.mu.Unlock()

	switch g.state {
	case StateNormalPlay:
		g.tickNormal(delta)
	case StateBossWarning:
		// 警告 3 秒後進入 BOSS 戰
		if time.Since(g.bossWarningAt).Seconds() >= 3 {
			g.spawnBoss()
		}
	case StateBossBattle:
		g.tickBoss(delta)
	case StateBonusGame:
		g.tickBonus(delta)
	}
}

func (g *Game) tickNormal(delta float64) {
	// 目標物生成
	g.spawnTimer += delta
	if g.spawnTimer >= SpawnInterval {
		g.spawnTimer = 0
		if len(g.targets) < MaxTargets {
			g.spawnTarget()
		}
	}

	// 移除超時目標物
	for id, t := range g.targets {
		if t.IsExpired() {
			delete(g.targets, id)
			g.hub.Broadcast(protocol.MsgTargetKill, protocol.TargetKillPayload{
				InstanceID: t.InstanceID,
				DefID:      t.Def.ID,
				Multiplier: t.Multiplier,
				Reward:     0,
				LaborGain:  0,
				KillerID:   "",
			})
		}
	}

	// BOSS 觸發計時
	g.nextBossIn -= delta
	if g.nextBossIn <= 0 && len(g.players) > 0 {
		g.startBossWarning()
	}

	// AUTO 射擊
	for _, p := range g.players {
		if p.IsAuto {
			g.autoFire(p)
		}
	}
}

func (g *Game) tickBoss(delta float64) {
	if g.boss == nil {
		g.setState(StateNormalPlay)
		return
	}
	g.bossTimer -= delta
	if g.bossTimer <= 0 {
		// BOSS 超時
		g.boss = nil
		g.setState(StateNormalPlay)
		g.hub.Broadcast(protocol.MsgBossEvent, protocol.BossEventPayload{
			Event: "timeout",
		})
		g.nextBossIn = 180 + rand.Float64()*120
	}
}

func (g *Game) tickBonus(delta float64) {
	g.bonusTimer -= delta
	g.hub.Broadcast(protocol.MsgBonusEvent, protocol.BonusEventPayload{
		Event:    "tick",
		TimeLeft: g.bonusTimer,
	})
	if g.bonusTimer <= 0 {
		g.endBonus()
	}
}

// ── 目標物生成 ────────────────────────────────────────────────

func (g *Game) spawnTarget() {
	// 取平均 BetLevel 決定難度
	avgBet := 1
	if len(g.players) > 0 {
		total := 0
		for _, p := range g.players {
			total += p.BetLevel
		}
		avgBet = total / len(g.players)
	}

	def := pickTargetDef(avgBet)
	t := NewTarget(def, SpawnX, spawnY())
	g.targets[t.InstanceID] = t
	g.hub.Broadcast(protocol.MsgTargetSpawn, g.targetSpawnPayload(t))
}

func (g *Game) targetSpawnPayload(t *Target) protocol.TargetSpawnPayload {
	return protocol.TargetSpawnPayload{
		InstanceID: t.InstanceID,
		DefID:      t.Def.ID,
		Name:       t.Def.Name,
		Type:       string(t.Def.Type),
		X:          t.X,
		Y:          t.Y,
		HP:         t.HP,
		MaxHP:      t.MaxHP,
		Speed:      t.Def.Speed,
		Lifetime:   t.Def.Lifetime,
		Behavior:   string(t.Def.Behavior),
		Multiplier: t.Multiplier,
	}
}

// ── 攻擊處理 ──────────────────────────────────────────────────

func (g *Game) handleAttack(playerID string, req protocol.AttackRequest) {
	g.mu.Lock()
	defer g.mu.Unlock()

	p, ok := g.players[playerID]
	if !ok {
		return
	}

	bet := p.GetBetDef()
	if !p.SpendCoins(bet.BetCost) {
		g.hub.Send(playerID, protocol.MsgError, protocol.ErrorPayload{
			Code:    "insufficient_coins",
			Message: "Not enough coins",
		})
		return
	}

	// 找目標
	var t *Target
	if req.TargetID != "" {
		t = g.targets[req.TargetID]
		if t == nil && g.boss != nil && g.boss.InstanceID == req.TargetID {
			t = g.boss
		}
	}

	if t == nil {
		// 未命中
		g.hub.Send(playerID, protocol.MsgAttackResult, protocol.AttackResultPayload{
			TargetID:    req.TargetID,
			IsHit:       false,
			CharacterID: p.GetCharacterID(),
		})
		return
	}

	// 命中
	isKill := t.TryKill(bet.BetCost)
	result := protocol.AttackResultPayload{
		TargetID:    t.InstanceID,
		IsHit:       true,
		IsKill:      isKill,
		Damage:      bet.AttackPower,
		CharacterID: p.GetCharacterID(),
		Multiplier:  t.Multiplier,
	}

	if isKill {
		reward := int(float64(bet.BetCost) * t.Multiplier)
		result.Reward = reward
		result.LaborGain = t.Def.LaborGain

		p.AddCoins(reward)
		laborFull := p.AddLabor(t.Def.LaborGain)

		// 廣播擊破
		if t.Def.Type == data.TypeBoss {
			g.handleBossKill(playerID, p, t, reward)
		} else {
			delete(g.targets, t.InstanceID)
			g.hub.Broadcast(protocol.MsgTargetKill, protocol.TargetKillPayload{
				InstanceID: t.InstanceID,
				DefID:      t.Def.ID,
				Multiplier: t.Multiplier,
				Reward:     reward,
				LaborGain:  t.Def.LaborGain,
				KillerID:   playerID,
			})

			// 幸運特殊魚觸發
			killerName := p.ID // 使用 playerID 作為名稱（可擴展為真實名稱）
			switch {
			case isLuckyChainLightningFish(t.Def.ID):
				g.tryLuckyChainLightning(playerID, killerName)
			case isLuckyCrabTorpedoFish(t.Def.ID):
				g.tryLuckyCrabTorpedo(playerID, killerName)
			case isLuckyVortexFish(t.Def.ID):
				g.tryLuckyVortex(playerID, killerName)
			case isLuckyGoldenDragonFish(t.Def.ID):
				g.tryLuckyGoldenDragon(playerID, killerName)
			case isLuckyThunderLobsterFish(t.Def.ID):
				g.tryLuckyThunderLobster(playerID, killerName)
			case isLuckyAwakenedPhoenixFish(t.Def.ID):
				g.tryLuckyAwakenedPhoenix(playerID, killerName)
			case isLuckyShockwaveBombFish(t.Def.ID):
				g.tryLuckyShockwaveBomb(playerID, killerName)
			case isLuckyDrillTorpedoFish(t.Def.ID):
				g.tryLuckyDrillTorpedo(playerID, killerName)
			case isLuckyTimeFreezeFish(t.Def.ID):
				g.tryLuckyTimeFreeze(playerID, killerName)
			case isLuckyChainExplosionFish(t.Def.ID):
				g.tryLuckyChainExplosion(playerID, killerName)
			// DAY-295 新增
			case isLuckyChainLongKingFish(t.Def.ID):
				g.luckyChainLongKing.tryLuckyChainLongKing(g, playerID, killerName)
			case isLuckyDragonShotgunFish(t.Def.ID):
				g.luckyDragonShotgun.tryLuckyDragonShotgun(g, playerID, killerName)
			case isLuckyRocketCannonFish(t.Def.ID):
				g.luckyRocketCannon.tryLuckyRocketCannon(g, playerID, killerName)
			case isLuckyDeepWhirlpoolFish(t.Def.ID):
				g.luckyDeepWhirlpool.tryLuckyDeepWhirlpool(g, playerID, killerName)
			case isLuckyVampireMultFish(t.Def.ID):
				g.luckyVampireMult.tryLuckyVampireMult(g, playerID, killerName)
			// DAY-296 新增
			case isLuckyMirrorFish(t.Def.ID):
				g.tryLuckyMirrorFish(playerID, killerName)
			case isLuckyGoldenRainFish(t.Def.ID):
				g.tryLuckyGoldenRain(playerID, killerName)
			case isLuckyFreezeBombFish(t.Def.ID):
				g.tryLuckyFreezeBomb(playerID, killerName)
			case isLuckyThunderStormFish(t.Def.ID):
				g.tryLuckyThunderStorm(playerID, killerName)
			case isLuckyLuckyWheelFish(t.Def.ID):
				g.tryLuckyLuckyWheel(playerID, killerName)
			}
			if g.luckyChainExplosion.isChainExplosionActive(playerID) {
				g.notifyChainExplosionKill(playerID, killerName, t.X, t.Y)
			}
			// 凍結期間擊破計數
			if g.luckyTimeFreeze.isTimeFreezeActive() {
				g.luckyTimeFreeze.notifyFreezeKill(playerID)
			}
			// 吸血鬼模式：每次擊破吸收倍率
			if g.luckyVampireMult.isVampireActive(playerID) {
				g.luckyVampireMult.notifyVampireKill(g, playerID)
			}
			if g.luckyAwakenedPhoenix.isAwakenedPhoenixActive(playerID) {
				powerUpMult, powerReward, isDone := g.luckyAwakenedPhoenix.consumeAwakenedPhoenixShot(playerID, true, bet.BetCost)
				if powerReward > 0 {
					p.AddCoins(powerReward)
					reward += powerReward
					result.Reward = reward
					g.notifyAwakenedPhoenixShot(playerID, killerName, powerUpMult, powerReward, isDone)
				}
			}
		}

		// 獎勵通知（單播）
		g.hub.Send(playerID, protocol.MsgReward, protocol.RewardPayload{
			Source:     "target",
			Amount:     reward,
			Multiplier: t.Multiplier,
			NewBalance: p.Coins,
		})

		// 勞動值滿 → 觸發 Bonus
		if laborFull && g.state == StateNormalPlay {
			g.triggerBonus(playerID)
		}
	} else {
		// 未擊破，更新 HP（視覺反饋）
		t.HP = max(1, t.HP-bet.AttackPower/5)
		g.hub.Broadcast(protocol.MsgTargetUpdate, protocol.TargetUpdatePayload{
			InstanceID: t.InstanceID,
			HP:         t.HP,
			MaxHP:      t.MaxHP,
			X:          t.X,
			Y:          t.Y,
			IsFleeing:  t.IsFleeing,
		})
	}

	g.hub.Send(playerID, protocol.MsgAttackResult, result)
	g.sendPlayerUpdate(playerID)
}

// ── BOSS 系統 ─────────────────────────────────────────────────

func (g *Game) startBossWarning() {
	g.setState(StateBossWarning)
	g.bossWarningAt = time.Now()
	g.hub.Broadcast(protocol.MsgBossEvent, protocol.BossEventPayload{Event: "warning"})
	log.Printf("[Game] Boss warning!")
}

func (g *Game) spawnBoss() {
	def, ok := data.GetTarget("B001")
	if !ok {
		return
	}
	g.boss = NewTarget(def, GameWidth/2, GameHeight/2)
	g.bossTimer = 60.0
	g.setState(StateBossBattle)

	// 清除部分普通目標（BOSS 期間最多 8 個）
	count := 0
	for id := range g.targets {
		if count >= MaxBossTargets {
			delete(g.targets, id)
		}
		count++
	}

	g.hub.Broadcast(protocol.MsgBossEvent, protocol.BossEventPayload{
		Event:      "spawn",
		InstanceID: g.boss.InstanceID,
		HP:         g.boss.HP,
		MaxHP:      g.boss.MaxHP,
	})
	g.hub.Broadcast(protocol.MsgTargetSpawn, g.targetSpawnPayload(g.boss))
	log.Printf("[Game] Boss spawned!")
}

func (g *Game) handleBossKill(playerID string, p *Player, boss *Target, reward int) {
	// 依剩餘時間計算倍率
	mult := bossRewardMult(g.bossTimer)
	finalReward := int(float64(p.GetBetDef().BetCost) * mult)
	p.AddCoins(finalReward - reward) // 補差額

	g.boss = nil
	g.setState(StateNormalPlay)
	g.nextBossIn = 180 + rand.Float64()*120

	g.hub.Broadcast(protocol.MsgBossEvent, protocol.BossEventPayload{
		Event:      "kill",
		InstanceID: boss.InstanceID,
		Reward:     finalReward,
		Multiplier: mult,
	})
	log.Printf("[Game] Boss killed by %s! Reward: %d (%.0fx)", playerID, finalReward, mult)
}

func bossRewardMult(timeLeft float64) float64 {
	switch {
	case timeLeft > 50:
		return 500
	case timeLeft > 40:
		return 400
	case timeLeft > 30:
		return 300
	case timeLeft > 20:
		return 200
	case timeLeft > 10:
		return 150
	default:
		return 100
	}
}

func (g *Game) triggerBoss() {
	if g.state != StateNormalPlay {
		return
	}
	g.startBossWarning()
}

// ── Bonus 系統 ────────────────────────────────────────────────

func (g *Game) triggerBonus(playerID string) {
	if g.state != StateNormalPlay {
		return
	}
	g.mu.Lock()
	p, ok := g.players[playerID]
	if !ok {
		g.mu.Unlock()
		return
	}
	p.EntryBetCost = p.GetBetDef().BetCost
	p.ResetLabor()
	g.bonusPlayer = playerID
	g.bonusTimer = 15.0
	g.bonusScore = make(map[string]int)
	g.bonusWeedHP = make(map[string]int)
	g.setState(StateBonusGame)
	g.mu.Unlock()

	g.hub.Broadcast(protocol.MsgBonusEvent, protocol.BonusEventPayload{
		Event:    "start",
		TimeLeft: 15.0,
	})
	g.sendPlayerUpdate(playerID)
	log.Printf("[Game] Bonus triggered by %s", playerID)
}

func (g *Game) handleBonusClick(playerID string, req protocol.BonusClickRequest) {
	if g.state != StateBonusGame {
		return
	}
	g.mu.Lock()
	defer g.mu.Unlock()

	p, ok := g.players[playerID]
	if !ok {
		return
	}

	// 找 Bonus 目標物定義
	var bonusDef *data.BonusTargetDef
	for i := range data.BonusTargets {
		if data.BonusTargets[i].ID == req.TargetID {
			bonusDef = &data.BonusTargets[i]
			break
		}
	}
	if bonusDef == nil {
		return
	}

	// BG002 硬雜草需要連點 2 次
	if bonusDef.Special == "hard" {
		hp := g.bonusWeedHP[req.TargetID]
		if hp == 0 {
			hp = 2
		}
		hp--
		g.bonusWeedHP[req.TargetID] = hp
		if hp > 0 {
			return // 還沒打完
		}
	}

	score := bonusDef.Score
	g.bonusScore[playerID] += score

	g.hub.Broadcast(protocol.MsgBonusEvent, protocol.BonusEventPayload{
		Event:    "click",
		TimeLeft: g.bonusTimer,
		Score:    g.bonusScore[playerID],
	})

	_ = p // 避免 unused
}

func (g *Game) endBonus() {
	p, ok := g.players[g.bonusPlayer]
	if !ok {
		g.setState(StateNormalPlay)
		return
	}

	score := g.bonusScore[g.bonusPlayer]
	mult := 20.0 + float64(score)*0.375
	if mult < 20 {
		mult = 20
	}
	if mult > 50 {
		mult = 50
	}
	reward := int(float64(p.EntryBetCost) * mult)
	p.AddCoins(reward)

	g.setState(StateBonusResult)
	g.hub.Broadcast(protocol.MsgBonusEvent, protocol.BonusEventPayload{
		Event:      "result",
		Score:      score,
		Multiplier: mult,
		Reward:     reward,
	})
	g.sendPlayerUpdate(g.bonusPlayer)
	log.Printf("[Game] Bonus ended. Score: %d, Mult: %.1f, Reward: %d", score, mult, reward)

	// 2 秒後回到正常遊戲
	time.AfterFunc(2*time.Second, func() {
		g.mu.Lock()
		g.setState(StateNormalPlay)
		g.mu.Unlock()
	})
}

// ── AUTO 射擊 ─────────────────────────────────────────────────

func (g *Game) autoFire(p *Player) {
	if g.state != StateNormalPlay && g.state != StateBossBattle {
		return
	}
	// 找最高價值目標
	var best *Target
	bestScore := -1.0
	for _, t := range g.targets {
		score := t.Multiplier * 2.0
		if t.X < 400 {
			score += 20
		}
		if score > bestScore {
			bestScore = score
			best = t
		}
	}
	if g.boss != nil {
		best = g.boss
	}
	if best == nil {
		return
	}
	// 模擬攻擊
	g.handleAttack(p.ID, protocol.AttackRequest{
		TargetID: best.InstanceID,
		ClickX:   best.X,
		ClickY:   best.Y,
	})
}

// ── 輔助函數 ──────────────────────────────────────────────────

func (g *Game) setState(s GameState) {
	g.state = s
	g.hub.Broadcast(protocol.MsgGameState, protocol.GameStatePayload{
		State:     string(s),
		Timestamp: time.Now().UnixMilli(),
	})
}

func (g *Game) handleAutoToggle(playerID string) {
	g.mu.Lock()
	p, ok := g.players[playerID]
	if ok {
		p.IsAuto = !p.IsAuto
	}
	g.mu.Unlock()
	if ok {
		g.sendPlayerUpdate(playerID)
	}
}

func (g *Game) handleBetChange(playerID string, level int) {
	g.mu.Lock()
	p, ok := g.players[playerID]
	if ok {
		if level < 1 {
			level = 1
		}
		if level > 10 {
			level = 10
		}
		p.BetLevel = level
	}
	g.mu.Unlock()
	if ok {
		g.sendPlayerUpdate(playerID)
	}
}

func (g *Game) handleLock(playerID string, targetID string) {
	g.mu.Lock()
	p, ok := g.players[playerID]
	if ok {
		p.LockTargetID = targetID
	}
	g.mu.Unlock()
	if ok {
		g.sendPlayerUpdate(playerID)
	}
}

func (g *Game) sendPlayerUpdate(playerID string) {
	g.mu.RLock()
	p, ok := g.players[playerID]
	g.mu.RUnlock()
	if !ok {
		return
	}
	bet := p.GetBetDef()
	g.hub.Send(playerID, protocol.MsgPlayerUpdate, protocol.PlayerUpdatePayload{
		ID:              p.ID,
		Coins:           p.Coins,
		BetLevel:        p.BetLevel,
		BetCost:         bet.BetCost,
		CharacterID:     p.GetCharacterID(),
		CharacterName:   p.GetCharacterName(),
		LaborValue:      p.LaborValue,
		IsAuto:          p.IsAuto,
		LockTargetID:    p.LockTargetID,
		ProjectileSpeed: bet.ProjectileSpeed,
		FireRate:        bet.FireRate,
	})
}

func (g *Game) GetState() string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return string(g.state)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
