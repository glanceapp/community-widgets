package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
	"text/template"
	"time"

	"gopkg.in/yaml.v3"
)

const widgetsDir = "../widgets"

var registeredWidgets = make(map[string]*widget)

type widget struct {
	Title       string    `json:"title" yaml:"title"`
	Description string    `json:"description" yaml:"description"`
	Author      string    `json:"author" yaml:"author"`
	Directory   string    `json:"directory" yaml:"-"`
	ReadmeHash  string    `json:"readme_hash" yaml:"-"`
	TimeAdded   time.Time `json:"time_added" yaml:"-"`
	TimeUpdated time.Time `json:"time_updated" yaml:"-"`
}

type extension struct {
	Title       string `yaml:"title"`
	URL         string `yaml:"url"`
	Author      string `yaml:"author"`
	Description string `yaml:"description"`
}

func main() {
	// Load already registered widgets from widgets.json
	loadRegisteredWidgets()

	// Scan widgets directory for available widgets
	widgetDirs, err := os.ReadDir(widgetsDir)
	if err != nil {
		log.Fatalf("Failed to read widgets directory: %v", err)
	}

	for _, entry := range widgetDirs {
		if !entry.IsDir() {
			continue
		}

		loadMetaFileForWidget(entry.Name())
	}

	// Prepare sorted lists
	sortedByTitle := make([]widget, 0, len(registeredWidgets))
	for _, w := range registeredWidgets {
		sortedByTitle = append(sortedByTitle, *w)
	}

	slices.SortStableFunc(sortedByTitle, func(a, b widget) int {
		return strings.Compare(a.Title, b.Title)
	})

	sortedByTimeAdded := make([]widget, len(sortedByTitle))
	copy(sortedByTimeAdded, sortedByTitle)

	slices.SortStableFunc(sortedByTimeAdded, func(a, b widget) int {
		return b.TimeAdded.Compare(a.TimeAdded)
	})

	if len(sortedByTimeAdded) > 5 {
		sortedByTimeAdded = sortedByTimeAdded[:5]
	}

	// Save updated widgets.json
	data, err := json.MarshalIndent(sortedByTitle, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal widgets to JSON: %v", err)
	}

	err = os.WriteFile("widgets.json", data, 0644)
	if err != nil {
		log.Fatalf("Failed to write widgets.json: %v", err)
	}

	// Load extension widgets
	contents, err := os.ReadFile("../widgets/extensions.yml")
	if err != nil {
		log.Fatalf("Failed to read extensions.yml: %v", err)
	}

	extensions := []extension{}
	err = yaml.Unmarshal(contents, &extensions)
	if err != nil {
		log.Fatalf("Failed to parse extensions.yml: %v", err)
	}

	slices.SortStableFunc(extensions, func(a, b extension) int {
		return strings.Compare(a.Title, b.Title)
	})

	// Generate README.md
	readmeTemplate := parseTemplate("../README_template.md")

	readmeTemplateData := struct {
		WidgetsSortedByTimeAdded []widget
		WidgetsSortedByTitle     []widget
		ExtensionsSortedByTitle  []extension
	}{
		WidgetsSortedByTimeAdded: sortedByTimeAdded,
		WidgetsSortedByTitle:     sortedByTitle,
		ExtensionsSortedByTitle:  extensions,
	}

	readmeContents := readmeTemplate(readmeTemplateData)
	err = os.WriteFile("../README.md", readmeContents, 0644)
	if err != nil {
		log.Fatalf("Failed to write README.md: %v", err)
	}
}

func loadRegisteredWidgets() {
	entries, err := os.ReadFile("widgets.json")
	if err != nil {
		log.Fatalf("Failed to read widgets.json: %v", err)
	}

	var widgets []widget
	err = json.Unmarshal(entries, &widgets)
	if err != nil {
		log.Fatalf("Failed to parse widgets.json: %v", err)
	}

	for i := range widgets {
		w := &widgets[i]

		if w.Directory == "" {
			continue
		}

		if _, err := os.Stat(widgetPath(w.Directory, "")); os.IsNotExist(err) {
			continue
		}

		registeredWidgets[w.Directory] = w
	}
}

func loadMetaFileForWidget(widgetDir string) {
	contents, err := os.ReadFile(widgetPath(widgetDir, "meta.yml"))
	if err != nil {
		log.Fatalf("Failed to read meta file for widget %s: %v", widgetDir, err)
	}

	var meta widget
	err = yaml.Unmarshal(contents, &meta)
	if err != nil {
		log.Fatalf("Failed to parse meta file for widget %s: %v", widgetDir, err)
	}

	if meta.Title == "" {
		log.Fatalf("Widget %s is missing title in meta.yml", widgetDir)
	}

	if meta.Description == "" {
		log.Fatalf("Widget %s is missing description in meta.yml", widgetDir)
	}

	readmeHash := computeFileHash(widgetPath(widgetDir, "README.md"))

	w, ok := registeredWidgets[widgetDir]
	if !ok {
		meta.TimeAdded = time.Now().UTC()
		meta.TimeUpdated = time.Now().UTC()
		meta.ReadmeHash = readmeHash
		meta.Directory = widgetDir
		registeredWidgets[widgetDir] = &meta
	} else {
		w.Title = meta.Title
		w.Description = meta.Description
		w.Author = meta.Author

		if w.ReadmeHash != readmeHash {
			w.ReadmeHash = readmeHash
			w.TimeUpdated = time.Now().UTC()
		}
	}
}

func widgetPath(widgetDir, file string) string {
	return fmt.Sprintf("%s/%s/%s", widgetsDir, widgetDir, file)
}

func computeFileHash(path string) string {
	contents, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to open file %s: %v", path, err)
	}

	h := sha1.New()
	h.Write(contents)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func parseTemplate(path string) func(data any) []byte {
	contents, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read template file %s: %v", path, err)
	}

	funcs := template.FuncMap{
		"toLowerFirst": func(s string) string {
			if len(s) < 2 {
				return s
			}

			firstLowercased := strings.ToLower(string(s[0]))
			return firstLowercased + s[1:]
		},
		"trimSuffix": func(suffix, str string) string {
			return strings.TrimSuffix(str, suffix)
		},
	}

	tmpl, err := template.New("").Funcs(funcs).Parse(string(contents))
	if err != nil {
		log.Fatalf("Failed to parse template file %s: %v", path, err)
	}

	return func(data any) []byte {
		var sb strings.Builder
		err := tmpl.Execute(&sb, data)
		if err != nil {
			log.Fatalf("Failed to execute template %s: %v", path, err)
		}
		return []byte(sb.String())
	}
}
