## StreakPanel.gd - DAY-083
## 連擊系統 UI：顯示當前連擊數和倍率加成
extends Node2D

const PANEL_W := 120
const PANEL_H := 48

var _font: FontFile
var _bg: ColorRect
var _streak_label: Label
var _mult_label: Label
var _level_label: Label
var _current := 0
var _mult := 1.0
var _level_name := ""
var _level_color := Color.WHITE

func setup(font: FontFile) -> void:
_font = font
_build_ui()
_connect_signals()
hide()

func _build_ui() -> void:
_bg = ColorRect.new()
_bg.color = Color(0.05, 0.05, 0.1, 0.85)
_bg.size = Vector2(PANEL_W, PANEL_H)
add_child(_bg)

_streak_label = Label.new()
_streak_label.position = Vector2(8, 4)
_streak_label.add_theme_color_override("font_color", Color.WHITE)
if _font:
_streak_label.add_theme_font_override("font", _font)
_streak_label.add_theme_font_size_override("font_size", 20)
add_child(_streak_label)

_level_label = Label.new()
_level_label.position = Vector2(8, 26)
_level_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.2))
if _font:
_level_label.add_theme_font_override("font", _font)
_level_label.add_theme_font_size_override("font_size", 12)
add_child(_level_label)

_mult_label = Label.new()
_mult_label.position = Vector2(70, 26)
_mult_label.add_theme_color_override("font_color", Color(0.5, 1.0, 0.5))
if _font:
_mult_label.add_theme_font_override("font", _font)
_mult_label.add_theme_font_size_override("font_size", 12)
add_child(_mult_label)

func _connect_signals() -> void:
if GameManager.has_signal("streak_updated"):
GameManager.streak_updated.connect(_on_streak_updated)
if GameManager.has_signal("streak_reset"):
GameManager.streak_reset.connect(_on_streak_reset)

func _on_streak_updated(data: Dictionary) -> void:
var current: int = data.get("current", 0)
var mult: float = data.get("mult_bonus", 1.0)
var level_name: String = data.get("level_name", "")
var level_color_hex: String = data.get("level_color", "#FFFFFF")
var is_new_level: bool = data.get("is_new_level", false)

_current = current
_mult = mult
_level_name = level_name
_level_color = Color(level_color_hex)

_streak_label.text = "x%d" % current
_streak_label.add_theme_color_override("font_color", _level_color)
_level_label.text = level_name if current >= 3 else ""
_mult_label.text = "x%.1f" % mult if mult > 1.0 else ""

if current >= 3:
show()
else:
hide()

if is_new_level and current >= 3:
_play_level_up_anim()

func _on_streak_reset(data: Dictionary) -> void:
_current = 0
_mult = 1.0
hide()

func _play_level_up_anim() -> void:
var tween = create_tween()
tween.tween_property(self, "scale", Vector2(1.3, 1.3), 0.1)
tween.tween_property(self, "scale", Vector2(1.0, 1.0), 0.15)