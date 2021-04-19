package boba

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type Viewer interface {
	View() string
}

type viewer struct {
	val interface{}
}

func (v viewer) View() string {
	return fmt.Sprint(v.val)
}

type SelectedFunc func(msg tea.Msg) (tea.Model, tea.Cmd)

type List struct {
	items  []Viewer
	cursor int
}

func NewList() *List {
	return &List{
		items:  []Viewer{},
		cursor: 0,
	}
}

func (m *List) Clear() {
	m.items = nil
}

func (m *List) AddItem(item Viewer) {
	m.items = append(m.items, item)
}

func (m *List) Add(item interface{}) {
	v := viewer{item}
	m.items = append(m.items, v)
}

func (m *List) CurrentItem() Viewer {
	return m.items[m.cursor]
}

func (m *List) Init() tea.Cmd {
	return nil
}

func (m *List) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		}
	}

	return m, nil
}

func (m *List) View() string {
	s := ""
	for idx, item := range m.items {
		cursor := "  "
		if idx == m.cursor {
			cursor = "> "
		}
		s += fmt.Sprintf("%s%s\n", cursor, item.View())
	}

	return s
}
