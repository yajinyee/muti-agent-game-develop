# 多人房間架構設計文件

> 版本：v1.0  
> 日期：2026-05-19（DAY-019）  
> 狀態：設計草稿（未來功能，當前 Prototype 為單房間）

---

## 概覽

當前 Prototype 是單房間架構（room-001 硬編碼）。  
本文件規劃從單房間升級到多房間的架構設計，支援多個獨立遊戲房間同時運行。

---

## 當前架構（單房間）

```
Client → WebSocket → Hub → Game（room-001）
                      ↓
                   Analytics
```

**限制：**
- 所有玩家共用同一個遊戲狀態
- 無法隔離不同投注等級的玩家
- 無法支援不同主題的遊戲房間

---

## 目標架構（多房間）

```
Client → WebSocket → Hub → RoomManager
                              ↓
                    ┌─────────┼─────────┐
                    ↓         ↓         ↓
                 Room-001  Room-002  Room-003
                 (低投注)  (高投注)  (VIP)
                    ↓         ↓         ↓
                 Game      Game      Game
                    ↓
                 Analytics（每房間獨立）
```

---

## 核心元件設計

### 1. RoomManager（新增）

```go
// internal/room/manager.go
type RoomManager struct {
    rooms   map[string]*Room
    mu      sync.RWMutex
    maxRooms int
}

type Room struct {
    ID       string
    Name     string
    Config   RoomConfig
    Game     *game.Game
    Players  map[string]bool
    CreatedAt time.Time
}

type RoomConfig struct {
    MinBetLevel  int     // 最低投注等級
    MaxBetLevel  int     // 最高投注等級
    MaxPlayers   int     // 最大玩家數（建議 8-16）
    Theme        string  // 主題（chiikawa/default）
    RTPTarget    float64 // 目標 RTP（0.92-0.96）
}
```

### 2. Hub 升級（支援房間路由）

```go
// 現有 Hub 加入房間概念
type Hub struct {
    clients    map[string]*Client
    rooms      map[string]map[string]bool  // roomID → clientIDs
    mu         sync.RWMutex
    // ... 其他欄位不變
}

// 廣播到特定房間
func (h *Hub) BroadcastToRoom(roomID string, msg *Message) {
    h.mu.RLock()
    clients := h.rooms[roomID]
    h.mu.RUnlock()
    for clientID := range clients {
        h.Send(clientID, msg)
    }
}
```

### 3. WebSocket 連線帶房間 ID

```
ws://localhost:7777/ws?player_id=xxx&room_id=room-001
```

- `room_id` 為空時：自動分配到人數最少的房間
- `room_id` 指定時：加入指定房間（滿員則拒絕）

### 4. HTTP 房間管理 API

```
GET  /rooms              → 列出所有房間（ID、名稱、人數、狀態）
POST /rooms              → 建立新房間（需要管理員 token）
GET  /rooms/{id}         → 查詢特定房間狀態
DELETE /rooms/{id}       → 關閉房間（踢出所有玩家）
GET  /rooms/{id}/leaderboard → 房間排行榜
```

---

## 實作計畫

### Phase 1：最小可行多房間（2 週）

1. **建立 RoomManager**（`internal/room/manager.go`）
   - 建立/刪除房間
   - 玩家加入/離開房間
   - 自動分配房間（負載均衡）

2. **升級 Hub**（`internal/ws/hub.go`）
   - 加入 `rooms` map
   - `BroadcastToRoom()` 方法
   - `JoinRoom()` / `LeaveRoom()` 方法

3. **升級 main.go**
   - WebSocket 連線解析 `room_id` 參數
   - 加入 `/rooms` HTTP API
   - 預建 3 個預設房間（低/中/高投注）

4. **升級 Client（Godot）**
   - 連線 URL 加入 `room_id` 參數
   - 房間選擇 UI（大廳畫面）

### Phase 2：房間持久化（1 週）

- 房間狀態寫入 Redis（支援 Server 重啟後恢復）
- 玩家 session 跨房間保留金幣

### Phase 3：動態房間（1 週）

- 玩家可以建立私人房間
- 房間密碼保護
- 觀戰模式（只讀 WebSocket）

---

## 效能考量

| 指標 | 單房間 | 多房間（10個） |
|------|--------|---------------|
| Goroutine 數量 | ~15 | ~150 |
| 記憶體使用 | ~50MB | ~200MB |
| WebSocket 連線 | 無限制 | 每房間 16 人上限 |
| 廣播延遲 | <1ms | <1ms（房間隔離） |

**結論：** 10 個房間的資源消耗在合理範圍內，GTX 1650 的機器可以輕鬆支撐。

---

## 遷移策略

1. **向後相容**：`room_id` 為空時預設 `room-001`，現有 Client 無需修改
2. **漸進式**：先在 Server 端實作，Client 端可以後續升級
3. **測試**：用 `stress_test.py` 模擬多房間並發測試

---

## 參考資料

- [Go WebSocket Multi-Room Architecture](https://oneuptime.com/blog/post/2026-01-25-websocket-chat-app-rooms-go/view)
- [Scalable WebSocket Architecture](https://dev.to/hpx7/scalable-websocket-architecture-4e1n)
- [Go Game Server Best Practices 2025](https://generalistprogrammer.com/tutorials/go-game-development-complete-server-side-guide-2025)

Content was rephrased for compliance with licensing restrictions.
