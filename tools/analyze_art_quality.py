"""
美術品質分析工具
分析所有 sprite 的像素密度、顏色豐富度、邊緣清晰度
"""
from PIL import Image, ImageFilter
import os
import math

def analyze_sprite(path):
    img = Image.open(path).convert('RGBA')
    w, h = img.size
    pixels = list(img.getdata())
    
    non_transparent = [p for p in pixels if p[3] > 10]
    density = len(non_transparent) / (w * h) * 100
    
    # 顏色豐富度（唯一顏色數）
    unique_colors = len(set((p[0]//16, p[1]//16, p[2]//16) for p in non_transparent))
    
    # 邊緣清晰度（用 Laplacian 方差）
    gray = img.convert('L')
    lap = gray.filter(ImageFilter.FIND_EDGES)
    lap_pixels = list(lap.getdata())
    sharpness = sum(p*p for p in lap_pixels) / len(lap_pixels) if lap_pixels else 0
    
    return {
        'size': f'{w}x{h}',
        'density': density,
        'unique_colors': unique_colors,
        'sharpness': math.sqrt(sharpness),
    }

def score_sprite(info):
    score = 0
    # 密度分（0-40）
    score += min(40, info['density'] * 0.7)
    # 顏色分（0-30）
    score += min(30, info['unique_colors'] * 0.5)
    # 清晰度分（0-30）
    score += min(30, info['sharpness'] * 0.3)
    return score

print("=" * 60)
print("美術品質分析報告")
print("=" * 60)

# 角色 sprites
print("\n【角色 Sprites】")
chars_dir = 'client/chiikawa-pixel/assets/sprites/characters'
for f in sorted(os.listdir(chars_dir)):
    if f.endswith('.png') and not f.endswith('.import') and 'ref' not in f and 'backup' not in f:
        path = os.path.join(chars_dir, f)
        info = analyze_sprite(path)
        score = score_sprite(info)
        status = '✅' if score >= 60 else '⚠️' if score >= 40 else '❌'
        print(f"  {status} {f}: {info['size']}, 密度={info['density']:.1f}%, 顏色={info['unique_colors']}, 清晰={info['sharpness']:.1f}, 分={score:.0f}")

# 目標物 sprites
print("\n【目標物 Sprites】")
targets_dir = 'client/chiikawa-pixel/assets/sprites/targets'
low_quality = []
for f in sorted(os.listdir(targets_dir)):
    if f.endswith('.png') and not f.endswith('.import') and 'sheet' not in f and '_swim' not in f:
        path = os.path.join(targets_dir, f)
        info = analyze_sprite(path)
        score = score_sprite(info)
        status = '✅' if score >= 60 else '⚠️' if score >= 40 else '❌'
        print(f"  {status} {f}: {info['size']}, 密度={info['density']:.1f}%, 顏色={info['unique_colors']}, 分={score:.0f}")
        if score < 60:
            low_quality.append((f, score, info))

print(f"\n【需要優化的目標物】（分數 < 60）")
for f, score, info in sorted(low_quality, key=lambda x: x[1]):
    print(f"  ❌ {f}: 分={score:.0f} (密度={info['density']:.1f}%, 顏色={info['unique_colors']})")

# 背景
print("\n【背景 Sprites】")
bg_dir = 'client/chiikawa-pixel/assets/sprites/backgrounds'
for f in sorted(os.listdir(bg_dir)):
    if f.endswith('.png') and not f.endswith('.import'):
        path = os.path.join(bg_dir, f)
        info = analyze_sprite(path)
        score = score_sprite(info)
        status = '✅' if score >= 60 else '⚠️'
        print(f"  {status} {f}: {info['size']}, 顏色={info['unique_colors']}, 分={score:.0f}")
