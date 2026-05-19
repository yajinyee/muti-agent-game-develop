# 自動 commit + push 腳本
param([string]$Message = "chore: auto update $(Get-Date -Format 'yyyy-MM-dd HH:mm')")

Set-Location "d:\Kiro"

# 設定 GIT_TMPDIR（修復 Windows 上 git add 的 temp file 問題）
$env:GIT_TMPDIR = "d:\Kiro\.git\tmp"

# 修復 git objects 權限（避免 Permission denied）
icacls "d:\Kiro\.git\objects" /grant "$env:USERNAME`:F" /T /Q | Out-Null
icacls "d:\Kiro\.git\objects" /inheritance:e /T /Q | Out-Null

# Stage 所有變更（逐一 add 避免批次失敗）
$files = git status --short | ForEach-Object { $_.Substring(3).Trim() }
foreach ($file in $files) {
    git add "$file" 2>&1 | Out-Null
}

# 確認有變更才 commit
$status = git status --short
if ($status) {
    git commit -m $Message 2>&1 | Out-Null
    git push origin main 2>&1 | Out-Null
    Write-Host "Pushed: $Message"
} else {
    Write-Host "No changes"
}
