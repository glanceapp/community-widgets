# Overview
For feature requests or issues, contact @SleebySky on the project's Discord server. Thanks to @Haondt for his assistance on this widget. 
# Widget

![](./drive_health.png)
```yml
- type: custom-api
  title: "Drive Health"
  url: http://${SCRUTINY_URL}/api/summary
  cache: 1h
  method: GET
  options:
    filter_archived: false 
  template: |
    {{- $filterArchived := .Options.filter_archived }}
    {{- $total := 0 }}
    {{- range $wwn, $rec := (.JSON.Get "data.summary").Map }}
      {{- $archived := $rec.Get "device.archived" }}
      {{- if or (not $filterArchived) (and $archived.Exists (not $archived.Bool)) }}
        {{- $total = add $total 1 }}
      {{- end }}
    {{- end }}

    {{- $count:= 0}}

    {{- range $wwn, $rec := (.JSON.Get "data.summary").Map }}
      {{- $archived := $rec.Get "device.archived" }}
      {{- if or (not $filterArchived) (and $archived.Exists (not $archived.Bool)) }}
        {{- $count = add $count 1 }}

        {{- $device := $rec.Get "device.device_name" }}
        {{- $model := $rec.Get "device.model_name" }}
        {{- $hoursFloat := ($rec.Get "smart.power_on_hours").Num }}
        {{- $daysFloat := div $hoursFloat 24 }}
        {{- $days := printf "%.0f" $daysFloat }}
        {{- $capacity := ($rec.Get "device.capacity").Num }}
        {{- $capacity_tb := div $capacity 1000000000000 }}
        {{- $capacity_print := printf "%.0f" $capacity_tb }}
        {{- $status := ($rec.Get "device.device_status").Num }}

        {{- $tempHistory := $rec.Get "temp_history" }}
        {{- $lastIndex := sub (len $tempHistory.Array) 1 }}
        {{- $latestTemp := index $tempHistory.Array $lastIndex }}
        {{- $latestTempValue := $latestTemp.Get "temp" }}

        <div style="margin-top: 1rem; display: flex; justify-content: space-between; align-items: center;">
          <a href="http://${SCRUTINY_URL}/web/device/{{ $wwn }}" target="_blank"> 
            <div>
              <strong class="color-highlight" style="text-transform: uppercase; font-size: 1.5rem;">
                /DEV/{{ $device }} · {{ $capacity_print }}TB · {{ $latestTempValue.Int }}°C
              </strong><br>
              <p class="color-highlight" style="margin: 0;">
                {{ $model }}
              </p>
              <div class="color-primary" style="font-size: 1.2rem;">
                Powered on for {{ $days }} days
              </div>
            </div>
          </a>

          <div style="display: flex; align-items: center; gap: 6px;">
            <a href="http://${SCRUTINY_URL}/web/device/{{ $wwn }}" target="_blank"> 
              {{- if eq $status 0.0 }}
                <span title = "Passed Health Checks" style="color: green; font-size: 18px; cursor: pointer;">●</span>
                <span style="color: green;"></span>
              {{- else }}
                <span title = "Failed Health Checks" style="color: red; font-size: 18px; cursor: pointer;">●</span>
                <span style="color: red;"></span>
              {{- end }}
            </a>
          </div>
        </div>
        {{- if lt $count $total }}
          <hr class="color-secondary" style="margin: 1rem 0; border: none; border-top: 1px solid currentColor;" />
        {{- end }}
      {{- end }}
    {{- end }}
```

## Environment Variables
- `SCRUTINY_URL` - Address to your scrutiny instance. 