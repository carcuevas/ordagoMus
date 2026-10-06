package sonido

import (
	"encoding/binary"
	"testing"
)

func TestWAV(t *testing.T) {
	for s := Sonido(0); s < nSonidos; s++ {
		x := sintetizar(s)
		if len(x) == 0 {
			t.Fatalf("%s: sin muestras", nombres[s])
		}
		w := wav(x)
		if string(w[:4]) != "RIFF" || string(w[8:16]) != "WAVEfmt " || string(w[36:40]) != "data" {
			t.Fatalf("%s: cabecera mal", nombres[s])
		}
		if n := binary.LittleEndian.Uint32(w[40:44]); int(n) != len(x)*2 || len(w) != 44+len(x)*2 {
			t.Fatalf("%s: tamaño %d", nombres[s], n)
		}
	}
}

// Sin Preparar, Tocar no hace nada (los tests de la interfaz no suenan).
func TestTocarSinPreparar(t *testing.T) {
	Tocar(Reparto)
	Tocar(-1)
	Tocar(nSonidos)
}
