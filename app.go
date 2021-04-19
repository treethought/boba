package boba

import (
	tea "github.com/charmbracelet/bubbletea"
)

type UpdateFunc func(tea.Msg) (tea.Model, tea.Cmd)

type MessageStateChange struct {
	name string
}

type App struct {
	current      tea.Model
	currentName  string
	views        map[string]tea.Model
	initFunc     func() tea.Cmd
	delegateFunc UpdateFunc
}

func NewApp() *App {
	return &App{
		initFunc:     func() tea.Cmd { return nil },
		views:        make(map[string]tea.Model),
		delegateFunc: func(tea.Msg) (tea.Model, tea.Cmd) { return nil, nil },
	}
}

func ChangeState(name string) tea.Cmd {
	return func() tea.Msg {
		return MessageStateChange{name: name}
	}
}

func (a *App) Add(name string, m tea.Model) {
	a.setModel(name, m)
}

func (a *App) setModel(name string, m tea.Model) {
	a.views[name] = m
}

func (a *App) SetFocus(name string) {
	a.current = a.views[name]
	a.currentName = name
}

func (a *App) getFocused() tea.Model {
	return a.current
}

func (a *App) SetInit(f func() tea.Cmd) {
	a.initFunc = f
}

func (a *App) SetDelgate(f UpdateFunc) {
	a.delegateFunc = f
}

func (a *App) Init() tea.Cmd {
	if a.current == nil {
		return a.initFunc()
	}
	return a.current.Init()
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if a.current == nil {
		return a, nil
	}

	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return a, tea.Quit
		}

	case MessageStateChange:
		a.SetFocus(msg.name)
		return a, nil
	}

	m, cmd := a.delegateFunc(msg)
	if cmd != nil {
		return a, cmd
	}

	current := a.getFocused()

	m, cmd = current.Update(msg)
	a.setModel(a.currentName, m)
	return a, cmd
}

func (a *App) View() string {
	current := a.getFocused()
	if current != nil {
		return current.View()
	}
	return ""
}
