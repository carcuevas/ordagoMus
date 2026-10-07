package tui

import (
	"io"
	"math"
	"sync"

	"ordagomus/arte"
	"ordagomus/internal/cards"
	"ordagomus/internal/kitty"
	"ordagomus/internal/personajes"
)

// tam es el tamaño de una carta ilustrada, en celdas.
type tam struct{ c, r int }

// graficos guarda si el terminal muestra las cartas ilustradas y a qué tamaños.
// Cada tamaño tiene su propia colocación en kitty, así que cambiarlo al
// redimensionar la ventana no obliga a volver a mandar las imágenes.
type graficos struct {
	on     bool
	big    []tam // tamaños posibles de la mano del jugador
	small  []tam // y de las manos de los demás
	tb, ts int   // tamaños en uso (los elige la mesa según el terminal)
	cw, ch int   // tamaño de la celda en píxeles
}

var gfx graficos

// Columnas de cada tamaño, de menor a mayor.
var (
	colsBig   = []int{11, 13, 15, 17, 19, 21}
	colsSmall = []int{7, 8, 9, 10, 11}
)

func (g graficos) cartaBig() tam {
	if len(g.big) == 0 {
		return tam{11, 8}
	}
	return g.big[g.tb]
}

func (g graficos) cartaSmall() tam {
	if len(g.small) == 0 {
		return tam{7, 5}
	}
	return g.small[g.ts]
}

// Las colocaciones caben en la tabla de diacríticos de kitty (< 30): en el
// navegador se mandan como tercer diacrítico de cada celda.
func pidBig(t int) uint32   { return pidTallas + uint32(t) }
func pidSmall(t int) uint32 { return pidTallas + 10 + uint32(t) }

const (
	idBase    = 0x4d0000
	idReverso = idBase + 99
	pidTallas = 5 // y siguientes: una colocación por tamaño
	pidCara   = 3
	pidFicha  = 4
	idCaras   = 0x4e0000
	idFichas  = 0x4f0000
	cartaAlto = 379.0 / 240.0 // proporción de las imágenes
)

func cardID(c cards.Card) uint32 { return idBase + uint32(c.Suit)*13 + uint32(c.Rank) }

// caraID es la imagen de una cara: personaje, variante (expresión o seña) y espejo.
func caraID(pj, v int, espejo bool) uint32 {
	id := idCaras + uint32(pj*64+v*2)
	if espejo {
		id++
	}
	return id
}

func filasPara(cols, cellW, cellH, lo, hi int) int {
	r := int(math.Round(float64(cols*cellW) * cartaAlto / float64(cellH)))
	return max(lo, min(hi, r))
}

// PrepararGraficos envía la baraja al terminal (debe llamarse ya dentro de la
// pantalla alternativa) y devuelve la función que la borra al salir.
func PrepararGraficos(w io.Writer) func() {
	cw, ch := kitty.CellSize()
	gfx = graficos{on: true, cw: cw, ch: ch}
	for _, c := range colsBig {
		gfx.big = append(gfx.big, tam{c, filasPara(c, cw, ch, 6, kitty.MaxCells)})
	}
	for _, c := range colsSmall {
		gfx.small = append(gfx.small, tam{c, filasPara(c, cw, ch, 4, kitty.MaxCells)})
	}

	var ids []uint32
	send := func(id uint32, png []byte) {
		kitty.Transmit(w, id, png)
		for i, t := range gfx.big {
			kitty.Place(w, id, pidBig(i), t.c, t.r)
		}
		for i, t := range gfx.small {
			kitty.Place(w, id, pidSmall(i), t.c, t.r)
		}
		ids = append(ids, id)
	}
	for _, c := range cards.NewDeck() {
		send(cardID(c), arte.Carta(c))
	}
	send(idReverso, arte.Reverso())

	// Las caras se pintan en paralelo: son muchas variantes.
	type imagen struct {
		id  uint32
		png []byte
	}
	var caras []imagen
	for i := range personajes.Todos {
		for v := 0; v < nVariante; v++ {
			for _, e := range []bool{false, true} {
				caras = append(caras, imagen{id: caraID(i, v, e)})
			}
		}
	}
	var wg sync.WaitGroup
	for k := range caras {
		wg.Add(1)
		go func() {
			defer wg.Done()
			pj := personajes.Todos[(caras[k].id-idCaras)/64]
			ex, sena := desdeVariante(int((caras[k].id - idCaras) % 64 / 2))
			caras[k].png = pngCara(pj.ID, ex, sena, caras[k].id%2 == 1)
		}()
	}
	wg.Wait()
	for _, c := range caras {
		kitty.Transmit(w, c.id, c.png)
		kitty.Place(w, c.id, pidCara, caraW, caraFilas)
		ids = append(ids, c.id)
	}
	for _, f := range fichasImg() {
		kitty.Transmit(w, f.id, f.png)
		kitty.Place(w, f.id, pidFicha, f.cols, f.rows)
		ids = append(ids, f.id)
	}
	return func() { kitty.DeleteAll(w, ids) }
}
