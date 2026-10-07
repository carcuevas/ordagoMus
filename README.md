# Órdago · Mus

Juego de mus para terminal inspirado en el antiguo *Órdago* para DOS de Pedro W. Torrecilla. Las reglas son las del
reglamento de Bizkaia: 8 reyes y 8 ases, "ganan las 31 de mano", los duples los gana
el par mayor, señas de Euskadi y primera mano a mus corrido y sin señas.

## Jugar

Hay binarios para Linux (amd64 y arm64) en [Releases](https://github.com/carcuevas/ordagoMus/releases):

```sh
curl -LO https://github.com/carcuevas/ordagoMus/releases/latest/download/ordago-linux-$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')
chmod +x ordago-linux-* && ./ordago-linux-*
```

O compilando desde el código:

```sh
go build -o ordago ./cmd/ordago
./ordago             # menú principal
./ordago -sim 50     # 50 partidas entre bots (para ajustar la IA)
```

El juego está en castellano, inglés, checo y euskera (**Opciones → Idioma**), y en **Cómo se
juega** están las reglas en los cuatro idiomas. Las palabras del mus —*mus*, *envido*, *grande*,
*órdago*…— se dejan siempre en castellano.

Desde el menú eliges compañero y rivales entre 12 personajes de toda España, cada uno con su
historia, sus frases y su forma de jugar. En **Opciones** se configuran tu nombre, 8 o 4 reyes,
los tantos, los juegos y vacas, y la velocidad del ordenador (*Paso a paso* espera a que
pulses espacio). Las opciones se guardan en `~/.config/ordago/ajustes.json`.

**Caras:** como en el Órdago original, cada personaje tiene su retrato pixelado (medios bloques,
se ve en cualquier terminal). Pone cara según lo que juega —sonríe al envidar, se enfada en el
órdago, se le cae la cara al no querer— y al acabar la mano. Y lo mejor: cuando ves la seña de tu
compañero (o pillas la de un rival) se la ves hacer en la cara: morderse el labio, sacar la lengua,
guiñar un ojo... A Xosé, como es imposible de leer, solo se le ve hablar. En kitty las caras son imágenes en alta
resolución (el retrato suavizado con Scale2x y los ojos, cejas y boca dibujados encima), para que los
gestos se distingan bien.

**Señas:** en Opciones → Señas se elige el modo. *Escritas*: se pasan solas y se escribe lo que pasa
cada uno. *De verdad*: no se escribe nada; el gesto dura un instante (100–1000 ms) en la cara de tu
compañero, y tú pasas las tuyas con `s` y un número (solo las que llevas: no se miente). Tu compañero
puede no verla —si la ve, te lo confirma— y un rival puede pillarla; si te la pillan, te enteras.

**Amarracos:** el marcador y los envites se ven con amarracos (5 tantos) y piedras (1), como en la mesa.

El terminal debe tener al menos 110 columnas.

**Cartas ilustradas:** en kitty (y Ghostty) las cartas se ven con la baraja de Heraclio Fournier
de 1878, que es de dominio público (ver `arte/LEEME.md`). Usa el protocolo gráfico de kitty con
marcadores Unicode, así que sigue siendo una aplicación de consola. En otros terminales, o dentro
de tmux, se usan las cartas de texto. Se puede cambiar en Opciones → Cartas.

### En el navegador

El mismo binario sirve el juego por el navegador, con la misma interfaz (baraja ilustrada, caras
y sonidos incluidos):

```sh
./ordago -servidor localhost:8080            # y abre http://localhost:8080
./ordago -servidor :8080 -clave unaclave     # para la red: entra con http://equipo:8080/?clave=unaclave
```

Cada pestaña juega su propia partida contra los bots, en un proceso aparte. Los ajustes se
guardan en el navegador. Opciones: `-sesiones N` (partidas a la vez, 4 por defecto) y `-clave`
(o la variable `ORDAGO_CLAVE`). Las sesiones se cierran tras 30 minutos sin tocar una tecla.
Si escuchas fuera de `localhost`, pon clave; y si es por internet, ponlo detrás de un proxy con
HTTPS. La página usa [xterm.js](https://xtermjs.org/) (MIT), incluido en el binario.

Teclas en la mesa: `m` mus · `c` corto · `1-4` marcar descarte · `enter` confirmar ·
`p` paso · `e` envido / dos más · `n` envido de N · `q` quiero · `x` no quiero ·
`o` órdago · `s` pasar seña (señas de verdad) · `espacio` siguiente jugada (paso a paso) · `esc` abandonar.

## Estructura

| Paquete | Qué hace |
|---|---|
| `internal/cards` | Baraja española de 40 cartas |
| `internal/rules` | Valoración de grande, chica, pares, juego, punto y señas |
| `internal/game` | Motor: máquina de estados pura (acciones → eventos), `View` por asiento |
| `internal/ai` | Bots: simulación Monte Carlo que tiene en cuenta señas y lo cantado, modulada por un `Perfil` de habilidades |
| `internal/personajes` | Los 12 personajes: historia, perfil, vista/disimulo para las señas y frases |
| `internal/ajustes` | Opciones guardadas en disco |
| `internal/frases` | Refranero y jerga de mus |
| `internal/tui` | Menús y mesa en terminal (Bubble Tea + Lip Gloss) |
| `internal/kitty` | Imágenes en el terminal con el protocolo gráfico de kitty |
| `internal/i18n` | Traducciones (castellano, inglés, checo y euskera) |
| `internal/sonido` | Efectos de sonido sintetizados |
| `internal/servidor` | Modo navegador: HTTP + websocket, una partida por pestaña en un pseudoterminal |
| `internal/canal` | Canal lateral de cada partida web (sonidos y ajustes) hacia el navegador |
| `arte` | Baraja Fournier 1878 incrustada en el binario |

El motor no sabe quién es humano y quién máquina: cada asiento solo recibe su `View`
(sus cartas + lo público). Eso es lo que permitirá jugar en red sin cambiar las reglas.

## Hoja de ruta

1. **Fase 1 (hecha):** 1 humano + 3 personajes en terminal, menú, opciones y fichas.
2. **Fase 2:** que el jugador elija qué señas pasar (hecho: señas de verdad), envido a las dos (grande y chica a la vez),
   cortar el mus fuera de turno, más frases y bots con más carácter.
3. **Fase 3:** red. Un servidor que tenga el `game.Game` y clientes que manden acciones y
   reciban eventos filtrados por asiento (`Event.To`). Mesas con 2, 3 o 4 humanos y bots
   rellenando los huecos.
4. **Fase 4:** más gráficos en consola: retratos de los personajes (hechos), tapete y animaciones.
