package boba

import (
	tea "github.com/charmbracelet/bubbletea"
)

type UpdateFunc func(tea.Msg) (tea.Model, tea.Cmd)

type MessageStateChange struct {
	name string
}

type bobaModel struct {
	tea.Model
	control bool
}

type App struct {
	current      tea.Model
	currentName  string
	views        map[string]bobaModel
	initFunc     func() tea.Cmd
	delegateFunc UpdateFunc
}

func NewApp() *App {
	return &App{
		initFunc:     func() tea.Cmd { return nil },
		views:        make(map[string]bobaModel),
		delegateFunc: func(tea.Msg) (tea.Model, tea.Cmd) { return nil, nil },
	}
}

func ChangeState(name string) tea.Cmd {
	return func() tea.Msg {
		return MessageStateChange{name: name}
	}
}

func (a *App) Get(name string) tea.Model {
	return a.getModel(name)
}

func (a *App) Add(name string, m tea.Model) {
	mod := bobaModel{
		Model:   m,
		control: true,
	}
	a.setModel(name, mod)
}

func (a *App) Register(name string, m tea.Model) {
	mod := bobaModel{
		Model:   m,
		control: false,
	}
	a.setModel(name, mod)
}

func (a *App) setModel(name string, m bobaModel) {
	a.views[name] = m
}

func (a *App) getModel(name string) tea.Model {
	m, ok := a.views[name]
	if !ok {
		return nil
	}
	return m
}

// SetFocus sets the model which will become focused and receive messages
// this method updates the application directly, if you would like to change focus
// via the MessageStateChange tea.Msg, you may use the ChangeState tea.Cmd
func (a *App) SetFocus(name string) {
	a.current = a.views[name]
	a.currentName = name
}

func (a *App) getFocused() tea.Model {
	return a.current
}

// SetInit sets the function to be called on initialization if no model is in focus
func (a *App) SetInit(f func() tea.Cmd) {
	a.initFunc = f
}

// SetDelegate sets the message handler to be run before messages are passed to the focused model.
// This function can be used for managing your own key bindings and state handling.
//
// If the function returns a nil tea.Cmd, then will continue and pass the message
// through to the currently focused model.
func (a *App) SetDelegate(f UpdateFunc) {
	a.delegateFunc = f
}

func (a *App) Init() tea.Cmd {
	if a.current == nil {
		return a.initFunc()
	}
	return a.current.Init()
}

// Update handles tea.Msgs to perform updates to application model.
// Messages are inspected and handled in the following manner.
//
// 1. If the msg is an exit key press (crtl-c) boba exits by invoking tea.Quit
// 2. If the msg is a boba.MessageStateChange, then the application's focus
// is changed to the reference model and no tea.Cmd is returned
// 3. The message is then passed to the user defined delegate function to perform custom message handling,
// and returns the resulting tea.Cmd if not nil
// 4. Otherwise, the message is passed onto the currently focused model
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

	// current := a.getFocused()

	cmds := []tea.Cmd{}
	// var mod tea.Model
	for n, m := range a.views {
		if !m.control {
			continue
		}
		_, cmd := m.Update(msg)
		if n == a.currentName {
			// mod = m
		}
		cmds = append(cmds, cmd)
	return a, tea.Batch(cmds...)
}

// View returns the view to be rendered by calling the currently focused model's View() function.
func (a *App) View() string {
	current := a.getFocused()
	if current != nil {
		return current.View()
	}
	return ""
}
