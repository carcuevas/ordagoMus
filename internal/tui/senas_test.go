package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/muesli/termenv"

	"ordagomus/internal/ajustes"
	"ordagomus/internal/game"
	"ordagomus/internal/personajes"
	"ordagomus/internal/rules"
)

func tecla(m *Mesa, k string) {
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(k)})
}

// mesaDeVerdad monta una mesa con señas de verdad y juega hasta que se corta.
func mesaDeVerdad(t *testing.T, seed uint64, img bool) *Mesa {
	t.Helper()
	lipgloss.SetColorProfile(termenv.ANSI256)
	aj := ajustes.Defecto()
	aj.Velocidad = ajustes.PasoAPaso
	aj.Senas = ajustes.SenasDeVerdad
	pj := personajes.Todos
	gfx.on = img
	t.Cleanup(func() { gfx.on = false })
	m := NuevaMesa(aj, pj[1], [2]*personajes.Personaje{pj[6], pj[3]}, seed)
	m.width, m.height = 110, 50
	m.g.Apply(m.g.ToAct(), game.Action{Kind: game.ActCorto})
	m.drain()
	return m
}

func TestSenasDeVerdadEnLaMesa(t *testing.T) {
	vistas := 0
	for seed := uint64(1); seed < 60; seed++ {
		m := mesaDeVerdad(t, seed, false)
		v := m.g.View(m.human)
		if !v.SenasAbiertas {
			t.Fatal("después de cortar se pasan señas")
		}
		if strings.Contains(strings.Join(m.log, "\n"), "(seña)") {
			t.Fatal("de verdad no se escribe lo que pasan los demás")
		}
		tecla(m, "s")
		if !m.menuSena || !strings.Contains(m.View(), "Seña a tu compañero") {
			t.Fatal("no se abre el menú de señas")
		}
		llevo := v.Cfg.Senas(v.Hand)
		if len(llevo) == 0 {
			continue
		}
		tecla(m, string(rune('1'+int(llevo[0]))))
		out := ansi.Strip(m.View())
		if !strings.Contains(out, "Tus señas: "+llevo[0].Significado()) {
			t.Fatalf("no sale la seña que has pasado:\n%s", out)
		}
		if m.misSenas[0].vista {
			vistas++
			if !strings.Contains(out, "la ha visto") {
				t.Fatal("falta la confirmación del compañero")
			}
		}
		anchoMax(t, out)
	}
	if vistas == 0 {
		t.Fatal("el compañero no ve nunca las señas")
	}
}

// Los gestos de verdad duran entre 100 y 1000 ms y empiezan con una espera.
func TestGestosRapidos(t *testing.T) {
	m := mesaDeVerdad(t, 1, false)
	m.limpiarGestos()
	m.encolarGestos(2, []rules.Sena{rules.SenaDuples, rules.Sena31})
	g := m.gestos[2]
	if len(g) != 4 || g[0].sena != nil || g[2].sena != nil {
		t.Fatalf("cola de gestos: %+v", g)
	}
	for _, x := range []gesto{g[1], g[3]} {
		if x.dur < 100e6 || x.dur > 1000e6 {
			t.Fatalf("gesto de %v", x.dur)
		}
	}
}

func TestMarcadorConAmarracos(t *testing.T) {
	for _, img := range []bool{false, true} {
		m := mesaDeVerdad(t, 5, img)
		anchoMax(t, m.View())
		f := ansi.Strip(fichas(17, false, 60))
		if strings.Count(f, "▄██▄") != 3 || strings.Count(f, "●") != 2 {
			t.Fatalf("17 tantos deberían ser 3 amarracos y 2 piedras:\n%s", f)
		}
	}
}

// anchoMax falla si alguna línea no cabe en 110 columnas y dice cuál.
func anchoMax(t *testing.T, out string) {
	t.Helper()
	peor, linea := 0, ""
	for _, l := range strings.Split(ansi.Strip(out), "\n") {
		if w := ansi.StringWidth(strings.TrimRight(l, " ")); w > peor {
			peor, linea = w, l
		}
	}
	if peor > 110 {
		t.Fatalf("línea de %d columnas: %q", peor, strings.TrimRight(linea, " "))
	}
}
