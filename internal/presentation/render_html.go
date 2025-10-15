package presentation

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"slices"
)

const baseTemplateDir = "templates"

var pageTemplates []string = []string{
	"home",
	"login",
	"register",
	"group",
	"dashboard",
}

var partialTemplates []string = []string{
	"groups",
}

var baseTemplates []string = []string{
	"base",
	"navbar",
}

type CompliedTemplates map[string]*template.Template

type Templates struct {
	templates *CompliedTemplates
}

func NewHtmlPresentationService() *Templates {
	templates, err := compileTemplates(baseTemplates, pageTemplates, partialTemplates)
	if err != nil {
		log.Fatalln("Failed to compile templates:", err.Error())
	}
	return &Templates{templates}
}

func (t *Templates) Present(w http.ResponseWriter, r *http.Request, name string, payload any) error {
	tmpl, exists := (*t.templates)[name]
	if !exists {
		return errors.New(name + " template not found")
	}

	if slices.Contains(partialTemplates, name) {
		log.Println("This is a partial")
		return tmpl.ExecuteTemplate(w, name, payload)
	}

	tmplName := "base"

	if r.Header.Get("HX-Request") == "true" {
		tmplName = "content"
	}

	return tmpl.ExecuteTemplate(w, tmplName, payload)
}

func compileTemplates(bases []string, pages []string, partials []string) (*CompliedTemplates, error) {
	var templates = make(CompliedTemplates)
	for _, p := range pages {
		var allPages = make([]string, 0)
		for _, b := range bases {
			allPages = append(allPages, fmt.Sprintf("%s/%s.html", baseTemplateDir, b))
		}
		var err error
		allPages = append(allPages, fmt.Sprintf("%s/%s.html", baseTemplateDir, p))
		templates[p], err = template.ParseFiles(
			allPages...,
		)
		if err != nil {
			return nil, err
		}
	}

	var err error
	for _, partial := range partials {
		partialFile := fmt.Sprintf("%s/%s.html", baseTemplateDir, partial)
		templates[partial], err = template.ParseFiles(partialFile)
		if err != nil {
			return nil, err
		}
	}
	return &templates, nil
}
