package main

import (
	"html/template"
	"io/fs"
	"path/filepath"

	"github.com/PenguGG0/go-cli/ui"
)

type templateData struct {
	Body template.HTML
}

func newTemplateCache() (map[string]*template.Template, error) {
	// Create a map to hold the templates act as the cache
	templateCache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/mdp/*.gohtml")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// Get the base name of the file (error.g., "base.gohtml")
		name := filepath.Base(page)

		patterns := []string{
			"html/mdp/base.gohtml",
		}

		tmpl := template.New(name)
		tmpl, err = tmpl.ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		// Add the parsed template to the map of templates
		templateCache[name] = tmpl
	}

	return templateCache, nil
}
