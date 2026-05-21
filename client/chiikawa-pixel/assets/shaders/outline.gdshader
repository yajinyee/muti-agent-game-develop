// outline.gdshader
// 像素輪廓效果 — 讓目標物有清晰的像素輪廓，提升辨識度
// 用法：套用到 Sprite2D，設定 outline_color 和 outline_width
// 參考：Godot Asset Library #2829 + cyanilux.com sprite outline breakdown

shader_type canvas_item;

uniform vec4 outline_color : source_color = vec4(0.0, 0.0, 0.0, 1.0);
uniform float outline_width : hint_range(0.5, 4.0) = 1.0;
uniform float outline_alpha : hint_range(0.0, 1.0) = 1.0;

void fragment() {
    vec4 original = texture(TEXTURE, UV);
    
    // 如果原始像素是透明的，檢查周圍是否有非透明像素（輪廓判定）
    if (original.a < 0.1) {
        vec2 tex_size = vec2(textureSize(TEXTURE, 0));
        vec2 pixel = outline_width / tex_size;
        
        // 8方向採樣（像素輪廓）
        float neighbor_alpha = 0.0;
        neighbor_alpha += texture(TEXTURE, UV + vec2(pixel.x, 0.0)).a;
        neighbor_alpha += texture(TEXTURE, UV + vec2(-pixel.x, 0.0)).a;
        neighbor_alpha += texture(TEXTURE, UV + vec2(0.0, pixel.y)).a;
        neighbor_alpha += texture(TEXTURE, UV + vec2(0.0, -pixel.y)).a;
        neighbor_alpha += texture(TEXTURE, UV + vec2(pixel.x, pixel.y)).a;
        neighbor_alpha += texture(TEXTURE, UV + vec2(-pixel.x, pixel.y)).a;
        neighbor_alpha += texture(TEXTURE, UV + vec2(pixel.x, -pixel.y)).a;
        neighbor_alpha += texture(TEXTURE, UV + vec2(-pixel.x, -pixel.y)).a;
        
        if (neighbor_alpha > 0.1) {
            // 輪廓像素
            COLOR = vec4(outline_color.rgb, outline_color.a * outline_alpha);
        } else {
            COLOR = vec4(0.0);
        }
    } else {
        COLOR = original;
    }
}
