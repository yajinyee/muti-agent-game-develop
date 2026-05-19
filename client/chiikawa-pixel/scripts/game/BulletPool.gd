## BulletPool.gd
## 子彈 Object Pool — 避免高頻建立/刪除節點造成 GC 壓力
## 設計原則：子彈節點永遠在遊戲場景中，Pool 只管理「可用清單」
##
## 正確架構（避免 reparent 的 scale/position 問題）：
##   - 所有子彈在遊戲開始時一次性加入遊戲場景
##   - 不使用時 visible=false，使用時 visible=true
##   - Pool 只是一個「可用節點清單」，不持有節點的父子關係
##
## 使用方式（由 Cannon.gd 呼叫）：
##   BulletPool.init_pool(parent_node)   # 遊戲場景 _ready 時初始化
##   var bullet = BulletPool.acquire(char_id, texture)
##   # ... 設定位置、動畫 ...
##   BulletPool.release(bullet)

extends Node

const POOL_SIZE_PER_CHAR = 8  # 每個角色預建立 8 個（3角色 × 8 = 24 個子彈）

# 子彈 pool：{char_id: [Node2D, ...]}（可用清單）
var _available: Dictionary = {}
# 正在使用中的子彈
var _active: Array = []
# 子彈的父節點（遊戲場景）
var _parent: Node = null
# 是否已初始化
var _initialized: bool = false

## init_pool — 在遊戲場景中預建立所有子彈節點
## 必須在遊戲場景的 _ready() 中呼叫，傳入遊戲場景節點作為父節點
func init_pool(parent: Node) -> void:
	if _initialized and is_instance_valid(_parent) and _parent == parent:
		return  # 已初始化，跳過

	# 清理舊的 pool（場景重載時）
	_cleanup()

	_parent = parent
	_initialized = true

	for char_id in ["chiikawa", "hachiware", "usagi"]:
		_available[char_id] = []
		for i in POOL_SIZE_PER_CHAR:
			var bullet = _create_bullet(char_id)
			bullet.visible = false
			parent.add_child(bullet)
			_available[char_id].append(bullet)

	print("[BulletPool] Initialized: %d bullets in scene" % (POOL_SIZE_PER_CHAR * 3))

## acquire — 從 pool 取出一個子彈節點
## 如果 pool 空了，動態建立新的並加入場景
func acquire(char_id: String, texture: Texture2D = null) -> Node2D:
	if not _initialized or not is_instance_valid(_parent):
		# Pool 未初始化，回傳臨時節點（降級模式）
		return _create_fallback_bullet(char_id, texture)

	var pool = _available.get(char_id, [])
	var bullet: Node2D

	if pool.size() > 0:
		bullet = pool.pop_back()
	else:
		# Pool 耗盡，動態建立並加入場景（不常發生）
		bullet = _create_bullet(char_id)
		_parent.add_child(bullet)
		print("[BulletPool] Pool exhausted for %s, created new bullet" % char_id)

	# 重置狀態
	bullet.visible = true
	bullet.modulate = Color.WHITE
	bullet.rotation = 0.0
	bullet.scale = Vector2.ONE
	bullet.z_index = 10

	# 更新 texture
	var sprite = bullet.get_node_or_null("Sprite")
	if is_instance_valid(sprite) and texture != null:
		sprite.texture = texture
		sprite.visible = true
	elif is_instance_valid(sprite):
		sprite.visible = (texture != null)

	_active.append(bullet)
	return bullet

## release — 歸還子彈到 pool（不 queue_free，只是隱藏）
func release(bullet: Node2D) -> void:
	if not is_instance_valid(bullet):
		return

	# 停止所有 tween（避免殘留動畫）
	var tweens = bullet.get_meta("active_tweens", []) as Array
	for t in tweens:
		if t != null and is_instance_valid(t):
			t.kill()
	bullet.set_meta("active_tweens", [])

	# 重置視覺狀態
	bullet.visible = false
	bullet.position = Vector2(-9999, -9999)  # 移到畫面外（確保不可見）
	bullet.rotation = 0.0
	bullet.scale = Vector2.ONE
	bullet.modulate = Color.WHITE

	_active.erase(bullet)

	# 歸還到對應的 pool
	var char_id = bullet.get_meta("char_id", "chiikawa") as String
	if _available.has(char_id):
		_available[char_id].append(bullet)

## release_all — 緊急回收所有子彈（場景切換時使用）
func release_all() -> void:
	for bullet in _active.duplicate():
		release(bullet)

## register_tween — 讓子彈追蹤自己的 tween（release 時自動 kill）
func register_tween(bullet: Node2D, tween: Tween) -> void:
	if not is_instance_valid(bullet):
		return
	var tweens = bullet.get_meta("active_tweens", []) as Array
	tweens.append(tween)
	bullet.set_meta("active_tweens", tweens)

## get_stats — 取得 pool 統計（供 PerformanceMonitor 使用）
func get_stats() -> Dictionary:
	var total_available = 0
	for char_id in _available:
		total_available += _available[char_id].size()
	return {
		"active": _active.size(),
		"pooled": total_available,
		"total": total_available + _active.size(),
		"initialized": _initialized
	}

## _create_bullet — 建立一個子彈節點（帶 Sprite2D 子節點）
func _create_bullet(char_id: String) -> Node2D:
	var bullet = Node2D.new()
	bullet.set_meta("char_id", char_id)
	bullet.set_meta("active_tweens", [])
	bullet.set_meta("pooled", true)  # 標記為 pool 管理的節點

	var sprite = Sprite2D.new()
	sprite.name = "Sprite"
	sprite.texture_filter = CanvasItem.TEXTURE_FILTER_NEAREST
	sprite.scale = Vector2(1.5, 1.5)
	bullet.add_child(sprite)

	return bullet

## _create_fallback_bullet — Pool 未初始化時的降級節點（會被 queue_free）
func _create_fallback_bullet(char_id: String, texture: Texture2D) -> Node2D:
	var bullet = Node2D.new()
	bullet.set_meta("char_id", char_id)
	bullet.set_meta("active_tweens", [])
	bullet.set_meta("pooled", false)  # 標記為非 pool 節點，需要 queue_free

	if texture != null:
		var sprite = Sprite2D.new()
		sprite.name = "Sprite"
		sprite.texture = texture
		sprite.texture_filter = CanvasItem.TEXTURE_FILTER_NEAREST
		sprite.scale = Vector2(1.5, 1.5)
		bullet.add_child(sprite)

	return bullet

## _cleanup — 清理舊的 pool（場景重載時）
func _cleanup() -> void:
	_active.clear()
	_available.clear()
	_initialized = false
	_parent = null
