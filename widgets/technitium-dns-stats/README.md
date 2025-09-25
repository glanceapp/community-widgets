# Technitium DNS Stats

This Custom API implementation supports pulling basic stats from [Technitium DNS Server](https://technitium.com/dns/) based off a set timeframe:
- Current version
- Total queries
- Blocked queries
- Top clients table
- Top queries table
- Top blocked queries table

![Collapsed widget](preview.png)

![Table example](preview_table.png)

It will also check for updates to Technitium and display a notificaiton when an update is available.

## Requirements

- Technitium service account with appropriate permissions (I have my account in the "DNS Administrators" group)
- API token for said account (Create by clicking your account icon in the top-right of Technitium -> "Create API Token")

## Environment Variables

- `TECHNITIUM_IP_PORT`: The accessible IP of your Technitium instance. Default port is 5380.
- `TECHNITIUM_TOKEN`: Technitium API token.
- `TECHNITIUM_STATS_TIMEFRAME`: Timeframe to fetch stats. Currently supports the following values: `LastHour`, `LastDay`, `LastWeek`, `LastMonth`, `LastYear`.

## Configuration

```
- type: custom-api
    title: Technitium DNS Stats
    cache: 1h
    subrequests:
        update-check:
            url: http://${TECHNITIUM_IP_PORT}/api/user/checkForUpdate?token=${TECHNITIUM_TOKEN}
    url: http://${TECHNITIUM_IP_PORT}/api/dashboard/stats/get?token=${TECHNITIUM_TOKEN}&type=${TECHNITIUM_STATS_TIMEFRAME}&utc=true
    method: GET
    template: |
    <ul class="list list-gap-10 collapsible-container" data-collapse-after="3">
    {{ $updateCheck := .Subrequest "update-check" }}
        <li>
        <strong class="color-highlight">Version:</strong>
        {{ $updateCheck.JSON.String "response.currentVersion" }}
        {{ if ($updateCheck.JSON.Bool "response.updateAvailable") }}
        <a href="https://github.com/TechnitiumSoftware/DnsServer/blob/master/CHANGELOG.md"><strong><span class="color-positive"> Update Available: {{ $updateCheck.response.updateVersion }}</span></strong></a>
        {{ end }}
        </li>
        <li>
        <strong class="color-highlight">Total Queries:</strong>
        {{ .JSON.String "response.stats.totalQueries" }}
        </li>
        <li>
        <strong class="color-highlight">Blocked Queries:</strong>
        {{ .JSON.String "response.stats.totalBlocked" }}
        </li>
        <li>
        <!-- Top Clients Table -->
        <h4 class="size-h5 color-highlight">Top Clients</h4>
        <table style="width: 100%; border-collapse: collapse; margin-bottom: 20px;">
            <thead>
            <tr style="border-bottom: 1px solid #ccc;">
                <th style="text-align: left; padding: 8px;">Name</th>
                <th style="text-align: left; padding: 8px;">Domain</th>
                <th style="text-align: left; padding: 8px;">Hits</th>
                <th style="text-align: left; padding: 8px;">Rate Limited</th>
            </tr>
            </thead>
            <tbody>
            {{ range $index, $client := .JSON.Array "response.topClients" }}
                <tr>
                <td style="padding: 8px;">{{ $client.String "name" }}</td>
                <td style="padding: 8px;">{{ $client.String "domain" }}</td>
                <td style="padding: 8px;">{{ $client.String "hits" }}</td>
                <td style="padding: 8px;">
                    {{ if $client.Bool "rateLimited" }}Yes{{ else }}No{{ end }}
                </td>
                </tr>
            {{ end }}
            </tbody>
        </table>
        </li>
        <li>
        <!-- Top Domains Table -->
        <h4 class="size-h5 color-highlight">Top Domains</h4>
        <table style="width: 100%; border-collapse: collapse; margin-bottom: 20px;">
            <thead>
            <tr style="border-bottom: 1px solid #ccc;">
                <th style="text-align: left; padding: 8px;">Domain</th>
                <th style="text-align: left; padding: 8px;">Hits</th>
            </tr>
            </thead>
            <tbody>
            {{ range $index, $domain := .JSON.Array "response.topDomains" }}
                <tr>
                <td style="padding: 8px;">{{ $domain.String "name" }}</td>
                <td style="padding: 8px;">{{ $domain.String "hits" }}</td>
                </tr>
            {{ end }}
            </tbody>
        </table>
        </li>
        <li>
        <!-- Top Blocked Domains Table -->
        <h4 class="size-h5 color-highlight">Top Blocked Domains</h4>
        <table style="width: 100%; border-collapse: collapse;">
            <thead>
            <tr style="border-bottom: 1px solid #ccc;">
                <th style="text-align: left; padding: 8px;">Domain</th>
                <th style="text-align: left; padding: 8px;">Hits</th>
            </tr>
            </thead>
            <tbody>
            {{ range $index, $blocked := .JSON.Array "response.topBlockedDomains" }}
                <tr>
                <td style="padding: 8px;">{{ $blocked.String "name" }}</td>
                <td style="padding: 8px;">{{ $blocked.String "hits" }}</td>
                </tr>
            {{ end }}
            </tbody>
        </table>
        </li>
    </ul>
```
