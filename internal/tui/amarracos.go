package tui

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"ordagomus/internal/kitty"
)

// Los tantos se llevan como en la mesa de verdad: cada amarraco vale cinco
// piedras. El amarraco es una ficha gorda con cinco puntos, como el cinco del
// dado; la piedra, un garbanzo.

const (
	idAmarraco = idFichas + 1
	idPiedra   = idFichas + 2
	amarracoC  = 4 // celdas que ocupa un amarraco (×2 filas)
	piedraC    = 2 // celdas de una piedra (×1 fila)
)

type fichaImg struct {
	id         uint32
	cols, rows int
	png        []byte
}

// fichasImg dibuja el amarraco y la piedra a la medida de la celda para que
// salgan redondos.
func fichasImg() []fichaImg {
	cw, ch := max(gfx.cw, 4), max(gfx.ch, 8)
	return []fichaImg{
		{idAmarraco, amarracoC, 2, codifica(dibujaAmarraco(amarracoC*cw, 2*ch))},
		{idPiedra, piedraC, 1, codifica(dibujaPiedra(piedraC*cw, ch))},
	}
}

func codifica(img image.Image) []byte {
	var buf bytes.Buffer
	png.Encode(&buf, img)
	return buf.Bytes()
}

func dibujaAmarraco(w, h int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	l := lienzo{img}
	cx, cy := float64(w)/2, float64(h)/2
	r := min(float64(w), float64(h))/2 - 1
	l.elipse(cx+r*0.06, cy+r*0.08, r, r, hex("#3a240a"))           // sombra
	l.elipse(cx, cy, r, r, hex("#8a5a14"))                         // canto
	l.elipse(cx, cy, r*0.86, r*0.86, hex("#d9a136"))               // cara
	l.elipse(cx, cy, r*0.66, r*0.66, hex("#b98222"))               // aro
	l.elipse(cx, cy, r*0.6, r*0.6, hex("#e6b54c"))                 // centro
	l.elipse(cx-r*0.35, cy-r*0.45, r*0.22, r*0.12, hex("#f6dc94")) // brillo
	p := r * 0.1
	for _, d := range [][2]float64{{0, 0}, {-1, -1}, {1, -1}, {-1, 1}, {1, 1}} {
		l.elipse(cx+d[0]*r*0.3, cy+d[1]*r*0.3, p, p, hex("#5a3608"))
	}
	return img
}

func dibujaPiedra(w, h int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	l := lienzo{img}
	cx, cy := float64(w)/2, float64(h)/2
	r := min(float64(w), float64(h))/2 - 1
	l.elipse(cx+r*0.08, cy+r*0.1, r*0.95, r*0.8, hex("#2a2a26"))
	l.elipse(cx, cy, r*0.95, r*0.8, hex("#b8b2a2"))
	l.elipse(cx-r*0.1, cy-r*0.1, r*0.7, r*0.55, hex("#d8d2c2"))
	l.elipse(cx-r*0.35, cy-r*0.3, r*0.2, r*0.14, hex("#f4f0e6"))
	return img
}

var (
	styleAmarraco = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	stylePiedra   = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
)

// amarraco y piedras devuelven siempre dos líneas para poder juntarlos.
func amarraco(img bool) string {
	if img {
		return strings.Join(kitty.Lines(idAmarraco, pidFicha, amarracoC, 2), "\n")
	}
	return styleAmarraco.Render("▄██▄") + "\n" + styleAmarraco.Render("▀██▀")
}

func piedras(n int, img bool) string {
	fila := func(k int) string {
		if img {
			var b strings.Builder
			for i := 0; i < k; i++ {
				b.WriteString(kitty.Lines(idPiedra, pidFicha, piedraC, 1)[0])
			}
			return b.String() + strings.Repeat(" ", (2-k)*piedraC)
		}
		return stylePiedra.Render(strings.TrimSpace(strings.Repeat("● ", k))) + strings.Repeat(" ", 4-max(2*k-1, 0))
	}
	return fila(min(n, 2)) + "\n" + fila(max(n-2, 0))
}

// fichas pinta n tantos en amarracos y piedras en dos líneas. Si no caben en
// ancho, se dibuja un amarraco y se pone cuántos son.
func fichas(n int, img bool, ancho int) string {
	if n <= 0 {
		return styleDim.Render("    \n    ")
	}
	a, p := n/5, n%5
	var cols []string
	if a > 0 && a*(amarracoC+1)+piedraC*2 > ancho {
		cols = append(cols, amarraco(img), styleAmarraco.Render(fmt.Sprintf(" ×%d ", a))+"\n")
	} else {
		for i := 0; i < a; i++ {
			cols = append(cols, amarraco(img), " \n ")
		}
	}
	if p > 0 {
		cols = append(cols, piedras(p, img))
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, cols...)
}
