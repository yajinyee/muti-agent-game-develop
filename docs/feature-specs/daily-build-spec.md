# Daily Build 流程規格

> 版本：1.0.0  
> 維護者：Game Director + Go Server Agent  
> 最後更新：2026-05-17

---

## 概覽

本規格定義吉伊卡哇：像素大討伐的每日 Build 流程，確保每日都有可玩的穩定版本。

---

## Branch 命名規範

### Agent 工作分支
```
agent/<agent-name>/<task-id>-<title>

範例：
agent/animation-agent/ANIM-001-chiikawa-idle-8frames
agent/go-server-agent/SRV-003-boss-phase2-logic
agent/godot-client-agent/UI-007-bet-panel-redesign
agent/balance-agent/BAL-002-rtp-calibration
```

### 每日整合分支
```
integration/daily-YYYYMMDD

範例：
integration/daily-20260517
integration/daily-20260518
```

### 可玩發布分支
```
release/playable-YYYYMMDD

範例：
release/playable-20260517
```

---

## Merge 條件

### 從 agent 分支 merge 到 integration/daily

所有條件必須全部通過：

| 條件 | 檢查方式 | 通過標準 |
|------|---------|---------|
| Build Check | `go build ./...` | 零錯誤 |
| Lint | `go vet ./...` | 零警告 |
| Test | `go test ./...` | 全部通過 |
| Work Report | 對應報告存在 | 報告完整 |
| Risk Classification | 評估變更影響範圍 | 風險 <= 中 |
| Rollback Plan | 說明如何回滾 | 計畫存在 |

### 從 integration/daily merge 到 release/playable

額外條件：

| 條件 | 檢查方式 | 通過標準 |
|------|---------|---------|
| QA Report | `py tools/qa_check.py` | 全部 P0 通過 |
| Build Stability | QA 分數 | >= 95 |
| Regression Risk | QA 分數 | <= 10 |
| Art Director 審核 | 視覺審查 | 通過 |

---

## 每日輸出物清單

每日 Build 完成後，必須產出以下所有項目：

| 輸出物 | 路徑 | 說明 |
|-------|------|------|
| Playable Build | `builds/daily/YYYYMMDD/` | HTML5 可玩版本 |
| Demo Video | `builds/daily/YYYYMMDD/demo.mp4` | 30 秒遊戲影片 |
| Screenshots | `builds/daily/YYYYMMDD/screenshots/` | 關鍵畫面截圖 |
| Quality Score | `reports/quality/quality-score-YYYYMMDD.md` | 8 項品質分數 |
| QA Report | `reports/qa/qa-report-YYYYMMDD.md` | 完整 QA 報告 |
| Bug List | `reports/qa/bug-list-YYYYMMDD.md` | 當日發現的 Bug |
| Next Day Plan | `tasks/today-plan.md`（更新）| 明日任務計畫 |
| Retro | `reports/nightly/nightly-report-YYYYMMDD.md` | 今日回顧 |

---

## 每日 Build 時間表

| 時間 | 動作 | 負責 Agent |
|------|------|-----------|
| 09:00 | 讀取昨日 nightly report | Game Director |
| 09:15 | 更新 today-plan.md | Game Director |
| 09:30 | 各 Agent 開始工作 | 全體 |
| 16:00 | 各 Agent 提交工作分支 | 全體 |
| 16:30 | 執行 QA 自動測試 | QA Agent |
| 17:00 | 整合到 integration/daily | Game Director |
| 17:30 | 生成 Playable Build | Godot Client Agent |
| 18:00 | 輸出所有報告 | 各 Agent |
| 18:30 | 撰寫 nightly report | Game Director |

---

## 自動化腳本

```powershell
# 執行每日 Build
powershell -File tools/daily_build.ps1

# 執行後確認
# 1. 查看 reports/quality/quality-score-YYYYMMDD.md
# 2. 查看 reports/qa/qa-report-YYYYMMDD.md
# 3. 確認 builds/daily/YYYYMMDD/ 目錄存在
```

---

## 緊急處理

### Build 失敗
1. 立即停止 merge
2. Go Server Agent 確認問題
3. 修復後重新執行 Build Check
4. 記錄到 failed-attempts/

### QA 分數低於門檻
1. 不產出 release/playable 分支
2. 記錄問題到 bug-list
3. 明日優先修復
4. 通知 Game Director

### 緊急 Rollback
```powershell
# 回滾到上一個穩定版本
git checkout release/playable-<上一個日期>
go build ./...
# 確認編譯通過後部署
```
