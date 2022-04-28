package boba

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	boxStyle = lipgloss.NewStyle().Padding(0).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.AdaptiveColor{Light: "#B793FF", Dark: "#AD58B4"})

	nodeStyle = lipgloss.NewStyle()
)

type BoxNode struct {
	tea.Model
	// Height represents the percent of node's parent height to fill
	SizeX float64
	// Width represents the percent of node's parent width to fill
	SizeY float64
	style lipgloss.Style
}

func (n *BoxNode) Init() tea.Cmd {
	return nil
}

func (n *BoxNode) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	m, cmd := n.Model.Update(msg)

	node := &BoxNode{
		Model: m,
		SizeY: n.SizeY,
		SizeX: n.SizeX,
		style: n.style,
	}
	return node, cmd
}

type joinFunc func(n BoxNode, join func(pos lipgloss.Position, strs ...string) string)

type Box struct {
	nodes []*BoxNode
	// Percentage of Window height
	SizeX int
	// Percentage of Window width
	SizeY int

	width  int
	height int
	ready  bool
	// joinFunc func(pos lipgloss.Position, strs ...string) string
	orientation string
	style       lipgloss.Style
}

func NewBox(orientation string, x, y int) *Box {
	return &Box{
		nodes:       []*BoxNode{},
		SizeX:       x,
		SizeY:       y,
		orientation: orientation,
		style:       boxStyle,
	}
}

// AddNode adds a new node with h and w percent of parent height and width
func (m *Box) AddNode(n tea.Model, h, w int) {
	node := &BoxNode{
		Model: n,
		SizeY: float64(h) / float64(100),
		SizeX: float64(w) / float64(100),
		style: nodeStyle,
	}
	m.nodes = append(m.nodes, node)
}

func (m *Box) AddNodeWithStyle(n tea.Model, h, w int, style lipgloss.Style) {
	node := &BoxNode{
		Model: n,
		SizeY: float64(h) / float64(100),
		SizeX: float64(w) / float64(100),
		style: style,
	}
	m.nodes = append(m.nodes, node)
}

func (m *Box) Init() tea.Cmd {
	cmds := []tea.Cmd{}
	for _, n := range m.nodes {
		cmds = append(cmds, n.Init())
	}

	return tea.Batch(cmds...)
}

func (m Box) getNodeSize(n *BoxNode) (w int, y int) {

	x, y := n.style.GetFrameSize()

	targetWidth := int(float64(m.width)*(n.SizeX)) - x
	targetLines := int(float64(m.height)*(n.SizeY)) - y
	return targetWidth, targetLines

}

func (m *Box) updateNodes(msg tea.Msg) (mod *Box, cmds []tea.Cmd) {
	nodes := []*BoxNode{}
	for _, n := range m.nodes {
		nmod, cmd := n.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

		nNode, ok := nmod.(*BoxNode)
		if !ok {
			continue
		}
		nodes = append(nodes, nNode)
	}
	m.nodes = nodes
	return m, cmds
}

func (m *Box) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	// case tea.KeyMsg:
	// 	return m, nil

	case tea.WindowSizeMsg:
		x, y := m.style.GetFrameSize()

		m.width = int(float64(msg.Width)*(float64(m.SizeX/100))) - x
		m.height = int(float64(msg.Height)*(float64(m.SizeY/100))) - y

		cmds = append(cmds, m.resizeNodes()...)

		if !m.ready {
			m.ready = true
		}
	default:
		m, cmds = m.updateNodes(msg)
	}

	return m, tea.Batch(cmds...)
}

func (m *Box) resizeNodes() (cmds []tea.Cmd) {
	for i, n := range m.nodes {
		x, y := m.getNodeSize(n)

		msg := tea.WindowSizeMsg{
			Width:  x,
			Height: y,
		}
		mod, cmd := n.Update(msg)
		nm, ok := mod.(*BoxNode)
		if ok {
			m.nodes[i] = nm
		}
		cmds = append(cmds, cmd)
	}
	return cmds
}

func (m *Box) View() string {

	out := ""

	for _, n := range m.nodes {

		x, y := n.style.GetFrameSize()

		targetWidth, targetLines := m.getNodeSize(n)

		nodeContent := n.View()

		s := strings.ReplaceAll(nodeContent, "\r\n", "\n") // normalize line endings

		s = n.style.
			Width(targetWidth - x).
			// Height(targetLines - y).
			// MaxWidth(targetLines - x).
			MaxHeight(targetLines - y).
			Render(s)

		if m.orientation == "horizontal" {
			out = lipgloss.JoinHorizontal(lipgloss.Center, out, s)

		}
		if m.orientation == "vertical" {
			out = lipgloss.JoinVertical(lipgloss.Center, out, s)
		}
	}
	// x, y := m.style.GetFrameSize()
	return m.style.
		// Width(m.width - x).Height(m.height - y).
		// MaxWidth(m.width - x).MaxHeight(m.height - y).
		Render(out)
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
