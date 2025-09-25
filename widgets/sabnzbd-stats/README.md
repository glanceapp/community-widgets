![](preview.png)


```yaml
- type: custom-api
  title: SABnzbd Status
  cache: 30s
  url: http://YOUR_SERVER:PORT/api?output=json&apikey=${SABNZBD_API_KEY}&mode=queue
  headers:
    Accept: application/json
  template: |
    <div class="p-2">
      <div class="flex justify-between mb-2">
        <div class="flex-1">
          <div class="size-h6">SPEED</div>
          <div class="color-highlight size-h3">{{ if eq (.JSON.String "queue.status") "Downloading" }}{{ .JSON.String "queue.speed" }}B/s{{ else }}Paused{{ end }}</div>
        </div>
        <div class="flex-1">
          <div class="size-h6">TIME LEFT</div>
          <div class="color-highlight size-h3">{{ if eq (.JSON.String "queue.status") "Downloading" }}{{ .JSON.String "queue.timeleft" }}{{ else }}--:--:--{{ end }}</div>
        </div>
      </div>
      <div class="flex justify-between mb-2">
        <div class="flex-1">
          <div class="size-h6">QUEUE</div>
          <div class="color-highlight size-h3">{{ .JSON.Int "queue.noofslots" }} items</div>
        </div>
        <div class="flex-1">
          <div class="size-h6">SIZE</div>
          <div class="color-highlight size-h3">{{ .JSON.Float "queue.mb" | printf "%.1f" }}MB</div>
        </div>
      </div>
    </div>
```
## Environment variables

- `SABNZBD_API_KEY` - the API key of SABnzbd can be found under `Settings -> General`
