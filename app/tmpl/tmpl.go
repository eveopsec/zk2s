// Package tmpl loads the template file(s) and exposes them for execution.
package tmpl

import (
	"log"
	"text/template"
)

// T is the template collection, used in execution of templates.
var T *template.Template

func Init(tmplFilePath string) error {
	var err error
	if tmplFilePath == "" {
		log.Println("[WARN] Template file path not specified in configuration; using default 'template.tmpl'...")
		tmplFilePath = "response.tmpl"
	}
	T, err = template.ParseGlob(tmplFilePath)
	return err
}
