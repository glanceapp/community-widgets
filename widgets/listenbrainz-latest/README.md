![](preview.png)

```yaml
- type: custom-api
  title: Listening to
  cache: 5m
  frameless: true
  url: https://api.listenbrainz.org/1/user/${LISTENBRAINZ_USER}/listens?count=1
  template: |
    {{ $listen := .JSON.Get "payload.listens.0" }}
    {{ $mbid := $listen.String "track_metadata.mbid_mapping.release_mbid" }}

    <a href="https://listenbrainz.org/user/${LISTENBRAINZ_USER}" rel="noreferer noopener" class="widget-content-frame container-lb">
      {{ if $mbid }}
        {{ $cover := newRequest (print "https://coverartarchive.org/release/" $mbid) | getResponse }}
        <div class="cover-wrapper-lb">
          <img class="cover-img-lb" src="{{ $cover.JSON.String "images.0.thumbnails.small" }}" />
        </div>
      {{ end }}

      <div style="padding: var(--widget-content-padding)">
        <p class="color-highlight">{{ $listen.String "track_metadata.track_name" }}</p>
        <p>{{ $listen.String "track_metadata.artist_name" }}</p>
      </div>
    </a>

    <style>
      .container-lb {
         display: grid;
         grid-template-columns: auto 1fr;
         overflow: hidden;
      }
      .cover-wrapper-lb {
         position: relative;
         height: 100%;
         aspect-ratio: 1 / 1;
      }
      .cover-img-lb {
         position: absolute;
         top: 0;
         left: 0;
         width: 100%;
         height: 100%;
      }
    </style>
```

## Environment variables

- `LISTENBRAINZ_USER` - the username of the ListenBrainz user
