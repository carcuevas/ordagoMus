package i18n

// Los personajes (oficio, historia, estilo y frases) y el refranero de mus
// (paquete frases). Las frases en gallego, catalán, valenciano o euskera se
// dejan como están: son parte del personaje. «¡Órdago!» no se registra aquí:
// no se traduce y ya lo dice el juego.

func init() {
	registrar(map[string]trad{
		// ── Paco ──
		"camarero jubilado de Triana": {"retired waiter from Triana", "číšník v důchodu z Triany", "Trianako zerbitzari erretiratua"},
		"Cuarenta años sirviendo cañas en un bar de Triana donde se jugaba al mus " +
			"en la mesa del fondo. Aprendió mirando y ahora no hay quien le cierre la boca ni le " +
			"adivine las cartas. Dice que ganó un jamón en un campeonato del 83 y lo cuenta cada vez que gana un órdago.": {
			"Forty years pouring beers in a Triana bar where mus was played at the back table. " +
				"He learned by watching, and now there's no shutting him up or guessing his cards. " +
				"He swears he won a ham at a tournament back in '83, and tells the story every time he wins an órdago.",
			"Čtyřicet let točil pivo v baru v Trianě, kde se u zadního stolu hrál mus. " +
				"Naučil se to koukáním a teď ho nikdo neumlčí ani mu neodhadne karty. " +
				"Tvrdí, že v roce 83 vyhrál na turnaji šunku, a vypráví to pokaždé, když vyhraje órdago.",
			"Berrogei urte eman zituen garagardoak zerbitzatzen Trianako taberna batean, atzeko mahaian musean jokatzen zen tokian. " +
				"Begira ikasi zuen, eta orain ez dago ahoa isilaraziko dionik, ezta kartak asmatuko dizkionik ere. " +
				"83ko txapelketa batean urdaiazpiko bat irabazi zuela dio, eta órdago bat irabazten duen bakoitzean kontatzen du.",
		},
		"Farolero de libro. Envida con cualquier cosa y te lo dice con una sonrisa.": {
			"A textbook bluffer. He'll envido with anything and tell you so with a smile.",
			"Učebnicový blafér. Dá envido s čímkoli a ještě se u toho usmívá.",
			"Liburuetako faroleroa. Edozerekin egiten du envido, eta irribarrez esaten dizu.",
		},
		"¡Ozú, qué cartas me han dao!":                  {"Good Lord, what cards they've dealt me!", "Panebože, to jsou mi karty!", "Jesus, a zer kartak eman dizkidaten!"},
		"Esto lo gano yo con los ojos cerraos, quillo.": {"I'll win this one with my eyes shut, mate.", "Tohle vyhraju i se zavřenýma očima, kamaráde.", "Hau begiak itxita irabaziko dut, lagun."},
		"Ni me mires, que me pongo colorao.":            {"Don't even look at me, I'll go all red.", "Ani se na mě nedívej, nebo zrudnu.", "Ez begiratu ere egin, gorritu egingo naiz eta."},
		"Una vez gané un jamón con peores cartas.":      {"I once won a ham with worse cards than these.", "Jednou jsem s horšíma kartama vyhrál šunku.", "Behin urdaiazpiko bat irabazi nuen hauek baino karta okerragoekin."},
		"Arsa, que esto se anima.":                      {"Olé, now we're getting going!", "Olé, teď to začíná být zajímavé!", "Olé, hau berotzen ari da!"},
		"Tranquilo, que no llevo nada... o sí.":         {"Relax, I've got nothing... or have I?", "Klid, nic nemám... nebo jo?", "Lasai, ez daukat ezer... edo bai?"},

		// ── Doña Remedios ──
		"maestra jubilada": {"retired schoolteacher", "učitelka v důchodu", "irakasle erretiratua"},
		"Dio clase de matemáticas durante cuarenta años en un instituto de Valladolid. " +
			"Juega los martes en el hogar del jubilado y nadie recuerda la última vez que perdió una chica. " +
			"No farolea jamás: dice que mentir está feo, aunque sea al mus.": {
			"She taught maths for forty years at a secondary school in Valladolid. " +
				"She plays on Tuesdays at the pensioners' club, and nobody can remember the last time she lost a chica. " +
				"She never bluffs: lying is ugly, she says, even at mus.",
			"Čtyřicet let učila matematiku na gymnáziu ve Valladolidu. " +
				"Hraje v úterý v klubu důchodců a nikdo si nepamatuje, kdy naposledy prohrála chicu. " +
				"Nikdy neblafuje: říká, že lhát se nemá, ani při musu.",
			"Berrogei urtez matematika irakatsi zuen Valladolideko institutu batean. " +
				"Asteartetan jubilatuen etxean jokatzen du, eta inork ez du gogoratzen noiz galdu zuen azken aldiz chica bat. " +
				"Ez du inoiz farolik egiten: gezurra esatea itsusia dela dio, musean bada ere.",
		},
		"La reina de la chica. Prudentísima, casi nunca se tira un farol.": {
			"Queen of the chica. Extremely cautious, she almost never bluffs.",
			"Královna chicy. Nesmírně opatrná, skoro nikdy neblafuje.",
			"Chicaren erregina. Oso zuhurra, ia inoiz ez du farolik egiten.",
		},
		"A la chica, hijo, a la chica.":                            {"To the chica, dear, to the chica.", "Na chicu, chlapče, na chicu.", "Chicara, seme, chicara."},
		"Con paciencia se gana el cielo y el mus.":                 {"Patience wins you heaven, and mus too.", "Trpělivost růže přináší, i v musu.", "Pazientziaz irabazten dira zerua eta musa."},
		"Eso en mis tiempos no se hacía.":                          {"We didn't do that in my day.", "Za mých časů se tohle nedělalo.", "Nire garaian hori ez zen egiten."},
		"Los ases, como los buenos alumnos: callados y a lo suyo.": {"Aces are like good pupils: quiet and minding their own business.", "Esa jsou jako hodní žáci: potichu a dělají, co mají.", "Batekoak, ikasle onak bezala: isilik eta beren lanean."},
		"Mira que os lo tengo dicho.":                              {"How many times have I told you?", "Kolikrát jsem vám to už říkala?", "Zenbat aldiz esan behar dizuet?"},

		// ── Iñaki ──
		"ex pelotari": {"former pelota player", "bývalý pelotari", "pelotari ohia"},
		"Fue delantero de pelota a mano hasta que la muñeca dijo basta. Ahora juega al mus en " +
			"una sociedad gastronómica de Bilbao, donde el mus se toma tan en serio como el txakoli. " +
			"Va a la grande como iba a la pared: de frente y sin miedo.": {
			"He played up front in hand pelota until his wrist said enough. Now he plays mus at " +
				"a gastronomic society in Bilbao, where mus is taken as seriously as the txakoli. " +
				"He goes for the grande the way he went for the wall: head-on and fearless.",
			"Hrál baskickou pelotu holou rukou v přední řadě, dokud mu zápěstí neřeklo dost. Teď hraje mus " +
				"v gastronomickém spolku v Bilbau, kde se mus bere stejně vážně jako txakoli. " +
				"Na grande jde, jako chodil ke zdi: zpříma a beze strachu.",
			"Esku pilotako aurrelaria izan zen, eskumuturrak nahikoa esan zion arte. Orain " +
				"Bilboko elkarte gastronomiko batean jokatzen du musean, eta han musa txakolina bezain serio hartzen da. " +
				"Grandera paretara joaten zen bezala joaten da: aurrez aurre eta beldurrik gabe.",
		},
		"Grande y órdago. Cuando lleva reyes no se corta un pelo.": {
			"Grande and órdago. When he's holding kings he doesn't hold back.",
			"Grande a órdago. Když má krále, nebere si servítky.",
			"Grande eta órdago. Erregeak dituenean, ez du atzera egiten.",
		},
		"¡Aupa!":                                      {"Aupa!", "Aupa!", "Aupa!"},
		"Órdago y a otra cosa, mariposa.":             {"Órdago, and on to the next thing.", "Órdago a jede se dál.", "Órdago, eta hurrengora!"},
		"En Bilbao esto se juega así, de frente.":     {"In Bilbao we play it like this, straight on.", "V Bilbau se to hraje takhle, zpříma.", "Bilbon honela jokatzen da, aurrez aurre."},
		"Con reyes, a la grande, que para eso están.": {"With kings, go to the grande; that's what they're for.", "S králi na grande, od toho tu jsou.", "Erregeekin, grandera, horretarako daude eta."},
		"¡Venga, ahí va!":                             {"Come on, here goes!", "Tak jdem na to!", "Tira, hor doa!"},

		// ── Montse ──
		"actuaria de seguros": {"insurance actuary", "pojistná matematička", "aseguru-aktuarioa"},
		"Calcula riesgos para una aseguradora y aplica lo mismo en la mesa. Aprendió a jugar " +
			"en la facultad y desde entonces lleva una hoja de cálculo mental de cada partida. " +
			"Lo que más le gusta: los pares, porque ahí los números cuadran.": {
			"She calculates risks for an insurance company and does the same at the table. She learned to play " +
				"at university and has kept a mental spreadsheet of every game ever since. " +
				"Her favourite: the pares, because that's where the numbers add up.",
			"Počítá rizika pro pojišťovnu a totéž dělá u stolu. Hrát se naučila " +
				"na fakultě a od té doby si o každé partii vede tabulku v hlavě. " +
				"Nejraději má pares, protože tam čísla sedí.",
			"Arriskuak kalkulatzen ditu aseguru-etxe batentzat, eta gauza bera egiten du mahaian. Fakultatean " +
				"ikasi zuen jokatzen, eta harrezkero partida bakoitzaren kalkulu-orri bat darama buruan. " +
				"Gehien gustatzen zaiona: pares, han zenbakiak bat datozelako.",
		},
		"Fría y calculadora. Brilla en los pares y aprovecha cada seña que ve.": {
			"Cool and calculating. She shines at pares and makes the most of every seña she spots.",
			"Chladná a vypočítavá. Exceluje v pares a využije každou seña, kterou zahlédne.",
			"Hotza eta kalkulatzailea. Pares lancean distira egiten du, eta ikusten duen seña oro aprobetxatzen du.",
		},
		"Estadísticamente, no.":                    {"Statistically, no.", "Statisticky vzato ne.", "Estatistikoki, ez."},
		"Las probabilidades no mienten.":           {"The odds don't lie.", "Pravděpodobnost nelže.", "Probabilitateek ez dute gezurrik esaten."},
		"Endavant.":                                {"Endavant.", "Endavant.", "Endavant."},
		"Un 62 por ciento de que lo lleves. Paso.": {"62 per cent chance you've got it. Paso.", "Na 62 procent to máš. Paso.", "Ehuneko 62ko aukera daukazu. Paso."},
		"Esto ya lo tenía calculado.":              {"I'd already worked this out.", "Tohle jsem měla spočítané.", "Hau dagoeneko kalkulatuta neukan."},

		// ── Xosé ──
		"ganadero": {"cattle farmer", "chovatel dobytka", "abeltzaina"},
		"Tiene vacas rubias en una aldea cerca de Lugo y juega al mus en la feria de los jueves. " +
			"Nadie sabe nunca lo que lleva: si le preguntas si sube o baja la escalera, contesta que depende. " +
			"Sus señas son tan finas que ni su compañero está seguro de haberlas visto.": {
			"He keeps Galician blond cattle in a hamlet near Lugo and plays mus at the Thursday fair. " +
				"Nobody ever knows what he's holding: ask him if he's going up or down the stairs and he'll say it depends. " +
				"His señas are so subtle that even his partner isn't sure he's seen them.",
			"Chová galicijský skot rubia ve vísce u Luga a mus hraje na čtvrtečním trhu. " +
				"Nikdo nikdy neví, co má v ruce: zeptejte se ho, jestli jde po schodech nahoru, nebo dolů, a odpoví, že podle toho. " +
				"Jeho señas jsou tak jemné, že si ani jeho spoluhráč není jistý, jestli je viděl.",
			"Behi rubiak ditu Lugo ondoko herrixka batean, eta osteguneko azokan jokatzen du musean. " +
				"Inork ez daki inoiz zer daraman: eskaileretan gora ala behera doan galdetzen badiozu, araberakoa dela erantzungo dizu. " +
				"Bere señak hain dira finak, ezen bikotekidea bera ere ez baitago ziur ikusi dituenik.",
		},
		"Imposible de leer. Equilibrado en todo, y a él no le pillas una seña.": {
			"Impossible to read. Balanced in everything, and you'll never catch one of his señas.",
			"Nečitelný. Ve všem vyrovnaný a nechytíš mu ani jednu seña.",
			"Ezin da irakurri. Orekatua denetan, eta ez diozu seña bakar bat ere harrapatuko.",
		},
		"Depende.":                           {"It depends.", "Podle toho.", "Araberakoa."},
		"Pues según se mire...":              {"Well, depends how you look at it...", "No, jak se to vezme...", "Tira, nola begiratzen zaion..."},
		"Non é por nada, pero quero.":        {"Non é por nada, pero quero.", "Non é por nada, pero quero.", "Non é por nada, pero quero."},
		"Eu nin sí nin non.":                 {"Eu nin sí nin non.", "Eu nin sí nin non.", "Eu nin sí nin non."},
		"Isto vai ser que non... ou que si.": {"Isto vai ser que non... ou que si.", "Isto vai ser que non... ou que si.", "Isto vai ser que non... ou que si."},

		// ── Manolo ──
		"camionero": {"lorry driver", "řidič kamionu", "kamioilaria"},
		"Ha recorrido España entera en camión y ha jugado al mus en todas las áreas de servicio " +
			"entre Zaragoza y Algeciras. Es más tozudo que una mula: si le envidas, lo quiere. " +
			"Le han ganado muchas veces, pero nunca le han hecho callar.": {
			"He has driven his lorry the length and breadth of Spain and played mus at every service station " +
				"between Zaragoza and Algeciras. He's as stubborn as a mule: envido him and he'll say quiero. " +
				"He's been beaten plenty of times, but never shut up.",
			"Projel s kamionem celé Španělsko a mus hrál na všech odpočívadlech " +
				"mezi Zaragozou a Algecirasem. Je tvrdohlavější než mula: když mu dáš envido, řekne quiero. " +
				"Porazili ho mockrát, ale umlčet ho nedokázal nikdo.",
			"Espainia osoa zeharkatu du kamioian, eta Zaragoza eta Algeciras arteko " +
				"zerbitzugune guztietan jokatu du musean. Mandoa baino burugogorragoa da: envido egiten badiozu, quiero esango dizu. " +
				"Askotan irabazi diote, baina inoiz ez dute isilarazi.",
		},
		"Cabezota: casi nunca dice que no quiere. Valiente y algo temerario.": {
			"Pig-headed: he hardly ever says no quiero. Brave and a bit reckless.",
			"Paličák: skoro nikdy neřekne no quiero. Odvážný a trochu lehkomyslný.",
			"Burugogorra: ia inoiz ez du no quiero esaten. Ausarta eta pixka bat arduragabea.",
		},
		"¡Que sí, que lo quiero, cucha!":                  {"Yes, I said quiero, listen!", "Ale jo, quiero, slyšíš!", "Baietz, quiero, entzun!"},
		"Más tozudo que un maño, dicen. Y a mucha honra.": {"Stubborn as an Aragonese, they say. And proud of it.", "Prý jsem tvrdohlavý jako Aragonec. A jsem na to hrdý.", "Aragoiarra bezain burugogorra, diotenez. Eta harro nago."},
		"A mí no me achanta nadie.":                       {"Nobody scares me off.", "Mě nikdo nezastraší.", "Ni inork ez nau ikaratzen."},
		"¡Pilarica, échame una mano!":                     {"Our Lady of the Pillar, give me a hand!", "Panenko Maria Pilarská, pomoz mi!", "Pilarreko Ama, lagun iezadazu!"},
		"Lo quiero y punto.":                              {"Quiero, end of story.", "Quiero, a basta.", "Quiero, eta kitto."},

		// ── Carmen ──
		"guardia civil retirada": {"retired Civil Guard officer", "příslušnice Guardia Civil ve výslužbě", "guardia zibil erretiratua"},
		"Treinta años en la Benemérita le dejaron un ojo que no se le escapa nada. Juega en un " +
			"bar de Lavapiés y los habituales dicen que es mejor no pasar señas cuando ella mira. " +
			"Al juego, la 31 parece perseguirla.": {
			"Thirty years in the Civil Guard left her with an eye that misses nothing. She plays in a " +
				"bar in Lavapiés, and the regulars say you'd better not pass señas while she's watching. " +
				"At juego, la 31 seems to follow her around.",
			"Třicet let u Guardia Civil jí zanechalo oko, kterému nic neunikne. Hraje v " +
				"baru v Lavapiés a štamgasti říkají, že když se ona dívá, je lepší nedávat señas. " +
				"Při juego jako by ji la 31 pronásledovala.",
			"Guardia Zibilean emandako hogeita hamar urteek ezer ihes egiten ez dion begia utzi zioten. " +
				"Lavapiéseko taberna batean jokatzen du, eta ohikoek diote hobe dela señarik ez egitea bera begira dagoenean. " +
				"Juegoan, la 31 atzetik dabilkiola dirudi.",
		},
		"Especialista en el juego y la que más señas pilla. Agresiva.": {
			"A juego specialist and the best at catching señas. Aggressive.",
			"Specialistka na juego, která chytí nejvíc señas. Agresivní.",
			"Juegoan aditua, eta señak gehien harrapatzen dituena. Erasokorra.",
		},
		"Aquí mando yo.":                           {"I'm in charge here.", "Tady velím já.", "Hemen nik agintzen dut."},
		"A sus órdenes... y órdago.":               {"Yes, sir... and órdago.", "Rozkaz... a órdago.", "Zure aginduetara... eta órdago."},
		"Te he visto, majo.":                       {"I saw that, sunshine.", "Viděla jsem tě, zlatíčko.", "Ikusi zaitut, motel."},
		"Esa seña la he visto hasta yo sin gafas.": {"Even I caught that seña without my glasses.", "Tu seña jsem viděla i bez brýlí.", "Seña hori betaurrekorik gabe ere ikusi dut."},
		"Circulen, que aquí no hay nada que ver.":  {"Move along, nothing to see here.", "Rozejděte se, tady není nic k vidění.", "Aurrera, hemen ez dago ezer ikusteko."},

		// ── Pepe ──
		"estudiante de Magisterio": {"trainee teacher", "student učitelství", "Magisteritzako ikaslea"},
		"Aprendió a jugar el verano pasado con su abuelo en el pueblo y se ha venido arriba. " +
			"Confunde a veces la grande con la chica, pero tiene una suerte que da rabia. " +
			"Ideal para empezar, o para que te toque de rival.": {
			"He learned to play last summer with his grandad in the village and now he's on a roll. " +
				"He sometimes mixes up the grande and the chica, but he's annoyingly lucky. " +
				"Ideal to start with, or to get as an opponent.",
			"Hrát se naučil loni v létě s dědou na vesnici a teď je v ráži. " +
				"Občas si plete grande s chicou, ale má štěstí, až to štve. " +
				"Ideální na začátek, nebo jako soupeř.",
			"Joan den udan ikasi zuen jokatzen aitonarekin herrian, eta gora etorri da. " +
				"Batzuetan grandea eta chica nahasten ditu, baina amorrua ematen duen zortea du. " +
				"Hasteko aproposa, edo aurkari gisa egokitzeko.",
		},
		"Principiante. Se da mus con todo y valora mal las jugadas.": {
			"Beginner. He says mus with anything and misjudges his hands.",
			"Začátečník. Dává mus se vším a špatně odhaduje karty.",
			"Hasiberria. Edozerekin ematen du mus, eta jokaldiak gaizki baloratzen ditu.",
		},
		"¿Esto es bueno? Digo... ¡envido!":       {"Is this good? I mean... envido!", "Je to dobrý? Teda... envido!", "Ona da hau? Esan nahi dut... envido!"},
		"Mi abuelo me dijo que siempre a reyes.": {"My grandad told me: always go for kings.", "Děda říkal, že vždycky na krále.", "Aitonak esan zidan beti erregeetara."},
		"Ay, que no sé si quiero.":               {"Oh, I don't know if it's quiero.", "Jéje, nevím, jestli quiero.", "Ai, ez dakit quiero den."},
		"¿Los treses eran reyes o no?":           {"Threes count as kings, right? Or not?", "Trojky byly králové, nebo ne?", "Hirukoak erregeak ziren, ala ez?"},
		"¡Toma ya!":                              {"Take that!", "Tumáš!", "Hartu hori!"},

		// ── Toni ──
		"pirotécnico fallero": {"Fallas pyrotechnician", "pyrotechnik na Fallas", "Fallasetako piroteknikoa"},
		"Monta mascletàs en Fallas y juega al mus igual: mucho ruido y que no se sepa por dónde " +
			"va a salir. Se marca faroles tremendos al punto y luego se ríe a carcajadas.": {
			"He sets up mascletàs during Fallas and plays mus the same way: lots of noise and nobody knows where " +
				"it'll go off. He pulls outrageous bluffs at punto and then roars with laughter.",
			"O Fallas staví mascletà a mus hraje stejně: spousta randálu a nikdo neví, kde " +
				"to bouchne. Při punto si dává šílené blafy a pak se hlasitě chechtá.",
			"Fallasetan mascletàk muntatzen ditu, eta musean berdin jokatzen du: zarata handia, eta inork ez daki nondik " +
				"aterako den. Punto lancean farol ikaragarriak egiten ditu, eta gero barre algaraka hasten da.",
		},
		"Explosivo y farolero. Muy bueno al punto.": {
			"Explosive and a bluffer. Very good at punto.",
			"Výbušný blafér. Velmi dobrý v punto.",
			"Leherkorra eta faroleroa. Oso ona punto lancean.",
		},
		"¡Che, qué mano!": {"Che, what a hand!", "Che, to je ruka!", "Che, hau eskua!"},
		"Esto es más fácil que hacer una paella... bueno, no.": {"This is easier than making a paella... well, no.", "Tohle je snazší než uvařit paellu... no, vlastně ne.", "Hau paella bat egitea baino errazagoa da... beno, ez."},
		"¡Pum! Como una mascletà.":                             {"Boom! Like a mascletà.", "Bum! Jako mascletà.", "Danba! Mascletà bat bezala."},
		"Al punto no me gana ni mi suegra.":                    {"Not even my mother-in-law beats me at punto.", "V punto mě neporazí ani tchyně.", "Punto lancean amaginarrebak ere ez dit irabazten."},
		"¡Vinga, va!":                                          {"Vinga, va!", "Vinga, va!", "Vinga, va!"},

		// ── Yeray ──
		"profesor de surf": {"surf instructor", "instruktor surfování", "surf-irakaslea"},
		"Entre ola y ola juega al mus en un chiringuito de Las Canteras. No hay nada que le altere: " +
			"se da mus con medio mazo y espera tranquilo a que lleguen las cartas, como espera las olas.": {
			"Between waves he plays mus at a beach bar on Las Canteras. Nothing ruffles him: " +
				"he says mus with half the deck and calmly waits for the cards to come, the way he waits for the waves.",
			"Mezi vlnami hraje mus v plážovém baru na Las Canteras. Nic ho nerozhází: " +
				"dává mus s půlkou balíčku a v klidu čeká, až karty přijdou, jako čeká na vlny.",
			"Olatu batetik bestera musean jokatzen du Las Canteraseko hondartza-taberna batean. Ezerk ez du asaldatzen: " +
				"sortaren erdiarekin ematen du mus, eta lasai itxaroten die kartei, olatuei itxaroten dien bezala.",
		},
		"Tranquilo y musero. Sin prisa, pero con buen ojo para el juego.": {
			"Laid-back and quick to say mus. No hurry, but a good eye for juego.",
			"Klidný a rád dává mus. Bez spěchu, ale s dobrým okem na juego.",
			"Lasaia eta museroa. Presarik gabe, baina juegorako begi ona du.",
		},
		"Tranqui, mi niño.":                          {"Easy, my lad.", "V klidu, kámo.", "Lasai, motel."},
		"Sin estrés, que el mus es pa' disfrutarlo.": {"No stress, mus is for enjoying.", "Bez stresu, mus je od toho, aby se užíval.", "Estresik gabe, musa gozatzeko da eta."},
		"Ya vendrá la ola buena.":                    {"The good wave will come.", "Však ona přijde ta pravá vlna.", "Etorriko da olatu ona."},
		"¡Fuerte jugada, chacho!":                    {"What a play, mate!", "To je tah, kámo!", "Hau jokaldia, motel!"},
		"Ños, qué cartas.":                           {"Blimey, what cards.", "Teda, to jsou karty.", "Ene, hau kartak."},

		// ── Miren ──
		"cocinera de sociedad gastronómica": {"cook at a gastronomic society", "kuchařka v gastronomickém spolku", "elkarte gastronomikoko sukaldaria"},
		"Cocinó durante treinta años para los socios de un txoko de Azpeitia, donde las mujeres no podían ni sentarse a la mesa. " +
			"Mientras removía el bacalao al pil-pil aprendió todas las señas de todos los socios, y el día que por fin la dejaron jugar " +
			"les ganó tres partidas seguidas. La llaman Sorgina, la bruja, porque dicen que te lee las cartas en la cara.": {
			"For thirty years she cooked for the members of a txoko in Azpeitia, where women weren't even allowed to sit at the table. " +
				"While stirring the bacalao al pil-pil she learned every seña of every member, and the day they finally let her play " +
				"she beat them three games in a row. They call her Sorgina, the witch, because they say she reads your cards in your face.",
			"Třicet let vařila členům txoka v Azpeitii, kde ženy nesměly ani usednout ke stolu. " +
				"Zatímco míchala tresku pil-pil, naučila se všechny señas všech členů, a v den, kdy ji konečně pustili hrát, " +
				"je porazila třikrát za sebou. Říkají jí Sorgina, čarodějnice, protože prý vám karty přečte z obličeje.",
			"Hogeita hamar urtez Azpeitiko txoko bateko bazkideentzat sukaldatu zuen, emakumeak mahaian esertzerik ere ez zuten garaian. " +
				"Bakailaoa pil-pilean eragiten zuen bitartean bazkide guztien seña guztiak ikasi zituen, eta azkenean jokatzen utzi zioten egunean " +
				"hiru partida jarraian irabazi zizkien. Sorgina deitzen diote, kartak aurpegian irakurtzen dizkizula diotelako.",
		},
		"Zorra vieja: pilla las señas de todos, no deja ver las suyas y farolea cuando menos te lo esperas.": {
			"A wily old fox: she catches everyone's señas, never gives hers away, and bluffs when you least expect it.",
			"Stará liška: chytí señas všech, svoje nikdy neprozradí a blafuje, když to nejmíň čekáš.",
			"Azeri zaharra: denen señak harrapatzen ditu, bereak ez ditu erakusten, eta gutxien espero duzunean egiten du farol.",
		},
		"Ene, ene... qué cartas más bonitas.":         {"Ene, ene... what pretty cards.", "Ene, ene... to jsou ale hezké karty.", "Ene, ene... zer karta politak."},
		"Yo no he visto nada, ¿eh? Nada.":             {"I haven't seen a thing, eh? Not a thing.", "Já nic neviděla, jo? Vůbec nic.", "Nik ez dut ezer ikusi, e? Ezer ez."},
		"Isilik, que el mus se juega callado.":        {"Isilik: mus is played in silence.", "Isilik, mus se hraje potichu.", "Isilik, musean isilik jokatzen da eta."},
		"Ondo, ondo. Muy bien, maitia.":               {"Ondo, ondo. Very good, maitia.", "Ondo, ondo. Moc dobře, maitia.", "Ondo, ondo. Oso ondo, maitia."},
		"Esa seña la hacía mejor mi difunto Joxe.":    {"My late Joxe did that seña better.", "Tu seña dělal můj nebožtík Joxe líp.", "Seña hori hobeto egiten zuen nire Joxe zenak."},
		"Con paciencia y un buen pil-pil, todo sale.": {"With patience and a good pil-pil, everything comes together.", "S trpělivostí a pořádným pil-pil se všechno povede.", "Pazientziaz eta pil-pil on batekin, dena ateratzen da."},

		// ── Ane ──
		"remera de trainera": {"trainera rower", "veslařka na traineře", "traineruko arraunlaria"},
		"Rema en una trainera de la Parte Vieja y entrena cada mañana en la bahía de la Concha. " +
			"Aprendió a jugar en los bares de pintxos con los veteranos del club, que la llaman Txiki porque es la más pequeña " +
			"de la tripulación... y la que más grita. Juega como rema: a ritmo, sin parar y sin mirar atrás.": {
			"She rows in a trainera from the Old Town and trains every morning in La Concha bay. " +
				"She learned to play in the pintxo bars with the club's veterans, who call her Txiki because she's the smallest " +
				"in the crew... and the loudest. She plays the way she rows: in rhythm, non-stop and never looking back.",
			"Vesluje na traineře ze Starého města a každé ráno trénuje v zátoce La Concha. " +
				"Hrát se naučila v pintxo barech s veterány klubu, kteří jí říkají Txiki, protože je nejmenší " +
				"z posádky... a nejhlasitější. Hraje, jak vesluje: v rytmu, bez přestávky a bez ohlížení.",
			"Alde Zaharreko trainera batean egiten du arraun, eta goizero Kontxako badian entrenatzen da. " +
				"Pintxo-tabernetan ikasi zuen jokatzen, klubeko beteranoekin; Txiki deitzen diote, tripulazioko txikiena " +
				"delako... eta ozenen oihu egiten duena. Arraunean bezala jokatzen du: erritmoan, gelditu gabe eta atzera begiratu gabe.",
		},
		"Corta el mus a la primera y aprieta en pares y juego. Valiente, pero se le notan las señas.": {
			"She cuts the mus at the first chance and pushes hard at pares and juego. Brave, but her señas show.",
			"Mus utne hned napoprvé a tlačí v pares a juego. Odvážná, ale señas jsou na ní vidět.",
			"Lehen aldian mozten du musa, eta gogor estutzen du pares eta juegoan. Ausarta, baina señak nabaritzen zaizkio.",
		},
		"¡Aupa ahí!":                           {"Aupa, there you go!", "Aupa, tak do toho!", "Aupa hor!"},
		"Kaixo, ¿jugamos o qué?":               {"Kaixo, are we playing or what?", "Kaixo, hrajeme, nebo co?", "Kaixo, jokatuko dugu ala zer?"},
		"No hay mus, que se enfría el pintxo.": {"No hay mus, my pintxo's getting cold.", "No hay mus, ať mi nevystydne pintxo.", "No hay mus, pintxoa hozten ari da eta."},
		"¡Ondo! Como en la Kontxa.":            {"Ondo! Just like at La Kontxa.", "Ondo! Jako v Kontxe.", "Ondo! Kontxan bezala."},
		"¡Hemen gaude!":                        {"Hemen gaude!", "Hemen gaude!", "Hemen gaude!"},
		"Agur, y gracias por los tantos.":      {"Agur, and thanks for the tantos.", "Agur, a díky za tantos.", "Agur, eta eskerrik asko tantoengatik."},

		// ── Refranero (paquete frases) ──
		"Mus.": {"Mus.", "Mus.", "Mus."},
		"Mus, que estas no valen ni para el café.": {"Mus, these aren't worth a cup of coffee.", "Mus, tyhle nestojí ani za kafe.", "Mus, hauek ez dute kafe baterako ere balio."},
		"Mus hasta que salgan los reyes.":          {"Mus until the kings turn up.", "Mus, dokud nepřijdou králové.", "Mus, erregeak atera arte."},
		"Dame mus, compañero.":                     {"Give me mus, partner.", "Dej mi mus, parťáku.", "Emadazu mus, bikote."},
		"Mus corrido.":                             {"Mus corrido.", "Mus corrido.", "Mus corrido."},
		"Corta con buenas, compañero.":             {"Cut it if you've got good ones, partner.", "Utni to, jestli máš dobré, parťáku.", "Moztu onak badituzu, bikote."},
		"Con esto no se va ni a la esquina. Mus.":  {"You couldn't get round the corner with these. Mus.", "S tímhle nedojdu ani za roh. Mus.", "Honekin ez zoaz ezta izkinaraino ere. Mus."},

		"No hay mus.":                   {"No hay mus.", "No hay mus.", "No hay mus."},
		"Corto.":                        {"Corto.", "Corto.", "Corto."},
		"Quita, que llevo cosas.":       {"Out of the way, I've got something.", "Uhni, něco mám.", "Kendu, gauzak ditut eta."},
		"Aquí se corta.":                {"This is where it gets cut.", "Tady se utíná.", "Hemen mozten da."},
		"Hasta aquí ha llegado el mus.": {"That's the end of the mus.", "Tady mus končí.", "Hona arte iritsi da musa."},

		"A reyes.":            {"Going for kings.", "Na krále.", "Erregeetara."},
		"Dame de las buenas.": {"Give me some good ones.", "Dej mi nějaké dobré.", "Emadazu onetakoak."},
		"Esta vez sí.":        {"This time for sure.", "Tentokrát jo.", "Oraingoan bai."},
		"Que venga la 31.":    {"Come on, la 31.", "Ať přijde la 31.", "Etor dadila la 31."},

		"Paso.":                   {"Paso.", "Paso.", "Paso."},
		"Paso, por mi parte.":     {"Paso, as far as I'm concerned.", "Paso, za mě.", "Paso, nire aldetik."},
		"Paso, a ver qué hacéis.": {"Paso, let's see what you do.", "Paso, uvidíme, co uděláte vy.", "Paso, ea zer egiten duzuen."},
		"Llegó a mí... paso.":     {"My turn... paso.", "Jsem na řadě... paso.", "Niri iritsi zait... paso."},

		"Envido.":                          {"Envido.", "Envido.", "Envido."},
		"Ahí van dos.":                     {"Two, there you go.", "Tady máte dva.", "Hor doaz bi."},
		"Dos más, por no callar.":          {"Two more, just for something to say.", "O dva víc, ať jen nemlčím.", "Beste bi, isilik ez egoteagatik."},
		"Envido, que hay que animar esto.": {"Envido, this needs livening up.", "Envido, chce to oživit.", "Envido, hau alaitu behar da eta."},
		"El que no envida no gana.":        {"If you don't envido, you don't win.", "Kdo nedá envido, nevyhraje.", "Envido egiten ez duenak ez du irabazten."},
		"Las de Hontanares.":               {"The Hontanares special.", "Po hontanaresku.", "Hontanareskoak."},

		"Quiero.":         {"Quiero.", "Quiero.", "Quiero."},
		"Lo quiero.":      {"Lo quiero.", "Lo quiero.", "Lo quiero."},
		"Veo.":            {"I'll see that.", "Vidím.", "Ikusten dut."},
		"Queremos.":       {"Queremos.", "Queremos.", "Queremos."},
		"Con lo que sea.": {"Whatever it takes.", "Ať je to cokoli.", "Dena delakoarekin."},

		"No quiero.":                          {"No quiero.", "No quiero.", "No quiero."},
		"Para vosotros.":                      {"All yours.", "Je to vaše.", "Zuentzat."},
		"Llevaos el tanto.":                   {"Take the tanto.", "Vezměte si tanto.", "Eraman tantoa."},
		"No queremos, más se perdió en Cuba.": {"No queremos; worse things happen at sea.", "No queremos, to není konec světa.", "No queremos, ez da munduaren amaiera."},

		"¡Órdago, y que salga el sol por Antequera!":       {"Órdago, and come what may!", "Órdago, a ať se děje, co se děje!", "Órdago, eta gertatzen dena gertatzen dela!"},
		"¡Ahí va el órdago!":                               {"Here comes the órdago!", "A je tu órdago!", "Hor doa órdagoa!"},
		"¡Órdago! A ver quién tiene lo que hay que tener.": {"Órdago! Let's see who's got the nerve.", "Órdago! Uvidíme, kdo na to má.", "Órdago! Ea nork duen behar den adorea."},
		"Se acabó lo que se daba: ¡órdago!":                {"That's all, folks: órdago!", "A dost bylo: órdago!", "Bukatu da festa: órdago!"},

		"¿Os rendís?":                                       {"Do you give up?", "Vzdáváte se?", "Amore ematen duzue?"},
		"Esto es mus de los de antes.":                      {"This is old-school mus.", "Tohle je mus jako za starých časů.", "Hau garai bateko musa da."},
		"La mano es la mano.":                               {"Mano is mano.", "Mano je mano.", "Mano mano da."},
		"A estas cartas no les gana ni el cura del pueblo.": {"Not even the village priest could beat these cards.", "Na tyhle karty by nevyzrál ani vesnický farář.", "Karta hauei ez die herriko apaizak ere irabaziko."},
		"Manitas de plata, ¿eh?":                            {"Golden hands, eh?", "Zlaté ručičky, co?", "Urrezko eskuak, e?"},

		"Manitas de plata...":            {"Golden hands...", "Zlaté ručičky...", "Urrezko eskuak..."},
		"Vaya cartitas os han tocado.":   {"Nice little cards you got there.", "Pěkné kartičky vám přišly.", "Hara zer kartatxo egokitu zaizkizuen."},
		"Ya vendrán tiempos mejores.":    {"Better days will come.", "Však přijdou lepší časy.", "Etorriko dira garai hobeak."},
		"Más se perdió en Cuba.":         {"Worse things happen at sea.", "Horší věci se staly.", "Ez da munduaren amaiera."},
		"No hay mal que cien años dure.": {"Nothing bad lasts forever.", "Všechno zlé jednou skončí.", "Ez dago ehun urte irauten duen gaitzik."},
	})
}
