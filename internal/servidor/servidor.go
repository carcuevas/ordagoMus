// Package servidor sirve el juego por el navegador: cada pestaña abre un
// websocket y el servidor arranca para ella una partida propia (el mismo
// binario, en un pseudoterminal), cuya salida pinta xterm.js en la página.
// Así la interfaz es exactamente la de la consola.
//
// Cada sesión es un proceso aparte: el idioma, los ajustes y los gráficos son
// suyos y no se mezclan con los de otras pestañas.
package servidor

import (
	"bufio"
	"context"
	"crypto/subtle"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/coder/websocket"
	"github.com/creack/pty"

	"ordagomus/internal/canal"
	"ordagomus/internal/sonido"
)

//go:embed web
var webFS embed.FS

// Config son las opciones del servidor.
type Config struct {
	Addr     string        // dirección de escucha, p. ej. "localhost:8080"
	Sesiones int           // partidas simultáneas como máximo
	Clave    string        // si no está vacía, hace falta para jugar
	Inactivo time.Duration // se cierra la sesión tras este tiempo sin teclas
}

const (
	maxMensaje   = 64 << 10 // lo más que se acepta del navegador en un mensaje
	maxAjustes   = 16 << 10
	cookieClave  = "ordago_clave"
	esperaInicio = 10 * time.Second
	maxCols      = 500
	maxFilas     = 200
)

type servidor struct {
	cfg     Config
	activas atomic.Int32
	exe     string
}

// Escuchar arranca el servidor y no vuelve hasta que falla o se cancela ctx.
func Escuchar(ctx context.Context, cfg Config) error {
	if cfg.Sesiones <= 0 {
		cfg.Sesiones = 4
	}
	if cfg.Inactivo <= 0 {
		cfg.Inactivo = 30 * time.Minute
	}
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	s := &servidor{cfg: cfg, exe: exe}

	mux := http.NewServeMux()
	estaticos, _ := fs.Sub(webFS, "web")
	mux.Handle("GET /", s.cabeceras(http.FileServerFS(estaticos)))
	mux.HandleFunc("GET /sonidos/{nombre}", s.sonido)
	mux.HandleFunc("GET /ws", s.ws)

	ln, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		return err
	}
	hs := &http.Server{
		Handler:           s.entrada(mux),
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       2 * time.Minute,
		MaxHeaderBytes:    16 << 10,
	}
	go func() {
		<-ctx.Done()
		sc, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		hs.Shutdown(sc)
	}()
	log.Printf("Órdago · Mus en http://%s (hasta %d partidas a la vez)", ln.Addr(), cfg.Sesiones)
	if err := hs.Serve(ln); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// EsLocal dice si la dirección solo escucha en el propio equipo.
func EsLocal(addr string) bool {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return false
	}
	if host == "localhost" {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

// entrada comprueba la clave (si la hay). La primera vez se pasa como
// ?clave=… y se guarda en una cookie, para que no se quede en la barra.
func (s *servidor) entrada(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.cfg.Clave == "" {
			next.ServeHTTP(w, r)
			return
		}
		if c := r.URL.Query().Get("clave"); c != "" {
			if !s.claveOK(c) {
				http.Error(w, "clave incorrecta", http.StatusForbidden)
				return
			}
			http.SetCookie(w, &http.Cookie{
				Name: cookieClave, Value: c, Path: "/", HttpOnly: true,
				SameSite: http.SameSiteStrictMode, Secure: r.TLS != nil,
			})
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		if c, err := r.Cookie(cookieClave); err != nil || !s.claveOK(c.Value) {
			http.Error(w, "hace falta la clave: abre /?clave=…", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *servidor) claveOK(c string) bool {
	return subtle.ConstantTimeCompare([]byte(c), []byte(s.cfg.Clave)) == 1
}

func (s *servidor) cabeceras(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline'; "+
			"img-src 'self' data: blob:; media-src 'self'; connect-src 'self'; frame-ancestors 'none'; base-uri 'none'")
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("Referrer-Policy", "no-referrer")
		next.ServeHTTP(w, r)
	})
}

func (s *servidor) sonido(w http.ResponseWriter, r *http.Request) {
	nombre := strings.TrimSuffix(r.PathValue("nombre"), ".wav")
	data, ok := sonido.WAV(nombre)
	if !ok {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "audio/wav")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Write(data)
}

// inicio es el primer mensaje del navegador: tamaño del terminal, de la
// celda en píxeles y los ajustes guardados en la página.
type inicio struct {
	T       string  `json:"t"`
	Cols    int     `json:"cols"`
	Filas   int     `json:"rows"`
	CeldaW  float64 `json:"cw"`
	CeldaH  float64 `json:"ch"`
	Ajustes string  `json:"ajustes"`
}

func (m inicio) tam() *pty.Winsize {
	cols := min(max(m.Cols, 20), maxCols)
	filas := min(max(m.Filas, 10), maxFilas)
	cw := min(max(m.CeldaW, 4), 64)
	ch := min(max(m.CeldaH, 8), 128)
	return &pty.Winsize{Cols: uint16(cols), Rows: uint16(filas), X: uint16(float64(cols) * cw), Y: uint16(float64(filas) * ch)}
}

func (s *servidor) ws(w http.ResponseWriter, r *http.Request) {
	if n := s.activas.Add(1); int(n) > s.cfg.Sesiones {
		s.activas.Add(-1)
		http.Error(w, "hay demasiadas partidas en marcha; prueba dentro de un rato", http.StatusServiceUnavailable)
		return
	}
	defer s.activas.Add(-1)

	// Accept rechaza por defecto los orígenes de otros sitios.
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		return
	}
	defer c.CloseNow()
	c.SetReadLimit(maxMensaje)

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	ini, err := leerInicio(ctx, c)
	if err != nil {
		c.Close(websocket.StatusPolicyViolation, "inicio no válido")
		return
	}
	if err := s.sesion(ctx, c, ini); err != nil {
		log.Printf("sesión %s: %v", r.RemoteAddr, err)
	}
}

func leerInicio(ctx context.Context, c *websocket.Conn) (inicio, error) {
	ictx, cancel := context.WithTimeout(ctx, esperaInicio)
	defer cancel()
	tipo, data, err := c.Read(ictx)
	if err != nil {
		return inicio{}, err
	}
	var ini inicio
	if tipo != websocket.MessageText || json.Unmarshal(data, &ini) != nil || ini.T != "inicio" {
		return inicio{}, errors.New("primer mensaje no válido")
	}
	if len(ini.Ajustes) > maxAjustes || (ini.Ajustes != "" && !json.Valid([]byte(ini.Ajustes))) {
		ini.Ajustes = ""
	}
	return ini, nil
}

// sesion arranca la partida en un pseudoterminal y la conecta al websocket.
func (s *servidor) sesion(ctx context.Context, c *websocket.Conn, ini inicio) error {
	dir, err := os.MkdirTemp("", "ordago-sesion-") // 0700
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	if ini.Ajustes != "" {
		cfgDir := filepath.Join(dir, "ordago")
		if err := os.MkdirAll(cfgDir, 0o700); err == nil {
			os.WriteFile(filepath.Join(cfgDir, "ajustes.json"), []byte(ini.Ajustes), 0o600)
		}
	}

	lado, ladoHijo, err := os.Pipe()
	if err != nil {
		return err
	}
	defer lado.Close()

	cmd := exec.Command(s.exe)
	cmd.Dir = dir
	// Entorno mínimo: nada del servidor (claves incluidas) pasa a la partida.
	cmd.Env = []string{
		canal.EnvWeb + "=1",
		"HOME=" + dir,
		"XDG_CONFIG_HOME=" + dir,
		"TMPDIR=" + dir,
		"TERM=xterm-256color",
		"COLORTERM=truecolor",
		"LANG=C.UTF-8",
	}
	cmd.ExtraFiles = []*os.File{ladoHijo} // descriptor 3: canal.Fd
	tty, err := pty.StartWithSize(cmd, ini.tam())
	ladoHijo.Close()
	if err != nil {
		return err
	}
	defer tty.Close()
	defer func() {
		// pty arranca el hijo en su propia sesión: se mata el grupo entero.
		syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		cmd.Wait()
	}()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	ultima := atomic.Int64{}
	ultima.Store(time.Now().Unix())

	// Partida → navegador.
	go func() {
		defer cancel()
		buf := make([]byte, 32<<10)
		for {
			n, err := tty.Read(buf)
			if n > 0 {
				if c.Write(ctx, websocket.MessageBinary, buf[:n]) != nil {
					return
				}
			}
			if err != nil {
				return
			}
		}
	}()
	// Canal lateral (sonidos, ajustes) → navegador, como texto.
	go func() {
		sc := bufio.NewScanner(lado)
		sc.Buffer(make([]byte, 4096), maxAjustes*2)
		for sc.Scan() {
			if c.Write(ctx, websocket.MessageText, sc.Bytes()) != nil {
				return
			}
		}
	}()
	// Inactividad.
	go func() {
		t := time.NewTicker(time.Minute)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				if time.Since(time.Unix(ultima.Load(), 0)) > s.cfg.Inactivo {
					c.Close(websocket.StatusNormalClosure, "sesión cerrada por inactividad")
					cancel()
					return
				}
			}
		}
	}()

	// Navegador → partida: teclas (binario) y cambios de tamaño (texto).
	for {
		tipo, data, err := c.Read(ctx)
		if err != nil {
			if ctx.Err() != nil || websocket.CloseStatus(err) != -1 || errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		ultima.Store(time.Now().Unix())
		switch tipo {
		case websocket.MessageBinary:
			if _, err := tty.Write(data); err != nil {
				return nil
			}
		case websocket.MessageText:
			var m inicio
			if json.Unmarshal(data, &m) == nil && m.T == "tam" {
				pty.Setsize(tty, m.tam())
			}
		}
	}
}

// Aviso explica por qué conviene una clave al escuchar fuera del equipo.
func Aviso(addr string) string {
	return fmt.Sprintf("Atención: %s acepta conexiones de otros equipos y no hay clave. "+
		"Usa -clave (o ORDAGO_CLAVE) o escucha en localhost.", addr)
}
