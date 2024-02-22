package welcome

import (
	"github.com/charmbracelet/bubbles/textinput"
)

type errMsg error

type model struct {
	textInput textinput.Model
	err       error
}
