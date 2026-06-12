package handlers

import (
	"encoding/json"
	"net/http"
	"text/template"
	"bytes"
)

type MefetchInput struct {
	Username string `json:"username"`
	Hostname string `json:"hostname"`
	OS 		 string `json:"os"`
	Kernel	 string `json:"kernel"`
	Shell	 string `json:"shell"`
	Editor	 string `json:"editor"`
	Colors	 string `json:"colors"`
}

var mefetchTemplate = `
	{{.Username}}@{{.Hostname}}
	--------------------------
	OS		{{.OS}}
	Kernel	{{.Kernel}}
	Shell	{{.Shell}}
	Editor	{{.Editor}}
	Colors	{{.Colors}}
`

func Generate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input MefetchInput
	var result string

	ok := run(w,
		func() error { return decodeInput(r, &input)},
		func() error { return renderMefetch(input, &result)},
	)

	if !ok {
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(result))
}

func decodeInput(r *http.Request, input *MefetchInput) error {
	return json.NewDecoder(r.Body).Decode(input)
}

func renderMefetch(input MefetchInput, result *string) error {
	tmpl, err := template.New("mefetch").Parse(mefetchTemplate)
	if err != nil {
		return err
	}

	var buf bytes.Buffer

	if err = tmpl.Execute(&buf, input); err != nil {
		return err
	}

	*result = buf.String()
	return nil
}