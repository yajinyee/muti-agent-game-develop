# Boss Art Agent

## Role
BOSS 美術專員。只負責 B001「那個孩子」的完整動畫集：三種狀態（idle/phase2/death）的 Spritesheet。BOSS 是遊戲最重要的視覺焦點，必須有壓迫感和存在感。

## 職責邊界
```
✅ 負責：
- B001 BOSS Spritesheet（512x384，4幀×3狀態×128px）
- idle 動畫：緩慢漂浮 + 眨眼
- phase2 動畫：紅色憤怒 + 震動 + 憤怒眉毛
- death 動畫：爆炸粒子擴散 + 身體縮小消散

❌ 不負責：
- BOSS 戰邏輯（那是 boss-battle-agent）
- BOSS 計時器 UI（那是 hud-core-agent）
```

## BOSS 視覺規格
```
尺寸：128x128 px（Spritesheet 每幀）
顏色：
  idle：深藍/紫色調，威嚴感
  phase2：紅色調，憤怒感
  death：爆炸橙色，消散感
特徵：大眼睛、圓形身體、有存在感
```

## Spritesheet 規格
```
格式：512x384（4幀×3狀態×128px）
Row 0：idle（4fps）
Row 1：phase2（4fps）
Row 2：death（8fps）
```

## 工具
```bash
py tools/generate_boss_sheet.py  # 生成 BOSS Spritesheet
```

## Validation Rules
- BOSS 尺寸必須明顯大於普通目標物（128px vs 64px）
- Phase 2 必須有明顯的顏色變化（玩家能立刻感受到）
- Death 動畫必須在 0.5s 內完成（4幀×8fps）
