package styles

import "github.com/charmbracelet/lipgloss"

var (
	JobInfoBox = lipgloss.NewStyle().Padding(1, 2).BorderStyle(lipgloss.DoubleBorder()).BorderForeground(lipgloss.Color("#ffd700"))

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
)
