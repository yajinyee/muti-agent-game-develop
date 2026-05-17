#!/usr/bin/env python3
"""
Animation Pipeline Tool
吉伊卡哇：像素大討伐 動畫品質檢查與生成工具

使用方式：
  py tools/animation_pipeline.py --audit          # 審查所有動畫
  py tools/animation_pipeline.py --gif chiikawa   # 生成 GIF 預覽
  py tools/animation_pipeline.py --check <path>   # 檢查單個 spritesheet
  py tools/animation_pipeline.py --audit --report # 審查並輸出報告
"""

import argparse
import json
import os
import sys
from datetime import datetime
from pathlib import Path

try:
    from PIL import Image
    import numpy as np
except ImportError:
    print("[ERROR] 缺少依賴套件，請執行：pip install Pillow numpy")
    sys.exit(1)

# ─── 常數設定 ───────────────────────────────────────────────────────────────

PROJECT_ROOT = Path(__file__).parent.parent
SPRITES_DIR = PROJECT_ROOT / "client" / "chiikawa-pixel" / "assets" / "sprites" / "characters"
REPORTS_DIR = PROJECT_ROOT / "reports" / "animation"
PREVIEW_DIR = REPORTS_DIR / "preview"

CHARACTERS = ["chiikawa", "hachiware", "usagi"]
ANIMATION_STATES = ["idle", "attack", "hit", "hurt", "bigwin", "skill", "bonus", "fail"]

# 各動畫狀態的預期幀數
EXPECTED_FRAMES = {
    "idle": (4, 8),
    "attack": (6, 8),
    "hit": (4, 4),
    "hurt": (3, 4),
    "bigwin": (8, 12),
    "skill": (8, 8),
    "bonus": (6, 8),
    "fail": (4, 4),
}

# 各動畫狀態的 FPS
ANIMATION_FPS = {
    "idle": 8,
    "attack": 12,
    "hit": 10,
    "hurt": 10,
    "bigwin": 10,
    "skill": 12,
    "bonus": 10,
    "fail": 8,
}

# 一致性檢查容差
TOLERANCES = {
    "anchor_tolerance": 2,       # px
    "bottom_alignment_tolerance": 3,  # px
    "silhouette_similarity": 85,  # %
    "head_ratio_tolerance": 5,    # %
    "weapon_position_tolerance": 4,  # px
    "color_drift_max": 5.0,       # ΔE
    "deformation_max": 10.0,      # %
    "jitter_max": 2,              # px
}


# ─── 核心功能 ────────────────────────────────────────────────────────────────

def check_frame_consistency(sheet_path: str) -> dict:
    """
    檢查 spritesheet 各幀一致性
    
    Args:
        sheet_path: spritesheet 圖片路徑
    
    Returns:
        dict: 包含各項檢查結果和整體分數的字典
    """
    sheet_path = Path(sheet_path)
    
    if not sheet_path.exists():
        return {
            "error": f"檔案不存在：{sheet_path}",
            "score": 0,
            "passed": False
        }
    
    try:
        sheet = Image.open(sheet_path).convert("RGBA")
    except Exception as e:
        return {
            "error": f"無法開啟圖片：{e}",
            "score": 0,
            "passed": False
        }
    
    sheet_w, sheet_h = sheet.size
    
    # 嘗試推斷幀數（假設幀是正方形或已知比例）
    # 如果寬度是高度的整數倍，則幀數 = 寬/高
    if sheet_w % sheet_h == 0:
        frame_count = sheet_w // sheet_h
        frame_w = sheet_h
        frame_h = sheet_h
    else:
        # 嘗試常見幀數
        for fc in [4, 6, 8, 12, 3]:
            if sheet_w % fc == 0:
                frame_count = fc
                frame_w = sheet_w // fc
                frame_h = sheet_h
                break
        else:
            frame_count = 1
            frame_w = sheet_w
            frame_h = sheet_h
    
    # 切割各幀
    frames = []
    for i in range(frame_count):
        frame = sheet.crop((i * frame_w, 0, (i + 1) * frame_w, frame_h))
        frames.append(frame)
    
    results = {
        "file": str(sheet_path),
        "sheet_size": f"{sheet_w}x{sheet_h}",
        "frame_count": frame_count,
        "frame_size": f"{frame_w}x{frame_h}",
        "checks": {},
        "issues": [],
        "score": 0,
        "passed": False
    }
    
    checks = results["checks"]
    issues = results["issues"]
    
    # ── 檢查 1：canvas_size（所有幀尺寸相同）──
    sizes = [f.size for f in frames]
    canvas_ok = len(set(sizes)) == 1
    checks["canvas_size"] = {
        "passed": canvas_ok,
        "detail": f"所有幀尺寸：{set(sizes)}"
    }
    if not canvas_ok:
        issues.append("幀尺寸不一致")
    
    # ── 檢查 2：transparent_bg（背景透明）──
    bg_issues = 0
    for i, frame in enumerate(frames):
        arr = np.array(frame)
        # 檢查四個角落是否透明
        corners = [arr[0, 0], arr[0, -1], arr[-1, 0], arr[-1, -1]]
        for corner in corners:
            if corner[3] > 10:  # alpha > 10 視為不透明
                bg_issues += 1
                break
    
    bg_ok = bg_issues == 0
    checks["transparent_bg"] = {
        "passed": bg_ok,
        "detail": f"背景不透明幀數：{bg_issues}/{frame_count}"
    }
    if not bg_ok:
        issues.append(f"{bg_issues} 幀背景不透明")
    
    # ── 檢查 3：anchor_point（底部中心一致）──
    anchors = []
    for frame in frames:
        bbox = frame.getbbox()
        if bbox:
            anchor_x = (bbox[0] + bbox[2]) // 2
            anchor_y = bbox[3]
            anchors.append((anchor_x, anchor_y))
        else:
            anchors.append(None)
    
    valid_anchors = [a for a in anchors if a is not None]
    if len(valid_anchors) >= 2:
        base_ax, base_ay = valid_anchors[0]
        max_drift = max(
            abs(ax - base_ax) + abs(ay - base_ay)
            for ax, ay in valid_anchors
        )
        anchor_ok = max_drift <= TOLERANCES["anchor_tolerance"]
    else:
        max_drift = 0
        anchor_ok = True
    
    checks["anchor_point"] = {
        "passed": anchor_ok,
        "detail": f"最大偏移：±{max_drift}px（容差：±{TOLERANCES['anchor_tolerance']}px）"
    }
    if not anchor_ok:
        issues.append(f"anchor point 偏移過大（{max_drift}px）")
    
    # ── 檢查 4：bottom_alignment（底部對齊）──
    bottoms = []
    for frame in frames:
        bbox = frame.getbbox()
        if bbox:
            bottoms.append(bbox[3])
    
    if len(bottoms) >= 2:
        bottom_range = max(bottoms) - min(bottoms)
        bottom_ok = bottom_range <= TOLERANCES["bottom_alignment_tolerance"]
    else:
        bottom_range = 0
        bottom_ok = True
    
    checks["bottom_alignment"] = {
        "passed": bottom_ok,
        "detail": f"底部偏差範圍：{bottom_range}px（容差：{TOLERANCES['bottom_alignment_tolerance']}px）"
    }
    if not bottom_ok:
        issues.append(f"底部對齊偏差過大（{bottom_range}px）")
    
    # ── 檢查 5：silhouette（輪廓相似度）──
    if len(frames) >= 2:
        base_arr = np.array(frames[0])
        base_mask = (base_arr[:, :, 3] > 10).astype(float)
        
        similarities = []
        for frame in frames[1:]:
            arr = np.array(frame)
            mask = (arr[:, :, 3] > 10).astype(float)
            
            intersection = np.sum(base_mask * mask)
            union = np.sum(np.clip(base_mask + mask, 0, 1))
            
            if union > 0:
                iou = intersection / union * 100
            else:
                iou = 100.0
            similarities.append(iou)
        
        min_similarity = min(similarities) if similarities else 100.0
        silhouette_ok = min_similarity >= TOLERANCES["silhouette_similarity"]
    else:
        min_similarity = 100.0
        silhouette_ok = True
    
    checks["silhouette"] = {
        "passed": silhouette_ok,
        "detail": f"最低輪廓相似度：{min_similarity:.1f}%（門檻：{TOLERANCES['silhouette_similarity']}%）"
    }
    if not silhouette_ok:
        issues.append(f"輪廓相似度過低（{min_similarity:.1f}%）")
    
    # ── 檢查 6：deformation（形變檢查）──
    pixel_counts = []
    for frame in frames:
        arr = np.array(frame)
        count = np.sum(arr[:, :, 3] > 10)
        pixel_counts.append(count)
    
    if len(pixel_counts) >= 2 and max(pixel_counts) > 0:
        deform_ratio = (max(pixel_counts) - min(pixel_counts)) / max(pixel_counts) * 100
        deform_ok = deform_ratio <= TOLERANCES["deformation_max"]
    else:
        deform_ratio = 0.0
        deform_ok = True
    
    checks["deformation"] = {
        "passed": deform_ok,
        "detail": f"像素數量差異：{deform_ratio:.1f}%（容差：{TOLERANCES['deformation_max']}%）"
    }
    if not deform_ok:
        issues.append(f"形變過大（{deform_ratio:.1f}%）")
    
    # ── 檢查 7：jitter（抖動檢查）──
    bboxes = []
    for frame in frames:
        bbox = frame.getbbox()
        if bbox:
            cx = (bbox[0] + bbox[2]) // 2
            cy = (bbox[1] + bbox[3]) // 2
            bboxes.append((cx, cy))
    
    if len(bboxes) >= 2:
        max_jitter = 0
        for i in range(1, len(bboxes)):
            dx = abs(bboxes[i][0] - bboxes[i-1][0])
            dy = abs(bboxes[i][1] - bboxes[i-1][1])
            max_jitter = max(max_jitter, dx, dy)
        jitter_ok = max_jitter <= TOLERANCES["jitter_max"]
    else:
        max_jitter = 0
        jitter_ok = True
    
    checks["jitter"] = {
        "passed": jitter_ok,
        "detail": f"最大抖動：±{max_jitter}px（容差：±{TOLERANCES['jitter_max']}px）"
    }
    if not jitter_ok:
        issues.append(f"幀間抖動過大（{max_jitter}px）")
    
    # ── 計算整體分數 ──
    check_weights = {
        "canvas_size": 20,
        "transparent_bg": 15,
        "anchor_point": 15,
        "bottom_alignment": 15,
        "silhouette": 15,
        "deformation": 10,
        "jitter": 10,
    }
    
    total_weight = sum(check_weights.values())
    earned_weight = sum(
        w for k, w in check_weights.items()
        if checks.get(k, {}).get("passed", False)
    )
    
    score = int(earned_weight / total_weight * 100)
    results["score"] = score
    results["passed"] = score >= 85
    
    return results


def generate_preview_gif(sheet_path: str, char_name: str, state: str) -> str:
    """
    從 spritesheet 生成預覽 GIF
    
    Args:
        sheet_path: spritesheet 路徑
        char_name: 角色名稱
        state: 動畫狀態
    
    Returns:
        str: 輸出 GIF 路徑
    """
    sheet_path = Path(sheet_path)
    
    if not sheet_path.exists():
        print(f"[ERROR] 找不到 spritesheet：{sheet_path}")
        return ""
    
    try:
        sheet = Image.open(sheet_path).convert("RGBA")
    except Exception as e:
        print(f"[ERROR] 無法開啟圖片：{e}")
        return ""
    
    sheet_w, sheet_h = sheet.size
    
    # 推斷幀數
    if sheet_w % sheet_h == 0:
        frame_count = sheet_w // sheet_h
        frame_w = sheet_h
    else:
        for fc in [4, 6, 8, 12, 3]:
            if sheet_w % fc == 0:
                frame_count = fc
                frame_w = sheet_w // fc
                break
        else:
            frame_count = 1
            frame_w = sheet_w
    
    # 切割幀
    frames = []
    for i in range(frame_count):
        frame = sheet.crop((i * frame_w, 0, (i + 1) * frame_w, sheet_h))
        frames.append(frame)
    
    if not frames:
        print("[ERROR] 無法切割幀")
        return ""
    
    # 確保輸出目錄存在
    PREVIEW_DIR.mkdir(parents=True, exist_ok=True)
    
    output_path = PREVIEW_DIR / f"{char_name}_{state}.gif"
    fps = ANIMATION_FPS.get(state, 8)
    duration = int(1000 / fps)
    
    # 轉換為調色板模式（GIF 格式需要）
    palette_frames = []
    for frame in frames:
        # 建立白色背景（GIF 不支援真正的透明度）
        bg = Image.new("RGBA", frame.size, (200, 200, 200, 255))
        bg.paste(frame, mask=frame.split()[3])
        p_frame = bg.convert("P", palette=Image.ADAPTIVE, colors=255)
        palette_frames.append(p_frame)
    
    palette_frames[0].save(
        output_path,
        save_all=True,
        append_images=palette_frames[1:],
        duration=duration,
        loop=0,
        optimize=False
    )
    
    print(f"[OK] GIF 已生成：{output_path}（{frame_count} 幀，{fps} FPS）")
    return str(output_path)


def run_animation_audit() -> dict:
    """
    審查所有角色動畫，輸出報告
    
    Returns:
        dict: 完整審查結果
    """
    print("=" * 60)
    print("Animation Audit — 吉伊卡哇：像素大討伐")
    print(f"時間：{datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
    print("=" * 60)
    
    audit_results = {
        "timestamp": datetime.now().isoformat(),
        "characters": {},
        "summary": {
            "total_animations": 0,
            "passed": 0,
            "failed": 0,
            "missing": 0,
            "average_score": 0,
        }
    }
    
    all_scores = []
    
    for char in CHARACTERS:
        char_dir = SPRITES_DIR  # 所有角色 sprite 都在同一個 characters 目錄
        char_results = {}
        
        print(f"\n[角色] {char}")
        print("-" * 40)
        
        for state in ANIMATION_STATES:
            # 尋找 spritesheet（在 characters 目錄中）
            sheet_candidates = [
                char_dir / f"{char}_{state}.png",
                char_dir / f"{char}_{state}_sheet.png",
                PROJECT_ROOT / "client" / "chiikawa-pixel" / "assets" / "sprites" / "sheets" / f"{char}_{state}.png",
            ]
            
            sheet_path = None
            for candidate in sheet_candidates:
                if candidate.exists():
                    sheet_path = candidate
                    break
            
            if sheet_path is None:
                print(f"  [{state}] [WARN]  找不到 spritesheet")
                char_results[state] = {
                    "status": "missing",
                    "score": 0,
                    "issues": ["spritesheet 不存在"]
                }
                audit_results["summary"]["missing"] += 1
                audit_results["summary"]["total_animations"] += 1
                continue
            
            # 執行一致性檢查
            result = check_frame_consistency(str(sheet_path))
            
            if "error" in result:
                print(f"  [{state}] [FAIL] 錯誤：{result['error']}")
                char_results[state] = {
                    "status": "error",
                    "score": 0,
                    "issues": [result["error"]]
                }
                audit_results["summary"]["failed"] += 1
            elif result["passed"]:
                print(f"  [{state}] [OK] 通過（分數：{result['score']}/100，{result['frame_count']} 幀）")
                char_results[state] = {
                    "status": "passed",
                    "score": result["score"],
                    "frame_count": result["frame_count"],
                    "issues": result.get("issues", [])
                }
                audit_results["summary"]["passed"] += 1
                all_scores.append(result["score"])
            else:
                print(f"  [{state}] [FAIL] 未通過（分數：{result['score']}/100）")
                for issue in result.get("issues", []):
                    print(f"           問題：{issue}")
                char_results[state] = {
                    "status": "failed",
                    "score": result["score"],
                    "frame_count": result.get("frame_count", 0),
                    "issues": result.get("issues", [])
                }
                audit_results["summary"]["failed"] += 1
                all_scores.append(result["score"])
            
            audit_results["summary"]["total_animations"] += 1
        
        audit_results["characters"][char] = char_results
    
    # 計算平均分數
    if all_scores:
        audit_results["summary"]["average_score"] = round(sum(all_scores) / len(all_scores), 1)
    
    # 輸出摘要
    print("\n" + "=" * 60)
    print("審查摘要")
    print("=" * 60)
    s = audit_results["summary"]
    print(f"總動畫數：{s['total_animations']}")
    print(f"通過：{s['passed']} [OK]")
    print(f"未通過：{s['failed']} [FAIL]")
    print(f"缺失：{s['missing']} [WARN]")
    print(f"平均分數：{s['average_score']}/100")
    
    quality_label = "優秀" if s['average_score'] >= 90 else \
                    "良好" if s['average_score'] >= 85 else \
                    "需改善" if s['average_score'] >= 70 else "不合格"
    print(f"整體評級：{quality_label}")
    
    return audit_results


def generate_animation_report(results: dict) -> str:
    """
    輸出 Markdown 格式的動畫審查報告
    
    Args:
        results: run_animation_audit() 的輸出
    
    Returns:
        str: 報告檔案路徑
    """
    REPORTS_DIR.mkdir(parents=True, exist_ok=True)
    
    date_str = datetime.now().strftime("%Y-%m-%d")
    report_path = REPORTS_DIR / f"animation-audit-report-{date_str}.md"
    
    s = results["summary"]
    
    lines = [
        "# Animation Audit Report",
        "",
        f"**日期**：{date_str}",
        f"**執行者**：Animation Agent",
        f"**整體分數**：{s['average_score']}/100",
        "",
        "---",
        "",
        "## 摘要",
        "",
        f"| 項目 | 數量 |",
        f"|------|------|",
        f"| 總動畫數 | {s['total_animations']} |",
        f"| 通過 | {s['passed']} [OK] |",
        f"| 未通過 | {s['failed']} [FAIL] |",
        f"| 缺失 | {s['missing']} [WARN] |",
        f"| 平均分數 | {s['average_score']}/100 |",
        "",
        "---",
        "",
        "## 各角色詳細結果",
        "",
    ]
    
    for char, char_results in results.get("characters", {}).items():
        lines.append(f"### {char}")
        lines.append("")
        lines.append("| 動畫狀態 | 狀態 | 分數 | 幀數 | 問題 |")
        lines.append("|---------|------|------|------|------|")
        
        for state, data in char_results.items():
            status_icon = "[OK]" if data["status"] == "passed" else \
                          "[WARN]" if data["status"] == "missing" else "[FAIL]"
            score = data.get("score", 0)
            frame_count = data.get("frame_count", "-")
            issues = "; ".join(data.get("issues", [])) or "無"
            lines.append(f"| {state} | {status_icon} {data['status']} | {score} | {frame_count} | {issues} |")
        
        lines.append("")
    
    lines.extend([
        "---",
        "",
        "## 改善建議",
        "",
        "1. 缺失的 spritesheet 需要優先生成",
        "2. 分數低於 85 的動畫需要重新製作",
        "3. 有 jitter 問題的動畫需要重新對齊 anchor point",
        "4. 有 deformation 問題的動畫需要檢查生成參數",
        "",
        "---",
        "",
        f"*報告生成時間：{results.get('timestamp', date_str)}*",
    ])
    
    with open(report_path, "w", encoding="utf-8") as f:
        f.write("\n".join(lines))
    
    print(f"\n[OK] 報告已輸出：{report_path}")
    return str(report_path)


# ─── 主程式 ──────────────────────────────────────────────────────────────────

def main():
    parser = argparse.ArgumentParser(
        description="Animation Pipeline Tool — 吉伊卡哇：像素大討伐"
    )
    parser.add_argument("--audit", action="store_true", help="審查所有角色動畫")
    parser.add_argument("--gif", metavar="CHAR", help="生成指定角色的所有 GIF 預覽")
    parser.add_argument("--check", metavar="PATH", help="檢查單個 spritesheet")
    parser.add_argument("--report", action="store_true", help="輸出 Markdown 報告（搭配 --audit）")
    parser.add_argument("--state", metavar="STATE", help="指定動畫狀態（搭配 --gif）")
    
    args = parser.parse_args()
    
    if args.check:
        # 檢查單個 spritesheet
        print(f"檢查 spritesheet：{args.check}")
        result = check_frame_consistency(args.check)
        
        if "error" in result:
            print(f"[ERROR] {result['error']}")
            sys.exit(1)
        
        print(f"\n結果：{'[OK] 通過' if result['passed'] else '[FAIL] 未通過'}")
        print(f"分數：{result['score']}/100")
        print(f"幀數：{result['frame_count']}")
        print(f"尺寸：{result['sheet_size']}")
        print("\n各項檢查：")
        for check_name, check_data in result.get("checks", {}).items():
            icon = "[OK]" if check_data["passed"] else "[FAIL]"
            print(f"  {icon} {check_name}: {check_data['detail']}")
        
        if result.get("issues"):
            print("\n問題清單：")
            for issue in result["issues"]:
                print(f"  - {issue}")
        
        # 輸出 JSON 結果
        json_path = REPORTS_DIR / f"consistency-check-{Path(args.check).stem}.json"
        REPORTS_DIR.mkdir(parents=True, exist_ok=True)
        with open(json_path, "w", encoding="utf-8") as f:
            json.dump(result, f, ensure_ascii=False, indent=2)
        print(f"\n[OK] JSON 結果已儲存：{json_path}")
    
    elif args.gif:
        # 生成 GIF 預覽
        char_name = args.gif
        states_to_process = [args.state] if args.state else ANIMATION_STATES
        
        print(f"生成 {char_name} 的 GIF 預覽...")
        
        for state in states_to_process:
            # 在 characters 目錄中尋找
            sheet_candidates = [
                SPRITES_DIR / f"{char_name}_{state}.png",
                PROJECT_ROOT / "client" / "chiikawa-pixel" / "assets" / "sprites" / "sheets" / f"{char_name}_{state}.png",
            ]
            
            sheet_path = None
            for candidate in sheet_candidates:
                if candidate.exists():
                    sheet_path = candidate
                    break
            
            if sheet_path is None:
                print(f"  [{state}] [WARN]  找不到 spritesheet，跳過")
                continue
            
            generate_preview_gif(str(sheet_path), char_name, state)
    
    elif args.audit:
        # 審查所有動畫
        results = run_animation_audit()
        
        if args.report:
            generate_animation_report(results)
    
    else:
        parser.print_help()
        print("\n範例：")
        print("  py tools/animation_pipeline.py --audit")
        print("  py tools/animation_pipeline.py --audit --report")
        print("  py tools/animation_pipeline.py --gif chiikawa")
        print("  py tools/animation_pipeline.py --check client/chiikawa-pixel/assets/sprites/chars/chiikawa/chiikawa_idle.png")


if __name__ == "__main__":
    main()


