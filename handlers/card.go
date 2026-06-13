package handlers

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

const (
	colStart  = 308
	colValue  = 845
	charWidth = 7
	fieldsTop = 47
)

type Field struct {
	Label string
	Dots  string
	Value string
	Color string
	Type  string
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

var svgTemplate = `<svg xmlns="http://www.w3.org/2000/svg" width="860" height="{{svgHeight .Fields}}" viewBox="0 0 860 {{svgHeight .Fields}}">
  <style>
    .bg        { fill: {{.Background}}; }
    .border    { fill: none; stroke: #30363d; stroke-width: 1; }
    .title     { fill: #ffffff; font-family: monospace; font-size: 12px; font-weight: bold; }
    .titledash { fill: #444d56; font-family: monospace; font-size: 12px; }
    .divider   { fill: #ffffff; font-family: monospace; font-size: 12px; font-weight: bold; }
    .dividash  { fill: #444d56; font-family: monospace; font-size: 12px; }
    .sep       { fill: #30363d; }
    .key       { fill: {{.KeyColor}}; font-family: monospace; font-size: 11px; }
    .val       { fill: {{.TextColor}}; font-family: monospace; font-size: 11px; }
    .ascii     { fill: {{.TextColor}}; font-family: monospace; font-size: 5.8px; }
    .dots      { fill: #444d56; font-family: monospace; font-size: 11px; }
    .green     { fill: #3fb950; font-family: monospace; font-size: 11px; }
    .red       { fill: #f85149; font-family: monospace; font-size: 11px; }
  </style>

  <rect class="bg" width="860" height="{{svgHeight .Fields}}" rx="8" />
  <rect class="border" x="0.5" y="0.5" width="859" height="{{svgHeightBorder .Fields}}" rx="8" />

  {{range $i, $line := .AsciiLines}}
  <text class="ascii" x="150" y="{{asciiY $i}}" text-anchor="middle" xml:space="preserve">{{$line}}</text>
  {{end}}

  <rect class="sep" x="290" y="12" width="1" height="{{sepHeight .Fields}}" />

  <text class="title" x="308" y="30" textLength="{{rowWidth}}" lengthAdjust="spacingAndGlyphs">{{titleDisplay .Username .Hostname}} <tspan class="titledash">{{titleDashes .Username .Hostname}}</tspan></text>

  {{range $i, $field := .Fields}}
  {{if eq $field.Type "divider"}}
  <text class="divider" x="308" y="{{fieldY $i}}" textLength="{{rowWidth}}" lengthAdjust="spacingAndGlyphs">- {{$field.Label}} - <tspan class="dividash">{{dividerDashes $field.Label}}</tspan></text>
  {{else if eq $field.Type "spacer"}}
  {{else if eq $field.Type "halfspacer"}}
  {{else}}
  <text class="key" x="308" y="{{fieldY $i}}" textLength="{{rowWidth}}" lengthAdjust="spacingAndGlyphs">·{{$field.Label}}: <tspan class="dots">{{$field.Dots}}</tspan><tspan class="{{valClass $field.Color}}">{{$field.Value}}</tspan></text>
  {{end}}
  {{end}}

</svg>`

func renderCard(input CardInput, result *string) error {
	funcMap := template.FuncMap{
		"asciiY": func(i int) int {
			cardHeight := fieldsTop + fieldsHeight(input.Fields) + 6
			top := (cardHeight - ((len(input.AsciiLines) - 1) * 7)) / 2
			if top < 14 {
				top = 14
			}
			return top + (i * 7)
		},
		"fieldY": func(i int) int {
			y := fieldsTop
			for j := 0; j < i && j < len(input.Fields); j++ {
				y += rowHeight(input.Fields[j])
			}
			return y
		},
		"svgHeight": func(fields []Field) int {
			return fieldsTop + fieldsHeight(fields) + 6
		},
		"svgHeightBorder": func(fields []Field) int {
			return fieldsTop + fieldsHeight(fields) + 5
		},
		"sepHeight": func(fields []Field) int {
			return fieldsTop + fieldsHeight(fields) - 6
		},
		"titleDisplay": func(username, hostname string) string {
			return titleStr(username, hostname)
		},
		"titleDashes": func(username, hostname string) string {
			title := titleStr(username, hostname)
			totalCols := (colValue - colStart) / charWidth
			dashes := totalCols - len(title) - 1
			if dashes < 3 {
				dashes = 3
			}
			return strings.Repeat("-", dashes)
		},
		"dividerDashes": func(label string) string {
			totalCols := (colValue - colStart) / charWidth
			dashes := totalCols - len(label) - 5
			if dashes < 3 {
				dashes = 3
			}
			return strings.Repeat("-", dashes)
		},
		"rowWidth": func() int { return colValue - colStart },
		"valClass": func(color string) string {
			switch color {
			case "green":
				return "green"
			case "red":
				return "red"
			default:
				return "val"
			}
		},
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

func titleStr(username, hostname string) string {
	if username == "" && hostname == "" {
		return "user@host"
	}
	if hostname == "" {
		return username
	}
	return username + "@" + hostname
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

func parseFields(values []string) []Field {
	var fields []Field
	for _, f := range values {
		if f == "---space---" {
			fields = append(fields, Field{Type: "spacer"})
		} else if strings.HasPrefix(f, "---:") {
			fields = append(fields, makeDivider(strings.TrimPrefix(f, "---:")))
		} else {
			parts := strings.SplitN(f, ":", 2)
			if len(parts) == 2 {
				fields = append(fields, makeField(parts[0], parts[1], ""))
			}
		}
	}
	return fields
}

func appendGitHubStats(fields []Field, username string, showStats bool) []Field {
	if username == "" || !showStats {
		return fields
	}
	stats, err := FetchGitHubStats(username)
	if err != nil {
		return fields
	}
	if len(fields) > 0 {
		fields = append(fields, Field{Type: "halfspacer"})
	}
	return append(fields,
		makeDivider("GitHub Stats"),
		makeField("Repos",         fmt.Sprintf("%d", stats.TotalRepos),   ""),
		makeField("Commits",       fmt.Sprintf("%d", stats.TotalCommits), ""),
		makeField("Lines Added",   fmt.Sprintf("%d", stats.LinesAdded),   "green"),
		makeField("Lines Deleted", fmt.Sprintf("%d", stats.LinesDeleted), "red"),
	)
}

func makeField(label, value, color string) Field {
	return Field{
		Label: label,
		Dots:  makeDots(label, value),
		Value: value,
		Color: color,
		Type:  "field",
	}
}

func makeDivider(label string) Field {
	return Field{
		Label: label,
		Type:  "divider",
	}
}

func makeDots(label, value string) string {
	totalCols := (colValue - colStart) / charWidth
	labelCols := len(label) + 3
	valueCols := len(value)
	dots := totalCols - labelCols - valueCols
	if dots < 3 {
		dots = 3
	}
	return strings.Repeat(".", dots)
}

func rowHeight(f Field) int {
	if f.Type == "halfspacer" {
		return 9
	}
	return 17
}

func fieldsHeight(fields []Field) int {
	total := 0
	for _, f := range fields {
		total += rowHeight(f)
	}
	return total
}

func maxAsciiRows(fields []Field) int {
	cardHeight := fieldsTop + fieldsHeight(fields) + 6
	rows := (cardHeight - 14) / 7
	if rows < 1 {
		rows = 1
	}
	return rows
}

func trimAsciiLines(lines []string, fields []Field) []string {
	maxLines := maxAsciiRows(fields)
	if len(lines) <= maxLines {
		return lines
	}
	return lines[:maxLines]
}