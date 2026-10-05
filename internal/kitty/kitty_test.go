package kitty

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

// La maquetación depende de que cada marcador ocupe exactamente una columna.
func TestAnchoMarcadores(t *testing.T) {
	for _, l := range Lines(0x2a0001, 1, 12, 9) {
		if w := ansi.StringWidth(l); w != 12 {
			t.Fatalf("ansi: ancho %d, want 12", w)
		}
		if w := lipgloss.Width(l); w != 12 {
			t.Fatalf("lipgloss: ancho %d, want 12", w)
		}
	}
}
