package messagegroup

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

type errMsg error

type model struct {
	viewport      viewport.Model
	messages      []string
	textarea      textarea.Model
	senderStyle   lipgloss.Style
	recieverStyle lipgloss.Style
	errorStyle    lipgloss.Style
	err           error
}

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render

