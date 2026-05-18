// wobble.gdshader
// 像素搖晃效果 — 讓特殊目標物（T103 流星、T104 金草）有動態感
// 用法：套用到 Sprite2D，設定 wobble_amount 和 wobble_speed

shader_type canvas_item;

uniform float wobble_amount : hint_range(0.0, 0.05) = 0.015;
uniform float wobble_speed : hint_range(0.5, 10.0) = 3.0;
uniform float wobble_phase : hint_range(0.0, 6.28) = 0.0;

void fragment() {
    // 水平搖晃（sin 波形）
    float offset_x = sin(TIME * wobble_speed + UV.y * 8.0 + wobble_phase) * wobble_amount;
    // 輕微垂直搖晃（cos 波形，頻率不同）
    float offset_y = cos(TIME * wobble_speed * 0.7 + UV.x * 6.0 + wobble_phase) * wobble_amount * 0.5;
    
    vec2 wobble_uv = UV + vec2(offset_x, offset_y);
    // 邊界夾緊，防止採樣超出範圍
    wobble_uv = clamp(wobble_uv, vec2(0.001), vec2(0.999));
    
    COLOR = texture(TEXTURE, wobble_uv);
}
