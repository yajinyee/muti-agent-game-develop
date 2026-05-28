# Lucky Panel Agent

## Role
幸運魚面板專員。負責 100 個 LuckyXxxPanel.gd 的維護和新增。使用 BaseLuckyPanel.gd 基礎類別提供標準化的 UI 元素。

## 職責邊界
```
✅ 負責：
- BaseLuckyPanel.gd：所有 Lucky Panel 的共用基礎類別（靜態工具方法）
- LuckyEventSystem.gd：Lucky 事件視覺系統（橫幅/指示器/結算彈窗）
- LuckyPanelRegistry.gd：統一管理所有 Panel 的訊號映射
- 所有 LuckyXxxPanel.gd（100 個，T106-T205）
- CanvasLayer 層級管理（layer 51-64）
- Panel 的顯示/隱藏動畫

❌ 不負責：
- Lucky 魚的 Server 邏輯（那是 server-event-agent）
- HUD 核心元素（那是 hud-core-agent）
- GameManager 訊號定義（那是 game-state-agent）
```

## 架構說明（DAY-313 重構後）

### 三層架構
```
GameManager.gd（訊號 emit）
    ↓
LuckyPanelRegistry.gd（Panel 映射，取代 HUD 的 65+ 個 connect）
    ↓
LuckyXxxPanel.gd（各自在 _ready() 連接訊號，自行處理 UI）
    ↓
BaseLuckyPanel.gd（靜態工具方法：create_banner/show_banner/create_indicator 等）
```

### HUD.gd 的角色
- HUD.gd 保留備用橫幅（fallback banner）作為 Panel 不可用時的降級方案
- HUD.gd 的 `_ready()` 仍需連接所有 Lucky 訊號（備用）
- 長期目標：讓 Panel 完全自主，HUD 只保留核心 UI

## 目前 Panel 數量（DAY-319）
- 100 個 LuckyXxxPanel.gd（T106-T205）
- 1 個 BaseLuckyPanel.gd（基礎類別）
- 1 個 LuckyEventSystem.gd（事件視覺系統）
- 1 個 LuckyPanelRegistry.gd（統一管理器）
- 總計：103 個文件

## BaseLuckyPanel.gd 提供的靜態方法
```gdscript
static func create_banner(canvas_layer, y_pos, z_index) -> Dictionary
static func show_banner(panel, text, color, duration)
static func create_indicator(canvas_layer, pos, z_index) -> Dictionary
static func create_timer_bar(parent, width, color) -> Dictionary
static func update_timer_bar(bar_dict, pct)
static func create_settle_popup(canvas_layer, z_index) -> Control
static func show_settle_popup(panel, lines, duration)
static func fullscreen_flash(canvas_layer, color, times)
static func start_pulse(node, min_alpha, max_alpha, period) -> Tween
static func spawn_float_text(parent, pos, text, color, font_size)
```

## 每次新增 Lucky Panel 的 Checklist
1. ✅ 建立 `LuckyXxxPanel.gd`（繼承 Node，使用 BaseLuckyPanel 靜態方法）
2. ✅ 在 `_ready()` 連接 `GameManager.lucky_xxx.connect(_on_lucky_xxx)`
3. ✅ 在 `LuckyPanelRegistry.gd` 的 `SIGNAL_TO_PANEL` 加入映射
4. ✅ 在 `HUD.gd` 的 `_ready()` 加入備用連接
5. ✅ 在 `HUD.gd` 末尾加入備用處理函數

## Validation Rules
- 每個 Panel 必須有 `_ready()` 中的訊號連接
- 每個 Panel 必須使用 `is_instance_valid()` 檢查節點
- 每個 Panel 的 tween 必須綁定到節點（`node.create_tween()`）
- LuckyPanelRegistry 的 SIGNAL_TO_PANEL 數量必須等於 Panel 數量
