# Super-Tanks (Doc)

Projektarbete på kursen Operativsystem och processorienterad programmering
(1DT096) våren 2019, Uppsala universitet.

Super-Tanks är ett online 2D multiplayerspel som är baserat på tidigare spel
som Pocket Tanks och Tanks. Likt tidigare spel så kontrollerar varje spelare
en stridsvagn som kan röra sig i sidled samt skjuta projektiler som rör sig i
kastparabler. Om en spelares stridsvagn blir träffad av en projektil så tappar
stridsvagnen ett liv av tre. Om projektilen bara träffar marken så bildas ett hål
i terrängen.

Olikt dess föregångare så kommer Super-Tanks spelas i realtid istället för att
använda ett turbaserat system. Detta blir en utmaning i att utveckla nya banor
och projektiler som får spelet att kännas balanserat även i realtid.

Det här spelet behövs eftersom liknande spel körs via Adobe Flash vilket slu-
tar stödjas av Adobe i slutet av 2020. Målet är att Super-Tanks ska vara en
vidareutveckling av dess föregångare, med bättre grafik och modernare spelme-
kaniker, och även mer lättillgängligt i form av open source utveckling.

## Bygga projektet

Kompilering/körning kommandon
* "make run"................för att bygga och köra programmet.

Test kommandon
* "make test".................för att köra tester för klient och server.

Rengör kommandon
* "make clean".................för att städar bort kompilerade filer.


## Katalogstruktur
<pre>
.
├── backend
│   └── src
│       ├── client.py
│       ├── data.go
│       ├── server.go
│       └── tank_test.go
├── frontend
│   ├── img
│   │   ├── BerryA.png
│   │   ├── BerryB.png
│   │   ├── Bullet.png
│   │   ├── GrapeA.png
│   │   ├── GrapeB.png
│   │   ├── LimeA.png
│   │   ├── LimeB.png
│   │   ├── OrangeA.png
│   │   ├── OrangeB.png
│   │   ├── Pipe.png
│   │   └── Woods.png
│   └── src
│       ├── client.py
│       ├── __init__.py
│       ├── menuhelpers.py
│       ├── projectile.py
│       ├── resources.py
│       ├── run.py
│       └── tank.py
├── makefile
├── meta
│   ├── gruppkontrakt.md
│   ├── images
│   │   ├── erik.jpg
│   │   ├── joel.jpg
│   │   ├── marcus.png
│   │   ├── rasmus.jpg
│   │   └── simon.jpg
│   └── medlemmar.md
└── README.md
</pre>
