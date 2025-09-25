![](preview.png)

```yaml
- type: custom-api
  title: AzireVPN Connection Status
  cache: 1h
  url: https://c.azi.re/v1/network-details
  options:
    hide_ip: "${AZIREVPN_HIDE_IP}"
  template: |
    {{ $hideIP := .Options.StringOr "hide_ip" "false" }}

    <ul class="list list-gap-10 collapsible-container">
    {{ range .JSON.Array "" }}
      {{ $isEnabled := .String "vpn_enabled" }}

      <li class="flex items-center gap-12">
        <span class="color-highlight text-truncate block grow">Status</span>
        {{ if eq $isEnabled "true" }}
          <div class="margin-left-auto shrink-0">
            <div class="monitor-site-status-icon-compact" title="Connected">
              <svg fill="var(--color-positive)" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M10 18a8 8 0 1 0 0-16 8 8 0 0 0 0 16Zm3.857-9.809a.75.75 0 0 0-1.214-.882l-3.483 4.79-1.88-1.88a.75.75 0 1 0-1.06 1.061l2.5 2.5a.75.75 0 0 0 1.137-.089l4-5.5Z" clip-rule="evenodd"></path>
              </svg>
            </div>
          </div>
          <span class="color-highlight">Connected</span>
        {{ else }}
          <div class="margin-left-auto shrink-0">
            <div class="monitor-site-status-icon-compact" title="Disconnected">
              <svg fill="var(--color-negative)" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M8.485 2.495c.673-1.167 2.357-1.167 3.03 0l6.28 10.875c.673 1.167-.17 2.625-1.516 2.625H3.72c-1.347 0-2.189-1.458-1.515-2.625L8.485 2.495ZM10 5a.75.75 0 0 1 .75.75v3.5a.75.75 0 0 1-1.5 0v-3.5A.75.75 0 0 1 10 5Zm0 9a1 1 0 1 0 0-2 1 1 0 0 0 0 2Z" clip-rule="evenodd"></path>
              </svg>
            </div>
          </div>
          <span class="color-negative">Disconnected</span>
        {{ end }}
      </li>

      <li class="flex items-center gap-12">
        <span class="color-highlight text-truncate block grow">IP</span>

        {{ if eq $hideIP "true" }}
          <span class="text-muted">Redacted</span>
        {{ else }}
          <span>{{ .String "ip" }}</span>
        {{ end }}
        
      </li>

      <li class="flex items-center gap-12">
        <span class="color-highlight text-truncate block grow">Geo</span>
        <span>{{ .String "geo.city" }}, {{ .String "geo.country" }}</span>
      </li>
    {{ end }}
    </ul>
```
## Environment variables
- `AZIREVPN_HIDE_IP` boolean