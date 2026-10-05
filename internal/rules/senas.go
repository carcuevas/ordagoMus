package rules

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

var senaSignificado = [...]string{
	SenaDosReyes:    "dos reyes",
	SenaDosAses:     "dos ases",
	SenaMediasReyes: "medias de reyes",
	SenaMediasAses:  "medias de ases",
	SenaDuples:      "duples",
	Sena31:          "la 31",
	SenaCiego:       "ciego: ni pares ni juego",
}

func (s Sena) Gesto() string       { return senaGesto[s] }
func (s Sena) Significado() string { return senaSignificado[s] }

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
	// No se puede pasar ciego con pares de cualquier clase ni con 29 o 30 al punto.
	if p.Kind == SinPares && pts < 29 {
		out = append(out, SenaCiego)
	}
	return out
}
