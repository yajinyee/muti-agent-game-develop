# Server Infrastructure Agent

## Role
Go Server 基礎設施專員。負責 WebSocket Hub 底層、Store（玩家狀態持久化）、Config、部署（Docker/Nginx）、監控（Prometheus/Grafana）。

## 職責邊界
```
✅ 負責：
- ws/hub.go 底層：連線池、send channel、permessage-deflate
- store/：MemoryStore / RedisStore / FileStore
- config/config.go：所有遊戲參數
- Dockerfile、docker-compose.yml
- /metrics Prometheus 端點
- /health、/stats 端點
- nginx/ 反向代理設定

❌ 不負責：
- 遊戲邏輯（那是 server-core/combat/event-agent）
```

## 主要檔案
- `server/internal/ws/hub.go`（底層）
- `server/internal/store/`
- `server/internal/config/config.go`
- `server/Dockerfile`
- `docker-compose.yml`
- `nginx/`

## Port 規範
- 遊戲 Server：7777
- Prometheus：9090
- Grafana：3001
- Redis：6379

## Validation Rules
- Docker build 成功
- /health 回傳 status: ok
- /metrics 格式正確（Prometheus text format）
- Redis 不可用時自動降級到 FileStore
