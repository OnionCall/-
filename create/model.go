package create

import (
	"github.com/charmbracelet/bubbles/textinput"
)

// type errMsg error

type model struct {
	focusIndex int
	inputs     []textinput.Model
}
