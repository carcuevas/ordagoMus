// Package ajustes guarda las opciones del jugador entre partidas.
package ajustes

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

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

func (v Velocidad) String() string { return velocidadNombres[v] }

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

func (m ModoSenas) String() string { return modoNombres[m] }

type Ajustes struct {
	Nombre    string
	Reglas    rules.Config
	Velocidad Velocidad
	Imagenes  bool      // baraja ilustrada si el terminal la admite
	Companero string    // ID del personaje
	Rivales   [2]string // "" = al azar
	Senas     ModoSenas
}

func Defecto() Ajustes {
	return Ajustes{Nombre: "Tú", Reglas: rules.DefaultConfig(), Velocidad: Normal, Imagenes: true, Companero: "remedios"}
}

func ruta() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "ordago", "ajustes.json"), nil
}

// Cargar devuelve los ajustes guardados o los de por defecto.
func Cargar() Ajustes {
	a := Defecto()
	p, err := ruta()
	if err != nil {
		return a
	}
	data, err := os.ReadFile(p)
	if err != nil {
		return a
	}
	if json.Unmarshal(data, &a) != nil || a.Reglas.Tantos <= 0 {
		return Defecto()
	}
	if a.Senas != SenasDeVerdad {
		a.Senas = SenasEscritas
	}
	return a
}

func (a Ajustes) Guardar() error {
	p, err := ruta()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	data, _ := json.MarshalIndent(a, "", "  ")
	return os.WriteFile(p, data, 0o644)
}
