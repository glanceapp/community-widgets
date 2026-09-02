# Nebula Videos Widget

Glance widget to display the latest [Nebula](https://nebula.tv/) videos from selected creators/channels.

## Carousel Style

![widget screenshot](preview.png)

```yaml
- type: custom-api
  title: Nebula
  title-url: https://nebula.tv/
  cache: 1h
  frameless: true
  url: https://content.api.nebula.app/video_episodes/
  parameters:
    ordering: -published_at
    page_size: "100"
  options:
    max_items: 12
    channel-ids:
      - video_channel:9cea6296-223e-4c7e-a245-d96db75de32f # Jet Lag: The Game
      - video_channel:8f3a2a56-3f9f-4ce0-b105-ede41688d84b # Half as Interesting
      - video_channel:fe4d9c1c-017b-494c-9afc-e79e6859b211 # Wendover Productions
  template: |
    {{ if ne .Response.StatusCode 200 }}
      <div class="widget-content-frame padding-widget color-negative">Failed to fetch Nebula videos.</div>
    {{ else }}
      {{ $channels := .Options.JSON "channel-ids" }}
      {{ $maxItems := (index .Options "max_items") }}
      {{ if not $maxItems }}{{ $maxItems = 12 }}{{ end }}
      {{ $shown := 0 }}
      <div class="carousel-container">
        <div class="cards-horizontal carousel-items-container">
          {{ range $video := .JSON.Array "results" }}
            {{ if and (lt $shown $maxItems) (ne "" (findMatch (concat "(^|[^A-Za-z0-9:_-])" ($video.String "channel_id") "([^A-Za-z0-9:_-]|$)") $channels)) }}
              <div class="card widget-content-frame thumbnail-parent">
                {{ if $video.Exists "images.thumbnail.src" }}
                  <img class="video-thumbnail thumbnail" loading="lazy" src="{{ $video.String "images.thumbnail.src" }}" alt="">
                {{ else }}
                  <svg class="video-thumbnail" viewBox="0 0 120 120" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <rect width="120" height="120" fill="#EFF1F3"/>
                    <path fill-rule="evenodd" clip-rule="evenodd" d="M33.2503 38.4816C33.2603 37.0472 34.4199 35.8864 35.8543 35.875H83.1463C84.5848 35.875 85.7503 37.0431 85.7503 38.4816V80.5184C85.7403 81.9528 84.5807 83.1136 83.1463 83.125H35.8543C34.4158 83.1236 33.2503 81.957 33.2503 80.5184V38.4816ZM80.5006 41.1251H38.5006V77.8751L62.8921 53.4783C63.9172 52.4536 65.5788 52.4536 66.6039 53.4783L80.5006 67.4013V41.1251ZM43.75 51.6249C43.75 54.5244 46.1005 56.8749 49 56.8749C51.8995 56.8749 54.25 54.5244 54.25 51.6249C54.25 48.7254 51.8995 46.3749 49 46.3749C46.1005 46.3749 43.75 48.7254 43.75 51.6249Z" fill="#687787"/>
                  </svg>
                {{ end }}
                <div class="margin-top-10 margin-bottom-widget flex flex-column grow padding-inline-widget">
                  <a class="text-truncate-2-lines margin-bottom-auto color-primary-if-not-visited" href="{{ $video.String "share_url" }}" target="_blank" rel="noreferrer">{{ $video.String "title" }}</a>
                  <ul class="list-horizontal-text flex-nowrap margin-top-7">
                    <li class="shrink-0" {{ $video.String "published_at" | parseTime "rfc3339" | toRelativeTime }}></li>
                    <li class="min-width-0">
                      <a class="block text-truncate" href="{{ concat "https://nebula.tv/" ($video.String "channel_slug") }}" target="_blank" rel="noreferrer">{{ $video.String "channel_title" }}</a>
                    </li>
                  </ul>
                </div>
              </div>
              {{ $shown = add $shown 1 }}
            {{ end }}
          {{ end }}
        </div>
      </div>
    {{ end }`
```

## List Style

<img src="preview-list.png" alt="widget screenshot" height="400">

```yaml
- type: custom-api
  title: Nebula
  title-url: https://nebula.tv/
  cache: 1h
  frameless: true
  url: https://content.api.nebula.app/video_episodes/
  parameters:
    ordering: -published_at
    page_size: "100"
  options:
    max_items: 12
    collapse_after: 5
    channel-ids:
        - video_channel:9cea6296-223e-4c7e-a245-d96db75de32f # Jet Lag: The Game
        - video_channel:8f3a2a56-3f9f-4ce0-b105-ede41688d84b # Half as Interesting
        - video_channel:fe4d9c1c-017b-494c-9afc-e79e6859b211 # Wendover Productions
  template: |
    {{ if ne .Response.StatusCode 200 }}
      <div class="widget-content-frame padding-widget color-negative">Failed to fetch Nebula videos.</div>
    {{ else }}
      {{ $channels := .Options.JSON "channel-ids" }}
      {{ $maxItems := (index .Options "max_items") }}
      {{ $collapseAfter := (index .Options "collapse_after") }}
      {{ if not $maxItems }}{{ $maxItems = 12 }}{{ end }}
      {{ if not $collapseAfter }}{{ $collapseAfter = 5 }}{{ end }}
      {{ $shown := 0 }}

      <ul class="list list-gap-10 collapsible-container" data-collapse-after="{{ $collapseAfter }}" style="list-style: none; padding: 0; margin: 0;">
        {{ range $video := .JSON.Array "results" }}
          {{ if and (lt $shown $maxItems) (ne "" (findMatch (concat "(^|[^A-Za-z0-9:_-])" ($video.String "channel_id") "([^A-Za-z0-9:_-]|$)") $channels)) }}
            <a href="{{ $video.String "share_url" }}" target="_blank" rel="noreferrer" style="text-decoration: none;">
              <li style="padding: 10px 0; border-bottom: 1px solid var(--border-color);">
                <div style="display: flex; gap: 12px; align-items: flex-start;">
                  <div style="flex-shrink: 0; width: 120px; height: 68px; border-radius: 6px; overflow: hidden;">
                    {{ if $video.Exists "images.thumbnail.src" }}
                      <img src="{{ $video.String "images.thumbnail.src" }}" alt="" style="width: 100%; height: 100%; object-fit: cover;">
                    {{ else }}
                      <svg viewBox="0 0 120 120" fill="none" xmlns="http://www.w3.org/2000/svg" style="width: 100%; height: 100%;">
                        <rect width="120" height="120" fill="#EFF1F3"/>
                        <path fill-rule="evenodd" clip-rule="evenodd" d="M33.2503 38.4816C33.2603 37.0472 34.4199 35.8864 35.8543 35.875H83.1463C84.5848 35.875 85.7503 37.0431 85.7503 38.4816V80.5184C85.7403 81.9528 84.5807 83.1136 83.1463 83.125H35.8543C34.4158 83.1236 33.2503 81.957 33.2503 80.5184V38.4816ZM80.5006 41.1251H38.5006V77.8751L62.8921 53.4783C63.9172 52.4536 65.5788 52.4536 66.6039 53.4783L80.5006 67.4013V41.1251ZM43.75 51.6249C43.75 54.5244 46.1005 56.8749 49 56.8749C51.8995 56.8749 54.25 54.5244 54.25 51.6249C54.25 48.7254 51.8995 46.3749 49 46.3749C46.1005 46.3749 43.75 48.7254 43.75 51.6249Z" fill="#687787"/>
                      </svg>
                    {{ end }}
                  </div>
                  <div style="flex: 1; min-width: 0; display: flex; flex-direction: column; gap: 6px;">
                    <div class="text-truncate-2-lines color-primary">{{ $video.String "title" }}</div>
                    <div class="size-h6 color-subdue" style="display: flex; gap: 6px; align-items: center;">
                      <span>{{ $video.String "published_at" | parseTime "rfc3339" | toRelativeTime }}</span>
                      <span>•</span>
                      <span class="text-truncate">{{ $video.String "channel_title" }}</span>
                    </div>
                  </div>
                </div>
              </li>
            </a>
            {{ $shown = add $shown 1 }}
          {{ end }}
        {{ end }}
      </ul>
    {{ end }}
```

### Configuring

- `max_items`: maximum number of videos to render (default: `12`)
- `collapse_after`: list-only collapse threshold (default: `5`)
- `channel-ids`: only include videos from these channel IDs

#### How to find a channel ID
1. Go to the channel's page on Nebula (e.g. https://nebula.tv/channel/half-as-interesting)
2. Open dev tools and go to the Network tab
3. Refresh the page and look for a request to `https://content.api.nebula.app/video_channels/` (you can filter by "video_channels" in the Network tab)
4. Click on that request and look at the "Response" tab, where you'll see a request which includes the channel ID (e.g. `video_channel:8f3a2a56-3f9f-4ce0-b105-ede41688d84b` for Half as Interesting). It can look something like this: ![channel_id_screenshot](find_channel_id_devtools.png)

## Notes
- The Nebula endpoint for the latest videos only fetches 100 at most, meaning that sometimes you won't hit the maximum item count (depending on how many channel ids you have and frequency of video releases)
- This endpoint is not officially supported and there is no public Nebula API, so it may break at any time if they change their internal APIs. Please use this API carefully and do not make excessive requests to it.
