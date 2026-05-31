"""
generate_combo_sfx_day341.py — DAY-341 Combo 里程碑音效生成
生成 4 個 Combo 里程碑音效：combo_5, combo_10, combo_20, combo_30
每個音效都比前一個更震撼
"""
import struct
import math
import os

OUT_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\audio\sfx"

def write_wav(filename, samples, sample_rate=44100):
    """寫入 WAV 檔案"""
    num_samples = len(samples)
    # 正規化到 16-bit
    max_val = max(abs(s) for s in samples) if samples else 1.0
    if max_val < 0.001:
        max_val = 1.0
    int_samples = [int(s / max_val * 32767) for s in samples]
    
    with open(filename, 'wb') as f:
        # RIFF header
        f.write(b'RIFF')
        f.write(struct.pack('<I', 36 + num_samples * 2))
        f.write(b'WAVE')
        # fmt chunk
        f.write(b'fmt ')
        f.write(struct.pack('<I', 16))
        f.write(struct.pack('<H', 1))   # PCM
        f.write(struct.pack('<H', 1))   # mono
        f.write(struct.pack('<I', sample_rate))
        f.write(struct.pack('<I', sample_rate * 2))
        f.write(struct.pack('<H', 2))   # block align
        f.write(struct.pack('<H', 16))  # bits per sample
        # data chunk
        f.write(b'data')
        f.write(struct.pack('<I', num_samples * 2))
        for s in int_samples:
            f.write(struct.pack('<h', max(-32768, min(32767, s))))
    print(f"  ✅ {os.path.basename(filename)} ({num_samples} samples, {num_samples/sample_rate:.2f}s)")

def sine(freq, t, sr=44100):
    return math.sin(2 * math.pi * freq * t / sr)

def envelope(t, total, attack=0.01, decay=0.05, sustain=0.7, release=0.2):
    """ADSR 包絡"""
    a = attack * total
    d = decay * total
    s_end = (1.0 - release) * total
    r = release * total
    if t < a:
        return t / a
    elif t < a + d:
        return 1.0 - (t - a) / d * (1.0 - sustain)
    elif t < s_end:
        return sustain
    else:
        return sustain * (1.0 - (t - s_end) / r)

SR = 44100

# ── combo_5.wav：輕快上升音階（5連擊）──────────────────────────
def gen_combo_5():
    """5連擊：輕快的三音上升（C-E-G）"""
    samples = []
    notes = [523.25, 659.25, 783.99]  # C5, E5, G5
    note_dur = int(SR * 0.08)
    gap = int(SR * 0.02)
    
    for freq in notes:
        for i in range(note_dur):
            t = i / SR
            env = math.exp(-t * 15)
            s = sine(freq, i) * env * 0.6
            s += sine(freq * 2, i) * env * 0.2  # 泛音
            samples.append(s)
        samples.extend([0.0] * gap)
    
    # 尾音
    tail_dur = int(SR * 0.15)
    for i in range(tail_dur):
        t = i / SR
        env = math.exp(-t * 8)
        samples.append(sine(1046.5, i) * env * 0.4)  # C6
    
    return samples

# ── combo_10.wav：更強的上升音階 + 和弦（10連擊）──────────────
def gen_combo_10():
    """10連擊：快速音階 + 和弦爆發"""
    samples = []
    # 快速音階
    notes = [523.25, 587.33, 659.25, 698.46, 783.99]  # C-D-E-F-G
    note_dur = int(SR * 0.06)
    gap = int(SR * 0.01)
    
    for freq in notes:
        for i in range(note_dur):
            env = math.exp(-i / SR * 20)
            s = sine(freq, i) * env * 0.5
            s += sine(freq * 1.5, i) * env * 0.15
            samples.append(s)
        samples.extend([0.0] * gap)
    
    # 和弦爆發
    chord_dur = int(SR * 0.3)
    chord_freqs = [783.99, 987.77, 1174.66]  # G5, B5, D6
    for i in range(chord_dur):
        t = i / SR
        env = math.exp(-t * 5)
        s = 0.0
        for freq in chord_freqs:
            s += sine(freq, i) * env * 0.25
        samples.append(s)
    
    return samples

# ── combo_20.wav：電子感上升 + 強力和弦（20連擊）──────────────
def gen_combo_20():
    """20連擊：電子感音效 + 強力和弦 + 震撼低音"""
    samples = []
    
    # 電子感上升掃頻
    sweep_dur = int(SR * 0.2)
    for i in range(sweep_dur):
        t = i / SR
        freq = 300 + (1200 - 300) * (t / 0.2) ** 2
        env = t / 0.2 * 0.8
        s = sine(freq, i) * env
        s += sine(freq * 2, i) * env * 0.3
        # 方波成分（電子感）
        s += (1.0 if sine(freq, i) > 0 else -1.0) * env * 0.1
        samples.append(s)
    
    # 強力和弦爆發
    chord_dur = int(SR * 0.4)
    chord_freqs = [523.25, 659.25, 783.99, 1046.5]  # C-E-G-C 大三和弦
    for i in range(chord_dur):
        t = i / SR
        env = math.exp(-t * 4)
        s = 0.0
        for freq in chord_freqs:
            s += sine(freq, i) * env * 0.2
        # 低音衝擊
        if i < int(SR * 0.05):
            s += sine(130.81, i) * (1.0 - i / (SR * 0.05)) * 0.5
        samples.append(s)
    
    return samples

# ── combo_30.wav：最強里程碑（30連擊 MAX）──────────────────────
def gen_combo_30():
    """30連擊 MAX：震撼的多層音效 + 勝利號角"""
    samples = []
    
    # 低音衝擊
    impact_dur = int(SR * 0.05)
    for i in range(impact_dur):
        t = i / SR
        env = 1.0 - t / 0.05
        s = sine(80, i) * env * 0.8
        s += sine(160, i) * env * 0.4
        # 噪音衝擊
        import random
        s += (random.random() * 2 - 1) * env * 0.3
        samples.append(s)
    
    # 上升掃頻（更快更強）
    sweep_dur = int(SR * 0.15)
    for i in range(sweep_dur):
        t = i / SR
        freq = 200 + (2000 - 200) * (t / 0.15) ** 1.5
        env = t / 0.15
        s = sine(freq, i) * env * 0.7
        s += sine(freq * 2, i) * env * 0.3
        s += sine(freq * 3, i) * env * 0.1
        samples.append(s)
    
    # 勝利號角（大三和弦 + 八度）
    fanfare_dur = int(SR * 0.5)
    fanfare_notes = [
        (523.25, 0.0, 0.5),    # C5 全程
        (659.25, 0.05, 0.5),   # E5 稍後
        (783.99, 0.10, 0.5),   # G5 再後
        (1046.5, 0.15, 0.5),   # C6 最後
    ]
    for i in range(fanfare_dur):
        t = i / SR
        env = math.exp(-t * 3)
        s = 0.0
        for freq, delay, vol in fanfare_notes:
            if t >= delay:
                local_t = t - delay
                local_env = math.exp(-local_t * 4)
                s += sine(freq, int((t - delay) * SR)) * local_env * vol * 0.25
        # 低音持續
        s += sine(130.81, i) * math.exp(-t * 6) * 0.3
        samples.append(s)
    
    # 閃光尾音
    tail_dur = int(SR * 0.2)
    for i in range(tail_dur):
        t = i / SR
        env = math.exp(-t * 10)
        s = sine(2093.0, i) * env * 0.3  # C7 高音
        s += sine(1568.0, i) * env * 0.2  # G6
        samples.append(s)
    
    return samples

def main():
    print("=== DAY-341 Combo 里程碑音效生成 ===")
    os.makedirs(OUT_DIR, exist_ok=True)
    
    combos = [
        ("combo_5.wav", gen_combo_5, "5連擊輕快音階"),
        ("combo_10.wav", gen_combo_10, "10連擊和弦爆發"),
        ("combo_20.wav", gen_combo_20, "20連擊電子感"),
        ("combo_30.wav", gen_combo_30, "30連擊勝利號角"),
    ]
    
    for filename, gen_func, desc in combos:
        print(f"\n生成 {filename} ({desc})...")
        samples = gen_func()
        path = os.path.join(OUT_DIR, filename)
        write_wav(path, samples)
    
    print("\n=== 全部完成 ===")
    print(f"輸出目錄：{OUT_DIR}")

if __name__ == "__main__":
    main()
