package main

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

var connections []*websocket.Conn
var lastConnectionSentIndex = 0

func main() {
	app := fiber.New()

	// WebSocket route
	app.Get("/ws", websocket.New(handleWebSocket))

	// Start the server
	go func() {
		err := app.Listen(":9999")
		if err != nil {
			panic(err)
		}
	}()
	go handleJobs()
	select {}
}

func handleJobs() {
	// Iterate over an array from one to infinity
	fmt.Println("handling jobs")
	for i := 1; ; i++ {
		time.Sleep(time.Second)

		fmt.Println("Running job", i)

		// Check if connections is empty
		if len(connections) == 0 {
			fmt.Println("Could not send job", i, "to any clients")
			continue
		}

		// Get xth element from connections
		c := connections[lastConnectionSentIndex]

		// Echo the message back to the client

		msg := []byte(fmt.Sprintf("Job %v", i))

		// Send message to client (c)
		err := c.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			fmt.Println("Error writing message:", err)
			break
		}

		// Increment lastConnectionSentIndex
		lastConnectionSentIndex++
		if lastConnectionSentIndex >= len(connections) {
			lastConnectionSentIndex = 0
		}

	}
}

// WebSocket handler function
func handleWebSocket(c *websocket.Conn) {
	fmt.Println("Client connected!")
	connections = append(connections, c)

	// Handle WebSocket close event
	defer func() {
		fmt.Println("Client disconnected.")
		err := c.Close()

		for i, conn := range connections {
			if conn == c {
				connections = append(connections[:i], connections[i+1:]...)
				break
			}
		}

		if err != nil {
			fmt.Println("Error closing connection:", err)
		}
	}()

	// Read messages from the client
	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}
	}
}
