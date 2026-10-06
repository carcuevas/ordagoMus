// Package cards modela la baraja española de 40 cartas.
package cards

import (
	"fmt"
	"math/rand/v2"

	"ordagomus/internal/i18n"
)

type Suit int

const (
	Oros Suit = iota
	Copas
	Espadas
	Bastos
)

var suitNames = [...]string{"oros", "copas", "espadas", "bastos"}

func (s Suit) String() string { return i18n.T(suitNames[s]) }

// Card es una carta: Rank 1-7, 10 (sota), 11 (caballo), 12 (rey).
type Card struct {
	Rank int
	Suit Suit
}

// Ranks son los números de la baraja española sin ochos ni nueves.
var Ranks = []int{1, 2, 3, 4, 5, 6, 7, 10, 11, 12}

// Label es el nombre corto de la carta: "As", "Sota", "Rey", "5"...
func (c Card) Label() string {
	switch c.Rank {
	case 1:
		return i18n.T("As")
	case 10:
		return i18n.T("Sota")
	case 11:
		return i18n.T("Caballo")
	case 12:
		return i18n.T("Rey")
	}
	return fmt.Sprint(c.Rank)
}

// ShortLabel cabe en una carta dibujada en terminal.
func (c Card) ShortLabel() string {
	if c.Rank == 11 {
		return i18n.T("Cab")
	}
	return c.Label()
}

func (c Card) String() string { return i18n.Tf("%s de %s", c.Label(), c.Suit.String()) }

// NewDeck devuelve la baraja ordenada de 40 cartas.
func NewDeck() []Card {
	deck := make([]Card, 0, 40)
	for s := Oros; s <= Bastos; s++ {
		for _, r := range Ranks {
			deck = append(deck, Card{Rank: r, Suit: s})
		}
	}
	return deck
}

func Shuffle(deck []Card, rng *rand.Rand) {
	rng.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })
}
