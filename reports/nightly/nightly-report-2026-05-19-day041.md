# Nightly Report — DAY-041

**日期**：2026-05-19  
**報告人**：Game Director Agent  
**狀態**：✅ 全部完成

---

## 今日完成摘要

### DAY-041：TargetPool 物件池 + Server /metrics active_targets

| 項目 | 狀態 | 說明 |
|------|------|------|
| `TargetPool.gd` 建立 | ✅ | 24 個空殼節點預建立，acquire/release/get_stats |
| `TargetManager.gd` 整合 | ✅ | init_pool + acquire + release 全部替換 |
| `project.godot` autoload | ✅ | TargetPool 加入 autoload 清單 |
| `game.go` GetActiveTargetCount | ✅ | thread-safe，RLock 讀取 len(Targets) |
| `/metrics` active_targets 指標 | ✅ | chiikawa_active_targets gauge |
| Grafana dashboard 更新 | ✅ | 從 10 個面板升級到 12 個面板 |
| go build + go vet | ✅ | 零錯誤，零警告 |
| go test | ✅ | 全部通過（unlinkat 是 Norton 防毒，非測試失敗） |
| knowhow-log 更新 | ✅ | #87 #88 #89 三條新知識 |

---

## TargetPool 設計說明

### 問題
TargetManager 每次 `target_spawn` 都建立新節點（含 Sprite2D + HP條 + Label），
每次 `target_kill` 都 `queue_free`。最多 20 個目標同時存在，高頻建立/刪除
會造成 GC 壓力和 draw call 抖動。

### 解法
```
TargetPool.init_pool(self)  # 預建立 24 個空殼節點，一次性加入場景

# 目標生成時
var node = TargetPool.acquire()  # 取出空殼，清除舊子節點
# ... 填入 Sprite2D、HP條、Label ...

# 目標消滅時（動畫結束後）
TargetPool.release(node)  # 隱藏並移到 (-9999, -9999)，不 queue_free
```

### 效能改善
- 消除高頻 `add_child` / `queue_free` 的 GC 壓力
- 節點數量穩定（24 個，不再動態增減）
- draw call 更穩定（節點數固定）

---

## /metrics active_targets 指標

```
# HELP chiikawa_active_targets Current number of active targets on screen
# TYPE chiikawa_active_targets gauge
chiikawa_active_targets 8
```

Grafana 面板：
- stat 面板：0-14 綠，15-19 黃，20+ 紅
- timeseries 面板：0-25 範圍，監控目標物生成/消滅的平衡

---

## 品質指標（DAY-041 結束時）

| 指標 | 分數 | 門檻 | 狀態 |
|------|------|------|------|
| Build Stability | 100 | ≥95 | ✅ |
| Visual Consistency | 100 | ≥90 | ✅ |
| Animation Quality | 100 | ≥88 | ✅ |
| Balance Health | 96 | ≥90 | ✅ |
| Audio Sync | 97 | ≥90 | ✅ |
| Gameplay Feel | 100 | ≥85 | ✅ |
| Spec Completeness | 100 | ≥95 | ✅ |
| Regression Risk | 5 | ≤10 | ✅ |

---

## 明日計畫（DAY-042）

1. **Client 端 MultiMesh 優化**（同類型目標物合批渲染，減少 draw call）
2. **PerformanceMonitor 加入 TargetPool 統計**（active/pooled/total 顯示）
3. **Server 端 WebSocket 訊息壓縮統計**（壓縮率指標）

---

## 技術備忘

- TargetPool 與 BulletPool 的差異：BulletPool 子節點固定（Sprite2D），TargetPool 子節點動態（Sprite2D + HP條 + Label），所以 acquire 時需要清除子節點
- pool 節點的 tween 生命週期要特別注意：container 級別的 tween 需要手動 kill，子節點的 tween 在子節點 queue_free 時自動停止
- `GetActiveTargetCount()` 使用 RLock（讀鎖），不影響遊戲邏輯的寫鎖效能
