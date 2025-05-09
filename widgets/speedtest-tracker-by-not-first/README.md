![widget screenshot](preview.png)

```yaml
- type: custom-api
  cache: 1h
  title: Internet Speed
  title-url: ${SPEEDTEST_TRACKER_URL}
  url: ${SPEEDTEST_TRACKER_URL}/api/v1/results/latest
  headers:
    Authorization: Bearer ${SPEEDTEST_TRACKER_API_TOKEN}
    Accept: application/json
  template: |
    <div class="flex justify-between text-center margin-block-3">
        <div>
            <div class="color-highlight size-h3">{{ .JSON.Float "data.download_bits" | mul 0.000001 | printf "%.1f" }}</div>
            <div class="size-h6">DOWNLOAD</div>
        </div>
        <div>
            <div class="color-highlight size-h3">{{ .JSON.Float "data.upload_bits" | mul 0.000001 | printf "%.1f" }}</div>
            <div class="size-h6">UPLOAD</div>
        </div>
        <div>
            <div class="color-highlight size-h3">{{ .JSON.Float "data.ping" | printf "%.0f ms" }}</div>
            <div class="size-h6">PING</div>
        </div>
    </div>
```

## Environment variables

- `SPEEDTEST_TRACKER_URL` - The URL of the Speedtest Tracker instance (e.g `http://my.speedtest-tracker.instance`)
- `SPEEDTEST_TRACKER_API_TOKEN` - Your Speedtest Tracker API token
