package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type LocationData struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Address   string  `json:"address"`
	StudentID string  `json:"studentId"`
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// Allow all origins for simplicity, should restrict this in production
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	clients      = make(map[*websocket.Conn]bool)
	clientsMutex sync.Mutex

	latestLocation     *LocationData
	latestLocationLock sync.RWMutex
)

// Handle location updates from the React Native app
func HandleLocationUpdate(w http.ResponseWriter, r *http.Request) {
	studentID := r.Header.Get("X-Student-ID")
	if studentID == "" {
		http.Error(w, "Unauthorized: Missing student ID", http.StatusUnauthorized)
		return
	}

	if !isValidStudentID(studentID) {
		http.Error(w, "Unauthorized: Invalid student ID", http.StatusUnauthorized)
		return
	}

	var locationData LocationData
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&locationData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	locationData.StudentID = studentID

	latestLocationLock.Lock()
	latestLocation = &locationData
	latestLocationLock.Unlock()

	broadcastLocation(locationData)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Location data received",
	})
}

// Validate student ID
func isValidStudentID(id string) bool {
	return id != "" && len(id) >= 5
}

// Handle WebSocket connections
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade to websocket: %v", err)
		return
	}

	clientsMutex.Lock()
	clients[conn] = true
	clientsMutex.Unlock()

	latestLocationLock.RLock()
	if latestLocation != nil {
		conn.WriteJSON(latestLocation)
	}
	latestLocationLock.RUnlock()

	defer func() {
		clientsMutex.Lock()
		delete(clients, conn)
		clientsMutex.Unlock()
		conn.Close()
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

// Broadcast location updates to all connected clients
func broadcastLocation(data LocationData) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	for client := range clients {
		err := client.WriteJSON(data)
		if err != nil {
			log.Printf("Error sending to client: %v", err)
			client.Close()
			delete(clients, client)
		}
	}
}
