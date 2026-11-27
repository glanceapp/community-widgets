# RackNerd

![RackNerd VPS status with sample usage](preview.png)

Show one RackNerd VPS: hostname, power state, addresses, operating system, and bandwidth, disk, and memory usage. Glance calls the per-server client API directly. The preview uses a fictional server.

The widget requests `info` only. It does not boot, reboot, or shut down the VPS.

## Set up access

1. In the RackNerd panel, open the VPS and enable its API. Copy that server's API key and API hash. Each key controls only that VPS.
2. Pass the key and hash to the **Glance process or container**. With Compose, list them under the Glance service's `environment` or `env_file`.

   | Variable | Value |
   | --- | --- |
   | `RACKNERD_API_KEY` | API key for this VPS |
   | `RACKNERD_API_HASH` | API hash for this VPS |

For another VPS, add a second copy of the widget and point it at a different pair of variables. You can also override `panel-url` or `endpoint` in `options` if your VPS uses a different SolusVM panel host.

## What the numbers cover

One request asks for status, addresses, bandwidth, disk, and memory. Each meter puts the API percent beside the label, draws a full-width bar (highlighting at `>= 80%`), and shows used / total underneath. Meters with missing or zero capacity are skipped automatically. Those byte counters use decimal units (powers of 1,000), with one decimal place. The image template name stays on the address line as a tooltip. A rejected key is shown as “API authentication failed.” The same key can also power the VPS on or off from the RackNerd panel; this widget never sends those actions.

The card caches results for five minutes.

## Glance configuration

Tested with Glance v0.8.8. Copy the YAML below and include it in the page’s widgets list.

```yaml
- type: custom-api
  title: RackNerd
  title-url: https://ctrl.racknerd.com/
  cache: 5m
  options:
    api-key: "${RACKNERD_API_KEY}"
    api-hash: "${RACKNERD_API_HASH}"
    # panel-url: https://ctrl.racknerd.com/
    # endpoint: https://ctrl.racknerd.com/api/client/command.php
  template: |
    {{ define "rnSize" }}
      {{- $n := . -}}
      {{- $l := len $n -}}
      {{- if or (eq $n "") (le $l 3) -}}{{ if eq $n "" }}0{{ else }}{{ $n }}{{ end }} B
      {{- else -}}
        {{- $cut := 3 -}}{{- $unit := "KB" -}}
        {{- if gt $l 12 -}}{{- $cut = 12 -}}{{- $unit = "TB" -}}
        {{- else if gt $l 9 -}}{{- $cut = 9 -}}{{- $unit = "GB" -}}
        {{- else if gt $l 6 -}}{{- $cut = 6 -}}{{- $unit = "MB" -}}
        {{- end -}}
        {{- $whole := sub $l $cut -}}
        {{- slice $n 0 $whole -}}.{{- slice $n $whole (add $whole 1) }} {{ $unit -}}
      {{- end -}}
    {{ end }}
    {{ define "rnUsage" }}
      {{- $label := findSubmatch "^([^:]+)" . -}}
      {{- $raw := trimPrefix (printf "%s:" $label) . -}}
      {{- if $raw -}}
        {{- $used := findSubmatch "^[^,]+,([0-9]+)" $raw -}}
        {{- $total := findSubmatch "^([0-9]+)" $raw -}}
        {{- $pct := findSubmatch "^[^,]+,[^,]+,[^,]+,([0-9.]+)" $raw -}}
        {{- if and (ne $total "") (ne $total "0") (ne $pct "") -}}
          {{- $high := findMatch "^(8[0-9]|9[0-9]|[1-9][0-9]{2,})(\\.|$)" $pct -}}
          <div class="margin-top-15">
            <div class="flex justify-between gap-10">
              <span>{{ $label }}</span>
              <span class="{{ if $high }}color-negative{{ else }}color-highlight{{ end }}" style="white-space: nowrap;">{{ $pct }}%</span>
            </div>
            <div style="height: 6px; margin-top: 7px; border-radius: 99px; background: var(--color-separator); overflow: hidden;">
              <div style="height: 100%; width: {{ $pct }}%; max-width: 100%; border-radius: inherit; background: {{ if $high }}var(--color-negative){{ else }}var(--color-primary){{ end }};"></div>
            </div>
            <div class="size-h6 color-subdue margin-top-5">{{ template "rnSize" $used }} / {{ template "rnSize" $total }}</div>
          </div>
        {{- end -}}
      {{- end -}}
    {{ end }}
    {{ if or (eq (.Options.StringOr "api-key" "") "") (eq (.Options.StringOr "api-hash" "") "") }}
      <p class="color-subdue">Add the RackNerd API key and hash to show this VPS.</p>
    {{ else }}
      {{ $panelUrl := .Options.StringOr "panel-url" "https://ctrl.racknerd.com/" }}
      {{ $endpoint := .Options.StringOr "endpoint" "https://ctrl.racknerd.com/api/client/command.php" }}
      {{ $body := printf "key=%s&hash=%s&action=info&rdtype=json&status=1&ipaddr=1&bw=1&hdd=1&mem=1"
           (urlquery (.Options.StringOr "api-key" "")) (urlquery (.Options.StringOr "api-hash" "")) }}
      {{ $res := newRequest $endpoint
           | withHeader "Content-Type" "application/x-www-form-urlencoded"
           | withHeader "Accept" "application/json"
           | withHeader "User-Agent" "glance-community-racknerd/1.0"
           | withStringBody $body
           | getResponse }}
      {{ if eq $res.Response.StatusCode 0 }}
        <p class="color-subdue">RackNerd temporarily unavailable.</p>
      {{ else if ne $res.Response.StatusCode 200 }}
        <p class="color-negative">RackNerd returned HTTP {{ $res.Response.StatusCode }}.</p>
      {{ else if eq ($res.JSON.String "status") "error" }}
        <p class="color-negative">
          {{- if findMatch "(?i)invalid (key|hash)" ($res.JSON.String "statusmsg") -}}
            API authentication failed.
          {{- else -}}
            {{ $res.JSON.String "statusmsg" }}
          {{- end -}}
        </p>
      {{ else }}
        {{ $status := $res.JSON.String "vmstat" }}
        {{ $ip := $res.JSON.String "ipaddress" }}
        {{ $ips := $res.JSON.String "ipaddr" }}
        <div class="flex justify-between gap-10">
          <a class="color-highlight text-truncate" href="{{ $panelUrl }}" target="_blank" rel="noopener noreferrer">{{ $res.JSON.String "hostname" }}</a>
          {{ if eq $status "online" }}
            <span class="color-positive" style="white-space: nowrap;">online</span>
          {{ else if or (eq $status "offline") (eq $status "disabled") }}
            <span class="color-negative" style="white-space: nowrap;">{{ $status }}</span>
          {{ else }}
            <span style="white-space: nowrap;">{{ $status }}</span>
          {{ end }}
        </div>
        {{ if or $ip ($res.JSON.String "os") }}
          <div class="size-h6 color-subdue margin-top-5" title="{{ $res.JSON.String "template" }}">{{ $ip }}{{ if and $ip ($res.JSON.String "os") }} · {{ end }}{{ $res.JSON.String "os" }}</div>
        {{ end }}
        {{ if and $ips (ne $ips $ip) }}
          <div class="size-h6 color-subdue">{{ $ips }}</div>
        {{ end }}
        {{ template "rnUsage" (printf "Bandwidth:%s" ($res.JSON.String "bw")) }}
        {{ template "rnUsage" (printf "Disk:%s" ($res.JSON.String "hdd")) }}
        {{ template "rnUsage" (printf "Memory:%s" ($res.JSON.String "mem")) }}
      {{ end }}
    {{ end }}
```
