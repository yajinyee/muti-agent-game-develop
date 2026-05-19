## BulletPool.gd
## 子彈 Object Pool — 避免高頻建立/刪除節點造成 GC 壓力
## 參考：Godot 4 Object Pooling 最佳實踐（uhiyama-lab.com, 2025）
##
## 使用方式：
##   var bullet = BulletPool.acquire(char_id, texture)
##   # ... 設定位置、動畫 ...
##   BulletPool.release(bullet)  # 歸還到 pool（不 queue_free）
##
## 效能提升：
##   - 避免每次射擊都 Node2D.new() + Sprite2D.new() + add_child()
##   - 避免每次命中都 queue_free()（觸發 SceneTree 重建）
##   - HTML5 環境下 GC 壓力顯著降低

extends Node

const POOL_SIZE = 20  # 預建立 20 個子彈（捕魚機最多同時 10 個玩家，每人 2 顆在飛）

# 子彈 pool：{char_id: [Node2D, ...]}
var _pools: Dictionary = {}
# 正在使用中的子彈
var _active: Array = []

func _ready() -> void:
	# 預建立各角色的子彈節點
	for char_id in ["chiikawa", "hachiware", "usagi"]:
		_pools[char_id] = []
		for i in POOL_SIZE:
			var bullet = _create_bullet(char_id)
			bullet.visible = false
			add_child(bullet)
			_pools[char_id].append(bullet)

## acquire — 從 pool 取出一個子彈節點
## 如果 pool 空了，動態建立新的（不限制上限，避免遊戲卡頓）
func acquire(char_id: String, texture: Texture2D = null) -> Node2D:
	var pool = _pools.get(char_id, [])
	var bullet: Node2D

	if pool.size() > 0:
		bullet = pool.pop_back()
	else:
		# Pool 耗盡，動態建立（不常發生）
		bullet = _create_bullet(char_id)
		add_child(bullet)

	# 重置狀態
	bullet.visible = true
	bullet.modulate = Color.WHITE
	bullet.rotation = 0.0
	bullet.scale = Vector2.ONE
	bullet.z_index = 10

	# 更新 texture（如果有提供）
	var sprite = bullet.get_node_or_null("Sprite")
	if is_instance_valid(sprite) and texture != null:
		sprite.texture = texture

	_active.append(bullet)
	return bullet

## release — 歸還子彈到 pool（不 queue_free，只是隱藏）
func release(bullet: Node2D) -> void:
	if not is_instance_valid(bullet):
		return

	bullet.visible = false
	bullet.position = Vector2.ZERO
	bullet.rotation = 0.0
	bullet.scale = Vector2.ONE
	bullet.modulate = Color.WHITE

	# 停止所有 tween（避免殘留動畫）
	var tweens = bullet.get_meta("tweens", [])
	for t in tweens:
		if is_instance_valid(t):
			t.kill()
	bullet.set_meta("tweens", [])

	_active.erase(bullet)

	# 找到對應的 pool 歸還
	var char_id = bullet.get_meta("char_id", "chiikawa")
	if _pools.has(char_id):
		_pools[char_id].append(bullet)

## release_all — 緊急回收所有子彈（場景切換時使用）
func release_all() -> void:
	for bullet in _active.duplicate():
		release(bullet)

## _create_bullet — 建立一個子彈節點（帶 Sprite2D 子節點）
func _create_bullet(char_id: String) -> Node2D:
	var bullet = Node2D.new()
	bullet.set_meta("char_id", char_id)
	bullet.set_meta("tweens", [])

	# Sprite 子節點（texture 在 acquire 時設定）
	var sprite = Sprite2D.new()
	sprite.name = "Sprite"
	sprite.texture_filter = CanvasItem.TEXTURE_FILTER_NEAREST
	sprite.scale = Vector2(1.5, 1.5)
	bullet.add_child(sprite)

	return bullet

## get_stats — 取得 pool 統計（供 PerformanceMonitor 使用）
func get_stats() -> Dictionary:
	var total_pooled = 0
	for char_id in _pools:
		total_pooled += _pools[char_id].size()
	return {
		"active": _active.size(),
		"pooled": total_pooled,
		"total": total_pooled + _active.size()
	}
