![](preview.png)

```yaml
- type: custom-api
  title: Gatus
  cache: 1m
  url: ${GATUS_URL}/api/v1/endpoints/statuses
  template: |
    {{ $up := 0 }}
    {{ $down := 0 }}
    {{ $totalChecks := 0 }}
    {{ $successfulChecks := 0 }}
    {{ range .JSON.Array "" }}
      {{ $results := .Array "results" }}
      {{ $totalChecks = add $totalChecks (len $results) }}
      {{ range $results }}
        {{ if .Bool "success" }}
          {{ $successfulChecks = add $successfulChecks 1 }}
        {{ end }}
      {{ end }}
      {{ if gt (len $results) 0 }}
        {{ $latestResult := index $results (sub (len $results) 1) }}
        {{ if $latestResult.Bool "success" }}
          {{ $up = add $up 1 }}
        {{ else }}
          {{ $down = add $down 1 }}
        {{ end }}
      {{ end }}
    {{ end }}
    {{ $uptime := 0 }}
    {{ if gt $totalChecks 0 }}
      {{ $uptime = div (mul $successfulChecks 1000) $totalChecks }}
    {{ end }}

    {{ define "popover" }}
      <div data-popover-html>
        <div style="margin: 5px;">
          <p class="margin-top-12" style="overflow-y: auto; text-align: justify; max-height: 20rem;">
            {{ range .Array "conditionResults" }}
              {{ if (.Bool "success") }}
                <p class="color-positive">{{ .String "condition" }}</p>
              {{ else }}
                <p class="color-negative">{{ .String "condition" }}</p>
              {{ end }}
            {{ end }}
          </p>
        </div>
      </div>
    {{ end }}

    <div class="widget-small-content-bounds">
      <div style="display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 10px; text-align: center;">
        <div>
          <p class="size-h3 color-highlight">{{ $up }}</p>
          <p class="size-h6">SITES UP</p>
        </div>
        <div>
          <p class="size-h3 color-highlight">{{ $down }}</p>
          <p class="size-h6">SITES DOWN</p>
        </div>
        <div>
          <p class="size-h3 color-highlight">{{ div $uptime 10.0 | printf "%.1f" }}%</p>
          <p class="size-h6">UPTIME</p>
        </div>
      </div>
      <ul class="margin-top-10 list list-gap-10 collapsible-container" data-collapse-after="5">
        {{ range .JSON.Array "" }}
          {{ $results := .Array "results" }}
          {{ $latestResult := index $results (sub (len $results) 1) }}
          {{ if not ($latestResult.Bool "success") }}
              <li class="flex items-center justify-between">
                <p class="size-h3 color-highlight">{{ .String "name" }}</p>
                <div data-popover-type="html" data-popover-position="above" data-popover-show-delay="500" style="align-content: center;">
                  {{ template "popover" $latestResult }}
                  <svg width="20" height="20" fill="var(--color-negative)" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M8.485 2.495c.673-1.167 2.357-1.167 3.03 0l6.28 10.875c.673 1.167-.17 2.625-1.516 2.625H3.72c-1.347 0-2.189-1.458-1.515-2.625L8.485 2.495ZM10 5a.75.75 0 0 1 .75.75v3.5a.75.75 0 0 1-1.5 0v-3.5A.75.75 0 0 1 10 5Zm0 9a1 1 0 1 0 0-2 1 1 0 0 0 0 2Z" clip-rule="evenodd" />
                  </svg>
                </div>
              </li>
          {{ end }}
        {{ end }}
        {{ range .JSON.Array "" }}
          {{ $results := .Array "results" }}
          {{ $latestResult := index $results (sub (len $results) 1) }}
          {{ if $latestResult.Bool "success" }}
              <li class="flex items-center justify-between">
                <p class="size-h3 color-highlight">{{ .String "name" }}</p>
                <div data-popover-type="html" data-popover-position="above" data-popover-show-delay="500" style="align-content: center;">
                  {{ template "popover" $latestResult }}
                  <svg width="20" height="20" fill="var(--color-positive)" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M10 18a8 8 0 1 0 0-16 8 8 0 0 0 0 16Zm3.857-9.809a.75.75 0 0 0-1.214-.882l-3.483 4.79-1.88-1.88a.75.75 0 1 0-1.06 1.061l2.5 2.5a.75.75 0 0 0 1.137-.089l4-5.5Z" clip-rule="evenodd" />
                  </svg>
                </div>
              </li>
          {{ end }}
        {{ end }}
      </ul>
    </div>

```

## Environment variables

- `GATUS_URL`: The url of your Gatus instace (no slash on the end). For example, `http://gatus:8080`
