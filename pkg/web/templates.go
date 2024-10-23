package web

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"html/template"
	"path/filepath"
)

const baseTemplateName = "_layout.gohtml"

func (v *views) renderTemplate(c echo.Context, templateName string, viewModel interface{}) error {
	var err error

	tmpl, err := template.ParseFS(fs,
		filepath.Join("templates", templateName),
		filepath.Join("templates", baseTemplateName),
	)
	if err != nil {
		return errors.Wrap(err, "failed to parse template")
	}

	return tmpl.ExecuteTemplate(
		c.Response().Writer,
		baseTemplateName,
		viewModel,
	)
}
