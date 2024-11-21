package templates

import (
	"text/template"

	nginfier "github.com/vanvanni/ginfier"
	"github.com/vanvanni/ginfier/internal/logger"
)

func ReverseHost() (*template.Template, error) {
	logger.Info("Loading: templates/reverse-host.tmpl")
	tmplContent, err := nginfier.Templates().ReadFile("templates/reverse-host.tmpl")
	if err != nil {
		return nil, err
	}

	return template.New("nginx.conf").Parse(string(tmplContent))
}
