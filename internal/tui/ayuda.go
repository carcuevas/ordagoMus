package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Instrucciones de juego en varios idiomas. Las palabras del mus (mus, envido,
// grande, órdago…) no se traducen nunca: así se cantan en la mesa.

type seccion struct{ titulo, texto string }

type idioma struct {
	nombre    string
	titulo    string
	secciones []seccion
}

var idiomas = []idioma{
	{
		nombre: "Español",
		titulo: "CÓMO SE JUEGA AL MUS",
		secciones: []seccion{
			{"La mesa", `Cuatro jugadores en dos parejas; tu compañero se sienta enfrente. Se juega con la baraja española de 40 cartas (sin ochos ni nueves).
Con 8 reyes (la opción por defecto) los treses cuentan como reyes y los doses como ases.
Se reparten cuatro cartas a cada uno. El primero en hablar es "la mano" y el último "el postre". En caso de empate gana siempre quien esté más cerca de la mano.
Gana el juego la pareja que llega antes a los tantos (40 por defecto). Varios juegos hacen una "vaca" y varias vacas la partida.`},
			{"El mus", `Antes de jugar, cada uno, empezando por la mano, dice "mus" si quiere cambiar cartas o "no hay mus" para cortar.
Si los cuatro dan "mus", cada uno descarta las cartas que quiera (de una a cuatro) y recibe otras nuevas. Se puede seguir dando "mus" hasta que alguien corte.
Cuando alguien corta, se juega con las cartas que hay.
La primera mano se juega a "mus corrido" y sin señas hasta que alguien corte.`},
			{"Los lances", `Se juegan cuatro lances, siempre en este orden:
1. "Grande": gana la mano con las cartas más altas. Orden: rey, caballo, sota, 7, 6, 5, 4, (3), (2), as.
2. "Chica": gana la mano con las cartas más bajas.
3. "Pares": primero cada uno dice "pares sí" o "pares no". Hay tres jugadas: "pareja" (dos cartas iguales, vale 1), "medias" (tres iguales, vale 2) y "duples" (dos parejas o cuatro iguales, vale 3).
4. "Juego": las figuras valen 10 y el resto su número. Hay "juego" si la mano suma 31 o más. El orden es 31, 32, 40, 37, 36, 35, 34, 33. "La 31" vale 3; cualquier otro "juego" vale 2.
Si nadie tiene "juego" se juega al "punto": gana quien más se acerque a 30.
En "pares" y "juego" solo hablan los que llevan la jugada. Si solo la lleva una pareja, no hay envite y se la apunta al final.`},
			{"Envidar", `En cada lance, por turno desde la mano, puedes decir:
· "paso": no apuestas.
· "envido": apuestas 2 tantos (o "envido" de más tantos).
· "N más": subes el envite de los rivales.
· "quiero" / "no quiero": aceptas o rechazas el envite.
· "órdago": te juegas el juego entero. Si te dicen "quiero", se enseñan las cartas y quien gane el lance gana el juego.
Si no te quieren el envite, te llevas el "deje": 1 tanto si era el primer envite, o lo que había antes de la última subida.
Si el lance queda "en paso", "grande", "chica" y "punto" dan 1 tanto al que lo gane; "pares" y "juego" se cobran lo que valen.`},
			{"El recuento", `Al terminar la mano se enseñan las cartas y se cuenta en orden: "grande", "chica", "pares", "juego" (o "punto").
Cada pareja cobra los envites queridos que haya ganado y, además, el valor de sus "pares" y su "juego".
El juego termina en cuanto una pareja llega a los tantos, aunque falte por contar.`},
			{"Las señas", `Al repartir, cada jugador le pasa a su compañero la seña de lo que lleva. Nunca se miente:
· "dos reyes": se muerde el labio inferior.
· "dos ases": saca la punta de la lengua.
· "medias de reyes": tuerce la comisura de los labios.
· "medias de ases": saca la lengua hacia un lado.
· "duples": levanta las cejas.
· "la 31": guiña un ojo.
· "ciego" (ni "pares" ni "juego"): cierra los dos ojos.
Ojo: a veces un rival pilla la seña. Si os pillan una, te enteras.
En Opciones → Señas hay dos modos:
· "Escritas": las señas se pasan solas y se escribe lo que pasa cada uno.
· "De verdad": nadie te dice nada; el gesto dura un instante (entre 0,1 y 1 segundo) en la cara de tu compañero, así que hay que estar atento. Las tuyas las pasas tú con s y un número. Puede que tu compañero no la vea (cuando la ve, te lo confirma) y puede que un rival la pille; si no te confirma, repítela, pero cada vez es otra ocasión para que te la pillen.
Las piedras de los tantos se llevan en amarracos: cada amarraco vale cinco.`},
			{"Teclas", `m "mus" · c "no hay mus" (cortar) · 1-4 marcar descarte · enter confirmar
p "paso" · e "envido" / "dos más" · n "envido" de N · q "quiero" · x "no quiero" · o "órdago"
s pasar una seña (señas de verdad) · espacio siguiente jugada (paso a paso) · esc abandonar la partida`},
		},
	},
	{
		nombre: "English",
		titulo: "HOW TO PLAY MUS",
		secciones: []seccion{
			{"The table", `Four players in two teams; your partner sits opposite you. The game uses the 40-card Spanish deck (no eights or nines).
With 8 kings (the default) threes count as kings and twos count as aces.
Everyone is dealt four cards. The first player to speak is "la mano" and the last one is "el postre". Ties are always won by the player closest to "la mano".
A game is won by the first team to reach the target score (40 "tantos" by default). Several games make a "vaca", and several "vacas" make the match.`},
			{"The mus", `Before playing, starting with "la mano", each player says "mus" to exchange cards or "no hay mus" to cut.
If all four say "mus", each player discards as many cards as they like (one to four) and gets new ones. "Mus" can continue until someone cuts.
Once someone cuts, the hands are played as they are.
The first hand is played as "mus corrido" and with no "señas" until someone cuts.`},
			{"The four rounds", `There are four rounds ("lances"), always in this order:
1. "Grande": the highest cards win. Order: king, knight, jack, 7, 6, 5, 4, (3), (2), ace.
2. "Chica": the lowest cards win.
3. "Pares": first everyone declares "pares sí" or "pares no". There are three combinations: "pareja" (two equal cards, worth 1), "medias" (three equal cards, worth 2) and "duples" (two pairs or four equal cards, worth 3).
4. "Juego": face cards are worth 10, the rest their number. You have "juego" if your hand adds up to 31 or more. The order is 31, 32, 40, 37, 36, 35, 34, 33. "La 31" is worth 3; any other "juego" is worth 2.
If nobody has "juego", "punto" is played instead: the hand closest to 30 wins.
In "pares" and "juego" only players holding the combination take part. If only one team has it, there is no betting and that team scores it at the end.`},
			{"Betting", `In each round, taking turns from "la mano", you can say:
· "paso": you don't bet.
· "envido": you bet 2 "tantos" (or "envido" a larger amount).
· "N más": you raise the opponents' bet by N.
· "quiero" / "no quiero": you accept or refuse the bet.
· "órdago": you bet the whole game. If they answer "quiero", the cards are shown and whoever wins that round wins the game.
If your bet is refused you take the "deje": 1 "tanto" if it was the opening bet, or the amount on the table before the last raise.
If a round ends "en paso" (nobody bet), "grande", "chica" and "punto" give 1 "tanto" to the winner; "pares" and "juego" score their value.`},
			{"Scoring", `At the end of the hand the cards are shown and the score is counted in order: "grande", "chica", "pares", "juego" (or "punto").
Each team collects the accepted bets it won, plus the value of its "pares" and its "juego".
The game ends as soon as a team reaches the target score, even if there is still something left to count.`},
			{"The señas", `After the deal, each player secretly signals their hand to their partner with a "seña". Nobody can lie:
· "dos reyes": bites the lower lip.
· "dos ases": shows the tip of the tongue.
· "medias de reyes": twists the corner of the mouth.
· "medias de ases": sticks the tongue out to one side.
· "duples": raises the eyebrows.
· "la 31": winks.
· "ciego" (neither "pares" nor "juego"): closes both eyes.
Careful: sometimes an opponent catches the "seña". If they catch one of yours, you are told.
Options → Señas has two modes:
· "Escritas" (written): the "señas" are passed automatically and you can read what everyone signals.
· "De verdad" (for real): nothing is written; the gesture flashes on your partner's face for an instant (0.1 to 1 second), so keep your eyes open. You pass your own with s and a number. Your partner may miss it (you get a confirmation when they see it) and an opponent may catch it; if there is no confirmation, repeat it, but every time is another chance to be caught.
Points are kept with stones in "amarracos": each "amarraco" is worth five.`},
			{"Keys", `m "mus" · c "no hay mus" (cut) · 1-4 mark discard · enter confirm
p "paso" · e "envido" / "dos más" · n "envido" N · q "quiero" · x "no quiero" · o "órdago"
s pass a "seña" (real señas) · space next move (step by step) · esc leave the game`},
		},
	},
	{
		nombre: "Čeština",
		titulo: "JAK SE HRAJE MUS",
		secciones: []seccion{
			{"Stůl", `Hrají čtyři hráči ve dvou dvojicích; spoluhráč sedí naproti vám. Hraje se se španělskými kartami (40 karet, bez osmiček a devítek).
S 8 králi (výchozí nastavení) se trojky počítají jako králové a dvojky jako esa.
Každý dostane čtyři karty. Hráč, který mluví první, je "la mano", poslední je "el postre". Při shodě vyhrává vždy ten, kdo je blíž k "la mano".
Hru vyhrává dvojice, která jako první dosáhne cílového počtu bodů (výchozí je 40 "tantos"). Několik her tvoří "vaca" a několik "vacas" tvoří celou partii.`},
			{"Mus", `Před hrou každý, počínaje "la mano", řekne "mus", pokud chce měnit karty, nebo "no hay mus", pokud chce hru zastavit.
Pokud všichni čtyři řeknou "mus", každý odhodí libovolný počet karet (jednu až čtyři) a dostane nové. "Mus" se může opakovat, dokud ho někdo nezastaví.
Jakmile ho někdo zastaví, hraje se s kartami, které jsou v ruce.
První rozdání se hraje jako "mus corrido" a bez "señas", dokud ho někdo nezastaví.`},
			{"Čtyři kola", `Hrají se čtyři kola ("lances"), vždy v tomto pořadí:
1. "Grande": vyhrávají nejvyšší karty. Pořadí: král, kůň, spodek, 7, 6, 5, 4, (3), (2), eso.
2. "Chica": vyhrávají nejnižší karty.
3. "Pares": nejprve každý ohlásí "pares sí" nebo "pares no". Existují tři kombinace: "pareja" (dvě stejné karty, má hodnotu 1), "medias" (tři stejné, hodnota 2) a "duples" (dva páry nebo čtyři stejné, hodnota 3).
4. "Juego": figury mají hodnotu 10, ostatní karty své číslo. "Juego" máte, pokud je součet 31 nebo víc. Pořadí je 31, 32, 40, 37, 36, 35, 34, 33. "La 31" má hodnotu 3, každé jiné "juego" hodnotu 2.
Pokud nikdo nemá "juego", hraje se "punto": vyhrává ruka nejblíže 30.
V "pares" a "juego" hrají jen ti, kdo kombinaci mají. Pokud ji má jen jedna dvojice, nesází se a dvojice si ji započítá na konci.`},
			{"Sázky", `V každém kole, postupně od "la mano", můžete říct:
· "paso": nesázíte.
· "envido": sázíte 2 "tantos" (nebo "envido" s vyšší částkou).
· "N más": zvýšíte sázku soupeřů o N.
· "quiero" / "no quiero": sázku přijmete nebo odmítnete.
· "órdago": vsadíte celou hru. Pokud soupeři odpoví "quiero", karty se ukážou a kdo vyhraje toto kolo, vyhrává celou hru.
Když vaši sázku odmítnou, získáte "deje": 1 "tanto", pokud to byla první sázka, jinak částku, která ležela na stole před posledním zvýšením.
Pokud kolo skončí "en paso" (nikdo nesázel), "grande", "chica" a "punto" dají vítězi 1 "tanto"; "pares" a "juego" se počítají podle své hodnoty.`},
			{"Počítání bodů", `Na konci rozdání se karty ukážou a body se počítají v tomto pořadí: "grande", "chica", "pares", "juego" (nebo "punto").
Každá dvojice získá přijaté sázky, které vyhrála, a navíc hodnotu svých "pares" a svého "juego".
Hra končí, jakmile některá dvojice dosáhne cílového počtu bodů, i když ještě zbývá něco dopočítat.`},
			{"Señas", `Po rozdání každý hráč tajně ukáže spoluhráči, co má v ruce, pomocí "seña". Nikdo nemůže lhát:
· "dos reyes": kousne se do spodního rtu.
· "dos ases": vystrčí špičku jazyka.
· "medias de reyes": zkřiví koutek úst.
· "medias de ases": vystrčí jazyk do strany.
· "duples": zvedne obočí.
· "la 31": mrkne.
· "ciego" (ani "pares", ani "juego"): zavře obě oči.
Pozor: soupeř někdy "seña" zahlédne. Když zahlédne tu tvou, dozvíš se to.
V Nastavení → Señas jsou dva režimy:
· "Escritas" (psané): "señas" se předávají samy a je napsáno, co kdo ukazuje.
· "De verdad" (doopravdy): nic se nepíše; gesto se na tváři spoluhráče mihne jen na okamžik (0,1 až 1 sekundu), takže dávej pozor. Své "señas" předáváš sám klávesou s a číslem. Spoluhráč je nemusí zahlédnout (když ano, potvrdí ti to) a soupeř je může zachytit; bez potvrzení ji zopakuj, ale pokaždé je to další šance, že ji někdo zahlédne.
Body se počítají kamínky v "amarracos": každý "amarraco" má hodnotu pět.`},
			{"Klávesy", `m "mus" · c "no hay mus" (zastavit) · 1-4 označit kartu k odhození · enter potvrdit
p "paso" · e "envido" / "dos más" · n "envido" N · q "quiero" · x "no quiero" · o "órdago"
s předat "seña" (doopravdy) · mezerník další tah (krok za krokem) · esc opustit hru`},
		},
	},
}

// lineasAyuda devuelve el texto del idioma ya formateado y partido en líneas
// para el ancho w.
func lineasAyuda(id idioma, w int) []string {
	wrap := lipgloss.NewStyle().Width(w)
	var out []string
	for _, s := range id.secciones {
		out = append(out, styleLance.Render(s.titulo))
		for _, p := range strings.Split(s.texto, "\n") {
			out = append(out, strings.Split(wrap.Render(p), "\n")...)
		}
		out = append(out, "")
	}
	return out
}

func (a *App) keyAyuda(k string) {
	switch k {
	case "left", "h", "shift+tab":
		a.idioma = (a.idioma + len(idiomas) - 1) % len(idiomas)
		a.scroll = 0
	case "right", "l", "tab":
		a.idioma = (a.idioma + 1) % len(idiomas)
		a.scroll = 0
	case "1", "2", "3":
		a.idioma = int(k[0] - '1')
		a.scroll = 0
	case "up", "k":
		a.scroll--
	case "down", "j", " ":
		a.scroll++
	case "pgup":
		a.scroll -= 10
	case "pgdown":
		a.scroll += 10
	case "home", "g":
		a.scroll = 0
	case "esc", "q":
		a.pant, a.cursor = pMenu, 5
	}
}

func (a *App) viewAyuda() string {
	id := idiomas[a.idioma]
	var tabs []string
	for i, x := range idiomas {
		t := " " + string(rune('1'+i)) + " " + x.nombre + " "
		if i == a.idioma {
			tabs = append(tabs, styleTitle.Render(t))
		} else {
			tabs = append(tabs, styleDim.Render(t))
		}
	}
	w := min(max(a.width, 80)-6, 100)
	lines := lineasAyuda(id, w)
	alto := max(a.height, 24) - 8
	maxScroll := max(len(lines)-alto, 0)
	a.scroll = min(max(a.scroll, 0), maxScroll)
	vis := lines[a.scroll:min(a.scroll+alto, len(lines))]
	pos := ""
	if a.scroll < maxScroll {
		pos = styleDim.Render("  ↓ …")
	}
	help := key("←→", "idioma / language / jazyk") + "  " + key("↑↓", "moverse") + "  " + key("esc", "volver")
	return styleTitle.Render(id.titulo) + "  " + strings.Join(tabs, " ") + "\n\n" +
		strings.Join(vis, "\n") + "\n" + pos + "\n" + help
}
