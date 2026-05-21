## PixelTheme.gd
## 像素風格 UI Theme 生成器
## 用 GDScript 動態建立 Theme，讓所有 HUD 元素套用一致的像素風格
## 使用方式：var theme = PixelTheme.create()

extends RefCounted
class_name PixelTheme

## 建立並返回像素風格 Theme
static func create() -> Theme:
	var theme = Theme.new()

	# ── 字體大小 ──────────────────────────────────────────
	theme.set_default_font_size(14)

	# ── Button 樣式 ──────────────────────────────────────
	_setup_button(theme)

	# ── ProgressBar 樣式 ─────────────────────────────────
	_setup_progress_bar(theme)

	# ── Label 樣式 ───────────────────────────────────────
	_setup_label(theme)

	return theme

## Button 像素風格
static func _setup_button(theme: Theme) -> void:
	# Normal 狀態（深藍底 + 亮邊框）
	var normal = _make_pixel_panel(
		Color(0.05, 0.10, 0.28, 0.92),   # 背景：深海藍
		Color(0.30, 0.55, 0.90, 0.85),   # 邊框：亮藍
		Color(0.15, 0.25, 0.50, 0.60),   # 陰影：暗藍
		2
	)
	theme.set_stylebox("normal", "Button", normal)

	# Hover 狀態（稍亮）
	var hover = _make_pixel_panel(
		Color(0.10, 0.18, 0.40, 0.95),
		Color(0.50, 0.75, 1.00, 0.90),
		Color(0.20, 0.35, 0.65, 0.65),
		2
	)
	theme.set_stylebox("hover", "Button", hover)

	# Pressed 狀態（下沉感：背景更暗，邊框金色）
	var pressed = _make_pixel_panel(
		Color(0.03, 0.07, 0.20, 0.95),
		Color(0.90, 0.75, 0.20, 0.90),   # 金色邊框（按下時）
		Color(0.10, 0.18, 0.38, 0.70),
		2
	)
	theme.set_stylebox("pressed", "Button", pressed)

	# Focus 狀態（金色外框）
	var focus = StyleBoxFlat.new()
	focus.bg_color = Color(0, 0, 0, 0)
	focus.border_color = Color(0.90, 0.75, 0.20, 0.80)
	focus.border_width_left = 1
	focus.border_width_right = 1
	focus.border_width_top = 1
	focus.border_width_bottom = 1
	theme.set_stylebox("focus", "Button", focus)

	# 文字顏色
	theme.set_color("font_color", "Button", Color(0.90, 0.95, 1.00))
	theme.set_color("font_hover_color", "Button", Color(1.00, 1.00, 1.00))
	theme.set_color("font_pressed_color", "Button", Color(1.00, 0.90, 0.30))
	theme.set_color("font_disabled_color", "Button", Color(0.50, 0.50, 0.55))

## ProgressBar 像素風格（勞動值條）
static func _setup_progress_bar(theme: Theme) -> void:
	# 背景（深色凹槽）
	var bg = StyleBoxFlat.new()
	bg.bg_color = Color(0.05, 0.05, 0.12, 0.90)
	bg.border_color = Color(0.20, 0.30, 0.55, 0.80)
	bg.border_width_left = 1
	bg.border_width_right = 1
	bg.border_width_top = 1
	bg.border_width_bottom = 1
	bg.corner_radius_top_left = 2
	bg.corner_radius_top_right = 2
	bg.corner_radius_bottom_left = 2
	bg.corner_radius_bottom_right = 2
	theme.set_stylebox("background", "ProgressBar", bg)

	# 填充（漸層綠→黃→紅，由 GDScript 動態控制顏色）
	var fill = StyleBoxFlat.new()
	fill.bg_color = Color(0.20, 0.85, 0.30, 0.95)  # 預設綠色
	fill.corner_radius_top_left = 2
	fill.corner_radius_top_right = 2
	fill.corner_radius_bottom_left = 2
	fill.corner_radius_bottom_right = 2
	theme.set_stylebox("fill", "ProgressBar", fill)

## Label 像素風格
static func _setup_label(theme: Theme) -> void:
	theme.set_color("font_color", "Label", Color(0.92, 0.95, 1.00))
	theme.set_color("font_shadow_color", "Label", Color(0.0, 0.0, 0.1, 0.6))
	theme.set_constant("shadow_offset_x", "Label", 1)
	theme.set_constant("shadow_offset_y", "Label", 1)

## 建立像素風格面板 StyleBoxFlat
## bg: 背景色, border: 邊框色, shadow: 陰影色, bw: 邊框寬度
static func _make_pixel_panel(bg: Color, border: Color, shadow: Color, bw: int) -> StyleBoxFlat:
	var sb = StyleBoxFlat.new()
	sb.bg_color = bg
	sb.border_color = border
	sb.border_width_left = bw
	sb.border_width_right = bw
	sb.border_width_top = bw
	sb.border_width_bottom = bw
	# 像素風格：不圓角（方形邊框）
	sb.corner_radius_top_left = 0
	sb.corner_radius_top_right = 0
	sb.corner_radius_bottom_left = 0
	sb.corner_radius_bottom_right = 0
	# 陰影（右下偏移，增加立體感）
	sb.shadow_color = shadow
	sb.shadow_size = 2
	sb.shadow_offset = Vector2(1, 1)
	# 內邊距
	sb.content_margin_left = 6
	sb.content_margin_right = 6
	sb.content_margin_top = 3
	sb.content_margin_bottom = 3
	return sb

## 建立像素風格 HUD 面板背景（用於 TopBar / BottomBar）
static func make_hud_panel(alpha: float = 0.88) -> StyleBoxFlat:
	var sb = StyleBoxFlat.new()
	sb.bg_color = Color(0.03, 0.06, 0.18, alpha)
	sb.border_color = Color(0.20, 0.35, 0.65, 0.70)
	sb.border_width_top = 1
	sb.border_width_bottom = 1
	sb.border_width_left = 0
	sb.border_width_right = 0
	return sb
