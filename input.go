package boba

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type SubmittedFunc func(value string) tea.Cmd

type Input struct {
	input    textinput.Model
	onSubmit SubmittedFunc
}

func NewInput(submitted SubmittedFunc) *Input {

	return &Input{
		input:    textinput.NewModel(),
		onSubmit: submitted,
	}
}

func (m Input) SetSubmitHandler(f SubmittedFunc) Input {
	m.onSubmit = f
	return m
}

func (m *Input) Focus() {
	m.input.Focus()
}

func (m *Input) Value() string {
	return m.input.Value()
}

func (m *Input) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Input) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.input.Focus()
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			val := m.Value()
			m.input.SetValue("")
			return m, m.onSubmit(val)
		}
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m *Input) View() string {
	return m.input.View()
}
