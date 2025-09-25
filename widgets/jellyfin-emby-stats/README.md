## Preview

![](preview.png)

## Configuration

```yaml
- type: custom-api
  title: "Jellyfin/Emby Stats"
  base-url: ${JELLYFIN_URL}
  options:
    url: ${JELLYFIN_URL}
    key: ${JELLYFIN_API_KEY}

  template: |
    {{ $url := .Options.StringOr "url" "" }}
    {{ $key := .Options.StringOr "key" "" }}

    {{- if or (eq $url "") (eq $key "") -}}
    
      <p>Error: The URL or API Key was not configured in the widget options.</p>
      
    {{- else -}}

      {{- $requestUrl := printf "%s/emby/Items/Counts?api_key=%s" $url $key -}}
      {{- $jellyfinData := newRequest $requestUrl | getResponse -}}

      {{- if eq $jellyfinData.Response.StatusCode 200 -}}
        <div class="flex flex-column gap-5">
          <div class="flex justify-between text-center">
            
            <div>
              <div class="color-highlight size-h3">{{ $jellyfinData.JSON.Int "MovieCount" | formatNumber }}</div>
              <div class="size-h5 uppercase">Movies</div>
            </div>

            <div>
              <div class="color-highlight size-h3">{{ $jellyfinData.JSON.Int "SeriesCount" | formatNumber }}</div>
              <div class="size-h5 uppercase">TV Shows</div>
            </div>

            <div>
              <div class="color-highlight size-h3">{{ $jellyfinData.JSON.Int "EpisodeCount" | formatNumber }}</div>
              <div class="size-h5 uppercase">Episodes</div>
            </div>

            <div>
              <div class="color-highlight size-h3">{{ $jellyfinData.JSON.Int "SongCount" | formatNumber }}</div>
              <div class="size-h5 uppercase">Songs</div>
            </div>

          </div>
        </div>
      {{- else -}}
        <p>Failed: {{ $jellyfinData.Response.Status }}</p>
      {{- end -}}
    {{- end -}}
```