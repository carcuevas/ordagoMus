package i18n

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"
)

// Recorre el código y comprueba que todo texto literal que pasa por T o Tf
// está traducido a todos los idiomas, y que los formatos conservan los verbos.
func TestCatalogoCompleto(t *testing.T) {
	raiz := filepath.Join("..", "..")
	fset := token.NewFileSet()
	usados := map[string]string{}
	filepath.WalkDir(raiz, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(p, ".go") || strings.HasSuffix(p, "_test.go") {
			return nil
		}
		f, err := parser.ParseFile(fset, p, nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		ast.Inspect(f, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok || len(call.Args) == 0 {
				return true
			}
			var nombre string
			switch fn := call.Fun.(type) {
			case *ast.SelectorExpr:
				if x, ok := fn.X.(*ast.Ident); ok && x.Name == "i18n" {
					nombre = fn.Sel.Name
				}
			case *ast.Ident:
				if strings.HasSuffix(p, filepath.Join("i18n", "i18n.go")) {
					return true
				}
			}
			if nombre != "T" && nombre != "Tf" {
				return true
			}
			if lit, ok := call.Args[0].(*ast.BasicLit); ok && lit.Kind == token.STRING {
				s, _ := strconv.Unquote(lit.Value)
				usados[s] = fset.Position(lit.Pos()).String()
			}
			return true
		})
		return nil
	})
	if len(usados) == 0 {
		t.Fatal("no se encontró ningún texto traducible")
	}
	for s, donde := range usados {
		if f := Falta(s); len(f) > 0 {
			t.Errorf("%s: falta traducir %q a %v", donde, s, f)
		}
	}
}

var verbo = regexp.MustCompile(`%[-+# 0-9.]*[a-zA-Z%]`)

func TestFormatos(t *testing.T) {
	for es, tr := range catalogo {
		want := verbo.FindAllString(es, -1)
		slices.Sort(want)
		for i, s := range tr {
			got := verbo.FindAllString(s, -1)
			slices.Sort(got)
			if s != "" && !slices.Equal(got, want) && !strings.Contains(s, "%[") {
				t.Errorf("%q → %s %q: verbos %v, quiero %v", es, Idioma(i+1), s, got, want)
			}
		}
	}
}

func TestEn(t *testing.T) {
	catalogo["__prueba"] = trad{"test", "zkouška", "proba"}
	defer delete(catalogo, "__prueba")
	defer Poner(Actual())
	for i, w := range []string{"__prueba", "test", "zkouška", "proba"} {
		Poner(Idioma(i))
		if got := T("__prueba"); got != w {
			t.Errorf("%s: %q", Idioma(i), got)
		}
	}
	Poner(EN)
	if T("sin traducir") != "sin traducir" {
		t.Error("lo que falta debe salir en castellano")
	}
}
