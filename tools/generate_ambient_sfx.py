"""
generate_ambient_sfx.py — 生成海底環境音效
- underwater_ambient.wav：海底環境音（低頻水流 + 氣泡 + 遠距離水聲）
- bubble_pop.wav：氣泡破裂音（BubbleLayer 視覺對應）

用法：py tools/generate_ambient_sfx.py
"""
import numpy as np
import wave
import os

SAMPLE_RATE = 44100
OUTPUT_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\audio\sfx"


def write_wav(filename: str, samples: np.ndarray):
    """寫入 WAV 檔案"""
    os.makedirs(os.path.dirname(filename), exist_ok=True)
    samples = np.clip(samples, -1.0, 1.0)
    samples_int = (samples * 32767).astype(np.int16)
    with wave.open(filename, 'w') as f:
        f.setnchannels(1)
        f.setsampwidth(2)
        f.setframerate(SAMPLE_RATE)
        f.writeframes(samples_int.tobytes())


def gen_underwater_ambient(duration: float = 8.0) -> np.ndarray:
    """
    海底環境音：
    - 低頻水流（60-120 Hz 正弦波，緩慢調製）
    - 隨機氣泡（短暫高頻脈衝）
    - 遠距離水聲（帶通濾波白噪音）
    - 整體音量很低（沉浸感，不搶主音效）
    """
    n = int(SAMPLE_RATE * duration)
    t = np.linspace(0, duration, n, False)
    rng = np.random.default_rng(42)

    # 1. 低頻水流（60 Hz 基頻 + 緩慢 LFO 調製）
    lfo = 0.5 + 0.5 * np.sin(2 * np.pi * 0.15 * t)  # 0.15 Hz LFO
    water_flow = np.sin(2 * np.pi * 65 * t) * 0.08 * lfo
    water_flow += np.sin(2 * np.pi * 90 * t) * 0.05 * lfo
    water_flow += np.sin(2 * np.pi * 120 * t) * 0.03 * lfo

    # 2. 遠距離水聲（帶通濾波白噪音，200-800 Hz）
    raw_noise = rng.uniform(-1, 1, n)
    # 簡單 IIR 帶通濾波（低通 + 高通近似）
    # 低通：y[n] = 0.05 * x[n] + 0.95 * y[n-1]
    lp = np.zeros(n)
    lp[0] = raw_noise[0] * 0.05
    for i in range(1, n):
        lp[i] = 0.05 * raw_noise[i] + 0.95 * lp[i-1]
    # 高通：y[n] = x[n] - lp[n]（去除 DC 和極低頻）
    hp = raw_noise - lp
    distant_water = hp * 0.04

    # 3. 隨機氣泡（每 0.5-2 秒一個，短暫高頻脈衝）
    bubbles = np.zeros(n)
    bubble_times = []
    t_pos = 0.3
    while t_pos < duration - 0.1:
        bubble_times.append(t_pos)
        t_pos += rng.uniform(0.4, 1.8)

    for bt in bubble_times:
        idx = int(bt * SAMPLE_RATE)
        bubble_dur = int(0.04 * SAMPLE_RATE)  # 40ms 氣泡
        if idx + bubble_dur >= n:
            break
        # 氣泡：上升音調 (400→800 Hz) + 快速衰減
        bt_t = np.linspace(0, 0.04, bubble_dur, False)
        freqs = np.linspace(400, 800, bubble_dur)
        phase = np.cumsum(2 * np.pi * freqs / SAMPLE_RATE)
        bubble_wave = np.sin(phase)
        # 快速衰減包絡
        env = np.exp(-bt_t * 60)
        amplitude = rng.uniform(0.03, 0.08)
        bubbles[idx:idx+bubble_dur] += bubble_wave * env * amplitude

    # 4. 合成
    ambient = water_flow + distant_water + bubbles

    # 5. 整體淡入淡出（避免點擊聲）
    fade_samples = int(0.3 * SAMPLE_RATE)
    ambient[:fade_samples] *= np.linspace(0, 1, fade_samples)
    ambient[-fade_samples:] *= np.linspace(1, 0, fade_samples)

    return ambient


def gen_bubble_pop() -> np.ndarray:
    """
    氣泡破裂音（對應 BubbleLayer 視覺氣泡）
    - 短暫（0.15 秒）
    - 上升音調 + 快速衰減
    - 輕柔，不搶主音效
    """
    duration = 0.15
    n = int(SAMPLE_RATE * duration)
    t = np.linspace(0, duration, n, False)

    # 上升音調 (300 → 900 Hz)
    freqs = np.linspace(300, 900, n)
    phase = np.cumsum(2 * np.pi * freqs / SAMPLE_RATE)
    wave = np.sin(phase)

    # 快速衰減包絡（attack 極短，快速 decay）
    env = np.exp(-t * 35)
    env[:int(0.003 * SAMPLE_RATE)] = np.linspace(0, 1, int(0.003 * SAMPLE_RATE))

    # 加一點噪音讓氣泡聲更真實
    rng = np.random.default_rng(7)
    noise = rng.uniform(-1, 1, n) * 0.15 * env

    result = (wave * env * 0.25 + noise)
    return result


if __name__ == "__main__":
    print("🌊 生成海底環境音效...")

    # 生成 8 秒環境音（循環播放）
    ambient = gen_underwater_ambient(8.0)
    path = os.path.join(OUTPUT_DIR, "underwater_ambient.wav")
    write_wav(path, ambient)
    print(f"  ✅ underwater_ambient.wav (8.0s, {len(ambient)/SAMPLE_RATE:.1f}s)")

    # 生成氣泡破裂音
    bubble = gen_bubble_pop()
    path = os.path.join(OUTPUT_DIR, "bubble_pop.wav")
    write_wav(path, bubble)
    print(f"  ✅ bubble_pop.wav ({len(bubble)/SAMPLE_RATE:.2f}s)")

    print("\n✅ 環境音效生成完畢！")
    print("整合方式：")
    print("  1. AudioManager 加入 SFX.BUBBLE_POP 和 BGM.UNDERWATER_AMBIENT")
    print("  2. BackgroundManager 在海底狀態播放 underwater_ambient（循環）")
    print("  3. BubbleLayer 氣泡消失時播放 bubble_pop（低音量）")
