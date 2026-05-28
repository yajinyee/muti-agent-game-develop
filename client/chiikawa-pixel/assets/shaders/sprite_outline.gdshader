// sprite_outline.gdshader
// 目標物輪廓 Shader — 提升視覺清晰度
// 用法：高倍率目標物（≥30x）套用此 Shader，讓玩家一眼識別高價值目標
// hit-effect-agent / visual-clarity-agent 負責維護

shader_type canvas_item;

uniform vec4 outline_color : source_color = vec4(1.0, 0.85, 0.0, 1.0); // 預設金色輪廓
uniform float outline_width : hint_range(0.5, 4.0) = 1.5;
uniform float pulse_speed : hint_range(0.0, 5.0) = 2.0;   // 脈動速度（0=靜態）
uniform float pulse_min : hint_range(0.0, 1.0) = 0.4;     // 脈動最低透明度
uniform bool enable_pulse = true;

void fragment() {
	vec4 tex = texture(TEXTURE, UV);

	// 計算輪廓：對周圍 8 個方向取樣，若有透明像素則為輪廓
	vec2 texel = outline_width / vec2(textureSize(TEXTURE, 0));
	float alpha_sum = 0.0;
	alpha_sum += texture(TEXTURE, UV + vec2(-texel.x, 0.0)).a;
	alpha_sum += texture(TEXTURE, UV + vec2( texel.x, 0.0)).a;
	alpha_sum += texture(TEXTURE, UV + vec2(0.0, -texel.y)).a;
	alpha_sum += texture(TEXTURE, UV + vec2(0.0,  texel.y)).a;
	alpha_sum += texture(TEXTURE, UV + vec2(-texel.x, -texel.y)).a;
	alpha_sum += texture(TEXTURE, UV + vec2( texel.x, -texel.y)).a;
	alpha_sum += texture(TEXTURE, UV + vec2(-texel.x,  texel.y)).a;
	alpha_sum += texture(TEXTURE, UV + vec2( texel.x,  texel.y)).a;

	// 輪廓條件：自身透明 + 周圍有不透明像素
	bool is_outline = (tex.a < 0.1) && (alpha_sum > 0.1);

	if (is_outline) {
		float pulse_alpha = 1.0;
		if (enable_pulse) {
			// 脈動效果：sin 波動透明度
			pulse_alpha = mix(pulse_min, 1.0, (sin(TIME * pulse_speed) + 1.0) * 0.5);
		}
		COLOR = vec4(outline_color.rgb, outline_color.a * pulse_alpha);
	} else {
		COLOR = tex;
	}
}
