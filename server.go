package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func main() {
	app := fiber.New()

	// WebSocket route
	app.Get("/ws", websocket.New(handleWebSocket))

	// Start the server
	err := app.Listen(":9999")
	fmt.Println("Server running!")

	if err != nil {
		panic(err)
	}
}

// WebSocket handler function
func handleWebSocket(c *websocket.Conn) {
	fmt.Println("Client connected!")

	// Read messages from the client
	go func() {
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				fmt.Println("Client disconnected.")
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
	}()
}
