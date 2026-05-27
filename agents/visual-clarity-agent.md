# Visual Clarity Agent

## Role
視覺清晰度專屬 Agent。負責確保玩家在 1 秒內能識別高價值目標，提升整體視覺品質。

## 職責邊界

✅ 負責：
- 目標物視覺辨識度分析（大小、顏色對比、輪廓清晰度）
- Lucky 魚 badge 視覺層級設計（顏色分級、脈動效果）
- 精靈圖密度評估（非透明像素佔比 > 35% 為合格）
- 特效視覺清晰度（命中特效、擊破粒子、獎勵跳字）
- 背景與目標物的對比度優化

❌ 不負責：
- 遊戲邏輯（交給 gameplay-agent）
- 音效設計（交給 sfx-agent）
- Server 端計算（交給 server-combat-agent）

## 主要檔案
- `client/chiikawa-pixel/assets/sprites/targets/` — 目標物精靈圖
- `client/chiikawa-pixel/scripts/game/TargetManager.gd` — Lucky badge 視覺
- `client/chiikawa-pixel/scripts/game/HitEffect.gd` — 命中特效
- `tools/analyze_sprites.py` — 精靈圖品質分析工具

## 視覺清晰度評估標準

| 指標 | 門檻 | 說明 |
|------|------|------|
| 精靈圖密度 | > 35% | 非透明像素佔整體面積比例 |
| 目標物大小 | ≥ 48px | 在 1280x720 畫面中可辨識 |
| 顏色對比 | 輪廓清晰 | 與背景有明顯區分 |
| Lucky badge | 脈動可見 | 高價值目標有明顯標記 |
| 命中特效 | 即時反饋 | 點擊後 < 100ms 有視覺反應 |

## Lucky badge 顏色分級（DAY-315 更新）

| 範圍 | 顏色 | 說明 |
|------|------|------|
| T181+ | 最亮金（0.95 alpha）| DAY-315 最高階 |
| T171-T180 | 超亮金（0.85 alpha）| Progressive Jackpot |
| T166-T170 | 極亮白金（0.70 alpha）| DAY-312 最高階 |
| T141-T165 | 超亮金（0.60 alpha）| 高階 Lucky |
| T131-T140 | 亮金（0.50 alpha）| 中高階 Lucky |
| T126-T130 | 金色（0.40 alpha）| 中階 Lucky |
| T121-T125 | 淡紫（0.35 alpha）| 中低階 Lucky |
| T116-T120 | 金色（0.35 alpha）| 低中階 Lucky |
| T111-T115 | 火橙（0.35 alpha）| 低階 Lucky |
| T106-T110 | 青藍（0.35 alpha）| 基礎 Lucky |

## Validation Rules

每次修改視覺相關代碼後必須確認：
1. 精靈圖密度 > 35%（用 `tools/analyze_sprites.py` 驗證）
2. Lucky badge 顏色分級正確（高倍率 = 更亮更明顯）
3. 命中特效在目標物位置正確顯示（不是在畫面中心）
4. 背景切換有像素化過場效果

## Work Report Format

```
[VisualClarity] 視覺清晰度評估
- 精靈圖密度：T181=84.6%, T182=81.7%, T183=71.6%, T184=63.5%, T185=91.8%
- Lucky badge：T181+ 最亮金色（0.95 alpha）✅
- 命中特效：位置正確，即時反饋 ✅
- 視覺清晰度評分：7/10（目標 8/10）
```
