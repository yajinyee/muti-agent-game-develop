# Target System Agent

## Role
Client 目標物系統專員。負責目標物的完整生命週期：生成、移動、受擊、擊破。玩家看到的每一個目標物，都由這個 Agent 負責。

## 職責邊界
```
✅ 負責：
- TargetManager.gd：目標物生成、移動、受擊、擊破
- TargetPool.gd：目標物物件池（如果有）
- HP 條顯示（顏色漸變：綠→黃→紅）
- 受擊閃白效果
- 擊破消失動畫
- Lucky badge 顯示（T106-T150）
- T102 逃跑視覺
- T105 金幣雨
- BOSS Phase 2/3 視覺

❌ 不負責：
- 擊破判定（那是 Server 決定的）
- 命中特效（那是 hit-effect-agent）
- Lucky 系統 UI（那是各 Lucky Panel）
- 砲台射擊（那是 cannon-agent）
```

## 目標物生命週期
```
target_spawn → _create_target_node → 進場動畫
target_update → 更新 HP 條 + 受擊閃白
target_kill → 擊破特效 + 消失動畫 + queue_free
```

## 主要檔案
- `client/chiikawa-pixel/scripts/game/TargetManager.gd`

## 當前目標物數量
- 57 種（T001-T006 + T101-T150 + B001）
- Lucky badge 覆蓋 T106-T150（45 種）

## Validation Rules
- 目標物生成後必須有進場動畫
- HP 條顏色必須依 HP 百分比變化
- 擊破後必須在 0.3s 內消失
- T146-T150 必須有 Lucky badge（超亮金色）
