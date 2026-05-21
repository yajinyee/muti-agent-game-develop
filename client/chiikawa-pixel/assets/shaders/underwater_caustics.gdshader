// underwater_caustics.gdshader
// 海底焦散光效果（水面折射光斑）
// 套用在海底背景 Sprite2D 上，增加動態海底感

shader_type canvas_item;

uniform float time_scale: hint_range(0.1, 3.0) = 0.5;
uniform float caustic_intensity: hint_range(0.0, 0.3) = 0.12;
uniform float caustic_scale: hint_range(1.0, 10.0) = 5.0;
uniform vec4 caustic_color: source_color = vec4(0.55, 0.88, 1.0, 1.0);

// 簡單的 hash 函數（用於生成偽隨機值）
float hash(vec2 p) {
	p = fract(p * vec2(234.34, 435.345));
	p += dot(p, p + 34.23);
	return fract(p.x * p.y);
}

// 2D 噪聲（用於模擬焦散光斑）
float noise(vec2 p) {
	vec2 i = floor(p);
	vec2 f = fract(p);
	f = f * f * (3.0 - 2.0 * f);  // smoothstep
	
	float a = hash(i);
	float b = hash(i + vec2(1.0, 0.0));
	float c = hash(i + vec2(0.0, 1.0));
	float d = hash(i + vec2(1.0, 1.0));
	
	return mix(mix(a, b, f.x), mix(c, d, f.x), f.y);
}

// 焦散光斑（多層噪聲疊加）
float caustics(vec2 uv, float t) {
	float c = 0.0;
	// 三層不同速度和方向的噪聲疊加
	c += noise(uv * caustic_scale + vec2(t * 0.3, t * 0.2)) * 0.5;
	c += noise(uv * caustic_scale * 1.7 + vec2(-t * 0.2, t * 0.35)) * 0.3;
	c += noise(uv * caustic_scale * 2.3 + vec2(t * 0.15, -t * 0.25)) * 0.2;
	// 讓光斑更集中（提高對比度）
	c = pow(c, 2.0);
	return c;
}

void fragment() {
	vec4 base_color = texture(TEXTURE, UV);
	
	float t = TIME * time_scale;
	float c = caustics(UV, t);
	
	// 只在畫面上半部分顯示焦散（水面光線從上方射入）
	float depth_fade = clamp(1.0 - UV.y * 1.5, 0.0, 1.0);
	
	// 疊加焦散光效果
	vec3 final_color = base_color.rgb + caustic_color.rgb * c * caustic_intensity * depth_fade;
	
	COLOR = vec4(final_color, base_color.a);
}
