{{- define "entry" }}{{ if .Title }}<div><p align="center">[{{ .Title }}](widgets/{{ .Directory }}/README.md)<br>by [@{{ .Author }}](https://github.com/{{ .Author }})<p><p align="center">[![](widgets/{{ .Directory }}/{{ .Preview }})](/widgets/{{ .Directory }}/README.md)</p></div>{{ end }}{{- end }}
| | | |
| - | - | - |
{{- range .WidgetsGroupedSortedByTitle }}
| {{ template "entry" (index . 0) }} | {{ template "entry" (index . 1) }} | {{ template "entry" (index . 2) }} |
{{- end }}
