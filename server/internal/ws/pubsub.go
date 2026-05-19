// Package ws — Redis Pub/Sub 廣播層（DAY-060）
//
// 解決水平擴展問題：多個 Server 實例時，Server A 的廣播無法到達 Server B 的客戶端。
// 解法：每個 Hub 訂閱 Redis channel，收到訊息後廣播給本地客戶端。
//
// 架構：
//   Hub.BroadcastWithPubSub(msg) → 1) 廣播給本地客戶端
//                                  2) Publish 到 Redis channel
//   Redis channel → 其他 Server 實例的 Hub.localBroadcast(msg)
//
// Channel 命名：game:broadcast:{roomID}
//
// 注意：無 Redis 時自動降級為純本地廣播，不影響現有功能。
package ws

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	pubsubChannelPrefix = "game:broadcast:"
	pubsubTimeout       = 3 * time.Second
)

// PubSubBroker Redis Pub/Sub 廣播代理
// 負責跨 Server 實例的訊息廣播
type PubSubBroker struct {
	client    *redis.Client
	roomID    string
	channel   string
	hub       *Hub
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	serverID  string // 本 Server 實例 ID（避免收到自己發的訊息）
}

// pubsubEnvelope Redis channel 上的訊息格式
// 包含 serverID 用於過濾自己發的訊息
type pubsubEnvelope struct {
	ServerID string          `json:"server_id"`
	Payload  json.RawMessage `json:"payload"`
}

// NewPubSubBroker 建立 Redis Pub/Sub 代理
// redisURL 為空時回傳 nil（降級模式）
func NewPubSubBroker(redisURL, roomID, serverID string, hub *Hub) *PubSubBroker {
	if redisURL == "" {
		return nil
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Printf("[PubSub] Invalid Redis URL: %v — falling back to local broadcast", err)
		return nil
	}

	client := redis.NewClient(opt)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		log.Printf("[PubSub] Redis ping failed: %v — falling back to local broadcast", err)
		return nil
	}

	brokerCtx, brokerCancel := context.WithCancel(context.Background())
	broker := &PubSubBroker{
		client:   client,
		roomID:   roomID,
		channel:  pubsubChannelPrefix + roomID,
		hub:      hub,
		ctx:      brokerCtx,
		cancel:   brokerCancel,
		serverID: serverID,
	}

	log.Printf("[PubSub] Connected to Redis, channel: %s (serverID: %s)", broker.channel, serverID)
	return broker
}

// Start 啟動訂閱 goroutine
func (b *PubSubBroker) Start() {
	b.wg.Add(1)
	go b.subscribeLoop()
}

// Stop 停止訂閱，釋放資源
func (b *PubSubBroker) Stop() {
	b.cancel()
	b.wg.Wait()
	b.client.Close()
	log.Printf("[PubSub] Stopped (channel: %s)", b.channel)
}

// Publish 發布訊息到 Redis channel（供其他 Server 實例接收）
func (b *PubSubBroker) Publish(msg *Message) {
	payload, err := json.Marshal(msg)
	if err != nil {
		log.Printf("[PubSub] Marshal error: %v", err)
		return
	}

	envelope := pubsubEnvelope{
		ServerID: b.serverID,
		Payload:  payload,
	}
	data, err := json.Marshal(envelope)
	if err != nil {
		log.Printf("[PubSub] Envelope marshal error: %v", err)
		return
	}

	ctx, cancel := context.WithTimeout(b.ctx, pubsubTimeout)
	defer cancel()

	if err := b.client.Publish(ctx, b.channel, data).Err(); err != nil {
		log.Printf("[PubSub] Publish error: %v", err)
	}
}

// subscribeLoop 訂閱 Redis channel，收到訊息後廣播給本地客戶端
func (b *PubSubBroker) subscribeLoop() {
	defer b.wg.Done()

	pubsub := b.client.Subscribe(b.ctx, b.channel)
	defer pubsub.Close()

	log.Printf("[PubSub] Subscribed to channel: %s", b.channel)

	ch := pubsub.Channel()
	for {
		select {
		case <-b.ctx.Done():
			return
		case redisMsg, ok := <-ch:
			if !ok {
				log.Printf("[PubSub] Channel closed, reconnecting...")
				// 嘗試重新訂閱
				time.Sleep(1 * time.Second)
				pubsub = b.client.Subscribe(b.ctx, b.channel)
				ch = pubsub.Channel()
				continue
			}

			var envelope pubsubEnvelope
			if err := json.Unmarshal([]byte(redisMsg.Payload), &envelope); err != nil {
				log.Printf("[PubSub] Unmarshal error: %v", err)
				continue
			}

			// 過濾自己發的訊息（避免重複廣播）
			if envelope.ServerID == b.serverID {
				continue
			}

			var msg Message
			if err := json.Unmarshal(envelope.Payload, &msg); err != nil {
				log.Printf("[PubSub] Payload unmarshal error: %v", err)
				continue
			}

			// 廣播給本地客戶端（不再 publish 到 Redis，避免無限循環）
			b.hub.localBroadcast(&msg)
		}
	}
}

// localBroadcast 只廣播給本地客戶端（不 publish 到 Redis）
// 供 subscribeLoop 使用，避免無限循環
func (h *Hub) localBroadcast(msg *Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("[WS] localBroadcast marshal error: %v", err)
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, client := range h.clients {
		select {
		case client.send <- data:
			h.MsgSent.Add(1)
			h.BytesSentRaw.Add(int64(len(data)))
		default:
			h.MsgDropped.Add(1)
		}
	}
}

// BroadcastWithPubSub 廣播訊息（本地 + Redis pub/sub）
// broker 為 nil 時降級為純本地廣播
func (h *Hub) BroadcastWithPubSub(msg *Message, broker *PubSubBroker) {
	// 1. 廣播給本地客戶端
	h.Broadcast(msg)

	// 2. 如果有 Redis broker，publish 到 Redis channel（供其他 Server 實例接收）
	if broker != nil {
		broker.Publish(msg)
	}
}
