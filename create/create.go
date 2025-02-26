package create

// A simple example demonstrating the use of multiple text input components
// from the Bubbles component library.

import (
	"fmt"
	// "log"
	"os"
	"strings"

	// "github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	messagegroup "github.com/OnionCall/squa/message-group"

	//"github.com/OnionCall/squa/messagegroup"
	"github.com/OnionCall/squa/entities"
	"github.com/OnionCall/squa/services"

	// "github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff8c00"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle  = focusedStyle.Copy()
	noStyle      = lipgloss.NewStyle()

	focusedButton = focusedStyle.Copy().Render("[ CREATE ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("CREATE"))

	passwordsMatch = true
)

func Execute() {
	services.Clear()
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}

func initialModel() model {
	m := model{
		inputs: make([]textinput.Model, 3),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "display name"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "set password..."
			t.CharLimit = 30
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '*'
		case 2:
			t.Placeholder = "confirm password..."
			t.CharLimit = 30
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '*'
		}

		m.inputs[i] = t
	}

	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	u := entities.UserDetails{}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			os.Exit(0)
			return m, tea.Quit

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if m.inputs[1].Value() != m.inputs[2].Value() {
				passwordsMatch = false
			} else { // I wish I didn't need to use an else here, but it seems to be the easiest way
				passwordsMatch = true
			}

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.inputs) && passwordsMatch {
				u.DisplayName = m.inputs[0].Value()

				if len(strings.TrimSpace(m.inputs[0].Value())) == 0 {
					u.DisplayName = services.GenerateDefaultName()
				}

				groupUuid := services.GenerateUuid()
				g := entities.MessageGroup {
					GroupUuid: groupUuid,
					GroupKey:  m.inputs[1].Value(),
				}

				groupId := g.CreateGroup()

				u.GroupId = groupId
				u.CreateUser()

				services.Clear()
				messagegroup.Execute(groupUuid)

				return m, tea.Quit
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m model) View() string {
	var b strings.Builder

	b.WriteString("Enter a display name, then enter and confirm password. \nIf no password is required, leave them blank\n\n")
	if !passwordsMatch {
		b.WriteString("Passwords do not match\n")
	}

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	return b.String()
}
