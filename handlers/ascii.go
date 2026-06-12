package handlers

import (
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"net/http"
	"strings"
)

var baseAsciiChars = []string {
	"@", "#", "S", "%", "?", "*", "+", ";", ":", ",", ".", " ",
}

func Ascii(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var img image.Image
	var result string

	ok := run(w,
		func() error { return parseImage(r, &img) },
		func() error { return convertToAscii(img, &result) },
	)

	if !ok {
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(result))
}

func parseImage(r *http.Request, img *image.Image) error {
	r.ParseMultipartForm(10 << 20)

	file, _, err := r.FormFile("image")

	if err != nil {
		return err
	}

	defer file.Close()

	decoded, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	*img = decoded
	return nil
}

func toGrayScale(c color.Color) uint8 {
	r, g, b, _  := c.RGBA()
	gray := 0.299 * float64(r) + 0.587 * float64(g) + 0.114 * float64(b)
	return uint8(gray / 256)
}

func convertToAscii(img image.Image, result *string) error {
	targetWidth := 60
	bounds := img.Bounds()
	imgWidth := bounds.Max.X
	imgHeight := bounds.Max.Y

	ratio := float64(imgHeight) / float64(imgWidth)
	targetHeight := int(math.Round(float64(targetWidth) * ratio / 2))

	var sb strings.Builder

	for y := 0; y < targetHeight; y++ {
		for x := 0; x < targetWidth; x++ {
			imgX := int(float64(x) / float64(targetWidth) * float64(imgWidth))
			imgY := int(float64(y) / float64(targetHeight) * float64(imgHeight))

			pixel := img.At(imgX, imgY)
			gray := toGrayScale(pixel)

			index := int(float64(gray) / 255 * float64(len(baseAsciiChars)-1))
			sb.WriteString(baseAsciiChars[index])
		}

		sb.WriteString("\n")
	}

	*result = sb.String()
	return nil
}