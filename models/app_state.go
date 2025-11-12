package models

import (
	"github.com/idsulik/go-collections/v3/stack"
)

type AppState struct {
	Config      *Config
	ScreenStack stack.Stack[string]
}
