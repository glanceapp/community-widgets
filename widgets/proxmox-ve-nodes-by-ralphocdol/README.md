![](preview.png)

```yaml
- type: custom-api
  title: PROXMOX-VE NODES
  title-url: https://${PROXMOXVE_URL}
  cache: 1m
  url: https://${PROXMOXVE_URL}/api2/json/nodes
  allow-insecure: true
  headers:
    Accept: application/json
    Authorization: PVEAPIToken=${PROXMOXVE_KEY}
  options:
    uptime-mode: duration  # relative | duration
    collapse-after: 5
  template: |
    {{ $uptimeMode := .Options.StringOr "uptime-mode" "relative" }}
    {{ $collapseAfter := .Options.IntOr "collapseAfter" 3 }}
    <ul class="list list-gap-14 collapsible-container" data-collapse-after="{{ $collapseAfter }}">
      {{ range .JSON.Array "data" }}
      {{ if ne (.String "type") "node" }} {{ continue }} {{ end }}
      <li>
        <div class="server" proxmox-node="{{ .String "node" }}">
          <div class="server-info">
              <div class="server-details">
                  <div class="server-name color-highlight size-h3 text-truncate">
                    {{ .String "node" }}
                  </div>
                  <div>
                      {{ if eq $.Response.StatusCode 200 }}
                        <div data-popover-type="html">
                          <div data-popover-html>
                            {{ if eq $uptimeMode "relative" }}
                              <span {{ now.Add (duration (concat "-" (.String "uptime") "s")) | toRelativeTime }}></span>
                            {{ else }}
                              <span>{{ duration (concat (.String "uptime") "s") }}</span>
                            {{ end }}
                          </div>
                          <div>
                            {{ if ne $uptimeMode "relative" }}
                              <span {{ now.Add (duration (concat "-" (.String "uptime") "s")) | toRelativeTime }}></span>
                            {{ else }}
                              <span>{{ duration (concat (.String "uptime") "s") }}</span>
                            {{ end }}
                          </div>
                        </div>
                      {{ else }}
                        {{ .String "status" }}
                      {{ end }}
                  </div>
              </div>
              <svg class="server-icon" stroke="var(--color-{{ if eq (.String "status") "online" }}positive{{ else }}negative{{ end }})" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M21.75 17.25v-.228a4.5 4.5 0 0 0-.12-1.03l-2.268-9.64a3.375 3.375 0 0 0-3.285-2.602H7.923a3.375 3.375 0 0 0-3.285 2.602l-2.268 9.64a4.5 4.5 0 0 0-.12 1.03v.228m19.5 0a3 3 0 0 1-3 3H5.25a3 3 0 0 1-3-3m19.5 0a3 3 0 0 0-3-3H5.25a3 3 0 0 0-3 3m16.5 0h.008v.008h-.008v-.008Zm-3 0h.008v.008h-.008v-.008Z" />
              </svg>
          </div>
          <div class="server-stats">
              <div class="flex-1">
                  <div class="flex items-end size-h5">
                      <div>CPU</div>
                      {{ if (ge (mul (.Float "cpu") 100 | toInt ) 80) }}
                      <svg class="server-spicy-cpu-icon" fill="var(--color-negative)" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" >
                          <path fill-rule="evenodd" d="M8.074.945A4.993 4.993 0 0 0 6 5v.032c.004.6.114 1.176.311 1.709.16.428-.204.91-.61.7a5.023 5.023 0 0 1-1.868-1.677c-.202-.304-.648-.363-.848-.058a6 6 0 1 0 8.017-1.901l-.004-.007a4.98 4.98 0 0 1-2.18-2.574c-.116-.31-.477-.472-.744-.28Zm.78 6.178a3.001 3.001 0 1 1-3.473 4.341c-.205-.365.215-.694.62-.59a4.008 4.008 0 0 0 1.873.03c.288-.065.413-.386.321-.666A3.997 3.997 0 0 1 8 8.999c0-.585.126-1.14.351-1.641a.42.42 0 0 1 .503-.235Z" clip-rule="evenodd" />
                      </svg>
                      {{ end }}
                      <div class="color-highlight margin-left-auto text-very-compact">
                        <span>{{ mul (.Float "cpu") 100 | toInt | formatNumber }}</span> <span class="color-base">%</span>
                      </div>
                  </div>
                  <div data-popover-type="html">
                      <div data-popover-html>
                        <div>
                          <div class="flex">
                            <div class="size-h5">LOAD</div>
                            <div class="value-separator"></div>
                            <div class="color-highlight text-very-compact"><span>{{ mul (.Float "cpu") 100 | toInt | formatNumber }}</span> <span class="color-base size-h5">%</span></div>
                          </div>
                          <div class="flex margin-top-3">
                            <div class="size-h5">CORES</div>
                            <div class="value-separator"></div>
                            <div class="color-highlight text-very-compact">{{ .Int "maxcpu" }}</div>
                          </div>
                        </div>
                      </div>
                      <div class="progress-bar progress-bar-combined">
                          <div class="progress-value{{ if ge (mul (.Float "cpu") 100 | toInt) 80 }} progress-value-notice{{ end }}" style="--percent: {{ mul (.Float "cpu") 100 | toInt | formatNumber }}"></div>
                      </div>
                  </div>
              </div>
              <div class="flex-1">
                  <div class="flex justify-between items-end size-h5">
                      <div>RAM</div>
                      <div class="color-highlight text-very-compact">
                        <span>{{ mul (div (.Int "mem" | toFloat) (.Int "maxmem" | toFloat)) 100 | toInt }}</span> <span class="color-base">%</span>
                      </div>
                  </div>
                  <div data-popover-type="html">
                      <div data-popover-html>
                        <div>
                          <div class="flex">
                              <div class="size-h5">RAM</div>
                              <div class="value-separator"></div>
                              <div class="color-highlight text-very-compact">
                                <span>{{ div (.Int "mem" | toFloat) 1073741824 | toInt | formatNumber }}GB</span>
                                <span class="color-base size-h5">/</span>
                                <span>{{ div (.Int "maxmem" | toFloat) 1073741824 | toInt | formatNumber }}GB</span>
                              </div>
                          </div>
                        </div>
                      </div>
                      <div class="progress-bar progress-bar-combined">
                        <div class="progress-value{{ if ge (mul (div (.Int "mem" | toFloat) (.Int "maxmem" | toFloat)) 100 | toInt) 80 }} progress-value-notice{{ end }}" style="--percent: {{ mul (div (.Int "mem" | toFloat) (.Int "maxmem" | toFloat)) 100 | toInt }}"></div>
                      </div>
                  </div>
              </div>
              <div class="flex-1">
                  <div class="flex justify-between items-end size-h5">
                      <div>DISK</div>
                      <div class="color-highlight text-very-compact">
                        <span>{{ mul (div (.Int "disk" | toFloat) (.Int "maxdisk" | toFloat)) 100 | toInt }}</span> <span class="color-base">%</span>
                      </div>
                  </div>
                  <div data-popover-type="html">
                    <div data-popover-html>
                      <div>
                        <div class="flex">
                            <div class="size-h5">DISK</div>
                            <div class="value-separator"></div>
                            <div class="color-highlight text-very-compact">
                                <span>{{ div (.Int "disk" | toFloat) 1073741824 | toInt | formatNumber }}GB</span>
                                <span class="color-base size-h5">/</span>
                                <span>{{ div (.Int "maxdisk" | toFloat) 1073741824 | toInt | formatNumber }}GB</span>
                            </div>
                        </div>
                      </div>
                    </div>
                    <div class="progress-bar progress-bar-combined">
                      <div class="progress-value{{ if ge (mul (div (.Int "disk" | toFloat) (.Int "maxdisk" | toFloat)) 100 | toInt) 80 }} progress-value-notice{{ end }}" style="--percent: {{ mul (div (.Int "disk" | toFloat) (.Int "maxdisk" | toFloat)) 100 | toInt }}"></div>
                    </div>
                  </div>
              </div>
          </div>
        </div>
      </li>
      {{ end }}
    </ul>
```

## Environment variables

#### `PROXMOXVE_URL`
the URL of the Proxmox VE server

#### `PROXMOXVE_KEY`
You need to generate an API token for it, follow the steps below if you don't have one

#### [See proxmox-ve-stats widget](/widgets/proxmox-ve-stats-by-ralphocdol#proxmoxve_key) on how to make one.

-------
## Credits
[titembaataar](https://github.com/titembaatar) - For testing with multiple nodes
