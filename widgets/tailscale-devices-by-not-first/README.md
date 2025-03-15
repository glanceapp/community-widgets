![widget example](preview.png)

```yaml
- type: custom-api
  cache: 1d
  title: Tailscale Devices
  url: https://api.tailscale.com/api/v2/tailnet/-/devices
  headers:
    Authorization: Bearer ${TAILSCALE_API_KEY}
  template: |
    <ul class="list list-gap-10 collapsible-container" data-collapse-after="4">
      {{ range .JSON.Array "devices" }}
      <li>
        <div class="flex items-center gap-10">
          <span class="size-h4 block text-truncate color-primary"
                data-popover-type="text"
                data-popover-text="{{ .String "addresses.0" }}">
            {{ .String "hostname" }}
          </span>
        </div>
        <ul class="list-horizontal-text">
          <li>{{ .String "os" }}</li>
          <li>{{ .String "user" }}</li>
        </ul>
      </li>
      {{ end }}
    </ul>
```

## Environment variables

- `TAILSCALE_API_KEY` - a Tailscale API Key with permissions to access tailnet device list
