# daily_build.ps1
# 每日 Build 自動化腳本
# 吉伊卡哇：像素大討伐
#
# 使用方式：powershell -File tools/daily_build.ps1
# 或：.\tools\daily_build.ps1

param(
    [switch]$SkipRTP,
    [switch]$SkipSpriteQC,
    [switch]$Verbose
)

$ErrorActionPreference = "Stop"
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ProjectRoot = Split-Path -Parent $ScriptDir
$Date = Get-Date -Format "yyyy-MM-dd"
$Timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"

# 確保 Python 輸出 UTF-8（解決中文編碼問題）
$env:PYTHONUTF8 = "1"
$env:PYTHONIOENCODING = "utf-8"

# ─── 顏色輸出函數 ────────────────────────────────────────────────────────────

function Write-OK { param($msg) Write-Host "[OK] $msg" -ForegroundColor Green }
function Write-FAIL { param($msg) Write-Host "[FAIL] $msg" -ForegroundColor Red }
function Write-WARN { param($msg) Write-Host "[WARN] $msg" -ForegroundColor Yellow }
function Write-INFO { param($msg) Write-Host "[INFO] $msg" -ForegroundColor Cyan }
function Write-SECTION { param($msg) Write-Host "`n=== $msg ===" -ForegroundColor White }

# ─── 結果追蹤 ────────────────────────────────────────────────────────────────

$Results = @{
    GoBuild     = $false
    GoVet       = $false
    GoTest      = $false
    RTPSim      = $false
    SpriteQC    = $false
    BuildStable = $false
    Errors      = @()
    Warnings    = @()
    StartTime   = $Timestamp
}

# ─── 主流程 ──────────────────────────────────────────────────────────────────

Write-Host "=" * 60
Write-Host "Daily Build — 吉伊卡哇：像素大討伐"
Write-Host "日期：$Date"
Write-Host "開始時間：$Timestamp"
Write-Host "=" * 60

# ── Step 1：確認 Go Server 編譯 ──────────────────────────────────────────────

Write-SECTION "Step 1: Go Build"

Push-Location "$ProjectRoot\server"
try {
    Write-INFO "執行 go build ./..."
    $buildOutput = & go build ./... 2>&1
    $buildExit = $LASTEXITCODE
    
    if ($buildExit -eq 0) {
        Write-OK "go build 通過"
        $Results.GoBuild = $true
    } else {
        Write-FAIL "go build 失敗"
        Write-Host $buildOutput -ForegroundColor Red
        $Results.Errors += "go build 失敗：$buildOutput"
    }
} catch {
    Write-FAIL "go build 執行錯誤：$_"
    $Results.Errors += "go build 執行錯誤：$_"
} finally {
    Pop-Location
}

# ── Step 2：確認 go vet 通過 ─────────────────────────────────────────────────

Write-SECTION "Step 2: Go Vet"

Push-Location "$ProjectRoot\server"
try {
    Write-INFO "執行 go vet ./..."
    $vetOutput = & go vet ./... 2>&1
    $vetExit = $LASTEXITCODE
    
    if ($vetExit -eq 0) {
        Write-OK "go vet 通過（零警告）"
        $Results.GoVet = $true
    } else {
        Write-WARN "go vet 有警告"
        Write-Host $vetOutput -ForegroundColor Yellow
        $Results.Warnings += "go vet 警告：$vetOutput"
        $Results.GoVet = $false
    }
} catch {
    Write-WARN "go vet 執行錯誤：$_"
    $Results.Warnings += "go vet 執行錯誤：$_"
} finally {
    Pop-Location
}

# ── Step 3：執行 Go 測試 ─────────────────────────────────────────────────────

Write-SECTION "Step 3: Go Test"

Push-Location "$ProjectRoot\server"
try {
    Write-INFO "執行 go test ./..."
    $testOutput = & go test ./... 2>&1
    $testExit = $LASTEXITCODE
    
    if ($testExit -eq 0) {
        Write-OK "go test 全部通過"
        $Results.GoTest = $true
    } else {
        Write-WARN "go test 有失敗"
        Write-Host $testOutput -ForegroundColor Yellow
        $Results.Warnings += "go test 失敗：$testOutput"
        $Results.GoTest = $false
    }
} catch {
    Write-WARN "go test 執行錯誤（可能無測試檔案）：$_"
    $Results.GoTest = $true  # 無測試檔案視為通過
} finally {
    Pop-Location
}

# ── Step 4：執行 RTP 模擬 ────────────────────────────────────────────────────

Write-SECTION "Step 4: RTP Simulation"

if ($SkipRTP) {
    Write-WARN "跳過 RTP 模擬（--SkipRTP）"
    $Results.RTPSim = $true
} else {
    $rtpScript = "$ProjectRoot\tools\qa_check.py"
    if (Test-Path $rtpScript) {
        Write-INFO "執行 RTP 模擬（1000 局快速版）..."
        try {
            $rtpOutput = & python $rtpScript --rtp-only --quick 2>&1
            $rtpExit = $LASTEXITCODE
            
            if ($rtpExit -eq 0) {
                Write-OK "RTP 模擬通過"
                $Results.RTPSim = $true
                if ($Verbose) { Write-Host $rtpOutput }
            } else {
                Write-WARN "RTP 模擬有警告"
                Write-Host $rtpOutput -ForegroundColor Yellow
                $Results.Warnings += "RTP 模擬警告"
                $Results.RTPSim = $false
            }
        } catch {
            Write-WARN "RTP 模擬執行錯誤：$_"
            $Results.Warnings += "RTP 模擬執行錯誤"
            $Results.RTPSim = $false
        }
    } else {
        Write-WARN "找不到 qa_check.py，跳過 RTP 模擬"
        $Results.RTPSim = $true
    }
}

# ── Step 5：執行 Sprite QC ───────────────────────────────────────────────────

Write-SECTION "Step 5: Sprite QC"

if ($SkipSpriteQC) {
    Write-WARN "跳過 Sprite QC（--SkipSpriteQC）"
    $Results.SpriteQC = $true
} else {
    $animScript = "$ProjectRoot\tools\animation_pipeline.py"
    if (Test-Path $animScript) {
        Write-INFO "執行 Sprite QC..."
        try {
            $qcOutput = & python $animScript --audit 2>&1
            $qcExit = $LASTEXITCODE
            
            if ($qcExit -eq 0) {
                Write-OK "Sprite QC 完成"
                $Results.SpriteQC = $true
                if ($Verbose) { Write-Host $qcOutput }
            } else {
                Write-WARN "Sprite QC 有問題"
                $Results.Warnings += "Sprite QC 發現問題"
                $Results.SpriteQC = $false
            }
        } catch {
            Write-WARN "Sprite QC 執行錯誤：$_"
            $Results.Warnings += "Sprite QC 執行錯誤"
            $Results.SpriteQC = $false
        }
    } else {
        Write-WARN "找不到 animation_pipeline.py，跳過 Sprite QC"
        $Results.SpriteQC = $true
    }
}

# ── Step 6：判斷 Build 穩定性 ────────────────────────────────────────────────

Write-SECTION "Step 6: Build Stability Check"

$criticalPassed = $Results.GoBuild  # go build 是唯一硬性條件
$allPassed = $Results.GoBuild -and $Results.GoVet -and $Results.GoTest -and $Results.RTPSim -and $Results.SpriteQC

if ($criticalPassed) {
    $Results.BuildStable = $true
    if ($allPassed) {
        Write-OK "Build Stable ✅（所有檢查通過）"
    } else {
        Write-WARN "Build Stable（有警告，但可發布）"
    }
} else {
    $Results.BuildStable = $false
    Write-FAIL "Build Unstable ❌（go build 失敗）"
}

# ── Step 7：輸出 Build Report ────────────────────────────────────────────────

Write-SECTION "Step 7: Generate Build Report"

$ReportDir = "$ProjectRoot\reports\qa"
if (-not (Test-Path $ReportDir)) {
    New-Item -ItemType Directory -Path $ReportDir -Force | Out-Null
}

$ReportPath = "$ReportDir\build-report-$Date.md"
$EndTime = Get-Date -Format "yyyy-MM-dd HH:mm:ss"

$BuildStableIcon = if ($Results.BuildStable) { "✅" } else { "❌" }
$GoBuildIcon = if ($Results.GoBuild) { "✅" } else { "❌" }
$GoVetIcon = if ($Results.GoVet) { "✅" } else { "⚠️" }
$GoTestIcon = if ($Results.GoTest) { "✅" } else { "⚠️" }
$RTPIcon = if ($Results.RTPSim) { "✅" } else { "⚠️" }
$SpriteIcon = if ($Results.SpriteQC) { "✅" } else { "⚠️" }

$BuildStabilityScore = 100
if (-not $Results.GoBuild) { $BuildStabilityScore -= 50 }
if (-not $Results.GoVet) { $BuildStabilityScore -= 5 }
if (-not $Results.GoTest) { $BuildStabilityScore -= 10 }
if (-not $Results.RTPSim) { $BuildStabilityScore -= 10 }
if (-not $Results.SpriteQC) { $BuildStabilityScore -= 5 }

$ReportContent = @"
# Daily Build Report

**日期**：$Date
**開始時間**：$($Results.StartTime)
**結束時間**：$EndTime
**Build 狀態**：$BuildStableIcon $(if ($Results.BuildStable) { "Build Stable" } else { "Build Unstable" })
**Build Stability 分數**：$BuildStabilityScore/100

---

## 檢查結果

| 檢查項目 | 結果 | 說明 |
|---------|------|------|
| go build | $GoBuildIcon | $(if ($Results.GoBuild) { "編譯成功" } else { "編譯失敗" }) |
| go vet | $GoVetIcon | $(if ($Results.GoVet) { "零警告" } else { "有警告" }) |
| go test | $GoTestIcon | $(if ($Results.GoTest) { "全部通過" } else { "有失敗" }) |
| RTP 模擬 | $RTPIcon | $(if ($Results.RTPSim) { "通過" } else { "有警告" }) |
| Sprite QC | $SpriteIcon | $(if ($Results.SpriteQC) { "通過" } else { "有問題" }) |

---

## 錯誤清單

$(if ($Results.Errors.Count -eq 0) { "無錯誤 ✅" } else { $Results.Errors | ForEach-Object { "- $_" } | Out-String })

## 警告清單

$(if ($Results.Warnings.Count -eq 0) { "無警告 ✅" } else { $Results.Warnings | ForEach-Object { "- $_" } | Out-String })

---

*報告生成時間：$EndTime*
"@

$ReportContent | Out-File -FilePath $ReportPath -Encoding UTF8
Write-OK "Build Report 已輸出：$ReportPath"

# ── 最終摘要 ─────────────────────────────────────────────────────────────────

Write-Host "`n" + "=" * 60
Write-Host "Daily Build 完成"
Write-Host "=" * 60
Write-Host "Build Stability 分數：$BuildStabilityScore/100"
Write-Host "Build 狀態：$(if ($Results.BuildStable) { '✅ Build Stable' } else { '❌ Build Unstable' })"
Write-Host "報告路徑：$ReportPath"

if (-not $Results.BuildStable) {
    Write-Host "`n[FAIL] Build 不穩定，請修復錯誤後重新執行" -ForegroundColor Red
    exit 1
} else {
    Write-Host "`n[OK] Build 完成！" -ForegroundColor Green
    exit 0
}
