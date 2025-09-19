<img src="preview.png" height="1000px" />

```yaml
  - type: custom-api
    title: London Transport Status
    cache: 15m
    options:
      lines: [] # optional, e.g. ['windrush', 'jubilee', 'elizabeth']. default is all lines.
    template: |
      <div class="flex flex-column gap-5">
        {{ with .Options.lines }}
          {{ range . }}
            {{ with newRequest (printf "https://api.tfl.gov.uk/line/%s/status" .) | withHeader "User-Agent" "Glance-Dashboard/1.0" | getResponse }}
              {{ if eq .Response.StatusCode 200 }}
                {{ range .JSON.Array "" }}{{ template "line" . }}{{ end }}
              {{ end }}
            {{ end }}
          {{ end }}
        {{ else }}
          {{ with newRequest "https://api.tfl.gov.uk/line/mode/tube/status" | withHeader "User-Agent" "Glance-Dashboard/1.0" | getResponse }}
            {{ if eq .Response.StatusCode 200 }}{{ range .JSON.Array "" }}{{ template "line" . }}{{ end }}{{ end }}
          {{ end }}
          {{ with newRequest "https://api.tfl.gov.uk/line/mode/elizabeth-line/status" | withHeader "User-Agent" "Glance-Dashboard/1.0" | getResponse }}
            {{ if eq .Response.StatusCode 200 }}{{ range .JSON.Array "" }}{{ template "line" . }}{{ end }}{{ end }}
          {{ end }}
          {{ with newRequest "https://api.tfl.gov.uk/line/mode/overground/status" | withHeader "User-Agent" "Glance-Dashboard/1.0" | getResponse }}
            {{ if eq .Response.StatusCode 200 }}{{ range .JSON.Array "" }}{{ template "line" . }}{{ end }}{{ end }}
          {{ end }}
        {{ end }}
      </div>

      {{ define "line" }}
      {{ $id := .String "id" }}{{ $status := (index (.Array "lineStatuses") 0).String "statusSeverityDescription" }}
      <div class="flex justify-between items-center" style="padding: 8px 0; border-bottom: 1px solid rgba(255,255,255,0.1);">
        <div class="flex items-center gap-10">
          <div style="width: 20px; height: 4px; border-radius: 2px; background-color: {{ if eq $id "bakerloo" }}#B36305{{ else if eq $id "central" }}#E32017{{ else if eq $id "circle" }}#FFD300{{ else if eq $id "district" }}#00782A{{ else if eq $id "hammersmith-city" }}#F3A9BB{{ else if eq $id "jubilee" }}#A0A5A9{{ else if eq $id "metropolitan" }}#9B0056{{ else if eq $id "northern" }}#000000{{ else if eq $id "piccadilly" }}#003688{{ else if eq $id "victoria" }}#0098D4{{ else if eq $id "waterloo-city" }}#95CDBA{{ else if eq $id "elizabeth" }}#7156A5{{ else if eq $id "lioness" }}#E9B908{{ else if eq $id "mildmay" }}#0019A8{{ else if eq $id "windrush" }}#D05F0E{{ else if eq $id "weaver" }}#8F9499{{ else if eq $id "suffragette" }}#64A334{{ else if eq $id "liberty" }}#B05F8C{{ else }}#666{{ end }};"></div>
          <span class="size-h4">{{ .String "name" }}</span>
        </div>
        <div class="text-right">
          <span class="size-h5 {{ if eq $status "Good Service" }}color-positive{{ else if eq $status "Minor Delays" }}color-primary{{ else }}color-negative{{ end }}">{{ $status }}</span>
        </div>
      </div>
      {{ end }}
```

## Options

- `lines`: Optionally specify which lines to display. By default, all lines are shown.
  - Valid values: `["bakerloo", "central", "circle", "district", "hammersmith-city", "jubilee", "metropolitan", "northern", "piccadilly", "victoria", "waterloo-city", "elizabeth", "liberty", "lioness", "mildmay", "suffragette", "weaver", "windrush" ]`
