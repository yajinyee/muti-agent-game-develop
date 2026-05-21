// jackpot_meter.gdshader
// Jackpot 進度條發光 Shader（DAY-118）
// 讓進度條有動態發光效果，接近觸發時更強烈
shader_type canvas_item;

uniform float fill_ratio : hint_range(0.0, 1.0) = 0.5;  // 填充比例（0-1）
uniform vec4 bar_color : source_color = vec4(1.0, 0.85, 0.2, 1.0);  // 進度條顏色
uniform float glow_intensity : hint_range(0.0, 3.0) = 1.0;  // 發光強度
uniform float time_offset : hint_range(0.0, 6.28) = 0.0;    // 時間偏移（動畫用）

void fragment() {
    vec2 uv = UV;
    
    // 進度條填充
    float filled = step(uv.x, fill_ratio);
    
    // 基礎顏色
    vec4 base_color = bar_color * filled;
    
    // 高光效果（頂部 30% 更亮）
    float highlight = smoothstep(0.7, 1.0, uv.y) * 0.4 * filled;
    base_color.rgb += highlight;
    
    // 動態發光（sin 波動）
    float pulse = 0.5 + 0.5 * sin(TIME * 3.0 + time_offset);
    float glow = glow_intensity * pulse * filled;
    
    // 接近滿格時（>80%）加強發光
    float near_full = smoothstep(0.8, 1.0, fill_ratio);
    glow += near_full * 1.5 * pulse;
    
    base_color.rgb += bar_color.rgb * glow * 0.3;
    
    // 邊緣光暈（左右邊緣）
    float edge_glow = (1.0 - abs(uv.x - fill_ratio) * 20.0);
    edge_glow = max(0.0, edge_glow) * filled * glow_intensity;
    base_color.rgb += bar_color.rgb * edge_glow * 0.5;
    
    // 未填充部分：深色背景
    vec4 bg_color = vec4(0.05, 0.05, 0.1, 0.8);
    COLOR = mix(bg_color, base_color, filled) + vec4(base_color.rgb * edge_glow * 0.2, 0.0);
    COLOR.a = 1.0;
}
