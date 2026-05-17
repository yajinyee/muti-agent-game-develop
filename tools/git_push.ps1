# 自動 commit + push 腳本
param([string]$Message = "chore: auto update $(Get-Date -Format 'yyyy-MM-dd HH:mm')")

Set-Location "d:\Kiro"

# 修復 git objects 權限（避免 Permission denied）
icacls "d:\Kiro\.git\objects" /grant "$env:USERNAME`:F" /T /Q | Out-Null
icacls "d:\Kiro\.git\objects" /inheritance:e /T /Q | Out-Null

# Stage 所有變更
git add . 2>&1 | Out-Null

# 確認有變更才 commit
$status = git status --short
if ($status) {
    git commit -m $Message 2>&1 | Out-Null
    git push origin master 2>&1 | Out-Null
    Write-Host "Pushed: $Message"
} else {
    Write-Host "No changes"
}
