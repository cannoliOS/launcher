package ui

import (
	"cannoliOS/models"
	"cannoliOS/state"
	"cannoliOS/utils"

	"github.com/UncleJunVIP/gabagool/pkg/gabagool"
	"github.com/UncleJunVIP/gabagool/pkg/gabagool/i18n"
)

type MainMenu struct {
	Data     interface{}
	Position models.Position
}

func (m MainMenu) Name() models.ScreenName {
	return models.MainMenu
}

func (m MainMenu) Draw() (models.ScreenReturn, error) {
	var menuItems []gabagool.MenuItem

	gameMenuItems := buildGameDirectoryMenuItems()

	menuItems = append(menuItems, gameMenuItems...)

	options := gabagool.DefaultListOptions("cannoli_OS", menuItems)

	selectedIndex, visibleStartIndex := 0, 0 //TODO replace me with actual stack state
	options.SelectedIndex = selectedIndex
	options.VisibleStartIndex = visibleStartIndex
	options.DisableBackButton = true
	options.EnableMultiSelect = false
	options.EnableAction = true

	options.FooterHelpItems = []gabagool.FooterHelpItem{
		{ButtonName: "X", HelpText: i18n.GetString("tools")},
		{ButtonName: "A", HelpText: i18n.GetString("select")},
	}

	sel, _ := gabagool.List(options)

	if sel.IsSome() && sel.Unwrap().ActionTriggered {
		return models.ScreenReturn{
			Code: models.Action,
		}, nil
	} else if sel.IsSome() && !sel.Unwrap().ActionTriggered && sel.Unwrap().SelectedIndex != -1 {
		md := sel.Unwrap().SelectedItem.Metadata
		return models.ScreenReturn{
			Output: md,
			Position: models.Position{
				SelectedIndex:    sel.Unwrap().SelectedIndex,
				SelectedPosition: sel.Unwrap().VisiblePosition,
			},
			Code: models.Select,
		}, nil
	}

	return models.ScreenReturn{
		Code: models.Canceled,
	}, nil
}

func buildGameDirectoryMenuItems() []gabagool.MenuItem {
	fb := utils.NewFileBrowser()

	if err := fb.CWD(utils.GetRomPath(), state.Get().Config.HideEmptyDirectories); err != nil {
		utils.GetLogger().Error("Failed to fetch ROM directories", "error", err)
		utils.ShowMessage("Error fetching ROM directories", 5000)
	}

	var menuItems []gabagool.MenuItem
	for _, item := range fb.Items {
		if item.IsDirectory {
			gameDirectory := item.ToDirectory()
			tagless, _ := utils.ItemNameCleaner(gameDirectory.DisplayName, true)
			menuItems = append(menuItems, gabagool.MenuItem{
				Text:     tagless,
				Selected: false,
				Focused:  false,
				Metadata: gameDirectory,
			})
		}
	}

	return menuItems
}
