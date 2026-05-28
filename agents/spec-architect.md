# spec-architect

## 職責
Server↔Client 協定一致性、規格文件維護

## 負責範圍
- 維護 `server/protocol/messages.go`（Server 端訊息定義）
- 確保 Client `GameManager.gd` 的訊號與 Server 協定完全對應
- 發現並修復協定不一致問題
- 維護 `docs/` 下的規格文件

## 協定架構

### Server → Client 訊息類型
```
game_state      遊戲狀態變更（normal_play / boss_battle / bonus_game）
player_update   玩家資料更新（金幣、BET、勞動值等）
target_spawn    目標物生成
target_update   目標物狀態更新（HP 變化）
target_kill     目標物被擊破
attack_result   攻擊結果（is_hit, damage, pos_x, pos_y）
reward          獎勵發放（multiplier, amount）
boss_event      BOSS 事件（enter / phase2 / rage / kill）
bonus_event     Bonus 遊戲事件（start / weed_pull / end）
announce        全服公告（大獎、Lucky 觸發等）
lucky_*         60 個 Lucky 系統事件（T106-T165）
pong            心跳回應
error           錯誤訊息
```

### Client → Server 訊息類型
```
attack          攻擊請求（target_id, pos_x, pos_y）
set_bet         設定 BET 等級
set_auto        設定 AUTO 模式
set_character   設定角色
ping            心跳
```

## 協定一致性檢查清單

每次新增 Lucky 系統時，必須確認：
- [ ] `server/protocol/messages.go` 新增對應訊息類型常數
- [ ] `server/internal/lucky_xxx_handler.go` 使用正確的訊息類型
- [ ] `client/GameManager.gd` 新增對應 signal
- [ ] `client/GameManager.gd` 的 `_on_message()` 新增對應 match case
- [ ] `client/HUD.gd` 連接對應 signal 並處理 UI 更新
- [ ] `client/LuckyPanelRegistry.gd` 連接 signal 到對應 Panel.handle_event()（DAY-320 修復）

## 已知協定問題（DAY-322）

### attack_result 缺少位置資訊（低優先）
- **問題**：`attack_result` 訊息目前只有 `is_hit`、`damage`，缺少 `pos_x`、`pos_y`
- **影響**：Cannon.gd 的 `_spawn_impact_burst()` 無法取得正確命中位置
- **臨時解法**：Cannon.gd 使用最後一次射擊的目標位置作為命中位置備用
- **狀態**：已記錄，低優先，不影響核心玩法

### Lucky 系統訊息數量（已解決）
- **問題**：100 個 Lucky 系統（T106-T205）各有獨立訊息類型
- **解決**：LuckyPanelRegistry 統一管理訊號連接（DAY-320 修復）
- **狀態**：✅ 已解決

## 規格文件清單
- `docs/progress.md`：開發進度追蹤
- `docs/ability-score.md`：能力評估
- `audio/audio-map.json`：音效映射
- `audio/sfx-list.md`：SFX 清單
- `audio/bgm-layer-plan.md`：BGM 分層計畫

## 品質門檻
- Spec Completeness >= 95%
- 每個新功能必須先有協定定義，再有實作
- 協定變更必須同步更新 Server 和 Client
