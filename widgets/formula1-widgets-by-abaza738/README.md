## Preview
![preview](preview.png)

## Config

Choose one or more of the following widgets.

```yaml
- type: custom-api
  title: Next Race
  cache: 2h
  url: https://f1api.dev/api/current/next
  template: |
    <div class="flex flex-column gap-10">
      {{ $session := index (.JSON.Array "race") 0 }}
      <p class="size-h5">
        Round {{ .JSON.String "round" }}
      </p>

      <div class="margin-block-4">
        <p class="color-highlight">{{ $session.String "raceName" }}</p>
        <p class="color-primary">
          <span>Race</span>
          {{ $datetime := concat ($session.String "schedule.race.date") "T" ($session.String "schedule.race.time") }}
          <span
            class="color-highlight"
            title="{{ $session.String "schedule.race.date" }}"
            {{ parseRelativeTime "rfc3339" $datetime }}
          ></span>
        </p>
        <p class="size-h5">{{ $session.String "schedule.race.date" }} at {{ $session.String "schedule.race.time" }}</p>
      </div>

      <ul class="size-h5 attachments">
        <li>{{ $session.String "circuit.country" }}</li>
        <li>{{ $session.String "circuit.city" }}</li>
        <li>{{ $session.String "circuit.circuitName" }}</li>
      </ul>
    </div>

- type: custom-api
  title: Last Race Results
  cache: 1d
  url: https://f1api.dev/api/current/last/race
  template: |
    <div class="flex flex-column gap-10">
      <p class="size-h5">{{ .JSON.String "races.raceName" }}</p>
      <ul class="list collapsible-container" data-collapse-after="5">
        {{ range $i, $v := .JSON.Array "races.results" }}
        <li class="flex items-center {{ if eq $i 0 }}color-primary{{ else if eq $i 1 }}color-highlight{{ end }}">
          <span class="grow min-width-0">{{ .String "position" }}. {{ .String "driver.shortName" }}</span>
          <span class="shrink-0 text-right">{{ .String "time" }}</span>
        </li>
        {{ end }}
      </ul>
    </div>

- type: custom-api
  title: Drivers Standings
  cache: 1d
  url: https://f1api.dev/api/current/drivers-championship
  template: |
    <ul class="list collapsible-container" data-collapse-after="5">
      {{ range $i, $v := .JSON.Array "drivers_championship" }}
      <li class="flex items-center {{ if eq $i 0 }}color-primary{{ else if eq $i 1 }}color-highlight{{ end }}">
        <div class="grow min-width-0">
          <span>{{ .String "position" }}. {{ .String "driver.shortName" }}</span>
          <span class="size-h6">({{ .String "team.teamName" }})</span>
        </div>
        <span class="shrink-0 text-right">{{ .String "points" }}</span>
      </li>
      {{ end }}
    </ul>

- type: custom-api
  title: Constructors Standings
  cache: 1d
  url: https://f1api.dev/api/current/constructors-championship
  template: |
    <ul class="list collapsible-container" data-collapse-after="5">
      {{ range $i, $v := .JSON.Array "constructors_championship" }}
      <li class="flex items-center {{ if eq $i 0 }}color-primary{{ else if eq $i 1 }}color-highlight{{ end }}">
        <span class="grow min-width-0">{{ .String "position" }}. {{ .String "team.teamName" }}</span>
        <span class="shrink-0 text-right">{{ .String "points" }}</span>
      </li>
      {{ end }}
    </ul>
```
