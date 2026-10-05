// Package ai decide las jugadas de los jugadores que maneja el ordenador.
// Solo ve lo mismo que vería un jugador en la mesa (game.View).
package ai

import (
	"math"
	"math/rand/v2"
	"slices"

	"ordagomus/internal/cards"
	"ordagomus/internal/game"
	"ordagomus/internal/rules"
)

// Perfil es la forma de jugar de un bot. Todos los valores van de 0 a 1;
// 0.5 es un jugador corriente.
type Perfil struct {
	Lance    [5]float64 // acierto valorando cada lance (Grande, Chica, Pares, Juego, Punto)
	Farol    float64    // ganas de envidar sin cartas
	Valentia float64    // lo alto que envida con buenas cartas
	Ordago   float64    // afición a echar y aceptar órdagos
	Musero   float64    // tendencia a darse mus en vez de cortar
	Lectura  float64    // cuánto aprovecha las señas y lo cantado
	Querer   float64    // tendencia a aceptar envites
}

func PerfilMedio() Perfil {
	return Perfil{
		Lance: [5]float64{0.7, 0.7, 0.7, 0.7, 0.7},
		Farol: 0.3, Valentia: 0.5, Ordago: 0.5, Musero: 0.5, Lectura: 0.7, Querer: 0.5,
	}
}

type Bot struct {
	P   Perfil
	rng *rand.Rand
}

func New(seed uint64, p Perfil) *Bot {
	return &Bot{P: p, rng: rand.New(rand.NewPCG(seed, seed*7+1))}
}

func (b *Bot) Decide(v game.View) game.Action {
	switch v.Phase {
	case game.PhaseMus:
		return b.mus(v)
	case game.PhaseDescarte:
		return game.Action{Kind: game.ActDescarte, Discard: b.descarte(v)}
	case game.PhaseApuesta:
		return b.apuesta(v)
	}
	return game.Action{Kind: game.ActPaso}
}

// WinProb estima, simulando repartos compatibles con lo que se sabe (cartas
// propias, señas vistas y lo cantado en pares/juego), la probabilidad de que la
// pareja del asiento gane el lance.
func WinProb(v game.View, l rules.Lance, samples int, useInfo bool, rng *rand.Rand) float64 {
	unseen := make([]cards.Card, 0, 36)
	for _, c := range cards.NewDeck() {
		if !slices.Contains(v.Hand[:], c) {
			unseen = append(unseen, c)
		}
	}
	try := func(constrained bool) (float64, int) {
		wins, n := 0, 0
		for attempt := 0; attempt < samples*25 && n < samples; attempt++ {
			rng.Shuffle(len(unseen), func(i, j int) { unseen[i], unseen[j] = unseen[j], unseen[i] })
			var hands [4]rules.Hand
			hands[v.Seat] = v.Hand
			k := 0
			ok := true
			for s := 0; s < 4 && ok; s++ {
				if s == v.Seat {
					continue
				}
				copy(hands[s][:], unseen[k:k+4])
				k += 4
				if constrained {
					ok = consistent(v, s, hands[s])
				}
			}
			if !ok {
				continue
			}
			n++
			if w := v.Cfg.Winner(l, hands, v.Mano); w >= 0 && game.Team(w) == game.Team(v.Seat) {
				wins++
			}
		}
		if n == 0 {
			return 0, 0
		}
		return float64(wins) / float64(n), n
	}
	p, n := try(useInfo)
	if useInfo && n < samples/10 {
		p, _ = try(false)
	}
	return p
}

// estima es WinProb visto con los ojos del bot: los torpes se equivocan más y
// los despistados no aprovechan las señas.
func (b *Bot) estima(v game.View, l rules.Lance, samples int) float64 {
	p := WinProb(v, l, samples, b.rng.Float64() < b.P.Lectura, b.rng)
	err := b.rng.NormFloat64() * 0.45 * (1 - b.P.Lance[l])
	return math.Min(1, math.Max(0, p+err))
}

func consistent(v game.View, s int, h rules.Hand) bool {
	if v.ParesDecl[s] >= 0 && (v.ParesDecl[s] == 1) != v.Cfg.HasPares(h) {
		return false
	}
	if v.JuegoDecl[s] >= 0 && (v.JuegoDecl[s] == 1) != v.Cfg.HasJuego(h) {
		return false
	}
	ss, ok := v.Senas[s]
	if !ok {
		return true
	}
	todas := v.Cfg.Senas(h)
	if v.SenaCompleta[s] {
		return slices.Equal(ss, todas)
	}
	// Seña a medias (pasada a mano): lo visto tiene que estar, lo demás no se sabe.
	for _, x := range ss {
		if !slices.Contains(todas, x) {
			return false
		}
	}
	return true
}

func (b *Bot) mus(v game.View) game.Action {
	cfg := v.Cfg
	p := cfg.Pares(v.Hand)
	pts := cfg.Points(v.Hand)
	reyes := 0
	for _, c := range v.Hand {
		if cfg.Rank(c) == 12 {
			reyes++
		}
	}
	cut := false
	switch {
	case pts == 31, p.Kind >= rules.Medias:
		cut = true
	case p.Kind == rules.Pareja && p.High == 12 && cfg.HasJuego(v.Hand):
		cut = true
	case reyes >= 2 && (cfg.HasJuego(v.Hand) || p.Kind > rules.SinPares):
		cut = b.rng.Float64() > b.P.Musero*0.4
	default:
		// Fuerza media en los lances, con lo que diga el compañero por señas.
		g := b.estima(v, rules.Grande, 150)
		c := b.estima(v, rules.Chica, 150)
		j := b.estima(v, rules.Juego, 150)
		s := (g + c + j) / 3
		cut = s > 0.5+b.P.Musero*0.2
	}
	if cut {
		return game.Action{Kind: game.ActCorto}
	}
	return game.Action{Kind: game.ActMus}
}

// descarte: se va a reyes y conserva los pares; nunca se queda con las cuatro.
// Quien es bueno a chica guarda también los ases.
func (b *Bot) descarte(v game.View) []int {
	cfg := v.Cfg
	count := map[int]int{}
	reyes := 0
	for _, c := range v.Hand {
		r := cfg.Rank(c)
		count[r]++
		if r == 12 {
			reyes++
		}
	}
	chiquero := b.P.Lance[rules.Chica] > 0.85 && reyes == 0
	var drop []int
	for i, c := range v.Hand {
		r := cfg.Rank(c)
		keep := r == 12 || count[r] >= 2 || (chiquero && r == 1)
		if !keep {
			drop = append(drop, i)
		}
	}
	// Los novatos a veces se guardan lo que no deben.
	if torpeza := 1 - (b.P.Lance[rules.Grande]+b.P.Lance[rules.Pares])/2; len(drop) > 1 && b.rng.Float64() < torpeza-0.3 {
		drop = drop[:len(drop)-1]
	}
	if len(drop) == 0 {
		// Suelta la carta más baja.
		worst := 0
		for i, c := range v.Hand {
			if cfg.Rank(c) < cfg.Rank(v.Hand[worst]) {
				worst = i
			}
		}
		drop = []int{worst}
	}
	return drop
}

func (b *Bot) apuesta(v game.View) game.Action {
	bet := v.Bet
	team := game.Team(v.Seat)
	mine, theirs := v.Score[team], v.Score[1-team]
	target := v.Cfg.Tantos
	pf := b.P
	p := b.estima(v, bet.Lance, 300) + (pf.Valentia-0.5)*0.1
	r := b.rng.Float64()
	ordagoChance := 0.12 * pf.Ordago

	if bet.Open {
		desesperados := target-theirs <= 5 && target-mine > 15
		switch {
		case p > 0.93 && (r < ordagoChance || target-theirs <= 8):
			return game.Action{Kind: game.ActOrdago}
		case desesperados && p > 0.65-pf.Ordago*0.2:
			return game.Action{Kind: game.ActOrdago}
		case p > 0.8:
			return game.Action{Kind: game.ActEnvido, Amount: 2 + b.rng.IntN(1+int(pf.Valentia*5))}
		case p > 0.6:
			return game.Action{Kind: game.ActEnvido, Amount: 2}
		case p < 0.35 && r < pf.Farol*0.3:
			return game.Action{Kind: game.ActEnvido, Amount: 2 + b.rng.IntN(2)}
		}
		return game.Action{Kind: game.ActPaso}
	}

	// Si no querer les da el juego, hay que querer.
	if theirs+bet.Deje >= target {
		return game.Action{Kind: game.ActQuiero}
	}
	if bet.Ordago {
		thr := 0.72
		if target-theirs <= 8 {
			thr = 0.5 // nos van a ganar igualmente
		}
		if target-mine <= 8 && target-theirs > 15 {
			thr = 0.85 // vamos sobrados, no hace falta arriesgar
		}
		thr -= (pf.Ordago - 0.5) * 0.2
		if p >= thr {
			return game.Action{Kind: game.ActQuiero}
		}
		return game.Action{Kind: game.ActNoQuiero}
	}
	switch {
	case p > 0.93 && (r < ordagoChance || target-theirs <= 8):
		return game.Action{Kind: game.ActOrdago}
	case p > 0.78 && bet.Amount < 10:
		return game.Action{Kind: game.ActEnvido, Amount: 2 + b.rng.IntN(1+int(pf.Valentia*4))}
	case p < 0.3 && r < pf.Farol*0.15 && bet.Amount <= 4:
		return game.Action{Kind: game.ActEnvido, Amount: 2}
	}
	thr := 0.42 + 0.03*float64(min(bet.Amount, 10)) - (pf.Querer-0.5)*0.3
	if p >= thr {
		return game.Action{Kind: game.ActQuiero}
	}
	return game.Action{Kind: game.ActNoQuiero}
}
