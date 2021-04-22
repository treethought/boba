package boba

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)


type Input struct {
	textinput.Model
	onSubmit func(v string) tea.Cmd
}

func NewInput() *Input {
	return &Input{
		onSubmit: func(v string) tea.Cmd {return nil},
		Model:    textinput.NewModel(),
	}
}

func (m *Input) SetSubmitHandler(f func(val string) tea.Cmd) {
	m.onSubmit = f
}

func (m *Input) Init() tea.Cmd {
	return textinput.Blink
}

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
