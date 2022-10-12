package model

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pja237/slurmcommander/internal/command"
	"github.com/pja237/slurmcommander/internal/keybindings"
	"github.com/pja237/slurmcommander/internal/model/tabs/jobfromtemplate"
	"github.com/pja237/slurmcommander/internal/model/tabs/jobtab"
	"github.com/pja237/slurmcommander/internal/slurm"
)

type errMsg error

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var (
		brk          bool = false
		activeTable  *table.Model
		activeFilter *textinput.Model
	)

	// This shortens the testing for table movement keys
	switch m.ActiveTab {
	case tabJobs:
		activeTable = &m.SqueueTable
		activeFilter = &m.JobTab.Filter
	case tabJobHist:
		activeTable = &m.SacctTable
		activeFilter = &m.JobHistTab.Filter
	case tabJobFromTemplate:
		activeTable = &m.TemplatesTable
	case tabCluster:
		activeTable = &m.SinfoTable
		activeFilter = &m.JobClusterTab.Filter
	}

	// Filter is turned on, take care of this first
	// TODO: revisit this for filtering on multiple tabs
	switch {
	case m.FilterSwitch != -1:
		m.Log.Printf("Update: In filter %d\n", m.FilterSwitch)
		switch msg := msg.(type) {

		case tea.KeyMsg:
			switch msg.Type {
			// TODO: when filter is set/cleared, trigger refresh with new filtered data
			case tea.KeyEnter:
				// finish & apply entering filter
				m.FilterSwitch = -1
				m.lastKey = "ENTER"
				brk = true
			case tea.KeyEsc:
				// abort entering filter
				m.FilterSwitch = -1
				activeFilter.SetValue("")
				m.lastKey = "ESC"
				brk = true
			}
			if brk {
				// TODO: this is a "fix" for crashing-after-filter when Cursor() goes beyond list end
				// TODO: don't feel good about this... what if list is empty? no good. revisit
				// NOTE: This doesn't do what i image it should, cursor remains -1 when table is empty situation?
				// Explanation in clamp function: https://github.com/charmbracelet/bubbles/blob/13f52d678d315676568a656b5211b8a24a54a885/table/table.go#L296
				activeTable.SetCursor(0)
				//m.Log.Printf("ActiveTable = %v\n", activeTable)
				m.Log.Printf("Update: Filter set, setcursor(0), activetable.Cursor==%d\n", activeTable.Cursor())
				switch m.ActiveTab {
				case tabJobs:
					return m, command.QuickGetSqueue()
				case tabJobHist:
					return m, command.QuickGetSacct()
				case tabCluster:
					return m, command.QuickGetSinfo()
				default:
					return m, nil
				}
			}
		}

		m.DebugMsg += "f"
		tmp, cmd := activeFilter.Update(msg)
		*activeFilter = tmp
		return m, cmd

	case m.JobTab.MenuOn:
		m.Log.Printf("Update: In Menu\n")
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
				// TODO: This is just temporarily here, instead of this, depending on the MenuChoice turn on Info if selected
				m.JobTab.InfoOn = true
				i, ok := m.JobTab.Menu.SelectedItem().(jobtab.MenuItem)
				if ok {
					m.JobTab.MenuChoice = jobtab.MenuItem(i)
					// TODO: here, we should call a Cmd function that does something depending on the MenuChoice and returns a msg
					// e.g. for "CANCEL", call scancelCmd and return msg type cancelSent on which queue refresh should be triggered
					// Also, get and pass in the Selected() jobID
					m.JobTab.MenuChoice.ExecMenuItem(m.JobTab.SelectedJob, m.Log)
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

		m.Log.Printf("Update: In EditTemplate: %#v\n", msg)
		switch msg := msg.(type) {
		case tea.KeyMsg:
			m.Log.Printf("Update: m.EditTemplate case tea.KeyMsg\n")
			switch msg.Type {
			case tea.KeyEsc:
				m.EditTemplate = false
				tabKeys[m.ActiveTab].SetupKeys()
				//if m.TemplateEditor.Focused() {
				//	m.TemplateEditor.Blur()
				//} else {
				//	m.EditTemplate = false
				//}

			case tea.KeyCtrlS:
				// TODO:
				// 1. Exit editor
				// 2. Save content to file
				// 3. Notify user about generated filename from 2.
				// 4. Submit job
				m.Log.Printf("EditTemplate: Ctrl+s pressed\n")
				m.EditTemplate = false
				tabKeys[m.ActiveTab].SetupKeys()
				jobfromtemplate.SaveToFile(m.JobFromTemplateTab.TemplatesTable.SelectedRow()[0], m.JobFromTemplateTab.TemplateEditor.Value(), m.Log)
				return m, nil

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

	// Get initial job template list
	case jobfromtemplate.TemplatesListRows:
		m.Log.Printf("Update: Got TemplatesListRows msg: %#v\n", msg)
		if msg != nil {
			// if it's not empty, append to table
			m.JobFromTemplateTab.TemplatesTable.SetRows(msg)
		}
		return m, nil

	// getting initial template text
	case jobfromtemplate.TemplateText:
		m.Log.Printf("Update: Got TemplateText msg: %#v\n", msg)
		// HERE: we initialize the new textarea editor and flip the EditTemplate switch to ON
		jobfromtemplate.EditorKeyMap.SetupKeys()
		m.EditTemplate = true
		m.TemplateEditor = textarea.New()
		m.TemplateEditor.SetWidth(m.winW - 30)
		m.TemplateEditor.SetHeight(m.winH - 30)
		m.TemplateEditor.SetValue(string(msg))
		m.TemplateEditor.Focus()
		return m, jobfromtemplate.EditorOn()

	// Windows resize
	case tea.WindowSizeMsg:
		m.winW = msg.Width
		m.winH = msg.Height
		m.Log.Printf("Update: got WindowSizeMsg: %d %d\n", msg.Width, msg.Height)
		// Tabs :  3
		// Header  3
		// TABLE:  X
		// Debug:  5
		// Filter: 3
		// Help :  1
		// ---
		// TOTAL:  15
		m.SqueueTable.SetHeight(m.winH - 30)
		m.SacctTable.SetHeight(m.winH - 30)
		m.SinfoTable.SetHeight(m.winH - 30)

	// JobTab update
	case slurm.SqueueJSON:
		m.Log.Printf("U(): got SqueueJSON\n")
		if len(msg.Jobs) != 0 {
			m.Squeue = msg

			// TODO:
			// fix: if after filtering m.table.Cursor|SelectedRow > lines in table, Info crashes trying to fetch nonexistent row
			rows, sqf := msg.FilterSqueueTable(m.JobTab.Filter.Value())
			m.JobTab.SqueueTable.SetRows(rows)
			m.JobTab.SqueueFiltered = sqf
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
		m.Log.Printf("U(): got SinfoJSON\n")
		if len(msg.Nodes) != 0 {
			m.Sinfo = msg
			//slurm.SinfoTabRows = nil
			//for _, v := range msg.Nodes {
			//	slurm.SinfoTabRows = append(slurm.SinfoTabRows, table.Row{*v.Name, *v.State, strconv.Itoa(*v.Cpus), strconv.FormatInt(*v.IdleCpus, 10), strconv.Itoa(*v.RealMemory), strconv.Itoa(*v.FreeMemory), strings.Join(*v.StateFlags, ",")})
			//}
			rows, sif := msg.FilterSinfoTable(m.JobClusterTab.Filter.Value())
			m.JobClusterTab.SinfoTable.SetRows(rows)
			m.JobClusterTab.SinfoFiltered = sif
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
		m.Log.Printf("U(): got SacctList\n")
		// fill out model
		m.DebugMsg += "H"
		m.JobHistTab.SacctList = msg
		//m.JobHistTab.SacctTable.SetRows(msg.FilterSacctTable(m.JobHistTab.Filter.Value()))
		rows, saf := msg.FilterSacctTable(m.JobHistTab.Filter.Value())
		m.JobHistTab.SacctTable.SetRows(rows)
		m.JobHistTab.SacctListFiltered = saf
		//m.LogF.WriteString(fmt.Sprintf("U(): got Filtered rows %#v\n", msg.FilterSacctTable(m.JobHistTab.Filter.Value())))
		return m, nil

	// Job Details tab update
	case slurm.SacctJob:
		m.Log.Printf("U(): got SacctJob\n")
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
			m.FilterSwitch = FilterSwitch(m.ActiveTab)
			m.DebugMsg += "/"
			return m, nil

		// ENTER
		case key.Matches(msg, keybindings.DefaultKeyMap.Enter):
			switch m.ActiveTab {

			// Job Queue tab: Open Job menu
			case tabJobs:
				n := m.JobTab.SqueueTable.Cursor()
				m.Log.Printf("Update ENTER key @ jobqueue table\n")
				// Job Queue menu on
				if n == -1 || len(m.JobTab.SqueueFiltered.Jobs) == 0 {
					m.Log.Printf("Update ENTER key @ jobqueue table, no jobs selected/empty table\n")
					return m, nil
				}
				m.JobTab.MenuOn = true
				m.JobTab.SelectedJob = m.JobTab.SqueueTable.SelectedRow()[0]
				// TODO: fill out menu with job options
				m.JobTab.SelectedJobState = m.JobTab.SqueueTable.SelectedRow()[4]
				menu := jobtab.MenuList[m.JobTab.SelectedJobState]
				m.Log.Printf("MENU %#v\n", jobtab.MenuList[m.JobTab.SelectedJobState])
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
				m.JobTab.Menu.Styles.Title = lipgloss.NewStyle().Background(lipgloss.Color("#0057b7")).Foreground(lipgloss.Color("#ffd700"))
				return m, nil

			// Job History tab: Select Job from history and open its Details tab
			case tabJobHist:
				n := m.JobHistTab.SacctTable.Cursor()
				m.Log.Printf("Update ENTER key @ jobhist table, cursor=%d, len=%d\n", n, len(m.JobHistTab.SacctListFiltered))
				if n == -1 || len(m.JobHistTab.SacctListFiltered) == 0 {
					m.Log.Printf("Update ENTER key @ jobhist table, no jobs selected/empty table\n")
					return m, nil
				}
				m.ActiveTab = tabJobDetails
				tabKeys[m.ActiveTab].SetupKeys()
				m.JobDetailsTab.SelJobID = m.JobHistTab.SacctTable.SelectedRow()[0]
				m.DebugMsg += "<-"
				return m, command.SingleJobGetSacct(m.JobDetailsTab.SelJobID)

			// Job from Template tab: Open template for editing
			case tabJobFromTemplate:
				m.Log.Printf("Update ENTER key @ jobfromtemplate table\n")
				// return & handle editing there
				return m, jobfromtemplate.GetTemplate(m.JobFromTemplateTab.TemplatesTable.SelectedRow()[2], m.Log)
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
