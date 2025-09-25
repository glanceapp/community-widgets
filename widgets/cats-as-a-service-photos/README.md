## Screenshots
#### Normal
![Preview widget image](./preview.png)

#### No tags found
![No tag found preview](./preview-error.png)

```yaml
  - type: custom-api
    title: Cat hall of Fame (From CaaS)
    cache: 1m
    url: https://cataas.com/api/cats?limit=${NUMBER_OF_CATS}&tags=${CAT_TAGS}
    headers:
        Accept: application/json
    template: |
        {{if eq (len (.JSON.Array "")) 0}}
        <div class="widget-error-header">
            <p class="color-negative size-h3">Uh oh</p>
            <svg class="widget-error-icon" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126ZM12 15.75h.007v.008H12v-.008Z"></path>
            </svg>
        </div>
        <p style="font-size: var(--font-size-h4);">No cat was found with those tags, perhaps you should try with fewer tags! ฅ՞•ﻌ•՞ฅ</p>
        {{else}}
        <div class="cards-grid">
            {{range .JSON.Array ""}}
            <div style="display: flex; justify-content:center;">
                <img src="https://cataas.com/cat/{{.String "id"}}?type=medium" style="object-fit: cover; border-radius: 10px;" class="card">
            </div>
            {{end}} 
        </div>
        {{end}}
```

## Environment variables

- `NUMBER_OF_CATS` - total number of cats per fetch, keep this value low for faster fetching and lower bandwidth, recommended at `12`
- `CAT_TAGS` - the cat tag(s), seperated by comma(s), for example `cute,fluffy,orange`. Note that with too much tags, there won't be sufficient (or maybe none) images match your query.
