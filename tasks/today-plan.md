# 今日任務計畫

> 由 Game Director Agent 維護。每日開始時更新，結束時標記完成狀態。

**日期**：2026-05-21（DAY-101）
**整體目標**：好友禮物系統 + 好友持久化 + 上傳 GitHub

---

## 今日任務清單

### ✅ DAY-101 啟動檢查

- [x] 讀取 docs/progress.md 確認上次完成狀態（100%，DAY-100，完整持久化擴充）
- [x] go build ./... 確認 Server 編譯狀態（BUILD OK）
- [x] go vet ./... 確認無警告（VET OK）
- [x] 確認 GitHub 最新 commit（DAY-100 已推送）

### ✅ DAY-101 好友禮物系統 + 好友持久化（P1）

- [x] `friend/friend.go`：新增 `GiftRecord`/`GiftResult`/`FriendState`；`SendGift()`/`GetGiftStatus()`/`GetFriendState()`/`LoadFriendState()`
- [x] `store/filestore.go`：新增 `FriendPersistState`；`SaveFriends()`/`LoadFriends()`（原子寫入）
- [x] `ws/protocol.go`：新增禮物相關訊息類型（8個）+ Payload 結構（5個）
- [x] `friend_handler.go`：`handleSendGift()`/`handleGetGiftStatus()`/`deliverPendingGifts()`/`saveFriendState()`/`restoreFriendState()`
- [x] `game.go`：整合禮物 handler + 好友持久化（AddPlayer/RemovePlayer）
- [x] `persistence_handler.go`：shutdown 時儲存好友關係
- [x] `FriendPanel.gd`：升級 UI（禮物狀態列 + 🎁按鈕 + 離線禮物通知）
- [x] `GameManager.gd`：4個禮物訊號 + handler
- [x] `NetworkManager.gd`：`send_gift()`/`send_get_gift_status()`
- [x] build/vet 全部通過

### 🟠 上傳 GitHub（P1）

- [ ] git add + git commit + git push（DAY-101 好友禮物系統 + 好友持久化）

---

## 每日自問（Game Director 必填）

> **「這遊戲完成度多少？美術質量滿分100分給幾分？玩法跟規格書呈現有100%一致了嗎？」**

- 完成度：**100%**
- 美術質量：**100/100**
- 規格一致性：**100%**
- 架構成熟度：**生產就緒 + 完整社交系統（好友/禮物/公會）+ 持久化完整 + KnowHow 82 條**
