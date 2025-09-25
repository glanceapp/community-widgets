* [Introduction](#introduction)
* [Preview](#preview)
* [Environment Variables](#environment-variables)
* [Secrets](#secrets)
* [Options](#options)
* [Widget YAML](#widget-yaml)
* [🍻 Cheers](#cheers)

## Introduction
This is a widget for various media servers to display currently active sessions.

> [!NOTE]
>
> The widget has been updated to `glance v0.8.0`.
> Ensure you update to at least this version.

Tested with `Plex`, `Tautulli`, `Jellyfin`, and `Emby`.
If you encounter any issues, please open an issue, tag me, and I’ll investigate further.

The appearance is consistent across all media servers.  
Customisation can be applied using the `options:` field. See [Options](#options) for more details.

## Preview
### Full-Size Column
![Preview](preview.png)

### Small-Size Column
![Preview Small](preview-small.png)

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
* `TAUTULLI_KEY` - The Tautulli API key, available in `Settings` -> `Web Interface `-> `API key`

### Jellyfin
* `JELLYFIN_URL` - The Jellyfin URL, e.g., `http://<ip_address>:<port>` or `https://<domain>`
* `JELLYFIN_KEY` - The Jellyfin API key, available in `Administration` -> `Dashboard` -> `API Keys`

### Emby
* `EMBY_URL` - The Emby URL, e.g., `http://<ip_address>:<port>` or `https://<domain>`
* `EMBY_KEY` - The Emby API key, available in `⚙️ icon in top-right` -> `Advanced` -> `API Keys`

## Secrets
Since `v0.8.0`, you can use Docker secrets instead of environment variables. See [v0.8.0 Release Notes](https://github.com/glanceapp/glance/releases/tag/v0.8.0#g-rh-5) for more information.  
If you do, replace `${YOUR_API_KEY}` with `${secret:your-api-key-secret}`.

## Options
Since `v0.8.0`, you can use the `options:` field to customise the widget.  
See [v0.8.0 Release Notes](https://github.com/glanceapp/glance/releases/tag/v0.8.0#g-rh-15) for more information.

> [!CAUTION]
>
> Enabling thumbnails **will** expose your token/API keys in the HTML.
> Do **not** enable this in production or on internet-exposed services.

Default options are:
```yaml
options:
  # Required options
  media-server: "plex"     # Your media server; "plex", "tautulli", "jellyfin", "emby"
  base-url: ${PLEX_URL}    # Your environment-variables for the URL
  api-key: ${PLEX_TOKEN}   # Your environment-variables for the API key/token. Can a secret as well `${secret:plex-token}`

  # Optional options
  small-column: false      # `true` if using the widget in a small column
  compact: true            # `false` for a more spread-out layout
  play-state: "indicator"  # `"text"` for plain text, `"indicator"` for a pulsing dot
  show-thumbnail: false    # `true` to show thumbnails
  full-thumbnail: false    # `true` to show full-size thumbnails in compact mode
  show-paused: false       # `true` to display paused items
  show-progress-bar: false # `true` to display the progress bar
  show-progress-info: true # `false` to hide progress info; requires `show-progress-bar`
  time-format: "15:04"     # `"03:04pm"` if you prefer 12h format
```

> [!NOTE]
>
> The progress bar is cosmetic, using CSS animation and time calculations.
> It is **not** dynamic and does not auto-refresh when reaching 100%.

## Widget YAML
```yaml
- type: custom-api
  title: Media Server
  cache: 5m
  options:
    media-server: "plex"
    base-url: ${PLEX_URL}
    api-key: ${PLEX_TOKEN}
    small-column: false
    compact: true
    play-state: "indicator"
    show-thumbnail: false
    full-thumbnail: false
    show-paused: false
    show-progress-bar: false
    show-progress-info: true
    time-format: "15:04"
  template: |
    {{ $mediaServer := .Options.StringOr "media-server" "" }}
    {{ $baseURL := .Options.StringOr "base-url" "" }}
    {{ $apiKey := .Options.StringOr "api-key" "" }}

    {{ define "errorMsg" }}
      <div class="widget-error-header">
        <div class="color-negative size-h3">ERROR</div>
        <svg class="widget-error-icon" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5">
          <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126ZM12 15.75h.007v.008H12v-.008Z"></path>
        </svg>
      </div>
      <p class="break-all">{{ . }}</p>
    {{ end }}

    {{ if or
      (eq $mediaServer "")
      (eq $baseURL "")
      (eq $apiKey "")
    }}
      {{ template "errorMsg" "Some required options are not set" }}
    {{ else }}

    {{ $isSmallColumn := .Options.BoolOr "small-column" false }}
    {{ $isCompact := .Options.BoolOr "compact" true }}
    {{ $playState := .Options.StringOr "play-state" "indicator" }}
    {{ $showThumbnail := .Options.BoolOr "show-thumbnail" false }}
    {{ $fullThumbnail := .Options.BoolOr "full-thumbnail" false }}
    {{ $showPaused := .Options.BoolOr "show-paused" false }}
    {{ $showProgressBar := .Options.BoolOr "show-progress-bar" false }}
    {{ $showProgressInfo := .Options.BoolOr "show-progress-info" true }}
    {{ $timeFormat := .Options.StringOr "time-format" "15:04" }}

    {{ $userID := "" }}
    {{ $sessionsRequestURL := "" }}
    {{ $sessionsCall := "" }}
    {{ $sessions := "" }}
    {{ $activeSessions := 0 }}

    {{ if eq $mediaServer "plex" }}
      {{ $sessionsRequestURL = concat $baseURL "/status/sessions" }}
      {{ $sessionsCall = newRequest $sessionsRequestURL
        | withHeader "Accept" "application/json"
        | withHeader "X-Plex-Token" $apiKey
        | getResponse }}

      {{ if $sessionsCall.JSON.Exists "MediaContainer" }}
        {{ $sessions = $sessionsCall.JSON.Array "MediaContainer.Metadata" }}
        {{ $activeSessions = len $sessions }}
      {{ else }}
        {{ template "errorMsg" (concat "Could not fetch " $mediaServer " API.") }}
      {{ end }}

    {{ else if eq $mediaServer "tautulli" }}
      {{ $sessionsRequestURL = concat $baseURL "/api/v2" }}
      {{ $sessionsCall = newRequest $sessionsRequestURL
        | withParameter "apikey" $apiKey
        | withParameter "cmd" "get_activity"
        | withHeader "Accept" "application/json"
        | getResponse }}

      {{ if eq $sessionsCall.Response.StatusCode 200 }}
        {{ $sessions = $sessionsCall.JSON.Array "response.data.sessions" }}
        {{ $activeSessions = len $sessions }}
      {{ else }}
        {{ template "errorMsg" (concat "Could not fetch " $mediaServer " API.") }}
      {{ end }}

    {{ else if or (eq $mediaServer "jellyfin") (eq $mediaServer "emby") }}
      {{ $sessionsRequestURL = concat $baseURL "/Sessions" }}
      {{ $sessionsCall = newRequest $sessionsRequestURL
        | withParameter "api_key" $apiKey
        | withParameter "activeWithinSeconds" "30"
        | withHeader "Accept" "application/json"
        | getResponse }}

      {{ if eq $sessionsCall.Response.StatusCode 200 }}
        {{ $sessions = $sessionsCall.JSON.Array "" }}
        {{ if eq $mediaServer "emby" }}
          {{ range $session := $sessions }}
            {{ if $session.Bool "PlayState.CanSeek" }}
              {{ $activeSessions = 1 }}
              {{ break }}
            {{ end }}
          {{ end }}
        {{ else }}
          {{ $activeSessions = len $sessions }}
        {{ end }}
      {{ else }}
        {{ template "errorMsg" (concat "Could not fetch " $mediaServer " API.") }}
      {{ end }}
    {{ end }}

    {{ if and (eq $sessionsCall.Response.StatusCode 200) (eq $activeSessions 0) }}
      <p>Nothing is playing right now.</p>
    {{ else }}

      <style>
        .media-server-session-container--grid {
          display: grid;
          grid-template-columns: repeat(2, 1fr);
        }
        @media (max-width: 768px) {
          .media-server-session-container--grid {
            display: flex;
            flex-direction: column;
          }
        }
        .media-server-progress-container {
          height: 1rem;
          max-width: 32rem;
          border: 1px solid var(--color-text-base);
          border-radius: var(--border-radius);
          overflow: hidden;
        }
        .media-server-progress-bar {
          height: 100%;
          background: var(--color-primary);
          border-radius: 3px;
          transition: width 1s linear;
        }
        @keyframes progress-animation { to { width: 100%; } }
      </style>

      <div class="gap-10 {{ if $isSmallColumn }}flex flex-column{{ else }}media-server-session-container--grid{{ end }}">
      {{ range $i, $session := $sessions }}
        {{ $isClient := true }}
        {{ $isPlaying := false }}
        {{ $state := "playing" }}
        {{ $isMovie := false }}
        {{ $isShows := false }}
        {{ $isMusic := false }}
        {{ $userName := "" }}
        {{ $title := "" }}
        {{ $showTitle := "" }}
        {{ $showSeason := "" }}
        {{ $showEpisode := "" }}
        {{ $artist := "" }}
        {{ $albumTitle := "" }}
        {{ $thumbURL := "" }}
        {{ $now := now | formatTime "15:04:05" | parseLocalTime "15:04:05" }}
        {{ $duration := 0 }}
        {{ $offset := 0 }}
        {{ $remainingSeconds := 0 }}

        {{ if eq $mediaServer "plex" }}
          {{ $isPlaying = eq ($session.String "Player.state") "playing" }}
          {{ if not $isPlaying }}
            {{ $state = "paused"}}
          {{ end }}

          {{ $mediaType := $session.String "type" }}
          {{ $isMovie = eq $mediaType "movie" }}
          {{ $isShows = eq $mediaType "episode" }}
          {{ $isMusic = eq $mediaType "track" }}

          {{ $userName = $session.String "User.title" }}
          {{ $title = $session.String "title" }}
          {{ $showTitle = $session.String "grandparentTitle" }}
          {{ $showSeason = $session.String "parentIndex" }}
          {{ $showEpisode = $session.String "index" }}
          {{ $artist = $session.String "grandparentTitle" }}
          {{ $albumTitle = $session.String "parentTitle" }}

          {{ $thumbID := $session.String "thumb" }}
          {{ if or $isShows $isMusic }}
            {{ $thumbID = $session.String "parentThumb" }}
          {{ end }}
          {{ $thumbURL = concat $baseURL $thumbID "?X-Plex-Token=" $apiKey }}

          {{ $duration = $session.Float "duration" }}
          {{ $offset = $session.Float "viewOffset" }}
          {{ $remainingSeconds = div (sub $duration $offset) 1000 | toInt }}
        {{ else if eq $mediaServer "tautulli" }}
          {{ $isPlaying = eq ($session.String "state") "playing" }}
          {{ if not $isPlaying }}
            {{ $state = "paused"}}
          {{ end }}

          {{ $mediaType := $session.String "media_type" }}
          {{ $isMovie = eq $mediaType "movie" }}
          {{ $isShows = eq $mediaType "episode" }}
          {{ $isMusic = eq $mediaType "track" }}

          {{ $userName = $session.String "user" }}
          {{ $title = $session.String "title" }}
          {{ $showTitle = $session.String "grandparent_title" }}
          {{ $showSeason = $session.String "parent_media_index" }}
          {{ $showEpisode = $session.String "media_index" }}
          {{ $artist = $session.String "grandparent_title" }}
          {{ $albumTitle = $session.String "parent_title" }}

          {{ $thumbID := $session.String "thumb" }}
          {{ if or $isShows $isMusic }}
            {{ $thumbID = $session.String "parent_thumb" }}
          {{ end }}
          {{ $thumbURL = concat $baseURL "/api/v2?apikey=" $apiKey "&cmd=pms_image_proxy&img=" $thumbID }}

          {{ $duration = $session.Float "duration" }}
          {{ $offset = $session.Float "view_offset" }}
          {{ $remainingSeconds = div (sub $duration $offset) 1000 | toInt }}
        {{ else if or (eq $mediaServer "jellyfin") (eq $mediaServer "emby") }}
          {{ if eq $mediaServer "emby" }}
            {{ $isClient = $session.Bool "PlayState.CanSeek" }}
          {{ end }}
          {{ $isPlaying = and ($session.Exists "NowPlayingItem") (not ($session.Bool "PlayState.IsPaused")) }}
          {{ if not $isPlaying }}
            {{ $state = "paused"}}
          {{ end }}

          {{ $mediaType := $session.String "NowPlayingItem.Type" }}
          {{ $isMovie = eq $mediaType "Movie" }}
          {{ $isShows = eq $mediaType "Episode" }}
          {{ $isMusic = eq $mediaType "Audio" }}

          {{ $userName = $session.String "UserName" }}
          {{ $title = $session.String "NowPlayingItem.Name" }}
          {{ $showTitle = $session.String "NowPlayingItem.SeriesName" }}
          {{ $showSeason = $session.String "NowPlayingItem.ParentIndexNumber" }}
          {{ $showEpisode = $session.String "NowPlayingItem.IndexNumber" }}
          {{ $artist = $session.String "NowPlayingItem.AlbumArtist" }}
          {{ $albumTitle = $session.String "NowPlayingItem.Album" }}

          {{ $thumbID := $session.String "NowPlayingItem.Id" }}
          {{ if $isShows }}
            {{ $thumbID = $session.String "NowPlayingItem.ParentId" }}
          {{ end }}
          {{ $thumbURL = concat $baseURL "/Items/" $thumbID "/Images/Primary?api_key=" $apiKey }}

          {{ $duration = $session.Float "NowPlayingItem.RunTimeTicks" }}
          {{ $offset = $session.Float "PlayState.PositionTicks" }}
          {{ $remainingSeconds = div (sub $duration $offset) 10000000 | toInt }}
        {{ end }}

        {{ $progress := mul 100 (div $offset $duration) | toInt }}
        {{ $endTime := $now.Add (duration (printf "%ds" $remainingSeconds)) | formatTime $timeFormat }}

        {{ $showInfoFormat := concat "Season " $showSeason " Episode " $showEpisode}}
        {{ if $isCompact }}
          {{ $showInfoFormat = concat "S" $showSeason "E" $showEpisode}}
        {{ end }}

        {{ if and $isClient (or $isPlaying $showPaused) }}
        <div class="card gap-5">
          <div class="flex items-center gap-10 size-h3">
            <span class="color-primary">{{ $userName }}</span>
            {{ if eq $playState "text" }}
              <span {{ if $isPlaying }}class="color-primary"{{ end }}>
                [{{ $state }}]
              </span>
            {{ else if eq $playState "indicator" }}
              <style>
                .media-server-indicator {
                  height: .7rem;
                  width: .7rem;
                  border-radius: 100%;
                  {{ if $isPlaying }}
                    animation: pulse 5s infinite;
                    background: var(--color-primary);
                  {{ else }}
                    background: var(--color-text-base-muted);
                  {{ end }}
                }
                @keyframes pulse {
                  0% { box-shadow: 0 0 0 0 var(--color-text-base); }
                  40% { box-shadow: 0 0 0 4px transparent; }
                  100% { box-shadow: 0 0 0 4px transparent; }
                }
              </style>
              <div class="media-server-indicator"></div>
            {{ end }}
          </div>

          <hr class="margin-bottom-5" />

          <div class="flex items-center gap-10" style="align-items: stretch;">
            {{ if $showThumbnail }}
              <img src="{{ $thumbURL | safeURL }}"
                alt="{{ $title }} thumbnail"
                class="shrink-0"
                loading="lazy"
                style="
                  max-width: 7.5rem;
                  border: 2px solid var(--color-primary);
                  border-radius: var(--border-radius);
                  object-fit: cover;
                  {{ if and $isCompact (not $fullThumbnail) }} aspect-ratio: 1; {{ end }} "
              />
            {{ end }}
            <ul class="flex flex-column grow justify-evenly" style="width: 0;">
              {{ if $isCompact }}
                {{ if $isShows }}
                  <ul class="list-horizontal-text flex-nowrap">
                    <li class="shrink-0">{{ concat "S" $showSeason "E" $showEpisode }}</li>
                    <li class="text-truncate">{{ $showTitle }}</li>
                  </ul>
                {{ else if $isMusic }}
                  <ul class="list-horizontal-text flex-nowrap">
                    <li class="shrink-0">{{ $artist }}</li>
                    <li class="text-truncate">{{ $albumTitle }}</li>
                  </ul>
                {{ end }}
                <li class="text-truncate">{{ $title }}</li>
              {{ else }}
                {{ if $isShows }}
                  <li>{{ $showTitle }}</li>
                  <li>{{ concat "S" $showSeason "E" $showEpisode }}</li>
                {{ else if $isMusic }}
                  <li>{{ $artist }}</li>
                  <li>{{ $albumTitle }}</li>
                {{ end }}
                <li>{{ $title }}</li>
              {{ end }}

              <li>
                {{ if and $isPlaying $showProgressBar }}
                  <div class="flex gap-10 items-center">
                    <div class="media-server-progress-container grow">
                      <div
                        class ="media-server-progress-bar"
                        data-progress="{{ $progress }}"
                        data-remaining="{{ $remainingSeconds }}"
                        style="
                          width: {{ $progress }}%;
                          animation: progress-animation {{ $remainingSeconds }}s linear forwards;"
                      >
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
      {{ end }}
    </div>
    {{ end }}
    {{ end }}
```

## 🍻 Cheers
* [svilenmarkov](https://github.com/svilenmarkov)
* [erkston](https://github.com/erkston)
* [ralphocdol](https://github.com/ralphocdol)

