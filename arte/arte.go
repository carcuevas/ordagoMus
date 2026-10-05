// Package arte contiene la baraja ilustrada: Heraclio Fournier, 1878
// (dominio público, Museo Fournier de Naipes de Vitoria-Gasteiz, vía Wikimedia Commons).
package arte

import (
	"embed"
	"fmt"

	"ordagomus/internal/cards"
)

//go:embed cartas/*.png
var fs embed.FS

// Carta devuelve el PNG de la carta.
func Carta(c cards.Card) []byte {
	data, err := fs.ReadFile(fmt.Sprintf("cartas/%s_%d.png", c.Suit, c.Rank))
	if err != nil {
		panic(err)
	}
	return data
}

// Reverso devuelve el PNG del dorso de las cartas.
func Reverso() []byte {
	data, err := fs.ReadFile("cartas/reverso.png")
	if err != nil {
		panic(err)
	}
	return data
}
