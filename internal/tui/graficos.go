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

// graficos guarda si el terminal muestra las cartas ilustradas y a qué tamaño.
type graficos struct {
	on           bool
	bigC, bigR   int
	smallC, smal int
	cw, ch       int // tamaño de la celda en píxeles
}

var gfx graficos

const (
	idBase    = 0x4d0000
	idReverso = idBase + 99
	pidBig    = 1
	pidSmall  = 2
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
	gfx = graficos{on: true, bigC: 11, smallC: 7, cw: cw, ch: ch}
	gfx.bigR = filasPara(gfx.bigC, cw, ch, 6, 14)
	gfx.smal = filasPara(gfx.smallC, cw, ch, 4, 10)

	var ids []uint32
	send := func(id uint32, png []byte) {
		kitty.Transmit(w, id, png)
		kitty.Place(w, id, pidBig, gfx.bigC, gfx.bigR)
		kitty.Place(w, id, pidSmall, gfx.smallC, gfx.smal)
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
