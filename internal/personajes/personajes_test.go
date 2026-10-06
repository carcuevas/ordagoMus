package personajes

import (
	"math/rand/v2"
	"testing"

	"ordagomus/internal/i18n"
)

func TestTraducidos(t *testing.T) {
	for _, p := range Todos {
		textos := append([]string{p.Oficio, p.Historia, p.Estilo}, p.Frases...)
		for _, s := range textos {
			if f := i18n.Falta(s); len(f) > 0 {
				t.Errorf("%s: falta traducir %q a %v", p.ID, s, f)
			}
		}
	}
}

func TestIDsUnicos(t *testing.T) {
	vistos := map[string]bool{}
	for _, p := range Todos {
		if vistos[p.ID] {
			t.Errorf("ID repetido: %s", p.ID)
		}
		vistos[p.ID] = true
		if Buscar(p.ID) != p {
			t.Errorf("Buscar(%q) no lo encuentra", p.ID)
		}
		if len(p.Frases) == 0 {
			t.Errorf("%s no tiene frases", p.ID)
		}
	}
}

func TestFraseTraducida(t *testing.T) {
	defer i18n.Poner(i18n.Actual())
	i18n.Poner(i18n.EN)
	p := Buscar("remedios")
	rng := rand.New(rand.NewPCG(1, 2))
	for range 20 {
		if f := p.Frase(rng); i18n.Falta(f) == nil {
			t.Fatalf("la frase %q ha salido en castellano", f)
		}
	}
}
