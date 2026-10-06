package tui

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"ordagomus/internal/ai"
	"ordagomus/internal/ajustes"
	"ordagomus/internal/frases"
	"ordagomus/internal/game"
	"ordagomus/internal/i18n"
	"ordagomus/internal/personajes"
	"ordagomus/internal/rules"
	"ordagomus/internal/sonido"
)

type botMsg struct{ gen int }
type volverMenuMsg struct{}
type gestoMsg struct{ gen, seat, g int }

// gestoPausa es lo que dura cada gesto de una seña escrita.
const gestoPausa = 1300 * time.Millisecond

// gesto es un instante de la cara: una seña o, con sena nil, la cara normal
// (lo que tarda en hacerla o entre una y otra).
type gesto struct {
	sena *rules.Sena
	dur  time.Duration
}

// senaHecha es una seña que has pasado tú y lo que ha pasado con ella.
type senaHecha struct {
	sena     rules.Sena
	vista    bool
	pillada  []string
	repetida int
}

var mesaGen int

// Mesa es la pantalla de juego: un humano en el asiento 0 y tres personajes.
type Mesa struct {
	g     *game.Game
	human int
	pjs   [4]*personajes.Personaje // nil para el humano
	bots  [4]*ai.Bot
	rng   *rand.Rand
	vel   ajustes.Velocidad
	img   bool // cartas ilustradas
	gen   int

	modo      ajustes.ModoSenas
	log       []string
	bubble    [4]string
	flavor    [4]string
	senas     []string
	expr      [4]expresion
	gestos    [4][]gesto // lo que se le está viendo hacer en la cara
	gestoVivo [4]bool
	gestoGen  [4]int
	marca     [4]string // "te ha visto", "te ha pillado la seña"...
	misSenas  []senaHecha
	menuSena  bool
	fase      game.Phase
	sel       [4]bool
	input     string
	typing    bool
	confirm   bool
	status    string
	startSc   [2]int
	newLance  bool
	sonido    bool
	tallas    [2]int // tamaño de las cartas ilustradas: propias y ajenas
	tallaPara [2]int // tamaño del terminal para el que se eligió

	width  int
	height int
}

func NuevaMesa(aj ajustes.Ajustes, companero *personajes.Personaje, rivales [2]*personajes.Personaje, seed uint64) *Mesa {
	pjs := [4]*personajes.Personaje{nil, rivales[0], companero, rivales[1]}
	names := [4]string{aj.Nombre}
	if aj.Nombre == ajustes.Defecto().Nombre {
		names[0] = i18n.T("Tú") // el nombre por defecto, en el idioma elegido
	}
	for i := 1; i < 4; i++ {
		names[i] = pjs[i].Corto()
	}
	g := game.New(aj.Reglas, names, seed)
	g.Vista[0], g.Disimulo[0] = 0.2, 0.5
	g.Manda[0] = true // los envites de los rivales los contestas tú
	m := &Mesa{g: g, pjs: pjs, vel: aj.Velocidad, img: gfx.on && aj.Imagenes, modo: aj.Senas, sonido: aj.Sonido, rng: rand.New(rand.NewPCG(seed, 99))}
	for i := 1; i < 4; i++ {
		m.bots[i] = ai.New(seed+uint64(i)*1000, pjs[i].Perfil)
		g.Vista[i], g.Disimulo[i] = pjs[i].Vista, pjs[i].Disimulo
	}
	if m.modo == ajustes.SenasDeVerdad {
		// Las tuyas las pasas tú, y tu compañero puede no verlas.
		g.Manual[0] = true
		g.Atento[2] = min(0.97, 0.75+pjs[2].Vista)
	}
	mesaGen++
	m.gen = mesaGen
	m.fase = game.PhaseFinMano // para que suene el primer reparto
	m.drain()
	return m
}

func (m *Mesa) Init() tea.Cmd { return m.schedule() }

func (m *Mesa) botTurn() bool {
	s := m.g.ToAct()
	return s >= 0 && s != m.human
}

func (m *Mesa) schedule() tea.Cmd {
	cmds := []tea.Cmd{m.botTick()}
	for s := range m.gestos {
		cmds = append(cmds, m.gestoTick(s))
	}
	return tea.Batch(cmds...)
}

func entre(rng *rand.Rand, lo, hi int) time.Duration {
	return time.Duration(lo+rng.IntN(hi-lo+1)) * time.Millisecond
}

// encolarGestos pone en la cara de s las señas que se le ven. Escritas, cada
// gesto dura lo mismo; de verdad, como en la mesa: tarda un poco en hacerlas y
// cada una dura un instante, entre 500 y 1000 ms.
func (m *Mesa) encolarGestos(s int, ss []rules.Sena) {
	for i := range ss {
		x := &ss[i]
		if m.modo == ajustes.SenasEscritas {
			m.gestos[s] = append(m.gestos[s], gesto{x, gestoPausa})
			continue
		}
		espera := entre(m.rng, 150, 450)
		if i == 0 {
			espera = entre(m.rng, 400, 2000)
		}
		m.gestos[s] = append(m.gestos[s], gesto{nil, espera}, gesto{x, entre(m.rng, 500, 1000)})
	}
}

// gestoTick arranca el reloj de los gestos del asiento si tiene alguno pendiente.
func (m *Mesa) gestoTick(s int) tea.Cmd {
	if m.gestoVivo[s] || len(m.gestos[s]) == 0 {
		return nil
	}
	m.gestoVivo[s] = true
	msg := gestoMsg{gen: m.gen, seat: s, g: m.gestoGen[s]}
	return tea.Tick(m.gestos[s][0].dur, func(time.Time) tea.Msg { return msg })
}

func (m *Mesa) siguienteGesto(s int) tea.Cmd {
	m.gestoVivo[s] = false
	if len(m.gestos[s]) > 0 {
		m.gestos[s] = m.gestos[s][1:]
	}
	return m.gestoTick(s)
}

// limpiarGestos corta los gestos en curso (los relojes viejos se ignoran).
func (m *Mesa) limpiarGestos() {
	for s := range m.gestos {
		m.gestos[s], m.gestoVivo[s] = nil, false
		m.gestoGen[s]++
	}
}

func (m *Mesa) botTick() tea.Cmd {
	if !m.botTurn() || m.vel == ajustes.PasoAPaso {
		return nil
	}
	d := m.vel.Pausa()
	if m.newLance {
		d = d * 8 / 5 // respiro al empezar un lance
	}
	gen := m.gen
	return tea.Tick(d, func(time.Time) tea.Msg { return botMsg{gen: gen} })
}

func (m *Mesa) addLog(s string) {
	m.log = append(m.log, s)
	if len(m.log) > 300 {
		m.log = m.log[len(m.log)-300:]
	}
}

// drain pasa los eventos del motor al registro y a los bocadillos.
func (m *Mesa) drain() {
	m.newLance = false
	if f := m.g.Phase(); m.fase == game.PhaseDescarte && f == game.PhaseMus {
		// Cartas nuevas: las señas de antes ya no valen.
		m.marca, m.misSenas, m.senas = [4]string{}, nil, nil
	}
	toca := sinSonido
	if f := m.g.Phase(); f == game.PhaseMus && m.fase != game.PhaseMus {
		toca = sonido.Reparto
	}
	m.fase = m.g.Phase()
	defer func() {
		if m.sonido && toca != sinSonido {
			sonido.Tocar(toca)
		}
	}()
	for _, e := range m.g.Events() {
		if e.To != -1 && e.To != m.human {
			continue
		}
		switch e.Kind {
		case game.EvAccion:
			m.bubble[e.Seat] = e.Text
			m.flavor[e.Seat] = ""
			if pj := m.pjs[e.Seat]; pj != nil && e.Action != game.ActDeclara {
				if m.rng.Float64() < 0.15 {
					m.flavor[e.Seat] = pj.Frase(m.rng)
				} else if f := frases.Para(e.Action, m.rng, 0.25); f != "" && !mismo(f, e.Text) {
					m.flavor[e.Seat] = f
				}
			}
			m.addLog(m.g.Names[e.Seat] + ": " + e.Text)
			m.expr[e.Seat] = m.exprPara(e.Seat, e.Action)
			switch e.Action {
			case game.ActCorto:
				toca = masUrgente(toca, sonido.Corte)
			case game.ActEnvido:
				toca = masUrgente(toca, sonido.Envite)
			case game.ActOrdago:
				toca = masUrgente(toca, sonido.Ordago)
			}
		case game.EvLance:
			m.bubble, m.flavor = [4]string{}, [4]string{}
			m.expr = [4]expresion{}
			m.newLance = true
			toca = masUrgente(toca, sonido.Lance)
			m.addLog(styleLance.Render("── " + e.Text + " ──"))
		case game.EvSena:
			if e.Seat != m.human {
				// la seña se ve en la cara: un gesto detrás de otro
				m.encolarGestos(e.Seat, append([]rules.Sena(nil), e.Senas...))
			}
			if m.modo == ajustes.SenasEscritas {
				m.senas = append(m.senas, e.Text)
				m.addLog(styleSena.Render(i18n.Tf("(seña) %s", e.Text)))
			}
		case game.EvSenaTurno:
			m.misSenas = nil
		case game.EvSenaHecha:
			m.addLog(styleSena.Render(e.Text))
			m.hecha(e.Senas[0]).repetida++
		case game.EvSenaVista:
			m.addLog(styleTantos.Render(e.Text))
			m.hecha(e.Senas[0]).vista = true
			m.marca[e.Seat] = styleTantos.Render(i18n.T("✓ vio tu seña"))
			if pj := m.pjs[e.Seat]; pj != nil && pj.Disimulo < 0.9 {
				m.expr[e.Seat] = exContento
			}
		case game.EvSenaPillada:
			m.addLog(styleErr.Render(e.Text))
			m.status = e.Text
			m.marca[e.Seat] = styleErr.Render(i18n.T("◉ pilló una seña"))
			if e.De == m.human {
				for _, x := range e.Senas {
					h := m.hecha(x)
					if !slices.Contains(h.pillada, m.g.Names[e.Seat]) {
						h.pillada = append(h.pillada, m.g.Names[e.Seat])
					}
				}
			}
		case game.EvTantos:
			m.addLog(styleTantos.Render(e.Text))
		case game.EvFinMano:
			toca = masUrgente(toca, sonido.Recuento)
			m.addLog(e.Text)
			m.comentarioFinal()
		default:
			m.addLog(e.Text)
		}
	}
}

const sinSonido sonido.Sonido = -1

// prioridadSonido: si en una tanda pasan varias cosas (cortan y empieza la
// grande, por ejemplo), suena solo la más importante.
var prioridadSonido = map[sonido.Sonido]int{
	sinSonido: 0, sonido.Lance: 1, sonido.Reparto: 2, sonido.Corte: 3,
	sonido.Envite: 4, sonido.Recuento: 5, sonido.Ordago: 6,
}

func masUrgente(a, b sonido.Sonido) sonido.Sonido {
	if prioridadSonido[b] > prioridadSonido[a] {
		return b
	}
	return a
}

// hecha devuelve el apunte de una seña tuya (lo crea si hace falta).
func (m *Mesa) hecha(x rules.Sena) *senaHecha {
	for i := range m.misSenas {
		if m.misSenas[i].sena == x {
			return &m.misSenas[i]
		}
	}
	m.misSenas = append(m.misSenas, senaHecha{sena: x})
	return &m.misSenas[len(m.misSenas)-1]
}

// exprPara es la cara que pone el personaje al hacer la jugada. Al que no se
// le lee la cara (Xosé) solo se le ve hablar.
func (m *Mesa) exprPara(s int, a game.ActionKind) expresion {
	ex := exNormal
	switch a {
	case game.ActMus, game.ActCorto, game.ActQuiero, game.ActDeclara:
		ex = exHabla
	case game.ActEnvido:
		ex = exContento
	case game.ActOrdago:
		ex = exEnfadado
	case game.ActNoQuiero:
		ex = exTriste
	}
	if pj := m.pjs[s]; pj != nil && pj.Disimulo >= 0.9 && ex != exHabla {
		ex = exNormal
	}
	return ex
}

func mismo(a, b string) bool {
	return strings.EqualFold(strings.Trim(a, ".¡!"), strings.Trim(b, ".¡!"))
}

func (m *Mesa) comentarioFinal() {
	v := m.g.View(m.human)
	d0, d1 := v.Score[0]-m.startSc[0], v.Score[1]-m.startSc[1]
	if v.JuegoWinner >= 0 {
		d0, d1 = 0, 0
		if v.JuegoWinner == 0 {
			d0 = 1
		} else {
			d1 = 1
		}
	}
	if d0 == d1 {
		return
	}
	winTeam := 0
	if d1 > d0 {
		winTeam = 1
	}
	for s, pj := range m.pjs {
		switch {
		case pj == nil || pj.Disimulo >= 0.9:
		case game.Team(s) == winTeam:
			m.expr[s] = exContento
		case pj.Perfil.Valentia >= 0.8:
			m.expr[s] = exEnfadado
		default:
			m.expr[s] = exTriste
		}
	}
	if m.rng.Float64() > 0.7 {
		return
	}
	s := 1 + m.rng.IntN(3)
	if game.Team(s) == winTeam && m.rng.Float64() < 0.5 {
		m.flavor[s] = m.pjs[s].Frase(m.rng)
	} else {
		m.flavor[s] = frases.Final(game.Team(s) == winTeam, m.rng)
	}
}

func (m *Mesa) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case botMsg:
		if msg.gen != m.gen {
			return nil
		}
		return m.botStep()
	case gestoMsg:
		if msg.gen != m.gen || msg.g != m.gestoGen[msg.seat] {
			return nil
		}
		return m.siguienteGesto(msg.seat)
	case tea.KeyMsg:
		return m.key(msg.String())
	}
	return nil
}

func (m *Mesa) botStep() tea.Cmd {
	if !m.botTurn() {
		return nil
	}
	s := m.g.ToAct()
	if err := m.g.Apply(s, m.bots[s].Decide(m.g.View(s))); err != nil {
		m.status = i18n.Tf("Error del bot: %s", err.Error())
		return nil
	}
	m.drain()
	return m.schedule()
}

func (m *Mesa) act(a game.Action) tea.Cmd {
	if err := m.g.Apply(m.human, a); err != nil {
		m.status = err.Error()
		return nil
	}
	m.status = ""
	m.sel = [4]bool{}
	m.menuSena = false
	m.drain()
	return m.schedule()
}

// senaKey hace la seña número k (1-7) del menú.
func (m *Mesa) senaKey(k int) tea.Cmd {
	if err := m.g.HacerSena(m.human, rules.Sena(k-1)); err != nil {
		m.status = err.Error()
		return nil
	}
	m.status = ""
	m.drain()
	return m.schedule()
}

func volver() tea.Msg { return volverMenuMsg{} }

func (m *Mesa) key(k string) tea.Cmd {
	if m.confirm {
		m.confirm = false
		if k == "s" || k == "S" {
			return volver
		}
		return nil
	}
	v := m.g.View(m.human)
	if m.menuSena && !v.SenasAbiertas {
		m.menuSena = false
	}
	if m.menuSena && !m.typing {
		switch {
		case k == "esc" || k == "s":
			m.menuSena = false
			return nil
		case len(k) == 1 && k[0] >= '1' && k[0] <= '7':
			return m.senaKey(int(k[0] - '0'))
		}
	}
	if k == "s" && !m.typing && m.modo == ajustes.SenasDeVerdad && v.SenasAbiertas {
		m.menuSena = true
		return nil
	}
	if k == "esc" && !m.typing {
		m.confirm = true
		return nil
	}

	if m.typing {
		switch {
		case k == "esc":
			m.typing, m.input = false, ""
		case k == "backspace" && len(m.input) > 0:
			m.input = m.input[:len(m.input)-1]
		case k == "enter":
			n, err := strconv.Atoi(m.input)
			m.typing, m.input = false, ""
			if err != nil || n < 1 {
				m.status = i18n.T("Número no válido")
				return nil
			}
			return m.act(game.Action{Kind: game.ActEnvido, Amount: n})
		case len(k) == 1 && k[0] >= '0' && k[0] <= '9' && len(m.input) < 2:
			m.input += k
		}
		return nil
	}

	switch v.Phase {
	case game.PhaseFinMano:
		if k == "enter" || k == " " {
			m.g.NextHand()
			m.bubble, m.flavor = [4]string{}, [4]string{}
			m.expr, m.marca = [4]expresion{}, [4]string{}
			m.limpiarGestos()
			m.senas, m.misSenas = nil, nil
			m.startSc = m.g.View(m.human).Score
			m.addLog(strings.Repeat("─", 30))
			m.drain()
			return m.schedule()
		}
		return nil
	case game.PhaseFinPartida:
		if k == "enter" {
			return volver
		}
		return nil
	}
	if v.Turn != m.human {
		if k == " " && m.vel == ajustes.PasoAPaso {
			return m.botStep()
		}
		return nil
	}

	switch v.Phase {
	case game.PhaseMus:
		switch k {
		case "m":
			return m.act(game.Action{Kind: game.ActMus})
		case "c":
			return m.act(game.Action{Kind: game.ActCorto})
		}
	case game.PhaseDescarte:
		switch k {
		case "1", "2", "3", "4":
			i := int(k[0] - '1')
			m.sel[i] = !m.sel[i]
		case "t":
			all := m.sel != [4]bool{true, true, true, true}
			m.sel = [4]bool{all, all, all, all}
		case "enter":
			var idx []int
			for i, s := range m.sel {
				if s {
					idx = append(idx, i)
				}
			}
			if len(idx) == 0 {
				m.status = i18n.T("Hay que pedir al menos una carta")
				return nil
			}
			return m.act(game.Action{Kind: game.ActDescarte, Discard: idx})
		}
	case game.PhaseApuesta:
		legal := map[game.ActionKind]bool{}
		for _, a := range v.Legal {
			legal[a] = true
		}
		switch {
		case k == "p" && legal[game.ActPaso]:
			return m.act(game.Action{Kind: game.ActPaso})
		case k == "e" && legal[game.ActEnvido]:
			return m.act(game.Action{Kind: game.ActEnvido, Amount: 2})
		case k == "n" && legal[game.ActEnvido]:
			m.typing = true
		case k == "o" && legal[game.ActOrdago]:
			return m.act(game.Action{Kind: game.ActOrdago})
		case k == "q" && legal[game.ActQuiero]:
			return m.act(game.Action{Kind: game.ActQuiero})
		case k == "x" && legal[game.ActNoQuiero]:
			return m.act(game.Action{Kind: game.ActNoQuiero})
		}
	}
	return nil
}

// ───────────────────────────── dibujo ─────────────────────────────

func (m *Mesa) seatTag(v game.View, s int) string {
	var tags []string
	if s == v.Mano {
		tags = append(tags, "mano")
	}
	if s == (v.Mano+3)%4 {
		tags = append(tags, "postre")
	}
	label := v.Names[s]
	if pj := m.pjs[s]; pj != nil {
		label = pj.Completo()
	}
	name := styleName.Render(label)
	if v.Turn == s {
		name = styleTurn.Render("▶ " + label)
	}
	if pj := m.pjs[s]; pj != nil {
		tags = append([]string{pj.Origen}, tags...)
	}
	if len(tags) > 0 {
		name += styleDim.Render(" · " + strings.Join(tags, " · "))
	}
	return name
}

func (m *Mesa) speech(s, w int) []string {
	wrap := lipgloss.NewStyle().Width(w)
	var lines []string
	bocadillo := ""
	if m.bubble[s] != "" {
		bocadillo = styleBubble.Render("« "+m.bubble[s]+" »") + "  "
	}
	lines = append(lines, bocadillo+m.marca[s])
	if m.flavor[s] != "" {
		lines = append(lines, styleFlavor.Render(wrap.Render("“"+m.flavor[s]+"”")))
	} else {
		lines = append(lines, "")
	}
	return lines
}

// cara dibuja al personaje con la expresión del momento o la seña que hace.
func (m *Mesa) cara(s int, espejo bool) string {
	pj := m.pjs[s]
	if pj == nil {
		return ""
	}
	var sena *rules.Sena
	if len(m.gestos[s]) > 0 {
		sena = m.gestos[s][0].sena
	}
	return dibujarCara(pj.ID, m.expr[s], sena, espejo, m.img)
}

// player es un jugador del ordenador: nombre arriba y, debajo, la cara junto a
// las cartas y lo que dice. El de la derecha (espejo) tiene la cara a la
// derecha para que mire a la mesa.
func (m *Mesa) player(v game.View, s int, espejo bool) string {
	hand := compactHand(v.Hands[s], !v.Revealed, m.img)
	lines := []string{hand}
	if v.Revealed {
		lines = append(lines, styleDim.Render(v.Cfg.Describe(v.Hands[s])))
	}
	lines = append(lines, m.speech(s, max(lipgloss.Width(hand), 28))...)
	align := lipgloss.Left
	if espejo {
		align = lipgloss.Right
	}
	col := lipgloss.JoinVertical(align, lines...)
	row := lipgloss.JoinHorizontal(lipgloss.Top, m.cara(s, espejo), " ", col)
	if espejo {
		row = lipgloss.JoinHorizontal(lipgloss.Top, col, " ", m.cara(s, espejo))
	}
	return lipgloss.JoinVertical(lipgloss.Center, m.seatTag(v, s), row)
}

func (m *Mesa) View() string {
	if m.img {
		m.ajustarTalla()
		gfx.tb, gfx.ts = m.tallas[0], m.tallas[1]
	}
	body, w := m.cuerpo()

	// El registro ocupa lo que quede de pantalla (mínimo 4 líneas).
	logN := 6
	if m.height > 0 {
		logN = max(m.height-lipgloss.Height(body)-2, 4)
	}
	start := max(len(m.log)-logN, 0)
	logBox := styleBox.Width(w - 4).Render(strings.Join(m.log[start:], "\n"))
	return lipgloss.JoinVertical(lipgloss.Left, body, logBox)
}

// logMin es lo mínimo que se quiere ver del registro (con su marco).
const logMin = 6

// ajustarTalla elige las cartas más grandes con las que la mesa cabe en el
// terminal: las de los demás crecen si sobra ancho y las tuyas si sobra alto.
// Solo se recalcula cuando cambia el tamaño del terminal.
func (m *Mesa) ajustarTalla() {
	if len(gfx.big) == 0 || m.width == 0 || m.height == 0 || m.tallaPara == [2]int{m.width, m.height} {
		return
	}
	m.tallaPara = [2]int{m.width, m.height}
	cabe := func() bool {
		// Se mide en la fase con más altura: cartas enseñadas y recuento.
		body, _ := m.cuerpo()
		return lipgloss.Height(body)+logMin+recuentoExtra <= m.height && lipgloss.Width(body) <= max(m.width, 110)
	}
	gfx.tb, gfx.ts = 0, 0
	for gfx.ts = len(gfx.small) - 1; gfx.ts > 0 && !cabe(); gfx.ts-- {
	}
	for gfx.tb = len(gfx.big) - 1; gfx.tb > 0 && !cabe(); gfx.tb-- {
	}
	m.tallas = [2]int{gfx.tb, gfx.ts}
}

// recuentoExtra es lo que crece la mesa al enseñar las cartas (descripción de
// cada mano y recuento), para que la talla no tenga que cambiar a media mano.
const recuentoExtra = 7

// cuerpo dibuja todo menos el registro; devuelve también el ancho usado.
func (m *Mesa) cuerpo() (string, int) {
	v := m.g.View(m.human)
	me := game.Team(m.human)
	w := max(m.width, 110)

	reyes := i18n.T("4 reyes")
	if v.Cfg.OchoReyes {
		reyes = i18n.T("8 reyes")
	}
	head := styleTitle.Render("ÓRDAGO · Mus") + styleDim.Render(i18n.Tf("  a %d tantos · %s · juegos al mejor de %d · vacas al mejor de %d · velocidad %s",
		v.Cfg.Tantos, reyes, v.Cfg.JuegosPorVaca, v.Cfg.Vacas, m.vel))

	// El marcador se lleva con amarracos (5) y piedras (1), como en la mesa.
	scoreLine := func(t int, label string) string {
		txt := fmt.Sprintf("%-9s %-28s\n%-9s %s", label, m.g.TeamName(t), "",
			styleDim.Render(i18n.Tf("juegos %d · vacas %d", v.Juegos[t], v.Vacas[t])))
		tantos := styleName.Render(fmt.Sprintf("%2d", v.Score[t])) + "\n" + styleDim.Render(i18n.Tf("de %d", v.Cfg.Tantos))
		return lipgloss.JoinHorizontal(lipgloss.Top, txt, "  ", tantos, "  ", fichas(v.Score[t], m.img, 60))
	}
	score := styleBox.Render(scoreLine(me, i18n.T("Nosotros")) + "\n" + scoreLine(1-me, i18n.T("Ellos")))

	north := (m.human + 2) % 4
	east := (m.human + 1) % 4
	west := (m.human + 3) % 4

	pw, pe := m.player(v, west, false), m.player(v, east, true)
	sideW := max(46, lipgloss.Width(pw), lipgloss.Width(pe))
	middle := lipgloss.JoinHorizontal(lipgloss.Center,
		lipgloss.PlaceHorizontal(sideW, lipgloss.Center, pw),
		lipgloss.PlaceHorizontal(max(w-2*sideW, 18), lipgloss.Center, m.centerInfo(v)),
		lipgloss.PlaceHorizontal(sideW, lipgloss.Center, pe),
	)

	showIdx := v.Phase == game.PhaseDescarte && v.Turn == m.human
	south := []string{m.seatTag(v, m.human), bigHand(v.Hand, v.Cfg, m.sel, showIdx, m.img), styleDim.Render(v.Cfg.Describe(v.Hand))}
	if a := v.Cfg.Apodo(v.Hand); a != "" {
		south = append(south, styleFlavor.Render(a))
	}
	if m.bubble[m.human] != "" {
		south = append(south, styleBubble.Render("« "+m.bubble[m.human]+" »"))
	}
	if t := m.tusSenas(); t != "" {
		south = append(south, t)
	}

	table := lipgloss.JoinVertical(lipgloss.Center,
		m.player(v, north, false), middle, lipgloss.JoinVertical(lipgloss.Center, south...))
	table = lipgloss.PlaceHorizontal(w, lipgloss.Center, table)

	var extra []string
	if len(m.senas) > 0 {
		extra = append(extra, styleSena.Render(strings.Join(m.senas, "\n")))
	}
	if v.Revealed && len(v.Recuento) > 0 {
		extra = append(extra, styleTantos.Render(i18n.T("Recuento:")+"\n  "+strings.Join(v.Recuento, "\n  ")))
	}

	parts := []string{head, score, table}
	parts = append(parts, extra...)
	parts = append(parts, m.help(v))
	if m.status != "" {
		parts = append(parts, styleErr.Render(m.status))
	}
	return lipgloss.JoinVertical(lipgloss.Left, parts...), w
}

func (m *Mesa) centerInfo(v game.View) string {
	var lines []string
	switch v.Phase {
	case game.PhaseMus:
		lines = append(lines, styleLance.Render("¿MUS?"))
	case game.PhaseDescarte:
		lines = append(lines, styleLance.Render("DESCARTE"))
	case game.PhaseApuesta:
		lines = append(lines, styleLance.Render(strings.ToUpper(i18n.Tf("A %s", v.Bet.Lance))))
		switch {
		case v.Bet.Open:
			lines = append(lines, i18n.T("sin envite"))
		case v.Bet.Ordago:
			lines = append(lines, styleErr.Render("¡ÓRDAGO!"))
		default:
			lines = append(lines, i18n.Tf("envite: %d", v.Bet.Amount), fichas(v.Bet.Amount, m.img, 18),
				styleDim.Render(i18n.Tf("si no se quiere: %d", v.Bet.Deje)))
		}
	case game.PhaseFinMano:
		lines = append(lines, styleLance.Render(i18n.T("FIN DE MANO")))
	case game.PhaseFinPartida:
		lines = append(lines, styleLance.Render(i18n.T("FIN DE PARTIDA")))
	}
	lines = append(lines, "")
	for _, r := range v.Results {
		lines = append(lines, styleDim.Render(resultado(r)))
	}
	return lipgloss.JoinVertical(lipgloss.Center, lines...)
}

func resultado(r game.LanceResult) string {
	switch r.Estado {
	case game.EnPaso:
		return i18n.Tf("%s: en paso", r.Lance)
	case game.Querido:
		return i18n.Tf("%s: %d queridos", r.Lance, r.Amount)
	case game.NoQuerido:
		return i18n.Tf("%s: no querido", r.Lance)
	case game.SinEnvite:
		return i18n.Tf("%s: sin envite", r.Lance)
	case game.NoJugado:
		return i18n.Tf("%s: nadie", r.Lance)
	case game.OrdagoQuerido:
		return i18n.Tf("%s: ¡órdago!", r.Lance)
	}
	return ""
}

func (m *Mesa) help(v game.View) string {
	var opts []string
	salir := styleDim.Render("[esc] " + i18n.T("abandonar"))
	switch {
	case m.confirm:
		return styleErr.Render(i18n.T("¿Abandonar la partida y volver al menú?")+" ") + key("s", i18n.T("sí")) + "  " + key("n", i18n.T("no"))
	case m.typing:
		return styleTurn.Render(i18n.T("¿Cuántos tantos?")+" ") + m.input + "_  " + key("enter", i18n.T("aceptar")) + "  " + key("esc", i18n.T("cancelar"))
	case m.menuSena:
		return m.ayudaSenas(v)
	case v.Phase == game.PhaseFinMano:
		opts = append(opts, key("enter", i18n.T("siguiente mano")))
	case v.Phase == game.PhaseFinPartida:
		w := i18n.T("¡Habéis ganado la partida!")
		if v.PartidaWinner != game.Team(m.human) {
			w = i18n.T("Habéis perdido la partida. Más se perdió en Cuba.")
		}
		return styleTurn.Render(w) + "  " + key("enter", i18n.T("volver al menú"))
	case v.Turn != m.human:
		txt := styleDim.Render(i18n.Tf("Habla %s...", v.Names[v.Turn]))
		if m.vel == ajustes.PasoAPaso {
			txt = styleTurn.Render(i18n.Tf("Le toca a %s.", v.Names[v.Turn])+" ") + key(i18n.T("espacio"), i18n.T("que hable"))
		}
		if m.puedeSenar(v) {
			txt += "  " + key("s", i18n.T("pasar seña"))
		}
		return txt + "   " + salir
	case v.Phase == game.PhaseMus:
		opts = append(opts, key("m", "Mus"), key("c", "Corto (no hay mus)"))
	case v.Phase == game.PhaseDescarte:
		opts = append(opts, key("1-4", i18n.T("marcar cartas")), key("t", i18n.T("todas")), key("enter", i18n.T("descartarse")))
	case v.Phase == game.PhaseApuesta:
		for _, a := range v.Legal {
			switch a {
			case game.ActPaso:
				opts = append(opts, key("p", "Paso"))
			case game.ActEnvido:
				if v.Bet.Open {
					opts = append(opts, key("e", "Envido"), key("n", "Envido N"))
				} else {
					opts = append(opts, key("e", i18n.T("Dos más")), key("n", i18n.T("N más")))
				}
			case game.ActQuiero:
				opts = append(opts, key("q", "Quiero"))
			case game.ActNoQuiero:
				opts = append(opts, key("x", "No quiero"))
			case game.ActOrdago:
				opts = append(opts, key("o", "Órdago"))
			}
		}
	}
	prefix := ""
	if v.Turn == m.human {
		prefix = styleTurn.Render(i18n.T("Te toca:") + " ")
	}
	if m.puedeSenar(v) {
		opts = append(opts, key("s", i18n.T("pasar seña")))
	}
	return prefix + strings.Join(opts, "  ") + "   " + salir
}

func (m *Mesa) puedeSenar(v game.View) bool {
	return m.modo == ajustes.SenasDeVerdad && v.SenasAbiertas
}

// ayudaSenas es el menú para pasar una seña: las que llevas, resaltadas.
func (m *Mesa) ayudaSenas(v game.View) string {
	llevo := v.Cfg.Senas(v.Hand)
	var opts []string
	for x := rules.SenaDosReyes; x <= rules.SenaCiego; x++ {
		nombre, _, _ := strings.Cut(x.Significado(), ":")
		txt := key(fmt.Sprint(int(x)+1), nombre)
		if !slices.Contains(llevo, x) {
			txt = styleDim.Render(fmt.Sprintf("[%d] %s", int(x)+1, nombre))
		}
		opts = append(opts, txt)
	}
	nota := i18n.T("No llevas nada que tenga seña: no hagas nada.")
	if len(llevo) > 0 {
		nota = i18n.T("Es un instante: tu compañero puede no verla y un rival puede pillarla.")
	}
	titulo := i18n.T("Seña a tu compañero:") + " "
	sangria := strings.Repeat(" ", lipgloss.Width(titulo))
	return styleTurn.Render(titulo) + strings.Join(opts[:4], "  ") + "\n" +
		sangria + strings.Join(opts[4:], "  ") + "\n" +
		sangria + styleDim.Render(nota+" ") + key("esc", i18n.T("cerrar"))
}

// tusSenas resume, en el modo de verdad, las señas que has pasado.
func (m *Mesa) tusSenas() string {
	if m.modo != ajustes.SenasDeVerdad || len(m.misSenas) == 0 {
		return ""
	}
	partner := m.g.Names[game.Partner(m.human)]
	var out []string
	for _, h := range m.misSenas {
		txt := styleSena.Render(h.sena.Significado())
		if h.vista {
			txt += styleTantos.Render(" ✓ " + i18n.Tf("%s la ha visto", partner))
		} else {
			txt += styleDim.Render(" " + i18n.T("(sin confirmar: repítela)"))
		}
		switch len(h.pillada) {
		case 0:
		case 1:
			txt += styleErr.Render(" · " + i18n.Tf("¡pillada por %s!", h.pillada[0]))
		default:
			txt += styleErr.Render(" · " + i18n.Tf("¡pillada por %s y %s!", strings.Join(h.pillada[:len(h.pillada)-1], ", "), h.pillada[len(h.pillada)-1]))
		}
		out = append(out, txt)
	}
	titulo := i18n.T("Tus señas:") + " "
	return styleDim.Render(titulo) + strings.Join(out, "\n"+strings.Repeat(" ", lipgloss.Width(titulo)))
}
