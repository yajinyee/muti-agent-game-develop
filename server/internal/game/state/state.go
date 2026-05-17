// Package state 定義遊戲狀態機
package state

// GameState 遊戲狀態
type GameState string

const (
	StateLoading            GameState = "loading"
	StateLobby              GameState = "lobby"
	StateNormalPlay         GameState = "normal_play"
	StateSpecialTargetEvent GameState = "special_target_event"
	StateBossWarning        GameState = "boss_warning"
	StateBossBattle         GameState = "boss_battle"
	StateBossResult         GameState = "boss_result"
	StateBonusReady         GameState = "bonus_ready"
	StateBonusGame          GameState = "bonus_game"
	StateBonusResult        GameState = "bonus_result"
)

// ValidTransitions 合法的狀態轉換
var ValidTransitions = map[GameState][]GameState{
	StateLoading:            {StateLobby},
	StateLobby:              {StateNormalPlay},
	StateNormalPlay:         {StateSpecialTargetEvent, StateBossWarning, StateBonusReady},
	StateSpecialTargetEvent: {StateNormalPlay, StateBossWarning},
	StateBossWarning:        {StateBossBattle},
	StateBossBattle:         {StateBossResult, StateNormalPlay}, // 加入直接回 NormalPlay（BOSS 超時）
	StateBossResult:         {StateNormalPlay},
	StateBonusReady:         {StateBonusGame},
	StateBonusGame:          {StateBonusResult},
	StateBonusResult:        {StateNormalPlay},
}

// CanTransition 檢查狀態轉換是否合法
func CanTransition(from, to GameState) bool {
	nexts, ok := ValidTransitions[from]
	if !ok {
		return false
	}
	for _, s := range nexts {
		if s == to {
			return true
		}
	}
	return false
}
