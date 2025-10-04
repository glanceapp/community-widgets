## Komodo Container Manager - Server Monitor
Monitor the status and resource use of your Komodo servers

<img width="315" height="726" alt="komodo-servers-preview" src="https://github.com/user-attachments/assets/5566c276-4133-4f03-88ff-2041e3e2e728" />

### Widget Yaml

<details>

  <summary>View YAML</summary>
  
```yaml
- type: custom-api
  title: Komodo Servers
  method: POST
  cache: 5m
  allow-insecure: false # Change to true if your Komodo core uses a self-signed certificate
  url: ${KOMODO_URL}/read
  body-type: json
  body:
    type: ListServers
    params: {}
  headers:
    Content-Type: application/json
    X-Api-Key: ${KOMODO_API_KEY}
    X-Api-Secret: ${KOMODO_API_SECRET}
  options:
    base-url: ${KOMODO_URL}
    api-key: ${KOMODO_API_KEY}
    api-secret: ${KOMODO_API_SECRET}
  template: |
    {{ $urlBase := .Options.StringOr "base-url" "" }}
    {{ $apiKey := .Options.StringOr "api-key" "" }}
    {{ $apiSecret := .Options.StringOr "api-secret" "" }}
    <style>
      .list-horizontal-text.no-bullets-komodo-servers > *:not(:last-child)::after {
          content: none !important;
      }
      .list-horizontal-text.no-bullets-komodo-servers > *:not(:last-child) {
        margin-right: 1em;
      }
    </style>
    
    {{ $servers := .JSON.Array "" }}
    {{ $total := len $servers }}
    {{ $Ok := 0 }}
    {{ $NotOk := 0 }}
    {{ $Disabled := 0 }}
    {{ $unknown := 0 }}
    {{ range $servers }}
      {{ $state := .String "info.state" }}
      {{ if eq $state "Ok" }}{{ $Ok = add $Ok 1 }}{{ end }}
      {{ if eq $state "NotOk" }}{{ $NotOk = add $NotOk 1 }}{{ end }}
      {{ if eq $state "Disabled" }}{{ $Disabled = add $Disabled 1 }}{{ end }}
      {{ if or (eq $state "") (eq $state "unknown") }}{{ $unknown = add $unknown 1 }}{{ end }}
    {{ end }}

    <div style="display:flex; align-items:center; gap:12px; margin-bottom: 1rem;">
      <div style="display:flex; align-items:center; gap:12px;">
        <div style="width:30px; height:30px; flex-shrink:0; display:flex; justify-content:center; align-items:center; overflow:hidden;">
          <img src="https://cdn.jsdelivr.net/gh/homarr-labs/dashboard-icons/svg/komodo.svg" height="30" style="object-fit:contain;">
        </div>
        <div style="flex-grow:1; min-width:0;">
          <h4 class="size-h4 block text-truncate color-highlight">Komodo Servers</h4>
          <ul class="list list-horizontal-text no-bullets-komodo-servers items-center">
            <li data-popover-type="text" data-popover-text="{{ $Ok }}/{{ $total }} Ok">
              <p style="display:inline-flex;align-items:center;">
                <svg xmlns="http://www.w3.org/2000/svg" 
                viewBox="0 0 256 256" 
                style="height:2rem; margin-right:0.5rem; stroke:var(--color-positive); fill:var(--color-positive);">
                  <path d="M224,64H32A16,16,0,0,0,16,80v96a16,16,0,0,0,16,16H224a16,16,0,0,0,16-16V80A16,16,0,0,0,224,64Zm0,112H32V80H224v96Zm-24-48a12,12,0,1,1-12-12A12,12,0,0,1,200,128Z" />
                </svg>
                {{ $Ok }}
              </p>
            </li>
            <li data-popover-type="text" data-popover-text="{{ $Disabled }}/{{ $total }} disabled">
              <p style="display:inline-flex;align-items:center;">
                <svg xmlns="http://www.w3.org/2000/svg" 
                viewBox="0 0 256 256" 
                style="height:2rem; margin-right:0.5rem; stroke:var(--color-progress-value); fill:var(--color-progress-value);">
                  <path d="M224,64H32A16,16,0,0,0,16,80v96a16,16,0,0,0,16,16H224a16,16,0,0,0,16-16V80A16,16,0,0,0,224,64Zm0,112H32V80H224v96Zm-24-48a12,12,0,1,1-12-12A12,12,0,0,1,200,128Z" />
                </svg>
                {{ $Disabled }}
              </p>
            </li>
            <li data-popover-type="text" data-popover-text="{{ $NotOk }}/{{ $total }} Not OK; {{ $unknown }}/{{ $total }} unknown" >
              <p style="display:inline-flex;align-items:center;">
                <svg xmlns="http://www.w3.org/2000/svg" 
                viewBox="0 0 256 256" 
                style="height:2rem; margin-right:0.5rem; stroke:var(--color-negative); fill:var(--color-negative);">
                  <path d="M224,64H32A16,16,0,0,0,16,80v96a16,16,0,0,0,16,16H224a16,16,0,0,0,16-16V80A16,16,0,0,0,224,64Zm0,112H32V80H224v96Zm-24-48a12,12,0,1,1-12-12A12,12,0,0,1,200,128Z" />
                </svg>
                {{ add ( $NotOk ) ( $unknown ) }}
              </p>
            </li>
          </ul>
        </div>
      </div>
    </div>
    <ul class="list dynamic-columns list-gap-15 list-with-separator">
      {{ range .JSON.Array "" }}
        {{ $theID := .String "id" }}
        <li style="display:flex; align-items:center; gap:12px; margin-top: .5rem;">
          <div style="flex-grow:1; min-width:0;">
            <a class="size-h4 block text-truncate color-primary" href={{ concat $urlBase "servers/" ( .String "id" ) }} target="_blank">
              {{
                $systemStats := newRequest (concat $urlBase "/read")
                  | withHeader "Content-Type" "application/json"
                  | withHeader "X-Api-Key" $apiKey
                  | withHeader "X-Api-Secret" $apiSecret
                  | withStringBody (printf `{"type":"GetSystemStats","params":{"server":"%s"}}` $theID)
                  | getResponse
              }}
              
              {{ $state := .String "info.state" }}
              
              <span
                role="status"
                aria-label="Server Status: {{ $state }}"
                style="
                  width: 8px; 
                  height: 8px; 
                  border-radius: 50%; 
                  background-color: 
                    {{ if eq $state "Ok" }} var(--color-positive); 
                    {{ else if eq $state "Disabled" }} var(--color-text-subdue); 
                    {{ else }} var(--color-negative); 
                    {{ end }} 
                  display: inline-block; 
                  vertical-align: middle;"
                data-popover-type="text"
                data-popover-text="Server Status: {{ $state }}"
              ></span>
              {{ .String "name"}}
            </a>
            <span class="size-h6 color-subdue">{{ .String "info.address" }}</span>
            {{ if and $systemStats.Response (ge $systemStats.Response.StatusCode 200) (lt $systemStats.Response.StatusCode 300) }}
              {{ $stats := $systemStats.JSON }}
              <ul class="list-horizontal-text no-bullets-komodo-servers">
                <li data-popover-type="text" data-popover-text="CPU Load: {{ printf "%.2f" ( $stats.Float "cpu_perc" ) }}%">
                  <p style="display:inline-flex;align-items:center;">
                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 256 256" fill="currentColor" style="height:1.5rem;vertical-align:middle;margin-right:0.5rem;">
                      <path d="M152,96H104a8,8,0,0,0-8,8v48a8,8,0,0,0,8,8h48a8,8,0,0,0,8-8V104A8,8,0,0,0,152,96Zm-8,48H112V112h32Zm88,0H216V112h16a8,8,0,0,0,0-16H216V56a16,16,0,0,0-16-16H160V24a8,8,0,0,0-16,0V40H112V24a8,8,0,0,0-16,0V40H56A16,16,0,0,0,40,56V96H24a8,8,0,0,0,0,16H40v32H24a8,8,0,0,0,0,16H40v40a16,16,0,0,0,16,16H96v16a8,8,0,0,0,16,0V216h32v16a8,8,0,0,0,16,0V216h40a16,16,0,0,0,16-16V160h16a8,8,0,0,0,0-16Zm-32,56H56V56H200v95.87s0,.09,0,.13,0,.09,0,.13V200Z" />
                    </svg>
                    {{ printf "%.2f" ( $stats.Float "cpu_perc" ) }}<span class="color-subdue">&nbsp;%</span>
                  </p>
                </li>
                <li data-popover-type="text" data-popover-text="Memory: {{ printf "%.2f"  ( $stats.Float "mem_used_gb" ) }}/{{ printf "%.2f" ( $stats.Float "mem_total_gb" ) }}GB">
                  <p style="display:inline-flex;align-items:center;">
                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 256 256" fill="currentColor" style="height:1.5rem;vertical-align:middle;margin-right:0.5rem;">
                      <path d="M232,56H24A16,16,0,0,0,8,72V200a8,8,0,0,0,16,0V184H40v16a8,8,0,0,0,16,0V184H72v16a8,8,0,0,0,16,0V184h16v16a8,8,0,0,0,16,0V184h16v16a8,8,0,0,0,16,0V184h16v16a8,8,0,0,0,16,0V184h16v16a8,8,0,0,0,16,0V184h16v16a8,8,0,0,0,16,0V72A16,16,0,0,0,232,56ZM24,72H232v96H24Zm88,80a8,8,0,0,0,8-8V96a8,8,0,0,0-8-8H48a8,8,0,0,0-8,8v48a8,8,0,0,0,8,8ZM56,104h48v32H56Zm88,48h64a8,8,0,0,0,8-8V96a8,8,0,0,0-8-8H144a8,8,0,0,0-8,8v48A8,8,0,0,0,144,152Zm8-48h48v32H152Z" />
                    </svg>
                    {{ printf "%.2f" (mul ( div ( $stats.Float "mem_used_gb" ) ( $stats.Float "mem_total_gb" ) ) 100 ) }} <span class="color-subdue">&nbsp;%</span>
                  </p>
                </li>
                {{ range $stats.Array "disks"}}
                  <li data-popover-type="text" data-popover-text="Mount: {{ .String "mount" }}; Format: {{ .String "file_system" }}; Used: {{ printf "%.2f" ( .Float "used_gb" ) }}/{{ printf "%.2f" ( .Float "total_gb" ) }}GB">
                    <p style="display:inline-flex;align-items:center;">
                      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 256 256" fill="currentColor" style="height:2rem;vertical-align:middle;margin-right:0.5rem;">
                        <path d="M208,136H48a16,16,0,0,0-16,16v48a16,16,0,0,0,16,16H208a16,16,0,0,0,16-16V152A16,16,0,0,0,208,136Zm0,64H48V152H208v48Zm0-160H48A16,16,0,0,0,32,56v48a16,16,0,0,0,16,16H208a16,16,0,0,0,16-16V56A16,16,0,0,0,208,40Zm0,64H48V56H208v48ZM192,80a12,12,0,1,1-12-12A12,12,0,0,1,192,80Zm0,96a12,12,0,1,1-12-12A12,12,0,0,1,192,176Z" />
                      </svg>
                      {{ printf "%.2f" (mul ( div ( .Float "used_gb" ) ( .Float "total_gb" ) ) 100 ) }} <span class="color-subdue">&nbsp;%</span>
                    </p>
                  </li>
                {{ end }}
              </ul>
            {{ else }}
              <span class="color-negative size-h6">Stats unavailable</span>
            {{ end }}
          </div>
        </li>
      {{ end }}
    </ul>
```
</details>

### Environment Variables
There are three environment variables that need to be set (or replaced in the YAML with hardcoded values if you don't feel like messing with the environment variables):
- `KOMODO_URL` - The URL for your Komodo Core deployment, including port but without trailing slash, e.g.: 'https://192.168.1.2:9120', or 'https://my-komodo.example.com'
- `KOMODO_API_KEY` - API Key generated from inside your Komodo Core dashboard, begins with 'K-...'. See instructions at bottom.
- `KOMODO_API_SECRET` - API Secret generated from inside your Komodo Core dashboard, begins with 'S-...'. See instructions at bottom.

> [!NOTE]  
>  - The Options fields will be auto-set through the environment variables, so you don't need to set anything there.
>  - 'allow-insecure:' is set to 'false' by default for security, however if you use self-signed certificates in your Komodo instance, you will need to change it to 'true' for the API requests to succeed.
>  - **Expected API Call Volume:** `(1 + N)` total calls per load/refresh, where `N` is your number of servers. Relatively snappy unless you have a lot of servers.
>  - Cache time is set to 5 minutes, but a value of 1m will be totally fine, given the relatively lightweight API calls made to a local server.
>  - This widget has only been tested with Glance v0.8.4. I can't promise it will work with earlier versions.



## Komodo Container Manager - Stack Monitor
Monitor the status and resource usage of your Komodo stacks and services.
> [!WARNING]
> This is a relatively heavy widget that makes `1 + ( 2 x N )` API calls on every non-cached load/refresh, where `N = the total number of stacks managed by Komodo`.
> In testing, this occasionally caused load times in excess of four minutes. Cache time has been set to 30 minutes to balance freshness vs. load time. You can set it lower, 
> but anything below 5m is unlikely to matter.
> Hopefully this can be improved and optimized, but Komodo's API is currently very recursion-heavy, so until that changes it's inlikely to get much faster. 

#### Preview ( Collapsed )
<img width="1248" height="357" alt="image" src="https://github.com/user-attachments/assets/97829606-fe69-4e28-9fe9-5c23383a64be" />

#### Preview ( Expanded )
<img width="1251" height="997" alt="image" src="https://github.com/user-attachments/assets/43fc1ae2-ad88-4b20-8bca-cfe6f83d4c59" />

#### Preview ( Expanded w/ service details in popover )
<img width="1251" height="1000" alt="image" src="https://github.com/user-attachments/assets/b1f9aea9-e20d-466c-87b4-00bbabfb2c6b" />

### Widget Yaml

<details>

  <summary>View YAML</summary>
  
```yaml
- type: custom-api
  title: Komodo Stacks
  method: POST
  cache: 30m
  allow-insecure: true
  url: http://monitor.sandg.site:9120/read
  body-type: json
  body:
    type: ListStacks
    params: {}
  headers:
    Content-Type: application/json
    X-Api-Key: ${KOMODO_API_KEY}
    X-Api-Secret: ${KOMODO_API_SECRET}
  options:
    base-url: ${KOMODO_URL}
    api-key: ${KOMODO_API_KEY}
    api-secret: ${KOMODO_API_SECRET}
  template: |
    {{ $urlBase := .Options.StringOr "base-url" "" }}
    {{ $apiKey := .Options.StringOr "api-key" "" }}
    {{ $apiSecret := .Options.StringOr "api-secret" "" }}

    <style>
      .list-horizontal-text.no-bullets-komodo-stacks > *:not(:last-child)::after {
          content: none !important;
      }
      .list-horizontal-text.no-bullets-komodo-stacks > *:not(:last-child) {
        margin-right: 1rem;
      }
      .status-badge-positive-komodo-stacks {
        border-radius: var(--border-radius);
        display: inline-block;
        vertical-align: middle;
        padding-left: 1rem;
        color: var(--color-positive);
      }
      .status-badge-neutral-komodo-stacks {
        border-radius: var(--border-radius);
        display: inline-block;
        vertical-align: middle;
        padding-left: 1rem;
        color: var(--color-text-highlight);
      }
      .status-badge-negative-komodo-stacks {
        border-radius: var(--border-radius);
        display: inline-block;
        vertical-align: middle;
        padding-left: 1rem;
        color: var(--color-negative);
      }
      .stack-icon-komodo-stacks {
        width: 10%;
        margin-right: .5rem;
        justify-content: center;
        align-items: center;
      }
      .stack-down-stopped-komodo-stacks {
        position: relative;
      }
      .stack-down-stopped-komodo-stacks {
        position: relative;
      }
      .stack-down-stopped-komodo-stacks::before {
        content: '';
        position: absolute;
        inset: 0;
        background-color: var(--color-widget-background);
        opacity: .5;
        pointer-events: none;
      }
    </style>

    {{/* ======================= API CALLS/VARS TO-DO: optimize/reduce calls ======================= */}}

    {{
      $serverDeets := newRequest (concat $urlBase "/read")
        | withHeader "Content-Type" "application/json"
        | withHeader "X-Api-Key" $apiKey
        | withHeader "X-Api-Secret" $apiSecret
        | withStringBody `{"type":"ListServers","params":{ }}`
        | getResponse
    }}

    {{ $stacks := .JSON.Array "" }}
    {{ $total := len $stacks }}
    {{ $running := 0 }}
    {{ $stopped := 0 }}
    {{ $down := 0 }}
    {{ $unhealthy := 0 }}
    {{ $unknown := 0 }}
    {{ range $stacks }}
      {{ $state := .String "info.state" }}
      {{ if eq $state "running" }}{{ $running = add $running 1 }}{{ end }}
      {{ if eq $state "stopped" }}{{ $stopped = add $stopped 1 }}{{ end }}
      {{ if eq $state "down" }}{{ $down = add $down 1 }}{{ end }}
      {{ if eq $state "unhealthy" }}{{ $unhealthy = add $unhealthy 1 }}{{ end }}
      {{ if or (eq $state "") (eq $state "unknown") }}{{ $unknown = add $unknown 1 }}{{ end }}
    {{ end }}

    <div style="display:flex; align-items:center; gap:12px; margin-bottom: 1rem;">
      <div style="display:flex; align-items:center; gap:12px;">
        <div style="width:40px; height:40px; flex-shrink:0; display:flex; justify-content:center; align-items:center; overflow:hidden;">
          <img src="https://cdn.jsdelivr.net/gh/homarr-labs/dashboard-icons/svg/komodo.svg" width="30" height="30" style="object-fit:contain;">
        </div>
        <div style="flex-grow:1; min-width:0;">
          <h4 class="size-h4 block text-truncate color-highlight" >Komodo Stacks</h4>
          <ul class="list-horizontal-text no-bullets-komodo-stacks items-center" style="vertical-align:middle;">
            <li data-popover-type="text" data-popover-text="{{ $running }}/{{ $total }} running">
              <p style="display:inline-flex;align-items:center;">
                <svg xmlns="http://www.w3.org/2000/svg" 
                viewBox="0 0 256 256" 
                style="height:1.5rem; vertical-align:middle; margin-right:0.5rem; stroke:var(--color-positive); fill:var(--color-positive);">
                  <path d="M224,64H32A16,16,0,0,0,16,80v96a16,16,0,0,0,16,16H224a16,16,0,0,0,16-16V80A16,16,0,0,0,224,64Zm0,112H32V80H224v96Zm-24-48a12,12,0,1,1-12-12A12,12,0,0,1,200,128Z" />
                </svg>
                {{ $running }}
              </p>
            </li>
            <li data-popover-type="text" data-popover-text="{{ $stopped }}/{{ $total }} stopped; {{ $down }}/{{ $total }} down;">
              <p style="display:inline-flex;align-items:center;">
                <svg xmlns="http://www.w3.org/2000/svg" 
                viewBox="0 0 256 256" 
                style="height:1.5rem; vertical-align:middle; margin-right:0.5rem; stroke:var(--color-progress-value); fill:var(--color-progress-value);">
                  <path d="M224,64H32A16,16,0,0,0,16,80v96a16,16,0,0,0,16,16H224a16,16,0,0,0,16-16V80A16,16,0,0,0,224,64Zm0,112H32V80H224v96Zm-24-48a12,12,0,1,1-12-12A12,12,0,0,1,200,128Z" />
                </svg>
                {{ add ( $stopped ) ( $down ) }}
              </p>
            </li>
            <li data-popover-type="text" data-popover-text="{{ $unhealthy }}/{{ $total }} unhealthy; {{ $unknown }}/{{ $total }} unknown">
              <p style="display:inline-flex;align-items:center;">
                <svg xmlns="http://www.w3.org/2000/svg" 
                viewBox="0 0 256 256" 
                style="height:1.5rem; vertical-align:middle; margin-right:0.5rem; stroke:var(--color-negative); fill:var(--color-negative);">
                  <path d="M224,64H32A16,16,0,0,0,16,80v96a16,16,0,0,0,16,16H224a16,16,0,0,0,16-16V80A16,16,0,0,0,224,64Zm0,112H32V80H224v96Zm-24-48a12,12,0,1,1-12-12A12,12,0,0,1,200,128Z" />
                </svg>
                {{ add ( $unknown ) ( $unhealthy ) }}
              </p>
            </li>
          </ul>
        </div>
      </div>
    </div>
    <div class="margin-block-2" style="margin-top: 1rem">
        <ul class="list dynamic-columns list-gap-15 list-with-separator collapsible-container" data-collapse-after="6">
          {{ range $stacks }}

            {{ $stackID := .String "id" }}
            {{ $listItem := . }}
            {{ $serverID := .String "info.server_id" }}
            
            {{
              $stackDeets := newRequest (concat $urlBase "/read")
                | withHeader "Content-Type" "application/json"
                | withHeader "X-Api-Key" $apiKey
                | withHeader "X-Api-Secret" $apiSecret
                | withStringBody (printf `{"type":"GetStack","params":{"stack":"%s"}}` $stackID)
                | getResponse
            }} 
            
            {{
              $serviceDeets := newRequest (concat $urlBase "/read")
                | withHeader "Content-Type" "application/json"
                | withHeader "X-Api-Key" $apiKey
                | withHeader "X-Api-Secret" $apiSecret
                | withStringBody (printf `{"type":"ListStackServices","params":{"stack":"%s"}}` $stackID)
                | getResponse
            }}
            
            {{ $stack := $stackDeets.JSON }}
            {{ $stackID := .String "id" }}
            {{ $state := .String "info.state" }}
            {{ $serverID := .String "info.server_id" }}
            {{ $serverName := "" }}
            {{ $services := $serviceDeets.JSON }}

            {{ range $serverDeets.JSON.Array "" }}
              {{ if eq (.String "id") $serverID }}
                {{ $serverName = .String "name" }}
              {{ end }}
            {{ end }}

            {{ $stackUpdate := false }}
            {{ if or (.Bool "info.services.0.update_available") (gt (len (unique "update_available" (.Array "info.services"))) 1) }}
              {{ $stackUpdate = true }}
            {{ end }}
            
            <li style="border-radius: var(--border-radius); padding: 1rem 1.5rem;border: 1px solid {{ if $stackUpdate }}var(--color-positive);" {{ else if ( or ( eq $state "stopped" ) ( eq $state "down" ) ( eq $state "unknown" ) ) }} var(--color-separator);" class="stack-down-stopped-komodo-stacks" {{ else }} var(--color-progress-value);" {{ end }} data-popover-type="html">
              <div class="grow min-width-0" style="display:flex; align-items:center; gap:12px; vertical-align:middle;" >
                <div class="stack-icon-komodo-stacks">
                  {{ $env := $stack.String "config.environment" }}
                  {{ $iconLink := findSubmatch `ICON_LINK=(.+)` $env }}
                  {{ if and (ne $iconLink "") (ne $iconLink "ZgotmplZ") }}
                    <img src="{{ $iconLink }}" style="height:3rem;vertical-align:middle;margin-right:0.5rem;object-fit:contain;">
                  {{ else }}
                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 448 512" style="height:3rem;vertical-align:middle;margin-right:0.5rem;fill:currentColor;">
                      <path d="M0 48C0 21.5 21.5 0 48 0L400 0c26.5 0 48 21.5 48 48l0 320c0 26.5-21.5 48-48 48L48 416c-26.5 0-48-21.5-48-48L0 48zM32 480l0-16 384 0 0 16c0 17.7-14.3 32-32 32L64 512c-17.7 0-32-14.3-32-32zM96 64C78.3 64 64 78.3 64 96l0 128c0 17.7 14.3 32 32 32l256 0c17.7 0 32-14.3 32-32l0-128c0-17.7-14.3-32-32-32L96 64zM88 352a24 24 0 1 0 0-48 24 24 0 1 0 0 48zm136-24c0 13.3 10.7 24 24 24l112 0c13.3 0 24-10.7 24-24s-10.7-24-24-24l-112 0c-13.3 0-24 10.7-24 24z"/>
                    </svg>
                  {{ end }}
                </div>
                <div>
                  <div data-popover-html>

                    {{ if eq (len (.Array "info.services")) 0 }}
                      <div>
                        <div>
                          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 256 256" style="fill: var(--color-primary); height: 5rem;">
                            <path d="M239.71,125l-16.42-88a16,16,0,0,0-19.61-12.58l-.31.09L150.85,40h-45.7L52.63,24.56l-.31-.09A16,16,0,0,0,32.71,37.05L16.29,125a15.77,15.77,0,0,0,9.12,17.52A16.26,16.26,0,0,0,32.12,144,15.48,15.48,0,0,0,40,141.84V184a40,40,0,0,0,40,40h96a40,40,0,0,0,40-40V141.85a15.5,15.5,0,0,0,7.87,2.16,16.31,16.31,0,0,0,6.72-1.47A15.77,15.77,0,0,0,239.71,125ZM32,128h0L48.43,40,90.5,52.37Zm144,80H136V195.31l13.66-13.65a8,8,0,0,0-11.32-11.32L128,180.69l-10.34-10.35a8,8,0,0,0-11.32,11.32L120,195.31V208H80a24,24,0,0,1-24-24V123.11L107.92,56h40.15L200,123.11V184A24,24,0,0,1,176,208Zm48-80L165.5,52.37,207.57,40,224,128ZM104,140a12,12,0,1,1-12-12A12,12,0,0,1,104,140Zm72,0a12,12,0,1,1-12-12A12,12,0,0,1,176,140Z" />
                          </svg>
                        </div>
                        <div class="flex">
                          <div>Container Name</div>
                          <div class="value-separator"></div>
                          <div class="color-highlight text-very-compact">Many empty</div>
                        </div>
                        <div class="flex">
                          <div>Image</div>
                          <div class="value-separator"></div>
                          <div class="color-highlight text-very-compact">Such no services</div>
                        </div>
                        <div class="flex">
                          <div>Wow</div>
                          <div class="value-separator"></div>
                          <div class="color-highlight text-very-compact">Wow</div>
                        </div>
                        <div class="flex">
                          <div class="color-highlight text-very-compact">Stack returned no services, something other than an array, or was otherwise broken.</div>
                        </div>
                      </div>
                    {{ else }}
                      {{ range $services.Array "" }}
                        {{ $containerState := .String "container.state" }}

                        <div class="margin-top-10 size-h6 align-center text-compact" style="border-bottom: 1px solid var(--color-separator);padding-bottom: 0.25rem;">
                          <p class="color-highlight" style="display:inline-flex;align-items:center;">
                            {{ if .Bool "update_available" }}
                              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 256 256" stroke-width="1"  class="size-6" style="height:1.5rem;vertical-align:middle;margin-right:0.5rem;stroke: var(--color-positive);fill: var(--color-positive);">
                                <path d="M248,128a87.34,87.34,0,0,1-17.6,52.81,8,8,0,1,1-12.8-9.62A71.34,71.34,0,0,0,232,128a72,72,0,0,0-144,0,8,8,0,0,1-16,0,88,88,0,0,1,3.29-23.88C74.2,104,73.1,104,72,104a48,48,0,0,0,0,96H96a8,8,0,0,1,0,16H72A64,64,0,1,1,81.29,88.68,88,88,0,0,1,248,128Zm-69.66,42.34L160,188.69V128a8,8,0,0,0-16,0v60.69l-18.34-18.35a8,8,0,0,0-11.32,11.32l32,32a8,8,0,0,0,11.32,0l32-32a8,8,0,0,0-11.32-11.32Z" />
                              </svg>
                            {{ else }}
                              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 256 256" stroke-width="1"  class="size-6" style="height:1.5rem;vertical-align:middle;margin-right:0.5rem;stroke: var(--color-progress-value);fill: var(--color-progress-value);">
                                <path d="M160,40A88.09,88.09,0,0,0,81.29,88.67,64,64,0,1,0,72,216h88a88,88,0,0,0,0-176Zm0,160H72a48,48,0,0,1,0-96c1.1,0,2.2,0,3.29.11A88,88,0,0,0,72,128a8,8,0,0,0,16,0,72,72,0,1,1,72,72Zm37.66-93.66a8,8,0,0,1,0,11.32l-48,48a8,8,0,0,1-11.32,0l-24-24a8,8,0,0,1,11.32-11.32L144,148.69l42.34-42.35A8,8,0,0,1,197.66,106.34Z" />
                              </svg>
                            {{ end }}
                            {{ .String "service" }}&emsp;
                            <span class="
                              {{ if (or ( eq $containerState "running") ( eq $containerState "created") ) }} 
                                color-positive 
                              {{ else if (or ( eq $containerState "paused") ( eq $containerState "down") ( eq $containerState "stopped") ) }} 
                                color-subdue 
                              {{ else }} 
                                color-negative 
                              {{ end }}" 
                              style="display:inline-flex;align-items:center;">
                                {{ .String "container.status" }}
                            </span>&emsp;
                          </p>
                          <div class="flex">
                            <div>Container Name</div>
                            <div class="value-separator"></div>
                            <div class="color-highlight text-very-compact">{{ .String "container.name" }}</div>
                          </div>
                          <div class="flex">
                            <div>Image</div>
                            <div class="value-separator"></div>
                            <div class="color-highlight text-very-compact">{{ .String "image" }}</div>
                          </div>
                          <div class="flex margin-top-5">
                            <ul class="list-horizontal-text no-bullets-komodo-servers justify-evenly">
                              <li>
                                <p style="display:inline-flex;align-items:center;">
                                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 256 256" fill="currentColor" style="height:1.5rem;vertical-align:middle;margin-right:0.5rem;">
                                    <path d="M152,96H104a8,8,0,0,0-8,8v48a8,8,0,0,0,8,8h48a8,8,0,0,0,8-8V104A8,8,0,0,0,152,96Zm-8,48H112V112h32Zm88,0H216V112h16a8,8,0,0,0,0-16H216V56a16,16,0,0,0-16-16H160V24a8,8,0,0,0-16,0V40H112V24a8,8,0,0,0-16,0V40H56A16,16,0,0,0,40,56V96H24a8,8,0,0,0,0,16H40v32H24a8,8,0,0,0,0,16H40v40a16,16,0,0,0,16,16H96v16a8,8,0,0,0,16,0V216h32v16a8,8,0,0,0,16,0V216h40a16,16,0,0,0,16-16V160h16a8,8,0,0,0,0-16Zm-32,56H56V56H200v95.87s0,.09,0,.13,0,.09,0,.13V200Z" />
                                  </svg>
                                  <span class="color-paragraph">{{ .String "container.stats.cpu_perc" }}</span>
                                </p>
                              </li>
                              <li>
                                <p style="display:inline-flex;align-items:center;">
                                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 256 256" fill="currentColor" style="height:1.5rem;vertical-align:middle;margin-right:0.5rem;">
                                    <path d="M232,56H24A16,16,0,0,0,8,72V200a8,8,0,0,0,16,0V184H40v16a8,8,0,0,0,16,0V184H72v16a8,8,0,0,0,16,0V184h16v16a8,8,0,0,0,16,0V184h16v16a8,8,0,0,0,16,0V184h16v16a8,8,0,0,0,16,0V184h16v16a8,8,0,0,0,16,0V184h16v16a8,8,0,0,0,16,0V72A16,16,0,0,0,232,56ZM24,72H232v96H24Zm88,80a8,8,0,0,0,8-8V96a8,8,0,0,0-8-8H48a8,8,0,0,0-8,8v48a8,8,0,0,0,8,8ZM56,104h48v32H56Zm88,48h64a8,8,0,0,0,8-8V96a8,8,0,0,0-8-8H144a8,8,0,0,0-8,8v48A8,8,0,0,0,144,152Zm8-48h48v32H152Z" />
                                  </svg>
                                  <span class="color-paragraph">{{ .String "container.stats.mem_perc" }}</span>
                                </p>
                              </li>
                              <li>
                                <p style="display:inline-flex;align-items:center;">
                                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 640 640" fill="currentColor" style="height:1.5rem;vertical-align:middle;margin-right:0.5rem;">
                                    <path d="M176 544C96.5 544 32 479.5 32 400C32 336.6 73 282.8 129.9 263.5C128.6 255.8 128 248 128 240C128 160.5 192.5 96 272 96C327.4 96 375.5 127.3 399.6 173.1C413.8 164.8 430.4 160 448 160C501 160 544 203 544 256C544 271.7 540.2 286.6 533.5 299.7C577.5 320 608 364.4 608 416C608 486.7 550.7 544 480 544L176 544zM264 224C241.9 224 224 241.9 224 264L224 296C224 318.1 241.9 336 264 336L280 336C302.1 336 320 318.1 320 296L320 264C320 241.9 302.1 224 280 224L264 224zM256 264C256 259.6 259.6 256 264 256L280 256C284.4 256 288 259.6 288 264L288 296C288 300.4 284.4 304 280 304L264 304C259.6 304 256 300.4 256 296L256 264zM368 224C359.2 224 352 231.2 352 240C352 248.8 359.2 256 368 256L368 320C368 328.8 375.2 336 384 336C392.8 336 400 328.8 400 320L400 240C400 231.2 392.8 224 384 224L368 224zM256 368C247.2 368 240 375.2 240 384C240 392.8 247.2 400 256 400L256 464C256 472.8 263.2 480 272 480C280.8 480 288 472.8 288 464L288 384C288 375.2 280.8 368 272 368L256 368zM376 368C353.9 368 336 385.9 336 408L336 440C336 462.1 353.9 480 376 480L392 480C414.1 480 432 462.1 432 440L432 408C432 385.9 414.1 368 392 368L376 368zM368 408C368 403.6 371.6 400 376 400L392 400C396.4 400 400 403.6 400 408L400 440C400 444.4 396.4 448 392 448L376 448C371.6 448 368 444.4 368 440L368 408z"/>
                                  </svg>
                                  <span class="color-paragraph">{{ .String "container.stats.net_io" }}</span>
                                </p>
                              </li>
                            </ul>
                          </div>
                        </div>
                      {{ end }}
                    {{ end}}
                  </div>
                <div>
                  <a class="size-h3 flex inline-block items-center text-truncate color-primary" href={{ concat $urlBase "stacks/" ( .String "id" ) }} target="_blank">
                    {{ if $stackUpdate }}
                      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 256 256" stroke-width="1" style="height:2rem;vertical-align:middle;margin-right:1rem;stroke: var(--color-positive);fill: var(--color-positive);">
                        <path d="M248,128a87.34,87.34,0,0,1-17.6,52.81,8,8,0,1,1-12.8-9.62A71.34,71.34,0,0,0,232,128a72,72,0,0,0-144,0,8,8,0,0,1-16,0,88,88,0,0,1,3.29-23.88C74.2,104,73.1,104,72,104a48,48,0,0,0,0,96H96a8,8,0,0,1,0,16H72A64,64,0,1,1,81.29,88.68,88,88,0,0,1,248,128Zm-69.66,42.34L160,188.69V128a8,8,0,0,0-16,0v60.69l-18.34-18.35a8,8,0,0,0-11.32,11.32l32,32a8,8,0,0,0,11.32,0l32-32a8,8,0,0,0-11.32-11.32Z" />
                      </svg>
                    {{ else }}
                      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 256 256" stroke-width="1" style="height:2rem;vertical-align:middle;margin-right:1rem;stroke: var(--color-progress-value);fill: var(--color-progress-value);">
                        <path d="M160,40A88.09,88.09,0,0,0,81.29,88.67,64,64,0,1,0,72,216h88a88,88,0,0,0,0-176Zm0,160H72a48,48,0,0,1,0-96c1.1,0,2.2,0,3.29.11A88,88,0,0,0,72,128a8,8,0,0,0,16,0,72,72,0,1,1,72,72Zm37.66-93.66a8,8,0,0,1,0,11.32l-48,48a8,8,0,0,1-11.32,0l-24-24a8,8,0,0,1,11.32-11.32L144,148.69l42.34-42.35A8,8,0,0,1,197.66,106.34Z" />
                      </svg>
                    {{ end }}
                    {{ .String "name"}}
                    <span class="
                      {{ if (or ( eq $state "deploying") ( eq $state "running") ( eq $state "created") ) }}
                        status-badge-positive-komodo-stacks
                      {{ else if (or ( eq $state "paused") ( eq $state "down") ( eq $state "stopped") ) }}
                        status-badge-neutral-komodo-stacks
                      {{ else }}
                        status-badge-negative-komodo-stacks
                      {{ end }} 
                      size-h6 ">
                      {{ $state }}
                    </span>
                  </a>
                  <div>
                    <a class="size-h5 color-subdue" href={{ concat $urlBase "servers/" ( $serverID ) }} target="_blank">Server: <span class="color-base">{{ $serverName }}</span></a>
                  </div>
                  <div class="flex margin-top-5 size-h6">
                    <div>Last Update: </div>
                    <div class="value-separator"></div>
                    <div class="color-highlight text-very-compact">{{ formatTime "DateTime" ( parseTime "unix" (printf "%d" (div ( $stack.Int "updated_at" ) 1000 ) ) ) }}<span class="color-base size-h5"></span></div>
                  </div>
                </div>
              </div>
            </div>
          </li>
        {{ end }}
      </ul>
    </div>
```
</details>

### Environment Variables
There are three environment variables that need to be set (or replaced in the YAML with hardcoded values if you don't feel like messing with the environment variables):
- `KOMODO_URL` - The URL for your Komodo Core deployment, including port but without trailing slash, e.g.: 'https://192.168.1.2:9120', or 'https://my-komodo.example.com'
- `KOMODO_API_KEY` - API Key generated from inside your Komodo Core dashboard, begins with 'K-...'. See instructions at bottom.
- `KOMODO_API_SECRET` - API Secret generated from inside your Komodo Core dashboard, begins with 'S-...'. See instructions at bottom.

> [!NOTE]  
>  - The Options fields will be auto-set through the environment variables, so you don't need to set anything there.
>  - 'allow-insecure:' is set to 'false' by default for security, however if you use self-signed certificates in your Komodo instance, you will need to change it to 'true' for the API requests to succeed.
>  - **Expected API Call Volume:** `1 + (2 x N)` total calls per load/refresh, where `N` is your number of stacks. This can balloon *very* quickly, and there's a pretty hard limit to optimization opportunities with Komodo's current API. This is very much a 'use at your own risk' widget, and if you have a large number of stacks (>20) or a lot of containers per stack, it may crash your client browser (but is unlikely to cause problems for the server -- tests with Komodo running on an i7-6700T saw barely a 5% increase in CPU utilization and no significant memory hit with 50 stacks of 4 containers each). 
>  - This widget has only been tested with Glance v0.8.4. I can't promise it will work with earlier versions.

### Komodo API Key/Secret Instructions
1. #### Click 'Settings' in your Komodo dashboard.<img width="1646" height="817" alt="image" src="https://github.com/user-attachments/assets/6b0022de-d3e9-40e5-8dd7-e55633e5845c" />
2. #### Click 'New API Key ⊕'<img width="1646" height="817" alt="image" src="https://github.com/user-attachments/assets/79559aee-388f-4af2-b9d6-3fc4f5d24c79" />
3. #### Enter a descriptive, memorable name for your API key, set expiry to 'never', and click 'Submit ✓'<img width="1637" height="819" alt="image" src="https://github.com/user-attachments/assets/e5b92473-588b-461a-98c1-b729e682e569" />
4. #### Copy the API key and secret provided by Komodo<img width="1537" height="823" alt="image" src="https://github.com/user-attachments/assets/4eb8af0d-5f62-43d9-8cb9-2f56c9902b1a" />
