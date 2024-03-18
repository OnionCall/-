package welcome

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/onioncall/squa/common"
	"github.com/onioncall/squa/create"
	"github.com/onioncall/squa/join"
)

var (
	style       = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff8c00"))
	cursorStyle = style.Copy()
)

func Execute() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Type Create or Join"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	ti.Cursor.Style = cursorStyle

	return model{
		textInput: ti,
		err:       nil,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	lcInput := strings.ToLower(m.textInput.Value())

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:

			if lcInput == "create" || lcInput == "c" {
				create.Execute()
			} else if lcInput == "join" || lcInput == "j" {
				join.Execute()
			} else {
				//back to welcome screen
				Execute()
			}

			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	environmentText := ""
	if common.Environment == "http://localhost:8040" {
		environmentText = "Development Mode\n"
	}
	
	return fmt.Sprintf(
		"%sWelcome to squa, a way to send and receive messages from the terminal.\nCreate or Join a message group\n\n%s\n",
		environmentText,
		m.textInput.View(),
	) + "\n"
}
