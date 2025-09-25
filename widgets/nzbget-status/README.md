<img src=preview.png width="400">

```yaml
- type: custom-api
  title: NZBGet Status
  cache: 5s
  url: ${NZBGET_BASE_URL}/jsonrpc/status
  options:
    api-base-url: ${NZBGET_BASE_URL} # base API URL for NZBGet
    queue-enabled: true              # enable or disable the queue
    queue-items: 3                   # number of queue items to show by default
    slow-speed: 0                    # under this KB/s, speed shows negative color. 0 to disable
  template: |

    {{ $apiBaseUrl := .Options.StringOr "api-base-url" "" }}
    {{ $queueEnabled := .Options.BoolOr "queue-enabled" false }}
    {{ $queueItems := .Options.IntOr "queue-items" 3 }}

    <style>
      .progress-container-nzbget {
        border: 2px solid var(--color-background);
        outline: 1px solid var(--color-text-subdue);
        position: relative;
        border-radius: 1rem;
        height: 1rem;
        width: 100%;
        margin-top: 0.2rem;
      }
      .progress-bar-nzbget {
        background: var(--color-positive);
        border-radius: 1rem;
        position: absolute;
        top: 0;
        bottom: 0;
        left: 0;
      }
      .queue-item-nzbget {
        background: var(--color-widget-background-highlight);
        border-radius: 1rem;
        line-height: 1.8rem;
        margin-bottom: 1rem;
        padding: 0.5rem 1rem;
      }
      .queue-item-title-nzbget {
        overflow-x: hidden;
        white-space: nowrap;
        text-overflow: clip;
      }
    </style>

    {{ $truncateFloat := `\d+\.\d\d` }}
    {{ $standby := .JSON.Bool "result.ServerStandBy" }}
    {{ $paused  := .JSON.Bool "result.DownloadPaused" }}

    {{ $rateB     := .JSON.Float "result.DownloadRate" }}
    {{ $rateKB    := div $rateB 1024 }}
    {{ $rateMB    := div $rateKB 1024 }}
    {{ $rateKBStr := $rateKB | formatNumber | findMatch $truncateFloat }}
    {{ $rateMBStr := $rateMB | formatNumber | findMatch $truncateFloat }}

    {{ $remainingMB := .JSON.Float "result.RemainingSizeMB" }}
    {{ $remainingGB := div $remainingMB 1024 | formatNumber | findMatch $truncateFloat  }}

    {{ $timeLeftTotalSecs := div $remainingMB $rateMB }}
    {{ $timeLeftDays      := div $timeLeftTotalSecs 86400 | toInt }}
    {{ $timeLeftTotalSecs  = mod (toInt $timeLeftTotalSecs) 86400 | toFloat }}
    {{ $timeLeftHours     := div $timeLeftTotalSecs 3600 | toInt }}
    {{ $timeLeftTotalSecs  = mod (toInt $timeLeftTotalSecs) 3600 | toFloat }}
    {{ $timeLeftMins      := div $timeLeftTotalSecs 60 | toInt }}
    {{ $timeLeftSecs      := mod (toInt $timeLeftTotalSecs) 60 }}

    {{ $speedColor        := "color-highlight" }}
    {{ $threshold         := .Options.IntOr "slow-speed" 0 }}
    {{ $thresholdEnabled  := gt $threshold 0}}
    {{ $thresholdExceeded := gt (toInt $rateKB) $threshold }}

    {{ if and $thresholdEnabled (not $thresholdExceeded) }}
      {{ $speedColor = "color-negative" }}
    {{ end }}

    <div>
      {{ if $paused }}
          <p class="size-h5 text-center">Downloads are paused.</p>
      {{ else if $standby }}
          <p class="size-h5 text-center">No downloads queued.</p>
      {{ else }}
          <ul class="flex justify-between mb-2 text-center">
            <li class="flex-1">
              <div class="{{ $speedColor }} size-h3">
                {{ if lt (toInt $rateKB) 1024 }}
                  {{ $rateKBStr }}<span class="size-h6">KB/s</span>
                {{ else }}
                  {{ $rateMBStr }}<span class="size-h6">MB/s</span>
                {{ end }}
              </div>
              <div class="size-h6">SPEED</div>
            </li>
            <li class="flex-1">
              <div class="color-highlight size-h3 flex justify-center"  style="gap: 0.3rem;">
                {{ if gt $timeLeftDays 0 }}
                  <div class="flex items-end">
                    <span class="size-h3">{{ $timeLeftDays }}</span>
                    <span class="size-h6" style="margin-bottom: 1px">d</span>
                  </div>
                {{ end }}
                {{ if gt $timeLeftHours 0 }}
                  <div class="flex items-end">
                    <span class="size-h3">{{ $timeLeftHours }}</span>
                    <span class="size-h6" style="margin-bottom: 1px">h</span>
                  </div>
                {{ end }}
                {{ if and (gt $timeLeftMins 0) (eq $timeLeftDays 0) }}
                  <div class="flex items-end">
                    <span class="size-h3">{{ $timeLeftMins }}</span>
                    <span class="size-h6" style="margin-bottom: 1px">m</span>
                  </div>
                {{ end }}
                {{ if and (eq $timeLeftHours 0) (eq $timeLeftMins 0) }}
                  <div class="flex items-end">
                    <span class="size-h3">{{ $timeLeftSecs }}</span>
                    <span class="size-h6" style="margin-bottom: 1px">s</span>
                  </div>
                {{ end }}
              </div>
              <div class="size-h6">ETA</div>
            </li>
            <li class="flex-1">
              <div class="flex justify-center color-highlight" style="align-items: baseline;">
                {{ if lt (toInt $remainingMB) 1024 }}
                  <span class="size-h3">{{ $remainingMB }}</span>
                  <span class="size-h6">MB</span>
                {{ else }}
                  <span class="size-h3">{{ $remainingGB }}</span>
                  <span class="size-h6">GB</span>
                {{ end }}
              </div>
              <div class="size-h6">REMAINING</div>
            </li>
          </ul>

      {{ end }}

      {{ if $queueEnabled }}
        {{ $groups := newRequest (concat $apiBaseUrl "/jsonrpc/listgroups") | getResponse }}
        {{ $result := $groups.JSON.Array "result" }}
        {{ $collapsible := gt (len $result) $queueItems }}

        {{ if gt (len $result) 0 }}
          {{ if $collapsible }}
            <p class="size-h6 color-base text-center" style="margin: 1rem 0 -1rem;">{{ len $result }} items in queue</p>
          {{ end }}
          <ul class="list collapsible-container" data-collapse-after="{{ $queueItems }}" style="margin: 2rem 0 {{ if $collapsible }}-1rem{{ else }}0{{ end }};">
          {{ range $result }}
            <li class="queue-item-nzbget">
              <p class="queue-item-title-nzbget size-h6">{{ .String "NZBName" }}</p>
              {{ template "status" . }}
            </li>
          {{ end }}
          </ul>
        {{ end }}
      {{ end }}

    </div>

    {{ define "status" }}
      {{ $status := .String "Status" }}
      <p class="size-h6 flex flex-column">
        {{ if eq $status "QUEUED" }}
          <span class="color-subdue" title="Queued for download">Queued ({{ .Int "FileSizeMB" }} MB)</span>
        {{ else if eq $status "PAUSED" }}
          <span>Paused</span>
        {{ else if eq $status "DOWNLOADING" }}
          {{ $downloaded := .Float "DownloadedSizeMB" }}
          {{ $total := .Float "FileSizeMB" }}
          {{ $completed := div $downloaded $total | mul 100 }}
          <span class="color-highlight" title="Item is being downloaded">Downloading ({{ $downloaded }}/{{ $total }} MB)</span>
          <div class="progress-container-nzbget">
            <div class="progress-bar-nzbget" style="width: {{ $completed }}%;"></div>
          </div>
        {{ else if eq $status "FETCHING" }}
          <span title="NZB-file is being fetched from URL (Kind=URL)">Fetching</span>
        {{ else if eq $status "PP_QUEUED" }}
          <span title="queued for post-processing (completely downloaded)">Queued (Post-processing) </span>
        {{ else if eq $status "LOADING_PARS" }}
          <span title="Stage of par-check">Loading pars</span>
        {{ else if eq $status "VERIFYING_SOURCES" }}
          <span title="Stage of par-check">Verifying sources</span>
        {{ else if eq $status "REPAIRING" }}
          <span title="Stage of par-check">Repairing</span>
        {{ else if eq $status "VERIFYING_REPAIRED" }}
          <span title="Stage of par-check">Verifying repaired</span>
        {{ else if eq $status "RENAMING" }}
          <span title="Processed by par-renamer">Renaming</span>
        {{ else if eq $status "UNPACKING" }}
          <span title="Being unpacked">Unpacking</span>
        {{ else if eq $status "MOVING" }}
          <span title="Moving files from intermediate directory into destination directory">Moving </span>
        {{ else if eq $status "EXECUTING_SCRIPT" }}
          <span title="Executing post-processing script">Post-processing</span>
        {{ else if eq $status "PP_FINISHED" }}
          <span title="Post-processing is finished, the item is about to be moved to history.">Post-processing finished</span>
        {{ else }}
          Unknown status
        {{ end }}
      </p>
    {{ end }}
```

### Environment variables

- `NZBGET_BASE_URL` - the base URL for NZBGet that the widget can access. It may require authentication which you can provide in a format like `http://username:password@server:port`. See the [NZBGet API documentation](https://nzbget.net/api/) for more details.

### User options

| Option           | Default  | Description                                                                               |
| -----------------| ---------|------------------------------------------------------------------------------------------ |
| queue-enabled    | true     | Displays the queue. Disable it for a more minimal look.                                   |
| queue-items      | 3        | Shows this many queue items above the fold.                                               |
| slow-speed       | 0        | Below this speed (in KB/s), uses negative color to warn about slow download. 0 to disable.|

###### Minimal look with `queue-enabled: false`

<img src="preview_minimal.png" width="400">

###### Slow speed warning with `slow-speed: 100`

<img src="preview_slow_speed.png" width="400">
