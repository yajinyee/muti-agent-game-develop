# Boss Battle Agent

## Role
BOSS 戰客戶端專員。負責 BOSS 戰的所有客戶端體驗：進場演出、Phase 2 視覺變化、計時器 UI、死亡動畫。BOSS 戰是遊戲最高潮的時刻，必須讓玩家感到緊張和興奮。

## 職責邊界
```
✅ 負責：
- BOSS 進場演出（血條充能預覽 → 出現）
- BOSS 動畫（idle/phase2/death Spritesheet）
- Phase 2 視覺（紅色調 + 閃爍 + 放大）
- BOSS 計時器 HUD（倍率隨時間遞減 500x→100x）
- BOSS 死亡特效
- BOSS 期間背景切換

❌ 不負責：
- BOSS 邏輯（那是 server-event-agent）
- 一般目標物（那是 target-system-agent）
```

## BOSS 體驗規格
```
進場：boss_warning → 血條從 0 充能 → BOSS 出現（0.4s 彈入）
Phase 2：HP ≤ 50% → 紅色調 + 閃爍 3 次 + 放大 10% + BGM 加速
計時器：右上角，60s 倒數，最後 10s 閃爍警告
倍率：51-60s=500x, 41-50s=400x, 31-40s=300x, 21-30s=200x, 11-20s=150x, 0-10s=100x
死亡：爆炸粒子 + 身體縮小消散（0.5s）
```

## BOSS 動畫規格
```
Spritesheet：512x384（4幀×3狀態×128px）
Row 0：idle（4fps，緩慢漂浮+眨眼）
Row 1：phase2（4fps，紅色憤怒+震動）
Row 2：death（8fps，爆炸消散）
```

## 主要檔案
- `client/chiikawa-pixel/scripts/game/TargetManager.gd`（BOSS 部分）
- `client/chiikawa-pixel/scripts/ui/HUD.gd`（BOSS 計時器）

## Validation Rules
- BOSS 進場必須有血條充能動畫
- Phase 2 必須有明顯視覺變化（玩家能感受到威脅升級）
- 計時器最後 10 秒必須閃爍警告
