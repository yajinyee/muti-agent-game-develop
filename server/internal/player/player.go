// Package player 管理玩家狀態
package player

import (
	"sync"
	"time"

	"digital-twin/server/internal/data"
	"digital-twin/server/internal/game/achievement"
	"digital-twin/server/internal/game/codex"
	"digital-twin/server/internal/game/streak"
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

	// 稱號系統（DAY-068）
	Titles *achievement.TitleTracker

	// 砲台外觀系統（DAY-071）
	EquippedSkin string   // 當前裝備的外觀 ID（預設 "default"）
	OwnedSkins   []string // 已擁有的外觀 ID 列表

	// 魚類圖鑑系統（DAY-081）
	Codex *codex.Manager

	// 連擊系統（DAY-082）
	Streak *streak.Manager

	// 房間難度系統（DAY-091）
	RoomDifficulty string // 當前所在房間難度（"beginner"/"intermediate"/"advanced"/"vip"）
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
		Titles:       achievement.NewTitleTracker(),
		EquippedSkin: "default",
		OwnedSkins:   []string{"default"},
		Codex:        codex.NewManager(),
		Streak:       streak.NewManager(),
		RoomDifficulty: "beginner", // 預設初級房間
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

// AddCoins 直接增加金幣（不觸發成就，用於系統獎勵如賽季通行證）（DAY-072）
func (p *Player) AddCoins(amount int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Coins += amount
	if p.Coins > p.MaxCoins {
		p.MaxCoins = p.Coins
	}
}

// GetCoins 取得當前金幣數（thread-safe）
func (p *Player) GetCoins() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Coins
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
	title := p.Titles.GetActiveTitle()
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
		// 稱號（DAY-068）
		TitleID:    string(title.ID),
		TitleName:  title.Name,
		TitleIcon:  title.Icon,
		TitleColor: title.Color,
		// 砲台外觀（DAY-071）
		EquippedSkin: p.EquippedSkin,
		OwnedSkins:   append([]string{}, p.OwnedSkins...),
	}
}

// TryUnlockAchievement 嘗試解鎖指定成就（用於外部觸發，如 BOSS 擊殺、Bonus 觸發）
func (p *Player) TryUnlockAchievement(id achievement.AchievementID) *achievement.AchievementUnlock {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.Achievements.TryUnlock(id)
}

// OnAchievementUnlocked 成就解鎖後檢查稱號（DAY-068）
// 回傳新解鎖的稱號定義，若無新稱號則回傳 nil
func (p *Player) OnAchievementUnlocked(id achievement.AchievementID) *achievement.TitleDef {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.Titles.OnAchievementUnlocked(id, len(p.Achievements.Unlocked))
}

// SetTitle 設定顯示稱號（DAY-068）
func (p *Player) SetTitle(titleID achievement.TitleID) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.Titles.SetActiveTitle(titleID)
}

// GetAchievements 取得已解鎖成就列表（DAY-069）
func (p *Player) GetAchievements() []achievement.AchievementUnlock {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Achievements.UnlockedList()
}

// GetLoginInfo 取得登入資訊（DAY-069）
func (p *Player) GetLoginInfo() (streak int, maxStreak int) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.LoginStreak, p.MaxLoginStreak
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
	title := p.Titles.GetActiveTitle()
	return LeaderboardSnapshot{
		PlayerID:    p.ID,
		DisplayName: p.DisplayName,
		Score:       p.SessionScore,
		MaxCoins:    p.MaxCoins,
		KillCount:   p.KillCount,
		TitleID:     string(title.ID),
		TitleName:   title.Name,
		TitleIcon:   title.Icon,
		TitleColor:  title.Color,
	}
}

// LeaderboardSnapshot 排行榜快照資料
type LeaderboardSnapshot struct {
	PlayerID    string
	DisplayName string
	Score       int
	MaxCoins    int
	KillCount   int
	// 稱號（DAY-068）
	TitleID    string
	TitleName  string
	TitleIcon  string
	TitleColor string
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
	// 稱號（DAY-068）
	TitleID    string `json:"title_id"`
	TitleName  string `json:"title_name"`
	TitleIcon  string `json:"title_icon"`
	TitleColor string `json:"title_color"`
	// 砲台外觀（DAY-071）
	EquippedSkin string   `json:"equipped_skin"`
	OwnedSkins   []string `json:"owned_skins"`
}

// BuySkin 購買砲台外觀（DAY-071）
// 回傳 true=購買成功，false=金幣不足或已擁有
func (p *Player) BuySkin(skinID string, price int) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	// 已擁有
	for _, s := range p.OwnedSkins {
		if s == skinID {
			return false
		}
	}

	// 金幣不足
	if p.Coins < price {
		return false
	}

	p.Coins -= price
	p.OwnedSkins = append(p.OwnedSkins, skinID)
	return true
}

// EquipSkin 裝備砲台外觀（DAY-071）
// 回傳 true=裝備成功，false=未擁有
func (p *Player) EquipSkin(skinID string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, s := range p.OwnedSkins {
		if s == skinID {
			p.EquippedSkin = skinID
			return true
		}
	}
	return false
}

// GetSkinInfo 取得外觀資訊（DAY-071）
func (p *Player) GetSkinInfo() (equippedSkin string, ownedSkins []string) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	owned := make([]string, len(p.OwnedSkins))
	copy(owned, p.OwnedSkins)
	return p.EquippedSkin, owned
}

// GetRoomDifficulty 取得當前房間難度（DAY-091）
func (p *Player) GetRoomDifficulty() string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if p.RoomDifficulty == "" {
		return "beginner"
	}
	return p.RoomDifficulty
}

// SetRoomDifficulty 設定房間難度（DAY-091）
func (p *Player) SetRoomDifficulty(diff string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.RoomDifficulty = diff
}

// DeductCoins 扣除金幣（用於進場費等，DAY-091）
// 回傳是否成功（金幣足夠）和扣除後餘額
func (p *Player) DeductCoins(amount int) (int, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.Coins < amount {
		return p.Coins, false
	}
	p.Coins -= amount
	return p.Coins, true
}
