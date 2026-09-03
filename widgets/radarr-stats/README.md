# Radarr Stats

Total movie library count, missing movies, and active download queue.

![](preview.png)

```yaml
- type: custom-api
  refresh-interval: 5s
  title: Radarr Stats
  cache: 2s
  url: ${RADARR_URL}/api/v3/movie?apikey=${RADARR_API_KEY}
  subrequests:
    missing:
      url: ${RADARR_URL}/api/v3/wanted/missing?apikey=${RADARR_API_KEY}
    queue:
      url: ${RADARR_URL}/api/v3/queue?apikey=${RADARR_API_KEY}
  template: |
    <div class="flex justify-between text-center">
      <div>
        <div class="color-highlight size-h3">{{ len (.JSON.Array "") }}</div>
        <div class="size-h6">MOVIES</div>
      </div>
      <div>
        <div class="color-highlight size-h3">{{ (.Subrequest "missing").JSON.Int "totalRecords" }}</div>
        <div class="size-h6">MISSING</div>
      </div>
      <div>
        <div class="color-highlight size-h3">{{ (.Subrequest "queue").JSON.Int "totalRecords" }}</div>
        <div class="size-h6">QUEUED</div>
      </div>
    </div>
```

## Environment variables

- `RADARR_URL` - The full URL of your Radarr instance, e.g. `http://192.168.1.100:7878`
- `RADARR_API_KEY` - Your Radarr API key, found in Settings -> General -> Security
