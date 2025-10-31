## Preview

![](preview.png)

## Configuration

```yaml
- type: custom-api
  title: Proxmox-VE Resources
  cache: 1m
  url: ${PROXMOX_URL}/api2/json/cluster/resources
  allow-insecure: true
  headers:
    Accept: application/json
    Authorization: ${PROXMOX_TOKEN}
  options:
    # Enable or disable charts for RAM, CPU, and Disk usage
    cpu_graph: true 
    disk_graph: true
    ram_graph: true
  template: |
    <style>
      .proxmox-container { font-size: 0.9em; }
      .resource-header, .resource-row { display: flex; border-collapse: collapse; align-items: center; gap: 10px; padding: 4px 0;}
      .resource-row { text-transform: lowercase;}
      .resource-header { text-transform: uppercase; font-size: 1em; font-weight: bold; color: var(--color-primary); border-bottom: 1px solid var(--color-primary); margin-bottom: 1px;}
      .col { min-width: 0; margin-right: 15px; }
      .col-item { flex: 1.5; color: var(--color-primary);}
      .col-type { flex: 1; }
      .col-ram { flex: 2; }
      .col-cpu { flex: 2; }
      .col-disk { flex: 2; }
      .col-status { flex: 0.5; text-align: center; }
      .server-type { padding: 2px 6px; border-radius: 4px; text-transform: lowercase; }
      .stat-header { display: flex; justify-content: space-between; align-items: flex-end; font-size: 0.85em;}
      .stat-percent { font-weight: 500;}
      .status-indicator { width: 8px; height: 8px; border-radius: 50%; display: inline-block; margin-left: 4px;}
      .indicator-positive { background-color: var(--color-positive); }
      .indicator-negative { background-color: var(--color-negative); }
    </style>
    
    {{ if eq .Response.StatusCode 200 }}
      <div class="proxmox-container">
        {{- $filteredResources := .JSON.Array "data.#(type!=\"node\" and type!=\"sdn\")#" -}}
        {{- $sortedResources := sortByString "type" "asc" $filteredResources -}}
        {{- $bytesInGB := toFloat 1073741824 -}}

        <div class="resource-header">
          <div class="col col-item">ITEM</div>
          <div class="col col-type">TYPE</div>
          {{ if .Options.BoolOr "ram_graph" true }}<div class="col col-ram">RAM</div>{{ end }}
          {{ if .Options.BoolOr "cpu_graph" true }}<div class="col col-cpu">CPU</div>{{ end }}
          {{ if .Options.BoolOr "disk_graph" true }}<div class="col col-disk">DISK</div>{{ end }}
          <div class="col col-status">STATUS</div>
        </div>

        {{- range $sortedResources -}}
          {{- $type := .String "type" -}}
            {{ if and (ne $type "node") (ne $type "sdn") }}
              <div class="resource-row">
                
                <div class="col col-item" title="{{ if .Exists "name" }}{{ .String "name" }}{{ else }}{{ .String "storage" }}{{ end }}">
                  {{ if .Exists "name" }}{{ .String "name" }}{{ else }}{{ .String "storage" }}{{ end }}
                </div>
              
                <div class="col col-type">
                  <span class="server-type">{{ $type }}</span>
                </div>

                {{ if $.Options.BoolOr "ram_graph" true }}
                  <div class="col col-ram">
                    {{- if and (or (eq $type "qemu") (eq $type "lxc")) (gt (.Float "maxmem") 0.0) -}}
                      <div data-popover-type="html">
                        <div data-popover-html>Usage: {{ printf "%.1f" (div (.Float "mem") $bytesInGB) }}GB<br>Total: {{ printf "%.1f" (div (.Float "maxmem") $bytesInGB) }}GB</div>
                        <div class="stat-header"><span class="stat-percent">{{ mul (div (.Float "mem") (.Float "maxmem")) 100 | toInt }}%</span></div>
                        <div class="progress-bar"><div class="progress-value{{ if ge (mul (div (.Float "mem") (.Float "maxmem")) 100 | toInt) 80 }} progress-value-notice{{ end }}" style="--percent: {{ mul (div (.Float "mem") (.Float "maxmem")) 100 | toInt }}"></div></div>
                      </div>
                    {{- end -}}
                  </div>
                {{ end }}

                {{ if $.Options.BoolOr "cpu_graph" true }}
                  <div class="col col-cpu">
                    {{- if and (or (eq $type "qemu") (eq $type "lxc")) -}}
                      {{- $percentCPU := mul (.Float "cpu") 100 | toInt -}}
                      <div data-popover-type="html">
                        <div data-popover-html>Load: {{ $percentCPU }}%<br>Cores: {{ .Int "maxcpu" }}</div>
                        <div class="stat-header"><span class="stat-percent">{{ $percentCPU }}%</span></div>
                        <div class="progress-bar"><div class="progress-value{{ if ge $percentCPU 80 }} progress-value-notice{{ end }}" style="--percent: {{ $percentCPU }}"></div></div>
                      </div>
                    {{- end -}}
                  </div>
                {{ end }}
                
                {{ if $.Options.BoolOr "disk_graph" true }}
                  <div class="col col-disk">
                    {{- if and (or (eq $type "lxc") (eq $type "storage")) (gt (.Float "maxdisk") 0.0) -}}
                      <div data-popover-type="html">
                        <div data-popover-html>Usage: {{ printf "%.1f" (div (.Float "disk") $bytesInGB) }}GB<br>Total: {{ printf "%.1f" (div (.Float "maxdisk") $bytesInGB) }}GB</div>
                        <div class="stat-header"><span class="stat-percent">{{ mul (div (.Float "disk") (.Float "maxdisk")) 100 | toInt }}%</span></div>
                        <div class="progress-bar"><div class="progress-value{{ if ge (mul (div (.Float "disk") (.Float "maxdisk")) 100 | toInt) 80 }} progress-value-notice{{ end }}" style="--percent: {{ mul (div (.Float "disk") (.Float "maxdisk")) 100 | toInt }}"></div></div>
                      </div>
                    {{- end -}}
                  </div>
                {{ end }}

                <div class="col col-status">
                  {{- $status := .String "status" -}}
                  <span class="status-indicator {{ if or (eq $status "running") (eq $status "available") (eq $status "ok") (eq $status "online")}}indicator-positive{{ else }}indicator-negative{{ end }}"
                        data-popover-type="text" data-popover-text="{{ $status }}">
                  </span>
                </div>
              </div>
            {{ end }}
        {{- end -}}
      </div>
    {{ else }}
      <p>Failed: {{ .Response.Status }}</p>
    {{ end }}
```
## Environment variables

- `PROXMOX_URL` - The URL for your proxmox host, including port but without trailing slash, e.g.: 'https://192.168.1.2:8006', or 'https://proxmox.example.com'
- `PROXMOX_TOKEN` - read only token in format `PVEAPIToken=<username>@pam!<tokenID>=<secret>`. (instructions to set it up can be found in [gethomepage's documentation](https://github.com/gethomepage/homepage/blob/main/docs/widgets/services/proxmox.md)

    token example:
    ```
    PVEAPIToken=api@pam!shared=some-random-secret-value-here
    ```
