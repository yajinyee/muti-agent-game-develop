@echo off
chcp 65001 >nul
echo === Starting ComfyUI ===
echo.

set COMFYUI_DIR=C:\ComfyUI\ComfyUI_windows_portable

if not exist "%COMFYUI_DIR%" (
    echo ERROR: ComfyUI not found at %COMFYUI_DIR%
    echo Run setup first: py tools/setup_comfyui.py
    pause
    exit /b 1
)

echo ComfyUI location: %COMFYUI_DIR%
echo.
echo Starting with NVIDIA GPU...
echo Access at: http://127.0.0.1:8188
echo.

cd /d "%COMFYUI_DIR%"
call run_nvidia_gpu.bat
