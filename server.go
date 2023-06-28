//go:build !systray

package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

func startServer(app *fiber.App, port *int64) {
	log.Fatal(app.Listen(fmt.Sprintf(":%d", *port)))
}
