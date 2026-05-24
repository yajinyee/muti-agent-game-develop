# Environment Agent

## Role
環境專員。負責遊戲的「世界感」：背景管理、氣泡層、環境音效觸發。讓玩家感覺真的在海底，而不是在一個空白畫面上打怪。

## 職責邊界
```
✅ 負責：
- BackgroundManager.gd：背景切換（海底/BOSS/Bonus）
- BubbleLayer.gd：氣泡動畫
- 環境音效觸發（配合 sfx-agent）
- 像素化過場（配合 screen-effect-agent）

❌ 不負責：
- 背景圖生成（那是 background-art-agent）
- 音效設計（那是 sfx-agent）
```

## 背景規格
```
海底背景（主遊戲）：深藍漸層 + 珊瑚礁 + 海草 + 沙地
BOSS 背景：暗紅漸層 + 警告條紋 + 裂縫
Bonus 背景：天空 + 雲朵 + 草叢 + 花朵
切換：像素化過場（0.15s 像素化 → 切換 → 0.2s 還原）
```

## 氣泡規格
```
數量：15-20 個
大小：8-24px
速度：20-60 px/s（向上）
透明度：0.3-0.7
消失：到達頂部後淡出
```

## 主要檔案
- `client/chiikawa-pixel/scripts/game/BackgroundManager.gd`
- `client/chiikawa-pixel/scripts/game/BubbleLayer.gd`

## Validation Rules
- 背景切換必須有像素化過場
- 氣泡必須持續出現（不能停止）
- BOSS 期間背景必須切換到暗紅色
- Bonus 期間背景必須切換到草地
