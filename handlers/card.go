package handlers

import (
	"text/template"
	"bytes"
)

type Field struct {
	Label string
	Dots  string
	Value string
}

type CardInput struct {
	Username   string
	Hostname   string
	Background string
	KeyColor   string
	TextColor  string
	AsciiLines []string
	Fields     []Field
}

var svgTemplate = `<svg xmlns="http://www.w3.org/2000/svg" width="900" height="520" viewBox="0 0 900 520">
  <style>
    .bg     { fill: {{.Background}}; }
    .border { fill: none; stroke: #30363d; stroke-width: 1; }
    .title  { fill: #ffffff; font-family: monospace; font-size: 13px; font-weight: bold; }
    .sep    { fill: #30363d; }
    .key    { fill: {{.KeyColor}}; font-family: monospace; font-size: 11px; }
    .val    { fill: {{.TextColor}}; font-family: monospace; font-size: 11px; }
    .ascii  { fill: {{.TextColor}}; font-family: monospace; font-size: 6.5px; }
    .dots   { fill: #444d56; font-family: monospace; font-size: 11px; }
    .green  { fill: #3fb950; font-family: monospace; font-size: 11px; }
    .red    { fill: #f85149; font-family: monospace; font-size: 11px; }
    .dim    { fill: #444d56; font-family: monospace; font-size: 11px; }
  </style>

  <rect class="bg" width="900" height="520" rx="8" />
  <rect class="border" x="0.5" y="0.5" width="899" height="519" rx="8" />

  {{range $i, $line := .AsciiLines}}
  <text class="ascii" x="20" y="{{asciiY $i}}">{{$line}}</text>
  {{end}}

  <rect class="sep" x="310" y="15" width="1" height="490" />

  <text class="title" x="330" y="35">{{.Username}}@{{.Hostname}}</text>
  <rect class="sep" x="330" y="45" width="550" height="1" />

  {{range $i, $field := .Fields}}
  <text class="key"  x="330" y="{{fieldY $i}}">{{$field.Label}}:</text>
  <text class="dots" x="430" y="{{fieldY $i}}">{{$field.Dots}}</text>
  <text class="val"  x="560" y="{{fieldY $i}}">{{$field.Value}}</text>
  {{end}}
</svg>`

func renderCard(input CardInput, result *string) error {
	funcMap := template.FuncMap{
		"asciiY": func(i int) int { return 20 + (i * 8) },
		"fieldY": func(i int) int { return 65 + (i * 20) },
	}

	tmpl, err := template.New("card").Funcs(funcMap).Parse(svgTemplate)

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

func makeDots(label string) string {
	total := 30
	dots := total - len(label)
	if dots < 3 {
		dots = 3
	}
	result := ""
	for i := 0; i < dots; i++ {
		result += "."
	}
	return result
}

func defaultInput(input CardInput) CardInput {
	if input.Background == "" {
		input.Background = "#0d1117"
	}
	if input.KeyColor == "" {
		input.KeyColor = "#58a6ff"
	}
	if input.TextColor == "" {
		input.TextColor = "#cdd9e5"
	}
	return input
}