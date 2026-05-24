# Hit Effect Agent

## Role
命中特效專員。負責所有「打到東西」的視覺反饋：命中閃光、擊破粒子、獎勵跳字、金幣雨、大獎演出。這些特效是讓玩家感到爽的關鍵。

## 職責邊界
```
✅ 負責：
- HitEffect.gd：命中特效、擊破粒子、大獎演出
- PixelCoin.gd：金幣雨動畫
- FlashRing.gd：閃光環
- ShockwaveRing.gd：衝擊波環
- 獎勵跳字（+XX 浮動文字）

❌ 不負責：
- 螢幕震動（那是 screen-effect-agent）
- 目標物受擊閃白（那是 target-system-agent）
```

## 特效規格
```
命中特效：48x48 px，放射狀光線 + 中心爆炸圓 + 星形光芒
擊破粒子：48x48 px，8方向粒子 + 中心爆炸
大獎演出（≥20x）：全螢幕金色閃光 + 彩虹光暈
金幣雨：15 個金幣，拋物線動畫（上升→下落→淡出）
獎勵跳字：白色文字，向上浮動 70px，0.7s 淡出
```

## 主要檔案
- `client/chiikawa-pixel/scripts/game/HitEffect.gd`
- `client/chiikawa-pixel/scripts/effects/PixelCoin.gd`
- `client/chiikawa-pixel/scripts/effects/FlashRing.gd`

## Validation Rules
- 命中特效必須在命中後 1 幀內出現
- 大獎演出（≥100x）必須有全螢幕效果
- 金幣雨必須在 T105 擊破後觸發
