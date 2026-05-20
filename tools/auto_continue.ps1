param(
    [int]$IntervalHours = 2,
    [switch]$RunOnce,
    [int]$MaxSessions = 5   # Kiro chat 最多保留幾個 session
)

$LOG_FILE    = "d:\Kiro\tools\auto_continue.log"
$MSG_FILE    = "d:\Kiro\tools\auto_continue_message.txt"

# Kiro sessions.json 路徑（workspace ID = ZDpcS2lybw__ 對應 d:\Kiro）
$SESSIONS_JSON = "C:\Users\$env:USERNAME\AppData\Roaming\Kiro\User\globalStorage\kiro.kiroagent\workspace-sessions\ZDpcS2lybw__\sessions.json"
$SESSIONS_DIR  = "C:\Users\$env:USERNAME\AppData\Roaming\Kiro\User\globalStorage\kiro.kiroagent\workspace-sessions\ZDpcS2lybw__"

Add-Type -AssemblyName System.Windows.Forms
$src = 'using System; using System.Runtime.InteropServices; public class WinAPI2 { [DllImport("user32.dll")] public static extern bool SetForegroundWindow(IntPtr h); [DllImport("user32.dll")] public static extern bool ShowWindow(IntPtr h, int n); [DllImport("user32.dll")] public static extern bool IsIconic(IntPtr h); }'
Add-Type -TypeDefinition $src -ErrorAction SilentlyContinue

function Write-Log {
    param([string]$msg)
    $ts = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
    $line = "[$ts] $msg"
    Write-Host $line
    Add-Content -Path $LOG_FILE -Value $line -Encoding UTF8
}

function Trim-OldSessions {
    <#
    .SYNOPSIS
        直接操作 sessions.json，保留最新 $MaxSessions 個，刪除舊的 session 資料夾。
    .NOTES
        Kiro 必須在修改後重新載入，但 sessions.json 是 Kiro 啟動時讀取的，
        執行中修改不影響目前顯示，但下次啟動後就會生效。
        為了讓目前執行中的 Kiro 也能感知，我們在送訊息前先清理，
        讓 session 列表保持乾淨。
    #>

    if (-not (Test-Path $SESSIONS_JSON)) {
        Write-Log "sessions.json 不存在：$SESSIONS_JSON"
        return
    }

    $content = [System.IO.File]::ReadAllText($SESSIONS_JSON, [System.Text.Encoding]::UTF8)
    $sessions = $content | ConvertFrom-Json

    $total = $sessions.Count
    Write-Log "目前 session 總數：$total（上限 $MaxSessions）"

    if ($total -le $MaxSessions) {
        Write-Log "Session 數量正常，不需要清理"
        return
    }

    # 按 dateCreated 降序排列（最新在前）
    $sorted = $sessions | Sort-Object { [long]$_.dateCreated } -Descending

    # 保留最新 MaxSessions 個
    $keep   = $sorted | Select-Object -First $MaxSessions
    $remove = $sorted | Select-Object -Skip  $MaxSessions

    Write-Log "保留 $($keep.Count) 個，刪除 $($remove.Count) 個舊 session"

    # 刪除舊 session 的資料夾
    foreach ($s in $remove) {
        $sessionFolder = Join-Path $SESSIONS_DIR $s.sessionId
        if (Test-Path $sessionFolder) {
            Remove-Item $sessionFolder -Recurse -Force -ErrorAction SilentlyContinue
            Write-Log "  已刪除 session 資料夾：$($s.sessionId.Substring(0,8))... ($($s.title.Substring(0, [Math]::Min(30, $s.title.Length))))"
        }
    }

    # 更新 sessions.json（只保留 keep 的部分，維持原始順序）
    $keepIds = $keep | ForEach-Object { $_.sessionId }
    $newSessions = $sessions | Where-Object { $keepIds -contains $_.sessionId }

    $newJson = $newSessions | ConvertTo-Json -Depth 10
    # 確保是陣列格式（單一元素時 ConvertTo-Json 可能輸出物件）
    if ($newSessions.Count -eq 1) {
        $newJson = "[$newJson]"
    }
    [System.IO.File]::WriteAllText($SESSIONS_JSON, $newJson, [System.Text.Encoding]::UTF8)

    Write-Log "sessions.json 已更新，剩餘 $($keep.Count) 個 session"
}

function Send-ToKiro {
    # 從檔案讀取訊息（支援中文）
    $text = [System.IO.File]::ReadAllText($MSG_FILE, [System.Text.Encoding]::UTF8).Trim()

    $proc = Get-Process -Name "Kiro" -ErrorAction SilentlyContinue `
        | Where-Object { $_.MainWindowTitle -ne "" } `
        | Select-Object -First 1

    if (-not $proc) {
        Write-Log "ERROR: Kiro IDE not found. Please open Kiro IDE."
        return $false
    }

    Write-Log "Found Kiro IDE (PID=$($proc.Id))"

    # ── 清理多餘的舊 session ─────────────────────────────────────────────────
    Trim-OldSessions

    # ── 聚焦 Kiro 視窗 ───────────────────────────────────────────────────────
    $hwnd = $proc.MainWindowHandle
    if ([WinAPI2]::IsIconic($hwnd)) {
        [WinAPI2]::ShowWindow($hwnd, 9)
        Start-Sleep -Milliseconds 500
    }
    [WinAPI2]::SetForegroundWindow($hwnd)
    Start-Sleep -Milliseconds 800

    # ── Ctrl+Shift+L：聚焦 Kiro chat input ──────────────────────────────────
    [System.Windows.Forms.SendKeys]::SendWait("^+l")
    Start-Sleep -Milliseconds 600

    # ── 貼上訊息並送出 ───────────────────────────────────────────────────────
    [System.Windows.Forms.Clipboard]::SetText($text)
    [System.Windows.Forms.SendKeys]::SendWait("^v")
    Start-Sleep -Milliseconds 400

    [System.Windows.Forms.SendKeys]::SendWait("{ENTER}")
    Start-Sleep -Milliseconds 200

    Write-Log "Sent message to Kiro chat"
    return $true
}

# ── 主程式 ────────────────────────────────────────────────────────────────────

Write-Log "=== Kiro Auto Continue (MaxSessions=$MaxSessions) ==="

if ($RunOnce) {
    $ok = Send-ToKiro
    if ($ok) { exit 0 } else { exit 1 }
} else {
    Write-Log "Loop mode: every $IntervalHours hours. Press Ctrl+C to stop."
    Send-ToKiro
    while ($true) {
        $next = (Get-Date).AddHours($IntervalHours)
        Write-Log "Next trigger: $($next.ToString('yyyy-MM-dd HH:mm:ss'))"
        Start-Sleep -Seconds ($IntervalHours * 3600)
        Write-Log "Triggering..."
        Send-ToKiro
    }
}
