package render

import (
	"bytes"
	"fmt"
	"github.com/zubsingh/bookings/cmd/pkg/config"
	"github.com/zubsingh/bookings/cmd/pkg/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var functionstemp = template.FuncMap{}

var app *config.AppConfig

func SetConfig(ac *config.AppConfig) {
	app = ac
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

func RenderTemplate(rw http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template

	if app.UseCache {
		// get the config cache from app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not able to get template from template cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td)
	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(rw)
	if err != nil {
		fmt.Println("Error writing template to browser", err)
	}

	// parsedTemplate, _ := template.ParseFiles(tmpl)
	// err := parsedTemplate.Execute(rw, nil)
	// if err != nil {
	// 	fmt.Println("error parsing: ", err)
	// 	return
	// }
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.html")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		//fmt.Println("Page is currently ", page)
		ts, err := template.New(name).Funcs(functionstemp).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.html")

		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			//fmt.Println(matches)
			ts, err = ts.ParseGlob("./templates/*.layout.html")
			if err != nil {
				//fmt.Println(err)
				return myCache, err
			}
		}

		myCache[name] = ts
		//fmt.Println(myCache)
	}
	//fmt.Println(myCache)
	return myCache, nil
}
