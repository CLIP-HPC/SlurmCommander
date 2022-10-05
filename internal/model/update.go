package model

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pja237/slurmcommander/internal/command"
	"github.com/pja237/slurmcommander/internal/keybindings"
	"github.com/pja237/slurmcommander/internal/model/tabs/jobfromtemplate"
	"github.com/pja237/slurmcommander/internal/model/tabs/jobtab"
	"github.com/pja237/slurmcommander/internal/slurm"
)

type errMsg error

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var (
		brk         bool = false
		activeTable *table.Model
	)

	// This shortens the testing for table movement keys
	switch m.ActiveTab {
	case tabJobs:
		activeTable = &m.SqueueTable
	case tabJobHist:
		activeTable = &m.SacctTable
	case tabJobFromTemplate:
		activeTable = &m.TemplatesTable
	case tabCluster:
		activeTable = &m.SinfoTable
	}

	// Filter is turned on, take care of this first
	// TODO: revisit this for filtering on multiple tabs
	switch {
	case m.FilterOn:
		switch msg := msg.(type) {

		case tea.KeyMsg:
			switch msg.Type {
			// TODO: when filter is set/cleared, trigger refresh with new filtered data
			case tea.KeyEnter:
				// finish entering filter
				m.FilterOn = false
				m.lastKey = "ENTER"
				brk = true
				// TODO: this is a "fix" for crashing-after-filter when Cursor() goes beyond list end
				// TODO: don't feel good about this... what if list is empty? no good. revisit
				m.SqueueTable.SetCursor(0)
			case tea.KeyEsc:
				// abort entering filter
				m.FilterOn = false
				m.Filter.SetValue("")
				m.lastKey = "ESC"
				brk = true
			}
			if brk {
				// TODO:
				// for now slash is enabled only in Jobs tab, but later we might enable it for cluster and others as well
				// so keep this part in...
				switch m.ActiveTab {
				case tabJobs:
					//return m, command.TimedGetSqueue()
					return m, command.QuickGetSqueue()
				case tabCluster:
					return m, command.TimedGetSinfo()
				}
			}
		}

		m.DebugMsg += "f"
		tmp, cmd := m.Filter.Update(msg)
		m.Filter = tmp
		return m, cmd

	case m.JobTab.MenuOn:
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.JobTab.Menu.SetWidth(msg.Width)
			return m, nil

		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "ctrl+c":
				//m.quitting = true
				return m, tea.Quit

			case "enter":
				m.JobTab.MenuOn = false
				m.JobTab.InfoOn = true
				i, ok := m.JobTab.Menu.SelectedItem().(jobtab.MenuItem)
				if ok {
					m.JobTab.MenuChoice = jobtab.MenuItem(i)
				}
				//return m, tea.Quit
				return m, nil
			}
		}

		var cmd tea.Cmd
		m.JobTab.Menu, cmd = m.JobTab.Menu.Update(msg)
		return m, cmd

	case m.EditTemplate:
		// TODO: move this code to a function/method
		var cmds []tea.Cmd
		var cmd tea.Cmd

		m.LogF.WriteString("U(): m.EditTemplate\n")
		switch msg := msg.(type) {
		case tea.KeyMsg:
			m.LogF.WriteString("U(): m.EditTemplate case tea.KeyMsg\n")
			switch msg.Type {
			case tea.KeyEsc:
				m.EditTemplate = false
				tabKeys[m.ActiveTab].SetupKeys()
				//if m.TemplateEditor.Focused() {
				//	m.TemplateEditor.Blur()
				//} else {
				//	m.EditTemplate = false
				//}
			case tea.KeyCtrlC:
				return m, tea.Quit
			default:
				if !m.TemplateEditor.Focused() {
					cmd = m.TemplateEditor.Focus()
					cmds = append(cmds, cmd)
				}
			}

		// We handle errors just like any other message
		case errMsg:
			//m.err = msg
			return m, nil
		}

		m.TemplateEditor, cmd = m.TemplateEditor.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)

	}

	switch msg := msg.(type) {

	// TODO: https://pkg.go.dev/github.com/charmbracelet/bubbletea#WindowSizeMsg
	// ToDo:
	// prevent updates for non-selected tabs

	// getting initial template text
	case jobfromtemplate.TemplateText:
		m.LogF.WriteString("U(): msg TemplateText")
		m.TemplateEditor = textarea.New()
		//m.TemplateEditor.Placeholder = string(msg)
		m.TemplateEditor.SetValue(string(msg))
		//m.TemplateEditor.Focus()

	// Windows resize
	case tea.WindowSizeMsg:
		m.winW = msg.Width
		m.winH = msg.Height
		m.SqueueTable.SetHeight(m.winH - 30)
		m.SacctTable.SetHeight(m.winH - 30)
		m.SinfoTable.SetHeight(m.winH - 30)

	// JobTab update
	case slurm.SqueueJSON:
		m.LogF.WriteString("U(): got SqueueJSON\n")
		if len(msg.Jobs) != 0 {
			m.Squeue = msg
			slurm.SqueueTabRows = nil
			// DONE: if there is no filter set, there is no need for all the string searching
			for _, v := range msg.Jobs {
				app := false
				if m.Filter.Value() != "" {
					switch {
					case strings.Contains(strconv.Itoa(*v.JobId), m.Filter.Value()):
						app = true
					case strings.Contains(*v.Name, m.Filter.Value()):
						app = true
					case strings.Contains(*v.Account, m.Filter.Value()):
						app = true
					case strings.Contains(*v.UserName, m.Filter.Value()):
						app = true
					case strings.Contains(*v.JobState, m.Filter.Value()):
						app = true
					}
				} else {
					app = true
				}
				if app {
					slurm.SqueueTabRows = append(slurm.SqueueTabRows, table.Row{strconv.Itoa(*v.JobId), *v.Name, *v.Account, *v.UserName, *v.JobState})
				}
			}
			// TODO:
			// fix: if after filtering m.table.Cursor|SelectedRow > lines in table, Info crashes trying to fetch nonexistent row
			m.SqueueTable.SetRows(slurm.SqueueTabRows)
			//m.SqueueTable.UpdateViewport()
		}
		m.UpdateCnt++
		// if active window != this, don't trigger new refresh
		m.DebugMsg += "J"
		if m.ActiveTab == tabJobs {
			m.DebugMsg += "2"
			return m, command.TimedGetSqueue()
		} else {
			m.DebugMsg += "3"
			return m, nil
		}

	// Cluster tab update
	case slurm.SinfoJSON:
		m.LogF.WriteString("U(): got SinfoJSON\n")
		if len(msg.Nodes) != 0 {
			m.Sinfo = msg
			slurm.SinfoTabRows = nil
			for _, v := range msg.Nodes {
				slurm.SinfoTabRows = append(slurm.SinfoTabRows, table.Row{*v.Name, *v.State, strconv.Itoa(*v.Cpus), strconv.FormatInt(*v.IdleCpus, 10), strconv.Itoa(*v.RealMemory), strconv.Itoa(*v.FreeMemory), strings.Join(*v.StateFlags, ",")})
			}
			m.SinfoTable.SetRows(slurm.SinfoTabRows)
		}
		m.UpdateCnt++
		// if active window != this, don't trigger new refresh
		m.DebugMsg += "C"
		if m.ActiveTab == tabCluster {
			m.DebugMsg += "4"
			return m, command.TimedGetSinfo()
		} else {
			m.DebugMsg += "5"
			return m, nil
		}

	// Job History tab update
	case slurm.SacctList:
		m.LogF.WriteString("U(): got SacctList\n")
		// fill out model
		m.DebugMsg += "H"
		m.JobHistTab.SacctList = msg
		// We do it only once on the start? or tick it?
		for _, v := range m.JobHistTab.SacctList {
			slurm.SacctTabRows = append(slurm.SacctTabRows, table.Row{v[0], v[1], v[2], v[3], v[4]})
		}
		m.JobHistTab.SacctTable.SetRows(slurm.SacctTabRows)
		return m, nil

	// Job Details tab update
	case slurm.SacctJob:
		m.LogF.WriteString("U(): got SacctJob\n")
		m.DebugMsg += "D"
		m.JobDetailsTab.SacctJob = msg
		return m, nil

	// TODO: find a way to simplify this mess below...
	// Keys pressed
	case tea.KeyMsg:
		switch {

		// UP
		// TODO: what if it's a list?
		case key.Matches(msg, keybindings.DefaultKeyMap.Up):
			activeTable.MoveUp(1)
			m.lastKey = "up"

		// DOWN
		case key.Matches(msg, keybindings.DefaultKeyMap.Down):
			activeTable.MoveDown(1)
			m.lastKey = "down"

		// PAGE DOWN
		case key.Matches(msg, keybindings.DefaultKeyMap.PageDown):
			activeTable.MoveDown(activeTable.Height())
			m.lastKey = "pgdown"

		// PAGE UP
		case key.Matches(msg, keybindings.DefaultKeyMap.PageUp):
			activeTable.MoveUp(activeTable.Height())
			m.lastKey = "pgup"

		// 1..6 Tab Selection keys
		case key.Matches(msg, keybindings.DefaultKeyMap.TtabSel):
			k, _ := strconv.Atoi(msg.String())
			m.ActiveTab = uint(k) - 1
			tabKeys[m.ActiveTab].SetupKeys()
			m.DebugMsg += "Ts"
			m.lastKey = msg.String()
			// TODO: needs triggering of the TimedGet*() like TAB key below
			return m, nil

		// TAB
		case key.Matches(msg, keybindings.DefaultKeyMap.Tab):
			// switch tab
			m.ActiveTab = (m.ActiveTab + 1) % uint(len(tabs))
			// setup keys
			tabKeys[m.ActiveTab].SetupKeys()
			m.lastKey = "tab"

			switch m.ActiveTab {
			case tabJobs:
				m.DebugMsg += "Tj"
				return m, command.TimedGetSqueue()

			case tabCluster:
				m.DebugMsg += "Tc"
				return m, command.TimedGetSinfo()
			}

		// SLASH
		case key.Matches(msg, keybindings.DefaultKeyMap.Slash):
			m.FilterOn = true
			m.DebugMsg += "/"
			return m, nil

		// ENTER
		case key.Matches(msg, keybindings.DefaultKeyMap.Enter):
			switch m.ActiveTab {
			case tabJobs:
				// Job Queue menu on
				m.JobTab.MenuOn = true
				m.JobTab.SelectedJob = m.JobTab.SqueueTable.SelectedRow()[0]
				// TODO: fill out menu with job options
				m.JobTab.SelectedJobState = m.JobTab.SqueueTable.SelectedRow()[4]
				menu := jobtab.MenuList[m.JobTab.SelectedJobState]
				m.LogF.WriteString(fmt.Sprintf("MENU %#v\n", jobtab.MenuList[m.JobTab.SelectedJobState]))
				m.JobTab.Menu = list.New(menu, list.NewDefaultDelegate(), 10, 10)
				//m.JobTab.Menu.Styles = list.DefaultStyles()
				m.JobTab.Menu.Title = "Job actions"
				m.JobTab.Menu.SetShowStatusBar(true)
				m.JobTab.Menu.SetFilteringEnabled(false)
				m.JobTab.Menu.SetShowHelp(false)
				m.JobTab.Menu.SetShowPagination(false)
				m.JobTab.Menu.SetHeight(30)
				m.JobTab.Menu.SetWidth(30)
				m.JobTab.Menu.SetSize(30, 30)

				return m, nil
			case tabJobHist:
				m.LogF.WriteString("U(): ENTER @ jobhist list\n")
				m.ActiveTab = tabJobDetails
				tabKeys[m.ActiveTab].SetupKeys()
				m.JobDetailsTab.SelJobID = m.JobHistTab.SacctTable.SelectedRow()[0]
				m.DebugMsg += "<-"
				return m, command.SingleJobGetSacct(m.JobDetailsTab.SelJobID)
			case tabJobFromTemplate:
				m.LogF.WriteString("U(): ENTER @ jobfromtemplate list\n")
				m.EditTemplate = true
				m.TemplateEditor = textarea.New()
				//m.TemplateEditor.Placeholder = string(jobfromtemplate.TemplateSample)
				m.TemplateEditor.Focus()
				m.TemplateEditor.SetValue(string(jobfromtemplate.TemplateSample))
				m.TemplateEditor.SetWidth(m.winW - 30)
				m.TemplateEditor.SetHeight(m.winH - 30)
				jobfromtemplate.EditorKeyMap.SetupKeys()
				// return & handle editing there
				return m, jobfromtemplate.GetTemplate("blabla")
			}

		// Info - toggle on/off
		case key.Matches(msg, keybindings.DefaultKeyMap.Info):
			if m.InfoOn {
				m.InfoOn = false
			} else {
				m.InfoOn = true
			}
			m.DebugMsg += "I"
			return m, nil

		// QUIT
		case key.Matches(msg, keybindings.DefaultKeyMap.Quit):
			fmt.Println("Quit key pressed")
			return m, tea.Quit
		}
	}

	return m, nil
}
