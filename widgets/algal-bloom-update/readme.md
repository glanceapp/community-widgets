# Algal Bloom Update Widget
South Australia is experiencing a toxic algal bloom caused by Karenia species. While some blooms are harmless, this one can affect fish and other marine life.

The bloom has been driven by a mix of factors, including river floodwaters adding nutrients, summer upwellings bringing nutrients to the surface, and warmer-than-normal ocean temperatures. Climate change has likely made these conditions more favourable.

Algal blooms are complex and evolve over time. Scientists are studying this bloom to improve early detection and monitoring.

This Glance` widget displays the latest algal bloom information for a selected beach using data from [Beachsafe](https://beachsafe.org.au).

## Features

- Shows whether the beach has been cleaned.
- Indicates the presence of abnormal foam.
- Indicates abnormal water colour in the water.
- Last update timestamp displayed dynamically using relative time.
- Colour-coded indicators for easy visual recognition:
  - **Positive (green)**: Good conditions
  - **Negative (red)**: Potential issues


## API URL for beach

Use the developer tools inside your browser of choice to get the api endpoint for your beach.

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
