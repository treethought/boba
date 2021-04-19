package boba

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type SubmittedFunc func(value string) tea.Cmd

type Input struct {
	textinput.Model
	onSubmit SubmittedFunc
}

func NewInput(submitted SubmittedFunc) *Input {
	return &Input{
		onSubmit: submitted,
		Model:    textinput.NewModel(),
	}
}

func (m Input) SetSubmitHandler(f SubmittedFunc) Input {
	m.onSubmit = f
	return m
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
