// streak_handler.go - DAY-083
package game

import (
"log"
"time"

"digital-twin/server/internal/player"
"digital-twin/server/internal/ws"
)

func (g *Game) notifyStreakKill(p *player.Player) float64 {
if p.Streak == nil {
return 1.0
}
currentStreak, multBonus, isNewLevel := p.Streak.RecordKill()
snap := p.Streak.GetSnapshot()
if err := g.Hub.Send(p.ID, &ws.Message{
Type: ws.MsgStreakUpdate,
Payload: ws.StreakUpdatePayload{
Current:    currentStreak,
MultBonus:  multBonus,
LevelName:  snap.LevelName,
LevelColor: snap.LevelColor,
IsNewLevel: isNewLevel,
MaxStreak:  snap.MaxStreak,
},
}); err != nil {
log.Printf("[Streak] send update error: %v", err)
}
if isNewLevel && currentStreak >= 3 {
log.Printf("[Streak] player=%s streak=%d (%s, mult=%.2f)", p.ID, currentStreak, snap.LevelName, multBonus)
}
return multBonus
}

func (g *Game) tickStreakTimeout() {
g.mu.RLock()
players := make([]*player.Player, 0, len(g.Players))
for _, p := range g.Players {
players = append(players, p)
}
g.mu.RUnlock()
for _, p := range players {
if p.Streak == nil {
continue
}
if reset := p.Streak.CheckTimeout(); reset {
snap := p.Streak.GetSnapshot()
if err := g.Hub.Send(p.ID, &ws.Message{
Type: ws.MsgStreakReset,
Payload: ws.StreakResetPayload{
FinalStreak: 0,
MaxStreak:   snap.MaxStreak,
},
}); err != nil {
log.Printf("[Streak] send reset error: %v", err)
}
}
}
}

func (g *Game) startStreakTicker() {
go func() {
ticker := time.NewTicker(1 * time.Second)
defer ticker.Stop()
for {
select {
case <-ticker.C:
g.tickStreakTimeout()
case <-g.stopCh:
return
}
}
}()
}