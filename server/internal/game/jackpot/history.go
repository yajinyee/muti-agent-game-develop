// history.go — Jackpot 中獎歷史記錄（DAY-048e）
// 保存最近 10 筆中獎記錄，供 /jackpot HTTP 端點和 Client 顯示
package jackpot

import (
	"sync"
	"time"
)

// WinRecord 單筆中獎記錄
type WinRecord struct {
	Level      Level     `json:"level"`
	Amount     int       `json:"amount"`
	WinnerID   string    `json:"winner_id"`
	WinnerName string    `json:"winner_name"`
	WonAt      time.Time `json:"won_at"`
	WonAtMs    int64     `json:"won_at_ms"` // Unix milliseconds（Client 用）
}

// History Jackpot 中獎歷史
type History struct {
	mu      sync.RWMutex
	records []WinRecord
	maxSize int
}

// NewHistory 建立歷史記錄（保存最近 maxSize 筆）
func NewHistory(maxSize int) *History {
	return &History{
		records: make([]WinRecord, 0, maxSize),
		maxSize: maxSize,
	}
}

// Add 加入一筆中獎記錄
func (h *History) Add(win *JackpotWin, winnerName string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	record := WinRecord{
		Level:      win.Level,
		Amount:     win.Amount,
		WinnerID:   win.WinnerID,
		WinnerName: winnerName,
		WonAt:      win.WonAt,
		WonAtMs:    win.WonAt.UnixMilli(),
	}

	// 插入到最前面（最新的在前）
	h.records = append([]WinRecord{record}, h.records...)

	// 超過上限時截斷
	if len(h.records) > h.maxSize {
		h.records = h.records[:h.maxSize]
	}
}

// GetRecent 取得最近 n 筆記錄
func (h *History) GetRecent(n int) []WinRecord {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if n > len(h.records) {
		n = len(h.records)
	}
	result := make([]WinRecord, n)
	copy(result, h.records[:n])
	return result
}

// Count 取得總記錄數
func (h *History) Count() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.records)
}
