package ajustes

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLimpiarNombre(t *testing.T) {
	casos := map[string]string{
		"Ana":                       "Ana",
		"  Ana  ":                   "Ana",
		"\x1b]52;c;ZXZpbA==\x07Ana": "]52;c;ZXZpbA==An",
		"a\x1b[2Jb":                 "a[2Jb",
		"ñññññññññññññññññññññ": "ññññññññññññññññ",
	}
	for in, want := range casos {
		if got := LimpiarNombre(in); got != want {
			t.Errorf("LimpiarNombre(%q) = %q, quiero %q", in, got, want)
		}
	}
}

func TestCargarValida(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("HOME", dir)
	p, err := ruta()
	if err != nil {
		t.Fatal(err)
	}
	os.MkdirAll(filepath.Dir(p), 0o700)
	malo := `{"Nombre":"\u001b[31mX","Reglas":{"Tantos":-5,"JuegosPorVaca":0,"Vacas":99},"Velocidad":42,"Senas":7}`
	if err := os.WriteFile(p, []byte(malo), 0o600); err != nil {
		t.Fatal(err)
	}
	a := Cargar()
	def := Defecto()
	if a.Nombre != "[31mX" || a.Reglas != def.Reglas || a.Velocidad != def.Velocidad || a.Senas != SenasEscritas {
		t.Fatalf("ajustes sin validar: %+v", a)
	}
	_ = a.Velocidad.String() // no debe hacer panic
}

func TestGuardarPermisos(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("HOME", dir)
	if err := Defecto().Guardar(); err != nil {
		t.Fatal(err)
	}
	p, _ := ruta()
	fi, err := os.Stat(p)
	if err != nil {
		t.Fatal(err)
	}
	if perm := fi.Mode().Perm(); perm != 0o600 {
		t.Errorf("permisos %o, quiero 600", perm)
	}
	if got := Cargar(); got != Defecto() {
		t.Errorf("ida y vuelta: %+v", got)
	}
}
