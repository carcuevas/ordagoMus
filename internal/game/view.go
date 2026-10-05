package game

import "ordagomus/internal/rules"

// BetView es el envite en curso tal como lo ve la mesa.
type BetView struct {
	Lance   rules.Lance
	Open    bool
	Amount  int
	Deje    int
	Ordago  bool
	BetTeam int
	Players []int
}

// View es lo que un asiento puede saber de la partida: sus cartas y lo público.
// Es lo único que reciben la IA y (en el futuro) los clientes de red.
type View struct {
	Seat     int
	Cfg      rules.Config
	Names    [4]string
	Hand     rules.Hand
	Hands    [4]rules.Hand // todas, solo cuando Revealed
	Revealed bool

	Phase Phase
	Lance rules.Lance
	Mano  int
	Turn  int
	Legal []ActionKind
	Bet   *BetView

	Score, Juegos, Vacas [2]int

	ParesDecl, JuegoDecl [4]int               // -1 sin declarar, 0 no, 1 sí
	Senas                map[int][]rules.Sena // señas conocidas: compañero y pilladas
	SenaCompleta         [4]bool              // la seña conocida es entera (si no, puede faltar algo)
	SenasAbiertas        bool                 // se pueden pasar señas a mano

	Results       []LanceResult
	Recuento      []string
	JuegoWinner   int
	PartidaWinner int
	FirstHand     bool
}

func (g *Game) View(seat int) View {
	v := View{
		Seat:          seat,
		Cfg:           g.Cfg,
		Names:         g.Names,
		Hand:          g.hands[seat],
		Phase:         g.phase,
		Lance:         g.lance,
		Mano:          g.mano,
		Turn:          g.ToAct(),
		Legal:         g.Legal(seat),
		Score:         g.score,
		Juegos:        g.juegos,
		Vacas:         g.vacas,
		ParesDecl:     g.paresDecl,
		JuegoDecl:     g.juegoDecl,
		Senas:         map[int][]rules.Sena{},
		Results:       append([]LanceResult(nil), g.results...),
		Recuento:      append([]string(nil), g.recuento...),
		JuegoWinner:   g.juegoWinner,
		PartidaWinner: g.partidaWinner,
		FirstHand:     g.firstHand,
	}
	if g.phase == PhaseFinMano || g.phase == PhaseFinPartida {
		v.Revealed = true
		v.Hands = g.hands
	} else {
		v.Hands[seat] = g.hands[seat]
	}
	if g.senasDadas {
		v.SenasAbiertas = g.SenasAbiertas()
		for s := 0; s < 4; s++ {
			if s != seat && g.conoce[seat][s] {
				v.Senas[s] = append([]rules.Sena(nil), g.vistas[seat][s]...)
				v.SenaCompleta[s] = g.completa[s]
			}
		}
	}
	if b := g.bet; b != nil {
		v.Bet = &BetView{
			Lance: b.lance, Open: b.open, Amount: b.amount, Deje: b.deje,
			Ordago: b.ordago, BetTeam: b.betTeam, Players: append([]int(nil), b.players...),
		}
	}
	return v
}
