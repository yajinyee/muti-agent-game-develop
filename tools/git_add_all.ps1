# 逐個 add 檔案，每次失敗後修復權限再重試
Set-Location "d:\Kiro"
$env:TMPDIR = "d:\Kiro\.git\tmp"
$env:TMP = "d:\Kiro\.git\tmp"
$env:TEMP = "d:\Kiro\.git\tmp"

function Fix-GitPermissions {
    icacls "d:\Kiro\.git\objects" /grant "$env:USERNAME`:F" /T /Q | Out-Null
    icacls "d:\Kiro\.git\objects" /inheritance:e /T /Q | Out-Null
    icacls "d:\Kiro\.git\tmp" /grant "$env:USERNAME`:F" /Q | Out-Null
}

# 取得所有未追蹤和修改的檔案
$files = git status --short | ForEach-Object { $_.Substring(3).Trim() }
$total = $files.Count
$done = 0
$failed = 0

foreach ($file in $files) {
    Fix-GitPermissions
    $result = git add "$file" 2>&1
    if ($LASTEXITCODE -eq 0) {
        $done++
    } else {
        # 再試一次
        Fix-GitPermissions
        $result2 = git add "$file" 2>&1
        if ($LASTEXITCODE -eq 0) {
            $done++
        } else {
            Write-Warning "Failed: $file"
            $failed++
        }
    }
    if ($done % 20 -eq 0) {
        Write-Host "Progress: $done/$total (failed: $failed)"
    }
}

Write-Host "Done: $done/$total staged, $failed failed"
