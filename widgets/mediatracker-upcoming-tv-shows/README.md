## Description

This widget fetches upcoming movies, TV shows or video games from an instance of [MediaTracker](https://github.com/bonukai/MediaTracker) and shows them as cards in horizontal carousel. Each card contains a poster, an episode number (for a TV show), an episode name (for a TV show), a name (of a TV show, movie or a video game) and its release date.

TV Shows:
![Preview of the widget (TV Shows)](preview_tv.png)
Movies:
![Preview of the widget (Movies)](preview_movies.png)
Video Games:
![Preview of the widget (Video Games)](preview_video-games.png)
No upcoming media:
![Preview of the widget (Nothing fetched)](preview_nocards.png)


The widget requires modifications to the configuration file in order to work. Under the `options` field in the yaml configuration file, one needs to supply MediaTracker URL, API key and some other options. One can either supply them directly or as environment variables. 

## List of Options

- (required) `base-url` requires URL pointing to your MediaTracker, be sure to include "http://" or "https://"
- (required) `api-key` requires MediaTracker API key that can be retrieved in MediaTracker under "ACCOUNTNAME -> Application tokens -> Add Token"
- (optional) `media-type` requires "tv", "movie" or "video_game". The default value is "tv".
- (optional) `only-on-watchlist` can be either "true" or "false". Determines whether only media items from the watchlist are fetched. The default value is "false".
- (optional) `items` value limits number of cards. The default value is 10.

## Widget Configuration YAML
```yaml
- type: custom-api
  title: Upcoming # change to your liking
  frameless: true
  cache: 1h
  options:
    base-url: ${MEDIATRACKER_URL}       # URL pointing to your MediaTracker
    api-key: ${MEDIATRACKER_API_KEY}    # can be retrieved in Account Name -> Application tokens
    media-type: tv                      # tv, movie or video_game
    only-on-watchlist: "false"          # "true" or "false"
    items: 10                           # integer number
  template: |
    {{ $baseURL := .Options.StringOr "base-url" "" }}
    {{ $apiKey := .Options.StringOr "api-key" "" }}
    {{ $mediaType := .Options.StringOr "media-type" "tv" }}
    {{ $onlyOnWatchlist := .Options.StringOr "only-on-watchlist" "false" }}
    {{ $numberOfItems := .Options.IntOr "items" 10 }}
    {{ $orderBy := "releaseDate" }}
    {{ $sortOrder := "asc" }}
    {{ $videoGameReleaseDate := "TBA" }}
    
    {{ if or (eq $baseURL "") (eq $apiKey "") }}
        <p class="color-negative">MediaTracker base-url or API token not set.</p>
    {{ else }}
      {{ if eq $mediaType "tv" }}
        {{ $orderBy = "nextAiring" }}
      {{ end }}
      {{ $requestURL := concat $baseURL "/api/items?mediaType=" $mediaType "&orderBy=" $orderBy "&sortOrder=" $sortOrder "&onlyOnWatchlist" $onlyOnWatchlist "&token=" $apiKey}}
      {{ if or (eq $mediaType "video_game") (eq $mediaType "movie") }}
        {{ $requestURL = concat $requestURL "&onlyWithNextAiring=true" }}
      {{ end }}
      
      {{ $items := newRequest $requestURL | getResponse }}
      {{ if eq $items.Response.StatusCode 200 }}
        {{ $arr := $items.JSON.Array "" }}
        {{ $len := len $arr }}
        {{ $shown := 0 }}
        
        {{ if gt (len $arr) 0 }}
            <div class="cards-horizontal carousel-items-container">
              {{ range $i, $_ := $arr }}
                {{ if and (.Exists "nextAiring") (ne (.String "nextAiring") "") (lt $shown $numberOfItems) }}
                     {{ $el := index $arr $i }}
                     
                     <a class="card widget-content-frame"
                        href="{{ $baseURL }}/#/details/{{ $el.Int "id" }}"
                        style="
                            flex: 0 0 25vh; 
                            min-width: 170px;
                            min-heght: 150px; 
                            display:flex; 
                            flex-direction:column; 
                            box-sizing:border-box;
                            text-decoration:none; 
                            color:inherit;"
                     >
                     <div style="position: relative;">
                        {{ if and ($el.Exists "posterSmall") (ne ($el.String "posterSmall") "") }}
                            <img src="{{ $baseURL }}{{ $el.String "posterSmall" }}&token={{ $apiKey }}"
                                alt="{{ $el.String "title" }}"
                                loading="lazy"
                                class="media-server-thumbnail shrink-0 loaded finished-transition"
                                style="
                                    width:100%; 
                                    display:block; 
                                    border-radius:var(--border-radius) var(--border-radius) 0 0; 
                                    object-fit:cover;" 
                            />
                         {{ else }}
                            <div style="width:100%; background:#222; display:flex; align-items:center; justify-content:center; color:#ddd; aspect-ratio:2/3;">
                              No image
                            </div>
                         {{ end }}
                     </div>
                     <div class="grow padding-inline-widget margin-top-10 margin-bottom-10">
                          <ul class="flex flex-column justify-evenly margin-bottom-3" style="height:100%; gap: 4px;">
                            {{ if eq $mediaType "tv" }}
                                <ul class="list-horizontal-text flex-nowrap" style="padding:0; margin:0;">
                                <li class="color-primary shrink-0"> S{{ printf "%02d" (.Int "upcomingEpisode.seasonNumber") }}E{{ printf "%02d" (.Int "upcomingEpisode.episodeNumber") }}</li>
                                <li class="text-truncate">{{ .String "upcomingEpisode.title" }}</li>
                                </ul>
                            {{ end }}
                            <li class="text-truncate color-primary" style="overflow:hidden; text-overflow:ellipsis; white-space:nowrap;" title="{{ .String "title" }}">
                              {{ .String "title" }}
                            </li>
                            {{ if eq $mediaType "video_game" }}
                                <li style="font-size:0.85em; opacity:0.7; white-space:nowrap; overflow:hidden; text-overflow:ellipsis;">                                    
                                    {{ $videoGameReleaseDate = .String "releaseDate" | parseTime "RFC3339" }}
                                    {{ if eq ($videoGameReleaseDate | formatTime "02-01") "31-12" }}
                                        {{ $videoGameReleaseDate | formatTime "2006" }}
                                    {{ else }}
                                        {{ $videoGameReleaseDate | formatTime "Jan 02, 2006" }}
                                    {{ end }}
                                </li>
                            {{ else}}
                                <li style="font-size:0.85em; opacity:0.7; white-space:nowrap; overflow:hidden; text-overflow:ellipsis;">
                                    {{ .String "nextAiring" | parseTime "2006-01-02" | formatTime "Jan 02, 2006" }}
                                </li>
                            {{ end }}
                          </ul>
                       </div>
                     </a>
                {{ $shown = add $shown 1 }}
                {{ end }}
              {{ end }}
            </div>
        {{ else }}
            <div class="widget-content-frame" style="flex=0 0 25vh; display:flex">
              <div class="grow padding-inline-widget margin-top-10 margin-bottom-10">
                No upcoming entries on the watchlist.
              </div>
            </div>
        {{ end }}
      {{ else }}
         <div class="widget-content-frame" style="flex=0 0 25vh; display:flex">
              <div class="grow padding-inline-widget margin-top-10 margin-bottom-10">
                Failed to fetch items (status {{ $items.Response.StatusCode }})
              </div>
         </div>
      {{ end }}
    {{ end }}
```

