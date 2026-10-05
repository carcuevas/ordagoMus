package tui

import (
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"ordagomus/internal/personajes"
)

// ORDAGO_HOJA=dir go test -run Hoja ./internal/tui guarda una hoja con todas
// las caras para verlas de un vistazo.
func TestHojaDeCaras(t *testing.T) {
	dir := os.Getenv("ORDAGO_HOJA")
	if dir == "" {
		t.Skip("ORDAGO_HOJA no definido")
	}
	const sep = 6
	hoja := image.NewRGBA(image.Rect(0, 0, nVariante*(imgW+sep), len(personajes.Todos)*(imgH+sep)))
	draw.Draw(hoja, hoja.Bounds(), image.NewUniform(hex("#1e2a22")), image.Point{}, draw.Src)
	for f, p := range personajes.Todos {
		for v := 0; v < nVariante; v++ {
			ex, s := desdeVariante(v)
			img := caras[p.ID].imagenCara(p.ID, ex, s, false)
			r := image.Rect(v*(imgW+sep), f*(imgH+sep), v*(imgW+sep)+imgW, f*(imgH+sep)+imgH)
			draw.Draw(hoja, r, img, image.Point{}, draw.Over)
		}
	}
	out, err := os.Create(filepath.Join(dir, "caras.png"))
	if err != nil {
		t.Fatal(err)
	}
	defer out.Close()
	png.Encode(out, hoja)
}
