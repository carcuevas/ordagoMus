// Package canal es el canal lateral entre una sesión web del juego (el proceso
// hijo que arranca el servidor) y el navegador: por él van los sonidos y los
// ajustes, aparte de la salida del terminal para no mezclarse con el dibujo.
//
// El servidor abre el canal como descriptor 3 del hijo y le pone
// ORDAGO_WEB=1. Cada mensaje es una línea JSON.
package canal

import (
	"encoding/json"
	"os"
	"sync"
)

// EnvWeb es la variable de entorno que marca una sesión web.
const EnvWeb = "ORDAGO_WEB"

// Fd es el descriptor del canal en el proceso hijo.
const Fd = 3

var (
	once sync.Once
	mu   sync.Mutex
	f    *os.File
)

// Web dice si el proceso es una sesión web.
func Web() bool { return os.Getenv(EnvWeb) == "1" }

// Mensaje es lo que viaja por el canal.
type Mensaje struct {
	T string `json:"t"`           // "sonido" o "ajustes"
	N string `json:"n,omitempty"` // nombre del sonido
	D string `json:"d,omitempty"` // ajustes en JSON
}

// Enviar manda un mensaje al navegador; fuera de una sesión web no hace nada.
func Enviar(m Mensaje) {
	if !Web() {
		return
	}
	once.Do(func() { f = os.NewFile(Fd, "canal") })
	if f == nil {
		return
	}
	b, err := json.Marshal(m)
	if err != nil {
		return
	}
	mu.Lock()
	defer mu.Unlock()
	f.Write(append(b, '\n'))
}
