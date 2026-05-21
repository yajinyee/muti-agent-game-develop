# Audio Pipeline 完整規格

> 版本：1.0.0  
> 維護者：Audio Director  
> 最後更新：2026-05-17

---

## 概覽

本規格定義吉伊卡哇：像素大討伐的音效生產與整合流程，確保 Audio Sync >= 90 的品質門檻。

---

## 事件驅動音效架構

```
遊戲事件（GDScript）
    │
    ▼
AudioEventBus（全域信號匯流排）
    │
    ├─► SFXPlayer（音效播放器）
    │       └─ AudioStreamPlayer（多個，支援同時播放）
    │
    └─► BGMController（BGM 控制器）
            └─ AudioStreamPlayer（單一，支援淡入淡出）
```

### 核心原則
1. **事件驅動**：所有音效由遊戲事件觸發，不使用輪詢
2. **幀同步**：攻擊音效必須在對應動畫幀觸發（sync_event 欄位）
3. **Bus 分離**：SFX 和 BGM 使用不同 Audio Bus，可獨立調整音量
4. **非阻塞**：音效播放不阻塞遊戲邏輯

---

## 完整事件音效表

| 事件 ID | 觸發條件 | 音效檔案 | Bus | 音量(dB) | 同步幀 |
|---------|---------|---------|-----|---------|-------|
| `attack.chiikawa` | 吉伊卡哇發射子彈 | attack_fire.wav | SFX | -4 | attack_frame_2 |
| `attack.hachiware` | 小八發射子彈 | attack_fire_hachiware.wav | SFX | -4 | attack_frame_2 |
| `attack.usagi` | 烏薩奇發射子彈 | attack_fire_usagi.wav | SFX | -3 | attack_frame_2 |
| `hit.normal` | 子彈命中普通目標 | hit.wav | SFX | -6 | 命中判定幀 |
| `kill.normal` | 擊殺普通目標 | kill.wav | SFX | -5 | 擊殺判定幀 |
| `kill.bigwin` | 擊殺 20x 以上目標 | big_win.wav | SFX | -2 | 擊殺判定幀 |
| `boss.warning` | BOSS 即將出現（倒數 5 秒）| boss_warning.wav | SFX | -3 | 立即 |
| `boss.enter` | BOSS 正式登場 | boss_enter.wav | BGM | -8 | 立即 |
| `boss.phase2` | BOSS 進入 Phase 2（HP < 50%）| boss_enter.wav | BGM | -6 | 立即（淡入）|
| `bonus.ready` | Bonus 計量條滿 | bonus_ready.wav | SFX | -3 | 立即 |
| `bonus.start` | Bonus 遊戲開始 | bonus_game.wav | BGM | -10 | 立即（淡入）|
| `reward.bag` | 獎勵袋掉落 | reward_bag.wav | SFX | -5 | 掉落動畫第 1 幀 |
| `ui.click` | UI 按鈕點擊 | weed_pull.wav | UI | -8 | 立即 |
| `coin.drop` | 硬幣掉落動畫 | coin_drop.wav | SFX | -7 | 掉落動畫第 1 幀 |

---

## BGM Layer 設計

### 遊戲狀態與 BGM 對應

| 遊戲狀態 | BGM Layer | 檔案 | 音量(dB) | 淡入時間 | 淡出時間 |
|---------|-----------|------|---------|---------|---------|
| 主遊戲（一般）| `normal` | main_game.wav | -12 | 1.0s | 0.5s |
| 主遊戲（Fever）| `fever` | main_game.wav | -10 | 0.3s | 0.3s |
| BOSS 出現 | `boss` | boss_enter.wav | -8 | 0.5s | 1.0s |
| BOSS Phase 2 | `boss_phase2` | boss_enter.wav | -6 | 0.3s | 0.5s |
| Bonus 遊戲 | `bonus` | bonus_game.wav | -10 | 0.5s | 1.0s |
| 大獎演出 | `bigwin_sting` | big_win.wav | -4 | 0.0s | 2.0s |
| 結算畫面 | `result_jingle` | big_win.wav | -8 | 0.5s | 1.0s |

### BGM 切換邏輯

```gdscript
# BGM 狀態機
enum BGMState {
    NORMAL,
    FEVER,
    BOSS,
    BOSS_PHASE2,
    BONUS,
    BIGWIN_STING,
    RESULT_JINGLE
}

# 切換規則（優先級由高到低）
# BIGWIN_STING > BOSS_PHASE2 > BOSS > BONUS > FEVER > NORMAL
```

---

## audio-map.json 格式規範

```json
{
  "version": "1.0.0",
  "events": {
    "<event_id>": {
      "file": "res://assets/audio/sfx/<filename>.wav",
      "bus": "SFX | BGM | UI",
      "volume_db": <number>,
      "sync_event": "<godot_signal_or_animation_frame>",
      "description": "<中文說明>"
    }
  },
  "bgm_layers": {
    "<layer_name>": {
      "file": "res://assets/audio/bgm/<filename>.wav",
      "volume_db": <number>,
      "fade_in": <seconds>,
      "fade_out": <seconds>
    }
  }
}
```

### 欄位說明

| 欄位 | 類型 | 必填 | 說明 |
|------|------|------|------|
| `file` | string | ✅ | Godot 資源路徑（res://開頭）|
| `bus` | string | ✅ | Audio Bus 名稱（SFX/BGM/UI）|
| `volume_db` | number | ✅ | 音量（dB），負數為衰減 |
| `sync_event` | string | ✅ | 同步觸發的 Godot 信號或動畫幀 |
| `description` | string | ✅ | 中文說明 |
| `fade_in` | number | BGM 必填 | 淡入時間（秒）|
| `fade_out` | number | BGM 必填 | 淡出時間（秒）|

---

## Audio Bus 設定

| Bus 名稱 | 用途 | 預設音量 | 效果器 |
|---------|------|---------|-------|
| Master | 主輸出 | 0 dB | Limiter |
| BGM | 背景音樂 | -3 dB | Compressor |
| SFX | 音效 | 0 dB | - |
| UI | UI 音效 | -3 dB | - |

---

## 品質門檻

| 指標 | 門檻 | 說明 |
|------|------|------|
| Audio Sync | >= 90 | 音效與動畫幀的同步精度 |
| 音效覆蓋率 | 100% | 所有定義事件都有對應音效 |
| 音量一致性 | ±3 dB | 同類音效音量差異不超過 3 dB |
