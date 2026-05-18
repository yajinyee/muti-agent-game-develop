# setup_scheduler.ps1
# 用 Windows Task Scheduler 建立每小時自動執行的排程
# 排程期間：2026-05-19 ~ 2026-08-19（三個月）
#
# 執行方式（不需要管理員）：
#   powershell -ExecutionPolicy Bypass -File "d:\Kiro\tools\setup_scheduler.ps1"

$TASK_NAME  = "KiroAutoContinue"
$SCRIPT     = "d:\Kiro\tools\auto_continue.ps1"
$INTERVAL   = 1          # 小時
$END_DATE   = "2026-08-19T23:59:00"   # 三個月後到期

# 先移除舊排程（若存在）
Unregister-ScheduledTask -TaskName $TASK_NAME -Confirm:$false -ErrorAction SilentlyContinue

# Action：每次觸發時執行 auto_continue.ps1 -RunOnce
$action = New-ScheduledTaskAction `
    -Execute "powershell.exe" `
    -Argument "-ExecutionPolicy Bypass -WindowStyle Hidden -File `"$SCRIPT`" -RunOnce"

# Trigger：從現在起，每 1 小時重複一次，到 2026-08-19 為止
$startTime = (Get-Date).AddMinutes(2)   # 2 分鐘後第一次
$endTime   = [datetime]$END_DATE

$trigger = New-ScheduledTaskTrigger `
    -Once `
    -At $startTime `
    -RepetitionInterval  (New-TimeSpan -Hours $INTERVAL) `
    -RepetitionDuration  ($endTime - $startTime)

# 設定
$settings = New-ScheduledTaskSettingsSet `
    -ExecutionTimeLimit  (New-TimeSpan -Minutes 10) `
    -MultipleInstances   IgnoreNew `
    -StartWhenAvailable `
    -WakeToRun:$false

# 以目前使用者執行（互動式登入）
$principal = New-ScheduledTaskPrincipal `
    -UserId    "$env:USERDOMAIN\$env:USERNAME" `
    -LogonType Interactive `
    -RunLevel  Limited

# 註冊
Register-ScheduledTask `
    -TaskName   $TASK_NAME `
    -Action     $action `
    -Trigger    $trigger `
    -Settings   $settings `
    -Principal  $principal `
    -Description "每 $INTERVAL 小時自動在 Kiro IDE chat 送出訊息，到 $END_DATE 為止" `
    -Force | Out-Null

# ── 顯示結果 ──────────────────────────────────────────────────────────────────

$task = Get-ScheduledTask -TaskName $TASK_NAME -ErrorAction SilentlyContinue
if ($task) {
    Write-Host ""
    Write-Host "✅ 排程建立完成！" -ForegroundColor Green
    Write-Host ""
    Write-Host "  名稱    ：$TASK_NAME"
    Write-Host "  間隔    ：每 $INTERVAL 小時"
    Write-Host "  第一次  ：$($startTime.ToString('yyyy-MM-dd HH:mm:ss'))（約 2 分鐘後）"
    Write-Host "  到期日  ：$END_DATE"
    Write-Host "  訊息內容：$(Get-Content 'd:\Kiro\tools\auto_continue_message.txt' -Encoding UTF8)"
    Write-Host ""

    # 計算總執行次數
    $totalHours = [math]::Round(($endTime - $startTime).TotalHours)
    Write-Host "  預計執行次數：約 $totalHours 次（$([math]::Round($totalHours/24)) 天 × 24 次/天）" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "管理指令：" -ForegroundColor Cyan
    Write-Host "  立即測試：Start-ScheduledTask -TaskName '$TASK_NAME'"
    Write-Host "  查看狀態：Get-ScheduledTask -TaskName '$TASK_NAME' | Select-Object State, LastRunTime, NextRunTime"
    Write-Host "  暫停排程：Disable-ScheduledTask -TaskName '$TASK_NAME'"
    Write-Host "  恢復排程：Enable-ScheduledTask -TaskName '$TASK_NAME'"
    Write-Host "  刪除排程：Unregister-ScheduledTask -TaskName '$TASK_NAME' -Confirm:`$false"
    Write-Host ""
    Write-Host "Log 位置：d:\Kiro\tools\auto_continue.log" -ForegroundColor Yellow
} else {
    Write-Host "❌ 排程建立失敗，請確認權限" -ForegroundColor Red
}
