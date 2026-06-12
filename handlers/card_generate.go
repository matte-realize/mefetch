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

	var fields []Field
	for _, f := range r.Form["field"] {
		parts := strings.SplitN(f, ":", 2)
		if len(parts) == 2 {
			fields = append(fields, makeField(parts[0], parts[1], ""))
		}
	}

	input := defaultInput(CardInput{
		Username:   r.FormValue("username"),
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