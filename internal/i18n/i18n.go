// Package i18n traduce los textos del juego. El castellano es el idioma de
// origen: cada texto se escribe en castellano en el código y se busca su
// traducción en el catálogo (las claves son el propio texto en castellano).
// Si falta una traducción se muestra en castellano.
//
// Las palabras del mus (mus, envido, grande, órdago…) no se traducen: así se
// cantan en la mesa.
package i18n

import (
	"fmt"
	"sync/atomic"
)

type Idioma int

const (
	ES Idioma = iota // castellano
	EN               // inglés
	CS               // checo
	EU               // euskera
	NIdiomas
)

var nombres = [NIdiomas]string{"Castellano", "English", "Čeština", "Euskara"}

// String es el nombre del idioma en ese mismo idioma.
func (i Idioma) String() string {
	if i < 0 || i >= NIdiomas {
		return nombres[ES]
	}
	return nombres[i]
}

func (i Idioma) Valido() bool { return i >= 0 && i < NIdiomas }

var actual atomic.Int32

// Poner cambia el idioma de la interfaz (los valores no válidos se ignoran).
func Poner(i Idioma) {
	if i.Valido() {
		actual.Store(int32(i))
	}
}

func Actual() Idioma { return Idioma(actual.Load()) }

// trad son las traducciones de un texto: inglés, checo y euskera.
type trad [NIdiomas - 1]string

var catalogo = map[string]trad{}

// registrar añade entradas al catálogo; lo usan los ficheros cat_*.go.
func registrar(m map[string]trad) {
	for k, v := range m {
		if _, ok := catalogo[k]; ok {
			panic("i18n: texto repetido en el catálogo: " + k)
		}
		catalogo[k] = v
	}
}

// En devuelve el texto en el idioma pedido.
func En(i Idioma, es string) string {
	if i == ES || !i.Valido() {
		return es
	}
	if t, ok := catalogo[es]; ok && t[i-1] != "" {
		return t[i-1]
	}
	return es
}

// T traduce un texto al idioma actual.
func T(es string) string { return En(Actual(), es) }

// Tf traduce un formato y lo rellena como fmt.Sprintf.
func Tf(es string, args ...any) string { return fmt.Sprintf(T(es), args...) }

// Falta devuelve los idiomas en los que el texto no está traducido.
func Falta(es string) []Idioma {
	t, ok := catalogo[es]
	var out []Idioma
	for i := EN; i < NIdiomas; i++ {
		if !ok || t[i-1] == "" {
			out = append(out, i)
		}
	}
	return out
}
