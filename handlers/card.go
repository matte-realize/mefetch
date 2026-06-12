package handlers

import (
	"bytes"
	"strings"
	"text/template"
)

type Field struct {
	Label string
	Dots  string
	Value string
	Color string
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

var svgTemplate = `<svg xmlns="http://www.w3.org/2000/svg" width="860" height="490" viewBox="0 0 860 490">
  <style>
    .bg     { fill: {{.Background}}; }
    .border { fill: none; stroke: #30363d; stroke-width: 1; }
    .title  { fill: #ffffff; font-family: monospace; font-size: 12px; font-weight: bold; }
    .sep    { fill: #30363d; }
    .key    { fill: {{.KeyColor}}; font-family: monospace; font-size: 11px; }
    .val    { fill: {{.TextColor}}; font-family: monospace; font-size: 11px; }
    .ascii  { fill: {{.TextColor}}; font-family: monospace; font-size: 5.8px; }
    .dots   { fill: #444d56; font-family: monospace; font-size: 11px; }
    .green  { fill: #3fb950; font-family: monospace; font-size: 11px; }
    .red    { fill: #f85149; font-family: monospace; font-size: 11px; }
  </style>

  <rect class="bg" width="860" height="490" rx="8" />
  <rect class="border" x="0.5" y="0.5" width="859" height="489" rx="8" />

  <!-- ASCII art left side -->
  {{range $i, $line := .AsciiLines}}
  <text class="ascii" x="15" y="{{asciiY $i}}">{{$line}}</text>
  {{end}}

  <!-- Divider -->
  <rect class="sep" x="290" y="12" width="1" height="466" />

  <!-- Title -->
  <text class="title" x="308" y="30">{{.Username}}@{{.Hostname}}</text>

  <!-- Title separator line stretching to right edge -->
  <text class="sep" x="308" y="44" fill="#444d56" font-family="monospace" font-size="11px">----------------------------------------------------------------</text>

  <!-- Fields -->
  {{range $i, $field := .Fields}}
  <text class="key"  x="308" y="{{fieldY $i}}">{{$field.Label}}:</text>
  <text class="dots" x="410" y="{{fieldY $i}}">{{$field.Dots}}</text>
  {{if eq $field.Color "green"}}
  <text class="green" x="600" y="{{fieldY $i}}">{{$field.Value}}</text>
  {{else if eq $field.Color "red"}}
  <text class="red"   x="600" y="{{fieldY $i}}">{{$field.Value}}</text>
  {{else}}
  <text class="val"   x="600" y="{{fieldY $i}}">{{$field.Value}}</text>
  {{end}}
  {{end}}

</svg>`

func renderCard(input CardInput, result *string) error {
	funcMap := template.FuncMap{
		"asciiY": func(i int) int { return 14 + (i * 7) },
		"fieldY": func(i int) int { return 58 + (i * 17) },
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

func makeField(label, value, color string) Field {
	return Field{
		Label: label,
		Dots:  makeDots(label),
		Value: value,
		Color: color,
	}
}

func makeDots(label string) string {
	total := 30
	dots := total - len(label)
	if dots < 3 {
		dots = 3
	}
	result := strings.Repeat(".", dots)
	return result
}