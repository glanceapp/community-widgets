![](preview.png)

```yaml
- type: custom-api
  title: Steam Specials
  cache: 12h
  url: https://store.steampowered.com/api/featuredcategories?cc=us
  options:
    show-thumbnails: true
    limit: 5
  template: |
    {{ $showThumbnails := .Options.BoolOr "show-thumbnails" false }}
    {{ $limit := .Options.IntOr "limit" 0 }}
    <ul class="list list-gap-10 collapsible-container" data-collapse-after="5">
    {{ range $i, $v := .JSON.Array "specials.items" }}
      {{ if or (eq $limit 0) (lt $i $limit) }}
        {{ $header := $v.String "header_image" }}
        {{ $urlPrefix := "https://store.steampowered.com/sub/" }}
        {{ if findMatch "/steam/apps/" $header }}
          {{ $urlPrefix = "https://store.steampowered.com/app/" }}
        {{ end }}
        {{ if $showThumbnails }}
          <li style="display: flex; align-items: center; gap: 1rem;">
            <img src="{{ $v.String "small_capsule_image" }}" alt="{{ $v.String "name" }}"
                style="width: 120px; height: auto; border-radius: 4px; flex-shrink: 0;">
            <div style="min-width: 0;">
              <a class="size-h4 color-highlight block text-truncate" href="{{ $urlPrefix }}{{ $v.Int "id" }}/">{{ $v.String "name" }}</a>
              <ul class="list-horizontal-text">
                <li>{{ $v.Int "final_price" | toFloat | mul 0.01 | printf "€%.2f" }}</li>
                {{ $discount := $v.Int "discount_percent" }}
                <li{{ if ge $discount 40 }} class="color-positive"{{ end }}>{{ $discount }}% off</li>
              </ul>
            </div>
          </li>
        {{ else }}
          <li>
            <a class="size-h4 color-highlight block text-truncate" href="{{ $urlPrefix }}{{ $v.Int "id" }}/">{{ $v.String "name" }}</a>
            <ul class="list-horizontal-text">
              <li>{{ $v.Int "final_price" | toFloat | mul 0.01 | printf "€%.2f" }}</li>
              {{ $discount := $v.Int "discount_percent" }}
              <li{{ if ge $discount 40 }} class="color-positive"{{ end }}>{{ $discount }}% off</li>
            </ul>
          </li>
        {{ end }}
      {{ end }}
    {{ end }}
    </ul>
```

## Changing currency

You can change the currency by changing the `cc` parameter in the URL to your country's [2 character code](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2). For example, to get the prices in euros, you can change the URL to `https://store.steampowered.com/api/featuredcategories?cc=eu` and then change the `$` symbol to `€` in the template.
