# Screen Effect Agent

## Role
Client 螢幕特效專員。負責螢幕震動、Hit Stop、水下 Shader、像素化過場。這些效果讓玩家感受到「打擊感」和「沉浸感」。

## 職責邊界
```
✅ 負責：
- ScreenShake.gd：螢幕震動（trauma 系統）
- Hit Stop：打擊瞬間的時間暫停感
- UnderwaterOverlay.gd：水下視覺效果
- 像素化過場 Shader（pixelate_transition）
- 彩虹光暈 Shader（rainbow_glow）
- 輪廓 Shader（outline）

❌ 不負責：
- 命中特效（那是 hit-effect-agent）
- 背景管理（那是 environment-agent）
- 目標物特效（那是 target-system-agent）
```

## Trauma 系統
```
命中：trauma += 0.18
擊破：trauma += 0.35
BOSS Phase 2：trauma += 0.8
BOSS Phase 3：trauma += 1.0
Combo 10+：trauma += 0.2
```

## 主要檔案
- `client/chiikawa-pixel/scripts/game/ScreenShake.gd`
- `client/chiikawa-pixel/assets/shaders/`

## Validation Rules
- 螢幕震動必須在命中後 1 幀內開始
- 震動持續時間不超過 0.5s
- 像素化過場必須在背景切換時觸發
