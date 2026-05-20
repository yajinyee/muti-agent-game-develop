// Package guild 公會系統（DAY-074）
// 玩家可以建立/加入公會，共同完成公會任務，獲得公會獎勵
// 公會最多 20 人，有會長/副會長/成員三個職位
package guild

import (
	"fmt"
	"sync"
	"time"
)

// MemberRole 公會成員職位
type MemberRole string

const (
	RoleLeader    MemberRole = "leader"    // 會長
	RoleOfficer   MemberRole = "officer"   // 副會長
	RoleMember    MemberRole = "member"    // 成員
)

// MaxGuildMembers 公會最大成員數
const MaxGuildMembers = 20

// GuildTaskType 公會任務類型
type GuildTaskType string

const (
	TaskKillTargets GuildTaskType = "kill_targets" // 擊破目標數
	TaskKillBoss    GuildTaskType = "kill_boss"    // 擊殺 BOSS 次數
	TaskEarnCoins   GuildTaskType = "earn_coins"   // 賺取金幣
)

// GuildTask 公會任務
type GuildTask struct {
	ID          string        `json:"id"`
	Type        GuildTaskType `json:"type"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Icon        string        `json:"icon"`
	Target      int           `json:"target"`
	Current     int           `json:"current"`
	Reward      int           `json:"reward"` // 每人獎勵金幣
	Completed   bool          `json:"completed"`
	ResetAt     time.Time     `json:"reset_at"` // 每日 UTC+8 00:00 重置
}

// GuildMember 公會成員
type GuildMember struct {
	PlayerID    string     `json:"player_id"`
	DisplayName string     `json:"display_name"`
	Role        MemberRole `json:"role"`
	JoinedAt    time.Time  `json:"joined_at"`
	IsOnline    bool       `json:"is_online"`
	Contribution int       `json:"contribution"` // 本週貢獻積分
}

// Guild 公會
type Guild struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Icon        string                 `json:"icon"`
	Level       int                    `json:"level"`       // 公會等級（1-5）
	Exp         int                    `json:"exp"`         // 公會經驗值
	Members     map[string]*GuildMember `json:"members"`    // playerID → member
	Tasks       []*GuildTask           `json:"tasks"`
	CreatedAt   time.Time              `json:"created_at"`
	TotalKills  int                    `json:"total_kills"`  // 累積擊破數
	TotalCoins  int                    `json:"total_coins"`  // 累積賺取金幣
}

// Manager 公會系統管理器
type Manager struct {
	mu           sync.RWMutex
	guilds       map[string]*Guild  // guildID → guild
	playerGuild  map[string]string  // playerID → guildID
	nextGuildID  int
}

// New 建立新的公會管理器
func New() *Manager {
	return &Manager{
		guilds:      make(map[string]*Guild),
		playerGuild: make(map[string]string),
		nextGuildID: 1,
	}
}

// CreateGuild 建立公會
// 回傳 guildID, error
func (m *Manager) CreateGuild(leaderID, leaderName, guildName, description string) (string, error) {
	if guildName == "" {
		return "", fmt.Errorf("公會名稱不能為空")
	}
	if len(guildName) > 20 {
		return "", fmt.Errorf("公會名稱不能超過 20 字")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// 玩家已在公會
	if _, ok := m.playerGuild[leaderID]; ok {
		return "", fmt.Errorf("你已經在公會中，請先退出")
	}

	// 公會名稱重複
	for _, g := range m.guilds {
		if g.Name == guildName {
			return "", fmt.Errorf("公會名稱已被使用")
		}
	}

	guildID := fmt.Sprintf("guild-%04d", m.nextGuildID)
	m.nextGuildID++

	now := time.Now()
	guild := &Guild{
		ID:          guildID,
		Name:        guildName,
		Description: description,
		Icon:        "⚔️",
		Level:       1,
		Exp:         0,
		Members:     make(map[string]*GuildMember),
		Tasks:       m.generateDailyTasks(),
		CreatedAt:   now,
	}

	// 加入會長
	guild.Members[leaderID] = &GuildMember{
		PlayerID:    leaderID,
		DisplayName: leaderName,
		Role:        RoleLeader,
		JoinedAt:    now,
		IsOnline:    true,
	}

	m.guilds[guildID] = guild
	m.playerGuild[leaderID] = guildID
	return guildID, nil
}

// JoinGuild 加入公會
func (m *Manager) JoinGuild(playerID, playerName, guildID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 玩家已在公會
	if _, ok := m.playerGuild[playerID]; ok {
		return fmt.Errorf("你已經在公會中，請先退出")
	}

	guild, ok := m.guilds[guildID]
	if !ok {
		return fmt.Errorf("公會不存在")
	}

	if len(guild.Members) >= MaxGuildMembers {
		return fmt.Errorf("公會已滿（最多 %d 人）", MaxGuildMembers)
	}

	guild.Members[playerID] = &GuildMember{
		PlayerID:    playerID,
		DisplayName: playerName,
		Role:        RoleMember,
		JoinedAt:    time.Now(),
		IsOnline:    true,
	}
	m.playerGuild[playerID] = guildID
	return nil
}

// LeaveGuild 退出公會
// 如果是會長且有其他成員，自動轉讓給副會長或最早加入的成員
func (m *Manager) LeaveGuild(playerID string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	guildID, ok := m.playerGuild[playerID]
	if !ok {
		return "", fmt.Errorf("你不在任何公會中")
	}

	guild := m.guilds[guildID]
	member := guild.Members[playerID]

	// 如果是會長
	if member.Role == RoleLeader {
		if len(guild.Members) == 1 {
			// 最後一個成員，解散公會
			delete(m.guilds, guildID)
			delete(m.playerGuild, playerID)
			return guildID, nil
		}
		// 轉讓給副會長或最早加入的成員
		newLeaderID := m.findNewLeader(guild, playerID)
		if newLeaderID != "" {
			guild.Members[newLeaderID].Role = RoleLeader
		}
	}

	delete(guild.Members, playerID)
	delete(m.playerGuild, playerID)
	return guildID, nil
}

// KickMember 踢出成員（只有會長/副會長可以踢普通成員，會長可以踢副會長）
func (m *Manager) KickMember(operatorID, targetID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	guildID, ok := m.playerGuild[operatorID]
	if !ok {
		return fmt.Errorf("你不在任何公會中")
	}

	targetGuildID, ok := m.playerGuild[targetID]
	if !ok || targetGuildID != guildID {
		return fmt.Errorf("目標玩家不在你的公會中")
	}

	guild := m.guilds[guildID]
	operator := guild.Members[operatorID]
	target := guild.Members[targetID]

	// 不能踢自己
	if operatorID == targetID {
		return fmt.Errorf("不能踢出自己")
	}

	// 不能踢會長
	if target.Role == RoleLeader {
		return fmt.Errorf("不能踢出會長")
	}

	// 副會長只能踢普通成員
	if operator.Role == RoleOfficer && target.Role == RoleOfficer {
		return fmt.Errorf("副會長不能踢出其他副會長")
	}

	// 普通成員不能踢人
	if operator.Role == RoleMember {
		return fmt.Errorf("你沒有踢人的權限")
	}

	delete(guild.Members, targetID)
	delete(m.playerGuild, targetID)
	return nil
}

// PromoteMember 升職成員（會長可以升副會長，副會長可以升普通成員為副會長）
func (m *Manager) PromoteMember(operatorID, targetID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	guildID, ok := m.playerGuild[operatorID]
	if !ok {
		return fmt.Errorf("你不在任何公會中")
	}

	targetGuildID, ok := m.playerGuild[targetID]
	if !ok || targetGuildID != guildID {
		return fmt.Errorf("目標玩家不在你的公會中")
	}

	guild := m.guilds[guildID]
	operator := guild.Members[operatorID]
	target := guild.Members[targetID]

	if operator.Role == RoleMember {
		return fmt.Errorf("你沒有升職的權限")
	}

	if target.Role == RoleOfficer || target.Role == RoleLeader {
		return fmt.Errorf("目標玩家已是副會長或以上職位")
	}

	target.Role = RoleOfficer
	return nil
}

// AddContribution 增加成員貢獻（擊破目標、賺取金幣時呼叫）
func (m *Manager) AddContribution(playerID string, amount int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	guildID, ok := m.playerGuild[playerID]
	if !ok {
		return
	}

	guild := m.guilds[guildID]
	if member, ok := guild.Members[playerID]; ok {
		member.Contribution += amount
	}
}

// UpdateTaskProgress 更新公會任務進度
// 回傳完成的任務列表（用於通知）
func (m *Manager) UpdateTaskProgress(playerID string, taskType GuildTaskType, amount int) []*GuildTask {
	m.mu.Lock()
	defer m.mu.Unlock()

	guildID, ok := m.playerGuild[playerID]
	if !ok {
		return nil
	}

	guild := m.guilds[guildID]
	now := time.Now()

	var completed []*GuildTask
	for _, task := range guild.Tasks {
		if task.Type != taskType || task.Completed {
			continue
		}
		// 檢查是否需要重置
		if now.After(task.ResetAt) {
			task.Current = 0
			task.Completed = false
			task.ResetAt = nextMidnightUTC8()
		}
		task.Current += amount
		if task.Current >= task.Target {
			task.Current = task.Target
			task.Completed = true
			// 公會獲得經驗值
			guild.Exp += task.Target / 10
			m.checkGuildLevelUp(guild)
			completed = append(completed, task)
		}
	}

	// 更新公會統計
	switch taskType {
	case TaskKillTargets:
		guild.TotalKills += amount
	case TaskEarnCoins:
		guild.TotalCoins += amount
	}

	return completed
}

// SetOnlineStatus 設定成員在線狀態
func (m *Manager) SetOnlineStatus(playerID string, isOnline bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	guildID, ok := m.playerGuild[playerID]
	if !ok {
		return
	}

	guild := m.guilds[guildID]
	if member, ok := guild.Members[playerID]; ok {
		member.IsOnline = isOnline
	}
}

// GetPlayerGuildID 取得玩家所在公會 ID（空字串=不在公會）
func (m *Manager) GetPlayerGuildID(playerID string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.playerGuild[playerID]
}

// GetGuild 取得公會資訊
func (m *Manager) GetGuild(guildID string) *Guild {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.guilds[guildID]
}

// GetPlayerGuild 取得玩家所在公會
func (m *Manager) GetPlayerGuild(playerID string) *Guild {
	m.mu.RLock()
	defer m.mu.RUnlock()

	guildID, ok := m.playerGuild[playerID]
	if !ok {
		return nil
	}
	return m.guilds[guildID]
}

// GetAllGuilds 取得所有公會列表（用於搜尋）
func (m *Manager) GetAllGuilds() []*Guild {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]*Guild, 0, len(m.guilds))
	for _, g := range m.guilds {
		result = append(result, g)
	}
	return result
}

// GetGuildMemberIDs 取得公會所有成員 ID
func (m *Manager) GetGuildMemberIDs(guildID string) []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	guild, ok := m.guilds[guildID]
	if !ok {
		return nil
	}

	ids := make([]string, 0, len(guild.Members))
	for id := range guild.Members {
		ids = append(ids, id)
	}
	return ids
}

// GetGuildCount 取得公會數量
func (m *Manager) GetGuildCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.guilds)
}

// generateDailyTasks 生成每日公會任務
func (m *Manager) generateDailyTasks() []*GuildTask {
	resetAt := nextMidnightUTC8()
	return []*GuildTask{
		{
			ID:          "guild_kill_100",
			Type:        TaskKillTargets,
			Name:        "公會討伐",
			Description: "公會成員合力擊破 100 個目標",
			Icon:        "⚔️",
			Target:      100,
			Current:     0,
			Reward:      500,
			Completed:   false,
			ResetAt:     resetAt,
		},
		{
			ID:          "guild_boss_3",
			Type:        TaskKillBoss,
			Name:        "公會 BOSS 挑戰",
			Description: "公會成員合力擊殺 3 隻 BOSS",
			Icon:        "👹",
			Target:      3,
			Current:     0,
			Reward:      1000,
			Completed:   false,
			ResetAt:     resetAt,
		},
		{
			ID:          "guild_coins_10000",
			Type:        TaskEarnCoins,
			Name:        "公會財富積累",
			Description: "公會成員合力賺取 10000 金幣",
			Icon:        "💰",
			Target:      10000,
			Current:     0,
			Reward:      800,
			Completed:   false,
			ResetAt:     resetAt,
		},
	}
}

// checkGuildLevelUp 檢查公會是否升級（非 thread-safe）
func (m *Manager) checkGuildLevelUp(guild *Guild) {
	// 升級所需經驗：等級 × 100
	for guild.Level < 5 {
		needed := guild.Level * 100
		if guild.Exp >= needed {
			guild.Exp -= needed
			guild.Level++
		} else {
			break
		}
	}
}

// findNewLeader 找到新會長（非 thread-safe）
// 優先選副會長，其次選最早加入的成員
func (m *Manager) findNewLeader(guild *Guild, excludeID string) string {
	// 先找副會長
	for id, member := range guild.Members {
		if id != excludeID && member.Role == RoleOfficer {
			return id
		}
	}
	// 再找最早加入的成員
	var earliest *GuildMember
	var earliestID string
	for id, member := range guild.Members {
		if id == excludeID {
			continue
		}
		if earliest == nil || member.JoinedAt.Before(earliest.JoinedAt) {
			earliest = member
			earliestID = id
		}
	}
	return earliestID
}

// nextMidnightUTC8 計算下一個 UTC+8 00:00
func nextMidnightUTC8() time.Time {
	loc := time.FixedZone("UTC+8", 8*60*60)
	now := time.Now().In(loc)
	next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, loc)
	return next.UTC()
}
