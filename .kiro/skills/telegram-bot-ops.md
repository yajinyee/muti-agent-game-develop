---
name: telegram-bot-ops
description: Telegram Bot 的啟動、停止、除錯、驗證標準流程。當需要操作或除錯 Bot 時使用。
---

# Telegram Bot 操作 Skill

## Bot 資訊
- Token: 8793843542:AAEss_pJAo9LHt-ZNSioN6bFUKOW-HeJtWo
- Username: @igs_arcade_father_bot
- 名稱: 陳總(數位分身)
- Privacy Mode: 已關閉（可收到群組所有訊息）

## 啟動流程

```bash
npx ts-node src/telegram/bot.ts
```

### 啟動成功的標誌
```
[陳總Bot] ✅ 啟動完成！
[陳總Bot] AI 後端: Kiro CLI + claude-opus-4.6
[陳總Bot] Bot username: @igs_arcade_father_bot
```

### 啟動失敗排查
1. `ECONNREFUSED` → 網路問題或 Telegram API 被封
2. `409 Conflict` → 有另一個 Bot 實例在跑，先停掉
3. `401 Unauthorized` → Token 錯誤
4. `Kiro CLI 不可用` → 確認 kiro-cli 已安裝且已登入

## 停止流程
- 在 terminal 按 Ctrl+C
- 或用 Kiro IDE 的 controlPwshProcess stop

## 除錯流程

### Bot 收到訊息但沒回應
1. 檢查 log 有沒有 `[KiroBridge] 請求已寫入` 或 `回應錯誤`
2. 確認 Kiro CLI 能正常回應：
```powershell
cmd /c "chcp 65001 >nul && echo 回覆OK | `"C:\Program Files\Kiro-Cli\kiro-cli.exe`" chat --no-interactive --trust-all-tools"
```
3. 確認 Bot 的 Privacy Mode 已關閉

### Bot 回應但內容不對
1. 確認 systemPrompt 有被帶入（檢查 prompt 暫存檔內容）
2. 確認 cleanOutput 沒有把有效內容過濾掉
3. 確認中文編碼正確（cmd.exe + chcp 65001）

### Bot 回應亂碼
- 原因：沒有用 cmd.exe + chcp 65001
- 解決：確認 executeKiroCli 使用 `shell: "cmd.exe"` 且有 `chcp 65001 >nul`

## 驗證清單
- [ ] `npx tsc --noEmit` 編譯通過
- [ ] Bot 啟動 log 顯示 ✅
- [ ] 在群組 @ Bot 能收到回應
- [ ] 回應是繁體中文
- [ ] 回應帶有陳總人格特徵
- [ ] 圖片/文件能被下載處理
