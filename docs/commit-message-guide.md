# Commit Message 規範

> 所有 Agent 在 commit 時必須遵守此規範。

---

## 格式

```
<type>(<scope>): <簡短說明>（50字以內）

[空行]

[詳細說明]
- 做了什麼
- 為什麼這樣做
- 影響範圍

[可選：QA 結果]
QA: 8/8 通過 | Build: OK | Risk: Low
```

---

## Type 類型

| Type | 用途 | 範例 |
|------|------|------|
| `feat` | 新功能 | `feat(server): 加入 BOSS 自動觸發機制` |
| `fix` | 修復 bug | `fix(animation): 修復 usagi bigwin 幀一致性問題` |
| `chore` | 維護工作 | `chore(memory): 更新 project-memory.md 品質分數` |
| `docs` | 文件更新 | `docs(spec): 補齊 WebSocket 協定說明` |
| `refactor` | 重構 | `refactor(hub): 廣播改用非阻塞 select` |
| `perf` | 效能優化 | `perf(ws): 啟用 WebSocket permessage-deflate 壓縮` |
| `release` | 每日整合 | `release: DAY-003 integration（8/8 QA 通過）` |
| `merge` | Branch 合併 | `merge: ANIM-001 修復 Animation Quality 87→100` |

---

## Scope 範圍

| Scope | 說明 |
|-------|------|
| `server` | Go Server 相關 |
| `client` | Godot Client 相關 |
| `animation` | 動畫相關 |
| `art` | 美術資產相關 |
| `audio` | 音效相關 |
| `spec` | 規格文件相關 |
| `tools` | Python 工具腳本 |
| `memory` | memory/ 文件 |
| `skills` | skills/ 文件 |
| `qa` | QA 相關 |
| `ws` | WebSocket 相關 |

---

## 好的範例

```
feat(server): 加入 BOSS 自動觸發機制（規格書 28.1）

每 3-5 分鐘自動觸發 BOSS，不再只依賴手動按鈕。
- 在 Game struct 加入 nextBossAt 計時器
- updateNormalPlay() 中檢查時間並觸發 triggerBoss()
- 初始延遲 180-300 秒（隨機）

影響：game.go
QA: go build OK | go vet OK | Risk: Low
```

```
fix(animation): 修復 qa_check.py Animation Quality 硬編碼問題

問題：Animation Quality 被硬編碼為 87，不反映實際狀態
原因：animation_pipeline.py 的 importlib 載入環境問題
修復：改為直接 import，從實際 audit 結果取值

結果：Animation Quality 87 → 100
QA: 8/8 通過
```

---

## 壞的範例（不要這樣寫）

```
❌ update files
❌ fix bug
❌ chore: DAY-003 Nightly Report + Memory 更新（規格一致性 97%）  ← 太簡略
❌ feat: 加功能
```

---

## 每日 Release Commit 格式

```
release: DAY-XXX integration（品質摘要）

[今日完成]
- Agent/Task 1：說明
- Agent/Task 2：說明

[品質分數]
Build: 100 | Animation: 100 | Balance: 96 | Spec: 97
8/8 QA 指標通過 ✅

[修改的檔案]
- path/to/file1：說明
- path/to/file2：說明

[下一步]
DAY-XXX+1 計畫：...
```

*最後更新：2026-05-19（DAY-003）*
