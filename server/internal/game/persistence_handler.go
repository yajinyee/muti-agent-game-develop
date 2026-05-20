// persistence_handler.go — 完整玩家資料持久化（DAY-098）
// 集中管理所有子系統的儲存和恢復邏輯
// 支援：基礎資料 / VIP / 賽季 / 圖鑑 / 統計
package game

import (
	"log"
	"time"

	"digital-twin/server/internal/game/codex"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/store"
)

// saveFullPlayerState 儲存玩家完整狀態到 FileStore（DAY-098）
// 在 RemovePlayer 時呼叫，也可定期呼叫
func (g *Game) saveFullPlayerState(p *player.Player) {
	fs, ok := g.store.(*store.FileStore)
	if !ok {
		// 非 FileStore（Redis 或 Memory），使用舊的基礎儲存
		state := &store.PlayerState{
			PlayerID:       p.ID,
			DisplayName:    p.DisplayName,
			Coins:          int64(p.Coins),
			Labor:          p.LaborValue,
			BetLevel:       p.BetLevel,
			SessionScore:   int64(p.SessionScore),
			MaxCoins:       int64(p.MaxCoins),
			KillCount:      p.KillCount,
			RoomID:         g.ID,
			LoginStreak:    p.LoginStreak,
			MaxLoginStreak: p.MaxLoginStreak,
			LastLoginDate:  p.LastLoginDate,
			EquippedSkin:   p.EquippedSkin,
			OwnedSkins:     p.OwnedSkins,
		}
		if err := g.store.SavePlayer(state); err != nil {
			log.Printf("[Persist] Failed to save player %s: %v", p.ID, err)
		}
		return
	}

	full := &store.FullPlayerState{
		// 基礎資料
		PlayerID:       p.ID,
		DisplayName:    p.DisplayName,
		Coins:          p.Coins,
		MaxCoins:       p.MaxCoins,
		BetLevel:       p.BetLevel,
		WeaponLevel:    p.WeaponLevel,
		KillCount:      p.KillCount,
		TotalBet:       p.TotalBet,
		TotalReward:    p.TotalReward,
		AttackCount:    p.AttackCount,
		SessionScore:   p.SessionScore,
		RoomDifficulty: p.RoomDifficulty,
		LastSeen:       time.Now(),
		// 登入資訊
		LastLoginDate:  p.LastLoginDate,
		LoginStreak:    p.LoginStreak,
		MaxLoginStreak: p.MaxLoginStreak,
		// 砲台外觀
		EquippedSkin: p.EquippedSkin,
		OwnedSkins:   p.OwnedSkins,
	}

	// VIP 系統（DAY-078）
	if g.VIP != nil {
		totalSpend, vipLevel, lastWeeklyAt := g.VIP.GetData(p.ID)
		full.VIPTotalSpend = totalSpend
		full.VIPLevel = vipLevel
		full.VIPLastWeeklyAt = lastWeeklyAt
	}

	// 賽季通行證（DAY-072）
	if g.Season != nil {
		points, level, claimed := g.Season.GetData(p.ID)
		full.SeasonPoints = points
		full.SeasonLevel = level
		full.SeasonClaimed = claimed
	}

	// 魚類圖鑑（DAY-081）
	if p.Codex != nil {
		entries := p.Codex.GetSnapshot()
		full.CodexEntries = make([]store.CodexEntryState, 0, len(entries))
		for _, e := range entries {
			full.CodexEntries = append(full.CodexEntries, store.CodexEntryState{
				TargetID:      e.TargetID,
				Unlocked:      e.Unlocked,
				UnlockedAt:    e.UnlockedAt,
				KillCount:     e.KillCount,
				MaxMultiplier: e.MaxMultiplier,
			})
		}
	}

	// 玩家統計（DAY-096）
	if p.Stats != nil {
		snap := p.Stats.Snapshot()
		full.StatsTotalSessions = snap.TotalSessions
		full.StatsTotalPlayTime = snap.TotalPlayTimeSec
		full.StatsTotalShots = snap.TotalShots
		full.StatsTotalKills = snap.TotalKills
		full.StatsTotalBet = snap.TotalBet
		full.StatsTotalReward = snap.TotalReward
		full.StatsTotalBonuses = snap.TotalBonuses
		full.StatsTotalBossKills = snap.TotalBossKills
		full.StatsBestMultiplier = snap.BestMultiplier
		full.StatsBestStreak = snap.BestStreak
		full.StatsBestSession = snap.BestSessionScore
		full.StatsBestBonus = snap.BestBonusReward
		full.StatsMaxCoins = snap.MaxCoins
		full.StatsJackpotWins = snap.JackpotWins
		full.StatsJackpotMini = snap.JackpotMiniWins
		full.StatsJackpotMinor = snap.JackpotMinorWins
		full.StatsJackpotMajor = snap.JackpotMajorWins
		full.StatsJackpotGrand = snap.JackpotGrandWins
		full.StatsJackpotPayout = snap.TotalJackpotPayout
		full.StatsHitCount = snap.TotalKills // HitCount = TotalKills（命中即擊破）
		full.StatsMissCount = snap.TotalShots - snap.TotalKills
		if snap.FirstPlayAt > 0 {
			full.StatsFirstPlayAt = time.UnixMilli(snap.FirstPlayAt)
		}
		if snap.LastPlayAt > 0 {
			full.StatsLastPlayAt = time.UnixMilli(snap.LastPlayAt)
		}
	}

	if err := fs.SaveFull(full); err != nil {
		log.Printf("[Persist] Failed to save full player %s: %v", p.ID, err)
	} else {
		log.Printf("[Persist] Player %s saved: coins=%d, vip=%d, season=%d, codex=%d entries",
			p.ID, full.Coins, full.VIPLevel, full.SeasonPoints, len(full.CodexEntries))
	}

	// 更新排行榜
	g.store.UpdateLeaderboard(p.ID, int64(p.SessionScore))
}

// restoreFullPlayerState 從 FileStore 恢復玩家完整狀態（DAY-098）
// 在 AddPlayer 時呼叫
func (g *Game) restoreFullPlayerState(p *player.Player) {
	fs, ok := g.store.(*store.FileStore)
	if !ok {
		// 非 FileStore，使用舊的基礎恢復
		if saved, err := g.store.LoadPlayer(p.ID); err == nil && saved != nil {
			p.Coins = int(saved.Coins)
			p.MaxCoins = int(saved.MaxCoins)
			p.KillCount = saved.KillCount
			if saved.BetLevel >= 1 && saved.BetLevel <= 10 {
				p.BetLevel = saved.BetLevel
			}
			if saved.DisplayName != "" {
				p.DisplayName = saved.DisplayName
			}
			p.LoginStreak = saved.LoginStreak
			p.MaxLoginStreak = saved.MaxLoginStreak
			p.LastLoginDate = saved.LastLoginDate
			if saved.EquippedSkin != "" {
				p.EquippedSkin = saved.EquippedSkin
			}
			if len(saved.OwnedSkins) > 0 {
				p.OwnedSkins = saved.OwnedSkins
			}
		}
		return
	}

	full, err := fs.LoadFull(p.ID)
	if err != nil {
		log.Printf("[Persist] Failed to load player %s: %v", p.ID, err)
		return
	}
	if full == nil {
		log.Printf("[Persist] New player %s, starting fresh", p.ID)
		return
	}

	// 恢復基礎資料
	p.Coins = full.Coins
	p.MaxCoins = full.MaxCoins
	if full.BetLevel >= 1 && full.BetLevel <= 10 {
		p.BetLevel = full.BetLevel
	}
	if full.WeaponLevel >= 1 && full.WeaponLevel <= 3 {
		p.WeaponLevel = full.WeaponLevel
	}
	p.KillCount = full.KillCount
	p.TotalBet = full.TotalBet
	p.TotalReward = full.TotalReward
	p.AttackCount = full.AttackCount
	if full.DisplayName != "" {
		p.DisplayName = full.DisplayName
	}
	if full.RoomDifficulty != "" {
		p.RoomDifficulty = full.RoomDifficulty
	}
	// 登入資訊
	p.LoginStreak = full.LoginStreak
	p.MaxLoginStreak = full.MaxLoginStreak
	p.LastLoginDate = full.LastLoginDate
	// 砲台外觀
	if full.EquippedSkin != "" {
		p.EquippedSkin = full.EquippedSkin
	}
	if len(full.OwnedSkins) > 0 {
		p.OwnedSkins = full.OwnedSkins
	}

	// 恢復 VIP 系統（DAY-078）
	if g.VIP != nil && full.VIPTotalSpend > 0 {
		g.VIP.LoadState(p.ID, full.VIPTotalSpend, full.VIPLevel, full.VIPLastWeeklyAt)
	}

	// 恢復賽季通行證（DAY-072）
	if g.Season != nil && full.SeasonPoints > 0 {
		g.Season.LoadState(p.ID, full.SeasonPoints, full.SeasonLevel, full.SeasonClaimed)
	}

	// 恢復魚類圖鑑（DAY-081）
	if p.Codex != nil && len(full.CodexEntries) > 0 {
		entries := make([]*codex.Entry, 0, len(full.CodexEntries))
		for _, e := range full.CodexEntries {
			entries = append(entries, &codex.Entry{
				TargetID:      e.TargetID,
				Unlocked:      e.Unlocked,
				UnlockedAt:    e.UnlockedAt,
				KillCount:     e.KillCount,
				MaxMultiplier: e.MaxMultiplier,
			})
		}
		p.Codex.LoadState(entries)
	}

	// 恢復玩家統計（DAY-096）
	if p.Stats != nil {
		p.Stats.LoadState(
			full.StatsTotalSessions, full.StatsTotalPlayTime,
			full.StatsTotalShots, full.StatsTotalKills,
			full.StatsTotalBet, full.StatsTotalReward,
			full.StatsTotalBonuses, full.StatsTotalBossKills,
			full.StatsBestMultiplier, full.StatsBestStreak,
			full.StatsBestSession, full.StatsBestBonus, full.StatsMaxCoins,
			full.StatsJackpotWins, full.StatsJackpotMini, full.StatsJackpotMinor,
			full.StatsJackpotMajor, full.StatsJackpotGrand, full.StatsJackpotPayout,
			full.StatsHitCount, full.StatsMissCount,
			full.StatsFirstPlayAt, full.StatsLastPlayAt,
		)
	}

	log.Printf("[Persist] Player %s restored: coins=%d, vip=%d, season=%d, codex=%d entries",
		p.ID, full.Coins, full.VIPLevel, full.SeasonPoints, len(full.CodexEntries))
}
