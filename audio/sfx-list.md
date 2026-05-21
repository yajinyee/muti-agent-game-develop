# 完整音效清單

> 維護者：Audio Director  
> 最後更新：2026-05-17  
> 總音效數：14 個

---

## SFX 音效（短音效）

### 攻擊類

| 編號 | 事件 ID | 檔案 | 觸發條件 | 音量 | 同步要求 | 狀態 |
|------|---------|------|---------|------|---------|------|
| SFX-001 | attack.chiikawa | attack_fire.wav | 吉伊卡哇按下攻擊鍵 | -4 dB | attack 動畫第 2 幀 | ✅ |
| SFX-002 | attack.hachiware | attack_fire_hachiware.wav | 小八按下攻擊鍵 | -4 dB | attack 動畫第 2 幀 | ✅ |
| SFX-003 | attack.usagi | attack_fire_usagi.wav | 烏薩奇按下攻擊鍵 | -3 dB | attack 動畫第 2 幀 | ✅ |

**同步說明**：攻擊音效必須在角色 attack 動畫的第 2 幀觸發（子彈發射的視覺幀），不能在按鍵時立即觸發。

### 命中/擊殺類

| 編號 | 事件 ID | 檔案 | 觸發條件 | 音量 | 同步要求 | 狀態 |
|------|---------|------|---------|------|---------|------|
| SFX-004 | hit.normal | hit.wav | 子彈命中任何目標 | -6 dB | 命中判定幀（立即）| ✅ |
| SFX-005 | kill.normal | kill.wav | 擊殺 1x-19x 目標 | -5 dB | 擊殺判定幀（立即）| ✅ |
| SFX-006 | kill.bigwin | big_win.wav | 擊殺 20x 以上目標 | -2 dB | 擊殺判定幀（立即）| ✅ |

**注意**：kill.bigwin 和 kill.normal 互斥，同一次擊殺只觸發一個。

### BOSS 類

| 編號 | 事件 ID | 檔案 | 觸發條件 | 音量 | 同步要求 | 狀態 |
|------|---------|------|---------|------|---------|------|
| SFX-007 | boss.warning | boss_warning.wav | BOSS 出現前 5 秒 | -3 dB | 立即 | ✅ |

**注意**：boss.warning 是 SFX，boss.enter 是 BGM（見 BGM 區段）。

### Bonus 類

| 編號 | 事件 ID | 檔案 | 觸發條件 | 音量 | 同步要求 | 狀態 |
|------|---------|------|---------|------|---------|------|
| SFX-008 | bonus.ready | bonus_ready.wav | Bonus 計量條達到 100% | -3 dB | 立即 | ✅ |

### 獎勵類

| 編號 | 事件 ID | 檔案 | 觸發條件 | 音量 | 同步要求 | 狀態 |
|------|---------|------|---------|------|---------|------|
| SFX-009 | reward.bag | reward_bag.wav | 獎勵袋物件生成 | -5 dB | 掉落動畫第 1 幀 | ✅ |
| SFX-010 | coin.drop | coin_drop.wav | 硬幣獎勵動畫 | -7 dB | 掉落動畫第 1 幀 | ✅ |

### UI 類

| 編號 | 事件 ID | 檔案 | 觸發條件 | 音量 | 同步要求 | 狀態 |
|------|---------|------|---------|------|---------|------|
| SFX-011 | ui.click | weed_pull.wav | 任何 UI 按鈕點擊 | -8 dB | 立即 | ✅ |

---

## BGM 音效（背景音樂）

| 編號 | 事件 ID | 檔案 | 觸發條件 | 音量 | 淡入/淡出 | 狀態 |
|------|---------|------|---------|------|---------|------|
| BGM-001 | boss.enter | boss_enter.wav | BOSS 正式登場 | -8 dB | 0.5s / 1.0s | ✅ |
| BGM-002 | boss.phase2 | boss_enter.wav | BOSS HP < 50% | -6 dB | 0.3s / 0.5s | ✅ |
| BGM-003 | bonus.start | bonus_game.wav | Bonus 遊戲開始 | -10 dB | 0.5s / 1.0s | ✅ |

---

## 音效觸發優先級

當多個音效同時觸發時，依以下優先級處理：

```
1. kill.bigwin（最高優先，覆蓋其他音效）
2. boss.warning / boss.enter
3. bonus.ready / bonus.start
4. kill.normal
5. hit.normal
6. attack.*（最低優先）
```

---

## 同時播放限制

| Bus | 最大同時播放數 | 說明 |
|-----|-------------|------|
| SFX | 8 | 超過時，最舊的音效停止 |
| BGM | 1 | 同時只有一個 BGM |
| UI | 3 | UI 音效可疊加 |

---

## 音效品質規格

| 規格 | 要求 |
|------|------|
| 格式 | WAV（PCM 16-bit）|
| 取樣率 | 44100 Hz |
| 聲道 | Mono（SFX）/ Stereo（BGM）|
| 最大時長 | SFX <= 3 秒，BGM 無限制 |
| 正規化 | 峰值 <= -1 dBFS |
