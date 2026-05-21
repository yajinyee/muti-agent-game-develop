# 吉伊卡哇：像素大討伐 — 遊戲規格知識庫

> 來源：完整遊戲規格書 v1.0（48頁）
> 用途：開發參考、AI 知識庫、快速查閱

---

## 1. 遊戲定位

- **類型**：捕魚機 / 休閒射擊 / IP 包裝型
- **美術**：16-bit Retro Pixel Art
- **架構**：Go Server + Godot Client（WebSocket 連線）
- **核心體驗**：可愛角色討伐怪物，獲得勞動報酬袋

### 包裝轉換對照
| 傳統捕魚機 | 本作 |
|-----------|------|
| 砲台 | 吉伊卡哇 / 小八 / 烏薩奇 |
| 砲彈 | 討伐棒劍氣 |
| 魚群 | 怪物、雜草、小蟲 |
| 金幣掉落 | 勞動報酬袋、銅幣、金幣 |
| BOSS | 那個孩子 / 巨型奇美拉 Anoko |
| 能量條 | 勞動值 |
| Bonus Game | 瘋狂拔草 Weeding Frenzy |

---

## 2. 角色系統

### 吉伊卡哇（LV1-3）
- 攻擊色：粉紅色劍氣
- 攻速：普通（2.0-2.1 shots/sec）
- 定位：低投注、新手友善
- Kill Modifier：1.00 / Labor Modifier：1.10
- 大獎演出：驚慌跳起，「YaDa」字卡

### 小八（LV4-7）
- 攻擊色：藍色劍氣
- 攻速：稍快（2.2-2.5 shots/sec）
- 定位：中投注、平衡型
- Kill Modifier：1.00 / Fire Rate Modifier：1.08
- 大獎演出：高舉討伐棒，「尖尖哇嘎乃」字卡

### 烏薩奇（LV8-10）
- 攻擊色：黃色旋轉殘影
- 攻速：快（2.7-3.0 shots/sec）
- 定位：高投注、高爆發
- Kill Modifier：0.98 / Fire Rate Modifier：1.20
- 大獎演出：高速旋轉跳起，「Yaha」字卡

---

## 3. 投注等級表

| LV | 角色 | Bet Cost | Attack Power | Fire Rate | 投射物速度 |
|----|------|----------|--------------|-----------|-----------|
| 1 | 吉伊卡哇 | 1 | 1 | 2.0 | 700 |
| 2 | 吉伊卡哇 | 2 | 2 | 2.0 | 720 |
| 3 | 吉伊卡哇 | 3 | 3 | 2.1 | 740 |
| 4 | 小八 | 5 | 5 | 2.2 | 780 |
| 5 | 小八 | 10 | 10 | 2.3 | 800 |
| 6 | 小八 | 20 | 20 | 2.4 | 820 |
| 7 | 小八 | 30 | 30 | 2.5 | 850 |
| 8 | 烏薩奇 | 50 | 50 | 2.7 | 900 |
| 9 | 烏薩奇 | 80 | 80 | 2.9 | 940 |
| 10 | 烏薩奇 | 100 | 100 | 3.0 | 980 |

---

## 4. 目標物完整 Paytable

### 基礎目標（2x-10x）
| ID | 名稱 | 倍率 | HP | 出現權重 | 移動速度 | 停留時間 | 勞動值 |
|----|------|------|----|---------|---------|---------|-------|
| T001 | 像素雜草 | 2x | 3 | 180 | 0 | 20s | 1 |
| T002 | 綠色小蟲 | 3x | 5 | 160 | 40 | 18s | 1 |
| T003 | 紅色小蟲 | 5x | 8 | 130 | 55 | 16s | 1 |
| T004 | 藍色小蟲 | 6x | 10 | 110 | 65 | 15s | 2 |
| T005 | 會走路的布丁 | 8x | 16 | 90 | 35 | 20s | 2 |
| T006 | 巨大蘑菇 | 10x | 22 | 70 | 25 | 22s | 3 |

### 特殊目標（15x-50x）
| ID | 名稱 | 倍率 | HP | 出現權重 | 移動速度 | 停留時間 | 勞動值 | 特殊行為 |
|----|------|------|----|---------|---------|---------|-------|---------|
| T101 | 擬態型怪物 | 15x-30x | 35 | 35 | 50 | 14s | 5 | 死亡變回原形 |
| T102 | 寶箱怪 | 25x | 55 | 22 | 70 | 10s | 6 | 受擊後加速逃跑 |
| T103 | 流星 Star | 20x-50x | 20 | 18 | 220 | 4s | 5 | 快速通過 |
| T104 | 金色雜草 | 30x | 45 | 12 | 0 | 8s | 15 | 大量勞動值 |
| T105 | 巨大金幣魚 | 50x | 90 | 8 | 80 | 8s | 10 | 擊破後金幣雨 |

### 流星倍率權重
| 倍率 | 權重 |
|------|------|
| 20x | 50 |
| 30x | 30 |
| 40x | 15 |
| 50x | 5 |

### BOSS（100x-500x）
| ID | 名稱 | HP | 出現權重 | 停留時間 | 勞動值 |
|----|------|-----|---------|---------|-------|
| B001 | 那個孩子 | 3000 | 事件觸發 | 60s | 30 |

---

## 5. 擊破判定（混合制）

### 公式
```
Kill Chance = Base_RTP_Factor × Bet_Cost ÷ Target_Reward_Value
           ≈ 0.92 ÷ Multiplier（簡化版）
```

### 保底機制
```
Required_Hits = ceil(Target_Multiplier ÷ Bet_Cost × Difficulty_Factor)
```
| 目標類型 | Difficulty Factor |
|---------|-----------------|
| 基礎目標 | 0.3 - 0.5 |
| 特殊目標 | 0.6 - 0.9 |
| BOSS | 1.2 - 2.0 |

### 單次命中擊破率參考
| 倍率 | 擊破率 |
|------|-------|
| 2x | 46.0% |
| 5x | 18.4% |
| 10x | 9.2% |
| 25x | 3.7% |
| 50x | 1.8% |
| 100x | 0.9% |
| 500x | 0.18% |

---

## 6. 勞動值系統

- 上限：100
- 觸發 Bonus 後歸零

| 目標 | 勞動值 |
|------|-------|
| 像素雜草 | +1 |
| 小蟲類 | +1 |
| 藍色小蟲 | +2 |
| 布丁 | +2 |
| 巨大蘑菇 | +3 |
| 擬態型怪物 | +5 |
| 寶箱怪 | +6 |
| 流星 | +5 |
| 金色雜草 | +15 |
| BOSS | +30 |

---

## 7. Bonus Game：瘋狂拔草

### 流程
勞動值滿 → Bonus Ready → 草地場景 → 15秒倒數 → 點擊拔草 → 結算

### Bonus 目標表
| ID | 名稱 | 點擊分數 | 出現權重 | 特殊效果 |
|----|------|---------|---------|---------|
| BG001 | 普通雜草 | 1 | 180 | 無 |
| BG002 | 硬雜草 | 3 | 80 | 需連點2次 |
| BG003 | 發光雜草 | 8 | 35 | 增加倍率 |
| BG004 | 金色雜草 | 20 | 10 | 觸發巨大金幣 |
| BG005 | 搗亂怪草 | -5 | 20 | 扣分或暫停0.3秒 |

### 倍率計算
```
Bonus_Score = 普通×1 + 硬×3 + 發光×8 + 金色×20 - 搗亂×5
Bonus_Multiplier = clamp(20 + Score×0.375, 20, 50)
Bonus_Reward = Entry_Bet_Cost × Bonus_Multiplier
```
> 注意：原始規格為 clamp(50+Score×2, 50, 150)，已依 RTP 校正調整為 20-50x（Prototype 版）

---

## 8. BOSS 戰規格

### 觸發
- 時間觸發：每3-5分鐘
- Bonus 後提高出現率
- Prototype：手動觸發按鈕

### BOSS 參數
| 參數 | 數值 |
|------|------|
| HP | 3000 |
| 基礎倍率 | 100x |
| 最高倍率 | 500x |
| 出場時間 | 60秒 |
| Phase 2 門檻 | HP ≤ 50% |

### 獎勵（依擊殺剩餘時間）
| 剩餘時間 | 倍率 |
|---------|------|
| 0-10秒 | 100x |
| 11-20秒 | 150x |
| 21-30秒 | 200x |
| 31-40秒 | 300x |
| 41-50秒 | 400x |
| 51-60秒 | 500x |

---

## 9. 目標生成系統

### 生成參數
| 參數 | 值 |
|------|---|
| Spawn Interval | 0.8秒 |
| Max Targets | 18 |
| Basic Ratio | 85% |
| Special Ratio | 15% |
| BOSS期間 Max | 8 |

### 動態難度（依 Bet Level）
| Bet 區間 | 基礎目標 | 特殊目標 | 高倍率 |
|---------|---------|---------|-------|
| LV1-3 | 90% | 9% | 1% |
| LV4-7 | 82% | 15% | 3% |
| LV8-10 | 75% | 20% | 5% |

---

## 10. 遊戲狀態機

```
Loading → Lobby → NormalPlay → SpecialTargetEvent
                              → BossWarning → BossBattle → BossResult
                              → BonusReady → BonusGame → BonusResult
                              → NormalPlay（循環）
```

---

## 11. RTP 設計

- 目標 RTP：94%
- 波動定位：中高波動

| 獎勵來源 | RTP 佔比 |
|---------|---------|
| 基礎目標 | 45% |
| 特殊目標 | 25% |
| Bonus Game | 15% |
| BOSS | 15% |

---

## 12. WebSocket 訊息協定（已完整實作）

### Client → Server
| 訊息類型 | Payload | 說明 |
|---------|---------|------|
| `attack` | `{target_id, click_x, click_y}` | 玩家攻擊 |
| `lock` | `{target_id}` | 鎖定目標（空=解除）|
| `auto_toggle` | `{}` | 切換自動攻擊 |
| `bet_change` | `{bet_level}` | 切換投注等級 |
| `bonus_click` | `{target_id, click_x, click_y}` | Bonus 拔草點擊 |
| `ping` | `{}` | 心跳 |
| `trigger_boss` | `{}` | 手動觸發 BOSS（Prototype）|
| `trigger_bonus` | `{}` | 手動觸發 Bonus（Prototype）|

### Server → Client
| 訊息類型 | Payload | 說明 |
|---------|---------|------|
| `game_state` | `{state, timestamp}` | 遊戲狀態變更 |
| `target_spawn` | `{instance_id, def_id, name, type, x, y, hp, max_hp, speed, lifetime, behavior}` | 目標生成 |
| `target_update` | `{instance_id, hp, max_hp, x, y, phase, is_fleeing}` | 目標狀態更新 |
| `target_kill` | `{instance_id, def_id, multiplier, reward, labor_gain, killer_id}` | 目標擊破 |
| `attack_result` | `{target_id, is_hit, is_kill, damage, reward, labor_gain, character_id, multiplier}` | 攻擊結果（單播）|
| `reward` | `{source, amount, multiplier, new_balance}` | 獎勵發放（單播）|
| `boss_event` | `{event, instance_id, phase, hp, max_hp, reward, multiplier}` | BOSS 事件 |
| `bonus_event` | `{event, time_left, score, multiplier, reward}` | Bonus 事件 |
| `player_update` | `{id, coins, bet_level, bet_cost, character_id, character_name, labor_value, is_auto, lock_target_id, projectile_speed, fire_rate}` | 玩家狀態（單播）|
| `error` | `{code, message}` | 錯誤訊息（單播）|
| `pong` | `{}` | 心跳回應 |

---

## 13. 資料表設計

詳見規格書第36章，核心表：
- `character_table`：角色設定
- `target_table`：目標物設定
- `bet_table`：投注等級
- `bonus_table`：Bonus 設定
- `event_config_table`：活動設定
