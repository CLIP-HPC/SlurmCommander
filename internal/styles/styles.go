package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Blue
	blue = lipgloss.Color("#0057b7")
	// Yellow
	yellow = lipgloss.Color("#ffd700")
	// Red
	red = lipgloss.Color("#ff0000")

	JobInfoBox = lipgloss.NewStyle().Padding(1, 1).BorderStyle(lipgloss.DoubleBorder()).BorderForeground(blue)

	Tab = lipgloss.NewStyle().
		Border(TabTabBorder, true).
		BorderForeground(TabColor).
		Padding(0, 1)

	TabColor = lipgloss.AdaptiveColor{Light: "#0057B7", Dark: "#0057B7"}

	TabActiveTab       = Tab.Copy().Border(TabActiveTabBorder, true)
	TabActiveTabBorder = lipgloss.ThickBorder()
	//TabActiveTabBorder= lipgloss.Border{
	//	Top:         "─",
	//	Top:         "/",
	//	Bottom:      " ",
	//	Left:        "│",
	//	Right:       "│",
	//	TopLeft:     "╭",
	//	TopRight:    "╮",
	//	BottomLeft:  "┘",
	//	BottomRight: "└",
	//}

	TabTabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}

	TabGap = Tab.Copy().
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false)

	//MenuTitleStyle    = lipgloss.NewStyle().Background(blue).Foreground(yellow)
	MenuBoxStyle      = lipgloss.NewStyle().Padding(1, 2).BorderStyle(lipgloss.DoubleBorder()).BorderForeground(blue)
	MenuTitleStyle    = lipgloss.NewStyle().Foreground(yellow)
	MenuNormalTitle   = lipgloss.NewStyle().Foreground(blue)
	MenuSelectedTitle = lipgloss.NewStyle().Foreground(yellow).Background(blue)
	MenuNormalDesc    = lipgloss.NewStyle().Foreground(yellow).Background(blue)
	MenuSelectedDesc  = lipgloss.NewStyle().Foreground(yellow)

	NotFound = lipgloss.NewStyle().Foreground(red)
)
