#!/usr/bin/env python3
"""
qa_check_day336.py — DAY-336 QA 驗證腳本
驗證 HUD.gd 重構（HUDLuckySignals.gd 擴充 DAY-304~319）
"""
import os
import re

ROOT = r"d:\Kiro"
CLIENT = os.path.join(ROOT, "client", "chiikawa-pixel")
SCRIPTS_UI = os.path.join(CLIENT, "scripts", "ui")

errors = []
passed = 0

def check(condition, msg):
    global passed
    if condition:
        passed += 1
    else:
        errors.append(f"FAIL: {msg}")

def read(path):
    try:
        with open(path, encoding="utf-8") as f:
            return f.read()
    except:
        return ""

# ── 1. HUD.gd 行數應該減少 ──────────────────────────────────
hud = read(os.path.join(SCRIPTS_UI, "HUD.gd"))
hud_lines = len(hud.splitlines())
check(hud_lines < 2500, f"HUD.gd 行數應 < 2500（目前 {hud_lines}）")
check(hud_lines > 2000, f"HUD.gd 行數應 > 2000（目前 {hud_lines}，不應過度刪除）")

# ── 2. HUDLuckySignals.gd 應該存在且有足夠行數 ──────────────
lucky_sig = read(os.path.join(SCRIPTS_UI, "HUDLuckySignals.gd"))
lucky_lines = len(lucky_sig.splitlines())
check(lucky_lines > 1000, f"HUDLuckySignals.gd 應 > 1000 行（目前 {lucky_lines}）")

# ── 3. HUDLuckySignals 應包含 DAY-304~319 的訊號連接 ────────
day304_signals = [
    "lucky_electric_eel", "lucky_angler_fish", "lucky_black_hole",
    "lucky_bounty_hunter", "lucky_tsunami"
]
for sig in day304_signals:
    check(sig in lucky_sig, f"HUDLuckySignals 應包含 {sig} 連接")

day305_signals = [
    "lucky_dragon_wrath_v2", "lucky_humpback_whale", "lucky_legend_dragon",
    "lucky_guild_war", "lucky_quality_fish"
]
for sig in day305_signals:
    check(sig in lucky_sig, f"HUDLuckySignals 應包含 {sig} 連接")

day306_signals = [
    "lucky_tornado", "lucky_earthquake", "lucky_volcano",
    "lucky_cosmic_ray", "lucky_divine_dragon"
]
for sig in day306_signals:
    check(sig in lucky_sig, f"HUDLuckySignals 應包含 {sig} 連接")

day307_signals = [
    "lucky_quantum", "lucky_supernova", "lucky_infinite",
    "lucky_genesis", "lucky_rebirth"
]
for sig in day307_signals:
    check(sig in lucky_sig, f"HUDLuckySignals 應包含 {sig} 連接")

day308_signals = [
    "lucky_awakened_croc", "lucky_vampire_v2", "lucky_super_awaken",
    "lucky_giant_prize", "lucky_immortal_boss"
]
for sig in day308_signals:
    check(sig in lucky_sig, f"HUDLuckySignals 應包含 {sig} 連接")

day309_signals = [
    "lucky_ice_phoenix", "lucky_dragon_fury", "lucky_mult_cascade",
    "lucky_awaken_boss_v2", "lucky_ultimate_judgment"
]
for sig in day309_signals:
    check(sig in lucky_sig, f"HUDLuckySignals 應包含 {sig} 連接")

day310_signals = [
    "lucky_combo_burst", "lucky_time_bomb", "lucky_elemental_fusion",
    "lucky_treasure_hunter", "lucky_myth_awaken"
]
for sig in day310_signals:
    check(sig in lucky_sig, f"HUDLuckySignals 應包含 {sig} 連接")

day312_signals = [
    "lucky_star_portal", "lucky_dragon_soul", "lucky_spacetime_rift",
    "lucky_holy_judgment", "lucky_big_bang"
]
for sig in day312_signals:
    check(sig in lucky_sig, f"HUDLuckySignals 應包含 {sig} 連接")

day313_signals = ["lucky_jackpot_pool"]
for sig in day313_signals:
    check(sig in lucky_sig, f"HUDLuckySignals 應包含 {sig} 連接")

day314_signals = [
    "lucky_multiverse", "lucky_time_loop", "lucky_fate_wheel",
    "lucky_divine_realm", "lucky_final_power"
]
for sig in day314_signals:
    check(sig in lucky_sig, f"HUDLuckySignals 應包含 {sig} 連接")

day315_signals = [
    "lucky_mutation", "lucky_arctic_storm", "lucky_fisher_wild",
    "lucky_risk_level", "lucky_cosmic_pulse"
]
for sig in day315_signals:
    check(sig in lucky_sig, f"HUDLuckySignals 應包含 {sig} 連接")

day316_signals = [
    "lucky_mirror_universe", "lucky_gravity_field", "lucky_time_acceleration",
    "lucky_nebula_vortex", "lucky_cosmic_judgment"
]
for sig in day316_signals:
    check(sig in lucky_sig, f"HUDLuckySignals 應包含 {sig} 連接")

day317_signals = [
    "lucky_pvp_battle", "lucky_skill_chain", "lucky_global_explosion",
    "lucky_spacetime_fold", "lucky_cosmic_end"
]
for sig in day317_signals:
    check(sig in lucky_sig, f"HUDLuckySignals 應包含 {sig} 連接")

day318_signals = [
    "lucky_dragon_king", "lucky_eternal_cycle", "lucky_chaos_explosion",
    "lucky_divine_revival", "lucky_genesis_epoch"
]
for sig in day318_signals:
    check(sig in lucky_sig, f"HUDLuckySignals 應包含 {sig} 連接（fallback）")

day319_signals = [
    "lucky_energy_storm", "lucky_crystal_resonance", "lucky_fate_judgment",
    "lucky_time_reversal", "lucky_cosmic_singularity"
]
for sig in day319_signals:
    check(sig in lucky_sig, f"HUDLuckySignals 應包含 {sig} 連接（fallback）")

# ── 4. HUDLuckySignals 應包含 DAY-292~303 的訊號連接 ────────
day292_303_signals = [
    "lucky_chain_lightning", "lucky_crab_torpedo", "lucky_vortex",
    "lucky_golden_dragon", "lucky_thunder_lobster",
    "lucky_awakened_phoenix", "lucky_shockwave_bomb",
    "lucky_drill_torpedo", "lucky_time_freeze", "lucky_chain_explosion",
    "lucky_chain_long_king", "lucky_dragon_shotgun", "lucky_rocket_cannon",
    "lucky_deep_whirlpool", "lucky_vampire_mult",
    "lucky_mirror_fish", "lucky_golden_rain", "lucky_freeze_bomb",
    "lucky_thunder_storm", "lucky_lucky_wheel",
    "lucky_jackpot_fish", "lucky_coop_fish", "lucky_time_warp",
    "lucky_chain_meteor", "lucky_crash_fish"
]
for sig in day292_303_signals:
    check(sig in lucky_sig, f"HUDLuckySignals 應包含 {sig} 連接（DAY-292~303）")

# ── 5. HUD.gd 應包含 HUDLuckySignals 初始化 ─────────────────
check("HUDLuckySignals" in hud, "HUD.gd 應引用 HUDLuckySignals")
check("_lucky_signals" in hud, "HUD.gd 應有 _lucky_signals 變數")
check("_init_lucky_signals" in hud, "HUD.gd 應有 _init_lucky_signals 函數")
check("connect_all_lucky_signals" in hud, "HUD.gd 應呼叫 connect_all_lucky_signals")

# ── 6. HUDLuckySignals 應有委派輔助函數 ─────────────────────
check("_show_banner" in lucky_sig, "HUDLuckySignals 應有 _show_banner")
check("_show_event" in lucky_sig, "HUDLuckySignals 應有 _show_event")
check("_update_indicator" in lucky_sig, "HUDLuckySignals 應有 _update_indicator")
check("_hide_indicator" in lucky_sig, "HUDLuckySignals 應有 _hide_indicator")
check("_show_reward" in lucky_sig, "HUDLuckySignals 應有 _show_reward")
check("_on_lucky_fallback" in lucky_sig, "HUDLuckySignals 應有 _on_lucky_fallback")

# ── 7. Server 目標物數量確認 ─────────────────────────────────
tables_path = os.path.join(ROOT, "server", "internal", "data", "tables.go")
tables = read(tables_path)
check("T001" in tables, "Server tables.go 應包含 T001")
check("T253" in tables, "Server tables.go 應包含 T253")
check("B001" in tables, "Server tables.go 應包含 B001")

# ── 8. T001-T006 視覺升級確認（DAY-335）────────────────────
targets_dir = os.path.join(CLIENT, "assets", "sprites", "targets")
for t in ["T001_grass.png", "T002_bug_g.png", "T003_bug_r.png",
          "T004_bug_b.png", "T005_pudding.png", "T006_mushroom.png"]:
    path = os.path.join(targets_dir, t)
    check(os.path.exists(path), f"T001-T006 精靈圖應存在：{t}")

# ── 9. Shader 資產確認 ───────────────────────────────────────
shaders_dir = os.path.join(CLIENT, "assets", "shaders")
for shader in ["hit_flash.gdshader", "sprite_outline.gdshader", "tier_glow.gdshader"]:
    path = os.path.join(shaders_dir, shader)
    check(os.path.exists(shader) or os.path.exists(path), f"Shader 應存在：{shader}")

# ── 結果 ─────────────────────────────────────────────────────
total = passed + len(errors)
print(f"\n{'='*60}")
print(f"DAY-336 QA 結果：{passed}/{total} 通過")
print(f"{'='*60}")
if errors:
    print("\n失敗項目：")
    for e in errors:
        print(f"  {e}")
else:
    print("\n✅ 全部通過！")
