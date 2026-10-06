package rules

import "ordagomus/internal/i18n"

// Sena es una de las señas admitidas en Euskadi (reglamento de Bizkaia, art. II).
type Sena int

const (
	SenaDosReyes Sena = iota
	SenaDosAses
	SenaMediasReyes
	SenaMediasAses
	SenaDuples
	Sena31
	SenaCiego
)

var senaGesto = [...]string{
	SenaDosReyes:    "se muerde el labio inferior",
	SenaDosAses:     "saca la punta de la lengua",
	SenaMediasReyes: "tuerce la comisura de los labios",
	SenaMediasAses:  "saca la lengua hacia un lado",
	SenaDuples:      "levanta las cejas",
	Sena31:          "guiña un ojo",
	SenaCiego:       "cierra los dos ojos",
}

// senaGestoTu es el gesto en segunda persona, para quien lo hace.
var senaGestoTu = [...]string{
	SenaDosReyes:    "te muerdes el labio inferior",
	SenaDosAses:     "sacas la punta de la lengua",
	SenaMediasReyes: "tuerces la comisura de los labios",
	SenaMediasAses:  "sacas la lengua hacia un lado",
	SenaDuples:      "levantas las cejas",
	Sena31:          "guiñas un ojo",
	SenaCiego:       "cierras los dos ojos",
}

var senaSignificado = [...]string{
	SenaDosReyes:    "dos reyes",
	SenaDosAses:     "dos ases",
	SenaMediasReyes: "medias de reyes",
	SenaMediasAses:  "medias de ases",
	SenaDuples:      "duples",
	Sena31:          "la 31",
	SenaCiego:       "ciego: sin juego y sin pares que valgan",
}

// Gesto es el gesto en tercera persona: "se muerde el labio inferior".
func (s Sena) Gesto() string { return i18n.T(senaGesto[s]) }

// GestoTu es el gesto dicho a quien lo hace: "te muerdes el labio inferior".
func (s Sena) GestoTu() string { return i18n.T(senaGestoTu[s]) }

// Significado es lo que quiere decir la seña. En el ciego va "ciego: …"
// (la interfaz se queda con lo que va antes de los dos puntos).
func (s Sena) Significado() string { return i18n.T(senaSignificado[s]) }

// Senas devuelve la seña completa que corresponde a la mano. Es obligatorio
// pasarla entera y nunca se puede mentir.
func (c Config) Senas(h Hand) []Sena {
	var out []Sena
	p := c.Pares(h)
	switch {
	case p.Kind == Duples:
		out = append(out, SenaDuples)
	case p.Kind == Medias && p.High == 12:
		out = append(out, SenaMediasReyes)
	case p.Kind == Medias && p.High == 1:
		out = append(out, SenaMediasAses)
	case p.Kind == Pareja && p.High == 12:
		out = append(out, SenaDosReyes)
	case p.Kind == Pareja && p.High == 1:
		out = append(out, SenaDosAses)
	}
	pts := c.Points(h)
	if pts == 31 {
		out = append(out, Sena31)
	}
	// Ciego: sin juego y sin nada que señar a pares. Una pareja baja (de seises,
	// de cuatros...) no tiene seña propia y no vale casi nada, así que también
	// se va ciego con ella.
	paresBajos := p.Kind == SinPares || (p.Kind == Pareja && p.High != 12 && p.High != 1)
	if paresBajos && pts < 31 {
		out = append(out, SenaCiego)
	}
	return out
}
