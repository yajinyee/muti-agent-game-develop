// shockwave_distortion.gdshader
// 衝擊波螢幕扭曲效果（帶色差）
// 來源：gameidea.org/2025/01/20/shockwave-distortion-shader-2d-space-wrap/
// 套用在全畫面 ColorRect 上，用 GDScript 控制 center 和 radius
// Content was rephrased for compliance with licensing restrictions

shader_type canvas_item;

// 扭曲參數
uniform float strength : hint_range(0.0, 0.15, 0.001) = 0.06;
uniform vec2 center = vec2(0.5, 0.5);
uniform float radius : hint_range(0.0, 1.0, 0.001) = 0.0;

// 視覺效果參數
uniform float aberration : hint_range(0.0, 1.0, 0.001) = 0.35;
uniform float width : hint_range(0.0, 0.15, 0.0001) = 0.05;
uniform float feather : hint_range(0.0, 1.0, 0.001) = 0.12;

uniform sampler2D SCREEN_TEXTURE : hint_screen_texture, filter_linear_mipmap;

void fragment() {
    vec2 screen_uv = SCREEN_UV;
    
    // 修正長寬比（讓效果是圓形而非橢圓）
    float aspect_ratio = SCREEN_PIXEL_SIZE.y / SCREEN_PIXEL_SIZE.x;
    vec2 aspect_uv = (screen_uv - vec2(0.0, 0.5)) / vec2(1.0, aspect_ratio) + vec2(0.0, 0.5);
    
    // 計算到中心的距離
    vec2 dist_from_center = aspect_uv - center;
    float dist_len = length(dist_from_center);
    
    // 建立環形遮罩（smoothstep 讓邊緣平滑）
    float outer_edge = smoothstep(radius - feather, radius, dist_len);
    float inner_edge = smoothstep(radius - width - feather, radius - width, dist_len);
    float mask = (1.0 - outer_edge) * inner_edge;
    
    // 計算扭曲方向和偏移量
    vec2 distort_dir = normalize(dist_from_center);
    vec2 distort_offset = distort_dir * strength * mask;
    
    // 套用扭曲
    vec2 distorted_uv = aspect_uv - distort_offset;
    vec2 final_uv = mix(screen_uv, distorted_uv, mask);
    
    // 色差效果（紅藍通道分離，增加衝擊感）
    vec2 color_sep = distort_offset * aberration * mask;
    vec4 red_ch = texture(SCREEN_TEXTURE, final_uv + color_sep);
    vec4 base_ch = texture(SCREEN_TEXTURE, final_uv);
    vec4 blue_ch = texture(SCREEN_TEXTURE, final_uv - color_sep);
    
    COLOR = vec4(red_ch.r, base_ch.g, blue_ch.b, 1.0);
}
