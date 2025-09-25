![](preview.png)

```yaml
- type: custom-api
  title: Steam Specials
  cache: 12h
  url: https://store.steampowered.com/api/featuredcategories?cc=us
  template: |
    <ul class="list list-gap-10 collapsible-container" data-collapse-after="5">
    {{ range .JSON.Array "specials.items" }}
      {{ $header := .String "header_image" }}
      {{ $urlPrefix := "https://store.steampowered.com/sub/" }}
      {{ if findMatch "/steam/apps/" $header }}
        {{ $urlPrefix = "https://store.steampowered.com/app/" }}
      {{ end }}
      <li>
        <a class="size-h4 color-highlight block text-truncate" href="{{ $urlPrefix }}{{ .Int "id" }}/">{{ .String "name" }}</a>
        <ul class="list-horizontal-text">
          <li>{{ .Int "final_price" | toFloat | mul 0.01 | printf "$%.2f" }}</li>
          {{ $discount := .Int "discount_percent" }}
          <li{{ if ge $discount 40 }} class="color-positive"{{ end }}>{{ $discount }}% off</li>
        </ul>
      </li>
    {{ end }}
    </ul>
```

## Changing currency

You can change the currency by changing the `cc` parameter in the URL to your country's [2 character code](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2). For example, to get the prices in euros, you can change the URL to `https://store.steampowered.com/api/featuredcategories?cc=eu` and then change the `$` symbol to `€` in the template.
