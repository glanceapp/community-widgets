# Algal Bloom Update Widget

This Glance `custom-api` widget displays the latest algal bloom information for a selected beach using data from [Beachsafe](https://beachsafe.org.au). It provides a quick overview of water quality, including beach cleaning status, abnormal foam, and abnormal water colour. The widget also shows the last updated timestamp in relative time.

## Features

- Displays the beach name prominently.
- Shows whether the beach has been cleaned.
- Indicates the presence of abnormal foam.
- Indicates abnormal water colour in the water.
- Last update timestamp displayed dynamically using relative time.
- Colour-coded indicators for easy visual recognition:
  - **Positive (green)**: Good conditions
  - **Negative (red)**: Potential issues

## How It Works

This widget fetches JSON data from the Beachsafe API and parses the `algal_bloom` object. The relevant fields are:

- `beach.algal_bloom.clean_rating` – Whether the beach has been cleaned (true/false)
- `beach.algal_bloom.beach_foam` – Whether abnormal foam is present (true/false)
- `beach.algal_bloom.water_discolour` – Whether abnormal water colour is present (true/false)
- `beach.algal_bloom.updated_at` – Timestamp of the last update (RFC3339)

The template uses Go-style conditional rendering to display each field with the appropriate colour.

## Widget YAML

```yaml
- type: custom-api
  title: Algal Bloom Update
  url: https://beachsafe.org.au/api/v4/beach/brighton-3 
  method: GET
  headers:
    "User-Agent": Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36
  template: |
    {{ $b := .JSON }}
    <div class="size-h3">{{ $b.String "beach.title" }}</div>
    <div>Updated <span {{ $b.String "beach.algal_bloom.updated_at" | parseTime "rfc3339" | toRelativeTime }}></span> ago</div><br>
    <div class="size-h4 color-{{ if $b.Bool "beach.algal_bloom.clean_rating" }}positive{{ else }}negative{{ end }}">
      Beach cleaned: {{ if $b.Bool "beach.algal_bloom.clean_rating" }}Yes{{ else }}No{{ end }}
    </div>
    <div class="size-h4 color-{{ if $b.Bool "beach.algal_bloom.beach_foam" }}negative{{ else }}positive{{ end }}">
      Abnormal foam: {{ if $b.Bool "beach.algal_bloom.beach_foam" }}Yes{{ else }}No{{ end }}
    </div>
    <div class="size-h4 color-{{ if $b.Bool "beach.algal_bloom.water_discolour" }}negative{{ else }}positive{{ end }}">
      Abnormal water colour: {{ if $b.Bool "beach.algal_bloom.water_discolour" }}Yes{{ else }}No{{ end }}
    </div>
