// rainbow_glow.gdshader
// 彩虹光暈效果 — 大獎演出時讓角色/目標物有彩虹色輪廓
// 用法：套用到 Sprite2D，設定 glow_intensity 和 glow_speed

shader_type canvas_item;

uniform float glow_intensity : hint_range(0.0, 1.0) = 0.8;
uniform float glow_speed : hint_range(0.5, 5.0) = 2.0;
uniform float outline_width : hint_range(0.5, 3.0) = 1.5;

// HSV 轉 RGB（用於彩虹色）
vec3 hsv2rgb(vec3 c) {
    vec4 K = vec4(1.0, 2.0 / 3.0, 1.0 / 3.0, 3.0);
    vec3 p = abs(fract(c.xxx + K.xyz) * 6.0 - K.www);
    return c.z * mix(K.xxx, clamp(p - K.xxx, 0.0, 1.0), c.y);
}

void fragment() {
    vec4 original = texture(TEXTURE, UV);
    
    if (original.a < 0.1) {
        // 輪廓區域：彩虹色
        vec2 tex_size = vec2(textureSize(TEXTURE, 0));
        vec2 pixel = outline_width / tex_size;
        
        float neighbor_alpha = 0.0;
        neighbor_alpha += texture(TEXTURE, UV + vec2(pixel.x, 0.0)).a;
        neighbor_alpha += texture(TEXTURE, UV + vec2(-pixel.x, 0.0)).a;
        neighbor_alpha += texture(TEXTURE, UV + vec2(0.0, pixel.y)).a;
        neighbor_alpha += texture(TEXTURE, UV + vec2(0.0, -pixel.y)).a;
        
        if (neighbor_alpha > 0.1) {
            // 彩虹色（依時間和位置變化）
            float hue = fract(TIME * glow_speed * 0.1 + UV.x * 0.5 + UV.y * 0.3);
            vec3 rainbow = hsv2rgb(vec3(hue, 1.0, 1.0));
            COLOR = vec4(rainbow, glow_intensity);
        } else {
            COLOR = vec4(0.0);
        }
    } else {
        // 原始像素：輕微彩虹色調疊加
        float hue = fract(TIME * glow_speed * 0.1 + UV.x * 0.3);
        vec3 rainbow = hsv2rgb(vec3(hue, 0.6, 1.0));
        COLOR = vec4(mix(original.rgb, rainbow, glow_intensity * 0.3), original.a);
    }
}
