# Skill：ComfyUI 一致性 Spritesheet 生成（2025 最新技術）

> 來源：apatero.com, medium.com（IPAdapter + ControlNet）  
> 研究日期：2026-05-17（DAY-002）  
> 記錄者：Research Agent

---

## 2025 最新方法：三件套工作流

根據最新研究，生成一致性角色 Spritesheet 的最佳方法是：

```
SDXL（品質）+ IPAdapter FaceID（身份一致性）+ ControlNet（姿勢控制）
```

### 為什麼比 SD 1.5 + LoRA 更好

| 方法 | 一致性 | 品質 | 速度 |
|------|-------|------|------|
| SD 1.5 + LoRA（現有）| 中 | 中 | 快 |
| SDXL + IPAdapter + ControlNet | 高 | 高 | 慢 |
| 訓練角色 LoRA（最佳）| 最高 | 高 | 中 |

### 現有方法的問題

目前用 SD 1.5 + Pixel Art LoRA 生成的問題：
- 每次生成的角色外觀不完全一致（頭部比例、顏色略有差異）
- 沒有姿勢控制，attack/bigwin 姿勢不穩定

### 改善方案

**短期（不需要重新訓練）：**
1. 用 ComfyUI 的 Background Removal node 替代 Python flood fill 去背
2. 用 Image Grid node 直接生成 spritesheet 格式
3. 固定 seed + 固定 prompt 確保一致性

**中期（需要訓練）：**
1. 訓練吉伊卡哇角色 LoRA（20-30 張參考圖）
2. 用 ControlNet OpenPose 控制姿勢
3. 生成 idle/attack/bigwin 各 4 幀

### ComfyUI Workflow 改善

```python
# 現有 workflow（comfyui_generate.py）改善點：
# 1. 加入 seed 固定（確保同角色不同姿勢一致）
# 2. 加入 cfg_scale 調整（7.5 → 8.0 提升細節）
# 3. 加入 steps 增加（28 → 35 提升品質）
# 4. 考慮升級到 SDXL 模型

IMPROVED_WORKFLOW = {
    "seed": FIXED_SEED_PER_CHARACTER,  # 每個角色固定 seed
    "steps": 35,                        # 增加步數
    "cfg": 8.0,                         # 提升 CFG
    "sampler_name": "dpmpp_2m",         # 更好的 sampler
}
```

### 授權注意事項

- SD 1.5：CreativeML Open RAIL-M（可商用，有限制）
- Pixel Art LoRA：需確認原始授權
- IPAdapter：Apache 2.0（可商用）
- ControlNet：Apache 2.0（可商用）

## 下一步行動

1. 在 `tools/comfyui_generate.py` 加入固定 seed 機制
2. 研究 ComfyUI Background Removal node
3. 評估是否值得訓練角色 LoRA

## 相關檔案
- `tools/comfyui_generate.py`
- `tools/batch_process_ai.py`
- `.kiro/skills/comfyui-pixel-art.md`

*Content was rephrased for compliance with licensing restrictions*
