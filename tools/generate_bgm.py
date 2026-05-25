#!/usr/bin/env python3
"""
generate_bgm.py — 程式生成遊戲 BGM WAV 檔案
bgm-agent 負責維護
使用純 Python 標準庫，不需要任何外部依賴
輸出到 client/chiikawa-pixel/assets/audio/bgm/
"""

import struct
import math
import os
import random

OUTPUT_DIR = os.path.join(
    os.path.dirname(os.path.dirname(os.path.abspath(__file__))),
    "client", "chiikawa-pixel", "assets", "audio", "bgm"
)

SAMPLE_RATE = 22050
CHANNELS = 1
BITS = 16
MAX_AMP = 32767


def write_wav(filename: str, samples: list[float]) -> None:
    os.makedirs(OUTPUT_DIR, exist_ok=True)
    path = os.path.join(OUTPUT_DIR, filename)
    num_samples = len(samples)
    data_size = num_samples * 2
    with open(path, "wb") as f:
        f.write(b"RIFF")
        f.write(struct.pack("<I", 36 + data_size))
        f.write(b"WAVE")
        f.write(b"fmt ")
        f.write(struct.pack("<I", 16))
        f.write(struct.pack("<H", 1))
        f.write(struct.pack("<H", CHANNELS))
        f.write(struct.pack("<I", SAMPLE_RATE))
        f.write(struct.pack("<I", SAMPLE_RATE * CHANNELS * 2))
        f.write(struct.pack("<H", CHANNELS * 2))
        f.write(struct.pack("<H", BITS))
        f.write(b"data")
        f.write(struct.pack("<I", data_size))
        for s in samples:
            val = int(max(-1.0, min(1.0, s)) * MAX_AMP)
            f.write(struct.pack("<h", val))
    print(f"  ✅ {filename} ({num_samples/SAMPLE_RATE:.1f}s)")


def sine(freq: float, t: float, phase: float = 0.0) -> float:
    return math.sin(2 * math.pi * freq * t + phase)


def note_freq(semitone: int, base: float = 261.63) -> float:
    """從半音數計算頻率（C4 = 261.63 Hz）"""
    return base * (2 ** (semitone / 12.0))


# ── 音符序列輔助 ──────────────────────────────────────────────

def gen_melody_samples(notes: list, bpm: float, duration_beats: float) -> list[float]:
    """
    notes: [(semitone, beats, volume), ...]
    semitone: 相對 C4 的半音數（None = 休止符）
    """
    beat_dur = 60.0 / bpm
    total_dur = duration_beats * beat_dur
    n = int(SAMPLE_RATE * total_dur)
    samples = [0.0] * n

    t_offset = 0.0
    for note_info in notes:
        semitone, beats, vol = note_info
        note_dur = beats * beat_dur
        note_n = int(SAMPLE_RATE * note_dur)
        start_i = int(t_offset * SAMPLE_RATE)
        if semitone is not None:
            freq = note_freq(semitone)
            for j in range(note_n):
                if start_i + j >= n:
                    break
                t = j / SAMPLE_RATE
                # ADSR 包絡
                attack = min(0.02, note_dur * 0.1)
                release = min(0.05, note_dur * 0.3)
                if t < attack:
                    env = t / attack
                elif t > note_dur - release:
                    env = (note_dur - t) / release
                else:
                    env = 1.0
                s = sine(freq, t) * 0.5 + sine(freq * 2, t) * 0.2 + sine(freq * 3, t) * 0.1
                samples[start_i + j] += s * env * vol
        t_offset += note_dur

    return samples


def gen_bass_samples(notes: list, bpm: float, duration_beats: float) -> list[float]:
    beat_dur = 60.0 / bpm
    total_dur = duration_beats * beat_dur
    n = int(SAMPLE_RATE * total_dur)
    samples = [0.0] * n

    t_offset = 0.0
    for note_info in notes:
        semitone, beats, vol = note_info
        note_dur = beats * beat_dur
        note_n = int(SAMPLE_RATE * note_dur)
        start_i = int(t_offset * SAMPLE_RATE)
        if semitone is not None:
            freq = note_freq(semitone - 12)  # 低八度
            for j in range(note_n):
                if start_i + j >= n:
                    break
                t = j / SAMPLE_RATE
                attack = min(0.01, note_dur * 0.05)
                release = min(0.08, note_dur * 0.4)
                if t < attack:
                    env = t / attack
                elif t > note_dur - release:
                    env = (note_dur - t) / release
                else:
                    env = 1.0
                s = sine(freq, t) * 0.7 + sine(freq * 2, t) * 0.15
                samples[start_i + j] += s * env * vol
        t_offset += note_dur

    return samples


def mix(a: list[float], b: list[float], vol_a: float = 1.0, vol_b: float = 1.0) -> list[float]:
    n = max(len(a), len(b))
    result = [0.0] * n
    for i in range(len(a)):
        result[i] += a[i] * vol_a
    for i in range(len(b)):
        result[i] += b[i] * vol_b
    # 軟限幅
    peak = max(abs(s) for s in result) if result else 1.0
    if peak > 0.9:
        factor = 0.9 / peak
        result = [s * factor for s in result]
    return result


def loop_samples(samples: list[float], target_duration: float) -> list[float]:
    """將樣本循環到目標時長"""
    target_n = int(SAMPLE_RATE * target_duration)
    result = []
    while len(result) < target_n:
        result.extend(samples)
    return result[:target_n]


# ── BGM 生成 ──────────────────────────────────────────────────

def gen_main_game_bgm() -> list[float]:
    """主遊戲 BGM：輕快的海底風格，8 秒循環"""
    bpm = 120
    # 主旋律（C 大調，輕快）
    melody = [
        (0, 0.5, 0.7),   # C4
        (2, 0.5, 0.7),   # D4
        (4, 0.5, 0.7),   # E4
        (5, 0.5, 0.7),   # F4
        (7, 1.0, 0.8),   # G4
        (5, 0.5, 0.6),   # F4
        (4, 0.5, 0.6),   # E4
        (2, 1.0, 0.6),   # D4
        (0, 0.5, 0.7),   # C4
        (4, 0.5, 0.7),   # E4
        (7, 0.5, 0.8),   # G4
        (9, 0.5, 0.8),   # A4
        (7, 1.0, 0.7),   # G4
        (5, 0.5, 0.6),   # F4
        (4, 0.5, 0.6),   # E4
        (0, 1.0, 0.7),   # C4
    ]
    bass = [
        (0, 1.0, 0.5),   # C
        (7, 1.0, 0.5),   # G
        (5, 1.0, 0.5),   # F
        (7, 1.0, 0.5),   # G
        (0, 1.0, 0.5),   # C
        (7, 1.0, 0.5),   # G
        (5, 1.0, 0.5),   # F
        (7, 1.0, 0.5),   # G
    ]
    total_beats = sum(n[1] for n in melody)
    mel = gen_melody_samples(melody, bpm, total_beats)
    bas = gen_bass_samples(bass, bpm, total_beats)
    combined = mix(mel, bas, 0.7, 0.4)
    return loop_samples(combined, 8.0)


def gen_boss_enter_bgm() -> list[float]:
    """BOSS 進場 BGM：緊張的低頻，6 秒循環"""
    bpm = 140
    melody = [
        (0, 0.25, 0.8),
        (0, 0.25, 0.8),
        (-1, 0.5, 0.9),
        (0, 0.25, 0.8),
        (-2, 0.25, 0.7),
        (-3, 1.0, 0.9),
        (0, 0.25, 0.8),
        (0, 0.25, 0.8),
        (-1, 0.5, 0.9),
        (-3, 0.5, 0.8),
        (-5, 1.0, 1.0),
        (None, 0.5, 0),
    ]
    bass = [
        (-12, 0.5, 0.8),
        (-12, 0.5, 0.8),
        (-13, 0.5, 0.9),
        (-15, 0.5, 0.9),
        (-12, 0.5, 0.8),
        (-12, 0.5, 0.8),
        (-13, 0.5, 0.9),
        (-17, 0.5, 1.0),
    ]
    total_beats = sum(n[1] for n in melody)
    mel = gen_melody_samples(melody, bpm, total_beats)
    bas = gen_bass_samples(bass, bpm, total_beats)
    combined = mix(mel, bas, 0.6, 0.5)
    return loop_samples(combined, 6.0)


def gen_boss_rage_bgm() -> list[float]:
    """BOSS 狂暴 BGM：更快更激烈，5 秒循環"""
    bpm = 170
    melody = [
        (0, 0.25, 0.9),
        (-1, 0.25, 0.9),
        (0, 0.25, 0.9),
        (-2, 0.25, 0.8),
        (-3, 0.5, 1.0),
        (-5, 0.5, 1.0),
        (-3, 0.25, 0.9),
        (-2, 0.25, 0.9),
        (-1, 0.5, 0.9),
        (0, 0.5, 1.0),
        (-3, 0.5, 0.9),
        (-5, 0.5, 1.0),
    ]
    bass = [
        (-12, 0.25, 0.9),
        (-12, 0.25, 0.9),
        (-13, 0.25, 0.9),
        (-15, 0.25, 0.9),
        (-12, 0.25, 0.9),
        (-12, 0.25, 0.9),
        (-13, 0.25, 0.9),
        (-17, 0.25, 1.0),
        (-12, 0.25, 0.9),
        (-12, 0.25, 0.9),
        (-13, 0.25, 0.9),
        (-15, 0.25, 0.9),
    ]
    total_beats = sum(n[1] for n in melody)
    mel = gen_melody_samples(melody, bpm, total_beats)
    bas = gen_bass_samples(bass, bpm, total_beats)
    combined = mix(mel, bas, 0.6, 0.5)
    return loop_samples(combined, 5.0)


def gen_bonus_game_bgm() -> list[float]:
    """Bonus 遊戲 BGM：歡快活潑，7 秒循環"""
    bpm = 150
    melody = [
        (0, 0.5, 0.8),   # C
        (4, 0.5, 0.8),   # E
        (7, 0.5, 0.9),   # G
        (12, 0.5, 0.9),  # C5
        (9, 0.5, 0.8),   # A
        (7, 0.5, 0.8),   # G
        (5, 0.5, 0.7),   # F
        (4, 0.5, 0.7),   # E
        (2, 0.5, 0.7),   # D
        (4, 0.5, 0.8),   # E
        (5, 0.5, 0.8),   # F
        (7, 0.5, 0.9),   # G
        (9, 0.5, 0.9),   # A
        (7, 0.5, 0.8),   # G
        (5, 0.5, 0.7),   # F
        (0, 1.0, 0.8),   # C
    ]
    bass = [
        (0, 1.0, 0.5),
        (7, 1.0, 0.5),
        (5, 1.0, 0.5),
        (4, 1.0, 0.5),
        (0, 1.0, 0.5),
        (7, 1.0, 0.5),
        (5, 1.0, 0.5),
        (0, 1.0, 0.5),
    ]
    total_beats = sum(n[1] for n in melody)
    mel = gen_melody_samples(melody, bpm, total_beats)
    bas = gen_bass_samples(bass, bpm, total_beats)
    combined = mix(mel, bas, 0.7, 0.4)
    return loop_samples(combined, 7.0)


def main():
    print("🎵 生成遊戲 BGM WAV 檔案...")
    print(f"   輸出目錄：{OUTPUT_DIR}")
    print()

    bgm_list = [
        ("main_game.wav",   gen_main_game_bgm()),
        ("boss_enter.wav",  gen_boss_enter_bgm()),
        ("boss_rage.wav",   gen_boss_rage_bgm()),
        ("bonus_game.wav",  gen_bonus_game_bgm()),
    ]

    for filename, samples in bgm_list:
        write_wav(filename, samples)

    print()
    print(f"✅ 完成！共生成 {len(bgm_list)} 個 BGM 檔案")


if __name__ == "__main__":
    main()
