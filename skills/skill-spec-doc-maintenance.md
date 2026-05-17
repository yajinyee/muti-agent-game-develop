# Skill：規格文件維護原則

> 記錄日期：2026-05-19（DAY-003）  
> 記錄者：Spec Architect Agent / Skill Librarian

---

## 問題描述

規格文件（game-spec.md）和實作之間容易出現不一致：
1. 功能已實作但文件標記「待定義」
2. 數值因 RTP 校正調整但文件未更新
3. memory 文件出現重複區塊

## 解決方案

### 1. 每日 Spec Completeness 檢查清單

```
□ WebSocket 協定是否完整記錄？
□ 所有已實作功能是否在規格書中有對應描述？
□ 數值（倍率、RTP、機率）是否與實作一致？
□ memory 文件是否有重複或過時內容？
□ 規格書的「待定義」標記是否已清除？
```

### 2. 規格更新流程

```
實作變更
    ↓
Spec Architect 更新 game-spec.md
    ↓
更新 memory/project-memory.md 的規格一致性分數
    ↓
QA Agent 重新計算 Spec Completeness
```

### 3. 常見不一致類型

| 類型 | 範例 | 修復方式 |
|------|------|---------|
| 待定義未清除 | WebSocket 協定標記「待定義」| 補充完整實作內容 |
| 數值不一致 | Bonus 倍率規格 50-150x，實作 20-50x | 加注釋說明調整原因 |
| 重複內容 | memory 有兩個「專案基本資訊」| 合併，保留最新版本 |
| 缺少實作說明 | 規格有功能但沒有實作細節 | 補充實作方式 |

## 預防措施

- 每次修改數值時，同步更新 game-spec.md
- 每次新增功能時，在規格書加入對應章節
- memory 文件每日只更新一次，避免重複

## 相關檔案
- `docs/game-spec.md`
- `memory/project-memory.md`
- `docs/acceptance-criteria.md`

*Content was rephrased for compliance with licensing restrictions*
