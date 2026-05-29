# Buy Bonus Agent

## 職責
負責「Buy Bonus」機制的設計、實作與維護。
Buy Bonus 是 2026 年捕魚機業界最新趨勢：玩家可以直接花費籌碼購買特定 Bonus 效果，無需等待隨機觸發。

## 業界依據
- **BGaming Fishing Club 2（2026-04）**：Fishing Net（×60 stake）+ TNT Bonus（×100 stake），兩種 Bonus 可直接購買
- **Reflex Gaming Big Game Fishing Rapid Riches（2026-05）**：Rapid Riches 快速獎勵機制
- **BGaming Shark & Spark Hold & Win（2026-05）**：Pearl 倍率符號 + Cascading Wins

## 核心機制

### Fishing Net（漁網）
- 觸發：撒網捕獲全場所有目標
- 每個目標獎勵：×60.0
- 完美條件：捕獲 ≥ 5 個目標
- 完美獎勵：全服 ×38.5 加成 77 秒

### TNT Bonus（TNT 爆炸）
- 觸發：3 秒倒數後水下大爆炸
- 傷害：全場 HP -80%
- 每個目標獎勵：×100.0
- 完美條件：炸毀 ≥ 3 個目標
- 完美獎勵：全服 ×39.0 加成 78 秒

### Disturbance System（擾動系統）
- 觸發：基於玩家最近 30 秒擊破數（擾動值 1-30）
- 倍率：擾動值越高倍率越高（×5.0 → ×50.0）
- 完美條件：擾動值 ≥ 20
- 完美獎勵：全服 ×39.5 加成 79 秒

### Pearl Multiplier（珍珠倍率）
- 觸發：為場上所有目標分配珍珠倍率（×1-×100）
- 珍珠倍率權重：×1（30%）→ ×100（1%）
- 完美條件：收集 ≥ 5 個珍珠
- 完美獎勵：全服 ×40.0 加成 80 秒（里程碑）

### Rapid Riches（快速暴富）
- 觸發：5 秒快速連擊模式
- 每次擊破獎勵：×200.0
- 完美條件：5 秒內連擊 ≥ 10 次
- 完美獎勵：全服 ×41.0 加成 82 秒（新史上最高）

## 主要檔案
- `server/internal/game/lucky_fishing_net_handler.go`
- `server/internal/game/lucky_tnt_bonus_handler.go`
- `server/internal/game/lucky_disturbance_handler.go`
- `server/internal/game/lucky_pearl_multiplier_handler.go`
- `server/internal/game/lucky_rapid_riches_handler.go`
- `client/chiikawa-pixel/scripts/ui/LuckyFishingNetPanel.gd`
- `client/chiikawa-pixel/scripts/ui/LuckyTNTBonusPanel.gd`
- `client/chiikawa-pixel/scripts/ui/LuckyDisturbancePanel.gd`
- `client/chiikawa-pixel/scripts/ui/LuckyPearlMultiplierPanel.gd`
- `client/chiikawa-pixel/scripts/ui/LuckyRapidRichesPanel.gd`

## 設計原則
1. **視覺反饋要強烈**：TNT 倒數要有明顯的視覺倒數，爆炸要有震動效果
2. **擾動值要可見**：玩家要能看到自己的擾動值，激勵他們保持活躍
3. **珍珠倍率要顯示在目標物上**：高倍率珍珠（≥50x）要特別標示
4. **Rapid Riches 計時要精確**：5 秒計時要即時更新，讓玩家感受到緊迫感

## 未來擴展方向
- Buy Bonus 按鈕 UI（玩家可以直接花費籌碼購買）
- Bonus 組合機制（同時觸發多個 Bonus）
- 限時 Bonus 活動（特定時段 Bonus 效果加倍）

## 最後更新
DAY-325（2026-05-29）：初始建立，T216-T220 五個 Buy Bonus 機制
