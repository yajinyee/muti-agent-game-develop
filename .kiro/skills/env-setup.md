---
name: env-setup
description: 本專案的完整環境安裝與驗證流程。當需要在新機器上建立開發環境，或確認環境是否正確時使用。
---

# 環境安裝 Skill

## 已驗證環境（2026-05-07，Windows 11，D:\Kiro）

| 工具 | 版本 | 安裝方式 | 狀態 |
|------|------|----------|------|
| Go | 1.26.2 | `winget install --id GoLang.Go --source winget` | ✅ |
| Node.js | 24.15.0 (LTS) | `winget install --id OpenJS.NodeJS.LTS --source winget` | ✅ |
| TypeScript | 6.0.3 | `npm install -g typescript ts-node` | ✅ |
| ts-node | 10.9.2 | 同上 | ✅ |
| Git | 2.54.0 | `winget install --id Git.Git --source winget` | ✅ |
| Kiro CLI | 2.2.2 | `irm 'https://cli.kiro.dev/install.ps1' \| iex` | ✅ |

## 安裝順序

```powershell
# 1. Go
winget install --id GoLang.Go --source winget --silent --accept-package-agreements --accept-source-agreements

# 2. Node.js
winget install --id OpenJS.NodeJS.LTS --source winget --silent --accept-package-agreements --accept-source-agreements

# 3. Git
winget install --id Git.Git --source winget --silent --accept-package-agreements --accept-source-agreements

# 4. 重新載入 PATH（新終端機或執行以下）
$env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")

# 5. 解除 PowerShell 執行原則（npm 需要）
Set-ExecutionPolicy -Scope CurrentUser -ExecutionPolicy RemoteSigned -Force

# 6. TypeScript + ts-node（需先關閉 strict-ssl，企業網路憑證問題）
npm config set strict-ssl false
npm install -g typescript ts-node

# 7. Kiro CLI
irm 'https://cli.kiro.dev/install.ps1' | iex

# 8. Kiro CLI 登入
kiro-cli login
```

## 驗證指令

```powershell
$env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")
go version
node --version
npm --version
tsc --version
git --version
kiro-cli --version
```

## Kiro CLI Headless 測試

```powershell
# 寫入測試 prompt
[System.IO.File]::WriteAllText("D:\Kiro\data\tmp\test.txt", "請回答：1+1=?", [System.Text.Encoding]::UTF8)

# 執行並輸出到檔案
cmd /c "chcp 65001 >nul && type D:\Kiro\data\tmp\test.txt | kiro-cli chat --no-interactive --trust-all-tools > D:\Kiro\data\tmp\out.txt 2>&1"

# 讀取結果
Get-Content D:\Kiro\data\tmp\out.txt -Encoding UTF8
```

期望看到 `> 2` 或類似回應，代表認證與 headless 模式正常。

## 已知問題

- **Docker Desktop**：需要系統管理員權限，用一般帳號 winget 安裝會失敗（exit code 4294967291）。需以系統管理員身份執行安裝。
- **npm strict-ssl**：企業網路有自簽憑證，需 `npm config set strict-ssl false`
- **PowerShell 執行原則**：預設會擋 npm.ps1，需設定 `RemoteSigned`
- **Kiro CLI 路徑**：安裝在 `C:\Program Files\Kiro-Cli\kiro-cli.exe`，**不是** `$env:LOCALAPPDATA\Kiro-Cli\`
