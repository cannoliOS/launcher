package utils

import (
	"time"

	"github.com/UncleJunVIP/gabagool/pkg/gabagool"
	"github.com/UncleJunVIP/gabagool/pkg/gabagool/constants"
)

func ShowMessage(message string, delay time.Duration) {
	gabagool.ProcessMessage(message, gabagool.ProcessMessageOptions{}, func() (interface{}, error) {
		time.Sleep(delay)
		return nil, nil
	})
}

func ShowConfirmation(message string) bool {
	return ShowCustomConfirmation(message, "Cancel", "Confirm", constants.VirtualButtonA)
}

func ShowCustomConfirmation(message string, cancelText string, confirmText string, confirmButton constants.VirtualButton) bool {
	result, err := gabagool.ConfirmationMessage(message, []gabagool.FooterHelpItem{
		{ButtonName: "B", HelpText: cancelText},
		{ButtonName: "A", HelpText: confirmText},
	}, gabagool.MessageOptions{ConfirmButton: confirmButton})

	return err == nil && result.IsSome()
}
