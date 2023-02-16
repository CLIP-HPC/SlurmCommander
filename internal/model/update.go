package model

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/CLIP-HPC/SlurmCommander/internal/command"
	"github.com/CLIP-HPC/SlurmCommander/internal/keybindings"
	"github.com/CLIP-HPC/SlurmCommander/internal/model/tabs/clustertab"
	"github.com/CLIP-HPC/SlurmCommander/internal/model/tabs/jobfromtemplate"
	"github.com/CLIP-HPC/SlurmCommander/internal/model/tabs/jobhisttab"
	"github.com/CLIP-HPC/SlurmCommander/internal/model/tabs/jobtab"
	"github.com/CLIP-HPC/SlurmCommander/internal/slurm"
	"github.com/CLIP-HPC/SlurmCommander/internal/generic"
	"github.com/CLIP-HPC/SlurmCommander/internal/styles"
	"github.com/CLIP-HPC/SlurmCommander/internal/table"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type errMsg error

type activeTabType interface {
	AdjTableHeight(int, *log.Logger)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var (
		brk                 bool = false
		activeTab           activeTabType
		activeTable         *table.Model
		activeFilter        *textinput.Model
		activeFilterOn      *bool
		activeUserInputs    *generic.UserInputs
		activeUserInputsOn  *bool
		activeJDViewport    bool
	)

	// This shortens the testing for table movement keys
	switch m.ActiveTab {
	case tabJobs:
		activeTab = &m.JobTab
		activeTable = &m.JobTab.SqueueTable
		activeFilter = &m.JobTab.Filter
		activeFilterOn = &m.JobTab.FilterOn
	case tabJobHist:
		activeTab = &m.JobHistTab
		activeTable = &m.JobHistTab.SacctTable
		activeFilter = &m.JobHistTab.Filter
		activeFilterOn = &m.JobHistTab.FilterOn
		activeUserInputs = &m.JobHistTab.UserInputs
		activeUserInputsOn = &m.JobHistTab.UserInputsOn
	case tabJobDetails:
		// here we're in the special situation, we need to pass on keys to viewport
		activeJDViewport = true
	case tabJobFromTemplate:
		activeTable = &m.JobFromTemplateTab.TemplatesTable
	case tabCluster:
		activeTab = &m.ClusterTab
		activeTable = &m.ClusterTab.SinfoTable
		activeFilter = &m.ClusterTab.Filter
		activeFilterOn = &m.ClusterTab.FilterOn
	}

	// Filter is turned on, take care of this first
	// TODO: revisit this for filtering on multiple tabs
	switch {
	case activeJDViewport:
		// catch only up/down keys, leave the rest to fallthrough
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keybindings.DefaultKeyMap.Up),
				key.Matches(msg, keybindings.DefaultKeyMap.Down),
				key.Matches(msg, keybindings.DefaultKeyMap.PageUp),
				key.Matches(msg, keybindings.DefaultKeyMap.PageDown):
				m.Log.Printf("VIEWPORT: up/down msg\n")
				var cmd tea.Cmd
				m.JobDetailsTab.ViewPort, cmd = m.JobDetailsTab.ViewPort.Update(msg)
				return m, cmd
			}
		}

	case activeFilterOn != nil && *activeFilterOn:
		m.Log.Printf("Filter is ON")
		switch msg := msg.(type) {

		case tea.KeyMsg:
			switch msg.Type {
			// TODO: when filter is set/cleared, trigger refresh with new filtered data
			case tea.KeyEnter:
				// finish & apply entering filter
				*activeFilterOn = false
				brk = true

			case tea.KeyEsc:
				// abort entering filter
				*activeFilterOn = false
				brk = true
			}

			if brk {
				// TODO: this is a "fix" for crashing-after-filter when Cursor() goes beyond list end
				// TODO: don't feel good about this... what if list is empty? no good. revisit
				// NOTE: This doesn't do what i image it should, cursor remains -1 when table is empty situation?
				// Explanation in clamp function: https://github.com/charmbracelet/bubbles/blob/13f52d678d315676568a656b5211b8a24a54a885/table/table.go#L296
				activeTable.SetCursor(0)
				activeTab.AdjTableHeight(m.winH, m.Log)
				//m.Log.Printf("ActiveTable = %v\n", activeTable)
				m.Log.Printf("Update: Filter set, setcursor(0), activetable.Cursor==%d\n", activeTable.Cursor())
				switch m.ActiveTab {
				case tabJobs:
					rows, sqf, err := m.JobTab.Squeue.FilterSqueueTable(m.JobTab.Filter.Value(), m.Log)
					if err != nil {
						m.Globals.ErrorHelp = err.ErrHelp
						m.Globals.ErrorMsg = err.OrigErr
						m.JobTab.Filter.SetValue("")
					} else {
						m.JobTab.SqueueTable.SetRows(*rows)
						m.JobTab.SqueueFiltered = *sqf
						m.JobTab.GetStatsFiltered(m.Log)
					}
					return m, nil

				case tabJobHist:
					rows, saf, err := m.JobHistTab.SacctHist.FilterSacctTable(m.JobHistTab.Filter.Value(), m.Log)
					if err != nil {
						m.Globals.ErrorHelp = err.ErrHelp
						m.Globals.ErrorMsg = err.OrigErr
						m.JobHistTab.Filter.SetValue("")
					} else {
						m.JobHistTab.SacctTable.SetRows(*rows)
						m.JobHistTab.SacctHistFiltered = *saf
						m.JobHistTab.GetStatsFiltered(m.Log)
					}
					return m, nil

				case tabCluster:
					m.ClusterTab.GetStatsFiltered(m.Log)
					return m, clustertab.QuickGetSinfo(m.Log)

				default:
					return m, nil
				}
			}
		}

		tmp, cmd := activeFilter.Update(msg)
		*activeFilter = tmp
		return m, cmd

	case activeUserInputsOn != nil && *activeUserInputsOn:
		m.Log.Printf("UserInputs is ON")
		switch msg := msg.(type) {

		case tea.KeyMsg:
			switch msg.Type {
			// TODO: when filter is set/cleared, trigger refresh with new filtered data
			case tea.KeyEnter:
				// finish & apply entering filter
				*activeUserInputsOn = false
				brk = true

			case tea.KeyEsc:
				// abort entering filter
				*activeUserInputsOn = false
				brk = true

			case tea.KeyUp, tea.KeyDown, tea.KeyTab:
				s := msg.String()
				m.JobHistTab.UserInputs.Params[m.JobHistTab.UserInputs.FocusIndex].Blur()

				if s == "up" {
					m.JobHistTab.UserInputs.FocusIndex--
				} else {
					m.JobHistTab.UserInputs.FocusIndex++
				}

				if m.JobHistTab.UserInputs.FocusIndex >= len(m.JobHistTab.UserInputs.Params) {
					m.JobHistTab.UserInputs.FocusIndex = 0
				} else if m.JobHistTab.UserInputs.FocusIndex < 0 {
					m.JobHistTab.UserInputs.FocusIndex = len(m.JobHistTab.UserInputs.Params)-1
				}

				m.JobHistTab.UserInputs.Params[m.JobHistTab.UserInputs.FocusIndex].Focus()
			}

			if brk {
				// TODO: this is a "fix" for crashing-after-filter when Cursor() goes beyond list end
				// TODO: don't feel good about this... what if list is empty? no good. revisit
				// NOTE: This doesn't do what i image it should, cursor remains -1 when table is empty situation?
				// Explanation in clamp function: https://github.com/charmbracelet/bubbles/blob/13f52d678d315676568a656b5211b8a24a54a885/table/table.go#L296
				activeTable.SetCursor(0)
				activeTab.AdjTableHeight(m.winH, m.Log)
				m.Log.Printf("Update: Param set, setcursor(0), activetable.Cursor==%d\n", activeTable.Cursor())
				switch m.ActiveTab {
				case tabJobHist:
					chngd := false
					t, _ := strconv.ParseUint(m.JobHistTab.UserInputs.Params[0].Value(), 10, 32)

					if m.JobHistTab.UserInputs.Params[1].Value() != m.JobHistTab.JobHistStart {
						m.JobHistTab.JobHistStart = m.JobHistTab.UserInputs.Params[1].Value()
						chngd = true
					}
					if m.JobHistTab.UserInputs.Params[2].Value() != m.JobHistTab.JobHistEnd {
						m.JobHistTab.JobHistEnd = m.JobHistTab.UserInputs.Params[2].Value()
						chngd = true
					}
					if uint(t) != m.JobHistTab.JobHistTimeout {
						m.JobHistTab.JobHistTimeout = uint(t)
						chngd = true
					}

					if chngd {
						m.Log.Println ("Refreshing JobHist View")
						m.JobHistTab.HistFetched = false
						return m, jobhisttab.GetSacctHist(strings.Join(m.Globals.UAccounts, ","),
										m.JobHistTab.JobHistStart,
										m.JobHistTab.JobHistEnd,
										m.JobHistTab.JobHistTimeout,
										m.Log)
					}

				default:
					return m, nil
				}
			}
		}

		var tmp textinput.Model
		var cmd tea.Cmd
		for i := range m.JobHistTab.UserInputs.Params {
			tmp, cmd = activeUserInputs.Params[i].Update(msg)
			*&activeUserInputs.Params[i] = tmp
		}
		return m, cmd

	case m.JobTab.MenuOn:
		m.Log.Printf("Update: In Menu\n")
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.JobTab.Menu.SetWidth(msg.Width)
			return m, nil

		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "esc":
				m.JobTab.MenuOn = false
				return m, nil
			case "ctrl+c":
				//m.quitting = true
				m.JobTab.MenuOn = false
				//return m, tea.Quit
				return m, nil

			case "enter":
				m.JobTab.MenuOn = false
				i, ok := m.JobTab.Menu.SelectedItem().(jobtab.MenuItem)
				if ok {
					m.JobTab.MenuChoice = jobtab.MenuItem(i)
					if m.JobTab.MenuChoice.GetAction() == "INFO" {
						// TODO: IF Stats==ON AND NxM, turn it of, can't have both on below NxM
						m.JobTab.InfoOn = true
						if m.JobTab.StatsOn && m.Globals.winH < 60 {
							m.Log.Printf("Toggle InfoBox: Height %d too low (<60). Turn OFF Stats\n", m.Globals.winH)
							// We have to turn off stats otherwise screen will break at this Height!
							m.JobTab.StatsOn = false
							// TODO: send a message via ErrMsg
						}
					}
					// host is needed for ssh command
					activeTab.AdjTableHeight(m.winH, m.Log)
					host := m.JobTab.SqueueFiltered.Jobs[m.JobTab.SqueueTable.Cursor()].BatchHost
					retCmd := m.JobTab.MenuChoice.ExecMenuItem(m.JobTab.SelectedJob, *host, m.Log)
					return m, retCmd
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
				jobfromtemplate.EditorKeyMap.DisableKeys()
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
				jobfromtemplate.EditorKeyMap.DisableKeys()
				tabKeys[m.ActiveTab].SetupKeys()
				name, err := jobfromtemplate.SaveToFile(m.JobFromTemplateTab.TemplatesTable.SelectedRow()[0], m.JobFromTemplateTab.TemplateEditor.Value(), m.Log)
				if err != nil {
					m.Log.Printf("ERROR saving to file!\n")
					return m, nil
				}
				return m, command.CallSbatch(name, m.Log)

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

	// ERROR msg
	case command.ErrorMsg:
		m.Log.Printf("ERROR msg, from: %s\n", msg.From)
		m.Log.Printf("ERROR msg, original error: %q\n", msg.OrigErr)
		m.Globals.ErrorMsg = msg.OrigErr
		m.Globals.ErrorHelp = msg.ErrHelp
		// cases when this is BAD and we can't continue
		switch msg.From {
		case "GetUserName", "GetUserAssoc":
			return m, tea.Quit
		}
		return m, nil

	// Ssh finished
	case command.SshCompleted:
		m.Log.Printf("Got SshCompleted msg, value: %#v\n", msg)
		return m, nil

	// UAccounts fetched
	case command.UserAssoc:
		m.Log.Printf("Got UserAssoc msg, value: %#v\n", msg)
		// TODO: consider changing this to string and do a join(",") to be ready to pass around
		m.Globals.UAccounts = append(m.Globals.UAccounts, msg...)
		m.Log.Printf("Appended UserAssoc msg go Globals, value now: %#v\n", m.Globals.UAccounts)
		// Now we trigger a sacctHist
		return m, jobhisttab.GetSacctHist(strings.Join(m.Globals.UAccounts, ","),
							m.JobHistTab.JobHistStart,
							m.JobHistTab.JobHistEnd,
							m.JobHistTab.JobHistTimeout,
							m.Log)

	// UserName fetched
	case command.UserName:
		m.Log.Printf("Got UserNAme msg, save %q to Globals.\n", msg)
		m.Globals.UserName = string(msg)
		// now, call GetUserAssoc()
		return m, command.GetUserAssoc(m.Globals.UserName, m.Log)

	// Shold executed
	case command.SBatchSent:
		m.Log.Printf("Got SBatchSent msg on file %q\n", msg.JobFile)
		return m, nil

	// Shold executed
	case command.SHoldSent:
		m.Log.Printf("Got SHoldSent msg on job %q\n", msg.Jobid)
		return m, jobtab.TimedGetSqueue(m.Log)

	// Scancel executed
	case command.ScancelSent:
		m.Log.Printf("Got ScancelSent msg on job %q\n", msg.Jobid)
		return m, jobtab.TimedGetSqueue(m.Log)

	// Srequeue executed
	case command.SRequeueSent:
		m.Log.Printf("Got SRequeueSent msg on job %q\n", msg.Jobid)
		return m, jobtab.TimedGetSqueue(m.Log)

	// Get initial job template list
	case jobfromtemplate.TemplatesListRows:
		m.Log.Printf("Update: Got TemplatesListRows msg: %#v\n", msg)
		if msg != nil {
			// if it's not empty, append to table
			m.JobFromTemplateTab.TemplatesTable.SetRows(msg)
			m.JobFromTemplateTab.TemplatesList = msg
		}
		return m, nil

	// getting initial template text
	case jobfromtemplate.TemplateText:
		m.Log.Printf("Update: Got TemplateText msg: %#v\n", msg)
		// HERE: we initialize the new textarea editor and flip the EditTemplate switch to ON
		tabKeys[m.ActiveTab].DisableKeys()
		jobfromtemplate.EditorKeyMap.SetupKeys()
		m.EditTemplate = true
		m.TemplateEditor = textarea.New()
		m.TemplateEditor.SetWidth(m.winW - 15)
		m.TemplateEditor.SetHeight(m.winH - 15)
		m.TemplateEditor.SetValue(string(msg))
		m.TemplateEditor.Focus()
		m.TemplateEditor.CharLimit = 0
		return m, jobfromtemplate.EditorOn()

	// Windows resize
	case tea.WindowSizeMsg:
		m.Log.Printf("Update: got WindowSizeMsg: %d %d\n", msg.Width, msg.Height)
		// TODO: if W<195 || H<60 we can't really run without breaking view, so quit and inform user
		// 187x44 == 13" MacBook Font 14 iTerm (HUGE letters!)
		if msg.Height < 43 || msg.Width < 185 {
			m.Log.Printf("FATAL: Window too small to run without breaking view. Have %dx%d. Need at least 185x43.\n", msg.Width, msg.Height)
			m.Globals.SizeErr = fmt.Sprintf("FATAL: Window too small to run without breaking view. Have %dx%d. Need at least 185x43.\nIncrease your terminal window and/or decrease font size.", msg.Width, msg.Height)
			return m, tea.Quit
		}
		m.winW = msg.Width
		m.winH = msg.Height
		// TODO: set also maxheight/width here on change?
		styles.MainWindow = styles.MainWindow.Height(m.winH - 10)
		styles.MainWindow = styles.MainWindow.Width(m.winW - 15)
		styles.HelpWindow = styles.HelpWindow.Width(m.winW - 15)
		styles.JobStepBoxStyle = styles.JobStepBoxStyle.Width(m.winW - 20)
		// InfoBox
		w := ((m.Globals.winW - 25) / 3) * 3
		styles.JobInfoInBox = styles.JobInfoInBox.Width(w / 3).Height(5)
		styles.JobInfoInBottomBox = styles.JobInfoInBottomBox.Width(w + 4).Height(5)

		// Adjust ALL tables
		m.JobTab.AdjTableHeight(m.winH, m.Log)
		m.JobHistTab.AdjTableHeight(m.winH, m.Log)
		m.ClusterTab.AdjTableHeight(m.winH, m.Log)

		// Fix jobdetails viewport
		m.JobDetailsTab.ViewPort.Width = m.winW - 15
		m.JobDetailsTab.ViewPort.Height = m.winH - 15

		// Adjust StatBoxes
		m.Log.Printf("CTB Width = %d\n", styles.ClusterTabStats.GetWidth())
		styles.ClusterTabStats = styles.ClusterTabStats.Width(m.winW - clustertab.SinfoTabWidth)
		m.Log.Printf("CTB Width = %d\n", styles.ClusterTabStats.GetWidth())

	// JobTab update
	case jobtab.SqueueJSON:
		m.Log.Printf("U(): got SqueueJSON\n")
		if len(msg.Jobs) != 0 {
			m.Squeue = msg

			// TODO:
			// fix: if after filtering m.table.Cursor|SelectedRow > lines in table, Info crashes trying to fetch nonexistent row
			rows, sqf, err := msg.FilterSqueueTable(m.JobTab.Filter.Value(), m.Log)
			if err != nil {
				m.Globals.ErrorHelp = err.ErrHelp
				m.Globals.ErrorMsg = err.OrigErr
				m.JobTab.Filter.SetValue("")
			} else {
				m.JobTab.SqueueTable.SetRows(*rows)
				m.JobTab.SqueueFiltered = *sqf
				m.JobTab.GetStatsFiltered(m.Log)
			}
		}
		m.UpdateCnt++
		// if active window != this, don't trigger new refresh
		if m.ActiveTab == tabJobs {
			return m, jobtab.TimedGetSqueue(m.Log)
		} else {
			return m, nil
		}

	// Cluster tab update
	case clustertab.SinfoJSON:
		m.Log.Printf("U(): got SinfoJSON\n")
		if len(msg.Nodes) != 0 {
			m.Sinfo = msg
			rows, sif, err := msg.FilterSinfoTable(m.ClusterTab.Filter.Value(), m.Log)
			if err != nil {
				m.Globals.ErrorHelp = err.ErrHelp
				m.Globals.ErrorMsg = err.OrigErr
				m.ClusterTab.Filter.SetValue("")
			} else {
				m.ClusterTab.SinfoTable.SetRows(*rows)
				m.ClusterTab.SinfoFiltered = *sif
				m.ClusterTab.GetStatsFiltered(m.Log)
			}
		}
		m.UpdateCnt++
		// if active window != this, don't trigger new refresh
		if m.ActiveTab == tabCluster {
			return m, clustertab.TimedGetSinfo(m.Log)
		} else {
			return m, nil
		}

	// Job Details tab update
	case slurm.SacctSingleJobHist:
		m.Log.Printf("Got SacctSingleJobHist\n")
		m.JobDetailsTab.SacctSingleJobHist = msg
		return m, nil

	// Job History tab update - NEW, with wrapped failure message
	case jobhisttab.JobHistTabMsg:
		m.Log.Printf("Got SacctJobHist len: %d\n", len(msg.Jobs))
		m.JobHistTab.SacctHist = msg.SacctJSON
		m.JobHistTab.HistFetchFail = msg.HistFetchFail
		// Filter and create filtered table
		rows, saf, err := msg.FilterSacctTable(m.JobHistTab.Filter.Value(), m.Log)
		if err != nil {
			m.Globals.ErrorHelp = err.ErrHelp
			m.Globals.ErrorMsg = err.OrigErr
			m.JobHistTab.Filter.SetValue("")
		} else {
			m.JobHistTab.SacctTable.SetRows(*rows)
			m.JobHistTab.SacctHistFiltered = *saf
			m.JobHistTab.GetStatsFiltered(m.Log)
		}
		if !m.JobHistTab.HistFetchFail {
			m.JobHistTab.HistFetched = true
		}

		return m, nil

	// Keys pressed
	case tea.KeyMsg:
		switch {

		// Counters
		case key.Matches(msg, keybindings.DefaultKeyMap.Count):
			// Depends at which tab we're at
			m.Log.Printf("Toggle Counters pressed at %d\n", m.ActiveTab)
			switch m.ActiveTab {
			case tabJobs:
				m.JobTab.InfoOn = false
				toggleSwitch(&m.JobTab.CountsOn)
			case tabJobHist:
				toggleSwitch(&m.JobHistTab.CountsOn)
			case tabCluster:
				toggleSwitch(&m.ClusterTab.CountsOn)
			}
			activeTab.AdjTableHeight(m.winH, m.Log)
			return m, nil

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
			tabKeys[m.ActiveTab].DisableKeys()
			m.ActiveTab = uint(k) - 1
			tabKeys[m.ActiveTab].SetupKeys()
			m.lastKey = msg.String()

			// clear error states
			m.Globals.ErrorHelp = ""
			m.Globals.ErrorMsg = nil

			switch m.ActiveTab {
			case tabJobs:
				return m, jobtab.TimedGetSqueue(m.Log)
			case tabCluster:
				return m, clustertab.TimedGetSinfo(m.Log)
			default:
				return m, nil
			}

		// TAB
		case key.Matches(msg, keybindings.DefaultKeyMap.Tab):
			tabKeys[m.ActiveTab].DisableKeys()
			// switch tab
			m.ActiveTab = (m.ActiveTab + 1) % uint(len(tabs))
			// setup keys
			tabKeys[m.ActiveTab].SetupKeys()
			m.lastKey = "tab"

			// clear error states
			m.Globals.ErrorHelp = ""
			m.Globals.ErrorMsg = nil

			switch m.ActiveTab {
			case tabJobs:
				return m, jobtab.TimedGetSqueue(m.Log)
			case tabCluster:
				return m, clustertab.TimedGetSinfo(m.Log)
			default:
				return m, nil
			}

		// Shift+TAB
		case key.Matches(msg, keybindings.DefaultKeyMap.ShiftTab):
			tabKeys[m.ActiveTab].DisableKeys()
			// switch tab
			if m.ActiveTab == 0 {
				m.ActiveTab = uint(len(tabs) - 1)
			} else {
				m.ActiveTab -= 1
			}
			// setup keys
			tabKeys[m.ActiveTab].SetupKeys()
			m.lastKey = "tab"

			// clear error states
			m.Globals.ErrorHelp = ""
			m.Globals.ErrorMsg = nil

			switch m.ActiveTab {
			case tabJobs:
				return m, jobtab.TimedGetSqueue(m.Log)
			case tabCluster:
				return m, clustertab.TimedGetSinfo(m.Log)
			default:
				return m, nil
			}

		// SLASH
		case key.Matches(msg, keybindings.DefaultKeyMap.Slash):
			m.Log.Printf("Filter key pressed\n")
			switch m.ActiveTab {
			case tabJobs:
				m.JobTab.FilterOn = true
			case tabJobHist:
				m.JobHistTab.FilterOn = true
			case tabCluster:
				m.ClusterTab.FilterOn = true
			}
			activeTab.AdjTableHeight(m.winH, m.Log)
			return m, nil

		// t
		case key.Matches(msg, keybindings.DefaultKeyMap.TimeRange):
			m.Log.Printf("time-range key pressed\n")
			switch m.ActiveTab {
			case tabJobHist:
				m.JobHistTab.UserInputsOn = true
			}
			activeTab.AdjTableHeight(m.winH, m.Log)
			return m, nil

		// ENTER
		case key.Matches(msg, keybindings.DefaultKeyMap.Enter):
			switch m.ActiveTab {

			// Job Queue tab: Open Job menu
			case tabJobs:
				// Check if there is anything in the filtered table and if cursor is on a valid item
				n := m.JobTab.SqueueTable.Cursor()
				m.Log.Printf("Update ENTER key @ jobqueue table\n")
				if n == -1 || len(m.JobTab.SqueueFiltered.Jobs) == 0 {
					m.Log.Printf("Update ENTER key @ jobqueue table, no jobs selected/empty table\n")
					return m, nil
				}
				// IF Info==ON AND NxM, turn it of, can't have both on below NxM
				if m.JobTab.InfoOn && m.Globals.winH < 60 {
					m.Log.Printf("Toggle MenuBox: Height %d too low (<60). Turn OFF Info\n", m.Globals.winH)
					m.JobTab.InfoOn = false
				}
				// If yes, turn on menu
				m.JobTab.MenuOn = true
				m.JobTab.SelectedJob = m.JobTab.SqueueTable.SelectedRow()[0]
				m.JobTab.SelectedJobState = m.JobTab.SqueueTable.SelectedRow()[4]
				// Create new menu
				m.JobTab.Menu = jobtab.NewMenu(m.JobTab.SelectedJobState, m.Log)
				return m, nil

			// Job History tab: Select Job from history and open its Details tab
			case tabJobHist:
				n := m.JobHistTab.SacctTable.Cursor()
				m.Log.Printf("Update ENTER key @ jobhist table, cursor=%d, len=%d\n", n, len(m.JobHistTab.SacctHistFiltered.Jobs))
				if n == -1 || len(m.JobHistTab.SacctHistFiltered.Jobs) == 0 {
					m.Log.Printf("Update ENTER key @ jobhist table, no jobs selected/empty table\n")
					return m, nil
				}
				tabKeys[m.ActiveTab].DisableKeys()
				m.ActiveTab = tabJobDetails
				tabKeys[m.ActiveTab].SetupKeys()
				m.JobDetailsTab.SelJobIDNew = n
				// clear error states
				m.Globals.ErrorHelp = ""
				m.Globals.ErrorMsg = nil

				// new job selected, fill out viewport
				m.JobDetailsTab.SetViewportContent(&m.JobHistTab, m.Log)
				return m, nil

			// Job from Template tab: Open template for editing
			case tabJobFromTemplate:
				m.Log.Printf("Update ENTER key @ jobfromtemplate table\n")
				// return & handle editing there
				if len(m.JobFromTemplateTab.TemplatesList) != 0 {
					return m, jobfromtemplate.GetTemplate(m.JobFromTemplateTab.TemplatesTable.SelectedRow()[2], m.Log)
				} else {
					return m, nil
				}
			}

		// Refresh the View
		case key.Matches(msg, keybindings.DefaultKeyMap.Refresh):
			switch m.ActiveTab {
			case tabJobHist:
				m.Log.Println ("Refreshing JobHist View")
				m.JobHistTab.HistFetched = false
				return m, jobhisttab.GetSacctHist(strings.Join(m.Globals.UAccounts, ","),
				                                  m.JobHistTab.JobHistStart,
								  m.JobHistTab.JobHistEnd,
								  m.JobHistTab.JobHistTimeout,
								  m.Log)
			}
			return m, nil

		// Info - toggle on/off
		case key.Matches(msg, keybindings.DefaultKeyMap.Info):
			m.Log.Println("Toggle InfoBox")

			// TODO: IF Stats==ON AND NxM, turn it of, can't have both on below NxM
			if m.JobTab.StatsOn && m.Globals.winH < 60 {
				m.Log.Printf("Toggle InfoBox: Height %d too low (<60). Turn OFF Stats\n", m.Globals.winH)
				// We have to turn off stats otherwise screen will break at this Height!
				m.JobTab.StatsOn = false
				// TODO: send a message via ErrMsg
			}

			m.JobTab.CountsOn = false
			toggleSwitch(&m.JobTab.InfoOn)
			m.JobTab.AdjTableHeight(m.Globals.winH, m.Log)
			return m, nil

		// Stats - toggle on/off
		case key.Matches(msg, keybindings.DefaultKeyMap.Stats):
			switch m.ActiveTab {
			case tabJobs:
				m.Log.Printf("JobTab toggle from: %v\n", m.JobTab.StatsOn)
				toggleSwitch(&m.JobTab.StatsOn)
				// IF Info==ON AND NxM, turn it of, can't have both on below NxM
				if m.JobTab.InfoOn && m.Globals.winH < 60 {
					m.Log.Printf("Toggle StatsBox: Height %d too low (<60). Turn OFF Info\n", m.Globals.winH)
					m.JobTab.InfoOn = false
				}
				m.Log.Printf("JobTab toggle to: %v\n", m.JobTab.StatsOn)
			case tabJobHist:
				m.Log.Printf("JobHistTab toggle from: %v\n", m.JobHistTab.StatsOn)
				toggleSwitch(&m.JobHistTab.StatsOn)
				m.Log.Printf("JobHistTab toggle to: %v\n", m.JobHistTab.StatsOn)
			case tabCluster:
				m.Log.Printf("JobCluster toggle from: %v\n", m.ClusterTab.StatsOn)
				toggleSwitch(&m.ClusterTab.StatsOn)
				m.Log.Printf("JobCluster toggle to: %v\n", m.ClusterTab.StatsOn)
			}
			return m, nil

		// QUIT
		case key.Matches(msg, keybindings.DefaultKeyMap.Quit):
			m.Log.Printf("Quit key pressed\n")
			return m, tea.Quit
		}
	}

	return m, nil
}

func toggleSwitch(b *bool) {
	if *b {
		*b = false
	} else {
		*b = true
	}
}
