// Package analytics 遊戲數據埋點模組
// 追蹤玩家行為、遊戲事件、財務指標，輸出到 JSON 日誌
package analytics

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

// EventType 事件類型
type EventType string

const (
	// 玩家行為事件
	EventAttack      EventType = "attack"       // 玩家攻擊
	EventKill        EventType = "kill"          // 擊破目標
	EventBetChange   EventType = "bet_change"    // 切換投注
	EventLockTarget  EventType = "lock_target"   // 鎖定目標
	EventAutoToggle  EventType = "auto_toggle"   // 切換自動模式

	// 遊戲系統事件
	EventBossSpawn   EventType = "boss_spawn"    // BOSS 出現
	EventBossKill    EventType = "boss_kill"     // BOSS 擊敗
	EventBonusStart  EventType = "bonus_start"   // Bonus 開始
	EventBonusEnd    EventType = "bonus_end"     // Bonus 結束
	EventSpecialTarget EventType = "special_target" // 特殊目標出現

	// 財務事件
	EventReward      EventType = "reward"        // 獎勵發放
	EventPlayerJoin  EventType = "player_join"   // 玩家加入
	EventPlayerLeave EventType = "player_leave"  // 玩家離開
)

// Event 單筆事件記錄
type Event struct {
	Timestamp int64             `json:"ts"`         // Unix milliseconds
	EventType EventType         `json:"event"`      // 事件類型
	PlayerID  string            `json:"player_id"`  // 玩家 ID
	RoomID    string            `json:"room_id"`    // 房間 ID
	Data      map[string]interface{} `json:"data"` // 事件附加資料
}

// SessionStats 單次遊戲會話統計（記憶體中，玩家離開時輸出）
type SessionStats struct {
	PlayerID      string    `json:"player_id"`
	RoomID        string    `json:"room_id"`
	JoinTime      time.Time `json:"join_time"`
	LeaveTime     time.Time `json:"leave_time"`
	DurationSec   float64   `json:"duration_sec"`

	// 攻擊統計
	TotalAttacks  int64   `json:"total_attacks"`
	TotalHits     int64   `json:"total_hits"`
	TotalKills    int64   `json:"total_kills"`
	HitRate       float64 `json:"hit_rate"`       // hits / attacks

	// 財務統計
	TotalBet      int64   `json:"total_bet"`      // 總投入金幣
	TotalReward   int64   `json:"total_reward"`   // 總獲得金幣
	ActualRTP     float64 `json:"actual_rtp"`     // reward / bet
	MaxSingleWin  int64   `json:"max_single_win"` // 單次最高獎勵

	// 遊戲事件統計
	BossCount     int64   `json:"boss_count"`     // 遭遇 BOSS 次數
	BossKillCount int64   `json:"boss_kill_count"`// 擊敗 BOSS 次數
	BonusCount    int64   `json:"bonus_count"`    // 觸發 Bonus 次數

	// 投注等級分布（LV1-LV10 各用了多少次攻擊）
	BetLevelDist  map[int]int64 `json:"bet_level_dist"`

	// 目標類型擊破分布
	TargetKillDist map[string]int64 `json:"target_kill_dist"`
}

// RoomStats 房間整體統計（持續累積）
type RoomStats struct {
	RoomID        string    `json:"room_id"`
	StartTime     time.Time `json:"start_time"`
	LastUpdate    time.Time `json:"last_update"`

	// 玩家統計
	TotalPlayers   int64 `json:"total_players"`   // 歷史總玩家數
	PeakPlayers    int64 `json:"peak_players"`    // 同時在線峰值
	CurrentPlayers int64 `json:"current_players"` // 當前在線

	// 攻擊統計
	TotalAttacks int64 `json:"total_attacks"`
	TotalHits    int64 `json:"total_hits"`
	TotalKills   int64 `json:"total_kills"`

	// 財務統計
	TotalBet    int64   `json:"total_bet"`
	TotalReward int64   `json:"total_reward"`
	OverallRTP  float64 `json:"overall_rtp"`

	// 遊戲事件統計
	BossSpawnCount int64 `json:"boss_spawn_count"`
	BossKillCount  int64 `json:"boss_kill_count"`
	BonusCount     int64 `json:"bonus_count"`
}

// Tracker 數據追蹤器（Singleton）
type Tracker struct {
	roomID   string
	logFile  *os.File
	logMu    sync.Mutex
	sessions map[string]*SessionStats // playerID → session
	sessMu   sync.RWMutex
	room     RoomStats
	roomMu   sync.RWMutex  // 保護 room 欄位

	// 原子計數器（高頻事件用，避免鎖競爭）
	attackCount atomic.Int64
	hitCount    atomic.Int64
	killCount   atomic.Int64
	betTotal    atomic.Int64
	rewardTotal atomic.Int64
}

var (
	instance *Tracker
	once     sync.Once
)

// Init 初始化追蹤器（Server 啟動時呼叫）
func Init(roomID string, logDir string) *Tracker {
	once.Do(func() {
		instance = &Tracker{
			roomID:   roomID,
			sessions: make(map[string]*SessionStats),
		}
		instance.room = RoomStats{
			RoomID:    roomID,
			StartTime: time.Now(),
		}

		// 建立日誌目錄
		if err := os.MkdirAll(logDir, 0755); err != nil {
			log.Printf("[analytics] Failed to create log dir: %v", err)
			return
		}

		// 開啟日誌檔案（按日期命名）
		logPath := fmt.Sprintf("%s/events-%s.jsonl",
			logDir, time.Now().Format("2006-01-02"))
		f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("[analytics] Failed to open log file: %v", err)
			return
		}
		instance.logFile = f
		log.Printf("[analytics] Logging to %s", logPath)
	})
	return instance
}

// Get 取得追蹤器實例
func Get() *Tracker {
	return instance
}

// newTracker 建立新的 Tracker 實例（測試用，不走 singleton）
func newTracker(roomID string, logDir string) *Tracker {
	t := &Tracker{
		roomID:   roomID,
		sessions: make(map[string]*SessionStats),
	}
	t.room = RoomStats{
		RoomID:    roomID,
		StartTime: time.Now(),
	}
	if logDir != "" {
		_ = os.MkdirAll(logDir, 0755)
		logPath := fmt.Sprintf("%s/events-%s.jsonl", logDir, time.Now().Format("2006-01-02"))
		f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err == nil {
			t.logFile = f
		}
	}
	return t
}

// Track 記錄一個事件
func (t *Tracker) Track(eventType EventType, playerID string, data map[string]interface{}) {
	if t == nil {
		return
	}
	event := Event{
		Timestamp: time.Now().UnixMilli(),
		EventType: eventType,
		PlayerID:  playerID,
		RoomID:    t.roomID,
		Data:      data,
	}
	t.writeEvent(event)
	t.updateStats(event)
}

// writeEvent 寫入 JSONL 日誌（每行一個 JSON 事件）
func (t *Tracker) writeEvent(event Event) {
	if t.logFile == nil {
		return
	}
	b, err := json.Marshal(event)
	if err != nil {
		return
	}
	t.logMu.Lock()
	defer t.logMu.Unlock()
	t.logFile.Write(b)
	t.logFile.WriteString("\n")
}

// updateStats 更新記憶體中的統計數據
func (t *Tracker) updateStats(event Event) {
	t.roomMu.Lock()
	t.room.LastUpdate = time.Now()
	t.roomMu.Unlock()

	switch event.EventType {
	case EventPlayerJoin:
		t.sessMu.Lock()
		t.sessions[event.PlayerID] = &SessionStats{
			PlayerID:       event.PlayerID,
			RoomID:         t.roomID,
			JoinTime:       time.Now(),
			BetLevelDist:   make(map[int]int64),
			TargetKillDist: make(map[string]int64),
		}
		t.sessMu.Unlock()
		t.roomMu.Lock()
		t.room.TotalPlayers++
		t.room.CurrentPlayers++
		if t.room.CurrentPlayers > t.room.PeakPlayers {
			t.room.PeakPlayers = t.room.CurrentPlayers
		}
		t.roomMu.Unlock()

	case EventPlayerLeave:
		t.sessMu.Lock()
		if sess, ok := t.sessions[event.PlayerID]; ok {
			sess.LeaveTime = time.Now()
			sess.DurationSec = sess.LeaveTime.Sub(sess.JoinTime).Seconds()
			if sess.TotalAttacks > 0 {
				sess.HitRate = float64(sess.TotalHits) / float64(sess.TotalAttacks)
			}
			if sess.TotalBet > 0 {
				sess.ActualRTP = float64(sess.TotalReward) / float64(sess.TotalBet)
			}
			// 輸出 session 摘要到日誌
			t.writeEvent(Event{
				Timestamp: time.Now().UnixMilli(),
				EventType: "session_summary",
				PlayerID:  event.PlayerID,
				RoomID:    t.roomID,
				Data: map[string]interface{}{
					"stats": sess,
				},
			})
			delete(t.sessions, event.PlayerID)
		}
		t.sessMu.Unlock()
		t.roomMu.Lock()
		t.room.CurrentPlayers--
		t.roomMu.Unlock()

	case EventAttack:
		t.attackCount.Add(1)
		t.roomMu.Lock()
		t.room.TotalAttacks++
		if betCost, ok := event.Data["bet_cost"].(int); ok {
			t.room.TotalBet += int64(betCost)
		}
		t.roomMu.Unlock()
		// 更新 session
		t.sessMu.RLock()
		if sess, ok := t.sessions[event.PlayerID]; ok {
			sess.TotalAttacks++
			if betLevel, ok := event.Data["bet_level"].(int); ok {
				sess.BetLevelDist[betLevel]++
			}
			if betCost, ok := event.Data["bet_cost"].(int); ok {
				sess.TotalBet += int64(betCost)
			}
		}
		t.sessMu.RUnlock()

	case EventKill:
		t.killCount.Add(1)
		t.roomMu.Lock()
		t.room.TotalKills++
		t.roomMu.Unlock()
		t.sessMu.RLock()
		if sess, ok := t.sessions[event.PlayerID]; ok {
			sess.TotalKills++
			if defID, ok := event.Data["def_id"].(string); ok {
				sess.TargetKillDist[defID]++
			}
		}
		t.sessMu.RUnlock()

	case EventReward:
		t.sessMu.RLock()
		if sess, ok := t.sessions[event.PlayerID]; ok {
			if amount, ok := event.Data["amount"].(int); ok {
				sess.TotalReward += int64(amount)
				if int64(amount) > sess.MaxSingleWin {
					sess.MaxSingleWin = int64(amount)
				}
			}
		}
		t.sessMu.RUnlock()
		if amount, ok := event.Data["amount"].(int); ok {
			t.rewardTotal.Add(int64(amount))
			t.roomMu.Lock()
			t.room.TotalReward += int64(amount)
			if t.room.TotalBet > 0 {
				t.room.OverallRTP = float64(t.room.TotalReward) / float64(t.room.TotalBet)
			}
			t.roomMu.Unlock()
		}

	case EventBossSpawn:
		t.roomMu.Lock()
		t.room.BossSpawnCount++
		t.roomMu.Unlock()
		t.sessMu.RLock()
		if sess, ok := t.sessions[event.PlayerID]; ok {
			sess.BossCount++
		}
		t.sessMu.RUnlock()

	case EventBossKill:
		t.roomMu.Lock()
		t.room.BossKillCount++
		t.roomMu.Unlock()
		t.sessMu.RLock()
		if sess, ok := t.sessions[event.PlayerID]; ok {
			sess.BossKillCount++
		}
		t.sessMu.RUnlock()

	case EventBonusStart:
		t.roomMu.Lock()
		t.room.BonusCount++
		t.roomMu.Unlock()
		t.sessMu.RLock()
		if sess, ok := t.sessions[event.PlayerID]; ok {
			sess.BonusCount++
		}
		t.sessMu.RUnlock()
	}
}

// GetRoomStats 取得房間統計（供 /analytics 端點使用）
func (t *Tracker) GetRoomStats() RoomStats {
	if t == nil {
		return RoomStats{}
	}
	t.roomMu.RLock()
	stats := t.room
	t.roomMu.RUnlock()
	// 計算整體 RTP
	if stats.TotalBet > 0 {
		stats.OverallRTP = float64(stats.TotalReward) / float64(stats.TotalBet)
	}
	return stats
}

// GetSessionStats 取得特定玩家的 session 統計
func (t *Tracker) GetSessionStats(playerID string) *SessionStats {
	if t == nil {
		return nil
	}
	t.sessMu.RLock()
	defer t.sessMu.RUnlock()
	if sess, ok := t.sessions[playerID]; ok {
		// 回傳副本
		copy := *sess
		return &copy
	}
	return nil
}

// Close 關閉追蹤器（Server 關閉時呼叫）
func (t *Tracker) Close() {
	if t == nil || t.logFile == nil {
		return
	}
	// 輸出最終房間統計
	stats := t.GetRoomStats()
	t.writeEvent(Event{
		Timestamp: time.Now().UnixMilli(),
		EventType: "room_summary",
		PlayerID:  "system",
		RoomID:    t.roomID,
		Data: map[string]interface{}{
			"stats": stats,
		},
	})
	t.logMu.Lock()
	t.logFile.Close()
	t.logMu.Unlock()
	log.Printf("[analytics] Tracker closed, final RTP: %.2f%%", stats.OverallRTP*100)
}
