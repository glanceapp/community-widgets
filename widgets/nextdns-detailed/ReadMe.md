![](preview.png)


```yaml
- type: custom-api
  title: NextDNS
  title-url: https://my.nextdns.io/${NEXTDNS_PROFILE_ID}/analytics
  cache: 1m
  options:
    profile: ${NEXTDNS_PROFILE_ID}
    api-key: ${NEXTDNS_API_KEY}
  template: |
    <style>
      .ndns-totals { display: flex; justify-content: space-between; text-align: center; margin-bottom: 1rem; margin-top: 0.5rem; }
      .ndns-totals-value { font-size: 2rem; font-weight: 700; color: var(--color-text-highlight); line-height: 1.2; }
      .ndns-totals-label { font-size: var(--font-size-h6); opacity: 0.6; margin-top: 0.25rem; text-transform: uppercase; }
      .text-red { color: var(--color-negative); }
      .text-primary { color: var(--color-primary); }

      .ndns-list { list-style: none; margin: 0; padding: 0; display: flex; flex-direction: column; gap: 0.25rem; font-size: var(--font-size-base); }
      .ndns-item { display: flex; justify-content: space-between; align-items: center; padding: 0.2rem 0; }

      .ndns-details { margin-top: 1rem; }
      .ndns-details summary { cursor: pointer; font-size: var(--font-size-base); opacity: 0.8; user-select: none; }

      .ndns-name { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; flex: 1; margin-right: 1rem; display: flex; align-items: center; gap: 8px; color: var(--color-text-paragraph); }
      .ndns-badge { text-align: right; font-family: monospace; opacity: 0.7; font-size: var(--font-size-base); }
      .ndns-device-txt { font-size: var(--font-size-h6); opacity: 0.5; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; max-width: 80px; text-align: right; margin-left: 0.25rem; }

      .ndns-dot { width: 8px; height: 8px; border-radius: 50%; display: inline-block; flex-shrink: 0; }
      .dot-green { background: var(--color-positive); box-shadow: 0 0 5px var(--color-positive); }
      .dot-red { background: var(--color-negative); box-shadow: 0 0 5px var(--color-negative); }
      .dot-blue { background: var(--color-primary); }
    </style>

    {{ $profile := .Options.StringOr "profile" "" }}
    {{ $apiKey  := .Options.StringOr "api-key" "" }}
    {{ $baseURL := print "https://api.nextdns.io/profiles/" $profile }}

    {{ $statusResp  := newRequest (print $baseURL "/analytics/status?from=-30d") | withHeader "X-Api-Key" $apiKey | getResponse }}
    {{ $reasonsResp := newRequest (print $baseURL "/analytics/reasons?from=-30d&limit=5&lang=en") | withHeader "X-Api-Key" $apiKey | getResponse }}
    {{ $devicesResp := newRequest (print $baseURL "/analytics/devices?from=-30d&limit=5") | withHeader "X-Api-Key" $apiKey | getResponse }}
    {{ $destResp    := newRequest (print $baseURL "/analytics/destinations?type=countries&from=-30d&limit=5&lang=en") | withHeader "X-Api-Key" $apiKey | getResponse }}
    {{ $logsResp    := newRequest (print $baseURL "/logs?lang=en&limit=10") | withHeader "X-Api-Key" $apiKey | getResponse }}

    <div>
      {{ if $statusResp.JSON }}
        {{ $total := 0 }}{{ $blocked := 0 }}
        {{ range $statusResp.JSON.Array "data" }}
          {{ $total = add $total (.Int "queries") }}
          {{ if eq (.String "status") "blocked" }}
            {{ $blocked = add $blocked (.Int "queries") }}
          {{ end }}
        {{ end }}
        <div class="ndns-totals">
          <div>
            <div class="ndns-totals-value">{{ $total | formatNumber }}</div>
            <div class="ndns-totals-label">QUERIES</div>
          </div>
          <div>
            <div class="ndns-totals-value text-red">{{ $blocked | formatNumber }}</div>
            <div class="ndns-totals-label">BLOCKED</div>
          </div>
          <div>
            <div class="ndns-totals-value text-primary">
              {{ if gt $total 0 }}
                {{ div (mul $blocked 100.0) $total | printf "%.1f" }}%
              {{ else }}
                0.0%
              {{ end }}
            </div>
            <div class="ndns-totals-label">RATE</div>
          </div>
        </div>
      {{ end }}

      {{ if $reasonsResp.JSON }}
        <details class="ndns-details">
          <summary>Block Reasons</summary>
          <ul class="ndns-list" style="margin-top: 0.625rem;">
            {{ range $reasonsResp.JSON.Array "data" }}
              <li class="ndns-item">
                <div class="ndns-name"><span class="ndns-dot dot-red"></span>{{ .String "name" }}</div>
                <div class="ndns-badge">{{ .Int "queries" }}</div>
              </li>
            {{ end }}
          </ul>
        </details>
      {{ end }}

      {{ if $devicesResp.JSON }}
        <details class="ndns-details">
          <summary>Active Devices</summary>
          <ul class="ndns-list" style="margin-top: 0.625rem;">
            {{ range $devicesResp.JSON.Array "data" }}
              <li class="ndns-item">
                <div class="ndns-name"><span class="ndns-dot dot-blue"></span>{{ .String "name" }}</div>
                <div class="ndns-badge">{{ .Int "queries" }}</div>
              </li>
            {{ end }}
          </ul>
        </details>
      {{ end }}

      {{ if $destResp.JSON }}
        <details class="ndns-details">
          <summary>Destinations</summary>
          <ul class="ndns-list" style="margin-top: 0.625rem;">
            {{ range $destResp.JSON.Array "data" }}
              <li class="ndns-item">
                <div class="ndns-name"><span class="ndns-dot dot-blue"></span>{{ .String "name" }}</div>
                <div class="ndns-badge">{{ .Int "queries" }}</div>
              </li>
            {{ end }}
          </ul>
        </details>
      {{ end }}

      {{ if $logsResp.JSON }}
        <details class="ndns-details" open>
          <summary>Live Traffic</summary>
          <ul class="ndns-list" style="margin-top: 0.625rem;">
            {{ range $logsResp.JSON.Array "data" }}
              <li class="ndns-item">
                <div class="ndns-name" title="{{ .String "domain" }}">
                  <span class="ndns-dot {{ if eq (.String "status") "blocked" }}dot-red{{ else }}dot-green{{ end }}"></span>
                  {{ .String "domain" }}
                </div>
                {{ $devName := .String "deviceName" }}
                {{ if ne $devName "" }}
                  <div class="ndns-device-txt">{{ $devName }}</div>
                {{ end }}
              </li>
            {{ end }}
          </ul>
        </details>
      {{ end }}
    </div>
```

## Environment variables

- `NEXTDNS_PROFILE_ID` - Your unique profile ID found in NextDNS control panel
- `NEXTDNS_API_KEY` - Configure and unique api key in your user profile
