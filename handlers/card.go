package handlers

import (
	"net/http"
	"text/template"
	"bytes"
)

type CardInput struct {
	Username string
	Hostname string
	OS 		 string
	Kernel	 string
	Shell	 string
	Editor	 string
	Uptime   string
	Colors	 string
}

var svgTemplate = `<svg xmlns="http://www.w3.org/2000/svg" width="495" height="195" viewBox="0 0 495 195">
  <style>
    .bg    { fill: #0d1117; }
    .border { fill: none; stroke: #30363d; stroke-width="1"; }
    .title  { fill: #ffffff; font-family: monospace; font-size: 14px; font-weight: bold; }
    .sep    { fill: #30363d; }
    .key    { fill: #58a6ff; font-family: monospace; font-size: 12px; }
    .val    { fill: #cdd9e5; font-family: monospace; font-size: 12px; }
  </style>

  <rect class="bg" width="495" height="195" rx="6" />

  <rect class="border" x="0.5" y="0.5" width="494" height="194" rx="6" />

  <text class="title" x="25" y="35">{{.Username}}@{{.Hostname}}</text>

  <rect class="sep" x="25" y="45" width="445" height="1" />

  <text class="key" x="25" y="65">OS</text>
  <text class="val" x="120" y="65">{{.OS}}</text>

  <text class="key" x="25" y="85">Kernel</text>
  <text class="val" x="120" y="85">{{.Kernel}}</text>

  <text class="key" x="25" y="105">Shell</text>
  <text class="val" x="120" y="105">{{.Shell}}</text>

  <text class="key" x="25" y="125">Editor</text>
  <text class="val" x="120" y="125">{{.Editor}}</text>

  <text class="key" x="25" y="145">Uptime</text>
  <text class="val" x="120" y="145">{{.Uptime}}</text>

  <text class="key" x="25" y="170">Colors</text>
  <rect x="120" y="158" width="18" height="14" rx="2" fill="#ff5f56"/>
  <rect x="142" y="158" width="18" height="14" rx="2" fill="#ffbd2e"/>
  <rect x="164" y="158" width="18" height="14" rx="2" fill="#27c93f"/>
  <rect x="186" y="158" width="18" height="14" rx="2" fill="#58a6ff"/>
  <rect x="208" y="158" width="18" height="14" rx="2" fill="#cdd9e5"/>
  <rect x="230" y="158" width="18" height="14" rx="2" fill="#ffffff"/>
</svg>`

func Card(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	input := CardInput{
		Username: r.URL.Query().Get("username"),
		Hostname: r.URL.Query().Get("hostname"),
		OS:       r.URL.Query().Get("os"),
		Kernel:   r.URL.Query().Get("kernel"),
		Shell:    r.URL.Query().Get("shell"),
		Editor:   r.URL.Query().Get("editor"),
		Uptime:   r.URL.Query().Get("uptime"),
	}

	var result string
	ok := run(w,
		func() error { return renderCard(input, &result) },
	)
	if !ok {
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Write([]byte(result))
}

func renderCard(input CardInput, result *string) error {
	tmpl, err := template.New("card").Parse(svgTemplate)

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