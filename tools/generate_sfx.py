#!/usr/bin/env python3
"""
generate_sfx.py — 程式生成遊戲音效 WAV 檔案
sfx-agent 負責維護
使用純 Python 標準庫，不需要任何外部依賴
輸出到 client/chiikawa-pixel/assets/audio/sfx/
"""

import struct
import math
import os
import random

OUTPUT_DIR = os.path.join(
    os.path.dirname(os.path.dirname(os.path.abspath(__file__))),
    "client", "chiikawa-pixel", "assets", "audio", "sfx"
)

SAMPLE_RATE = 22050
CHANNELS = 1
BITS = 16
MAX_AMP = 32767


def write_wav(filename: str, samples: list[float]) -> None:
    """將 float 樣本列表（-1.0 ~ 1.0）寫入 WAV 檔案"""
    os.makedirs(OUTPUT_DIR, exist_ok=True)
    path = os.path.join(OUTPUT_DIR, filename)
    num_samples = len(samples)
    data_size = num_samples * 2  # 16-bit = 2 bytes per sample
    with open(path, "wb") as f:
        # RIFF header
        f.write(b"RIFF")
        f.write(struct.pack("<I", 36 + data_size))
        f.write(b"WAVE")
        # fmt chunk
        f.write(b"fmt ")
        f.write(struct.pack("<I", 16))       # chunk size
        f.write(struct.pack("<H", 1))        # PCM
        f.write(struct.pack("<H", CHANNELS))
        f.write(struct.pack("<I", SAMPLE_RATE))
        f.write(struct.pack("<I", SAMPLE_RATE * CHANNELS * 2))  # byte rate
        f.write(struct.pack("<H", CHANNELS * 2))  # block align
        f.write(struct.pack("<H", BITS))
        # data chunk
        f.write(b"data")
        f.write(struct.pack("<I", data_size))
        for s in samples:
            val = int(max(-1.0, min(1.0, s)) * MAX_AMP)
            f.write(struct.pack("<h", val))
    print(f"  ✅ {filename} ({num_samples} samples, {num_samples/SAMPLE_RATE:.2f}s)")


def sine(freq: float, t: float) -> float:
    return math.sin(2 * math.pi * freq * t)


def envelope(t: float, total: float, attack: float = 0.01, release: float = 0.1) -> float:
    if t < attack:
        return t / attack
    if t > total - release:
        return (total - t) / release
    return 1.0


def gen_attack_fire(char_id: str = "chiikawa") -> list[float]:
    """射擊音效：短促的「啵」聲"""
    duration = 0.12
    n = int(SAMPLE_RATE * duration)
    samples = []
    base_freq = {"chiikawa": 880, "hachiware": 660, "usagi": 1100}.get(char_id, 880)
    for i in range(n):
        t = i / SAMPLE_RATE
        env = envelope(t, duration, 0.005, 0.08)
        freq = base_freq * (1.0 - t / duration * 0.3)  # 頻率下滑
        s = sine(freq, t) * 0.6
        s += sine(freq * 2.1, t) * 0.2  # 泛音
        noise = (random.random() * 2 - 1) * 0.1
        samples.append((s + noise) * env * 0.7)
    return samples


def gen_hit() -> list[float]:
    """命中音效：短促打擊聲"""
    duration = 0.08
    n = int(SAMPLE_RATE * duration)
    samples = []
    for i in range(n):
        t = i / SAMPLE_RATE
        env = envelope(t, duration, 0.002, 0.06)
        freq = 300 * (1.0 - t / duration * 0.5)
        s = sine(freq, t) * 0.5
        noise = (random.random() * 2 - 1) * 0.4
        samples.append((s + noise) * env * 0.8)
    return samples


def gen_kill() -> list[float]:
    """擊破音效：爆炸感的「碰」聲"""
    duration = 0.25
    n = int(SAMPLE_RATE * duration)
    samples = []
    for i in range(n):
        t = i / SAMPLE_RATE
        env = envelope(t, duration, 0.003, 0.18)
        freq = 200 * math.exp(-t * 8)  # 快速衰減頻率
        s = sine(freq, t) * 0.4
        s += sine(freq * 1.5, t) * 0.3
        noise = (random.random() * 2 - 1) * 0.5
        samples.append((s + noise) * env * 0.9)
    return samples


def gen_big_win() -> list[float]:
    """大獎音效：上升的「叮叮叮」"""
    duration = 0.6
    n = int(SAMPLE_RATE * duration)
    samples = []
    notes = [523, 659, 784, 1047]  # C5 E5 G5 C6
    note_dur = duration / len(notes)
    for i in range(n):
        t = i / SAMPLE_RATE
        note_idx = min(int(t / note_dur), len(notes) - 1)
        note_t = t - note_idx * note_dur
        env = envelope(note_t, note_dur, 0.01, note_dur * 0.5)
        freq = notes[note_idx]
        s = sine(freq, t) * 0.5
        s += sine(freq * 2, t) * 0.2
        s += sine(freq * 3, t) * 0.1
        samples.append(s * env * 0.8)
    return samples


def gen_coin_drop() -> list[float]:
    """金幣音效：清脆的「叮」聲"""
    duration = 0.3
    n = int(SAMPLE_RATE * duration)
    samples = []
    for i in range(n):
        t = i / SAMPLE_RATE
        env = math.exp(-t * 12)
        freq = 1200 + 200 * math.exp(-t * 5)
        s = sine(freq, t) * 0.6
        s += sine(freq * 2.76, t) * 0.2  # 金屬泛音
        samples.append(s * env * 0.7)
    return samples


def gen_boss_warning() -> list[float]:
    """BOSS 警告音效：低沉的警報聲"""
    duration = 1.2
    n = int(SAMPLE_RATE * duration)
    samples = []
    for i in range(n):
        t = i / SAMPLE_RATE
        # 兩段警報
        cycle = t % 0.4
        env = envelope(cycle, 0.4, 0.05, 0.15)
        freq = 220 + 30 * math.sin(2 * math.pi * 2 * t)  # 顫音
        s = sine(freq, t) * 0.5
        s += sine(freq * 0.5, t) * 0.3  # 低頻加強
        samples.append(s * env * 0.8)
    return samples


def gen_boss_enter() -> list[float]:
    """BOSS 進場音效：震撼的低頻衝擊"""
    duration = 0.8
    n = int(SAMPLE_RATE * duration)
    samples = []
    for i in range(n):
        t = i / SAMPLE_RATE
        env = envelope(t, duration, 0.01, 0.5)
        # 低頻衝擊
        freq = 80 * math.exp(-t * 2)
        s = sine(freq, t) * 0.6
        s += sine(freq * 2, t) * 0.3
        noise = (random.random() * 2 - 1) * 0.3 * math.exp(-t * 5)
        samples.append((s + noise) * env * 0.9)
    return samples


def gen_bonus_ready() -> list[float]:
    """Bonus 準備音效：歡快的上升音階"""
    duration = 0.5
    n = int(SAMPLE_RATE * duration)
    samples = []
    notes = [392, 494, 587, 784]  # G4 B4 D5 G5
    note_dur = duration / len(notes)
    for i in range(n):
        t = i / SAMPLE_RATE
        note_idx = min(int(t / note_dur), len(notes) - 1)
        note_t = t - note_idx * note_dur
        env = envelope(note_t, note_dur, 0.01, note_dur * 0.4)
        freq = notes[note_idx]
        s = sine(freq, t) * 0.5
        s += sine(freq * 2, t) * 0.15
        samples.append(s * env * 0.8)
    return samples


def gen_bonus_game() -> list[float]:
    """Bonus 遊戲中音效：歡快短音"""
    duration = 0.15
    n = int(SAMPLE_RATE * duration)
    samples = []
    for i in range(n):
        t = i / SAMPLE_RATE
        env = envelope(t, duration, 0.01, 0.08)
        freq = 660 + 200 * t / duration
        s = sine(freq, t) * 0.5
        samples.append(s * env * 0.7)
    return samples


def gen_weed_pull() -> list[float]:
    """拔草音效：「嗖」的拔出聲"""
    duration = 0.18
    n = int(SAMPLE_RATE * duration)
    samples = []
    for i in range(n):
        t = i / SAMPLE_RATE
        env = envelope(t, duration, 0.005, 0.12)
        freq = 400 + 600 * (t / duration)  # 頻率上升
        s = sine(freq, t) * 0.3
        noise = (random.random() * 2 - 1) * 0.5 * (1 - t / duration)
        samples.append((s + noise) * env * 0.7)
    return samples


def main():
    print("🎵 生成遊戲音效 WAV 檔案...")
    print(f"   輸出目錄：{OUTPUT_DIR}")
    print()

    sfx_list = [
        ("attack_fire.wav",          gen_attack_fire("chiikawa")),
        ("attack_fire_hachiware.wav", gen_attack_fire("hachiware")),
        ("attack_fire_usagi.wav",     gen_attack_fire("usagi")),
        ("hit.wav",                   gen_hit()),
        ("kill.wav",                  gen_kill()),
        ("big_win.wav",               gen_big_win()),
        ("coin_drop.wav",             gen_coin_drop()),
        ("boss_warning.wav",          gen_boss_warning()),
        ("boss_enter.wav",            gen_boss_enter()),
        ("bonus_ready.wav",           gen_bonus_ready()),
        ("bonus_game.wav",            gen_bonus_game()),
        ("weed_pull.wav",             gen_weed_pull()),
    ]

    for filename, samples in sfx_list:
        write_wav(filename, samples)

    print()
    print(f"✅ 完成！共生成 {len(sfx_list)} 個音效檔案")


if __name__ == "__main__":
    main()
