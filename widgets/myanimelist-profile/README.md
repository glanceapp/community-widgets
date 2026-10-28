## MyAnimeList Statistics

A widget to display your MyAnimeList statistics.

<div style="display: flex;">
  <img src=preview.png width="300">
  <img src=preview_theme.png width="300">
</div>

### Widget yaml

<details>

<summary>Default style</summary>

```yaml
- type: custom-api
  title: MAL Profile
  url: https://jikan.lucashdo.com/v1/users/${MAL_USER}/full
  cache: 1d
  template: |
    {{ if not (.JSON.Exists "data") }}
    <p class="color-negative">Error reading response</p>
    {{ else }}

    {{ $data := "data" }}
    {{ $profile := printf "%s.profile" $data }}
    {{ $anime := printf "%s.statistics.anime" $data }}
    {{ $manga := printf "%s.statistics.manga" $data }}


    <a class="block margin-bottom-10" href="{{ .JSON.String (printf "%s.profileUrl" $profile) }}" target="_blank" rel="noreferrer">
      <div class="flex items-center gap-12">
        <span class="text-truncate grow size-title-dynamic color-highlight" style="font-weight: 900;">{{ .JSON.String (printf "%s.canonicalUsername" $profile) }}</span>
        <img src="{{ .JSON.String (printf "%s.avatarUrl" $profile) }}" alt="Avatar" style="width: 32px; height: 32px; border-radius: 50%;" />
      </div>
    </a>

    <ul class="dynamic-columns list-gap-8">
      <div class="flex items-center gap-12"><div class="grow">Anime Days</div><div>{{ .JSON.Float (printf "%s.daysWatched" $anime) }}</div></div>
      <div class="flex items-center gap-12"><div class="grow">Anime Score</div><div>{{ .JSON.Float (printf "%s.meanScore" $anime) }}/10</div></div>
      <div class="flex items-center gap-12"><div class="grow">Episodes</div><div>{{ .JSON.Int (printf "%s.episodesWatched" $anime) }}</div></div>
      <div class="flex items-center gap-12"><div class="grow">Completed</div><div>{{ .JSON.Int (printf "%s.completed" $anime) }}</div></div>

      <div class="flex items-center gap-12"><div class="grow">Manga Days</div><div>{{ .JSON.Float (printf "%s.daysRead" $manga) }}</div></div>
      <div class="flex items-center gap-12"><div class="grow">Manga Score</div><div>{{ .JSON.Float (printf "%s.meanScore" $manga) }}/10</div></div>
      <div class="flex items-center gap-12"><div class="grow">Chapters</div><div>{{ .JSON.Int (printf "%s.chaptersRead" $manga) }}</div></div>
      <div class="flex items-center gap-12"><div class="grow">Completed</div><div>{{ .JSON.Int (printf "%s.completed" $manga) }}</div></div>

      {{ $watching := printf "%s.watching" $anime }}
      {{ if .JSON.Exists $watching }}
      <div class="flex items-center gap-12"><div class="grow">Currently Watching</div><div>{{ .JSON.Int $watching }}</div></div>
      {{ end }}

      {{ $reading := printf "%s.reading" $manga }}
      {{ if .JSON.Exists $reading }}
      <div class="flex items-center gap-12"><div class="grow">Currently Reading</div><div>{{ .JSON.Int $reading }}</div></div>
      {{ end }}
    </ul>
    {{ end }}
```

</details>


<details>

<summary>Theme color style</summary>

```yaml
- type: custom-api
  title: MAL Profile
  url: https://jikan.lucashdo.com/v1/users/${MAL_USER}/full
  cache: 1d
  template: |
    {{ if not (.JSON.Exists "data") }}
    <p class="color-negative">Error reading response</p>
    {{ else }}

    {{ $data := "data" }}
    {{ $profile := printf "%s.profile" $data }}
    {{ $anime := printf "%s.statistics.anime" $data }}
    {{ $manga := printf "%s.statistics.manga" $data }}


    <a class="block margin-bottom-10" href="{{ .JSON.String (printf "%s.profileUrl" $profile) }}" target="_blank" rel="noreferrer">
      <div class="flex items-center gap-12">
        <span class="text-truncate grow size-title-dynamic color-highlight" style="font-weight: 900;">{{ .JSON.String (printf "%s.canonicalUsername" $profile) }}</span>
        <img src="{{ .JSON.String (printf "%s.avatarUrl" $profile) }}" alt="Avatar" style="width: 32px; height: 32px; border-radius: 50%;" />
      </div>
    </a>

    <ul class="dynamic-columns list-gap-8">
      <div class="flex items-center gap-12"><div class="color-highlight grow color-primary">Anime Days</div><div>{{ .JSON.Float (printf "%s.daysWatched" $anime) }}</div></div>
      <div class="flex items-center gap-12"><div class="color-highlight grow color-primary">Anime Score</div><div>{{ .JSON.Float (printf "%s.meanScore" $anime) }}/10</div></div>
      <div class="flex items-center gap-12"><div class="color-highlight grow color-primary">Episodes</div><div>{{ .JSON.Int (printf "%s.episodesWatched" $anime) }}</div></div>
      <div class="flex items-center gap-12"><div class="color-highlight grow color-primary">Completed</div><div>{{ .JSON.Int (printf "%s.completed" $anime) }}</div></div>

      <div class="flex items-center gap-12"><div class="color-highlight grow color-primary">Manga Days</div><div>{{ .JSON.Float (printf "%s.daysRead" $manga) }}</div></div>
      <div class="flex items-center gap-12"><div class="color-highlight grow color-primary">Manga Score</div><div>{{ .JSON.Float (printf "%s.meanScore" $manga) }}/10</div></div>
      <div class="flex items-center gap-12"><div class="color-highlight grow color-primary">Chapters</div><div>{{ .JSON.Int (printf "%s.chaptersRead" $manga) }}</div></div>
      <div class="flex items-center gap-12"><div class="color-highlight grow color-primary">Completed</div><div>{{ .JSON.Int (printf "%s.completed" $manga) }}</div></div>

      {{ $watching := printf "%s.watching" $anime }}
      {{ if .JSON.Exists $watching }}
      <div class="flex items-center gap-12"><div class="color-highlight grow color-primary">Currently Watching</div><div>{{ .JSON.Int $watching }}</div></div>
      {{ end }}

      {{ $reading := printf "%s.reading" $manga }}
      {{ if .JSON.Exists $reading }}
      <div class="flex items-center gap-12"><div class="color-highlight grow color-primary">Currently Reading</div><div>{{ .JSON.Int $reading }}</div></div>
      {{ end }}
    </ul>
    {{ end }}
```

</details>


### Environment variable

- `MAL_USER` - your MyAnimeList username.

## Credits

[WaYy](https://github.com/WaYyTempest) - Widget creator
[Swiirz](https://github.com/aschnepp) - Widget api fixer

[Jikan](https://github.com/jikan-me) - Creator of api jikan
[Jikan fork](https://github.com/LucasHenriqueDiniz/jikan-edge) - Fork of api jikan made public
