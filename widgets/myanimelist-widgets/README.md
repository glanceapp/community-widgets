

## Carousel Style

![widget screenshot](preview-carousel.png)

```yaml
- type: custom-api
  title: Seasonal Anime
  frameless: true
  cache: 1d
  url: https://api.jikan.moe/v4/seasons/now
  template: |
    {{ $arr := .JSON.Array "data" }}
    <div style="overflow-x: auto; padding: 8px 0;">
      <div class="cards-horizontal carousel-items-container" style="display: flex; gap: 16px; padding: 0 8px;">
        {{ range $i, $el := $arr }}
          {{ if lt $i 15 }}
            {{ $image := $el.String "images.jpg.image_url" }}
            {{ $score := $el.Float "score" }}
            {{ $type := $el.String "type" }}
            {{ $episodes := $el.Int "episodes" }}
            
            <a href="{{ $el.String "url" }}" target="_blank" 
              class="card widget-content-frame" 
              style="flex: 0 0 auto; width: 150px; min-height: 260px; display: flex; flex-direction: column; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 12px rgba(0,0,0,0.1); background: var(--card-bg); text-decoration: none;">
              
              <!-- Thumbnail with Score -->
              <div style="height: 190px; overflow: hidden; position: relative;">
                <img src="{{ $image }}" alt="{{ $el.String "title" }}"
                    style="width: 100%; height: 100%; object-fit: cover; object-position: center;"
                    onerror="this.style.display='none'">
                {{ if gt $score 0.0 }}
                <div style="position: absolute; top: 8px; right: 8px; background: rgba(0,0,0,0.9); color: white; padding: 3px 8px; border-radius: 12px; font-weight: bold; backdrop-filter: blur(4px); font-size: 11px;">
                  ⭐ {{ $score }}
                </div>
                {{ end }}
              </div>
              
              <!-- Content -->
              <div style="padding: 12px; flex-grow: 1; display: flex; flex-direction: column; justify-content: space-between; min-height: 70px;">
                <div class="size-base color-primary" 
                    style="line-height: 1.3; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; min-height: 2.6em; margin-bottom: 8px;">
                  {{ $el.String "title" }}
                </div>
                
                <div class="size-h6" style="display: flex; flex-direction: column; gap: 4px;">
                  <div style="display: flex; align-items: center; gap: 6px;">
                    <span>Type:</span>
                    <span>{{ $type }}</span>
                  </div>
                  <div style="display: flex; align-items: center; gap: 6px;">
                    <span>{{ if $episodes }}{{ $episodes }} Episodes{{ else }}Airing{{ end }}</span>
                  </div>
                </div>
              </div>
            </a>
          {{ end }}
        {{ end }}
      </div>
    </div>
```
## List Style

![widget screenshot](preview-list.png)

```yaml
- type: custom-api
  title: Seasonal Anime
  cache: 1d
  url: https://api.jikan.moe/v4/seasons/now
  options:
    max_items: 8
    collapse_after: 5
  template: |
    {{ $arr := .JSON.Array "data" }}
    {{ $maxItems := (index .Options "max_items") }}
    {{ $collapseAfter := (index .Options "collapse_after") }}
    {{ if not $maxItems }}{{ $maxItems = 25 }}{{ end }}
    
    <ul class="list list-gap-10 collapsible-container" data-collapse-after="{{ $collapseAfter }}" style="list-style: none; padding: 0; margin: 0;">
      {{ range $i, $el := $arr }}
        {{ if lt $i $maxItems }}
          {{ $image := $el.String "images.jpg.image_url" }}
          {{ $type := $el.String "type" }}
          {{ $episodes := $el.Int "episodes" }}
          {{ $score := $el.Float "score" }}
          
          <a href="{{ $el.String "url" }}" target="_blank" style="text-decoration: none;">
            <li style="padding: 10px 0; border-bottom: 1px solid var(--border-color);">
              <div style="display: flex; gap: 12px; align-items: flex-start;">
                
                <!-- Thumbnail -->
                <div style="flex-shrink: 0; width: 80px; height: 112px;">
                  <img src="{{ $image }}" alt=""
                       style="width: 100%; height: 100%; object-fit: cover; border-radius: 4px;"
                       onerror="this.style.display='none'">
                </div>
                
                <!-- Content -->
                <div style="flex: 1; min-width: 0; display: flex; flex-direction: column; height: 112px; justify-content: space-between;">
                  
                  <!-- Title -->
                  <div class="size-base color-primary" 
                       style="line-height: 1.3; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; min-height: 2.6em; margin-bottom: 4px;">
                    {{ $el.String "title" }}
                  </div>
                  
                  <!-- Metadata -->
                  <div class="size-h6" style="display: flex; flex-direction: column; gap: 2px;">
                    <div style="display: flex; align-items: center; gap: 6px;">
                      <span>Type:</span>
                      <span>{{ $type }}</span>
                    </div>
                    
                    <div style="display: flex; align-items: center; gap: 6px;">
                      <span>Episode:</span>
                      <span>{{ if $episodes }}{{ $episodes }}{{ else }}Airing{{ end }}</span>
                    </div>
                    
                    <div style="display: flex; align-items: center; gap: 6px;">
                      <span>Rating:</span>
                      <span>{{ if gt $score 0.0 }}{{ $score }}{{ else }}N/A{{ end }}</span>
                    </div>
                  </div>
                  
                </div>
              </div>
            </li>
          </a>
        {{ end }}
      {{ end }}
    </ul>
```
### Configuring:
- `max_items`: set maximum amount of items to show when expanding list
- `collapse_after`: collapse after `x` amount of anime

### Notes
- Ratings are only displayed when available (score > 0.0)
- Episode count shows "Airing" for currently airing anime without final episode count

## Credits
These widgets use the [Jikan Api](https://jikan.moe/), an unofficial MyAnimelst API
