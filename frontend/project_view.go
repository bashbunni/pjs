package frontend

import (
	"log"

	"github.com/bashbunni/project-management/entry"
	"github.com/bashbunni/project-management/project"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// TODO: add multi-page navigation

func (m model) handleProjectList(msg tea.Msg, cmds []tea.Cmd, cmd tea.Cmd) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case updateProjectListMsg:
		projects, err := m.pr.GetAllProjects()
		if err != nil {
			log.Fatal(err)
		}
		items := projectsToItems(projects)
		m.list.SetItems(items)
		m.mode = ""
	case renameProjectMsg:
		projects, err := m.pr.GetAllProjects()
		if err != nil {
			log.Fatal(err)
		}
		items := projectsToItems(projects)
		m.list.SetItems(items)
		m.mode = ""
	case tea.KeyMsg:
		if m.input.Focused() {
			if key.Matches(msg, m.keymap.enter) {
				if m.mode == "create" {
					cmds = append(cmds, createProjectCmd(m.input.Value(), m.pr))
				}
				if m.mode == "edit" {
					cmds = append(cmds, renameProjectCmd(m.getActiveProjectID(), m.pr, m.input.Value()))
				}
				m.input.SetValue("")
				m.mode = ""
				m.input.Blur()
			}
			if key.Matches(msg, m.keymap.back) {
				m.input.SetValue("")
				m.mode = ""
				m.input.Blur()
			}
			// only log keypresses for the input field when it's focused
			m.input, cmd = m.input.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			switch {
			case key.Matches(msg, m.keymap.create):
				m.mode = "create"
				m.input.Focus()
				cmds = append(cmds, textinput.Blink)
			case msg.String() == "ctrl+c":
				return m, tea.Quit
			case key.Matches(msg, m.keymap.enter):
				m.initEntries()
				m.state = "viewEntries"
			case key.Matches(msg, m.keymap.rename):
				m.mode = "edit"
				m.input.Focus()
				cmds = append(cmds, textinput.Blink)
			case key.Matches(msg, m.keymap.delete):
				items := m.list.Items()
				activeItem := items[m.list.Index()]
				cmds = append(cmds, deleteProjectCmd(activeItem.(project.Project).ID, m.pr))
			}
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, cmd)
		}
	}
	return m, tea.Batch(cmds...)
}

func initProjectView(items []list.Item, input textinput.Model, pr *project.GormRepository, er *entry.GormRepository) tea.Model {
	m := model{state: "viewProjectList", list: list.NewModel(items, list.NewDefaultDelegate(), 0, 0), input: input, pr: pr, er: er, keymap: keymap{
		create: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "create"),
		),
		enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		rename: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "rename"),
		),
		delete: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete"),
		),
		back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
	},
	}
	m.list.Title = "projects"
	m.list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			m.keymap.create,
			m.keymap.rename,
			m.keymap.delete,
			m.keymap.back,
		}
	}

	m.viewport = viewport.New(8, 8)
	return m
}

// TODO: use generics
// convert []model.Project to []list.Item
func projectsToItems(projects []project.Project) []list.Item {
	items := make([]list.Item, len(projects))
	for i, proj := range projects {
		items[i] = list.Item(proj)
	}
	return items
}

func (m model) getActiveProjectID() uint {
	items := m.list.Items()
	activeItem := items[m.list.Index()]
	return activeItem.(project.Project).ID
}
