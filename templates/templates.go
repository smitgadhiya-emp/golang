package templates

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
)

//go:embed *.html
var files embed.FS

// tmpl holds every *.html file in this directory, each addressable by its
// file name (e.g. "welcome.html", "otp.html").
var tmpl = template.Must(template.ParseFS(files, "*.html"))

// Render executes the named template with the given data and returns the
// resulting HTML string, ready to be used as an email body.
func Render(name string, data any) (string, error) {
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, name, data); err != nil {
		return "", fmt.Errorf("render email template %q: %w", name, err)
	}
	return buf.String(), nil
}
