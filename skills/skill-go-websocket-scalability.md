# Skill：Go WebSocket 擴展性最佳實踐

> 來源：[leapcell.io - Building a Scalable Go WebSocket Service](https://www.leapcell.io/blog/building-a-scalable-go-websocket-service-for-thousands-of-concurrent-connections)  
> 研究日期：2026-05-17（DAY-002）  
> 記錄者：Research Agent

---

## 核心架構：Hub + Client + Channel

### Client 結構
```go
type Client struct {
    conn *websocket.Conn
    send chan []byte  // 緩衝 channel，防止慢消費者阻塞
}
```

### Hub 結構（集中管理所有連線）
```go
type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
}
```

### 關鍵設計原則

1. **每個連線一個 goroutine**：readPump + writePump 分離
2. **緩衝 channel**：`make(chan []byte, 256)` 防止慢客戶端阻塞廣播
3. **非阻塞廣播**：用 `select { case client.send <- msg: default: close(client.send) }` 踢掉慢客戶端
4. **Ping/Pong 心跳**：每 10 秒發送 Ping 保持連線

### 廣播優化（防止慢客戶端影響其他人）
```go
case message := <-h.broadcast:
    for client := range h.clients {
        select {
        case client.send <- message:
            // 成功
        default:
            // send channel 滿了，踢掉這個客戶端
            close(client.send)
            delete(h.clients, client)
        }
    }
```

## 擴展性優化

| 優化項目 | 方法 | 效果 |
|---------|------|------|
| 訊息序列化 | Protocol Buffers 替代 JSON | 減少 30-50% 訊息大小 |
| 水平擴展 | 多 Server + Redis PubSub | 支援百萬連線 |
| 緩衝大小 | `ReadBufferSize: 1024, WriteBufferSize: 1024` | 減少記憶體分配 |
| 優雅關閉 | signal.Notify + context.WithCancel | 避免連線中斷 |

## 與現有專案的關聯

現有 `server/internal/ws/hub.go` 已實作 Hub 模式，但可以改善：
- 增加 `ReadBufferSize/WriteBufferSize` 設定
- 廣播時加入非阻塞 select（防止慢客戶端）
- 加入 Ping/Pong 心跳機制

## 相關檔案
- `server/internal/ws/hub.go`
- `server/internal/ws/protocol.go`

*Content was rephrased for compliance with licensing restrictions*
