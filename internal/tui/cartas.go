package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"ordagomus/internal/cards"
	"ordagomus/internal/i18n"
	"ordagomus/internal/kitty"
	"ordagomus/internal/rules"
)

// Las cartas se dibujan línea a línea con un único estilo por tramo: anidar
// estilos dentro de otros (p. ej. tachado) rompe los códigos de color.

var (
	suitColor = map[cards.Suit]lipgloss.Color{
		cards.Oros: "136", cards.Copas: "160", cards.Espadas: "25", cards.Bastos: "28",
	}
	suitBorder = map[cards.Suit]lipgloss.Color{
		cards.Oros: "220", cards.Copas: "203", cards.Espadas: "75", cards.Bastos: "114",
	}
	suitSymbol = map[cards.Suit]string{
		cards.Oros: "◉", cards.Copas: "∪", cards.Espadas: "†", cards.Bastos: "♣",
	}
	// Los textos de las cartas tienen que caber en la carta (bigW y compactW),
	// por eso van aquí, a medida, y no en el catálogo.
	rankName = [i18n.NIdiomas]map[int]string{
		i18n.ES: {1: "AS", 2: "DOS", 3: "TRES", 4: "CUATRO", 5: "CINCO", 6: "SEIS", 7: "SIETE",
			10: "SOTA", 11: "CABALLO", 12: "REY"},
		i18n.EN: {1: "ACE", 2: "TWO", 3: "THREE", 4: "FOUR", 5: "FIVE", 6: "SIX", 7: "SEVEN",
			10: "JACK", 11: "KNIGHT", 12: "KING"},
		i18n.CS: {1: "ESO", 2: "DVOJKA", 3: "TROJKA", 4: "ČTYŘKA", 5: "PĚTKA", 6: "ŠESTKA", 7: "SEDMIČKA",
			10: "SPODEK", 11: "JEZDEC", 12: "KRÁL"},
		i18n.EU: {1: "BATEKOA", 2: "BIKOA", 3: "HIRUKOA", 4: "LAUKOA", 5: "BOSTEKOA", 6: "SEIKOA", 7: "ZAZPIKOA",
			10: "SOTA", 11: "ZALDUNA", 12: "ERREGEA"},
	}
	compactRank = [i18n.NIdiomas]map[int]string{
		i18n.ES: {1: "As", 10: "Sot", 11: "Cab", 12: "Rey"},
		i18n.EN: {1: "A", 10: "J", 11: "Kn", 12: "K"},
		i18n.CS: {1: "Eso", 10: "Sp", 11: "Jez", 12: "Kr"},
		i18n.EU: {1: "Bat", 10: "Sot", 11: "Zal", 12: "Err"},
	}
	// valeComo es lo que vale el tres o el dos a 8 reyes: en la carta grande y
	// debajo de la ilustrada.
	valeComo = [i18n.NIdiomas][2]map[int]string{
		i18n.ES: {{3: "= REY", 2: "= AS"}, {3: "=rey", 2: "=as"}},
		i18n.EN: {{3: "= KING", 2: "= ACE"}, {3: "=king", 2: "=ace"}},
		i18n.CS: {{3: "= KRÁL", 2: "= ESO"}, {3: "=král", 2: "=eso"}},
		i18n.EU: {{3: "= ERREGE", 2: "= BATEKO"}, {3: "=errege", 2: "=bateko"}},
	}
)

// recortar deja s en w caracteres como mucho.
func recortar(s string, w int) string {
	if r := []rune(s); len(r) > w {
		return string(r[:w])
	}
	return s
}

const (
	bigW     = 9 // ancho interior de la carta grande
	compactW = 5
)

func center(s string, w int) string {
	pad := max(w-lipgloss.Width(s), 0)
	l := pad / 2
	return strings.Repeat(" ", l) + s + strings.Repeat(" ", pad-l)
}

func spread(left, right string, w int) string {
	pad := max(w-lipgloss.Width(left)-lipgloss.Width(right), 1)
	return left + strings.Repeat(" ", pad) + right
}

func boxed(lines []string, inner, border lipgloss.Style, w int) string {
	out := []string{border.Render("╭" + strings.Repeat("─", w) + "╮")}
	for _, l := range lines {
		out = append(out, border.Render("│")+inner.Render(l)+border.Render("│"))
	}
	out = append(out, border.Render("╰"+strings.Repeat("─", w)+"╯"))
	return strings.Join(out, "\n")
}

// bigCard es la carta del jugador: grande, con fondo claro y el valor efectivo
// (a 8 reyes, el tres es rey y el dos es as).
func bigCard(c cards.Card, cfg rules.Config, selected bool) string {
	sym := suitSymbol[c.Suit]
	num := fmt.Sprint(c.Rank)
	middle := sym + " " + sym + " " + sym
	if cfg.OchoReyes && (c.Rank == 3 || c.Rank == 2) {
		middle = valeComo[i18n.Actual()][0][c.Rank]
	}
	lines := []string{
		spread(num, sym, bigW),
		center(recortar(rankName[i18n.Actual()][c.Rank], bigW), bigW),
		center(middle, bigW),
		center(recortar(strings.ToUpper(c.Suit.String()), bigW), bigW),
		spread(sym, num, bigW),
	}
	bg := lipgloss.Color("231")
	border := lipgloss.NewStyle().Foreground(suitBorder[c.Suit])
	if selected {
		bg = lipgloss.Color("224")
		border = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	}
	inner := lipgloss.NewStyle().Background(bg).Foreground(suitColor[c.Suit]).Bold(true)
	return boxed(lines, inner, border, bigW)
}

func compactCard(c cards.Card, hidden bool) string {
	if hidden {
		back := lipgloss.NewStyle().Background(lipgloss.Color("52")).Foreground(lipgloss.Color("124"))
		return boxed([]string{"▚▚▚▚▚", "▚▚▚▚▚"}, back, lipgloss.NewStyle().Foreground(lipgloss.Color("124")), compactW)
	}
	r, ok := compactRank[i18n.Actual()][c.Rank]
	if !ok {
		r = fmt.Sprint(c.Rank)
	}
	suit := recortar(c.Suit.String(), compactW)
	lines := []string{spread(r, suitSymbol[c.Suit], compactW), center(suit, compactW)}
	inner := lipgloss.NewStyle().Background(lipgloss.Color("231")).Foreground(suitColor[c.Suit]).Bold(true)
	return boxed(lines, inner, lipgloss.NewStyle().Foreground(suitBorder[c.Suit]), compactW)
}

func imageCard(id uint32, pid uint32, cols, rows int) string {
	return strings.Join(kitty.Lines(id, pid, cols, rows), "\n")
}

// bigHand dibuja la mano del jugador con las marcas de descarte y los números.
// Con img usa la baraja ilustrada.
func bigHand(h rules.Hand, cfg rules.Config, sel [4]bool, showIdx, img bool) string {
	cols := make([]string, 4)
	w := bigW + 2
	t := gfx.cartaBig()
	if img {
		w = t.c + 1
	}
	for i, c := range h {
		mark := strings.Repeat(" ", w)
		if sel[i] {
			mark = styleErr.Bold(true).Render(center(i18n.T("✗ FUERA"), w))
		}
		card := bigCard(c, cfg, sel[i])
		if img {
			card = imageCard(cardID(c), pidBig(gfx.tb), t.c, t.r)
		}
		parts := []string{mark, card}
		var label []string
		if showIdx {
			label = append(label, fmt.Sprintf("[%d]", i+1))
		}
		if img && cfg.OchoReyes && (c.Rank == 3 || c.Rank == 2) {
			label = append(label, valeComo[i18n.Actual()][1][c.Rank])
		}
		if len(label) > 0 {
			parts = append(parts, styleKeys.Render(center(strings.Join(label, " "), w)))
		} else if img {
			parts = append(parts, "")
		}
		cols[i] = lipgloss.JoinVertical(lipgloss.Left, parts...)
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, cols...)
}

func compactHand(h rules.Hand, hidden, img bool) string {
	cols := make([]string, 4)
	t := gfx.cartaSmall()
	for i, c := range h {
		switch {
		case img && hidden:
			cols[i] = imageCard(idReverso, pidSmall(gfx.ts), t.c, t.r)
		case img:
			cols[i] = imageCard(cardID(c), pidSmall(gfx.ts), t.c, t.r)
		default:
			cols[i] = compactCard(c, hidden)
		}
		if img && i < 3 {
			cols[i] = lipgloss.JoinHorizontal(lipgloss.Top, cols[i], " ")
		}
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, cols...)
}
