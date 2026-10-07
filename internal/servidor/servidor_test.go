package servidor

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEsLocal(t *testing.T) {
	for addr, want := range map[string]bool{
		"localhost:8080": true, "127.0.0.1:8080": true, "[::1]:8080": true,
		":8080": false, "0.0.0.0:8080": false, "192.168.1.5:80": false, "basura": false,
	} {
		if got := EsLocal(addr); got != want {
			t.Errorf("EsLocal(%q) = %v", addr, got)
		}
	}
}

func TestTamAcotado(t *testing.T) {
	ws := inicio{Cols: 100000, Filas: -3, CeldaW: 1e9, CeldaH: 0}.tam()
	if ws.Cols != maxCols || ws.Rows != 10 || ws.X != uint16(maxCols*64) || ws.Y != 10*8 {
		t.Fatalf("tamaño sin acotar: %+v", ws)
	}
}

func TestClave(t *testing.T) {
	s := &servidor{cfg: Config{Clave: "secreta"}}
	h := s.entrada(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("dentro")) }))

	pedir := func(url string, cookie *http.Cookie) *httptest.ResponseRecorder {
		r := httptest.NewRequest("GET", url, nil)
		if cookie != nil {
			r.AddCookie(cookie)
		}
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		return w
	}
	if w := pedir("/", nil); w.Code != http.StatusUnauthorized {
		t.Errorf("sin clave: %d", w.Code)
	}
	if w := pedir("/?clave=mala", nil); w.Code != http.StatusForbidden {
		t.Errorf("clave mala: %d", w.Code)
	}
	w := pedir("/?clave=secreta", nil)
	if w.Code != http.StatusSeeOther {
		t.Fatalf("clave buena: %d", w.Code)
	}
	c := w.Result().Cookies()
	if len(c) != 1 || !c[0].HttpOnly || c[0].SameSite != http.SameSiteStrictMode {
		t.Fatalf("cookie: %+v", c)
	}
	if w := pedir("/ws", c[0]); w.Code != 200 || w.Body.String() != "dentro" {
		t.Errorf("con cookie: %d", w.Code)
	}
	if w := pedir("/ws", &http.Cookie{Name: cookieClave, Value: "otra"}); w.Code != http.StatusUnauthorized {
		t.Errorf("cookie mala: %d", w.Code)
	}
}

func TestSonidosYCabeceras(t *testing.T) {
	s := &servidor{cfg: Config{Sesiones: 1}}
	r := httptest.NewRequest("GET", "/sonidos/ordago.wav", nil)
	r.SetPathValue("nombre", "ordago.wav")
	w := httptest.NewRecorder()
	s.sonido(w, r)
	if w.Code != 200 || !strings.HasPrefix(w.Body.String(), "RIFF") {
		t.Errorf("sonido: %d", w.Code)
	}
	r.SetPathValue("nombre", "../../etc/passwd")
	w = httptest.NewRecorder()
	s.sonido(w, r)
	if w.Code != 404 {
		t.Errorf("sonido inexistente: %d", w.Code)
	}
	w = httptest.NewRecorder()
	s.cabeceras(http.NotFoundHandler()).ServeHTTP(w, httptest.NewRequest("GET", "/", nil))
	if !strings.Contains(w.Header().Get("Content-Security-Policy"), "frame-ancestors 'none'") {
		t.Error("falta la CSP")
	}
}

func TestLimiteDeSesiones(t *testing.T) {
	s := &servidor{cfg: Config{Sesiones: 1}}
	s.activas.Store(1) // ya hay una en marcha
	w := httptest.NewRecorder()
	s.ws(w, httptest.NewRequest("GET", "/ws", nil))
	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("con el cupo lleno: %d", w.Code)
	}
	if s.activas.Load() != 1 {
		t.Errorf("el contador no vuelve: %d", s.activas.Load())
	}
}
