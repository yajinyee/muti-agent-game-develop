## TargetPool.gd
## 目標物 Object Pool — 避免高頻建立/刪除節點造成 GC 壓力
## 設計原則：目標節點永遠在遊戲場景中，Pool 只管理「可用清單」
##
## 背景：TargetManager 每次 target_spawn 都建立新節點（含 Sprite2D + HP條 + Label），
##       每次 target_kill 都 queue_free。最多 20 個目標同時存在，高頻建立/刪除
##       會造成 GC 壓力和 draw call 抖動。
##
## 解法：預建立 POOL_SIZE 個「空殼節點」，acquire 時重置狀態並填入資料，
##       release 時隱藏並移到畫面外，下次直接重用。
##
## 使用方式（由 TargetManager.gd 呼叫）：
##   TargetPool.init_pool(parent_node)   # 遊戲場景 _ready 時初始化
##   var node = TargetPool.acquire()     # 取出空殼節點
##   # ... 填入 Sprite、HP條等子節點 ...
##   TargetPool.release(node)            # 歸還（不 queue_free）

extends Node

const POOL_SIZE = 24  # 預建立 24 個（最大目標數 20 + 緩衝 4）

# 可用節點清單
var _available: Array = []
# 正在使用中的節點
var _active: Array = []
# 父節點（遊戲場景）
var _parent: Node = null
# 是否已初始化
var _initialized: bool = false

## init_pool — 在遊戲場景中預建立所有目標節點容器
## 必須在遊戲場景的 _ready() 中呼叫
func init_pool(parent: Node) -> void:
	if _initialized and is_instance_valid(_parent) and _parent == parent:
		return  # 已初始化，跳過

	_cleanup()
	_parent = parent
	_initialized = true

	for i in POOL_SIZE:
		var node = _create_empty_container(i)
		node.visible = false
		node.position = Vector2(-9999, -9999)
		parent.add_child(node)
		_available.append(node)

	print("[TargetPool] Initialized: %d target containers in scene" % POOL_SIZE)

## acquire — 從 pool 取出一個空殼節點
## 如果 pool 空了，動態建立新的並加入場景
func acquire() -> Node2D:
	if not _initialized or not is_instance_valid(_parent):
		# Pool 未初始化，回傳臨時節點（降級模式）
		var fallback = Node2D.new()
		fallback.set_meta("pooled", false)
		return fallback

	var node: Node2D
	if _available.size() > 0:
		node = _available.pop_back()
	else:
		# Pool 耗盡，動態建立（不常發生）
		node = _create_empty_container(_active.size() + POOL_SIZE)
		_parent.add_child(node)
		push_warning("[TargetPool] Pool exhausted, created new container (total: %d)" % (_active.size() + 1))

	# 重置基本狀態
	node.visible = true
	node.modulate = Color.WHITE
	node.rotation = 0.0
	node.scale = Vector2.ONE
	node.z_index = 5

	# 清除所有子節點（上次使用留下的 Sprite、HP條、Label 等）
	for child in node.get_children():
		child.queue_free()

	# 清除所有 meta（上次使用留下的資料）
	for key in node.get_meta_list():
		node.remove_meta(key)

	# 重新標記為 pool 管理
	node.set_meta("pooled", true)

	_active.append(node)
	return node

## release — 歸還節點到 pool（不 queue_free，只是隱藏）
func release(node: Node2D) -> void:
	if not is_instance_valid(node):
		return

	# 非 pool 管理的節點直接 queue_free
	if not node.get_meta("pooled", false):
		node.queue_free()
		return

	# 停止所有 tween（避免殘留動畫影響下次使用）
	# 注意：子節點的 tween 會在 queue_free 時自動停止
	# 但 container 本身的 tween 需要手動停止
	# GDScript 4 沒有直接停止所有 tween 的 API，
	# 改用 kill_tweens meta 追蹤
	var tweens = node.get_meta("active_tweens", []) as Array
	for t in tweens:
		if t != null and is_instance_valid(t):
			t.kill()

	# 清除所有子節點（Sprite2D、HP條、Label、LockFrame 等）
	for child in node.get_children():
		child.queue_free()

	# 重置視覺狀態
	node.visible = false
	node.position = Vector2(-9999, -9999)
	node.rotation = 0.0
	node.scale = Vector2.ONE
	node.modulate = Color.WHITE
	node.z_index = 5

	# 清除所有 meta
	for key in node.get_meta_list():
		node.remove_meta(key)

	# 重新標記
	node.set_meta("pooled", true)

	_active.erase(node)
	_available.append(node)

## release_all — 緊急回收所有節點（場景切換時使用）
func release_all() -> void:
	for node in _active.duplicate():
		release(node)

## register_tween — 讓節點追蹤自己的 tween（release 時自動 kill）
func register_tween(node: Node2D, tween: Tween) -> void:
	if not is_instance_valid(node):
		return
	var tweens = node.get_meta("active_tweens", []) as Array
	tweens.append(tween)
	node.set_meta("active_tweens", tweens)

## get_stats — 取得 pool 統計（供 PerformanceMonitor 顯示）
func get_stats() -> Dictionary:
	return {
		"active": _active.size(),
		"pooled": _available.size(),
		"total": _active.size() + _available.size(),
		"initialized": _initialized
	}

## get_active_count — 取得當前活躍目標數量
func get_active_count() -> int:
	return _active.size()

## _create_empty_container — 建立一個空殼 Node2D 容器
## 不包含任何子節點，acquire 後由 TargetManager 填入
func _create_empty_container(index: int) -> Node2D:
	var node = Node2D.new()
	node.name = "TargetContainer_%d" % index
	node.set_meta("pooled", true)
	node.set_meta("active_tweens", [])
	return node

## _cleanup — 清理舊的 pool（場景重載時）
func _cleanup() -> void:
	_active.clear()
	_available.clear()
	_initialized = false
	_parent = null
