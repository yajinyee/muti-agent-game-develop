# setup_scheduler.ps1
# 用 Windows Task Scheduler 建立每兩小時自動執行的排程
#
# 執行方式（不需要管理員）：
#   powershell -ExecutionPolicy Bypass -File "d:\Kiro\tools\setup_scheduler.ps1"

$TASK_NAME  = "KiroAutoContinue"
$SCRIPT     = "d:\Kiro\tools\auto_continue.ps1"
$INTERVAL   = 2   # 小時

# 先移除舊排程
Unregister-ScheduledTask -TaskName $TASK_NAME -Confirm:$false -ErrorAction SilentlyContinue

# Action：每次觸發時執行 auto_continue.ps1 -RunOnce
$action = New-ScheduledTaskAction `
    -Execute "powershell.exe" `
    -Argument "-ExecutionPolicy Bypass -WindowStyle Hidden -File `"$SCRIPT`" -RunOnce"

# Trigger：從現在起，每 2 小時重複一次
$startTime = (Get-Date).AddMinutes(2)   # 2 分鐘後第一次
$trigger = New-ScheduledTaskTrigger `
    -Once `
    -At $startTime `
    -RepetitionInterval (New-TimeSpan -Hours $INTERVAL) `
    -RepetitionDuration ([TimeSpan]::MaxValue)

# 設定
$settings = New-ScheduledTaskSettingsSet `
    -ExecutionTimeLimit (New-TimeSpan -Minutes 5) `
    -MultipleInstances IgnoreNew `
    -StartWhenAvailable `
    -WakeToRun:$false

# 以目前使用者執行（互動式登入，才能操作 UI）
$principal = New-ScheduledTaskPrincipal `
    -UserId "$env:USERDOMAIN\$env:USERNAME" `
    -LogonType Interactive `
    -RunLevel Limited

# 註冊
Register-ScheduledTask `
    -TaskName $TASK_NAME `
    -Action $action `
    -Trigger $trigger `
    -Settings $settings `
    -Principal $principal `
    -Description "每 $INTERVAL 小時自動在 Kiro IDE chat 送出「延續上一次規劃跟當前專案的原則與架構，繼續執行運作」" `
    -Force | Out-Null

Write-Host ""
Write-Host "✅ 排程建立完成！" -ForegroundColor Green
Write-Host ""
Write-Host "  名稱：$TASK_NAME"
Write-Host "  間隔：每 $INTERVAL 小時"
Write-Host "  第一次：$($startTime.ToString('HH:mm:ss'))（約 2 分鐘後）"
Write-Host ""
Write-Host "管理指令：" -ForegroundColor Cyan
Write-Host "  立即測試：Start-ScheduledTask -TaskName '$TASK_NAME'"
Write-Host "  查看狀態：Get-ScheduledTask -TaskName '$TASK_NAME' | Select-Object State"
Write-Host "  暫停排程：Disable-ScheduledTask -TaskName '$TASK_NAME'"
Write-Host "  恢復排程：Enable-ScheduledTask -TaskName '$TASK_NAME'"
Write-Host "  刪除排程：Unregister-ScheduledTask -TaskName '$TASK_NAME' -Confirm:`$false"
Write-Host ""
Write-Host "Log 位置：d:\Kiro\tools\auto_continue.log" -ForegroundColor Yellow
