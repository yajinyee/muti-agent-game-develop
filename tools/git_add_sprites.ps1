# 逐一 stage 所有修改的 sprite 檔案
$files = @(
    "client/chiikawa-pixel/assets/sprites/backgrounds/boss_bg.png",
    "client/chiikawa-pixel/assets/sprites/sheets/targets_sheet.json",
    "client/chiikawa-pixel/assets/sprites/sheets/targets_sheet.png",
    "client/chiikawa-pixel/assets/sprites/targets/B001_boss.png",
    "client/chiikawa-pixel/assets/sprites/targets/BG001_weed_normal.png",
    "client/chiikawa-pixel/assets/sprites/targets/BG002_weed_hard.png",
    "client/chiikawa-pixel/assets/sprites/targets/BG003_weed_glow.png",
    "client/chiikawa-pixel/assets/sprites/targets/BG004_weed_gold.png",
    "client/chiikawa-pixel/assets/sprites/targets/BG005_weed_evil.png",
    "client/chiikawa-pixel/assets/sprites/targets/T001_grass.png",
    "client/chiikawa-pixel/assets/sprites/targets/T001_grass_swim.png",
    "client/chiikawa-pixel/assets/sprites/targets/T002_bug_g.png",
    "client/chiikawa-pixel/assets/sprites/targets/T002_bug_g_swim.png",
    "client/chiikawa-pixel/assets/sprites/targets/T003_bug_r.png",
    "client/chiikawa-pixel/assets/sprites/targets/T003_bug_r_swim.png",
    "client/chiikawa-pixel/assets/sprites/targets/T004_bug_b.png",
    "client/chiikawa-pixel/assets/sprites/targets/T004_bug_b_swim.png",
    "client/chiikawa-pixel/assets/sprites/targets/T005_pudding.png",
    "client/chiikawa-pixel/assets/sprites/targets/T005_pudding_swim.png",
    "client/chiikawa-pixel/assets/sprites/targets/T006_mushroom.png",
    "client/chiikawa-pixel/assets/sprites/targets/T006_mushroom_swim.png",
    "client/chiikawa-pixel/assets/sprites/targets/T101_mimic.png",
    "client/chiikawa-pixel/assets/sprites/targets/T101_mimic_swim.png",
    "client/chiikawa-pixel/assets/sprites/targets/T102_chest_swim.png",
    "client/chiikawa-pixel/assets/sprites/targets/T103_meteor.png",
    "client/chiikawa-pixel/assets/sprites/targets/T103_meteor_swim.png",
    "client/chiikawa-pixel/assets/sprites/targets/T104_gold_grass.png",
    "client/chiikawa-pixel/assets/sprites/targets/T104_gold_grass_swim.png",
    "client/chiikawa-pixel/assets/sprites/targets/T105_coin_fish.png",
    "client/chiikawa-pixel/assets/sprites/targets/T105_coin_fish_swim.png",
    "docs/progress.md",
    "tools/analyze_art_quality.py",
    "tools/enhance_targets_v2.py",
    "tools/enhance_remaining.py",
    "tools/enhance_swim_frames.py"
)

$success = 0
$failed = 0

foreach ($f in $files) {
    $result = git add $f 2>&1
    if ($LASTEXITCODE -eq 0) {
        $success++
    } else {
        Write-Host "FAILED: $f - $result"
        $failed++
        # 重試一次
        Start-Sleep -Milliseconds 200
        git add $f 2>&1 | Out-Null
        if ($LASTEXITCODE -eq 0) {
            $success++
            $failed--
            Write-Host "  RETRY OK: $f"
        }
    }
}

Write-Host "Done: $success OK, $failed failed"
