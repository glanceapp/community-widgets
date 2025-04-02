
![](preview.png)

## Using RSS Bridge
This only supports 1 channel unless you use `By playlist` or `By search result` context or if you manage to merge the `items` array.
```yml
- type: custom-api
  title: Youtube
  url: ${RSS_BRIDGE}
  parameters:
    action: display
    bridge: YoutubeBridge
    context: By username
    u: LinusTechTips
    duration_min: 3
    duration_max: nil
    format: Json
  frameless: true
  headers:
    content-type: application/json
  cache: 10m
  template: |
    {{ $arrayList := .JSON.Array "items" }}
    {{ $closeOutsideIframe := false }} {{/* Closes the popup anywhere around the iframe */}}
    {{ $popupWidth := "85%" }}
    {{ $popupHeight := "85%" }}
    {{ $author := .JSON.String "title" | trimSuffix " - YouTube" }}
    {{ $authorUrl := .JSON.String "home_page_url" | trimSuffix "/videos" }}
    <div class="cards-grid collapsible-container" data-collapse-after-rows="3">
      {{ if eq (len ($arrayList)) 0 }}
      <div>Nothing to show</div>
      {{ else }}
        {{ range $arrayList }}

        {{ $title := .String "title" }}
        {{ $youtubeUrl := .String "url" }}
        {{ $youtubeId := .String "_rssbridge.id" }}
        {{ $youtubeEmbedUrl := concat "https://www.youtube-nocookie.com/embed/" $youtubeId }}
        {{ $timestamp := .String "date_modified" | parseTime "rfc3339" }}
        {{ $thumbnail := findSubmatch "img src=\"([^\"]+)" (.String "content_html") }}
        
        <div class="card widget-content-frame thumbnail-parent">
          <div onclick="
            const classPrefix = 'popup-youtube-';
            const embedClass = classPrefix + 'embed';
            const dimBgClass = classPrefix + 'dim';
            const closeBtnClass = classPrefix + 'close';
            const popupEmbedExistingDiv = document.querySelector('.' + embedClass);
            if (!popupEmbedExistingDiv) {
              const dimBackground = document.createElement('div');
              dimBackground.className = dimBgClass;
              dimBackground.style = 'position: fixed; top: 0; left: 0; width: 100%; height: 100%; background: rgba(0, 0, 0, 0.5);';c
              document.body.appendChild(dimBackground);
              
              const popupEmbedDiv = document.createElement('div');
              popupEmbedDiv.className = embedClass;
              popupEmbedDiv.style = 'position: fixed; top: 50%; left: 50%; transform: translate(-50%, -50%); width: {{ $popupWidth }}; height: {{ $popupHeight }};';

              const popupIframe = document.createElement('iframe');
              popupIframe.className = classPrefix + 'iframe';
              popupIframe.src = {{ $youtubeEmbedUrl }};
              popupIframe.loading = 'lazy';
              popupIframe.allowFullscreen = true;
              popupIframe.style ='outline: 2px solid var(--color-primary); outline-offset: -0.1rem; border-radius: 20px; width: 100%; height: 100%; border: none; background-color: var(--color-popover-background);';
              
              popupEmbedDiv.appendChild(popupIframe);
              document.body.appendChild(popupEmbedDiv);

              const closeBtn = document.createElement('span');
              closeBtn.className = closeBtnClass;
              closeBtn.innerText = '&times';
              closeBtn.style = 'cursor: pointer; font-size: 3rem; width: 30px; height: 30px; display: flex; align-items: center; justify-content: center; position: absolute; top: 3rem; right: 3rem;';
              document.body.appendChild(closeBtn);
              document.body.style.overflowY = 'hidden';
              
              const closePopup = () => {
                popupEmbedDiv.remove();
                dimBackground.remove();
                closeBtn.remove();
                document.body.style.overflowY = 'scroll';
              }

              {{ if $closeOutsideIframe }} dimBackground.onclick = closePopup; {{ end }}
              closeBtn.onclick = closePopup;
            } else {
              document.querySelector('.' + dimBgClass).remove();
              document.querySelector('.' + closeBtnClass).remove();
              document.body.style.overflowY = 'scroll';
              popupEmbedExistingDiv.remove();
            }
          " style="cursor: pointer;">
            <span>
              {{ if ne $thumbnail "" }}
              <img class="video-thumbnail thumbnail" loading="lazy" src="{{ $thumbnail }}" alt="">
              {{ else }}
              <svg class="video-thumbnail thumbnail" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="var(--color-text-subdue)">
                <path stroke-linecap="round" stroke-linejoin="round" d="m2.25 15.75 5.159-5.159a2.25 2.25 0 0 1 3.182 0l5.159 5.159m-1.5-1.5 1.409-1.409a2.25 2.25 0 0 1 3.182 0l2.909 2.909m-18 3.75h16.5a1.5 1.5 0 0 0 1.5-1.5V6a1.5 1.5 0 0 0-1.5-1.5H3.75A1.5 1.5 0 0 0 2.25 6v12a1.5 1.5 0 0 0 1.5 1.5Zm10.5-11.25h.008v.008h-.008V8.25Zm.375 0a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Z" />
              </svg>
              {{ end }}
            </span>
          </div>
          <div class="margin-bottom-widget padding-inline-widget flex flex-column grow">
            <div class="margin-top-10 margin-bottom-auto">
              <a href="{{ $youtubeUrl }}#" class="color-primary-if-not-visited text-truncate-2-lines" target="_blank" rel="noreferrer">
                {{ $title }}
              </a>
            </div>
            <ul class="list-horizontal-text flex-nowrap margin-top-7">
              <li class="shrink-0" {{ $timestamp | toRelativeTime }}></li>
              <li class="min-width-0">
                <a href="{{ $authorUrl }}" class="color-primary-if-not-visited text-truncate-2-lines" target="_blank" rel="noreferrer">
                  {{ $author }}
                </a>
              </li>
            </ul>
          </div>
        </div>
        {{ end }}
      {{ end }}
    </div>
```

Public Instance list: https://rss-bridge.github.io/rss-bridge/General/Public_Hosts.html

## FreshRSS
If you're using [FreshRSS](https://github.com/FreshRSS/FreshRSS) as a backend then replace the top part with

```yml
- type: custom-api
  url: ${FRESHRSS_URL}/api/query.php
  parameters:
    user: yourusername
    f: greader
    t: seeBelowOnHowToGenerateOne
  frameless: true
  title: Youtube
  css-class: widget-no-title
  headers:
    content-type: application/json
  cache: 30m
  template: |
    {{ $arrayList := .JSON.Array "items" }}
    {{ $closeOutsideIframe := false }} {{/* Closes the popup anywhere around the iframe */}}
    {{ $popupWidth := "85%" }}
    {{ $popupHeight := "85%" }}
    <div class="cards-grid collapsible-container" data-collapse-after-rows="3">
      {{ if eq (len ($arrayList)) 0 }}
      <div>Nothing to show</div>
      {{ else }}
        {{ range $arrayList }}

        {{ $title := .String "title" }}
        {{ $author := .String "origin.title" }}
        {{ $authorUrl := .String "origin.htmlUrl" | trimSuffix "/videos" }}
        {{ $youtubeUrl := .String "canonical.0.href" }}
        {{ $youtubeId := $youtubeUrl | trimPrefix "https://www.youtube.com/watch?v=" }}
        {{ $youtubeEmbedUrl := concat "https://www.youtube-nocookie.com/embed/" $youtubeId }}
        {{ $timestamp := .String "published" | parseTime "unix" }}
        {{ $thumbnail := findSubmatch "img src=\"([^\"]+)" (.String "content.content") }}
        
        <div class="card widget-content-frame thumbnail-parent">
```

### To generate User Queries
1. go to https://your-freshrss-domain.com/i/
2. Select a Category on the left side
3. Select the bookmark icon and choose Bookmark current query see screenshot
4. This will create a query Query n°1
5. Enable sharing by HTML & RSS then Submit
6. Copy `Shareable link to the GReader JSON`

## YouTube embed proxy
You can replace `https://www.youtube-nocookie.com/embed/` with your instance that supports embedding.