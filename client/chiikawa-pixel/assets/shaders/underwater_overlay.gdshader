// underwater_overlay.gdshader
// 全螢幕水下視覺效果（套用在最上層透明 ColorRect）
// 效果：輕微色差 + 水波扭曲 + 藍色調 + 深度霧氣
// 讓整個遊戲畫面有水下折射感，提升海底沉浸感

shader_type canvas_item;

// 色差強度（水下折射讓紅藍通道輕微分離）
uniform float aberration_strength : hint_range(0.0, 0.008, 0.0001) = 0.0025;

// 水波扭曲（非常輕微，不影響遊玩）
uniform float wave_strength : hint_range(0.0, 0.005, 0.0001) = 0.0015;
uniform float wave_speed : hint_range(0.1, 2.0) = 0.4;
uniform float wave_frequency : hint_range(1.0, 10.0) = 3.5;

// 藍色調（水下光線過濾）
uniform float blue_tint : hint_range(0.0, 0.15) = 0.04;

// 深度霧氣（畫面底部略暗，模擬水深）
uniform float depth_fog : hint_range(0.0, 0.25) = 0.08;

// 整體透明度（0 = 完全不影響，1 = 完整效果）
uniform float effect_alpha : hint_range(0.0, 1.0) = 1.0;

uniform sampler2D SCREEN_TEXTURE : hint_screen_texture, filter_linear_mipmap;

void fragment() {
	vec2 uv = SCREEN_UV;
	float t = TIME * wave_speed;

	// 水波扭曲（非常輕微的 sin 波，讓畫面有水下漂浮感）
	float wave_x = sin(uv.y * wave_frequency * 6.28 + t) * wave_strength;
	float wave_y = cos(uv.x * wave_frequency * 4.71 + t * 0.7) * wave_strength * 0.5;
	vec2 distorted_uv = uv + vec2(wave_x, wave_y);

	// 色差：紅通道向左偏移，藍通道向右偏移（水下折射效果）
	float ca = aberration_strength;
	vec4 red_sample   = texture(SCREEN_TEXTURE, distorted_uv + vec2(-ca, 0.0));
	vec4 green_sample = texture(SCREEN_TEXTURE, distorted_uv);
	vec4 blue_sample  = texture(SCREEN_TEXTURE, distorted_uv + vec2(ca, 0.0));

	// 合成色差結果
	vec4 color = vec4(red_sample.r, green_sample.g, blue_sample.b, green_sample.a);

	// 藍色調（水下光線過濾：紅色被水吸收，藍色保留）
	color.r *= (1.0 - blue_tint * 0.8);
	color.g *= (1.0 - blue_tint * 0.3);
	color.b = min(color.b + blue_tint * 0.5, 1.0);

	// 深度霧氣（畫面底部略暗，模擬水深增加）
	float depth = uv.y;
	float fog = depth * depth_fog;
	color.rgb *= (1.0 - fog);
	// 底部加藍色霧氣
	color.rgb = mix(color.rgb, vec3(0.05, 0.12, 0.28), fog * 0.4);

	// 輸出（透明，只影響顏色）
	COLOR = vec4(color.rgb, 0.0);  // alpha=0 讓 ColorRect 本身透明，只用 shader 修改螢幕
}
