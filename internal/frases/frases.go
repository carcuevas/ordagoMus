// Package frases es el refranero y la jerga de mus que sueltan los jugadores
// del ordenador. Para añadir frases basta con ampliar las listas.
package frases

import (
	"math/rand/v2"

	"ordagomus/internal/game"
)

var porAccion = map[game.ActionKind][]string{
	game.ActMus: {
		"Mus.", "Mus, que estas no valen ni para el café.", "Mus hasta que salgan los reyes.",
		"Dame mus, compañero.", "Mus corrido.", "Corta con buenas, compañero.", "Con esto no se va ni a la esquina. Mus.",
	},
	game.ActCorto: {
		"No hay mus.", "Corto.", "Quita, que llevo cosas.", "Aquí se corta.",
		"Hasta aquí ha llegado el mus.",
	},
	game.ActDescarte: {
		"A reyes.", "Dame de las buenas.", "Esta vez sí.", "Que venga la 31.",
	},
	game.ActPaso: {
		"Paso.", "Paso, por mi parte.", "Paso, a ver qué hacéis.", "Llegó a mí... paso.",
	},
	game.ActEnvido: {
		"Envido.", "Ahí van dos.", "Dos más, por no callar.", "Envido, que hay que animar esto.",
		"El que no envida no gana.", "Las de Hontanares.",
	},
	game.ActQuiero: {
		"Quiero.", "Lo quiero.", "Veo.", "Queremos.", "Con lo que sea.",
	},
	game.ActNoQuiero: {
		"No quiero.", "Para vosotros.", "Llevaos el tanto.", "No queremos, más se perdió en Cuba.",
	},
	game.ActOrdago: {
		"¡Órdago!", "¡Órdago, y que salga el sol por Antequera!", "¡Ahí va el órdago!",
		"¡Órdago! A ver quién tiene lo que hay que tener.", "Se acabó lo que se daba: ¡órdago!",
	},
}

var ganamos = []string{
	"¿Os rendís?", "Esto es mus de los de antes.", "La mano es la mano.",
	"A estas cartas no les gana ni el cura del pueblo.", "Manitas de plata, ¿eh?",
}

var perdemos = []string{
	"Manitas de plata...", "Vaya cartitas os han tocado.", "Ya vendrán tiempos mejores.",
	"Más se perdió en Cuba.", "No hay mal que cien años dure.",
}

func pick(rng *rand.Rand, xs []string) string { return xs[rng.IntN(len(xs))] }

// Para devuelve una frase para la acción, o "" si el jugador prefiere callar.
func Para(k game.ActionKind, rng *rand.Rand, prob float64) string {
	xs := porAccion[k]
	if len(xs) == 0 || rng.Float64() >= prob {
		return ""
	}
	return pick(rng, xs)
}

// Final devuelve un comentario al acabar la mano según si se ha ganado.
func Final(gano bool, rng *rand.Rand) string {
	if gano {
		return pick(rng, ganamos)
	}
	return pick(rng, perdemos)
}
