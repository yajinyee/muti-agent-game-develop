// tier_glow.gdshader
// 目標物分層光暈 Shader — 依倍率等級顯示不同強度光暈
// 讓玩家在 1 秒內識別高價值目標
// visual-clarity-agent / hit-effect-agent 負責維護

shader_type canvas_item;

// 光暈顏色（由 GDScript 依倍率動態設定）
uniform vec4 glow_color : source_color = vec4(1.0, 0.85, 0.0, 1.0);
uniform float glow_radius : hint_range(1.0, 8.0) = 3.0;
uniform float glow_intensity : hint_range(0.0, 3.0) = 1.5;
uniform float pulse_speed : hint_range(0.0, 4.0) = 1.5;
uniform float pulse_amplitude : hint_range(0.0, 0.5) = 0.2;

void fragment() {
	vec4 tex = texture(TEXTURE, UV);
	vec2 texel = 1.0 / vec2(textureSize(TEXTURE, 0));

	// 高斯模糊近似：對周圍多個方向取樣累加
	float glow_alpha = 0.0;
	int samples = 8;
	for (int i = 0; i < samples; i++) {
		float angle = float(i) * TAU / float(samples);
		for (float r = 1.0; r <= glow_radius; r += 1.0) {
			vec2 offset = vec2(cos(angle), sin(angle)) * texel * r;
			glow_alpha += texture(TEXTURE, UV + offset).a;
		}
	}
	glow_alpha = clamp(glow_alpha / (float(samples) * glow_radius), 0.0, 1.0);

	// 脈動效果
	float pulse = 1.0 + pulse_amplitude * sin(TIME * pulse_speed);
	float final_glow = glow_alpha * glow_intensity * pulse;

	// 混合：原始像素 + 光暈
	vec4 glow_contribution = vec4(glow_color.rgb, glow_color.a * final_glow * (1.0 - tex.a));
	COLOR = tex + glow_contribution * (1.0 - tex.a);
	COLOR.a = max(tex.a, glow_alpha * glow_color.a * 0.6);
}
