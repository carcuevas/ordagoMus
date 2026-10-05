// Package rules contiene la valoración de jugadas del mus (reglamento de Bizkaia):
// grande, chica, pares, juego y punto, además de las señas.
package rules

import (
	"fmt"
	"slices"
	"strings"

	"ordagomus/internal/cards"
)

// Config son las opciones de la partida.
type Config struct {
	OchoReyes     bool // los treses valen como reyes y los doses como ases
	Tantos        int  // tantos para ganar un juego
	JuegosPorVaca int  // una vaca se juega al mejor de N juegos
	Vacas         int  // la partida se juega al mejor de N vacas
}

func DefaultConfig() Config {
	return Config{OchoReyes: true, Tantos: 40, JuegosPorVaca: 3, Vacas: 1}
}

type Hand [4]cards.Card

type Lance int

const (
	Grande Lance = iota
	Chica
	Pares
	Juego
	Punto
)

var lanceNames = [...]string{"Grande", "Chica", "Pares", "Juego", "Punto"}

func (l Lance) String() string { return lanceNames[l] }

// Rank es el valor efectivo de la carta: con 8 reyes el 3 es rey y el 2 es as.
func (c Config) Rank(card cards.Card) int {
	if c.OchoReyes {
		switch card.Rank {
		case 3:
			return 12
		case 2:
			return 1
		}
	}
	return card.Rank
}

// Value es lo que vale la carta para el juego: figuras 10, el resto su número.
func (c Config) Value(card cards.Card) int {
	r := c.Rank(card)
	if r >= 10 {
		return 10
	}
	return r
}

func (c Config) ranks(h Hand) []int {
	rs := make([]int, 4)
	for i, card := range h {
		rs[i] = c.Rank(card)
	}
	return rs
}

// Points es la suma de la mano para juego o punto.
func (c Config) Points(h Hand) int {
	sum := 0
	for _, card := range h {
		sum += c.Value(card)
	}
	return sum
}

func (c Config) HasJuego(h Hand) bool { return c.Points(h) >= 31 }

// JuegoValue: la 31 vale 3 tantos, cualquier otro juego 2.
func JuegoValue(points int) int {
	if points == 31 {
		return 3
	}
	return 2
}

// juegoOrder: 31, 32, 40, 37, 36, 35, 34, 33.
func juegoOrder(points int) int {
	switch points {
	case 31:
		return 8
	case 32:
		return 7
	case 40:
		return 6
	}
	return points - 32 // 37→5 ... 33→1
}

type ParesKind int

const (
	SinPares ParesKind = iota
	Pareja
	Medias
	Duples
)

type ParesInfo struct {
	Kind      ParesKind
	High, Low int // rango del par (en duples, el mayor y el menor)
}

func (p ParesInfo) Value() int { return int(p.Kind) }

// Pares calcula la jugada de pares de la mano.
func (c Config) Pares(h Hand) ParesInfo {
	count := map[int]int{}
	for _, r := range c.ranks(h) {
		count[r]++
	}
	var pairs []int
	for r, n := range count {
		switch n {
		case 4:
			return ParesInfo{Kind: Duples, High: r, Low: r}
		case 3:
			return ParesInfo{Kind: Medias, High: r}
		case 2:
			pairs = append(pairs, r)
		}
	}
	switch len(pairs) {
	case 2:
		slices.Sort(pairs)
		return ParesInfo{Kind: Duples, High: pairs[1], Low: pairs[0]}
	case 1:
		return ParesInfo{Kind: Pareja, High: pairs[0]}
	}
	return ParesInfo{}
}

func (c Config) HasPares(h Hand) bool { return c.Pares(h).Kind != SinPares }

func compareInts(a, b int) int {
	switch {
	case a > b:
		return 1
	case a < b:
		return -1
	}
	return 0
}

// Compare devuelve >0 si a gana a b en el lance, <0 si pierde y 0 si empatan
// (el empate lo gana quien esté más cerca de la mano).
func (c Config) Compare(l Lance, a, b Hand) int {
	switch l {
	case Grande:
		ra, rb := c.ranks(a), c.ranks(b)
		slices.Sort(ra)
		slices.Sort(rb)
		for i := 3; i >= 0; i-- {
			if d := compareInts(ra[i], rb[i]); d != 0 {
				return d
			}
		}
	case Chica:
		ra, rb := c.ranks(a), c.ranks(b)
		slices.Sort(ra)
		slices.Sort(rb)
		for i := 0; i < 4; i++ {
			if d := compareInts(rb[i], ra[i]); d != 0 {
				return d
			}
		}
	case Pares:
		pa, pb := c.Pares(a), c.Pares(b)
		if d := compareInts(int(pa.Kind), int(pb.Kind)); d != 0 {
			return d
		}
		if d := compareInts(pa.High, pb.High); d != 0 {
			return d
		}
		return compareInts(pa.Low, pb.Low)
	case Juego:
		return compareInts(juegoOrder(c.Points(a)), juegoOrder(c.Points(b)))
	case Punto:
		return compareInts(c.Points(a), c.Points(b))
	}
	return 0
}

// Eligible indica si la mano puede jugar el lance (pares y juego exigen tenerlos).
func (c Config) Eligible(l Lance, h Hand) bool {
	switch l {
	case Pares:
		return c.HasPares(h)
	case Juego:
		return c.HasJuego(h)
	}
	return true
}

// Winner devuelve el asiento ganador del lance, recorriendo desde la mano para
// que los empates los gane el más cercano a ella. -1 si nadie puede jugarlo.
func (c Config) Winner(l Lance, hands [4]Hand, mano int) int {
	best := -1
	for i := 0; i < 4; i++ {
		s := (mano + i) % 4
		if !c.Eligible(l, hands[s]) {
			continue
		}
		if best < 0 || c.Compare(l, hands[s], hands[best]) > 0 {
			best = s
		}
	}
	return best
}

var plural = map[int]string{
	1: "ases", 2: "doses", 3: "treses", 4: "cuatros", 5: "cincos", 6: "seises",
	7: "sietes", 10: "sotas", 11: "caballos", 12: "reyes",
}

func (p ParesInfo) String() string {
	switch p.Kind {
	case Pareja:
		return "pareja de " + plural[p.High]
	case Medias:
		return "medias de " + plural[p.High]
	case Duples:
		if p.High == p.Low {
			return "duples de " + plural[p.High] + " (cuatro iguales)"
		}
		return fmt.Sprintf("duples de %s y %s", plural[p.High], plural[p.Low])
	}
	return "sin pares"
}

// Describe resume la mano: "R R C 5 · pareja de reyes · 35 (juego)".
func (c Config) Describe(h Hand) string {
	rs := c.ranks(h)
	slices.Sort(rs)
	slices.Reverse(rs)
	var short []string
	for _, r := range rs {
		short = append(short, rankLetter(r))
	}
	pts := c.Points(h)
	j := fmt.Sprintf("%d al punto", pts)
	if pts >= 31 {
		j = fmt.Sprintf("juego de %d", pts)
	}
	return fmt.Sprintf("%s · %s · %s", strings.Join(short, " "), c.Pares(h), j)
}

func rankLetter(r int) string {
	switch r {
	case 1:
		return "A"
	case 10:
		return "S"
	case 11:
		return "C"
	case 12:
		return "R"
	}
	return fmt.Sprint(r)
}

// Apodos clásicos de algunas manos.
func (c Config) Apodo(h Hand) string {
	p := c.Pares(h)
	pts := c.Points(h)
	switch {
	case p.Kind == Duples && p.High == 12 && p.Low == 12:
		return "¡Cuatro reyes, la piara!"
	case p.Kind == Medias && p.High == 12 && pts == 31:
		return "¡Solomillo! Tres reyes y la 31."
	case p.Kind == Medias && p.High == 1 && c.Rank(h[0])+c.Rank(h[1])+c.Rank(h[2])+c.Rank(h[3]) == 15:
		return "¡Besugo! Tres ases y un rey."
	}
	return ""
}
