# 音效記憶 — 吉伊卡哇：像素大討伐

> 記錄所有音效相關的設定、同步參數、已知問題。由 Audio Director 維護。

**最後更新**：2025-01-01  
**更新者**：Audio Director

---

## 音效設計原則

### 整體風格
- **風格**：輕快、可愛、日系遊戲音效
- **參考**：任天堂 DS 遊戲音效、吉伊卡哇動畫音效
- **禁止**：過於激烈、暴力、恐怖的音效

### 音量標準
| 類型 | 目標音量 | 格式 |
|------|---------|------|
| 背景音樂（BGM）| -14 LUFS | WAV 44100Hz 16-bit |
| 遊戲音效（SFX）| -6 LUFS | WAV 44100Hz 16-bit |
| UI 音效 | -8 LUFS | WAV 44100Hz 16-bit |

---

## 音效清單與設定

### 攻擊音效
| 檔案 | 角色 | 觸發條件 | 延遲設定 | 音量 | 狀態 |
|------|------|---------|---------|------|------|
| attack_fire.wav | 吉伊卡哇 | 攻擊動畫第 3 幀 | 0ms | -6 LUFS | ✅ |
| attack_fire_hachiware.wav | 小八 | 攻擊動畫第 3 幀 | 0ms | -6 LUFS | ✅ |
| attack_fire_usagi.wav | 烏薩奇 | 攻擊動畫第 3 幀 | 0ms | -6 LUFS | ✅ |

### 命中/擊殺音效
| 檔案 | 觸發條件 | 延遲設定 | 音量 | 狀態 |
|------|---------|---------|------|------|
| hit.wav | 目標物被擊中（Server 確認）| 0ms | -6 LUFS | ✅ |
| kill.wav | 目標物消滅 | 0ms | -5 LUFS | ✅ |

### 獎勵音效
| 檔案 | 觸發條件 | 延遲設定 | 音量 | 狀態 |
|------|---------|---------|------|------|
| coin_drop.wav | 每次獲得金幣 | 0ms | -8 LUFS | ✅ |
| reward_bag.wav | 特殊獎勵觸發 | 0ms | -5 LUFS | ✅ |
| big_win.wav | 高倍率獎勵（>= 50x）| 0ms | -3 LUFS | ✅ |

### BOSS 音效序列
| 檔案 | 觸發條件 | 延遲設定 | 音量 | 狀態 |
|------|---------|---------|------|------|
| boss_warning.wav | BOSS 出現前 3 秒 | 0ms | -4 LUFS | ✅ |
| boss_enter.wav | BOSS 正式出現 | 0ms | -3 LUFS | ✅ |

**BOSS 音效序列說明**：
1. Server 發送 `boss_spawn` 訊息
2. Client 立即播放 `boss_warning.wav`（3 秒警告）
3. 3 秒後播放 `boss_enter.wav` + 顯示 BOSS

### Bonus 音效
| 檔案 | 觸發條件 | 延遲設定 | 音量 | 狀態 |
|------|---------|---------|------|------|
| bonus_ready.wav | Bonus 觸發確認 | 0ms | -4 LUFS | ✅ |
| bonus_game.wav | 進入 Bonus 場景 | 500ms（場景切換後）| -14 LUFS | ✅ |

### 背景音樂
| 檔案 | 使用場景 | 循環 | 音量 | 狀態 |
|------|---------|------|------|------|
| main_game.wav | 主遊戲場景 | 是（無縫）| -14 LUFS | ✅ |
| bonus_game.wav | Bonus 遊戲場景 | 是（無縫）| -14 LUFS | ✅ |

### 特殊音效
| 檔案 | 觸發條件 | 說明 | 狀態 |
|------|---------|------|------|
| weed_pull.wav | 特定目標物（草類）| 拔草動作音效 | ✅ |

---

## Godot AudioStreamPlayer 設定

### Bus 架構
```
Master Bus
├── BGM Bus（音量 -6dB，Reverb 輕微）
├── SFX Bus（音量 0dB，無效果）
└── UI Bus（音量 -3dB，無效果）
```

### 關鍵設定
```gdscript
# 背景音樂設定
$BGMPlayer.bus = "BGM"
$BGMPlayer.volume_db = -14.0
$BGMPlayer.stream.loop = true

# 音效設定
$SFXPlayer.bus = "SFX"
$SFXPlayer.volume_db = -6.0

# 攻擊音效（角色專屬）
func play_attack_sound(character: String):
    match character:
        "chiikawa": $AttackPlayer.stream = attack_fire
        "hachiware": $AttackPlayer.stream = attack_fire_hachiware
        "usagi": $AttackPlayer.stream = attack_fire_usagi
    $AttackPlayer.play()
```

---

## 同步測試記錄

### 最近一次測試結果
| 音效 | 測量誤差 | 門檻 | 狀態 |
|------|---------|------|------|
| attack_fire.wav | 待測 | <50ms | ⏳ |
| hit.wav | 待測 | <50ms | ⏳ |
| boss_warning.wav | 待測 | <100ms | ⏳ |

---

## 已知音效問題

### 待修復
- 無已知問題

### 已修復
- 無記錄

---

## 版本記錄

| 版本 | 日期 | 變更 |
|------|------|------|
| 1.0 | 2025-01-01 | 初始記錄 |
