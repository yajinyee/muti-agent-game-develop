// hit_flash.gdshader
// 目標物受擊閃白效果
// 用法：將此 shader 套用到 Sprite2D，設定 flash_amount (0.0-1.0)

shader_type canvas_item;

uniform float flash_amount : hint_range(0.0, 1.0) = 0.0;
uniform vec4 flash_color : source_color = vec4(1.0, 1.0, 1.0, 1.0);

void fragment() {
    vec4 original = texture(TEXTURE, UV);
    // 只對非透明像素套用閃白
    if (original.a > 0.01) {
        COLOR = mix(original, vec4(flash_color.rgb, original.a), flash_amount);
    } else {
        COLOR = original;
    }
}
