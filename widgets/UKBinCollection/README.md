# GlanceUKBinWidget
UK Bin Collection Widget for Glance Dashboard - Credit to Robert Bradley for the original project https://github.com/robbrad/UKBinCollectionData

![Fixture widget preview](./preview.png)

# UK Bin Collection

A custom widget for [Glance](https://github.com/glanceapp/glance) that displays upcoming rubbish and recycling collection dates for UK households. 

It dynamically highlights the different bin types using native Glance theme colors and provides a clean, vertical layout with standardized icons.

## Prerequisites

This widget does not scrape council websites directly. Instead, it relies on the fantastic [UKBinCollectionData](https://github.com/robbrad/UKBinCollectionData) project to parse the data and serve it as a local API.

Before using this widget, you must have the **UK Bin Collection API Server** running (typically via Docker). See their official repository for setup instructions.

## Environment Variables

To keep your configuration clean and avoid hardcoding local IP addresses, this widget requires one environment variable to be defined in your Glance `docker-compose.yml` or `.env` file:

| Variable | Description | Example |
| :--- | :--- | :--- |
| `UK_BIN_COLLECTION_URL` | The full URL to your specific council's endpoint on your local API server, including any required query parameters (like `uprn` or `postcode`). | `http://X.X.X.X:8080/api/bin_collection/CouncilName?uprn=XXXXXXXXXXXXXXXXX` |

## Configuration

Add the following to your `home.yml` file:

```yaml
- type: custom-api
  title: Bin Collections
  url: ${UK_BIN_COLLECTION_URL}
  cache: 24h
  template: |
    <div class="space-y-3">
      {{ range .JSON.Array "bins" }}
      <div class="block p-3 rounded bg-white/5">
        
        <div class="block text-left w-full mb-3">
          <div class="text-[10px] font-bold uppercase opacity-60 mb-1">
            {{ if eq (.String "type") "Grey bin, small electrical items and food bin" }}Rubbish & Food
            {{ else }}Recycling & Food{{ end }}
          </div>
          
          <div class="text-sm font-bold whitespace-normal leading-snug mb-3 
            {{ if eq (.String "type") "Grey bin, small electrical items and food bin" }}color-subdue
            {{ else if eq (.String "type") "Food bin, green bin, textiles and garden waste bin" }}color-positive
            {{ else }}color-highlight{{ end }}">
            {{ .String "type" }}
          </div>
          
          <div class="inline-block px-2 py-1 color-primary rounded text-xs font-bold tracking-wide">
            {{ .String "collectionDate" }}
          </div>
        </div>
        
        <div class="block text-3xl text-left w-full">
          {{ if eq (.String "type") "Grey bin, small electrical items and food bin" }}&#128465;
          {{ else if eq (.String "type") "Food bin, green bin, textiles and garden waste bin" }}&#9851;
          {{ else }}&#127810;{{ end }}
        </div>
      </div>
      {{ end }}
    </div>
