package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
)

// subjects is a mapping between language and the
// corresponding email subject.
var subjects map[string]string

// templates are the available email templates split by
// language.
var templates *template.Template

// _templateFiles hold all templates for notifications.
//
// DO NOT USE.
//
//go:embed templates/*.html
var _templateFiles embed.FS

// _subjectsFile is a JSON containing a mapping between
// language and the corresponding email subject.
//
// DO NOT USE.
//
//go:embed templates/subjects.json
var _subjectsFile []byte

func init() {
	// funcsMap are all the functions required for the
	// templates to render correctly.
	var funcsMap = template.FuncMap{
		"round": func(x float32) string { return fmt.Sprintf("%.2f", x) },
	}

	subjects = map[string]string{}
	if err := json.Unmarshal(_subjectsFile, &subjects); err != nil {
		panic(err)
	}

	templates = template.Must(
		template.New("emails").
			Funcs(funcsMap).
			ParseFS(_templateFiles, "templates/*.html"),
	)
}
