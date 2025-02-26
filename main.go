package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// list down the entities in the game
type Player struct {
	ID    string
	X, Y  float64
	Color string          `json:"color"`
	Conn  *websocket.Conn `json:"-"` // Don't send connection info in JSON
}

func randomColor() string {
	colors := []string{"red", "green", "blue", "yellow", "purple", "orange"}
	return colors[rand.Intn(len(colors))]
}

// Game state
var players = make(map[string]*Player)
var playersMutex sync.Mutex

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Websocket upgrade error:", err)
		return
	}

	defer conn.Close()

	// generate random ID
	playerID := fmt.Sprintf("%d", time.Now().UnixNano())
	// newPlayer := &Player{ID: playerID, X: 400, Y: 400}
	newPlayer := &Player{ID: playerID, X: 400, Y: 400, Conn: conn, Color: randomColor()}

	playersMutex.Lock()
	players[playerID] = newPlayer
	playersMutex.Unlock()

	log.Println("New Player connected:", playerID)

	//game loop
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			playersMutex.Lock()
			delete(players, playerID)
			playersMutex.Unlock()
			fmt.Println("Player disconnected:", playerID)
			break
		}

		// Handle movement
		if string(msg) == "up" {
			newPlayer.Y -= 5
		} else if string(msg) == "down" {
			newPlayer.Y += 5
		} else if string(msg) == "left" {
			newPlayer.X -= 5
		} else if string(msg) == "right" {
			newPlayer.X += 5
		}

		// Send game state to all clients
		broadcastState()
	}
}

// Send game state to all players
func broadcastState() {
	playersMutex.Lock()
	defer playersMutex.Unlock()

	// Convert player state to JSON
	gameState, err := json.Marshal(players)
	if err != nil {
		log.Println("Error marshalling game state:", err)
		return
	}

	log.Println("Broadcasting game state:", string(gameState)) // Debugging

	// Send game state to all players
	for _, player := range players {
		err := player.Conn.WriteMessage(websocket.TextMessage, gameState)
		if err != nil {
			log.Println("Error sending game state to", player.ID, ":", err)
			player.Conn.Close()
			delete(players, player.ID)
		}
	}
}

func main() {

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs) // Serve static files from "static" folder

	http.HandleFunc("/ws", handleConnections)
	log.Println("WebSocket server running on ws://localhost:8080/ws")
	http.ListenAndServe(":8080", nil)
}
