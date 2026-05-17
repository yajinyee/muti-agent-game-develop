# Research Notes：Go WebSocket 優化（2025）

**研究日期**：2026-05-17  
**研究者**：Research Agent  
**來源**：leapcell.io, hemaks.org, gorilla/websocket docs

---

## 關鍵發現

### 1. gorilla/websocket 仍是 2025 年最佳選擇
- 18,000+ GitHub stars
- 生產環境廣泛使用
- 支援 permessage-deflate 壓縮（需手動啟用）

### 2. 現有 Hub 架構已符合最佳實踐
現有 `server/internal/ws/hub.go` 的 Hub 模式是正確的。

### 3. 可改善的點

**A. 廣播非阻塞（防止慢客戶端）**
```go
// 現有可能的問題：如果 client.send 滿了會阻塞
// 改善：加入 default case
select {
case client.send <- message:
default:
    close(client.send)
    delete(h.clients, client)
}
```

**B. WebSocket 壓縮（減少頻寬）**
```go
upgrader := websocket.Upgrader{
    EnableCompression: true,  // 啟用 permessage-deflate
}
```

**C. 讀寫 Buffer 大小**
```go
upgrader := websocket.Upgrader{
    ReadBufferSize:  4096,  // 遊戲訊息通常 < 1KB，但 4KB 更安全
    WriteBufferSize: 4096,
}
```

### 4. 遊戲 Server 特殊考量
- 遊戲訊息頻率高（10 FPS = 每秒 10 個廣播）
- 每個廣播需要發送給所有玩家
- 建議：廣播 channel 緩衝設為 100+

## 授權
- gorilla/websocket：BSD-2-Clause（可商用）

## 行動建議
- 低優先：加入 EnableCompression
- 中優先：確認廣播有非阻塞 select
- 低優先：調整 Buffer 大小
