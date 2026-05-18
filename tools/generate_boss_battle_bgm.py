"""
generate_boss_battle_bgm.py — 生成 BOSS 戰持續 BGM
- boss_battle.wav：BOSS Phase 1 戰鬥 BGM（緊張、重複循環、低頻驅動）
  - 設計：低頻方波 bass + 緊張旋律 + 打擊節奏感
  - 長度：8 秒（循環播放）

用法：py tools/generate_boss_battle_bgm.py
"""
import numpy as np
import wave
import os

SAMPLE_RATE = 44100
BGM_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\audio\bgm"


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
    n = int(SAMPLE_RATE * duration)
    t = np.linspace(0, duration, n, False)
    return np.sign(np.sin(2 * np.pi * freq * t)) * volume


def triangle_wave(freq: float, duration: float, volume: float = 0.5) -> np.ndarray:
    n = int(SAMPLE_RATE * duration)
    t = np.linspace(0, duration, n, False)
    w = 2 * np.abs(2 * (t * freq - np.floor(t * freq + 0.5))) - 1
    return w * volume


def noise_burst(duration: float, volume: float = 0.3) -> np.ndarray:
    rng = np.random.default_rng(42)
    return rng.uniform(-1, 1, int(SAMPLE_RATE * duration)) * volume


def envelope(samples: np.ndarray, attack: float = 0.005, decay: float = 0.05,
             sustain: float = 0.8, release: float = 0.05) -> np.ndarray:
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


def silence(duration: float) -> np.ndarray:
    return np.zeros(int(SAMPLE_RATE * duration))


def gen_boss_battle_bgm(duration: float = 8.0) -> np.ndarray:
    """
    BOSS 戰 BGM（Phase 1）— 緊張、重複循環
    設計：
    1. 低頻 bass 驅動（110 Hz 方波，每拍 0.25s，強弱交替）
    2. 緊張旋律線（三角波，不和諧音程製造壓迫感）
    3. 打擊節奏（噪音短脈衝，模擬鼓點）
    4. 整體 8 秒可無縫循環
    """
    n = int(SAMPLE_RATE * duration)
    result = np.zeros(n)

    # ── 1. Bass 驅動（每 0.25s 一拍，強弱交替）──────────────────────────────
    # 音符序列：A2(110) D3(147) A2(110) E3(165) A2(110) D3(147) G2(98) A2(110)
    bass_notes = [110, 147, 110, 165, 110, 147, 98, 110,
                  110, 147, 110, 165, 110, 147, 98, 110,
                  110, 147, 110, 165, 110, 147, 98, 110,
                  110, 147, 110, 165, 110, 147, 98, 110]
    beat_dur = duration / len(bass_notes)

    for i, freq in enumerate(bass_notes):
        start = int(i * beat_dur * SAMPLE_RATE)
        seg_n = int(beat_dur * SAMPLE_RATE * 0.85)  # 85% 佔空比，留一點間隙
        if start + seg_n > n:
            break
        # 強拍（偶數）音量更大
        vol = 0.35 if i % 2 == 0 else 0.22
        seg = square_wave(freq, beat_dur * 0.85, vol)
        seg = envelope(seg, 0.003, 0.04, 0.75, 0.05)
        result[start:start+len(seg)] += seg

    # ── 2. 緊張旋律線（三角波，每 0.5s 一音）────────────────────────────────
    # 使用不和諧音程（小二度、增四度）製造壓迫感
    melody_notes = [
        (220, 0.4),   # A3
        (233, 0.4),   # Bb3（小二度，不和諧）
        (220, 0.4),   # A3
        (311, 0.4),   # Eb4（增四度，魔鬼音程）
        (220, 0.4),   # A3
        (207, 0.4),   # Ab3
        (196, 0.4),   # G3
        (220, 0.8),   # A3（長音，製造懸念）
        (220, 0.4),
        (233, 0.4),
        (220, 0.4),
        (311, 0.4),
        (220, 0.4),
        (207, 0.4),
        (185, 0.4),   # F#3
        (220, 0.8),
    ]
    pos = 0.0
    for freq, dur in melody_notes:
        start = int(pos * SAMPLE_RATE)
        seg_n = int(dur * SAMPLE_RATE * 0.9)
        if start + seg_n > n:
            break
        seg = triangle_wave(freq, dur * 0.9, 0.18)
        seg = envelope(seg, 0.01, 0.05, 0.7, 0.1)
        result[start:start+len(seg)] += seg
        pos += dur

    # ── 3. 打擊節奏（噪音短脈衝，每 0.5s 一次）──────────────────────────────
    drum_interval = 0.5
    drum_positions = np.arange(0, duration, drum_interval)
    for dp in drum_positions:
        start = int(dp * SAMPLE_RATE)
        drum_n = int(0.06 * SAMPLE_RATE)
        if start + drum_n > n:
            break
        # 低頻噪音（模擬鼓）
        drum = noise_burst(0.06, 0.25)
        # 快速衰減
        env_arr = np.exp(-np.linspace(0, 8, drum_n))
        drum = drum * env_arr
        result[start:start+drum_n] += drum

    # ── 4. 高頻緊張顫音（每 2 秒一次，製造緊張感）───────────────────────────
    tremolo_positions = [1.5, 3.5, 5.5, 7.5]
    for tp in tremolo_positions:
        start = int(tp * SAMPLE_RATE)
        trem_dur = 0.3
        trem_n = int(trem_dur * SAMPLE_RATE)
        if start + trem_n > n:
            break
        # 快速顫音（440 Hz 方波，每 0.02s 開關）
        t = np.linspace(0, trem_dur, trem_n, False)
        trem_env = (np.sin(2 * np.pi * 25 * t) > 0).astype(float)  # 25 Hz 開關
        trem = square_wave(440, trem_dur, 0.12) * trem_env
        trem = envelope(trem, 0.005, 0.05, 0.6, 0.1)
        result[start:start+trem_n] += trem

    # ── 5. 淡入淡出（無縫循環）────────────────────────────────────────────────
    fade = int(0.15 * SAMPLE_RATE)
    result[:fade] *= np.linspace(0, 1, fade)
    result[-fade:] *= np.linspace(1, 0, fade)

    # 正規化
    peak = np.max(np.abs(result))
    if peak > 0:
        result = result / peak * 0.82

    return result


if __name__ == "__main__":
    print("🎵 生成 BOSS 戰 BGM...")

    samples = gen_boss_battle_bgm(8.0)
    path = os.path.join(BGM_DIR, "boss_battle.wav")
    write_wav(path, samples)
    print(f"  ✅ boss_battle.wav ({len(samples)/SAMPLE_RATE:.1f}s) — 緊張低頻 bass + 不和諧旋律 + 打擊節奏")
    print(f"  大小：{os.path.getsize(path)/1024:.1f} KB")
    print()
    print("整合方式：")
    print("  1. AudioManager BGM 枚舉加入 BOSS_BATTLE")
    print("  2. BackgroundManager boss_battle 狀態播放 BOSS_BATTLE BGM")
    print("  3. boss_event phase_change 時切換到 BOSS_RAGE BGM")
