package messagegroup

// A simple program demonstrating the text area component from the Bubbles
// component library.

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
	"github.com/onioncall/squa/services"
	"github.com/onioncall/squa/entities"
	"golang.org/x/term"
)

func Execute(groupUuid uuid.UUID) {
	services.Clear()
	go entities.MessagesService()
	p := tea.NewProgram(initialModel(groupUuid))

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func initialModel(groupUuid uuid.UUID) model {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(70)
	ta.SetHeight(5)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	_, tHeight, _ := term.GetSize(int(os.Stdout.Fd()))

	vp := viewport.New(70, tHeight-8)

	welcomeMessage := fmt.Sprintf("Welcome to message group %s!\nType a message and press Enter to send.", groupUuid.String())

	vp.SetContent(welcomeMessage)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return model{
		textarea:      ta,
		messages:      []string{},
		viewport:      vp,
		senderStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("29")), //lets do 5 for other chats
		recieverStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("3")),
		err:           nil,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			u := entities.User
			u.DeactivateUser();
			log.Fatal("Goodbye!")
			return m, tea.Quit

		case tea.KeyEnter:
			message := entities.DisplayMessage{
				DisplayName:     entities.User.DisplayName,
				MessageContents: m.textarea.Value(),
			}
			
			message.SendMessage()
			m.messages = append(m.messages, m.senderStyle.Render(entities.User.DisplayName+": ")+m.textarea.Value())
			m.viewport.SetContent(strings.Join(m.messages, "\n"))
			m.textarea.Reset()
			m.viewport.GotoBottom()
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	if len(entities.UnrecievedMessages) > 0 {
		for _, message := range entities.UnrecievedMessages {
			m.messages = append(m.messages, m.recieverStyle.Render(message.DisplayName+": ")+message.MessageContents)
			m.viewport.SetContent(strings.Join(m.messages, "\n"))
			m.viewport.GotoBottom()
		}
		entities.UnrecievedMessages = []entities.DisplayMessage{}
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s\n\n%s",
		m.viewport.View(),
		m.textarea.View(),
	) + "\n\n"
}
