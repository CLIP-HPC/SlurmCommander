package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pja237/slurmcommander-dev/internal/command"
	"github.com/pja237/slurmcommander-dev/internal/keybindings"
	"github.com/pja237/slurmcommander-dev/internal/model/tabs/clustertab"
	"github.com/pja237/slurmcommander-dev/internal/model/tabs/jobfromtemplate"
	"github.com/pja237/slurmcommander-dev/internal/model/tabs/jobhisttab"
	"github.com/pja237/slurmcommander-dev/internal/model/tabs/jobtab"
	"github.com/pja237/slurmcommander-dev/internal/slurm"
	"github.com/pja237/slurmcommander-dev/internal/styles"
	"github.com/pja237/slurmcommander-dev/internal/table"
)

type errMsg error

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var (
		brk            bool = false
		activeTable    *table.Model
		activeFilter   *textinput.Model
		activeFilterOn *bool
	)

	// This shortens the testing for table movement keys
	switch m.ActiveTab {
	case tabJobs:
		activeTable = &m.JobTab.SqueueTable
		activeFilter = &m.JobTab.Filter
		activeFilterOn = &m.JobTab.FilterOn
	case tabJobHist:
		activeTable = &m.JobHistTab.SacctTable
		activeFilter = &m.JobHistTab.Filter
		activeFilterOn = &m.JobHistTab.FilterOn
	case tabJobFromTemplate:
		activeTable = &m.JobFromTemplateTab.TemplatesTable
	case tabCluster:
		activeTable = &m.JobClusterTab.SinfoTable
		activeFilter = &m.JobClusterTab.Filter
		activeFilterOn = &m.JobClusterTab.FilterOn
	}

	// Filter is turned on, take care of this first
	// TODO: revisit this for filtering on multiple tabs
	switch {
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
				activeFilter.SetValue("")
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
					m.JobClusterTab.GetStatsFiltered(m.Log)
					return m, clustertab.QuickGetSinfo(m.Log)
				default:
					return m, nil
				}
			}
		}

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
				// TODO: This is just temporarily here, instead of this, depending on the MenuChoice turn on Info if selected
				m.JobTab.InfoOn = true
				i, ok := m.JobTab.Menu.SelectedItem().(jobtab.MenuItem)
				if ok {
					m.JobTab.MenuChoice = jobtab.MenuItem(i)
					// host is needed for ssh command
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
		//return m, nil
		m.Log.Printf("Appended UserAssoc msg go Globals, calling GetSacctHist()\n")
		return m, jobhisttab.GetSacctHist(strings.Join(m.Globals.UAccounts, ","), m.Globals.JobHistStart, m.Globals.JobHistTimeout, m.Log)

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
		m.Log.Printf("Update: got WindowSizeMsg: %d %d\n", msg.Width, msg.Height)
		// TODO: if W<195 || H<60 we can't really run without breaking view, so quit and inform user
		if msg.Height < 60 || msg.Width < 195 {
			m.Log.Printf("FATAL: Window too small to run without breaking view. Have %dx%d. Need at least 195x60.\n", msg.Width, msg.Height)
			//m.Globals.SizeErr = fmt.Sprintf("FATAL: Window too small to run without breaking view. Have %dx%d. Need at least 195x60.\nIncrease your terminal window and/or decrease font size.", msg.Width, msg.Height)
			//return m, tea.Quit
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
			rows, sif, err := msg.FilterSinfoTable(m.JobClusterTab.Filter.Value(), m.Log)
			if err != nil {
				m.Globals.ErrorHelp = err.ErrHelp
				m.Globals.ErrorMsg = err.OrigErr
				m.JobClusterTab.Filter.SetValue("")
			} else {
				m.JobClusterTab.SinfoTable.SetRows(*rows)
				m.JobClusterTab.SinfoFiltered = *sif
				m.JobClusterTab.GetStatsFiltered(m.Log)
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
		// TODO: Here we don't tick refresh, because of potentially long sacct calls, make it manually triggered
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
				toggleSwitch(&m.JobClusterTab.CountsOn)
			}
			return m, nil

		// UP
		// TODO: what if it's a list?
		case key.Matches(msg, keybindings.DefaultKeyMap.Up):
			activeTable.MoveUp(1)
			m.lastKey = "up"

		// DOWN
		case key.Matches(msg, keybindings.DefaultKeyMap.Down):
			t := time.Now()
			m.Log.Printf("Update: Move down\n")
			activeTable.MoveDown(1)
			m.Log.Printf("Update: Move down finished in: %.3f sec\n", time.Since(t).Seconds())
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

		// SLASH
		case key.Matches(msg, keybindings.DefaultKeyMap.Slash):
			switch {
			case m.ActiveTab == tabJobs:
				m.JobTab.FilterOn = true
			case m.ActiveTab == tabJobHist:
				m.JobHistTab.FilterOn = true
			case m.ActiveTab == tabCluster:
				m.JobClusterTab.FilterOn = true
			}
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
				m.ActiveTab = tabJobDetails
				tabKeys[m.ActiveTab].SetupKeys()
				// TODO: this we change to directly address data from SacctHistFiltered instead of
				// calling another sacct Cmd
				//m.JobDetailsTab.SelJobID = m.JobHistTab.SacctTable.SelectedRow()[0]
				//return m, command.SingleJobGetSacct(m.JobDetailsTab.SelJobID, m.Globals.JobHistStart, m.Log)
				m.JobDetailsTab.SelJobIDNew = n
				// clear error states
				m.Globals.ErrorHelp = ""
				m.Globals.ErrorMsg = nil
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

		// Info - toggle on/off
		case key.Matches(msg, keybindings.DefaultKeyMap.Info):
			m.Log.Printf("Toggle InfoBox\n")
			m.JobTab.CountsOn = false
			toggleSwitch(&m.JobTab.InfoOn)
			return m, nil

		// Stats - toggle on/off
		case key.Matches(msg, keybindings.DefaultKeyMap.Stats):
			switch m.ActiveTab {
			case tabJobs:
				m.Log.Printf("JobTab toggle from: %v\n", m.JobTab.StatsOn)
				toggleSwitch(&m.JobTab.StatsOn)
				m.Log.Printf("JobTab toggle to: %v\n", m.JobTab.StatsOn)
			case tabJobHist:
				m.Log.Printf("JobHistTab toggle from: %v\n", m.JobHistTab.StatsOn)
				toggleSwitch(&m.JobHistTab.StatsOn)
				m.Log.Printf("JobHistTab toggle to: %v\n", m.JobHistTab.StatsOn)
			case tabCluster:
				m.Log.Printf("JobCluster toggle from: %v\n", m.JobClusterTab.StatsOn)
				toggleSwitch(&m.JobClusterTab.StatsOn)
				m.Log.Printf("JobCluster toggle to: %v\n", m.JobClusterTab.StatsOn)
			}
			return m, nil

		// QUIT
		case key.Matches(msg, keybindings.DefaultKeyMap.Quit):
			fmt.Println("Quit key pressed")
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
