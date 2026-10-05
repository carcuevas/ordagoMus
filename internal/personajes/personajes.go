// Package personajes son los jugadores del ordenador: de toda España, cada uno
// con su historia, su forma de jugar y sus frases.
package personajes

import (
	"math/rand/v2"

	"ordagomus/internal/ai"
)

type Personaje struct {
	ID       string
	Nombre   string
	Apodo    string
	Origen   string
	Edad     int
	Oficio   string
	Historia string
	Estilo   string // resumen de cómo juega
	Perfil   ai.Perfil
	Vista    float64 // facilidad para pillar señas a los rivales
	Disimulo float64 // dificultad para que le pillen las suyas
	Frases   []string
}

// Corto es como se le llama en la mesa.
func (p *Personaje) Corto() string { return p.Nombre }

// Completo incluye el apodo: Paco «el Tito».
func (p *Personaje) Completo() string {
	if p.Apodo == "" {
		return p.Nombre
	}
	return p.Nombre + " «" + p.Apodo + "»"
}

func (p *Personaje) Frase(rng *rand.Rand) string {
	return p.Frases[rng.IntN(len(p.Frases))]
}

func perfil(g, c, pa, j, pu, farol, valentia, ordago, musero, lectura, querer float64) ai.Perfil {
	return ai.Perfil{
		Lance: [5]float64{g, c, pa, j, pu},
		Farol: farol, Valentia: valentia, Ordago: ordago, Musero: musero, Lectura: lectura, Querer: querer,
	}
}

var Todos = []*Personaje{
	{
		ID: "paco", Nombre: "Paco", Apodo: "el Tito", Origen: "Sevilla", Edad: 67,
		Oficio: "camarero jubilado de Triana",
		Historia: "Cuarenta años sirviendo cañas en un bar de Triana donde se jugaba al mus " +
			"en la mesa del fondo. Aprendió mirando y ahora no hay quien le cierre la boca ni le " +
			"adivine las cartas. Dice que ganó un jamón en un campeonato del 83 y lo cuenta cada vez que gana un órdago.",
		Estilo: "Farolero de libro. Envida con cualquier cosa y te lo dice con una sonrisa.",
		Perfil: perfil(0.65, 0.55, 0.65, 0.65, 0.55, 0.9, 0.7, 0.6, 0.6, 0.6, 0.55),
		Vista:  0.2, Disimulo: 0.5,
		Frases: []string{
			"¡Ozú, qué cartas me han dao!", "Esto lo gano yo con los ojos cerraos, quillo.",
			"Ni me mires, que me pongo colorao.", "Una vez gané un jamón con peores cartas.",
			"Arsa, que esto se anima.", "Tranquilo, que no llevo nada... o sí.",
		},
	},
	{
		ID: "remedios", Nombre: "Doña Remedios", Origen: "Valladolid", Edad: 74,
		Oficio: "maestra jubilada",
		Historia: "Dio clase de matemáticas durante cuarenta años en un instituto de Valladolid. " +
			"Juega los martes en el hogar del jubilado y nadie recuerda la última vez que perdió una chica. " +
			"No farolea jamás: dice que mentir está feo, aunque sea al mus.",
		Estilo: "La reina de la chica. Prudentísima, casi nunca se tira un farol.",
		Perfil: perfil(0.6, 0.97, 0.7, 0.7, 0.85, 0.05, 0.3, 0.2, 0.4, 0.9, 0.4),
		Vista:  0.15, Disimulo: 0.7,
		Frases: []string{
			"A la chica, hijo, a la chica.", "Con paciencia se gana el cielo y el mus.",
			"Eso en mis tiempos no se hacía.", "Los ases, como los buenos alumnos: callados y a lo suyo.",
			"Mira que os lo tengo dicho.",
		},
	},
	{
		ID: "inaki", Nombre: "Iñaki", Apodo: "Goiko", Origen: "Bilbao", Edad: 45,
		Oficio: "ex pelotari",
		Historia: "Fue delantero de pelota a mano hasta que la muñeca dijo basta. Ahora juega al mus en " +
			"una sociedad gastronómica de Bilbao, donde el mus se toma tan en serio como el txakoli. " +
			"Va a la grande como iba a la pared: de frente y sin miedo.",
		Estilo: "Grande y órdago. Cuando lleva reyes no se corta un pelo.",
		Perfil: perfil(0.92, 0.5, 0.75, 0.7, 0.6, 0.3, 0.85, 0.9, 0.3, 0.7, 0.6),
		Vista:  0.15, Disimulo: 0.6,
		Frases: []string{
			"¡Aupa!", "Órdago y a otra cosa, mariposa.", "En Bilbao esto se juega así, de frente.",
			"Con reyes, a la grande, que para eso están.", "¡Venga, ahí va!",
		},
	},
	{
		ID: "montse", Nombre: "Montse", Apodo: "", Origen: "Barcelona", Edad: 38,
		Oficio: "actuaria de seguros",
		Historia: "Calcula riesgos para una aseguradora y aplica lo mismo en la mesa. Aprendió a jugar " +
			"en la facultad y desde entonces lleva una hoja de cálculo mental de cada partida. " +
			"Lo que más le gusta: los pares, porque ahí los números cuadran.",
		Estilo: "Fría y calculadora. Brilla en los pares y aprovecha cada seña que ve.",
		Perfil: perfil(0.8, 0.8, 0.95, 0.85, 0.8, 0.1, 0.45, 0.3, 0.5, 1.0, 0.45),
		Vista:  0.25, Disimulo: 0.6,
		Frases: []string{
			"Estadísticamente, no.", "Las probabilidades no mienten.", "Endavant.",
			"Un 62% de que lo lleves. Paso.", "Esto ya lo tenía calculado.",
		},
	},
	{
		ID: "xose", Nombre: "Xosé", Apodo: "o Fino", Origen: "Lugo", Edad: 58,
		Oficio: "ganadero",
		Historia: "Tiene vacas rubias en una aldea cerca de Lugo y juega al mus en la feria de los jueves. " +
			"Nadie sabe nunca lo que lleva: si le preguntas si sube o baja la escalera, contesta que depende. " +
			"Sus señas son tan finas que ni su compañero está seguro de haberlas visto.",
		Estilo: "Imposible de leer. Equilibrado en todo, y a él no le pillas una seña.",
		Perfil: perfil(0.75, 0.75, 0.75, 0.75, 0.75, 0.45, 0.5, 0.4, 0.5, 0.8, 0.5),
		Vista:  0.2, Disimulo: 0.95,
		Frases: []string{
			"Depende.", "Pues según se mire...", "Non é por nada, pero quero.",
			"Eu nin sí nin non.", "Isto vai ser que non... ou que si.",
		},
	},
	{
		ID: "manolo", Nombre: "Manolo", Apodo: "el Maño", Origen: "Zaragoza", Edad: 52,
		Oficio: "camionero",
		Historia: "Ha recorrido España entera en camión y ha jugado al mus en todas las áreas de servicio " +
			"entre Zaragoza y Algeciras. Es más tozudo que una mula: si le envidas, lo quiere. " +
			"Le han ganado muchas veces, pero nunca le han hecho callar.",
		Estilo: "Cabezota: casi nunca dice que no quiere. Valiente y algo temerario.",
		Perfil: perfil(0.7, 0.6, 0.7, 0.7, 0.6, 0.35, 0.8, 0.7, 0.45, 0.5, 0.95),
		Vista:  0.15, Disimulo: 0.3,
		Frases: []string{
			"¡Que sí, que lo quiero, cucha!", "Más tozudo que un maño, dicen. Y a mucha honra.",
			"A mí no me achanta nadie.", "¡Pilarica, échame una mano!", "Lo quiero y punto.",
		},
	},
	{
		ID: "carmen", Nombre: "Carmen", Apodo: "la Sargenta", Origen: "Madrid", Edad: 61,
		Oficio: "guardia civil retirada",
		Historia: "Treinta años en la Benemérita le dejaron un ojo que no se le escapa nada. Juega en un " +
			"bar de Lavapiés y los habituales dicen que es mejor no pasar señas cuando ella mira. " +
			"Al juego, la 31 parece perseguirla.",
		Estilo: "Especialista en el juego y la que más señas pilla. Agresiva.",
		Perfil: perfil(0.7, 0.6, 0.7, 0.95, 0.7, 0.4, 0.8, 0.6, 0.4, 0.8, 0.55),
		Vista:  0.45, Disimulo: 0.5,
		Frases: []string{
			"Aquí mando yo.", "A sus órdenes... y órdago.", "Te he visto, majo.",
			"Esa seña la he visto hasta yo sin gafas.", "Circulen, que aquí no hay nada que ver.",
		},
	},
	{
		ID: "pepe", Nombre: "Pepe", Apodo: "el Novato", Origen: "Badajoz", Edad: 23,
		Oficio: "estudiante de Magisterio",
		Historia: "Aprendió a jugar el verano pasado con su abuelo en el pueblo y se ha venido arriba. " +
			"Confunde a veces la grande con la chica, pero tiene una suerte que da rabia. " +
			"Ideal para empezar, o para que te toque de rival.",
		Estilo: "Principiante. Se da mus con todo y valora mal las jugadas.",
		Perfil: perfil(0.35, 0.35, 0.35, 0.4, 0.35, 0.15, 0.4, 0.4, 0.8, 0.3, 0.5),
		Vista:  0.05, Disimulo: 0.1,
		Frases: []string{
			"¿Esto es bueno? Digo... ¡envido!", "Mi abuelo me dijo que siempre a reyes.",
			"Ay, que no sé si quiero.", "¿Los treses eran reyes o no?", "¡Toma ya!",
		},
	},
	{
		ID: "toni", Nombre: "Toni", Apodo: "", Origen: "Valencia", Edad: 33,
		Oficio: "pirotécnico fallero",
		Historia: "Monta mascletàs en Fallas y juega al mus igual: mucho ruido y que no se sepa por dónde " +
			"va a salir. Se marca faroles tremendos al punto y luego se ríe a carcajadas.",
		Estilo: "Explosivo y farolero. Muy bueno al punto.",
		Perfil: perfil(0.65, 0.65, 0.65, 0.65, 0.9, 0.8, 0.65, 0.5, 0.5, 0.6, 0.5),
		Vista:  0.2, Disimulo: 0.4,
		Frases: []string{
			"¡Che, qué mano!", "Esto es más fácil que hacer una paella... bueno, no.",
			"¡Pum! Como una mascletà.", "Al punto no me gana ni mi suegra.", "¡Vinga, va!",
		},
	},
	{
		ID: "yeray", Nombre: "Yeray", Apodo: "", Origen: "Las Palmas", Edad: 29,
		Oficio: "profesor de surf",
		Historia: "Entre ola y ola juega al mus en un chiringuito de Las Canteras. No hay nada que le altere: " +
			"se da mus con medio mazo y espera tranquilo a que lleguen las cartas, como espera las olas.",
		Estilo: "Tranquilo y musero. Sin prisa, pero con buen ojo para el juego.",
		Perfil: perfil(0.65, 0.65, 0.7, 0.8, 0.7, 0.2, 0.4, 0.35, 0.9, 0.65, 0.45),
		Vista:  0.2, Disimulo: 0.6,
		Frases: []string{
			"Tranqui, mi niño.", "Sin estrés, que el mus es pa' disfrutarlo.",
			"Ya vendrá la ola buena.", "¡Fuerte jugada, chacho!", "Ños, qué cartas.",
		},
	},
}

func Buscar(id string) *Personaje {
	for _, p := range Todos {
		if p.ID == id {
			return p
		}
	}
	return nil
}
