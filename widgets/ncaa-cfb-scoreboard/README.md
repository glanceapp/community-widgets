![](preview.png)

```yaml
- type: custom-api
  title: NCAA CFB Scoreboard
  url: "https://site.api.espn.com/apis/site/v2/sports/football/college-football/scoreboard"
  cache: 1m
  template: |
    {{/* NCAA CFB Scoreboard */}}
    {{ $events := .JSON.Array "events" }}

    {{ if eq (len $events) 0 }}
    <div>No games scheduled today.</div>
    {{ else }}
    {{ $wid := printf "cfb-%d" now.UnixNano }}
    <div id="{{ $wid }}" style="position:relative;padding-bottom:{{ if gt (len $events) 6 }}25px{{ else }}0px{{ end }};">
        <input type="checkbox" id="{{ $wid }}-cb" style="display:none"/>
        <ul class="visible" style="list-style:none;padding:0;margin:0;">
        {{ $count := 0 }}

        {{ range $i, $ev := $events }}
            {{ if and (lt $count 6) (gt (len ($ev.Array "competitions.0.competitors")) 1) (eq ($ev.String "competitions.0.status.type.name") "STATUS_IN_PROGRESS") }}
            {{ template "cfb-game" $ev }}
            {{ $count = add $count 1 }}
            {{ end }}
        {{ end }}

        {{ range $i, $ev := $events }}
            {{ if and (lt $count 6) (gt (len ($ev.Array "competitions.0.competitors")) 1) (eq ($ev.String "competitions.0.status.type.name") "STATUS_SCHEDULED") }}
            {{ template "cfb-game" $ev }}
            {{ $count = add $count 1 }}
            {{ end }}
        {{ end }}

        {{ range $i, $ev := $events }}
            {{ if and (lt $count 6) (gt (len ($ev.Array "competitions.0.competitors")) 1) (and (ne ($ev.String "competitions.0.status.type.name") "STATUS_IN_PROGRESS") (ne ($ev.String "competitions.0.status.type.name") "STATUS_SCHEDULED")) }}
            {{ template "cfb-game" $ev }}
            {{ $count = add $count 1 }}
            {{ end }}
        {{ end }}
        </ul>

        <ul class="all" style="list-style:none;padding:0;margin:0;display:none">
        {{ range $i, $ev := $events }}
            {{ if and (gt (len ($ev.Array "competitions.0.competitors")) 1) (eq ($ev.String "competitions.0.status.type.name") "STATUS_IN_PROGRESS") }}
            {{ template "cfb-game" $ev }}
            {{ end }}
        {{ end }}

        {{ range $i, $ev := $events }}
            {{ if and (gt (len ($ev.Array "competitions.0.competitors")) 1) (eq ($ev.String "competitions.0.status.type.name") "STATUS_SCHEDULED") }}
            {{ template "cfb-game" $ev }}
            {{ end }}
        {{ end }}

        {{ range $i, $ev := $events }}
            {{ if and (gt (len ($ev.Array "competitions.0.competitors")) 1) (and (ne ($ev.String "competitions.0.status.type.name") "STATUS_IN_PROGRESS") (ne ($ev.String "competitions.0.status.type.name") "STATUS_SCHEDULED")) }}
            {{ template "cfb-game" $ev }}
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

    {{ define "cfb-game" }}
    {{ $comp := index (.Array "competitions") 0 }}
    {{ $away := index ($comp.Array "competitors") 0 }}
    {{ $home := index ($comp.Array "competitors") 1 }}

    {{ $awayRec := "" }}
    {{ if gt (len ($away.Array "records")) 0 }}{{ $awayRec = (index ($away.Array "records") 0).String "summary" }}{{ end }}
    {{ $homeRec := "" }}
    {{ if gt (len ($home.Array "records")) 0 }}{{ $homeRec = (index ($home.Array "records") 0).String "summary" }}{{ end }}

    {{ $awaySchedule := "" }}
    {{ range $away.Array "team.links" }}
        {{ if eq (.String "text") "Schedule" }}
        {{ $awaySchedule = .String "href" }}
        {{ end }}
    {{ end }}

    {{ $homeSchedule := "" }}
    {{ range $home.Array "team.links" }}
        {{ if eq (.String "text") "Schedule" }}
        {{ $homeSchedule = .String "href" }}
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

    {{ $leagueLogo := "https://a.espncdn.com/redesign/assets/img/icons/ESPN-icon-football-college.png" }}
    {{ $awayLogo := $away.String "team.logo" }}
    {{ if eq $awayLogo "" }}{{ $awayLogo = $leagueLogo }}{{ end }}
    {{ $homeLogo := $home.String "team.logo" }}
    {{ if eq $homeLogo "" }}{{ $homeLogo = $leagueLogo }}{{ end }}

    {{ $awayBox := "" }}
    {{ $homeBox := "" }}
    {{ if gt (len ($away.Array "linescores")) 0 }}
        {{ range $i2, $ls := $away.Array "linescores" }}
        {{ $qn := add $i2 1 }}
        {{ $val := $ls.String "displayValue" }}
        {{ $awayBox = printf "%sQ%d:%s " $awayBox $qn $val }}
        {{ end }}
    {{ end }}
    {{ if gt (len ($home.Array "linescores")) 0 }}
        {{ range $i2, $ls := $home.Array "linescores" }}
        {{ $qn := add $i2 1 }}
        {{ $val := $ls.String "displayValue" }}
        {{ $homeBox = printf "%sQ%d:%s " $homeBox $qn $val }}
        {{ end }}
    {{ end }}

    {{ $leadersText := "" }}
    {{ range $comp.Array "leaders" }}
        {{ $name := .String "name" }}
        {{ if or (eq $name "passing") (eq $name "rushing") (eq $name "receiving") }}
        {{ range .Array "leaders" }}
            {{ $ath := .String "athlete.shortName" }}
            {{ $disp := .String "displayValue" }}
            {{ if ne $ath "" }}{{ $leadersText = printf "%s%s: %s\n" $leadersText $ath $disp }}{{ end }}
        {{ end }}
        {{ end }}
    {{ end }}

    {{ $state := $comp.String "status.type.name" }}
    {{/* alias schedule link variable names to match nfl template naming */}}
    {{ $awayScheduleLink := $awaySchedule }}
    {{ $homeScheduleLink := $homeSchedule }}

    {{ $tooltip := printf "%s Box: %s\n%s Box: %s\n\nGame Leaders:\n%s" ($away.String "team.abbreviation") $awayBox ($home.String "team.abbreviation") $homeBox $leadersText }}

    <li style="display:flex;align-items:center;white-space:nowrap;gap:12px;padding:6px 0;border-bottom:1px solid var(--glance-divider);cursor:default" {{ if ne $state "STATUS_SCHEDULED" }}title="{{ $tooltip }}"{{ end }}>
        <span style="display:flex;align-items:center;flex:1;justify-content:space-between;">
        {{ if $awayScheduleLink }}
            <a href="{{ $awayScheduleLink }}" target="_blank" style="text-decoration:none;color:inherit;display:flex;align-items:center;">
            <img src="{{ $awayLogo }}" alt="{{ $away.String "team.abbreviation" }}" style="width:24px;height:24px;margin-right:6px;cursor:pointer;"/>
            </a>
        {{ else }}
            <img src="{{ $awayLogo }}" alt="{{ $away.String "team.abbreviation" }}" style="width:24px;height:24px;margin-right:6px;"/>
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
            <img src="{{ $homeLogo }}" alt="{{ $home.String "team.abbreviation" }}" style="width:24px;height:24px;margin-right:6px;cursor:pointer;"/>
            </a>
        {{ else }}
            <img src="{{ $homeLogo }}" alt="{{ $home.String "team.abbreviation" }}" style="width:24px;height:24px;margin-right:6px;"/>
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
