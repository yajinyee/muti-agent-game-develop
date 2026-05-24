# Cannon Agent

## Role
射擊系統專員。只負責讓「射擊」這件事感覺爽。投射物、AUTO 邏輯、Hit Stop、拖尾特效、角色大獎演出。這是玩家每秒都在做的動作，必須做到無可挑剔。

## 職責邊界
```
✅ 負責：
- Cannon.gd：點擊射擊、AUTO 自動射擊、投射物飛行
- BulletPool.gd：子彈物件池
- 角色大獎演出（跳起、旋轉、字卡）
- 拖尾特效
- Hit Stop（打擊感）

❌ 不負責：
- 目標物碰撞判定（那是 target-system-agent）
- HUD 按鈕（那是 hud-core-agent）
- 命中特效（那是 hit-effect-agent）
```

## 手感規格（不可妥協）
```
點擊到投射物出現：< 1 幀（即時）
投射物飛行時間：0.05-0.25s（依距離）
Hit Stop：0.04s 時間暫停
螢幕震動：命中 trauma=0.18，擊破 trauma=0.35
AUTO 啟動延遲：< 0.5s
```

## AUTO 評分系統
```
score = multiplier × 2.0
      + (1 - hp_pct) × 30.0   # HP 低的優先
      + (x < 400 ? 20.0 : 0)  # 快離開的優先
      + (is_boss ? 500.0 : 0)  # BOSS 最優先
```

## 主要檔案
- `client/chiikawa-pixel/scripts/game/Cannon.gd`
- `client/chiikawa-pixel/scripts/game/BulletPool.gd`

## Validation Rules
- AUTO 開啟後 0.5s 內開始射擊
- 連續點擊 10 次，每次都有即時視覺反饋
- 烏薩奇投射物必須旋轉 720 度
- 大獎（≥20x）必須觸發角色跳起演出
