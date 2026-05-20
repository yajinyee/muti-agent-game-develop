package friend

import (
	"testing"
)

func TestNew(t *testing.T) {
	m := New()
	if m == nil {
		t.Fatal("New() returned nil")
	}
}

func TestSendRequest_Success(t *testing.T) {
	m := New()
	ok := m.SendRequest("player1", "player2")
	if !ok {
		t.Fatal("SendRequest should succeed")
	}
}

func TestSendRequest_SelfRequest(t *testing.T) {
	m := New()
	ok := m.SendRequest("player1", "player1")
	if ok {
		t.Fatal("SendRequest to self should fail")
	}
}

func TestSendRequest_DuplicateRequest(t *testing.T) {
	m := New()
	m.SendRequest("player1", "player2")
	ok := m.SendRequest("player1", "player2")
	if ok {
		t.Fatal("Duplicate SendRequest should fail")
	}
}

func TestSendRequest_MutualRequest(t *testing.T) {
	m := New()
	// A 發請求給 B
	m.SendRequest("player1", "player2")
	// B 發請求給 A（應該自動接受）
	ok := m.SendRequest("player2", "player1")
	if !ok {
		t.Fatal("Mutual request should auto-accept")
	}
	// 應該已成為好友
	if !m.IsFriend("player1", "player2") {
		t.Fatal("Should be friends after mutual request")
	}
}

func TestAcceptRequest_Success(t *testing.T) {
	m := New()
	m.SendRequest("player1", "player2")
	ok := m.AcceptRequest("player1", "player2")
	if !ok {
		t.Fatal("AcceptRequest should succeed")
	}
	if !m.IsFriend("player1", "player2") {
		t.Fatal("Should be friends after accept")
	}
	if !m.IsFriend("player2", "player1") {
		t.Fatal("Friendship should be bidirectional")
	}
}

func TestAcceptRequest_NoRequest(t *testing.T) {
	m := New()
	ok := m.AcceptRequest("player1", "player2")
	if ok {
		t.Fatal("AcceptRequest without request should fail")
	}
}

func TestRejectRequest_Success(t *testing.T) {
	m := New()
	m.SendRequest("player1", "player2")
	ok := m.RejectRequest("player1", "player2")
	if !ok {
		t.Fatal("RejectRequest should succeed")
	}
	if m.IsFriend("player1", "player2") {
		t.Fatal("Should not be friends after reject")
	}
}

func TestRemoveFriend_Success(t *testing.T) {
	m := New()
	m.SendRequest("player1", "player2")
	m.AcceptRequest("player1", "player2")
	ok := m.RemoveFriend("player1", "player2")
	if !ok {
		t.Fatal("RemoveFriend should succeed")
	}
	if m.IsFriend("player1", "player2") {
		t.Fatal("Should not be friends after remove")
	}
	if m.IsFriend("player2", "player1") {
		t.Fatal("Friendship removal should be bidirectional")
	}
}

func TestGetFriendIDs(t *testing.T) {
	m := New()
	m.SendRequest("player1", "player2")
	m.AcceptRequest("player1", "player2")
	m.SendRequest("player1", "player3")
	m.AcceptRequest("player1", "player3")

	ids := m.GetFriendIDs("player1")
	if len(ids) != 2 {
		t.Errorf("expected 2 friends, got %d", len(ids))
	}
}

func TestGetPendingRequests(t *testing.T) {
	m := New()
	m.SendRequest("player1", "player2")
	m.SendRequest("player3", "player2")

	pending := m.GetPendingRequests("player2")
	if len(pending) != 2 {
		t.Errorf("expected 2 pending requests, got %d", len(pending))
	}
}

func TestGetFriendCount(t *testing.T) {
	m := New()
	m.SendRequest("player1", "player2")
	m.AcceptRequest("player1", "player2")

	count := m.GetFriendCount("player1")
	if count != 1 {
		t.Errorf("expected 1 friend, got %d", count)
	}
}

func TestSendRequest_AlreadyFriends(t *testing.T) {
	m := New()
	m.SendRequest("player1", "player2")
	m.AcceptRequest("player1", "player2")
	ok := m.SendRequest("player1", "player2")
	if ok {
		t.Fatal("SendRequest to existing friend should fail")
	}
}
