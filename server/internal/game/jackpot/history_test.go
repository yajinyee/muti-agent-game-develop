package jackpot

import (
	"testing"
	"time"
)

func TestHistory_AddAndGet(t *testing.T) {
	h := NewHistory(10)

	win := &JackpotWin{
		Level:    LevelMini,
		Amount:   500,
		WinnerID: "player1",
		WonAt:    time.Now(),
	}
	h.Add(win, "Player One")

	records := h.GetRecent(5)
	if len(records) != 1 {
		t.Fatalf("Expected 1 record, got %d", len(records))
	}
	if records[0].Level != LevelMini {
		t.Errorf("Level = %s, want mini", records[0].Level)
	}
	if records[0].WinnerName != "Player One" {
		t.Errorf("WinnerName = %s, want Player One", records[0].WinnerName)
	}
}

func TestHistory_MaxSize(t *testing.T) {
	h := NewHistory(3)

	for i := 0; i < 5; i++ {
		win := &JackpotWin{
			Level:    LevelMini,
			Amount:   500 + i,
			WinnerID: "player1",
			WonAt:    time.Now(),
		}
		h.Add(win, "Player")
	}

	records := h.GetRecent(10)
	if len(records) != 3 {
		t.Errorf("Expected 3 records (max), got %d", len(records))
	}
	// 最新的在最前面
	if records[0].Amount != 504 {
		t.Errorf("Latest record amount = %d, want 504", records[0].Amount)
	}
}

func TestHistory_NewestFirst(t *testing.T) {
	h := NewHistory(10)

	for i := 0; i < 3; i++ {
		win := &JackpotWin{
			Level:    LevelMini,
			Amount:   i * 100,
			WinnerID: "player1",
			WonAt:    time.Now(),
		}
		h.Add(win, "Player")
	}

	records := h.GetRecent(3)
	// 最新的（amount=200）應該在最前面
	if records[0].Amount != 200 {
		t.Errorf("First record amount = %d, want 200 (newest)", records[0].Amount)
	}
}

func TestHistory_Empty(t *testing.T) {
	h := NewHistory(10)
	records := h.GetRecent(5)
	if len(records) != 0 {
		t.Errorf("Expected 0 records, got %d", len(records))
	}
}
