## Screenshots
#### Normal
![Preview widget](./preview.png)

```yaml
- type: custom-api
  title: Steam Top Sellers
  cache: 12h
  url: https://store.steampowered.com/api/featuredcategorcc=us
  template: |
    <ul class="list list-gap-10 collapsible-containdata-collapse-after="15">
    {{ range .JSON.Array "top_sellers.items" }}
      {{ if ne (.String "name") "Steam Deck" }}
      <li style="display: flex; align-items: center; gap: 1">
        <img src="{{ .String "small_capsule_image" }}" a{{ .String "name" }}" style="width: 120px; heigauto; border-radius: 4px; flex-shrink: 0;">
        <div style="min-width: 0;">
          <a class="size-h4 color-highlight bltext-truncate" href="https://store.steampowered.app/{{ .Int "id" }}/">{{ .String "name" }}</a>
          <ul class="list-horizontal-text">
            <li>{{ div (.Int "final_price" | toFloat) 10printf "$%.2f" }}</li>
            {{ $discount := .Int "discount_percent" }}
            {{ if gt $discount 0 }}
            <li{{ if ge $discount 40 }} class="color-posit{{ end }}>{{ $discount }}% off</li>
            {{ end }}
          </ul>
        </div>
      </li>
      {{ end }}
    {{ end }}
    </ul>
```
