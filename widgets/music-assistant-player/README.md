# Music Assistant Player

![Music Assistant widget in a small column](preview.png)

![Music Assistant widget with big-cover in a full column](preview-big-cover.png)

Shows what your [Music Assistant](https://www.music-assistant.io) players are playing and lets you control them from Glance: play/pause, previous/next, seek, volume, shuffle, repeat, favourite and play any upcoming track. Players with something in their queue get a tab; the one that is playing opens first.

Works with Music Assistant alone, no Home Assistant needed.

> [!CAUTION]
> With controls enabled (the default), the Music Assistant token is part of the page, because the browser sends the commands straight to Music Assistant. Anyone who can open your dashboard can read it. Use a dedicated user restricted to the players you show (step 1), keep Glance on your LAN or behind [Glance authentication](https://github.com/glanceapp/glance/blob/main/docs/configuration.md#authentication), or set `show-controls: false` to keep the token out of the page entirely.

## Requirements

- Glance v0.8.0 or newer (tested with v0.8.6)
- Music Assistant 2.10 or newer, with a long-lived access token

## Setup

1. **Create a user for the widget.** In Music Assistant open *Settings → User management*, add a user (for example `glance`) with the role **User** and set its *Player filter* to the players you want on the dashboard. Then open the ⋮ menu next to the user, choose *Manage access tokens* and create a token. It is valid for one year and does not renew itself. This user can control playback, edit the library (favourites, playlists) and read the configuration, but it cannot change settings, users or music sources. The **Guest** role does not work: Music Assistant does not create long-lived tokens for guests.

2. **Add two environment variables to Glance:**

   ```bash
   MA_URL=http://192.168.1.10:8095   # your Music Assistant server, no trailing slash
   MA_TOKEN=eyJhbGciOi...            # the token from step 1
   ```

3. **Add the widget** to a column, either by pasting the YAML below or by saving [`music-assistant.yml`](https://github.com/vednolacni/glance-music-assistant/blob/main/widget/music-assistant.yml) next to your `glance.yml` and adding `- $include: music-assistant.yml`.

4. **If Glance runs in Docker** and reaches Music Assistant under a different address than your browser does (for example a container name), set `public-url` to the address your browser can open. It is used for the controls and the cover art.

## Options

| Option | Default | What it does |
|---|---|---|
| `id` | `ma` | Unique id; give each copy a different one if you add the widget more than once to a page |
| `players` | `""` | Comma separated player ids to show; empty shows every player that has a queue and is in the user's player filter |
| `player-id` | `""` | Player to open when nothing is playing |
| `public-url` | `${MA_URL}` | Music Assistant address as the browser sees it |
| `show-controls` | `true` | `false` makes the widget read-only and leaves the token out of the page |
| `show-favorite` | `true` | Heart button that adds the current track to the favourites or removes it; `false` hides it |
| `queue-length` | `4` | Upcoming tracks to list; `0` hides the list |
| `max-players` | `4` | Number of player tabs, 1 to 6 |
| `big-cover` | `false` | Large cover like the Music Assistant player screen. When the widget is at least 720px wide (a column of about 760px, such as a full column), up next moves next to the cover; a longer `queue-length` (for example `8`) fills that side |
| `auto-refresh` | `true` | Update when the track ends and every `refresh-interval` seconds |
| `refresh-interval` | `30` | Seconds between updates while `auto-refresh` is on |

Player ids are shown at the top of each player's settings page in Music Assistant.

## Troubleshooting

- **"Cannot connect to …" under the controls**: the browser cannot reach `public-url`. If Glance is served over HTTPS, the browser only allows a secure connection to Music Assistant, so Music Assistant has to be behind HTTPS as well (for example your reverse proxy), and `public-url` must use `https://`.
- **"Music Assistant login failed"**: the token is wrong, expired (after one year) or revoked. Create a new one and update `MA_TOKEN`.
- **"Long-lived tokens cannot be created for guest accounts"** when creating the token: change the user's role to **User** (step 1).
- **"You do not have permission to perform this action"**: the user may not control that player. The widget only shows players in the user's player filter, so this happens when the filter changed after the page was loaded; reload the page.
- **Error `401 Unauthorized` in the widget**: Glance itself could not log in to Music Assistant; check `MA_URL` and `MA_TOKEN`.
- **"Nothing is playing" while music plays**: the player has no Music Assistant queue (for example playback started from another app), it is not in the user's player filter, or it is filtered out by `players`.
- **Music note instead of the cover**: Music Assistant could not fetch the image (for example a radio station logo whose address moved).

## How it works

Glance reads `player_queues/all`, `players/all` (the players the user may use), the upcoming tracks and the favourite state of each current track from the Music Assistant HTTP API on the server side. The buttons send commands over the Music Assistant WebSocket API, because the HTTP API does not allow requests from other sites in the browser (no CORS) while WebSockets are not limited that way. After a command, only this widget is re-rendered, without reloading the page.

`cache: 1s` is deliberate: the widget asks Glance for its fresh state about 1.5 seconds after a command, at the end of a track and every `refresh-interval` seconds, and a longer cache would show the old state. Music Assistant is only queried at those moments and when the page loads, not every second.

Source, tests and issues: [vednolacni/glance-music-assistant](https://github.com/vednolacni/glance-music-assistant).

Icons are from Google's Material Icons (Apache License 2.0).

## Widget YAML

<!-- widget-yaml:start -->
```yaml
- type: custom-api
  title: Music Assistant
  cache: 1s
  url: ${MA_URL}/api
  method: POST
  headers:
    Authorization: Bearer ${MA_TOKEN}
  body-type: json
  body:
    message_id: glance
    command: player_queues/all
  options:
    id: ma                  # unique id, change it if you add the widget more than once to a page
    players: ""             # comma separated player ids to show, empty shows every player with a queue
    player-id: ""           # player to open when nothing is playing
    public-url: ${MA_URL}   # Music Assistant address as the browser sees it (controls and cover art)
    show-controls: true     # false makes the widget read-only and keeps the token out of the page
    show-favorite: true     # heart button that adds or removes the current track from the favourites
    queue-length: 4         # upcoming tracks to list, 0 hides the list
    max-players: 4          # number of player tabs, 1 to 6
    big-cover: false        # large cover like the Music Assistant player screen, up next beside it in a wide column
    auto-refresh: true      # update when the track ends and every refresh-interval seconds
    refresh-interval: 30
  template: |
    {{- $id := .Options.StringOr "id" "ma" -}}
    {{- $server := trimSuffix "/" (.Options.StringOr "public-url" "${MA_URL}") -}}
    {{- $controls := .Options.BoolOr "show-controls" true -}}
    {{- $favorites := .Options.BoolOr "show-favorite" true -}}
    {{- $big := .Options.BoolOr "big-cover" false -}}
    {{- $queueLength := .Options.IntOr "queue-length" 4 -}}
    {{- $maxPlayers := .Options.IntOr "max-players" 4 -}}
    {{- if gt $maxPlayers 6 }}{{ $maxPlayers = 6 }}{{ end -}}
    {{- $preferred := .Options.StringOr "player-id" "" -}}
    {{- $allowed := printf ",%s," (replaceAll " " "" (.Options.StringOr "players" "")) -}}
    {{- $refresh := 0 -}}
    {{- if .Options.BoolOr "auto-refresh" true }}{{ $refresh = .Options.IntOr "refresh-interval" 30 }}{{ end -}}

    {{- /* Players the token's user may use: Music Assistant applies the user's player filter to players/all, not to player_queues/all */ -}}
    {{- $access := newRequest "${MA_URL}/api"
          | withHeader "Authorization" "Bearer ${MA_TOKEN}"
          | withHeader "Content-Type" "application/json"
          | withStringBody `{"message_id":"glance","command":"players/all"}`
          | getResponse -}}
    {{- $usable := "" -}}
    {{- if eq $access.Response.StatusCode 200 }}{{ $usable = "," }}{{ range $access.JSON.Array "" }}{{ $usable = printf "%s%s," $usable (.String "player_id") }}{{ end }}{{ end -}}

    {{- /* Pass 1: players to show (have a current track) and the one to open first */ -}}
    {{- $shown := "," -}}
    {{- $count := 0 -}}
    {{- $selected := "" -}}
    {{- $first := "" -}}
    {{- $preferredShown := false -}}
    {{- range .JSON.Array "" -}}
      {{- $qid := .String "queue_id" -}}
      {{- if and (.Bool "available") (ne (.String "current_item.queue_item_id") "") (lt $count $maxPlayers) (or (eq $allowed ",,") (ne (replaceAll (printf ",%s," $qid) "" $allowed) $allowed)) (or (eq $usable "") (ne (replaceAll (printf ",%s," $qid) "" $usable) $usable)) -}}
        {{- $shown = printf "%s%s," $shown $qid -}}
        {{- $count = add $count 1 -}}
        {{- if eq $first "" }}{{ $first = $qid }}{{ end -}}
        {{- if and (eq $selected "") (eq (.String "state") "playing") }}{{ $selected = $qid }}{{ end -}}
        {{- if eq $qid $preferred }}{{ $preferredShown = true }}{{ end -}}
      {{- end -}}
    {{- end -}}
    {{- if eq $selected "" }}{{ if $preferredShown }}{{ $selected = $preferred }}{{ else }}{{ $selected = $first }}{{ end }}{{ end }}

    <style>
      .root-music-assistant { position: relative; container-type: inline-size; }
      .boot-music-assistant { position: absolute; width: 1px; height: 1px; opacity: 0; pointer-events: none; }
      .radio-music-assistant { position: absolute; width: 1px; height: 1px; margin: 0; opacity: 0; pointer-events: none; }
      .panel-music-assistant { height: 0; overflow: hidden; visibility: hidden; }
      .radio-music-assistant:nth-of-type(1):checked ~ .panels-music-assistant > .panel-music-assistant:nth-child(1),
      .radio-music-assistant:nth-of-type(2):checked ~ .panels-music-assistant > .panel-music-assistant:nth-child(2),
      .radio-music-assistant:nth-of-type(3):checked ~ .panels-music-assistant > .panel-music-assistant:nth-child(3),
      .radio-music-assistant:nth-of-type(4):checked ~ .panels-music-assistant > .panel-music-assistant:nth-child(4),
      .radio-music-assistant:nth-of-type(5):checked ~ .panels-music-assistant > .panel-music-assistant:nth-child(5),
      .radio-music-assistant:nth-of-type(6):checked ~ .panels-music-assistant > .panel-music-assistant:nth-child(6) {
        height: auto; overflow: visible; visibility: visible;
      }
      .tabs-music-assistant { display: flex; flex-wrap: wrap; gap: 0.6rem; margin-bottom: 1.4rem; }
      .tab-music-assistant {
        display: inline-flex; align-items: center; gap: 0.6rem; padding: 0.3rem 1rem; cursor: pointer;
        border: 1px solid var(--color-separator); border-radius: 10rem; font-size: var(--font-size-h6);
        color: var(--color-text-subdue); white-space: nowrap;
      }
      .tab-music-assistant:hover { color: var(--color-text-highlight); }
      .tab-playing-music-assistant::before {
        content: ""; width: 0.6rem; height: 0.6rem; border-radius: 50%; background: var(--color-primary);
      }
      .radio-music-assistant:nth-of-type(1):checked ~ .tabs-music-assistant > .tab-music-assistant:nth-child(1),
      .radio-music-assistant:nth-of-type(2):checked ~ .tabs-music-assistant > .tab-music-assistant:nth-child(2),
      .radio-music-assistant:nth-of-type(3):checked ~ .tabs-music-assistant > .tab-music-assistant:nth-child(3),
      .radio-music-assistant:nth-of-type(4):checked ~ .tabs-music-assistant > .tab-music-assistant:nth-child(4),
      .radio-music-assistant:nth-of-type(5):checked ~ .tabs-music-assistant > .tab-music-assistant:nth-child(5),
      .radio-music-assistant:nth-of-type(6):checked ~ .tabs-music-assistant > .tab-music-assistant:nth-child(6) {
        color: var(--color-text-highlight); border-color: var(--color-primary);
      }
      .radio-music-assistant:nth-of-type(1):focus-visible ~ .tabs-music-assistant > .tab-music-assistant:nth-child(1),
      .radio-music-assistant:nth-of-type(2):focus-visible ~ .tabs-music-assistant > .tab-music-assistant:nth-child(2),
      .radio-music-assistant:nth-of-type(3):focus-visible ~ .tabs-music-assistant > .tab-music-assistant:nth-child(3),
      .radio-music-assistant:nth-of-type(4):focus-visible ~ .tabs-music-assistant > .tab-music-assistant:nth-child(4),
      .radio-music-assistant:nth-of-type(5):focus-visible ~ .tabs-music-assistant > .tab-music-assistant:nth-child(5),
      .radio-music-assistant:nth-of-type(6):focus-visible ~ .tabs-music-assistant > .tab-music-assistant:nth-child(6) {
        outline: 2px solid var(--color-primary); outline-offset: 2px;
      }
      .now-music-assistant { display: flex; align-items: center; gap: 1.4rem; }
      .art-music-assistant {
        position: relative; flex-shrink: 0; width: 7.2rem; aspect-ratio: 1; display: grid; place-items: center; overflow: hidden;
        border-radius: var(--border-radius); background: var(--color-widget-background-highlight); color: var(--color-text-subdue);
      }
      .art-music-assistant img { position: absolute; inset: 0; width: 100%; height: 100%; object-fit: cover; display: block; }
      .art-music-assistant svg { width: 40%; height: 40%; fill: currentColor; }
      .meta-music-assistant { min-width: 0; flex: 1; }
      .title-music-assistant { font-size: var(--font-size-h3); color: var(--color-text-highlight); }
      @container (min-width: 460px) {
        .art-music-assistant { width: 11.2rem; }
        .title-music-assistant { font-size: var(--font-size-h2); }
      }
      .big-music-assistant .now-music-assistant { flex-direction: column; align-items: stretch; text-align: center; gap: 1.2rem; }
      .big-music-assistant .art-music-assistant { width: min(100%, 32rem); margin-inline: auto; }
      .big-music-assistant .title-music-assistant { font-size: var(--font-size-h2); }
      .progress-music-assistant { display: flex; align-items: center; gap: 1rem; margin-top: 1.4rem; }
      .time-music-assistant { font-size: var(--font-size-h6); color: var(--color-text-subdue); font-variant-numeric: tabular-nums; min-width: 3.4em; }
      .time-music-assistant:last-child { text-align: right; }
      .bar-music-assistant { position: relative; flex: 1; height: 0.4rem; border-radius: 0.2rem; background: var(--color-progress-border); }
      .bar-music-assistant[data-ma-cmd] { cursor: pointer; }
      .bar-music-assistant[data-ma-cmd]::before { content: ""; position: absolute; inset: -0.8rem 0; }
      .fill-music-assistant { height: 100%; width: var(--ma-from); border-radius: inherit; background: var(--color-primary); }
      .playing-music-assistant .fill-music-assistant { animation: progress-music-assistant var(--ma-left) linear forwards; }
      @keyframes progress-music-assistant { from { width: var(--ma-from); } to { width: 100%; } }
      .transport-music-assistant { display: flex; align-items: center; justify-content: center; gap: 1.6rem; margin-top: 1rem; }
      .secondary-music-assistant { display: flex; align-items: center; justify-content: space-evenly; max-width: 30rem; margin: 0.6rem auto 0; }
      .big-music-assistant .progress-music-assistant { max-width: 56rem; margin-inline: auto; }
      .btn-music-assistant {
        appearance: none; display: grid; place-items: center; padding: 0.6rem; border: 0; border-radius: 50%;
        background: none; color: var(--color-text-base); cursor: pointer;
      }
      .btn-music-assistant svg { width: 2.2rem; height: 2.2rem; fill: currentColor; }
      .btn-music-assistant:hover { color: var(--color-text-highlight); background: var(--color-widget-background-highlight); }
      .btn-music-assistant:focus-visible, .row-music-assistant:focus-visible, .bar-music-assistant:focus-visible { outline: 2px solid var(--color-primary); outline-offset: 2px; }
      .play-music-assistant { width: 4.6rem; height: 4.6rem; background: var(--color-primary); color: var(--color-widget-background); }
      .play-music-assistant svg { width: 2.8rem; height: 2.8rem; }
      .play-music-assistant:hover { background: var(--color-primary); color: var(--color-widget-background); filter: brightness(1.12); }
      .secondary-music-assistant .btn-music-assistant { color: var(--color-text-subdue); }
      .secondary-music-assistant .btn-music-assistant svg { width: 1.9rem; height: 1.9rem; }
      .secondary-music-assistant .on-music-assistant { color: var(--color-primary); }
      .queue-music-assistant { margin-top: 1.4rem; padding-top: 1.2rem; border-top: 1px solid var(--color-separator); }
      .queue-head-music-assistant { display: flex; justify-content: space-between; margin-bottom: 0.6rem; font-size: var(--font-size-h6); color: var(--color-text-subdue); }
      @container (min-width: 720px) {
        .big-music-assistant .panel-music-assistant:has(> .queue-music-assistant) { display: grid; grid-template-columns: minmax(0, 1fr) minmax(0, 1fr); align-content: start; }
        .big-music-assistant .panel-music-assistant:has(> .queue-music-assistant) > :not(.queue-music-assistant) { grid-column: 1; width: 100%; }
        .big-music-assistant .panel-music-assistant:has(> .queue-music-assistant) > .queue-music-assistant {
          grid-column: 2; grid-row: 1 / span 5; margin: 0 0 0 3.2rem; padding: 0 0 0 3.2rem; border-top: 0; border-left: 1px solid var(--color-separator);
        }
      }
      .row-music-assistant { display: flex; align-items: center; gap: 1rem; margin-inline: -0.5rem; padding: 0.5rem; border-radius: var(--border-radius); }
      .row-music-assistant[data-ma-cmd] { cursor: pointer; }
      .row-music-assistant[data-ma-cmd]:hover { background: var(--color-widget-background-highlight); }
      .thumb-music-assistant {
        flex-shrink: 0; width: 3.6rem; height: 3.6rem; object-fit: cover; display: block;
        border-radius: calc(var(--border-radius) - 1px); background: var(--color-widget-background-highlight);
      }
      .empty-music-assistant { display: grid; justify-items: center; gap: 0.4rem; padding-block: 1.6rem; text-align: center; }
      .empty-music-assistant svg { width: 3.2rem; height: 3.2rem; fill: var(--color-text-subdue); margin-bottom: 0.6rem; }
      .root-music-assistant[data-ma-busy] .btn-music-assistant,
      .root-music-assistant[data-ma-busy] [data-ma-cmd] { opacity: 0.55; pointer-events: none; }
      .root-music-assistant[data-ma-error]::after {
        content: attr(data-ma-error); display: block; margin-top: 1rem; font-size: var(--font-size-h6); color: var(--color-negative);
      }
    </style>

    <div class="root-music-assistant{{ if $big }} big-music-assistant{{ end }}" data-ma-root="{{ $id }}" data-ma-server="{{ $server }}" data-ma-refresh="{{ $refresh }}"{{ if $controls }} data-ma-token="${MA_TOKEN}"{{ end }}>
    {{- if or $controls (gt $refresh 0) }}
      <img class="boot-music-assistant" alt="" aria-hidden="true" src="data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7" onload="
        if (!window.maPlayer) window.maPlayer = {
          timers: {},
          waiting: {},
          find: function (id) {
            var roots = document.querySelectorAll('[data-ma-root]');
            for (var i = 0; i < roots.length; i++) if (roots[i].getAttribute('data-ma-root') === id) return roots[i];
            return null;
          },
          send: function (root, command, args, done) {
            var server = root.getAttribute('data-ma-server') || '', socket = null, finished = false, timer = 0;
            while (server.slice(-1) === '/') server = server.slice(0, -1);
            var finish = function (error) {
              if (finished) return;
              finished = true;
              clearTimeout(timer);
              try { socket.close(); } catch (e) {}
              done(error);
            };
            timer = setTimeout(function () { finish('Music Assistant did not answer'); }, 8000);
            try {
              socket = new WebSocket((server.indexOf('https:') === 0 ? 'wss:' : 'ws:') + server.slice(server.indexOf(':') + 1) + '/ws');
            } catch (e) { return finish('Cannot connect to ' + server); }
            socket.onerror = function () { finish('Cannot connect to ' + server); };
            socket.onclose = function () { finish('Music Assistant closed the connection'); };
            socket.onopen = function () {
              socket.send(JSON.stringify({ message_id: 'auth', command: 'auth', args: { token: root.getAttribute('data-ma-token') } }));
            };
            socket.onmessage = function (message) {
              var data = null;
              try { data = JSON.parse(message.data); } catch (e) { return; }
              if (data.message_id === 'auth') {
                if (data.error_code != null) return finish('Music Assistant login failed: ' + (data.details || data.error_code));
                socket.send(JSON.stringify({ message_id: 'cmd', command: command, args: args }));
              } else if (data.message_id === 'cmd') {
                finish(data.error_code != null ? (data.details || 'Command failed') : null);
              }
            };
          },
          cmd: function (el, event) {
            var root = el.closest('[data-ma-root]'), panel = el.closest('[data-ma-queue]');
            if (!root || !panel || root.hasAttribute('data-ma-busy')) return;
            var command = el.getAttribute('data-ma-cmd'), args = {};
            try { args = JSON.parse(el.getAttribute('data-ma-args') || '{}'); } catch (e) {}
            if (command.indexOf('player_queues/') === 0) args.queue_id = panel.getAttribute('data-ma-queue');
            else if (command.indexOf('players/') === 0) args.player_id = panel.getAttribute('data-ma-queue');
            if (el.hasAttribute('data-ma-seek')) {
              var box = el.getBoundingClientRect();
              if (!box.width) return;
              args.position = Math.max(0, Math.round((event.clientX - box.left) / box.width * Number(el.getAttribute('data-ma-seek'))));
            }
            root.setAttribute('data-ma-busy', '');
            root.removeAttribute('data-ma-error');
            window.maPlayer.send(root, command, args, function (error) {
              root.removeAttribute('data-ma-busy');
              if (error) return root.setAttribute('data-ma-error', error);
              window.maPlayer.refresh(root.getAttribute('data-ma-root'), 1500);
            });
          },
          refresh: function (id, delay) {
            var player = window.maPlayer;
            clearTimeout(player.timers[id]);
            player.timers[id] = setTimeout(function () {
              if (document.hidden) {
                if (!player.waiting[id]) {
                  player.waiting[id] = true;
                  document.addEventListener('visibilitychange', function resume() {
                    if (document.hidden) return;
                    document.removeEventListener('visibilitychange', resume);
                    player.waiting[id] = false;
                    player.refresh(id, 0);
                  });
                }
                return;
              }
              var current = player.find(id);
              if (!current) return;
              if (typeof pageData !== 'object') return location.reload();
              var checked = current.querySelector('.radio-music-assistant:checked');
              var keep = checked ? checked.getAttribute('data-qid') : '';
              fetch(pageData.baseURL + '/api/pages/' + pageData.slug + '/content/')
                .then(function (response) { if (!response.ok) throw new Error(response.status); return response.text(); })
                .then(function (html) {
                  var fresh = null, doc = new DOMParser().parseFromString(html, 'text/html');
                  var roots = doc.querySelectorAll('[data-ma-root]');
                  for (var i = 0; i < roots.length; i++) if (roots[i].getAttribute('data-ma-root') === id) fresh = roots[i];
                  current = player.find(id);
                  if (!fresh || !current) throw new Error('widget not found');
                  var radios = fresh.querySelectorAll('.radio-music-assistant');
                  for (var j = 0; j < radios.length; j++) {
                    if (radios[j].getAttribute('data-qid') !== keep) continue;
                    for (var k = 0; k < radios.length; k++) radios[k].removeAttribute('checked');
                    radios[j].setAttribute('checked', '');
                  }
                  current.replaceWith(document.importNode(fresh, true));
                })
                .catch(function () { location.reload(); });
            }, delay);
          },
          arm: function (img) {
            var root = img.closest('[data-ma-root]');
            if (!root) return;
            var wait = Number(root.getAttribute('data-ma-refresh')) || 0;
            if (!wait) return;
            var playing = root.querySelectorAll('[data-ma-remaining]');
            for (var i = 0; i < playing.length; i++) {
              var left = Number(playing[i].getAttribute('data-ma-remaining'));
              if (left >= 0 && left + 2 < wait) wait = left + 2;
            }
            window.maPlayer.refresh(root.getAttribute('data-ma-root'), wait * 1000);
          }
        };
        window.maPlayer.arm(this);">
    {{- end }}

    {{- if eq $count 0 }}
      <div class="empty-music-assistant">
        {{ template "ma-icon-note" }}
        <div class="color-highlight">Nothing is playing</div>
        <div class="size-h6 color-subdue">Start music in Music Assistant and it shows up here.</div>
        <a class="size-h6 color-primary" href="{{ $server }}" target="_blank" rel="noreferrer">Open Music Assistant</a>
      </div>
    {{- else }}

    {{- /* One radio per player: the checked one decides which panel and tab are shown */ -}}
    {{- range .JSON.Array "" }}
      {{- $qid := .String "queue_id" }}
      {{- if ne (replaceAll (printf ",%s," $qid) "" $shown) $shown }}
      <input type="radio" class="radio-music-assistant" name="ma-{{ $id }}" id="ma-{{ $id }}-{{ $qid }}" data-qid="{{ $qid }}" aria-label="{{ .String "display_name" }}"{{ if eq $qid $selected }} checked{{ end }}>
      {{- end }}
    {{- end }}

    {{- if gt $count 1 }}
      <div class="tabs-music-assistant">
      {{- range .JSON.Array "" }}
        {{- $qid := .String "queue_id" }}
        {{- if ne (replaceAll (printf ",%s," $qid) "" $shown) $shown }}
        <label class="tab-music-assistant{{ if eq (.String "state") "playing" }} tab-playing-music-assistant{{ end }}" for="ma-{{ $id }}-{{ $qid }}">{{ .String "display_name" }}</label>
        {{- end }}
      {{- end }}
      </div>
    {{- end }}

      <div class="panels-music-assistant">
      {{- range .JSON.Array "" }}
        {{- $qid := .String "queue_id" }}
        {{- if ne (replaceAll (printf ",%s," $qid) "" $shown) $shown }}
        {{- $state := .String "state" }}
        {{- $duration := .Int "current_item.duration" }}
        {{- $elapsed := .Float "elapsed_time" }}
        {{- if eq $state "playing" }}
          {{- $elapsed = add $elapsed ((now).Sub (parseTime "unix" (printf "%.0f" (.Float "elapsed_time_last_updated")))).Seconds }}
        {{- end }}
        {{- if lt $elapsed 0.0 }}{{ $elapsed = 0.0 }}{{ end }}
        {{- if and (gt $duration 0) (gt $elapsed (toFloat $duration)) }}{{ $elapsed = toFloat $duration }}{{ end }}
        {{- $elapsedSeconds := toInt $elapsed }}
        {{- $remaining := sub $duration $elapsedSeconds }}
        {{- $percent := 0.0 }}
        {{- if gt $duration 0 }}{{ $percent = mul (div $elapsed (toFloat $duration)) 100.0 }}{{ end }}
        {{- $title := .String "current_item.streamdetails.stream_metadata.title" }}
        {{- if eq $title "" }}{{ $title = .String "current_item.media_item.name" }}{{ end }}
        {{- if eq $title "" }}{{ $title = .String "current_item.name" }}{{ end }}
        {{- $artist := .String "current_item.streamdetails.stream_metadata.artist" }}
        {{- $album := .String "current_item.streamdetails.stream_metadata.album" }}
        {{- if eq $album "" }}{{ $album = .String "current_item.media_item.album.name" }}{{ end }}
        {{- $image := .String "current_item.image.proxy_id" }}
        {{- $shuffle := .Bool "shuffle_enabled" }}
        {{- $repeat := .String "repeat_mode" }}
        {{- $nextRepeat := "all" }}
        {{- if eq $repeat "all" }}{{ $nextRepeat = "one" }}{{ else if eq $repeat "one" }}{{ $nextRepeat = "off" }}{{ end }}
        {{- /* Favourite state: queue items keep the flag they had when queued, the library has the current one */ -}}
        {{- $favorite := "" }}
        {{- $favoriteArgs := "" }}
        {{- $uri := .String "current_item.media_item.uri" }}
        {{- if and $controls $favorites $uri }}
          {{- $item := newRequest "${MA_URL}/api"
                | withHeader "Authorization" "Bearer ${MA_TOKEN}"
                | withHeader "Content-Type" "application/json"
                | withStringBody (printf `{"message_id":"glance","command":"music/item_by_uri","args":{"uri":%q}}` $uri)
                | getResponse }}
          {{- if eq $item.Response.StatusCode 200 }}
            {{- if and ($item.JSON.Bool "favorite") (eq ($item.JSON.String "provider") "library") }}
              {{- $favorite = "on" }}
              {{- $favoriteArgs = printf `{"media_type": %q, "library_item_id": %q}` ($item.JSON.String "media_type") ($item.JSON.String "item_id") }}
            {{- else }}
              {{- $favorite = "off" }}
              {{- $favoriteArgs = printf `{"item": %q}` ($item.JSON.String "uri") }}
            {{- end }}
          {{- end }}
        {{- end }}
        <div class="panel-music-assistant{{ if eq $state "playing" }} playing-music-assistant{{ end }}" data-ma-queue="{{ $qid }}"{{ if and (eq $state "playing") (gt $duration 0) }} data-ma-remaining="{{ $remaining }}"{{ end }}>
          <div class="now-music-assistant">
            <div class="art-music-assistant">
              {{- template "ma-icon-note" }}
              {{- if $image }}<img src="{{ $server }}/imageproxy/{{ $image }}?size={{ if $big }}512{{ else }}256{{ end }}" alt="" onerror="this.remove()">{{ end -}}
            </div>
            <div class="meta-music-assistant">
              <div class="title-music-assistant text-truncate" title="{{ $title }}">{{ $title }}</div>
              <div class="color-base text-truncate">
                {{- if $artist }}{{ $artist }}{{ else }}{{ range $i, $a := .Array "current_item.media_item.artists" }}{{ if $i }}, {{ end }}{{ $a.String "name" }}{{ end }}{{ end -}}
              </div>
              {{- if $album }}
              <div class="size-h6 color-subdue text-truncate">{{ $album }}</div>
              {{- end }}
            </div>
          </div>

          {{- if gt $duration 0 }}
          <div class="progress-music-assistant">
            <span class="time-music-assistant">{{ template "ma-time" $elapsedSeconds }}</span>
            <div class="bar-music-assistant"{{ if $controls }} data-ma-cmd="player_queues/seek" data-ma-seek="{{ $duration }}" title="Click to seek" onclick="window.maPlayer && window.maPlayer.cmd(this, event)"{{ end }}>
              <div class="fill-music-assistant" style="--ma-from: {{ printf "%.2f" $percent }}%; --ma-left: {{ $remaining }}s"></div>
            </div>
            <span class="time-music-assistant">{{ template "ma-time" $duration }}</span>
          </div>
          {{- end }}

          {{- if $controls }}
          <div class="transport-music-assistant">
            <button type="button" class="btn-music-assistant" data-ma-cmd="player_queues/previous" aria-label="Previous track" title="Previous track" onclick="window.maPlayer && window.maPlayer.cmd(this, event)">{{ template "ma-icon-previous" }}</button>
            <button type="button" class="btn-music-assistant play-music-assistant" data-ma-cmd="player_queues/play_pause" aria-label="{{ if eq $state "playing" }}Pause{{ else }}Play{{ end }}" title="{{ if eq $state "playing" }}Pause{{ else }}Play{{ end }}" onclick="window.maPlayer && window.maPlayer.cmd(this, event)">{{ if eq $state "playing" }}{{ template "ma-icon-pause" }}{{ else }}{{ template "ma-icon-play" }}{{ end }}</button>
            <button type="button" class="btn-music-assistant" data-ma-cmd="player_queues/next" aria-label="Next track" title="Next track" onclick="window.maPlayer && window.maPlayer.cmd(this, event)">{{ template "ma-icon-next" }}</button>
          </div>
          <div class="secondary-music-assistant">
            <button type="button" class="btn-music-assistant" data-ma-cmd="players/cmd/volume_down" aria-label="Volume down" title="Volume down" onclick="window.maPlayer && window.maPlayer.cmd(this, event)">{{ template "ma-icon-volume-down" }}</button>
            <button type="button" class="btn-music-assistant" data-ma-cmd="players/cmd/volume_up" aria-label="Volume up" title="Volume up" onclick="window.maPlayer && window.maPlayer.cmd(this, event)">{{ template "ma-icon-volume-up" }}</button>
            <button type="button" class="btn-music-assistant{{ if $shuffle }} on-music-assistant{{ end }}" data-ma-cmd="player_queues/shuffle" data-ma-args='{"shuffle_enabled": {{ not $shuffle }}}' aria-label="Shuffle" aria-pressed="{{ $shuffle }}" title="Shuffle {{ if $shuffle }}on{{ else }}off{{ end }}" onclick="window.maPlayer && window.maPlayer.cmd(this, event)">{{ template "ma-icon-shuffle" }}</button>
            <button type="button" class="btn-music-assistant{{ if ne $repeat "off" }} on-music-assistant{{ end }}" data-ma-cmd="player_queues/repeat" data-ma-args='{"repeat_mode": "{{ $nextRepeat }}"}' aria-label="Repeat: {{ $repeat }}" title="Repeat: {{ $repeat }}" onclick="window.maPlayer && window.maPlayer.cmd(this, event)">{{ if eq $repeat "one" }}{{ template "ma-icon-repeat-one" }}{{ else }}{{ template "ma-icon-repeat" }}{{ end }}</button>
            {{- if $favorite }}
            <button type="button" class="btn-music-assistant{{ if eq $favorite "on" }} on-music-assistant{{ end }}" data-ma-cmd="music/favorites/{{ if eq $favorite "on" }}remove_item{{ else }}add_item{{ end }}" data-ma-args="{{ $favoriteArgs }}" aria-label="Favourite" aria-pressed="{{ eq $favorite "on" }}" title="{{ if eq $favorite "on" }}Remove from favourites{{ else }}Add to favourites{{ end }}" onclick="window.maPlayer && window.maPlayer.cmd(this, event)">{{ if eq $favorite "on" }}{{ template "ma-icon-favorite" }}{{ else }}{{ template "ma-icon-favorite-border" }}{{ end }}</button>
            {{- end }}
          </div>
          {{- end }}

          {{- $index := .Int "current_index" }}
          {{- $upcoming := sub (.Int "items") (add $index 1) }}
          {{- if and (gt $queueLength 0) (gt $upcoming 0) }}
          {{- $next := newRequest "${MA_URL}/api"
                | withHeader "Authorization" "Bearer ${MA_TOKEN}"
                | withHeader "Content-Type" "application/json"
                | withStringBody (printf `{"message_id":"glance","command":"player_queues/items","args":{"queue_id":%q,"offset":%d,"limit":%d}}` $qid (add $index 1) $queueLength)
                | getResponse }}
          {{- if eq $next.Response.StatusCode 200 }}
          <div class="queue-music-assistant">
            <div class="queue-head-music-assistant"><span>Up next</span><span>{{ $upcoming }}</span></div>
            <ul>
            {{- range $next.JSON.Array "" }}
              {{- $name := .String "media_item.name" }}
              {{- if eq $name "" }}{{ $name = .String "name" }}{{ end }}
              <li class="row-music-assistant"{{ if $controls }} role="button" tabindex="0" title="Play now" data-ma-cmd="player_queues/play_index" data-ma-args='{"index": "{{ .String "queue_item_id" }}"}' onclick="window.maPlayer && window.maPlayer.cmd(this, event)" onkeydown="if (event.key === 'Enter' || event.key === ' ') { event.preventDefault(); this.click(); }"{{ end }}>
                {{- if .String "image.proxy_id" }}
                <img class="thumb-music-assistant" src="{{ $server }}/imageproxy/{{ .String "image.proxy_id" }}?size=80" alt="" onerror="var s = document.createElement('span'); s.className = this.className; this.replaceWith(s);">
                {{- else }}
                <span class="thumb-music-assistant"></span>
                {{- end }}
                <div class="min-width-0 grow">
                  <div class="color-highlight text-truncate">{{ $name }}</div>
                  <div class="size-h6 color-subdue text-truncate">{{ range $i, $a := .Array "media_item.artists" }}{{ if $i }}, {{ end }}{{ $a.String "name" }}{{ end }}</div>
                </div>
                {{- if gt (.Int "duration") 0 }}
                <span class="size-h6 color-subdue">{{ template "ma-time" (.Int "duration") }}</span>
                {{- end }}
              </li>
            {{- end }}
            </ul>
          </div>
          {{- end }}
          {{- end }}
        </div>
        {{- end }}
      {{- end }}
      </div>
    {{- end }}
    </div>

    {{- define "ma-time" }}{{ if ge . 3600 }}{{ printf "%d:%02d:%02d" (div . 3600) (div (mod . 3600) 60) (mod . 60) }}{{ else }}{{ printf "%d:%02d" (div . 60) (mod . 60) }}{{ end }}{{ end }}
    {{- define "ma-icon-note" }}<svg viewBox="0 0 24 24" aria-hidden="true"><path d="M12 3v10.55c-.59-.34-1.27-.55-2-.55-2.21 0-4 1.79-4 4s1.79 4 4 4 4-1.79 4-4V7h4V3h-6z"/></svg>{{ end }}
    {{- define "ma-icon-play" }}<svg viewBox="0 0 24 24" aria-hidden="true"><path d="M8 5v14l11-7z"/></svg>{{ end }}
    {{- define "ma-icon-pause" }}<svg viewBox="0 0 24 24" aria-hidden="true"><path d="M6 19h4V5H6v14zm8-14v14h4V5h-4z"/></svg>{{ end }}
    {{- define "ma-icon-previous" }}<svg viewBox="0 0 24 24" aria-hidden="true"><path d="M6 6h2v12H6zm3.5 6l8.5 6V6z"/></svg>{{ end }}
    {{- define "ma-icon-next" }}<svg viewBox="0 0 24 24" aria-hidden="true"><path d="M6 18l8.5-6L6 6v12zM16 6v12h2V6h-2z"/></svg>{{ end }}
    {{- define "ma-icon-shuffle" }}<svg viewBox="0 0 24 24" aria-hidden="true"><path d="M10.59 9.17L5.41 4 4 5.41l5.17 5.17 1.42-1.41zM14.5 4l2.04 2.04L4 18.59 5.41 20 17.96 7.46 20 9.5V4h-5.5zm.33 9.41l-1.41 1.41 3.13 3.13L14.5 20H20v-5.5l-2.04 2.04-3.13-3.13z"/></svg>{{ end }}
    {{- define "ma-icon-repeat" }}<svg viewBox="0 0 24 24" aria-hidden="true"><path d="M7 7h10v3l4-4-4-4v3H5v6h2V7zm10 10H7v-3l-4 4 4 4v-3h12v-6h-2v4z"/></svg>{{ end }}
    {{- define "ma-icon-repeat-one" }}<svg viewBox="0 0 24 24" aria-hidden="true"><path d="M7 7h10v3l4-4-4-4v3H5v6h2V7zm10 10H7v-3l-4 4 4 4v-3h12v-6h-2v4zm-4-2V9h-1l-2 1v1h1.5v4H13z"/></svg>{{ end }}
    {{- define "ma-icon-favorite" }}<svg viewBox="0 0 24 24" aria-hidden="true"><path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"/></svg>{{ end }}
    {{- define "ma-icon-favorite-border" }}<svg viewBox="0 0 24 24" aria-hidden="true"><path d="M16.5 3c-1.74 0-3.41.81-4.5 2.09C10.91 3.81 9.24 3 7.5 3 4.42 3 2 5.42 2 8.5c0 3.78 3.4 6.86 8.55 11.54L12 21.35l1.45-1.32C18.6 15.36 22 12.28 22 8.5 22 5.42 19.58 3 16.5 3zm-4.4 15.55l-.1.1-.1-.1C7.14 14.24 4 11.39 4 8.5 4 6.5 5.5 5 7.5 5c1.54 0 3.04.99 3.57 2.36h1.87C13.46 5.99 14.96 5 16.5 5c2 0 3.5 1.5 3.5 3.5 0 2.89-3.14 5.74-7.9 10.05z"/></svg>{{ end }}
    {{- define "ma-icon-volume-down" }}<svg viewBox="0 0 24 24" aria-hidden="true"><path d="M18.5 12c0-1.77-1.02-3.29-2.5-4.03v8.05c1.48-.73 2.5-2.25 2.5-4.02zM5 9v6h4l5 5V4L9 9H5z"/></svg>{{ end }}
    {{- define "ma-icon-volume-up" }}<svg viewBox="0 0 24 24" aria-hidden="true"><path d="M3 9v6h4l5 5V4L7 9H3zm13.5 3c0-1.77-1.02-3.29-2.5-4.03v8.05c1.48-.73 2.5-2.25 2.5-4.02zM14 3.23v2.06c2.89.86 5 3.54 5 6.71s-2.11 5.85-5 6.71v2.06c4.01-.91 7-4.49 7-8.77s-2.99-7.86-7-8.77z"/></svg>{{ end }}
```
<!-- widget-yaml:end -->
