# Nightly Report — DAY-055

**日期**：2026-05-20  
**執行者**：Game Director（自主觸發）  
**狀態**：✅ 完成

---

## 今日完成事項

### 觀戰者系統完整實作
- `hub.go`：新增 `BroadcastToPlayers()` — 只廣播給 RolePlayer，跳過 RoleSpectator
- `main.go`：觀戰者連線時廣播 `spectator_join` 給所有玩家（spectator_id + spectator_count）
- `GameManager.gd`：新增 `spectator_joined` 訊號 + `_handle_spectator_join()` 處理函數
- `HUD.gd`：連接 `spectator_joined` 訊號，用成就通知系統顯示「👁️ 有人在觀戰！」
- `hub_test.go`：新增 `TestBroadcastToPlayers`（確認只有 player 收到，spectator 不收）

### 觀戰者離開通知（DAY-055b）
- `hub.go`：新增 `OnSpectatorDisconnect` 回調欄位
- `main.go`：設定 `OnSpectatorDisconnect` 廣播 `spectator_leave` 訊息
- `GameManager.gd`：新增 `spectator_left` 訊號 + `_handle_spectator_leave()` 處理函數
- `HUD.gd`：連接 `spectator_left` 訊號，顯示「👁️ 觀戰者離開了」通知

### TopBar 觀戰者計數（DAY-055c）
- `HUD.gd`：TopBar 加入觀戰者計數標籤（👁️ 0）
- `_update_spectator_count_label()` 函數

### Grafana 升級（DAY-055d）
- Grafana dashboard 從 25 個面板升級到 26 個面板（觀戰者計數 stat）
- `tools/add_spectator_panel.py`：自動化工具

### KnowHow 更新
- KnowHow #108：Hub 回調擴充模式（OnSpectatorDisconnect）
- KnowHow #109：觀戰者離開通知的 UX 設計

---

## 品質分數

| 指標 | 分數 | 狀態 |
|------|------|------|
| Build Stability | 100 | ✅ |
| Visual Consistency | 100 | ✅ |
| Balance Health | 96 | ✅ |
| Animation Quality | 100 | ✅ |
| Audio Sync | 100 | ✅ |
| Gameplay Feel | 100 | ✅ |
| Spec Completeness | 100 | ✅ |
| Regression Risk | 5 | ✅ |

---

## 明日計畫
- goleak goroutine 洩漏偵測整合
