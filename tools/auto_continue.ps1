param(
    [int]$IntervalHours = 2,
    [switch]$RunOnce
)

$LOG_FILE = "d:\Kiro\tools\auto_continue.log"
$MSG_FILE = "d:\Kiro\tools\auto_continue_message.txt"

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

function Send-ToKiro {
    # 從檔案讀取訊息（支援中文）
    $text = [System.IO.File]::ReadAllText($MSG_FILE, [System.Text.Encoding]::UTF8).Trim()

    $proc = Get-Process -Name "Kiro" -ErrorAction SilentlyContinue | Where-Object { $_.MainWindowTitle -ne "" } | Select-Object -First 1
    if (-not $proc) { Write-Log "ERROR: Kiro IDE not found. Please open Kiro IDE."; return $false }

    $hwnd = $proc.MainWindowHandle
    Write-Log "Found Kiro IDE (PID=$($proc.Id))"

    if ([WinAPI2]::IsIconic($hwnd)) { [WinAPI2]::ShowWindow($hwnd, 9); Start-Sleep -Milliseconds 500 }
    [WinAPI2]::SetForegroundWindow($hwnd)
    Start-Sleep -Milliseconds 800

    # Ctrl+Shift+L: focus Kiro chat input (official shortcut)
    [System.Windows.Forms.SendKeys]::SendWait("^+l")
    Start-Sleep -Milliseconds 600

    # Paste via clipboard (handles Chinese/Unicode correctly)
    [System.Windows.Forms.Clipboard]::SetText($text)
    [System.Windows.Forms.SendKeys]::SendWait("^v")
    Start-Sleep -Milliseconds 400

    # Enter to send
    [System.Windows.Forms.SendKeys]::SendWait("{ENTER}")
    Start-Sleep -Milliseconds 200

    Write-Log "Sent message to Kiro chat"
    return $true
}

Write-Log "=== Kiro Auto Continue ==="
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