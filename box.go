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

func (n *BoxNode) Update(msg tea.Msg) (*BoxNode, tea.Cmd) {

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
	SizeX float64
	// Percentage of Window width
	SizeY float64

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
		SizeY:       float64(x) / float64(100),
		SizeX:       float64(y) / float64(100),
		orientation: orientation,
		style:       boxStyle,
	}
}

// AddNode adds a new node with h and w percent of parent height and width
func (m *Box) AddNode(n tea.Model, w, h int) {
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

	targetWidth := int(float64(m.width) * (n.SizeX))
	targetLines := int(float64(m.height) * (n.SizeY))
	return targetWidth, targetLines

}

func (m *Box) updateNodes(msg tea.Msg) (mod *Box, cmds []tea.Cmd) {
	for _, n := range m.nodes {
		_, cmd := n.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
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
		m.width = int(float64(msg.Width) * m.SizeX)
		m.height = int(float64(msg.Height) * m.SizeY)

		cmds = append(cmds, m.resizeNodes()...)

		if !m.ready {
			m.ready = true
		}
	default:
		_, cmds = m.updateNodes(msg)
	}

	return m, tea.Batch(cmds...)
}

func (m *Box) resizeNodes() (cmds []tea.Cmd) {
	for _, n := range m.nodes {
		x, y := m.getNodeSize(n)

		msg := tea.WindowSizeMsg{
			Width:  x,
			Height: y,
		}
		_, cmd := n.Update(msg)
		cmds = append(cmds, cmd)
	}
	return cmds
}

func (m *Box) View() string {

	out := ""

	for _, n := range m.nodes {

		x, y := m.getNodeSize(n)
		if m.width == 0 || m.height == 0 {
			return ""
		}

		nodeContent := n.View()

		s := strings.ReplaceAll(nodeContent, "\r\n", "\n") // normalize line endings

		fx, fy := n.style.GetFrameSize()
		s = n.style.
			Width(x - fx).
			Height(y - fy).
			MaxWidth(x).
			MaxHeight(y).
			Render(s)

		if m.orientation == "horizontal" {
			out = lipgloss.JoinHorizontal(lipgloss.Center, out, s)

		}
		if m.orientation == "vertical" {
			out = lipgloss.JoinVertical(lipgloss.Center, out, s)
		}
	}
	x, y := m.style.GetFrameSize()
	return m.style.
		Width(m.width - x).Height(m.height - y).
		MaxWidth(m.width).MaxHeight(m.height).
		Render(out)
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
