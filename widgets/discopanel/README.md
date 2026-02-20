# Discopanel Widget
Lists all servers in your Discopanel host. At current the functionality of this widget does not require authentication even when auth on the host is enabled. When Discopanel 2.0.0 releases this may change. Server icon functionality requires using a url for the picture. These can be obtained from the curseforge or modrinth page directly. I imagine that using direct file paths would also work, you just might have to adjust how the picture file path is read in glance.

# Environment Variables
`DISCOPANEL_HOST` = URL/IP of Discopanel host. Requires http/https prefix. 

# Preview
<img width="926" height="677" alt="Preview Image" src="preview.png"/>

# YAML
```yaml
- type: custom-api
  title: Discopanel
  cache: 30s
  url: ${DISCOPANEL_HOST}/discopanel.v1.ServerService/ListServers
  headers:
    Connect-Protocol-Version: 1
    Connect-Timeout-Ms: 2000
    Content-Type: application/json
  method: POST
  body-type: json
  body:
    fullStats: true
  skip-json-validation: true
  template: |
    {{ range .JSON.Array "servers" }}
      {{ $mem_usage := .Float "memoryUsage" }}
      {{ $cpu := .Float "cpuPercent" }}
      {{ $players := .Int "playersOnline" }}
      {{ $tps := .Float "tps" }}
      {{ $id := .String "id" }}
      {{ $status := .String "status" }}

      {{ $serverConfig := newRequest "${DISCOPANEL_HOST}/discopanel.v1.ConfigService/GetServerConfig"
        | withHeader "Connect-Protocol-Version" "1"
        | withHeader "Connect-Timeout-Ms" "2000"
        | withHeader "Content-Type" "application/json"
        | withStringBody (printf `{"serverId": %q}` $id)
        | getResponse
      }}
      {{ $icon := $serverConfig.JSON.String "categories.1.properties.6.value" }}

      <div style="margin-bottom: 1rem;">
        <div style="display: flex; align-items: center; gap: 1rem;">
          <a href="${DISCOPANEL_HOST}/servers/{{ $id }}"><img style="width=64px; height=64px;" src="{{ $icon }}" alt="no picture" width="64" height="64"/></a>
          <div style="flex-grow:1; min-width:0;">
            <span><a href="${DISCOPANEL_HOST}/servers/{{ $id }}" class="color-highlight size-h1">{{ .String "name" }}</a></span>
            {{ if eq "SERVER_STATUS_RUNNING" $status }}
              <span style="display: inline-block; background-color: var(--color-positive); width: 12px; height: 12px; border-radius: 50%;"></span>
            {{ else if eq "SERVER_STATUS_STARTING" $status }}
              <span style="display: inline-block; background-color: var(--color-separator); width: 12px; height: 12px; border-radius: 50%;"></span>
            {{ else }}
              <span style="display: inline-block; background-color: var(--color-negative); width: 12px; height: 12px; border-radius: 50%;"></span>
            {{ end }}
            <p class="size-h5">{{ .String "description" }}</p>
            <div class="size-h6">
              <span>Memory: {{ formatNumber $mem_usage }}MB</span>
              <span>CPU: {{ printf "%.2f" $cpu }}%</span>
              <span>Players: {{ $players }}</span>
              <span>TPS: {{ $tps }}/20</span>
            </div>
          </div>
        </div>
      </div>
    {{ end }}

```
