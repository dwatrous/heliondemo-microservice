package main

import (
	"bytes"
	"fmt"
	"html/template"
)

type FormElement struct {
	Input   string            `json:"input", bson:"input"`
	Name    string            `json:"name", bson:"name"`
	Options map[string]string `json:"options", bson:"options"`
}

var (
	formBase = `<!DOCTYPE html>
<html>
	<head>
	</head>
	<body>
		<div id="survey">
			<form method="post">
			{{ range . }}
				{{ printinput . }}
			{{ end }}
			<input type="submit" value="Submit" />
			</form>
		</div>
	</body>
</html>
`

	selectBase = `<select name="{{ .Name }}">
{{ range $label, $value := .Options }}
	<option value="{{ $value }}">{{ $label }}</option>
{{ end }}
</select>
`

	textBase = `<input type="text" name="{{ .Name }}">`

	formTmpl   *template.Template
	selectTmpl *template.Template
	textTmpl   *template.Template
)

func init() {
	formTmpl = template.Must(template.New("form").Funcs(template.FuncMap{"printinput": printInput}).Parse(formBase))
	selectTmpl = template.Must(template.New("select").Parse(selectBase))
	textTmpl = template.Must(template.New("text").Parse(textBase))
}

func printInput(f *FormElement) (template.HTML, error) {
	var err error
	buf := new(bytes.Buffer)

	switch f.Input {
	case "text":
		err = textTmpl.Execute(buf, f)
	case "select":
		err = selectTmpl.Execute(buf, f)
	default:
		return "", fmt.Errorf("invalid form input type %q", f.Input)
	}

	return template.HTML(buf.String()), err
}
