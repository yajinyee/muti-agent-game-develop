"""
generate_bonus_sfx_v2.py — 升級 Bonus 音效
- bonus_ready.wav v2：更有興奮感的 Bonus 觸發音效
  - 快速上升音階 + 和弦爆發 + 金屬撞擊感
- bonus_trigger.wav：Bonus 遊戲開始時的短促爆發音
  - 用於 BonusGame 開始瞬間

用法：py tools/generate_bonus_sfx_v2.py
"""
import numpy as np
import wave
import os

SAMPLE_RATE = 44100
OUTPUT_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\audio\sfx"


def write_wav(filename: str, samples: np.ndarray):
    os.makedirs(os.path.dirname(filename), exist_ok=True)
    samples = np.clip(samples, -1.0, 1.0)
    samples_int = (samples * 32767).astype(np.int16)
    with wave.open(filename, 'w') as f:
        f.setnchannels(1)
        f.setsampwidth(2)
        f.setframerate(SAMPLE_RATE)
        f.writeframes(samples_int.tobytes())


def square_wave(freq: float, duration: float, volume: float = 0.5) -> np.ndarray:
    t = np.linspace(0, duration, int(SAMPLE_RATE * duration), False)
    return np.sign(np.sin(2 * np.pi * freq * t)) * volume


def triangle_wave(freq: float, duration: float, volume: float = 0.5) -> np.ndarray:
    t = np.linspace(0, duration, int(SAMPLE_RATE * duration), False)
    w = 2 * np.abs(2 * (t * freq - np.floor(t * freq + 0.5))) - 1
    return w * volume


def noise(duration: float, volume: float = 0.3) -> np.ndarray:
    return np.random.default_rng(99).uniform(-1, 1, int(SAMPLE_RATE * duration)) * volume


def envelope(samples: np.ndarray, attack: float = 0.01, decay: float = 0.1,
             sustain: float = 0.7, release: float = 0.1) -> np.ndarray:
    n = len(samples)
    a = int(attack * SAMPLE_RATE)
    d = int(decay * SAMPLE_RATE)
    r = int(release * SAMPLE_RATE)
    s = max(0, n - a - d - r)
    env = np.zeros(n)
    if a > 0: env[:a] = np.linspace(0, 1, a)
    if d > 0: env[a:a+d] = np.linspace(1, sustain, d)
    if s > 0: env[a+d:a+d+s] = sustain
    if r > 0: env[a+d+s:a+d+s+r] = np.linspace(sustain, 0, r)
    return samples * env


def pitch_slide(f0: float, f1: float, duration: float, volume: float = 0.5) -> np.ndarray:
    n = int(SAMPLE_RATE * duration)
    freqs = np.linspace(f0, f1, n)
    phase = np.cumsum(2 * np.pi * freqs / SAMPLE_RATE)
    return np.sign(np.sin(phase)) * volume


def concat(*arrays) -> np.ndarray:
    return np.concatenate(arrays)


def silence(duration: float) -> np.ndarray:
    return np.zeros(int(SAMPLE_RATE * duration))


def gen_bonus_ready_v2() -> np.ndarray:
    """
    Bonus Ready v2 — 更有興奮感
    設計：
    1. 快速上升音階（C5→G5→C6→E6，每音 0.07s，加速感）
    2. 短暫靜音（0.03s，製造期待感）
    3. 和弦爆發（C6+E6+G6 三音同時，0.4s，帶 noise 衝擊）
    4. 尾音上揚（G6→C7，0.15s，結束感）
    """
    # 1. 快速上升音階（比 v1 更快，更有衝勁）
    notes_up = [
        (523, 0.07),   # C5
        (659, 0.06),   # E5
        (784, 0.06),   # G5
        (1047, 0.07),  # C6
        (1319, 0.06),  # E6
    ]
    parts = []
    for freq, dur in notes_up:
        s = square_wave(freq, dur * 0.85, 0.38)
        parts.append(envelope(s, 0.003, 0.02, 0.5, dur * 0.15 - 0.003))
        parts.append(silence(dur * 0.15))

    # 2. 短暫靜音（期待感）
    parts.append(silence(0.04))

    # 3. 和弦爆發（三音 + 噪音衝擊）
    chord_dur = 0.45
    chord_n = int(SAMPLE_RATE * chord_dur)
    c6 = square_wave(1047, chord_dur, 0.28)
    e6 = square_wave(1319, chord_dur, 0.22)
    g6 = square_wave(1568, chord_dur, 0.18)
    n_short = noise(chord_dur * 0.08, 0.35)  # 短暫衝擊噪音
    n_padded = np.zeros(chord_n)
    n_padded[:len(n_short)] = n_short
    # 確保長度一致
    min_len = min(len(c6), len(e6), len(g6), chord_n)
    chord = c6[:min_len] + e6[:min_len] + g6[:min_len] + n_padded[:min_len]
    parts.append(envelope(chord, 0.005, 0.08, 0.65, 0.35))

    # 4. 尾音上揚（G6→C7）
    tail = pitch_slide(1568, 2093, 0.18, 0.3)
    parts.append(envelope(tail, 0.003, 0.04, 0.4, 0.13))

    result = concat(*parts)
    # 整體音量正規化
    peak = np.max(np.abs(result))
    if peak > 0:
        result = result / peak * 0.85
    return result


def gen_bonus_trigger() -> np.ndarray:
    """
    Bonus 遊戲開始瞬間的爆發音（0.3 秒）
    設計：
    - 噪音衝擊（0.05s）+ 上升滑音（C5→C7，0.2s）+ 和弦尾音（0.05s）
    - 比 bonus_ready 更短促、更有爆發感
    """
    # 噪音衝擊
    n = noise(0.05, 0.5)
    n_env = envelope(n, 0.002, 0.03, 0.3, 0.018)

    # 快速上升滑音（C5→C7，兩個八度）
    slide = pitch_slide(523, 2093, 0.18, 0.45)
    slide_env = envelope(slide, 0.003, 0.05, 0.5, 0.127)

    # 和弦尾音（C7+E7，短暫）
    c7 = square_wave(2093, 0.07, 0.25)
    e7 = square_wave(2637, 0.07, 0.18)
    chord_tail = envelope(c7 + e7, 0.002, 0.02, 0.4, 0.048)

    result = concat(n_env, slide_env, chord_tail)
    peak = np.max(np.abs(result))
    if peak > 0:
        result = result / peak * 0.88
    return result


def gen_bonus_end_fanfare() -> np.ndarray:
    """
    Bonus 結束結算音效（0.6 秒）
    設計：
    - 下降音階（C7→C5）+ 最後和弦（勝利感）
    """
    notes_down = [
        (2093, 0.08),  # C7
        (1568, 0.07),  # G6
        (1319, 0.07),  # E6
        (1047, 0.08),  # C6
        (784, 0.07),   # G5
        (523, 0.07),   # C5
    ]
    parts = []
    for freq, dur in notes_down:
        s = triangle_wave(freq, dur * 0.8, 0.35)
        parts.append(envelope(s, 0.003, 0.02, 0.45, dur * 0.2 - 0.003))
        parts.append(silence(dur * 0.2))

    # 最後和弦（C5+G5+C6，勝利感）
    chord_dur = 0.25
    c5 = triangle_wave(523, chord_dur, 0.28)
    g5 = triangle_wave(784, chord_dur, 0.22)
    c6 = triangle_wave(1047, chord_dur, 0.18)
    chord = c5 + g5 + c6
    parts.append(envelope(chord, 0.005, 0.06, 0.6, 0.185))

    result = concat(*parts)
    peak = np.max(np.abs(result))
    if peak > 0:
        result = result / peak * 0.82
    return result


if __name__ == "__main__":
    print("🎵 升級 Bonus 音效...")

    # bonus_ready.wav v2（覆蓋舊版）
    samples = gen_bonus_ready_v2()
    path = os.path.join(OUTPUT_DIR, "bonus_ready.wav")
    write_wav(path, samples)
    print(f"  ✅ bonus_ready.wav v2 ({len(samples)/SAMPLE_RATE:.2f}s) — 快速音階 + 和弦爆發")

    # bonus_trigger.wav（新增）
    samples = gen_bonus_trigger()
    path = os.path.join(OUTPUT_DIR, "bonus_trigger.wav")
    write_wav(path, samples)
    print(f"  ✅ bonus_trigger.wav ({len(samples)/SAMPLE_RATE:.2f}s) — Bonus 開始爆發音")

    # bonus_end.wav（新增）
    samples = gen_bonus_end_fanfare()
    path = os.path.join(OUTPUT_DIR, "bonus_end.wav")
    write_wav(path, samples)
    print(f"  ✅ bonus_end.wav ({len(samples)/SAMPLE_RATE:.2f}s) — Bonus 結算音效")

    print("\n✅ Bonus 音效升級完畢！")
    print("整合方式：")
    print("  1. AudioManager 加入 SFX.BONUS_TRIGGER 和 SFX.BONUS_END")
    print("  2. BonusGame.gd 在 bonus_start 時播放 BONUS_TRIGGER")
    print("  3. BonusGame.gd 在 bonus_end 時播放 BONUS_END")
