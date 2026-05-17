# 失敗記錄：ComfyUI GPU 初始化失敗

**日期**：2026-05-15  
**記錄者**：Sprite Generation Agent  
**狀態**：✅ 已解決

---

## 問題描述

嘗試啟動 ComfyUI 時，GPU 初始化失敗，無法使用 CUDA 加速。

### 錯誤訊息

```
RuntimeError: CUDA error: no kernel image is available for execution on the device
CUDA kernel errors might be asynchronously reported at some other API call, so the stacktrace below might be incorrect.
For debugging consider passing CUDA_LAUNCH_BLOCKING=1.

Traceback (most recent call last):
  File "main.py", line 45, in <module>
    import torch
  ...
  File "torch/cuda/__init__.py", line 289, in _check_capability
    raise RuntimeError(...)
```

### 環境資訊

| 項目 | 版本 |
|------|------|
| PyTorch | 2.11.0+cu130 |
| CUDA（PyTorch 要求）| 13.0 |
| NVIDIA 驅動版本 | 555.85 |
| 驅動支援的最高 CUDA | 12.5 |
| GPU | NVIDIA RTX 系列 |

---

## 根本原因分析

**核心問題**：PyTorch 版本與 NVIDIA 驅動版本不相容。

```
PyTorch 2.11.0+cu130 需要 CUDA 13.0
CUDA 13.0 需要 NVIDIA 驅動 >= 570.xx
當前驅動 555.85 只支援 CUDA 12.5
```

### 相容性矩陣

| PyTorch 版本 | CUDA 版本 | 最低驅動版本 |
|-------------|---------|------------|
| 2.11.0+cu130 | 13.0 | 570.xx |
| 2.5.0+cu124 | 12.4 | 550.xx |
| 2.4.0+cu121 | 12.1 | 530.xx |
| 2.3.0+cu118 | 11.8 | 520.xx |

---

## 解決過程

### 嘗試 1：降級 PyTorch（失敗）

```bash
pip install torch==2.4.0+cu121 --index-url https://download.pytorch.org/whl/cu121
```

結果：ComfyUI 的某些節點需要 PyTorch >= 2.5，降級後節點報錯。

### 嘗試 2：更新 NVIDIA 驅動（成功）

1. 前往 [NVIDIA 驅動下載頁面](https://www.nvidia.com/Download/index.aspx)
2. 選擇對應 GPU 型號
3. 下載驅動版本 **596.49**（支援 CUDA 13.0）
4. 安裝並重啟

```
安裝後驗證：
nvidia-smi
→ Driver Version: 596.49   CUDA Version: 13.0
```

5. 重新啟動 ComfyUI → 成功！

---

## 解決方案

**更新 NVIDIA 驅動到 596.49**（或更新版本）

```powershell
# 確認當前驅動版本
nvidia-smi

# 確認 CUDA 版本
nvcc --version

# 確認 PyTorch CUDA 版本
python -c "import torch; print(torch.version.cuda)"
```

---

## 教訓

### 核心教訓：安裝前確認 PyTorch 版本與驅動相容性

**安裝 PyTorch 前必做的檢查**：

```python
# check_cuda_compatibility.py
import subprocess
import re

def get_driver_version():
    """取得 NVIDIA 驅動版本"""
    try:
        output = subprocess.check_output(['nvidia-smi'], text=True)
        match = re.search(r'Driver Version: (\d+\.\d+)', output)
        if match:
            return float(match.group(1))
    except:
        return None

def get_required_driver(cuda_version: str) -> float:
    """根據 CUDA 版本返回最低驅動要求"""
    requirements = {
        "13.0": 570.0,
        "12.4": 550.0,
        "12.1": 530.0,
        "11.8": 520.0,
    }
    return requirements.get(cuda_version, 999.0)

# 使用範例
driver = get_driver_version()
print(f"當前驅動版本：{driver}")

# 如果要安裝 PyTorch 2.11.0+cu130
required = get_required_driver("13.0")
if driver and driver < required:
    print(f"❌ 驅動版本不足！需要 >= {required}，當前 {driver}")
    print("請先更新 NVIDIA 驅動")
else:
    print("✅ 驅動版本符合要求")
```

### 標準安裝流程

```
1. 確認 GPU 型號
2. 查詢 GPU 支援的最高 CUDA 版本
3. 確認當前驅動版本（nvidia-smi）
4. 根據驅動版本選擇對應的 PyTorch 版本
5. 安裝 PyTorch
6. 驗證：python -c "import torch; print(torch.cuda.is_available())"
```

---

## 相關資源

- [PyTorch 官方安裝頁面](https://pytorch.org/get-started/locally/)（選擇正確的 CUDA 版本）
- [NVIDIA CUDA 與驅動相容性表](https://docs.nvidia.com/cuda/cuda-toolkit-release-notes/index.html)
- [NVIDIA 驅動下載](https://www.nvidia.com/Download/index.aspx)

---

*記錄時間：2026-05-15*  
*解決時間：2026-05-15（當日解決）*
