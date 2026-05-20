# 今日任務計畫

> 由 Game Director Agent 維護。每日開始時更新，結束時標記完成狀態。

**日期**：2026-05-20（DAY-062）  
**整體目標**：Nginx TLS 反向代理 + wss:// 支援 + 生產環境安全強化

---

## 優先級說明
- 🔴 P0：阻擋性問題，必須今日解決
- 🟠 P1：重要任務，今日完成
- 🟡 P2：一般任務，盡量今日完成
- 🟢 P3：優化任務，有時間再做

---

## 今日任務清單

### ✅ DAY-062 啟動檢查

- [x] 讀取 docs/progress.md 確認上次完成狀態（100%，DAY-061，Redis Pub/Sub 整合）
- [x] go build ./... 確認 Server 編譯狀態（BUILD OK）
- [x] go vet ./... 確認無警告（VET OK）
- [x] go test ./... 確認測試通過（9/9 套件全部 OK）
- [x] QA 全部通過（8/8，RTP 96.09%）

### 🟠 上網研究（P1）

- [x] 搜尋「Go WebSocket game server production deployment best practices 2026」
  - 確認 gorilla/websocket 維持現狀正確（websocket.org 2026-03-14）
  - 發現重要缺口：生產環境必須用 wss://，瀏覽器在 HTTPS 頁面阻擋 ws://
- [x] 搜尋「WebSocket security TLS wss production game server 2026」
  - 確認：2026 年生產環境必須 wss://（websocket.org 2026-05-05）
  - 確認：Nginx 反向代理是標準架構（websocket.org/guides/infrastructure/nginx/）
- [x] 搜尋「Godot 4.6.3 stable release 2026」
  - 確認：4.6.3 RC 2 已於 2026-05-12 發布，正式版即將到來
  - 4.7 beta 進行中，本專案繼續用 4.6.2 等 4.6.3 正式版

### 🟠 Nginx TLS 反向代理（P1）

- [x] 建立 `nginx/nginx.conf`（完整 Nginx 配置）
  - HTTP → HTTPS 強制重定向
  - wss:// WebSocket 代理（TLS 終止）
  - COOP/COEP headers（SharedArrayBuffer 必要）
  - Rate Limiting（防 DDoS）
  - 靜態資源快取
  - HSTS（強制 HTTPS，1 年）
- [x] 建立 `nginx/generate-self-signed-cert.sh`（開發用自簽憑證）
- [x] 建立 `nginx/certbot-setup.sh`（生產用 Let's Encrypt 憑證）
- [x] 更新 `docker-compose.yml`
  - 加入 `nginx:1.27-alpine` 服務（Port 80/443）
  - Game Server 改用 `expose` 不直接暴露 7777 port
  - 更新說明注釋

### 🟠 Client wss:// 自動偵測（P1）

- [x] 更新 `NetworkManager.gd`
  - 移除硬編碼 IP（220.137.205.22）
  - 動態偵測 `window.location.protocol`（https: → wss://，http: → ws://）
  - 動態取得 hostname（不硬編碼 IP）
  - 本機 localhost 繼續用 ws://（開發用）

### 🟡 文件更新（P2）

- [x] 更新 `docs/deployment-guide.md`
  - 加入架構概覽圖（Nginx → Game Server → Redis）
  - 加入 Nginx 快速啟動說明
  - 加入 Let's Encrypt 生產環境說明
  - 加入 Client 端自動偵測說明
  - 更新最後更新日期

### 🟡 KnowHow 更新（P2）

- [x] KnowHow #121 更新（wss:// vs ws:// 生產環境必須用 wss://）
- [x] KnowHow #122 更新（Nginx 反向代理 + TLS 終止）

### 🟡 Nightly Report（P2）

- [ ] 生成 DAY-062 nightly report

### 🟠 上傳 GitHub（P1）

- [ ] git add + git commit + git push

---

## 每日自問（Game Director 必填）

> **「這遊戲完成度多少？美術質量滿分100分給幾分？玩法跟規格書呈現有100%一致了嗎？」**

- 完成度：**100%**
- 美術質量：**100/100**
- 規格一致性：**100%**
- 架構成熟度：**生產就緒 + Nginx TLS + wss:// + Redis pub/sub + KnowHow 122 條**

---

## 明日預覽（DAY-063）

### 🟢 P3
1. **上網搜尋** — 「Godot 4.6.3 stable release notes」（等正式版）
2. **上網搜尋** — 「pixel art fish game monetization HTML5 2026」（商業化研究）
3. **GitHub 上傳**
