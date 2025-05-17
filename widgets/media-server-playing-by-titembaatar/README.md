* [Introduction](#introduction)
* [Preview](#preview)
    * [Full-Size Column](#full-size-column)
    * [Small-Size Column](#vertical-layout-for-small-size-column)
    * [Compact Mode](#compact-mode)
* [Environment Variables](#environment-variables)
* [Secrets](#secrets)
* [Options](#options)
* [Widget YAML](#widget-yaml)
    * [Plex YAML](#plex-yaml)
    * [Tautulli YAML](#tautulli-yaml)
    * [Jellyfin YAML](#jellyfin-yaml)

## Introduction
This is a collection of widgets for various media servers.

> [!NOTE]
>
> The widget has been updated to `glance v0.8.0`.
> Ensure you update to at least this version.

Tested with `Plex`, `Tautulli`, and `Jellyfin`. For `Emby`, the API appears (from what I glanced)
to be identical to `Jellyfin`. You should be able to use the [Jellyfin Widget](#jellyfin-yaml) with
an `Emby` URL and API key. If you encounter any issues, please open an issue, tag me, and I’ll investigate further.

The appearance is consistent across all media servers.  
Customisation can be applied using the `options:` field. See [Options](#options) for more details.

## Preview
### Full-Size Column
![Preview](preview.png)

### Vertical Layout for Small-Size Column
![Preview small](preview-small.png)

### Compact Mode
![Preview Compact](preview-compact.png)

## Environment Variables

> [!IMPORTANT]
>
> For URLs, you **MUST** include `http://` or `https://`.
> Do **NOT** include a trailing `/` at the end of your URLs.

### Plex
* `PLEX_URL` - The Plex URL, e.g., `http://<ip_address>:<port>` or `https://<domain>`
* `PLEX_TOKEN` - The Plex token; follow [this guide](https://support.plex.tv/articles/204059436-finding-an-authentication-token-x-plex-token/) if you need help obtaining it.

### Tautulli
* `TAUTULLI_URL` - The Tautulli URL, e.g., `http://<ip_address>:<port>` or `https://<domain>`
* `TAUTULLI_KEY` - The Tautulli API key, found in `Settings` -> `Web Interface` -> `API Key`

### Jellyfin
* `JELLYFIN_URL` - The Jellyfin URL, e.g., `http://<ip_address>:<port>` or `https://<domain>`
* `JELLYFIN_KEY` - The Jellyfin API key, available in `Administration` -> `Dashboard` -> `API Keys`

## Secrets
Since `v0.8.0`, you can use Docker secrets instead of environment variables.  
If you do, replace `${YOUR_API_KEY}` with `${secret:your-api-key-secret}` throughout the code.

## Options
Since `v0.8.0`, you can use the `options:` field to customise the widget.

> [!CAUTION]
>
> Enabling thumbnails **will** expose your token/API keys in the HTML.
> Do **not** enable this in production or on internet-exposed services.

Default options are:
```yaml
options:
  small-column: false      # `true` if using the widget in a small column
  compact: true            # `false` for a more spread-out layout
  play-state: "text"       # Options: "text" (plain text), "indicator" (pulsing dot when playing)
  show-thumbnail: false    # `true` to show thumbnails
  show-paused: false       # `true` to display paused items
  show-progress-bar: false # `true` to display the progress bar
  show-progress-info: true # `false` to hide progress info; requires `show-progress-bar`
```

> [!NOTE]
>
> The progress bar is cosmetic, using CSS animation and time calculations.
> It is **not** dynamic and does not auto-refresh when reaching 100%.

## Widget YAML
### Plex YAML
```yaml
- type: custom-api
  title: plex
  cache: 5m
  url: ${PLEX_URL}/status/sessions
  headers:
    Accept: application/json
    X-Plex-Token: ${PLEX_TOKEN}
  template: |
    {{ $isSmallColumn := .Options.BoolOr "small-column" false }}
    {{ $isCompact := .Options.BoolOr "compact" true }}
    {{ $playState := .Options.StringOr "play-state" "text" }}
    {{ $showThumbnail := .Options.BoolOr "show-thumbnail" false }}
    {{ $showPaused := .Options.BoolOr "show-paused" false }}
    {{ $showProgressBar := .Options.BoolOr "show-progress-bar" false }}
    {{ $showProgressInfo := .Options.BoolOr "show-progress-info" true }}

    {{ if eq .Response.StatusCode 200 }}
      {{ $sessions := .JSON.Array "MediaContainer.Metadata" }}
      {{ if eq (len $sessions) 0 }}
        <p>Nothing is playing right now.</p>
      {{ else }}
        {{ if $isSmallColumn }}
          <div class="flex flex-column gap-10">
        {{ else }}
          <div class="gap-10" style="display: grid; grid-template-columns: repeat(2, 1fr);">
        {{ end }}
          {{ range $i, $session := $sessions }}
            {{/* WIDGET VARIABLES BEGIN */}}
              {{ $mediaType := $session.String "type" }}
              {{ $isPlaying := eq ($session.String "Player.state") "playing" }}
              {{ $state := $session.String "Player.state" }}

              {{ $isMovie := eq $mediaType "movie" }}
              {{ $isShows := eq $mediaType "episode" }}
              {{ $isMusic := eq $mediaType "track" }}

              {{ $userName := $session.String "User.title" }}
              {{ $movieTitle := $session.String "title" }}
              {{ $showTitle := $session.String "grandparentTitle" }}
              {{ $showSeason := $session.String "parentIndex" }}
              {{ $showEpisode := $session.String "index" }}
              {{ $episodeTitle := $session.String "title" }}
              {{ $artist := $session.String "grandparentTitle" }}
              {{ $albumTitle := $session.String "parentTitle" }}
              {{ $songTitle := $session.String "title" }}
              {{ $default := $session.String "title" }}

              {{ $thumbPath := $session.String "thumb" }}
              {{ if or $isShows $isMusic }}
                {{ $thumbPath = $session.String "parentThumb" }}
              {{ end }}
              {{ $thumbURL := concat "${PLEX_URL}" $thumbPath "?X-Plex-Token=${PLEX_TOKEN}" }}

              {{ $duration := $session.Float "duration" }}
              {{ $offset := $session.Float "viewOffset" }}
              {{ $progress := mul 100 (div $offset $duration) | toInt }}
              {{ $remainingSeconds := div (sub $duration $offset) 1000 | toInt }}
              {{ $now := now | formatTime "15:04:05" | parseLocalTime "15:04:05" }}
              {{ $endTime := $now.Add (duration (printf "%ds" $remainingSeconds)) | formatTime "15:04" }}
            {{/* WIDGET VARIABLES END */}}
            {{/* WIDGET TEMPLATE BEGIN */}}
            {{ if or $isPlaying $showPaused }}
              <div class="card gap-5">
                <div class="flex items-center gap-10 size-h3">
                  <span class="color-primary">{{ $userName }}</span>
                  {{ if eq $playState "text" }}
                    <span {{ if $isPlaying }}class="color-primary"{{ end }}>
                      [{{ $state }}]
                    </span>
                  {{ else if eq $playState "indicator" }}
                    <style>
                      @keyframes pulse {
                        0% { box-shadow: 0 0 0 0 var(--color-text-base); }
                        40% { box-shadow: 0 0 0 4px transparent; }
                        100% { box-shadow: 0 0 0 4px transparent; }
                      }
                    </style>
                    <div style="
                      height: .7rem;
                      width: .7rem;
                      {{ if $isPlaying }}
                        animation: pulse 5s infinite;
                        background: var(--color-primary);
                      {{ else }}
                        background: var(--color-text-base-muted);
                      {{ end }}
                      border-radius: 100%;">
                    </div>
                  {{ end }}
                </div>

                <hr class="margin-bottom-5" />

                <div class="flex items-center gap-10" style="align-items: stretch;">
                  {{ if eq $showThumbnail true }}
                    <img src="{{ $thumbURL }}"
                      alt="{{ $default }} thumbnail"
                      class="shrink-0"
                      loading="lazy"
                      style="
                        max-width: 7.5rem;
                        border: 2px solid var(--color-primary);
                        border-radius: var(--border-radius);
                        object-fit: cover;
                        {{ if $isCompact }}aspect-ratio: 1;{{ end }}"/>
                  {{ end }}
                  <ul class="flex flex-column grow justify-evenly" style="width: 0;">
                    {{ if $isMovie }}
                      <li>{{ $movieTitle }}</li>
                    {{ else if $isShows }}
                      {{ if $isCompact }}
                        <ul class="list-horizontal-text flex-nowrap">
                      {{ end }}
                          <li class="text-truncate">{{ $showTitle }}</li>
                          <li class="text-truncate shrink-0">{{ concat "S" $showSeason "E" $showEpisode }}</li>
                      {{ if $isCompact }}
                        </ul>
                      {{ end }}
                      <li class="text-truncate">{{ $episodeTitle }}</li>
                    {{ else if $isMusic }}
                      {{ if $isCompact }}
                        <ul class="list-horizontal-text flex-nowrap">
                      {{ end }}
                          <li class="text-truncate shrink-0">{{ $artist }}</li>
                          <li class="text-truncate">{{ $albumTitle }}</li>
                      {{ if $isCompact }}
                        </ul>
                      {{ end }}
                      <li class="text-truncate">{{ $songTitle }}</li>
                    {{ else }}
                      <li>{{ $default }}</li>
                    {{ end }}

                    <li>
                      {{ if and $isPlaying $showProgressBar }}
                        <div class="flex gap-10 items-center">
                          <div class="grow" style="
                            height: 1rem;
                            max-width: 32rem;
                            border: 1px solid var(--color-text-base);
                            border-radius: var(--border-radius);
                            overflow: hidden;">
                            <style>
                              @keyframes progress-animation { to { width: 100%; } }
                            </style>
                            <div
                              data-progress="{{ $progress }}"
                              data-remaining="{{ $remainingSeconds }}"
                              style="
                                height: 100%;
                                width: {{ $progress }}%;
                                background: var(--color-primary);
                                border-radius: 3px;
                                transition: width 1s linear;
                                animation: progress-animation {{ $remainingSeconds }}s linear forwards;">
                            </div>
                          </div>
                          {{ if $showProgressInfo }}
                            <p>
                              {{ if and (not $isCompact) (not $isSmallColumn) }}
                                ends at 
                              {{ end }}
                              {{ $endTime }}
                            </p>
                          {{ end }}
                        </div>
                      {{ end }}
                    </li>
                  </ul>
                </div>
              </div>
            {{ end }}
            {{/* WIDGET TEMPLATE END */}}
          {{ end }}
        </div>
      {{ end }}
    {{ else }}
      <p>Failed to fetch Plex sessions</p>
    {{ end }}
```

### Tautulli YAML
```yaml
- type: custom-api
  title: tautulli
  cache: 5m
  url: ${TAUTULLI_URL}/api/v2
  parameters:
    apikey: ${TAUTULLI_KEY}
    cmd: get_activity
  template: |
    {{ $isSmallColumn := .Options.BoolOr "small-column" false }}
    {{ $isCompact := .Options.BoolOr "compact" true }}
    {{ $playState := .Options.StringOr "play-state" "text" }}
    {{ $showThumbnail := .Options.BoolOr "show-thumbnail" false }}
    {{ $showPaused := .Options.BoolOr "show-paused" false }}
    {{ $showProgressBar := .Options.BoolOr "show-progress-bar" false }}
    {{ $showProgressInfo := .Options.BoolOr "show-progress-info" true }}

    {{ if eq .Response.StatusCode 200 }}
      {{ $sessions := .JSON.Array "response.data.sessions" }}
      {{ if eq (len $sessions) 0 }}
        <p>Nothing is playing right now.</p>
      {{ else }}
        {{ if $isSmallColumn }}
          <div class="flex flex-column gap-10">
        {{ else }}
          <div class="gap-10" style="display: grid; grid-template-columns: repeat(2, 1fr);">
        {{ end }}
          {{ range $i, $session := $sessions }}
            {{/* WIDGET VARIABLES BEGIN */}}
              {{ $mediaType := $session.String "media_type" }}
              {{ $isPlaying := eq ($session.String "state") "playing" }}
              {{ $state := $session.String "state" }}

              {{ $isMovie := eq $mediaType "movie" }}
              {{ $isShows := eq $mediaType "episode" }}
              {{ $isMusic := eq $mediaType "track" }}

              {{ $userName := $session.String "user" }}
              {{ $movieTitle := $session.String "title" }}
              {{ $showTitle := $session.String "grandparent_title" }}
              {{ $showSeason := $session.String "parent_media_index" }}
              {{ $showEpisode := $session.String "media_index" }}
              {{ $episodeTitle := $session.String "title" }}
              {{ $artist := $session.String "grandparent_title" }}
              {{ $albumTitle := $session.String "parent_title" }}
              {{ $songTitle := $session.String "title" }}
              {{ $default := $session.String "title" }}

              {{ $thumbPath := $session.String "thumb" }}
              {{ if or $isShows $isMusic }}
                {{ $thumbPath = $session.String "parent_thumb" }}
              {{ end }}
              {{ $thumbURL := concat "${TAUTULLI_URL}/api/v2?apikey=${TAUTULLI_KEY}&cmd=pms_image_proxy&img=" $thumbPath }}

              {{ $duration := $session.Float "duration" }}
              {{ $offset := $session.Float "view_offset" }}
              {{ $progress := $session.Int "progress_percent" }}
              {{ $remainingSeconds := div (sub $duration $offset) 1000 | toInt }}
              {{ $now := now | formatTime "15:04:05" | parseLocalTime "15:04:05" }}
              {{ $endTime := $now.Add (duration (printf "%ds" $remainingSeconds)) | formatTime "15:04" }}
            {{/* WIDGET VARIABLES END */}}
            {{/* WIDGET TEMPLATE BEGIN */}}
            {{ if or $isPlaying $showPaused }}
              <div class="card gap-5">
                <div class="flex items-center gap-10 size-h3">
                  <span class="color-primary">{{ $userName }}</span>
                  {{ if eq $playState "text" }}
                    <span {{ if $isPlaying }}class="color-primary"{{ end }}>
                      [{{ $state }}]
                    </span>
                  {{ else if eq $playState "indicator" }}
                    <style>
                      @keyframes pulse {
                        0% { box-shadow: 0 0 0 0 var(--color-text-base); }
                        40% { box-shadow: 0 0 0 4px transparent; }
                        100% { box-shadow: 0 0 0 4px transparent; }
                      }
                    </style>
                    <div style="
                      height: .7rem;
                      width: .7rem;
                      {{ if $isPlaying }}
                        animation: pulse 5s infinite;
                        background: var(--color-primary);
                      {{ else }}
                        background: var(--color-text-base-muted);
                      {{ end }}
                      border-radius: 100%;">
                    </div>
                  {{ end }}
                </div>

                <hr class="margin-bottom-5" />

                <div class="flex items-center gap-10" style="align-items: stretch;">
                  {{ if eq $showThumbnail true }}
                    <img src="{{ $thumbURL }}"
                      alt="{{ $default }} thumbnail"
                      class="shrink-0"
                      loading="lazy"
                      style="
                        max-width: 7.5rem;
                        border: 2px solid var(--color-primary);
                        border-radius: var(--border-radius);
                        object-fit: cover;
                        {{ if $isCompact }}aspect-ratio: 1;{{ end }}"/>
                  {{ end }}
                  <ul class="flex flex-column grow justify-evenly" style="width: 0;">
                    {{ if $isMovie }}
                      <li>{{ $movieTitle }}</li>
                    {{ else if $isShows }}
                      {{ if $isCompact }}
                        <ul class="list-horizontal-text flex-nowrap">
                      {{ end }}
                          <li class="text-truncate">{{ $showTitle }}</li>
                          <li class="text-truncate shrink-0">{{ concat "S" $showSeason "E" $showEpisode }}</li>
                      {{ if $isCompact }}
                        </ul>
                      {{ end }}
                      <li class="text-truncate">{{ $episodeTitle }}</li>
                    {{ else if $isMusic }}
                      {{ if $isCompact }}
                        <ul class="list-horizontal-text flex-nowrap">
                      {{ end }}
                          <li class="text-truncate shrink-0">{{ $artist }}</li>
                          <li class="text-truncate">{{ $albumTitle }}</li>
                      {{ if $isCompact }}
                        </ul>
                      {{ end }}
                      <li class="text-truncate">{{ $songTitle }}</li>
                    {{ else }}
                      <li>{{ $default }}</li>
                    {{ end }}

                    <li>
                      {{ if and $isPlaying $showProgressBar }}
                        <div class="flex gap-10 items-center">
                          <div class="grow" style="
                            height: 1rem;
                            max-width: 32rem;
                            border: 1px solid var(--color-text-base);
                            border-radius: var(--border-radius);
                            overflow: hidden;">
                            <style>
                              @keyframes progress-animation { to { width: 100%; } }
                            </style>
                            <div
                              data-progress="{{ $progress }}"
                              data-remaining="{{ $remainingSeconds }}"
                              style="
                                height: 100%;
                                width: {{ $progress }}%;
                                background: var(--color-primary);
                                border-radius: 3px;
                                transition: width 1s linear;
                                animation: progress-animation {{ $remainingSeconds }}s linear forwards;">
                            </div>
                          </div>
                          {{ if $showProgressInfo }}
                            <p>
                              {{ if and (not $isCompact) (not $isSmallColumn) }}
                                ends at 
                              {{ end }}
                              {{ $endTime }}
                            </p>
                          {{ end }}
                        </div>
                      {{ end }}
                    </li>
                  </ul>
                </div>
              </div>
            {{ end }}
            {{/* WIDGET TEMPLATE END */}}
          {{ end }}
        </div>
      {{ end }}
    {{ else }}
      <p>Failed to fetch Tautulli sessions</p>
    {{ end }}
```

### Jellyfin YAML
```yaml
- type: custom-api
  title: jellyfin
  cache: 5m
  url: ${JELLYFIN_URL}/Sessions
  parameters:
    api_key: ${JELLYFIN_KEY}
    activeWithinSeconds: 30
  template: |
    {{ $isSmallColumn := .Options.BoolOr "small-column" false }}
    {{ $isCompact := .Options.BoolOr "compact" true }}
    {{ $playState := .Options.StringOr "play-state" "text" }}
    {{ $showThumbnail := .Options.BoolOr "show-thumbnail" false }}
    {{ $showPaused := .Options.BoolOr "show-paused" false }}
    {{ $showProgressBar := .Options.BoolOr "show-progress-bar" false }}
    {{ $showProgressInfo := .Options.BoolOr "show-progress-info" true }}

    {{ if eq .Response.StatusCode 200 }}
      {{ $sessions := .JSON.Array "" }}
      {{ if eq (len $sessions) 0 }}
        <p>Nothing is playing right now.</p>
      {{ else }}
        {{ if $isSmallColumn }}
          <div class="flex flex-column gap-10">
        {{ else }}
          <div class="gap-10" style="display: grid; grid-template-columns: repeat(2, 1fr);">
        {{ end }}
          {{ range $i, $session := $sessions }}
            {{/* WIDGET VARIABLES BEGIN */}}
              {{ $mediaType := $session.String "NowPlayingItem.Type" }}
              {{ $isPaused := $session.Bool "PlayState.IsPaused" }}
              {{ $isPlaying := not $isPaused }}
              {{ $state := "playing" }}

              {{ if $isPaused }}
                {{ $state = "paused" }}
              {{ end }}

              {{ $isMovie := eq $mediaType "Movie" }}
              {{ $isShows := eq $mediaType "Episode" }}
              {{ $isMusic := eq $mediaType "Audio" }}

              {{ $userName := $session.String "UserName" }}
              {{ $movieTitle := $session.String "NowPlayingItem.Name" }}
              {{ $showTitle := $session.String "NowPlayingItem.SeriesName" }}
              {{ $showSeason := $session.String "NowPlayingItem.ParentIndexNumber" }}
              {{ $showEpisode := $session.String "NowPlayingItem.IndexNumber" }}
              {{ $episodeTitle := $session.String "NowPlayingItem.Name" }}
              {{ $artist := $session.String "NowPlayingItem.AlbumArtist" }}
              {{ $albumTitle := $session.String "NowPlayingItem.Album" }}
              {{ $songTitle := $session.String "NowPlayingItem.Name" }}
              {{ $default := $session.String "NowPlayingItem.Name" }}

              {{ $thumbID := $session.String "NowPlayingItem.Id" }}
              {{ if $isShows }}
                {{ $thumbID = $session.String "NowPlayingItem.ParentId" }}
              {{ end }}
              {{ $thumbURL := concat "${JELLYFIN_URL}/Items/" $thumbID "/Images/Primary?api_key=${JELLYFIN_KEY}" }}

              {{ $duration := $session.Float "NowPlayingItem.RunTimeTicks" }}
              {{ $offset := $session.Float "PlayState.PositionTicks" }}
              {{ $progress := mul 100 (div $offset $duration) | toInt }}
              {{ $remainingSeconds := div (sub $duration $offset) 10000000 | toInt }}
              {{ $now := now | formatTime "15:04:05" | parseLocalTime "15:04:05" }}
              {{ $endTime := $now.Add (duration (printf "%ds" $remainingSeconds)) | formatTime "15:04" }}
            {{/* WIDGET VARIABLES END */}}
            {{/* WIDGET TEMPLATE BEGIN */}}
            {{ if or $isPlaying $showPaused }}
              <div class="card gap-5">
                <div class="flex items-center gap-10 size-h3">
                  <span class="color-primary">{{ $userName }}</span>
                  {{ if eq $playState "text" }}
                    <span {{ if $isPlaying }}class="color-primary"{{ end }}>
                      [{{ $state }}]
                    </span>
                  {{ else if eq $playState "indicator" }}
                    <style>
                      @keyframes pulse {
                        0% { box-shadow: 0 0 0 0 var(--color-text-base); }
                        40% { box-shadow: 0 0 0 4px transparent; }
                        100% { box-shadow: 0 0 0 4px transparent; }
                      }
                    </style>
                    <div style="
                      height: .7rem;
                      width: .7rem;
                      {{ if $isPlaying }}
                        animation: pulse 5s infinite;
                        background: var(--color-primary);
                      {{ else }}
                        background: var(--color-text-base-muted);
                      {{ end }}
                      border-radius: 100%;">
                    </div>
                  {{ end }}
                </div>

                <hr class="margin-bottom-5" />

                <div class="flex items-center gap-10" style="align-items: stretch;">
                  {{ if eq $showThumbnail true }}
                    <img src="{{ $thumbURL }}"
                      alt="{{ $default }} thumbnail"
                      class="shrink-0"
                      loading="lazy"
                      style="
                        max-width: 7.5rem;
                        border: 2px solid var(--color-primary);
                        border-radius: var(--border-radius);
                        object-fit: cover;
                        {{ if $isCompact }}aspect-ratio: 1;{{ end }}"/>
                  {{ end }}
                  <ul class="flex flex-column grow justify-evenly" style="width: 0;">
                    {{ if $isMovie }}
                      <li>{{ $movieTitle }}</li>
                    {{ else if $isShows }}
                      {{ if $isCompact }}
                        <ul class="list-horizontal-text flex-nowrap">
                      {{ end }}
                          <li class="text-truncate">{{ $showTitle }}</li>
                          <li class="text-truncate shrink-0">{{ concat "S" $showSeason "E" $showEpisode }}</li>
                      {{ if $isCompact }}
                        </ul>
                      {{ end }}
                      <li class="text-truncate">{{ $episodeTitle }}</li>
                    {{ else if $isMusic }}
                      {{ if $isCompact }}
                        <ul class="list-horizontal-text flex-nowrap">
                      {{ end }}
                          <li class="text-truncate shrink-0">{{ $artist }}</li>
                          <li class="text-truncate">{{ $albumTitle }}</li>
                      {{ if $isCompact }}
                        </ul>
                      {{ end }}
                      <li class="text-truncate">{{ $songTitle }}</li>
                    {{ else }}
                      <li>{{ $default }}</li>
                    {{ end }}

                    <li>
                      {{ if and $isPlaying $showProgressBar }}
                        <div class="flex gap-10 items-center">
                          <div class="grow" style="
                            height: 1rem;
                            max-width: 32rem;
                            border: 1px solid var(--color-text-base);
                            border-radius: var(--border-radius);
                            overflow: hidden;">
                            <style>
                              @keyframes progress-animation { to { width: 100%; } }
                            </style>
                            <div
                              data-progress="{{ $progress }}"
                              data-remaining="{{ $remainingSeconds }}"
                              style="
                                height: 100%;
                                width: {{ $progress }}%;
                                background: var(--color-primary);
                                border-radius: 3px;
                                transition: width 1s linear;
                                animation: progress-animation {{ $remainingSeconds }}s linear forwards;">
                            </div>
                          </div>
                          {{ if $showProgressInfo }}
                            <p>
                              {{ if and (not $isCompact) (not $isSmallColumn) }}
                                ends at 
                              {{ end }}
                              {{ $endTime }}
                            </p>
                          {{ end }}
                        </div>
                      {{ end }}
                    </li>
                  </ul>
                </div>
              </div>
            {{ end }}
            {{/* WIDGET TEMPLATE END */}}
          {{ end }}
        </div>
      {{ end }}
    {{ else }}
      <p>Failed to fetch Jellyfin sessions</p>
    {{ end }}
```
