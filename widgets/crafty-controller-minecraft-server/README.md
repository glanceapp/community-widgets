## Previews

#### When Server is offline:
![Server offline preview](./Preview_offline.png)

#### When server is starting up:
![Server starting preview](./Preview_starting.png)

#### When server is running:
![Server running preview](./Preview_online.png)

Displays server's custom icon, inferred IP, MOTD message and some statistics. The display of MOTD can be disabled in the options section of the template, and the widget will look like this:

![Server running with disabled MOTD preview](./Preview_online_disabledMOTD.png)

## Yaml
```yaml
- type: custom-api
  title: Crafty Controller - Minecraft Server
  cache: 5s
  allow-insecure: true
  url: ${CRAFTY_URL}/api/v2/servers/${CRAFTY_SERVER_ID}/stats
  headers:
    Authorization: Bearer ${CRAFTY_API_TOKEN}
    Accept: application/json
  options:
    display-MOTD: true
  template: |
    {{/* Optional config options */}}
    {{ $displayMOTD := .Options.BoolOr "display-MOTD" true }}

    {{/* Use the built-in parsed JSON response */}}
    {{ $is_running := .JSON.Bool "data.running" }}
    {{ $online_players := .JSON.Int "data.online" | formatNumber }}
    {{ $max_players := .JSON.Int "data.max" | formatNumber }}
    {{ $name := .JSON.String "data.world_name" }}
    {{ $size := .JSON.String "data.world_size" }}
    {{ $version := .JSON.String "data.version" }}
    {{ $icon := .JSON.String "data.icon" }}
    {{ $server_ip := .JSON.String "data.server_id.server_ip" }}
    {{ $server_port := .JSON.String "data.server_id.server_port" }}
    {{ $motd := .JSON.String "data.desc" }}

    {{ $server_addr := "" }}
    {{ if and ($is_running) (eq $server_ip "127.0.0.1") }}
      {{ $server_addr = printf "%s:%s" (replaceMatches "https?://" "" "${CRAFTY_URL}") $server_port }}
    {{ else if $is_running }}
      {{ $server_addr = printf "%s:%s" $server_ip $server_port }}
    {{ end }}

    {{ $starting := false }}
    {{ if and ($is_running) (eq $max_players "0") (eq $version "False") }}
      {{ $starting = true }}
    {{ end }}

    {{/* Corrected boolean states */}}
    {{ $updating := .JSON.Bool "data.updating" }}
    {{ $importing := .JSON.Bool "data.importing" }}
    {{ $crashed := .JSON.Bool "data.crashed" }}

    <div style="display:flex; align-items:center; gap:12px;">
      <!-- Server Icon -->
      <div style="width:40px; height:40px; flex-shrink:0; border-radius:4px; display:flex; justify-content:center; align-items:center; overflow:hidden;">
        {{ if eq $icon "" }}
          <img src="https://cdn.jsdelivr.net/gh/selfhst/icons/png/minecraft.png" style="width:100%; height:100%; object-fit:contain;">
        {{ else }}
          <img src="data:image/png;base64, {{ $icon }}" style="width:100%; height:100%; object-fit:contain;">
        {{ end }}
      </div>

      <!-- Server Info -->
      <div style="display:flex; flex-direction:column;">
        <div style="display:flex; align-items:center; gap:6px;">
          <span class="size-h4 block text-truncate color-primary">{{ $name }}</span>
          {{ if and ($is_running) (not $starting) (ne $server_addr "") }}
            <span class="size-h6 color-secondary">- {{ $server_addr }}</span>
          {{ end }}
        </div>

        {{ if and ($is_running) (not $starting) }}
          {{ if and (ne $motd "") ($displayMOTD) }}
            <div style="font-size:0.9em; color:var(--color-secondary);">
              {{ replaceMatches "§." "" $motd }}
            </div>
          {{ end }}
          <div style="font-size:0.9em; color:var(--color-secondary);">
            {{ $version }} - {{ $online_players }}/{{ $max_players }} players - {{ $size }}
          </div>
        {{ else if $starting }}
          <div style="font-size:0.9em; color:var(--color-secondary);">Server is starting up…</div>
        {{ else if $importing }}
          <div style="font-size:0.9em; color:var(--color-secondary);">Server is being imported…</div>
        {{ else if $updating }}
          <div style="font-size:0.9em; color:var(--color-secondary);">Server is being updated…</div>
        {{ else if $crashed }}
          <div style="font-size:0.9em; color:var(--color-secondary);">Server has crashed!</div>
        {{ else }}
          <div style="font-size:0.9em; color:var(--color-secondary);">Server is offline</div>
        {{ end }}
      </div>
    </div>
```


## Environment variables
- `CRAFTY_URL` - the URL or [IP]:[Port] to your Crafty Controller instance.
- `CRAFTY_SERVER_ID` - ID of the minecraft server you want to display info about. You can find it in the URL of the server details in Crafty Controller's dashboard. (Ex.; `https://[CRAFTY_URL]/panel/server_detail?id=[CRAFTY_SERVER_ID]`)
- `CRAFTY_API_TOKEN` - Your crafty API token. Requires `Full Access` permission.

## 🍻 Cheers

* @not-first as the author of [minecraft-server](/widgets/minecraft-server/README.md) community widget from which I took a lot of inspiration.
* @erkston on [Glance Discord](https://discord.com/invite/7KQ7Xa9kJd) for the solution on Self-Signed Certificate issue.
