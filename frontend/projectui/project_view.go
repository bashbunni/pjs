package projectui

import (
	"log"

	"github.com/bashbunni/project-management/entry"
	"github.com/bashbunni/project-management/frontend/constants"
	"github.com/bashbunni/project-management/project"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// TODO: add multi-page navigation

type SelectMsg struct {
	ActiveProjectID uint
}

type Model struct {
	mode  string
	list  list.Model
	input textinput.Model
	pr    *project.GormRepository
}

func New(pr *project.GormRepository, er *entry.GormRepository) tea.Model {
	input := textinput.New()
	input.Prompt = "$ "
	input.Placeholder = "Project name..."
	input.CharLimit = 250
	input.Width = 50

	items := newProjectList(pr)
	m := Model{mode: "", list: list.NewModel(items, list.NewDefaultDelegate(), 0, 0), input: input, pr: pr}
	m.list.Title = "projects"
	m.list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			constants.Keymap.Create,
			constants.Keymap.Rename,
			constants.Keymap.Delete,
			constants.Keymap.Back,
		}
	}
	return m
}

func newProjectList(pr *project.GormRepository) []list.Item {
	projects, err := pr.GetAllProjects()
	if err != nil {
		log.Fatal(err)
	}
	return projectsToItems(projects)

}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		top, right, bottom, left := constants.DocStyle.GetMargin()
		m.list.SetSize(msg.Width-left-right, msg.Height-top-bottom-1)
	case updateProjectListMsg:
		projects, err := m.pr.GetAllProjects()
		if err != nil {
			log.Fatal(err)
		}
		items := projectsToItems(projects)
		m.list.SetItems(items)
		m.mode = ""
	case renameProjectMsg:
		m.list.SetItems(msg)
		m.mode = ""
	case tea.KeyMsg:
		if m.input.Focused() {
			if key.Matches(msg, constants.Keymap.Enter) {
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
			if key.Matches(msg, constants.Keymap.Back) {
				m.input.SetValue("")
				m.mode = ""
				m.input.Blur()
			}
			// only log keypresses for the input field when it's focused
			m.input, cmd = m.input.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			switch {
			case key.Matches(msg, constants.Keymap.Create):
				m.mode = "create"
				m.input.Focus()
				cmd = textinput.Blink
			case msg.String() == "ctrl+c":
				return m, tea.Quit
			case key.Matches(msg, constants.Keymap.Enter):
				cmd = selectProjectCmd(m.getActiveProjectID())
			case key.Matches(msg, constants.Keymap.Rename):
				m.mode = "edit"
				m.input.Focus()
				cmd = textinput.Blink
			case key.Matches(msg, constants.Keymap.Delete):
				cmd = deleteProjectCmd(m.getActiveProjectID(), m.pr)
			default:
				m.list, cmd = m.list.Update(msg)
			}
			cmds = append(cmds, cmd)
		}
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.input.Focused() {
		return constants.DocStyle.Render(m.list.View() + "\n" + m.input.View())
	}
	return constants.DocStyle.Render(m.list.View() + "\n")
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

func (m Model) getActiveProjectID() uint {
	items := m.list.Items()
	activeItem := items[m.list.Index()]
	return activeItem.(project.Project).ID
}
