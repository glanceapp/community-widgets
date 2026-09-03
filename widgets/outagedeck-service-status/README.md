# Cloud & SaaS Status

![Cloud and SaaS status widgets](./preview.png)

A keyless `custom-api` widget that shows a provider's current status, active
incident count, service-level states, and source freshness. The example uses
GitHub; change `github` in both URLs to another [OutageDeck provider
slug](https://outagedeck.com/providers) to monitor a different dependency.

```yaml
  - type: custom-api
    title: GitHub status
    title-url: https://outagedeck.com/providers/github?utm_source=glance&utm_medium=community_widget&utm_campaign=glance_community_widget
    url: https://outagedeck.com/api/v1/providers/github
    cache: 10m
    template: |
      {{ $status := .JSON.String "data.currentStatus.code" }}
      <div class="flex justify-between items-center">
        <div class="size-h3 {{ if eq $status "operational" }}color-positive{{ else if or (eq $status "maintenance") (eq $status "unknown") }}color-primary{{ else }}color-negative{{ end }}">
          {{ .JSON.String "data.currentStatus.label" }}
        </div>
        <div class="size-h5 color-subdue">
          {{ .JSON.Int "data.counts.activeIncidents" }} active incidents
        </div>
      </div>
      <p class="margin-top-10">{{ .JSON.String "data.currentStatus.headline" }}</p>
      <ul class="list list-gap-4 margin-top-10">
      {{ range .JSON.Array "data.services" }}
        {{ $serviceStatus := .String "status" }}
        <li class="flex justify-between">
          <span>{{ .String "name" }}</span>
          <span class="{{ if eq $serviceStatus "operational" }}color-positive{{ else if or (eq $serviceStatus "maintenance") (eq $serviceStatus "unknown") }}color-primary{{ else }}color-negative{{ end }}">
            {{ if eq $serviceStatus "operational" }}Operational{{ else if eq $serviceStatus "degraded" }}Degraded{{ else if eq $serviceStatus "partial_outage" }}Partial outage{{ else if eq $serviceStatus "major_outage" }}Major outage{{ else if eq $serviceStatus "maintenance" }}Maintenance{{ else }}Unknown{{ end }}
          </span>
        </li>
      {{ end }}
      </ul>
      <p class="margin-top-10 size-h5 color-subdue">
        Source checked <span {{ .JSON.String "data.source.checkedAt" | parseTime "rfc3339" | toRelativeTime }}></span>
      </p>
```

The 10-minute cache makes six requests per hour for each copy of the widget.
OutageDeck's anonymous API allowance is 120 requests per hour, so one dashboard
can run up to 20 copies at that interval. No API key or credential is required.

I'm Kerolos, founder of [OutageDeck](https://outagedeck.com/), and I maintain
this widget.
