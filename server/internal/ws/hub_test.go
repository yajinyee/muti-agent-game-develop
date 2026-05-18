// Package ws 觀戰模式單元測試（DAY-023）
package ws

import (
	"testing"
)

// TestClientRole 確認 ClientRole 常數定義正確
func TestClientRole(t *testing.T) {
	if RolePlayer != "player" {
		t.Errorf("expected RolePlayer='player', got '%s'", RolePlayer)
	}
	if RoleSpectator != "spectator" {
		t.Errorf("expected RoleSpectator='spectator', got '%s'", RoleSpectator)
	}
}

// TestHubPlayerCount 確認 PlayerCount 只計算玩家，不含觀戰者
func TestHubPlayerCount(t *testing.T) {
	h := NewHub()

	// 手動加入玩家和觀戰者（不透過 WebSocket，直接操作 clients map）
	h.mu.Lock()
	h.clients["player-1"] = &Client{ID: "player-1", Role: RolePlayer, send: make(chan []byte, 1)}
	h.clients["player-2"] = &Client{ID: "player-2", Role: RolePlayer, send: make(chan []byte, 1)}
	h.clients["spectator-1"] = &Client{ID: "spectator-1", Role: RoleSpectator, send: make(chan []byte, 1)}
	h.mu.Unlock()

	if h.PlayerCount() != 2 {
		t.Errorf("expected PlayerCount=2, got %d", h.PlayerCount())
	}
	if h.SpectatorCount() != 1 {
		t.Errorf("expected SpectatorCount=1, got %d", h.SpectatorCount())
	}
	if h.ClientCount() != 3 {
		t.Errorf("expected ClientCount=3, got %d", h.ClientCount())
	}
}

// TestHubSpectatorCount 確認 SpectatorCount 只計算觀戰者
func TestHubSpectatorCount(t *testing.T) {
	h := NewHub()

	// 只有觀戰者
	h.mu.Lock()
	h.clients["s1"] = &Client{ID: "s1", Role: RoleSpectator, send: make(chan []byte, 1)}
	h.clients["s2"] = &Client{ID: "s2", Role: RoleSpectator, send: make(chan []byte, 1)}
	h.mu.Unlock()

	if h.PlayerCount() != 0 {
		t.Errorf("expected PlayerCount=0, got %d", h.PlayerCount())
	}
	if h.SpectatorCount() != 2 {
		t.Errorf("expected SpectatorCount=2, got %d", h.SpectatorCount())
	}
}

// TestHubOnConnectNotCalledForSpectator 確認觀戰者連線不觸發 OnConnect
func TestHubOnConnectNotCalledForSpectator(t *testing.T) {
	h := NewHub()
	connectCalled := 0
	h.OnConnect = func(clientID string) {
		connectCalled++
	}

	// 模擬觀戰者 Register
	spectator := &Client{ID: "spectator-x", Role: RoleSpectator, send: make(chan []byte, 1)}
	h.Register(spectator)

	if connectCalled != 0 {
		t.Errorf("OnConnect should NOT be called for spectator, but was called %d times", connectCalled)
	}

	// 模擬玩家 Register
	player := &Client{ID: "player-x", Role: RolePlayer, send: make(chan []byte, 1)}
	h.Register(player)

	if connectCalled != 1 {
		t.Errorf("OnConnect should be called once for player, but was called %d times", connectCalled)
	}
}

// TestHubOnDisconnectNotCalledForSpectator 確認觀戰者斷線不觸發 OnDisconnect
func TestHubOnDisconnectNotCalledForSpectator(t *testing.T) {
	h := NewHub()
	disconnectCalled := 0
	h.OnDisconnect = func(clientID string) {
		disconnectCalled++
	}

	// 先 Register 再 Unregister
	spectator := &Client{ID: "spectator-y", Role: RoleSpectator, send: make(chan []byte, 1)}
	h.Register(spectator)
	h.Unregister(spectator)

	if disconnectCalled != 0 {
		t.Errorf("OnDisconnect should NOT be called for spectator, but was called %d times", disconnectCalled)
	}

	// 玩家斷線應該觸發
	player := &Client{ID: "player-y", Role: RolePlayer, send: make(chan []byte, 1)}
	h.Register(player)
	h.Unregister(player)

	if disconnectCalled != 1 {
		t.Errorf("OnDisconnect should be called once for player, but was called %d times", disconnectCalled)
	}
}

// TestHubSendToSpectator 確認可以傳送訊息給觀戰者
func TestHubSendToSpectator(t *testing.T) {
	h := NewHub()

	spectator := &Client{ID: "spectator-z", Role: RoleSpectator, send: make(chan []byte, 10)}
	h.mu.Lock()
	h.clients[spectator.ID] = spectator
	h.mu.Unlock()

	msg := &Message{Type: MsgGameState, Payload: GameStatePayload{State: "normal_play"}}
	err := h.Send(spectator.ID, msg)
	if err != nil {
		t.Fatalf("Send to spectator failed: %v", err)
	}

	if len(spectator.send) != 1 {
		t.Errorf("expected 1 message in spectator send channel, got %d", len(spectator.send))
	}
}

// TestHubBroadcastReachesSpectator 確認廣播訊息也會傳給觀戰者
func TestHubBroadcastReachesSpectator(t *testing.T) {
	h := NewHub()

	player := &Client{ID: "p1", Role: RolePlayer, send: make(chan []byte, 10)}
	spectator := &Client{ID: "s1", Role: RoleSpectator, send: make(chan []byte, 10)}
	h.mu.Lock()
	h.clients[player.ID] = player
	h.clients[spectator.ID] = spectator
	h.mu.Unlock()

	msg := &Message{Type: MsgTargetSpawn, Payload: nil}
	h.Broadcast(msg)

	if len(player.send) != 1 {
		t.Errorf("player should receive broadcast, got %d messages", len(player.send))
	}
	if len(spectator.send) != 1 {
		t.Errorf("spectator should receive broadcast, got %d messages", len(spectator.send))
	}
}
