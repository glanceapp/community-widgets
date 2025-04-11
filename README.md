This is a collection of custom widgets for [Glance](https://github.com/glanceapp/glance) made by the community using the `custom-api` and `extension` widgets.

The `custom-api` widgets in this repository have been vetted by the maintainers of Glance, however they are not responsible for the maintenance of these widgets. The author of each widget is responsible for maintaining and responding to issues and pull requests related to that widget.

Before a widget is added to this repository it must receive a vouch which can be done by anyone. If you would like to vouch for a widget, please leave a comment on its corresponding pull request.

To add your widget to the list, please read the [contribution guidelines](CONTRIBUTING.md).

### What's the difference between a `custom-api` and `extension` widget?

Custom API widgets are much easier to setup and usually only require a copy-paste into your config. Extension widgets are a bit more involved and require running a separate server or Docker container.

## Custom API Widgets
* [Animals Display](widgets/animals-display-by-panonim/README.md) by @panonim - Display random images of animals such as cats or dogs.
* [Cats As A Service Photos](widgets/cats-as-a-service-photos-by-gugugiyu/README.md) by @gugugiyu - show a grid of cat photos from the Cat as a Service [API](https://cataas.com/)
* [Formula 1 Widgets](widgets/formula1-widgets-by-abaza738/README.md) by @abaza738 - show different stats from the current Formula 1 season from [F1 API](https://f1api.dev)
* [Ghostfolio stats](widgets/ghostfolio-stats-by-ziritione85/README.md) by @ziritione85 - show the status of your investments with today's results, YTD and totals
* [Immich stats](widgets/immich-stats-by-svilenmarkov/README.md) by @svilenmarkov - show the number of photos, videos and usage of your Immich server
* [Last.FM Recent Tracks](widgets/lastfm-recent-tracks-by-akasiek/README.md) by @Akasiek - show recent tracks scrobbled by a Last.FM user
* [Mealie Today's Meal](widgets/mealie-todays-meal-by-wtoa/README.md) by @wtoa - show today's meal based off the meal planner from [Mealie](https://mealie.io/)
* [Minecraft server](widgets/minecraft-server-by-not-first/README.md) by @not-first - show online status, icon, version and player count of a minecraft server
* [Mullvad VPN status](widgets/mullvad-vpn-status-by-delmonteaj/README.md) by @DelMonteAJ - show VPN connection status, IP, and location
* [NextDNS Stats](widgets/nextdns-stats-by-ziritione85/README.md) by @ziritione85 - show the basics stats of your nextdns stats. Total queries, total blocked and percentage. 
* [NHL Scores](widgets/nhl-scores-by-jo-nike/README.md) by @jo-nike - show the schedule, scores, situation for NHL games today
* [Proxmox VE stats](widgets/proxmox-ve-stats-by-ralphocdol/README.md) by @ralphocdol - show the number of nodes, LXCs, VMs and Storage of your Proxmox Virtual Environment server
* [Random Bible Verse](widgets/random-bible-verse-by-pypp/README.md) by @pypp - show a random bible verse
* [Random fact](widgets/random-fact-by-svilenmarkov/README.md) by @svilenmarkov - show a random fact
* [Speedtest tracker](widgets/speedtest-tracker-by-not-first/README.md) by @not-first - show the latest internet speed information from speedtest tracker
* [Spotify Now Playing](widgets/spotify-now-playing-by-needsadjustment/README.md) by @needsadjustment - show the currently playing Spotify song
* [St. Louis Fed US Mortgage Rates] by @ehaughee - Show mortgage rates from the St. Louis Federal Reserve's FRED API
* [Steam Recently Played Games](widgets/steam-recently-played-games-by-lunnosmp4/README.md) by @lunnosmp4 - show a list of recently played games by a Steam user
* [Steam specials](widgets/steam-specials-by-svilenmarkov/README.md) by @svilenmarkov - show a list of discounted games on Steam
* [Steam User](widgets/steam-user-by-lunnosmp4/README.md) by @lunnosmp4 - show information about a Steam user
* [Tailscale devices](widgets/tailscale-devices-by-not-first/README.md) by @not-first - show all devices inside to a Tailscale tailnet along with their connection status, update availability and IP
* [Technitium DNS Stats](widgets/technitium-dns-stats-by-eribbey/README.md) by @eribbey - show stats from Technitium DNS Server
* [Trending Mastodon Links](widgets/trending-mastodon-links-by-tomcasavant/README.md) by @tomcasavant - shows a list of trending links from a provided Mastodon instance
* [Trending Bluesky News](widgets/trending-bluesky-news-by-tomcasavant/README.md) by @tomcasavant - shows a list of trending news links from the [Trending News 2.0 Feed](https://bsky.app/profile/did:plc:kkf4naxqmweop7dv4l2iqqf5/feed/news-2-0) on Bluesky
* [Uptime Kuma](widgets/uptime-kuma-by-not-first/README.md) by @not-first - show the status of Uptime Kuma services
* [SABnzbd Status](widgets/sabnzbd-stats-by-Neo11Neo/README.md) by @Neo11Neo - show SABnzbd status
* [Vikunja Taskboard](widgets/vikunja-taskboard-by-gugugiyu/README.md) by @gugugiyu - Real time taskboard using the Vikunja API
* [What's Up Docker Monitor](widgets/wud-monitor-by-panonim/README.md) by @panonim - A customizable widget for displaying the wud data about containers.

## Extension Widgets

> [!WARNING]
>
> Extension widgets are not actively monitored by the maintainers of Glance, use them at your own risk.

* [linktiles](https://github.com/haondt/linktiles/) by @haondt - display your linkding bookmarks in a configurable mosaic
* [Restic snapshot](https://github.com/not-first/restic-glance-extension) by @not-first - show the most recent snapshot and storage stats of a restic repo
