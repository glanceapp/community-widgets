![](preview.png)

```yaml
- type: custom-api
  title: Random Fact
  cache: 3h
  url: https://uselessfacts.jsph.pl/api/v2/facts/random
  template: |
    <p class="size-h4 color-paragraph">{{ .JSON.String "text" }}</p>
```
