package game_test

import (
	"slices"
	"testing"

	"ordagomus/internal/ai"
	"ordagomus/internal/cards"
	"ordagomus/internal/game"
	"ordagomus/internal/rules"
)

// Juega partidas completas entre bots y comprueba que el motor no se atasca,
// que nadie se queda con cartas repetidas y que siempre hay un ganador.
func TestPartidasEntreBots(t *testing.T) {
	names := [4]string{"A", "B", "C", "D"}
	for seed := uint64(1); seed <= 40; seed++ {
		cfg := rules.DefaultConfig()
		cfg.OchoReyes = seed%2 == 0
		g := game.New(cfg, names, seed)
		var bots [4]*ai.Bot
		for i := range bots {
			bots[i] = ai.New(seed*10+uint64(i), ai.PerfilMedio())
		}
		for steps := 0; g.Phase() != game.PhaseFinPartida; steps++ {
			if steps > 20000 {
				t.Fatalf("seed %d: la partida no termina", seed)
			}
			if g.Phase() == game.PhaseFinMano {
				g.NextHand()
				continue
			}
			s := g.ToAct()
			v := g.View(s)
			checkHands(t, v)
			if err := g.Apply(s, bots[s].Decide(v)); err != nil {
				t.Fatalf("seed %d: %v", seed, err)
			}
			g.Events()
		}
		v := g.View(0)
		if v.PartidaWinner < 0 {
			t.Fatalf("seed %d: partida sin ganador", seed)
		}
	}
}

func checkHands(t *testing.T, v game.View) {
	t.Helper()
	seen := map[cards.Card]bool{}
	for _, c := range v.Hand {
		if seen[c] {
			t.Fatalf("carta repetida en la mano: %v", v.Hand)
		}
		seen[c] = true
	}
}

func TestNoQueridoYRecuento(t *testing.T) {
	g := game.New(rules.DefaultConfig(), [4]string{"A", "B", "C", "D"}, 7)
	m := g.Mano()
	// Mano corta y envida a grande; los dos rivales no quieren: 1 tanto de deje.
	must(t, g.Apply(m, game.Action{Kind: game.ActCorto}))
	must(t, g.Apply(m, game.Action{Kind: game.ActEnvido, Amount: 2}))
	must(t, g.Apply((m+1)%4, game.Action{Kind: game.ActNoQuiero}))
	must(t, g.Apply((m+3)%4, game.Action{Kind: game.ActNoQuiero}))
	if got := g.View(0).Score[game.Team(m)]; got != 1 {
		t.Fatalf("deje de grande: want 1 got %d", got)
	}
	// Reenvite no querido: se cobra lo envidado antes (2), no el 1.
	must(t, g.Apply(m, game.Action{Kind: game.ActEnvido, Amount: 2}))
	must(t, g.Apply((m+1)%4, game.Action{Kind: game.ActEnvido, Amount: 5}))
	must(t, g.Apply((m+2)%4, game.Action{Kind: game.ActNoQuiero})) // habla primero el siguiente
	must(t, g.Apply(m, game.Action{Kind: game.ActNoQuiero}))
	if got := g.View(0).Score[game.Team(m+1)]; got != 2 {
		t.Fatalf("deje de chica tras reenvite: want 2 got %d", got)
	}
}

func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

// Señas a mano: solo las verdaderas, el compañero confirma que las ve y si un
// rival la pilla se entera la pareja.
func TestSenasManuales(t *testing.T) {
	vistas, pilladas := 0, 0
	for seed := uint64(1); seed <= 300; seed++ {
		g := game.New(rules.DefaultConfig(), [4]string{"A", "B", "C", "D"}, seed)
		g.Manual[0] = true
		g.Atento[2] = 0.8
		g.Vista[1], g.Vista[3] = 1, 1
		if g.SenasAbiertas() {
			t.Fatal("en la primera mano no hay señas")
		}
		if err := g.HacerSena(0, rules.SenaCiego); err != game.ErrSinSenas {
			t.Fatalf("seña antes de cortar: %v", err)
		}
		must(t, g.Apply(g.Mano(), game.Action{Kind: game.ActCorto}))
		g.Events()
		cfg := g.Cfg
		mias := cfg.Senas(g.View(0).Hand)
		for x := rules.SenaDosReyes; x <= rules.SenaCiego; x++ {
			err := g.HacerSena(0, x)
			lleva := false
			for _, m := range mias {
				lleva = lleva || m == x
			}
			if lleva != (err == nil) {
				t.Fatalf("seed %d: seña %v con %v: %v", seed, x, mias, err)
			}
		}
		if g.HacerSena(2, rules.SenaCiego) != game.ErrSinSenas {
			t.Fatal("un asiento automático no pasa señas a mano")
		}
		vio := false
		for _, e := range g.Events() {
			switch e.Kind {
			case game.EvSenaVista:
				if e.To != 0 || e.Seat != 2 {
					t.Fatalf("confirmación mal dirigida: %+v", e)
				}
				vio = true
				vistas++
			case game.EvSenaPillada:
				if e.To != 0 && e.To != 2 {
					t.Fatalf("aviso de seña pillada a quien no toca: %+v", e)
				}
				pilladas++
			}
		}
		v := g.View(2)
		if got := v.Senas[0]; vio != (len(got) > 0) || v.SenaCompleta[0] {
			t.Fatalf("seed %d: el compañero sabe %v (vio %v, completa %v)", seed, got, vio, v.SenaCompleta[0])
		}
		if r := g.View(1); r.SenaCompleta[0] {
			t.Fatal("una seña pillada a mano no es completa")
		}
	}
	if vistas == 0 || pilladas == 0 {
		t.Fatalf("vistas %d, pilladas %d", vistas, pilladas)
	}
}

// Con Manda, cuando envidan los rivales contesta primero el que manda y su
// «no quiero» vale por la pareja.
func TestMandaContesta(t *testing.T) {
	probados := 0
	for seed := uint64(1); seed < 400 && probados < 20; seed++ {
		g := game.New(rules.DefaultConfig(), [4]string{"A", "B", "C", "D"}, seed)
		g.Manda[0] = true
		for g.Phase() == game.PhaseMus {
			g.Apply(g.ToAct(), game.Action{Kind: game.ActCorto})
		}
		// Busca un lance en el que envide un rival con los dos de la pareja 0 en juego.
		for g.Phase() == game.PhaseApuesta {
			s := g.ToAct()
			v := g.View(s)
			if game.Team(s) == 1 && slices.Contains(v.Legal, game.ActEnvido) && slices.Contains(v.Legal, game.ActPaso) {
				g.Apply(s, game.Action{Kind: game.ActEnvido, Amount: 2})
				if g.Phase() == game.PhaseApuesta && slices.Contains(g.View(g.ToAct()).Legal, game.ActQuiero) {
					if g.ToAct() != 0 {
						// solo contesta el compañero si el que manda no juega el lance
						break
					}
					g.Apply(0, game.Action{Kind: game.ActNoQuiero})
					if g.Phase() == game.PhaseApuesta && slices.Contains(g.View(g.ToAct()).Legal, game.ActQuiero) && g.ToAct() == 2 {
						t.Fatalf("semilla %d: tras el «no quiero» del que manda contesta el compañero", seed)
					}
					probados++
				}
				break
			}
			g.Apply(s, game.Action{Kind: game.ActPaso})
		}
	}
	if probados == 0 {
		t.Fatal("no se probó ningún envite")
	}
}
