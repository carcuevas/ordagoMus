package tui

import (
	"strings"
	"testing"

	"ordagomus/internal/ai"
	"ordagomus/internal/ajustes"
	"ordagomus/internal/game"
	"ordagomus/internal/personajes"
)

// Con «Te consulta», cuando el compañero quiere envidar la partida se para
// hasta que contestas; con «Decide solo» envida sin preguntar.
func TestConsultaAlCompanero(t *testing.T) {
	pj := personajes.Todos
	for _, solo := range []bool{false, true} {
		vistas := 0
		for seed := uint64(1); seed < 300 && vistas < 5; seed++ {
			aj := ajustes.Defecto()
			aj.Velocidad = ajustes.PasoAPaso
			aj.CompaneroSolo = solo
			m := NuevaMesa(aj, pj[1], [2]*personajes.Personaje{pj[2], pj[3]}, seed)
			yo := ai.New(seed, ai.PerfilMedio())
			for n := 0; n < 60 && m.g.Phase() != game.PhaseFinMano; n++ {
				if m.g.ToAct() == m.human {
					m.g.Apply(m.human, yo.Decide(m.g.View(m.human)))
					m.drain()
					continue
				}
				antes := m.g.ToAct()
				m.botStep()
				if m.propuesta == nil {
					continue
				}
				if solo {
					t.Fatal("con «Decide solo» no debe consultar")
				}
				if antes != game.Partner(m.human) || m.g.ToAct() != antes {
					t.Fatal("la consulta tiene que ser del compañero y con la partida parada")
				}
				if !m.esperandoConsulta() {
					t.Fatal("la ayuda no muestra la pregunta")
				}
				m.key("p") // si pasar no vale, se ignora y se contesta abajo
				if m.propuesta != nil {
					m.key("q")
				}
				if m.propuesta != nil || m.g.ToAct() == antes {
					t.Fatal("tras contestar, el compañero debe haber hablado")
				}
				vistas++
			}
		}
		if !solo && vistas == 0 {
			t.Fatal("no se vio ninguna consulta")
		}
	}
}

func (m *Mesa) esperandoConsulta() bool {
	return m.propuesta != nil && strings.Contains(m.help(m.g.View(m.human)), m.textoPropuesta(m.g.View(m.human)))
}
