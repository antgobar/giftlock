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
	templates, err := compileTemplates(baseTemplateDir, baseTemplates, pageTemplates, partialTemplates)
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
		return tmpl.ExecuteTemplate(w, name, payload)
	}

	tmplName := "base"

	if r.Header.Get("HX-Request") == "true" {
		tmplName = "content"
	}

	return tmpl.ExecuteTemplate(w, tmplName, payload)
}

func compileTemplates(dir string, bases []string, pages []string, partials []string) (*CompliedTemplates, error) {
	pageTemplates, err := compilePages(dir, bases, pages)
	if err != nil {
		return nil, err
	}

	partialTemates, err := compilePartials(dir, partials)
	if err != nil {
		return nil, err
	}

	var templates = make(CompliedTemplates)

	for k, v := range *pageTemplates {
		templates[k] = v
	}

	for k, v := range *partialTemates {
		templates[k] = v
	}

	return &templates, nil
}

func compilePages(dir string, bases []string, pages []string) (*CompliedTemplates, error) {
	var templates = make(CompliedTemplates)
	for _, p := range pages {
		var allPages = make([]string, 0)
		for _, b := range bases {
			allPages = append(allPages, fmt.Sprintf("%s/%s.html", dir, b))
		}
		var err error
		allPages = append(allPages, fmt.Sprintf("%s/%s.html", dir, p))
		templates[p], err = template.ParseFiles(
			allPages...,
		)
		if err != nil {
			return nil, err
		}
	}
	return &templates, nil
}

func compilePartials(baseTemplateDir string, partials []string) (*CompliedTemplates, error) {
	var templates = make(CompliedTemplates)
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
