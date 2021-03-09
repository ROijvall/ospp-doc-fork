package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const heightOfMap float64 = 800
const g float64 = -0.5             // Gravity
const mapSize float64 = 1200       // Width of map
const rad2deg float64 = 57.2957795 // Used in calcDeg
const firePower float64 = 20       // FIXME : Adjust this value to fit gamebalance
const explosionSize int = 50
const maxExplosionDmg int = 50
const maxVelocity float64 = 4.2
const jumpPower float64 = 9 // Bigger number = Bigger jump
const reactionHeight float64 = 350
const jumpCooldown int = 0 // amount of frames between jumps

//Gamestate holds all data needed to run the game
type gamestate struct {
	Terrain     map[int]*terrain           `json:"terrain"`
	Tanks       map[uint32]*dataTank       `json:"tanks"`
	Projectiles map[uint32]*dataProjectile `json:"projectile"` //spelar ingen roll vem som sköt projektilen, bara att de finns
	ID          uint32                     `json:"id"`         //för objekt utan tillhörande, ie projektiler
	Frame       int                        `json:"frame"`
	UniqueID    uint32                     `json:"uniqueid"`   //only read by the connecter upon first communication, should allow for client side identification
	AliveTanks  int                        `json:"alivetanks"` //maybe this can be solved by removing tanks from gamestate instead?
}

// Terrain
type terrain struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	//Material int     `json:"material"`
	//CanFall  bool    `json:"canfall"`

	/* Material:
	0 = Stone
	1 = Dirt
	2 = Sand
	*/
}

// dataTank holds all data for any given tank
type dataTank struct {
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	XVelocity float64 `json:"xVelocity"`
	YVelocity float64 `json:"yVelocity"`
	DegTank   float64 `json:"degTank"`
	DegCannon float64 `json:"degCannon"`
	Hp        int     `json:"hp"`
	Team      string  `json:"team"`
	Dir       int     `json:"dir"`
	LastFire  int     `json:"lastfire"`
	LastJump  int     `json:"lastjump"`
	Alive     bool    `json:"alive"`
	InAir     bool    `json:"InAir"`
}

// dataProjectile holds all data for any given projectile
type dataProjectile struct {
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	XVelocity float64 `json:"xVelocity"`
	YVelocity float64 `json:"yVelocity"`
	A         float64 `json:"a"`  // 	Initial angle
	V0        float64 `json:"v0"` // 	Initial velocity
	T         float64 `json:"t"`  //	Time
	ID        uint32  `json:"ID"` //	Time
}

// genTerrain returns a terrain according to x, y parameters
func genTerrain(x int, y int) *terrain {
	return &terrain{
		X: float64(x),
		Y: float64(y),
		//Material: 0,
		//CanFall:  false,
	}
}

// initTerrain will create terrain until mapSize, currently y is unchanged
func initTerrain(game *gamestate) {
	rand.New(rand.NewSource(42)) //42 best magical Constant
	rand.Seed(time.Now().UnixNano())
	x := 0

	y := heightOfMap - rand.Float64()*(heightOfMap/2) // the height to start at
	fmt.Println(y)
	//standardTerrain := y
	var dy float64 = 0
	var dyGoal float64 = 0
	var curveDensity float64 = 0

	for x < int(mapSize) {
		if curveDensity == 0 {
			dyGoal = 0.5 * (-0.5 + rand.Float64())
			curveDensity = 30
		}
		dy += dyGoal / 30
		y += dy
		if x == 0 {
			fmt.Println(int(y))
		}
		game.Terrain[x] = genTerrain(x, int(y))
		curveDensity--
		x++
		if y > heightOfMap-250 {
			dy -= 0.02
		}
		if y > heightOfMap-200 {
			dyGoal = -0.5
			dy -= 0.05
		}

		if y < reactionHeight+100 {
			dy += 0.01
		}
		if y < reactionHeight {
			dyGoal = 0.5
			dy += 0.05
		}
		if dy >= 0.33 {
			dy = 0.33
		}
		if dy <= -0.33 {
			dy = -0.33
		}

	}

}

// calcDegTank will look at a tanks x-value and compare it to the terrain beneath it and then change the tanks gradient
func calcDegTank(game *gamestate, tank *dataTank) {

	if tank.Y >= game.Terrain[int(tank.X)].Y { // possible crash
		yBefore := float64(0)
		yAfter := float64(0)
		if tank.X > 0 && tank.X < mapSize-1 {
			yBefore = game.Terrain[int(tank.X-1)].Y
			yAfter = game.Terrain[int(tank.X+1)].Y
		}
		if tank.X == 0 {
			yBefore = game.Terrain[int(tank.X)].Y
			yAfter = game.Terrain[int(tank.X+1)].Y
		}
		if tank.X == mapSize-1 {
			yBefore = game.Terrain[int(tank.X-1)].Y
			yAfter = game.Terrain[int(tank.X)].Y
		}

		k := float64((yAfter - yBefore) / 2)
		d := math.Atan(k)
		tank.DegTank = d * float64(rad2deg)
	}
}

// calcDeg will run calcDegTank for all tanks
func calcDeg(game *gamestate, tanks map[uint32]*dataTank) {
	for idx := range tanks {
		calcDegTank(game, tanks[idx])
	}
}

// initTank will create a new tank with a randomized spawn point
func initTank(team string, terrain map[int]*terrain) *dataTank {
	rand.Seed(time.Now().UnixNano())
	min := float64(0)
	max := float64(mapSize)
	x := math.RoundToEven(rand.Float64() * (max - min))
	return &dataTank{
		X:         float64(x),
		Y:         terrain[int(x)].Y,
		XVelocity: float64(0),
		YVelocity: float64(0),
		DegTank:   0,
		DegCannon: 45,
		Hp:        100,
		Team:      team,
		Dir:       1,
		Alive:     true,
	}
}

func initGamestate() *gamestate {
	return &gamestate{
		Terrain:     make(map[int]*terrain),
		Tanks:       make(map[uint32]*dataTank),
		Projectiles: make(map[uint32]*dataProjectile),
		ID:          0,
		Frame:       0,
		UniqueID:    0,
	}
}

// Adds a new tank to a given gamestate
func addTank(gamestate *gamestate, client uint32, team string) {
	gamestate.Tanks[client] = initTank(team, gamestate.Terrain)
	gamestate.AliveTanks++
}

// addProjectile adds a new projectile to a given gamestate according to a given tank
func addProjectile(gamestate *gamestate, tank *dataTank) {
	projectile := &dataProjectile{
		X:         tank.X,
		Y:         tank.Y + 1, // 1 to hinder it from hitting the ground directly after firing
		YVelocity: math.Sin(-((tank.DegCannon - tank.DegTank) * 0.0174532925)) * firePower,
		XVelocity: math.Cos(-((tank.DegCannon - tank.DegTank) * 0.0174532925)) * firePower,
		A:         tank.DegCannon,
		V0:        100,
		T:         1,
		ID:        gamestate.ID,
	}
	projLock.Lock()
	gamestate.Projectiles[gamestate.ID] = projectile
	projLock.Unlock()
	gamestate.ID++
}

// calculateExplosion will look through all tanks, reduce the hp of any tank within explosion radius
// and modify the terrain.
func calculateExplosion(x int, y int, radius int, gamestate *gamestate) {
	for i, tank := range gamestate.Tanks {
		dist := math.Sqrt((float64(x)-tank.X)*(float64(x)-tank.X) + (float64(y)-tank.Y)*(float64(y)-tank.Y))
		if dist < float64(radius) {
			changeHP(int(-float64(maxExplosionDmg)/(math.Sqrt(dist/float64(radius))+1)), gamestate.Tanks[i], gamestate)
		}
	}
	xMid := x
	xCurrent := x - radius
	if xCurrent < 0 {
		xCurrent = 0
	}
	xEnd := x + radius
	ySave := gamestate.Terrain[xCurrent].Y
	for xCurrent <= xEnd {
		distFromExp := math.Abs(float64(xCurrent - xMid))
		yPot := math.Sqrt(float64(-int(distFromExp*distFromExp)+radius*radius)) + float64(y) - 20 // seems to be a good offset
		if !(int(mapSize) > xCurrent && xCurrent >= 0) {

		} else if gamestate.Terrain[xCurrent].Y < yPot {
			if yPot < ySave && xCurrent < xMid {
				yPot = ySave
				gamestate.Terrain[xCurrent].Y = ySave
			} else if yPot > ySave && xCurrent > xMid {
				gamestate.Terrain[xCurrent].Y = ySave
				yPot = ySave
			} else {
				gamestate.Terrain[xCurrent].Y = yPot
			}
		}
		xCurrent++
		ySave = yPot

	}
}

func changeHP(change int, tank *dataTank, gamestate *gamestate) {
	tank.Hp += change
	if tank.Hp <= 0 {
		tank.Alive = false
		gamestate.AliveTanks--
		fmt.Print("tank died")
	}
	if tank.Hp > 100 {
		tank.Hp = 100
		fmt.Print("???")
	}
}

func changeTeam(team string, tank *dataTank) {
	if tank.Team != team {
		tank.Team = team
	}
}

// calculateProjectiles will iterate through every projectile and calculate their new position,
// if the new position is in terrain it will call calculateExplosion
func calculateProjectiles(gamestate *gamestate) {
	projLock.Lock()
	for _, projectile := range gamestate.Projectiles {
		projectile.X = projectile.X + projectile.XVelocity
		projectile.Y = projectile.Y + projectile.YVelocity
		//projectile.XVelocity = math.floor(projectile.XVelocity - projectile.XVelocity^2 * someConstant + wind) //where someConstant has a value that makes the function not do crazy things
		projectile.YVelocity = projectile.YVelocity - g //g should be tuned to fit the tick system so the function does not do crazy things
		if projectile.Y > heightOfMap || (projectile.X > mapSize) || (projectile.X < 0) {
			delete(gamestate.Projectiles, projectile.ID)
		} else if projectile.Y > gamestate.Terrain[int(projectile.X)].Y {
			calculateExplosion(int(projectile.X), int(gamestate.Terrain[int(projectile.X)].Y), explosionSize, gamestate)
			delete(gamestate.Projectiles, projectile.ID)
		}
	}
	projLock.Unlock()
}

func tankJump(tank *dataTank, gamestate *gamestate) {
	tank.LastJump = gamestate.Frame
	tank.InAir = true
	if tank.Y >= gamestate.Terrain[int(tank.X)].Y {
		tank.YVelocity = -jumpPower
	}
}

func tanksJump(gamestate *gamestate) {
	for _, tank := range gamestate.Tanks {
		if tank.Y < gamestate.Terrain[int(tank.X)].Y+(tank.YVelocity*-1) { // make sure this cannot go out of bounds
			tank.Y += tank.YVelocity
			tank.YVelocity -= g

			if tank.Y > gamestate.Terrain[int(tank.X)].Y {
				tank.YVelocity = 0
				tank.InAir = false
			}
		} else {
			tank.Y = gamestate.Terrain[int(tank.X)].Y + 1
			tank.YVelocity = 0
			tank.InAir = false
		}
	}
}

func calculateCollision(playingTank *dataTank, tanks map[uint32]*dataTank) bool {
	for _, tank := range tanks {
		if playingTank != tank {
			//difference := math.Sqrt(float64(math.Pow(playingTank.X-tank.X, 2) + math.Pow(playingTank.Y-tank.Y, 2)))
			//fmt.Println(playingTank.X, playingTank.Y, tank.X, tank.Y)
			//fmt.Println(difference)
			// if difference < 50 {
			// 	return true
			// }
		}
	}
	return false
}

func naturalDeacceleration(tanks map[uint32]*dataTank) {
	for _, tank := range tanks {
		if !tank.InAir {
			if tank.XVelocity <= 0.2 && tank.XVelocity >= -0.2 { // natural deacceleration, but what if it never hits 0?
				tank.XVelocity = 0
			} else if tank.XVelocity < 0 {
				tank.XVelocity += 0.2
			} else if tank.XVelocity > 0 {
				tank.XVelocity -= 0.2
			}
		}
	}
}

func tanksXMovement(gamestate *gamestate, tanks map[uint32]*dataTank) {
	slopeConst := 3 * maxVelocity // this makes sense because slopes are capped at 0.33 gradient
	naturalDeacceleration(tanks)
	for _, tank := range tanks {
		if tank.X+tank.XVelocity > 0 && mapSize > tank.X+tank.XVelocity && tank.XVelocity != 0 {
			if !tank.InAir {
				potentialMove := gamestate.Terrain[int(tank.X+tank.XVelocity)].Y // this has to be in bounds
				yDiff := potentialMove - tank.Y
				a := yDiff / slopeConst // should range between -0.33-0.33, more speed > larger penalty when climbing a slope
				if tank.XVelocity < 0 {
					if yDiff > 0 {
						tank.XVelocity += a
					}
				} else if tank.XVelocity > 0 {
					if yDiff > 0 {
						tank.XVelocity -= a
					}
				}
			}
			tank.X += tank.XVelocity
		}
	}
}

//handleInput is the only function that is called through the Server and will change gamestate according to input
func handleInput(input string, tank *dataTank, gamestate *gamestate) {
	if input != "" {
		inputs := strings.Split(input, ",") // every command is split into a single string which are then processed one by one
		for _, x := range inputs {
			y, _ := strconv.Atoi(x) // might not be entirely necessary
			switch y {
			case 0: //move right
				if !tank.InAir {
					tankLock.Lock()
					if tank.XVelocity < 0 {
						tank.XVelocity = 0
					} else if tank.XVelocity+0.8 > maxVelocity {
						tank.XVelocity = maxVelocity
					} else {
						tank.XVelocity += 0.8
					}
					tankLock.Unlock()
				}
			case 1:
				if !tank.InAir {
					tankLock.Lock()
					//if tank.XVelocity > -maxVelocity {
					if tank.XVelocity > 0 {
						tank.XVelocity = 0
					} else if tank.XVelocity-0.8 < -maxVelocity {
						tank.XVelocity = -maxVelocity
					} else {
						//print("adds negative velocity")
						tank.XVelocity -= 0.8
					}
					tankLock.Unlock()
				}

			case 2:
				if 0 <= tank.DegCannon && tank.DegCannon < 180 {
					tankLock.Lock()
					tank.DegCannon++
					tankLock.Unlock()
				}
			case 3:
				if 0 < tank.DegCannon && tank.DegCannon <= 180 {
					tankLock.Lock()
					tank.DegCannon--
					tankLock.Unlock()
				}
			case 4: //Jump
				if gamestate.Frame > tank.LastJump+jumpCooldown {
					tankLock.Lock()
					tankJump(tank, gamestate)
					tankLock.Unlock()
				}
			case 6:
				if gamestate.Frame > tank.LastFire+20 {
					tank.LastFire = gamestate.Frame
					addProjectile(gamestate, tank)
				}
			case 9:
				tankLock.Lock()
				fmt.Println("suicide initiated")
				changeHP(-100, tank, gamestate)
				tankLock.Unlock()
			// Cases 100+ are only for testing, will not be used in the game!
			case 100:
				if 0 <= tank.DegTank && tank.DegTank < 180 {
					tank.DegTank++
				}
			case 101:
				if 0 < tank.DegTank && tank.DegTank <= 180 {
					tank.DegTank--
				}
			case 102:
				calcDeg(gamestate, gamestate.Tanks)
			case 103:
				tank.X = float64(18)
				tank.Y = gamestate.Terrain[int(tank.X)].Y
			case 104:
				tank.Y++
			case 105:
				tank.Y--
			default:
				fmt.Println("Invalid input" + x)
			}
		}
	}
}
