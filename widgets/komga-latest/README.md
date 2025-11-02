## Description

This widget displays recently added/updated series of a Komga library as cards in horizontal carousel:

![Preview of the widget](preview.png)

The widget requires modifications to the configuration file in order to work. Under the `options` field in the yaml configuration file, one needs to supply Komga URL, API key and a library ID of a relevant library. One can either supply them directly or as environment variables. 

## List of Options

- (required) `base-url` requires URL pointing to your Komga, be sure to include "http://" or "https://"
- (required) `api-key` requires Komga API key that can be retrieved in Komga under "My Account -> API Keys -> +"
- (required) `library-id` requires ID of you Komga library that can be retrieved from library URL: ${KOMGA_URL}/libraries/${LIBRARY_ID}/series
- (optional) `items` requires integer number that will determine maximum amount of cards. The default value is 10
- (optional) `mode` determines whether the widget show recently added series or the series that recently had changes. Requires "new" or "updated". The default value is "new"

## Widget Configuration YAML
```yaml
- type: custom-api
  title: Recent Series                                        
  frameless: true
  cache: 15m
  options:                                                   
    base-url: ${KOMGA_URL}                                   # URL pointing to your Komga
    api-key: ${KOMGA_API_KEY}                                # retrieve from Komga in: My Account -> API Keys -> +
    library-id: ${LIBRARY_ID}                                # id of the Komga library (can be spotted in the url of that library)
    items: 10                                                # max number of cards to show in the widgets
    mode: new                                                # only "new" or "updated" are allowed
  template: |
    {{ $baseURL := .Options.StringOr "base-url" "" }}
    {{ $apiKey := .Options.StringOr "api-key" "" }}
    {{ $libraryId := .Options.StringOr "library-id" "" }}
    {{ $numberOfItems := .Options.IntOr "items" 10 }}
    {{ $mode := .Options.mode "library-id" "new" }}

    {{
      $auth := newRequest (concat $baseURL "/api/v1/login")
        | withHeader "X-API-Key" $apiKey
        | getResponse
    }}

    {{
      $content := newRequest (concat $baseURL "/api/v1/series/" $mode "?library_id=" $libraryId)
        | withHeader "X-API-Key" $apiKey
        | getResponse
    }}

    {{ $session := $auth.Response.Header.Get "Set-Cookie" }}
    
    {{ if or (eq $baseURL "") (eq $apiKey "") }}
        <div class="widget-content-frame" style="flex=0 0 25vh; display:flex">
          <div class="grow padding-inline-widget margin-top-10 margin-bottom-10 color-negative">
                Komga base-url or API token not set.
          </div>
        </div>
    {{ else }}
        <div class="cards-horizontal carousel-items-container">
          {{ $arr := $content.JSON.Array "content" }}
          {{ $len := len $arr }}
          {{ $shown := 0 }}
          
          {{ range $i, $_ := $arr }}
            {{ if lt $shown $numberOfItems }}
              {{ $el := index $arr $i }}
              
              <a class="card widget-content-frame"
                 href="{{ $baseURL }}/series/{{ $el.String "id" }}"
                 style="flex:0 0 25vh;min-width:170px; min-height: 150px; display:flex;flex-direction:column;box-sizing:border-box;text-decoration:none;color:inherit;">
                <div style="position: relative;">
                    <img src="{{ $baseURL }}/api/v1/series/{{ $el.String "id" }}/thumbnail"
                         alt="{{ $el.String "metadata.title" }}"
                         loading="lazy"
                         class="media-server-thumbnail shrink-0 loaded finished-transition"
                         style="width:100%; height: 37vh; min-height: 250px; display:block;border-radius:var(--border-radius) var(--border-radius) 0 0;object-fit:cover;" 
                    />
                </div>
                <div class="grow padding-inline-widget margin-top-10 margin-bottom-10">
                  <ul class="flex flex-column justify-evenly margin-bottom-3" style="height:100%; gap:4px;">
                    <li class="text-truncate color-primary" style="overflow:hidden;text-overflow:ellipsis;white-space:nowrap;" title="{{ $el.String "title" }}">
                      {{ $el.String "metadata.title" }} <!-- ({{ $el.String "metadata.releaseDate" | parseTime "2006-01-02" | formatTime "2006" }}) -->
                    </li>
                    <li style="font-size:0.85em;opacity:0.7;white-space:nowrap;overflow:hidden;text-overflow:ellipsis;">
                      Added: {{ $el.String "created" | parseTime "RFC3339" | formatTime "Jan 2, 2006" }}
                    </li>
                  </ul>
                </div>
              </a>
              {{ $shown = add $shown 1 }}
            {{ end }}
          {{ end }}
        </div>
    {{ end }}
```

