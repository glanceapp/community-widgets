# Healthchecks.io Glance Widget

A Glance widget for displaying the status of checks at healthchecks.io

### Code

``` yml
- type: custom-api
    title: HealthChecks.IO
    cache: 10m
    url: https://healthchecks.io/api/v3/checks/
    headers:
    X-Api-Key: ${HEALTHCHECKS_API_KEY}
    Content-Type: "application/json"
    template: |
    {{ if eq .Response.StatusCode 200 }}
        {{ $checks := .JSON.Array "checks" }}
        {{ if eq (len $checks) 0 }}
        <p class="color-subdue">No checks configured</p>
        {{ else }}
        <ul class="list list-gap-10 collapsible-container" data-collapse-after="6">
        {{ range $checks }}
        {{ $status := .String "status" }}
        <li class="flex items-center justify-between gap-10">
            <div class="min-width-0">
            <div class="size-h4 color-highlight text-truncate">{{ .String "name" }}</div>
            <div class="size-h6 color-subdue">
                {{ if ne (.String "last_ping") "" }}
                Last ping <span {{ .String "last_ping" | parseTime "rfc3339" | toRelativeTime }}></span>
                {{ else }}
                Never pinged
                {{ end }}
            </div>
            </div>
            <div class="shrink-0 size-h5 text-right">
            {{ if eq $status "up" }}
                <span class="color-positive">● Up</span>
            {{ else if eq $status "down" }}
                <span class="color-negative">● Down</span>
            {{ else if eq $status "grace" }}
                <span class="color-primary">● Grace</span>
            {{ else if eq $status "paused" }}
                <span class="color-subdue">● Paused</span>
            {{ else if eq $status "started" }}
                <span class="color-primary">● Running</span>
            {{ else }}
                <span class="color-subdue">● {{ $status }}</span>
            {{ end }}
            </div>
        </li>
        {{ end }}
        </ul>
        {{ end }}
    {{ else }}
        <p class="color-negative">Error fetching data ({{ .Response.StatusCode }})</p>
    {{ end }}
```
            
