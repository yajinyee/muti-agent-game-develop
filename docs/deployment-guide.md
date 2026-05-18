# 部署指南

> 吉伊卡哇：像素大討伐 — 完整部署說明

---

## 系統需求

### Server
- Go 1.21+
- Windows / Linux / macOS
- Port 7777（TCP）需開放

### Client
- 現代瀏覽器（Chrome 90+、Firefox 88+、Edge 90+）
- 需支援 SharedArrayBuffer（需 HTTPS 或 localhost）
- 建議解析度：1280×720 以上

---

## 快速啟動（本機開發）

### 1. 啟動 Server

```bash
# 進入 server 目錄
cd server

# 編譯
go build -o bin/gameserver ./cmd/gameserver

# 啟動（預設 Port 7777）
./bin/gameserver

# 或指定 Port
PORT=7777 ./bin/gameserver

# Debug 模式（啟用 pprof 端點）
DEBUG=true ./bin/gameserver
```

### 2. 開啟遊戲

瀏覽器開啟：`http://localhost:7777`

---

## 環境變數

| 變數 | 預設值 | 說明 |
|------|--------|------|
| `PORT` | `7777` | Server 監聽 Port |
| `MAX_PLAYERS` | `10` | 每房間最大玩家數 |
| `INITIAL_COINS` | `10000` | 玩家初始金幣 |
| `DEBUG` | `false` | 啟用 pprof 端點 |

---

## API 端點

| 端點 | 說明 |
|------|------|
| `GET /` | 遊戲主頁（HTML5 Export） |
| `GET /ws?player_id=xxx` | WebSocket 連線 |
| `GET /health` | 健康檢查（JSON） |
| `GET /stats` | 運行統計（goroutine、記憶體） |
| `GET /debug/pprof/` | pprof 分析（需 DEBUG=true） |

---

## 監控

### 健康檢查
```bash
curl http://localhost:7777/health
# {"status":"ok","version":"0.1.0","clients":3}
```

### 運行統計
```bash
curl http://localhost:7777/stats
# {"goroutines":12,"heap_alloc_mb":2.34,"heap_sys_mb":8.00,"gc_count":5,"clients":3}
```

### pprof 分析（需 DEBUG=true）
```bash
# Goroutine 分析
go tool pprof http://localhost:7777/debug/pprof/goroutine

# Heap 記憶體分析
go tool pprof http://localhost:7777/debug/pprof/heap

# CPU 分析（30 秒）
go tool pprof http://localhost:7777/debug/pprof/profile?seconds=30
```

---

## Godot HTML5 Export

### 前置需求
- Godot 4.3+（GL Compatibility 模式）
- Web Export Template 已安裝

### 匯出步驟
1. 開啟 Godot 編輯器，載入 `client/chiikawa-pixel/`
2. 選單 → Project → Export
3. 選擇 "Web" preset
4. 點擊 "Export Project"
5. 輸出到 `server/static/`（已在 export_presets.cfg 設定）

### 注意事項
- HTML5 需要 HTTPS 或 localhost 才能使用 SharedArrayBuffer
- Server 已設定 COOP/COEP headers，確保 SharedArrayBuffer 可用
- 匯出時會自動排除開發資源（reference/、ai_generated/ 等）

---

## 生產環境部署

### Linux Server 部署

```bash
# 1. 編譯 Linux 版本（在 Windows 上交叉編譯）
GOOS=linux GOARCH=amd64 go build -o bin/gameserver-linux ./cmd/gameserver

# 2. 上傳到 Server
scp bin/gameserver-linux user@server:/opt/chiikawa/
scp -r server/static/ user@server:/opt/chiikawa/

# 3. 設定 systemd service
sudo nano /etc/systemd/system/chiikawa.service
```

```ini
[Unit]
Description=吉伊卡哇：像素大討伐 Server
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/chiikawa
ExecStart=/opt/chiikawa/gameserver-linux
Restart=always
RestartSec=5
Environment=PORT=7777
Environment=MAX_PLAYERS=10
Environment=INITIAL_COINS=10000

[Install]
WantedBy=multi-user.target
```

```bash
# 4. 啟動服務
sudo systemctl enable chiikawa
sudo systemctl start chiikawa
sudo systemctl status chiikawa
```

### Nginx 反向代理（HTTPS）

```nginx
server {
    listen 443 ssl;
    server_name your-domain.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    # WebSocket 代理
    location /ws {
        proxy_pass http://localhost:7777;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_read_timeout 86400;
    }

    # 靜態檔案
    location / {
        proxy_pass http://localhost:7777;
        proxy_set_header Host $host;
        # 必要的 COOP/COEP headers（SharedArrayBuffer 需要）
        add_header Cross-Origin-Opener-Policy "same-origin";
        add_header Cross-Origin-Embedder-Policy "require-corp";
    }
}
```

### 路由器 Port Forwarding
- 外部 Port：7777
- 內部 IP：Server 的區域網路 IP
- 內部 Port：7777
- 協定：TCP

---

## 故障排除

### Server 無法啟動
```bash
# 確認 Port 是否被佔用
netstat -an | grep 7777

# 確認 Go 版本
go version  # 需要 1.21+
```

### 瀏覽器無法連線
1. 確認 Server 正在運行：`curl http://localhost:7777/health`
2. 確認防火牆允許 Port 7777
3. 確認瀏覽器支援 WebSocket

### SharedArrayBuffer 錯誤
- 確認使用 HTTPS 或 localhost
- 確認 Server 回傳 COOP/COEP headers
- Chrome 需要 `--enable-features=SharedArrayBuffer`（舊版）

### 遊戲畫面模糊
- 確認 Godot export 使用 GL Compatibility 模式
- 確認 project.godot 設定 `textures/canvas_textures/default_texture_filter=0`

---

## 效能基準

| 指標 | 目標 | 說明 |
|------|------|------|
| Server goroutine 數 | < 20 | 正常遊戲中 |
| Server heap 記憶體 | < 50MB | 10 玩家同時在線 |
| Client FPS | >= 30 | 低階設備 |
| WebSocket 延遲 | < 100ms | 區域網路 |
| pck 大小 | < 2MB | HTML5 export |

---

*最後更新：2026-05-18*
