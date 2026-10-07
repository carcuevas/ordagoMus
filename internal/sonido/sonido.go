// Package sonido toca efectos cortos que marcan la fase de la mano: el
// reparto, el corte del mus, el inicio de cada lance, los envites y el
// recuento.
//
// Los sonidos se sintetizan al arrancar (sin ficheros ni cgo, para que el
// binario siga siendo estático), se guardan como WAV en un directorio temporal
// privado y se reproducen con el reproductor del sistema que haya
// (pw-play, paplay o aplay). Si no hay ninguno, el juego sigue en silencio.
package sonido

import (
	"bytes"
	"encoding/binary"
	"math"
	"math/rand/v2"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"ordagomus/internal/canal"
)

type Sonido int

const (
	Reparto  Sonido = iota // se reparten cartas
	Corte                  // alguien corta el mus
	Lance                  // empieza un lance (grande, chica, pares, juego, punto)
	Envite                 // envido o reenvido
	Ordago                 // ¡órdago!
	Recuento               // se enseñan las cartas y se cuentan los tantos
	nSonidos
)

const frecuencia = 22050

var (
	mu       sync.Mutex
	player   string // ruta absoluta del reproductor
	ficheros [nSonidos]string
	activo   bool
)

// reproductores en orden de preferencia; se buscan en el PATH y se ejecutan
// sin shell, con la ruta del WAV como único argumento.
var reproductores = []string{"pw-play", "paplay", "aplay"}

// Preparar sintetiza los sonidos y busca un reproductor. Devuelve la función
// que borra los ficheros temporales; hay que llamarla al salir.
func Preparar() (limpiar func()) {
	nada := func() {}
	if canal.Web() {
		// En el navegador los sonidos los toca la página.
		mu.Lock()
		activo = true
		mu.Unlock()
		return nada
	}
	var bin string
	for _, r := range reproductores {
		if p, err := exec.LookPath(r); err == nil {
			bin = p
			break
		}
	}
	if bin == "" {
		return nada
	}
	dir, err := os.MkdirTemp("", "ordago-sonidos-") // 0700
	if err != nil {
		return nada
	}
	limpiar = func() {
		mu.Lock()
		activo = false
		mu.Unlock()
		os.RemoveAll(dir)
	}
	var fs [nSonidos]string
	for s := Sonido(0); s < nSonidos; s++ {
		p := filepath.Join(dir, nombres[s]+".wav")
		if err := os.WriteFile(p, wav(sintetizar(s)), 0o600); err != nil {
			limpiar()
			return nada
		}
		fs[s] = p
	}
	mu.Lock()
	player, ficheros, activo = bin, fs, true
	mu.Unlock()
	return limpiar
}

var nombres = [nSonidos]string{"reparto", "corte", "lance", "envite", "ordago", "recuento"}

// Nombres de los sonidos, para servirlos al navegador.
func Nombres() []string { return nombres[:] }

// WAV devuelve el sonido de ese nombre codificado en WAV.
func WAV(nombre string) ([]byte, bool) {
	for s, n := range nombres {
		if n == nombre {
			return wav(sintetizar(Sonido(s))), true
		}
	}
	return nil, false
}

// Tocar reproduce el sonido en segundo plano. No hace nada si no se llamó a
// Preparar o si no hay reproductor.
func Tocar(s Sonido) {
	if s < 0 || s >= nSonidos {
		return
	}
	mu.Lock()
	bin, f, ok := player, ficheros[s], activo
	mu.Unlock()
	if !ok {
		return
	}
	if canal.Web() {
		canal.Enviar(canal.Mensaje{T: "sonido", N: nombres[s]})
		return
	}
	cmd := exec.Command(bin, f)
	// nil = /dev/null: el reproductor no debe escribir en la interfaz.
	cmd.Stdin, cmd.Stdout, cmd.Stderr = nil, nil, nil
	if cmd.Start() == nil {
		go cmd.Wait()
	}
}

// ───────────── síntesis ─────────────

func muestras(seg float64) []float64 { return make([]float64, int(seg*frecuencia)) }

// tono suma a buf una nota con envolvente exponencial a partir de t0 segundos.
// parciales son múltiplos de f con su amplitud relativa.
func tono(buf []float64, t0, f, dur, amp, caida float64, parciales ...[2]float64) {
	if len(parciales) == 0 {
		parciales = [][2]float64{{1, 1}}
	}
	ini := int(t0 * frecuencia)
	n := int(dur * frecuencia)
	for i := 0; i < n && ini+i < len(buf); i++ {
		t := float64(i) / frecuencia
		env := math.Exp(-caida*t) * min(1, t*400) // ataque de 2,5 ms sin chasquido
		v := 0.0
		for _, p := range parciales {
			v += p[1] * math.Sin(2*math.Pi*f*p[0]*t)
		}
		buf[ini+i] += amp * env * v
	}
}

// roce suma ruido filtrado (paso bajo de un polo): el sonido de una carta al
// deslizarse sobre el tapete.
func roce(buf []float64, rng *rand.Rand, t0, dur, amp, suave float64) {
	ini := int(t0 * frecuencia)
	n := int(dur * frecuencia)
	y := 0.0
	for i := 0; i < n && ini+i < len(buf); i++ {
		t := float64(i) / float64(n)
		env := math.Sin(math.Pi*t) * (1 - t)
		y += suave * (rng.Float64()*2 - 1 - y)
		buf[ini+i] += amp * env * y
	}
}

func sintetizar(s Sonido) []float64 {
	rng := rand.New(rand.NewPCG(uint64(s), 7)) // siempre el mismo sonido
	switch s {
	case Reparto:
		// cuatro cartas que vuelan sobre el tapete
		b := muestras(0.62)
		for i := 0; i < 4; i++ {
			roce(b, rng, float64(i)*0.13, 0.11, 1.6, 0.35)
		}
		return b
	case Corte:
		// dos golpes secos en la mesa: «no hay mus»
		b := muestras(0.45)
		for i, f := range []float64{150, 120} {
			t := float64(i) * 0.16
			tono(b, t, f, 0.2, 0.7, 28, [2]float64{1, 1}, [2]float64{2.7, 0.3})
			roce(b, rng, t, 0.025, 0.6, 0.6)
		}
		return b
	case Lance:
		// campanilla suave
		b := muestras(0.7)
		tono(b, 0, 1175, 0.7, 0.35, 6, [2]float64{1, 1}, [2]float64{2.76, 0.25}, [2]float64{5.4, 0.08})
		return b
	case Envite:
		// dos notas que suben
		b := muestras(0.42)
		tono(b, 0, 659, 0.18, 0.4, 10, [2]float64{1, 1}, [2]float64{2, 0.3})
		tono(b, 0.12, 988, 0.3, 0.4, 9, [2]float64{1, 1}, [2]float64{2, 0.3})
		return b
	case Ordago:
		// golpe grave y arpegio que sube
		b := muestras(1.3)
		tono(b, 0, 98, 1.3, 0.5, 2.5, [2]float64{1, 1}, [2]float64{2.01, 0.5}, [2]float64{3.02, 0.3})
		for i, f := range []float64{392, 494, 587, 784} {
			tono(b, 0.08+float64(i)*0.1, f, 0.9, 0.22, 4, [2]float64{1, 1}, [2]float64{2, 0.25})
		}
		return b
	case Recuento:
		// amarracos que caen uno tras otro
		b := muestras(0.85)
		for i := 0; i < 5; i++ {
			f := 2200 + rng.Float64()*500
			tono(b, float64(i)*0.11, f, 0.3, 0.18, 18, [2]float64{1, 1}, [2]float64{1.48, 0.5}, [2]float64{2.31, 0.3})
		}
		return b
	}
	return nil
}

// wav codifica las muestras como PCM de 16 bits mono, normalizando el pico.
func wav(x []float64) []byte {
	pico := 0.0
	for _, v := range x {
		pico = max(pico, math.Abs(v))
	}
	escala := 0.0
	if pico > 0 {
		escala = 0.8 * 32767 / pico
	}
	datos := len(x) * 2
	var b bytes.Buffer
	b.WriteString("RIFF")
	binary.Write(&b, binary.LittleEndian, uint32(36+datos))
	b.WriteString("WAVEfmt ")
	for _, v := range []any{
		uint32(16), uint16(1), uint16(1), uint32(frecuencia), uint32(frecuencia * 2), uint16(2), uint16(16),
	} {
		binary.Write(&b, binary.LittleEndian, v)
	}
	b.WriteString("data")
	binary.Write(&b, binary.LittleEndian, uint32(datos))
	for _, v := range x {
		binary.Write(&b, binary.LittleEndian, int16(v*escala))
	}
	return b.Bytes()
}
