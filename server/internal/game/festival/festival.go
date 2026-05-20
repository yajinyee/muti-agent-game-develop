// Package festival 賽季節日活動系統（DAY-109）
// 節日主題活動：端午節/中秋節/新年/萬聖節
// 節日限定目標物、節日任務、節日 Jackpot 加成
package festival

import (
	"sync"
	"time"
)

// FestivalType 節日類型
type FestivalType string

const (
	FestivalDragonBoat  FestivalType = "dragon_boat"  // 端午節（6月）
	FestivalMidAutumn   FestivalType = "mid_autumn"   // 中秋節（9月）
	FestivalNewYear     FestivalType = "new_year"     // 農曆新年（1-2月）
	FestivalHalloween   FestivalType = "halloween"    // 萬聖節（10月）
	FestivalNone        FestivalType = "none"         // 無節日
)

// FestivalDef 節日定義
type FestivalDef struct {
	Type           FestivalType `json:"type"`
	Name           string       `json:"name"`
	Description    string       `json:"description"`
	Icon           string       `json:"icon"`
	Color          string       `json:"color"`
	BgTheme        string       `json:"bg_theme"`        // 背景主題
	// 效果加成
	JackpotMult    float64      `json:"jackpot_mult"`    // Jackpot 貢獻倍率
	RewardMult     float64      `json:"reward_mult"`     // 獎勵倍率加成
	BonusChanceAdd float64      `json:"bonus_chance_add"` // Bonus 觸發率加成
	// 限定目標物
	SpecialTargets []SpecialTarget `json:"special_targets"`
	// 限定任務
	Tasks []FestivalTask `json:"tasks"`
	// 限定稱號
	TitleID    string `json:"title_id"`
	TitleName  string `json:"title_name"`
	TitleColor string `json:"title_color"`
}

// SpecialTarget 節日限定目標物
type SpecialTarget struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Icon       string  `json:"icon"`
	Multiplier float64 `json:"multiplier"` // 獎勵倍率
	SpawnRate  float64 `json:"spawn_rate"` // 生成機率（0-1）
}

// FestivalTask 節日任務
type FestivalTask struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Target      int    `json:"target"`   // 目標數量
	RewardCoins int    `json:"reward_coins"`
}

// FestivalDefs 節日定義
var FestivalDefs = map[FestivalType]*FestivalDef{
	FestivalDragonBoat: {
		Type:           FestivalDragonBoat,
		Name:           "端午節慶典",
		Description:    "粽子飛舞！Jackpot 加倍，限定目標物大量出現！",
		Icon:           "🎋",
		Color:          "#2E8B57",
		BgTheme:        "dragon_boat",
		JackpotMult:    1.5,
		RewardMult:     1.2,
		BonusChanceAdd: 0.05,
		SpecialTargets: []SpecialTarget{
			{ID: "zongzi", Name: "粽子", Icon: "🫕", Multiplier: 8.0, SpawnRate: 0.15},
			{ID: "dragon", Name: "龍舟", Icon: "🐉", Multiplier: 25.0, SpawnRate: 0.04},
		},
		Tasks: []FestivalTask{
			{ID: "dt_kill_10", Description: "擊破 10 個粽子", Target: 10, RewardCoins: 5000},
			{ID: "dt_kill_dragon", Description: "擊破 3 條龍舟", Target: 3, RewardCoins: 15000},
			{ID: "dt_bonus_3", Description: "完成 3 次 Bonus", Target: 3, RewardCoins: 8000},
		},
		TitleID:    "festival_dragon_boat",
		TitleName:  "🎋端午勇士",
		TitleColor: "#2E8B57",
	},
	FestivalMidAutumn: {
		Type:           FestivalMidAutumn,
		Name:           "中秋節慶典",
		Description:    "月餅滿天飛！連擊加成，月兔限定目標！",
		Icon:           "🌕",
		Color:          "#FF8C00",
		BgTheme:        "mid_autumn",
		JackpotMult:    1.3,
		RewardMult:     1.3,
		BonusChanceAdd: 0.08,
		SpecialTargets: []SpecialTarget{
			{ID: "mooncake", Name: "月餅", Icon: "🥮", Multiplier: 10.0, SpawnRate: 0.12},
			{ID: "moon_rabbit", Name: "月兔", Icon: "🐰", Multiplier: 30.0, SpawnRate: 0.03},
		},
		Tasks: []FestivalTask{
			{ID: "ma_kill_15", Description: "擊破 15 個月餅", Target: 15, RewardCoins: 6000},
			{ID: "ma_kill_rabbit", Description: "擊破 5 隻月兔", Target: 5, RewardCoins: 20000},
			{ID: "ma_streak_10", Description: "達成 10 連擊", Target: 10, RewardCoins: 10000},
		},
		TitleID:    "festival_mid_autumn",
		TitleName:  "🌕中秋賞月者",
		TitleColor: "#FF8C00",
	},
	FestivalNewYear: {
		Type:           FestivalNewYear,
		Name:           "農曆新年慶典",
		Description:    "新年快樂！紅包雨降臨，Jackpot 大爆發！",
		Icon:           "🧧",
		Color:          "#FF0000",
		BgTheme:        "new_year",
		JackpotMult:    2.0,
		RewardMult:     1.5,
		BonusChanceAdd: 0.10,
		SpecialTargets: []SpecialTarget{
			{ID: "red_envelope", Name: "紅包", Icon: "🧧", Multiplier: 15.0, SpawnRate: 0.10},
			{ID: "golden_dragon", Name: "金龍", Icon: "🐲", Multiplier: 50.0, SpawnRate: 0.02},
		},
		Tasks: []FestivalTask{
			{ID: "ny_kill_20", Description: "擊破 20 個紅包", Target: 20, RewardCoins: 8888},
			{ID: "ny_kill_dragon", Description: "擊破 2 條金龍", Target: 2, RewardCoins: 28888},
			{ID: "ny_jackpot", Description: "觸發 1 次 Jackpot", Target: 1, RewardCoins: 18888},
		},
		TitleID:    "festival_new_year",
		TitleName:  "🧧新年財神",
		TitleColor: "#FF0000",
	},
	FestivalHalloween: {
		Type:           FestivalHalloween,
		Name:           "萬聖節慶典",
		Description:    "不給糖就搗蛋！南瓜怪物大出沒！",
		Icon:           "🎃",
		Color:          "#FF6600",
		BgTheme:        "halloween",
		JackpotMult:    1.2,
		RewardMult:     1.4,
		BonusChanceAdd: 0.06,
		SpecialTargets: []SpecialTarget{
			{ID: "pumpkin", Name: "南瓜怪", Icon: "🎃", Multiplier: 12.0, SpawnRate: 0.12},
			{ID: "ghost", Name: "幽靈", Icon: "👻", Multiplier: 20.0, SpawnRate: 0.06},
		},
		Tasks: []FestivalTask{
			{ID: "hw_kill_12", Description: "擊破 12 個南瓜怪", Target: 12, RewardCoins: 5500},
			{ID: "hw_kill_ghost", Description: "擊破 8 隻幽靈", Target: 8, RewardCoins: 12000},
			{ID: "hw_chain_5", Description: "完成 5 次連鎖爆炸", Target: 5, RewardCoins: 9000},
		},
		TitleID:    "festival_halloween",
		TitleName:  "🎃萬聖獵人",
		TitleColor: "#FF6600",
	},
}

// PlayerFestivalProgress 玩家節日任務進度
type PlayerFestivalProgress struct {
	FestivalType FestivalType       `json:"festival_type"`
	TaskProgress map[string]int     `json:"task_progress"` // taskID -> 當前進度
	TaskDone     map[string]bool    `json:"task_done"`     // taskID -> 是否已領取
	TitleClaimed bool               `json:"title_claimed"`
}

// Manager 節日活動管理器
type Manager struct {
	mu       sync.RWMutex
	current  FestivalType
	startAt  time.Time
	endAt    time.Time
	// 玩家進度：playerID -> progress
	progress map[string]*PlayerFestivalProgress
}

// New 建立節日管理器，根據當前日期自動判斷節日
func New() *Manager {
	m := &Manager{
		progress: make(map[string]*PlayerFestivalProgress),
	}
	m.detectFestival()
	return m
}

// detectFestival 根據當前月份判斷節日
func (m *Manager) detectFestival() {
	now := time.Now()
	month := now.Month()
	day := now.Day()

	var ft FestivalType
	var start, end time.Time
	year := now.Year()

	switch {
	case month == 1 || (month == 2 && day <= 15):
		// 農曆新年：1月1日 - 2月15日
		ft = FestivalNewYear
		start = time.Date(year, 1, 1, 0, 0, 0, 0, time.Local)
		end = time.Date(year, 2, 15, 23, 59, 59, 0, time.Local)
	case month == 6 && day >= 1 && day <= 14:
		// 端午節：6月1日 - 6月14日
		ft = FestivalDragonBoat
		start = time.Date(year, 6, 1, 0, 0, 0, 0, time.Local)
		end = time.Date(year, 6, 14, 23, 59, 59, 0, time.Local)
	case month == 9 && day >= 10 && day <= 30:
		// 中秋節：9月10日 - 9月30日
		ft = FestivalMidAutumn
		start = time.Date(year, 9, 10, 0, 0, 0, 0, time.Local)
		end = time.Date(year, 9, 30, 23, 59, 59, 0, time.Local)
	case month == 10 && day >= 20:
		// 萬聖節：10月20日 - 10月31日
		ft = FestivalHalloween
		start = time.Date(year, 10, 20, 0, 0, 0, 0, time.Local)
		end = time.Date(year, 10, 31, 23, 59, 59, 0, time.Local)
	default:
		ft = FestivalNone
		start = now
		end = now
	}

	m.current = ft
	m.startAt = start
	m.endAt = end
}

// IsActive 是否有節日活動進行中
func (m *Manager) IsActive() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.current == FestivalNone {
		return false
	}
	now := time.Now()
	return !now.Before(m.startAt) && now.Before(m.endAt)
}

// GetCurrent 取得當前節日類型
func (m *Manager) GetCurrent() FestivalType {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.current
}

// GetDef 取得當前節日定義
func (m *Manager) GetDef() *FestivalDef {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.current == FestivalNone {
		return nil
	}
	def, ok := FestivalDefs[m.current]
	if !ok {
		return nil
	}
	return def
}

// GetJackpotMult 取得 Jackpot 倍率加成
func (m *Manager) GetJackpotMult() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.current == FestivalNone {
		return 1.0
	}
	def, ok := FestivalDefs[m.current]
	if !ok {
		return 1.0
	}
	return def.JackpotMult
}

// GetRewardMult 取得獎勵倍率加成
func (m *Manager) GetRewardMult() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.current == FestivalNone {
		return 1.0
	}
	def, ok := FestivalDefs[m.current]
	if !ok {
		return 1.0
	}
	return def.RewardMult
}

// GetBonusChanceAdd 取得 Bonus 觸發率加成
func (m *Manager) GetBonusChanceAdd() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.current == FestivalNone {
		return 0.0
	}
	def, ok := FestivalDefs[m.current]
	if !ok {
		return 0.0
	}
	return def.BonusChanceAdd
}

// GetSpecialTargets 取得節日限定目標物列表
func (m *Manager) GetSpecialTargets() []SpecialTarget {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.current == FestivalNone {
		return nil
	}
	def, ok := FestivalDefs[m.current]
	if !ok {
		return nil
	}
	return def.SpecialTargets
}

// GetOrCreateProgress 取得或建立玩家節日進度
func (m *Manager) GetOrCreateProgress(playerID string) *PlayerFestivalProgress {
	m.mu.Lock()
	defer m.mu.Unlock()
	p, ok := m.progress[playerID]
	if !ok || p.FestivalType != m.current {
		// 新節日或新玩家，重置進度
		p = &PlayerFestivalProgress{
			FestivalType: m.current,
			TaskProgress: make(map[string]int),
			TaskDone:     make(map[string]bool),
		}
		m.progress[playerID] = p
	}
	return p
}

// RecordKill 記錄擊破節日目標物，回傳是否有任務進度更新
func (m *Manager) RecordKill(playerID string, targetID string) (taskUpdated bool, completedTaskID string) {
	if !m.IsActive() {
		return false, ""
	}
	def := m.GetDef()
	if def == nil {
		return false, ""
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	p := m.progress[playerID]
	if p == nil || p.FestivalType != m.current {
		p = &PlayerFestivalProgress{
			FestivalType: m.current,
			TaskProgress: make(map[string]int),
			TaskDone:     make(map[string]bool),
		}
		m.progress[playerID] = p
	}

	// 找到對應的任務（擊破類型）
	for _, task := range def.Tasks {
		if p.TaskDone[task.ID] {
			continue
		}
		// 判斷任務類型：kill_N 類型
		if isKillTask(task.ID, targetID) {
			p.TaskProgress[task.ID]++
			if p.TaskProgress[task.ID] >= task.Target {
				return true, task.ID
			}
			return true, ""
		}
	}
	return false, ""
}

// isKillTask 判斷任務是否為擊破特定目標的任務
func isKillTask(taskID, targetID string) bool {
	// 任務 ID 格式：{prefix}_kill_{target} 或 {prefix}_kill_{n}
	// 簡單判斷：任務 ID 包含 targetID 或是通用擊破任務
	if len(taskID) == 0 || len(targetID) == 0 {
		return false
	}
	// 通用擊破任務（kill_N 格式，任何節日目標都算）
	for _, st := range []string{"zongzi", "dragon", "mooncake", "moon_rabbit",
		"red_envelope", "golden_dragon", "pumpkin", "ghost"} {
		if targetID == st {
			// 檢查任務 ID 是否包含目標 ID
			if contains(taskID, st) {
				return true
			}
			// 通用擊破任務（如 dt_kill_10）
			if containsKillN(taskID) {
				return true
			}
		}
	}
	return false
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > len(sub) && (s[:len(sub)] == sub || s[len(s)-len(sub):] == sub || containsMiddle(s, sub)))
}

func containsMiddle(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func containsKillN(taskID string) bool {
	// 格式：xx_kill_N（N 是數字）
	for i := 0; i < len(taskID)-6; i++ {
		if taskID[i:i+5] == "kill_" {
			rest := taskID[i+5:]
			allDigit := len(rest) > 0
			for _, c := range rest {
				if c < '0' || c > '9' {
					allDigit = false
					break
				}
			}
			if allDigit {
				return true
			}
		}
	}
	return false
}

// RecordBonus 記錄完成 Bonus，回傳是否有任務進度更新
func (m *Manager) RecordBonus(playerID string) (taskUpdated bool, completedTaskID string) {
	if !m.IsActive() {
		return false, ""
	}
	def := m.GetDef()
	if def == nil {
		return false, ""
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	p := m.progress[playerID]
	if p == nil {
		p = &PlayerFestivalProgress{
			FestivalType: m.current,
			TaskProgress: make(map[string]int),
			TaskDone:     make(map[string]bool),
		}
		m.progress[playerID] = p
	}

	for _, task := range def.Tasks {
		if p.TaskDone[task.ID] {
			continue
		}
		if containsMiddle(task.ID, "bonus") {
			p.TaskProgress[task.ID]++
			if p.TaskProgress[task.ID] >= task.Target {
				return true, task.ID
			}
			return true, ""
		}
	}
	return false, ""
}

// RecordStreak 記錄連擊達成
func (m *Manager) RecordStreak(playerID string, streak int) (taskUpdated bool, completedTaskID string) {
	if !m.IsActive() {
		return false, ""
	}
	def := m.GetDef()
	if def == nil {
		return false, ""
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	p := m.progress[playerID]
	if p == nil {
		p = &PlayerFestivalProgress{
			FestivalType: m.current,
			TaskProgress: make(map[string]int),
			TaskDone:     make(map[string]bool),
		}
		m.progress[playerID] = p
	}

	for _, task := range def.Tasks {
		if p.TaskDone[task.ID] {
			continue
		}
		if containsMiddle(task.ID, "streak") {
			if streak >= task.Target && p.TaskProgress[task.ID] < task.Target {
				p.TaskProgress[task.ID] = task.Target
				return true, task.ID
			}
		}
	}
	return false, ""
}

// RecordJackpot 記錄 Jackpot 觸發
func (m *Manager) RecordJackpot(playerID string) (taskUpdated bool, completedTaskID string) {
	if !m.IsActive() {
		return false, ""
	}
	def := m.GetDef()
	if def == nil {
		return false, ""
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	p := m.progress[playerID]
	if p == nil {
		p = &PlayerFestivalProgress{
			FestivalType: m.current,
			TaskProgress: make(map[string]int),
			TaskDone:     make(map[string]bool),
		}
		m.progress[playerID] = p
	}

	for _, task := range def.Tasks {
		if p.TaskDone[task.ID] {
			continue
		}
		if containsMiddle(task.ID, "jackpot") {
			p.TaskProgress[task.ID]++
			if p.TaskProgress[task.ID] >= task.Target {
				return true, task.ID
			}
			return true, ""
		}
	}
	return false, ""
}

// RecordChain 記錄連鎖爆炸
func (m *Manager) RecordChain(playerID string) (taskUpdated bool, completedTaskID string) {
	if !m.IsActive() {
		return false, ""
	}
	def := m.GetDef()
	if def == nil {
		return false, ""
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	p := m.progress[playerID]
	if p == nil {
		p = &PlayerFestivalProgress{
			FestivalType: m.current,
			TaskProgress: make(map[string]int),
			TaskDone:     make(map[string]bool),
		}
		m.progress[playerID] = p
	}

	for _, task := range def.Tasks {
		if p.TaskDone[task.ID] {
			continue
		}
		if containsMiddle(task.ID, "chain") {
			p.TaskProgress[task.ID]++
			if p.TaskProgress[task.ID] >= task.Target {
				return true, task.ID
			}
			return true, ""
		}
	}
	return false, ""
}

// ClaimTaskReward 領取任務獎勵，回傳獎勵金幣數（0=失敗）
func (m *Manager) ClaimTaskReward(playerID, taskID string) int {
	if !m.IsActive() {
		return 0
	}
	def := m.GetDef()
	if def == nil {
		return 0
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	p := m.progress[playerID]
	if p == nil {
		return 0
	}
	if p.TaskDone[taskID] {
		return 0 // 已領取
	}

	for _, task := range def.Tasks {
		if task.ID == taskID {
			if p.TaskProgress[taskID] >= task.Target {
				p.TaskDone[taskID] = true
				return task.RewardCoins
			}
			return 0
		}
	}
	return 0
}

// IsAllTasksDone 是否所有任務都已完成
func (m *Manager) IsAllTasksDone(playerID string) bool {
	if !m.IsActive() {
		return false
	}
	def := m.GetDef()
	if def == nil {
		return false
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	p := m.progress[playerID]
	if p == nil {
		return false
	}
	for _, task := range def.Tasks {
		if !p.TaskDone[task.ID] {
			return false
		}
	}
	return true
}

// ClaimTitle 領取節日稱號（完成所有任務後）
func (m *Manager) ClaimTitle(playerID string) (titleID, titleName, titleColor string, ok bool) {
	if !m.IsActive() {
		return "", "", "", false
	}
	def := m.GetDef()
	if def == nil {
		return "", "", "", false
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	p := m.progress[playerID]
	if p == nil || p.TitleClaimed {
		return "", "", "", false
	}

	// 確認所有任務都完成
	for _, task := range def.Tasks {
		if !p.TaskDone[task.ID] {
			return "", "", "", false
		}
	}

	p.TitleClaimed = true
	return def.TitleID, def.TitleName, def.TitleColor, true
}

// GetSnapshot 取得節日快照（用於 WebSocket 廣播）
func (m *Manager) GetSnapshot(playerID string) FestivalSnapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.current == FestivalNone {
		return FestivalSnapshot{Type: string(FestivalNone), IsActive: false}
	}

	def, ok := FestivalDefs[m.current]
	if !ok {
		return FestivalSnapshot{Type: string(FestivalNone), IsActive: false}
	}

	now := time.Now()
	isActive := !now.Before(m.startAt) && now.Before(m.endAt)

	// 玩家任務進度
	var taskSnapshots []FestivalTaskSnapshot
	p := m.progress[playerID]
	for _, task := range def.Tasks {
		progress := 0
		done := false
		if p != nil {
			progress = p.TaskProgress[task.ID]
			done = p.TaskDone[task.ID]
		}
		taskSnapshots = append(taskSnapshots, FestivalTaskSnapshot{
			ID:          task.ID,
			Description: task.Description,
			Target:      task.Target,
			Progress:    progress,
			Done:        done,
			RewardCoins: task.RewardCoins,
		})
	}

	titleClaimed := false
	if p != nil {
		titleClaimed = p.TitleClaimed
	}

	return FestivalSnapshot{
		Type:           string(m.current),
		Name:           def.Name,
		Description:    def.Description,
		Icon:           def.Icon,
		Color:          def.Color,
		BgTheme:        def.BgTheme,
		IsActive:       isActive,
		EndAt:          m.endAt.UnixMilli(),
		TimeLeft:       time.Until(m.endAt).Seconds(),
		JackpotMult:    def.JackpotMult,
		RewardMult:     def.RewardMult,
		BonusChanceAdd: def.BonusChanceAdd,
		SpecialTargets: def.SpecialTargets,
		Tasks:          taskSnapshots,
		TitleID:        def.TitleID,
		TitleName:      def.TitleName,
		TitleColor:     def.TitleColor,
		TitleClaimed:   titleClaimed,
	}
}

// FestivalSnapshot 節日快照（用於 WebSocket 廣播）
type FestivalSnapshot struct {
	Type           string                 `json:"type"`
	Name           string                 `json:"name"`
	Description    string                 `json:"description"`
	Icon           string                 `json:"icon"`
	Color          string                 `json:"color"`
	BgTheme        string                 `json:"bg_theme"`
	IsActive       bool                   `json:"is_active"`
	EndAt          int64                  `json:"end_at"`
	TimeLeft       float64                `json:"time_left"`
	JackpotMult    float64                `json:"jackpot_mult"`
	RewardMult     float64                `json:"reward_mult"`
	BonusChanceAdd float64                `json:"bonus_chance_add"`
	SpecialTargets []SpecialTarget        `json:"special_targets"`
	Tasks          []FestivalTaskSnapshot `json:"tasks"`
	TitleID        string                 `json:"title_id"`
	TitleName      string                 `json:"title_name"`
	TitleColor     string                 `json:"title_color"`
	TitleClaimed   bool                   `json:"title_claimed"`
}

// FestivalTaskSnapshot 任務快照
type FestivalTaskSnapshot struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Target      int    `json:"target"`
	Progress    int    `json:"progress"`
	Done        bool   `json:"done"`
	RewardCoins int    `json:"reward_coins"`
}

// RemovePlayer 清理玩家資料（離線時）
func (m *Manager) RemovePlayer(playerID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.progress, playerID)
}

// RandFloat64 提供給外部使用的隨機數（避免 import cycle）
func RandFloat64() float64 {
	// 使用 math/rand 全局隨機數
	return randFloat64()
}
