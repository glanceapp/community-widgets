package main

import (
	"encoding/json"
	"log"
	"os"
	"regexp"
	"slices"
	"strings"
	"time"
)

const siteDataPath = "../site/data.json"

type siteWidget struct {
	Type        string     `json:"type"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Author      string     `json:"author"`
	Category    string     `json:"category"`
	Directory   string     `json:"directory,omitempty"`
	Preview     string     `json:"preview,omitempty"`
	URL         string     `json:"url,omitempty"`
	EnvVars     []string   `json:"env_vars,omitempty"`
	NeedsKey    bool       `json:"needs_key,omitempty"`
	TimeAdded   *time.Time `json:"time_added,omitempty"`
	TimeUpdated *time.Time `json:"time_updated,omitempty"`
}

type siteCategory struct {
	name    string
	pattern *regexp.Regexp
}

func category(name string, keywords ...string) siteCategory {
	return siteCategory{name, regexp.MustCompile(`\b(` + strings.Join(keywords, "|") + `)\b`)}
}

// Categories are inferred from the widget's directory, title and description,
// the first one with a matching keyword wins so the order matters
var siteCategories = []siteCategory{
	category("Sports", "afl", "football", "formula", "f1", "mlb", "nba", "ncaa", "nfl", "nhl", "tennis", "ufc", "soccer", "hockey"),
	category("Gaming", "steam", "epic", "minecraft", "retroachievements?", "chess", "tsumego", "romm", "pok[eé]mon", "discopanel", "lichess", "playstation", "xbox"),
	category("Fun & Daily", "animals?", "cats?", "calvin", "xkcd", "random", "astronomy", "quotes?"),
	category("Network & VPN", "tailscale", "netbird", "wg-easy", "wireguard", "mullvad", "gluetun", "azirevpn", "vpn", "cloudflared?", "nextdns", "dns", "unifi", "netalertx", "nginx", "speedtest", "pi-?hole", "adguard"),
	category("Media & Downloads", "jellyfin", "emby", "plex", "sonarr", "radarr", "lidarr", "prowlarr", "arr", "overseerr", "tautulli", "trakt", "spotify", "last\\.?fm", "listenbrainz", "audiobookshelf", "kavita", "komga", "calibre", "hardcover", "myanimelist", "anime", "nebula", "youtube", "immich", "mediatracker", "media", "movies?", "cineplex", "eventim", "qbittorrent", "nzbget", "sabnzbd", "transmission", "torrents?", "comic\\w*", "bilibili", "twitch"),
	category("Finance", "crypto\\w*", "bitcoin", "ghostfolio", "mortgage", "stocks?", "currency", "exchange rates?"),
	category("Developer", "github", "gitlab", "gitea", "forgejo", "leetcode", "deploy", "community widgets"),
	category("Weather & Local", "weather", "air quality", "algal", "sunrise", "bin", "tube", "transit\\w*", "aircraft", "dhl", "parcels?", "flights?"),
	category("News & Social", "lemmy", "bluesky", "mastodon", "wikipedia", "wikiwidgets", "news", "reddit", "rss"),
	category("Productivity & Home", "todoist", "vikunja", "nextcloud", "calendar", "raindrop", "linkwarden", "karakeep", "bookmarks?", "paperless\\w*", "papra", "mealie", "gotify", "countdown", "time bar", "duolingo", "wanikani", "dawarich", "home ?assistant", "tasks?"),
	category("Monitoring & Servers", "beszel", "proxmox", "truenas", "unraid", "synology", "glances", "docker", "portainer", "komodo", "containers?", "cup", "scrutiny", "ups", "grafana", "gatus", "uptime", "backrest", "restic", "kubernetes", "frigate", "syncthing", "termix", "pagerduty", "status", "monitor\\w*", "servers?", "metrics"),
}

var (
	markdownLinkPattern = regexp.MustCompile(`\[([^\]]+)\]\([^)]*\)`)
	envVarPattern       = regexp.MustCompile(`\$\{([A-Z][A-Z0-9_]*)\}`)
	credentialPattern   = regexp.MustCompile(`KEY|TOKEN|SECRET|PASS|AUTH|CREDENTIAL`)
)

func generateSiteData(widgets []widget, extensions []extension) {
	entries := make([]siteWidget, 0, len(widgets)+len(extensions))

	for _, w := range widgets {
		readme, err := os.ReadFile(widgetPath(w.Directory, "README.md"))
		if err != nil {
			log.Fatalf("Failed to read README for widget %s: %v", w.Directory, err)
		}

		envVars := extractEnvVars(string(readme))
		description := plainText(w.Description)

		entries = append(entries, siteWidget{
			Type:        "custom-api",
			Title:       w.Title,
			Description: description,
			Author:      w.Author,
			Category:    inferCategory(w.Directory, w.Title, description),
			Directory:   w.Directory,
			Preview:     w.Preview,
			EnvVars:     envVars,
			NeedsKey:    slices.ContainsFunc(envVars, credentialPattern.MatchString),
			TimeAdded:   &w.TimeAdded,
			TimeUpdated: &w.TimeUpdated,
		})
	}

	for _, e := range extensions {
		description := plainText(e.Description)

		entries = append(entries, siteWidget{
			Type:        "extension",
			Title:       e.Title,
			Description: description,
			Author:      e.Author,
			Category:    inferCategory(e.Title, description),
			URL:         e.URL,
		})
	}

	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal site data to JSON: %v", err)
	}

	err = os.WriteFile(siteDataPath, data, 0644)
	if err != nil {
		log.Fatalf("Failed to write %s: %v", siteDataPath, err)
	}
}

func inferCategory(fields ...string) string {
	text := strings.ToLower(strings.Join(fields, " "))

	for _, c := range siteCategories {
		if c.pattern.MatchString(text) {
			return c.name
		}
	}

	return "Other"
}

// Strips markdown links and emphasis so descriptions can be shown as plain text
func plainText(s string) string {
	s = markdownLinkPattern.ReplaceAllString(s, "$1")
	s = strings.NewReplacer("**", "", "`", "").Replace(s)
	return strings.TrimSpace(s)
}

func extractEnvVars(readme string) []string {
	vars := []string{}

	for _, match := range envVarPattern.FindAllStringSubmatch(readme, -1) {
		if !slices.Contains(vars, match[1]) {
			vars = append(vars, match[1])
		}
	}

	return vars
}
