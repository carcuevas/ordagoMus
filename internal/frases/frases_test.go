package frases

import (
	"testing"

	"ordagomus/internal/i18n"
)

// Lo que se dice igual en todos los idiomas y ya registra otro catálogo.
var sinCatalogo = map[string]bool{"¡Órdago!": true}

func TestTraducidas(t *testing.T) {
	listas := [][]string{ganamos, perdemos}
	for _, xs := range porAccion {
		listas = append(listas, xs)
	}
	for _, xs := range listas {
		for _, s := range xs {
			if sinCatalogo[s] {
				continue
			}
			if f := i18n.Falta(s); len(f) > 0 {
				t.Errorf("falta traducir %q a %v", s, f)
			}
		}
	}
}
