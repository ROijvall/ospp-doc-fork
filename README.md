# Super-Tanks (Doc)

This is an unofficial fork of a group project from 2019. The goal of the fork is simply to improve upon the original and to learn while doing.

Credit goes to github users: 
SpriCo
Tamrannman
SWallbing
joelhaile
ROijvall

Super-Tanks is an online 2D multiplayer game inspired by games like 
Pocket Tanks and Tanks. Each player controls a tank that can move 
laterally and shoot parabolic projectiles from an adjustable barrel. 
If a projectile hits the ground it affects the surrounding terrain. If
a projectile hits a player they take damage, the goal of the game is to be the last player surviving,

Unlike its predecessors Super-Tanks is played in real time which presents
additional challenges when it comes to balancing and map creation.

## Building the project

Currently the makefile is not updated. Golang and python and pygame are currently required to compile and run the project.
* go build server.go data.go (builds "server.exe" which can be launched)
* python run.py (launches client side application)


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
└── README.md
└── Todo.txt
└── .gitignore
</pre>
