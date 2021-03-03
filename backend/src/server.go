package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

//Server is the main data structure of the program holding all data
type Server struct {
	clients   map[Client]uint32
	broadcast chan bool
	clientID  uint32
	gamestate *gamestate
}

// Client is needed for having a goroutine listening to client and updating a channel
type Client struct {
	conn net.Conn
	ch   chan string
}

func createClient(conn net.Conn) Client {
	return Client{
		conn: conn,
		ch:   make(chan string),
	}
}

func initServer() *Server {
	return &Server{
		clients:   make(map[Client]uint32),
		broadcast: make(chan bool, 1),
		clientID:  0,
		gamestate: initGamestate(),
	}
}

// Lock as global variables named according to what map they're meant to lock.
var terrainLock = sync.RWMutex{}
var projLock = sync.RWMutex{}
var tankLock = sync.RWMutex{}
var clientLock = sync.RWMutex{}

func lockAllGamestate() {
	terrainLock.Lock()
	tankLock.Lock()
	projLock.Lock()
}

func unlockAllGamestate() {
	terrainLock.Unlock()
	tankLock.Unlock()
	projLock.Unlock()
}

// handleMessages will send the current gamestate to every connected client as JSON
func broadcastState(s *Server) {
	// make gamestate into JSON
	lockAllGamestate()
	newmsg, err := json.Marshal(*(s.gamestate))
	unlockAllGamestate()
	if err != nil {
		log.Print(err)
	}
	clientLock.RLock()
	for client := range s.clients {
		// send new string back to client
		client.conn.Write([]byte(newmsg))
	}
	clientLock.RUnlock()
}

//handleConnections runs once per tick per client and receives their input, if any, and sends it into handleInput
func handleConnections(client Client, s *Server, wg *sync.WaitGroup) {
	select {
	case msg := <-client.ch:
		handleInput(msg, s.gamestate.Tanks[s.clients[client]], s.gamestate)
		wg.Done()
	default:
		wg.Done()
		return
	}
}

//acceptConnections will accept any and all connections, add them to list of clients and spawn them a tank. Will then send updated gamestate to all clients
func acceptConnections(ln *net.Listener, s *Server) {
	for {
		conn, _ := (*ln).Accept()
		if len(s.clients) == 0 {
			initTerrain(s.gamestate)
		}
		client := createClient(conn) // Saves connection and adds a channel for clients to broadcast to
		clientLock.Lock()
		s.clients[client] = s.clientID // Add to map of clients
		clientLock.Unlock()
		tankLock.Lock()
		addTank(s.gamestate, s.clientID, "red") //Marcus - alla börjar som RED nu, bör vara ett val att när man som client startar att få välja vilket team man ska vara
		tankLock.Unlock()
		s.clientID++
		go listenToClient(client, s)
		lockAllGamestate()
		newmsg, err := json.Marshal(*(s.gamestate))
		unlockAllGamestate()
		if err != nil {
			log.Print(err)
		}
		clientLock.RLock()
		for client := range s.clients {
			// send updated gamestate with new player to all clients
			client.conn.Write([]byte(newmsg))
			log.Println("wrote to client: ", s.clientID)
		}
		clientLock.RUnlock()
	}
}

// listenToClient continually reads data from client and puts it into a channel we then read from
func listenToClient(client Client, s *Server) {
	for {
		message, err := bufio.NewReader(client.conn).ReadString('\n')
		if err != nil {
			log.Print("connection dead, ending goroutine")
			tankLock.Lock()
			delete(s.gamestate.Tanks, s.clients[client])
			tankLock.Unlock()
			clientLock.Lock()
			delete(s.clients, client)
			clientLock.Unlock()
			if len(s.clients) == 0 {
				s.gamestate = initGamestate()
			}
			return
		}
		message = strings.TrimRight(message, "\r\n")
		//intmsg, err := strconv.Atoi(message)
		if err != nil {
			log.Println(err)
		}
		client.ch <- message
	}
}

func main() {

	fmt.Println("Launching server...")
	s := initServer()
	//initTerrain(s.gamestate)
	t := time.NewTicker(50 * time.Millisecond)
	var wg sync.WaitGroup
	defer t.Stop()
	// listen on all interfaces
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Print(err)
		return
	}
	go acceptConnections(&ln, s)
	for {
		<-t.C
		s.gamestate.Frame++
		go calculateProjectiles(s.gamestate)
		go tanksJump(s.gamestate)
		go tanksXVelocity(s.gamestate, s.gamestate.Tanks)
		go calcDeg(s.gamestate, s.gamestate.Tanks)
		for client := range s.clients {
			wg.Add(1)
			go handleConnections(client, s, &wg)
		}
		wg.Wait()
		if len(s.clients) > 0 {
			go broadcastState(s)
		}
	}
}
