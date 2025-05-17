## Preview

![](preview.png)


## Configuration

```yaml
          - type: custom-api
            title: MLB Today
            url: "https://statsapi.mlb.com/api/v1/schedule/games/?sportId=1&hydrate=linescore(runners),team"
            cache: 15s
            template: |
              {{ $dates := .JSON.Array "dates" }}
              {{ if eq (len $dates) 0 }}
                <div style="padding:8px;text-align:center;">No games scheduled today.</div>
              {{ else }}
                {{ $games := (index $dates 0).Array "games" }}
                {{ $wid := printf "mlb-%d" now.UnixNano }}
                <div id="{{ $wid }}" style="position:relative;padding-bottom:{{ if gt (len $games) 6 }}25px{{ else }}0px{{ end }};">
                  <input type="checkbox" id="{{ $wid }}-cb" style="display:none"/>
                  <ul class="visible" style="list-style:none;padding:0;margin:0;">
                    {{ $count := 0 }}
                    {{ range $i,$g := $games }}
                      {{ if and (lt $count 6) (eq ($g.Get "status.abstractGameState").String "Live") }}
                        {{ template "row" $g }}
                        {{ $count = add $count 1 }}
                      {{ end }}
                    {{ end }}
                    
                    {{ range $i,$g := $games }}
                      {{ if and (lt $count 6) (eq ($g.Get "status.abstractGameState").String "Preview") }}
                        {{ template "row" $g }}
                        {{ $count = add $count 1 }}
                      {{ end }}
                    {{ end }}
                    
                    {{ range $i,$g := $games }}
                      {{ if and (lt $count 6) (not (or (eq ($g.Get "status.abstractGameState").String "Live") (eq ($g.Get "status.abstractGameState").String "Preview"))) }}
                        {{ template "row" $g }}
                        {{ $count = add $count 1 }}
                      {{ end }}
                    {{ end }}
                  </ul>
                  <ul class="all" style="list-style:none;padding:0;margin:0;display:none">

                    {{ range $i,$g := $games }}
                      {{ if eq ($g.Get "status.abstractGameState").String "Live" }}
                        {{ template "row" $g }}
                      {{ end }}
                    {{ end }}
                    
                    <!-- Then show PREVIEW games -->
                    {{ range $i,$g := $games }}
                      {{ if eq ($g.Get "status.abstractGameState").String "Preview" }}
                        {{ template "row" $g }}
                      {{ end }}
                    {{ end }}
                    
                    <!-- Finally show completed games -->
                    {{ range $i,$g := $games }}
                      {{ if not (or (eq ($g.Get "status.abstractGameState").String "Live") (eq ($g.Get "status.abstractGameState").String "Preview")) }}
                        {{ template "row" $g }}
                      {{ end }}
                    {{ end }}
                  </ul>
                  {{ if gt (len $games) 6 }}
                    <label for="{{ $wid }}-cb" style="position:absolute;bottom:0px;right:0;cursor:pointer;color:var(--glance-accent-color);font-size:16px;padding:5px;">
                      <span style="display:inline-block;transition:transform .2s">▼</span>
                    </label>
                  {{ end }}
                  <style>
                    #{{ $wid }}-cb:not(:checked) ~ .all { display:none!important }
                    #{{ $wid }}-cb:checked ~ .visible { display:none!important }
                    #{{ $wid }}-cb:checked ~ .all { display:block!important }
                    #{{ $wid }}-cb:checked ~ label span { transform:rotate(180deg)!important }
                  </style>
                </div>
              {{ end }}

              {{ define "row" }}
                {{ $awayComp := .Get "teams.away" }}
                {{ $homeComp := .Get "teams.home" }}
                {{ $status := .Get "status" }}
                {{ $lines := .Get "linescore" }}
                {{ $offense := $lines.Get "offense" }}
                {{ $r1 := ($offense.Get "first").Exists }}
                {{ $r2 := ($offense.Get "second").Exists }}
                {{ $r3 := ($offense.Get "third").Exists }}
                {{ $offTeamId := (($offense.Get "team").Get "id").Int }}
                {{ $abs := ($status.Get "abstractGameState").String }}
                {{ $live := eq $abs "Live" }}
                {{ $inningSt := ($lines.Get "inningState").String }}
                {{ $inningN := ($lines.Get "currentInning").Int }}

                {{ $awayTeam := $awayComp.Get "team" }}
                {{ $awayId := ($awayTeam.Get "id").Int }}
                {{ $awayRec := $awayComp.Get "leagueRecord" }}
                {{ $awayCode := ($awayTeam.Get "fileCode").String }}
                {{ $awayAbr := ($awayTeam.Get "abbreviation").String }}
                {{ $awaySc := ($awayComp.Get "score").Int }}

                {{ $homeTeam := $homeComp.Get "team" }}
                {{ $homeId := ($homeTeam.Get "id").Int }}
                {{ $homeRec := $homeComp.Get "leagueRecord" }}
                {{ $homeCode := ($homeTeam.Get "fileCode").String }}
                {{ $homeAbr := ($homeTeam.Get "abbreviation").String }}
                {{ $homeSc := ($homeComp.Get "score").Int }}

                                <li style="display:flex;align-items:center;gap:12px;padding:6px 0;border-bottom:1px solid var(--glance-divider);white-space:nowrap;"
                    {{ if ne $abs "Preview" }}title="{{ $awayAbr }} Box:{{ range $i, $inning := ($lines.Get "innings").Array }}{{ if lt $i 9 }}  Inning {{ $i | add 1 }}: {{ $inning.Get "away.runs" }}{{ else }}  Inning {{ $i | add 1 }}: {{ $inning.Get "away.runs" }}{{ end }}{{ end }}&#10;{{ $homeAbr }} Box:{{ range $i, $inning := ($lines.Get "innings").Array }}{{ if lt $i 9 }}  Inning {{ $i | add 1 }}: {{ $inning.Get "home.runs" }}{{ else }}  Inning {{ $i | add 1 }}: {{ $inning.Get "home.runs" }}{{ end }}{{ end }}"{{ end }}>
                  <span style="display:flex;align-items:center;width:80px;">
                    <img src='{{ if eq $awayCode "ana" }}https://a.espncdn.com/combiner/i?img=/i/teamlogos/mlb/500/laa.png&h=24&w=24{{ else }}https://a.espncdn.com/combiner/i?img=/i/teamlogos/mlb/500/{{ $awayCode }}.png&h=24&w=24{{ end }}' alt='{{ $awayAbr }}' style='width:24px;height:24px;margin-right:4px;'/>
                    <span style="display:flex;flex-direction:column;margin-right:4px;">
                      <span style="font-weight:bold;">{{ $awayAbr }}</span>
                      <span style="font-size:0.7em;color:var(--glance-muted-text)">({{ printf "%d-%d" ($awayRec.Get "wins").Int ($awayRec.Get "losses").Int }})</span>
                    </span>
                    {{ if ne $abs "Preview" }}<span style="margin-left:auto;font-size:0.95em;font-weight:700;">{{ $awaySc }}</span>{{ end }}
                                        {{ if and $live (eq $offTeamId $awayId) }}
                      <svg width="20" height="20" viewBox="0 0 24 24" style="margin-left:4px;flex-shrink:0;">
                        <rect x="6" y="6" width="12" height="12" fill="#808080" transform="rotate(45 12 12)"/>
                        {{ if $r1 }}<circle cx="17" cy="12" r="2.5" fill="white"/>{{ end }}
                        {{ if $r2 }}<circle cx="12" cy="7" r="2.5" fill="white"/>{{ end }}
                        {{ if $r3 }}<circle cx="7" cy="12" r="2.5" fill="white"/>{{ end }}
                      </svg>
                    {{ end }}
                  </span>

                  <span style="text-align:center;width:90px;display:inline-block;font-size:0.9em;">
                    {{ if eq $abs "Preview" }}
                      {{ $gameTime := (.String "gameDate" | parseTime "2006-01-02T15:04:05Z").In now.Location }}
                      <span style="color:var(--glance-muted-text)">{{ printf "%02d:%02d" $gameTime.Hour $gameTime.Minute }}</span>
                    {{ else if $live }}
                      {{ printf "%s %d" $inningSt $inningN }}
                    {{ else }}
                      {{ if or (eq ($status.Get "detailedState").String "Game Over") (eq ($status.Get "detailedState").String "Final") }}
                        <span style="font-size:1.1em;font-weight:500;">Final</span>
                        {{ if and (.Exists "seriesGameNumber") (.Exists "seriesDescription") }}
                          {{ $seriesDesc := .String "seriesDescription" }}
                          {{ if not (or (eq $seriesDesc "Regular Season") (eq $seriesDesc "Regular Season 1") (eq $seriesDesc "Regular Season 2") (eq $seriesDesc "Regular Season 3") (eq $seriesDesc "Regular Season 4")) }}
                            <div style="font-size:0.7em;color:var(--glance-accent-color);margin-top:2px;">
                              {{ $seriesDesc }} {{ .String "seriesGameNumber" }}
                            </div>
                          {{ end }}
                        {{ end }}
                      {{ else }}
                        {{ ($status.Get "detailedState").String }}
                      {{ end }}
                    {{ end }}
                  </span>

                  <span style="display:flex;align-items:center;width:80px;position:relative;">
                    <!-- Diamond positioned absolute to not affect flow -->
                    {{ if and $live (eq $offTeamId $homeId) }}
                      <svg width="20" height="20" viewBox="0 0 24 24" style="position:absolute;left:-24px;margin-right:0;flex-shrink:0;">
                        <rect x="6" y="6" width="12" height="12" fill="#808080" transform="rotate(45 12 12)"/>
                        {{ if $r1 }}<circle cx="17" cy="12" r="2.5" fill="white"/>{{ end }}
                        {{ if $r2 }}<circle cx="12" cy="7" r="2.5" fill="white"/>{{ end }}
                        {{ if $r3 }}<circle cx="7" cy="12" r="2.5" fill="white"/>{{ end }}
                      </svg>
                    {{ end }}
                    
                    <img src='{{ if eq $homeCode "ana" }}https://a.espncdn.com/combiner/i?img=/i/teamlogos/mlb/500/laa.png&h=24&w=24{{ else }}https://a.espncdn.com/combiner/i?img=/i/teamlogos/mlb/500/{{ $homeCode }}.png&h=24&w=24{{ end }}' alt='{{ $homeAbr }}' style='width:24px;height:24px;margin-right:4px;'/>
                    <span style="display:flex;flex-direction:column;margin-right:4px;">
                      <span style="font-weight:bold;">{{ $homeAbr }}</span>
                      <span style="font-size:0.7em;color:var(--glance-muted-text)">({{ printf "%d-%d" ($homeRec.Get "wins").Int ($homeRec.Get "losses").Int }})</span>
                    </span>
                    {{ if ne $abs "Preview" }}<span style="margin-left:auto;font-size:0.95em;font-weight:700;">{{ $homeSc }}</span>{{ end }}
                  </span>
                </li>
              {{ end }}
     ```