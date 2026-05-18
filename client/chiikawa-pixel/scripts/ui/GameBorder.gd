## GameBorder.gd
## 像素風格遊戲畫面邊框裝飾
## 海底主題：珊瑚、貝殼、泡泡裝飾邊框
## 作為 CanvasLayer（layer=2）疊加在遊戲畫面上，不影響遊戲邏輯

extends Node2D

const SCREEN_W = 1280.0
const SCREEN_H = 720.0
const BORDER_W = 12.0  # 邊框寬度（像素）

# 邊框顏色（深海主題）
const C_BORDER_DARK  = Color(0.05, 0.08, 0.20, 0.95)  # 深藍邊框
const C_BORDER_MID   = Color(0.08, 0.15, 0.35, 0.90)  # 中藍
const C_BORDER_LIGHT = Color(0.15, 0.30, 0.55, 0.85)  # 亮藍（高光）
const C_CORAL_1      = Color(0.90, 0.40, 0.30, 0.90)  # 珊瑚紅
const C_CORAL_2      = Color(0.95, 0.60, 0.20, 0.85)  # 珊瑚橙
const C_SHELL        = Color(0.95, 0.90, 0.75, 0.90)  # 貝殼米白
const C_SHELL_DARK   = Color(0.75, 0.65, 0.50, 0.85)  # 貝殼陰影
const C_PEARL        = Color(1.00, 0.98, 0.95, 0.95)  # 珍珠白
const C_SEAWEED      = Color(0.20, 0.65, 0.30, 0.85)  # 海草綠
const C_GOLD_TRIM    = Color(0.90, 0.75, 0.20, 0.90)  # 金色裝飾線

# 角落裝飾位置（四個角）
const CORNERS = [
	Vector2(0, 0),
	Vector2(SCREEN_W, 0),
	Vector2(0, SCREEN_H),
	Vector2(SCREEN_W, SCREEN_H),
]

var _time: float = 0.0

func _ready() -> void:
	z_index = 2  # 在遊戲元素上方，在 HUD 下方

func _process(delta: float) -> void:
	_time += delta
	queue_redraw()

func _draw() -> void:
	_draw_border_frame()
	_draw_corner_decorations()
	_draw_top_decoration()
	_draw_bottom_decoration()

## 主邊框（四邊深色框）
func _draw_border_frame() -> void:
	# 外層深色邊框
	draw_rect(Rect2(0, 0, SCREEN_W, BORDER_W), C_BORDER_DARK)
	draw_rect(Rect2(0, SCREEN_H - BORDER_W, SCREEN_W, BORDER_W), C_BORDER_DARK)
	draw_rect(Rect2(0, 0, BORDER_W, SCREEN_H), C_BORDER_DARK)
	draw_rect(Rect2(SCREEN_W - BORDER_W, 0, BORDER_W, SCREEN_H), C_BORDER_DARK)

	# 內層亮色邊框（高光效果）
	var inner = 2.0
	draw_line(Vector2(BORDER_W, BORDER_W), Vector2(SCREEN_W - BORDER_W, BORDER_W), C_BORDER_LIGHT, 1.5)
	draw_line(Vector2(BORDER_W, BORDER_W), Vector2(BORDER_W, SCREEN_H - BORDER_W), C_BORDER_LIGHT, 1.5)

	# 金色裝飾線（邊框內側）
	var gw = BORDER_W + 3
	draw_line(Vector2(gw, gw), Vector2(SCREEN_W - gw, gw), C_GOLD_TRIM, 1.0)
	draw_line(Vector2(gw, SCREEN_H - gw), Vector2(SCREEN_W - gw, SCREEN_H - gw), C_GOLD_TRIM, 1.0)
	draw_line(Vector2(gw, gw), Vector2(gw, SCREEN_H - gw), C_GOLD_TRIM, 1.0)
	draw_line(Vector2(SCREEN_W - gw, gw), Vector2(SCREEN_W - gw, SCREEN_H - gw), C_GOLD_TRIM, 1.0)

## 角落裝飾（珊瑚 + 貝殼）
func _draw_corner_decorations() -> void:
	# 左上角：珊瑚叢
	_draw_coral(Vector2(18, 18), 1)
	# 右上角：貝殼
	_draw_shell(Vector2(SCREEN_W - 18, 18), -1)
	# 左下角：貝殼
	_draw_shell(Vector2(18, SCREEN_H - 18), 1)
	# 右下角：珊瑚叢
	_draw_coral(Vector2(SCREEN_W - 18, SCREEN_H - 18), -1)

## 珊瑚裝飾（分支狀）
func _draw_coral(pos: Vector2, dir: float) -> void:
	var pulse = sin(_time * 1.2) * 0.5 + 0.5
	var c1 = C_CORAL_1.lerp(C_CORAL_2, pulse * 0.3)

	# 主幹
	draw_line(pos, pos + Vector2(0, -20), c1, 3.0)
	# 左分支
	draw_line(pos + Vector2(0, -8), pos + Vector2(dir * -8, -18), c1, 2.0)
	draw_circle(pos + Vector2(dir * -8, -18), 3.0, c1)
	# 右分支
	draw_line(pos + Vector2(0, -12), pos + Vector2(dir * 8, -22), c1, 2.0)
	draw_circle(pos + Vector2(dir * 8, -22), 3.0, c1)
	# 頂部
	draw_circle(pos + Vector2(0, -20), 4.0, c1)
	# 高光
	draw_circle(pos + Vector2(dir * -1, -21), 1.5, Color(1, 0.8, 0.7, 0.8))

## 貝殼裝飾
func _draw_shell(pos: Vector2, dir: float) -> void:
	# 貝殼主體（扇形）
	draw_arc(pos, 10.0, -PI * 0.8, PI * 0.8, 16, C_SHELL, 2.5)
	# 貝殼條紋
	for i in range(4):
		var angle = -PI * 0.6 + i * PI * 0.4
		var end_pos = pos + Vector2(cos(angle) * 9, sin(angle) * 9)
		draw_line(pos, end_pos, C_SHELL_DARK, 1.0)
	# 珍珠
	draw_circle(pos + Vector2(dir * 12, -2), 4.0, C_PEARL)
	draw_circle(pos + Vector2(dir * 12 - 1, -3), 1.5, Color(1, 1, 1, 0.9))

## 頂部裝飾（標題區域裝飾線）
func _draw_top_decoration() -> void:
	# 頂部裝飾點（等間距）
	var dot_count = 20
	for i in range(dot_count):
		var x = BORDER_W + 8 + i * ((SCREEN_W - BORDER_W * 2 - 16) / (dot_count - 1))
		var pulse = sin(_time * 1.5 + i * 0.4) * 0.5 + 0.5
		var r = 2.0 + pulse * 1.0
		var alpha = 0.5 + pulse * 0.4
		draw_circle(Vector2(x, BORDER_W * 0.5), r, Color(C_GOLD_TRIM.r, C_GOLD_TRIM.g, C_GOLD_TRIM.b, alpha))

## 底部裝飾（海草 + 氣泡）
func _draw_bottom_decoration() -> void:
	# 底部小海草（左右各一叢）
	for side in [-1, 1]:
		var base_x = SCREEN_W * 0.5 + side * 580
		var base_y = SCREEN_H - BORDER_W * 0.5
		for j in range(3):
			var ox = j * side * 4
			var sway = sin(_time * 0.8 + j * 0.5) * 3.0
			draw_line(
				Vector2(base_x + ox, base_y),
				Vector2(base_x + ox + sway, base_y - 8),
				C_SEAWEED, 1.5
			)
			draw_circle(Vector2(base_x + ox + sway, base_y - 8), 2.0, C_SEAWEED)

	# 底部裝飾點
	var dot_count = 20
	for i in range(dot_count):
		var x = BORDER_W + 8 + i * ((SCREEN_W - BORDER_W * 2 - 16) / (dot_count - 1))
		var pulse = sin(_time * 1.5 + i * 0.4 + PI) * 0.5 + 0.5
		var r = 2.0 + pulse * 1.0
		var alpha = 0.5 + pulse * 0.4
		draw_circle(Vector2(x, SCREEN_H - BORDER_W * 0.5), r, Color(C_GOLD_TRIM.r, C_GOLD_TRIM.g, C_GOLD_TRIM.b, alpha))
