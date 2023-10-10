//go:build systray

package main

import (
	"capuchin/app/icon"
	"fmt"
	"log"
	"os/exec"
	"runtime"

	"github.com/getlantern/systray"
	"github.com/gofiber/fiber/v2"
)

func startServer(app *fiber.App, port *int64) {
	systray.Run(func() {
		systray.SetIcon(icon.Data)
		systray.SetTitle(appName)
		systray.SetTooltip(appName)

		mOpen := systray.AddMenuItem("Open app", "Open the application in the browser")
		systray.AddSeparator()
		mExit := systray.AddMenuItem("Exit", "Close the application")

		go func() {
			for {
				select {
				case <-mOpen.ClickedCh:
					openUrl(fmt.Sprintf("http://127.0.0.1:%d", *port))
				case <-mExit.ClickedCh:
					fmt.Println("👽 Requesting quit")
					systray.Quit()
				}
			}
		}()

		openUrl(fmt.Sprintf("http://127.0.0.1:%d", *port))

		log.Fatal(app.Listen(fmt.Sprintf(":%d", *port)))
	}, nil)
}

func openUrl(url string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}

	args = append(args, url)

	err := exec.Command(cmd, args...).Start()
	if err != nil {
		log.Printf("❌ openUrl: %s\n", err.Error())
		return
	}
}
