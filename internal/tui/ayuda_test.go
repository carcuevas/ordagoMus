package tui

import (
	"strings"
	"testing"
)

func TestAyudaIdiomas(t *testing.T) {
	a := &App{width: 110, height: 40, pant: pAyuda}
	for i, id := range idiomas {
		a.keyAyuda(string(rune('1' + i)))
		a.keyAyuda("pgdown")
		a.keyAyuda("pgdown")
		a.keyAyuda("pgdown")
		a.keyAyuda("pgdown")
		if out := a.viewAyuda(); !strings.Contains(out, id.titulo) {
			t.Fatalf("%s: no sale el título", id.nombre)
		}
		todo := strings.Join(lineasAyuda(id, 100), "\n")
		// Las palabras del mus no se traducen.
		for _, w := range []string{`"mus"`, `"envido"`, `"grande"`, `"chica"`, `"pares"`, `"juego"`, `"punto"`, `"órdago"`, `"quiero"`, `"no quiero"`, `"paso"`} {
			if !strings.Contains(todo, w) {
				t.Errorf("%s: falta %s", id.nombre, w)
			}
		}
	}
	a.keyAyuda("esc")
	if a.pant != pMenu {
		t.Fatal("esc no vuelve al menú")
	}
}
