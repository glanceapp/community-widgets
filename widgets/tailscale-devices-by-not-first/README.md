![](preview.png)

```yaml
        - type: custom-api
          title: Tailscale Devices
          title-url: https://login.tailscale.com/admin/machines
          options:
            auth_mode: "oauth_proxy"  # choose "api_key" or "oauth_proxy"
          # Main URL (OAuth Proxy mode)
          url: http://tailscale-token-manager:1180/devices
          cache: 30s
          subrequests:
            api_data:
              url: https://api.tailscale.com/api/v2/tailnet/-/devices
              headers:
                Authorization: Bearer ${TAILSCALE_API_KEY}
          template: |
            {{ if eq .Options.auth_mode "api_key" }}
              {{ $data := .Subrequest "api_data" }}
              {{ if ne $data.Response.StatusCode 200 }}
                <div class="widget-error-header">
                  <div class="color-negative size-h3">ERROR</div>
                  <svg class="widget-error-icon" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126ZM12 15.75h.007v.008H12v-.008Z"></path>
                  </svg>
                </div>
                <p class="break-all">Failed to fetch Tailscale devices: {{ $data.Response.Status }}</p>
              {{ else }}
                {{/* Set to true if you'd like an indicator for online devices */}}
                {{ $enableOnlineIndicator := true }}

            <style>
              .device-info-container {
                position: relative;
                overflow: hidden;
                height: 1.5em;
              }
              .device-info {
                display: flex;
                transition: transform 0.2s ease, opacity 0.2s ease;
              }
              .device-ip {
                position: absolute;
                top: 0;
                left: 0;
                transform: translateY(-100%);
                opacity: 0;
                transition: transform 0.2s ease, opacity 0.2s ease;
              }
              .device-info-container:hover .device-info {
                transform: translateY(100%);
                opacity: 0;
              }
              .device-info-container:hover .device-ip {
                transform: translateY(0);
                opacity: 1;
              }
              .update-indicator {
                width: 8px;
                height: 8px;
                border-radius: 50%;
                background-color: var(--color-warning);
                display: inline-block;
                margin-left: 4px;
                vertical-align: middle;
              }
              .offline-indicator {
                width: 8px;
                height: 8px;
                border-radius: 50%;
                background-color: var(--color-negative);
                display: inline-block;
                margin-left: 4px;
                vertical-align: middle;
              }
              .online-indicator {
                width: 8px;
                height: 8px;
                border-radius: 50%;
                background-color: var(--color-positive);
                display: inline-block;
                margin-left: 4px;
                vertical-align: middle;
              }
              .device-name-container {
                display: flex;
                align-items: center;
                gap: 8px;
              }
              .indicators-container {
                display: flex;
                align-items: center;
                gap: 4px;
              }
            </style>
            <ul class="list list-gap-10 collapsible-container" data-collapse-after="4">
              {{ range $data.JSON.Array "devices" }}
              <li>
                <div class="flex items-center gap-10">
                  <div class="device-name-container grow">
                    <span class="size-h4 block text-truncate color-primary">
                      {{ findMatch "^([^.]+)" (.String "name") }}
                    </span>
                    <div class="indicators-container">
                      {{ if (.Bool "updateAvailable") }}
                      <span class="update-indicator" data-popover-type="text" data-popover-text="Update Available"></span>
                      {{ end }}

                      {{ $lastSeen := .String "lastSeen" | parseTime "rfc3339" }}
                      {{ if not ($lastSeen.After (offsetNow "-10s")) }}
                      {{ $lastSeenTimezoned := $lastSeen.In now.Location }}
                      <span class="offline-indicator" data-popover-type="text"
                        data-popover-text="Offline - Last seen {{ $lastSeenTimezoned.Format " Jan 2 3:04pm" }}"></span>
                      {{ end }}
                    </div>
                  </div>
                </div>
                <div class="device-info-container">
                  <ul class="list-horizontal-text device-info">
                    <li>{{ .String "os" }}</li>
                    <li>{{ .String "user" }}</li>
                  </ul>
                  <div class="device-ip">
                    {{ .String "addresses.0"}}
                  </div>
                </div>
              </li>
              {{ end }}
            </ul>
              {{ end }}
            {{ else }}
              {{ if ne .Response.StatusCode 200 }}
                <div class="widget-error-header">
                  <div class="color-negative size-h3">ERROR</div>
                  <svg class="widget-error-icon" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126ZM12 15.75h.007v.008H12v-.008Z"></path>
                  </svg>
                </div>
                <p class="break-all">Failed to fetch Tailscale devices: {{ .Response.Status }}</p>
              {{ else }}
                {{/* Set to true if you'd like an indicator for online devices */}}
                {{ $enableOnlineIndicator := true }}

                <style>
                  .device-info-container {
                    position: relative;
                    overflow: hidden;
                    height: 1.5em;
                  }
                  .device-info {
                    display: flex;
                    transition: transform 0.2s ease, opacity 0.2s ease;
                  }
                  .device-ip {
                    position: absolute;
                    top: 0;
                    left: 0;
                    transform: translateY(-100%);
                    opacity: 0;
                    transition: transform 0.2s ease, opacity 0.2s ease;
                  }
                  .device-info-container:hover .device-info {
                    transform: translateY(100%);
                    opacity: 0;
                  }
                  .device-info-container:hover .device-ip {
                    transform: translateY(0);
                    opacity: 1;
                  }
                  .update-indicator {
                    width: 8px;
                    height: 8px;
                    border-radius: 50%;
                    background-color: var(--color-primary);
                    display: inline-block;
                    margin-left: 4px;
                    vertical-align: middle;
                  }
                  .offline-indicator {
                    width: 8px;
                    height: 8px;
                    border-radius: 50%;
                    background-color: var(--color-negative);
                    display: inline-block;
                    margin-left: 4px;
                    vertical-align: middle;
                  }
                  .online-indicator {
                    width: 8px;
                    height: 8px;
                    border-radius: 50%;
                    background-color: var(--color-positive);
                    display: inline-block;
                    margin-left: 4px;
                    vertical-align: middle;
                  }
                  .device-name-container {
                    display: flex;
                    align-items: center;
                    gap: 8px;
                  }
                  .indicators-container {
                    display: flex;
                    align-items: center;
                    gap: 4px;
                  }
                </style>
                <ul class="list list-gap-10 collapsible-container" data-collapse-after="4">
                  {{ range .JSON.Array "devices" }}
                  <li>
                    <div class="flex items-center gap-10">
                      <div class="device-name-container grow">
                        <span class="size-h4 block text-truncate color-primary">
                          {{ findMatch "^([^.]+)" (.String "name") }}
                        </span>
                        <div class="indicators-container">
                          {{ if (.Bool "updateAvailable") }}
                          <span class="update-indicator" data-popover-type="text" data-popover-text="Update Available"></span>
                          {{ end }}
                          {{ $lastSeen := .String "lastSeen" | parseTime "rfc3339" }}
                          {{ if not ($lastSeen.After (offsetNow "-10s")) }}
                          {{ $lastSeenTimezoned := $lastSeen.In now.Location }}
                          <span class="offline-indicator" data-popover-type="text"
                            data-popover-text="Offline - Last seen {{ $lastSeenTimezoned.Format " Jan 2 3:04pm" }}"></span>
                          {{ end }}
                        </div>
                      </div>
                    </div>
                    <div class="device-info-container">
                      <ul class="list-horizontal-text device-info">
                        <li>{{ .String "os" }}</li>
                        <li>{{ .String "user" }}</li>
                      </ul>
                      <div class="device-ip">
                        {{ .String "addresses.0"}}
                      </div>
                    </div>
                  </li>
                  {{ end }}
                </ul>
              {{ end }}
            {{ end }}

```

## Configuration Options

### `auth_mode` (required)
- `"api_key"` - Use Tailscale API key authentication (default)
- `"oauth_proxy"` - Use OAuth proxy with automatic token refresh

## Usage Examples

### API Key Mode (Default)
```yaml
options:
  auth_mode: "api_key"
```

### OAuth Proxy Mode (Automatic Token Refresh)
```yaml
options:
  auth_mode: "oauth_proxy"
```

## Authentication Modes

### API Key Mode
- Uses your Tailscale API key directly
- Requires manual renewal every 90 days
- **Environment variable required:** `TAILSCALE_API_KEY`

### OAuth Proxy Mode
- Uses the [Tailscale Token Manager](https://github.com/5at0ri/tailscale-token-manager) container
- Automatically refreshes OAuth tokens - no manual updates needed
- **Prerequisites:**
  - Deploy the `tailscale-token-manager` container
  - Configure OAuth application in Tailscale admin console
  - Ensure the token manager is accessible at `tailscale-token-manager:1180`

## Environment Variables

- `TAILSCALE_API_KEY`: Your Tailscale API key (required for `auth_mode: "api_key"`)
- `TZ`: Container timezone for correct timestamps (optional)

## Credits
[5at0ri](https://github.com/5at0ri) - Incorporated oauth into the widget as an option and made the token manager (https://github.com/5at0ri/tailscale-token-manager)
