package ui

import (
	"cannoliOS/models"

	"github.com/UncleJunVIP/gabagool/pkg/gabagool"
	"github.com/UncleJunVIP/gabagool/pkg/gabagool/i18n"
)

type InGameMenu struct {
	Data     interface{}
	Position models.Position
	ROMPath  string
	GameName string
}

func (igm InGameMenu) Name() models.ScreenName {
	return models.InGameMenu
}

func (igm InGameMenu) Draw() (models.ScreenReturn, error) {
	menuItems := []gabagool.MenuItem{
		{
			Text:     i18n.GetString("resume"),
			Selected: false,
			Focused:  false,
			Metadata: "resume",
		},
		{
			Text:     i18n.GetString("save_state"),
			Selected: false,
			Focused:  false,
			Metadata: "save_state",
		},
		{
			Text:     i18n.GetString("load_state"),
			Selected: false,
			Focused:  false,
			Metadata: "load_state",
		},
		{
			Text:     i18n.GetString("reset_game"),
			Selected: false,
			Focused:  false,
			Metadata: "reset",
		},
		{
			Text:     i18n.GetString("settings"),
			Selected: false,
			Focused:  false,
			Metadata: "settings",
		},
		{
			Text:     i18n.GetString("quit"),
			Selected: false,
			Focused:  false,
			Metadata: "quit",
		},
	}

	title := "In-Game Menu"

	if igm.GameName != "" {
		title = igm.GameName
	}

	options := gabagool.DefaultListOptions(title, menuItems)

	options.SmallTitle = true
	options.SelectedIndex = igm.Position.SelectedIndex
	options.VisibleStartIndex = igm.Position.SelectedPosition

	options.FooterHelpItems = []gabagool.FooterHelpItem{
		{ButtonName: "B", HelpText: i18n.GetString("back")},
		{ButtonName: "A", HelpText: i18n.GetString("select")},
	}

	sel, err := gabagool.List(options)
	if err != nil {
		return models.ScreenReturn{
			Code: models.Canceled,
		}, err
	}

	if sel.IsSome() {
		result := sel.Unwrap()

		if result.SelectedIndex == -1 {
			return models.ScreenReturn{
				Code: models.Back,
			}, nil
		}

		selectedAction := result.SelectedItem.Metadata.(string)

		return models.ScreenReturn{
			Output: selectedAction,
			Position: models.Position{
				SelectedIndex:    result.SelectedIndex,
				SelectedPosition: result.VisiblePosition,
			},
			Code: models.Select,
		}, nil
	}

	return models.ScreenReturn{
		Code: models.Canceled,
	}, nil
}
