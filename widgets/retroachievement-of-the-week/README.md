# RetroAchievement of the Week

![RomM Widget Preview](https://github.com/milktoastrat/community-widgets/blob/main/widgets/retroachievement-of-the-week/retroachievement-of-the-week-example.png?raw=true)

A custom Glance widget that displays the current [RetroAchievements.org](https://retroachievements.org) "Achievement of the Week" — with full support for badge, game title, description, and links to view the achievement and game on RA.


```yaml
- type: custom-api
  title: RetroAchievement of the Week
  cache: 1h
  url: https://retroachievements.org/API/API_GetAchievementOfTheWeek.php?y=${RA_API}
  template: |
    <div style="display: flex; align-items: flex-start;">
      <div style="flex-shrink: 0; margin-right: 1em;">
        <a href="https://retroachievements.org/achievement/{{ .JSON.Int "Achievement.ID" }}" target="_blank" rel="noopener noreferrer">
          <img src="https://retroachievements.org{{ .JSON.String "Achievement.BadgeURL" }}" alt="Badge" style="width: 64px; height: 64px; border-radius: var(--border-radius);">
        </a>
      </div>
      <div>
        <div style="font-weight: bold; color: var(--color-accent);">
          <a href="https://retroachievements.org/achievement/{{ .JSON.Int "Achievement.ID" }}" style="color: inherit; text-decoration: none;" target="_blank" rel="noopener noreferrer">
            {{ .JSON.String "Achievement.Title" }}
          </a>
        </div>
        <div style="color: var(--color-subdue);">
          {{ .JSON.String "Achievement.Description" }}
        </div>
        <div style="margin-top: 0.5em;">
          <a href="https://retroachievements.org/game/{{ .JSON.Int "Game.ID" }}" style="color: var(--color-primary); font-family: monospace; font-weight: bold; text-decoration: none;" target="_blank" rel="noopener noreferrer">
            {{ .JSON.String "Game.Title" }}
          </a>
          <span style="color: var(--color-subdue); font-style: italic;">
            ({{ .JSON.String "Console.Title" }})
          </span>
        </div>
      </div>
    </div>
```
### Setup

1. **Get an API key** from [RetroAchievements Settings](https://retroachievements.org/settings) or your RA account settings.
2. Add a secret in your Glance `.env` file:

   ```env
   RA_API=your_api_key_here
