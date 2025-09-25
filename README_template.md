This is a collection of custom widgets for [Glance](https://github.com/glanceapp/glance) made by the community using the `custom-api` and `extension` widgets.

The `custom-api` widgets in this repository have been vetted by the maintainers of Glance, however they are not responsible for the maintenance of these widgets. The author of each widget is responsible for maintaining and responding to issues and pull requests related to that widget.

Before a widget is added to this repository it must receive a vouch which can be done by anyone. If you would like to vouch for a widget, please leave a comment or a 👍 on its corresponding pull request.

To add your widget to the list, please read the [contribution guidelines](CONTRIBUTING.md).

### What's the difference between a `custom-api` and `extension` widget?

Custom API widgets are much easier to setup and usually only require a copy-paste into your config. Extension widgets are a bit more involved and require running a separate server or Docker container.

## Custom API Widgets

A gallery with screenshots of all widgets can be found [here](GALLERY.md).

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

## Extension Widgets

> [!WARNING]
>
> Extension widgets are not actively monitored by the maintainers of Glance, use them at your own risk.
{{- range .ExtensionsSortedByTitle }}
* [{{ .Title }}]({{ .URL }}) - {{ .Description | toLowerFirst | trimSuffix "." }}{{ if .Author }} (by @{{ .Author }}){{ end }}
{{- end }}
