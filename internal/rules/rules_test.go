package rules

import (
	"testing"

	"ordagomus/internal/cards"
)

func h(ranks ...int) Hand {
	var out Hand
	for i, r := range ranks {
		out[i] = cards.Card{Rank: r, Suit: cards.Suit(i)}
	}
	return out
}

func TestOchoReyes(t *testing.T) {
	c := DefaultConfig()
	p := c.Pares(h(3, 12, 2, 1))
	if p.Kind != Duples || p.High != 12 || p.Low != 1 {
		t.Fatalf("3+R y 2+A deberían ser duples de reyes y ases, got %+v", p)
	}
	c.OchoReyes = false
	if c.Pares(h(3, 12, 2, 1)).Kind != SinPares {
		t.Fatal("con 4 reyes no hay pares")
	}
}

func TestPuntos(t *testing.T) {
	c := DefaultConfig()
	cases := map[int]Hand{
		31: h(12, 3, 11, 1), // 10+10+10+1
		40: h(12, 11, 10, 3),
		33: h(12, 11, 7, 6),
		25: h(12, 11, 4, 2), // con 8 reyes el 2 vale 1
	}
	for want, hand := range cases {
		if got := c.Points(hand); got != want {
			t.Errorf("%v: want %d got %d", hand, want, got)
		}
	}
}

func TestOrdenJuego(t *testing.T) {
	order := []int{31, 32, 40, 37, 36, 35, 34, 33}
	for i := 0; i+1 < len(order); i++ {
		if juegoOrder(order[i]) <= juegoOrder(order[i+1]) {
			t.Errorf("%d debería ganar a %d", order[i], order[i+1])
		}
	}
}

func TestComparar(t *testing.T) {
	c := DefaultConfig()
	tests := []struct {
		l    Lance
		a, b Hand
		want int
	}{
		{Grande, h(12, 12, 1, 1), h(12, 11, 11, 11), 1},
		{Grande, h(12, 7, 5, 4), h(12, 7, 5, 4), 0},
		{Chica, h(1, 1, 4, 12), h(1, 4, 4, 12), 1},
		{Pares, h(12, 12, 1, 1), h(11, 11, 11, 11), 1},    // ganan los duples con el par mayor
		{Pares, h(11, 11, 11, 11), h(12, 12, 10, 10), -1}, // caballos-caballos < reyes-sotas
		{Pares, h(5, 5, 5, 12), h(12, 12, 11, 11), -1},    // medias < duples
		{Pares, h(12, 12, 4, 5), h(11, 11, 11, 4), -1},    // pareja < medias
		{Juego, h(12, 12, 11, 1), h(12, 11, 10, 10), 1},   // 31 > 40
		{Juego, h(12, 11, 10, 10), h(12, 12, 7, 6), 1},    // 40 > 33
		{Punto, h(12, 11, 7, 1), h(12, 11, 5, 4), -1},     // 28 < 29
	}
	for i, tt := range tests {
		if got := c.Compare(tt.l, tt.a, tt.b); got != tt.want {
			t.Errorf("caso %d (%v): want %d got %d", i, tt.l, tt.want, got)
		}
	}
}

func TestEmpateGanaLaMano(t *testing.T) {
	c := DefaultConfig()
	same := h(12, 12, 11, 1)
	hands := [4]Hand{same, same, same, same}
	for mano := 0; mano < 4; mano++ {
		if w := c.Winner(Juego, hands, mano); w != mano {
			t.Errorf("con mano %d gana %d", mano, w)
		}
	}
}

func TestSenas(t *testing.T) {
	c := DefaultConfig()
	got := c.Senas(h(12, 3, 11, 1)) // pareja de reyes y 31
	if len(got) != 2 || got[0] != SenaDosReyes || got[1] != Sena31 {
		t.Errorf("want dos reyes + 31, got %v", got)
	}
	if s := c.Senas(h(1, 4, 5, 6)); len(s) != 1 || s[0] != SenaCiego {
		t.Errorf("want ciego, got %v", s)
	}
	if s := c.Senas(h(12, 11, 7, 1)); len(s) != 1 || s[0] != SenaCiego { // 28 sin pares
		t.Errorf("want ciego con 28, got %v", s)
	}
	if s := c.Senas(h(12, 11, 5, 4)); len(s) != 1 || s[0] != SenaCiego { // 29 al punto sin pares
		t.Errorf("want ciego con 29, got %v", s)
	}
	if s := c.Senas(h(6, 6, 1, 5)); len(s) != 1 || s[0] != SenaCiego { // pareja de seises
		t.Errorf("want ciego con pareja de seises, got %v", s)
	}
	if s := c.Senas(h(4, 4, 7, 1)); len(s) != 1 || s[0] != SenaCiego { // pareja de cuatros
		t.Errorf("want ciego con pareja de cuatros, got %v", s)
	}
	if s := c.Senas(h(6, 6, 10, 11)); len(s) != 0 { // pareja baja pero con juego (32)
		t.Errorf("con juego no hay ciego, got %v", s)
	}
	if s := c.Senas(h(6, 6, 6, 1)); len(s) != 0 { // medias de seises: no es ciego
		t.Errorf("con medias no hay ciego, got %v", s)
	}
	if s := c.Senas(h(1, 1, 5, 6)); len(s) != 1 || s[0] != SenaDosAses { // dos ases: su seña, no ciego
		t.Errorf("want dos ases, got %v", s)
	}
}
