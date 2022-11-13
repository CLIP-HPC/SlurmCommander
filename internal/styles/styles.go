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
	red = lipgloss.Color("#cc0000")
	//red = lipgloss.Color("#b30000")
	//green = lipgloss.Color("#009900")
	green = lipgloss.Color("#00b300")

	// Generic text color styles
	TextRed    = lipgloss.NewStyle().Foreground(red)
	TextYellow = lipgloss.NewStyle().Foreground(yellow)
	TextGreen  = lipgloss.NewStyle().Foreground(green)
	TextBlue   = lipgloss.NewStyle().Foreground(blue)

	// ErrorHelp Box
	//ErrorHelp = lipgloss.NewStyle().Foreground(red).Border(lipgloss.RoundedBorder()).BorderForeground(red)
	ErrorHelp = lipgloss.NewStyle().Foreground(red)

	Tab = lipgloss.NewStyle().
		Border(TabTabBorder, true).
		BorderForeground(TabColor).
		Padding(0, 1)

	TabColor = lipgloss.AdaptiveColor{Light: "#0057B7", Dark: "#0057B7"}

	TabActiveTab       = Tab.Copy().Border(TabActiveTabBorder, true).Foreground(yellow)
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

	StatsSeparatorTitle = lipgloss.NewStyle().Foreground(yellow).Background(blue)

	//MenuTitleStyle    = lipgloss.NewStyle().Background(blue).Foreground(yellow)
	MenuBoxStyle      = lipgloss.NewStyle().Padding(1, 2).BorderStyle(lipgloss.DoubleBorder()).BorderForeground(blue)
	MenuTitleStyle    = lipgloss.NewStyle().Foreground(yellow)
	MenuNormalTitle   = lipgloss.NewStyle().Foreground(blue)
	MenuSelectedTitle = lipgloss.NewStyle().Foreground(yellow).Background(blue)
	MenuNormalDesc    = lipgloss.NewStyle().Foreground(yellow).Background(blue)
	MenuSelectedDesc  = lipgloss.NewStyle().Foreground(yellow)

	// Main Window area
	//MainWindow = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(blue).MaxHeight(80)
	MainWindow = lipgloss.NewStyle().MaxHeight(80)
	HelpWindow = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(yellow).Height(2).MaxHeight(4)
	NotFound   = lipgloss.NewStyle().Foreground(red)

	// JobQueue tab, infobox
	//JobInfoBox   = lipgloss.NewStyle().Padding(1, 1).BorderStyle(lipgloss.DoubleBorder()).BorderForeground(blue)
	//JobInfoBox = lipgloss.NewStyle().BorderStyle(lipgloss.DoubleBorder()).BorderForeground(blue)
	JobInfoBox = lipgloss.NewStyle()
	//JobInfoInBox = lipgloss.NewStyle().Padding(1, 1).BorderStyle(lipgloss.RoundedBorder()).BorderForeground(blue)
	JobInfoInBox       = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(blue).MaxHeight(7)
	JobInfoInBottomBox = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(blue).MaxHeight(7)

	// Job steps
	JobStepBoxStyle = lipgloss.NewStyle().Padding(1, 2).BorderStyle(lipgloss.DoubleBorder()).BorderForeground(blue)
	//JobStepExitStatusRed = lipgloss.NewStyle().Padding(0, 0).BorderStyle(lipgloss.RoundedBorder()).BorderForeground(red)
	JobStepExitStatusRed   = lipgloss.NewStyle().Foreground(red)
	JobStepExitStatusGreen = lipgloss.NewStyle().Foreground(green)
	TresBox                = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(blue).Width(40)
)
