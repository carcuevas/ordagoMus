// Órdago: remake en terminal del clásico juego de mus.
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"ordagomus/internal/ai"
	"ordagomus/internal/game"
	"ordagomus/internal/kitty"
	"ordagomus/internal/rules"
	"ordagomus/internal/sonido"
	"ordagomus/internal/tui"
)

// version se fija al compilar con -ldflags "-X main.version=vX.Y.Z".
var version = "dev"

func main() {
	cfg := rules.DefaultConfig()
	semilla := flag.Uint64("semilla", 0, "semilla para repetir una partida (0 = aleatoria)")
	sim := flag.Int("sim", 0, "simular N partidas entre bots y mostrar estadísticas")
	verVersion := flag.Bool("version", false, "mostrar la versión y salir")
	flag.Parse()

	if *verVersion {
		fmt.Println("ordago", version)
		return
	}

	if *sim > 0 {
		seed := *semilla
		if seed == 0 {
			seed = uint64(time.Now().UnixNano())
		}
		simular(cfg, *sim, seed)
		return
	}

	defer sonido.Preparar()()

	opts := []tea.ProgramOption{tea.WithAltScreen()}
	if kitty.Supported() {
		// Las imágenes de kitty van ligadas a la pantalla en la que se envían,
		// así que entramos en la pantalla alternativa antes de mandar la baraja.
		fmt.Print("\x1b[?1049h")
		borrar := tui.PrepararGraficos(os.Stdout)
		defer func() {
			borrar()
			fmt.Print("\x1b[?1049l")
		}()
		opts = nil
	}
	p := tea.NewProgram(tui.NewApp(*semilla), opts...)
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func simular(cfg rules.Config, n int, seed uint64) {
	var wins [2]int
	manos, ordagos := 0, 0
	for i := 0; i < n; i++ {
		s := seed + uint64(i)
		g := game.New(cfg, [4]string{"N", "E", "S", "O"}, s)
		var bots [4]*ai.Bot
		for j := range bots {
			bots[j] = ai.New(s*4+uint64(j), ai.PerfilMedio())
		}
		for g.Phase() != game.PhaseFinPartida {
			if g.Phase() == game.PhaseFinMano {
				manos++
				for _, r := range g.View(0).Results {
					if r.Estado == game.OrdagoQuerido {
						ordagos++
					}
				}
				g.NextHand()
				continue
			}
			t := g.ToAct()
			if err := g.Apply(t, bots[t].Decide(g.View(t))); err != nil {
				fmt.Fprintln(os.Stderr, "error:", err)
				os.Exit(1)
			}
			g.Events()
		}
		wins[g.View(0).PartidaWinner]++
	}
	fmt.Printf("%d partidas, %d manos, %d órdagos queridos. Pareja 0: %d · Pareja 1: %d\n", n, manos, ordagos, wins[0], wins[1])
}
