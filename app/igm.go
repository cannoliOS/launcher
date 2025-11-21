package main

import (
	"cannoliOS/models"
	"cannoliOS/retroarch"
	"cannoliOS/ui"
	"cannoliOS/utils"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/UncleJunVIP/certifiable"
	gaba "github.com/UncleJunVIP/gabagool/pkg/gabagool"
	"github.com/UncleJunVIP/gabagool/pkg/gabagool/i18n"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	coolDownTime = 1 * time.Second
	quitTimeout  = 1750 * time.Millisecond
)

type menuAction string

const (
	Back       menuAction = ""
	ResumeGame            = "RESUME"
	SaveState             = "SAVE_STATE"
	LoadState             = "LOAD_STATE"
	Reset                 = "RESET"
	Settings              = "MENU_TOGGLE"
	Quit                  = "QUIT"
	Screenshot            = "SCREENSHOT"
)

var localIP = "127.0.0.1"

var romPath string
var gameName string

func init() {
	gaba.Init(gaba.Options{
		WindowTitle:    "In-Game Menu",
		ShowBackground: true,
		IsCannoli:      true,
		LogFilename:    "igm.log",
	})

	utils.LoadConfig()
}
func main() {
	defer gaba.Close()

	if len(os.Args) < 2 {
		log.Fatal("Usage: igm <rom file path>")
	}

	romPath = os.Args[1]
	gameName, _ = utils.ItemNameCleaner(filepath.Base(romPath), true)

	logger := utils.GetLogger()

	logger.Debug(fmt.Sprintf("Starting IGM for %s...", gameName))
	logger.Debug(fmt.Sprintf("ROM path: %s", romPath))

	gaba.HideWindow()

	retroArchCmd, err := retroarch.Launch(gameName, romPath)
	if err != nil {
		logger.Error("Failed to launch RetroArch", "error", err)
		return
	}

	retroArchExitChan := make(chan error, 1)
	shutdownChan := make(chan bool, 1)

	go func() {
		err := retroArchCmd.Wait()
		retroArchExitChan <- err
	}()

	go func() {
		err := <-retroArchExitChan
		if err != nil {
			logger.Debug("RetroArch process ended with error", "error", err)
		} else {
			logger.Debug("RetroArch process completed successfully")
		}
		shutdownChan <- true
	}()

	menuButtonHandler(shutdownChan)

	gaba.ProcessMessage(fmt.Sprintf("%s %s...", i18n.GetString("quitting"), gameName),
		gaba.ProcessMessageOptions{}, func() (interface{}, error) {
			logger.Debug("Waiting for retroarch process to exit...")
			select {
			case <-retroArchExitChan:
				logger.Debug("RetroArch process exited")
				time.Sleep(1750 * time.Millisecond)
			case <-time.After(quitTimeout):
				logger.Debug("Quit timeout reached... sigkilling.")
				retroarch.Terminate()
			}

			logger.Debug("Shutting down IGM...")

			gaba.HideWindow()

			os.Exit(0)
			return nil, nil
		})
}

func menuButtonHandler(shutdownChan chan bool) {
	var cooldownUntil time.Time

	for {
		select {
		case <-shutdownChan:
			return
		default:
		}

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.KeyboardEvent:
				if e.Keysym.Sym == sdl.K_1 {
					if e.State == sdl.PRESSED {
						log.Println("Button press detected, toggling menu...")
						cooldownUntil = time.Now().Add(coolDownTime)
						toggleMenu(shutdownChan)
					}
				}

			case *sdl.ControllerButtonEvent:
				if time.Now().Before(cooldownUntil) {
					continue
				}

				// TODO Fix this with virtual buttons
				//var menuButton gaba.Button = 5
				//
				//if gaba.Button(e.Button) == menuButton {
				//	if e.State == sdl.PRESSED {
				//		log.Println("Button press detected, toggling menu...")
				//		cooldownUntil = time.Now().Add(coolDownTime)
				//		toggleMenu(shutdownChan)
				//	}
				//}
			}
		}

		sdl.Delay(16)
	}
}

func toggleMenu(shutdownChan chan bool) {
	logger := utils.GetLogger()

	//retroarch.SendCommand(Screenshot, localIP, "55355")
	time.Sleep(500 * time.Millisecond)

	retroarch.Pause()

	time.Sleep(200 * time.Millisecond)

	gaba.ShowWindow()
	command, _ := igm()

	logger.Debug(fmt.Sprintf("In-game menu command: %s", command))

	if command == Quit {
		logger.Debug("Quit command received, initiating retroarch termination...")
		shutdownChan <- true
		return
	}

	retroarch.Resume()
	if command != Back {
		retroarch.SendCommand(string(command), localIP, "55355")
	}

	logger.Debug("Hiding IGM...")
	gaba.HideWindow()
}

func igm() (menuAction, string) {
	logger := utils.GetLogger()

	logger.Debug("Showing in-game menu for ROM", "game_name", gameName)

	currentScreen := ui.InGameMenu{
		Data:     nil,
		Position: models.Position{},
		GameName: gameName,
	}

	for {
		sr, err := currentScreen.Draw()
		if err != nil {
			logger.Error("Error drawing in-game menu", "error", err)
			continue
		}

		switch sr.Code {
		case models.Back, models.Canceled:
			logger.Debug("Menu cancelled or back pressed")
			return Back, ""

		case models.Select:
			action := sr.Output.(string)
			logger.Debug("In-game menu action", "action", action)

			switch action {
			case "resume":
				return ResumeGame, ""

			case "save_state":
				return SaveState, i18n.GetString("saving")

			case "load_state":
				return LoadState, i18n.GetString("loading")

			case "reset":
				return Reset, i18n.GetString("resetting")

			case "settings":
				return Settings, ""

			case "quit":
				return Quit, i18n.GetString("quitting")
			default:
				logger.Debug("Unhandled menu action", "action", action)
				continue
			}
		default:
			logger.Debug("Unhandled screen response code", "code", sr.Code)
		}
	}
}
