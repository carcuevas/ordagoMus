package i18n

// Textos del motor (internal/game), de las reglas (internal/rules), de la
// baraja (internal/cards) y de los nombres de algunos ajustes.
// Orden: inglés, checo, euskera.
func init() {
	registrar(map[string]trad{
		// game: errores
		"no es tu turno":                            {"it's not your turn", "nejsi na tahu", "ez da zure txanda"},
		"acción no permitida ahora":                 {"that action is not allowed now", "tento tah teď není povolen", "ekintza hori ez dago baimenduta orain"},
		"ahora no se pueden pasar señas":            {"señas can't be passed now", "teď se nedají předávat señas", "orain ezin da señarik pasatu"},
		"no llevas eso: con las señas no se miente": {"you don't have that: you never lie with señas", "to nemáš: se señas se nelže", "ez daukazu hori: señekin ez da gezurrik esaten"},
		"hay que pedir entre 1 y 4 cartas":          {"you must ask for between 1 and 4 cards", "musíš si říct o 1 až 4 karty", "1 eta 4 karta artean eskatu behar dira"},
		"descarte inválido":                         {"invalid discard", "neplatné odhození", "baztertze baliogabea"},
		"el envite mínimo es de 2":                  {"the minimum envite is 2", "nejmenší envite je 2", "gutxieneko envitea 2 da"},
		"hay que subir al menos 1":                  {"you must raise by at least 1", "musíš přihodit aspoň 1", "gutxienez 1 igo behar da"},

		// game: mesa
		"Se sortea la mano: la carta menor es para %s.": {
			"Drawing for mano: the lowest card goes to %s.",
			"Losuje se o mano: nejnižší kartu má %s.",
			"Mano zozketatzen da: karta txikiena %s jokalariak atera du."},
		"Reparte %s. Es mano %s.": {"%s deals. %s is mano.", "Rozdává %s. Mano je %s.", "%s jokalariak banatzen ditu kartak. Mano: %s."},
		"Primera mano: mus corrido y sin señas.": {
			"First hand: mus corrido and no señas.",
			"První rozdání: mus corrido a bez señas.",
			"Lehen eskua: mus corrido eta señarik gabe."},
		"Se acaba el mazo: se barajan los descartes.": {
			"The deck has run out: the discards are shuffled.",
			"Došel balíček: míchají se odhozené karty.",
			"Sorta amaitu da: baztertutako kartak nahasten dira."},
		"Ya puedes pasar tus señas.": {"You can pass your señas now.", "Teď můžeš předat své señas.", "Orain zure señak pasa ditzakezu."},
		"¡Seña pillada! %s":          {"Seña spotted! %s", "Prokouknutá seña! %s", "Seña harrapatuta! %s"},
		"¡%s te ha pillado la seña! (%s)": {
			"%s spotted your seña! (%s)",
			"Pozor, %s vidí tvou seña! (%s)",
			"%s jokalariak zure seña harrapatu du! (%s)"},
		"¡%s le ha pillado la seña a %s! (%s)": {
			"%s spotted %s's seña! (%s)",
			"Pozor, %s vidí seña hráče %s! (%s)",
			"%s jokalariak %s jokalariaren seña harrapatu du! (%s)"},
		"Haces la seña: %s (%s).":        {"Seña made: %s (%s).", "Seña: %s (%s).", "Seña egin duzu: %s (%s)."},
		"✓ %s te ha visto la seña (%s).": {"✓ %s saw your seña (%s).", "✓ %s tvou seña vidí (%s).", "✓ %s jokalariak zure seña ikusi du (%s)."},
		"%s no pasa seña.":               {"%s passes no seña.", "%s nepředává žádnou seña.", "%s jokalariak ez du señarik pasatzen."},
		"%s %s":                          {"%s %s", "%s %s", "%s jokalariak %s"},
		"Mus corrido. A descartarse.":    {"Mus corrido. Time to discard.", "Mus corrido. Odhazuje se.", "Mus corrido. Kartak baztertzera."},
		"%d carta":                       {"%d card", "%d karta", "%d karta"},
		"%d cartas":                      {"%d cards", "%d karty", "%d karta"},
		"Nadie lleva pares.":             {"Nobody has pares.", "Nikdo nemá pares.", "Inork ez du paresik."},
		"Nadie lleva juego: se juega al punto.": {
			"Nobody has juego: punto is played.",
			"Nikdo nemá juego: hraje se punto.",
			"Inork ez du juegorik: punto jokatuko da."},
		"Solo llevan pares %s: no hay envite.": {
			"Only %s have pares: no envite.",
			"Pares mají jen %s: envite nebude.",
			"Pares %s bikoteak bakarrik ditu: ez dago enviterik."},
		"Solo llevan juego %s: no hay envite.": {
			"Only %s have juego: no envite.",
			"Juego mají jen %s: envite nebude.",
			"Juego %s bikoteak bakarrik du: ez dago enviterik."},
		"%s: en paso.":  {"%s: everyone passes.", "%s: všichni dali paso.", "%s: guztiek paso."},
		"%s no querida": {"%s not accepted", "%s nepřijato", "%s, onartu gabe"},
		"%d más":        {"%d more", "o %d víc", "%d gehiago"},
		"¡Órdago querido! Se enseñan las cartas.": {
			"Órdago accepted! The cards are shown.",
			"Órdago přijato! Ukazují se karty.",
			"Órdagoa onartuta! Kartak erakusten dira."},
		"%s gana la %s con %s: el juego es para %s.": {
			"%s wins %s with %s: the juego goes to %s.",
			"%s vyhrává %s (%s): juego získávají %s.",
			"%s jokalariak irabazi du %s (%s): juegoa %s bikotearentzat da."},
		"%s: %d para %s (%d).":     {"%s: %d to %s (%d).", "%s: %d pro %s (%d).", "%s: %d %s bikotearentzat (%d)."},
		"Se enseñan las cartas.":   {"The cards are shown.", "Ukazují se karty.", "Kartak erakusten dira."},
		"%2d para %s (%s)":         {"%2d to %s (%s)", "%2d pro %s (%s)", "%2d %s bikotearentzat (%s)"},
		"%s en paso":               {"%s, everyone passed", "%s, všichni dali paso", "%s, guztiek paso"},
		"%s, envite querido de %d": {"%s, envite of %d accepted", "%s, přijaté envite za %d", "%s, envitea onartuta (%d)"},
		"tanto del punto":          {"tanto for punto", "tanto za punto", "punto lanceko tantoa"},
		"valor de %s":              {"value of %s", "hodnota za %s", "%s balioa"},
		"%s sin envite":            {"%s without envite", "%s bez envite", "%s, envite gabe"},
		"¡Juego para %s!":          {"Juego for %s!", "Juego pro %s!", "Juegoa %s bikotearentzat!"},
		"¡Vaca para %s!":           {"Vaca for %s!", "Vaca pro %s!", "Vaca %s bikotearentzat!"},
		"¡%s ganan la partida!":    {"%s win the match!", "%s vyhrávají partii!", "%s bikoteak irabazi du partida!"},
		"Nuevo juego.":             {"New juego.", "Nový juego.", "Juego berria."},

		// rules: cartas en plural (checo en genitivo, euskera en la forma
		// que va delante del nombre)
		"ases":     {"aces", "es", "bateko"},
		"doses":    {"twos", "dvojek", "biko"},
		"treses":   {"threes", "trojek", "hiruko"},
		"cuatros":  {"fours", "čtyřek", "lauko"},
		"cincos":   {"fives", "pětek", "bosteko"},
		"seises":   {"sixes", "šestek", "seiko"},
		"sietes":   {"sevens", "sedmiček", "zazpiko"},
		"sotas":    {"jacks", "spodků", "sota"},
		"caballos": {"knights", "koní", "zaldi"},
		"reyes":    {"kings", "králů", "errege"},

		"pareja de %s":                  {"pareja of %s", "pareja %s", "%s pareja"},
		"medias de %s":                  {"medias of %s", "medias %s", "%s medias"},
		"duples de %s (cuatro iguales)": {"duples of %s (four of a kind)", "duples %s (čtyři stejné)", "%s duples (lau berdin)"},
		"duples de %s y %s":             {"duples of %s and %s", "duples %s a %s", "%s eta %s duples"},
		"sin pares":                     {"no pares", "bez pares", "paresik gabe"},
		"%d al punto":                   {"%d for punto", "punto %d", "punto: %d"},
		"juego de %d":                   {"juego of %d", "juego %d", "%d-ko juegoa"},
		"¡Cuatro reyes, la piara!": {
			"Four kings: la piara (the herd)!",
			"Čtyři králové: la piara (stádo)!",
			"Lau errege: la piara (txerri-taldea)!"},
		"¡Solomillo! Tres reyes y la 31.": {
			"Solomillo (sirloin)! Three kings and la 31.",
			"Solomillo (svíčková)! Tři králové a la 31.",
			"Solomillo (azpizuna)! Hiru errege eta la 31."},
		"¡Besugo! Tres ases y un rey.": {
			"Besugo (sea bream)! Three aces and a king.",
			"Besugo (pražma)! Tři esa a král.",
			"Besugo (bisigua)! Hiru bateko eta errege bat."},

		// rules: señas (tercera persona, segunda persona y significado)
		"se muerde el labio inferior":      {"bites their lower lip", "se kousne do spodního rtu", "beheko ezpainari hozka egiten dio"},
		"saca la punta de la lengua":       {"sticks out the tip of their tongue", "vyplázne špičku jazyka", "mihiaren punta ateratzen du"},
		"tuerce la comisura de los labios": {"twists the corner of their mouth", "zkřiví koutek úst", "ahoaren ertza okertzen du"},
		"saca la lengua hacia un lado":     {"sticks their tongue out to one side", "vyplázne jazyk do strany", "mihia alde batera ateratzen du"},
		"levanta las cejas":                {"raises their eyebrows", "zvedne obočí", "bekainak altxatzen ditu"},
		"guiña un ojo":                     {"winks", "mrkne", "begi-keinua egiten du"},
		"cierra los dos ojos":              {"closes both eyes", "zavře obě oči", "bi begiak ixten ditu"},

		"te muerdes el labio inferior":      {"you bite your lower lip", "koušeš se do spodního rtu", "beheko ezpainari hozka egiten diozu"},
		"sacas la punta de la lengua":       {"you stick out the tip of your tongue", "vyplazuješ špičku jazyka", "mihiaren punta ateratzen duzu"},
		"tuerces la comisura de los labios": {"you twist the corner of your mouth", "křivíš koutek úst", "ahoaren ertza okertzen duzu"},
		"sacas la lengua hacia un lado":     {"you stick your tongue out to one side", "vyplazuješ jazyk do strany", "mihia alde batera ateratzen duzu"},
		"levantas las cejas":                {"you raise your eyebrows", "zvedáš obočí", "bekainak altxatzen dituzu"},
		"guiñas un ojo":                     {"you wink", "mrkáš", "begi-keinua egiten duzu"},
		"cierras los dos ojos":              {"you close both eyes", "zavíráš obě oči", "bi begiak ixten dituzu"},

		"dos reyes":       {"two kings", "dva králové", "bi errege"},
		"dos ases":        {"two aces", "dvě esa", "bi bateko"},
		"medias de reyes": {"medias of kings", "medias králů", "errege medias"},
		"medias de ases":  {"medias of aces", "medias es", "bateko medias"},
		"ciego: sin juego y sin pares que valgan": {
			"ciego: no juego and no pares worth anything",
			"ciego: bez juego a bez pares, které by za něco stály",
			"ciego: juegorik gabe eta ezer balio duen paresik gabe"},

		// Ya están en cat_tui.go (se comparten): "%s y %s", "A %s" y los
		// nombres de Velocidad y ModoSenas de internal/ajustes.

		// cards
		"oros":     {"coins", "mince", "urreak"},
		"copas":    {"cups", "poháry", "kopak"},
		"espadas":  {"swords", "meče", "ezpatak"},
		"bastos":   {"clubs", "kyje", "bastoiak"},
		"As":       {"Ace", "Eso", "Batekoa"},
		"Sota":     {"Jack", "Spodek", "Sota"},
		"Caballo":  {"Knight", "Kůň", "Zaldia"},
		"Rey":      {"King", "Král", "Erregea"},
		"Cab":      {"Kn", "Kůň", "Zal"},
		"%s de %s": {"%s of %s", "%s (%s)", "%s (%s)"},
	})
}
