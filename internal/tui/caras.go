package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"ordagomus/internal/kitty"
	"ordagomus/internal/personajes"
	"ordagomus/internal/rules"
)

// Las caras de los personajes, como en el Órdago de toda la vida: retratos
// pixelados de 14×16 que se dibujan con medios bloques (▀▄), así que se ven en
// cualquier terminal. Cambian de cara según la jugada y hacen el gesto de las
// señas cuando las ves.
//
// Todas comparten la misma rejilla: cejas en la fila 6, ojos en la 7 (x 4-5 y
// 8-9) y boca en la 11 (x 5-8). Miran hacia la derecha; el de la derecha de la
// mesa se dibuja en espejo para que mire al centro.

const (
	caraW = 14
	caraH = 16
	// CaraFilas es lo que ocupa una cara en el terminal.
	caraFilas = caraH / 2
)

type expresion int

const (
	exNormal expresion = iota
	exHabla
	exContento
	exEnfadado
	exTriste
)

type cara struct {
	px  [caraH]string
	pal map[byte]string
}

// paleta común; cada cara sobrescribe lo suyo.
var paletaBase = map[byte]string{
	'k': "#161616", // ojos y contornos
	'w': "#f4f2ee",
	'm': "#9c3b36", // boca
	'L': "#e8708a", // lengua
	'r': "#e08a7a", // mofletes
	's': "#f0c4a0",
	'S': "#cc9a74",
}

var caras = map[string]cara{
	"paco": {px: [caraH]string{
		"..............",
		"....ssssss....",
		"...ssssssss...",
		"..HssssssssH..",
		"..HHssssssHH..",
		"..HssssssssH..",
		"..ssssssssss..",
		".ssssssssssss.",
		".SsssssSssssS.",
		"..srssnnssrs..",
		"..ssssssssss..",
		"...ssmmmmss...",
		"...SssssssS...",
		"....SssssS....",
		"..CCCwkkwCCC..",
		".CCCCwwwwCCCC.",
	}, pal: map[byte]string{'s': "#d9a066", 'S': "#b07a48", 'n': "#d0705a", 'H': "#c4c4c4", 'b': "#9a9a9a", 'C': "#262626"}},

	"remedios": {px: [caraH]string{
		".....HHHH.....",
		"....HHHHHH....",
		"...HHHHHHHH...",
		"..HHHHHHHHHH..",
		"..HHssssssHH..",
		"..HssssssssH..",
		"..ssssssssss..",
		".ssgssggssgss.",
		".SssggssggssS.",
		".wsrssSSssrsw.",
		"..ssssssssss..",
		"...ssmmmmss...",
		"...SssssssS...",
		"....SssssS....",
		"..cccwwwwccc..",
		".ccccwGGwcccc.",
	}, pal: map[byte]string{'H': "#e6e6ee", 'b': "#b8b8c4", 'g': "#c8a040", 'c': "#7a4a9a", 'm': "#c0607a", 'G': "#e0b030"}},

	"inaki": {px: [caraH]string{
		"......T.......",
		"...TTTTTTTT...",
		"..TTTTTTTTTT..",
		".TTTTTTTTTTTT.",
		"..HssssssssH..",
		"..HssssssssH..",
		"..ssssssssss..",
		".ssssssssssss.",
		".SsssssSssssS.",
		"..ssssSSssss..",
		"..ssssssssss..",
		"...ddmmmmdd...",
		"...dddddddd...",
		"....SssssS....",
		"..wwwrrrrwww..",
		".wwwwwrrwwwww.",
	}, pal: map[byte]string{'T': "#1c2232", 'H': "#2a1a10", 'b': "#2a1a10", 's': "#e0a878", 'S': "#b88058", 'd': "#a87e5c", 'r': "#d02020"}},

	"montse": {px: [caraH]string{
		"....HHHHHH....",
		"...HHHHHHHH...",
		"..HHHHHHHHHH..",
		"..HHHHHHHHHH..",
		"..HHHHHsssHH..",
		"..HssssssssH..",
		".HssssssssssH.",
		".HsgssggssgsH.",
		".HssggssggssH.",
		".HsrssSSssrsH.",
		".HssssssssssH.",
		".HHssmmmmssHH.",
		".HHSssssssSHH.",
		".HHHSssssSHHH.",
		"..cccwwwwccc..",
		".ccccwwwwcccc.",
	}, pal: map[byte]string{'H': "#4a2c1a", 'b': "#3a2010", 's': "#f2c9a6", 'S': "#d4a07c", 'g': "#2e2e2e", 'c': "#30406a", 'm': "#c04050"}},

	"xose": {px: [caraH]string{
		"..............",
		"...TTTTTTTT...",
		"..TTTTTTTTTT..",
		"..TTTTTTTTTT..",
		"..UUUUUUUUUUUU",
		"..HssssssssH..",
		"..ssssssssss..",
		".ssssssssssss.",
		".SsssssSssssS.",
		"..srssSSssrs..",
		"..HssssssssH..",
		"...ddmmmmdd...",
		"...dddddddd...",
		"....SssssS....",
		"..cccRwRwccc..",
		".ccccwRwRcccc.",
	}, pal: map[byte]string{'T': "#6e6e5c", 'U': "#4a4a3c", 'H': "#9a9a9a", 'b': "#8a8a8a", 'd': "#b0a89c", 's': "#e8b090", 'S': "#c08868", 'r': "#e07870", 'c': "#3a5a30", 'w': "#e8e0d0", 'R': "#b03030"}},

	"manolo": {px: [caraH]string{
		"..............",
		"...RkRkRkRk...",
		"..kRkRkRkRkR..",
		"..RkRkRkRkRk..",
		"..HssssssssHRk",
		"..HssssssssH.R",
		"..ssssssssss..",
		".ssssssssssss.",
		".SsssssSssssS.",
		"..ssssSSssss..",
		"..sBBBBBBBBs..",
		"...ddmmmmdd...",
		"...dddddddd...",
		"....SssssS....",
		"..cccwwwwccc..",
		".cccccwwccccc.",
	}, pal: map[byte]string{'R': "#c42020", 'H': "#1a1a1a", 'b': "#101010", 'B': "#141414", 's': "#d89a70", 'S': "#b07850", 'd': "#8a6a54", 'c': "#3a6aa0"}},

	"carmen": {px: [caraH]string{
		"..............",
		"....TTTTTT....",
		"...TTTTTTTT...",
		"..TTTTTTTTTT..",
		".TTlTTTTTTTTT.",
		"..HssssssssH..",
		"..ssssssssss..",
		".ssssssssssss.",
		".SsssssSssssS.",
		"..srssSSssrs..",
		"..ssssssssss..",
		"...ssmmmmss...",
		"...SssssssS...",
		"....SssssS....",
		"..cccwwwwccc..",
		".ccccGccGcccc.",
	}, pal: map[byte]string{'T': "#141414", 'l': "#6a6a6a", 'H': "#c8b890", 'b': "#9a8a68", 'c': "#4a5a30", 'G': "#e0b030", 'm': "#c02030"}},

	"pepe": {px: [caraH]string{
		"....H.HH.H....",
		"...HHHHHHHHH..",
		"..HHHHHHHHHH..",
		"..HHHHHHHHHH..",
		"..HHsHHssHHH..",
		"..HssssssssH..",
		"..ssssssssss..",
		"ssssssssssssss",
		"SSsssssSssssSS",
		"..sfsfSSfsfs..",
		"..ssssssssss..",
		"...ssmmmmss...",
		"...SssssssS...",
		"....SssssS....",
		"..cccwccwccc..",
		".ccccwccwcccc.",
	}, pal: map[byte]string{'H': "#c07830", 'b': "#a05a20", 'f': "#c88a5a", 's': "#f8d0b0", 'S': "#d8a888", 'c': "#e07020"}},

	"toni": {px: [caraH]string{
		"..HH.HHH.HH...",
		"..HHHHHHHHHH..",
		".HHHhHHHhHHHH.",
		".HHHHHHHHHHHH.",
		"..HHssssssHH..",
		"..HssssssssH..",
		"..ssssssssss..",
		".ssssssssssss.",
		".SsssssSssssS.",
		"..sDssSSssss..",
		"..HssssssssH..",
		"..HBBmmmmBBH..",
		"...BBBBBBBB...",
		"....BBBBBB....",
		"..kkkYRYRkkk..",
		".kkkkYRYRkkkk.",
	}, pal: map[byte]string{'H': "#1a1410", 'h': "#3a2a20", 'b': "#1a1410", 'B': "#2a1e18", 'D': "#5a5a5a", 's': "#e0a878", 'S': "#b88058", 'Y': "#f0c020", 'R': "#c02020"}},

	"yeray": {px: [caraH]string{
		"....HHHHHH....",
		"..HHHHHHHHHH..",
		".HHHHHHHHHHHH.",
		".HggggHHggggH.",
		".HHssssssssHH.",
		".HssssssssssH.",
		".HssssssssssH.",
		".HssssssssssH.",
		".HsssssSssssH.",
		".HsrssSSssrsH.",
		".HHssssssssHH.",
		".HHssmmmmssHH.",
		"..HSssssssSH..",
		"....SssssS....",
		"..cYcswwscYc..",
		".ccccYsscYccc.",
	}, pal: map[byte]string{'H': "#e8c060", 'b': "#b89040", 'g': "#202020", 's': "#c88050", 'S': "#a06038", 'r': "#d07060", 'c': "#20a0a0", 'Y': "#f0d040", 'w': "#f8f0e0"}},

	// Miren: moño gris, pendientes de oro y rebeca verde de txoko.
	"miren": {px: [caraH]string{
		"......HHH.....",
		".....HHHHH....",
		"...HHHHHHHH...",
		"..HHHHHHHHHH..",
		"..HHssssssHH..",
		"..HssssssssH..",
		"..ssssssssss..",
		".ssssssssssss.",
		".SsssssSssssS.",
		".GsrssSSssrsG.",
		"..sSssssssSs..",
		"...ssmmmmss...",
		"...SssssssS...",
		"....SssssS....",
		"..cccwwwwccc..",
		".ccccwwGwcccc.",
	}, pal: map[byte]string{'H': "#a9a9a9", 'b': "#7e7e7e", 's': "#e4b08c", 'S': "#bf8a68", 'r': "#d88070", 'm': "#a04848", 'G': "#d4a020", 'c': "#2f4a3a"}},

	// Ane: coleta morena con goma roja y camiseta de remo a rayas.
	"ane": {px: [caraH]string{
		"....HHHHHH....",
		"...HHHHHHHH...",
		"..HHHHHHHHHH..",
		".HHHHHHHHHHHH.",
		".HHHssssssHHH.",
		"HRHssssssssH..",
		"HHssssssssss..",
		"Hssssssssssss.",
		"HSsssssSssssS.",
		"HHsrssSSssrs..",
		".Hssssssssss..",
		"...ssmmmmss...",
		"...SssssssS...",
		"....SssssS....",
		"..cwcwsswcwc..",
		".cwcwcsscwcwc.",
	}, pal: map[byte]string{'H': "#2a1810", 's': "#f0c09a", 'S': "#cf9a76", 'r': "#e88878", 'm': "#c44a5a", 'R': "#d02030", 'c': "#1f4fa0"}},
}

// pinta pone un píxel; fuera de la rejilla no hace nada.
func pinta(px *[caraH][]byte, x, y int, c byte) {
	if x >= 0 && x < caraW && y >= 0 && y < caraH {
		px[y][x] = c
	}
}

func ojos(px *[caraH][]byte, izq, der bool) {
	for _, o := range []struct {
		x       int
		abierto bool
	}{{4, izq}, {8, der}} {
		if o.abierto {
			pinta(px, o.x, 7, 'w')
			pinta(px, o.x+1, 7, 'k')
		} else {
			pinta(px, o.x, 7, 'k')
			pinta(px, o.x+1, 7, 'k')
		}
	}
}

func cejas(px *[caraH][]byte, y4, y5, y8, y9 int) {
	pinta(px, 4, y4, 'b')
	pinta(px, 5, y5, 'b')
	pinta(px, 8, y8, 'b')
	pinta(px, 9, y9, 'b')
}

func boca(px *[caraH][]byte, s string) {
	for i := 0; i < len(s); i++ {
		if s[i] != ' ' {
			pinta(px, 5+i, 11, s[i])
		}
	}
}

// bajo pinta debajo de la boca solo donde hay piel (no tapa barbas).
func bajo(px *[caraH][]byte, x int, c byte) {
	if strings.IndexByte("sSrdf", px[12][x]) >= 0 {
		px[12][x] = c
	}
}

// pixeles devuelve la cara con la expresión o, si se está haciendo, la seña.
func (c cara) pixeles(ex expresion, sena *rules.Sena) [caraH][]byte {
	var px [caraH][]byte
	for y, row := range c.px {
		px[y] = []byte(row)
	}
	// cara de base
	cejas(&px, 6, 6, 6, 6)
	ojos(&px, true, true)
	boca(&px, "mmmm")

	if sena != nil {
		switch *sena {
		case rules.SenaDosReyes: // se muerde el labio inferior
			boca(&px, "wwww")
			for x := 5; x <= 8; x++ {
				bajo(&px, x, 'm')
			}
		case rules.SenaDosAses: // saca la punta de la lengua
			boca(&px, "mLLm")
			bajo(&px, 6, 'L')
			bajo(&px, 7, 'L')
		case rules.SenaMediasReyes: // tuerce la comisura
			boca(&px, "mmm")
			pinta(&px, 8, 11, px[10][8])
			pinta(&px, 8, 10, 'm')
			pinta(&px, 9, 10, 'm')
		case rules.SenaMediasAses: // saca la lengua hacia un lado
			boca(&px, "mmmL")
			pinta(&px, 9, 11, 'L')
			bajo(&px, 9, 'L')
		case rules.SenaDuples: // levanta las cejas
			cejas(&px, 5, 5, 5, 5)
			for _, x := range []int{4, 5, 8, 9} {
				px[6][x] = 's'
			}
		case rules.Sena31: // guiña un ojo
			ojos(&px, true, false)
			pinta(&px, 9, 10, 'm')
		case rules.SenaCiego: // cierra los dos ojos
			ojos(&px, false, false)
		}
		return px
	}

	switch ex {
	case exHabla:
		boca(&px, "mkkm")
	case exContento:
		boca(&px, "wwww")
		pinta(&px, 4, 11, 'm')
		pinta(&px, 9, 11, 'm')
		for x := 5; x <= 8; x++ {
			bajo(&px, x, 'm')
		}
	case exEnfadado, exTriste:
		for _, x := range []int{4, 5, 8, 9} {
			px[6][x] = 's'
		}
		if ex == exEnfadado { // \ /
			cejas(&px, 5, 6, 6, 5)
		} else { // / \
			cejas(&px, 6, 5, 5, 6)
		}
		bajo(&px, 4, 'm')
		bajo(&px, 9, 'm')
	}
	return px
}

func (c cara) color(b byte) (string, bool) {
	if b == '.' {
		return "", false
	}
	if col, ok := c.pal[b]; ok {
		return col, true
	}
	if b == 'b' { // cejas: del color del pelo si no se dice otra cosa
		if col, ok := c.pal['H']; ok {
			return col, true
		}
	}
	col, ok := paletaBase[b]
	return col, ok
}

// dibujarCara devuelve la cara en caraFilas líneas de caraW columnas.
// espejo la gira para que mire a la izquierda. Con img, en kitty, es la
// imagen en alta resolución.
func dibujarCara(id string, ex expresion, sena *rules.Sena, espejo, img bool) string {
	c, ok := caras[id]
	if !ok {
		return ""
	}
	if img {
		for i, p := range personajes.Todos {
			if p.ID == id {
				return strings.Join(kitty.Lines(caraID(i, variante(ex, sena), espejo), pidCara, caraW, caraFilas), "\n")
			}
		}
	}
	px := c.pixeles(ex, sena)
	lines := make([]string, caraFilas)
	for f := range lines {
		var b strings.Builder
		for i := 0; i < caraW; i++ {
			x := i
			if espejo {
				x = caraW - 1 - i
			}
			arriba, ha := c.color(px[2*f][x])
			abajo, hb := c.color(px[2*f+1][x])
			st := lipgloss.NewStyle()
			switch {
			case ha && hb:
				b.WriteString(st.Foreground(lipgloss.Color(arriba)).Background(lipgloss.Color(abajo)).Render("▀"))
			case ha:
				b.WriteString(st.Foreground(lipgloss.Color(arriba)).Render("▀"))
			case hb:
				b.WriteString(st.Foreground(lipgloss.Color(abajo)).Render("▄"))
			default:
				b.WriteByte(' ')
			}
		}
		lines[f] = b.String()
	}
	return strings.Join(lines, "\n")
}
