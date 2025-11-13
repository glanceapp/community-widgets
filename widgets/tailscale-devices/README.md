# Tailscale Devices Widget

![](preview.png)

A widget for Glance that displays all devices in your Tailscale tailnet, showing their connection status, update availability, and IP addresses.

## Setup

### Option 1: API Key (Recommended)

The standard implementation using Tailscale API keys. Simple to set up but requires manual API key renewal every 90 days maximum.

1. **Generate API Key:**
   - Go to [Tailscale admin console](https://login.tailscale.com/admin/settings/keys)
   - Generate a new API key (maximum 90-day expiration)
   - Copy the key

2. **Configure Widget:**
   - Set the `TAILSCALE_API_KEY` environment variable to your generated key

### Option 2: OAuth Proxy (Alternative)

An alternative implementation that uses automatically refreshing OAuth tokens, eliminating the need for manual key renewal. Created by @5at0ri, requires [5at0ri's OAuth Proxy](https://github.com/5at0ri/tailscale-token-manager).

1. **Deploy OAuth Proxy:**
   - Set up OAuth client credentials in your [Tailscale admin console](https://login.tailscale.com/admin/settings/oauth)
   - Deploy the [OAuth Proxy container](https://github.com/5at0ri/tailscale-token-manager)
   - Configure it with your OAuth Client ID and Client Secret

2. **Configure Widget:**
   - Replace the `url` line with: `url: http://tailscale-token-manager:1180/devices`
   - Remove the `headers:` section entirely
   - Change `cache:` to `10s` for more frequent updates

## Template

```yaml
- type: custom-api
  title: Tailscale Devices
  title-url: https://login.tailscale.com/admin/machines
  url: https://api.tailscale.com/api/v2/tailnet/-/devices
  headers:
    Authorization: Bearer ${TAILSCALE_API_KEY}
  cache: 10m
  options:
    collapseAfter: 3
    # disableOfflineIndicator: true
    # disableUpdateIndicator: true
    # prioritiseTags: true
  template: |
    <style>
      .device-info-container-tailscale {
        position: relative;
        overflow: hidden;
        height: 1.5em;
      }

      .device-info-tailscale {
        display: flex;
        transition: transform 0.2s ease, opacity 0.2s ease;
      }

      .device-ip-tailscale {
        position: absolute;
        top: 0;
        left: 0;
        transform: translateY(-100%);
        opacity: 0;
        transition: transform 0.2s ease, opacity 0.2s ease;
      }

      .device-info-container-tailscale:hover .device-info-tailscale {
        transform: translateY(100%);
        opacity: 0;
      }

      .device-info-container-tailscale:hover .device-ip-tailscale {
        transform: translateY(0);
        opacity: 1;
      }

      .update-indicator-tailscale {
        width: 8px;
        height: 8px;
        border-radius: 50%;
        background-color: var(--color-primary);
        display: inline-block;
        margin-left: 4px;
        vertical-align: middle;
      }

      .offline-indicator-tailscale {
        width: 8px;
        height: 8px;
        border-radius: 50%;
        background-color: var(--color-negative);
        display: inline-block;
        margin-left: 4px;
        vertical-align: middle;
      }

      .device-name-container-tailscale {
        display: flex;
        align-items: center;
        gap: 8px;
      }

      .indicators-container-tailscale {
        display: flex;
        align-items: center;
        gap: 4px;
      }
    </style>
    <ul class="list list-gap-10 collapsible-container" data-collapse-after="{{ .Options.IntOr "collapseAfter" 3 }}">
      {{ range .JSON.Array "devices" }}
      <li>
        <div class="flex items-center gap-10">
          <div class="device-name-container-tailscale grow">
            <span class="size-h4 block text-truncate color-primary">
              {{ findMatch "^([^.]+)" (.String "name") }}
            </span>
            <div class="indicators-container-tailscale">
              {{ if and (not ($.Options.BoolOr "disableUpdateIndicator" false)) (.Bool "updateAvailable") }}
              <span class="update-indicator-tailscale" data-popover-type="text" data-popover-text="Update Available"></span>
              {{ end }}

              {{ if not ($.Options.BoolOr "disableOfflineIndicator" false) }}
              {{ $lastSeen := .String "lastSeen" | parseTime "rfc3339" }}
              {{ if not ($lastSeen.After (offsetNow "-10s")) }}
              {{ $lastSeenTimezoned := $lastSeen.In now.Location }}
              <span class="offline-indicator-tailscale" data-popover-type="text"
                data-popover-text="Offline - Last seen {{ $lastSeenTimezoned.Format " Jan 2 3:04pm" }}"></span>
              {{ end }}
              {{ end }}

            </div>
          </div>
        </div>
        <div class="device-info-container-tailscale">
          <ul class="list-horizontal-text device-info-tailscale">
            <li>{{ .String "os" }}</li>
            <li>
              {{ if and ($.Options.BoolOr "prioritiseTags" false) (.Exists "tags.0") }}
                {{ trimPrefix "tag:" (.String "tags.0") }}
              {{ else }}
                {{ .String "user" }}
              {{ end }}
            </li>
          </ul>
          <div class="device-ip-tailscale">
            {{ .String "addresses.0"}}
          </div>
        </div>
      </li>
      {{ end }}
    </ul>
```

## Available Options

| Option                    | Type    | Default | Description                                                                                  |
| ------------------------- | ------- | ------- | -------------------------------------------------------------------------------------------- |
| `collapseAfter`           | integer | `3`     | Number of devices to show before collapsing the list.                                        |
| `disableOfflineIndicator` | boolean | `false` | When set to `true`, hides the red dot indicator for offline devices.                         |
| `disableUpdateIndicator`  | boolean | `false` | When set to `true`, hides the blue dot indicator when updates are available.                 |
| `prioritiseTags`          | boolean | `false` | When set to `true`, displays the device's primary tag instead of the user (if a tag exists). |
