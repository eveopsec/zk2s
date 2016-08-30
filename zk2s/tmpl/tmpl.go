// Package tmpl loads the template file(s) and exposes them for execution.
package tmpl

import (
	"text/template"

	"github.com/urfave/cli"
)

// T is the template collection, used in execution of templates.
var T *template.Template

func Init(c *cli.Context) error {
	var err error
	pattern := c.String("template")
	T, err = template.ParseGlob(pattern)
	return err
}
