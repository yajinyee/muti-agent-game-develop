// water_surface.gdshader
// 水面波紋效果（套用在頂部水面 ColorRect 上）
// 模擬從水面往下看的波紋折射感

shader_type canvas_item;

uniform float wave_speed: hint_range(0.1, 3.0) = 0.8;
uniform float wave_amplitude: hint_range(0.001, 0.05) = 0.012;
uniform float wave_frequency: hint_range(1.0, 20.0) = 8.0;
uniform float shimmer_speed: hint_range(0.1, 5.0) = 1.5;
uniform vec4 water_color: source_color = vec4(0.15, 0.45, 0.75, 0.55);
uniform vec4 foam_color: source_color = vec4(0.75, 0.92, 1.0, 0.35);

// 簡單 hash
float hash(float n) {
	return fract(sin(n) * 43758.5453123);
}

// 1D 噪聲（用於波紋形狀）
float noise1d(float x) {
	float i = floor(x);
	float f = fract(x);
	f = f * f * (3.0 - 2.0 * f);
	return mix(hash(i), hash(i + 1.0), f);
}

void fragment() {
	vec2 uv = UV;
	float t = TIME * wave_speed;

	// 多層波紋疊加（不同頻率和速度）
	float wave1 = sin(uv.x * wave_frequency + t) * wave_amplitude;
	float wave2 = sin(uv.x * wave_frequency * 1.7 - t * 0.8) * wave_amplitude * 0.6;
	float wave3 = sin(uv.x * wave_frequency * 0.5 + t * 1.3) * wave_amplitude * 0.4;
	float wave_offset = wave1 + wave2 + wave3;

	// 水面亮度（模擬光線在水面的閃爍）
	float shimmer = noise1d(uv.x * 12.0 + TIME * shimmer_speed) * 0.5 + 0.5;
	shimmer += noise1d(uv.x * 7.0 - TIME * shimmer_speed * 0.7) * 0.3;
	shimmer = clamp(shimmer / 1.3, 0.0, 1.0);

	// 波峰泡沫（波紋頂部的白色泡沫線）
	float wave_crest = abs(wave1 + wave2);
	float foam = smoothstep(0.006, 0.012, wave_crest);

	// 深度漸變（頂部較亮，往下漸暗）
	float depth = 1.0 - uv.y;
	vec4 base = mix(water_color, water_color * 1.3, depth * shimmer);

	// 疊加泡沫
	vec4 final_color = mix(base, foam_color, foam * 0.7);

	// 整體透明度：頂部不透明，往下漸透明
	final_color.a = water_color.a * (1.0 - uv.y * 0.6);

	COLOR = final_color;
}
