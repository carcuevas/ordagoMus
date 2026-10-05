package tui

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"math"

	"ordagomus/internal/rules"
)

// Las caras en alta resolución para los terminales con imágenes (kitty): el
// retrato pixelado se suaviza con Scale2x tres veces (×8) y encima se dibujan,
// grandes y nítidos, los ojos, las cejas y la boca, que es donde se hacen las
// señas. Ocupan lo mismo que la cara de medios bloques (caraW×caraFilas celdas).

const (
	escala = 8
	imgW   = caraW * escala
	imgH   = caraH * escala
)

// variantes de cada cara: las expresiones y luego una por seña.
const (
	nExpr     = int(exTriste) + 1
	nVariante = nExpr + int(rules.SenaCiego) + 1
)

func variante(ex expresion, sena *rules.Sena) int {
	if sena != nil {
		return nExpr + int(*sena)
	}
	return int(ex)
}

func desdeVariante(v int) (expresion, *rules.Sena) {
	if v < nExpr {
		return expresion(v), nil
	}
	s := rules.Sena(v - nExpr)
	return exNormal, &s
}

// gafas: quien las lleva puestas y de qué color es la montura.
var gafas = map[string]string{"remedios": "#b8902c", "montse": "#262626"}

func hex(s string) color.RGBA {
	var c [3]uint8
	for i := range c {
		h := s[1+2*i : 3+2*i]
		var v uint8
		for _, r := range h {
			v <<= 4
			switch {
			case r >= '0' && r <= '9':
				v |= uint8(r - '0')
			case r >= 'a' && r <= 'f':
				v |= uint8(r-'a') + 10
			case r >= 'A' && r <= 'F':
				v |= uint8(r-'A') + 10
			}
		}
		c[i] = v
	}
	return color.RGBA{c[0], c[1], c[2], 255}
}

func oscurece(c color.RGBA, f float64) color.RGBA {
	return color.RGBA{uint8(float64(c.R) * f), uint8(float64(c.G) * f), uint8(float64(c.B) * f), 255}
}

// scale2x (EPX) dobla la rejilla redondeando las diagonales.
func scale2x(src [][]byte) [][]byte {
	h, w := len(src), len(src[0])
	at := func(x, y int) byte {
		if x < 0 || y < 0 || x >= w || y >= h {
			return '.'
		}
		return src[y][x]
	}
	out := make([][]byte, 2*h)
	for y := range out {
		out[y] = make([]byte, 2*w)
	}
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			p := src[y][x]
			a, b, c, d := at(x, y-1), at(x+1, y), at(x-1, y), at(x, y+1)
			e0, e1, e2, e3 := p, p, p, p
			if c == a && c != d && a != b {
				e0 = a
			}
			if a == b && a != c && b != d {
				e1 = b
			}
			if d == c && d != b && c != a {
				e2 = c
			}
			if b == d && b != a && d != c {
				e3 = d
			}
			out[2*y][2*x], out[2*y][2*x+1] = e0, e1
			out[2*y+1][2*x], out[2*y+1][2*x+1] = e2, e3
		}
	}
	return out
}

// lienzo dibuja formas sencillas sobre la imagen.
type lienzo struct{ *image.RGBA }

func (l lienzo) elipse(cx, cy, rx, ry float64, c color.RGBA) {
	for y := int(cy - ry - 1); y <= int(cy+ry+1); y++ {
		for x := int(cx - rx - 1); x <= int(cx+rx+1); x++ {
			dx, dy := (float64(x)+0.5-cx)/rx, (float64(y)+0.5-cy)/ry
			if dx*dx+dy*dy <= 1 {
				l.SetRGBA(x, y, c)
			}
		}
	}
}

// relleno pinta los píxeles cuyo centro cumple dentro.
func (l lienzo) relleno(dentro func(x, y float64) bool, c color.RGBA) {
	b := l.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			if dentro(float64(x)+0.5, float64(y)+0.5) {
				l.SetRGBA(x, y, c)
			}
		}
	}
}

func (l lienzo) rect(x0, y0, x1, y1 float64, c color.RGBA) {
	for y := int(y0); y < int(math.Ceil(y1)); y++ {
		for x := int(x0); x < int(math.Ceil(x1)); x++ {
			l.SetRGBA(x, y, c)
		}
	}
}

// curva es un trazo de grosor r por una curva cuadrática (x0,y0)→(x1,y1) con
// punto de control (qx,qy).
func (l lienzo) curva(x0, y0, qx, qy, x1, y1, r float64, c color.RGBA) {
	const pasos = 48
	for i := 0; i <= pasos; i++ {
		t := float64(i) / pasos
		x := (1-t)*(1-t)*x0 + 2*(1-t)*t*qx + t*t*x1
		y := (1-t)*(1-t)*y0 + 2*(1-t)*t*qy + t*t*y1
		l.elipse(x, y, r, r, c)
	}
}

func (l lienzo) linea(x0, y0, x1, y1, r float64, c color.RGBA) {
	l.curva(x0, y0, (x0+x1)/2, (y0+y1)/2, x1, y1, r, c)
}

// Sitio de las facciones en la imagen grande (ver la rejilla en caras.go).
const (
	ojoIzq = 5.0 * escala // centro x del ojo izquierdo
	ojoDer = 9.0 * escala
	ojoY   = 7.5 * escala
	cejaY  = 6.2 * escala
	bocaX  = 7.0 * escala
	bocaY  = 11.5 * escala
)

// imagenCara dibuja la cara en alta resolución.
func (c cara) imagenCara(id string, ex expresion, sena *rules.Sena, espejo bool) *image.RGBA {
	rej := make([][]byte, caraH)
	for y, row := range c.px {
		rej[y] = []byte(row)
		if _, ok := gafas[id]; ok && (y == 7 || y == 8) {
			for x := range rej[y] {
				if rej[y][x] == 'g' {
					rej[y][x] = 's'
				}
			}
		}
	}
	for i := 0; i < 3; i++ {
		rej = scale2x(rej)
	}
	img := image.NewRGBA(image.Rect(0, 0, imgW, imgH))
	for y, row := range rej {
		for x, b := range row {
			if col, ok := c.color(b); ok {
				img.SetRGBA(x, y, hex(col))
			}
		}
	}
	l := lienzo{img}
	col := func(b byte) color.RGBA {
		s, _ := c.color(b)
		return hex(s)
	}
	sombra := col('S')
	labio := col('m')
	cejaC := oscurece(col('b'), 0.8)
	tinta := hex("#151515")
	blanco := hex("#fbfaf6")
	dientes := hex("#fffdf4")
	boca := hex("#3a1216")
	lengua := hex("#ee6f8c")
	lenguaD := hex("#b8405e")

	// ── ojos ──
	abierto := func(cx, ry float64) {
		l.elipse(cx, ojoY, 8.5, ry, blanco)
		l.elipse(cx+2, ojoY+0.5, 4.2, min(4.4, ry), tinta)
		l.elipse(cx+3.3, ojoY-1.2, 1.3, 1.3, blanco) // brillo
		l.curva(cx-8.5, ojoY, cx, ojoY-ry*2.1, cx+8.5, ojoY, 1.4, tinta)
	}
	cerrado := func(cx float64) { // ‿ párpado cerrado, con pestañas
		l.curva(cx-8.5, ojoY-1, cx, ojoY+5, cx+8.5, ojoY-1, 1.9, tinta)
		for _, dx := range []float64{-4.5, 0, 4.5} {
			l.linea(cx+dx, ojoY+2.5, cx+dx*1.3, ojoY+6.5, 0.9, tinta)
		}
	}
	guino := func(cx float64) { // ^ ojo apretado y moflete arriba
		l.curva(cx-9, ojoY+1, cx, ojoY-6, cx+9, ojoY+1, 2.2, tinta)
		l.curva(cx-8, ojoY+8, cx, ojoY+3, cx+8, ojoY+8, 1.6, sombra)
	}
	ceja := func(cx, yIn, yOut, arco float64, dentroDerecha bool) {
		xIn, xOut := cx+8, cx-8
		if !dentroDerecha {
			xIn, xOut = cx-8, cx+8
		}
		l.curva(xOut, yOut, cx, (yIn+yOut)/2-arco, xIn, yIn, 2.6, cejaC)
	}
	cejas := func(dy, arco, inclina float64) {
		// inclina > 0: la parte de dentro baja (enfado); < 0: sube (pena).
		ceja(ojoIzq, cejaY+dy+inclina, cejaY+dy-inclina, arco, true)
		ceja(ojoDer, cejaY+dy+inclina, cejaY+dy-inclina, arco, false)
	}

	ojosRy := 5.0
	switch {
	case sena != nil && *sena == rules.SenaCiego:
		cejas(0, 2, 0)
		cerrado(ojoIzq)
		cerrado(ojoDer)
	case sena != nil && *sena == rules.Sena31:
		ceja(ojoIzq, cejaY-2, cejaY-2, 3, true)
		ceja(ojoDer, cejaY+3, cejaY+2, 1, false)
		abierto(ojoIzq, ojosRy)
		guino(ojoDer)
	case sena != nil && *sena == rules.SenaDuples:
		cejas(-12, 7, 0)
		for _, dy := range []float64{-24, -19} { // arrugas de la frente
			l.curva(ojoIzq-3, cejaY+dy+1, bocaX, cejaY+dy-2, ojoDer+3, cejaY+dy+1, 0.9, sombra)
		}
		abierto(ojoIzq, 7.5)
		abierto(ojoDer, 7.5)
	case sena == nil && ex == exEnfadado:
		cejas(1, 0, 4)
		abierto(ojoIzq, 4)
		abierto(ojoDer, 4)
	case sena == nil && ex == exTriste:
		cejas(-1, 0, -4)
		abierto(ojoIzq, 4.5)
		abierto(ojoDer, 4.5)
	case sena == nil && ex == exContento:
		cejas(-2, 3, 0)
		abierto(ojoIzq, 4.5)
		abierto(ojoDer, 4.5)
	default:
		cejas(0, 2, 0)
		abierto(ojoIzq, ojosRy)
		abierto(ojoDer, ojosRy)
	}
	if g, ok := gafas[id]; ok {
		m := hex(g)
		for _, cx := range []float64{ojoIzq, ojoDer} {
			l.curva(cx-10, ojoY-6, cx, ojoY-8, cx+10, ojoY-6, 1.3, m)
			l.curva(cx-10, ojoY+6, cx, ojoY+9, cx+10, ojoY+6, 1.3, m)
			l.linea(cx-10, ojoY-6, cx-10, ojoY+6, 1.3, m)
			l.linea(cx+10, ojoY-6, cx+10, ojoY+6, 1.3, m)
		}
		l.curva(ojoIzq+10, ojoY-2, bocaX, ojoY-5, ojoDer-10, ojoY-2, 1.2, m)
	}

	// ── boca ──
	cerrada := func(dy float64) {
		l.curva(bocaX-12, bocaY, bocaX, bocaY+dy, bocaX+12, bocaY, 1.9, labio)
	}
	switch {
	case sena != nil && *sena == rules.SenaDosReyes: // se muerde el labio de abajo
		l.curva(bocaX-13, bocaY-3, bocaX, bocaY-5, bocaX+13, bocaY-3, 1.6, labio)
		l.elipse(bocaX, bocaY+3.5, 12, 5, oscurece(labio, 1.15))
		l.curva(bocaX-12, bocaY+5, bocaX, bocaY+11, bocaX+12, bocaY+5, 1.2, oscurece(labio, 0.7))
		l.relleno(func(x, y float64) bool { // los dientes de arriba clavados en el labio
			dx := x - bocaX
			return math.Abs(dx) <= 8 && y >= bocaY-3 && y <= bocaY+3.5-(dx/8)*(dx/8)*2
		}, dientes)
		for _, dx := range []float64{-4, 0, 4} {
			l.rect(bocaX+dx-0.45, bocaY-3, bocaX+dx+0.45, bocaY+2.5, hex("#cfc6b4"))
		}
	case sena != nil && *sena == rules.SenaDosAses: // la punta de la lengua
		l.elipse(bocaX, bocaY, 8, 3.6, labio)
		l.elipse(bocaX, bocaY, 6, 2, boca)
		l.elipse(bocaX, bocaY+6, 6.5, 7.5, lenguaD)
		l.elipse(bocaX, bocaY+5.5, 5.5, 6.5, lengua)
		l.linea(bocaX, bocaY+2, bocaX, bocaY+8, 0.8, lenguaD)
	case sena != nil && *sena == rules.SenaMediasReyes: // tuerce la comisura
		l.curva(bocaX-12, bocaY+2, bocaX-2, bocaY+3, bocaX+4, bocaY+1, 2.3, labio)
		l.curva(bocaX+4, bocaY+1, bocaX+11, bocaY-1, bocaX+15, bocaY-10, 2.3, labio)
		l.curva(bocaX+12, bocaY-15, bocaX+21, bocaY-12, bocaX+20, bocaY-3, 1.4, sombra) // moflete
		l.curva(bocaX-12, bocaY+2, bocaX-13, bocaY+4, bocaX-12, bocaY+5, 1.2, oscurece(labio, 0.7))
	case sena != nil && *sena == rules.SenaMediasAses: // la lengua a un lado
		cerrada(1.5)
		l.elipse(bocaX+14, bocaY+5, 6.5, 7.5, lenguaD)
		l.elipse(bocaX+14, bocaY+4.5, 5.5, 6.5, lengua)
		l.linea(bocaX+13, bocaY+1, bocaX+15, bocaY+8, 0.8, lenguaD)
		l.elipse(bocaX+11, bocaY, 3.5, 2, labio)
	case sena != nil && *sena == rules.Sena31:
		l.curva(bocaX-11, bocaY+1, bocaX, bocaY+2, bocaX+12, bocaY-3, 1.9, labio) // media sonrisa
	case sena != nil:
		cerrada(1)
	case ex == exHabla:
		l.elipse(bocaX, bocaY+1, 8.5, 6, labio)
		l.elipse(bocaX, bocaY+1, 6.5, 4.2, boca)
		l.rect(bocaX-5, bocaY-3, bocaX+5, bocaY-1.2, dientes)
	case ex == exContento: // sonrisa abierta: una D tumbada con los dientes de arriba
		arriba := func(dx float64) float64 { return bocaY - 3 + 2*(dx/13)*(dx/13) }
		abajo := func(dx float64) float64 { return bocaY - 3 + 10*(1-(dx/13)*(dx/13)) }
		l.relleno(func(x, y float64) bool {
			dx := x - bocaX
			return math.Abs(dx) <= 13 && y >= arriba(dx) && y <= abajo(dx)
		}, boca)
		l.relleno(func(x, y float64) bool {
			dx := x - bocaX
			return math.Abs(dx) <= 11 && y >= arriba(dx) && y <= arriba(dx)+3 && y <= abajo(dx)
		}, dientes)
		l.curva(bocaX-13, bocaY-1, bocaX, bocaY+15, bocaX+13, bocaY-1, 1.5, labio)
		l.curva(bocaX-13, bocaY-1, bocaX, bocaY-3, bocaX+13, bocaY-1, 1.3, labio)
	case ex == exEnfadado:
		l.curva(bocaX-11, bocaY+3, bocaX, bocaY-4, bocaX+11, bocaY+3, 2.1, labio)
		l.rect(bocaX-6, bocaY-1.2, bocaX+6, bocaY+0.8, dientes)
	case ex == exTriste:
		l.curva(bocaX-10, bocaY+3, bocaX, bocaY-3, bocaX+10, bocaY+3, 1.8, labio)
	default:
		cerrada(1.5)
	}
	if espejo {
		b := img.Bounds()
		for y := 0; y < b.Dy(); y++ {
			for x := 0; x < b.Dx()/2; x++ {
				a, z := img.RGBAAt(x, y), img.RGBAAt(b.Dx()-1-x, y)
				img.SetRGBA(x, y, z)
				img.SetRGBA(b.Dx()-1-x, y, a)
			}
		}
	}
	return img
}

// pngCara es la imagen de la cara codificada en PNG.
func pngCara(id string, ex expresion, sena *rules.Sena, espejo bool) []byte {
	c, ok := caras[id]
	if !ok {
		return nil
	}
	var buf bytes.Buffer
	enc := png.Encoder{CompressionLevel: png.BestSpeed}
	enc.Encode(&buf, c.imagenCara(id, ex, sena, espejo))
	return buf.Bytes()
}
