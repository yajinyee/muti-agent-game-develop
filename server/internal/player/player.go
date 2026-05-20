// Package player 管理玩家狀態
package player

import (
	"sync"
	"time"

	"digital-twin/server/internal/data"
	"digital-twin/server/internal/game/achievement"
)

// Player 玩家狀態
type Player struct {
	mu sync.RWMutex

	ID          string
	Coins       int
	BetLevel    int
	LaborValue  int // 勞動值 0-100
	IsAuto      bool
	LockTargetID string // 鎖定目標 InstanceID
	SessionStart time.Time
	LastAttackAt time.Time

	// 統計
	TotalBet    int
	TotalReward int
	AttackCount int
	KillCount   int

	// 排行榜
	SessionScore int // 本局累積獎勵（用於排行榜）
	MaxCoins     int // 歷史最高金幣
	DisplayName  string // 顯示名稱（預設為 ID 前 8 碼）

	// 連擊系統（DAY-022）
	ComboCount   int       // 當前連擊數
	LastKillAt   time.Time // 上次擊破時間（用於判斷連擊是否中斷）

	// 成就系統
	Achievements *achievement.Tracker

	// 每日登入獎勵（DAY-065）
	LoginStreak    int    // 連續登入天數
	MaxLoginStreak int    // 歷史最高連續天數
	LastLoginDate  string // 最後登入日期（UTC+8，格式 "2006-01-02"）

	// 武器升級系統（DAY-067）
	WeaponLevel int // 武器等級 1/2/3（預設 1）
}

// NewPlayer 建立新玩家
func NewPlayer(id string, initialCoins int) *Player {
	// 顯示名稱：取 ID 前 8 碼，若 ID 太短就全用
	displayName := id
	if len(id) > 8 {
		displayName = id[:8]
	}
	return &Player{
		ID:           id,
		Coins:        initialCoins,
		BetLevel:     1,
		WeaponLevel:  1,
		LaborValue:   0,
		IsAuto:       false,
		SessionStart: time.Now(),
		MaxCoins:     initialCoins,
		DisplayName:  displayName,
		Achievements: achievement.NewTracker(),
	}
}

// GetBetDef 取得目前投注定義
func (p *Player) GetBetDef() *data.BetDef {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return data.GetBetDef(p.BetLevel)
}

// GetCharacter 取得目前角色
func (p *Player) GetCharacter() *data.CharacterDef {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return data.GetCharacterByBetLevel(p.BetLevel)
}

// CanAttack 是否可以攻擊（金幣足夠）
func (p *Player) CanAttack() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	bet := data.GetBetDef(p.BetLevel)
	return p.Coins >= bet.BetCost
}

// DeductBet 扣除投注金額，回傳是否成功
func (p *Player) DeductBet() (int, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	bet := data.GetBetDef(p.BetLevel)
	weapon := data.GetWeaponDef(p.WeaponLevel)
	totalCost := bet.BetCost + weapon.ExtraCost
	if p.Coins < totalCost {
		return 0, false
	}
	p.Coins -= totalCost
	p.TotalBet += totalCost
	p.AttackCount++
	p.LastAttackAt = time.Now()
	return bet.BetCost, true
}

// GetWeaponPowerMod 取得武器攻擊力加成係數（thread-safe）
func (p *Player) GetWeaponPowerMod() float64 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return data.GetWeaponDef(p.WeaponLevel).PowerMod
}

// UpgradeWeapon 升級武器，回傳是否成功（DAY-067）
// 武器升級不需要金幣，只是切換等級（費用在每次攻擊時扣除）
func (p *Player) UpgradeWeapon(level int) bool {
	if level < 1 || level > 3 {
		return false
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.WeaponLevel = level
	return true
}

// AddKill 增加擊破計數，回傳解鎖的成就（可能為 nil）
func (p *Player) AddKill() []*achievement.AchievementUnlock {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.KillCount++
	count := p.KillCount

	var unlocks []*achievement.AchievementUnlock
	// 首殺
	if count == 1 {
		if u := p.Achievements.TryUnlock(achievement.AchFirstKill); u != nil {
			unlocks = append(unlocks, u)
		}
	}
	// 累計擊破里程碑
	milestones := map[int]achievement.AchievementID{
		5:   achievement.AchKill5,
		20:  achievement.AchKill20,
		50:  achievement.AchKill50,
		100: achievement.AchKill100,
	}
	if id, ok := milestones[count]; ok {
		if u := p.Achievements.TryUnlock(id); u != nil {
			unlocks = append(unlocks, u)
		}
	}
	return unlocks
}

// AddReward 增加獎勵，回傳解鎖的成就（可能為空）
func (p *Player) AddReward(amount int) []*achievement.AchievementUnlock {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Coins += amount
	p.TotalReward += amount
	p.SessionScore += amount
	// 更新歷史最高金幣
	if p.Coins > p.MaxCoins {
		p.MaxCoins = p.Coins
	}

	var unlocks []*achievement.AchievementUnlock
	// 金幣里程碑
	if p.Coins >= 100000 {
		if u := p.Achievements.TryUnlock(achievement.AchCoins100k); u != nil {
			unlocks = append(unlocks, u)
		}
	} else if p.Coins >= 50000 {
		if u := p.Achievements.TryUnlock(achievement.AchCoins50k); u != nil {
			unlocks = append(unlocks, u)
		}
	}
	return unlocks
}

// AddLaborValue 增加勞動值，回傳是否觸發 Bonus
func (p *Player) AddLaborValue(amount int) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.LaborValue += amount
	if p.LaborValue >= data.LaborValueMax {
		p.LaborValue = data.LaborValueMax
		return true
	}
	return false
}

// ResetLaborValue 重置勞動值（Bonus 觸發後）
func (p *Player) ResetLaborValue() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.LaborValue = 0
}

// SetBetLevel 切換投注等級
func (p *Player) SetBetLevel(level int) bool {
	if level < 1 || level > 10 {
		return false
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.BetLevel = level
	return true
}

// SetLock 設定鎖定目標
func (p *Player) SetLock(targetID string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.LockTargetID = targetID
}

// SetAuto 設定自動攻擊
func (p *Player) SetAuto(enabled bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.IsAuto = enabled
}

// SetDisplayName 設定顯示名稱（DAY-021）
func (p *Player) SetDisplayName(name string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.DisplayName = name
}

// AddKillCombo 更新連擊計數，回傳當前連擊數和勞動值加成係數（DAY-022）
// 連擊判定：2 秒內連續擊破
// 加成：×2=+10%, ×3=+20%, ×4+=+30%
func (p *Player) AddKillCombo() (comboCount int, laborBonus float64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	now := time.Now()
	const comboWindow = 2.0 // 秒

	if !p.LastKillAt.IsZero() && now.Sub(p.LastKillAt).Seconds() <= comboWindow {
		p.ComboCount++
	} else {
		p.ComboCount = 1 // 重置連擊
	}
	p.LastKillAt = now

	// 計算勞動值加成
	switch {
	case p.ComboCount >= 4:
		laborBonus = 0.30
	case p.ComboCount == 3:
		laborBonus = 0.20
	case p.ComboCount == 2:
		laborBonus = 0.10
	default:
		laborBonus = 0.0
	}

	return p.ComboCount, laborBonus
}

// Snapshot 取得玩家狀態快照（用於傳送給 Client）
func (p *Player) Snapshot() PlayerSnapshot {
	p.mu.RLock()
	defer p.mu.RUnlock()
	char := data.GetCharacterByBetLevel(p.BetLevel)
	bet := data.GetBetDef(p.BetLevel)
	weapon := data.GetWeaponDef(p.WeaponLevel)
	return PlayerSnapshot{
		ID:              p.ID,
		Coins:           p.Coins,
		BetLevel:        p.BetLevel,
		BetCost:         bet.BetCost,
		CharacterID:     char.ID,
		CharacterName:   char.Name,
		LaborValue:      p.LaborValue,
		IsAuto:          p.IsAuto,
		LockTargetID:    p.LockTargetID,
		ProjectileSpeed: bet.ProjectileSpeed,
		FireRate:        bet.FireRate,
		// Session 統計（DAY-046）
		SessionScore:    p.SessionScore,
		KillCount:       p.KillCount,
		DisplayName:     p.DisplayName,
		// 武器升級（DAY-067）
		WeaponLevel:     p.WeaponLevel,
		WeaponName:      weapon.Name,
		WeaponIcon:      weapon.Icon,
		WeaponColor:     weapon.Color,
		WeaponExtraCost: weapon.ExtraCost,
	}
}

// TryUnlockAchievement 嘗試解鎖指定成就（用於外部觸發，如 BOSS 擊殺、Bonus 觸發）
func (p *Player) TryUnlockAchievement(id achievement.AchievementID) *achievement.AchievementUnlock {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.Achievements.TryUnlock(id)
}

// TryUnlockBigWin 嘗試解鎖大獎成就（依倍率判斷）
func (p *Player) TryUnlockBigWin(multiplier float64) []*achievement.AchievementUnlock {
	p.mu.Lock()
	defer p.mu.Unlock()
	var unlocks []*achievement.AchievementUnlock
	if multiplier >= 50 {
		if u := p.Achievements.TryUnlock(achievement.AchMegaWin); u != nil {
			unlocks = append(unlocks, u)
		}
	}
	if multiplier >= 20 {
		if u := p.Achievements.TryUnlock(achievement.AchBigWin); u != nil {
			unlocks = append(unlocks, u)
		}
	}
	return unlocks
}

// LeaderboardSnapshot 排行榜快照
func (p *Player) LeaderboardSnapshot() LeaderboardSnapshot {	p.mu.RLock()
	defer p.mu.RUnlock()
	return LeaderboardSnapshot{
		PlayerID:    p.ID,
		DisplayName: p.DisplayName,
		Score:       p.SessionScore,
		MaxCoins:    p.MaxCoins,
		KillCount:   p.KillCount,
	}
}

// LeaderboardSnapshot 排行榜快照資料
type LeaderboardSnapshot struct {
	PlayerID    string
	DisplayName string
	Score       int
	MaxCoins    int
	KillCount   int
}

// PlayerSnapshot 玩家狀態快照
type PlayerSnapshot struct {
	ID              string  `json:"id"`
	Coins           int     `json:"coins"`
	BetLevel        int     `json:"bet_level"`
	BetCost         int     `json:"bet_cost"`
	CharacterID     string  `json:"character_id"`
	CharacterName   string  `json:"character_name"`
	LaborValue      int     `json:"labor_value"`
	IsAuto          bool    `json:"is_auto"`
	LockTargetID    string  `json:"lock_target_id"`
	ProjectileSpeed float64 `json:"projectile_speed"`
	FireRate        float64 `json:"fire_rate"`
	// Session 統計（DAY-046，供 Client 端 Session Stats 面板顯示）
	SessionScore    int     `json:"session_score"`
	KillCount       int     `json:"kill_count"`
	DisplayName     string  `json:"display_name"`
	// 武器升級（DAY-067）
	WeaponLevel     int     `json:"weapon_level"`
	WeaponName      string  `json:"weapon_name"`
	WeaponIcon      string  `json:"weapon_icon"`
	WeaponColor     string  `json:"weapon_color"`
	WeaponExtraCost int     `json:"weapon_extra_cost"`
}
