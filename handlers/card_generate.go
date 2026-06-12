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

	var img image.Image
	var asciiArt string
	var result string

	fields := []Field{
		{Label: "OS",     Dots: makeDots("OS"),     Value: r.FormValue("os")},
		{Label: "Kernel", Dots: makeDots("Kernel"), Value: r.FormValue("kernel")},
		{Label: "Shell",  Dots: makeDots("Shell"),  Value: r.FormValue("shell")},
		{Label: "Editor", Dots: makeDots("Editor"), Value: r.FormValue("editor")},
		{Label: "Uptime", Dots: makeDots("Uptime"), Value: r.FormValue("uptime")},
	}

	var filtered []Field
	for _, f := range fields {
		if f.Value != "" {
			filtered = append(filtered, f)
		}
	}

	input := defaultInput(CardInput{
		Username:   r.FormValue("username"),
		Hostname:   r.FormValue("hostname"),
		Background: r.FormValue("background"),
		KeyColor:   r.FormValue("keycolor"),
		TextColor:  r.FormValue("textcolor"),
		Fields:     filtered,
	})

	ok := run(w,
		func() error {
			err := parseImage(r, &img)
			if err != nil {
				return nil
			}
			return convertToAscii(img, &asciiArt)
		},
		func() error {
			if asciiArt != "" {
				input.AsciiLines = strings.Split(asciiArt, "\n")
			}
			return renderCard(input, &result)
		},
	)
	if !ok {
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Write([]byte(result))
}