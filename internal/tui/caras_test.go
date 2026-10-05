package tui

import (
	"bytes"
	"testing"

	"ordagomus/internal/personajes"
	"ordagomus/internal/rules"
)

func TestCaras(t *testing.T) {
	for _, p := range personajes.Todos {
		c, ok := caras[p.ID]
		if !ok {
			t.Errorf("%s no tiene cara", p.ID)
			continue
		}
		senas := []*rules.Sena{nil}
		for s := rules.SenaDosReyes; s <= rules.SenaCiego; s++ {
			senas = append(senas, &s)
		}
		for ex := exNormal; ex <= exTriste; ex++ {
			for _, s := range senas {
				for y, row := range c.pixeles(ex, s) {
					if len(row) != caraW {
						t.Fatalf("%s: la fila %d mide %d", p.ID, y, len(row))
					}
					for _, b := range row {
						if _, ok := c.color(b); !ok && b != '.' {
							t.Errorf("%s: color %q sin definir", p.ID, b)
						}
					}
				}
			}
		}
	}
}

// Cada seña y cada expresión tiene que notarse en la cara.
func TestGestosSeVen(t *testing.T) {
	for id, c := range caras {
		normal := c.pixeles(exNormal, nil)
		for ex := exHabla; ex <= exTriste; ex++ {
			if igual(c.pixeles(ex, nil), normal) {
				t.Errorf("%s: la expresión %d no cambia la cara", id, ex)
			}
		}
		for s := rules.SenaDosReyes; s <= rules.SenaCiego; s++ {
			if igual(c.pixeles(exNormal, &s), normal) {
				t.Errorf("%s: la seña %q no se ve", id, s.Gesto())
			}
		}
	}
}

func igual(a, b [caraH][]byte) bool {
	return bytes.Equal(bytes.Join(a[:], nil), bytes.Join(b[:], nil))
}
