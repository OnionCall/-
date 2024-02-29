package join

import (
	"github.com/charmbracelet/bubbles/textinput"
	// "github.com/charmbracelet/lipgloss"
)

// type errMsg error

type model struct {
	focusIndex int
	inputs     []textinput.Model
}