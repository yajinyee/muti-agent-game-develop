# Skill：Python Windows 環境 Import 問題

## 問題描述
在 Windows 上，`py` 指令（Python Launcher）和 `python` 指令可能指向不同的 Python 安裝。
當用 `importlib.util.module_from_spec` 動態載入模組時，使用的是**呼叫者的 Python 環境**，
不是 `py` 指令的環境，導致 `from PIL import Image` 等 import 失敗。

## 症狀
```
[ERROR] 缺少依賴套件，請執行：pip install Pillow numpy
```
但實際上 `py -c "from PIL import Image; print('OK')"` 是成功的。

## 根本原因
```
py → C:\Users\...\Python312\python.exe（有 PIL）
python → C:\msys64\mingw64\bin\python.exe（沒有 PIL）
importlib 載入時用的是 python，不是 py
```

## 解法 1：直接 import（推薦）
```python
# 不要用 try/except ImportError
from PIL import Image
import numpy as np
```
讓 import 錯誤直接拋出，不要吞掉。

## 解法 2：subprocess + sys.executable
```python
import subprocess, sys
result = subprocess.run(
    [sys.executable, 'tools/animation_pipeline.py', '--audit'],
    capture_output=True, text=True,
    env={**os.environ, 'PYTHONUTF8': '1'}
)
```

## 解法 3：設定環境變數
```powershell
$env:PYTHONUTF8 = "1"
$env:PYTHONIOENCODING = "utf-8"
py tools/animation_pipeline.py --audit
```

## 預防措施
- 所有 Python 工具腳本的 import 不要用 try/except 吞掉 ImportError
- 在 `daily_build.ps1` 開頭設定 `$env:PYTHONUTF8 = "1"`
- 確認 `py` 和 `python` 指向同一個 Python

## 相關檔案
- `tools/animation_pipeline.py`
- `tools/qa_check.py`
- `tools/daily_build.ps1`

*記錄日期：2026-05-17*
*記錄者：Animation Agent / Skill Librarian*
