package tui

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"strings"
	"time"
	"unicode"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"ordagomus/internal/ajustes"
	"ordagomus/internal/i18n"
	"ordagomus/internal/personajes"
)

type pantalla int

const (
	pMenu pantalla = iota
	pOpciones
	pCompanero
	pRivales
	pFichas
	pAyuda
	pMesa
)

// App es el programa completo: menús y mesa.
type App struct {
	aj      ajustes.Ajustes
	pant    pantalla
	cursor  int
	mesa    *Mesa
	width   int
	height  int
	editing bool
	idioma  int // pestaña de las instrucciones
	scroll  int
	buf     string
	aviso   string
	seed    uint64 // 0 = aleatoria
}

func NewApp(seed uint64) *App {
	aj := ajustes.Cargar()
	i18n.Poner(aj.Idioma)
	return &App{aj: aj, seed: seed, idioma: int(aj.Idioma)}
}

func (a *App) Init() tea.Cmd { return nil }

const logo = `  ___  ____  ____    _    ____  ___
 / _ \|  _ \|  _ \  / \  / ___|/ _ \
| | | | |_) | | | |/ _ \| |  _| | | |
| |_| |  _ <| |_| / ___ \ |_| | |_| |
 \___/|_| \_\____/_/   \_\____|\___/`

var menuItems = []string{"Jugar partida", "Elegir compañero", "Elegir rivales", "Opciones", "Los jugadores", "Cómo se juega", "Salir"}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if ws, ok := msg.(tea.WindowSizeMsg); ok {
		a.width, a.height = ws.Width, ws.Height
	}
	if k, ok := msg.(tea.KeyMsg); ok && k.String() == "ctrl+c" {
		return a, tea.Quit
	}
	if a.pant == pMesa {
		if _, ok := msg.(volverMenuMsg); ok {
			a.pant, a.mesa, a.cursor = pMenu, nil, 0
			return a, nil
		}
		return a, a.mesa.Update(msg)
	}
	k, ok := msg.(tea.KeyMsg)
	if !ok {
		return a, nil
	}
	switch a.pant {
	case pMenu:
		return a, a.keyMenu(k.String())
	case pOpciones:
		a.keyOpciones(k.String())
	case pCompanero, pRivales, pFichas:
		a.keyLista(k.String())
	case pAyuda:
		a.keyAyuda(k.String())
	}
	return a, nil
}

func (a *App) move(k string, n int) {
	switch k {
	case "up", "k":
		a.cursor = (a.cursor + n - 1) % n
	case "down", "j":
		a.cursor = (a.cursor + 1) % n
	}
}

func (a *App) ir(p pantalla) {
	a.pant, a.cursor, a.aviso = p, 0, ""
}

func (a *App) keyMenu(k string) tea.Cmd {
	a.move(k, len(menuItems))
	if k != "enter" {
		return nil
	}
	switch a.cursor {
	case 0:
		return a.jugar()
	case 1:
		a.ir(pCompanero)
	case 2:
		a.ir(pRivales)
	case 3:
		a.ir(pOpciones)
	case 4:
		a.ir(pFichas)
	case 5:
		a.ir(pAyuda)
		a.scroll = 0
	case 6:
		return tea.Quit
	}
	return nil
}

func (a *App) jugar() tea.Cmd {
	seed := a.seed
	if seed == 0 {
		seed = uint64(time.Now().UnixNano())
	}
	rng := rand.New(rand.NewPCG(seed, 3))
	comp := personajes.Buscar(a.aj.Companero)
	if comp == nil {
		comp = personajes.Todos[0]
	}
	usados := []string{comp.ID}
	var riv [2]*personajes.Personaje
	for i, id := range a.aj.Rivales {
		if p := personajes.Buscar(id); p != nil && !slices.Contains(usados, id) {
			riv[i] = p
			usados = append(usados, id)
		}
	}
	for i := range riv {
		for riv[i] == nil {
			p := personajes.Todos[rng.IntN(len(personajes.Todos))]
			if !slices.Contains(usados, p.ID) {
				riv[i] = p
				usados = append(usados, p.ID)
			}
		}
	}
	a.mesa = NuevaMesa(a.aj, comp, riv, seed)
	a.mesa.width, a.mesa.height = a.width, a.height
	a.pant = pMesa
	return a.mesa.Init()
}

// ───────────── opciones ─────────────

const nOpciones = 10

func ciclo(xs []int, cur, dir int) int {
	i := slices.Index(xs, cur)
	if i < 0 {
		return xs[0]
	}
	return xs[(i+dir+len(xs))%len(xs)]
}

func (a *App) keyOpciones(k string) {
	if a.editing {
		switch {
		case k == "enter":
			if n := ajustes.LimpiarNombre(a.buf); n != "" {
				a.aj.Nombre = n
			}
			a.editing = false
		case k == "esc":
			a.editing = false
		case k == "backspace":
			if r := []rune(a.buf); len(r) > 0 {
				a.buf = string(r[:len(r)-1])
			}
		case len([]rune(k)) == 1 && !unicode.IsControl([]rune(k)[0]) && len([]rune(a.buf)) < ajustes.MaxNombre:
			a.buf += k
		}
		return
	}
	a.move(k, nOpciones)
	dir := 0
	switch k {
	case "left", "h":
		dir = -1
	case "right", "l", "enter", " ":
		dir = 1
	case "esc", "q":
		if err := a.aj.Guardar(); err != nil {
			a.aviso = i18n.Tf("No se pudieron guardar las opciones: %s", err.Error())
		}
		a.pant, a.cursor = pMenu, 3
		return
	}
	if dir == 0 {
		return
	}
	r := &a.aj.Reglas
	switch a.cursor {
	case 0:
		if k == "enter" {
			a.editing, a.buf = true, a.aj.Nombre
		}
	case 1:
		r.OchoReyes = !r.OchoReyes
	case 2:
		r.Tantos = ciclo(ajustes.OpcTantos, r.Tantos, dir)
	case 3:
		r.JuegosPorVaca = ciclo(ajustes.OpcImpar, r.JuegosPorVaca, dir)
	case 4:
		r.Vacas = ciclo(ajustes.OpcImpar, r.Vacas, dir)
	case 5:
		a.aj.Velocidad = ajustes.Velocidad((int(a.aj.Velocidad) + dir + 4) % 4)
	case 6:
		a.aj.Imagenes = !a.aj.Imagenes
	case 7:
		a.aj.Senas = 1 - a.aj.Senas
	case 8:
		a.aj.Sonido = !a.aj.Sonido
	case 9:
		a.aj.Idioma = (a.aj.Idioma + i18n.Idioma(dir) + i18n.NIdiomas) % i18n.NIdiomas
		i18n.Poner(a.aj.Idioma)
		a.idioma = int(a.aj.Idioma)
	}
}

func (a *App) viewOpciones() string {
	r := a.aj.Reglas
	reyes := i18n.T("4 reyes")
	if r.OchoReyes {
		reyes = i18n.T("8 reyes (treses = reyes, doses = ases)")
	}
	nombre := a.aj.Nombre
	if a.editing {
		nombre = a.buf + "_"
	}
	rows := [][2]string{
		{i18n.T("Tu nombre"), nombre},
		{i18n.T("Reyes"), reyes},
		{i18n.T("Tantos por juego"), fmt.Sprint(r.Tantos)},
		{i18n.T("Juegos por vaca (al mejor de)"), fmt.Sprint(r.JuegosPorVaca)},
		{i18n.T("Vacas por partida (al mejor de)"), fmt.Sprint(r.Vacas)},
		{i18n.T("Velocidad del ordenador"), a.aj.Velocidad.String()},
		{i18n.T("Cartas"), cartasOpcion(a.aj.Imagenes)},
		{"Señas", a.aj.Senas.String()},
		{i18n.T("Sonido"), siNo(a.aj.Sonido)},
		{"Idioma / Language / Jazyk / Hizkuntza", a.aj.Idioma.String()},
	}
	// La fila del idioma va en todos los idiomas a la vez y puede sobresalir.
	ancho := 34
	for _, row := range rows[:len(rows)-1] {
		ancho = max(ancho, lipgloss.Width(row[0]))
	}
	var lines []string
	for i, row := range rows {
		cur := "  "
		st := lipgloss.NewStyle()
		if i == a.cursor {
			cur, st = "▶ ", styleSel
		}
		lines = append(lines, cur+st.Render(rellenar(row[0], ancho)+" ‹ "+row[1]+" ›"))
	}
	help := key("↑↓", i18n.T("elegir")) + "  " + key("←→", i18n.T("cambiar")) + "  " + key("enter", i18n.T("editar nombre")) + "  " + key("esc", i18n.T("guardar y volver"))
	nota := styleDim.Render(i18n.T("Paso a paso: el ordenador no habla hasta que pulses espacio."))
	switch {
	case a.cursor == 7 && a.aj.Senas == ajustes.SenasDeVerdad:
		nota = styleDim.Render(i18n.T(notaSenasDeVerdad))
	case a.cursor == 7:
		nota = styleDim.Render(i18n.T("Señas escritas: se pasan solas y se escribe lo que pasa cada uno."))
	}
	return styleTitle.Render(i18n.T("OPCIONES")) + "\n\n" + strings.Join(lines, "\n") + "\n\n" + nota + "\n\n" + help
}

const notaSenasDeVerdad = "Señas de verdad: los gestos duran un instante (de 0,5 a 1 segundo) y nadie te dice\n" +
	"lo que son: mira la cara de tu compañero. Tú pasas las tuyas con [s]; puede que no las\n" +
	"vea (te avisa cuando sí) y puede que un rival las pille (también te enteras)."

func cartasOpcion(img bool) string {
	switch {
	case !gfx.on:
		return i18n.T("Texto (tu terminal no admite imágenes; prueba kitty o Ghostty)")
	case img:
		return i18n.T("Ilustradas (Heraclio Fournier, 1878)")
	}
	return i18n.T("Texto")
}

// ───────────── listas de personajes ─────────────

// En la pantalla de rivales la primera fila es "Al azar".
func (a *App) listaOffset() int {
	if a.pant == pRivales {
		return 1
	}
	return 0
}

func (a *App) keyLista(k string) {
	n := len(personajes.Todos) + a.listaOffset()
	a.move(k, n)
	if k == "esc" || k == "q" {
		_ = a.aj.Guardar()
		back := map[pantalla]int{pCompanero: 1, pRivales: 2, pFichas: 4}[a.pant]
		a.pant, a.cursor = pMenu, back
		return
	}
	if k != "enter" && k != " " {
		return
	}
	switch a.pant {
	case pCompanero:
		p := personajes.Todos[a.cursor]
		a.aj.Companero = p.ID
		for i, id := range a.aj.Rivales {
			if id == p.ID {
				a.aj.Rivales[i] = ""
			}
		}
		_ = a.aj.Guardar()
		a.aviso = i18n.Tf("%s será tu compañero.", p.Completo())
	case pRivales:
		if a.cursor == 0 {
			a.aj.Rivales = [2]string{}
			a.aviso = i18n.T("Los rivales se sortearán.")
			return
		}
		p := personajes.Todos[a.cursor-1]
		switch {
		case p.ID == a.aj.Companero:
			a.aviso = i18n.Tf("%s es tu compañero; no puede ser rival.", p.Completo())
		case a.aj.Rivales[0] == p.ID:
			a.aj.Rivales[0] = ""
		case a.aj.Rivales[1] == p.ID:
			a.aj.Rivales[1] = ""
		case a.aj.Rivales[0] == "":
			a.aj.Rivales[0] = p.ID
		case a.aj.Rivales[1] == "":
			a.aj.Rivales[1] = p.ID
		default:
			a.aviso = i18n.T("Ya hay dos rivales: quita uno antes.")
		}
		_ = a.aj.Guardar()
	}
}

func barra(x float64) string {
	n := int(x*10 + 0.5)
	return styleTantos.Render(strings.Repeat("█", n)) + styleDim.Render(strings.Repeat("░", 10-n))
}

func ficha(p *personajes.Personaje, w int, img bool) string {
	pf := p.Perfil
	wrap := lipgloss.NewStyle().Width(w)
	stats := [][2]any{
		{"Grande", pf.Lance[0]}, {"Chica", pf.Lance[1]}, {"Pares", pf.Lance[2]},
		{"Juego", pf.Lance[3]}, {"Punto", pf.Lance[4]},
		{i18n.T("Faroles"), pf.Farol}, {i18n.T("Valentía"), pf.Valentia}, {"Órdagos", pf.Ordago},
		{i18n.T("Se da mus"), pf.Musero}, {i18n.T("Lee la mesa"), pf.Lectura}, {i18n.T("Acepta envites"), pf.Querer},
		{i18n.T("Pilla señas"), p.Vista * 2}, {i18n.T("Disimulo"), p.Disimulo},
	}
	ancho := 14
	for _, st := range stats {
		ancho = max(ancho, lipgloss.Width(st[0].(string)))
	}
	var sl []string
	for i := 0; i < len(stats); i += 2 {
		line := rellenar(stats[i][0].(string), ancho) + " " + barra(stats[i][1].(float64))
		if i+1 < len(stats) {
			line += "   " + rellenar(stats[i+1][0].(string), ancho) + " " + barra(stats[i+1][1].(float64))
		}
		sl = append(sl, line)
	}
	cabecera := lipgloss.JoinVertical(lipgloss.Left,
		styleLogo.Render(p.Completo()),
		styleDim.Render(i18n.Tf("%s · %d años · %s", p.Origen, p.Edad, i18n.T(p.Oficio))),
		"",
		styleFlavor.Render(lipgloss.NewStyle().Width(max(w-caraW-2, 20)).Render("“"+p.Frases[0]+"”")),
	)
	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Top, dibujarCara(p.ID, exNormal, nil, false, img), "  ", cabecera),
		"",
		wrap.Render(i18n.T(p.Historia)),
		"",
		styleFlavor.Render(wrap.Render(i18n.Tf("Estilo: %s", i18n.T(p.Estilo)))),
		"",
		strings.Join(sl, "\n"),
	)
}

func (a *App) viewLista() string {
	off := a.listaOffset()
	var rows []string
	if a.pant == pRivales {
		mark := "  "
		if a.aj.Rivales == [2]string{} {
			mark = "✓ "
		}
		rows = append(rows, a.fila(0, mark+i18n.T("Al azar")))
	}
	for i, p := range personajes.Todos {
		mark := "  "
		switch {
		case a.pant == pCompanero && p.ID == a.aj.Companero:
			mark = "✓ "
		case a.pant == pRivales && slices.Contains(a.aj.Rivales[:], p.ID):
			mark = "✓ "
		case a.pant == pRivales && p.ID == a.aj.Companero:
			mark = "♥ "
		}
		rows = append(rows, a.fila(i+off, fmt.Sprintf("%s%-24s %s", mark, p.Completo(), styleDim.Render(p.Origen))))
	}
	list := strings.Join(rows, "\n")

	var right string
	if i := a.cursor - off; i >= 0 {
		right = ficha(personajes.Todos[i], max(min(a.width-56, 70), 40), gfx.on && a.aj.Imagenes)
	} else {
		right = styleDim.Render(i18n.T("Sin elegir: se sortean entre los que queden libres.\nSi eliges solo uno, el otro se sortea."))
	}

	title := map[pantalla]string{
		pCompanero: i18n.T("ELIGE COMPAÑERO"),
		pRivales:   i18n.T("ELIGE RIVALES (hasta dos)"),
		pFichas:    i18n.T("LOS JUGADORES"),
	}[a.pant]
	moverse, atras := key("↑↓", i18n.T("moverse")), key("esc", i18n.T("volver"))
	help := moverse + "  " + atras
	switch a.pant {
	case pCompanero:
		help = moverse + "  " + key("enter", i18n.T("elegir")) + "  " + atras
	case pRivales:
		help = moverse + "  " + key("enter", i18n.T("marcar/quitar")) + "  " + atras
	}
	body := lipgloss.JoinHorizontal(lipgloss.Top, styleBox.Render(list), "  ", styleBox.Render(right))
	out := styleTitle.Render(title) + "\n\n" + body + "\n" + help
	if a.aviso != "" {
		out += "\n" + styleTurn.Render(a.aviso)
	}
	return out
}

func (a *App) fila(i int, txt string) string {
	if i == a.cursor {
		return styleSel.Render("▶ ") + styleSel.Render(txt)
	}
	return "  " + txt
}

// ───────────── vista ─────────────

func (a *App) View() string {
	switch a.pant {
	case pMesa:
		return a.mesa.View()
	case pOpciones:
		return a.viewOpciones()
	case pAyuda:
		return a.viewAyuda()
	case pCompanero, pRivales, pFichas:
		return a.viewLista()
	}
	return a.viewMenu()
}

func (a *App) viewMenu() string {
	comp := personajes.Buscar(a.aj.Companero)
	compName := "?"
	if comp != nil {
		compName = comp.Completo()
	}
	var riv []string
	for _, id := range a.aj.Rivales {
		if p := personajes.Buscar(id); p != nil {
			riv = append(riv, p.Completo())
		}
	}
	rivName := i18n.T("al azar")
	switch len(riv) {
	case 1:
		rivName = i18n.Tf("%s y otro al azar", riv[0])
	case 2:
		rivName = i18n.Tf("%s y %s", riv[0], riv[1])
	}
	extra := map[int]string{1: compName, 2: rivName}

	ancho := 18
	for _, it := range menuItems {
		ancho = max(ancho, lipgloss.Width(i18n.T(it)))
	}
	var rows []string
	for i, it := range menuItems {
		txt := i18n.T(it)
		if e, ok := extra[i]; ok {
			txt = rellenar(txt, ancho) + " " + styleDim.Render(e)
		}
		rows = append(rows, a.fila(i, txt))
	}
	r := a.aj.Reglas
	reyes := 4
	if r.OchoReyes {
		reyes = 8
	}
	resumen := styleDim.Render(i18n.Tf("%s · a %d tantos · %d reyes · velocidad %s · señas %s", a.aj.Nombre, r.Tantos, reyes, a.aj.Velocidad.String(), strings.ToLower(a.aj.Senas.String())))
	out := lipgloss.JoinVertical(lipgloss.Left,
		styleLogo.Render(logo),
		styleFlavor.Render("   "+i18n.T("El mus de toda la vida. Inspirado en el Órdago de DOS de Pedro W. Torrecilla.")),
		"",
		strings.Join(rows, "\n"),
		"",
		resumen,
		"",
		key("↑↓", i18n.T("moverse"))+"  "+key("enter", i18n.T("aceptar"))+"  "+key("ctrl+c", i18n.T("salir")),
	)
	if a.aviso != "" {
		out += "\n" + styleErr.Render(a.aviso)
	}
	return lipgloss.Place(max(a.width, 80), max(a.height, 24), lipgloss.Center, lipgloss.Center, out)
}

func siNo(b bool) string {
	if b {
		return i18n.T("Sí")
	}
	return i18n.T("No")
}

// rellenar completa s con espacios hasta w columnas (como %-Ns, pero
// contando el ancho en pantalla).
func rellenar(s string, w int) string {
	return s + strings.Repeat(" ", max(w-lipgloss.Width(s), 0))
}
