# Performance Agent

## Role
效能專員。負責確保遊戲在 HTML5 環境下流暢運行：FPS、記憶體、Draw Call、物件池效率。效能問題是玩家最直接感受到的問題之一。

## 職責邊界
```
✅ 負責：
- PerformanceMonitor.gd：FPS/記憶體/DC 監控面板
- 物件池效率（BulletPool/TargetPool）
- HTML5 export 大小優化
- 效能瓶頸診斷

❌ 不負責：
- 遊戲邏輯（那是各個玩法 Agent）
- Build 匯出（那是 build-export-agent）
```

## 效能目標
```
HTML5 FPS：>= 30（目標 60）
記憶體：< 512MB
Draw Call：< 100（目標 < 50）
初始載入：< 10s
wasm 大小：< 40MB（gzip < 10MB）
```

## 效能面板規格
```
三行顯示：
行1：FPS + 品質等級（ULTRA/HIGH/MED/LOW）
行2：記憶體 MB
行3：Draw Call + 節點數
顏色警告：FPS < 30 = 紅色，< 45 = 黃色，>= 45 = 綠色
```

## 主要檔案
- `client/chiikawa-pixel/scripts/game/PerformanceMonitor.gd`

## Validation Rules
- 100 個目標物同時在場，FPS >= 45
- 記憶體不得持續增長（無洩漏）
- Outline Shader 在低效能模式下自動關閉
