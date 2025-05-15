![widget screenshot](preview.png)

```yaml
- type: custom-api
  cache: 1h
  title: Internet Speed
  title-url: ${SPEEDTEST_URL}
  url: ${SPEEDTEST_URL}/api/v1/results/latest
  headers:
    Authorization: Bearer ${SPEEDTEST_TRACKER_API_TOKEN}
    Accept: application/json
  subrequests:
    stats:
      url: ${SPEEDTEST_URL}/api/v1/stats
      headers:
        Authorization: Bearer ${SPEEDTEST_TRACKER_API_TOKEN}
        Accept: application/json
  options:
    showPercentageDiff: true
  template: |
    {{ $showPercentage := .Options.BoolOr "showPercentageDiff" true }}

    {{ $stats := .Subrequest "stats" }}
    <div class="flex justify-between text-center margin-block-3">
    <div>
        {{ $downloadChange := percentChange ($stats.JSON.Float "data.download.avg_bits") (.JSON.Float "data.download_bits")
        }}
        {{ if $showPercentage }}
        <div
        class="size-small {{ if gt $downloadChange 0.0 }}color-positive{{ else if lt $downloadChange 0.0 }}color-negative{{ else }}color-primary{{ end }}"
        style="display: inline-flex; align-items: center;">
        {{ $downloadChange | printf "%+.1f%%" }}
        {{ if gt $downloadChange 0.0 }}
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" fill="currentColor"
            style="height: 1em; margin-left: 0.25em;" class="size-4">
            <path fill-rule="evenodd"
            d="M9.808 4.057a.75.75 0 0 1 .92-.527l3.116.849a.75.75 0 0 1 .528.915l-.823 3.121a.75.75 0 0 1-1.45-.382l.337-1.281a23.484 23.484 0 0 0-3.609 3.056.75.75 0 0 1-1.07.01L6 8.06l-3.72 3.72a.75.75 0 1 1-1.06-1.061l4.25-4.25a.75.75 0 0 1 1.06 0l1.756 1.755a25.015 25.015 0 0 1 3.508-2.85l-1.46-.398a.75.75 0 0 1-.526-.92Z"
            clip-rule="evenodd" />
        </svg>
        {{ else if lt $downloadChange 0.0 }}
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" fill="currentColor"
            style="height: 1em; margin-left: 0.25em;" class="size-4">
            <path fill-rule="evenodd"
            d="M1.22 4.22a.75.75 0 0 1 1.06 0L6 7.94l2.761-2.762a.75.75 0 0 1 1.158.12 24.9 24.9 0 0 1 2.718 5.556l.729-1.261a.75.75 0 0 1 1.299.75l-1.591 2.755a.75.75 0 0 1-1.025.275l-2.756-1.591a.75.75 0 1 1 .75-1.3l1.097.634a23.417 23.417 0 0 0-1.984-4.211L6.53 9.53a.75.75 0 0 1-1.06 0L1.22 5.28a.75.75 0 0 1 0-1.06Z"
            clip-rule="evenodd" />
        </svg>
        {{ end }}
        </div>
        {{ end }}
        <div class="color-highlight size-h3">{{ .JSON.Float "data.download_bits" | mul 0.000001 | printf "%.1f" }}</div>
        <div class="size-h6">DOWNLOAD</div>
    </div>
    <div>
        {{ $uploadChange := percentChange ($stats.JSON.Float "data.upload.avg_bits") (.JSON.Float "data.upload_bits") }}
        {{ if $showPercentage }}
        <div
        class="size-small {{ if gt $uploadChange 0.0 }}color-positive{{ else if lt $uploadChange 0.0 }}color-negative{{ else }}color-primary{{ end }}"
        style="display: inline-flex; align-items: center;">
        {{ $uploadChange | printf "%+.1f%%" }}
        {{ if gt $uploadChange 0.0 }}
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" fill="currentColor"
            style="height: 1em; margin-left: 0.25em;" class="size-4">
            <path fill-rule="evenodd"
            d="M9.808 4.057a.75.75 0 0 1 .92-.527l3.116.849a.75.75 0 0 1 .528.915l-.823 3.121a.75.75 0 0 1-1.45-.382l.337-1.281a23.484 23.484 0 0 0-3.609 3.056.75.75 0 0 1-1.07.01L6 8.06l-3.72 3.72a.75.75 0 1 1-1.06-1.061l4.25-4.25a.75.75 0 0 1 1.06 0l1.756 1.755a25.015 25.015 0 0 1 3.508-2.85l-1.46-.398a.75.75 0 0 1-.526-.92Z"
            clip-rule="evenodd" />
        </svg>
        {{ else if lt $uploadChange 0.0 }}
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" fill="currentColor"
            style="height: 1em; margin-left: 0.25em;" class="size-4">
            <path fill-rule="evenodd"
            d="M1.22 4.22a.75.75 0 0 1 1.06 0L6 7.94l2.761-2.762a.75.75 0 0 1 1.158.12 24.9 24.9 0 0 1 2.718 5.556l.729-1.261a.75.75 0 0 1 1.299.75l-1.591 2.755a.75.75 0 0 1-1.025.275l-2.756-1.591a.75.75 0 1 1 .75-1.3l1.097.634a23.417 23.417 0 0 0-1.984-4.211L6.53 9.53a.75.75 0 0 1-1.06 0L1.22 5.28a.75.75 0 0 1 0-1.06Z"
            clip-rule="evenodd" />
        </svg>
        {{ end }}
        </div>
        {{ end }}
        <div class="color-highlight size-h3">{{ .JSON.Float "data.upload_bits" | mul 0.000001 | printf "%.1f" }}</div>
        <div class="size-h6">UPLOAD</div>
    </div>
    <div>
        {{ $pingChange := percentChange ($stats.JSON.Float "data.ping.avg") (.JSON.Float "data.ping") }}
        {{ if $showPercentage }}
        <div
        class="size-small {{ if gt $pingChange 0.0 }}color-negative{{ else if lt $pingChange 0.0 }}color-positive{{ else }}color-primary{{ end }}"
        style="display: inline-flex; align-items: center;">
        {{ $pingChange | printf "%+.1f%%" }}
        {{ if lt $pingChange 0.0 }}
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" fill="currentColor"
            style="height: 1em; margin-left: 0.25em;" class="size-4">
            <path fill-rule="evenodd"
            d="M1.22 4.22a.75.75 0 0 1 1.06 0L6 7.94l2.761-2.762a.75.75 0 0 1 1.158.12 24.9 24.9 0 0 1 2.718 5.556l.729-1.261a.75.75 0 0 1 1.299.75l-1.591 2.755a.75.75 0 0 1-1.025.275l-2.756-1.591a.75.75 0 1 1 .75-1.3l1.097.634a23.417 23.417 0 0 0-1.984-4.211L6.53 9.53a.75.75 0 0 1-1.06 0L1.22 5.28a.75.75 0 0 1 0-1.06Z"
            clip-rule="evenodd" />
        </svg>
        {{ else if gt $pingChange 0.0 }}
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" fill="currentColor"
            style="height: 1em; margin-left: 0.25em;" class="size-4">
            <path fill-rule="evenodd"
            d="M9.808 4.057a.75.75 0 0 1 .92-.527l3.116.849a.75.75 0 0 1 .528.915l-.823 3.121a.75.75 0 0 1-1.45-.382l.337-1.281a23.484 23.484 0 0 0-3.609 3.056.75.75 0 0 1-1.07.01L6 8.06l-3.72 3.72a.75.75 0 1 1-1.06-1.061l4.25-4.25a.75.75 0 0 1 1.06 0l1.756 1.755a25.015 25.015 0 0 1 3.508-2.85l-1.46-.398a.75.75 0 0 1-.526-.92Z"
            clip-rule="evenodd" />
        </svg>
        {{ end }}
        </div>
        {{ end }}
        <div class="color-highlight size-h3">{{ .JSON.Float "data.ping" | printf "%.0f ms" }}</div>
        <div class="size-h6">PING</div>
    </div>
    </div>
```

## Options
- `showPercentageDiff`: Determines if the percentage difference of the current test from the average is shown. Default is true

## Environment variables

- `SPEEDTEST_TRACKER_URL` - The URL of the Speedtest Tracker instance (e.g `http://my.speedtest-tracker.instance`)
- `SPEEDTEST_TRACKER_API_TOKEN` - Your Speedtest Tracker API token
