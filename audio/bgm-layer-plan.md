# BGM 分層計畫

> 維護者：Audio Director  
> 最後更新：2026-05-17

---

## 遊戲狀態機與 BGM 對應

```
遊戲狀態機：
                    ┌─────────────┐
                    │   LOADING   │ → 無 BGM
                    └──────┬──────┘
                           │ 載入完成
                    ┌──────▼──────┐
                    │   NORMAL    │ ← main_game.wav (-12dB)
                    └──────┬──────┘
                    │      │      │
              Fever │      │      │ BOSS Warning
                    │      │      │
             ┌──────▼──┐   │   ┌──▼──────┐
             │  FEVER  │   │   │  BOSS   │ ← boss_enter.wav (-8dB)
             └──────┬──┘   │   └──┬──────┘
                    │      │      │ HP < 50%
                    │      │      │
                    │      │   ┌──▼──────────┐
                    │      │   │ BOSS_PHASE2 │ ← boss_enter.wav (-6dB, pitch+10%)
                    │      │   └──┬──────────┘
                    │      │      │ BOSS 死亡
                    │      │      │
                    └──────┼──────┘
                           │
                    Bonus  │
                    ┌──────▼──────┐
                    │    BONUS    │ ← bonus_game.wav (-10dB)
                    └──────┬──────┘
                           │ Bonus 結束
                           │
                    ┌──────▼──────┐
                    │  BIGWIN     │ ← big_win.wav (-4dB, 一次性)
                    │  STING      │
                    └──────┬──────┘
                           │ 演出結束
                           │
                    ┌──────▼──────┐
                    │   RESULT    │ ← big_win.wav (-8dB, 一次性)
                    └──────┬──────┘
                           │ 確認
                           └──► NORMAL
```

---

## 各狀態 BGM 詳細設定

### NORMAL（一般遊戲）
```
檔案：main_game.wav
Bus：BGM
音量：-12 dB
淡入：1.0 秒
淡出：0.5 秒
循環：是
觸發條件：遊戲開始、BOSS 死亡後、Bonus 結束後
```

### FEVER（Fever 狀態）
```
檔案：main_game.wav（同一檔案，但參數不同）
Bus：BGM
音量：-10 dB（比 NORMAL 大 2dB）
音調：+5%（pitch_scale = 1.05）
淡入：0.3 秒（快速切換）
淡出：0.3 秒
循環：是
觸發條件：Fever 計量條達到 100%
退出條件：Fever 時間結束（30 秒）
```

### BOSS（BOSS 戰）
```
檔案：boss_enter.wav
Bus：BGM
音量：-8 dB
淡入：0.5 秒
淡出：1.0 秒
循環：是
觸發條件：BOSS 正式登場（boss.enter 事件）
退出條件：BOSS 死亡
```

### BOSS_PHASE2（BOSS 第二階段）
```
檔案：boss_enter.wav（同一檔案，但音調更高）
Bus：BGM
音量：-6 dB（比 BOSS 大 2dB）
音調：+10%（pitch_scale = 1.1）
淡入：0.3 秒
淡出：0.5 秒
循環：是
觸發條件：BOSS HP 降至 50% 以下
退出條件：BOSS 死亡
```

### BONUS（Bonus 遊戲）
```
檔案：bonus_game.wav
Bus：BGM
音量：-10 dB
淡入：0.5 秒
淡出：1.0 秒
循環：是
觸發條件：Bonus 遊戲開始（bonus.start 事件）
退出條件：Bonus 遊戲結束
```

### BIGWIN_STING（大獎演出）
```
檔案：big_win.wav
Bus：BGM（暫時覆蓋）
音量：-4 dB
淡入：0.0 秒（立即）
淡出：2.0 秒
循環：否（一次性）
觸發條件：擊殺 20x 以上目標
退出條件：音效播放完畢，自動回到 NORMAL
```

### RESULT_JINGLE（結算畫面）
```
檔案：big_win.wav
Bus：BGM
音量：-8 dB
淡入：0.5 秒
淡出：1.0 秒
循環：否（一次性）
觸發條件：進入結算畫面
退出條件：玩家確認結算
```

---

## GDScript 實作參考

```gdscript
# BGMController.gd
extends Node

enum BGMState {
    NONE,
    NORMAL,
    FEVER,
    BOSS,
    BOSS_PHASE2,
    BONUS,
    BIGWIN_STING,
    RESULT_JINGLE
}

# BGM 優先級（數字越大優先級越高）
const BGM_PRIORITY = {
    BGMState.NONE: 0,
    BGMState.NORMAL: 1,
    BGMState.FEVER: 2,
    BGMState.BONUS: 3,
    BGMState.BOSS: 4,
    BGMState.BOSS_PHASE2: 5,
    BGMState.BIGWIN_STING: 6,
    BGMState.RESULT_JINGLE: 7,
}

var current_state: BGMState = BGMState.NONE
var bgm_player: AudioStreamPlayer

func transition_to(new_state: BGMState, fade_time: float = 0.5):
    if BGM_PRIORITY[new_state] < BGM_PRIORITY[current_state]:
        return  # 不降低優先級
    
    # 淡出當前 BGM
    var tween = create_tween()
    tween.tween_property(bgm_player, "volume_db", -80, fade_time)
    await tween.finished
    
    # 切換到新 BGM
    current_state = new_state
    _play_bgm_for_state(new_state)

func _play_bgm_for_state(state: BGMState):
    match state:
        BGMState.NORMAL:
            bgm_player.stream = load("res://assets/audio/bgm/main_game.wav")
            bgm_player.volume_db = -12
            bgm_player.pitch_scale = 1.0
        BGMState.FEVER:
            bgm_player.volume_db = -10
            bgm_player.pitch_scale = 1.05
        BGMState.BOSS:
            bgm_player.stream = load("res://assets/audio/bgm/boss_enter.wav")
            bgm_player.volume_db = -8
            bgm_player.pitch_scale = 1.0
        BGMState.BOSS_PHASE2:
            bgm_player.volume_db = -6
            bgm_player.pitch_scale = 1.1
        BGMState.BONUS:
            bgm_player.stream = load("res://assets/audio/bgm/bonus_game.wav")
            bgm_player.volume_db = -10
            bgm_player.pitch_scale = 1.0
    
    bgm_player.play()
```

---

## 音樂切換測試清單

| 切換場景 | 預期行為 | 測試狀態 |
|---------|---------|---------|
| 遊戲開始 → NORMAL | main_game.wav 淡入 1 秒 | ✅ |
| NORMAL → BOSS | boss_enter.wav 淡入 0.5 秒 | ✅ |
| BOSS → BOSS_PHASE2 | 音調提高，音量增大 | ✅ |
| BOSS 死亡 → NORMAL | boss_enter.wav 淡出，main_game.wav 淡入 | ✅ |
| NORMAL → BONUS | bonus_game.wav 淡入 0.5 秒 | ✅ |
| BONUS 結束 → NORMAL | bonus_game.wav 淡出，main_game.wav 淡入 | ✅ |
| 擊殺 20x → BIGWIN_STING | big_win.wav 立即播放 | ✅ |
| BIGWIN_STING 結束 → NORMAL | 自動回到 main_game.wav | ✅ |
