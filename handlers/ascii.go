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
		func() error { return convertToAscii(img, &result, 0) },
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

const maxAsciiCols = 60

func convertToAscii(img image.Image, result *string, maxRows int) error {
	bounds, hasAlpha := contentBounds(img)
	imgWidth := bounds.Dx()
	imgHeight := bounds.Dy()
	if imgWidth == 0 || imgHeight == 0 {
		*result = ""
		return nil
	}

	ink := func(c color.Color) float64 {
		if hasAlpha {
			_, _, _, a := c.RGBA()
			return float64(a >> 8)
		}
		return float64(toGrayScale(c))
	}

	targetWidth := maxAsciiCols
	if maxRows > 0 {
		aspect := float64(imgWidth) / float64(imgHeight)
		targetWidth = int(math.Round(float64(maxRows) * aspect * 2))
		if targetWidth > maxAsciiCols {
			targetWidth = maxAsciiCols
		}
		if targetWidth < 1 {
			targetWidth = 1
		}
	}

	ratio := float64(imgHeight) / float64(imgWidth)
	targetHeight := int(math.Round(float64(targetWidth) * ratio / 2))
	if targetHeight < 1 {
		targetHeight = 1
	}
	if maxRows > 0 && targetHeight > maxRows {
		targetHeight = maxRows
	}

	var sb strings.Builder

	for y := 0; y < targetHeight; y++ {
		y0 := bounds.Min.Y + y*imgHeight/targetHeight
		y1 := bounds.Min.Y + (y+1)*imgHeight/targetHeight
		if y1 <= y0 {
			y1 = y0 + 1
		}
		for x := 0; x < targetWidth; x++ {
			x0 := bounds.Min.X + x*imgWidth/targetWidth
			x1 := bounds.Min.X + (x+1)*imgWidth/targetWidth
			if x1 <= x0 {
				x1 = x0 + 1
			}

			var sum, count float64
			for sy := y0; sy < y1; sy++ {
				for sx := x0; sx < x1; sx++ {
					sum += ink(img.At(sx, sy))
					count++
				}
			}
			value := sum / count
			index := int((1 - value/255) * float64(len(baseAsciiChars)-1))

			sb.WriteString(baseAsciiChars[index])
		}

		sb.WriteString("\n")
	}

	*result = sb.String()
	return nil
}

func contentBounds(img image.Image) (image.Rectangle, bool) {
	b := img.Bounds()
	minX, minY := b.Max.X, b.Max.Y
	maxX, maxY := b.Min.X, b.Min.Y
	found := false
	hasAlpha := false

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a>>8 < 16 {
				hasAlpha = true
				continue
			}
			found = true
			if x < minX {
				minX = x
			}
			if x > maxX {
				maxX = x
			}
			if y < minY {
				minY = y
			}
			if y > maxY {
				maxY = y
			}
		}
	}

	if !found {
		return b, hasAlpha
	}
	return image.Rect(minX, minY, maxX+1, maxY+1), hasAlpha
}