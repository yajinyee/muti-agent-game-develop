# Skill: Git Windows Permissions

> 版本：1.0.0  
> 來源：實戰踩坑記錄  
> 最後更新：2026-05-17  
> 適用系統：Windows 10/11

---

## 問題描述

在 Windows 上使用 Git 時，常遇到以下錯誤：

```
error: unable to create file .git/objects/...
Permission denied
```

或：

```
fatal: cannot create directory at '.git/objects/pack': Permission denied
```

### 常見觸發場景
1. 多個程序同時存取 `.git` 目錄（如 IDE + 命令列）
2. 防毒軟體鎖定 `.git` 目錄
3. 目錄繼承了錯誤的 ACL 權限
4. 使用 OneDrive/Dropbox 同步的目錄

---

## 解決方案

### 方案 1：icacls 修復權限（推薦）

```powershell
# 修復 .git 目錄的完整權限
# 在專案根目錄執行

$gitDir = ".git"
$currentUser = [System.Security.Principal.WindowsIdentity]::GetCurrent().Name

# 給予當前使用者完整控制權
icacls $gitDir /grant "${currentUser}:(OI)(CI)F" /T /C

# 啟用繼承
icacls $gitDir /inheritance:e /T /C

Write-Host "Git 目錄權限已修復"
```

### 方案 2：設定 Git tmpdir

```powershell
# 設定 Git 使用自訂 tmpdir（避免系統 temp 目錄權限問題）
$tmpDir = "$env:USERPROFILE\.git-tmp"
New-Item -ItemType Directory -Path $tmpDir -Force | Out-Null

git config --global core.tmpdir $tmpDir
Write-Host "Git tmpdir 已設定為：$tmpDir"
```

### 方案 3：完整修復腳本（git_add_all.ps1）

```powershell
# git_add_all.ps1
# 修復 Git 權限問題並執行 git add

param(
    [string]$CommitMessage = "Update files"
)

$ErrorActionPreference = "Continue"

# Step 1: 修復 .git 目錄權限
Write-Host "修復 .git 目錄權限..."
$currentUser = [System.Security.Principal.WindowsIdentity]::GetCurrent().Name
icacls ".git" /grant "${currentUser}:(OI)(CI)F" /T /C 2>$null
icacls ".git" /inheritance:e /T /C 2>$null

# Step 2: 設定 tmpdir
$tmpDir = "$env:USERPROFILE\.git-tmp"
New-Item -ItemType Directory -Path $tmpDir -Force | Out-Null
git config core.tmpdir $tmpDir

# Step 3: 執行 git add
Write-Host "執行 git add..."
git add -A

# Step 4: 確認 staging area
$status = git status --short
if ($status) {
    Write-Host "已暫存的變更："
    Write-Host $status
    
    # Step 5: 提交（如果有 commit message）
    if ($CommitMessage -ne "Update files") {
        git commit -m $CommitMessage
        Write-Host "已提交：$CommitMessage"
    }
} else {
    Write-Host "沒有需要暫存的變更"
}
```

### 方案 4：防毒軟體排除

如果是防毒軟體造成的問題：

1. Windows Defender：
   - 開啟「Windows 安全性」
   - 病毒與威脅防護 → 管理設定
   - 排除項目 → 新增排除項目 → 資料夾
   - 加入專案根目錄

2. 其他防毒軟體：
   - 在防毒軟體設定中加入專案目錄為「信任區域」

---

## 預防措施

### 1. 避免在 OneDrive/Dropbox 目錄中使用 Git

```powershell
# 確認目前目錄不在同步目錄中
$cwd = Get-Location
if ($cwd -like "*OneDrive*" -or $cwd -like "*Dropbox*") {
    Write-Warning "警告：目前在雲端同步目錄中，可能造成 Git 問題"
    Write-Warning "建議將專案移到 C:\Projects\ 等非同步目錄"
}
```

### 2. 設定 .gitattributes 避免換行符問題

```
# .gitattributes
* text=auto eol=lf
*.ps1 text eol=crlf
*.bat text eol=crlf
```

### 3. 設定 Git 全域設定

```powershell
# 設定 Windows 友善的 Git 設定
git config --global core.autocrlf true
git config --global core.longpaths true
git config --global core.fileMode false
```

---

## 快速診斷

```powershell
# 快速診斷 Git 權限問題
function Test-GitPermissions {
    param([string]$RepoPath = ".")
    
    $gitDir = Join-Path $RepoPath ".git"
    
    if (-not (Test-Path $gitDir)) {
        Write-Host "不是 Git 倉庫"
        return
    }
    
    # 測試寫入權限
    $testFile = Join-Path $gitDir "permission_test_$(Get-Random)"
    try {
        [System.IO.File]::WriteAllText($testFile, "test")
        Remove-Item $testFile -Force
        Write-Host "✅ .git 目錄寫入權限正常"
    } catch {
        Write-Host "❌ .git 目錄寫入權限異常：$_"
        Write-Host "建議執行：icacls .git /grant `"${env:USERNAME}:(OI)(CI)F`" /T /C"
    }
}

Test-GitPermissions
```

---

## 相關工具

- `tools/git_add_all.ps1` — 修復權限並執行 git add 的腳本
- 本 Skill 文件記錄了所有已知的 Windows Git 權限問題解決方案
