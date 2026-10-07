// Package ajustes guarda las opciones del jugador entre partidas.
package ajustes

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
	"unicode"

	"ordagomus/internal/canal"
	"ordagomus/internal/i18n"
	"ordagomus/internal/rules"
)

type Velocidad int

const (
	PasoAPaso Velocidad = iota
	Lenta
	Normal
	Rapida
)

var velocidadNombres = [...]string{"Paso a paso", "Lenta", "Normal", "Rápida"}

func (v Velocidad) String() string { return i18n.T(velocidadNombres[v]) }

// Pausa entre jugadas del ordenador (0 en paso a paso: espera a una tecla).
func (v Velocidad) Pausa() time.Duration {
	switch v {
	case Lenta:
		return 2500 * time.Millisecond
	case Normal:
		return 1600 * time.Millisecond
	case Rapida:
		return 700 * time.Millisecond
	}
	return 0
}

// ModoSenas: cómo se pasan las señas en la mesa.
type ModoSenas int

const (
	// SenasEscritas: se pasan solas y se escribe lo que pasa cada uno.
	SenasEscritas ModoSenas = iota
	// SenasDeVerdad: los gestos duran un instante en la cara, nadie te dice lo
	// que significan y tú pasas las tuyas a mano.
	SenasDeVerdad
)

var modoNombres = [...]string{"Escritas", "De verdad"}

func (m ModoSenas) String() string { return i18n.T(modoNombres[m]) }

type Ajustes struct {
	Nombre    string
	Reglas    rules.Config
	Velocidad Velocidad
	Imagenes  bool      // baraja ilustrada si el terminal la admite
	Companero string    // ID del personaje
	Rivales   [2]string // "" = al azar
	Senas     ModoSenas
	Sonido    bool // efectos al repartir, cortar, envidar, contar...
	Idioma    i18n.Idioma
}

func Defecto() Ajustes {
	return Ajustes{Nombre: "Tú", Reglas: rules.DefaultConfig(), Velocidad: Normal, Imagenes: true, Companero: "remedios", Sonido: true}
}

func ruta() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "ordago", "ajustes.json"), nil
}

// Opciones que se pueden elegir en el menú; lo que venga del fichero se
// ajusta a ellas.
var (
	OpcTantos = []int{20, 30, 40, 50}
	OpcImpar  = []int{1, 3, 5}
)

// MaxNombre es la longitud máxima del nombre del jugador, en caracteres.
const MaxNombre = 16

// maxFichero acota lo que se lee del fichero de ajustes.
const maxFichero = 64 << 10

// LimpiarNombre quita los caracteres de control (que podrían colar secuencias
// de escape en el terminal) y recorta el nombre a MaxNombre caracteres.
func LimpiarNombre(s string) string {
	s = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) || r == unicode.ReplacementChar {
			return -1
		}
		return r
	}, s)
	s = strings.TrimSpace(s)
	if r := []rune(s); len(r) > MaxNombre {
		s = strings.TrimSpace(string(r[:MaxNombre]))
	}
	return s
}

// Cargar devuelve los ajustes guardados o los de por defecto.
func Cargar() Ajustes {
	a := Defecto()
	p, err := ruta()
	if err != nil {
		return a
	}
	f, err := os.Open(p)
	if err != nil {
		return a
	}
	defer f.Close()
	data, err := io.ReadAll(io.LimitReader(f, maxFichero+1))
	if err != nil || len(data) > maxFichero {
		return a
	}
	if json.Unmarshal(data, &a) != nil {
		return Defecto()
	}
	a.validar()
	return a
}

// validar corrige los valores que no podrían salir del menú de opciones: el
// fichero es editable a mano y no debe poder tumbar ni ensuciar la interfaz.
func (a *Ajustes) validar() {
	def := Defecto()
	if a.Nombre = LimpiarNombre(a.Nombre); a.Nombre == "" {
		a.Nombre = def.Nombre
	}
	if !slices.Contains(OpcTantos, a.Reglas.Tantos) {
		a.Reglas.Tantos = def.Reglas.Tantos
	}
	if !slices.Contains(OpcImpar, a.Reglas.JuegosPorVaca) {
		a.Reglas.JuegosPorVaca = def.Reglas.JuegosPorVaca
	}
	if !slices.Contains(OpcImpar, a.Reglas.Vacas) {
		a.Reglas.Vacas = def.Reglas.Vacas
	}
	if a.Velocidad < PasoAPaso || a.Velocidad > Rapida {
		a.Velocidad = def.Velocidad
	}
	if a.Senas != SenasDeVerdad {
		a.Senas = SenasEscritas
	}
	if !a.Idioma.Valido() {
		a.Idioma = i18n.ES
	}
}

// Guardar escribe los ajustes de forma atómica (fichero temporal + rename) y
// solo legibles por el usuario.
func (a Ajustes) Guardar() error {
	p, err := ruta()
	if err != nil {
		return err
	}
	dir := filepath.Dir(p)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, ".ajustes-*.json")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmp.Name(), p); err != nil {
		return err
	}
	// En el navegador, la página los guarda para la próxima visita.
	canal.Enviar(canal.Mensaje{T: "ajustes", D: string(data)})
	return nil
}
