package main

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

var connections []*websocket.Conn

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
		fmt.Println("Running job", i)
		time.Sleep(time.Second)
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
		_, msg, err := c.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}
		fmt.Printf("Received message: %s\n", msg)

		// Echo the message back to the client
		err = c.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			fmt.Println("Error writing message:", err)
			break
		}
	}
}
