![](preview.png)

```yaml
- type: custom-api
  title: Uptime Kumas
  title-url: ${UPTIME_KUMA_URL}
  url: ${UPTIME_KUMA_URL}/api/status-page/${UPTIME_KUMA_STATUS_SLUG}
  subrequests:
    heartbeats:
      url: ${UPTIME_KUMA_URL}/api/status-page/heartbeat/${UPTIME_KUMA_STATUS_SLUG}
  cache: 10m
  template: |
    {{ $hb := .Subrequest "heartbeats" }}

    {{ if not (.JSON.Exists "publicGroupList") }}
    <p class="color-negative">Error reading response</p>
    {{ else if eq (len (.JSON.Array "publicGroupList")) 0 }}
    <p>No monitors found</p>
    {{ else }}

    <ul class="dynamic-columns list-gap-8">
      {{ range .JSON.Array "publicGroupList" }}
      {{ range .Array "monitorList" }}
      {{ $id := .String "id" }}
      {{ $hbArray := $hb.JSON.Array (print "heartbeatList." $id) }}
      <div class="flex items-center gap-12">
        <a class="size-title-dynamic color-highlight text-truncate block grow" href="${UPTIME_KUMA_URL}/dashboard/{{ $id }}"
          target="_blank" rel="noreferrer">
          {{ .String "name" }} </a>
    
        {{ if gt (len $hbArray) 0 }}
          {{ $latest := index $hbArray (sub (len $hbArray) 1) }}
          {{ if eq ($latest.Int "status") 1 }}
          <div>{{ $latest.Int "ping" }}ms</div>
          <div class="monitor-site-status-icon-compact" title="OK">
            <svg fill="var(--color-positive)" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20">
              <path fill-rule="evenodd"
                d="M10 18a8 8 0 1 0 0-16 8 8 0 0 0 0 16Zm3.857-9.809a.75.75 0 0 0-1.214-.882l-3.483 4.79-1.88-1.88a.75.75 0 1 0-1.06 1.061l2.5 2.5a.75.75 0 0 0 1.137-.089l4-5.5Z"
                clip-rule="evenodd"></path>
            </svg>
          </div>
          {{ else }}
          <div><span class="color-negative">DOWN</span></div>
          <div class="monitor-site-status-icon-compact" title="{{ if $latest.Exists "msg" }}{{ $latest.String "msg" }}{{ else
            }}Error{{ end }}">
            <svg fill="var(--color-negative)" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20">
              <path fill-rule="evenodd"
                d="M8.485 2.495c.673-1.167 2.357-1.167 3.03 0l6.28 10.875c.673 1.167-.17 2.625-1.516 2.625H3.72c-1.347 0-2.189-1.458-1.515-2.625L8.485 2.495ZM10 5a.75.75 0 0 1 .75.75v3.5a.75.75 0 0 1-1.5 0v-3.5A.75.75 0 0 1 10 5Zm0 9a1 1 0 1 0 0-2 1 1 0 0 0 0 2Z"
                clip-rule="evenodd"></path>
            </svg>
          </div>
          {{ end }}
        {{ else }}
          <div><span class="color-negative">No data</span></div>
          <div class="monitor-site-status-icon-compact" title="No data available">
            <svg fill="var(--color-negative)" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20">
              <path d="M10 18a8 8 0 1 1 0-16 8 8 0 0 1 0 16zm0-2a6 6 0 1 0 0-12 6 6 0 0 0 0 12zm-.75-8a.75.75 0 0 1 1.5 0v3a.75.75 0 0 1-1.5 0V8zm.75 5a1 1 0 1 1 0 2 1 1 0 0 1 0-2z"/>
            </svg>
          </div>
        {{ end }}
      </div>
      {{ end }}
      {{ end }}
    </ul>
    {{ end }}

```

## Environment variables

- `UPTIME_KUMA_URL`: The url of your uptime kuma instace (no slash on the end). For example, `http://my.uptime-kuma.instance`
- `UPTIME_KUMA_STATUS_SLUG`: The slug of your uptime kuma status page. For example, if your status page url was `http://my.uptime-kuma.instance/status/example` this variable would be set to `example`.
