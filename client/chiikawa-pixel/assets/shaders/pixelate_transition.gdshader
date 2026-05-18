// pixelate_transition.gdshader
// 像素化過場效果 — 場景切換時先像素化再還原
// 用法：套用到全畫面 ColorRect，控制 pixelate_amount (0.0=正常, 1.0=最像素化)
// 參考：ddrkirby.com shader-based transitions

shader_type canvas_item;

uniform float pixelate_amount : hint_range(0.0, 1.0) = 0.0;
uniform vec2 screen_size = vec2(1280.0, 720.0);

void fragment() {
    if (pixelate_amount < 0.001) {
        // 完全正常，直接採樣
        COLOR = texture(TEXTURE, UV);
        return;
    }
    
    // 像素化：把畫面分成大像素塊
    // pixelate_amount = 0 → 1px 塊（正常）
    // pixelate_amount = 1 → 64px 塊（最像素化）
    float block_size = mix(1.0, 64.0, pixelate_amount);
    
    // 計算當前像素所在的塊的中心 UV
    vec2 pixel_uv = UV * screen_size;
    vec2 block_uv = floor(pixel_uv / block_size) * block_size + block_size * 0.5;
    vec2 snapped_uv = block_uv / screen_size;
    
    // 夾緊 UV 範圍
    snapped_uv = clamp(snapped_uv, vec2(0.0), vec2(1.0));
    
    COLOR = texture(TEXTURE, snapped_uv);
}
