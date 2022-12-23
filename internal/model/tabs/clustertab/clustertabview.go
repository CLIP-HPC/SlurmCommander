package clustertab

import (
	"fmt"
	"log"
	"strings"

	"github.com/CLIP-HPC/SlurmCommander/internal/generic"
	"github.com/CLIP-HPC/SlurmCommander/internal/slurm"
	"github.com/CLIP-HPC/SlurmCommander/internal/styles"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
)

func (ct *ClusterTab) tabCluster() string {

	scr := ct.SinfoTable.View() + "\n"

	return scr
}

func (ct *ClusterTab) tabClusterBars(l *log.Logger) string {
	var (
		scr      string  = ""
		cpuPerc  float64 = 0
		cpuUsed  int64   = 0
		cpuAvail int     = 0
		memPerc  float64 = 0
		memUsed  int64   = 0
		memAvail int     = 0
		gpuPerc  float64 = 0
		gpuUsed  int     = 0
		gpuAvail int     = 0
	)

	sel := ct.SinfoTable.Cursor()
	l.Printf("ClusterTab Selected: %d\n", sel)
	l.Printf("ClusterTab len results: %d\n", len(ct.SinfoFiltered.Nodes))
	ct.CpuBar = progress.New(progress.WithGradient("#277BC0", "#FFCB42"))
	ct.MemBar = progress.New(progress.WithGradient("#277BC0", "#FFCB42"))
	ct.GpuBar = progress.New(progress.WithGradient("#277BC0", "#FFCB42"))
	if len(ct.SinfoFiltered.Nodes) > 0 && sel != -1 {
		cpuUsed = *ct.SinfoFiltered.Nodes[sel].AllocCpus
		cpuAvail = *ct.SinfoFiltered.Nodes[sel].Cpus
		cpuPerc = float64(cpuUsed) / float64(cpuAvail)
		memUsed = *ct.SinfoFiltered.Nodes[sel].AllocMemory
		memAvail = *ct.SinfoFiltered.Nodes[sel].RealMemory
		memPerc = float64(memUsed) / float64(memAvail)
		gpuAvail = *slurm.ParseGRES(*ct.SinfoFiltered.Nodes[sel].Gres)
		gpuUsed = *slurm.ParseGRES(*ct.SinfoFiltered.Nodes[sel].GresUsed)
		if gpuAvail > 0 {
			gpuPerc = float64(gpuUsed) / float64(gpuAvail)
		}
	}
	cpur := lipgloss.JoinVertical(lipgloss.Left, fmt.Sprintf("CPU used/total: %d/%d", cpuUsed, cpuAvail), ct.CpuBar.ViewAs(cpuPerc))
	memr := lipgloss.JoinVertical(lipgloss.Left, fmt.Sprintf("MEM used/total: %d/%d", memUsed, memAvail), ct.MemBar.ViewAs(memPerc))
	scr += lipgloss.JoinVertical(lipgloss.Top, cpur, memr)
	if gpuAvail > 0 {
		gpur := lipgloss.JoinVertical(lipgloss.Left, fmt.Sprintf("GPU used/total: %d/%d", gpuUsed, gpuAvail), ct.GpuBar.ViewAs(gpuPerc))
		scr = lipgloss.JoinHorizontal(lipgloss.Top, scr, fmt.Sprintf("%4s", ""), gpur)
	}
	scr += "\n\n"
	return scr
}

func (ct *ClusterTab) ClusterTabStats(l *log.Logger) string {
	var str string

	l.Printf("JobClusterTabStats called\n")

	sel := ct.SinfoTable.Cursor()
	//str += styles.StatsSeparatorTitle.Render(fmt.Sprintf("%-30s", "Nodes states (filtered):"))
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
		gpp string
		nps string
	)

	fmtStrCpu := "%-8s : %4d / %4d %2.0f%%\n"
	fmtStrMem := "%-8s : %s / %s %2.0f%%\n"
	fmtStrGpu := "%-8s : %4d / %4d %2.0f%%\n"
	fmtStrNPS := "%-15s : %4d\n"
	fmtTitle := "%-30s"

	cpp += styles.TextYellowOnBlue.Render(fmt.Sprintf(fmtTitle, "CPUs per Partition (used/total)"))
	cpp += "\n"
	for _, v := range ct.Breakdowns.CpuPerPart {
		cpp += fmt.Sprintf(fmtStrCpu, v.Name, v.Count, v.Total, float32(v.Count)/float32(v.Total)*100)
	}

	mpp += styles.TextYellowOnBlue.Render(fmt.Sprintf(fmtTitle, "Mem per Partition (used/total)"))
	mpp += "\n"
	for _, v := range ct.Breakdowns.MemPerPart {
		mpp += fmt.Sprintf(fmtStrMem, v.Name, humanize.Bytes(uint64(v.Count)*1024*1024), humanize.Bytes(uint64(v.Total)*1024*1024), float32(v.Count)/float32(v.Total)*100)
	}

	gpp += styles.TextYellowOnBlue.Render(fmt.Sprintf(fmtTitle, "GPUs per Partition (used/total)"))
	gpp += "\n"
	for _, v := range ct.Breakdowns.GpuPerPart {
		var gpuPerc float32 = 0.0
		if v.Total > 0 {
			gpuPerc = float32(v.Count) / float32(v.Total) * 100
		}
		gpp += fmt.Sprintf(fmtStrGpu, v.Name, v.Count, v.Total, gpuPerc)
	}

	nps += styles.TextYellowOnBlue.Render(fmt.Sprintf("%-30s", "Nodes per State"))
	nps += "\n"
	for _, v := range ct.Breakdowns.NodesPerState {
		nps += fmt.Sprintf(fmtStrNPS, v.Name, v.Count)
	}

	cpp = styles.CountsBox.Render(cpp)
	mpp = styles.CountsBox.Render(mpp)
	gpp = styles.CountsBox.Render(gpp)
	nps = styles.CountsBox.Render(nps)

	ret = lipgloss.JoinHorizontal(lipgloss.Top, cpp, mpp, gpp, nps)

	return ret
}

func (ct *ClusterTab) View(l *log.Logger) string {
	var (
		Header       strings.Builder
		MainWindow   strings.Builder
		FooterWindow strings.Builder
	)

	// Top Main
	Header.WriteString(fmt.Sprintf("Filter: %10.20s\tItems: %d\n\n", ct.Filter.Value(), len(ct.SinfoFiltered.Nodes)))
	Header.WriteString(ct.tabClusterBars(l))

	// Table is always there
	MainWindow.WriteString(ct.tabCluster())

	// Attach below table whatever is turned on
	switch {
	case ct.FilterOn:
		// filter
		FooterWindow.WriteString("\n")
		FooterWindow.WriteString("Filter value (search in joined: Name + Partition + State + StateFlags!):\n")
		FooterWindow.WriteString(fmt.Sprintf("%s\n", ct.Filter.View()))
		FooterWindow.WriteString("(Enter to apply, Esc to clear filter and abort, Regular expressions supported.\n")
		FooterWindow.WriteString(" Syntax details: https://golang.org/s/re2syntax)\n")
	case ct.CountsOn:
		FooterWindow.WriteString("\n")
		FooterWindow.WriteString(styles.JobInfoBox.Render(ct.getClusterCounts()))

	default:
		FooterWindow.WriteString("\n")
		//MainWindow.WriteString(generic.GenCountStr(ct.Stats.StateCnt, l))
	}

	// Lastly, if stats are on, horizontally join them to main
	switch {
	case ct.StatsOn:
		X := MainWindow.String()
		MainWindow.Reset()
		// TODO: make this Width() somewhere else (e.g. Update() on WindowSizeMsg)
		// Table Width == 118 chars, so .Width(m.winW-118)
		//MainWindow.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, X, styles.StatsBoxStyle.Width(50).Render(ct.ClusterTabStats(l))))
		l.Printf("CTB Width = %d\n", styles.ClusterTabStats.GetWidth())
		MainWindow.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, X, styles.ClusterTabStats.Render(ct.ClusterTabStats(l))))
	}

	return Header.String() + MainWindow.String() + FooterWindow.String()
}
