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
