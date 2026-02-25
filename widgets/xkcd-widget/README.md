# xkcd Display
Simple custom-api widget to display current xkcd comic

**xkcd**
```
- type: custom-api
  title: xkcd
  cache: 2m
  url: https://xkcd.com/info.0.json
  template: | 
    <a
      href="https://xkcd.com/" 
      target="_blank" 
      rel="noopener noreferrer"
    >
      {{ .JSON.String "title" }}
    </a>
    <img src="{{ .JSON.String "img" }}" title="{{ .JSON.String "alt" }}"></img>
```
<img src="preview.png">