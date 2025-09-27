<p align="center"><img src="assets/logo.png"></p>
<h1 align="center">Community Widgets</h1>

<p align="center">
  <a href="GALLERY.md">Gallery</a> •
  <a href="CONTRIBUTING.md">Contributing</a> •
  <a href="#how-to-use">How to use</a> •
  <a href="#faq">FAQ</a>
</p>

<p align="center">
A collection of custom widgets for <a href="https://github.com/glanceapp/glance">Glance</a> made by the<br> community using the <code>custom-api</code> and <code>extension</code> widgets
</p>

<br>

## Custom API Widgets
{{- define "widget-list-item" }}
* [{{ .Title }}](widgets/{{ .Directory }}/README.md) - {{ .Description | toLowerFirst | trimSuffix "." }}{{ if .Author }} (by @{{ .Author }}){{ end }}
{{- end }}

### Newly added
{{- range .WidgetsSortedByTimeAdded }}
{{- template "widget-list-item" . }}
{{- end }}

### All
{{- range .WidgetsSortedByTitle }}
{{- template "widget-list-item" . }}
{{- end }}

<br>

## Extension Widgets
> [!WARNING]
>
> Extension widgets are not actively monitored by the maintainers of Glance, use them at your own risk.
{{- range .ExtensionsSortedByTitle }}
* [{{ .Title }}]({{ .URL }}) - {{ .Description | toLowerFirst | trimSuffix "." }}{{ if .Author }} (by @{{ .Author }}){{ end }}
{{- end }}

<br>

## How to use
For simpler widgets you can simply copy their code into your `glance.yml` as you would with any other widget, then add environment variables for any URL's and API keys if necessary.

For more complex widgets that span across hundreds of lines, it may be trickier to get their indentation right, so it's easier to place them in a separate `yml` file, configure them there, then include that file in your `glance.yml` like such:

```yml
widgets:
  - $include: immich-stats.yml
```

<br>

## FAQ
<details>
<summary><strong>Are the widgets maintained?</strong></summary>

The maintainers of Glance are not responsible for the maintenance of these widgets. Instead, the author of each widget is responsible for maintaining and responding to issues and pull requests related to that widget.
</details>

<details>
<summary><strong>What's the difference between a <code>custom-api</code> and an <code>extension</code> widget?</strong></summary>

Custom API widgets are much easier to setup and usually only require a copy-paste into your config. Extension widgets are a bit more involved and require running a separate server or Docker container.
</details>

<details>
<summary><strong>Are the widgets safe to use?</strong></summary>

The `custom-api` widgets in this repository have been vetted by the maintainers of Glance so they are safe to use, however they may still have bugs, be visually inconsistent with the rest of Glance, or in some cases have poor performance.
</details>
