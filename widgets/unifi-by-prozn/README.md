# Unifi Widgets

## Required environment variables

`UNIFI_CONSOLE_IP` - IP Address of your Unifi gateway
`UNIFI_API_KEY` - Unifi local API key (Site Settings -> Control Plane -> Integrations)

Alternatively you can hard-code into the url: and X-API-KEY: fields.

## Small

![](small-preview.png)

```yaml
- type: custom-api
  title: Unifi
  cache: 1m
  allow-insecure: true
  url: https://${UNIFI_CONSOLE_IP}/proxy/network/api/stat/sites
  headers:
    X-API-KEY: ${UNIFI_API_KEY}
    Accept: application/json
  template: | 
    <style>
      .list-horizontal-text.no-bullets-unifi > *:not(:last-child)::after {
          content: none !important;
      }
      .list-horizontal-text.no-bullets-unifi > *:not(:last-child) {
        margin-right: 1em;
      }
    </style>
    {{ range .JSON.Array "data" }}
      <div style="display:flex; align-items:center; gap:12px;">
        <div style="width:40px; height:40px; flex-shrink:0; display:flex; justify-content:center; align-items:center; overflow:hidden;">
          <img src="https://cdn.jsdelivr.net/gh/selfhst/icons/svg/ubiquiti-unifi-light.svg" width="24" height="24" style="object-fit:contain;">
        </div>
        <div style="flex-grow:1; min-width:0;">
          <a class="size-h4 block text-truncate color-highlight">
            {{ .String "health.#(subsystem=\"wan\").gw_name" }}
            {{ if eq (.String "health.#(subsystem=\"wan\").status") "ok" }}
            <span
              style="width: 8px; height: 8px; border-radius: 50%; background-color: var(--color-positive); display: inline-block; vertical-align: middle;"
              data-popover-type="text"
              data-popover-text="WAN Connected"
            ></span>
            {{ else }}
            <span
              style="width: 8px; height: 8px; border-radius: 50%; background-color: var(--color-negative); display: inline-block; vertical-align: middle;"
              data-popover-type="text"
              data-popover-text="WAN Disconnected"
            ></span>
            {{ end }}
          </a>
          <ul class="list-horizontal-text no-bullets-unifi">
            <li data-popover-type="text" data-popover-text="Uptime: {{ printf "%.1f" (div (.Float "health.#(subsystem=\"wan\").gw_system-stats.uptime") 86400) }} Days">
              <p style="display:inline-flex;align-items:center;">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 122.88 99.54" fill="currentColor" stroke="currentColor" stroke-width="0"  class="size-6" style="height:1em;vertical-align:middle;margin-right:0.5em;">
                    <path d="M73.12,0c6.73,0,13.16,1.34,19.03,3.77c6.09,2.52,11.57,6.22,16.16,10.81l0,0c4.58,4.58,8.28,10.06,10.8,16.17 c2.43,5.87,3.77,12.3,3.77,19.02c0,6.73-1.34,13.16-3.77,19.03c-2.52,6.09-6.22,11.57-10.81,16.16l0,0 c-4.58,4.58-10.06,8.28-16.17,10.8c-5.87,2.43-12.3,3.77-19.02,3.77c-6.73,0-13.16-1.34-19.03-3.77 c-6.09-2.52-11.57-6.22-16.15-10.8l-0.01-0.01c-4.59-4.59-8.28-10.07-10.8-16.15c-0.78-1.89-1.45-3.83-2-5.82 c1.04,0.1,2.1,0.15,3.17,0.15c2.03,0,4.01-0.18,5.94-0.53c0.32,0.96,0.67,1.91,1.05,2.84c2.07,5,5.11,9.5,8.89,13.28 c3.78,3.78,8.29,6.82,13.28,8.89c4.81,1.99,10.1,3.1,15.66,3.1s10.84-1.1,15.66-3.1c5-2.07,9.5-5.11,13.28-8.89 c3.78-3.78,6.82-8.29,8.89-13.28c1.99-4.81,3.1-10.1,3.1-15.66s-1.1-10.84-3.1-15.66c-2.07-5-5.11-9.5-8.89-13.28 s-8.29-6.82-13.28-8.89c-4.81-1.99-10.1-3.1-15.66-3.1s-10.84,1.1-15.66,3.1c-0.43,0.18-0.86,0.37-1.28,0.56 c-1.64-2.58-3.62-4.92-5.89-6.95c1.24-0.64,2.51-1.23,3.8-1.77C59.97,1.34,66.39,0,73.12,0L73.12,0L73.12,0z M67.41,26.11 c0-1.22,0.5-2.32,1.29-3.12c0.8-0.8,1.9-1.29,3.12-1.29c1.22,0,2.32,0.49,3.12,1.29c0.8,0.8,1.29,1.9,1.29,3.12v23.22l17.35,10.29 c1.04,0.62,1.74,1.61,2.02,2.7c0.28,1.09,0.15,2.29-0.47,3.33v0.01l0,0c-0.62,1.04-1.61,1.74-2.7,2.02s-2.29,0.15-3.33-0.47h-0.01 l0,0L69.68,55.7c-0.67-0.37-1.22-0.91-1.62-1.55c-0.41-0.67-0.65-1.46-0.65-2.3V26.11L67.41,26.11L67.41,26.11z"/>
                    <path style="fill-rule:evenodd; clip-rule:evenodd;" d="M26.98,2.64c14.9,0,26.98,12.08,26.98,26.98c0,14.9-12.08,26.98-26.98,26.98S0,44.52,0,29.62 C0,14.72,12.08,2.64,26.98,2.64L26.98,2.64L26.98,2.64z M26.98,13.72l14.48,17.9h-8.99v9.52H21.49v-9.52H12.5L26.98,13.72 L26.98,13.72z"/>
                </svg>
                {{ printf "%.1f" (div (.Float "health.#(subsystem=\"wan\").gw_system-stats.uptime") 86400) }}
              </p>
            </li>
            <li data-popover-type="text" data-popover-text="Wired clients: {{ .Int "health.#(subsystem=\"lan\").num_user" }}">
              <p style="display:inline-flex;align-items:center;">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" stroke="currentColor" stroke-width="0"  class="size-6" style="height:1em;vertical-align:middle;margin-right:0.5em;">
                  <path fill="none" d="M0 0h24v24H0z"></path><path d="M7.77 6.76 6.23 5.48.82 12l5.41 6.52 1.54-1.28L3.42 12l4.35-5.24zM7 13h2v-2H7v2zm10-2h-2v2h2v-2zm-6 2h2v-2h-2v2zm6.77-7.52-1.54 1.28L20.58 12l-4.35 5.24 1.54 1.28L23.18 12l-5.41-6.52z"></path>
                </svg>
                {{ .Int "health.#(subsystem=\"lan\").num_user" }}
              </p>
            </li>
            <li data-popover-type="text" data-popover-text="Wireless clients: {{ .Int "health.#(subsystem=\"wlan\").num_user" }}">
              <p style="display:inline-flex;align-items:center;">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" stroke="currentColor" stroke-width="0"  class="size-6" style="height:1em;vertical-align:middle;margin-right:0.5em;">
                  <path d="M12 6c3.537 0 6.837 1.353 9.293 3.809l1.414-1.414C19.874 5.561 16.071 4 12 4c-4.071.001-7.874 1.561-10.707 4.395l1.414 1.414C5.163 7.353 8.463 6 12 6zm5.671 8.307c-3.074-3.074-8.268-3.074-11.342 0l1.414 1.414c2.307-2.307 6.207-2.307 8.514 0l1.414-1.414z"></path><path d="M20.437 11.293c-4.572-4.574-12.301-4.574-16.873 0l1.414 1.414c3.807-3.807 10.238-3.807 14.045 0l1.414-1.414z"></path><circle cx="12" cy="18" r="2"></circle>
                </svg>
                {{ .Int "health.#(subsystem=\"wlan\").num_user"  }}
              </p>
            </li>
          </ul>
        </div>
      </div>
    {{ end }}
```

## Details

![](details-preview.png)

```yaml
- type: custom-api
  title: Unifi
  cache: 1m
  allow-insecure: true
  url: https://${UNIFI_CONSOLE_IP}/proxy/network/api/stat/sites
  headers:
    X-API-KEY: ${UNIFI_API_KEY}
    Accept: application/json
  template: | 
    <style>
      .list-horizontal-text.no-bullets-unifi > *:not(:last-child)::after {
          content: none !important;
      }
      .list-horizontal-text.no-bullets-unifi > *:not(:last-child) {
        margin-right: 1em;
      }
    </style>
    {{ range .JSON.Array "data" }}
      <div style="display:flex; align-items:center; gap:12px;">
        <div style="width:40px; height:40px; flex-shrink:0; display:flex; justify-content:center; align-items:center; overflow:hidden;">
          <img src="https://cdn.jsdelivr.net/gh/selfhst/icons/svg/ubiquiti-unifi-light.svg" width="24" height="24" style="object-fit:contain;">
        </div>
        <div style="flex-grow:1; min-width:0;">
          <a class="size-h4 block text-truncate color-highlight">
            {{ .String "health.#(subsystem=\"wan\").gw_name" }}
            {{ if eq (.String "health.#(subsystem=\"wan\").status") "ok" }}
            <span
              style="width: 8px; height: 8px; border-radius: 50%; background-color: var(--color-positive); display: inline-block; vertical-align: middle;"
              data-popover-type="text"
              data-popover-text="WAN Connected"
            ></span>
            {{ else }}
            <span
              style="width: 8px; height: 8px; border-radius: 50%; background-color: var(--color-negative); display: inline-block; vertical-align: middle;"
              data-popover-type="text"
              data-popover-text="WAN Disconnected"
            ></span>
            {{ end }}
          </a>
          <ul class="list-horizontal-text no-bullets-unifi">
            <li data-popover-type="text" data-popover-text="Uptime: {{ printf "%.1f" (div (.Float "health.#(subsystem=\"wan\").gw_system-stats.uptime") 86400) }} Days">
              <p style="display:inline-flex;align-items:center;">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 122.88 99.54" fill="currentColor" stroke="currentColor" stroke-width="0"  class="size-6" style="height:1em;vertical-align:middle;margin-right:0.5em;">
                    <path d="M73.12,0c6.73,0,13.16,1.34,19.03,3.77c6.09,2.52,11.57,6.22,16.16,10.81l0,0c4.58,4.58,8.28,10.06,10.8,16.17 c2.43,5.87,3.77,12.3,3.77,19.02c0,6.73-1.34,13.16-3.77,19.03c-2.52,6.09-6.22,11.57-10.81,16.16l0,0 c-4.58,4.58-10.06,8.28-16.17,10.8c-5.87,2.43-12.3,3.77-19.02,3.77c-6.73,0-13.16-1.34-19.03-3.77 c-6.09-2.52-11.57-6.22-16.15-10.8l-0.01-0.01c-4.59-4.59-8.28-10.07-10.8-16.15c-0.78-1.89-1.45-3.83-2-5.82 c1.04,0.1,2.1,0.15,3.17,0.15c2.03,0,4.01-0.18,5.94-0.53c0.32,0.96,0.67,1.91,1.05,2.84c2.07,5,5.11,9.5,8.89,13.28 c3.78,3.78,8.29,6.82,13.28,8.89c4.81,1.99,10.1,3.1,15.66,3.1s10.84-1.1,15.66-3.1c5-2.07,9.5-5.11,13.28-8.89 c3.78-3.78,6.82-8.29,8.89-13.28c1.99-4.81,3.1-10.1,3.1-15.66s-1.1-10.84-3.1-15.66c-2.07-5-5.11-9.5-8.89-13.28 s-8.29-6.82-13.28-8.89c-4.81-1.99-10.1-3.1-15.66-3.1s-10.84,1.1-15.66,3.1c-0.43,0.18-0.86,0.37-1.28,0.56 c-1.64-2.58-3.62-4.92-5.89-6.95c1.24-0.64,2.51-1.23,3.8-1.77C59.97,1.34,66.39,0,73.12,0L73.12,0L73.12,0z M67.41,26.11 c0-1.22,0.5-2.32,1.29-3.12c0.8-0.8,1.9-1.29,3.12-1.29c1.22,0,2.32,0.49,3.12,1.29c0.8,0.8,1.29,1.9,1.29,3.12v23.22l17.35,10.29 c1.04,0.62,1.74,1.61,2.02,2.7c0.28,1.09,0.15,2.29-0.47,3.33v0.01l0,0c-0.62,1.04-1.61,1.74-2.7,2.02s-2.29,0.15-3.33-0.47h-0.01 l0,0L69.68,55.7c-0.67-0.37-1.22-0.91-1.62-1.55c-0.41-0.67-0.65-1.46-0.65-2.3V26.11L67.41,26.11L67.41,26.11z"/>
                    <path style="fill-rule:evenodd; clip-rule:evenodd;" d="M26.98,2.64c14.9,0,26.98,12.08,26.98,26.98c0,14.9-12.08,26.98-26.98,26.98S0,44.52,0,29.62 C0,14.72,12.08,2.64,26.98,2.64L26.98,2.64L26.98,2.64z M26.98,13.72l14.48,17.9h-8.99v9.52H21.49v-9.52H12.5L26.98,13.72 L26.98,13.72z"/>
                </svg>
                {{ printf "%.1f" (div (.Float "health.#(subsystem=\"wan\").gw_system-stats.uptime") 86400) }}
              </p>
            </li>
            <li data-popover-type="text" data-popover-text="Wired clients: {{ .Int "health.#(subsystem=\"lan\").num_user" }}">
              <p style="display:inline-flex;align-items:center;">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" stroke="currentColor" stroke-width="0"  class="size-6" style="height:1em;vertical-align:middle;margin-right:0.5em;">
                  <path fill="none" d="M0 0h24v24H0z"></path><path d="M7.77 6.76 6.23 5.48.82 12l5.41 6.52 1.54-1.28L3.42 12l4.35-5.24zM7 13h2v-2H7v2zm10-2h-2v2h2v-2zm-6 2h2v-2h-2v2zm6.77-7.52-1.54 1.28L20.58 12l-4.35 5.24 1.54 1.28L23.18 12l-5.41-6.52z"></path>
                </svg>
                {{ .Int "health.#(subsystem=\"lan\").num_user" }}
              </p>
            </li>
            <li data-popover-type="text" data-popover-text="Wireless clients: {{ .Int "health.#(subsystem=\"wlan\").num_user" }}">
              <p style="display:inline-flex;align-items:center;">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" stroke="currentColor" stroke-width="0"  class="size-6" style="height:1em;vertical-align:middle;margin-right:0.5em;">
                  <path d="M12 6c3.537 0 6.837 1.353 9.293 3.809l1.414-1.414C19.874 5.561 16.071 4 12 4c-4.071.001-7.874 1.561-10.707 4.395l1.414 1.414C5.163 7.353 8.463 6 12 6zm5.671 8.307c-3.074-3.074-8.268-3.074-11.342 0l1.414 1.414c2.307-2.307 6.207-2.307 8.514 0l1.414-1.414z"></path><path d="M20.437 11.293c-4.572-4.574-12.301-4.574-16.873 0l1.414 1.414c3.807-3.807 10.238-3.807 14.045 0l1.414-1.414z"></path><circle cx="12" cy="18" r="2"></circle>
                </svg>
                {{ .Int "health.#(subsystem=\"wlan\").num_user"  }}
              </p>
            </li>
          </ul>
        </div>
      </div>
      <div class="margin-block-2" style="margin-top: 1em">
        <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 10px;">
          <div data-popover-type="html">
            <div data-popover-html>
              <div>
                {{ range .Array "health.#(subsystem=\"wan\").uptime_stats.WAN.monitors" }}
                <div class="flex">
                  <div class="size-h5">{{ .String "target" }}</div>
                  <div class="value-separator"></div>
                  <div class="color-highlight text-very-compact"><span>{{ .Int "latency_average" }}</span> <span class="color-base size-h5">ms</span></div>
                </div>
                {{ end }}
              </div>
            </div>
            <div class="size-h5">Latency</div>
            <div class="size-h3 color-highlight">{{ .Int "health.#(subsystem=\"wan\").uptime_stats.WAN.latency_average" }}<span class="color-base"> ms</span></div>
          </div>
          <div>
            <div class="size-h5">WAN IP</div>
            <div class="size-h3 color-highlight">{{ .String "health.#(subsystem=\"wan\").wan_ip" }}</div>
          </div>
          <div>
            <div class="size-h5">Gateway CPU</div>
            <div class="size-h3 color-highlight">{{ .String "health.#(subsystem=\"wan\").gw_system-stats.cpu" }}<span class="color-base"> %</span></div>
          </div>
          <div>
            <div class="size-h5">Gateway RAM</div>
            <div class="size-h3 color-highlight">{{ .String "health.#(subsystem=\"wan\").gw_system-stats.mem" }}<span class="color-base"> %</span></div>
          </div>
        </div>
      </div>
    {{ end }}
```

## Progress Bars

![](progress-preview.png)

```yaml
- type: custom-api
  title: Unifi
  cache: 1m
  allow-insecure: true
  url: https://${UNIFI_CONSOLE_IP}/proxy/network/api/stat/sites
  headers:
    X-API-KEY: ${UNIFI_API_KEY}
    Accept: application/json
  template: | 
    <style>
      .list-horizontal-text.no-bullets-unifi > *:not(:last-child)::after {
          content: none !important;
      }
      .list-horizontal-text.no-bullets-unifi > *:not(:last-child) {
        margin-right: 1em;
      }
    </style>
    {{ range .JSON.Array "data" }}
      <div style="display:flex; align-items:center; gap:12px;">
        <div style="width:40px; height:40px; flex-shrink:0; display:flex; justify-content:center; align-items:center; overflow:hidden;">
          <img src="https://cdn.jsdelivr.net/gh/selfhst/icons/svg/ubiquiti-unifi-light.svg" width="24" height="24" style="object-fit:contain;">
        </div>
        <div style="flex-grow:1; min-width:0;">
          <a class="size-h4 block text-truncate color-highlight">
            {{ .String "health.#(subsystem=\"wan\").gw_name" }}
            {{ if eq (.String "health.#(subsystem=\"wan\").status") "ok" }}
            <span
              style="width: 8px; height: 8px; border-radius: 50%; background-color: var(--color-positive); display: inline-block; vertical-align: middle;"
              data-popover-type="text"
              data-popover-text="WAN Connected"
            ></span>
            {{ else }}
            <span
              style="width: 8px; height: 8px; border-radius: 50%; background-color: var(--color-negative); display: inline-block; vertical-align: middle;"
              data-popover-type="text"
              data-popover-text="WAN Disconnected"
            ></span>
            {{ end }}
          </a>
          <ul class="list-horizontal-text no-bullets-unifi">
            <li data-popover-type="text" data-popover-text="Uptime: {{ printf "%.1f" (div (.Float "health.#(subsystem=\"wan\").gw_system-stats.uptime") 86400) }} Days">
              <p style="display:inline-flex;align-items:center;">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 122.88 99.54" fill="currentColor" stroke="currentColor" stroke-width="0"  class="size-6" style="height:1em;vertical-align:middle;margin-right:0.5em;">
                    <path d="M73.12,0c6.73,0,13.16,1.34,19.03,3.77c6.09,2.52,11.57,6.22,16.16,10.81l0,0c4.58,4.58,8.28,10.06,10.8,16.17 c2.43,5.87,3.77,12.3,3.77,19.02c0,6.73-1.34,13.16-3.77,19.03c-2.52,6.09-6.22,11.57-10.81,16.16l0,0 c-4.58,4.58-10.06,8.28-16.17,10.8c-5.87,2.43-12.3,3.77-19.02,3.77c-6.73,0-13.16-1.34-19.03-3.77 c-6.09-2.52-11.57-6.22-16.15-10.8l-0.01-0.01c-4.59-4.59-8.28-10.07-10.8-16.15c-0.78-1.89-1.45-3.83-2-5.82 c1.04,0.1,2.1,0.15,3.17,0.15c2.03,0,4.01-0.18,5.94-0.53c0.32,0.96,0.67,1.91,1.05,2.84c2.07,5,5.11,9.5,8.89,13.28 c3.78,3.78,8.29,6.82,13.28,8.89c4.81,1.99,10.1,3.1,15.66,3.1s10.84-1.1,15.66-3.1c5-2.07,9.5-5.11,13.28-8.89 c3.78-3.78,6.82-8.29,8.89-13.28c1.99-4.81,3.1-10.1,3.1-15.66s-1.1-10.84-3.1-15.66c-2.07-5-5.11-9.5-8.89-13.28 s-8.29-6.82-13.28-8.89c-4.81-1.99-10.1-3.1-15.66-3.1s-10.84,1.1-15.66,3.1c-0.43,0.18-0.86,0.37-1.28,0.56 c-1.64-2.58-3.62-4.92-5.89-6.95c1.24-0.64,2.51-1.23,3.8-1.77C59.97,1.34,66.39,0,73.12,0L73.12,0L73.12,0z M67.41,26.11 c0-1.22,0.5-2.32,1.29-3.12c0.8-0.8,1.9-1.29,3.12-1.29c1.22,0,2.32,0.49,3.12,1.29c0.8,0.8,1.29,1.9,1.29,3.12v23.22l17.35,10.29 c1.04,0.62,1.74,1.61,2.02,2.7c0.28,1.09,0.15,2.29-0.47,3.33v0.01l0,0c-0.62,1.04-1.61,1.74-2.7,2.02s-2.29,0.15-3.33-0.47h-0.01 l0,0L69.68,55.7c-0.67-0.37-1.22-0.91-1.62-1.55c-0.41-0.67-0.65-1.46-0.65-2.3V26.11L67.41,26.11L67.41,26.11z"/>
                    <path style="fill-rule:evenodd; clip-rule:evenodd;" d="M26.98,2.64c14.9,0,26.98,12.08,26.98,26.98c0,14.9-12.08,26.98-26.98,26.98S0,44.52,0,29.62 C0,14.72,12.08,2.64,26.98,2.64L26.98,2.64L26.98,2.64z M26.98,13.72l14.48,17.9h-8.99v9.52H21.49v-9.52H12.5L26.98,13.72 L26.98,13.72z"/>
                </svg>
                {{ printf "%.1f" (div (.Float "health.#(subsystem=\"wan\").gw_system-stats.uptime") 86400) }}
              </p>
            </li>
            <li data-popover-type="text" data-popover-text="Wired clients: {{ .Int "health.#(subsystem=\"lan\").num_user" }}">
              <p style="display:inline-flex;align-items:center;">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" stroke="currentColor" stroke-width="0"  class="size-6" style="height:1em;vertical-align:middle;margin-right:0.5em;">
                  <path fill="none" d="M0 0h24v24H0z"></path><path d="M7.77 6.76 6.23 5.48.82 12l5.41 6.52 1.54-1.28L3.42 12l4.35-5.24zM7 13h2v-2H7v2zm10-2h-2v2h2v-2zm-6 2h2v-2h-2v2zm6.77-7.52-1.54 1.28L20.58 12l-4.35 5.24 1.54 1.28L23.18 12l-5.41-6.52z"></path>
                </svg>
                {{ .Int "health.#(subsystem=\"lan\").num_user" }}
              </p>
            </li>
            <li data-popover-type="text" data-popover-text="Wireless clients: {{ .Int "health.#(subsystem=\"wlan\").num_user" }}">
              <p style="display:inline-flex;align-items:center;">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" stroke="currentColor" stroke-width="0"  class="size-6" style="height:1em;vertical-align:middle;margin-right:0.5em;">
                  <path d="M12 6c3.537 0 6.837 1.353 9.293 3.809l1.414-1.414C19.874 5.561 16.071 4 12 4c-4.071.001-7.874 1.561-10.707 4.395l1.414 1.414C5.163 7.353 8.463 6 12 6zm5.671 8.307c-3.074-3.074-8.268-3.074-11.342 0l1.414 1.414c2.307-2.307 6.207-2.307 8.514 0l1.414-1.414z"></path><path d="M20.437 11.293c-4.572-4.574-12.301-4.574-16.873 0l1.414 1.414c3.807-3.807 10.238-3.807 14.045 0l1.414-1.414z"></path><circle cx="12" cy="18" r="2"></circle>
                </svg>
                {{ .Int "health.#(subsystem=\"wlan\").num_user"  }}
              </p>
            </li>
          </ul>
        </div>
      </div>
      <div class="margin-block-2" style="margin-top: 1em">
        <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 10px;">
          <div data-popover-type="html">
            <div data-popover-html>
              <div>
                {{ range .Array "health.#(subsystem=\"wan\").uptime_stats.WAN.monitors" }}
                <div class="flex">
                  <div class="size-h5">{{ .String "target" }}</div>
                  <div class="value-separator"></div>
                  <div class="color-highlight text-very-compact"><span>{{ .Int "latency_average" }}</span> <span class="color-base size-h5">ms</span></div>
                </div>
                {{ end }}
              </div>
            </div>
            <div class="size-h5">Latency</div>
            <div class="size-h3 color-highlight">{{ .Int "health.#(subsystem=\"wan\").uptime_stats.WAN.latency_average" }}<span class="color-base"> ms</span></div>
          </div>
          <div>
            <div class="size-h5">WAN IP</div>
            <div class="size-h3 color-highlight">{{ .String "health.#(subsystem=\"wan\").wan_ip" }}</div>
          </div>
          <div>
            <div class="flex items-end size-h5">
                <div>Gateway CPU</div>
                <div class="color-highlight margin-left-auto text-very-compact">
                  {{ printf "%.0f" (.Float "health.#(subsystem=\"wan\").gw_system-stats.cpu") }}<span class="color-base"> %</span>
                </div>
            </div>
            <div class="progress-bar progress-bar-combined">
                <div class="progress-value{{ if ge (.Float "health.#(subsystem=\"wan\").gw_system-stats.cpu" | toInt) 80 }} progress-value-notice{{ end }}" style="--percent: {{ .Float "health.#(subsystem=\"wan\").gw_system-stats.cpu" | toInt | formatNumber }}"></div>
            </div>
          </div>
          <div class="flex-1">
            <div class="flex items-end size-h5">
                <div>Gateway RAM</div>
                <div class="color-highlight margin-left-auto text-very-compact">
                  {{ printf "%.0f" (.Float "health.#(subsystem=\"wan\").gw_system-stats.mem") }}<span class="color-base">%</span>
                </div>
            </div>
            <div class="progress-bar progress-bar-combined">
                <div class="progress-value{{ if ge (.Float "health.#(subsystem=\"wan\").gw_system-stats.mem" | toInt) 80 }} progress-value-notice{{ end }}" style="--percent: {{ .Float "health.#(subsystem=\"wan\").gw_system-stats.mem" | toInt | formatNumber }}"></div>
            </div>
          </div>
        </div>
      </div>
    {{ end }}
```