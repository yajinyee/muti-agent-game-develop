---
name: dev-self-check
description: 開發自檢清單。每次修改程式碼後必須執行的驗證步驟。
---

# 開發自檢清單

每次修改程式碼後，必須跑完以下檢查才算完成。

## 1. 編譯檢查

### Go 模組
- [ ] `go build ./...` 零錯誤
- [ ] `go vet ./...` 零警告

### TypeScript（Bot/前端）
- [ ] `npx tsc --noEmit` 零錯誤

## 2. 編碼問題（Windows 特有）
- [ ] 檔案寫入使用 `fs.writeFileSync(path, content, "utf-8")`（Node.js 預設 UTF-8 no BOM）
- [ ] 不使用 PowerShell 的 `Set-Content` 寫中文檔案
- [ ] 呼叫外部程式使用 `cmd.exe` + `chcp 65001`
- [ ] 不使用 `&&` 在 PowerShell 中（用 `;` 代替）

## 3. Kiro CLI 整合
- [ ] prompt 暫存檔用 `fs.writeFileSync` 寫入
- [ ] 使用 `cmd.exe` shell（不是 PowerShell）
- [ ] 有 `chcp 65001 >nul` 前綴
- [ ] catch 中有提取 stdout+stderr 的邏輯
- [ ] cleanOutput 正確過濾 kiro-cli 雜訊但保留有效內容

> 詳細呼叫方式見 `kiro-cli-caller` Skill

## 4. Telegram Bot 邏輯
- [ ] `isBotMentioned` 檢查 text + caption + entities + reply
- [ ] 群組訊息有被 `recordMessage` 記錄
- [ ] 長訊息有分段發送（≤4000 字元）
- [ ] 圖片/文件有正確下載和處理
- [ ] 錯誤有 try-catch 且回覆使用者「暫時無法回應」

## 5. 邏輯追蹤（最容易遺漏）
- [ ] 修改後，心理追蹤完整執行路徑
- [ ] 使用者 @ Bot → handleMessage → isBotMentioned → processAndRespond → Kiro CLI → cleanOutput → sendMessage
- [ ] 每個環節都有 log 輸出可追蹤
- [ ] 超時有合理的 fallback

## 6. 啟動驗證
- [ ] 服務/Bot 啟動 log 顯示 ✅
- [ ] 在群組 @ Bot 發一則測試訊息
- [ ] 確認回應是繁體中文 + 陳總人格
- [ ] 確認回應時間在合理範圍（< 60 秒）
