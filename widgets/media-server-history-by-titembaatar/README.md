* [Introduction](#introduction)
* [Preview](#preview)
    * [Full Size Column](#full-size-column)
    * [Small Size Column](#small-size-column)
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
### Full Size Column
![Preview](preview.png)

### Small Size Column
![Preview Small](preview-small.png)

### Compact Mode
![Preview Compact](preview-compact.png)

## Environment Variables
### Plex
> [!IMPORTANT]
>
> For URLs, you **MUST** include `http://` or `https://`.
> Do **NOT** include a trailing `/` at the end of URLs.

* `PLEX_URL` - the Plex URL, can be `http://<ip_address>:<port>` or `https://<domain>`
* `PLEX_TOKEN` - the Plex token, follow [this guide](https://support.plex.tv/articles/204059436-finding-an-authentication-token-x-plex-token/) if you don't know how to get it

### Tautulli
* `TAUTULLI_URL` - the Tautulli URL, can be `http://<ip_address>:<port>` or `https://<domain>`
* `TAUTULLI_KEY` - the Tautulli API key, can be found in `Settings` -> `Web Interface `-> `API key`

### Jellyfin
* `JELLYFIN_URL` - the Jellyfin URL, can be `http://<ip_address>:<port>` or `https://<domain>`
* `JELLYFIN_KEY` - the Jellyfin API key, use or create one in `Administration` -> `Dashboard` -> `API Keys`

## Secrets
Since `v0.8.0`, you can use Docker secrets instead of environment variables.  
If you do, replace `${YOUR_API_KEY}` with `${secret:your-api-key-secret}` throughout the code.

## Options
Since `v0.8.0`, you can use the `options:` field to customise the widget.

> [!CAUTION]
>
> Displaying thumbnails **WILL** expose your token/API keys in the HTML.
> Do **NOT** enable this option if you are using Glance in production or exposing the service to the internet.

Default options are:
```yaml
options:
  small-column: false        # `true` if using the widget in a small column
  compact: true              # `false` for a more spread-out layout
  show-thumbnail: false      # `true` to show thumbnails
  thumbnail-aspect-ratio: "" # Options: "square", "portrait", "landscape", ""
  show-user: true            # `false` to hide username
  time-absolute: false       # `true` for absolute time format

  # Jellyfin-specific options
  user-name: "titem"                # Username for retrieving history
  media-type: "Movie,Episode,Audio" # Media types to display, capitalized and comma-separated
```

> [!IMPORTANT]
>
> The Jellyfin API only retrieves playback history for a specific user.
> You must set the `user-name` option to your Jellyfin username.

On top of options, there is some paramaters you can change:
```yaml
parameters:
  # Plex
  limit: 20                        # Modify this value for the length of the history
  # Tautulli
  length: 10                       # Modify this value for the length of the history
  media_type: movie,episode,track  # Media types: movie, episode, track, live, collection, playlist
```

## Widget YAML
### Plex YAML
```yaml
- type: custom-api
  title: plex history
  frameless: true
  cache: 5m
  url: ${PLEX_URL}/status/sessions/history/all
  parameters:
    limit: 10 # Modify this value for the length of the history
    sort: viewedAt:desc
  headers:
    Accept: application/json
    X-Plex-Token: ${PLEX_TOKEN}
  template: |
    {{ $isSmallColumn := .Options.BoolOr "small-column" false }}
    {{ $isCompact := .Options.BoolOr "compact" true }}
    {{ $showThumbnail := .Options.BoolOr "show-thumbnail" false }}
    {{ $thumbAspectRatio := .Options.StringOr "thumbnail-aspect-ratio" "" }}
    {{ $showUser := .Options.BoolOr "show-user" true }}
    {{ $timeAbsolute := .Options.BoolOr "time-absolute" false }}

    {{ $usersCall := newRequest ("${PLEX_URL}/accounts")
      | withHeader "Accept" "application/json"
      | withHeader "X-Plex-Token" "${PLEX_TOKEN}"
      | getResponse
    }}
    {{ $users := $usersCall.JSON.Array "MediaContainer.Account" }}

    {{ if eq .Response.StatusCode 200 }}
      {{ $history := .JSON.Array "MediaContainer.Metadata" }}

      {{ if eq (len $history) 0 }}
        <p>Nothing is playing. Start watching something!</p>
      {{ else }}
        <div class="carousel-container show-right-cutoff">
          <div class="cards-horizontal carousel-items-container">
            {{ range $n, $item := $history }}
              {{/* PLEX VARIABLES BEGIN */}}
                {{ $userName := "" }}
                {{ $userID := $item.Int "accountID" }}

                {{ range $n, $u := $users }}
                  {{ if eq $userID ($u.Int "id") }}
                    {{ $userName = $u.String "name" }}
                    {{ break }}
                  {{ end }}
                {{ end }}
              {{/* PLEX VARIABLES END */}}
              {{/* WIDGET VARIABLES BEGIN */}}
                {{ $mediaType := $item.String "type" }}
                {{ $isMovie := eq $mediaType "movie" }}
                {{ $isShows := eq $mediaType "episode" }}
                {{ $isMusic := eq $mediaType "track" }}

                {{ $movieTitle := $item.String "title" }}
                {{ $showTitle := $item.String "grandparentTitle" }}
                {{ $showSeason := $item.String "parentIndex" }}
                {{ $showEpisode := $item.String "index" }}
                {{ $episodeTitle := $item.String "title" }}
                {{ $artist := $item.String "grandparentTitle" }}
                {{ $albumTitle := $item.String "parentTitle" }}
                {{ $songTitle := $item.String "title" }}
                {{ $default := $item.String "title" }}

                {{ $thumbPath := $item.String "thumb" }}
                {{ if or $isShows $isMusic }}
                  {{ $thumbPath = $item.String "parentThumb" }}
                {{ end }}
                {{ $thumbURL := concat "${PLEX_URL}" $thumbPath "?X-Plex-Token=${PLEX_TOKEN}" }}

                {{ $playedAt := $item.String "viewedAt" | parseRelativeTime "unix" }}
                {{ if $timeAbsolute }}
                  {{ $t := $item.String "viewedAt" | parseLocalTime "unix" }}
                  {{ $playedAt = $t.Format "Jan 02 15:04" }}
                {{ end }}
              {{/* WIDGET VARIABLES END */}}
              {{/* WIDGET TEMPLATE BEGIN */}}
              <div class="card widget-content-frame">
                {{ if $showThumbnail }}
                  <img src="{{ $thumbURL }}"
                    alt="{{ $default }} thumbnail"
                    loading="lazy"
                    class="shrink-0"
                    style="
                      object-fit: cover;
                      {{ if eq $thumbAspectRatio "square" }} aspect-ratio: 1;
                      {{ else if eq $thumbAspectRatio "portrait" }} aspect-ratio: 3/4;
                      {{ else if eq $thumbAspectRatio "landscape" }} aspect-ratio: 4/3;
                      {{ else }} aspect-ratio: initial;
                      {{ end }}
                      border-radius: var(--border-radius) var(--border-radius) 0 0;"
                  />
                {{ end }}

                <div class="grow padding-inline-widget margin-top-10 margin-bottom-10">
                  <ul class="flex flex-column justify-evenly margin-bottom-3
                      {{ if $isSmallColumn }}size-h6{{ end }}"
                    style="height: 100%;">

                    {{ if $isCompact }}
                      <ul class="list-horizontal-text flex-nowrap">
                    {{ end }}
                      {{ if $showUser }}
                        <li class="color-primary text-truncate">{{ $userName }}</li>
                      {{ end }}

                      {{ if $timeAbsolute }}
                        <li>{{ $playedAt }}</li>
                      {{ else }}
                        <li class="shrink-0">
                          <span {{ $playedAt }}></span>
                          {{ if or (not $isCompact) (not $showUser) }}
                            <span> ago</span>
                          {{ end }}
                        </li>
                      {{ end }}
                    {{ if $isCompact }}
                      </ul>
                    {{ end }}

                    {{ if $isMovie }}
                      <li {{ if $isCompact }}class="text-truncate"{{ end }}>{{ $movieTitle }}</li>
                    {{ else if $isShows }}
                      {{ if $isCompact }}
                        <ul class="list-horizontal-text flex-nowrap">
                      {{ end }}
                          {{ if or $isSmallColumn $isCompact }}
                            <li>{{ concat "S" $showSeason "E" $showEpisode }}</li>
                          {{ else }}
                            <li class="text-truncate">{{ concat "Season " $showSeason " Episode " $showEpisode }}</li>
                          {{ end }}
                          <li class="text-truncate">{{ $showTitle }}</li>
                      {{ if $isCompact }}
                        </ul>
                      {{ end }}
                      <li class="text-truncate">{{ $episodeTitle }}</li>
                    {{ else if $isMusic }}
                      <li class="text-truncate">{{ $artist }}</li>
                      {{ if not $isCompact }}
                        <li class="text-truncate">{{ $albumTitle }}</li>
                      {{ end }}
                      <li class="text-truncate">{{ $songTitle }}</li>
                    {{ else }}
                      <li class="text-truncate">{{ $default }}</li>
                    {{ end }}
                  </ul>
                </div>
              </div>
              {{/* WIDGET TEMPLATE END */}}
            {{ end }}
          </div>
        </div>
      {{ end }}
    {{ else }}
      <p>Failed to fetch Plex history</p>
    {{ end }}
```

### Tautulli YAML
```yaml
- type: custom-api
  title: tautulli history
  frameless: true
  cache: 5m
  url: ${TAUTULLI_URL}/api/v2
  parameters:
    length: 10                       # Modify this value for the length of the history
    media_type: movie,episode,track  # movie,episode,track,live,collection,playlist
    apikey: ${TAUTULLI_KEY}
    cmd: get_history
  template: |
    {{ $isSmallColumn := .Options.BoolOr "small-column" false }}
    {{ $isCompact := .Options.BoolOr "compact" true }}
    {{ $showThumbnail := .Options.BoolOr "show-thumbnail" false }}
    {{ $thumbAspectRatio := .Options.StringOr "thumbnail-aspect-ratio" "" }}
    {{ $showUser := .Options.BoolOr "show-user" true }}
    {{ $timeAbsolute := .Options.BoolOr "time-absolute" false }}

    {{ if eq .Response.StatusCode 200 }}
      {{ $history := .JSON.Array "response.data.data" }}

      {{ if eq (len $history) 0 }}
        <div class="card widget-content-frame padding-widget">
          <p>Nothing is playing. Start watching something!</p>
        </div>
      {{ else }}
        <div class="carousel-container show-right-cutoff">
          <div class="cards-horizontal carousel-items-container">
            {{ range $n, $item := $history }}
              {{/* WIDGET VARIABLES BEGIN */}}
                {{ $userName := $item.String "user" }}
                {{ $mediaType := $item.String "media_type" }}
                {{ $isMovie := eq $mediaType "movie" }}
                {{ $isShows := eq $mediaType "episode" }}
                {{ $isMusic := eq $mediaType "track" }}

                {{ $movieTitle := $item.String "title" }}
                {{ $showTitle := $item.String "grandparent_title" }}
                {{ $showSeason := $item.String "parent_media_index" }}
                {{ $showEpisode := $item.String "media_index" }}
                {{ $episodeTitle := $item.String "title" }}
                {{ $artist := $item.String "grandparent_title" }}
                {{ $albumTitle := $item.String "parent_title" }}
                {{ $songTitle := $item.String "title" }}
                {{ $default := $item.String "title" }}

                {{ $thumbPath := $item.String "thumb" }}
                {{ $thumbURL := concat "${TAUTULLI_URL}/api/v2?apikey=${TAUTULLI_KEY}&cmd=pms_image_proxy&img=" $thumbPath }}

                {{ $playedAt := $item.String "date" | parseRelativeTime "unix" }}
                {{ if $timeAbsolute }}
                  {{ $t := $item.String "date" | parseLocalTime "unix" }}
                  {{ $playedAt = $t.Format "Jan02 15:04" }}
                {{ end }}
              {{/* WIDGET VARIABLES END */}}
              {{/* WIDGET TEMPLATE BEGIN */}}
              <div class="card widget-content-frame">
                {{ if $showThumbnail }}
                  <img src="{{ $thumbURL }}"
                    alt="{{ $default }} thumbnail"
                    loading="lazy"
                    class="shrink-0"
                    style="
                      object-fit: cover;
                      {{ if eq $thumbAspectRatio "square" }} aspect-ratio: 1;
                      {{ else if eq $thumbAspectRatio "portrait" }} aspect-ratio: 3/4;
                      {{ else if eq $thumbAspectRatio "landscape" }} aspect-ratio: 4/3;
                      {{ else }} aspect-ratio: initial;
                      {{ end }}
                      border-radius: var(--border-radius) var(--border-radius) 0 0;"
                  />
                {{ end }}

                <div class="grow padding-inline-widget margin-top-10 margin-bottom-10">
                  <ul class="flex flex-column justify-evenly margin-bottom-3
                      {{ if $isSmallColumn }}size-h6{{ end }}"
                    style="height: 100%;">

                    {{ if $isCompact }}
                      <ul class="list-horizontal-text flex-nowrap">
                    {{ end }}
                      {{ if $showUser }}
                        <li class="color-primary text-truncate">{{ $userName }}</li>
                      {{ end }}

                      {{ if $timeAbsolute }}
                        <li>{{ $playedAt }}</li>
                      {{ else }}
                        <li class="shrink-0">
                          <span {{ $playedAt }}></span>
                          {{ if or (not $isCompact) (not $showUser) }}
                            <span> ago</span>
                          {{ end }}
                        </li>
                      {{ end }}
                    {{ if $isCompact }}
                      </ul>
                    {{ end }}

                    {{ if $isMovie }}
                      <li {{ if $isCompact }}class="text-truncate"{{ end }}>{{ $movieTitle }}</li>
                    {{ else if $isShows }}
                      {{ if $isCompact }}
                        <ul class="list-horizontal-text flex-nowrap">
                      {{ end }}
                          {{ if or $isSmallColumn $isCompact }}
                            <li>{{ concat "S" $showSeason "E" $showEpisode }}</li>
                          {{ else }}
                            <li class="text-truncate">{{ concat "Season " $showSeason " Episode " $showEpisode }}</li>
                          {{ end }}
                          <li class="text-truncate">{{ $showTitle }}</li>
                      {{ if $isCompact }}
                        </ul>
                      {{ end }}
                      <li class="text-truncate">{{ $episodeTitle }}</li>
                    {{ else if $isMusic }}
                      <li class="text-truncate">{{ $artist }}</li>
                      {{ if not $isCompact }}
                        <li class="text-truncate">{{ $albumTitle }}</li>
                      {{ end }}
                      <li class="text-truncate">{{ $songTitle }}</li>
                    {{ else }}
                      <li class="text-truncate">{{ $default }}</li>
                    {{ end }}
                  </ul>
                </div>
              </div>
              {{/* WIDGET TEMPLATE END */}}
            {{ end }}
          </div>
        </div>
      {{ end }}
    {{ else }}
      <p>Failed to fetch Tautulli history</p>
    {{ end }}
```

### Jellyfin YAML
```yaml
- type: custom-api
  title: jellyfin history
  frameless: true
  cache: 5m
  url: ${JELLYFIN_URL}/Users
  options:
    user-name: "titem" # You need to change your username
  parameters:
    api_key: ${JELLYFIN_KEY}
  template: |
    {{ $userName := .Options.StringOr "user-name" "" }}
    {{ $mediaType := .Options.StringOr "media-type" "Movie,Episode" }}
    {{ $isSmallColumn := .Options.BoolOr "small-column" false }}
    {{ $isCompact := .Options.BoolOr "compact" true }}
    {{ $showThumbnail := .Options.BoolOr "show-thumbnail" false }}
    {{ $thumbAspectRatio := .Options.StringOr "thumbnail-aspect-ratio" "" }}
    {{ $showUser := .Options.BoolOr "show-user" true }}
    {{ $timeAbsolute := .Options.BoolOr "time-absolute" false }}

    {{ if eq .Response.StatusCode 200 }}
      {{ $userId := "" }}
      {{ $users := .JSON.Array "" }}

      {{ range $i, $user := $users }}
        {{ if eq ($user.String "Name") $userName }}
          {{ $userId = $user.String "Id" }}
          {{ break }}
        {{ end }}
      {{ end }}

      {{ $historyCall := newRequest (concat "${JELLYFIN_URL}/Users/" $userId "/Items")
        | withParameter "api_key" "${JELLYFIN_KEY}"
        | withParameter "Limit" "20"
        | withParameter "IncludeItemTypes" "$mediaType"
        | withParameter "Recursive" "true"
        | withParameter "isPlayed" "true"
        | withParameter "sortBy" "DatePlayed"
        | withParameter "sortOrder" "Descending"
        | withHeader "Accept" "application/json"
        | getResponse
      }}

      {{ $history := $historyCall.JSON.Array "Items" }}

      {{ if eq (len $history) 0 }}
        <p>Nothing is playing. Start watching something!</p>
      {{ else }}
        <div class="carousel-container show-right-cutoff">
          <div class="cards-horizontal carousel-items-container">
            {{ range $n, $item := $history }}
              {{/* WIDGET VARIABLES BEGIN */}}
                {{ $mediaType := $item.String "Type" }}
                {{ $isMovie := eq $mediaType "Movie" }}
                {{ $isShows := eq $mediaType "Episode" }}
                {{ $isMusic := eq $mediaType "Audio" }}

                {{ $movieTitle := $item.String "Name" }}
                {{ $showTitle := $item.String "SeriesName" }}
                {{ $showSeason := $item.String "ParentIndexNumber" }}
                {{ $showEpisode := $item.String "IndexNumber" }}
                {{ $episodeTitle := $item.String "Name" }}
                {{ $artist := $item.String "AlbumArtist" }}
                {{ $albumTitle := $item.String "Album" }}
                {{ $songTitle := $item.String "Name" }}
                {{ $default := $item.String "Name" }}

                {{ $thumbID := $item.String "Id" }}
                {{ if $isShows }}
                  {{ $thumbID = $item.String "SeasonId" }}
                {{ end }}
                {{ $thumbURL := concat "${JELLYFIN_URL}/Items/" $thumbID "/Images/Primary?api_key=${JELLYFIN_KEY}" }}

                {{ $playedAt := $item.String "UserData.LastPlayedDate" | parseRelativeTime "rfc3339" }}
                {{ if $timeAbsolute }}
                  {{ $t := $item.String "UserData.LastPlayedDate" | parseLocalTime "rfc3339" }}
                  {{ $playedAt = $t.Format "Jan 02 15:04" }}
                {{ end }}
              {{/* WIDGET VARIABLES END */}}
              {{/* WIDGET TEMPLATE BEGIN */}}
              <div class="card widget-content-frame">
                {{ if $showThumbnail }}
                  <img src="{{ $thumbURL }}"
                    alt="{{ $default }} thumbnail"
                    loading="lazy"
                    class="shrink-0"
                    style="
                      object-fit: cover;
                      {{ if eq $thumbAspectRatio "square" }} aspect-ratio: 1;
                      {{ else if eq $thumbAspectRatio "portrait" }} aspect-ratio: 3/4;
                      {{ else if eq $thumbAspectRatio "landscape" }} aspect-ratio: 4/3;
                      {{ else }} aspect-ratio: initial;
                      {{ end }}
                      border-radius: var(--border-radius) var(--border-radius) 0 0;"
                  />
                {{ end }}

                <div class="grow padding-inline-widget margin-top-10 margin-bottom-10">
                  <ul class="flex flex-column justify-evenly margin-bottom-3
                      {{ if $isSmallColumn }}size-h6{{ end }}"
                    style="height: 100%;">

                    {{ if $isCompact }}
                      <ul class="list-horizontal-text flex-nowrap">
                    {{ end }}
                      {{ if $showUser }}
                        <li class="color-primary text-truncate">{{ $userName }}</li>
                      {{ end }}

                      {{ if $timeAbsolute }}
                        <li>{{ $playedAt }}</li>
                      {{ else }}
                        <li class="shrink-0">
                          <span {{ $playedAt }}></span>
                          {{ if or (not $isCompact) (not $showUser) }}
                            <span> ago</span>
                          {{ end }}
                        </li>
                      {{ end }}
                    {{ if $isCompact }}
                      </ul>
                    {{ end }}

                    {{ if $isMovie }}
                      <li {{ if $isCompact }}class="text-truncate"{{ end }}>{{ $movieTitle }}</li>
                    {{ else if $isShows }}
                      {{ if $isCompact }}
                        <ul class="list-horizontal-text flex-nowrap">
                      {{ end }}
                          {{ if or $isSmallColumn $isCompact }}
                            <li>{{ concat "S" $showSeason "E" $showEpisode }}</li>
                          {{ else }}
                            <li class="text-truncate">{{ concat "Season " $showSeason " Episode " $showEpisode }}</li>
                          {{ end }}
                          <li class="text-truncate">{{ $showTitle }}</li>
                      {{ if $isCompact }}
                        </ul>
                      {{ end }}
                      <li class="text-truncate">{{ $episodeTitle }}</li>
                    {{ else if $isMusic }}
                      <li class="text-truncate">{{ $artist }}</li>
                      {{ if not $isCompact }}
                        <li class="text-truncate">{{ $albumTitle }}</li>
                      {{ end }}
                      <li class="text-truncate">{{ $songTitle }}</li>
                    {{ else }}
                      <li class="text-truncate">{{ $default }}</li>
                    {{ end }}
                  </ul>
                </div>
              </div>
              {{/* WIDGET TEMPLATE END */}}
            {{ end }}
          </div>
        </div>
      {{ end }}
    {{ else }}
      <p>Failed to fetch Jellyfin history</p>
    {{ end }}
```
