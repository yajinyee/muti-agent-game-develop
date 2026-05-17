# 修復 .git/objects 目錄的繼承權限問題
$gitObjects = "d:\Kiro\.git\objects"
$username = $env:USERNAME

# 取得擁有權
takeown /F $gitObjects /R /D Y | Out-Null

# 設定完整控制權限（含繼承）
icacls $gitObjects /grant "${username}:F" /T /Q | Out-Null
icacls $gitObjects /inheritance:e /Q | Out-Null

# 確認
Write-Host "Git objects permissions fixed for: $username"
Write-Host "Testing write..."
$testFile = Join-Path $gitObjects "test_write_$(Get-Random).tmp"
try {
    [System.IO.File]::WriteAllText($testFile, "test")
    Remove-Item $testFile
    Write-Host "Write test: OK"
} catch {
    Write-Host "Write test FAILED: $_"
}
