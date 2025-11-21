package main

import (
	"cannoliOS/models"
	"cannoliOS/ui"
	"cannoliOS/utils"
	"fmt"
	"os"
	"time"

	_ "github.com/UncleJunVIP/certifiable"
	gaba "github.com/UncleJunVIP/gabagool/pkg/gabagool"
	"github.com/idsulik/go-collections/v3/stack/arraystack"
	"github.com/veandco/go-sdl2/sdl"
)

func init() {
	gaba.Init(gaba.Options{
		WindowTitle:    "cannoli_OS",
		ShowBackground: true,
		IsCannoli:      true,
		LogFilename:    "cannoliOS.log",
	})

	logger := utils.GetLogger()

	err := utils.LoadConfig()
	if err != nil {
		logger.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	logger.Info("=== Cannoli OS Started ===")
}

func exit() {
	gaba.Close()
}

func main() {
	defer exit()

	logger := utils.GetLogger()
	var currentScreen models.Screen

	currentScreen = ui.MainMenu{
		Data:     nil,
		Position: models.Position{},
	}

	logger.Debug(fmt.Sprintf("Initial screen set to: %s", currentScreen.Name()))

	for {
		logger.Debug(fmt.Sprintf("Drawing screen: %s", currentScreen.Name()))
		sr, err := currentScreen.Draw()

		if err != nil {
			logger.Error(fmt.Sprintf("Error drawing screen %s: %v", currentScreen.Name()), "error", err)
			continue
		}

		logger.Debug(fmt.Sprintf("Screen %s returned code: %v", currentScreen.Name(), sr.Code))

		switch currentScreen.Name() {
		case models.MainMenu:
			logger.Debug("Processing MainMenu screen response")

			switch sr.Code {
			case models.Select:
				directory := sr.Output.(models.Directory)
				logger.Debug(fmt.Sprintf("Selected directory: %s (path: %s)", directory.DisplayName, directory.Path))
				currentScreen = ui.GameList{
					Directory:      directory,
					SearchFilter:   "",
					DirectoryStack: arraystack.New[models.Directory](5),
				}
				logger.Debug("Switched to GameList screen")
			case models.Action:
				logger.Debug("Action triggered in MainMenu")
			default:
				logger.Debug(fmt.Sprintf("Unhandled code in MainMenu: %v", sr.Code))
			}

		case models.GameList:
			logger.Debug("Processing GameList screen response")
			gl := currentScreen.(ui.GameList)

			if sr.Code == models.Back {
				logger.Debug(fmt.Sprintf("Back action triggered, directory stack size: %d", gl.DirectoryStack.Len()))

				if gl.DirectoryStack.Len() == 0 {
					logger.Debug("Returning to MainMenu from GameList")
					currentScreen = ui.MainMenu{
						Data:     nil,
						Position: models.Position{},
					}
				} else {
					prev, _ := gl.DirectoryStack.Pop()
					logger.Debug(fmt.Sprintf("Navigating back to directory: %s", prev.DisplayName))
					currentScreen = ui.GameList{
						Directory:      prev,
						SearchFilter:   "",
						DirectoryStack: gl.DirectoryStack,
					}
				}
			} else if sr.Code == models.Select && sr.Output.([]models.Item)[0].IsDirectory { // TODO this needs to be cleaned
				selectedItem := sr.Output.([]models.Item)[0]
				logger.Debug(fmt.Sprintf("Selected directory item: %s", selectedItem.Filename))

				gl.DirectoryStack.Push(gl.Directory)
				currentScreen = ui.GameList{
					Directory:      selectedItem.ToDirectory(),
					SearchFilter:   "",
					DirectoryStack: gl.DirectoryStack,
				}
				logger.Debug(fmt.Sprintf("Navigated into directory: %s", selectedItem.Filename))
			} else if sr.Code == models.Select {
				selectedItems := sr.Output.([]models.Item)
				logger.Debug(fmt.Sprintf("Selected %d game item for launch", len(selectedItems)))
				for i, item := range selectedItems {
					logger.Debug(fmt.Sprintf("  Item %d: %s", i+1, item.Filename))
				}

				if len(selectedItems) > 0 {
					selectedItem := selectedItems[0]
					romPath := selectedItem.Path

					gaba.HideWindow()
					utils.LaunchROM(selectedItem.DisplayName, romPath)
					gaba.ShowWindow()
				}

				logger.Debug("Returning to cannoliOS after game launch")

				for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
					// Discard all events
				}

				time.Sleep(500 * time.Millisecond)

				currentScreen = ui.GameList{
					Directory:      currentScreen.(ui.GameList).Directory,
					SearchFilter:   currentScreen.(ui.GameList).SearchFilter,
					DirectoryStack: currentScreen.(ui.GameList).DirectoryStack,
					Position:       sr.Position,
				}
			} else {
				logger.Debug(fmt.Sprintf("Unhandled code in GameList: %v", sr.Code))
			}

		default:
			logger.Debug(fmt.Sprintf("Unknown screen type: %s", currentScreen.Name()))
		}
	}
}
