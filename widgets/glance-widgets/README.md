## Glanses widget

Displays data from the glances server

<img src=glances-widgets.png width="700">

### Widget yaml

<details>

<summary>Click</summary>

```yaml
- type: custom-api
    title: System Monitor
    cache: 15m
    url: ${GLANCES_URL}/api/4/system
    options:
    base_url:    ${GLANCES_URL}
    net_iface:   "<you network>"
    disk_mounts: "^(/youdisk1|/youdisk2)$"  
    template: |
    {{ $base     := .Options.StringOr "base_url"  "http://localhost:61208" }}
    {{ $netIface := .Options.StringOr "net_iface" "enp5s0" }}
    
    {{ $cpu    := newRequest (printf "%s/api/4/cpu"     $base) | getResponse }}
    {{ $mem    := newRequest (printf "%s/api/4/mem"     $base) | getResponse }}
    {{ $swap   := newRequest (printf "%s/api/4/memswap" $base) | getResponse }}
    {{ $load   := newRequest (printf "%s/api/4/load"    $base) | getResponse }}
    {{ $fs     := newRequest (printf "%s/api/4/fs"      $base) | getResponse }}
    {{ $net    := newRequest (printf "%s/api/4/network" $base) | getResponse }}
    {{ $diskio := newRequest (printf "%s/api/4/diskio"  $base) | getResponse }}
    {{ $uptime := newRequest (printf "%s/api/4/uptime"  $base) | getResponse }}

    {{/* CPU Calculations */}}
    {{ $cpuPct    := $cpu.JSON.Float "total" }}
    {{ $cores     := $load.JSON.Int  "cpucore" }}
    {{ $load1     := $load.JSON.Float "min1" }}
    {{ $load15    := $load.JSON.Float "min15" }}
    
    {{ $load1pct  := 0 }}
    {{ $load15pct := 0 }}
    {{ if gt $cores 0 }}
        {{ $load1pct = div (mul $load1 100) $cores }}
        {{ $load15pct = div (mul $load15 100) $cores }}
    {{ end }}

    {{ $memPct    := $mem.JSON.Float  "percent" }}
    {{ $memUsedGB := div ($mem.JSON.Int  "used")  1073741824 }}
    {{ $memTotGB  := div ($mem.JSON.Int  "total") 1073741824 }}
    
    {{ $swapPct   := $swap.JSON.Float "percent" }}
    {{ $swapUsedMB := div ($swap.JSON.Int "used")  1048576 }}
    {{ $swapTotGB  := div ($swap.JSON.Int "total") 1073741824 }}

    {{/* Network Calculations */}}
    {{ $netRx := 0 }}
    {{ $netTx := 0 }}
    {{ range $net.JSON.Array "" }}
        {{ if eq (.String "interface_name") $netIface }}
        {{ $netRx = .Int "bytes_recv_rate_per_sec" }}
        {{ $netTx = .Int "bytes_sent_rate_per_sec" }}
        {{ end }}
    {{ end }}
    
    {{ $rxKB := div $netRx 1024 }}
    {{ $txKB := div $netTx 1024 }}
    {{ $rxMB := div $rxKB 1024 }}
    {{ $txMB := div $txKB 1024 }}

        <style>
        /* Вертикальные бары */
        .v-bars-row {
        display: flex;
        justify-content: space-around;
        align-items: flex-end;
        height: 50px;
        gap: 4px;
        }
        .v-bar-col {
        display: flex;
        flex-direction: column;
        align-items: center;
        width: 12%;
        height: 100%;
        justify-content: flex-end;
        }
        .v-bar-track {
        width: 2rem;
        height: 100%;
        background-color: var(--color-bg-alt);
        border-radius: 3px;
        position: relative;
        overflow: hidden;
        }
        .v-bar-fill {
        position: absolute;
        bottom: 0;
        left: 0;
        width: 100%;
        background-color: var(--color-progress-value);
        border-radius: 0 var(--half-border-radius) var(--half-border-radius) 0;
        }
        .v-bar-label {
        font-size: 9px;
        margin-top: 4px;
        color: var(--color-text-dim);
        text-align: center;
        }


        .stat-item {
        display: flex;
        flex-direction: row;
        align-items: flex-end;
        gap: 5px;
        }
        .stat-label {
        display: flex;
        flex-direction: column;
        height: 100%;
        justify-content: space-between;
        }

        .v-bar-track{
        border: 1px solid var(--color-progress-border);
        }
        </style>
    
        {{/* Stats Row 1: CPU, RAM, NET */}}

        <div class="server-stats" style="display: flex; flex-wrap: wrap; gap: 8px; justify-content: flex-start; align-items: flex-start; ">
            
            {{/* CPU */}}
            <div style="display:flex; align-items: flex-start; flex-direction: row;">
            <div class="size-h5" style="writing-mode: vertical-lr; cursor:default;">CPU</div>
            <div style="display:flex; align-items: flex-end; gap:3px; flex-direction: column;">
                <div class="v-bars-row">
                <div class="v-bar-col">
                    <div class="v-bar-track">
                    <div class="v-bar-fill" style="height:{{ printf "%.0f" $cpuPct }}%"></div>
                    </div>
                </div>
                </div>
                <div class="size-h5 color-highlight" style="flex-shrink:0; text-align:center;">{{ printf "%.0f" $cpuPct }}%</div>
            </div>
            </div>

            {{/* RAM */}}
            <div style="display:flex; align-items: flex-start; flex-direction: row;">
            <div class="size-h5" style="writing-mode: vertical-lr; cursor:default;">RAM</div>
            <div style="display:flex; align-items: flex-end; gap:3px; flex-direction: column;">
                <div class="v-bars-row">
                <div class="v-bar-col">
                    <div class="v-bar-track">
                    <div class="v-bar-fill" style="height:{{ printf "%.0f" $memPct }}%"></div>
                    </div>
                </div>
                </div>
                <div class="size-h5 color-highlight" style="flex-shrink:0; text-align:center;">{{ printf "%.0f" $memPct }}%</div>
            </div>
            </div>

            {{/* NETWORK */}}
            {{ $rxPct := div $netRx 1048576 }}{{ if gt $rxPct 100 }}{{ $rxPct = 100 }}{{ end }}
            {{ $txPct := div $netTx 1048576 }}{{ if gt $txPct 100 }}{{ $txPct = 100 }}{{ end }}
            <div style="display:flex; align-items: flex-start; flex-direction: row;">
            <div class="size-h5" style="writing-mode: vertical-lr; cursor:default;">NETWORK</div>
            <div style="display:flex; align-items:flex-end; gap:3px; flex-direction:column;">
            <div style="display:flex; gap:4px;">

                {{/* ↓ RX */}}
                <div style="display:flex; flex-direction:column; align-items:center;">
                <div class="v-bar-track" style="width:100%; height:50px;">
                    <div class="v-bar-fill" style="height:{{ $rxPct }}%"></div>
                </div>
                <div class="size-h5 color-highlight" style="white-space:nowrap;">
                    ↓{{ if gt $rxMB 0 }}{{ $rxMB }}MB{{ else if gt $rxKB 0 }}{{ $rxKB }}KB{{ else }}{{ $netRx }}B{{ end }}
                </div>
                </div>

                {{/* ↑ TX */}}
                <div style="display:flex; flex-direction:column; align-items:center;">
                <div class="v-bar-track" style="width:100%; height:50px;">
                    <div class="v-bar-fill" style="height:{{ $txPct }}%"></div>
                </div>
                <div class="size-h5 color-highlight" style="white-space:nowrap;">
                    ↑{{ if gt $txMB 0 }}{{ $txMB }}MB{{ else if gt $txKB 0 }}{{ $txKB }}KB{{ else }}{{ $netTx }}B{{ end }}
                </div>
                </div>
            </div>
            </div>
        </div>
    

        {{/* ───── Диски ───── */}}
        {{ $diskPattern := .Options.StringOr "disk_mounts" "^/rootfs$" }}

        {{ range $fs.JSON.Array "" }}
        {{ $mnt := .String "mnt_point" }}
        {{ if findMatch $diskPattern $mnt }}

        {{ $pct     := .Float "percent" }}
        {{ $usedGB  := div (.Int "used")  1073741824 }}
        {{ $sizeGB  := div (.Int "size")  1073741824 }}
        {{ $devName := .String "device_name" }}

        {{ $readBytes  := 0 }}
        {{ $writeBytes := 0 }}
        {{ range $diskio.JSON.Array "" }}
        {{ $dname := .String "disk_name" }}
        {{ if findMatch (printf "/dev/%s" $dname) $devName }}
            {{ $readBytes  = .Int "read_bytes_rate_per_sec" }}
            {{ $writeBytes = .Int "write_bytes_rate_per_sec" }}
        {{ end }}
        {{ end }}

        {{ $readMB  := div $readBytes  1048576 }}
        {{ $writeMB := div $writeBytes 1048576 }}
        {{ $readKB  := div $readBytes  1024 }}
        {{ $writeKB := div $writeBytes 1024 }}

        <div data-popover-type="html" data-popover-max-width="180px" style="display:flex; align-items: flex-start; flex-direction: row;">

        {{/* Название */}}
        <div class="size-h5" style="writing-mode: vertical-lr; cursor:default;">{{ $mnt }}</div>

        <div style="display:flex; align-items: flex-end; gap:3px; flex-direction: column;">
            <div data-popover-html="">
            <div class="size-h5 color-highlight">{{ $mnt }}</div>
            <div class="flex margin-top-3">
                <div class="size-h5">Used</div>
                <div class="value-separator"></div>
                <div class="color-highlight">{{ $usedGB }} / {{ $sizeGB }} GB</div>
            </div>
            <div class="flex margin-top-3">
                <div class="size-h5">Read</div>
                <div class="value-separator"></div>
                <div class="color-highlight">{{ if gt $readMB 0 }}{{ $readMB }}MB/s{{ else if gt $readKB 0 }}{{ $readKB }}KB/s{{ else }}—{{ end }}</div>
            </div>
            <div class="flex margin-top-3">
                <div class="size-h5">Write</div>
                <div class="value-separator"></div>
                <div class="color-highlight">{{ if gt $writeMB 0 }}{{ $writeMB }}MB/s{{ else if gt $writeKB 0 }}{{ $writeKB }}KB/s{{ else }}—{{ end }}</div>
            </div>
            </div>


            <div class="v-bars-row">
            <div class="v-bar-col">
                <div class="v-bar-track">
                <div class="v-bar-fill" style="height:{{ printf "%.0f" $pct }}%"></div>
                </div>
            </div>
            </div>


            {{/* Процент */}}
            <div class="size-h5 color-highlight" style="flex-shrink:0; text-align: center;">{{ printf "%.0f" $pct }}%</div>
        </div>
        </div>

    {{ end }}
    {{ end }}

    </div>
```
</details>
