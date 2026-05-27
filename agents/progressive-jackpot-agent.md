# Progressive Jackpot Agent

## 職責
負責四層累積獎池系統（Mini/Minor/Major/Grand Jackpot）的設計、實作和維護。

## 業界參考
- **Jili Jackpot Fishing**：四層 Progressive Jackpot，RTP 97%，最高 888x
- **Hard Rock Bet**：每次下注 10 cents 進入獎池，累積到大額後觸發
- **業界標準**：Mini（最常觸發）→ Minor → Major → Grand（最稀有）

## 系統架構

### 四層獎池設計
| 層級 | 起始倍率 | 觸發機率 | 冷卻時間 | 顏色 |
|------|---------|---------|---------|------|
| Mini  | 50x   | 0.5%   | 30 秒  | 綠色 |
| Minor | 200x  | 0.1%   | 60 秒  | 藍色 |
| Major | 1000x | 0.02%  | 120 秒 | 橙色 |
| Grand | 5000x | 0.005% | 300 秒 | 金色 |

### 累積機制
每次射擊按比例累積：
- Mini Pool：+0.01 × bet_cost
- Minor Pool：+0.005 × bet_cost
- Major Pool：+0.002 × bet_cost
- Grand Pool：+0.001 × bet_cost

### 目標物系列（T171-T175）
| 目標 | 名稱 | 觸發層級 | 倍率 |
|------|------|---------|------|
| T171 | Mini Jackpot 魚 | Mini（直接觸發） | 50x |
| T172 | Minor Jackpot 魚 | Minor（直接觸發） | 200x |
| T173 | Major Jackpot 魚 | Major（直接觸發） | 1000x |
| T174 | Grand Jackpot 魚 | Grand（直接觸發） | 5000x |
| T175 | Jackpot Trigger 魚 | 隨機（Mini 60%/Minor 30%/Major 8%/Grand 2%） | 200x |

## 主要檔案
- **Server**：`server/internal/game/lucky_jackpot_pool_handler.go`
- **Client Panel**：`client/chiikawa-pixel/scripts/ui/LuckyJackpot*Panel.gd`（5 個）
- **Registry**：`client/chiikawa-pixel/scripts/ui/LuckyPanelRegistry.gd`
- **美術**：`tools/generate_targets_day313.py`

## 設計原則
1. **漸進式獎勵**：Mini 最常觸發，Grand 最稀有，保持玩家期待感
2. **全服廣播**：任何玩家觸發 Jackpot 都廣播給全服，製造社群感
3. **即時更新**：每 5 秒廣播一次獎池狀態，讓玩家看到獎池在增長
4. **冷卻保護**：防止同一層級短時間內連續觸發，保護 RTP 平衡

## 視覺設計
- Mini：綠色，4 道光芒
- Minor：藍色，6 道光芒
- Major：橙色，8 道光芒，較大
- Grand：金色，12 道光芒，三層光環，全螢幕慶典演出
- Trigger：彩虹四色，代表四層可能性

## 品質門檻
- 精靈圖密度：≥ 35%
- 觸發廣播延遲：< 100ms
- 獎池更新頻率：每 5 秒
- Grand Jackpot 演出時長：≥ 4 秒
