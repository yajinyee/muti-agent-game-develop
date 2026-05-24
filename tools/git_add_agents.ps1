$env:GIT_TMPDIR = "d:\Kiro\.git\tmp"
$files = @(
    "agents/character-pixel-agent.md",
    "agents/environment-agent.md",
    "agents/game-state-agent.md",
    "agents/hit-effect-agent.md",
    "agents/hud-core-agent.md",
    "agents/lucky-panel-agent.md",
    "agents/network-agent.md",
    "agents/performance-agent.md",
    "agents/screen-effect-agent.md",
    "agents/screen-recorder-agent.md",
    "agents/server-combat-agent.md",
    "agents/server-core-agent.md",
    "agents/server-event-agent.md",
    "agents/server-infra-agent.md",
    "agents/sfx-agent.md",
    "agents/social-ui-agent.md",
    "agents/target-ai-agent.md",
    "agents/target-design-agent.md",
    "agents/target-pixel-agent.md",
    "agents/target-system-agent.md",
    "agents/ui-art-agent.md"
)
foreach ($f in $files) {
    Write-Host "Adding $f..."
    git add $f
    if ($LASTEXITCODE -eq 0) {
        git commit -m "agent v3.0 $f"
    } else {
        Write-Host "FAILED: $f"
    }
}
Write-Host "Done"
