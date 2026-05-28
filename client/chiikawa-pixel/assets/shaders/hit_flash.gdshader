// hit_flash.gdshader
// 受擊閃白 Shader — 比純 modulate 更精確的受擊反饋
// 只讓有像素的地方閃白，透明區域不受影響
// hit-effect-agent 負責維護

shader_type canvas_item;

uniform float flash_intensity : hint_range(0.0, 1.0) = 0.0; // 由 GDScript 控制（0=正常，1=全白）
uniform vec4 flash_color : source_color = vec4(1.0, 1.0, 1.0, 1.0);

void fragment() {
	vec4 tex = texture(TEXTURE, UV);
	// 只對有像素的地方套用閃白，保留透明度
	vec3 flashed = mix(tex.rgb, flash_color.rgb, flash_intensity);
	COLOR = vec4(flashed, tex.a);
}
