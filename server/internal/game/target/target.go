// Package target 管理遊戲中的目標物
package target

import (
	"math"
	"math/rand"
	"time"

	"digital-twin/server/internal/data"
)

// Quality 目標品質等級（DAY-070）
type Quality string

const (
	QualityNormal    Quality = "normal"    // 普通（70%）— 無加成
	QualityRare      Quality = "rare"      // 稀有（20%）— 倍率 +20%，藍色光暈
	QualityEpic      Quality = "epic"      // 史詩（8%）— 倍率 +50%，紫色光暈
	QualityLegendary Quality = "legendary" // 傳說（2%）— 倍率 +100%，金色光暈
)

// QualityMultiplierBonus 品質倍率加成
var QualityMultiplierBonus = map[Quality]float64{
	QualityNormal:    1.0,
	QualityRare:      1.2,
	QualityEpic:      1.5,
	QualityLegendary: 2.0,
}

// QualityColor 品質顏色（Client 端光暈顏色）
var QualityColor = map[Quality]string{
	QualityNormal:    "",        // 無光暈
	QualityRare:      "#4488FF", // 藍色
	QualityEpic:      "#AA44FF", // 紫色
	QualityLegendary: "#FFD700", // 金色
}

// rollQuality 隨機決定品質等級
// 普通 70%，稀有 20%，史詩 8%，傳說 2%
func rollQuality() Quality {
	r := rand.Intn(100)
	switch {
	case r < 70:
		return QualityNormal
	case r < 90:
		return QualityRare
	case r < 98:
		return QualityEpic
	default:
		return QualityLegendary
	}
}

// Target 遊戲中的目標物實例
type Target struct {
	InstanceID  string
	DefID       string
	Def         *data.TargetDef
	HP          int
	MaxHP       int
	HitCount    int // 已被命中次數（用於保底）
	Multiplier  float64
	X           float64
	Y           float64
	SpawnedAt   time.Time
	IsAlive     bool
	Phase       int // BOSS 用：0=Phase1, 1=Phase2
	IsFleeing   bool // 寶箱怪逃跑狀態
	// 品質等級（DAY-070）
	Quality     Quality
	QualityColor string
}

// NewTarget 建立新目標實例
func NewTarget(instanceID string, def *data.TargetDef, x, y float64) *Target {
	multiplier := def.MultiplierMin
	if def.MultiplierMax > def.MultiplierMin {
		// 流星等有倍率範圍的目標，依權重隨機
		if def.ID == "T103" {
			multiplier = rollMeteorMultiplier()
		} else {
			// 其他有範圍的目標，均勻隨機
			multiplier = def.MultiplierMin + rand.Float64()*(def.MultiplierMax-def.MultiplierMin)
		}
	}

	// 品質等級（DAY-070）：BOSS 和特殊目標固定 Legendary，普通目標隨機
	quality := QualityNormal
	if def.Type == data.TargetTypeBoss {
		quality = QualityLegendary
	} else if def.Type == data.TargetTypeSpecial {
		// 特殊目標有更高機率出現高品質
		r := rand.Intn(100)
		switch {
		case r < 40:
			quality = QualityNormal
		case r < 70:
			quality = QualityRare
		case r < 90:
			quality = QualityEpic
		default:
			quality = QualityLegendary
		}
	} else {
		quality = rollQuality()
	}

	// 套用品質倍率加成
	qualityBonus := QualityMultiplierBonus[quality]
	multiplier = multiplier * qualityBonus

	return &Target{
		InstanceID:   instanceID,
		DefID:        def.ID,
		Def:          def,
		HP:           def.HP,
		MaxHP:        def.HP,
		HitCount:     0,
		Multiplier:   multiplier,
		X:            x,
		Y:            y,
		SpawnedAt:    time.Now(),
		IsAlive:      true,
		Phase:        0,
		Quality:      quality,
		QualityColor: QualityColor[quality],
	}
}

// rollMeteorMultiplier 流星倍率加權隨機（規格書 26.3）
func rollMeteorMultiplier() float64 {
	total := 0
	for _, w := range data.MeteorMultiplierWeights {
		total += w.Weight
	}
	r := rand.Intn(total)
	cumulative := 0
	for _, w := range data.MeteorMultiplierWeights {
		cumulative += w.Weight
		if r < cumulative {
			return w.Multiplier
		}
	}
	return 20
}

// IsExpired 目標是否已超時
func (t *Target) IsExpired() bool {
	return time.Since(t.SpawnedAt).Seconds() > t.Def.Lifetime
}

// HPPercent 血量百分比
func (t *Target) HPPercent() float64 {
	if t.MaxHP == 0 {
		return 0
	}
	return float64(t.HP) / float64(t.MaxHP)
}

// RequiredHits 保底命中次數（規格書 25.4）
// 基礎目標：保底 = min(期望命中 × 3, Lifetime × maxFireRate × 0.8)
// 特殊目標：不設保底（純機率，避免高倍率保底導致 RTP 爆炸）
func (t *Target) RequiredHits(betCost int) int {
	// 特殊目標和 BOSS 不設保底
	if t.Def.Type == data.TargetTypeSpecial || t.Def.Type == data.TargetTypeBoss {
		return 99999
	}

	// 基礎目標：期望命中次數 × 3
	expected := t.Multiplier / data.BaseRTPFactor
	required := int(math.Ceil(expected * 3.0))

	// 上限：Lifetime 內最多攻擊次數 × 0.8（假設最快 fire_rate=3.0）
	maxHitsInLifetime := int(t.Def.Lifetime * 3.0 * 0.8)
	if maxHitsInLifetime < 2 {
		maxHitsInLifetime = 2
	}
	if required > maxHitsInLifetime {
		required = maxHitsInLifetime
	}
	if required < 2 {
		required = 2
	}
	return required
}

// KillChance 單次命中擊破機率（規格書 25.3 混合制）
func (t *Target) KillChance(betCost int, charKillMod float64) float64 {
	if betCost <= 0 {
		betCost = 1
	}
	// Kill Chance = BaseRTPFactor × BetCost ÷ (BetCost × Multiplier) × CharMod
	// 簡化：= BaseRTPFactor ÷ Multiplier × CharMod
	chance := data.BaseRTPFactor / t.Multiplier * charKillMod
	return math.Min(chance, 0.95) // 最高 95%
}

// TryKill 嘗試擊破（混合制：機率 + 保底）
// 回傳 (isKilled, damage)
func (t *Target) TryKill(betCost int, charKillMod float64) (bool, int) {
	if !t.IsAlive {
		return false, 0
	}

	t.HitCount++

	// 視覺 HP 傷害（讓玩家看到血量下降）
	damage := int(math.Max(1, float64(t.MaxHP)/float64(t.RequiredHits(betCost)+1)))
	t.HP = int(math.Max(0, float64(t.HP-damage)))

	// 保底：命中次數達到 RequiredHits 必定擊破
	if t.HitCount >= t.RequiredHits(betCost) {
		t.IsAlive = false
		t.HP = 0
		return true, damage
	}

	// 機率擊破
	chance := t.KillChance(betCost, charKillMod)
	if rand.Float64() < chance {
		t.IsAlive = false
		t.HP = 0
		return true, damage
	}

	return false, damage
}

// UpdateBossPhase 更新 BOSS 階段（規格書 28.4）
func (t *Target) UpdateBossPhase() bool {
	if t.Def.Type != data.TargetTypeBoss {
		return false
	}
	if t.Phase == 0 && t.HPPercent() <= 0.5 {
		t.Phase = 1
		return true // 觸發 Phase 2
	}
	return false
}

// SpawnSystem 目標生成系統
type SpawnSystem struct {
	rng *rand.Rand
}

// NewSpawnSystem 建立生成系統
func NewSpawnSystem() *SpawnSystem {
	return &SpawnSystem{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// SpawnWeights 依 Bet Level 取得生成權重（規格書 27.3）
type SpawnWeights struct {
	BasicRatio   float64
	SpecialRatio float64
	HighRatio    float64
}

func GetSpawnWeights(betLevel int) SpawnWeights {
	switch {
	case betLevel <= 3:
		return SpawnWeights{0.90, 0.09, 0.01}
	case betLevel <= 7:
		return SpawnWeights{0.82, 0.15, 0.03}
	default:
		return SpawnWeights{0.75, 0.20, 0.05}
	}
}

// PickTargetDef 依權重隨機選擇目標定義（規格書 27.3 三段動態難度）
// BasicRatio：基礎目標（T001-T006）
// SpecialRatio：一般特殊目標（T101-T103）
// HighRatio：高倍率特殊目標（T104 金色雜草 30x, T105 金幣魚 50x）
func (s *SpawnSystem) PickTargetDef(betLevel int, bonusSpecialRatio float64) *data.TargetDef {
	weights := GetSpawnWeights(betLevel)
	weights.SpecialRatio += bonusSpecialRatio

	r := s.rng.Float64()
	if r < weights.BasicRatio {
		return s.pickFromPool(getBasicPool())
	}
	// 在 Special 範圍內，依 HighRatio 決定是否選高倍率目標
	specialRange := weights.SpecialRatio + weights.HighRatio
	if specialRange > 0 {
		highThreshold := weights.BasicRatio + weights.HighRatio
		if r < highThreshold {
			highPool := getHighValuePool()
			if len(highPool) > 0 {
				return s.pickFromPool(highPool)
			}
		}
	}
	return s.pickFromPool(getSpecialPool())
}

func (s *SpawnSystem) pickFromPool(pool []*data.TargetDef) *data.TargetDef {
	total := 0
	for _, d := range pool {
		total += d.SpawnWeight
	}
	if total == 0 {
		return pool[0]
	}
	r := s.rng.Intn(total)
	cumulative := 0
	for _, d := range pool {
		cumulative += d.SpawnWeight
		if r < cumulative {
			return d
		}
	}
	return pool[0]
}

func getBasicPool() []*data.TargetDef {
	pool := make([]*data.TargetDef, 0)
	for _, d := range data.Targets {
		if d.Type == data.TargetTypeBasic {
			pool = append(pool, d)
		}
	}
	return pool
}

// getSpecialPool 一般特殊目標（T101-T103，不含高倍率 T104/T105）
func getSpecialPool() []*data.TargetDef {
	pool := make([]*data.TargetDef, 0)
	highValueIDs := map[string]bool{"T104": true, "T105": true}
	for _, d := range data.Targets {
		if d.Type == data.TargetTypeSpecial && !highValueIDs[d.ID] {
			pool = append(pool, d)
		}
	}
	return pool
}

// getHighValuePool 高倍率特殊目標（T104 金色雜草 30x, T105 金幣魚 50x）
// 規格書 27.3：LV8-10 玩家有 5% 機率遇到這類目標
func getHighValuePool() []*data.TargetDef {
	pool := make([]*data.TargetDef, 0)
	highValueIDs := map[string]bool{"T104": true, "T105": true}
	for _, d := range data.Targets {
		if d.Type == data.TargetTypeSpecial && highValueIDs[d.ID] {
			pool = append(pool, d)
		}
	}
	return pool
}
