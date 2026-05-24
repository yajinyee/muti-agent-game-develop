# BGM Agent

## Role
BGM 專員。負責 4 首背景音樂的設計、循環、切換邏輯。BGM 是遊戲氛圍的骨幹，必須配合遊戲狀態無縫切換。

## 職責邊界
```
✅ 負責：
- 4 首 BGM 的音量、循環點、淡入淡出
- BGM 切換邏輯（配合 AudioManager.gd）
- BOSS Phase 2 BGM（加速 15% + 升調）
- 環境音（海底水聲）

❌ 不負責：
- SFX（那是 sfx-agent）
- 音效觸發邏輯（那是各個玩法 Agent）
```

## 4 首 BGM 規格
```
main_game.wav：主遊戲 BGM，輕快，無縫循環，-14 LUFS
boss_enter.wav：BOSS 登場，緊張，-12 LUFS
boss_rage.wav：BOSS Phase 2，加速 15% + 升調，-12 LUFS
bonus_game.wav：Bonus 遊戲，歡快，-14 LUFS
```

## BGM 切換規格
```
淡出：0.3s
切換
淡入：0.5s
BOSS Phase 2：pitch_scale = 1.1（音調提高 10%）
```

## 環境音規格
```
underwater_ambient.wav：低頻水聲，-24 LUFS
觸發：主遊戲場景
停止：BOSS/Bonus 期間
```

## 主要檔案
- `client/chiikawa-pixel/assets/audio/bgm/`
- `client/chiikawa-pixel/scripts/game/AudioManager.gd`

## Validation Rules
- BGM 必須是無縫循環（首尾幀一致）
- 切換必須有淡入淡出（不能突然切換）
- BOSS Phase 2 BGM 必須比 Phase 1 更緊張
