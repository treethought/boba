package boba

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Input is a wrapper around the bubbles textinput component
// It exposes the ability to define a function tio be called on input submission
type Input struct {
	textinput.Model
	onSubmit func(v string) tea.Cmd
}

// NewInput returns a new Input model
func NewInput() *Input {
	return &Input{
		onSubmit: func(v string) tea.Cmd { return nil },
		Model:    textinput.NewModel(),
	}
}

// SetSubmitHandler sets the function to be called when the user input is submitted with Enter
func (m *Input) SetSubmitHandler(f func(val string) tea.Cmd) {
	m.onSubmit = f
}

func (m *Input) Init() tea.Cmd {
	return textinput.Blink
}

// Update calls the submit function if Enter was pressed and calls the SubmitHandler,
// otherwise the message is passed on to the undnerlying bubbles textinput
func (m *Input) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.Focus()
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			val := m.Value()
			m.SetValue("")
			m.Blur()
			return m, m.onSubmit(val)
		}
	}

	m.Model, cmd = m.Model.Update(msg)
	return m, cmd
}

func (m *Input) View() string {
	return m.Model.View()
}
