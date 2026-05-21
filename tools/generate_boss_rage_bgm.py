"""
generate_boss_rage_bgm.py
從 boss_enter.wav 生成 boss_rage.wav（加速 + 升調）
方法：修改 WAV frame rate（提高 frame rate = 加速播放 + 升調）
不需要任何第三方套件，純 Python 標準庫
"""

import wave
import struct
import os

INPUT_PATH = r"d:\Kiro\client\chiikawa-pixel\assets\audio\bgm\boss_enter.wav"
OUTPUT_PATH = r"d:\Kiro\client\chiikawa-pixel\assets\audio\bgm\boss_rage.wav"

# 加速倍率（1.15 = 加速 15%，同時升調約 2 個半音）
SPEED_FACTOR = 1.15


def generate_boss_rage():
    if not os.path.exists(INPUT_PATH):
        print(f"[ERROR] 找不到輸入檔案: {INPUT_PATH}")
        return False

    with wave.open(INPUT_PATH, 'rb') as src:
        params = src.getparams()
        frames = src.readframes(params.nframes)
        original_framerate = params.framerate
        nchannels = params.nchannels
        sampwidth = params.sampwidth

    # 新的 frame rate（提高 = 加速 + 升調）
    new_framerate = int(original_framerate * SPEED_FACTOR)

    print(f"[INFO] 輸入: {INPUT_PATH}")
    print(f"[INFO] 原始 frame rate: {original_framerate} Hz")
    print(f"[INFO] 新 frame rate: {new_framerate} Hz (×{SPEED_FACTOR})")
    print(f"[INFO] 效果: 加速 {(SPEED_FACTOR-1)*100:.0f}%，升調約 {SPEED_FACTOR*12:.1f} 半音")

    # 寫入新 WAV（只改 frame rate，不改音頻數據）
    with wave.open(OUTPUT_PATH, 'wb') as dst:
        dst.setnchannels(nchannels)
        dst.setsampwidth(sampwidth)
        dst.setframerate(new_framerate)
        dst.writeframes(frames)

    # 驗證輸出
    with wave.open(OUTPUT_PATH, 'rb') as check:
        check_params = check.getparams()
        duration = check_params.nframes / check_params.framerate
        print(f"[OK] 輸出: {OUTPUT_PATH}")
        print(f"[OK] 時長: {duration:.2f}s（原始 {check_params.nframes / original_framerate:.2f}s）")
        print(f"[OK] Frame rate: {check_params.framerate} Hz")

    return True


if __name__ == "__main__":
    success = generate_boss_rage()
    if success:
        print("\n✅ boss_rage.wav 生成成功！")
        print("   效果：比 boss_enter.wav 快 15%，音調更高，緊張感更強")
    else:
        print("\n❌ 生成失敗")
