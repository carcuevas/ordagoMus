package tui

import "github.com/charmbracelet/lipgloss"

var (
	styleTitle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("230")).Background(lipgloss.Color("22")).Padding(0, 1)
	styleLance  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("214"))
	styleSena   = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("141"))
	styleTantos = lipgloss.NewStyle().Foreground(lipgloss.Color("114"))
	styleDim    = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	styleName   = lipgloss.NewStyle().Bold(true)
	styleTurn   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("226"))
	styleBubble = lipgloss.NewStyle().Foreground(lipgloss.Color("231")).Bold(true)
	styleFlavor = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("180"))
	styleKeys   = lipgloss.NewStyle().Foreground(lipgloss.Color("117"))
	styleErr    = lipgloss.NewStyle().Foreground(lipgloss.Color("203"))
	styleBox    = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("240")).Padding(0, 1)
	styleSel    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("226"))
	styleLogo   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("214"))
)

func key(k, label string) string { return styleKeys.Render("["+k+"]") + " " + label }
