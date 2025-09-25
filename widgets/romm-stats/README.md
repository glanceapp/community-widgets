# RomM Stats Widget

![RomM Widget Preview](https://github.com/milktoastrat/community-widgets/blob/main/widgets/romm-stats/romm-stats-example.png?raw=true)

This custom-api widget displays key statistics from a RomM server using the `/api/stats` endpoint. No authentication is required by default.

```yaml
- type: custom-api
  title: RomM
  cache: 1d
  url: http://${ROMM_URL}/api/stats
  headers:
    Accept: application/json
  template: |
    {{ $bytes := .JSON.Int "TOTAL_FILESIZE_BYTES" | toFloat }}
    {{ $tb := div $bytes 1099511627776 }}
    {{ $gb := div $bytes 1073741824 | toInt }}

    <div class="flex justify-between text-center">
      <div>
        <div class="color-highlight size-h3">{{ .JSON.Int "PLATFORMS" | formatNumber }}</div>
        <div class="size-h6">PLATFORMS</div>
      </div>
      <div>
        <div class="color-highlight size-h3">{{ .JSON.Int "ROMS" | formatNumber }}</div>
        <div class="size-h6">ROMS</div>
      </div>
      <div>
        <div class="color-highlight size-h3">
          {{ if ge $tb 1.0 }}
            {{ printf "%.2f" $tb }}TB
          {{ else }}
            {{ $gb }}GB
          {{ end }}
        </div>
        <div class="size-h6">FILESIZE</div>
      </div>
    </div>
```
Add ${ROMM_URL} to .env or replace here with your RomM URL.
This widget adapts to show TB with 2 decimals when the filesize is ≥1TB, or GB rounded if below that.
