# CLAUDE.md — 智能體根文件

本文件是 AI 智能體的根準則，始終載入、始終啟用。

## 架構規則

- Skills 放在 `.kiro/skills/`，工具放在 `tools/`
- Hooks 放在 `.kiro/hooks/`
- Steering 放在 `.kiro/steering/`
- 子代理（Subagents）透過 Kiro 內建機制委派任務

## 命名規範

- TypeScript 檔案使用 camelCase
- 資料夾使用 kebab-case
- Skill 文件使用 kebab-case.md

## 測試期望

- Go：每次修改後執行 `go build ./...` + `go vet ./...`
- TypeScript：每次修改後執行 `npx tsc --noEmit`
- 服務啟動後確認 log 輸出正常，發送測試請求確認回應正確

## 倉庫地圖

```
.
├── CLAUDE.md                  # 根文件（本文件）
├── .kiro/
│   ├── steering/              # 第1層延伸：全域 Steering 規則
│   ├── skills/                # 第2層：Skills 知識層
│   └── hooks/                 # 第3層：Hooks 守護層
├── tools/                     # 工具腳本
├── src/                       # 主要程式碼
│   └── telegram/              # Telegram Bot
├── data/
│   ├── training/              # 訓練資料
│   └── tmp/                   # 暫存檔
└── docs/                      # 文件
```

## 核心原則摘要

詳見 `.kiro/steering/core-principles.md`，以下為快速參考：

1. 遇到障礙就建立 Skill
2. 不能放棄，持續嘗試
3. 每件事都要有驗證
4. **Go 優先**（Server/底層），TypeScript 用於 Bot/前端
5. 一次執行到底
6. Kiro CLI 整合驗證
