# WebSocket 協定變更政策

> 本文件規範所有 WebSocket 協定的變更流程。協定是 Go Server 與 Godot Client 之間的契約，任何變更都必須嚴格遵守本政策，否則可能導致通訊中斷。

---

## 現行協定版本

**版本**：1.0.0  
**連線端點**：`ws://[server]:7777/ws`  
**編碼**：UTF-8 JSON  
**心跳間隔**：30 秒

---

## 訊息格式規範

### 基礎結構
```json
{
  "type": "string",      // 必填：訊息類型
  "version": "1.0",      // 必填：協定版本
  "timestamp": 1234567890, // 必填：Unix 時間戳（毫秒）
  "data": {}             // 選填：訊息內容
}
```

### Client → Server 訊息類型

| type | 說明 | data 結構 |
|------|------|---------|
| `shoot` | 玩家射擊 | `{"target_id": "T001", "bet": 1, "x": 100, "y": 200}` |
| `ping` | 心跳 | `{}` |
| `bet_change` | 更改下注 | `{"bet": 5}` |
| `bonus_select` | Bonus 選擇 | `{"choice": 1}` |

### Server → Client 訊息類型

| type | 說明 | data 結構 |
|------|------|---------|
| `hit` | 命中 | `{"target_id": "T001", "reward": 5, "multiplier": 1}` |
| `miss` | 未命中 | `{"target_id": "T001"}` |
| `target_spawn` | 目標物生成 | `{"targets": [...]}` |
| `target_remove` | 目標物移除 | `{"target_id": "T001"}` |
| `boss_spawn` | BOSS 出現 | `{"boss_id": "B001", "hp": 1000, "duration": 30}` |
| `boss_hit` | BOSS 受傷 | `{"boss_id": "B001", "hp_remaining": 800, "reward": 10}` |
| `boss_defeat` | BOSS 擊敗 | `{"boss_id": "B001", "total_reward": 500}` |
| `bonus_trigger` | Bonus 觸發 | `{"bonus_type": "free_shot", "count": 10}` |
| `bonus_end` | Bonus 結束 | `{"total_reward": 100}` |
| `score_update` | 分數更新 | `{"score": 1000, "balance": 5000}` |
| `pong` | 心跳回應 | `{}` |
| `error` | 錯誤 | `{"code": 1001, "message": "..."}` |

---

## 錯誤碼定義

| 錯誤碼 | 說明 | 處理方式 |
|-------|------|---------|
| 1001 | 無效的訊息格式 | Client 記錄並忽略 |
| 1002 | 目標物不存在 | Client 移除該目標物 |
| 1003 | 餘額不足 | Client 顯示提示 |
| 1004 | 下注金額無效 | Client 重置下注 |
| 2001 | 伺服器內部錯誤 | Client 顯示錯誤並嘗試重連 |
| 2002 | 遊戲狀態異常 | Client 重新同步狀態 |
| 3001 | 連線超時 | Client 自動重連 |

---

## 協定版本管理

### 版本號規則（SemVer）
- **Major**（X.0.0）：破壞性變更（移除欄位、改變語義）
- **Minor**（1.X.0）：向後相容的新增（新增訊息類型、新增選填欄位）
- **Patch**（1.0.X）：修復（錯誤碼說明更新、文件修正）

### 版本相容性規則
- Server 必須支援當前 Major 版本的所有 Minor 版本
- Client 必須能處理未知的訊息類型（忽略而非崩潰）
- 廢棄的訊息類型必須保留至少一個 Major 版本

---

## 變更流程

### 步驟 1：提案
任何 Agent 可以提出協定變更需求，填寫以下表單：

```markdown
## 協定變更提案

**提案者**：[Agent 名稱]
**日期**：[DATE]
**變更類型**：Major / Minor / Patch
**影響範圍**：Server / Client / 雙側

### 變更內容
[詳細描述要新增/修改/移除的內容]

### 變更理由
[為什麼需要這個變更]

### 影響評估
- Server 端影響：[說明]
- Client 端影響：[說明]
- 向後相容性：[是/否，若否說明原因]

### 測試計畫
[如何驗證變更正確]
```

### 步驟 2：Spec Architect 審核
- 評估技術可行性
- 確認向後相容性
- 更新協定文件草稿

### 步驟 3：Game Director 批准
- 評估對遊戲體驗的影響
- 確認優先級
- 批准或否決

### 步驟 4：雙側同步實作
- Go Server Agent 更新 Server 端
- Godot Client Agent 更新 Client 端
- **必須同時完成，不得只更新一側**

### 步驟 5：QA 驗證
- 執行完整的 WebSocket 通訊測試
- 確認所有訊息類型正常運作
- 確認錯誤處理正確

### 步驟 6：文件更新
- Spec Architect 更新本文件
- 更新版本號
- 記錄變更歷史

---

## 禁止行為

1. **禁止單側更新**：不得只更新 Server 或只更新 Client
2. **禁止刪除現有欄位**（Minor 版本）：只能新增選填欄位
3. **禁止改變現有欄位語義**（Minor 版本）：語義變更必須升 Major
4. **禁止在未測試的情況下部署協定變更**
5. **禁止繞過審核流程**（緊急修復除外，但事後必須補文件）

---

## 緊急修復程序

若發現協定 Bug 導致遊戲無法運作：

1. Go Server Agent 立即修復並部署
2. Godot Client Agent 同步修復
3. 事後 24 小時內補齊文件
4. Spec Architect 更新版本號（Patch）

---

## 變更歷史

| 版本 | 日期 | 變更內容 | 提案者 |
|------|------|---------|-------|
| 1.0.0 | 2025-01-01 | 初始協定定義 | Spec Architect |
