# QA 自動化測試規格

> 版本：1.0.0  
> 維護者：QA Playtest Agent  
> 最後更新：2026-05-17

---

## 概覽

本規格定義吉伊卡哇：像素大討伐的自動化 QA 測試流程，確保每次 Build 都達到品質門檻。

---

## 自動測試流程

```
Step 1: Start Server
    └─ go build ./...（確認編譯）
    └─ 啟動 server（Port 7777）
    └─ 等待 Ready 信號（最多 10 秒）

Step 2: Open HTML5
    └─ 啟動 HTTP 靜態伺服器（Port 8080）
    └─ 開啟 Chrome headless（或 Playwright）
    └─ 載入 index.html

Step 3: Connect WS
    └─ 確認 WebSocket 連線建立
    └─ 確認 handshake 完成
    └─ 確認心跳正常（30 秒週期）

Step 4: Play
    └─ 模擬玩家操作（自動射擊 1000 局）
    └─ 記錄所有事件（攻擊/命中/擊殺/獎勵）
    └─ 觸發 BOSS 戰（至少 1 次）
    └─ 觸發 Bonus 遊戲（至少 1 次）

Step 5: Capture
    └─ 截圖關鍵畫面（idle/attack/boss/bonus/bigwin）
    └─ 錄製 30 秒遊戲影片
    └─ 記錄 Console 錯誤

Step 6: Analyze
    └─ 分析截圖（sprite 品質、UI 完整性）
    └─ 分析 RTP（實際 vs 理論）
    └─ 分析音效觸發時機
    └─ 計算各項品質分數

Step 7: Report
    └─ 生成 QA 報告（reports/qa/qa-report-YYYY-MM-DD.md）
    └─ 更新 memory/project-memory.md 品質分數
    └─ 如有嚴重問題，觸發警報
```

---

## 必檢問題清單

### 🔴 P0 — 崩潰類（任何一項失敗 = Build Unstable）

| 檢查項目 | 測試方法 | 通過條件 |
|---------|---------|---------|
| crash | 執行 1000 局，監控 crash | 0 次 crash |
| 連線穩定性 | 連線 30 分鐘，監控斷線 | 斷線次數 = 0 |
| 記憶體洩漏 | 執行 1000 局後檢查記憶體 | 增長 < 50MB |
| Go Server 崩潰 | 監控 Server 進程 | 進程存活 100% |

### 🟠 P1 — 功能類（任何一項 < 90% = Build Warning）

| 檢查項目 | 測試方法 | 通過條件 |
|---------|---------|---------|
| 攻擊功能 | 發射 100 顆子彈 | 命中率符合設計 |
| 怪物生成 | 觀察 100 秒 | 生成頻率符合設計 |
| 獎勵計算 | 驗證 100 次擊殺 | 獎勵金額 100% 正確 |
| 動畫播放 | 觸發所有動畫狀態 | 無卡幀/跳幀 |
| 音效觸發 | 觸發所有音效事件 | 觸發率 100% |
| UI 顯示 | 截圖比對 | 無 UI 錯位/遮擋 |

### 🟡 P2 — 特殊功能類

| 檢查項目 | 測試方法 | 通過條件 |
|---------|---------|---------|
| BOSS 戰 | 觸發 BOSS 完整流程 | 全流程無錯誤 |
| Bonus 遊戲 | 觸發 Bonus 完整流程 | 全流程無錯誤 |
| RTP 平衡 | 模擬 10000 局 | RTP 在 92-96% |

---

## 品質分數計算方式

### Build Stability（建置穩定性）
```
Build Stability = 100
- 如果 go build 失敗：-50
- 如果 go vet 有警告：-5 per warning
- 如果 Server 啟動失敗：-30
- 如果 WebSocket 連線失敗：-20
- 如果有 crash：-10 per crash
最低分：0
```

### Visual Consistency（視覺一致性）
```
Visual Consistency = 基礎分 91（已知美術質量）
+ sprite QC 通過率加成（最多 +5）
- 發現新的視覺問題（-2 per issue）
範圍：0-100
```

### Balance Health（數值健康度）
```
Balance Health = 100
- |實際RTP - 目標RTP| * 10（目標 RTP = 94%）
- 如果 RTP < 90%：額外 -20
- 如果 RTP > 98%：額外 -10
範圍：0-100
```

### Animation Quality（動畫品質）
```
Animation Quality = 平均 Consistency Score（所有角色所有動畫）
門檻：>= 88
```

### Audio Sync（音效同步）
```
Audio Sync = (正確觸發次數 / 總觸發次數) * 100
門檻：>= 90
```

### Gameplay Feel（遊戲手感）
```
Gameplay Feel = 主觀評分（QA Agent 評估）
評估維度：
- 射擊手感（30%）
- 命中反饋（25%）
- 獎勵演出（25%）
- 整體流暢度（20%）
門檻：>= 85
```

### Regression Risk（回歸風險）
```
Regression Risk = 新發現問題數 * 2 + 未修復舊問題數
門檻：<= 10
```

---

## 工具指令

```bash
# 執行完整 QA 檢查
py tools/qa_check.py

# 只執行 Build 檢查
py tools/qa_check.py --build-only

# 只執行 RTP 模擬
py tools/qa_check.py --rtp-only

# 只執行 Sprite QC
py tools/qa_check.py --sprite-only

# 輸出詳細報告
py tools/qa_check.py --verbose
```

---

## 報告格式

報告輸出至 `reports/qa/qa-report-YYYY-MM-DD.md`，格式見 Phase 5 實作。
