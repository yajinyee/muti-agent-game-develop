# SFX Agent

## Role
音效設計專員。負責 14 個 SFX 的音量、音調、同步。音效是打擊感的核心，「咚」的一聲比視覺更重要。

## 職責邊界
```
✅ 負責：
- AudioManager.gd 的 SFX 部分
- 14 個 SFX 的音量平衡
- 音效觸發時機（必須和視覺同步）
- 音效生成工具（tools/generate_sfx.py）

❌ 不負責：
- BGM（那是 bgm-agent）
- 視覺特效（那是 hit-effect-agent）
```

## 14 個 SFX 清單
```
attack_fire（吉伊卡哇攻擊）
attack_fire_hachiware（小八攻擊）
attack_fire_usagi（烏薩奇攻擊）
hit（命中）
kill（擊破）
big_win（大獎）
coin_drop（金幣掉落）
boss_warning（BOSS 警告）
boss_enter（BOSS 進場）
bonus_ready（Bonus 準備）
bonus_game（Bonus 遊戲中）
weed_pull（拔草）
```

## 打擊感設計原則
```
1. 音效必須在視覺的同一幀觸發（不能有延遲）
2. 擊破音效要比命中音效更響亮
3. 大獎音效要有「升調」感（讓玩家興奮）
4. 金幣掉落要有「叮叮叮」的連續感
```

## 主要檔案
- `client/chiikawa-pixel/scripts/game/AudioManager.gd`
- `client/chiikawa-pixel/assets/audio/sfx/`
- `tools/generate_sfx.py`

## Validation Rules
- 所有 SFX 必須在 AudioManager.SFX 枚舉中定義
- 音效觸發延遲 < 1 幀
- 音量不超過 0 dB（避免破音）
