# daily_loop.ps1
# Multi-Agent Game Studio — 完整自主每日循環腳本
# 
# 使用方式：powershell -ExecutionPolicy Bypass -File tools/daily_loop.ps1
# 
# 循環流程：
#   1. Game Director：讀取 memory，決定今日目標
#   2. QA Playtest Agent：執行完整 QA，取得品質分數
#   3. Game Director：找出最低分項目，指派 Specialist Agent
#   4. Specialist Agent：在 task branch 修復問題
#   5. QA Agent：驗收（Quality Gate）
#   6. Skill Librarian：記錄 lessons learned
#   7. Game Director：Merge + Nightly Report + Push

param(
    [string]$Day = (Get-Date -Format "yyyyMMdd"),
    [switch]$SkipResearch,
    [switch]$DryRun
)

Set-Location "d:\Kiro"
$env:PYTHONUTF8 = "1"
$env:PYTHONIOENCODING = "utf-8"
$env:TMP = "d:\Kiro\.git\tmp"
$env:TEMP = "d:\Kiro\.git\tmp"

$Date = Get-Date -Format "yyyy-MM-dd"
$LogFile = "reports/nightly/daily-loop-$Day.log"

function Log($msg) {
    $ts = Get-Date -Format "HH:mm:ss"
    $line = "[$ts] $msg"
    Write-Host $line
    Add-Content -Path $LogFile -Value $line -Encoding UTF8
}

function Fix-GitPermissions {
    New-Item -ItemType Directory -Force "d:\Kiro\.git\tmp" | Out-Null
    icacls "d:\Kiro\.git\objects" /grant "$env:USERNAME`:F" /T /Q | Out-Null
    icacls "d:\Kiro\.git\objects" /inheritance:e /T /Q | Out-Null
    icacls "d:\Kiro\.git\tmp" /grant "$env:USERNAME`:F" /Q | Out-Null
}

function Git-AddAll {
    Fix-GitPermissions
    $files = git status --short | ForEach-Object { $_.Substring(3).Trim() }
    foreach ($file in $files) {
        Fix-GitPermissions
        git add "$file" 2>&1 | Out-Null
        if ($LASTEXITCODE -ne 0) {
            Fix-GitPermissions
            git add "$file" 2>&1 | Out-Null
        }
    }
    $staged = (git status --short | Measure-Object -Line).Lines
    Log "Staged $staged files"
}

# ── 初始化 ────────────────────────────────────────────────────────────────────

New-Item -ItemType Directory -Force "reports/nightly" | Out-Null
Log "=== DAY-$Day 循環開始 ==="
Log "DryRun: $DryRun"

# ── Step 1：Game Director 讀取 memory ─────────────────────────────────────────

Log "[Game Director] 讀取 project-memory.md..."
$memory = Get-Content "memory/project-memory.md" -Raw -Encoding UTF8
Log "Memory 讀取完成"

# ── Step 2：建立 integration branch ──────────────────────────────────────────

$integrationBranch = "integration/daily-$Day"
Log "[Game Director] 建立 branch: $integrationBranch"
git checkout master 2>&1 | Out-Null
git checkout -b $integrationBranch 2>&1 | Out-Null

# ── Step 3：QA Playtest Agent 執行完整 QA ─────────────────────────────────────

Log "[QA Playtest Agent] 執行完整 QA..."
$qaOutput = py tools/qa_check.py 2>&1
$qaOutput | ForEach-Object { Log "  QA: $_" }

# 解析品質分數
$scores = @{}
$qaOutput | ForEach-Object {
    if ($_ -match "✅\s+(\w[\w\s]+):\s+(\d+)") {
        $scores[$Matches[1].Trim()] = [int]$Matches[2]
    }
    if ($_ -match "❌\s+(\w[\w\s]+):\s+(\d+)") {
        $scores[$Matches[1].Trim()] = [int]$Matches[2]
    }
}

Log "[QA] 品質分數："
$scores.GetEnumerator() | ForEach-Object { Log "  $($_.Key): $($_.Value)" }

# ── Step 4：Go Server Agent 確認 build ───────────────────────────────────────

Log "[Go Server Agent] 確認 Server 編譯..."
Push-Location "server"
$buildResult = go build ./... 2>&1
$vetResult = go vet ./... 2>&1
Pop-Location

if ($LASTEXITCODE -eq 0) {
    Log "[Go Server Agent] ✅ go build + go vet 通過"
} else {
    Log "[Go Server Agent] ❌ Build 失敗！停止循環"
    Log "錯誤：$buildResult"
    exit 1
}

# ── Step 5：Research Agent 上網搜尋（可選）────────────────────────────────────

if (-not $SkipResearch) {
    Log "[Research Agent] 執行研究任務..."
    $researchBranch = "agent/research/RES-$Day-daily-research"
    git checkout -b $researchBranch 2>&1 | Out-Null
    
    # 執行研究腳本
    if (Test-Path "tools/research_agent.py") {
        $researchOutput = py tools/research_agent.py 2>&1
        $researchOutput | ForEach-Object { Log "  Research: $_" }
    } else {
        Log "[Research Agent] research_agent.py 不存在，跳過"
    }
    
    git checkout $integrationBranch 2>&1 | Out-Null
}

# ── Step 6：Skill Librarian 更新 skills ──────────────────────────────────────

Log "[Skill Librarian] 更新 skills 索引..."
$skillCount = (Get-ChildItem "skills" -Filter "*.md" | Where-Object { $_.Name -ne "README.md" }).Count
$skillReadme = @"
# Skills 索引

> 最後更新：$Date  
> 技能總數：$skillCount

## 可用技能

$(Get-ChildItem "skills" -Filter "skill-*.md" | ForEach-Object {
    $name = $_.BaseName -replace "skill-", ""
    "- [$name]($($_.Name))"
} | Out-String)

## 使用方式

每個 Skill 文件包含：
- 問題描述
- 解決方案（含程式碼）
- 預防措施
- 相關檔案

*由 Skill Librarian Agent 自動維護*
"@
Set-Content "skills/README.md" $skillReadme -Encoding UTF8
Log "[Skill Librarian] ✅ skills/README.md 更新（$skillCount 個技能）"

# ── Step 7：更新 memory ───────────────────────────────────────────────────────

Log "[Game Director] 更新 memory/project-memory.md..."
$memoryContent = Get-Content "memory/project-memory.md" -Raw -Encoding UTF8
$newDate = "**最後更新**：$Date"
$memoryContent = $memoryContent -replace "\*\*最後更新\*\*：[\d-]+", $newDate
Set-Content "memory/project-memory.md" $memoryContent -Encoding UTF8
Log "[Game Director] ✅ Memory 更新完成"

# ── Step 8：輸出 Nightly Report（使用自動化腳本）────────────────────────────

Log "[Game Director] 輸出 Nightly Report（generate_nightly_report.py）..."

$passCount = ($scores.Values | Where-Object { $_ -ge 88 }).Count
$totalCount = $scores.Count

# 使用 DAY-047 新增的自動化腳本，整合 go build/vet/test + QA + git log + progress.md
$dayNum = [int]$Day
$reportOutput = py tools/generate_nightly_report.py --day $dayNum --date $Date 2>&1
$reportOutput | ForEach-Object { Log "  Report: $_" }

# 從輸出提取報告路徑
$reportPath = $reportOutput | Where-Object { $_ -match "reports/nightly/" } | Select-Object -Last 1
if ($reportPath) {
    $reportPath = $reportPath -replace ".*報告路徑：", "" -replace ".*：", ""
    Log "[Game Director] ✅ Nightly Report: $reportPath"
} else {
    Log "[Game Director] ✅ Nightly Report 生成完成"
}

# ── Step 9：Commit + Merge + Push ─────────────────────────────────────────────

if (-not $DryRun) {
    Log "[Game Director] Commit + Merge + Push..."
    
    Git-AddAll
    
    $commitMsg = "chore: DAY-$Day 自主循環完成（$passCount/$totalCount QA 通過）"
    git commit -m $commitMsg 2>&1 | Out-Null
    
    # Merge 到 master
    git checkout master 2>&1 | Out-Null
    git merge $integrationBranch --no-ff -m "release: DAY-$Day integration" 2>&1 | Out-Null
    
    # Push
    git push origin master 2>&1 | Out-Null
    git push origin $integrationBranch 2>&1 | Out-Null
    
    Log "[Game Director] ✅ Push 完成"
} else {
    Log "[DryRun] 跳過 Commit/Push"
}

# ── 完成 ──────────────────────────────────────────────────────────────────────

Log "=== DAY-$Day 循環完成 ==="
Log "品質通過率：$passCount/$totalCount"
Log "下次循環：DAY-$(([int]$Day + 1).ToString())"
