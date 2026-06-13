package handlers

import (
	"image"
	"net/http"
	"strings"
)

func CardGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.ParseMultipartForm(10 << 20)

	username := r.FormValue("username")
	showStats := r.FormValue("showstats") != "false"

	rawFields := r.Form["field"]
	fields := parseFields(rawFields)
	fields = appendGitHubStats(fields, username, showStats)

	input := defaultInput(CardInput{
		Username:   username,
		Hostname:   r.FormValue("hostname"),
		Background: r.FormValue("background"),
		KeyColor:   r.FormValue("keycolor"),
		TextColor:  r.FormValue("textcolor"),
		Fields:     fields,
	})

	var img image.Image
	var asciiArt string
	var result string

	ok := run(w,
		func() error {
			err := parseImage(r, &img)
			if err != nil {
				if a := r.FormValue("ascii"); a != "" {
					asciiArt = a
				}
				return nil
			}
			return convertToAscii(img, &asciiArt, maxAsciiRows(input.Fields))
		},
		func() error {
			if asciiArt != "" {
				lines := strings.Split(asciiArt, "\n")
				input.AsciiLines = trimAsciiLines(lines, input.Fields)
			}
			input.ConfigJSON = encodeConfig(cardConfig{
				Username:   username,
				Hostname:   input.Hostname,
				Background: input.Background,
				KeyColor:   input.KeyColor,
				TextColor:  input.TextColor,
				ShowStats:  showStats,
				Fields:     rawFields,
				Ascii:      strings.Join(input.AsciiLines, "\n"),
			})
			return renderCard(input, &result)
		},
	)
	if !ok {
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Write([]byte(result))
}