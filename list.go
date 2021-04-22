package boba

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// Viewer implements the View() method of a tea.Model
type Viewer interface {
	View() string
}

// viewer is a wrapping struct allowing non tea.Model's to be added as list items
type viewer struct {
	val interface{}
}

func (v viewer) View() string {
	return fmt.Sprint(v.val)
}

// List is a model that displays a list of items that can be navigated and selected.
// The only requirement for list items are that they satisfy the Viewer interface
type List struct {
	items        []Viewer
	cursor       int
	selectedFunc func(tea.Msg) tea.Cmd
}

// NewList returns a new List model
func NewList() *List {
	return &List{
		items:        []Viewer{},
		cursor:       0,
		selectedFunc: func(tea.Msg) tea.Cmd { return nil },
	}
}

// Clear removes all items from the List
func (m *List) Clear() {
	m.items = nil
}

// SetSelectedFunc sets the function to be called when an item is selected.
// It receives the selected item as a tea.Msg and must return a tea.Cmd
func (m *List) SetSelectedFunc(f func(tea.Msg) tea.Cmd) {
	m.selectedFunc = f
}

//AddItem adds a new item to the list
func (m *List) AddItem(item Viewer) {
	m.items = append(m.items, item)
}

// Add adds a generic item to the list which will be rendered with fmt.Sprint
func (m *List) Add(item interface{}) {
	v := viewer{item}
	m.items = append(m.items, v)
}

// CurrentItem returns the item at the location of the cursor
func (m *List) CurrentItem() Viewer {
	return m.items[m.cursor]
}

func (m *List) Init() tea.Cmd {
	return nil
}

// Update handles messages for moving the cursor and selecting an item
func (m *List) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.cursor > len(m.items) {
		m.cursor = 0
	}
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "enter":
			return m, m.selectedFunc(m.CurrentItem())

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
