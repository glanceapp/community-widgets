## *arr Releases Widget

A widget for Sonarr, Radarr, or Lidarr that shows upcoming releases, recent downloads, or missing items.

<img src=sonarrrecent.png width="700">

<img src=radarrmissing.png width="320"> <img src=lidarrupcoming.png width="385">

### Widget yaml

<details>

<summary>Click</summary>

```yaml
- type: custom-api
  title: Media Releases
  title-url: ${RADARR_URL}
  cache: 30m
  options:
    service: radarr          # sonarr, radarr, or lidarr
    type: upcoming           # upcoming, recent, or missing
    size: medium             # small, medium, large, or huge
    collapse-after: 3
    show-grabbed: false
    timezone: "-04"          # must have quotes
    interval: 20             # optional, days
    #sort-time: asc          # optional, asc or desc (default)
    #cover-proxy: "https://proxy.example.com/radarrcover" # optional
    api-base-url: ${RADARR_API_URL}
    key: ${RADARR_KEY}
    url: ${RADARR_URL}       # optional
  template: |
    {{ $collapseAfter := .Options.IntOr "collapse-after" 5 }}
    {{ $showGrabbed := .Options.BoolOr "show-grabbed" false }}
    {{ $tzOffset := .Options.StringOr "timezone" "+00" }} 
    {{ $tzTime := (printf "2006-01-02T15:04:05%s:00" $tzOffset) | parseTime "rfc3339" }}
    {{ $timezone := $tzTime.Location }}
    {{ $sortTime := .Options.StringOr "sort-time" "desc" }}
    {{ $coverProxy := .Options.StringOr "cover-proxy" "" }}
    {{ $size := .Options.StringOr "size" "medium" }}
    {{ $service := .Options.StringOr "service" "" }}
    {{ $type := .Options.StringOr "type" "" }}
    {{ $intervalH := .Options.IntOr "interval" 0 | mul 24 }}

    {{ $nowUTC := offsetNow "0h" }}
    {{ $now := ($nowUTC.In $timezone) | formatTime "rfc3339" }}
    {{ $posInterval := ((offsetNow (printf "+%dh" $intervalH)).In $timezone) | formatTime "rfc3339" }}
    {{ $negInterval := ((offsetNow (printf "-%dh" $intervalH)).In $timezone) | formatTime "rfc3339" }}

    {{ $apiBaseUrl := .Options.StringOr "api-base-url" "" }}
    {{ $key := .Options.StringOr "key" "" }}
    {{ $url := .Options.StringOr "url" $apiBaseUrl }}

    {{ if or (and (ne $service "sonarr") (ne $service "radarr") (ne $service "lidarr"))
             (and (ne $type "upcoming") (ne $type "recent") (ne $type "missing")) 
             (eq $apiBaseUrl "") (eq $key "") (eq $url "") }}
      <div class="widget-error-header">
          <div class="color-negative size-h3">ERROR</div>
          <svg class="widget-error-icon" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126ZM12 15.75h.007v.008H12v-.008Z"></path>
          </svg>
        </div>
        <p class="break-all">
          Some options are not set or malformed
            <table style="border-spacing: 1rem;">
              <tr><td>service</td><td>{{ $service }}</td><td>must be sonarr, radarr, or lidarr</td></tr>
              <tr><td>type</td><td>{{ $type }}</td><td>must be upcoming, recent, or missing</td></tr>
              <tr><td>api-base-url</td><td>{{ $apiBaseUrl }}</td><td>should include http(s):// and port if needed</td></tr>
              <tr><td>key</td><td>{{ $key }}</td><td></td></tr>
            </table> 
        </p>
    {{ else }}
      {{ $requestUrl := "" }}
  
      {{ if eq $service "sonarr" }}
        {{ if eq $type "recent" }}
          {{ $requestUrl = printf "%s/api/v3/history/since?includeSeries=true&includeEpisode=true&eventType=grabbed&date=%s" $apiBaseUrl $negInterval }}
        {{ else if eq $type "missing" }}
          {{ $requestUrl = printf "%s/api/v3/wanted/missing?page=1&pageSize=50&includeSeries=true" $apiBaseUrl }}
        {{ else if eq $type "upcoming" }}
          {{ $requestUrl = printf "%s/api/v3/calendar?includeSeries=true&start=%s&end=%s" $apiBaseUrl $now $posInterval }}
        {{ end }}
        
      {{ else if eq $service "radarr" }}
        {{ if eq $type "recent" }}
          {{ $requestUrl = printf "%s/api/v3/history/since?includeMovie=true&eventType=grabbed&date=%s" $apiBaseUrl $negInterval }}
        {{ else if eq $type "missing" }}
          {{ $requestUrl = printf "%s/api/v3/wanted/missing?page=1&pageSize=50" $apiBaseUrl }}
        {{ else if eq $type "upcoming" }}
          {{ $requestUrl = printf "%s/api/v3/calendar?start=%s&end=%s" $apiBaseUrl $now $posInterval }}
        {{ end }}
        
      {{ else if eq $service "lidarr" }}
        {{ if eq $type "recent" }}
          {{ $requestUrl = printf "%s/api/v1/history/since?includeArtist=true&includeAlbum=true&eventType=grabbed&date=%s" $apiBaseUrl $negInterval }}
        {{ else if eq $type "missing" }}
          {{ $requestUrl = printf "%s/api/v1/wanted/missing?page=1&pageSize=50&includeArtist=true" $apiBaseUrl }}
        {{ else if eq $type "upcoming" }}
          {{ $requestUrl = printf "%s/api/v1/calendar?includeArtist=true&start=%s&end=%s" $apiBaseUrl $now $posInterval }}
        {{ end }}
      {{ end }}
  
      {{ $data := newRequest $requestUrl
        | withHeader "Accept" "application/json"
        | withHeader "X-Api-Key" $key
        | getResponse }}
  
      <ul class="list list-gap-14 collapsible-container single-line-titles" data-collapse-after="{{ $collapseAfter }}">
        {{ $array := "" }}
        {{ $sortByField := "" }}
        {{ $itemDisplayed := false }}
        {{ $itemDate := "" }}
        {{ $isAvailable := false }}
  
        {{ if eq $type "missing" }}
          {{ $array = "records" }}
          {{ if eq $service "sonarr" }}
            {{ $sortByField = "airDateUtc" }}
          {{ else }}
            {{ $sortByField = "releaseDate" }}
          {{ end }}
        {{ else if eq $type "recent" }}
          {{ $sortByField = "date" }}
        {{ else }}
          {{ if eq $service "sonarr" }}
            {{ $sortByField = "airDateUtc" }}
          {{ else }}
            {{ $sortByField = "releaseDate" }}
          {{ end }}
        {{ end }}
  
        {{ range $data.JSON.Array $array | sortByTime $sortByField "rfc3339" $sortTime }}
           
          {{ if eq $service "sonarr" }}
            {{ $itemDate = .String "airDateUtc" }}
            {{ $isAvailable = true }}
          {{ else if eq $service "radarr"}}
            {{ $isAvailable = .Bool "isAvailable" }}
            {{ $itemDate = .String "releaseDate" }}
          {{ else if eq $service "lidarr"}}
            {{ $itemDate = .String "releaseDate" }}
            {{ $isAvailable = true }}
          {{ end }}
  
          {{ if or (eq $type "upcoming") (eq $type "recent") (and (or (and (gt $itemDate $negInterval) ((lt $itemDate $now ))) (eq $intervalH 0))  $isAvailable) }}
            {{ $itemDisplayed = true }}
            {{ $title := "" }}
            {{ $subtitle := "" }}
            {{ $coverUrl := "" }}
            {{ $status := "" }}
            {{ $coverBase := "" }}
            {{ $height := "" }}
            {{ $width := "" }}
            {{ $popoverTitle := "" }}
            {{ $popoverSubtitle := "" }}
            {{ $popoverSummary := "" }}
            {{ $summary := "" }}
            {{ $link := "" }}
            {{ $grabbed := false }}
            {{ $date := now }}
            {{ $datetype := "" }}
            {{ $seString := "" }}
            {{ $genres := "" }}
            {{ $buttonJustify := "left" }}
  
    
            {{ if eq $service "sonarr" }}
              {{ $series := "" }}
              {{ if eq $coverProxy "" }}
                {{ $coverBase = printf "%s/api/v3/mediacover" $apiBaseUrl }}
                {{ $coverUrl = printf "%s/%s/poster-500.jpg?apikey=%s" $coverBase (.String "seriesId") $key }}
              {{ else }}
                {{ $coverBase = $coverProxy }}
                {{ $coverUrl = printf "%s/%s/poster-500.jpg" $coverBase (.String "seriesId") }}
              {{ end }}
              {{ $title = .String "series.title" }}
              {{ $link = printf "%s/series/%s#" $url (.String "series.titleSlug") }}
              {{ $series = .Get "series" }}
              {{ $genres = $series.Get "genres" }}
             
              {{ if eq $type "recent" }}
                {{ $date = (.String "date" | parseTime "rfc3339") }}
                {{ $subtitle = .String "episode.title" }}
                {{ $summary = .String "episode.overview" }}
                {{ $seString = printf "S%02dE%02d" (.Int "episode.seasonNumber") (.Int "episode.episodeNumber") }}
                {{ $datetype = "Downloaded" }}
        
                {{ if $showGrabbed }}
                  {{ $popoverTitle = .String "episode.title" }}
                  {{ $popoverSubtitle = $seString }}
                  {{ $popoverSummary = .String "episode.overview" }}
                  {{ if .Bool "episode.hasFile" }}
                    {{ $grabbed = true }}
                  {{ end }}
                {{ else }}
                  {{ $popoverTitle = .String "series.title" }}
                  {{ $popoverSummary = .String "series.overview" }}
                {{ end }}
    
              {{ else if eq $type "missing" }}
                {{ $date = (.String "airDateUtc" | parseTime "rfc3339") }}
                {{ $subtitle = .String "title" }}
                {{ $summary = .String "overview" }}
                {{ $seString = printf "S%02dE%02d" (.Int "seasonNumber") (.Int "episodeNumber") }}
                {{ $datetype = "Aired" }}
                {{ if $showGrabbed }}
                  {{ $popoverTitle = .String "title" }}
                  {{ $popoverSubtitle = $seString }}
                  {{ $popoverSummary = .String "overview" }}
                  {{ if .Bool "episode.hasFile" }}
                    {{ $grabbed = true }}
                  {{ end }}
                {{ else }}
                  {{ $popoverTitle = .String "series.title" }}
                  {{ $popoverSummary = .String "series.overview" }}
                {{ end }}
              {{ else if eq $type "upcoming" }}
                {{ $date = (.String "airDateUtc" | parseTime "rfc3339") }}
                {{ $subtitle = .String "title" }}
                {{ $summary = .String "overview" }}
                {{ $seString = printf "S%02dE%02d" (.Int "seasonNumber") (.Int "episodeNumber") }}
                {{ $datetype = "Airs" }}
                {{ if $showGrabbed }}
                  {{ $popoverTitle = .String "title" }}
                  {{ $popoverSubtitle = $seString }}
                  {{ $popoverSummary = .String "overview" }}
                  {{ if .Bool "hasFile" }}
                    {{ $grabbed = true }}
                  {{ end }}
                {{ else }}
                  {{ $popoverTitle = .String "series.title" }}
                  {{ $popoverSummary = .String "series.overview" }}
                {{ end }}
              {{ end }}
        
        
            {{ else if eq $service "radarr" }}
              {{ $movie := "" }}
              {{ $status = .String "status" }}
              {{ if eq $coverProxy "" }}
                {{ $coverBase = printf "%s/api/v3/mediacover" $apiBaseUrl }}
              {{ else }}
                {{ $coverBase = $coverProxy }}
              {{ end }}
              {{ if eq $status "announced"}}
                {{ if ne (.String "inCinemas") "" }}
                  {{ $date = (.String "inCinemas" | parseTime "rfc3339") }}
                  {{ $datetype = "In Cinemas" }}
                {{ else }}
                  {{ $date = (.String "releaseDate" | parseTime "rfc3339") }}
                  {{ $datetype = "Releases" }}
                {{ end }}
              {{ else if eq $status "inCinemas"}}
                {{ $date = (.String "releaseDate" | parseTime "rfc3339") }}
                {{ $datetype = "Releases" }}
              {{ else if eq $status "released"}}
                {{ $date = (.String "releaseDate" | parseTime "rfc3339") }}
                {{ $datetype = "Released" }}
              {{ end }}
   
              {{ if eq $type "recent" }}
                {{ if eq $coverProxy "" }}
                  {{ $coverUrl = printf "%s/%s/poster-500.jpg?apikey=%s" $coverBase (.String "movie.id") $key }}
                {{ else }}
                  {{ $coverUrl = printf "%s/%s/poster-500.jpg" $coverBase (.String "movie.id") }}
                {{ end }}
                {{ $datetype = "Downloaded" }}
                {{ $date = (.String "date" | parseTime "rfc3339") }}
                {{ $title = .String "movie.title" }}
                {{ $subtitle = "" }}
                {{ $summary = .String "movie.overview" }}
                {{ $popoverTitle = .String "movie.title" }}
                {{ $popoverSubtitle = printf "%s %s" $datetype (($date.In $timezone) | formatTime "1/2/2006") }}
                {{ $popoverSummary = .String "movie.overview" }}
                {{ $link = printf "%s/movie/%s#" $url (.String "movie.titleSlug") }}
                {{ $movie = .Get "movie" }}
                {{ $genres = $movie.Get "genres" }}
                {{ if and $showGrabbed (gt (.Int "movie.movieFileId") 0) }}
                  {{ $grabbed = true }}
                {{ end }}
              {{ else }}
                {{ if eq $coverProxy "" }}
                  {{ $coverUrl = printf "%s/%s/poster-500.jpg?apikey=%s" $coverBase (.String "id") $key }}
                {{ else }}
                  {{ $coverUrl = printf "%s/%s/poster-500.jpg" $coverBase (.String "id") }}
                {{ end }}
                {{ $title = .String "title" }}
                {{ $subtitle = "" }}
                {{ $summary = .String "overview" }}
                {{ $link = printf "%s/movie/%s#" $url (.String "titleSlug") }}
                {{ $popoverTitle = .String "title" }}
                {{ $popoverSubtitle = printf "%s %s" $datetype (($date.In $timezone) | formatTime "1/2/2006") }}
                {{ $popoverSummary = .String "overview" }}
                {{ $genres = .Get "genres" }}
                {{ if eq $type "missing" }}
                  {{ if and $showGrabbed (.Bool "movie.hasFile") }}
                    {{ $grabbed = true }}
                  {{ end }}
                {{ else }}
                  {{ if and $showGrabbed (.Bool "hasFile") }}
                    {{ $grabbed = true }}
                  {{ end }}
                {{ end }}
              {{ end }}
  
        
            {{ else if eq $service "lidarr" }}
              {{ $artist := "" }}
              {{ $album := "" }}
              {{ $albumId := "" }}
              {{ if eq $coverProxy "" }}
                {{ $coverBase = printf "%s/api/v1/mediacover" $apiBaseUrl }}
              {{ else }}
                {{ $coverBase = $coverProxy }}
              {{ end }}
   
              {{ if eq $type "recent" }}
                {{ $album = .Get "album" }}
                {{ $artist = $album.Get "artist" }}
                {{ if eq $coverProxy "" }}
                  {{ $coverUrl = printf "%s/album/%s/cover-500.jpg?apikey=%s" $coverBase (.String "albumId") $key }}
                {{ else }}
                  {{ $coverUrl = printf "%s/album/%s/cover-500.jpg" $coverBase (.String "albumId") }}
                {{ end }}
                {{ $grabbed = true }}
                {{ $title = $album.String "title" }}
                {{ $date = (.String "date" | parseTime "rfc3339") }}
                {{ $datetype = "Downloaded" }}
                {{ $subtitle = $artist.String "artistName" }}
                {{ $genres = $artist.Get "genres" }}
                {{ $link = printf "%s/artist/%s#" $url ($artist.String "foreignArtistId") }}
                {{ $summary = $album.String "overview" }}
                {{ $popoverTitle = $album.String "title" }}
                {{ $popoverSubtitle = printf "%s %s" $datetype (($date.In $timezone) | formatTime "1/2/2006") }}
                {{ $popoverSummary = $artist.String "overview" }}
  
              {{ else }}
                {{ $artist = .Get "artist" }}
                {{ $album = .Get "album" }}
                {{ if eq $type "missing" }}
                  {{ $datetype = "Released" }}
                {{ else }}
                  {{ $datetype = "Releases" }}
                {{ end }}
                {{ range .Array "releases" }}
                  {{ $albumId = .String "albumId" }}
                  {{ break }}
                {{ end }}
                {{ if eq $coverProxy "" }}
                  {{ $coverUrl = printf "%s/album/%s/cover-500.jpg?apikey=%s" $coverBase $albumId $key }}
                {{ else }}
                  {{ $coverUrl = printf "%s/album/%s/cover-500.jpg" $coverBase $albumId }}
                {{ end }}
                {{ $date = (.String "releaseDate" | parseTime "rfc3339") }}
                {{ $title = .String "title" }}
                {{ $subtitle = $artist.String "artistName" }}
                {{ $summary = $album.String "overview" }}
                {{ $genres = $artist.Get "genres" }}
                {{ $link = printf "%s/artist/%s#" $url ($artist.String "foreignArtistId") }}
                {{ $popoverTitle = .String "title" }}
                {{ $popoverSubtitle = printf "%s %s" $datetype (($date.In $timezone) | formatTime "1/2/2006") }}
                {{ $popoverSummary = .String "artist.overview" }}
              {{ end }}
            {{ end }}
        
            {{ if eq $size "small" }}
              {{ $buttonJustify = "right" }}
              {{ $height = "9rem" }}
              {{ if eq $service "lidarr" }}
                {{ $width = "9rem" }}
              {{ else }}
                {{ $width = "6rem" }}
              {{ end }}
            {{ else if eq $size "medium" }}
              {{ $height = "12rem" }}
              {{ if eq $service "lidarr" }}
                {{ $width = "12rem" }}
              {{ else }}
                {{ $width = "8rem" }}
              {{ end }}
            {{ else if eq $size "large" }} 
              {{ $height = "15rem" }}
              {{ if eq $service "lidarr" }}
                {{ $width = "15rem" }}
              {{ else }}
                {{ $width = "10rem" }}
              {{ end }}
            {{ else if eq $size "huge" }} 
              {{ $height = "18rem" }}
              {{ if eq $service "lidarr" }}
                {{ $width = "18rem" }}
              {{ else }}
                {{ $width = "12rem" }}
              {{ end }}
            {{ end }}
    
            <li style="position: relative;">
              <div class="flex gap-10 items-start thumbnail-container thumbnail-parent">
                <div>
                  <div data-popover-type="html" data-popover-position="above" data-popover-show-delay="500" style="width: {{ $width  }}; height: {{ $height }}; align-content: center;">
                    <div data-popover-html>
                      <div style="margin: 5px;">
                        <strong class="size-h4 color-primary" title="{{ $popoverTitle }}">{{ $popoverTitle }}</strong>
                        <div class="size-h4 text-truncate text-very-compact color-subdue" title="{{ $popoverSubtitle }}">{{ $popoverSubtitle }}</div>
                        <p class="margin-top-20" style="overflow-y: auto; text-align: justify; max-height: 20rem;">
                          {{ if ne $popoverSummary "" }}
                            {{ $popoverSummary }}
                          {{ else }}
                            TBA
                          {{ end }}
                        </p>
                       {{ if gt (len ($genres.Array "")) 0 }}
                        <ul class="attachments margin-top-20">
                          {{ range $genres.Array "" }}
                            <li>{{ .String "" }}</li>
                          {{ end }}
                        </ul>
                        {{ end }}
                      </div>
                    </div>
                    <img class="thumbnail" src="{{ $coverUrl }}" alt="Cover for {{ .String "title" }}" loading="lazy" style="width: 100%; height: 100%; box-shadow: 0 10px 15px -3px rgb(0 0 0 / 0.1), 0 4px 6px -4px rgb(0 0 0 / 0.1); object-fit: cover; border-radius: 0.5rem;">
                  </div>
              </div>
                <div class="shrink min-width-0" style="height: 9rem; position: relative; padding-top: 5px; padding-right: 5px;">
                  <strong class="size-h4 block text-truncate color-primary" title="{{ $title }}">{{ $title }}</strong>
                  <div class="text-truncate text-very-compact" title="{{ $subtitle }}">{{ $subtitle }}</div>
                  <div class="text-very-compact text-truncate">
                    {{ if eq $service "sonarr" }}
                      <div>{{ $seString }} - {{ $datetype }} {{ ($date.In $timezone).Format "1/2 03:04PM" }}</div>
                    {{ else }}
                      <span>{{ $datetype }} {{ ($date.In $timezone).Format "1/2" }}</span>
                    {{ end }}
                  </div>
    
                  {{ if $showGrabbed }}
                    {{ if eq $buttonJustify "right" }}
                      </div>
                      <a href="{{ $link }}" class="bookmarks-link size-h4 color-primary" target="_blank" rel="noreferrer"
                        style="position: absolute; bottom: 1rem; right: 1rem;
                          {{ if $grabbed }} color: var(--color-positive); border: 1px solid var(--color-positive); {{ else }} color: var(--color-negative); border: 1px solid var(--color-negative); {{ end }}
                          font-weight: bold; padding: 2px 5px; border-radius: 3px; display: inline-block; margin-top: 5px; text-decoration: none;">
                      {{ if $grabbed }}Grabbed{{ else }}Missing{{ end }}
                      </a>
                    {{ else }}
                      <a href="{{ $link }}" class="bookmarks-link size-h4 color-primary" target="_blank" rel="noreferrer"
                        style="{{ if $grabbed }} color: var(--color-positive); border: 1px solid var(--color-positive); {{ else }} color: var(--color-negative); border: 1px solid var(--color-negative); {{ end }}
                          font-weight: bold; padding: 2px 5px; border-radius: 3px; display: inline-block; margin-top: 5px; text-decoration: none;">
                      {{ if $grabbed }}Grabbed{{ else }}Missing{{ end }}
                      </a>
                      </div>
                    {{ end }}
                  {{ else }}
                    <div class="{{ if eq $size "small" }}text-truncate{{ if eq $service "radarr" }} margin-top-5{{ end }}{{ else }}text-truncate-2-lines margin-top-5{{ end }}">
                      {{ $summary }}
                    </div>
                  {{ end }}
              </div>
            </li>
          {{ end }}
        {{ end }}
        {{ if not $itemDisplayed }}
          <li>No items found.</li>
        {{ end }}
      </ul>
    {{ end }}

```
</details>

### Environment variables
- `*ARR_API_URL` - Your *arr API URL, eg: `http://192.168.1.36:8989`. Used for option `api-base-url`.
- `*ARR_KEY` - Your *arr API key, it should be under Settings > General > Security. Used for option `key`.
- `*ARR_URL` - (optional) URL used for links, eg: `https://your-arr-domain.com`. Used for option `url` and `title-url`.

> [!NOTE]  
>  If you don't need a separate url for links or the titlebar there's no need to set the `url` option.

### User variables/options

| Options           | Short Description                                     |
| ----------------- | ----------------------------------------------------- |
| service           | sonarr, radarr, or lidarr                             |
| type              | upcoming, recent, or missing                          |
| size              | Thumbnail size: small, medium, large, or huge         |
| collapse-after    | Collapse after this many items                        |
| show-grabbed      | Show a clickable button displaying grabbed status     |
| timezone          | Timezone offset from UTC. Must have quotes.           |
| interval          | (optional, in days) only show items within this time  |
| sort-time         | (optional) controls the sorting direction             |
| cover-proxy       | (optional) Avoids exposing the API key                |
| api-base-url      | See Environmental Variables above                     |
| key               | See Environmental Variables above                     |
| url               | See Environmental Variables above                     |

- `service` - `sonarr`, `radarr`, or `lidarr`. All `types` are available for all three.

- `type`: 
  - `upcoming` displays items that are soon to release
  - `recent` displays items that have recently been grabbed
  - `missing` displays items that are considered available but haven't been grabbed

- `size` - size of the thumbnail for each item. `small` will move the Missing/Grabbed button to the right if `show-grabbed` is true

- `show-grabbed` - will show the grab status and removes the overview/description. Still available upon hovering the thumbnail.

- `timezone` - change `"-04"` to your timezone, eg: +9 will be `"+09"`. Must have quotes.

- `interval` - If set will only show items for the given interval in days. For "upcoming" it's days to look ahead, otherwise it's within past X days.

- `sort-time` - If set to `asc` (ascending time) will invert the default sort direction `desc` (descending time).

- `cover-proxy` - This can be done through your proxy manager of choice. For Nginx Proxy Manager specifically, inside the `Advanced` tab:
    ```nginx
    location /radarrcover/ {
        rewrite ^/radarrcover/(.*)$ /api/v3/mediacover/$1?apikey=your-api-key-here break;
        proxy_pass $forward_scheme://$server:$port;
    }
    ```
    then set your coverProxy:
    ```yaml
    cover-proxy: "https://proxy.example.com/radarrcover"
    ```
> [!TIP]  
> If using a cover-proxy for Lidarr, you must replace `v3` in the proxy rewrite url with `v1`.

## Widget reuse for multiple services/types

With v0.8.0 of Glance it's now possible to reuse a single widget template and declare multiple copies with different options. 

See [v0.8.0 Release Notes](https://github.com/glanceapp/glance/releases/tag/v0.8.0#g-rh-15) (expand the part about the `options`)

## Credits
[iwa](https://github.com/iwa) - [*arr Glance extension](https://github.com/glanceapp/glance/pull/112)

[ralphocdol](https://github.com/ralphocdol) - original Sonarr custom-api widget

[svilenmarkov](https://github.com/svilenmarkov) - 🐐