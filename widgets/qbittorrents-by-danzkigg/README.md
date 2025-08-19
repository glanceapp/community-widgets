## qBittorrent widget

At-a-***glance*** view your total download speed, and the number of torrents currently seeding and leeching.

There are 2 options for this widget:
1. ```detailed```: with a detailed list for each active torrent currently downloading that shows under the "Show more" button (hidden when there is no active torrents)
2. ```basic```: just stats from qBittorrent (Current download speed, number of torrents seeding and leeching)

### Prerequisites (Required)

This widget requires you to bypass authentication for your Glances server within qBittorrent's settings. This allows Glance to securely access the API without needing to handle a complex login process or wrapper.

1. Open the qBittorrent Web UI.

2. Go to Tools -> Options -> Web UI.

3. Under the Authentication section, check the box for "Bypass authentication for clients in whitelisted IP subnets".

4. In the text box, add the IP address or subnet of where **Glance** is running. For example, if your entire local network is 192.168.1.x, you can safely add 192.168.1.0/24.

5. Click Save.

### Setup
The widget uses an environment variable to know where to find your qBittorrent server.

You must define QBITTORRENT_URL in your environment. **Important: Do not include http:// in the value.**

- Variable: ```QBITTORRENT_URL```

- Example Value: ```192.168.1.8:8080```

If you don't want to use environment variables, you can just hard code your ip adress directly in the widget by replacing the ```${QBITTORRENT_URL}``` with your qBittorrent IP.

### Detailed Widget

![Detailed Widget Closed](./preview1.png) 
![Detailed Widget Closed](./preview2.png) 
![Detailed Widget no Active Downloads](./preview3.png) 

Use ```type: detailed``` in ```options```to get the detailed view

![Simple Widget](./preview4.png) 

Use ```type: basic``` in ```options```to get the basic view

Once the prerequisites and setup are complete, copy the code below and add it to your ```glance.yml``` file.

```yaml
- type: custom-api
  title: qBittorrent
  cache: 10s
  options:
    type: "detailed" # "basic" for basic view | "detailed" for detailed view
  subrequests:
    transfer:
      url: http://${QBITTORRENT_URL}/api/v2/transfer/info
    seeding:
      url: http://${QBITTORRENT_URL}/api/v2/torrents/info
      parameters:
        filter: seeding
    leeching:
      url: http://${QBITTORRENT_URL}/api/v2/torrents/info
      parameters:
        filter: downloading
  template: |
    {{ $transfer := .Subrequest "transfer" }}
    {{ $seeding := .Subrequest "seeding" }}
    {{ $leeching := .Subrequest "leeching" }}

    {{ if and (eq $transfer.Response.StatusCode 200) (eq $seeding.Response.StatusCode 200) (eq $leeching.Response.StatusCode 200) }}

      {{ $isDetailed := eq (.Options.StringOr "type" "detailed") "detailed" }}

      {{ if $isDetailed }}
      <!-- Detailed View -->
        <div class="list" style="--list-gap: 15px;">
        
          <div class="flex justify-between text-center">
            <div>
              {{ $speed := $transfer.JSON.Float "dl_info_speed" }}
              {{ if lt $speed 1048576.0 }}
                <div class="color-highlight size-h3">{{ printf "%.0f KiB/s" (div $speed 1024.0) }}</div>
              {{ else }}
                <div class="color-highlight size-h3">{{ printf "%.1f MiB/s" (div $speed 1048576.0) }}</div>
              {{ end }}
              <div class="size-h6">DOWNLOADING</div>
            </div>
            <div>
              <div class="color-highlight size-h3">{{ len ($seeding.JSON.Array "") }}</div>
              <div class="size-h6">SEEDING</div>
            </div>
            <div>
              <div class="color-highlight size-h3">{{ len ($leeching.JSON.Array "") }}</div>
              <div class="size-h6">LEECHING</div>
            </div>
          </div>

          {{ $downloadingTorrents := $leeching.JSON.Array "" }}
          {{ if gt (len $downloadingTorrents) 0 }}
            <ul class="list collapsible-container" data-collapse-after="0" style="--list-gap: 15px; margin-top: 15px;">
              {{ range $t := $downloadingTorrents }}
                <li class="flex items-center" style="gap: 10px;">
                  {{ $state := $t.String "state" }}
                  {{ $icon := "❔" }}
                  {{ if ge ($t.Int "completed") ($t.Int "size") }}{{ $icon = "✔" }}
                  {{ else if eq $state "downloading" "forcedDL" }}{{ $icon = "↓" }}
                  {{ else if eq $state "pausedDL" "stoppedDL" "pausedUP" "stalledDL" "stalledUP" "queuedDL" "queuedUP" }}{{ $icon = "❚❚" }}
                  {{ else if eq $state "error" "missingFiles" }}{{ $icon = "!" }}
                  {{ else if eq $state "checkingDL" "checkingUP" "allocating" }}{{ $icon = "…" }}
                  {{ else if eq $state "checkingResumeData" }}{{ $icon = "⟳" }}
                  {{ end }}
                  <div class="size-h4" style="flex-shrink: 0;">{{ $icon }}</div>
                            
                  <div style="flex-grow: 1; min-width: 0;">
                    <div class="text-truncate color-highlight">{{ $t.String "name" }}</div>
                    <div title="{{ $t.Float "progress" | mul 100 | printf "%.1f" }}%" style="background: rgba(128, 128, 128, 0.2); border-radius: 5px; height: 6px; margin-top: 5px; overflow: hidden;">
                      <div style="width: {{ $t.Float "progress" | mul 100 }}%; background-color: var(--color-positive); height: 100%; border-radius: 5px;"></div>
                    </div>
                  </div>

                  <div style="flex-shrink: 0; text-align: right; width: 80px;">
                    {{ $dlSpeed := $t.Float "dlspeed" }}
                    <div class="size-sm color-paragraph">
                      {{ if lt $dlSpeed 1024.0 }}--
                      {{ else if lt $dlSpeed 1048576.0 }}{{ printf "%.0f KiB/s" (div $dlSpeed 1024.0) }}
                      {{ else }}{{ printf "%.1f MiB/s" (div $dlSpeed 1048576.0) }}
                      {{ end }}
                    </div>
                    {{ $eta := .Int "eta" }}
                    <div class="size-sm color-paragraph">
                    {{ if eq $eta 8640000 }}∞
                    {{ else if gt $eta 3600 }}{{ printf "%dh %dm" (div $eta 3600) (mod (div $eta 60) 60) }}
                    {{ else if gt $eta 0 }}{{ printf "%dm" (div $eta 60) }}
                    {{ else }}--
                    {{ end }}
                    </div>
                  </div>
                </li>
              {{ end }}
            </ul>
          {{ end }}
        </div>

      {{ else }}

      <!-- Basic View -->
        <div class="flex justify-between text-center">
          <div>
            {{ $speed := $transfer.JSON.Float "dl_info_speed" }}
            {{ if lt $speed 1048576.0 }}
              <div class="color-highlight size-h3">{{ printf "%.0f KiB/s" (div $speed 1024.0) }}</div>
            {{ else }}
              <div class="color-highlight size-h3">{{ printf "%.1f MiB/s" (div $speed 1048576.0) }}</div>
            {{ end }}
            <div class="size-h6">DOWNLOADING</div>
          </div>
          <div>
            <div class="color-highlight size-h3">{{ len ($seeding.JSON.Array "") }}</div>
            <div class="size-h6">SEEDING</div>
          </div>
          <div>
            <div class="color-highlight size-h3">{{ len ($leeching.JSON.Array "") }}</div>
            <div class="size-h6">LEECHING</div>
          </div>
        </div>
      {{ end }}

    {{ else }}
      <div class="color-negative text-center">
        <p>Error fetching qBittorrent data.</p>
        <p class="size-sm">Check URL and authentication bypass settings.</p>
      </div>
    {{ end }}
```