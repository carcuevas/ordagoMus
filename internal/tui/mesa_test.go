package tui

import (
	"regexp"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"

	"ordagomus/internal/ajustes"
	"ordagomus/internal/game"
	"ordagomus/internal/personajes"
)

// Un código de color que no empieza por ESC sale en pantalla como basura.
var codigoSuelto = regexp.MustCompile(`[^\x1b]\[[0-9;]*m|\x1b\x1b`)

func TestDescarteSeVeLimpio(t *testing.T) {
	lipgloss.SetColorProfile(termenv.ANSI256)
	aj := ajustes.Defecto()
	aj.Velocidad = ajustes.PasoAPaso
	pj := personajes.Todos
	m := NuevaMesa(aj, pj[1], [2]*personajes.Personaje{pj[2], pj[3]}, 3)
	for m.g.Phase() == game.PhaseMus {
		m.g.Apply(m.g.ToAct(), game.Action{Kind: game.ActMus})
	}
	for m.g.ToAct() != 0 {
		m.g.Apply(m.g.ToAct(), game.Action{Kind: game.ActDescarte, Discard: []int{0}})
	}
	m.drain()
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("1")})
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("3")})
	out := m.View()
	if !strings.Contains(out, "FUERA") {
		t.Fatal("no se marcan las cartas a descartar")
	}
	if loc := codigoSuelto.FindStringIndex(out); loc != nil {
		t.Fatalf("código de color suelto en pantalla: %q", out[max(loc[0]-20, 0):loc[1]+20])
	}
}
