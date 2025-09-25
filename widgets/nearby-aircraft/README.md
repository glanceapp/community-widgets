![](preview.png)

```yaml
- type: custom-api
  title: Nearby Aircraft
  cache: 1m
  url: https://opensky-network.org/api/states/all?lamin=${MY_LAT_MIN}&lomin=${MY_LON_MIN}&lamax=${MY_LAT_MAX}&lomax=${MY_LON_MAX}
  template: |
    <div class="opensky-aircraft-stats">
      {{ $states := .JSON.Array "states" }}
      {{ $statesCount := len $states }}

      <div class="margin-block-2">
        <div class="size-h6">AIRCRAFT NEARBY: {{ $statesCount }} </div>
      </div>

      {{ if gt $statesCount 0 }}
        {{ $plane := index $states 0 }}
        {{ $callsign := $plane.String "1" }}
        {{ $icao24 := $plane.String "0" }}
        {{ $origin := $plane.String "2" }}
        {{ $onGround := $plane.Bool "8" }}

        <div class="margin-block-2">
          <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 10px;">
            <div>
              <div class="color-highlight size-h3">{{ if ne $callsign "" }}{{ $callsign }}{{ else }}{{ $icao24 }}{{ end }}</div>
              <div class="size-h6">CALLSIGN</div>
            </div>
            <div>
              <div class="color-highlight size-h3 text-truncate" style="max-width: 100%;">{{ if ne $origin "" }}{{ $origin }}{{ else }}-{{ end }}</div>
              <div class="size-h6">ORIGIN</div>
            </div>
          </div>
        </div>

        <div class="margin-block-2">
          <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 10px;">
            <div>
              {{ $altitude := $plane.Float "7" }}
              <div class="color-highlight size-h3">{{ if and (not $onGround) ($altitude | ne 0.0) }}{{ $altitude | printf "%.0f" }}m{{ else }}-{{ end }}</div>
              <div class="size-h6">ALTITUDE</div>
            </div>
            <div>
              {{ $speed := $plane.Float "9" }}
              <div class="color-highlight size-h3">{{ if and (not $onGround) ($speed | ne 0.0) }}{{ $speed | printf "%.0f" }}m/s{{ else }}-{{ end }}</div>
              <div class="size-h6">SPEED</div>
            </div>
          </div>
        </div>
      {{ end }}
    </div>
```

## Environment variables

- `MY_LAT_MIN` - Minimum latitude (southern boundary)
- `MY_LAT_MAX` - Maximum latitude (northern boundary)
- `MY_LON_MIN` - Minimum longitude (western boundary)
- `MY_LON_MAX` - Maximum longitude (eastern boundary)
