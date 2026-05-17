---
name: kiro-cli-caller
description: 如何正確呼叫 Kiro CLI headless mode 執行 AI 任務。包含路徑、編碼、錯誤處理、prompt 傳遞的完整 KnowHow。當需要透過 Kiro CLI 生成 AI 回應時使用。
---

# Kiro CLI 呼叫 Skill

## 環境資訊
- 安裝路徑：`C:\Program Files\Kiro-Cli\kiro-cli.exe`（已驗證 2026-05-07）
- 版本：2.2.2
- 認證：企業帳號已登入（browser OAuth，不需要 KIRO_API_KEY）
- 模型：claude-opus-4.6
- 驗證日期：2026-05-07 ✅

## 正確呼叫方式

### 1. Prompt 寫入暫存檔（UTF-8 without BOM）
```typescript
const fs = require("fs");
const promptFile = `./data/tmp/prompt_${Date.now()}.txt`;
fs.writeFileSync(promptFile, prompt, "utf-8"); // Node.js 預設就是 UTF-8 no BOM
```

### 2. 使用 cmd.exe + chcp 65001 確保中文正確
```typescript
const cmd = `chcp 65001 >nul && type "${promptFile}" | "${KIRO_CLI_PATH}" chat --no-interactive --trust-all-tools`;
const result = await execAsync(cmd, {
  timeout: 120000,
  shell: "cmd.exe",
  encoding: "utf-8",
  maxBuffer: 1024 * 1024 * 10,
});
```

### 3. 處理 exit code 1（正常情況）
Kiro CLI 會把警告訊息寫到 stderr，導致 PowerShell 認為是錯誤（exit code 1）。
但 stdout 中有正確的回應。必須在 catch 中嘗試提取回應：
```typescript
catch (error: any) {
  const combined = (error.stdout || "") + (error.stderr || "");
  const cleaned = cleanOutput(combined);
  if (cleaned && cleaned.length > 20) return cleaned; // 有效回應
  throw error; // 真正的錯誤
}
```

### 4. 清理輸出（移除 kiro-cli 雜訊）
必須移除：
- ANSI escape codes
- `> ` 前綴
- `All tools are now trusted` 警告
- `Agents can sometimes do unexpected things` 警告
- `Time: Xs` 計時
- spinner 字元（⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏▰▱）

## KnowHow（踩過的坑）

| 問題 | 原因 | 解決 |
|------|------|------|
| 中文變亂碼 | PowerShell 預設編碼不是 UTF-8 | 改用 cmd.exe + chcp 65001 |
| 回應為空 | exit code 1 進入 catch，沒提取 stdout | catch 中合併 stdout+stderr |
| prompt 被截斷 | shell 參數長度限制 | 寫入暫存檔用 type pipe |
| 人格沒帶入 | prompt 太短沒包含 systemPrompt | 確保完整 systemPrompt 寫入檔案 |
| Set-Content 編碼錯 | PowerShell 加 BOM 或用 Big5 | 用 Node.js fs.writeFileSync |
