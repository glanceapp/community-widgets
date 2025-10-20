# Lemmy

This widget displays the latest posts from any [Lemmy](https://join-lemmy.org/) community inside.  It uses the public Lemmy API to fetch posts from a given community and shows the post title, score, comment count, and time since publication — output is formatted in a collapsible list.  It has the same look and feel as glances reddit widget.

---

## Preview
![](preview.png)

## Widget

```yaml
- type: custom-api
  title: L/${LEMMY_COMMUNITY}
  cache: 10m
  url: https://lemmy.world/api/v3/post/list?community_name=${LEMMY_COMMUNITY}&type_=${LEMMY_POST_TYPE}&sort=${LEMMY_POST_SORT}&limit=${LEMMY_POST_LIMIT}
  options:
    COLLAPSE_AFTER_COUNT: 5
  headers:
    Accept: application/json
  template: |
    {{ $collapse_after_count := .Options.IntOr "COLLAPSE_AFTER_COUNT" 5 }}

    <ul class="list list-gap-14 collapsible-container" data-collapse-after="{{ $collapse_after_count }}">
    {{ range.JSON.Array "posts" }}
      <li>
        <a href="{{ .String "post.ap_id" }}" class="size-title-dynamic color-primary-if-not-visited" target="_blank" rel="noreferrer">
          {{ .String "post.name" }}
        </a>
        <ul class="list-horizontal-text flex-nowrap">
          <li class="min-width-0">{{ .String "post.published" | parseTime "rfc3339" | toRelativeTime }}</li>
          <li class="min-width-0">{{ .Int "counts.score" }} pts</li>
          <li class="min-width-0">{{ .Int "counts.comments" }} comments</li>
          <li class="min-width-0 text-truncate">
            <a class="visited-indicator" target="_blank" rel="noreferrer" href="{{ .String "post.url" }}">Source</a>
          </li>
        </ul>
      </li>
    {{ else }}
      <li class="text-compact color-muted">No posts found.</li>
    {{ end }}
    </ul>
```

---

## Environment Variables

The followinfg environment variables are required in your `.env` file:

| Variable | Description | Example |
|-----------|--------------|----------|
| `LEMMY_COMMUNITY` | The community name (without the `!` prefix). | `technology` |
| `LEMMY_POST_TYPE` | Post type filter: `All`, `Local`, | `All` |
| `LEMMY_POST_SORT` | Sort order for posts. | `Hot`, `Active`, `New`, `TopDay`, `TopWeek` |
| `LEMMY_POST_LIMIT` | Number of posts to display. | `10` |

---

## Notes

- **Customisable Collapse:**  
  The `COLLAPSE_AFTER_COUNT` option controls how many posts appear before the list collapses.

- **Cross-instance Support:**  
  Change the `url` base to any Lemmy instance (e.g., `https://lemmy.ml`, `https://aussie.zone`).

> [!NOTE]
> The url base was hardcoded to the lemmy.world instance to give people flexibility if they wanted to reuse the multiple communities.

- **Example:**  
  To show top posts from the `australia` community on `lemmy.world`:
  ```yaml
  LEMMY_COMMUNITY=australia
  LEMMY_POST_TYPE=All
  LEMMY_POST_SORT=Hot
  LEMMY_POST_LIMIT=10
  ```
