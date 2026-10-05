// Package kitty muestra imágenes en el terminal con el protocolo gráfico de
// kitty usando marcadores Unicode (U+10EEEE): la imagen se transmite una vez y
// luego se "pinta" escribiendo caracteres normales, de modo que la interfaz de
// texto puede maquetarla como cualquier otra cadena.
//
// https://sw.kovidgoyal.net/kitty/graphics-protocol/#unicode-placeholders
package kitty

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/sys/unix"
)

const placeholder = '\U0010EEEE'

// diacritics codifican la fila y la columna de cada celda (rowcolumn-diacritics.txt de kitty).
var diacritics = []rune{
	0x0305, 0x030D, 0x030E, 0x0310, 0x0312, 0x033D, 0x033E, 0x033F,
	0x0346, 0x034A, 0x034B, 0x034C, 0x0350, 0x0351, 0x0352, 0x0357,
	0x035B, 0x0363, 0x0364, 0x0365, 0x0366, 0x0367, 0x0368, 0x0369,
	0x036A, 0x036B, 0x036C, 0x036D, 0x036E, 0x036F,
}

// MaxCells es el mayor número de filas o columnas que puede ocupar una imagen.
const MaxCells = 30

// Supported indica si el terminal entiende los marcadores Unicode de kitty.
// Dentro de tmux no funcionan.
func Supported() bool {
	if os.Getenv("TMUX") != "" {
		return false
	}
	return os.Getenv("KITTY_WINDOW_ID") != "" || os.Getenv("TERM") == "xterm-kitty" ||
		os.Getenv("TERM_PROGRAM") == "ghostty"
}

// CellSize devuelve el tamaño en píxeles de una celda del terminal.
func CellSize() (w, h int) {
	ws, err := unix.IoctlGetWinsize(int(os.Stdout.Fd()), unix.TIOCGWINSZ)
	if err != nil || ws.Col == 0 || ws.Row == 0 || ws.Xpixel == 0 || ws.Ypixel == 0 {
		return 10, 20
	}
	return int(ws.Xpixel) / int(ws.Col), int(ws.Ypixel) / int(ws.Row)
}

// Transmit envía un PNG al terminal con el identificador id.
func Transmit(w io.Writer, id uint32, png []byte) error {
	data := base64.StdEncoding.EncodeToString(png)
	const chunk = 4096
	var b strings.Builder
	for i := 0; i < len(data); i += chunk {
		end := min(i+chunk, len(data))
		more := 0
		if end < len(data) {
			more = 1
		}
		if i == 0 {
			fmt.Fprintf(&b, "\x1b_Ga=t,f=100,t=d,i=%d,q=2,m=%d;%s\x1b\\", id, more, data[i:end])
		} else {
			fmt.Fprintf(&b, "\x1b_Gm=%d;%s\x1b\\", more, data[i:end])
		}
	}
	_, err := io.WriteString(w, b.String())
	return err
}

// Place crea (o sustituye) una colocación virtual de la imagen de cols×rows celdas.
func Place(w io.Writer, id, pid uint32, cols, rows int) error {
	_, err := fmt.Fprintf(w, "\x1b_Ga=p,U=1,i=%d,p=%d,c=%d,r=%d,q=2\x1b\\", id, pid, cols, rows)
	return err
}

// DeleteAll borra del terminal las imágenes indicadas.
func DeleteAll(w io.Writer, ids []uint32) {
	var b strings.Builder
	for _, id := range ids {
		fmt.Fprintf(&b, "\x1b_Ga=d,d=I,i=%d,q=2\x1b\\", id)
	}
	io.WriteString(w, b.String())
}

// Lines devuelve las líneas de texto que, en kitty, se ven como la imagen.
// El color de primer plano lleva el id de la imagen y el de subrayado el de la
// colocación.
func Lines(id, pid uint32, cols, rows int) []string {
	cols, rows = min(cols, MaxCells), min(rows, MaxCells)
	pre := fmt.Sprintf("\x1b[38;2;%d;%d;%dm\x1b[58;2;%d;%d;%dm",
		(id>>16)&0xff, (id>>8)&0xff, id&0xff, (pid>>16)&0xff, (pid>>8)&0xff, pid&0xff)
	lines := make([]string, rows)
	for r := 0; r < rows; r++ {
		var b strings.Builder
		b.WriteString(pre)
		for c := 0; c < cols; c++ {
			b.WriteRune(placeholder)
			b.WriteRune(diacritics[r])
			b.WriteRune(diacritics[c])
		}
		b.WriteString("\x1b[39;59m")
		lines[r] = b.String()
	}
	return lines
}
