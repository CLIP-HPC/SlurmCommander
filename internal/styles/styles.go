package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	Blue   = lipgloss.Color("#0057b7")
	Yellow = lipgloss.Color("#ffd700")
	Red    = lipgloss.Color("#cc0000")
	Green  = lipgloss.Color("#00b300")
	Orange = lipgloss.Color("#FFA500")

	Bluegrey = lipgloss.Color("#c2d1f0")

	// Generic text color styles
	TextRed          = lipgloss.NewStyle().Foreground(Red)
	TextYellow       = lipgloss.NewStyle().Foreground(Yellow)
	TextGreen        = lipgloss.NewStyle().Foreground(Green)
	TextBlue         = lipgloss.NewStyle().Foreground(Blue)
	TextOrange       = lipgloss.NewStyle().Foreground(Orange)
	TextBlueGrey     = lipgloss.NewStyle().Foreground(Bluegrey)
	TextYellowOnBlue = lipgloss.NewStyle().Foreground(Yellow).Background(Blue).Underline(true)

	// Table styles
	//SelectedRow = lipgloss.NewStyle().Background(Blue).Foreground(Yellow).Bold(false)
	SelectedRow = lipgloss.NewStyle().Background(Blue).Foreground(Yellow)

	// ErrorHelp Box
	//ErrorHelp = lipgloss.NewStyle().Foreground(red).Border(lipgloss.RoundedBorder()).BorderForeground(red)
	ErrorHelp = lipgloss.NewStyle().Foreground(Red)

	// TABS
	Tab = lipgloss.NewStyle().
		Border(TabTabBorder, true).
		BorderForeground(TabColor).
		Padding(0, 1)
	TabColor           = lipgloss.AdaptiveColor{Light: "#0057B7", Dark: "#0057B7"}
	TabActiveTab       = Tab.Copy().Border(TabActiveTabBorder, true).Foreground(Yellow)
	TabActiveTabBorder = lipgloss.ThickBorder()
	TabTabBorder       = lipgloss.Border{
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

	// (S)tats Box Style
	StatsBoxStyle       = lipgloss.NewStyle().Padding(0, 1).BorderStyle(lipgloss.DoubleBorder()).BorderForeground(Blue)
	StatsSeparatorTitle = lipgloss.NewStyle().Foreground(Yellow).Background(Blue)

	// JobDetails viewport box
	//JDviewportBox = lipgloss.NewStyle().Border(lipgloss.DoubleBorder(), true, false).BorderForeground(Yellow).Padding(1, 1)
	JDviewportBox = lipgloss.NewStyle()

	// ClusterTab Stats Box
	ClusterTabStats = StatsBoxStyle.Copy()

	//MenuTitleStyle    = lipgloss.NewStyle().Background(blue).Foreground(yellow)
	MenuBoxStyle      = lipgloss.NewStyle().Padding(1, 1).BorderStyle(lipgloss.DoubleBorder()).BorderForeground(Blue)
	MenuTitleStyle    = lipgloss.NewStyle().Foreground(Yellow)
	MenuNormalTitle   = lipgloss.NewStyle().Foreground(Blue)
	MenuSelectedTitle = lipgloss.NewStyle().Foreground(Yellow).Background(Blue)
	MenuNormalDesc    = lipgloss.NewStyle().Foreground(Yellow).Background(Blue)
	MenuSelectedDesc  = lipgloss.NewStyle().Foreground(Yellow)

	CountsBox = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 0).BorderForeground(Blue)

	// Main Window area
	MainWindow = lipgloss.NewStyle().MaxHeight(80)
	HelpWindow = lipgloss.NewStyle().Padding(0, 0).Border(lipgloss.RoundedBorder(), true, false, false).Height(2).MaxHeight(3).BorderForeground(Blue)

	// JobTemplates, template not found
	NotFound = lipgloss.NewStyle().Foreground(Red)

	// JobQueue tab, infobox
	JobInfoBox         = lipgloss.NewStyle()
	JobInfoInBox       = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(Blue).MaxHeight(7)
	JobInfoInBottomBox = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(Blue).MaxHeight(7)

	// JobDetails tab

	// Job steps
	JobStepBoxStyle        = lipgloss.NewStyle().Padding(1, 2).BorderStyle(lipgloss.RoundedBorder()).BorderForeground(Blue)
	JobStepExitStatusRed   = lipgloss.NewStyle().Foreground(Red)
	JobStepExitStatusGreen = lipgloss.NewStyle().Foreground(Green)

	//TresBox = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(blue).Width(40)
	TresBox = lipgloss.NewStyle()
)
