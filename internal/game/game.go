// Package game es el motor del mus: una máquina de estados pura que recibe
// acciones de los cuatro asientos y emite eventos. No sabe nada de la interfaz
// ni de quién es humano o máquina, para poder usarse luego en red.
//
// Asientos 0..3 en orden de juego (cada uno habla después del anterior, "por la
// derecha"). Pareja 0 = asientos 0 y 2; pareja 1 = asientos 1 y 3.
package game

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"slices"
	"strings"

	"ordagomus/internal/cards"
	"ordagomus/internal/rules"
)

type Phase int

const (
	PhaseMus Phase = iota
	PhaseDescarte
	PhaseApuesta
	PhaseFinMano
	PhaseFinPartida
)

type ActionKind int

const (
	ActMus ActionKind = iota
	ActCorto
	ActDescarte
	ActPaso
	ActEnvido // envite inicial o reenvite: Amount son los tantos que se añaden
	ActQuiero
	ActNoQuiero
	ActOrdago
	ActDeclara // "pares sí/no", "juego sí/no": lo canta el motor, no es jugable
)

type Action struct {
	Kind    ActionKind
	Amount  int
	Discard []int // índices (0-3) de las cartas a descartar
}

type Estado int

const (
	EnPaso Estado = iota
	Querido
	NoQuerido
	SinEnvite // solo una pareja tenía la jugada
	NoJugado  // nadie tenía pares
	OrdagoQuerido
)

type LanceResult struct {
	Lance  rules.Lance
	Estado Estado
	Amount int // tantos queridos
	Team   int // en NoQuerido / SinEnvite, la pareja que se lleva el lance
}

type EventKind int

const (
	EvInfo EventKind = iota
	EvAccion
	EvLance
	EvSena // To ve la seña que hace Seat
	EvTantos
	EvFinMano
	EvSenaTurno   // el asiento To ya puede pasar sus señas a mano
	EvSenaHecha   // To (Seat) hace una seña
	EvSenaVista   // el compañero (Seat) ha visto la seña de To
	EvSenaPillada // el rival Seat ha pillado la seña de To o de su compañero
)

// Event es lo que ocurre en la mesa. To = -1 lo ve todo el mundo; si no, solo
// ese asiento (las señas son privadas).
type Event struct {
	Kind   EventKind
	Seat   int
	To     int
	Text   string
	Action ActionKind
	Senas  []rules.Sena // en los eventos de señas, cuáles
	De     int          // en EvSenaPillada, de quién era la seña
}

type betting struct {
	lance      rules.Lance
	players    []int // participantes en orden desde la mano
	open       bool  // aún nadie ha envidado
	idx        int   // turno en la ronda de apertura
	amount     int   // envite sobre la mesa
	deje       int   // lo que se lleva el envidador si no se le quiere
	ordago     bool
	betTeam    int
	responders []int
}

type Game struct {
	Cfg   rules.Config
	Names [4]string

	rng      *rand.Rand
	deck     []cards.Card
	discards []cards.Card
	hands    [4]rules.Hand

	mano      int
	phase     Phase
	lance     rules.Lance
	turn      int
	musCount  int
	descartes [4][]int
	ndesc     int
	bet       *betting
	results   []LanceResult

	score  [2]int
	juegos [2]int
	vacas  [2]int

	firstHand  bool // primera mano: mus corrido y sin señas hasta que se corte
	senasDadas bool
	senas      [4][]rules.Sena
	vistas     [4][4][]rules.Sena // vistas[quien][a quien]: lo que le ha visto
	conoce     [4][4]bool         // quien sabe algo de la seña de a quien
	completa   [4]bool            // la seña del asiento se pasó entera (automática)

	paresDecl, juegoDecl [4]int // -1 sin declarar, 0 no, 1 sí

	juegoWinner, partidaWinner int
	recuento                   []string
	events                     []Event

	// Vista es la probabilidad de cada asiento de pillar una seña a un rival;
	// Disimulo la reduce para quien la pasa.
	Vista, Disimulo [4]float64
	// Manual son los asientos que pasan las señas a mano con HacerSena; el
	// resto las pasa enteras al repartir. Atento es la probabilidad de que un
	// asiento vea cada seña que le hace a mano su compañero.
	Manual [4]bool
	Atento [4]float64
}

var ErrNotYourTurn = errors.New("no es tu turno")
var ErrIllegal = errors.New("acción no permitida ahora")
var ErrSinSenas = errors.New("ahora no se pueden pasar señas")
var ErrSenaFalsa = errors.New("no llevas eso: con las señas no se miente")

func Team(seat int) int    { return seat % 2 }
func Partner(seat int) int { return (seat + 2) % 4 }
func next(seat int) int    { return (seat + 1) % 4 }

func New(cfg rules.Config, names [4]string, seed uint64) *Game {
	g := &Game{
		Cfg:           cfg,
		Names:         names,
		rng:           rand.New(rand.NewPCG(seed, seed^0x9e3779b97f4a7c15)),
		juegoWinner:   -1,
		partidaWinner: -1,
		firstHand:     true,
		Vista:         [4]float64{0.15, 0.15, 0.15, 0.15},
		Atento:        [4]float64{1, 1, 1, 1},
	}
	g.mano = g.rng.IntN(4)
	g.info("Se sortea la mano: la carta menor es para %s.", names[g.mano])
	g.startHand()
	return g
}

func (g *Game) TeamName(t int) string { return g.Names[t] + " y " + g.Names[t+2] }

func (g *Game) emit(e Event) { g.events = append(g.events, e) }

func (g *Game) info(format string, args ...any) {
	g.emit(Event{Kind: EvInfo, Seat: -1, To: -1, Text: fmt.Sprintf(format, args...)})
}

func (g *Game) say(seat int, k ActionKind, text string) {
	g.emit(Event{Kind: EvAccion, Seat: seat, To: -1, Action: k, Text: text})
}

// Events devuelve y vacía los eventos pendientes.
func (g *Game) Events() []Event {
	ev := g.events
	g.events = nil
	return ev
}

func (g *Game) Phase() Phase { return g.phase }
func (g *Game) Mano() int    { return g.mano }

// ToAct es el asiento que debe hablar, o -1 si nadie (fin de mano o partida).
func (g *Game) ToAct() int {
	switch g.phase {
	case PhaseMus, PhaseDescarte, PhaseApuesta:
		return g.turn
	}
	return -1
}

func (g *Game) Legal(seat int) []ActionKind {
	if seat != g.ToAct() {
		return nil
	}
	switch g.phase {
	case PhaseMus:
		return []ActionKind{ActMus, ActCorto}
	case PhaseDescarte:
		return []ActionKind{ActDescarte}
	case PhaseApuesta:
		if g.bet.open {
			return []ActionKind{ActPaso, ActEnvido, ActOrdago}
		}
		if g.bet.ordago {
			return []ActionKind{ActQuiero, ActNoQuiero}
		}
		return []ActionKind{ActQuiero, ActNoQuiero, ActEnvido, ActOrdago}
	}
	return nil
}

func (g *Game) startHand() {
	g.deck = cards.NewDeck()
	cards.Shuffle(g.deck, g.rng)
	g.discards = nil
	for i := 0; i < 4; i++ {
		for s := 0; s < 4; s++ {
			seat := (g.mano + s) % 4
			g.hands[seat][i] = g.draw()
		}
	}
	g.phase = PhaseMus
	g.turn = g.mano
	g.musCount = 0
	g.results = nil
	g.bet = nil
	g.recuento = nil
	g.senasDadas = false
	g.senas = [4][]rules.Sena{}
	g.vistas, g.conoce, g.completa = [4][4][]rules.Sena{}, [4][4]bool{}, [4]bool{}
	for i := range g.paresDecl {
		g.paresDecl[i], g.juegoDecl[i] = -1, -1
	}
	g.info("Reparte %s. Es mano %s.", g.Names[(g.mano+3)%4], g.Names[g.mano])
	if g.firstHand {
		g.info("Primera mano: mus corrido y sin señas.")
	} else {
		g.passSenas()
	}
}

func (g *Game) draw() cards.Card {
	if len(g.deck) == 0 {
		g.deck, g.discards = g.discards, nil
		cards.Shuffle(g.deck, g.rng)
		g.info("Se acaba el mazo: se barajan los descartes.")
	}
	c := g.deck[len(g.deck)-1]
	g.deck = g.deck[:len(g.deck)-1]
	return c
}

// passSenas: cada jugador pasa la seña completa a su compañero; a veces un
// rival la pilla. Los asientos manuales la pasan luego con HacerSena.
func (g *Game) passSenas() {
	g.senasDadas = true
	g.vistas, g.conoce, g.completa = [4][4][]rules.Sena{}, [4][4]bool{}, [4]bool{}
	for s := 0; s < 4; s++ {
		g.senas[s] = g.Cfg.Senas(g.hands[s])
		if g.Manual[s] {
			g.emit(Event{Kind: EvSenaTurno, Seat: s, To: s, Text: "Ya puedes pasar tus señas."})
			continue
		}
		g.completa[s] = true
		p := Partner(s)
		g.ve(p, s, g.senas[s])
		g.emit(Event{Kind: EvSena, Seat: s, To: p, Text: senaText(g.Names[s], g.senas[s]), Senas: g.senas[s]})
		if len(g.senas[s]) > 0 {
			g.pillar(s, g.senas[s])
		}
	}
}

// ve apunta que quien ha visto ss en la seña de s.
func (g *Game) ve(quien, s int, ss []rules.Sena) {
	g.conoce[quien][s] = true
	for _, x := range ss {
		if !slices.Contains(g.vistas[quien][s], x) {
			g.vistas[quien][s] = append(g.vistas[quien][s], x)
		}
	}
}

// pillar: cada rival puede ver la seña que hace s. Si la pilla, la pareja de s
// se entera (se le nota en la cara).
func (g *Game) pillar(s int, ss []rules.Sena) {
	txt := senaText(g.Names[s], ss)
	for _, rival := range []int{next(s), (s + 3) % 4} {
		if g.rng.Float64() >= g.Vista[rival]*(1-g.Disimulo[s]) {
			continue
		}
		g.ve(rival, s, ss)
		g.emit(Event{Kind: EvSena, Seat: s, To: rival, Text: "¡Seña pillada! " + txt, Senas: ss})
		que := senaSignificados(ss)
		g.emit(Event{Kind: EvSenaPillada, Seat: rival, To: s, Senas: ss, De: s,
			Text: fmt.Sprintf("¡%s te ha pillado la seña! (%s)", g.Names[rival], que)})
		g.emit(Event{Kind: EvSenaPillada, Seat: rival, To: Partner(s), Senas: ss, De: s,
			Text: fmt.Sprintf("¡%s le ha pillado la seña a %s! (%s)", g.Names[rival], g.Names[s], que)})
	}
}

// SenasAbiertas dice si ahora se pueden pasar señas (no en la primera mano
// hasta que se corte, ni con la mano acabada).
func (g *Game) SenasAbiertas() bool {
	return g.senasDadas && g.ToAct() >= 0
}

// HacerSena: un asiento manual le hace una seña a su compañero, que puede no
// verla. Se puede repetir, pero cada vez los rivales tienen otra ocasión de
// pillarla. Con las señas no se miente.
func (g *Game) HacerSena(seat int, x rules.Sena) error {
	if !g.Manual[seat] || !g.SenasAbiertas() {
		return ErrSinSenas
	}
	if !slices.Contains(g.senas[seat], x) {
		return ErrSenaFalsa
	}
	ss := []rules.Sena{x}
	g.emit(Event{Kind: EvSenaHecha, Seat: seat, To: seat, Senas: ss,
		Text: fmt.Sprintf("Haces la seña: %s (%s).", segunda(x.Gesto()), x.Significado())})
	p := Partner(seat)
	if g.rng.Float64() < g.Atento[p] {
		g.ve(p, seat, ss)
		g.emit(Event{Kind: EvSena, Seat: seat, To: p, Text: senaText(g.Names[seat], ss), Senas: ss})
		g.emit(Event{Kind: EvSenaVista, Seat: p, To: seat, Senas: ss,
			Text: fmt.Sprintf("✓ %s te ha visto la seña (%s).", g.Names[p], x.Significado())})
	}
	g.pillar(seat, ss)
	return nil
}

// segunda pasa el gesto a segunda persona: "se muerde" → "te muerdes".
func segunda(gesto string) string {
	verbo, resto, _ := strings.Cut(gesto, " ")
	switch verbo {
	case "se":
		v, r, _ := strings.Cut(resto, " ")
		return "te " + v + "s " + r
	}
	return verbo + "s " + resto
}

func senaSignificados(ss []rules.Sena) string {
	var out []string
	for _, x := range ss {
		out = append(out, x.Significado())
	}
	return strings.Join(out, " y ")
}

func senaText(name string, ss []rules.Sena) string {
	if len(ss) == 0 {
		return name + " no pasa seña."
	}
	var parts []string
	for _, s := range ss {
		parts = append(parts, fmt.Sprintf("%s (%s)", s.Gesto(), s.Significado()))
	}
	return name + " " + strings.Join(parts, " y ")
}

// Apply ejecuta la acción del asiento.
func (g *Game) Apply(seat int, a Action) error {
	if seat != g.ToAct() {
		return ErrNotYourTurn
	}
	switch g.phase {
	case PhaseMus:
		return g.applyMus(seat, a)
	case PhaseDescarte:
		return g.applyDescarte(seat, a)
	case PhaseApuesta:
		return g.applyApuesta(seat, a)
	}
	return ErrIllegal
}

func (g *Game) applyMus(seat int, a Action) error {
	switch a.Kind {
	case ActMus:
		g.say(seat, ActMus, "Mus")
		g.musCount++
		if g.musCount == 4 {
			g.info("Mus corrido. A descartarse.")
			g.phase = PhaseDescarte
			g.descartes = [4][]int{}
			g.ndesc = 0
			g.turn = g.mano
		} else {
			g.turn = next(seat)
		}
	case ActCorto:
		g.say(seat, ActCorto, "No hay mus")
		if g.firstHand {
			g.firstHand = false
			g.passSenas()
		}
		g.startLance(rules.Grande)
	default:
		return ErrIllegal
	}
	return nil
}

func (g *Game) applyDescarte(seat int, a Action) error {
	if a.Kind != ActDescarte || len(a.Discard) < 1 || len(a.Discard) > 4 {
		return fmt.Errorf("%w: hay que pedir entre 1 y 4 cartas", ErrIllegal)
	}
	seen := map[int]bool{}
	for _, i := range a.Discard {
		if i < 0 || i > 3 || seen[i] {
			return fmt.Errorf("%w: descarte inválido", ErrIllegal)
		}
		seen[i] = true
	}
	g.descartes[seat] = append([]int(nil), a.Discard...)
	n := len(a.Discard)
	g.say(seat, ActDescarte, fmt.Sprintf("%d %s", n, map[bool]string{true: "carta", false: "cartas"}[n == 1]))
	g.ndesc++
	g.turn = next(seat)
	if g.ndesc < 4 {
		return nil
	}
	for s := 0; s < 4; s++ {
		for _, i := range g.descartes[s] {
			g.discards = append(g.discards, g.hands[s][i])
		}
	}
	for k := 0; k < 4; k++ {
		s := (g.mano + k) % 4
		for _, i := range g.descartes[s] {
			g.hands[s][i] = g.draw()
		}
	}
	g.phase = PhaseMus
	g.musCount = 0
	g.turn = g.mano
	if !g.firstHand {
		g.passSenas()
	}
	return nil
}

func (g *Game) seatsFromMano(filter func(int) bool) []int {
	var out []int
	for k := 0; k < 4; k++ {
		s := (g.mano + k) % 4
		if filter == nil || filter(s) {
			out = append(out, s)
		}
	}
	return out
}

func (g *Game) startLance(l rules.Lance) {
	g.lance = l
	switch l {
	case rules.Grande, rules.Chica, rules.Punto:
		g.openBetting(l, g.seatsFromMano(nil))
	case rules.Pares, rules.Juego:
		decl := &g.paresDecl
		name := "Pares"
		if l == rules.Juego {
			decl = &g.juegoDecl
			name = "Juego"
		}
		for _, s := range g.seatsFromMano(nil) {
			has := g.Cfg.Eligible(l, g.hands[s])
			decl[s] = 0
			txt := name + " no"
			if has {
				decl[s] = 1
				txt = name + " sí"
			}
			g.say(s, ActDeclara, txt)
		}
		players := g.seatsFromMano(func(s int) bool { return decl[s] == 1 })
		teams := map[int]bool{}
		for _, s := range players {
			teams[Team(s)] = true
		}
		switch {
		case len(players) == 0 && l == rules.Pares:
			g.info("Nadie lleva pares.")
			g.results = append(g.results, LanceResult{Lance: l, Estado: NoJugado})
			g.nextLance()
		case len(players) == 0:
			g.info("Nadie lleva juego: se juega al punto.")
			g.startLance(rules.Punto)
		case len(teams) == 1:
			t := Team(players[0])
			g.info("Solo llevan %s %s: no hay envite.", strings.ToLower(name), g.TeamName(t))
			g.results = append(g.results, LanceResult{Lance: l, Estado: SinEnvite, Team: t})
			g.nextLance()
		default:
			g.openBetting(l, players)
		}
	}
}

func (g *Game) openBetting(l rules.Lance, players []int) {
	g.phase = PhaseApuesta
	g.bet = &betting{lance: l, players: players, open: true}
	g.turn = players[0]
	g.emit(Event{Kind: EvLance, Seat: -1, To: -1, Text: "A " + strings.ToUpper(l.String())})
}

func (g *Game) nextLance() {
	if g.juegoWinner >= 0 {
		g.endHand()
		return
	}
	switch g.lance {
	case rules.Grande:
		g.startLance(rules.Chica)
	case rules.Chica:
		g.startLance(rules.Pares)
	case rules.Pares:
		g.startLance(rules.Juego)
	default:
		g.endHand()
	}
}

func (g *Game) setResponders(bettor int) {
	b := g.bet
	b.betTeam = Team(bettor)
	b.responders = nil
	for k := 1; k < 4; k++ {
		s := (bettor + k) % 4
		if Team(s) != b.betTeam && contains(b.players, s) {
			b.responders = append(b.responders, s)
		}
	}
	g.turn = b.responders[0]
}

func contains(xs []int, x int) bool {
	for _, v := range xs {
		if v == x {
			return true
		}
	}
	return false
}

func (g *Game) applyApuesta(seat int, a Action) error {
	b := g.bet
	if b.open {
		switch a.Kind {
		case ActPaso:
			g.say(seat, ActPaso, "Paso")
			b.idx++
			if b.idx == len(b.players) {
				g.info("%s: en paso.", b.lance)
				g.results = append(g.results, LanceResult{Lance: b.lance, Estado: EnPaso})
				g.nextLance()
			} else {
				g.turn = b.players[b.idx]
			}
		case ActEnvido:
			if a.Amount < 2 {
				return fmt.Errorf("%w: el envite mínimo es de 2", ErrIllegal)
			}
			b.open = false
			b.amount, b.deje = a.Amount, 1
			if a.Amount == 2 {
				g.say(seat, ActEnvido, "Envido")
			} else {
				g.say(seat, ActEnvido, fmt.Sprintf("Envido %d", a.Amount))
			}
			g.setResponders(seat)
		case ActOrdago:
			b.open = false
			b.ordago, b.deje = true, 1
			g.say(seat, ActOrdago, "¡Órdago a la "+strings.ToLower(b.lance.String())+"!")
			g.setResponders(seat)
		default:
			return ErrIllegal
		}
		return nil
	}

	switch a.Kind {
	case ActQuiero:
		g.say(seat, ActQuiero, "Quiero")
		if b.ordago {
			g.resolveOrdago()
			return nil
		}
		g.results = append(g.results, LanceResult{Lance: b.lance, Estado: Querido, Amount: b.amount})
		g.nextLance()
	case ActNoQuiero:
		g.say(seat, ActNoQuiero, "No quiero")
		b.responders = b.responders[1:]
		if len(b.responders) > 0 {
			g.turn = b.responders[0]
			return nil
		}
		g.results = append(g.results, LanceResult{Lance: b.lance, Estado: NoQuerido, Team: b.betTeam})
		g.award(b.betTeam, b.deje, fmt.Sprintf("%s no querida", b.lance))
		g.nextLance()
	case ActEnvido:
		if b.ordago {
			return ErrIllegal
		}
		if a.Amount < 1 {
			return fmt.Errorf("%w: hay que subir al menos 1", ErrIllegal)
		}
		b.deje = b.amount
		b.amount += a.Amount
		g.say(seat, ActEnvido, fmt.Sprintf("%d más", a.Amount))
		g.setResponders(seat)
	case ActOrdago:
		if b.ordago {
			return ErrIllegal
		}
		b.deje = b.amount
		b.ordago = true
		g.say(seat, ActOrdago, "¡Órdago!")
		g.setResponders(seat)
	default:
		return ErrIllegal
	}
	return nil
}

func (g *Game) resolveOrdago() {
	l := g.bet.lance
	w := g.Cfg.Winner(l, g.hands, g.mano)
	t := Team(w)
	g.results = append(g.results, LanceResult{Lance: l, Estado: OrdagoQuerido, Team: t})
	g.info("¡Órdago querido! Se enseñan las cartas.")
	g.info("%s gana la %s con %s: el juego es para %s.", g.Names[w], strings.ToLower(l.String()), g.Cfg.Describe(g.hands[w]), g.TeamName(t))
	g.juegoWinner = t
	g.endHand()
}

func (g *Game) award(team, n int, reason string) {
	if n <= 0 || g.juegoWinner >= 0 {
		return
	}
	g.score[team] += n
	g.emit(Event{Kind: EvTantos, Seat: -1, To: -1, Text: fmt.Sprintf("%s: %d para %s (%d).", reason, n, g.TeamName(team), g.score[team])})
	if g.score[team] >= g.Cfg.Tantos {
		g.juegoWinner = team
	}
}

func (g *Game) teamValue(l rules.Lance, team int) int {
	v := 0
	for _, s := range []int{team, team + 2} {
		h := g.hands[s]
		switch l {
		case rules.Pares:
			v += g.Cfg.Pares(h).Value()
		case rules.Juego:
			if g.Cfg.HasJuego(h) {
				v += rules.JuegoValue(g.Cfg.Points(h))
			}
		}
	}
	return v
}

func (g *Game) endHand() {
	g.phase = PhaseFinMano
	g.bet = nil
	if g.juegoWinner < 0 {
		g.count()
	}
	g.emit(Event{Kind: EvFinMano, Seat: -1, To: -1, Text: "Se enseñan las cartas."})
	if g.juegoWinner >= 0 {
		g.finishJuego(g.juegoWinner)
	}
}

// count hace el recuento en orden: grande, chica, pares, juego o punto. El juego
// se acaba en cuanto una pareja llega a los tantos, aunque falte por contar.
func (g *Game) count() {
	add := func(team, n int, why string) {
		if team < 0 || n <= 0 || g.juegoWinner >= 0 {
			return
		}
		g.recuento = append(g.recuento, fmt.Sprintf("%2d para %s (%s)", n, g.TeamName(team), why))
		g.award(team, n, why)
	}
	for _, r := range g.results {
		l := r.Lance
		winner := g.Cfg.Winner(l, g.hands, g.mano)
		wt := -1
		if winner >= 0 {
			wt = Team(winner)
		}
		switch r.Estado {
		case EnPaso:
			switch l {
			case rules.Grande, rules.Chica, rules.Punto:
				add(wt, 1, l.String()+" en paso")
			default:
				add(wt, g.teamValue(l, wt), l.String()+" en paso")
			}
		case Querido:
			add(wt, r.Amount, fmt.Sprintf("%s, envite querido de %d", l, r.Amount))
			if l == rules.Punto {
				add(wt, 1, "tanto del punto")
			} else {
				add(wt, g.teamValue(l, wt), "valor de "+strings.ToLower(l.String()))
			}
		case NoQuerido:
			if l == rules.Punto {
				add(r.Team, 1, "tanto del punto")
			} else {
				add(r.Team, g.teamValue(l, r.Team), "valor de "+strings.ToLower(l.String()))
			}
		case SinEnvite:
			add(r.Team, g.teamValue(l, r.Team), l.String()+" sin envite")
		}
	}
}

func (g *Game) finishJuego(t int) {
	g.juegos[t]++
	g.info("¡Juego para %s!", g.TeamName(t))
	if g.juegos[t] >= g.Cfg.JuegosPorVaca/2+1 {
		g.vacas[t]++
		g.juegos = [2]int{}
		g.info("¡Vaca para %s!", g.TeamName(t))
		if g.vacas[t] >= g.Cfg.Vacas/2+1 {
			g.partidaWinner = t
			g.info("¡%s ganan la partida!", g.TeamName(t))
		}
	}
}

// NextHand reparte la siguiente mano (o cierra la partida).
func (g *Game) NextHand() {
	if g.phase != PhaseFinMano {
		return
	}
	if g.partidaWinner >= 0 {
		g.phase = PhaseFinPartida
		return
	}
	if g.juegoWinner >= 0 {
		g.score = [2]int{}
		g.juegoWinner = -1
		g.info("Nuevo juego.")
	}
	g.mano = next(g.mano)
	g.startHand()
}
