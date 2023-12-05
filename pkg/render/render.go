package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/cristian-moreno-ruiz/go-booking/pkg/config"
	"github.com/cristian-moreno-ruiz/go-booking/pkg/models"
)

var app *config.AppConfig

// NewTemplates sets config for render package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(data *models.TemplateData) *models.TemplateData {
	return data
}

func RenderTemplate(w http.ResponseWriter, tmpl string, data *models.TemplateData) {
	var cache map[string]*template.Template
	if app.UseCache {
		// get template cache from AppConfig
		cache = app.TemplateCache
	} else {
		cache, _ = CreateTemplateCache()
	}

	// get requested template from cache
	t, ok := cache[tmpl]

	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	// render the template
	buf := new(bytes.Buffer)
	data = AddDefaultData(data)
	err := t.Execute(buf, data)

	if err != nil {
		log.Println(err)
	}

	// Write to response writer
	_, err = buf.WriteTo(w)

	if err != nil {
		fmt.Println("error writing template to response: ", err)
		return
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	// myCache :=make(map[string]*template.Template)
	myCache := map[string]*template.Template{}

	// get all templates
	pages, err := filepath.Glob("./templates/*page.tmpl")

	if err != nil {
		return myCache, err
	}

	// range through all files
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)

		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}

/*
OLD APPROACH
*/
var tc = make(map[string]*template.Template)

func OldRenderTemplate(w http.ResponseWriter, t string) {
	var tmpl *template.Template
	var err error

	_, exists := tc[t]
	if !exists {
		log.Println("Creating and caching template")
		err = oldCreateTemplateCache(t)
		if err != nil {
			log.Println(err)
		}
	} else {
		log.Println("Using cached template")
	}

	tmpl = tc[t]
	if err != nil {
		log.Println(err)
	}
	err = tmpl.Execute(w, nil)
}

func oldCreateTemplateCache(t string) error {
	templates := []string{
		fmt.Sprintf("./templates/%s", t),
		"./templates/base.layout.tmpl",
	}

	tmpl, err := template.ParseFiles(templates...)

	if err != nil {
		return err
	}

	tc[t] = tmpl
	return nil
}
