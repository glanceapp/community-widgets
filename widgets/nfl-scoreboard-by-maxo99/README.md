![](preview.png)

```yaml
- type: custom-api
  title: NFL Scoreboard
  url: "https://site.api.espn.com/apis/site/v2/sports/football/nfl/scoreboard"
  cache: 15s
  template: |
    {{ $events := .JSON.Array "events" }}
      {{ if eq (len $events) 0 }}
        <div>No games scheduled today.</div>
      {{ else }}
        {{ $wid := printf "nfl-%d" now.UnixNano }}
        <div id="{{ $wid }}" style="position:relative;padding-bottom:{{ if gt (len $events) 6 }}25px{{ else }}0px{{ end }};">
          <input type="checkbox" id="{{ $wid }}-cb" style="display:none"/>
          <ul class="visible" style="list-style:none;padding:0;margin:0;">
            {{ $count := 0 }}

            {{ range $i, $g := $events }}
              {{ if and (lt $count 6) (eq ($g.String "competitions.0.status.type.name") "STATUS_IN_PROGRESS") }}
                {{ template "nfl-game" $g }}
                {{ $count = add $count 1 }}
              {{ end }}
            {{ end }}

            {{ range $i, $g := $events }}
              {{ if and (lt $count 6) (eq ($g.String "competitions.0.status.type.name") "STATUS_SCHEDULED") }}
                {{ template "nfl-game" $g }}
                {{ $count = add $count 1 }}
              {{ end }}
            {{ end }}

            {{ range $i, $g := $events }}
              {{ if and (lt $count 6) (not (or (eq ($g.String "competitions.0.status.type.name") "STATUS_IN_PROGRESS") (eq ($g.String "competitions.0.status.type.name") "STATUS_SCHEDULED"))) }}
                {{ template "nfl-game" $g }}
                {{ $count = add $count 1 }}
              {{ end }}
            {{ end }}
          </ul>

          <ul class="all" style="list-style:none;padding:0;margin:0;display:none">
            {{ range $i, $g := $events }}
              {{ if eq ($g.String "competitions.0.status.type.name") "STATUS_IN_PROGRESS" }}
                {{ template "nfl-game" $g }}
              {{ end }}
            {{ end }}

            {{ range $i, $g := $events }}
              {{ if eq ($g.String "competitions.0.status.type.name") "STATUS_SCHEDULED" }}
                {{ template "nfl-game" $g }}
              {{ end }}
            {{ end }}

            {{ range $i, $g := $events }}
              {{ if not (or (eq ($g.String "competitions.0.status.type.name") "STATUS_IN_PROGRESS") (eq ($g.String "competitions.0.status.type.name") "STATUS_SCHEDULED")) }}
                {{ template "nfl-game" $g }}
              {{ end }}
            {{ end }}
          </ul>

          {{ if gt (len $events) 6 }}
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

      {{ define "nfl-game" }}
        {{ $state := .String "competitions.0.status.type.name" }}
        {{ $away := index (.Array "competitions.0.competitors") 0 }}
        {{ $home := index (.Array "competitions.0.competitors") 1 }}
        {{ $awayRec := "" }}
        {{ if gt (len ($away.Array "records")) 0 }}{{ $awayRec = (index ($away.Array "records") 0).String "summary" }}{{ end }}
        {{ $homeRec := "" }}
        {{ if gt (len ($home.Array "records")) 0 }}{{ $homeRec = (index ($home.Array "records") 0).String "summary" }}{{ end }}

        {{ $awayScheduleLink := "" }}
        {{ range $away.Array "team.links" }}
          {{ if eq (.String "text") "Schedule" }}
            {{ $awayScheduleLink = .String "href" }}
          {{ end }}
        {{ end }}

        {{ $homeScheduleLink := "" }}
        {{ range $home.Array "team.links" }}
          {{ if eq (.String "text") "Schedule" }}
            {{ $homeScheduleLink = .String "href" }}
          {{ end }}
        {{ end }}

        {{ $gameLink := "" }}
        {{ range .Array "links" }}
          {{ $isGameLink := false }}
          {{ range .Array "rel" }}
            {{ if eq (printf "%v" .) "summary" }}
              {{ $isGameLink = true }}
            {{ end }}
          {{ end }}
          {{ if $isGameLink }}
            {{ $gameLink = .String "href" }}
          {{ end }}
        {{ end }}

        {{ $tooltip := "" }}
        {{ if ne $state "STATUS_SCHEDULED" }}
          {{ $awayBoxScore := "" }}
          {{ range $j,$ls := $away.Array "linescores" }}
            {{ if eq $j 0 }}{{ $awayBoxScore = printf "%s Q1: %s" $awayBoxScore ($ls.String "displayValue") }}
            {{ else if eq $j 1 }}{{ $awayBoxScore = printf "%s Q2: %s" $awayBoxScore ($ls.String "displayValue") }}
            {{ else if eq $j 2 }}{{ $awayBoxScore = printf "%s Q3: %s" $awayBoxScore ($ls.String "displayValue") }}
            {{ else if eq $j 3 }}{{ $awayBoxScore = printf "%s Q4: %s" $awayBoxScore ($ls.String "displayValue") }}
            {{ else }}{{ $awayBoxScore = printf "%s OT: %s" $awayBoxScore ($ls.String "displayValue") }}{{ end }}
          {{ end }}

          {{ $homeBoxScore := "" }}
          {{ range $j,$ls := $home.Array "linescores" }}
            {{ if eq $j 0 }}{{ $homeBoxScore = printf "%s Q1: %s" $homeBoxScore ($ls.String "displayValue") }}
            {{ else if eq $j 1 }}{{ $homeBoxScore = printf "%s Q2: %s" $homeBoxScore ($ls.String "displayValue") }}
            {{ else if eq $j 2 }}{{ $homeBoxScore = printf "%s Q3: %s" $homeBoxScore ($ls.String "displayValue") }}
            {{ else if eq $j 3 }}{{ $homeBoxScore = printf "%s Q4: %s" $homeBoxScore ($ls.String "displayValue") }}
            {{ else }}{{ $homeBoxScore = printf "%s OT: %s" $homeBoxScore ($ls.String "displayValue") }}{{ end }}
          {{ end }}

          {{ $gameLeaders := "" }}
          {{ range .Array "competitions.0.leaders" }}
            {{ if eq (.String "name") "passingYards" }}
              {{ range .Array "leaders" }}
                {{ $gameLeaders = printf "%sPassing: %s (%s)\n" $gameLeaders (.String "athlete.shortName") (.String "displayValue") }}
              {{ end }}
            {{ else if eq (.String "name") "rushingYards" }}
              {{ range .Array "leaders" }}
                {{ $gameLeaders = printf "%sRushing: %s (%s)\n" $gameLeaders (.String "athlete.shortName") (.String "displayValue") }}
              {{ end }}
            {{ else if eq (.String "name") "receivingYards" }}
              {{ range .Array "leaders" }}
                {{ $gameLeaders = printf "%sReceiving: %s (%s)\n" $gameLeaders (.String "athlete.shortName") (.String "displayValue") }}
              {{ end }}
            {{ end }}
          {{ end }}

          {{ $tooltip = printf "%s Box:%s\n%s Box:%s\n\nGame Leaders:\n%s" ($away.String "team.abbreviation") $awayBoxScore ($home.String "team.abbreviation") $homeBoxScore $gameLeaders }}
        {{ end }}

        <li style="display:flex;align-items:center;white-space:nowrap;gap:12px;padding:6px 0;border-bottom:1px solid var(--glance-divider);cursor:default" {{ if ne $state "STATUS_SCHEDULED" }}title="{{ $tooltip }}"{{ end }}>
          <span style="display:flex;align-items:center;flex:1;justify-content:space-between;">
            {{ if $awayScheduleLink }}
              <a href="{{ $awayScheduleLink }}" target="_blank" style="text-decoration:none;color:inherit;display:flex;align-items:center;">
                <img src="{{ if $away.String "team.logo" }}{{ $away.String "team.logo" }}{{ else }}https://a.espncdn.com/i/teamlogos/leagues/500/nfl.png{{ end }}" alt="{{ $away.String "team.abbreviation" }}" style="width:24px;height:24px;margin-right:6px;cursor:pointer;"/>
              </a>
            {{ else }}
              <img src="{{ if $away.String "team.logo" }}{{ $away.String "team.logo" }}{{ else }}https://a.espncdn.com/i/teamlogos/leagues/500/nfl.png{{ end }}" alt="{{ $away.String "team.abbreviation" }}" style="width:24px;height:24px;margin-right:6px;"/>
            {{ end }}
            <span style="display:flex;flex-direction:column;flex:1;">
              <span style="font-weight:bold;">{{ $away.String "team.abbreviation" }}</span>
              {{ if $awayRec }}<span style="font-size:0.7em;color:var(--glance-muted-text)">({{ $awayRec }})</span>{{ end }}
            </span>
            {{ if ne $state "STATUS_SCHEDULED" }}<span style="font-size:1.1em;font-weight:500;min-width:20px;text-align:right;">{{ $away.String "score" }}</span>{{ end }}
          </span>
          {{ if $gameLink }}
            <a href="{{ $gameLink }}" target="_blank" style="text-decoration:none;color:inherit;">
              <span style="display:flex;flex-direction:column;align-items:center;flex-shrink:0;margin:0 8px;cursor:pointer;">
                <span>
                  {{ if eq $state "STATUS_IN_PROGRESS" }}
                    {{ $period := .String "competitions.0.status.period" }}
                    <span style="font-size:0.85em;font-weight:500;">{{ if eq $period "1" }}1st{{ else if eq $period "2" }}2nd{{ else if eq $period "3" }}3rd{{ else if eq $period "4" }}4th{{ else }}OT{{ end }} {{ .String "competitions.0.status.displayClock" }}</span>
                  {{ else if eq $state "STATUS_SCHEDULED" }}
                    <span style="font-size:0.85em;">{{ .String "competitions.0.status.type.shortDetail" }}</span>
                  {{ else }}
                    <span style="font-size:0.85em;font-weight:500;">{{ .String "competitions.0.status.type.shortDetail" }}</span>
                  {{ end }}
                </span>

                {{ if .Exists "competitions.0.series" }}
                  <span style="font-size:0.7em;color:var(--glance-accent-color);margin-top:5px;">{{ .String "competitions.0.series.summary" }}</span>
                {{ end }}
              </span>
            </a>
          {{ else }}
            <span style="display:flex;flex-direction:column;align-items:center;flex-shrink:0;margin:0 8px;">
              <span>
                {{ if eq $state "STATUS_IN_PROGRESS" }}
                  {{ $period := .String "competitions.0.status.period" }}
                  <span style="font-size:0.85em;font-weight:500;">{{ if eq $period "1" }}1st{{ else if eq $period "2" }}2nd{{ else if eq $period "3" }}3rd{{ else if eq $period "4" }}4th{{ else }}OT{{ end }} {{ .String "competitions.0.status.displayClock" }}</span>
                {{ else if eq $state "STATUS_SCHEDULED" }}
                  <span style="font-size:0.85em;">{{ .String "competitions.0.status.type.shortDetail" }}</span>
                {{ else }}
                  <span style="font-size:0.85em;font-weight:500;">{{ .String "competitions.0.status.type.shortDetail" }}</span>
                {{ end }}
              </span>

              {{ if .Exists "competitions.0.series" }}
                <span style="font-size:0.7em;color:var(--glance-accent-color);margin-top:5px;">{{ .String "competitions.0.series.summary" }}</span>
              {{ end }}
            </span>
          {{ end }}
          <span style="display:flex;align-items:center;flex:1;justify-content:space-between;">
            {{ if $homeScheduleLink }}
              <a href="{{ $homeScheduleLink }}" target="_blank" style="text-decoration:none;color:inherit;display:flex;align-items:center;">
                <img src="{{ if $home.String "team.logo" }}{{ $home.String "team.logo" }}{{ else }}https://a.espncdn.com/i/teamlogos/leagues/500/nfl.png{{ end }}" alt="{{ $home.String "team.abbreviation" }}" style="width:24px;height:24px;margin-right:6px;cursor:pointer;"/>
              </a>
            {{ else }}
              <img src="{{ if $home.String "team.logo" }}{{ $home.String "team.logo" }}{{ else }}https://a.espncdn.com/i/teamlogos/leagues/500/nfl.png{{ end }}" alt="{{ $home.String "team.abbreviation" }}" style="width:24px;height:24px;margin-right:6px;"/>
            {{ end }}
            <span style="display:flex;flex-direction:column;flex:1;">
              <span style="font-weight:bold;">{{ $home.String "team.abbreviation" }}</span>
              {{ if $homeRec }}<span style="font-size:0.7em;color:var(--glance-muted-text)">({{ $homeRec }})</span>{{ end }}
            </span>
            {{ if ne $state "STATUS_SCHEDULED" }}<span style="font-size:1.1em;font-weight:500;min-width:20px;text-align:right;">{{ $home.String "score" }}</span>{{ end }}
          </span>
        </li>
      {{ end }}
```
