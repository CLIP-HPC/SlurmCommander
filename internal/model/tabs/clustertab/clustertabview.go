package clustertab

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"
	"github.com/pja237/slurmcommander-dev/internal/generic"
	"github.com/pja237/slurmcommander-dev/internal/styles"
)

func (ct *ClusterTab) tabCluster() string {

	scr := ct.SinfoTable.View() + "\n"

	return scr
}

func (ct *ClusterTab) tabClusterBars(l *log.Logger) string {
	var (
		scr     string = ""
		cpuPerc float64
		memPerc float64
	)

	sel := ct.SinfoTable.Cursor()
	l.Printf("ClusterTab Selected: %d\n", sel)
	l.Printf("ClusterTab len results: %d\n", len(ct.SinfoFiltered.Nodes))
	ct.CpuBar = progress.New(progress.WithGradient("#277BC0", "#FFCB42"))
	ct.MemBar = progress.New(progress.WithGradient("#277BC0", "#FFCB42"))
	if len(ct.SinfoFiltered.Nodes) > 0 && sel != -1 {
		cpuPerc = float64(*ct.SinfoFiltered.Nodes[sel].AllocCpus) / float64(*ct.SinfoFiltered.Nodes[sel].Cpus)
		memPerc = float64(*ct.SinfoFiltered.Nodes[sel].AllocMemory) / float64(*ct.SinfoFiltered.Nodes[sel].RealMemory)

		scr += fmt.Sprintf("CPU used/total: %d/%d\n", *ct.SinfoFiltered.Nodes[sel].AllocCpus, *ct.SinfoFiltered.Nodes[sel].Cpus)
		scr += ct.CpuBar.ViewAs(cpuPerc)
		scr += "\n"
		scr += fmt.Sprintf("MEM used/total: %d/%d\n", *ct.SinfoFiltered.Nodes[sel].AllocMemory, *ct.SinfoFiltered.Nodes[sel].RealMemory)
		scr += ct.MemBar.ViewAs(memPerc)
		scr += "\n\n"
	} else {
		cpuPerc = 0
		memPerc = 0
		scr += fmt.Sprintf("CPU used/total: %d/%d\n", 0, 0)
		scr += ct.CpuBar.ViewAs(cpuPerc)
		scr += "\n"
		scr += fmt.Sprintf("MEM used/total: %d/%d\n", 0, 0)
		scr += ct.MemBar.ViewAs(memPerc)
		scr += "\n\n"

	}

	return scr
}

func (ct *ClusterTab) ClusterTabStats(l *log.Logger) string {
	var str string

	l.Printf("JobClusterTabStats called\n")

	sel := ct.SinfoTable.Cursor()
	str += styles.StatsSeparatorTitle.Render(fmt.Sprintf("%-30s", "Nodes states (filtered):"))
	str += "\n"

	if len(ct.SinfoFiltered.Nodes) > 0 {
		//str += generic.GenCountStrVert(m.JobClusterTab.Stats.StateCnt, m.Log)
		str += generic.GenCountStrVert(ct.Stats.StateSimpleCnt, l)
	}

	str += styles.StatsSeparatorTitle.Render(fmt.Sprintf("%-30s", "Selected node:"))

	if len(ct.SinfoFiltered.Nodes) > 0 && sel != -1 {
		str += "\n"
		str += fmt.Sprintf("%-15s: %s\n", "Arch", *ct.SinfoFiltered.Nodes[sel].Architecture)
		str += fmt.Sprintf("%-15s: %s\n", "Features", *ct.SinfoFiltered.Nodes[sel].ActiveFeatures)
		str += fmt.Sprintf("%-15s: %s\n", "TRES", *ct.SinfoFiltered.Nodes[sel].Tres)
		if ct.SinfoFiltered.Nodes[sel].TresUsed != nil {
			str += fmt.Sprintf("%-15s: %s\n", "TRES Used", *ct.SinfoFiltered.Nodes[sel].TresUsed)
		} else {
			str += fmt.Sprintf("%-15s: %s\n", "TRES Used", "")
		}
		str += fmt.Sprintf("%-15s: %s\n", "GRES", *ct.SinfoFiltered.Nodes[sel].Gres)
		str += fmt.Sprintf("%-15s: %s\n", "GRES Used", *ct.SinfoFiltered.Nodes[sel].GresUsed)
		str += fmt.Sprintf("%-15s: %s\n", "Partitions", strings.Join(*ct.SinfoFiltered.Nodes[sel].Partitions, ","))
	}
	return str
}

func (ct *ClusterTab) getClusterCounts() string {
	var (
		ret string
		cpp string
		mpp string
		nps string
	)

	fmtStrCpu := "%-10s : %4d / %4d %2.0f%%\n"
	fmtStrMem := "%-10s : %8d / %8d %2.0f%%\n"
	fmtStrNPS := "%-15s : %4d\n"
	fmtTitle := "%-40s"

	cpp += styles.TextYellowOnBlue.Render(fmt.Sprintf(fmtTitle, "CPUs per Partition (used/total)"))
	cpp += "\n"
	for _, v := range ct.Breakdowns.CpuPerPart {
		cpp += fmt.Sprintf(fmtStrCpu, v.Name, v.Count, v.Total, float32(v.Count)/float32(v.Total)*100)
	}

	mpp += styles.TextYellowOnBlue.Render(fmt.Sprintf(fmtTitle, "Mem per Partition (used/total)"))
	mpp += "\n"
	for _, v := range ct.Breakdowns.MemPerPart {
		mpp += fmt.Sprintf(fmtStrMem, v.Name, v.Count, v.Total, float32(v.Count)/float32(v.Total)*100)
	}

	nps += styles.TextYellowOnBlue.Render(fmt.Sprintf(fmtTitle, "Nodes per State"))
	nps += "\n"
	for _, v := range ct.Breakdowns.NodesPerState {
		nps += fmt.Sprintf(fmtStrNPS, v.Name, v.Count)
	}

	cpp = styles.CountsBox.Render(cpp)
	mpp = styles.CountsBox.Render(mpp)
	nps = styles.CountsBox.Render(nps)

	ret = lipgloss.JoinHorizontal(lipgloss.Top, cpp, mpp, nps)

	return ret
}

func (ct *ClusterTab) View(l *log.Logger) string {
	var (
		MainWindow strings.Builder
	)

	// Top Main
	MainWindow.WriteString(fmt.Sprintf("Filter: %10.20s\tItems: %d\n\n", ct.Filter.Value(), len(ct.SinfoFiltered.Nodes)))
	MainWindow.WriteString(ct.tabClusterBars(l))

	// Mid Main: table || table+stats
	switch {
	case ct.StatsOn:
		MainWindow.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, ct.tabCluster(), styles.MenuBoxStyle.Render(ct.ClusterTabStats(l))))
	default:
		MainWindow.WriteString(ct.tabCluster())
	}

	// Low Main: nil || filter || counts
	switch {
	case ct.FilterOn:
		// filter
		MainWindow.WriteString("\n")
		// filter
		MainWindow.WriteString("\n")
		MainWindow.WriteString("Filter value (search across: Name, State, StateFlags!):\n")
		MainWindow.WriteString(fmt.Sprintf("%s\n", ct.Filter.View()))
		MainWindow.WriteString("(Enter to apply, Esc to clear filter and abort, Regular expressions supported, syntax details: https://golang.org/s/re2syntax)\n")
	case ct.CountsOn:
		MainWindow.WriteString("\n")
		MainWindow.WriteString(styles.JobInfoBox.Render(ct.getClusterCounts()))

	default:
		MainWindow.WriteString("\n")
		MainWindow.WriteString(generic.GenCountStr(ct.Stats.StateCnt, l))
	}

	return MainWindow.String()
}
