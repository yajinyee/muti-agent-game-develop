// pubsub_test.go — Redis Pub/Sub 廣播層測試（DAY-060）
package ws

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

// TestNewPubSubBroker_NoRedis 無 Redis URL 時回傳 nil（降級模式）
func TestNewPubSubBroker_NoRedis(t *testing.T) {
	hub := NewHub()
	broker := NewPubSubBroker("", "room-001", "server-1", hub)
	if broker != nil {
		t.Error("expected nil broker when redisURL is empty")
	}
}

// TestNewPubSubBroker_InvalidURL 無效 URL 時回傳 nil
func TestNewPubSubBroker_InvalidURL(t *testing.T) {
	hub := NewHub()
	broker := NewPubSubBroker("not-a-valid-url", "room-001", "server-1", hub)
	if broker != nil {
		t.Error("expected nil broker for invalid URL")
	}
}

// TestBroadcastWithPubSub_NilBroker broker 為 nil 時降級為純本地廣播
func TestBroadcastWithPubSub_NilBroker(t *testing.T) {
	hub := NewHub()

	// 建立一個假客戶端
	ch := make(chan []byte, 10)
	client := &Client{
		ID:   "test-client",
		Hub:  hub,
		send: ch,
		Role: RolePlayer,
	}
	hub.mu.Lock()
	hub.clients["test-client"] = client
	hub.mu.Unlock()

	msg := &Message{Type: "test_msg", Payload: json.RawMessage(`{"data":"hello"}`)}

	// broker 為 nil，應該只廣播給本地客戶端
	hub.BroadcastWithPubSub(msg, nil)

	select {
	case data := <-ch:
		var received Message
		if err := json.Unmarshal(data, &received); err != nil {
			t.Fatalf("unmarshal error: %v", err)
		}
		if received.Type != "test_msg" {
			t.Errorf("expected type 'test_msg', got '%s'", received.Type)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("expected message to be received within 100ms")
	}
}

// TestLocalBroadcast 本地廣播只發給本地客戶端
func TestLocalBroadcast(t *testing.T) {
	hub := NewHub()

	ch1 := make(chan []byte, 10)
	ch2 := make(chan []byte, 10)

	client1 := &Client{ID: "c1", Hub: hub, send: ch1, Role: RolePlayer}
	client2 := &Client{ID: "c2", Hub: hub, send: ch2, Role: RoleSpectator}

	hub.mu.Lock()
	hub.clients["c1"] = client1
	hub.clients["c2"] = client2
	hub.mu.Unlock()

	msg := &Message{Type: "local_test", Payload: json.RawMessage(`{}`)}
	hub.localBroadcast(msg)

	// 兩個客戶端都應該收到（localBroadcast 不區分角色）
	for _, ch := range []chan []byte{ch1, ch2} {
		select {
		case <-ch:
			// ok
		case <-time.After(100 * time.Millisecond):
			t.Error("expected message to be received")
		}
	}
}

// TestPubSubBroker_Redis 有 REDIS_URL 時執行整合測試
func TestPubSubBroker_Redis(t *testing.T) {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		t.Skip("REDIS_URL not set, skipping Redis integration test")
	}

	hub := NewHub()
	broker := NewPubSubBroker(redisURL, "test-room", "server-test", hub)
	if broker == nil {
		t.Fatal("expected non-nil broker with valid REDIS_URL")
	}
	defer broker.Stop()

	broker.Start()

	// 建立接收客戶端
	ch := make(chan []byte, 10)
	client := &Client{ID: "redis-test-client", Hub: hub, send: ch, Role: RolePlayer}
	hub.mu.Lock()
	hub.clients["redis-test-client"] = client
	hub.mu.Unlock()

	// 模擬另一個 Server 實例發布訊息（不同 serverID）
	otherBroker := NewPubSubBroker(redisURL, "test-room", "server-other", hub)
	if otherBroker == nil {
		t.Fatal("expected non-nil otherBroker")
	}
	defer otherBroker.Stop()

	// 等待訂閱建立
	time.Sleep(100 * time.Millisecond)

	msg := &Message{Type: "redis_test", Payload: json.RawMessage(`{"from":"other"}`)}
	otherBroker.Publish(msg)

	// 等待訊息傳遞
	select {
	case data := <-ch:
		var received Message
		if err := json.Unmarshal(data, &received); err != nil {
			t.Fatalf("unmarshal error: %v", err)
		}
		if received.Type != "redis_test" {
			t.Errorf("expected type 'redis_test', got '%s'", received.Type)
		}
	case <-time.After(2 * time.Second):
		t.Error("expected message from Redis pub/sub within 2s")
	}
}
