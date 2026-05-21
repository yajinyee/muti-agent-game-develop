package room

import (
	"testing"
)

func TestNewManager(t *testing.T) {
	m := NewManager()
	rooms := m.ListRooms()
	if len(rooms) != 3 {
		t.Errorf("expected 3 default rooms, got %d", len(rooms))
	}
}

func TestGetRoom(t *testing.T) {
	m := NewManager()

	r, ok := m.GetRoom("room-001")
	if !ok {
		t.Fatal("room-001 should exist")
	}
	if r.ID != "room-001" {
		t.Errorf("expected room-001, got %s", r.ID)
	}

	_, ok = m.GetRoom("room-999")
	if ok {
		t.Error("room-999 should not exist")
	}
}

func TestJoinAndLeaveRoom(t *testing.T) {
	m := NewManager()

	// 加入玩家
	r, err := m.JoinRoom("room-001", "player-1")
	if err != nil {
		t.Fatalf("join failed: %v", err)
	}
	if r.PlayerCount() != 1 {
		t.Errorf("expected 1 player, got %d", r.PlayerCount())
	}

	// 重複加入同一玩家（應該成功，因為 map 覆蓋）
	_, err = m.JoinRoom("room-001", "player-1")
	if err != nil {
		t.Fatalf("re-join failed: %v", err)
	}

	// 離開
	m.LeaveRoom("room-001", "player-1")
	if r.PlayerCount() != 0 {
		t.Errorf("expected 0 players after leave, got %d", r.PlayerCount())
	}
}

func TestRoomFull(t *testing.T) {
	m := NewManager()

	// 建立一個只能容納 2 人的房間
	r, err := m.CreateRoom(Config{
		Name:        "測試房間",
		MaxPlayers:  2,
		MinBetLevel: 1,
		MaxBetLevel: 10,
		Theme:       "chiikawa",
		RTPTarget:   0.94,
	})
	if err != nil {
		t.Fatalf("create room failed: %v", err)
	}

	// 加入 2 個玩家
	_, err = m.JoinRoom(r.ID, "p1")
	if err != nil {
		t.Fatalf("join p1 failed: %v", err)
	}
	_, err = m.JoinRoom(r.ID, "p2")
	if err != nil {
		t.Fatalf("join p2 failed: %v", err)
	}

	// 第 3 個玩家應該被拒絕
	_, err = m.JoinRoom(r.ID, "p3")
	if err == nil {
		t.Error("expected error when room is full")
	}
}

func TestFindLeastPopulated(t *testing.T) {
	m := NewManager()

	// 加入玩家到 room-001
	m.JoinRoom("room-001", "p1")
	m.JoinRoom("room-001", "p2")
	m.JoinRoom("room-001", "p3")

	// 加入玩家到 room-002
	m.JoinRoom("room-002", "p4")

	// 最少人的應該是 room-003（0人）
	r := m.FindLeastPopulated()
	if r.ID != "room-003" {
		t.Errorf("expected room-003 (0 players), got %s (%d players)", r.ID, r.PlayerCount())
	}
}

func TestFindPlayerRoom(t *testing.T) {
	m := NewManager()

	m.JoinRoom("room-002", "player-x")

	r := m.FindPlayerRoom("player-x")
	if r == nil {
		t.Fatal("should find player-x's room")
	}
	if r.ID != "room-002" {
		t.Errorf("expected room-002, got %s", r.ID)
	}

	// 不存在的玩家
	r = m.FindPlayerRoom("nobody")
	if r != nil {
		t.Error("should return nil for unknown player")
	}
}

func TestGetOrDefault(t *testing.T) {
	m := NewManager()

	// 空字串應該回傳 room-001
	r := m.GetOrDefault("")
	if r.ID != "room-001" {
		t.Errorf("expected room-001 for empty roomID, got %s", r.ID)
	}

	// 不存在的房間應該回傳人數最少的
	r = m.GetOrDefault("room-999")
	if r == nil {
		t.Error("should return a room even for unknown roomID")
	}
}

func TestDeleteRoom(t *testing.T) {
	m := NewManager()

	// 建立一個空房間
	r, _ := m.CreateRoom(Config{
		Name:       "臨時房間",
		MaxPlayers: 4,
		Theme:      "chiikawa",
		RTPTarget:  0.94,
	})

	// 刪除空房間應該成功
	err := m.DeleteRoom(r.ID)
	if err != nil {
		t.Fatalf("delete empty room failed: %v", err)
	}

	// 確認已刪除
	_, ok := m.GetRoom(r.ID)
	if ok {
		t.Error("room should be deleted")
	}

	// 刪除有玩家的房間應該失敗
	m.JoinRoom("room-001", "p1")
	err = m.DeleteRoom("room-001")
	if err == nil {
		t.Error("should fail to delete room with players")
	}
}
