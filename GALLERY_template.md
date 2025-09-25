{{- define "entry" }}{{ if .Title }}
    <p align="center"><a href="widgets/{{ .Directory }}/README.md">{{ .Title }}</a><br>by <a href="https://github.com/{{ .Author }}">@{{ .Author }}</a><p>
    <p align="center"><a href="widgets/{{ .Directory }}/README.md"><img src="widgets/{{ .Directory }}/{{ .Preview }}"></a></p>
{{ end }}{{- end }}
<table>
{{- range .WidgetsGroupedSortedByTitle }}
  <tr>
    <td valign="top">{{ template "entry" (index . 0) }}</td>
    <td valign="top">{{ template "entry" (index . 1) }}</td>
    <td valign="top">{{ template "entry" (index . 2) }}</td>
  </tr>
{{- end }}
</table>
