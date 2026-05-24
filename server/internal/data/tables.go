// Package data 定義遊戲所有靜態資料表（來自規格書）
package data

// TargetType 目標類型
type TargetType string

const (
	TargetTypeBasic   TargetType = "basic"
	TargetTypeSpecial TargetType = "special"
	TargetTypeBoss    TargetType = "boss"
	TargetTypeBonus   TargetType = "bonus"
)

// TargetDef 目標物定義（對應規格書 Target Table）
type TargetDef struct {
	ID              string
	Name            string
	Type            TargetType
	MultiplierMin   float64
	MultiplierMax   float64
	HP              int
	SpawnWeight     int
	Speed           float64 // pixels/sec
	Lifetime        float64 // seconds
	LaborGain       int
	DifficultyFactor float64
	SpecialBehavior string
}

// CharacterDef 角色定義
type CharacterDef struct {
	ID               string
	Name             string
	BetLevelMin      int
	BetLevelMax      int
	AttackColor      string
	KillModifier     float64
	FireRateModifier float64
	LaborModifier    float64
	VoiceText        string
}

// BetDef 投注等級定義
type BetDef struct {
	Level           int
	CharacterID     string
	BetCost         int
	AttackPower     int
	FireRate        float64 // shots/sec
	ProjectileSpeed float64 // pixels/sec
}

// BonusTargetDef Bonus Game 目標定義
type BonusTargetDef struct {
	ID           string
	Name         string
	ClickScore   int
	SpawnWeight  int
	SpecialEffect string
}

// ---- 靜態資料表 ----

// Targets 所有目標物（規格書 26.1）
// DifficultyFactor 修正說明（2026-05-12）：
//   正確公式：required_hits = ceil(multiplier / bet_cost × DifficultyFactor)
//   要讓保底在期望命中次數（= multiplier / RTP）的 1.5 倍觸發：
//   DifficultyFactor = 1.5 × bet_cost / RTP ≈ 1.5 × bet_cost / 0.94
//   對 bet_cost=10：DifficultyFactor ≈ 16
//   這樣 T001(2x,bet=10)：required = ceil(2/10×16) = ceil(3.2) = 4 次保底
//   期望命中 ≈ 1/0.47 = 2.1 次，保底 4 次，RTP ≈ 94% ✓
var Targets = map[string]*TargetDef{
	"T001": {ID: "T001", Name: "像素雜草", Type: TargetTypeBasic, MultiplierMin: 2, MultiplierMax: 2, HP: 3, SpawnWeight: 180, Speed: 0, Lifetime: 20, LaborGain: 1, DifficultyFactor: 16.0, SpecialBehavior: "static_sway"},
	"T002": {ID: "T002", Name: "綠色小蟲", Type: TargetTypeBasic, MultiplierMin: 3, MultiplierMax: 3, HP: 5, SpawnWeight: 160, Speed: 40, Lifetime: 18, LaborGain: 1, DifficultyFactor: 16.0, SpecialBehavior: "linear"},
	"T003": {ID: "T003", Name: "紅色小蟲", Type: TargetTypeBasic, MultiplierMin: 5, MultiplierMax: 5, HP: 8, SpawnWeight: 130, Speed: 55, Lifetime: 16, LaborGain: 1, DifficultyFactor: 16.0, SpecialBehavior: "jump"},
	"T004": {ID: "T004", Name: "藍色小蟲", Type: TargetTypeBasic, MultiplierMin: 6, MultiplierMax: 6, HP: 10, SpawnWeight: 110, Speed: 65, Lifetime: 15, LaborGain: 2, DifficultyFactor: 16.0, SpecialBehavior: "curve"},
	"T005": {ID: "T005", Name: "會走路的布丁", Type: TargetTypeBasic, MultiplierMin: 8, MultiplierMax: 8, HP: 16, SpawnWeight: 90, Speed: 35, Lifetime: 20, LaborGain: 2, DifficultyFactor: 16.0, SpecialBehavior: "sway"},
	"T006": {ID: "T006", Name: "巨大蘑菇", Type: TargetTypeBasic, MultiplierMin: 10, MultiplierMax: 10, HP: 22, SpawnWeight: 70, Speed: 25, Lifetime: 22, LaborGain: 3, DifficultyFactor: 16.0, SpecialBehavior: "sink"},
	"T101": {ID: "T101", Name: "擬態型怪物", Type: TargetTypeSpecial, MultiplierMin: 15, MultiplierMax: 30, HP: 35, SpawnWeight: 35, Speed: 50, Lifetime: 14, LaborGain: 5, DifficultyFactor: 16.0, SpecialBehavior: "mimic"},
	"T102": {ID: "T102", Name: "寶箱怪", Type: TargetTypeSpecial, MultiplierMin: 25, MultiplierMax: 25, HP: 55, SpawnWeight: 22, Speed: 70, Lifetime: 10, LaborGain: 6, DifficultyFactor: 16.0, SpecialBehavior: "flee"},
	"T103": {ID: "T103", Name: "流星", Type: TargetTypeSpecial, MultiplierMin: 20, MultiplierMax: 50, HP: 20, SpawnWeight: 18, Speed: 220, Lifetime: 4, LaborGain: 5, DifficultyFactor: 16.0, SpecialBehavior: "meteor"},
	"T104": {ID: "T104", Name: "金色雜草", Type: TargetTypeSpecial, MultiplierMin: 30, MultiplierMax: 30, HP: 45, SpawnWeight: 12, Speed: 0, Lifetime: 8, LaborGain: 15, DifficultyFactor: 16.0, SpecialBehavior: "static"},
	"T105": {ID: "T105", Name: "巨大金幣魚", Type: TargetTypeSpecial, MultiplierMin: 50, MultiplierMax: 50, HP: 90, SpawnWeight: 8, Speed: 80, Lifetime: 8, LaborGain: 10, DifficultyFactor: 16.0, SpecialBehavior: "coin_rain"},
	// T106 鑽頭龍蝦（DAY-142）— 業界依據：Royal Fishing JILI 2026「Drill Bit Lobster (80X) — penetrating drill through multiple fish, self-detonates at end of trajectory」
	// 擊破後觸發穿透鑽頭，沿水平方向穿透所有目標，到達邊緣後爆炸，連帶擊破爆炸範圍內目標
	"T106": {ID: "T106", Name: "鑽頭龍蝦", Type: TargetTypeSpecial, MultiplierMin: 60, MultiplierMax: 80, HP: 60, SpawnWeight: 6, Speed: 45, Lifetime: 10, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "drill_lobster"},
	// T107 炸彈蟹（DAY-143）— 業界依據：royal-fishing.uk 2026「Worth 70x, this explosive crustacean triggers multiple large-scale detonations.
	// Each bomb creates expanding capture zones for massive multi-target eliminations.」
	// 擊破後觸發 3 波爆炸，每波爆炸半徑 150px，每波間隔 400ms，連帶擊破爆炸範圍內所有目標
	// T107 炸彈蟹（DAY-143）— 業界依據：royal-fishing.uk 2026「Worth 70x, this explosive crustacean triggers multiple large-scale detonations.
	// Each bomb creates expanding capture zones for massive multi-target eliminations.」
	// 擊破後觸發 3 波爆炸，每波爆炸半徑 150px，每波間隔 400ms，連帶擊破爆炸範圍內所有目標
	"T107": {ID: "T107", Name: "炸彈蟹", Type: TargetTypeSpecial, MultiplierMin: 60, MultiplierMax: 70, HP: 70, SpawnWeight: 5, Speed: 35, Lifetime: 12, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "bomb_crab"},
	// T108 巨型章魚（DAY-144）— 業界依據：JILI Mega Fishing「Mega Octopus Wheel – Defeat that giant octopus and enter
	// the bonus wheel round where you have a chance to win massive guaranteed prizes up to 950x.」
	// 擊破後觸發個人轉盤（8格：50x-950x），玩家點擊停止，結果預先決定（公平性保證）
	"T108": {ID: "T108", Name: "巨型章魚", Type: TargetTypeSpecial, MultiplierMin: 80, MultiplierMax: 120, HP: 120, SpawnWeight: 3, Speed: 30, Lifetime: 15, LaborGain: 15, DifficultyFactor: 16.0, SpecialBehavior: "mega_octopus"},
	// T109 巨型鮟鱇魚（DAY-145）— 業界依據：jiligames.com 2026「Giant Anglerfish can shoot electricity to open treasure chests」
	// 擊破後觸發電擊，電流傳導到附近的寶箱目標（T102），強制開啟寶箱獲得額外獎勵
	"T109": {ID: "T109", Name: "巨型鮟鱇魚", Type: TargetTypeSpecial, MultiplierMin: 70, MultiplierMax: 90, HP: 90, SpawnWeight: 4, Speed: 25, Lifetime: 14, LaborGain: 14, DifficultyFactor: 16.0, SpecialBehavior: "anglerfish_shock"},
	// T110 巨型鹹水鱷魚（DAY-146）— 業界依據：jiligames.com 2026「giant crocodiles awaken to hunt fish on the fish farm to accumulate big prizes!」
	// + megafishinggame.top「Giant Saltwater Crocodile」
	// 擊破後觸發「鱷魚獵魚」模式：鱷魚在 8 秒內自動獵殺場上的普通目標（T001-T006），累積獎勵給觸發玩家
	"T110": {ID: "T110", Name: "巨型鹹水鱷魚", Type: TargetTypeSpecial, MultiplierMin: 100, MultiplierMax: 150, HP: 150, SpawnWeight: 2, Speed: 20, Lifetime: 18, LaborGain: 18, DifficultyFactor: 16.0, SpecialBehavior: "crocodile_hunt"},
	// T111 夢幻巨型獎勵魚（DAY-147）— 業界依據：jiligames.com 2026「The dreamy Giant Prize Fish lets you easily win great prizes, with the chance for 5x multipliers」
	// 擊破後觸發「夢幻獎勵模式」：觸發玩家在 10 秒內所有擊破獎勵 ×5，讓玩家感受到「夢幻大獎」的爽感
	// 設計：低 HP（容易擊破）+ 中等倍率（40-60x）+ 觸發後 10 秒 5x 加成，是「容易觸發的短期爆發」機制
	"T111": {ID: "T111", Name: "夢幻巨型獎勵魚", Type: TargetTypeSpecial, MultiplierMin: 40, MultiplierMax: 60, HP: 80, SpawnWeight: 4, Speed: 35, Lifetime: 12, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "giant_prize_fish"},
	// T112 千龍王（DAY-148）— 業界依據：Royal Fishing JILI 2026「ChainLong King — capture this golden dragon to trigger
	// the dual-ring roulette. The ChainLong King itself can award up to 1000X mega wins.」
	// 擊破後觸發「千龍王強化輪盤」：內環（5x-50x）× 外環（2x-20x）= 最高 1000x
	// 設計：超高倍率（150-1000x）+ 超高 HP（300）+ 極低生成權重（1）= 終極稀有目標
	// 千龍王是全遊戲最高倍率目標，擊破後觸發專屬強化輪盤，最高 1000x 是業界最高水準
	"T112": {ID: "T112", Name: "千龍王", Type: TargetTypeSpecial, MultiplierMin: 150, MultiplierMax: 1000, HP: 300, SpawnWeight: 1, Speed: 15, Lifetime: 20, LaborGain: 30, DifficultyFactor: 16.0, SpecialBehavior: "chainlong_king"},
	// T113 黃金水母（DAY-149）— 業界依據：Ocean King 3 2026「Electric Jellyfish chain shocks across multiple targets.
	// Devastating against clustered schools.」— 擊破後觸發「全場電擊」，對畫面上所有目標發動電擊
	// 比閃電鰻（T103，200px 範圍跳躍 5 次）更強：全場範圍，最多 8 個目標，40% 擊破機率
	// 設計：中等倍率（60-80x）+ 中等 HP（80）+ 低生成權重（3）= 稀有但可遇到的強力目標
	"T113": {ID: "T113", Name: "黃金水母", Type: TargetTypeSpecial, MultiplierMin: 60, MultiplierMax: 80, HP: 80, SpawnWeight: 3, Speed: 30, Lifetime: 12, LaborGain: 14, DifficultyFactor: 16.0, SpecialBehavior: "golden_jellyfish"},
	// T114 雷霆龍蝦（DAY-150）— 業界依據：royalfishingsite.com 2026「Thunderbolt Lobster feature —
	// 15 seconds of free play followed by automatic shooting」
	// 擊破後觸發「免費射擊模式」：15 秒內所有子彈不扣費，Server 自動幫玩家以當前 betLevel 射擊
	// 每 0.5 秒自動射擊一次，優先選高倍率目標，讓玩家感受到「免費狂射」的爽感
	// 設計：中等倍率（50-70x）+ 中等 HP（70）+ 中低生成權重（5）= 常見但有驚喜的特殊目標
	"T114": {ID: "T114", Name: "雷霆龍蝦", Type: TargetTypeSpecial, MultiplierMin: 50, MultiplierMax: 70, HP: 70, SpawnWeight: 5, Speed: 40, Lifetime: 12, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "thunderbolt_lobster"},
	// T115 彩虹鳳凰（DAY-151）— 業界依據：royal-fishing.co.uk 2026「Multicoloured phoenix (blue, pink, purple, orange)
	// with magical aura. Awaken Boss with 30x basic multiplier. Power Up attack delivers 6x-10x boost for rewards up to 300 times bet.」
	// 擊破後觸發「Power Up 模式」：玩家在 8 秒內所有攻擊獲得隨機 6x-10x 倍率加成，最高 300x
	// 設計：中高倍率（80-120x）+ 中等 HP（100）+ 低生成權重（3）= 稀有但令人興奮的特殊目標
	"T115": {ID: "T115", Name: "彩虹鳳凰", Type: TargetTypeSpecial, MultiplierMin: 80, MultiplierMax: 120, HP: 100, SpawnWeight: 3, Speed: 35, Lifetime: 14, LaborGain: 16, DifficultyFactor: 16.0, SpecialBehavior: "rainbow_phoenix"},
	// T116 吸血鬼（DAY-152）— 業界依據：jiligames.com 2026「The explicit multiplier of vampires increases the more you fight,
	// and there is a chance that you can enter the multiplier mode, up to X5.」
	// 每次被命中倍率增加：5次→×2.0覺醒，10次→×3.5狂暴，15次→×5.0血月（全服廣播）
	// 設計：基礎倍率（20-30x）+ 高 HP（120）+ 中等生成權重（6）= 需要多人合力打的成長型目標
	"T116": {ID: "T116", Name: "吸血鬼", Type: TargetTypeSpecial, MultiplierMin: 20, MultiplierMax: 30, HP: 120, SpawnWeight: 6, Speed: 45, Lifetime: 18, LaborGain: 10, DifficultyFactor: 16.0, SpecialBehavior: "vampire_grow"},
	// T117 水晶龍（DAY-153）— 業界依據：jiligames.com JILI Flying Dragon 2026「collect crystals to get the grand prize!
	// Kill the Underworld Dragon and win the prize!」
	// 擊破後掉落 5 個水晶碎片，全服玩家共同收集水晶（目標 50 個），達到目標後觸發「地獄龍大獎」
	// 按貢獻比例分配獎勵（最高 200x betLevel），全服廣播，增加社交合作感
	// 設計：中等倍率（30-50x）+ 中等 HP（80）+ 中等生成權重（5）= 常見但有合作感的特殊目標
	"T117": {ID: "T117", Name: "水晶龍", Type: TargetTypeSpecial, MultiplierMin: 30, MultiplierMax: 50, HP: 80, SpawnWeight: 5, Speed: 38, Lifetime: 15, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "crystal_dragon"},
	// T118 皇家閃電鰻（DAY-156）— 業界依據：royal-fishing.co.uk 2026「Creates chain lightning that shocks nearby fish
	// consecutively until targeting turns off. Devastating against clustered schools.」
	// 擊破後觸發持續連鎖電擊，每 200ms 跳一次（最多 15 跳），每跳 300px 範圍，60% 擊破機率
	// 比 T103 閃電鰻（一次性 5 跳/200px/50%）更強，是升級版連鎖電擊
	// 設計：中高倍率（40-60x）+ 中等 HP（90）+ 中等生成權重（4）= 稀有但連鎖效果強的特殊目標
	"T118": {ID: "T118", Name: "皇家閃電鰻", Type: TargetTypeSpecial, MultiplierMin: 40, MultiplierMax: 60, HP: 90, SpawnWeight: 4, Speed: 50, Lifetime: 16, LaborGain: 14, DifficultyFactor: 16.0, SpecialBehavior: "royal_chain_lightning"},
	// T119 黃金海龜（DAY-159）— 業界依據：Ocean King 系列「Time Stop」機制
	// 擊破後觸發「全場時間停止」8 秒，所有目標物暫停移動，玩家可以輕鬆瞄準
	// 是「輔助型特殊目標」，不直接給高獎勵，但讓玩家在 8 秒內大量擊破其他目標
	// 設計：中等倍率（30-50x）+ 中等 HP（60）+ 中等生成權重（5）= 常見但有輔助效果的特殊目標
	"T119": {ID: "T119", Name: "黃金海龜", Type: TargetTypeSpecial, MultiplierMin: 30, MultiplierMax: 50, HP: 60, SpawnWeight: 5, Speed: 20, Lifetime: 15, LaborGain: 10, DifficultyFactor: 16.0, SpecialBehavior: "golden_turtle"},
	// T120 幸運星魚（DAY-160）— 業界依據：捕魚機業界標準「倍率爆發」機制
	// 擊破後觸發「全場倍率翻倍」10 秒，所有目標物的獎勵倍率 ×2
	// 是「爆發型特殊目標」，讓玩家在 10 秒內所有擊破獎勵翻倍，製造「大豐收」的爽感
	// 設計：中等倍率（40-60x）+ 中等 HP（70）+ 低生成權重（4）= 稀有但爆發效果強的特殊目標
	"T120": {ID: "T120", Name: "幸運星魚", Type: TargetTypeSpecial, MultiplierMin: 40, MultiplierMax: 60, HP: 70, SpawnWeight: 4, Speed: 45, Lifetime: 14, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "lucky_star_fish"},
	// T121 黃金鯊魚（DAY-161）— 業界依據：King of Ocean 2026「sharks climb into x50-x300 zone」
	// + 捕魚機業界「rage/berserk mode」機制 — 擊破後觸發「全服狂暴模式」
	// 全場所有目標物獎勵倍率 ×1.5，持續 12 秒，全服廣播
	// 設計：全服共享（不是個人），任何玩家擊破都讓全服受益，製造「全場爆發」的社交爽感
	// 與幸運星魚（個人 ×2）不同：黃金鯊魚是全服 ×1.5，社交性更強
	"T121": {ID: "T121", Name: "黃金鯊魚", Type: TargetTypeSpecial, MultiplierMin: 50, MultiplierMax: 80, HP: 100, SpawnWeight: 3, Speed: 55, Lifetime: 16, LaborGain: 15, DifficultyFactor: 16.0, SpecialBehavior: "golden_shark_berserk"},
	// T122 金幣魚王（DAY-162）— 業界依據：King of Ocean 2026（Galaxsys）「Money Fish trigger instant payouts」
	// 擊破後立即給予玩家一筆即時獎勵（betLevel × 20-50 隨機），不走正常 kill 倍率計算
	// 是「保底即時獎勵」型特殊目標，讓玩家在任何 betLevel 都能獲得有感的即時金幣
	// 設計：中等倍率（30-50x）+ 中等 HP（60）+ 中等生成權重（5）= 常見但有即時爆金幣效果
	"T122": {ID: "T122", Name: "金幣魚王", Type: TargetTypeSpecial, MultiplierMin: 30, MultiplierMax: 50, HP: 60, SpawnWeight: 5, Speed: 50, Lifetime: 14, LaborGain: 10, DifficultyFactor: 16.0, SpecialBehavior: "money_fish_instant"},
	// T123 船長魚（DAY-163）— 業界依據：King of Ocean 2026（Galaxsys）「Captain Fish trigger bonus rounds」
	// 擊破後觸發「全服競速模式」，30 秒內全服玩家競爭擊破最多目標
	// 第一名獲得 betLevel × 30，第二名 × 15，第三名 × 8
	// 設計：競技型社交機制（與水晶龍的合作型形成對比），製造「全服競爭」的緊張爽感
	"T123": {ID: "T123", Name: "船長魚", Type: TargetTypeSpecial, MultiplierMin: 40, MultiplierMax: 60, HP: 80, SpawnWeight: 3, Speed: 40, Lifetime: 16, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "captain_fish_race"},
	// T124 深淵巨鯨（DAY-164）— 業界依據：Fishing Frenzy Chapter 3 2026「Boss Fish as endgame content」
	// + Ocean King 2026「Abyss Whale — massive HP boss requiring full server cooperation」
	// 超高 HP（500），需要全服玩家合力攻擊才能擊破
	// 擊破後觸發「深淵寶藏」，按傷害貢獻比例分配獎勵（最高 500x betLevel）
	// 設計：合作型終局內容，與船長魚（競技型）形成對比，製造「全服合力打 Boss」的緊張爽感
	// 注意：HP=500 是特殊 Boss 級別，遠高於普通特殊目標（60-120），需要多人合力
	"T124": {ID: "T124", Name: "深淵巨鯨", Type: TargetTypeSpecial, MultiplierMin: 80, MultiplierMax: 120, HP: 500, SpawnWeight: 1, Speed: 15, Lifetime: 45, LaborGain: 20, DifficultyFactor: 16.0, SpecialBehavior: "abyss_whale_boss"},
	// T125 黃金輪盤螃蟹（DAY-167）— 業界依據：King of Treasures Plus 2026「Roulette Crab — triggers Golden Roulette
	// bonus game, player hits SHOOT to stop wheel, wins the amount listed where it stops.」
	// 擊破後觸發個人黃金輪盤（8格：10x-200x），玩家點擊停止，結果預先決定（公平性保證）
	// 與千龍王輪盤（雙環，最高 1000x）不同：輪盤螃蟹是單環輪盤，更簡單直接，適合中等 betLevel 玩家
	// 設計：中等倍率（20-40x）+ 低 HP（50）+ 中高生成權重（6）= 常見且有輪盤爽感的特殊目標
	"T125": {ID: "T125", Name: "黃金輪盤螃蟹", Type: TargetTypeSpecial, MultiplierMin: 20, MultiplierMax: 40, HP: 50, SpawnWeight: 6, Speed: 45, Lifetime: 12, LaborGain: 10, DifficultyFactor: 16.0, SpecialBehavior: "roulette_crab"},
	// T126 獅子舞魚（DAY-168）— 業界依據：Fortune King Jackpot（TaDa Gaming 2026）「Lion Dance bonus — triggered by
	// special fish, delivers burst multiplier payouts with festive visual effects」
	// 擊破後觸發「獅子舞爆發」：全場隨機 3-5 個目標被「獅子舞光環」標記，玩家在 15 秒內擊破標記目標獲得 3x-10x 額外倍率加成
	// 設計：中等倍率（30-50x）+ 中等 HP（60）+ 中等生成權重（5）= 常見且有爆發爽感的特殊目標
	"T126": {ID: "T126", Name: "獅子舞魚", Type: TargetTypeSpecial, MultiplierMin: 30, MultiplierMax: 50, HP: 60, SpawnWeight: 5, Speed: 50, Lifetime: 14, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "lion_dance_burst"},
	// T127 漩渦魚（DAY-169）— 業界依據：Ocean King（Google Play 2026）「Vortex Fish — catching a Vortex Fish will
	// suck all fish of the same species in the area into a whirlpool, capturing them all at once.」
	// 擊破後觸發「漩渦吸引」：場上所有基礎目標（T001-T006）被吸入漩渦，全部擊破，獲得 0.55x 倍率獎勵
	// 設計差異：與黑洞（吸引所有目標）不同，漩渦魚是「同類型吸引」，讓玩家有策略性選擇
	// 設計：中等倍率（35-55x）+ 中等 HP（65）+ 中等生成權重（5）= 常見且有「一網打盡」爽感的特殊目標
	"T127": {ID: "T127", Name: "漩渦魚", Type: TargetTypeSpecial, MultiplierMin: 35, MultiplierMax: 55, HP: 65, SpawnWeight: 5, Speed: 55, Lifetime: 14, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "vortex_fish"},
	// T128 冰凍炸彈魚（DAY-170）— 業界依據：King of Ocean 2026「The freezing blast pauses an entire school for
	// a few seconds — useful when a high-tier creature is escaping the frame.」
	// 擊破後觸發「冰凍爆炸」：場上所有特殊目標（T101-T127）被冰凍 6 秒，停止移動
	// 設計差異：與黃金海龜（全場時間停止 8 秒）不同，冰凍炸彈魚只凍結特殊目標，讓玩家集中火力打高價值目標
	// 設計：中等倍率（40-60x）+ 中等 HP（70）+ 中等生成權重（5）= 常見且有「凍結高價值目標」策略感的特殊目標
	"T128": {ID: "T128", Name: "冰凍炸彈魚", Type: TargetTypeSpecial, MultiplierMin: 40, MultiplierMax: 60, HP: 70, SpawnWeight: 5, Speed: 45, Lifetime: 14, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "freeze_bomb"},
	// T129 冰釣魚（DAY-171）— 業界依據：Cozy Fishing Life（2026-05-10）「Winter Wheel — 8 segments x2-x10 multipliers
	// + bonus mode triggers」+ Ice Fishing Live（Evolution Gaming）「wheel triggers bonus fishing rounds」
	// 擊破後觸發「冰釣幸運輪盤」（8格：2x-10x 倍率加成），玩家在 5 秒內點擊停止
	// 觸發後玩家在 8 秒內所有擊破獎勵套用輪盤倍率，製造「黃金 8 秒」的爽感
	// 設計差異：與巨型章魚輪盤（950x 大獎）不同，冰釣輪盤是「倍率加成型」（2x-10x），讓玩家在短時間內所有擊破都有倍率加成
	// 設計：中等倍率（45-65x）+ 中等 HP（75）+ 中等生成權重（5）= 常見且有「黃金時間」爽感的特殊目標
	"T129": {ID: "T129", Name: "冰釣魚", Type: TargetTypeSpecial, MultiplierMin: 45, MultiplierMax: 65, HP: 75, SpawnWeight: 5, Speed: 50, Lifetime: 14, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "ice_fishing_wheel"},
	// T130 幸運彩蛋魚（DAY-172）— 業界依據：JILI Mega Fishing 2026「Giant Prize Fish lets you easily win great prizes,
	// with the chance for 5x multipliers」+ Ocean King 2026「Egg Fish drops golden eggs containing random rewards」
	// 擊破後掉落 1-5 個彩蛋（加權隨機），每個彩蛋隨機包含：金幣獎勵（50%）、倍率加成 ×2 持續 5 秒（30%）、特殊武器充能（20%）
	// 設計差異：與冰釣輪盤（玩家選擇停止）不同，彩蛋是「自動掉落+隨機開啟」，製造「每個彩蛋都是驚喜」的期待感
	// 設計：中高倍率（50-70x）+ 中等 HP（80）+ 中等生成權重（4）= 常見且有「開彩蛋」驚喜感的特殊目標
	"T130": {ID: "T130", Name: "幸運彩蛋魚", Type: TargetTypeSpecial, MultiplierMin: 50, MultiplierMax: 70, HP: 80, SpawnWeight: 4, Speed: 45, Lifetime: 14, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "lucky_egg_fish"},
	// T131 彩虹幸運魚（DAY-173）— 業界依據：Fisch Roblox 2026「Rainbow Leviathan — rare rainbow fish that triggers a luck boost event」
	// + Fish It 2026「Rainbow Throw — triggered by Prismatic enchant, increases luck for rare fish」
	// + Ocean King 2026「Rainbow Fish — when caught, all players receive a luck boost for 10 seconds」
	// 擊破後觸發「彩虹幸運時間」（10秒），全服所有玩家的擊破機率提升 20%（BASE_RTP × 1.2）
	// 設計差異：與幸運星魚（個人 ×2 倍率）不同，彩虹幸運魚是**全服共享的擊破機率加成**，製造「全服一起爽」的社交感
	// 設計：中高倍率（55-75x）+ 中等 HP（85）+ 稀有生成權重（3）= 稀有且有「全服幸運」社交感的特殊目標
	"T131": {ID: "T131", Name: "彩虹幸運魚", Type: TargetTypeSpecial, MultiplierMin: 55, MultiplierMax: 75, HP: 85, SpawnWeight: 3, Speed: 55, Lifetime: 14, LaborGain: 13, DifficultyFactor: 16.0, SpecialBehavior: "rainbow_lucky_fish"},
	// T132 海葵（DAY-174）— 業界依據：JILI Jackpot Fishing「Sea Anemone introduces unique effects —
	// tentacle attacks that spread to nearby fish」jackpotfishing-game.com 2026「Sea Anemone introduce
	// unique effects, such as chain lightning or explosive torpedoes, adding layers of strategy and excitement」
	// 擊破後觸手向 8 個方向延伸，每個方向命中最近的目標（300px 範圍），命中目標有 70% 機率擊破
	// 設計差異：與閃電鰻（連鎖跳躍，隨機目標）不同，海葵是「方向性觸手」（8方向固定延伸），更有視覺衝擊感
	// 設計：中高倍率（55-80x）+ 中等 HP（90）+ 稀有生成權重（3）= 稀有且有「觸手蔓延」視覺爽感的特殊目標
	"T132": {ID: "T132", Name: "海葵", Type: TargetTypeSpecial, MultiplierMin: 55, MultiplierMax: 80, HP: 90, SpawnWeight: 3, Speed: 40, Lifetime: 15, LaborGain: 13, DifficultyFactor: 16.0, SpecialBehavior: "sea_anemone"},
	// T133 幸運骰子魚（DAY-175）— 業界依據：Ocean King 3 Plus「Fast Bomb — randomly triggered bonus」
	// + 捕魚機業界「Dice Roll bonus — roll dice to determine reward multiplier」
	// + Fishing Carnival 2026「Dice Fish — catching triggers a dice roll, sum determines payout」
	// 擊破後觸發「幸運骰子」：擲 2 顆骰子（1-6），點數之和決定獎勵（2=20x/7=7x/12=50x/其他=點數x）
	// 設計差異：與輪盤（多格選擇）不同，骰子是「兩顆骰子點數之和」，機率分布符合真實骰子（7最常見）
	// 設計：中高倍率（60-85x）+ 中等 HP（85）+ 中等生成權重（4）= 常見且有「骰子期待感」的特殊目標
	"T133": {ID: "T133", Name: "幸運骰子魚", Type: TargetTypeSpecial, MultiplierMin: 60, MultiplierMax: 85, HP: 85, SpawnWeight: 4, Speed: 50, Lifetime: 14, LaborGain: 13, DifficultyFactor: 16.0, SpecialBehavior: "lucky_dice_fish"},
	// T134 火焰風暴魚（DAY-176）— 業界依據：Ocean King 3 Plus「Fire Storm feature — triggers a fire storm
	// that burns multiple fish simultaneously, creating chain combustion across the screen」
	// 擊破後觸發「火焰風暴」：場上隨機 4-8 個目標被火焰標記，15 秒內逐一燃燒擊破（每 1.5 秒一個），獎勵 0.6x 倍率
	// 設計差異：與漩渦魚（吸入基礎目標）不同，火焰風暴是「隨機標記任意目標」（包含特殊目標），
	// 且有「燃燒蔓延」的視覺過程（每 1.5 秒一個），製造「火焰逐漸蔓延」的戲劇感
	// 設計：中高倍率（65-90x）+ 中等 HP（90）+ 中等生成權重（3）= 稀有且有「火焰期待感」的特殊目標
	"T134": {ID: "T134", Name: "火焰風暴魚", Type: TargetTypeSpecial, MultiplierMin: 65, MultiplierMax: 90, HP: 90, SpawnWeight: 3, Speed: 45, Lifetime: 15, LaborGain: 14, DifficultyFactor: 16.0, SpecialBehavior: "fire_storm_fish"},
	// T135 黃金寶藏魚（DAY-177）— 業界依據：Ocean King 3 Plus「Golden Treasure feature — catching triggers
	// treasure chests to appear, players open them for random rewards」
	// + JILI Giant Anglerfish「shoot electricity to open treasure chests」
	// 擊破後觸發「黃金寶藏」：場上出現 3 個寶藏箱，玩家在 12 秒內點擊開啟，每個寶藏箱隨機包含：
	//   金幣獎勵（50%）、倍率加成 ×3 持續 8 秒（30%）、特殊武器充能（20%）
	// 設計差異：與幸運彩蛋魚（自動掉落+隨機開啟）不同，黃金寶藏是「玩家主動點擊開啟」，
	// 製造「我選擇開哪個箱子」的互動感；寶藏箱在場上可見，讓其他玩家也能看到「有人觸發了寶藏」
	// 設計：高倍率（70-100x）+ 較高 HP（95）+ 低生成權重（3）= 稀有且有「寶藏期待感」的特殊目標
	"T135": {ID: "T135", Name: "黃金寶藏魚", Type: TargetTypeSpecial, MultiplierMin: 70, MultiplierMax: 100, HP: 95, SpawnWeight: 3, Speed: 40, Lifetime: 15, LaborGain: 14, DifficultyFactor: 16.0, SpecialBehavior: "golden_treasure_fish"},
	// T136 美人魚（DAY-178）— 業界依據：Ocean King 3 Plus「The Mermaid feature — catching the Mermaid
	// triggers a healing event, restoring coins to the player and granting a brief luck boost」
	// 擊破後觸發「美人魚治癒」：為觸發玩家恢復 betLevel × 15 金幣（治癒機制），
	// 同時全服廣播「美人魚降臨」，讓其他玩家也能看到；觸發後 20 秒內擊破獎勵 +20% 幸運加成
	// 設計差異：與所有「攻擊型/倍率型」特殊目標不同，美人魚是唯一的「治癒型」目標，
	// 讓玩家在連續失敗後有「回血」的機會，製造「美人魚救了我」的情感連結
	// 設計：中倍率（45-65x）+ 中等 HP（70）+ 中等生成權重（4）= 常見且有「治癒期待感」的特殊目標
	"T136": {ID: "T136", Name: "美人魚", Type: TargetTypeSpecial, MultiplierMin: 45, MultiplierMax: 65, HP: 70, SpawnWeight: 4, Speed: 35, Lifetime: 14, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "mermaid_healing"},
	// T137 幸運草魚（DAY-179）— 業界依據：Ocean King 3 Plus「Lucky Shamrock Leprechaun Boss」
	// + Fisch Roblox 2026「Lucky Gold Pool — rainbow event triggers lucky fish spawns」
	// 擊破後觸發「幸運草爆發」：場上所有目標物獎勵 +50% 持續 10 秒（全服共享），
	// 同時隨機為 1-3 個玩家發放「幸運草金幣」（betLevel × 10-30x）
	// 設計差異：與黃金鯊魚（全服 ×1.5，12秒）不同，幸運草是「+50% 加成」（不是倍率乘法），
	// 且有「隨機發放金幣給玩家」的社交機制，讓被選中的玩家感到「幸運」
	// 設計：中高倍率（50-70x）+ 中等 HP（75）+ 中等生成權重（4）= 常見且有「幸運期待感」的特殊目標
	"T137": {ID: "T137", Name: "幸運草魚", Type: TargetTypeSpecial, MultiplierMin: 50, MultiplierMax: 70, HP: 75, SpawnWeight: 4, Speed: 45, Lifetime: 14, LaborGain: 13, DifficultyFactor: 16.0, SpecialBehavior: "lucky_clover_fish"},
	// T138 彩虹鯊魚（DAY-180）— 業界依據：JILI 2026 新特性「Rainbow Shark — triggers a rainbow burst
	// that randomly assigns 1.5x-3x multiplier bonuses to all targets on screen for 10 seconds」
	// 擊破後觸發「彩虹爆發」：場上所有存活目標隨機獲得 1.5x/2.0x/2.5x/3.0x 倍率加成標記，持續 10 秒
	// 設計差異：與黃金鯊魚（全服固定 ×1.5）不同，彩虹鯊魚是「每個目標倍率不同」（1.5x-3x），
	// 製造「哪個目標倍率最高？快去打！」的策略感；與幸運星魚（個人 ×2）不同，彩虹鯊魚是全服共享
	// 設計：中高倍率（55-75x）+ 中等 HP（80）+ 稀有生成權重（3）= 稀有且有「彩虹策略感」的特殊目標
	"T138": {ID: "T138", Name: "彩虹鯊魚", Type: TargetTypeSpecial, MultiplierMin: 55, MultiplierMax: 75, HP: 80, SpawnWeight: 3, Speed: 50, Lifetime: 14, LaborGain: 13, DifficultyFactor: 16.0, SpecialBehavior: "rainbow_shark_burst"},
	// T139 雷霆鯊魚（DAY-181）— 業界依據：JILI Jackpot Fishing「Thunder Shark brings unique abilities —
	// chain lightning that jumps between nearby fish, with no distance limit」
	// 擊破後觸發「雷霆連鎖閃電」：全場隨機跳躍（不限距離），最多 20 跳，每跳 75% 擊破機率
	// 設計差異：與 T103 閃電鰻（5跳/200px範圍）和 T118 皇家閃電鰻（15跳/300px範圍）不同，
	// 雷霆鯊魚是「全場無限距離隨機跳躍」，讓玩家看到閃電在全場「隨機亂跳」的混亂爽感
	// 設計：高倍率（60-80x）+ 中等 HP（85）+ 稀有生成權重（3）= 稀有且有「全場閃電」爽感的特殊目標
	"T139": {ID: "T139", Name: "雷霆鯊魚", Type: TargetTypeSpecial, MultiplierMin: 60, MultiplierMax: 80, HP: 85, SpawnWeight: 3, Speed: 55, Lifetime: 14, LaborGain: 14, DifficultyFactor: 16.0, SpecialBehavior: "thunder_shark_chain"},
	// T140 吸血鬼魚（DAY-182）— 業界依據：JILI 2026「The explicit multiplier of vampires increases
	// the more you fight, and there is a chance that you can enter the multiplier mode, up to X5」
	// 擊破後觸發「吸血鬼模式」：玩家每擊破一個目標，倍率累積 +0.1x（從 1.0x 開始），最高 5.0x，持續 15 秒
	// 設計差異：與幸運星魚（固定 ×2，10秒）不同，吸血鬼魚是「累積型倍率」（越打越高），
	// 製造「越打越爽」的正向反饋；玩家需要在 15 秒內盡量多打目標
	// 設計：高倍率（65-85x）+ 中等 HP（90）+ 稀有生成權重（3）= 稀有且有「累積爽感」的特殊目標
	"T140": {ID: "T140", Name: "吸血鬼魚", Type: TargetTypeSpecial, MultiplierMin: 65, MultiplierMax: 85, HP: 90, SpawnWeight: 3, Speed: 45, Lifetime: 14, LaborGain: 14, DifficultyFactor: 16.0, SpecialBehavior: "vampire_fish_escalating"},
	// T141 閃電魚（DAY-183）— 業界依據：Ocean King 3 Monster Awaken「Lightning Fish — Catching a
	// Lightning Fish will trigger a Lightning Chain. Lightning Chain will continue to catch fish
	// automatically until time runs out.」
	// 擊破後觸發「閃電自動連鎖」：系統自動每 0.5 秒選一個隨機目標發射閃電，持續 8 秒（最多 16 次）
	// 設計差異：與 T103 閃電鰻（手動觸發，5跳）和 T139 雷霆鯊魚（手動跳躍，20跳）不同，
	// 閃電魚是「全自動時間驅動連鎖」，玩家不需要操作，純粹享受「自動收割」的爽感
	// 設計：高倍率（65-85x）+ 中等 HP（85）+ 稀有生成權重（3）= 稀有且有「自動收割」爽感的特殊目標
	"T141": {ID: "T141", Name: "閃電魚", Type: TargetTypeSpecial, MultiplierMin: 65, MultiplierMax: 85, HP: 85, SpawnWeight: 3, Speed: 60, Lifetime: 14, LaborGain: 14, DifficultyFactor: 16.0, SpecialBehavior: "lightning_auto_chain"},
	// T142 隕石魚（DAY-184）— 業界依據：Royal Fishing JILI「Dragon Wrath — unleash a massive
	// meteorite attack across the centre screen, simultaneously targeting multiple fish」
	// 擊破後觸發「隕石雨」：5-10 顆隕石從天而降，每顆命中隨機目標，70% 擊破機率，獎勵 0.60x 倍率
	// 設計：高倍率（70-90x）+ 中等 HP（90）+ 稀有生成權重（3）= 稀有且有「天降神兵」爽感的特殊目標
	"T142": {ID: "T142", Name: "隕石魚", Type: TargetTypeSpecial, MultiplierMin: 70, MultiplierMax: 90, HP: 90, SpawnWeight: 3, Speed: 50, Lifetime: 14, LaborGain: 14, DifficultyFactor: 16.0, SpecialBehavior: "meteor_fish_shower"},
	// T143 鳳凰魚（DAY-185）— 業界依據：Ocean King 3 Plus「Phoenix Fish — when defeated, the
	// Phoenix Fish triggers a rebirth explosion that deals massive damage to all fish on screen,
	// with the Phoenix rising from the ashes to grant a 30-second luck boost」
	// 擊破後觸發「涅槃爆炸」：全場同時爆炸（普通 80%/特殊 50%/BOSS 20%），爆炸後全服 +30% 加成 30 秒
	// 設計：高倍率（75-95x）+ 中等 HP（95）+ 稀有生成權重（2）= 極稀有且有「全場清場+重生加成」雙重爽感
	"T143": {ID: "T143", Name: "鳳凰魚", Type: TargetTypeSpecial, MultiplierMin: 75, MultiplierMax: 95, HP: 95, SpawnWeight: 2, Speed: 45, Lifetime: 15, LaborGain: 15, DifficultyFactor: 16.0, SpecialBehavior: "phoenix_fish_rebirth"},
	// T144 龍龜不死 Boss（DAY-186）— 業界依據：Royal Fishing JILI「Immortal Boss mechanic —
	// Golden Toad and Ancient Crocodile bosses appear randomly and award consecutive wins
	// ranging from 50X to 150X until they leave the screen. This creates extended winning
	// sequences impossible in standard fish games.」
	// 龍龜不死機制：出現後不會被擊破（Immortal），每次命中給 50-150x betLevel 獎勵，
	// 直到 Lifetime 結束離開畫面，全服廣播每次命中
	// 設計差異：與普通 BOSS（需要擊破）完全不同，龍龜是「持續收割型」，
	// 玩家不需要擊破，只要命中就有獎勵，製造「穩定收益」的安心感；
	// 全服共享龍龜，所有玩家都可以打，製造「搶打龍龜」的競爭感
	// 設計：中等倍率（30-50x）+ 超高 HP（99999，不死）+ 極稀有生成權重（1）+ 30秒 Lifetime
	// HP=99999 確保永遠不會被擊破；MultiplierMin/Max 用於顯示，實際獎勵由 handler 計算
	"T144": {ID: "T144", Name: "龍龜不死", Type: TargetTypeSpecial, MultiplierMin: 30, MultiplierMax: 50, HP: 99999, SpawnWeight: 1, Speed: 25, Lifetime: 30, LaborGain: 20, DifficultyFactor: 16.0, SpecialBehavior: "immortal_boss"},
	// T145 連鎖爆炸魚（DAY-187）— 業界依據：Royal Fishing「chain reaction mechanic —
	// players can trigger multiple explosions to capture additional fish within a blast radius」
	// 擊破後在原位爆炸（200px 半徑），爆炸命中的目標 75% 機率擊破（0.65x 倍率），
	// 若命中的目標也是 T145 則繼續引爆（連鎖反應，最多 5 層）
	// 設計差異：與炸彈武器（玩家主動放置）不同，連鎖爆炸魚是「被動觸發的連鎖反應」；
	// 與漩渦魚（吸引同類）不同，連鎖爆炸魚是「位置驅動的爆炸傳播」；
	// 最多 5 層連鎖，讓玩家有「一顆引爆全場」的爽感，但不會無限連鎖（平衡 RTP）
	// 設計：中等倍率（25-45x）+ 中等 HP（60）+ 常見生成權重（4）= 常見且有「連鎖爆炸」爽感
	"T145": {ID: "T145", Name: "連鎖爆炸魚", Type: TargetTypeSpecial, MultiplierMin: 25, MultiplierMax: 45, HP: 60, SpawnWeight: 4, Speed: 55, Lifetime: 12, LaborGain: 10, DifficultyFactor: 16.0, SpecialBehavior: "chain_bomb"},
	// T146 巨型鱷魚（DAY-188）— 業界依據：JILI Mega Fishing「giant crocodiles awaken to hunt fish
	// on the fish farm to accumulate big prizes」
	// 出現後每 2 秒主動「獵食」場上一個目標（優先高倍率普通目標），
	// 獵食成功給全服玩家共享獎勵（0.4x 倍率），玩家擊破鱷魚本身獲得鱷魚倍率 + 累積獎池 50%
	// 設計差異：與不死 BOSS（玩家打它）不同，鱷魚是「它主動打其他魚」，
	// 製造「看著鱷魚在場上橫行霸道」的緊張感，玩家需要決策：讓它繼續獵食累積獎池，還是立刻擊破？
	// 設計：高倍率（40-80x）+ 高 HP（120）+ 稀有生成（2）= 稀有且有「獵食決策」策略感
	"T146": {ID: "T146", Name: "巨型鱷魚", Type: TargetTypeSpecial, MultiplierMin: 40, MultiplierMax: 80, HP: 120, SpawnWeight: 2, Speed: 30, Lifetime: 25, LaborGain: 18, DifficultyFactor: 16.0, SpecialBehavior: "crocodile_hunter"},
	// T147 時間炸彈魚（DAY-189）— 業界靈感：Ocean King 炸彈魚概念 + 倒數計時緊張感設計
	// 出現後顯示 10 秒倒數計時：
	//   - 倒數結束前玩家擊破 → 「拆彈成功」：全服 +25% 加成持續 15 秒
	//   - 倒數結束無人擊破 → 「炸彈爆炸」：全場目標 80% 擊破機率（0.5x 倍率，全服共享）
	// 設計差異：與連鎖爆炸魚（被動觸發）不同，時間炸彈魚是「主動倒數」，製造「搶時間」的緊張感；
	// 「拆彈成功」的加成獎勵讓玩家有「英雄感」，「炸彈爆炸」讓玩家有「清場爽感」
	// 設計：中高倍率（35-65x）+ 中等 HP（75）+ 常見生成（3）= 常見且有「倒數緊張感」
	"T147": {ID: "T147", Name: "時間炸彈魚", Type: TargetTypeSpecial, MultiplierMin: 35, MultiplierMax: 65, HP: 75, SpawnWeight: 3, Speed: 50, Lifetime: 15, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "time_bomb_fish"},
	// T148 三重幸運魚（DAY-190）— 業界靈感：TaDa Gaming TriLuck™ 2026「trigger three different feature specifications simultaneously」
	// 擊破後同時觸發三重效果：金幣雨（betLevel×20-50x）+ 倍率加成（+50%，12秒）+ 武器充能（龍怒/魚雷/軌道炮隨機一發）
	// 三個效果同時生效，製造「三重爽感」；全服廣播讓其他玩家看到「有人觸發了三重幸運」
	"T148": {ID: "T148", Name: "三重幸運魚", Type: TargetTypeSpecial, MultiplierMin: 40, MultiplierMax: 70, HP: 80, SpawnWeight: 3, Speed: 45, Lifetime: 14, LaborGain: 13, DifficultyFactor: 16.0, SpecialBehavior: "triple_lucky_fish"},
	// T149 魚群領袖（DAY-191）— 業界靈感：Ocean King 3 Plus「School of Fish — when one fish is caught, others scatter in panic」
	// 擊破後觸發「魚群驚嚇」：場上所有基礎目標（T001-T006）HP 降低 50%，持續 8 秒
	// 讓玩家在「魚群驚嚇」中快速收割，製造「緊張但有利」的感覺
	"T149": {ID: "T149", Name: "魚群領袖", Type: TargetTypeSpecial, MultiplierMin: 30, MultiplierMax: 50, HP: 65, SpawnWeight: 4, Speed: 55, Lifetime: 12, LaborGain: 11, DifficultyFactor: 16.0, SpecialBehavior: "school_panic"},
	// T150 搖滾骷髏魚（DAY-192）— 業界依據：JILI 2026「Rock Skeleton Concert — Rock Skeleton and Super Awakening Performance, up to 3,000x」
	// 擊破後觸發「演唱會模式」（15 秒）：每 1 秒音符炸彈命中 2-4 個目標（70% 擊破機率，0.60x 倍率）
	// 第 10 秒觸發「超級覺醒高潮」：全場所有目標 HP 降低 70%，持續 5 秒
	// 演唱會結束後：≥10 個擊破 → 全服 +30% 加成 10 秒（安可獎勵）
	"T150": {ID: "T150", Name: "搖滾骷髏魚", Type: TargetTypeSpecial, MultiplierMin: 35, MultiplierMax: 60, HP: 70, SpawnWeight: 3, Speed: 50, Lifetime: 13, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "rock_skeleton_concert"},
	// T151 電流水母（DAY-193）— 業界依據：King of Ocean 2026「electric jellyfish chains current between adjacent targets, paying multipliers from every link in the chain」
	// 擊破後建立「電流網路」：場上所有相鄰目標（200px 內）之間建立電流連接，每條連接 65% 擊破機率（0.55x 倍率）
	// 密集目標群形成更多連接，製造「越多魚越爽」的策略感；電流選擇較低 HP 的目標擊破
	"T151": {ID: "T151", Name: "電流水母", Type: TargetTypeSpecial, MultiplierMin: 30, MultiplierMax: 55, HP: 60, SpawnWeight: 4, Speed: 45, Lifetime: 12, LaborGain: 11, DifficultyFactor: 16.0, SpecialBehavior: "electric_jellyfish_network"},
	// T152 長龍王（DAY-194）— 業界依據：Royal Fishing JILI「ChainLong King — dual-ring roulette activates when captured.
	// You control when the pointer stops, multiplying inner and outer ring values together.
	// Maximum combination delivers 350X, whilst the ChainLong King itself can award up to 1000X mega wins.」
	// 擊破後觸發「雙環輪盤」互動（個人）：內環 5x/10x/20x/50x × 外環 1x/2x/3x/5x/7x = 最高 350x
	// 特殊：1% 機率觸發「千倍大獎」（1000x），跳過輪盤直接給獎
	"T152": {ID: "T152", Name: "長龍王", Type: TargetTypeSpecial, MultiplierMin: 50, MultiplierMax: 100, HP: 100, SpawnWeight: 2, Speed: 35, Lifetime: 18, LaborGain: 15, DifficultyFactor: 16.0, SpecialBehavior: "chainlong_king_roulette"},
	// T153 鑽頭龍蝦（DAY-195）— 業界依據：Royal Fishing JILI「Drill Bit Lobster (80X) — fires a penetrating drill
	// that passes through multiple fish before self-detonating, capturing everything in the explosion radius.
	// Mechanical marvel with penetrating drill projectiles.」
	// 擊破後發射「穿透鑽頭」：沿隨機方向穿透最多 5 個目標（80% 擊破機率，0.70x 倍率）
	// 穿透結束後在終點「自爆」（300px 半徑，75% 擊破機率，0.65x 倍率）
	"T153": {ID: "T153", Name: "鑽頭龍蝦", Type: TargetTypeSpecial, MultiplierMin: 60, MultiplierMax: 80, HP: 90, SpawnWeight: 3, Speed: 40, Lifetime: 15, LaborGain: 13, DifficultyFactor: 16.0, SpecialBehavior: "drill_lobster_penetrate"},
	// T154 巨型鮟鱇魚（DAY-196）— 業界依據：JILI Mega Fishing「Giant Anglerfish can shoot electricity
	// to open treasure chests, giant crocodiles awaken to hunt fish on the fish farm to accumulate big prizes!」
	// 出現後每 3 秒電擊一次（最多 8 次）：命中 T102 寶箱怪強制開箱（3-5x 倍率）；命中普通目標 70% 擊破（0.60x 倍率）
	// 5% 機率「超級電擊」全場所有目標同時受到電擊；玩家擊破獲得基礎倍率 + 累積電擊獎池 40%
	"T154": {ID: "T154", Name: "巨型鮟鱇魚", Type: TargetTypeSpecial, MultiplierMin: 55, MultiplierMax: 85, HP: 100, SpawnWeight: 2, Speed: 30, Lifetime: 28, LaborGain: 14, DifficultyFactor: 16.0, SpecialBehavior: "anglerfish_electric_zap"},
	// T155 神秘龍魚（DAY-197）— 業界依據：Ocean King 3「Mystic Dragon — Catch this fish to get 8 waves
	// and have more chances to kill any fish on the screen.」
	// 擊破後觸發「八波龍息攻擊」：每波隨機選 3-5 個目標（65% 擊破機率，0.55x 倍率）
	// 第 8 波「龍怒爆發」：全場所有目標（85% 擊破機率，0.70x 倍率）；全服共享獎勵
	"T155": {ID: "T155", Name: "神秘龍魚", Type: TargetTypeSpecial, MultiplierMin: 65, MultiplierMax: 90, HP: 95, SpawnWeight: 2, Speed: 35, Lifetime: 20, LaborGain: 15, DifficultyFactor: 16.0, SpecialBehavior: "mystic_dragon_waves"},
	// T156 幽靈魚真身（DAY-198）— 原創設計，靈感來自 Fisch Phantom Mutation（4x 倍率）
	// 生成時同時在場上生成 2-3 個「幻影分身」（T156C，HP=1）
	// 擊破真身觸發「幽靈爆發」：所有幻影分身同時爆炸（50% 擊破機率，0.50x 倍率）
	// 擊破幻影分身：給 1x betLevel 安慰獎
	"T156": {ID: "T156", Name: "幽靈魚", Type: TargetTypeSpecial, MultiplierMin: 45, MultiplierMax: 70, HP: 80, SpawnWeight: 3, Speed: 45, Lifetime: 18, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "ghost_fish_phantom"},
	// T156C 幽靈魚幻影分身（DAY-198）— 由 T156 生成時自動創建，HP=1，外觀與 T156 相同
	// 擊破給 1x betLevel 安慰獎；真身被擊破時同時爆炸
	"T156C": {ID: "T156C", Name: "幽靈魚（幻影）", Type: TargetTypeSpecial, MultiplierMin: 1, MultiplierMax: 1, HP: 1, SpawnWeight: 0, Speed: 45, Lifetime: 18, LaborGain: 0, DifficultyFactor: 16.0, SpecialBehavior: "ghost_fish_clone"},
	// T157 雷霆龍蝦（DAY-199）— 業界依據：Royal Fishing JILI「Thunderbolt Lobster feature —
	// provides 15 seconds of free play followed by automatic shooting from the Thunderbolt Turret.
	// Players can earn extra seconds during this period to extend gameplay and increase reward potential.」
	// 擊破後觸發「雷霆砲台模式」（15秒）：系統自動每 0.5 秒選最高價值目標射擊（85% 擊破機率，0.75x 倍率）
	// 每擊破一個目標 +0.5 秒（最多延長到 30 秒）；全服廣播「雷霆砲台啟動」
	"T157": {ID: "T157", Name: "雷霆龍蝦", Type: TargetTypeSpecial, MultiplierMin: 50, MultiplierMax: 80, HP: 85, SpawnWeight: 3, Speed: 40, Lifetime: 16, LaborGain: 13, DifficultyFactor: 16.0, SpecialBehavior: "thunderbolt_lobster_free_play"},
	// T158 冰鳳凰（DAY-200）— 業界依據：Royal Fishing JILI「Ice Phoenix Awaken Feature —
	// fixed jackpot mechanic that awards up to 300x the bet when players eliminate the Ice Phoenix boss.
	// Awaken Boss with 30x basic multiplier. Power Up attack delivers 6x-10x boost for rewards up to 300x.」
	// 擊破後觸發「冰鳳凰覺醒」：基礎獎勵 30x betLevel；Power Up 攻擊 3-5 個目標（6-10x 倍率，70% 擊破機率）
	// 5% 機率觸發「冰霜爆發」：全場所有目標（50% 擊破機率，0.60x 倍率）；最高組合 300x
	"T158": {ID: "T158", Name: "冰鳳凰", Type: TargetTypeSpecial, MultiplierMin: 80, MultiplierMax: 120, HP: 110, SpawnWeight: 2, Speed: 30, Lifetime: 20, LaborGain: 16, DifficultyFactor: 16.0, SpecialBehavior: "ice_phoenix_awaken"},
	// T159 連環炸彈蟹（DAY-201）— 業界依據：Royal Fishing JILI「Serial Bomb Crab (70x) —
	// orange crab with panda face and skull bomb designs. Triggers large-scale multiple explosions
	// across screen, capturing fish within each explosion range.
	// Each bomb creates expanding capture zones for massive multi-target eliminations.」
	// 擊破後觸發「連環爆炸」：3-5 顆炸彈依序爆炸（每顆間隔 600ms）
	// 每顆炸彈：250px 半徑，75% 擊破機率，0.65x 倍率；炸彈位置隨機分散在場上
	"T159": {ID: "T159", Name: "連環炸彈蟹", Type: TargetTypeSpecial, MultiplierMin: 55, MultiplierMax: 75, HP: 80, SpawnWeight: 3, Speed: 45, Lifetime: 14, LaborGain: 13, DifficultyFactor: 16.0, SpecialBehavior: "serial_bomb_crab"},

	// T160 深淵漩渦魚（DAY-202）— 業界依據：Ocean King 2「Vortex Fish — sucks all fish of the same
	// species into a whirlpool. Catching a Vortex Fish will suck all fish of the same species in the
	// area into a whirlpool.」+ SteamDB OceanFest 2026「Abyssal Vortex (persistent whirlpool)」
	// 擊破後在擊破位置生成「深淵漩渦」（持續 5 秒）：
	//   每 0.5 秒吸引脈衝（500px 半徑內目標向中心移動 180px）
	//   進入 100px 中心：80% 擊破機率，0.70x 倍率
	//   漩渦結束後深淵爆炸：300px 半徑，60% 擊破機率，0.55x 倍率
	"T160": {ID: "T160", Name: "深淵漩渦魚", Type: TargetTypeSpecial, MultiplierMin: 45, MultiplierMax: 70, HP: 75, SpawnWeight: 3, Speed: 50, Lifetime: 13, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "abyss_vortex_pull"},

	// T161 座頭鯨（DAY-203）— 業界依據：Royal Fishing JILI「Humpback Whale offers 90-150x with
	// 15x base multiplier. Awaken Boss mechanic — triggers wave attack that sweeps the screen.
	// The Humpback Whale's signature breach mechanic creates massive splash zones.」
	// 擊破後觸發「鯨歌覺醒」：基礎獎勵 15x betLevel；3 波波浪攻擊（每波 3 個目標，65% 擊破機率，0.60x 倍率）
	// 5% 機率觸發「深海巨浪」：全場所有目標（60% 擊破機率，0.65x 倍率）；最高組合 150x
	"T161": {ID: "T161", Name: "座頭鯨", Type: TargetTypeSpecial, MultiplierMin: 90, MultiplierMax: 150, HP: 120, SpawnWeight: 2, Speed: 25, Lifetime: 22, LaborGain: 18, DifficultyFactor: 16.0, SpecialBehavior: "humpback_whale_awaken"},

	// T162 自由旋轉魚（DAY-204）— 業界依據：Galaxsys King of Ocean 2026
	// 「Free Spin Fish, Captain Fish, and Money Fish trigger bonus rounds, extra multipliers, and instant payouts.」
	// 擊破後觸發「個人免費射擊模式」（10秒，不扣費）：
	//   系統每 0.6 秒自動選最高價值目標射擊（80% 擊破機率，0.80x 倍率）
	//   每擊破一個目標 +1 秒（最多延長到 20 秒）；個人冷卻 30 秒
	"T162": {ID: "T162", Name: "自由旋轉魚", Type: TargetTypeSpecial, MultiplierMin: 35, MultiplierMax: 60, HP: 70, SpawnWeight: 4, Speed: 50, Lifetime: 14, LaborGain: 11, DifficultyFactor: 16.0, SpecialBehavior: "free_spin_fish"},

	// T163 獎池龍（DAY-205）— 業界依據：JILI Jackpot Fishing
	// 「special targets like the Jackpot Fish and Jackpot Dragon offering chances at substantial prizes.
	//  With the potential for high payouts up to 1000 times the bet.」
	// 擊破後觸發「獎池抽獎」（個人）：加權隨機選擇 Jackpot 等級
	//   Mini(70%) / Minor(20%) / Major(8%) / Grand(2%)；個人冷卻 60 秒
	"T163": {ID: "T163", Name: "獎池龍", Type: TargetTypeSpecial, MultiplierMin: 50, MultiplierMax: 100, HP: 100, SpawnWeight: 2, Speed: 30, Lifetime: 18, LaborGain: 15, DifficultyFactor: 16.0, SpecialBehavior: "jackpot_dragon_draw"},
	// T164 彗星魚（DAY-206）— 業界依據：Ocean King 3 Plus「Comet Fish — streaks across the screen
	// leaving a trail of explosions, each explosion has a chance to capture fish in its radius.」
	// 生成後沿弧線軌跡飛越全場（1.5秒），沿途 7 個爆炸點（200px 半徑，70% 擊破，0.65x 倍率）
	// 最終超新星爆炸（400px 半徑，80% 擊破，0.75x 倍率）；玩家擊破可提前引爆；全服冷卻 40 秒
	"T164": {ID: "T164", Name: "彗星魚", Type: TargetTypeSpecial, MultiplierMin: 40, MultiplierMax: 65, HP: 75, SpawnWeight: 3, Speed: 60, Lifetime: 12, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "comet_fish_trail"},
	// T165 黃金波浪魚（DAY-207）— 業界依據：Ocean King 4 Brand New World
	// 「Golden Wave Fish — triggers a golden tidal wave that sweeps across the entire screen,
	//  temporarily boosting all multipliers by 2x for 8 seconds.」
	// 擊破後觸發「黃金波浪」：8 列掃場（每列 150ms，70% 擊破，0.60x 倍率）
	// 波浪結束後全服 ×2.0 倍率加成 8 秒（乘法，最強加成）；全服冷卻 50 秒
	"T165": {ID: "T165", Name: "黃金波浪魚", Type: TargetTypeSpecial, MultiplierMin: 45, MultiplierMax: 70, HP: 80, SpawnWeight: 3, Speed: 45, Lifetime: 14, LaborGain: 13, DifficultyFactor: 16.0, SpecialBehavior: "golden_wave_boost"},
	// T166 深海龍王（DAY-208）— 業界依據：Royal Fishing JILI「Dragon Wrath — accumulate wrath value
	//  through shooting, then unleash devastating meteor strikes across the entire screen.」
	// 擊破後觸發「龍王怒火蓄力模式」（12秒）：全服合力射擊累積龍怒值（+1/shot）
	// 達到 20 點 → 龍怒隕石雨（5顆，350px，80%，0.75x）；未達到 → 小型龍怒（3顆，250px，65%，0.60x）
	// 全服冷卻 45 秒；全服共享獎勵；製造「大家一起打才能觸發龍怒」的社群感
	"T166": {ID: "T166", Name: "深海龍王", Type: TargetTypeSpecial, MultiplierMin: 55, MultiplierMax: 85, HP: 90, SpawnWeight: 2, Speed: 35, Lifetime: 16, LaborGain: 14, DifficultyFactor: 16.0, SpecialBehavior: "dragon_king_charge"},
	// T167 幸運金幣魚（DAY-209）— 業界依據：Galaxsys King of Ocean 2026
	// 「Money Fish trigger instant payouts.」
	// 擊破後立即觸發「金幣爆發」：加權隨機即時獎勵 5x(50%)/10x(30%)/20x(15%)/50x(5%) × betLevel
	// 3% 機率觸發「黃金爆發」：全場所有目標 HP 降低 80%（持續 5 秒）；個人冷卻 15 秒
	"T167": {ID: "T167", Name: "幸運金幣魚", Type: TargetTypeSpecial, MultiplierMin: 20, MultiplierMax: 50, HP: 50, SpawnWeight: 5, Speed: 55, Lifetime: 12, LaborGain: 8, DifficultyFactor: 16.0, SpecialBehavior: "fortune_coin_burst"},
	// T168 幸運熱區魚（DAY-210）— 業界依據：Ocean King 4 Brand New World 2025
	// 「Golden Zone — a glowing area appears on screen, all fish inside receive a 2x multiplier bonus.
	//  Zone lasts 8 seconds then explodes, capturing all remaining fish within the zone.」
	// 擊破後在場上建立「幸運熱區」（半徑 280px，持續 8 秒）：
	//   1. 熱區內所有目標獲得 ×2.0 倍率加成（乘法）
	//   2. 每 1 秒「熱區脈衝」：熱區內目標 HP 降低 15%
	//   3. 8 秒後「熱區爆炸」：熱區內所有目標 75% 擊破機率（0.65x 倍率，全服共享）
	//   4. 個人冷卻 20 秒；全服冷卻 30 秒
	"T168": {ID: "T168", Name: "幸運熱區魚", Type: TargetTypeSpecial, MultiplierMin: 30, MultiplierMax: 55, HP: 65, SpawnWeight: 4, Speed: 50, Lifetime: 13, LaborGain: 11, DifficultyFactor: 16.0, SpecialBehavior: "lucky_hot_zone"},
	// T169 幸運三叉魚（DAY-211）— 業界依據：TaDa Gaming TriLuck™ Series 2026
	// 「Within the TriLuck™ Series, you can trigger three different feature specifications,
	//  ranging from win multipliers, jackpot bonuses, collecting all rewards, and more unique features.」
	// 擊破後觸發「三叉幸運儀式」（個人互動）：
	//   三個獨立轉盤同時旋轉，玩家依序點擊停止：
	//   - 轉盤 A（金幣）：即時金幣獎勵 10x/20x/30x/50x/100x × betLevel
	//   - 轉盤 B（倍率）：下一次擊破倍率加成 ×1.5/×2.0/×2.5/×3.0/×5.0（持續 15 秒）
	//   - 轉盤 C（特效）：HP削減/免費射擊/全服廣播/小型清場（隨機）
	//   個人冷卻 25 秒；超時 12 秒自動停止
	"T169": {ID: "T169", Name: "幸運三叉魚", Type: TargetTypeSpecial, MultiplierMin: 35, MultiplierMax: 60, HP: 70, SpawnWeight: 4, Speed: 48, Lifetime: 13, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "lucky_trident"},
	// T170 時間凍結魚（DAY-212）— 業界依據：Evolution Ice Fishing Live 2026「frozen paradise」概念
	// 業界原創設計：「時間停止」讓玩家在凍結期間免費射擊所有靜止目標
	// 擊破後觸發「時間凍結」（5秒）：
	//   1. 全場所有目標物靜止不動（Speed=0）
	//   2. 凍結期間玩家射擊命中率 +30%
	//   3. 凍結結束後「解凍爆炸」：所有被命中過的目標 HP -50% + 60% 擊破機率
	//   個人冷卻 20 秒；全服冷卻 35 秒
	"T170": {ID: "T170", Name: "時間凍結魚", Type: TargetTypeSpecial, MultiplierMin: 30, MultiplierMax: 55, HP: 65, SpawnWeight: 4, Speed: 52, Lifetime: 13, LaborGain: 11, DifficultyFactor: 16.0, SpecialBehavior: "time_freeze"},
	// T171 彩虹稜鏡魚（DAY-213）— 業界依據：Dive Down 2026「Rainbow is the strongest mutation with 3.0x multiplier」
	// 業界原創「稜鏡折射染色」機制：擊破後隨機選場上最多 5 個目標染色（紅/橙/黃/綠/藍）
	// 每種顏色對應不同倍率加成（×1.5/×2.0/×2.5/×3.0/×5.0）；染色持續 10 秒
	// 染色期間擊破對應顏色目標獲得對應倍率加成（乘法）
	// 10 秒後「彩虹爆炸」：所有仍存活的染色目標同時爆炸（70% 擊破機率，0.65x 倍率）
	// 個人冷卻 25 秒；全服廣播染色開始/爆炸結算
	"T171": {ID: "T171", Name: "彩虹稜鏡魚", Type: TargetTypeSpecial, MultiplierMin: 35, MultiplierMax: 65, HP: 75, SpawnWeight: 3, Speed: 50, Lifetime: 14, LaborGain: 13, DifficultyFactor: 16.0, SpecialBehavior: "rainbow_prism"},
	// T172 黃金累積魚（DAY-214）— 業界依據：Evolution Ice Fishing Live 2026「random multipliers 2x-10x」
	// 業界原創「全服累積爆發」機制：T172 出現後，每次任何玩家擊破任何目標，累積槽 +1（最多 20 點）
	// 累積槽滿 → 自動觸發「黃金爆發」：全場所有目標 HP -60% + 全服 ×2.0 倍率加成 8 秒
	// 玩家擊破黃金累積魚本身 → 「提前引爆」（不論累積多少）
	// 全服廣播累積進度（每 5 點）；全服冷卻 40 秒
	"T172": {ID: "T172", Name: "黃金累積魚", Type: TargetTypeSpecial, MultiplierMin: 40, MultiplierMax: 70, HP: 80, SpawnWeight: 3, Speed: 45, Lifetime: 15, LaborGain: 14, DifficultyFactor: 16.0, SpecialBehavior: "golden_accumulator"},
	// T173 幸運鏡像魚（DAY-215）— 業界原創「鏡像複製」機制
	// 設計：擊破 T173 後觸發「鏡像複製」：在場上隨機選 3 個目標，為每個目標建立「鏡像分身」
	//   - 鏡像分身 HP = 原目標 HP × 50%，倍率 = 原目標倍率 × 1.5
	//   - 鏡像分身持續 8 秒；擊破鏡像分身獲得 ×1.5 倍率加成
	//   - 8 秒後所有未被擊破的鏡像分身「鏡像爆炸」（60% 擊破機率，0.60x 倍率）
	//   - 個人冷卻 20 秒；全服廣播鏡像建立/爆炸
	"T173": {ID: "T173", Name: "幸運鏡像魚", Type: TargetTypeSpecial, MultiplierMin: 35, MultiplierMax: 60, HP: 70, SpawnWeight: 3, Speed: 48, Lifetime: 13, LaborGain: 13, DifficultyFactor: 16.0, SpecialBehavior: "lucky_mirror"},
	// T174 詛咒毒魚（DAY-216）— 業界原創「詛咒反轉」機制
	// 設計：T174 出現後，場上隨機 3 個目標被「詛咒標記」（紫色）
	//   - 詛咒目標被擊破：獎勵 ×2.5 倍率（高風險高報酬）
	//   - 詛咒目標逃跑：觸發「詛咒懲罰」— 下一次擊破任何目標獎勵 ×0.5（持續 5 秒）
	//   - 擊破 T174 本身：「解除詛咒」— 移除所有詛咒標記 + 解咒獎勵 10x betLevel
	//   - 個人冷卻 18 秒；全服廣播詛咒標記/解除/懲罰
	"T174": {ID: "T174", Name: "詛咒毒魚", Type: TargetTypeSpecial, MultiplierMin: 30, MultiplierMax: 55, HP: 65, SpawnWeight: 3, Speed: 52, Lifetime: 12, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "cursed_poison"},
	// T175 幸運拍賣魚（DAY-217）— 業界原創「全服競標」機制
	// 設計：T175 出現後，開啟「全服競標」（持續 8 秒）：
	//   - 任何玩家可以「出價」（消耗 betLevel × 5 籌碼），出價最高者獲得「大獎控制權」
	//   - 競標結束後，最高出價者獲得「大獎控制權」（5 秒內自動射擊最高價值目標，0.85x 倍率）
	//   - 競標失敗者退還 50% 出價籌碼；若無人競標，T175 自動逃跑
	//   - 個人冷卻 20 秒；全服廣播競標開始/出價/結算
	"T175": {ID: "T175", Name: "幸運拍賣魚", Type: TargetTypeSpecial, MultiplierMin: 35, MultiplierMax: 65, HP: 70, SpawnWeight: 3, Speed: 48, Lifetime: 15, LaborGain: 13, DifficultyFactor: 16.0, SpecialBehavior: "lucky_auction"},
	// T176 幸運進化魚（DAY-218）— 業界原創「三段進化」機制
	// 設計：T176 出現後，每次任何玩家命中它（不需要擊破），它會「進化」：
	//   - 進化 1（命中 3 次）：HP -30%，倍率 ×1.5；進化 2（命中 6 次）：HP -50%，倍率 ×2.5
	//   - 進化 3（命中 9 次）：HP -70%，倍率 ×4.0；3 秒後自動「終極爆發」
	//   - 玩家擊破進化魚本身：立即觸發「終極爆發」（全場 HP -60% + 全服 ×4.0 倍率加成 6 秒）
	"T176": {ID: "T176", Name: "幸運進化魚", Type: TargetTypeSpecial, MultiplierMin: 40, MultiplierMax: 70, HP: 90, SpawnWeight: 3, Speed: 45, Lifetime: 18, LaborGain: 14, DifficultyFactor: 16.0, SpecialBehavior: "lucky_evolution"},
	// T177 幸運連鎖感染魚（DAY-219）— 業界原創「病毒式蔓延」機制
	// 設計：擊破 T177 後觸發「感染標記」：
	//   - 場上隨機 2 個目標被「感染」（綠色標記）
	//   - 感染目標每 2 秒向相鄰目標（300px 內）傳播感染（最多蔓延 3 層，最多 8 個感染目標）
	//   - 感染目標被擊破：獎勵 ×2.0 倍率加成（乘法）
	//   - 12 秒後所有感染目標同時「感染爆發」（75% 擊破機率，0.65x 倍率，全服共享）
	//   - 個人冷卻 22 秒；全服廣播感染建立/蔓延/爆發
	"T177": {ID: "T177", Name: "幸運連鎖感染魚", Type: TargetTypeSpecial, MultiplierMin: 35, MultiplierMax: 60, HP: 70, SpawnWeight: 3, Speed: 50, Lifetime: 14, LaborGain: 13, DifficultyFactor: 16.0, SpecialBehavior: "lucky_infection"},
	// T178 幸運反彈魚（DAY-220）— 業界原創「子彈反彈」機制
	// 設計：擊破 T178 後觸發「反彈模式」（8秒）：
	//   - 玩家的每次射擊在命中目標後，子彈會「反彈」到最近的另一個目標
	//   - 反彈範圍：第1跳 200px，第2跳 150px，第3跳 100px（最多 3 跳）
	//   - 每次反彈命中：60% 擊破機率，0.55x 倍率
	//   - 個人冷卻 18 秒；全服廣播反彈開始/每次反彈/結束
	"T178": {ID: "T178", Name: "幸運反彈魚", Type: TargetTypeSpecial, MultiplierMin: 30, MultiplierMax: 55, HP: 65, SpawnWeight: 4, Speed: 52, Lifetime: 13, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "lucky_ricochet"},
	// T179 幸運黑洞魚（DAY-221）— 業界原創「重力黑洞」機制
	// 設計：擊破 T179 後在場上建立「重力黑洞」（持續 10 秒）：
	//   - 黑洞建立在場景中央附近（隨機偏移），半徑 350px
	//   - 黑洞範圍內所有目標每 1 秒被「吸引」（HP -10%，模擬重力傷害）
	//   - 黑洞範圍內目標被擊破：獎勵 ×2.0 倍率加成（乘法）
	//   - 10 秒後「奇點爆炸」：黑洞範圍內所有目標 85% 擊破機率（0.70x 倍率，全服共享）
	//   - 個人冷卻 22 秒；全服冷卻 35 秒
	"T179": {ID: "T179", Name: "幸運黑洞魚", Type: TargetTypeSpecial, MultiplierMin: 40, MultiplierMax: 70, HP: 80, SpawnWeight: 3, Speed: 45, Lifetime: 15, LaborGain: 14, DifficultyFactor: 16.0, SpecialBehavior: "lucky_black_hole"},
	// T180 幸運共鳴魚（DAY-222）— 業界原創「全服共鳴」機制
	// 設計：擊破 T180 後觸發「共鳴模式」（15 秒）：
	//   - 全服所有玩家的每次射擊都累積「共鳴能量」（+1/shot）
	//   - 共鳴能量達到 30 點 → 觸發「共鳴爆發」：全場 HP -50% + 全服 ×1.8 倍率加成 6 秒
	//   - 共鳴爆發獎勵按「貢獻比例」分配（射擊越多，分到越多）
	//   - 15 秒內未達到 30 點 → 觸發「小型共鳴」：全場 HP -25% + 全服 ×1.3 倍率加成 3 秒
	//   - 全服冷卻 40 秒
	"T180": {ID: "T180", Name: "幸運共鳴魚", Type: TargetTypeSpecial, MultiplierMin: 35, MultiplierMax: 65, HP: 75, SpawnWeight: 3, Speed: 48, Lifetime: 15, LaborGain: 13, DifficultyFactor: 16.0, SpecialBehavior: "lucky_resonance"},
	// T181 幸運傳送魚（DAY-223）— 業界原創「傳送混亂」機制
	// 設計：擊破 T181 後觸發「傳送漩渦」（10 秒）：
	//   - 場上所有目標物立即隨機傳送到新位置（瞬間移動）
	//   - 傳送後 3 秒內擊破任何目標：獎勵 ×2.5 倍率加成（「傳送混亂」加成）
	//   - 每 3 秒再次傳送（最多 4 次傳送）
	//   - 個人冷卻 20 秒
	"T181": {ID: "T181", Name: "幸運傳送魚", Type: TargetTypeSpecial, MultiplierMin: 30, MultiplierMax: 60, HP: 70, SpawnWeight: 3, Speed: 50, Lifetime: 14, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "lucky_teleport"},
	// T182 幸運分裂魚（DAY-224）— 業界原創「一魚分三」機制
	// 設計：擊破 T182 後觸發「分裂爆炸」：
	//   - T182 分裂成 3 個「分裂碎片」（HP = 原 HP × 30%，倍率 ×1.8）
	//   - 分裂碎片在場上存活 8 秒，被擊破獲得 ×1.8 倍率加成（乘法）
	//   - 8 秒後所有未被擊破的分裂碎片「二次爆炸」（65% 擊破機率，0.60x 倍率）
	//   - 個人冷卻 18 秒
	"T182": {ID: "T182", Name: "幸運分裂魚", Type: TargetTypeSpecial, MultiplierMin: 28, MultiplierMax: 55, HP: 65, SpawnWeight: 4, Speed: 52, Lifetime: 13, LaborGain: 11, DifficultyFactor: 16.0, SpecialBehavior: "lucky_split"},
	// T183 幸運充能魚（DAY-225）— 業界原創「射擊充能→爆發」機制
	// 設計：擊破 T183 後觸發「充能模式」（12 秒）：
	//   - 玩家的每次射擊都累積「充能值」（+1/shot）
	//   - 充能值達到 10 → 自動觸發「充能爆發」：下一次擊破獲得 ×5.0 倍率加成（一次性）
	//   - 充能爆發後重置，可再次累積（12 秒內可觸發多次）
	//   - 個人冷卻 22 秒
	"T183": {ID: "T183", Name: "幸運充能魚", Type: TargetTypeSpecial, MultiplierMin: 32, MultiplierMax: 62, HP: 72, SpawnWeight: 3, Speed: 49, Lifetime: 14, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "lucky_charge"},
	// T184 幸運鏈鎖爆炸魚（DAY-226）— 業界原創「連鎖爆炸」機制
	// 設計：擊破 T184 後，場上隨機 3 個目標被「引爆標記」：
	//   - 引爆標記目標被擊破後，立即引爆周圍 200px 內所有目標（60% 擊破機率，×1.5 倍率）
	//   - 被引爆的目標如果也有引爆標記，繼續連鎖（最多 3 層連鎖）
	//   - 連鎖爆炸獎勵給觸發者；個人冷卻 20 秒
	"T184": {ID: "T184", Name: "幸運鏈鎖爆炸魚", Type: TargetTypeSpecial, MultiplierMin: 30, MultiplierMax: 58, HP: 68, SpawnWeight: 3, Speed: 50, Lifetime: 14, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "lucky_chain_bomb"},
	// T185 幸運鏡像時空魚（DAY-227）— 業界原創「時間倒流」機制
	// 設計：擊破 T185 後觸發「時間倒流」（8 秒）：
	//   - 場上所有目標物 HP 回滿（MaxHP）
	//   - 時間倒流期間擊破任何目標獲得 ×2.0 倍率加成（乘法）
	//   - 8 秒後「時間崩潰」：所有目標 HP -40%
	//   - 個人冷卻 25 秒
	"T185": {ID: "T185", Name: "幸運鏡像時空魚", Type: TargetTypeSpecial, MultiplierMin: 35, MultiplierMax: 65, HP: 75, SpawnWeight: 3, Speed: 47, Lifetime: 15, LaborGain: 13, DifficultyFactor: 16.0, SpecialBehavior: "lucky_mirror_time"},
	// T186 幸運量子魚（DAY-228）— 業界原創「量子疊加態」機制
	// 設計：擊破 T186 後觸發「量子疊加」：
	//   - 場上隨機 4 個目標進入「量子態」（同時疊加高倍率 ×3.0 和低倍率 ×0.8）
	//   - 玩家「觀測」（射擊命中）量子態目標時，50% 機率坍縮為高倍率（×3.0），50% 機率坍縮為低倍率（×0.8）
	//   - 量子態持續 10 秒；10 秒後所有未被觀測的量子態目標「量子爆炸」（70% 擊破機率，倍率隨機 ×1.0-×4.0）
	//   - 個人冷卻 20 秒；全服廣播量子態建立/坍縮/爆炸
	"T186": {ID: "T186", Name: "幸運量子魚", Type: TargetTypeSpecial, MultiplierMin: 32, MultiplierMax: 60, HP: 70, SpawnWeight: 3, Speed: 48, Lifetime: 14, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "lucky_quantum"},
	// T187 幸運寄生魚（DAY-229）— 業界原創「寄生附著+跳躍」機制
	// 設計：擊破 T187 後觸發「寄生釋放」：
	//   - 場上隨機 3 個目標被「寄生蟲附著」（綠色標記）
	//   - 寄生目標每 2 秒自動損失 HP（-8%/次，最多 5 次）
	//   - 寄生目標被擊破時，寄生蟲「跳躍」到最近的目標繼續寄生（最多跳躍 2 次）
	//   - 玩家擊破寄生目標獲得 ×2.2 倍率加成（乘法）
	//   - 個人冷卻 22 秒；全服廣播寄生附著/跳躍/消散
	"T187": {ID: "T187", Name: "幸運寄生魚", Type: TargetTypeSpecial, MultiplierMin: 30, MultiplierMax: 58, HP: 68, SpawnWeight: 3, Speed: 50, Lifetime: 14, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "lucky_parasite"},
	// T188 幸運風暴魚（DAY-230）— 業界原創「風暴旋轉+位置混亂」機制
	// 設計：擊破 T188 後在場上建立「風暴中心」（持續 10 秒）：
	//   - 風暴範圍（半徑 320px）內所有目標每 1.5 秒被「風暴旋轉」（隨機傳送到範圍內新位置）
	//   - 風暴範圍內目標被擊破：獎勵 ×2.5 倍率加成（乘法）
	//   - 10 秒後「風暴爆發」：範圍內所有目標 80% 擊破機率（0.75x 倍率，全服共享）
	//   - 個人冷卻 22 秒；全服冷卻 35 秒
	"T188": {ID: "T188", Name: "幸運風暴魚", Type: TargetTypeSpecial, MultiplierMin: 32, MultiplierMax: 60, HP: 70, SpawnWeight: 3, Speed: 50, Lifetime: 14, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "lucky_storm"},
	// T189 幸運迴旋鏢魚（DAY-231）— 業界原創「迴旋鏢來回穿透」機制
	// 設計：擊破 T189 後觸發「迴旋鏢模式」（10 秒）：
	//   - 玩家的每次射擊發射「迴旋鏢子彈」，命中目標後折返，最多來回 3 次
	//   - 每次命中：70% 擊破機率，0.65x 倍率（個人獎勵）
	//   - 個人冷卻 18 秒
	"T189": {ID: "T189", Name: "幸運迴旋鏢魚", Type: TargetTypeSpecial, MultiplierMin: 28, MultiplierMax: 55, HP: 65, SpawnWeight: 4, Speed: 52, Lifetime: 13, LaborGain: 11, DifficultyFactor: 16.0, SpecialBehavior: "lucky_boomerang"},
	// T190 幸運磁力魚（DAY-232）— 業界原創「磁力聚集+磁力爆發」機制
	// 設計：擊破 T190 後觸發「磁力場」（12 秒）：
	//   - 場上所有目標物被「磁力吸引」，每 1.5 秒向場景中央移動（聚集效果）
	//   - 磁力場期間擊破任何目標獲得 ×1.8 倍率加成（乘法）
	//   - 12 秒後「磁力爆發」：中央區域（半徑 200px）所有目標 75% 擊破機率（0.80x 倍率，全服共享）
	//   - 個人冷卻 20 秒；全服冷卻 30 秒
	"T190": {ID: "T190", Name: "幸運磁力魚", Type: TargetTypeSpecial, MultiplierMin: 30, MultiplierMax: 58, HP: 68, SpawnWeight: 3, Speed: 50, Lifetime: 14, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "lucky_magnet"},
	// T191 幸運回聲魚（DAY-233）— 業界原創「回聲分身+層疊倍率」機制
	//   - 擊破 T191 後，玩家的下一次擊破會產生「回聲分身」
	//   - 分身 HP = 原 HP × 50%，倍率遞增：第1層 ×1.5 → 第2層 ×2.0 → 第3層 ×2.5
	//   - 最多 3 層連鎖回聲；個人冷卻 18 秒
	"T191": {ID: "T191", Name: "幸運回聲魚", Type: TargetTypeSpecial, MultiplierMin: 28, MultiplierMax: 55, HP: 65, SpawnWeight: 4, Speed: 52, Lifetime: 13, LaborGain: 11, DifficultyFactor: 16.0, SpecialBehavior: "lucky_echo"},
	// T192 幸運漩渦魚（DAY-234）— 業界原創「漩渦旋轉+反向射擊」機制
	//   - 擊破 T192 後觸發「漩渦模式」（8 秒）：場景中央半徑 300px 內目標每 2 秒繞中心旋轉 45 度
	//   - 漩渦模式期間擊破任何目標獲得 ×2.2 倍率加成（乘法）
	//   - 8 秒後「漩渦爆發」：漩渦範圍內所有目標 70% 擊破機率（0.75x 倍率，全服共享）
	//   - 個人冷卻 20 秒；全服冷卻 30 秒
	"T192": {ID: "T192", Name: "幸運漩渦魚", Type: TargetTypeSpecial, MultiplierMin: 32, MultiplierMax: 60, HP: 70, SpawnWeight: 3, Speed: 50, Lifetime: 14, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "lucky_vortex"},
	// T193 幸運時間炸彈魚（DAY-235）— 業界原創「倒數計時+提前引爆+連鎖爆炸」機制
	//   - 擊破 T193 後，場上隨機 4 個目標被「時間炸彈標記」（倒數 8 秒）
	//   - 倒數結束時自動爆炸（80% 擊破機率，×1.6 倍率，個人獎勵）
	//   - 玩家可以「提前引爆」（射擊命中炸彈目標）：立即爆炸 + 引爆周圍 150px 內目標（60% 機率，×1.2 倍率）
	//   - 提前引爆的目標獲得 ×2.0 倍率加成（比等待爆炸更高）
	//   - 個人冷卻 20 秒
	"T193": {ID: "T193", Name: "幸運時間炸彈魚", Type: TargetTypeSpecial, MultiplierMin: 30, MultiplierMax: 58, HP: 68, SpawnWeight: 3, Speed: 50, Lifetime: 14, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "lucky_time_bomb"},
	// T194 幸運鏡面世界魚（DAY-236）— 業界原創「全場鏡像反轉+鏡面崩潰」機制
	//   - 擊破 T194 後觸發「鏡面世界」（10 秒）：場上所有目標 X 座標以場景中央為軸鏡像反轉
	//   - 鏡面世界期間擊破任何目標獲得 ×2.3 倍率加成（乘法）
	//   - 10 秒後「鏡面崩潰」：所有目標 HP -35%（保留最少 1）
	//   - 個人冷卻 22 秒；全服冷卻 35 秒
	"T194": {ID: "T194", Name: "幸運鏡面世界魚", Type: TargetTypeSpecial, MultiplierMin: 32, MultiplierMax: 60, HP: 70, SpawnWeight: 3, Speed: 50, Lifetime: 14, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "lucky_mirror_world"},
	// T195 幸運冰凍世界魚（DAY-237）— 業界原創「全場冰凍+冰裂爆發」機制
	//   - 擊破 T195 後觸發「冰凍世界」（8 秒）：場上所有目標移動速度降低 80%
	//   - 冰凍世界期間擊破任何目標獲得 ×2.0 倍率加成（乘法）
	//   - 8 秒後「冰裂爆發」：所有目標 HP -50%（保留最少 1），速度恢復
	//   - 個人冷卻 20 秒；全服冷卻 30 秒
	"T195": {ID: "T195", Name: "幸運冰凍世界魚", Type: TargetTypeSpecial, MultiplierMin: 30, MultiplierMax: 58, HP: 68, SpawnWeight: 3, Speed: 50, Lifetime: 14, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "lucky_freeze_world"},
	// T196 幸運重力反轉魚（DAY-238）— 業界原創「重力反轉+上下顛倒移動+重力崩潰」機制
	//   - 擊破 T196 後觸發「重力反轉」（10 秒）：場上所有目標 Y 座標以場景中央（Y=300）為軸翻轉
	//   - 重力反轉期間擊破任何目標獲得 ×2.1 倍率加成（乘法）
	//   - 10 秒後「重力崩潰」：所有目標 HP -45%（保留最少 1），Y 座標恢復
	//   - 個人冷卻 22 秒；全服冷卻 32 秒
	"T196": {ID: "T196", Name: "幸運重力反轉魚", Type: TargetTypeSpecial, MultiplierMin: 31, MultiplierMax: 59, HP: 69, SpawnWeight: 3, Speed: 50, Lifetime: 14, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "lucky_gravity_flip"},
	// T197 幸運共鳴爆發魚（DAY-239）— 業界原創「多效疊加共鳴爆發」機制
	//   - 擊破 T197 後，偵測場上當前同時啟動的幸運效果數量
	//   - ≥2 個效果 → 共鳴爆發：所有效果倍率額外 ×1.5，持續 6 秒
	//   - 1 個效果 → 小型共鳴：該效果倍率 ×1.3，持續 4 秒
	//   - 0 個效果 → 基礎爆發：全場 HP -30%，個人 ×1.8 倍率加成 5 秒
	//   - 個人冷卻 25 秒；全服冷卻 40 秒
	"T197": {ID: "T197", Name: "幸運共鳴爆發魚", Type: TargetTypeSpecial, MultiplierMin: 33, MultiplierMax: 62, HP: 72, SpawnWeight: 3, Speed: 48, Lifetime: 14, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "lucky_synergy_burst"},
	// T198 幸運賭注魚（DAY-240）— 業界原創「玩家主動風險決策+賭注翻倍」機制
	//   - 擊破 T198 後，玩家面臨「賭注選擇」（10 秒決策時間）
	//   - 選擇 A（保守）：下一次擊破 ×2.0 倍率，100% 觸發
	//   - 選擇 B（激進）：下一次擊破 ×5.0 倍率，50% 觸發；失敗則 ×0.5 倍率
	//   - 選擇 C（瘋狂）：下一次擊破 ×10.0 倍率，25% 觸發；失敗則 ×0.3 倍率
	//   - 10 秒內未選擇 → 自動選擇 A；個人冷卻 30 秒
	"T198": {ID: "T198", Name: "幸運賭注魚", Type: TargetTypeSpecial, MultiplierMin: 35, MultiplierMax: 65, HP: 75, SpawnWeight: 3, Speed: 47, Lifetime: 14, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "lucky_bet"},
	// T199 幸運連鎖反應魚（DAY-241）— 業界原創「多米諾骨牌效應」機制
	//   - 擊破 T199 後，場上隨機選 1 個目標作為「連鎖起點」（標記持續 15 秒）
	//   - 玩家擊破連鎖起點後，自動引爆距離最近的目標（100% 擊破，×1.4 倍率）
	//   - 被引爆的目標再引爆下一個最近目標（×1.3 倍率）
	//   - 連鎖最多 8 層，每層倍率遞減 0.1（×1.4 → ×1.3 → ... → ×0.7）
	//   - 每層引爆間隔 400ms，製造「多米諾骨牌」的視覺爽感
	//   - 個人冷卻 25 秒
	"T199": {ID: "T199", Name: "幸運連鎖反應魚", Type: TargetTypeSpecial, MultiplierMin: 36, MultiplierMax: 67, HP: 76, SpawnWeight: 3, Speed: 46, Lifetime: 14, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "lucky_chain_reaction"},
	// T200 幸運分身魚（DAY-242）— 業界原創「三方向同時射擊」機制
	//   - 擊破 T200 後觸發「分身模式」（8 秒）：玩家的每次射擊同時產生 2 個「分身子彈」
	//   - 分身子彈分別向左右各偏移 30 度飛出，命中目標：60% 擊破機率，×0.7 倍率（個人獎勵）
	//   - 分身子彈搜尋範圍：偏移方向 300px 內最近目標
	//   - 個人冷卻 20 秒
	"T200": {ID: "T200", Name: "幸運分身魚", Type: TargetTypeSpecial, MultiplierMin: 37, MultiplierMax: 68, HP: 77, SpawnWeight: 3, Speed: 45, Lifetime: 14, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "lucky_clone"},
	// T201 幸運預言魚（DAY-243）— 業界原創「預言指定目標」機制
	//   - 擊破 T201 後，Server 隨機「預言」場上 1 個目標（標記持續 12 秒）
	//   - 玩家在 12 秒內擊破預言目標 → 獲得 ×3.5 倍率加成（「預言成真」）
	//   - 若預言目標自然消失 → 自動「預言轉移」到下一個目標（最多轉移 2 次）
	//   - 若 12 秒後仍未擊破 → 「預言失敗」，全場 HP -20%（安慰獎）
	//   - 個人冷卻 20 秒
	"T201": {ID: "T201", Name: "幸運預言魚", Type: TargetTypeSpecial, MultiplierMin: 38, MultiplierMax: 69, HP: 78, SpawnWeight: 3, Speed: 44, Lifetime: 14, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "lucky_prophecy"},
	// T202 幸運奪旗魚（DAY-244）— 業界原創「全服搶旗競爭」機制
	//   - 擊破 T202 後，場上最高倍率目標被「旗幟標記」（持續 15 秒）
	//   - 所有玩家射擊旗幟目標，每次命中累積「搶旗積分」（+1/命中）
	//   - 每 3 秒廣播即時排名；15 秒後積分最高者獲得 ×4.0 倍率加成
	//   - 第 2 名 ×2.0，第 3 名 ×1.5；無人命中 → 自動爆炸全服共享
	//   - 個人冷卻 25 秒；全服冷卻 40 秒
	"T202": {ID: "T202", Name: "幸運奪旗魚", Type: TargetTypeSpecial, MultiplierMin: 39, MultiplierMax: 70, HP: 79, SpawnWeight: 3, Speed: 43, Lifetime: 14, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "lucky_flag"},
	// T203 幸運幽靈魚（DAY-245）— 業界原創「幽靈殘影+死亡後復活攻擊」機制
	//   - 擊破 T203 後，玩家獲得「幽靈護盾」（12 秒）
	//   - 護盾期間，玩家每次擊破任何目標，目標留下「幽靈殘影」（持續 5 秒）
	//   - 幽靈殘影可被再次擊破（50% 機率，×1.5 倍率，個人獎勵）
	//   - 12 秒後「幽靈爆發」：所有場上幽靈殘影同時爆炸（100% 擊破，×2.0 倍率，個人獎勵）
	//   - 個人冷卻 22 秒；全服冷卻 35 秒
	"T203": {ID: "T203", Name: "幸運幽靈魚", Type: TargetTypeSpecial, MultiplierMin: 40, MultiplierMax: 72, HP: 80, SpawnWeight: 3, Speed: 42, Lifetime: 14, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "lucky_phantom"},
	// T204 幸運水晶球魚（DAY-246）— 業界原創「預測未來+命中率提升」機制
	//   - 擊破 T204 後，Server「預測」場上 3 個目標為「水晶預言目標」（持續 8 秒）
	//   - 玩家射擊水晶預言目標時，命中率提升至 100%（必中）
	//   - 每次必中擊破獲得 ×2.5 倍率加成（個人獎勵）
	//   - 8 秒後「水晶爆炸」：所有未擊破的水晶預言目標自動爆炸（×1.8 倍率，個人獎勵）
	//   - 個人冷卻 20 秒；全服冷卻 30 秒
	"T204": {ID: "T204", Name: "幸運水晶球魚", Type: TargetTypeSpecial, MultiplierMin: 41, MultiplierMax: 73, HP: 81, SpawnWeight: 3, Speed: 41, Lifetime: 14, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "lucky_crystal_ball"},
	// T205 幸運時光倒流魚（DAY-247）— 業界原創「時光倒流+過去擊破重現」機制
	//   - 擊破 T205 後，Server 重播玩家「過去 10 秒內」擊破的最多 5 個目標
	//   - 每個重播目標以 ×1.6 倍率給予個人獎勵（不需要再次射擊，直接結算）
	//   - 同時場上所有目標 HP 恢復到 60%（讓玩家有更多目標可打）
	//   - 重播動畫：每個目標間隔 400ms 依序「閃現→爆炸」，製造「時光倒流」的視覺感
	//   - 個人冷卻 25 秒；全服冷卻 40 秒
	"T205": {ID: "T205", Name: "幸運時光倒流魚", Type: TargetTypeSpecial, MultiplierMin: 42, MultiplierMax: 75, HP: 82, SpawnWeight: 3, Speed: 40, Lifetime: 14, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "lucky_time_rewind"},
	// T206 幸運龍捲風魚（DAY-248）— 業界原創「龍捲風吸引+螺旋爆發」機制
	//   - 擊破 T206 後，場景中央生成「龍捲風」（持續 12 秒）
	//   - 龍捲風每 2 秒「吸引」場上所有目標向中央螺旋移動（每次移動 80px，帶旋轉角度）
	//   - 龍捲風期間擊破任何目標獲得 ×2.2 倍率加成（乘法）
	//   - 12 秒後「龍捲風爆發」：中央 250px 範圍內所有目標 85% 擊破機率（×1.5 倍率，全服共享）
	//   - 個人冷卻 22 秒；全服冷卻 35 秒
	"T206": {ID: "T206", Name: "幸運龍捲風魚", Type: TargetTypeSpecial, MultiplierMin: 43, MultiplierMax: 77, HP: 83, SpawnWeight: 3, Speed: 39, Lifetime: 14, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "lucky_tornado"},
	// T207 幸運黑洞爆炸魚（DAY-249）— 業界原創「黑洞吸收+能量爆炸」機制
	//   - 擊破 T207 後，場景中央生成「黑洞」（持續 10 秒）
	//   - 黑洞每 1.5 秒「吸收」場上距離最近的目標（直接消滅，×1.2 倍率，個人獎勵）
	//   - 黑洞最多吸收 6 個目標，每吸收一個「能量充能 +1」
	//   - 10 秒後「黑洞爆炸」：能量值 × 場上目標數 × 0.8 倍率（全服共享）
	//   - 個人冷卻 20 秒；全服冷卻 30 秒
	"T207": {ID: "T207", Name: "幸運黑洞爆炸魚", Type: TargetTypeSpecial, MultiplierMin: 44, MultiplierMax: 79, HP: 84, SpawnWeight: 3, Speed: 38, Lifetime: 14, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "lucky_black_hole_explosion"},
	// T208 幸運鏡像分裂魚（DAY-250）— 業界原創「鏡像分裂+雙重目標」機制
	//   - 擊破 T208 後，場上隨機 4 個目標被「鏡像分裂」
	//   - 每個目標在其鏡像位置（X 軸對稱）生成一個「鏡像副本」
	//   - 鏡像副本 HP = 原目標 50%，倍率 = 原目標 × 0.6（個人獎勵）
	//   - 鏡像副本存活 15 秒，玩家擊破獲得個人獎勵
	//   - 15 秒後所有未擊破的鏡像副本「鏡像消融」：每個消融給全服 ×0.3 倍率共享獎勵
	//   - 個人冷卻 22 秒；全服冷卻 35 秒
	"T208": {ID: "T208", Name: "幸運鏡像分裂魚", Type: TargetTypeSpecial, MultiplierMin: 45, MultiplierMax: 81, HP: 85, SpawnWeight: 3, Speed: 37, Lifetime: 14, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "lucky_mirror_split"},
	// T209 幸運量子糾纏魚（DAY-251）— 業界原創「量子糾纏+同步爆炸+量子共鳴」機制
	//   - 擊破 T209 後，場上隨機 2 個目標被「量子糾纏」（持續 20 秒）
	//   - 任何玩家擊破其中一個 → 另一個立刻「同步爆炸」（×1.8 倍率，全服共享）
	//   - 若兩個在 1.5 秒內被不同玩家擊破 → 觸發「量子共鳴」：全服 ×3.5 倍率大獎
	//   - 20 秒後未擊破 → 「量子衰變」：兩個目標 HP -60%（安慰獎）
	//   - 個人冷卻 25 秒；全服冷卻 40 秒
	"T209": {ID: "T209", Name: "幸運量子糾纏魚", Type: TargetTypeSpecial, MultiplierMin: 46, MultiplierMax: 83, HP: 86, SpawnWeight: 3, Speed: 36, Lifetime: 14, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "lucky_quantum_entangle"},
	// T210 幸運武器進化魚（DAY-252）— 業界原創「武器進化+穿透+武器爆發」機制
	//   - 擊破 T210 後，玩家武器「進化」（持續 12 秒）
	//   - 等級 2：命中率 +30%，倍率 ×1.5（乘法）
	//   - 進化期間再次擊破 T210 → 等級 3：穿透效果，倍率 ×2.5
	//   - 進化結束時「武器爆發」：自動 3 連射（×1.2 倍率，個人獎勵）
	//   - 個人冷卻 18 秒；全服冷卻 25 秒
	"T210": {ID: "T210", Name: "幸運武器進化魚", Type: TargetTypeSpecial, MultiplierMin: 47, MultiplierMax: 85, HP: 87, SpawnWeight: 3, Speed: 35, Lifetime: 14, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "lucky_weapon_evo"},
	// T211 幸運星際隕石魚（DAY-253）— 業界原創「隕石雨+隨機轟炸+隕石連擊+最終隕石」機制
	//   - 擊破 T211 後，天空降下「隕石雨」（持續 8 秒）
	//   - 每 1 秒隨機轟炸場上 2 個目標（70% 擊破機率，×1.3 倍率，全服共享）
	//   - 若連續 3 次都命中同一個目標 → 「隕石連擊」：×3.0 倍率（全服大獎）
	//   - 8 秒後「最終隕石」：場上最高 HP 目標被 100% 擊破（×2.0 倍率，全服共享）
	//   - 個人冷卻 20 秒；全服冷卻 30 秒
	"T211": {ID: "T211", Name: "幸運星際隕石魚", Type: TargetTypeSpecial, MultiplierMin: 48, MultiplierMax: 87, HP: 88, SpawnWeight: 3, Speed: 34, Lifetime: 14, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "lucky_meteor_shower"},
	// T212 幸運龍王降臨魚（DAY-254）— 業界原創「龍王降臨+龍息攻擊+龍王護盾+龍王爆發」機制
	//   - 擊破 T212 後，「龍王降臨」（持續 15 秒）
	//   - 每 2 秒「龍息攻擊」：隨機選 3 個目標，80% 擊破機率，×1.4 倍率（全服共享）
	//   - 龍王降臨期間，觸發玩家獲得「龍王護盾」（下一次免費射擊）
	//   - 15 秒後「龍王爆發」：場上所有目標 HP -60%，觸發玩家獲得 ×3.0 倍率加成（個人，5 秒）
	//   - 個人冷卻 25 秒；全服冷卻 40 秒
	"T212": {ID: "T212", Name: "幸運龍王降臨魚", Type: TargetTypeSpecial, MultiplierMin: 49, MultiplierMax: 89, HP: 89, SpawnWeight: 3, Speed: 33, Lifetime: 14, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "lucky_dragon_king"},
	// T213 幸運時空裂縫魚（DAY-255）— 業界原創「時空裂縫+傳送吸入+裂縫崩塌」機制
	//   - 擊破 T213 後，場景中央出現「時空裂縫」（持續 18 秒）
	//   - 每 3 秒「裂縫吸入」：吸入距離裂縫最近的目標，傳送到隨機位置（×1.6 倍率，全服共享）
	//   - 最多吸入 5 個目標（達到上限後裂縫提前崩塌）
	//   - 18 秒後「裂縫崩塌」：場上所有目標 HP -50%，全服 AOE 獎勵（×2.5 倍率，全服共享）
	//   - 個人冷卻 22 秒；全服冷卻 35 秒
	"T213": {ID: "T213", Name: "幸運時空裂縫魚", Type: TargetTypeSpecial, MultiplierMin: 50, MultiplierMax: 91, HP: 90, SpawnWeight: 3, Speed: 32, Lifetime: 14, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "lucky_rift"},
	// T214 幸運全服充能魚（DAY-256）— 業界原創「全服共同充能→全服大爆發」機制
	//   - 擊破 T214 後，全服所有玩家共同累積「充能值」（每次任何玩家擊破任何目標 +1）
	//   - 充能值達到 20 時「全服大爆發」：全場所有目標 100% 擊破（×2.0 倍率，全服共享）
	//   - 若 30 秒內未達到 20 → 「充能失敗」：已累積充能值 × 0.5 倍率（安慰獎，全服共享）
	//   - 個人冷卻 30 秒；全服冷卻 50 秒
	"T214": {ID: "T214", Name: "幸運全服充能魚", Type: TargetTypeSpecial, MultiplierMin: 51, MultiplierMax: 93, HP: 91, SpawnWeight: 3, Speed: 31, Lifetime: 14, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "lucky_server_charge"},
	// T215 幸運公會戰魚（DAY-257）— 業界原創「全服分隊競爭→勝隊爆發」機制
	//   - 擊破 T215 後，全服玩家自動分成兩隊（紅隊/藍隊，依玩家 ID 奇偶分配）
	//   - 30 秒內競爭擊破數，每次擊破為己隊累積積分
	//   - 勝隊全員 ×2.5 倍率加成（5 秒）；敗隊 ×1.2 安慰獎；平局 ×1.8
	//   - 個人冷卻 35 秒；全服冷卻 55 秒
	"T215": {ID: "T215", Name: "幸運公會戰魚", Type: TargetTypeSpecial, MultiplierMin: 52, MultiplierMax: 95, HP: 92, SpawnWeight: 3, Speed: 30, Lifetime: 14, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "lucky_guild_war"},
	// T216 幸運閃電風暴魚（DAY-258）— 業界原創「閃電風暴+連鎖跳躍+超級閃電+全場電擊」機制
	//   - 擊破 T216 後，觸發「閃電風暴」（持續 12 秒）
	//   - 每 1.5 秒「閃電跳躍」：從隨機目標出發，連鎖跳躍到最近的 3 個目標（×1.3 倍率，全服共享）
	//   - 累計跳躍達到 5 跳 → 「超級閃電」：×3.0 倍率（全服大獎）
	//   - 12 秒後「閃電爆炸」：場上所有目標 HP -40%（全服共享）
	//   - 個人冷卻 20 秒；全服冷卻 32 秒
	"T216": {ID: "T216", Name: "幸運閃電風暴魚", Type: TargetTypeSpecial, MultiplierMin: 53, MultiplierMax: 97, HP: 93, SpawnWeight: 3, Speed: 29, Lifetime: 14, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "lucky_lightning_storm"},
	"T217": {ID: "T217", Name: "幸運星座命運魚", Type: TargetTypeSpecial, MultiplierMin: 54, MultiplierMax: 99, HP: 94, SpawnWeight: 3, Speed: 28, Lifetime: 14, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "lucky_zodiac_fate"},
	"T218": {ID: "T218", Name: "幸運寶藏獵人魚", Type: TargetTypeSpecial, MultiplierMin: 55, MultiplierMax: 101, HP: 95, SpawnWeight: 3, Speed: 27, Lifetime: 14, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "lucky_treasure_hunter"},
	"T219": {ID: "T219", Name: "幸運時間膠囊魚", Type: TargetTypeSpecial, MultiplierMin: 56, MultiplierMax: 103, HP: 96, SpawnWeight: 3, Speed: 26, Lifetime: 14, LaborGain: 12, DifficultyFactor: 16.0, SpecialBehavior: "lucky_time_capsule"},
	"T220": {ID: "T220", Name: "幸運累積大獎池魚", Type: TargetTypeSpecial, MultiplierMin: 57, MultiplierMax: 105, HP: 97, SpawnWeight: 2, Speed: 25, Lifetime: 16, LaborGain: 15, DifficultyFactor: 16.0, SpecialBehavior: "lucky_progressive_jackpot"},
	"T221": {ID: "T221", Name: "幸運元素融合魚", Type: TargetTypeSpecial, MultiplierMin: 58, MultiplierMax: 107, HP: 98, SpawnWeight: 2, Speed: 24, Lifetime: 16, LaborGain: 15, DifficultyFactor: 16.0, SpecialBehavior: "lucky_element_fusion"},
	"T222": {ID: "T222", Name: "幸運命運輪迴魚", Type: TargetTypeSpecial, MultiplierMin: 59, MultiplierMax: 109, HP: 99, SpawnWeight: 2, Speed: 23, Lifetime: 16, LaborGain: 15, DifficultyFactor: 16.0, SpecialBehavior: "lucky_karma_cycle"},
	// T223 幸運競速賽魚（DAY-265）— 業界原創「全服即時競速+排行榜爆發」機制
	//   - 擊破 T223 後，觸發「全服競速賽」（持續 30 秒）
	//   - 所有玩家競爭擊破數，每次擊破 +1 積分，每 5 秒廣播即時排行榜（前 3 名）
	//   - 結算：第 1 名 ×4.0、第 2 名 ×2.5、第 3 名 ×1.8、其他 ×1.2 安慰獎（5 秒加成）
	//   - 個人冷卻 40 秒；全服冷卻 60 秒
	// T224 幸運連鎖爆炸魚（DAY-266）— 業界原創「連鎖爆炸+空間擴散+三層引爆」機制
	//   - 擊破 T224 後，第 1 層隨機引爆 1 個目標（×2.0），200px 內 HP-50%，40% 機率二次引爆（×1.5）
	//   - 第 2 層：150px 內 HP-30%，25% 機率三次引爆（×1.2）；最多 3 層連鎖
	//   - 個人冷卻 22 秒；全服冷卻 35 秒
	"T224": {ID: "T224", Name: "幸運連鎖爆炸魚", Type: TargetTypeSpecial, MultiplierMin: 61, MultiplierMax: 113, HP: 101, SpawnWeight: 2, Speed: 21, Lifetime: 16, LaborGain: 15, DifficultyFactor: 16.0, SpecialBehavior: "lucky_chain_explosion"},

	// T225 幸運倍率疊加魚（DAY-267）— Fishing Fortune Multiplier Cascade 機制
	//   - 擊破 T225 後，觸發「倍率疊加模式」（持續 25 秒）
	//   - 玩家每次擊破任何目標，疊加倍率 +0.3x（從 1.0x 開始，最高 10.0x）
	//   - 每次擊破都用「當前疊加倍率」計算獎勵
	//   - 達到 10.0x 時觸發「倍率爆發」：最後一次擊破獲得 ×20.0 大獎（個人）
	//   - 25 秒後未達到 10.0x → 「倍率結算」：用最終疊加倍率計算最後一次擊破獎勵
	//   - 個人冷卻 32 秒；全服冷卻 50 秒
	"T225": {ID: "T225", Name: "幸運倍率疊加魚", Type: TargetTypeSpecial, MultiplierMin: 62, MultiplierMax: 115, HP: 102, SpawnWeight: 2, Speed: 20, Lifetime: 16, LaborGain: 15, DifficultyFactor: 16.0, SpecialBehavior: "lucky_multiplier_stack"},

	// T226 幸運倒數炸彈魚（DAY-268）— 業界原創「倒數充能+全服爆炸」機制
	//   - 擊破 T226 後，場上出現「倒數炸彈」（10 秒倒數）
	//   - 倒數期間，任何玩家每次擊破任何目標，炸彈充能 +1（最多 10 次）
	//   - 10 秒後炸彈爆炸：充能數 × ×1.5 倍率（全服共享 AOE）
	//   - 若充能達到 10 次，提前引爆：×3.0 倍率（全服大獎）
	//   - 個人冷卻 28 秒；全服冷卻 45 秒
	"T226": {ID: "T226", Name: "幸運倒數炸彈魚", Type: TargetTypeSpecial, MultiplierMin: 63, MultiplierMax: 117, HP: 103, SpawnWeight: 2, Speed: 19, Lifetime: 16, LaborGain: 15, DifficultyFactor: 16.0, SpecialBehavior: "lucky_countdown_bomb"},

	// T227 幸運輪盤魚（DAY-269）— Royal Fishing ChainLong King Wheel 個人輪盤版本
	//   - 擊破 T227 後，觸發「個人幸運輪盤」（5 個扇區：×2.0/×3.0/×5.0/×8.0/×0.5）
	//   - Server 隨機決定結果，廣播給觸發玩家
	//   - 結果倍率套用到觸發玩家接下來 8 秒的所有擊破（個人）
	//   - 個人冷卻 25 秒；全服冷卻 40 秒
	"T227": {ID: "T227", Name: "幸運輪盤魚", Type: TargetTypeSpecial, MultiplierMin: 64, MultiplierMax: 119, HP: 104, SpawnWeight: 2, Speed: 18, Lifetime: 16, LaborGain: 15, DifficultyFactor: 16.0, SpecialBehavior: "lucky_spin_wheel"},

	// T228 幸運鏡像對決魚（DAY-270）— 2026 年最熱門「PvP 鏡像對決」機制
	//   - 擊破 T228 後，觸發「鏡像對決」：Server 隨機選一個其他玩家作為對手
	//   - 雙方進入 15 秒對決期；對決期間，雙方每次擊破目標，對手也獲得相同獎勵的 50%（鏡像分享）
	//   - 15 秒後，擊破數多的玩家獲得「對決勝利」×2.0 加成（5 秒）
	//   - 擊破數少的玩家獲得「對決失敗」×1.2 安慰獎（5 秒）；平局：雙方各獲得 ×1.5 加成（5 秒）
	//   - 若無其他玩家，觸發「孤獨模式」：個人 ×1.5 加成 10 秒
	//   - 個人冷卻 30 秒；全服冷卻 50 秒
	"T228": {ID: "T228", Name: "幸運鏡像對決魚", Type: TargetTypeSpecial, MultiplierMin: 65, MultiplierMax: 121, HP: 105, SpawnWeight: 2, Speed: 17, Lifetime: 16, LaborGain: 15, DifficultyFactor: 16.0, SpecialBehavior: "lucky_mirror_duel"},

	// T229 幸運倍率重擲魚（DAY-271）— GONE Fishing 2026 最新「multiplier reroll」機制
	//   - 擊破 T229 後，觸發「倍率重擲」：Server 為觸發玩家的下一次擊破重擲倍率（最多 3 次，取最高值）
	//   - 每次重擲有 40% 機率提升倍率（×1.5 到 ×4.0 隨機）
	//   - 最終用最高倍率計算獎勵（個人）；至少保證 ×1.5
	//   - 個人冷卻 20 秒；全服冷卻 35 秒
	"T229": {ID: "T229", Name: "幸運倍率重擲魚", Type: TargetTypeSpecial, MultiplierMin: 66, MultiplierMax: 123, HP: 106, SpawnWeight: 2, Speed: 16, Lifetime: 16, LaborGain: 15, DifficultyFactor: 16.0, SpecialBehavior: "lucky_reroll"},

	// T230 幸運品質突變魚（DAY-272）— Fishing Frenzy Chapter 3 Quality Roll 系統 + Fisch Mutation 機制
	//   - 擊破 T230 後，觸發「品質突變」：Server 為觸發玩家的下一次擊破「品質突變」
	//   - Normal（40%）×1.0 / Rare（30%）×1.8 / Epic（18%）×3.5 / Legendary（9%）×6.0 / Mythic（3%）×10.0
	//   - 品質效果持續到下一次擊破（一次性）；Mythic 品質全服廣播
	//   - 個人冷卻 18 秒；全服冷卻 30 秒
	"T230": {ID: "T230", Name: "幸運品質突變魚", Type: TargetTypeSpecial, MultiplierMin: 67, MultiplierMax: 125, HP: 107, SpawnWeight: 2, Speed: 15, Lifetime: 16, LaborGain: 15, DifficultyFactor: 16.0, SpecialBehavior: "lucky_quality_mutation"},

	// T231 幸運共鳴波魚（DAY-273）— Royal Fishing / Jili 2026「連鎖閃電+群體攻擊」趨勢進化版
	//   - 擊破 T231 後，發出「共鳴波」（3 層同心圓，每層間隔 400ms）
	//   - 第 1 層（r=150px）：HP -20%，35% 機率引爆（×2.0）
	//   - 第 2 層（r=250px）：HP -15%，25% 機率引爆（×1.8）
	//   - 第 3 層（r=350px）：HP -10%，15% 機率引爆（×1.5）
	//   - 引爆數 ≥ 5 → 全服 ×1.5 加成 8 秒；個人冷卻 25 秒；全服冷卻 40 秒
	"T231": {ID: "T231", Name: "幸運共鳴波魚", Type: TargetTypeSpecial, MultiplierMin: 68, MultiplierMax: 127, HP: 108, SpawnWeight: 2, Speed: 14, Lifetime: 16, LaborGain: 15, DifficultyFactor: 16.0, SpecialBehavior: "lucky_resonance_wave"},

	// T232 幸運命運預言魚（DAY-274）— Lucky Fish by AbraCadabra（2026-05-16）crash mechanic 進化版
	//   - 擊破 T232 後，Server 預言「下一條被擊破的魚倍率門檻」（預言值 = 隨機 ×2.0 到 ×8.0）
	//   - 玩家在 20 秒內擊破任何目標：
	//     若實際倍率 ≥ 預言值 → 「預言成真」×3.0 加成（個人）
	//     若實際倍率 < 預言值 → 「預言落空」×1.2 安慰獎（個人）
	//   - 預言值越高，成真機率越低但更有挑戰感；個人冷卻 22 秒；全服冷卻 38 秒
	"T232": {ID: "T232", Name: "幸運命運預言魚", Type: TargetTypeSpecial, MultiplierMin: 69, MultiplierMax: 129, HP: 109, SpawnWeight: 2, Speed: 13, Lifetime: 16, LaborGain: 15, DifficultyFactor: 16.0, SpecialBehavior: "lucky_fortune_prophecy"},

	// T233 幸運幸運圖騰魚（DAY-275）— Fish It Luck Totem（2026）「全場幸運加成」機制進化版
	//   - 擊破 T233 後，場上出現「幸運圖騰」（持續 15 秒）
	//   - 圖騰期間：全服所有玩家每次擊破任何目標 → ×1.3 全服加成
	//   - 觸發玩家額外獲得 ×1.5 個人加成（疊加在全服加成上）
	//   - 15 秒後圖騰消失，廣播結算（總擊破數/總獎勵）
	//   - 個人冷卻 30 秒；全服冷卻 50 秒
	"T233": {ID: "T233", Name: "幸運幸運圖騰魚", Type: TargetTypeSpecial, MultiplierMin: 70, MultiplierMax: 131, HP: 110, SpawnWeight: 2, Speed: 12, Lifetime: 16, LaborGain: 15, DifficultyFactor: 16.0, SpecialBehavior: "lucky_luck_totem"},

	// T234 幸運黃金颶風魚（DAY-276）— Royal Fishing Jili 2026「AOE 旋風掃場」機制進化版
	//   - 擊破 T234 後，觸發「黃金颶風」（螺旋掃場，持續 6 秒）
	//   - 颶風以螺旋路徑掃過整個場地，路徑上所有目標 HP -30%
	//   - 颶風每掃過一個目標，觸發玩家獲得 ×1.5 倍率加成（累積，最高 ×8.0）
	//   - 6 秒後颶風結算：廣播掃過目標數/累積倍率/總獎勵
	//   - 個人冷卻 28 秒；全服冷卻 45 秒
	"T234": {ID: "T234", Name: "幸運黃金颶風魚", Type: TargetTypeSpecial, MultiplierMin: 71, MultiplierMax: 133, HP: 111, SpawnWeight: 2, Speed: 11, Lifetime: 16, LaborGain: 15, DifficultyFactor: 16.0, SpecialBehavior: "lucky_golden_hurricane"},

	// T235 幸運閃電錘魚（DAY-277）— Battle of Luck「Lucky Slammer」機制（2026）進化版
	//   - 擊破 T235 後，觸發「閃電錘」：瞬間選定場上 3-6 個目標
	//   - 對每個目標造成「閃電錘擊」（HP -60%，30% 機率直接擊破）
	//   - 每個被錘擊的目標，觸發玩家獲得 ×1.2 倍率加成（累積）
	//   - 被直接擊破的目標，額外給予 ×2.0 倍率獎勵（個人）
	//   - 全服廣播錘擊結果（錘擊數/擊破數/總倍率）
	//   - 個人冷卻 22 秒；全服冷卻 35 秒
	"T235": {ID: "T235", Name: "幸運閃電錘魚", Type: TargetTypeSpecial, MultiplierMin: 72, MultiplierMax: 135, HP: 112, SpawnWeight: 2, Speed: 10, Lifetime: 16, LaborGain: 15, DifficultyFactor: 16.0, SpecialBehavior: "lucky_lightning_hammer"},

	// T236 幸運時間裂縫魚（DAY-278）— 業界原創「時間裂縫+最高倍率重現+裂縫複製體」機制
	//   - 擊破 T236 後，Server 查找玩家過去 30 秒內「最高倍率的那次擊破記錄」
	//   - 立即給予觸發玩家 ×2.5 加成（個人，基於最高倍率目標的 bet × mult × 2.5）
	//   - 同時在場上生成一個「裂縫複製體」（同種類目標，HP 只有 30%，擊破給 ×3.0 大獎）
	//   - 若過去 30 秒無擊破記錄 → 給予 ×1.5 保底獎勵 + 生成隨機裂縫複製體
	//   - 全服廣播「時間裂縫重現了什麼目標/倍率」
	//   - 個人冷卻 20 秒；全服冷卻 32 秒
	"T236": {ID: "T236", Name: "幸運時間裂縫魚", Type: TargetTypeSpecial, MultiplierMin: 73, MultiplierMax: 137, HP: 113, SpawnWeight: 2, Speed: 9, Lifetime: 16, LaborGain: 15, DifficultyFactor: 16.0, SpecialBehavior: "lucky_time_rift"},

	"B001": {ID: "B001", Name: "那個孩子", Type: TargetTypeBoss, MultiplierMin: 100, MultiplierMax: 500, HP: 3000, SpawnWeight: 0, Speed: 20, Lifetime: 60, LaborGain: 30, DifficultyFactor: 16.0, SpecialBehavior: "boss_phases"},
}

// MeteorMultiplierWeights 流星倍率權重（規格書 26.3）
var MeteorMultiplierWeights = []struct {
	Multiplier float64
	Weight     int
}{
	{20, 50},
	{30, 30},
	{40, 15},
	{50, 5},
}

// Characters 角色定義（規格書 5章）
var Characters = map[string]*CharacterDef{
	"chiikawa": {ID: "chiikawa", Name: "吉伊卡哇", BetLevelMin: 1, BetLevelMax: 3, AttackColor: "pink", KillModifier: 1.00, FireRateModifier: 1.00, LaborModifier: 1.10, VoiceText: "YaDa"},
	"hachiware": {ID: "hachiware", Name: "小八", BetLevelMin: 4, BetLevelMax: 7, AttackColor: "blue", KillModifier: 1.00, FireRateModifier: 1.08, LaborModifier: 1.00, VoiceText: "尖尖哇嘎乃"},
	"usagi": {ID: "usagi", Name: "烏薩奇", BetLevelMin: 8, BetLevelMax: 10, AttackColor: "yellow", KillModifier: 0.98, FireRateModifier: 1.20, LaborModifier: 0.95, VoiceText: "Yaha"},
}

// BetLevels 投注等級表（規格書 6章 & 25.5）
var BetLevels = []*BetDef{
	{Level: 1, CharacterID: "chiikawa", BetCost: 1, AttackPower: 1, FireRate: 2.0, ProjectileSpeed: 700},
	{Level: 2, CharacterID: "chiikawa", BetCost: 2, AttackPower: 2, FireRate: 2.0, ProjectileSpeed: 720},
	{Level: 3, CharacterID: "chiikawa", BetCost: 3, AttackPower: 3, FireRate: 2.1, ProjectileSpeed: 740},
	{Level: 4, CharacterID: "hachiware", BetCost: 5, AttackPower: 5, FireRate: 2.2, ProjectileSpeed: 780},
	{Level: 5, CharacterID: "hachiware", BetCost: 10, AttackPower: 10, FireRate: 2.3, ProjectileSpeed: 800},
	{Level: 6, CharacterID: "hachiware", BetCost: 20, AttackPower: 20, FireRate: 2.4, ProjectileSpeed: 820},
	{Level: 7, CharacterID: "hachiware", BetCost: 30, AttackPower: 30, FireRate: 2.5, ProjectileSpeed: 850},
	{Level: 8, CharacterID: "usagi", BetCost: 50, AttackPower: 50, FireRate: 2.7, ProjectileSpeed: 900},
	{Level: 9, CharacterID: "usagi", BetCost: 80, AttackPower: 80, FireRate: 2.9, ProjectileSpeed: 940},
	{Level: 10, CharacterID: "usagi", BetCost: 100, AttackPower: 100, FireRate: 3.0, ProjectileSpeed: 980},
}

// BonusTargets Bonus Game 目標（規格書 29.3）
var BonusTargets = []*BonusTargetDef{
	{ID: "BG001", Name: "普通雜草", ClickScore: 1, SpawnWeight: 180, SpecialEffect: "none"},
	{ID: "BG002", Name: "硬雜草", ClickScore: 3, SpawnWeight: 80, SpecialEffect: "double_click"},
	{ID: "BG003", Name: "發光雜草", ClickScore: 8, SpawnWeight: 35, SpecialEffect: "multiplier_up"},
	{ID: "BG004", Name: "金色雜草", ClickScore: 20, SpawnWeight: 10, SpecialEffect: "coin_shower"},
	{ID: "BG005", Name: "搗亂怪草", ClickScore: -5, SpawnWeight: 20, SpecialEffect: "stun"},
}

// GetBetDef 取得投注等級定義
func GetBetDef(level int) *BetDef {
	if level < 1 || level > 10 {
		return BetLevels[0]
	}
	return BetLevels[level-1]
}

// GetCharacterByBetLevel 依投注等級取得角色
func GetCharacterByBetLevel(level int) *CharacterDef {
	bet := GetBetDef(level)
	return Characters[bet.CharacterID]
}

// BaseRTPFactor 基礎 RTP 係數（規格書 30章）
const BaseRTPFactor = 0.92

// LaborValueMax 勞動值上限
const LaborValueMax = 100

// SpawnInterval 目標生成間隔（秒）
const SpawnInterval = 0.8

// MaxTargetsOnScreen 畫面最大目標數
const MaxTargetsOnScreen = 18

// BossDuration BOSS 持續時間（秒）
const BossDuration = 60.0

// BonusDuration Bonus Game 持續時間（秒）
const BonusDuration = 15.0

// ---- 武器升級系統（DAY-067）----

// WeaponDef 武器定義
type WeaponDef struct {
	Level       int
	Name        string
	Icon        string    // 顯示圖示
	PowerMod    float64   // 攻擊力加成係數（1.0=無加成，1.25=+25%）
	ExtraCost   int       // 每次攻擊額外扣除的金幣（在 BetCost 之外）
	Color       string    // 投射物顏色（Client 端視覺）
	Description string
}

// Weapons 武器等級定義（DAY-067）
var Weapons = []*WeaponDef{
	{
		Level:       1,
		Name:        "標準砲",
		Icon:        "🔫",
		PowerMod:    1.00,
		ExtraCost:   0,
		Color:       "white",
		Description: "標準攻擊力，無額外費用",
	},
	{
		Level:       2,
		Name:        "強化砲",
		Icon:        "⚡",
		PowerMod:    1.25,
		ExtraCost:   50,  // 每次攻擊額外扣 50 金幣
		Color:       "cyan",
		Description: "攻擊力 +25%，每次攻擊額外消耗 50 金幣",
	},
	{
		Level:       3,
		Name:        "超級砲",
		Icon:        "🌟",
		PowerMod:    1.60,
		ExtraCost:   150, // 每次攻擊額外扣 150 金幣
		Color:       "gold",
		Description: "攻擊力 +60%，每次攻擊額外消耗 150 金幣",
	},
}

// GetWeaponDef 取得武器定義
func GetWeaponDef(level int) *WeaponDef {
	if level < 1 || level > len(Weapons) {
		return Weapons[0]
	}
	return Weapons[level-1]
}
