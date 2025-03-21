![widget screenshot](preview.png)

```yaml
- type: custom-api
  cache: 1h
  title: Internet Speed
  url: http://${SPEEDTEST_API_URL}/api/v1/results/latest
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

- `SPEEDTEST_API_URL` - the URL of the Speedtest Tracker instance API
- `SPEEDTEST_TRACKER_API_TOKEN` - your Speedtest Tracker API token
