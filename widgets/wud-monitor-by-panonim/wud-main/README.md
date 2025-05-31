# What's Up Docker Monitor v2.0 - Main
In this version I tried to focus on improving visuals of the widget. I also added overall info and menu. 

[![2.0 Preview](./wud_main_preview.png)](./wud_main_preview.png)

[![2.0 Preview Other Options](./wud_main_preview2.png)](./wud_main_preview2.png)

### Code - Menu
```yaml
{{ $showInfo := true }}        {{/* Show the overall info summary (Updates, Total, Running) */}}
{{ $showList := true }}        {{/* Show the container list section */}}
{{ $showAll := true }}         {{/* Show all containers (not just ones with updates) */}}
{{ $showUpdateKind := true }}  {{/* Toggle this to false to hide the Update Kind row */}}
{{ $showUpdateTypes := true }} {{/* Toggle this to false to hide the Update Types breakdown */}}
{{ $showLocalVersion := true }}{{/* Toggle this to true to show local version instead of update status */}}
{{ $hasUpdates := false }}     {{/* Set this to true to hide up-to-date message */}}
```
### Widget Code

```yaml
        - type: custom-api
          title: What's Up Docker?
          cache: 1h
          url: http://${WUD_URL}/api/containers/
          method: GET
          template: |
            {{/* You are using WUD Monitor v2.0*/}}
            <style>
              .vertical-separator {
                width: 3px;
                height: 3.5rem;
                margin: 0 1rem;
                border-radius: 2px;
                align-self: center;
                background: linear-gradient(
                  to bottom,
                  transparent,
                  rgba(255, 255, 255, 0.05) 20%,
                  rgba(255, 255, 255, 0.05) 80%,
                  transparent
                );
              }

              .container-info {
                display: flex;
                align-items: flex-start;
                justify-content: space-evenly;
                margin: 1rem 0;
                width: 100%;
              }

              .info-column {
                display: flex;
                flex-direction: column;
                align-items: center;
                min-width: 4rem;
              }

              .info-label {
                font-size: 1.20rem;
                opacity: 0.9;
                text-align: center;
              }

              .info-value {
                font-size: 1.4rem;
                margin-top: 0.25rem;
                text-align: center;
              }

              .info-count {
                font-size: 1.45rem !important;
              }

              /* Container Size: Full */
              @container widget (min-width: 550px) {
                .container-info {
                  justify-content: flex-start;
                  gap: 2rem;
                  padding: 0.5rem;
                }

                .vertical-separator {
                  align-self: flex-start;
                  margin: 0 2rem;
                }

                .info-column {
                  align-items: flex-start;
                }

                .info-label .info-value {
                  text-align: left;
                }
              }
            </style>

            {{/* ======================= MENU ======================= */}}
            {{ $showInfo := true }}         {{/* Show the overall info summary (Updates, Total, Running) */}}
            {{ $showList := true }}         {{/* Show the container list section */}}
            {{ $showAll := false }}         {{/* Show all containers (not just ones with updates) */}}
            {{ $showUpdateKind := true }}   {{/* Toggle this to false to hide the Update Kind row */}}
            {{ $showUpdateTypes := false }} {{/* Toggle this to false to hide the Update Types breakdown */}}
            {{ $showLocalVersion := true }} {{/* Toggle this to true to show local version instead of update status */}}
            {{ $hasUpdates := false }}      {{/* Set this to true to hide up-to-date message */}}

            {{ $containers := .JSON.Array "" }}
            {{ $total := len $containers }}
            {{ $updates := 0 }}
            {{ $running := 0 }}

            {{ range $containers }}
              {{ if .Bool "updateAvailable" }}
                {{ $updates = add $updates 1 }}
              {{ end }}
              {{ if eq (.String "status") "running" }}
                {{ $running = add $running 1 }}
              {{ end }}
            {{ end }}

            {{ if $showInfo }}
              <div class="container-info">
                <div class="info-column">
                  <div class="info-label info-count">Updates</div>
                  <div class="info-value info-count">{{ $updates }}</div>
                </div>

                <div class="vertical-separator"></div>

                <div class="info-column">
                  <div class="info-label info-count">Total</div>
                  <div class="info-value info-count">{{ $total }}</div>
                </div>

                <div class="vertical-separator"></div>

                <div class="info-column">
                  <div class="info-label info-count">Running</div>
                  <div class="info-value info-count">{{ $running }}</div>
                </div>
              </div>
              {{ if $showUpdateTypes }}
              <div class="flex items-center text-center">
                {{ $majorCount := 0 }}
                {{ $minorCount := 0 }}
                {{ $patchCount := 0 }}
                {{ $unknownCount := 0 }}
                {{ range $containers }}
                  {{ if .Bool "updateAvailable" }}
                    {{ $kind := .String "updateKind.semverDiff" }}
                    {{ if eq $kind "major" }}
                      {{ $majorCount = add $majorCount 1 }}
                    {{ else if eq $kind "minor" }}
                      {{ $minorCount = add $minorCount 1 }}
                    {{ else if eq $kind "patch" }}
                      {{ $patchCount = add $patchCount 1 }}
                    {{ else }}
                      {{ $unknownCount = add $unknownCount 1 }}
                    {{ end }}
                  {{ end }}
                {{ end }}
                <div class="flex-1">
                  <div class="size-h3">{{ $majorCount }}</div>
                  <div class="size-h6">MAJOR</div>
                </div>
                <div class="vertical-separator" style="width: 2px; margin: 0 0.25rem;"></div>
                <div class="flex-1">
                  <div class="size-h3">{{ $minorCount }}</div>
                  <div class="size-h6">MINOR</div>
                </div>
                <div class="vertical-separator" style="width: 2px; margin: 0 0.25rem;"></div>
                <div class="flex-1">
                  <div class="size-h3">{{ $patchCount }}</div>
                  <div class="size-h6">PATCH</div>
                </div>
                <div class="vertical-separator" style="width: 2px; margin: 0 0.25rem;"></div>
                <div class="flex-1">
                  <div class="size-h3">{{ $unknownCount }}</div>
                  <div class="size-h6">UNKNOWN</div>
                </div>
              </div>
              {{ end }}
            {{ end }}

            {{ if $showList }}
              <div style="height: 10px;"></div>
              <ul class="list list-gap-10 collapsible-container" data-collapse-after="3">
                {{ range $index, $container := $containers }}
                  {{ if $container.Bool "updateAvailable" }}
                    {{ $hasUpdates = true }}
                    <li>
                      <div class="flex flex-column">
                        {{ $registryName := $container.String "image.registry.name" }}
                        {{ $imageName := $container.String "image.name" }}
                        {{ $hubSource := "#" }}
                        {{ if eq $registryName "hub.public" }}
                          {{ $hubSource = concat "https://hub.docker.com/r/" $imageName }}
                        {{ else if eq $registryName "ghcr.public" }}
                          {{ $hubSource = concat "https://github.com/" $imageName }}
                        {{ end }}

                        <div class="size-h4">
                          {{ if eq ($container.String "status") "running" }}
                            <span class="status-dot color-positive">●</span>
                          {{ else }}
                            <span class="status-dot color-negative">●</span>
                          {{ end }}
                          <a class="color-highlight text-truncate" href="{{ $hubSource }}" target="_blank" rel="noreferrer">{{ $container.String "name" }}</a>
                        </div>

                        <div class="container-info">
                          <div class="info-column">
                            <div class="info-label">Watcher</div>
                            <div class="info-value">{{ $container.String "watcher" }}</div>
                          </div>

                          <div class="vertical-separator"></div>

                          {{ if $showLocalVersion }}
                            <div class="info-column">
                              <div class="info-label">Version</div>
                              <div class="info-value">
                          {{ $localValue := $container.String "updateKind.localValue" }}
                          {{ if ge (len $localValue) 7 }}
                            {{ $isSha256 := eq (slice $localValue 0 7) "sha256:" }}
                            {{ if $isSha256 }}
                              {{ $shortDigest := slice $localValue 7 11 }}
                              <span>dig: {{ $shortDigest }}</span>
                            {{ else }}
                              <span>{{ $localValue }}</span>
                            {{ end }}
                            {{ else }}
                              <span>{{ $localValue }}</span>
                            {{ end }}
                              </div>
                            </div>
                          {{ else }}
                            <div class="info-column">
                              <div class="info-label">Update</div>
                              <div class="info-value">
                                <span class="color-positive">Available</span>
                              </div>
                            </div>
                          {{ end }}
                          {{ if $showUpdateKind }}
                            <div class="vertical-separator"></div>
                            {{ $updateKind := $container.String "updateKind.semverDiff" }}
                            {{ if eq $updateKind "" }}
                              {{ $updateKind = "Unknown" }}
                            {{ end }}
                            {{ $tagValue := $container.String "updateKind.remoteValue" }}
                            {{ if eq $tagValue "" }}
                              {{ $tagValue = "Unknown" }}
                            {{ end }}
                            <div class="info-column">
                              <div class="info-label">{{ $updateKind }}</div>
                              <div class="info-value">
                                {{ if ge (len $tagValue) 7 }}
                                  {{ $isSha256 := eq (slice $tagValue 0 7) "sha256:" }}
                                  {{ if $isSha256 }}
                                    {{ $shortDigest := slice $tagValue 7 11 }}
                                    <span>dig: {{ $shortDigest }}</span>
                                  {{ else }}
                                    <span>{{ $tagValue }}</span>
                                  {{ end }}
                                  {{ else }}
                                    <span>{{ $tagValue }}</span>
                                  {{ end }}
                              </div>
                            </div>
                          {{ end }}
                        </div>
                      </div>
                    </li>
                  {{ end }}
                {{ end }}

                {{ if not $hasUpdates }}
                  <li class="flex items-center justify-center">
                    <span class="color-positive size-h4">All containers are up to date!</span>
                  </li>
                {{ end }}

                {{ if $showAll }}
                  {{ range $index, $container := $containers }}
                    {{ if not ($container.Bool "updateAvailable") }}
                      <li>
                        <div class="flex flex-column">
                          {{ $registryName := $container.String "image.registry.name" }}
                          {{ $imageName := $container.String "image.name" }}
                          {{ $hubSource := "#" }}
                          {{ if eq $registryName "hub.public" }}
                            {{ $hubSource = concat "https://hub.docker.com/r/" $imageName }}
                          {{ else if eq $registryName "ghcr.public" }}
                            {{ $hubSource = concat "https://github.com/" $imageName }}
                          {{ end }}

                          <div class="size-h4">
                            {{ if eq ($container.String "status") "running" }}
                              <span class="status-dot color-positive">●</span>
                            {{ else }}
                              <span class="status-dot color-negative">●</span>
                            {{ end }}
                            <a class="color-highlight text-truncate" href="{{ $hubSource }}" target="_blank" rel="noreferrer">{{ $container.String "name" }}</a>
                          </div>

                          <div class="container-info">
                            <div class="info-column">
                              <div class="info-label">Watcher</div>
                              <div class="info-value">{{ $container.String "watcher" }}</div>
                            </div>

                            <div class="vertical-separator"></div>

                            {{ if $showLocalVersion }}
                              <div class="info-column">
                                <div class="info-label">Version</div>
                                <div class="info-value">
                                  {{ $localValue := $container.String "updateKind.localValue" }}
                                  {{ if eq $localValue "" }}
                                    {{ $localValue = $container.String "updatekind.remoteValue" }}
                                    {{ if eq $localValue "" }}
                                      {{ $imageDigest := $container.String "image.digest.repo" }}
                                      {{ $isSha256 := eq (slice $imageDigest 0 7) "sha256:" }}
                                      {{ if $isSha256 }}
                                        {{ $shortDigest := slice $imageDigest 7 11 }}
                                        <span>dig: {{ $shortDigest }}</span>
                                      {{ else }}
                                        <span>-</span>
                                      {{ end }}
                                    {{ else }}
                                      <span>{{ $localValue }}</span>
                                    {{ end }}
                                  {{ else }}
                                    {{ $isSha256 := eq (slice $localValue 0 7) "sha256:" }}
                                    {{ if $isSha256 }}
                                      {{ $shortDigest := slice $localValue 7 11 }}
                                      <span>dig: {{ $shortDigest }}</span>
                                    {{ else }}
                                      <span>{{ $localValue }}</span>
                                    {{ end }}
                                  {{ end }}
                                </div>
                              </div>
                            {{ else }}
                              <div class="info-column">
                                <div class="info-label">Update</div>
                                <div class="info-value"><span class="color-negative">Not Found</span></div>
                              </div>
                            {{ end }}
                          </div>
                        </div>
                      </li>
                    {{ end }}
                  {{ end }}
                {{ end }}
              </ul>
            {{ end }}
```

<hr>

**Made by:** Artur Flis  
**Contact:** @blue.dev on the project’s Discord

<hr>

### Contributors v2.0

- [**Gaodes**](https://github.com/gaodes) – Helped with css and general layout.
