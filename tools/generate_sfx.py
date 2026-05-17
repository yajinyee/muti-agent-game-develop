"""
8-bit 音效生成器（規格書 10章）
用 numpy + scipy 生成 WAV 格式的復古音效
"""
import numpy as np
import wave
import struct
import os

OUTPUT_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\audio\sfx"
BGM_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\audio\bgm"

SAMPLE_RATE = 44100

def write_wav(filename: str, samples: np.ndarray, sample_rate: int = SAMPLE_RATE):
    """寫入 WAV 檔案"""
    os.makedirs(os.path.dirname(filename), exist_ok=True)
    # 正規化到 16-bit
    samples = np.clip(samples, -1.0, 1.0)
    samples_int = (samples * 32767).astype(np.int16)

    with wave.open(filename, 'w') as f:
        f.setnchannels(1)
        f.setsampwidth(2)
        f.setframerate(sample_rate)
        f.writeframes(samples_int.tobytes())

def square_wave(freq: float, duration: float, volume: float = 0.5) -> np.ndarray:
    """方波（8-bit 音效特徵）"""
    t = np.linspace(0, duration, int(SAMPLE_RATE * duration), False)
    wave = np.sign(np.sin(2 * np.pi * freq * t))
    return wave * volume

def triangle_wave(freq: float, duration: float, volume: float = 0.5) -> np.ndarray:
    """三角波（柔和 8-bit）"""
    t = np.linspace(0, duration, int(SAMPLE_RATE * duration), False)
    wave = 2 * np.abs(2 * (t * freq - np.floor(t * freq + 0.5))) - 1
    return wave * volume

def noise(duration: float, volume: float = 0.3) -> np.ndarray:
    """白噪音（爆炸/擊破音效）"""
    samples = int(SAMPLE_RATE * duration)
    return np.random.uniform(-1, 1, samples) * volume

def envelope(samples: np.ndarray, attack: float = 0.01, decay: float = 0.1,
             sustain: float = 0.7, release: float = 0.1) -> np.ndarray:
    """ADSR 包絡"""
    n = len(samples)
    a = int(attack * SAMPLE_RATE)
    d = int(decay * SAMPLE_RATE)
    r = int(release * SAMPLE_RATE)
    s = n - a - d - r

    env = np.zeros(n)
    if a > 0:
        env[:a] = np.linspace(0, 1, a)
    if d > 0:
        env[a:a+d] = np.linspace(1, sustain, d)
    if s > 0:
        env[a+d:a+d+s] = sustain
    if r > 0:
        env[a+d+s:] = np.linspace(sustain, 0, r)

    return samples * env

def pitch_slide(freq_start: float, freq_end: float, duration: float,
                volume: float = 0.5, wave_type: str = "square") -> np.ndarray:
    """音調滑動（上升/下降音效）"""
    t = np.linspace(0, duration, int(SAMPLE_RATE * duration), False)
    freqs = np.linspace(freq_start, freq_end, len(t))
    phase = np.cumsum(2 * np.pi * freqs / SAMPLE_RATE)
    if wave_type == "square":
        return np.sign(np.sin(phase)) * volume
    return np.sin(phase) * volume

def concat(*arrays) -> np.ndarray:
    """連接音效片段"""
    return np.concatenate(arrays)

def silence(duration: float) -> np.ndarray:
    return np.zeros(int(SAMPLE_RATE * duration))

# ---- 生成各種音效 ----

def gen_attack_fire():
    """攻擊發射音（輕快 8-bit 揮棒）"""
    s = pitch_slide(440, 880, 0.08, 0.4, "square")
    return envelope(s, attack=0.005, decay=0.05, sustain=0.3, release=0.025)

def gen_attack_fire_hachiware():
    """小八攻擊（清脆斬擊）"""
    s = pitch_slide(660, 1320, 0.06, 0.4, "square")
    return envelope(s, attack=0.003, decay=0.04, sustain=0.2, release=0.017)

def gen_attack_fire_usagi():
    """烏薩奇攻擊（高亢連擊）"""
    s1 = pitch_slide(880, 1760, 0.04, 0.4, "square")
    s2 = pitch_slide(1100, 2200, 0.04, 0.3, "square")
    return envelope(concat(s1, silence(0.02), s2), attack=0.002, decay=0.03, sustain=0.2, release=0.01)

def gen_hit():
    """命中音（8-bit 嗶波）"""
    s = square_wave(440, 0.05, 0.5)
    return envelope(s, attack=0.002, decay=0.03, sustain=0.3, release=0.018)

def gen_kill():
    """擊破音（像素爆裂）"""
    n = noise(0.1, 0.6)
    s = pitch_slide(880, 110, 0.15, 0.4, "square")
    combined = concat(n[:len(n)//2], s)
    return envelope(combined, attack=0.001, decay=0.05, sustain=0.2, release=0.05)

def gen_coin_drop():
    """金幣跳動"""
    s1 = square_wave(1047, 0.05, 0.4)  # C6
    s2 = square_wave(1319, 0.05, 0.4)  # E6
    s3 = square_wave(1568, 0.08, 0.4)  # G6
    return concat(
        envelope(s1, 0.002, 0.02, 0.3, 0.028),
        silence(0.02),
        envelope(s2, 0.002, 0.02, 0.3, 0.028),
        silence(0.02),
        envelope(s3, 0.002, 0.03, 0.3, 0.048)
    )

def gen_reward_bag():
    """報酬袋彈出"""
    s = pitch_slide(330, 660, 0.12, 0.5, "square")
    return envelope(s, attack=0.01, decay=0.05, sustain=0.4, release=0.06)

def gen_boss_warning():
    """BOSS 警告（低頻警報）"""
    s1 = square_wave(110, 0.3, 0.6)
    s2 = square_wave(165, 0.3, 0.6)
    s3 = square_wave(110, 0.3, 0.6)
    return concat(
        envelope(s1, 0.01, 0.1, 0.5, 0.19),
        silence(0.1),
        envelope(s2, 0.01, 0.1, 0.5, 0.19),
        silence(0.1),
        envelope(s3, 0.01, 0.1, 0.5, 0.19)
    )

def gen_bonus_ready():
    """Bonus Ready 提示（上升音階）"""
    notes = [523, 659, 784, 1047]  # C5 E5 G5 C6
    parts = []
    for note in notes:
        s = square_wave(note, 0.1, 0.4)
        parts.append(envelope(s, 0.005, 0.03, 0.4, 0.065))
        parts.append(silence(0.02))
    return concat(*parts)

def gen_weed_pull():
    """拔草音（快速拔起）"""
    s = pitch_slide(220, 440, 0.08, 0.5, "square")
    n = noise(0.03, 0.2)
    return envelope(concat(s, n), attack=0.002, decay=0.04, sustain=0.2, release=0.038)

def gen_big_win():
    """大獎音效（華麗勝利）"""
    # 上升音階 + 和弦
    notes = [523, 659, 784, 1047, 1319, 1568]
    parts = []
    for i, note in enumerate(notes):
        s = square_wave(note, 0.12, 0.4)
        parts.append(envelope(s, 0.005, 0.04, 0.4, 0.075))
        if i < len(notes) - 1:
            parts.append(silence(0.01))
    # 最後和弦
    chord = (square_wave(1047, 0.5, 0.3) +
             square_wave(1319, 0.5, 0.3) +
             square_wave(1568, 0.5, 0.3))
    parts.append(envelope(chord, 0.01, 0.1, 0.6, 0.39))
    return concat(*parts)

# ---- 生成簡單 BGM ----

def gen_bgm_main():
    """主遊戲 BGM（輕快 8-bit 冒險感）"""
    # 簡單的 8-bit 旋律
    melody = [
        (523, 0.2), (659, 0.2), (784, 0.2), (659, 0.2),
        (523, 0.2), (523, 0.2), (440, 0.4),
        (494, 0.2), (523, 0.2), (659, 0.2), (784, 0.2),
        (880, 0.4), (784, 0.2), (659, 0.2),
    ]
    parts = []
    for freq, dur in melody:
        s = square_wave(freq, dur * 0.8, 0.3)
        parts.append(envelope(s, 0.01, 0.05, 0.5, dur * 0.2 - 0.01))
        parts.append(silence(dur * 0.2))
    return concat(*parts)

def gen_bgm_boss():
    """BOSS BGM（詭異低頻）"""
    bass = [
        (110, 0.3), (110, 0.3), (123, 0.3), (110, 0.3),
        (98, 0.3), (110, 0.3), (98, 0.6),
    ]
    parts = []
    for freq, dur in bass:
        s = square_wave(freq, dur * 0.9, 0.4)
        parts.append(envelope(s, 0.01, 0.05, 0.6, dur * 0.1))
        parts.append(silence(dur * 0.1))
    return concat(*parts)

def gen_bgm_bonus():
    """Bonus BGM（快速歡樂）"""
    melody = [
        (784, 0.15), (880, 0.15), (988, 0.15), (1047, 0.15),
        (988, 0.15), (880, 0.15), (784, 0.3),
        (659, 0.15), (784, 0.15), (880, 0.15), (988, 0.15),
        (1047, 0.3), (880, 0.15), (784, 0.15),
    ]
    parts = []
    for freq, dur in melody:
        s = triangle_wave(freq, dur * 0.8, 0.35)
        parts.append(envelope(s, 0.005, 0.03, 0.5, dur * 0.2 - 0.005))
        parts.append(silence(dur * 0.2))
    return concat(*parts)

if __name__ == "__main__":
    print("🎵 生成 8-bit 音效...")

    sfx_list = [
        ("attack_fire.wav", gen_attack_fire()),
        ("attack_fire_hachiware.wav", gen_attack_fire_hachiware()),
        ("attack_fire_usagi.wav", gen_attack_fire_usagi()),
        ("hit.wav", gen_hit()),
        ("kill.wav", gen_kill()),
        ("coin_drop.wav", gen_coin_drop()),
        ("reward_bag.wav", gen_reward_bag()),
        ("boss_warning.wav", gen_boss_warning()),
        ("bonus_ready.wav", gen_bonus_ready()),
        ("weed_pull.wav", gen_weed_pull()),
        ("big_win.wav", gen_big_win()),
    ]

    print("\n[音效 SFX]")
    for filename, samples in sfx_list:
        path = os.path.join(OUTPUT_DIR, filename)
        write_wav(path, samples)
        duration = len(samples) / SAMPLE_RATE
        print(f"  ✅ {filename} ({duration:.2f}s)")

    print("\n[背景音樂 BGM]")
    bgm_list = [
        ("main_game.wav", gen_bgm_main()),
        ("boss_enter.wav", gen_bgm_boss()),
        ("bonus_game.wav", gen_bgm_bonus()),
    ]
    for filename, samples in bgm_list:
        path = os.path.join(BGM_DIR, filename)
        os.makedirs(BGM_DIR, exist_ok=True)
        write_wav(path, samples)
        duration = len(samples) / SAMPLE_RATE
        print(f"  ✅ {filename} ({duration:.2f}s)")

    print("\n✅ 所有音效生成完畢！")
