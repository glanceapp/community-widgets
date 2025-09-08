# Tailscale Devices Widget

A widget for Glance that displays all devices in your Tailscale tailnet, showing their connection status, update availability, and IP addresses.

## Option 1: API Key (Recommended)

The standard implementation using Tailscale API keys. Simple to set up but requires manual key renewal every 90 days maximum.

### Setup

1. **Generate API Key:**
   - Go to [Tailscale admin console](https://login.tailscale.com/admin/settings/keys)
   - Generate a new API key (maximum 90-day expiration)
   - Copy the key

2. **Configure Widget:**
   - Set the `TAILSCALE_API_KEY` environment variable to your generated key
   - `TZ`: For correct times, the widget uses the container's timezone. If not already supplied, you can use this variable to provide your timezone.

### Configuration

```yaml
- type: custom-api
  title: Tailscale Devices
  title-url: https://login.tailscale.com/admin/machines
  url: https://api.tailscale.com/api/v2/tailnet/-/devices
  headers:
    Authorization: Bearer ${TAILSCALE_API_KEY}
  cache: 10m
  template: |
    {{/* User Variables */}}
    {{/* Set to true if you'd like an indicator for online devices */}}
    {{ $enableOnlineIndicator := false }}

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
              {{ else if $enableOnlineIndicator }}
                <span class="online-indicator" data-popover-type="text" data-popover-text="Online"></span>
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
```

## Option 2: OAuth Proxy (Alternative)

An alternative implementation that uses automatically refreshing OAuth tokens, eliminating the need for manual key renewal. Created by @5at0ri, requires [5at0ri's OAuth Proxy](https://github.com/5at0ri/tailscale-token-manager).

### Setup

1. **Deploy OAuth Proxy:**
   - Set up OAuth client credentials in your [Tailscale admin console](https://login.tailscale.com/admin/settings/oauth)
   - Deploy the [OAuth Proxy container](https://github.com/5at0ri/tailscale-token-manager)
   - Configure it with your OAuth Client ID and Client Secret

2. **Configure Widget:**
   - Point the `url` to your OAuth proxy container endpoint
   - The proxy handles all authentication automatically

### Configuration

```yaml
- type: custom-api
  title: Tailscale Devices
  title-url: https://login.tailscale.com/admin/machines
  url: http://tailscale-token-manager:1180/devices  # URL to your OAuth proxy container
  cache: 10s
  template: |
    {{/* User Variables */}}
    {{/* Set to true if you'd like an indicator for online devices */}}
    {{ $enableOnlineIndicator := false }}

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
              {{ else if $enableOnlineIndicator }}
                <span class="online-indicator" data-popover-type="text" data-popover-text="Online"></span>
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
```

## Notes

- **API Key Mode:** Simpler setup, requires manual renewal every 90 days
- **OAuth Proxy Mode:** Automatic token refresh, requires additional container setup
- Both implementations provide identical device information and functionality
