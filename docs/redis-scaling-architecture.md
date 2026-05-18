# Redis 水平擴展架構設計

> 版本：v1.0  
> 日期：2026-05-19（DAY-026）  
> 狀態：設計文件（Phase 2/3 實作藍圖）  
> 負責：Go Server Agent + Game Director

---

## 概覽

當前 Prototype 是單 Server 架構，所有遊戲狀態存在記憶體中。  
本文件規劃引入 Redis 作為共享狀態層，支援：
- **Server 重啟後狀態恢復**（玩家金幣、房間狀態）
- **多 Server 實例水平擴展**（負載均衡）
- **跨房間玩家 Session 保留**

---

## 當前架構（單 Server）

```
Browser → WebSocket → Go Server（記憶體狀態）
                           ↓
                      Analytics（JSONL 日誌）
```

**問題：**
- Server 重啟 → 所有玩家金幣歸零
- 無法水平擴展（狀態不共享）
- 玩家換房間 → Session 中斷

---

## 目標架構（Redis 共享狀態）

```
Browser ─┐
Browser ─┼→ Load Balancer（Nginx）
Browser ─┘         ↓
              ┌─────┴─────┐
              ↓           ↓
         Go Server 1  Go Server 2   ← 無狀態，可水平擴展
              ↓           ↓
              └─────┬─────┘
                    ↓
                  Redis
              ┌─────┴─────┐
              ↓           ↓
         Player State  Room State
         (Hash)        (Hash + Pub/Sub)
```

---

## Redis 資料結構設計

### 1. 玩家狀態（Hash）

```
KEY: player:{player_id}
TTL: 24 小時（玩家離線後保留）

FIELDS:
  coins         → int64    玩家金幣
  labor         → int      勞動值（0-100）
  bet_level     → int      當前投注等級（1-10）
  display_name  → string   顯示名稱
  session_score → int64    本次 session 得分
  max_coins     → int64    歷史最高金幣
  kill_count    → int      擊殺數
  room_id       → string   當前所在房間
  last_seen     → int64    最後活躍時間（Unix ms）
```

**Go 操作範例：**
```go
// 儲存玩家狀態
func (r *RedisStore) SavePlayer(p *player.Player) error {
    key := fmt.Sprintf("player:%s", p.ID)
    return r.client.HMSet(ctx, key,
        "coins", p.Coins,
        "labor", p.Labor,
        "bet_level", p.BetLevel,
        "display_name", p.DisplayName,
        "last_seen", time.Now().UnixMilli(),
    ).Err()
}

// 讀取玩家狀態
func (r *RedisStore) LoadPlayer(playerID string) (*player.Player, error) {
    key := fmt.Sprintf("player:%s", playerID)
    vals, err := r.client.HGetAll(ctx, key).Result()
    // ... 解析 vals 到 Player struct
}
```

### 2. 房間狀態（Hash）

```
KEY: room:{room_id}:state
TTL: 永久（房間持續存在）

FIELDS:
  game_state    → string   遊戲狀態（idle/playing/boss/bonus）
  player_count  → int      當前玩家數
  boss_hp       → int      BOSS 當前 HP（BOSS 戰中）
  bonus_active  → bool     Bonus 是否進行中
  last_update   → int64    最後更新時間
```

### 3. 房間玩家列表（Set）

```
KEY: room:{room_id}:players
TTL: 永久

MEMBERS: player_id strings
```

### 4. 排行榜（Sorted Set）

```
KEY: leaderboard:daily:{YYYY-MM-DD}
TTL: 7 天

SCORE: session_score（越高越前）
MEMBER: player_id
```

**Go 操作範例：**
```go
// 更新排行榜
func (r *RedisStore) UpdateLeaderboard(playerID string, score int64) error {
    key := fmt.Sprintf("leaderboard:daily:%s", time.Now().Format("2006-01-02"))
    return r.client.ZAdd(ctx, key, redis.Z{
        Score:  float64(score),
        Member: playerID,
    }).Err()
}

// 取得前 10 名
func (r *RedisStore) GetTopPlayers(n int) ([]LeaderboardEntry, error) {
    key := fmt.Sprintf("leaderboard:daily:%s", time.Now().Format("2006-01-02"))
    results, err := r.client.ZRevRangeWithScores(ctx, key, 0, int64(n-1)).Result()
    // ... 轉換為 LeaderboardEntry
}
```

### 5. 跨 Server 廣播（Pub/Sub）

```
CHANNEL: room:{room_id}:broadcast

MESSAGE FORMAT: JSON（與現有 WebSocket 訊息格式相同）
```

**多 Server 廣播流程：**
```
Server 1 收到玩家攻擊
    ↓
Server 1 計算結果
    ↓
Server 1 PUBLISH room:room-001:broadcast {msg}
    ↓
Server 1 + Server 2 都 SUBSCRIBE 此 channel
    ↓
兩台 Server 都廣播給各自連線的玩家
```

---

## 實作計畫

### Phase 2：Redis 持久化（預計 1 週）

**目標：** Server 重啟後玩家金幣不歸零

#### 步驟 1：建立 Redis Store 模組

```
server/internal/store/
├── redis.go          ← Redis 連線管理
├── player_store.go   ← 玩家狀態 CRUD
├── room_store.go     ← 房間狀態 CRUD
└── leaderboard.go    ← 排行榜操作
```

#### 步驟 2：整合到現有 Player 管理

```go
// player/player.go 加入 Store 介面
type Store interface {
    SavePlayer(p *Player) error
    LoadPlayer(playerID string) (*Player, error)
    DeletePlayer(playerID string) error
}

// 玩家加入時：先從 Redis 讀取，沒有則建立新玩家
// 玩家離開時：儲存到 Redis
// 每 30 秒：定期儲存所有在線玩家狀態
```

#### 步驟 3：環境變數設定

```bash
REDIS_URL=redis://localhost:6379
REDIS_PASSWORD=（可選）
REDIS_DB=0
```

#### 步驟 4：降級策略（Redis 不可用時）

```go
// 如果 Redis 連線失敗，降級到記憶體模式
// 記錄警告日誌，但不中斷服務
if err := store.Connect(); err != nil {
    log.Printf("⚠️  Redis unavailable, using in-memory mode: %v", err)
    store = NewMemoryStore()
}
```

### Phase 3：水平擴展（預計 1 週）

**目標：** 支援多個 Server 實例同時運行

#### 步驟 1：Pub/Sub 廣播

```go
// 每個 Server 啟動時訂閱所有房間的廣播 channel
func (s *Server) subscribeRoomBroadcasts() {
    pubsub := s.redis.Subscribe(ctx, "room:*:broadcast")
    go func() {
        for msg := range pubsub.Channel() {
            // 解析訊息，廣播給本 Server 連線的玩家
            s.hub.BroadcastLocal(msg.Payload)
        }
    }()
}
```

#### 步驟 2：遊戲邏輯主節點選舉

```
問題：多個 Server 都在跑遊戲循環 → 重複計算

解法：Redis 分散式鎖（SETNX）
- 搶到鎖的 Server 成為「主節點」，負責遊戲邏輯
- 其他 Server 成為「轉發節點」，只負責 WebSocket 連線
- 主節點每 5 秒更新鎖的 TTL（心跳）
- 主節點崩潰後，其他節點搶鎖接管
```

```go
// 主節點選舉
func (s *Server) tryBecomeMaster(roomID string) bool {
    key := fmt.Sprintf("master:%s", roomID)
    ok, err := s.redis.SetNX(ctx, key, s.serverID, 10*time.Second).Result()
    return err == nil && ok
}

// 心跳（每 5 秒）
func (s *Server) renewMasterLock(roomID string) {
    key := fmt.Sprintf("master:%s", roomID)
    s.redis.Expire(ctx, key, 10*time.Second)
}
```

#### 步驟 3：Nginx 負載均衡設定

```nginx
upstream chiikawa_servers {
    # WebSocket 需要 sticky session（同一玩家連到同一 Server）
    ip_hash;
    server 127.0.0.1:7777;
    server 127.0.0.1:7778;
    server 127.0.0.1:7779;
}

server {
    listen 443 ssl;
    location /ws {
        proxy_pass http://chiikawa_servers;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
```

---

## 效能預估

| 指標 | 單 Server | 3 Server + Redis |
|------|-----------|-----------------|
| 最大玩家數 | ~50 | ~150 |
| Server 重啟恢復時間 | 0（狀態丟失） | < 1 秒 |
| 廣播延遲 | < 1ms | < 5ms（Pub/Sub 延遲） |
| Redis 記憶體使用 | - | ~10MB（1000 玩家） |
| 月費（Redis Cloud 免費版） | $0 | $0（30MB 免費） |

---

## 依賴套件

```go
// go.mod 新增
require (
    github.com/redis/go-redis/v9 v9.x.x
)
```

**安裝：**
```bash
go get github.com/redis/go-redis/v9
```

**本機開發（Docker）：**
```bash
docker run -d -p 6379:6379 redis:7-alpine
```

**Windows 本機（無 Docker）：**
```
下載 Redis for Windows：https://github.com/tporadowski/redis/releases
解壓縮後執行 redis-server.exe
```

---

## 實作優先級

| 功能 | 優先級 | 預計工時 | 說明 |
|------|--------|---------|------|
| Redis 連線管理 | P1 | 2h | 基礎設施 |
| 玩家狀態持久化 | P1 | 3h | 最高價值 |
| 排行榜 Redis 化 | P2 | 2h | 跨 Server 共享 |
| Pub/Sub 廣播 | P2 | 4h | 多 Server 必要 |
| 主節點選舉 | P3 | 4h | 複雜，後期再做 |
| Nginx 負載均衡 | P3 | 2h | 需要多台機器 |

**建議：** 先做 P1（玩家狀態持久化），讓 Prototype 更穩定，再考慮 P2/P3。

---

## 降級策略（重要）

Redis 不是必要依賴，Server 必須在沒有 Redis 的情況下正常運行：

```go
type Store interface {
    SavePlayer(p *player.Player) error
    LoadPlayer(playerID string) (*player.Player, error)
}

// 兩種實作：
// 1. RedisStore（有 Redis 時使用）
// 2. MemoryStore（沒有 Redis 時降級）

// 啟動時自動選擇：
func NewStore(redisURL string) Store {
    if redisURL == "" {
        log.Println("⚠️  REDIS_URL not set, using in-memory store")
        return NewMemoryStore()
    }
    store, err := NewRedisStore(redisURL)
    if err != nil {
        log.Printf("⚠️  Redis connection failed: %v, using in-memory store", err)
        return NewMemoryStore()
    }
    return store
}
```

---

## 參考資料

- [Redis Go Client (go-redis)](https://redis.uptrace.dev/)
- [Redis Pub/Sub for Real-time Games](https://redis.io/docs/manual/pubsub/)
- [Distributed Locking with Redis](https://redis.io/docs/manual/patterns/distributed-locks/)
- [WebSocket Sticky Sessions with Nginx](https://nginx.org/en/docs/http/ngx_http_upstream_module.html#ip_hash)

---

*最後更新：2026-05-19（DAY-026）*
