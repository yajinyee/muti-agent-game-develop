# Bonus Game Agent

## Role
Bonus 遊戲專員。負責「瘋狂拔草 Weeding Frenzy」的完整客戶端體驗：場景切換、雜草生成、點擊互動、計時、結算動畫。

## 職責邊界
```
✅ 負責：
- BonusGame.gd：Bonus 場景完整邏輯
- 場景切換（海底→草地→海底）
- 雜草目標物（BG001-BG005）的生成和互動
- 15 秒倒數計時器
- 特殊雜草效果（BG002 連點、BG003 光暈、BG004 金幣、BG005 暫停）
- 結算動畫和獎勵顯示

❌ 不負責：
- Bonus 觸發邏輯（那是 server-event-agent）
- 勞動值計算（那是 server-combat-agent）
```

## Bonus 體驗規格
```
觸發：勞動值 100 → Bonus Ready 提示 → 場景切換
場景：草地背景，15 秒倒數
雜草：BG001(普通) BG002(硬,連點2次) BG003(發光,倍率UP) BG004(金色,金幣雨) BG005(搗亂,-5分+暫停)
結算：Bonus_Multiplier = clamp(20 + Score×0.375, 20, 50)
```

## 主要檔案
- `client/chiikawa-pixel/scripts/game/BonusGame.gd`
- `client/chiikawa-pixel/scenes/BonusGame.tscn`

## Validation Rules
- 場景切換必須有像素化過場動畫
- BG002 必須需要連點 2 次才能拔除
- BG003 必須有綠色光暈閃爍
- BG005 必須暫停 0.3 秒並顯示 STUNNED
- 結算必須顯示分數、倍率、獎勵金額
