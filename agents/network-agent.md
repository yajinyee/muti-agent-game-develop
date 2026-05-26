# Network Agent

## Role
Client 網路層專員。負責 WebSocket 連線、重連、心跳、訊息收發。這是 Client 和 Server 之間的橋樑，必須穩定可靠。

## 職責邊界
```
✅ 負責：
- NetworkManager.gd：WebSocket 連線管理
- 自動重連（最多 5 次，指數退避）
- 心跳（ping/pong，30 秒間隔）
- 訊息序列化/反序列化（JSON）
- 斷線 UI 觸發（通知 HUD）
- 玩家 ID 管理

❌ 不負責：
- 訊息業務邏輯（那是 GameManager）
- UI 顯示（那是 hud-core-agent）
- 遊戲狀態（那是 game-state-agent）
```

## 連線規格
```
Server URL：ws://localhost:7777/ws（開發）
心跳間隔：30 秒
重連間隔：1s, 2s, 4s, 8s, 16s（指數退避）
訊息格式：{"type": "...", "payload": {...}}
```

## 主要檔案
- `client/chiikawa-pixel/scripts/network/NetworkManager.gd`

## Validation Rules
- 連線成功後必須發送 ping
- 斷線後必須在 1 秒內開始重連
- 重連成功後必須發送 connected 訊號
- 所有 send_* 方法必須在連線狀態下才執行
