package frontend

import (
	"fmt"
	"log"
	"os"

	"github.com/bashbunni/project-management/models"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TODO: render list of entries -> need to update m.list with entries instead of projects and change title to proj name
// TODO: stop cursor from moving in edit mode -> override list.Model

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

// implements tea.Model (Init, Update, View)
type model struct {
	list list.Model 
	input textinput.Model
	pr *models.GormProjectRepository
	er *models.GormEntryRepository
	keymap keymap
	mode string
	err error
}

type keymap struct {
	create key.Binding
	enter key.Binding
	rename key.Binding
	delete key.Binding
	back key.Binding
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case updateEntryListMsg:
		entries := make([]list.Item, 0, len(msg.entries))
		for _, e := range msg.entries {
			entries = append(entries, item{
				title: fmt.Sprintf("%d", e.ID),
			})
		}
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
					items := m.list.Items()
					activeItem := items[m.list.Index()]
					cmds = append(cmds, renameProjectCmd(activeItem.(models.Project).ID, m.pr, m.input.Value()))
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
					// TODO: update list
					return m, nil
				case key.Matches(msg, m.keymap.rename):
					m.mode = "edit"
					m.input.Focus()
					cmds = append(cmds, textinput.Blink)
				case key.Matches(msg, m.keymap.delete):
					items := m.list.Items()
					activeItem := items[m.list.Index()]
					cmds = append(cmds, deleteProjectCmd(activeItem.(models.Project).ID, m.pr))
			}
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, cmd)
		}

	case tea.WindowSizeMsg:
		top, right, bottom, left := docStyle.GetMargin()
		m.list.SetSize(msg.Width-left-right, msg.Height-top-bottom-1)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.input.Focused() {
		return docStyle.Render(m.list.View() + "\n" + m.input.View())
	}
	return docStyle.Render(m.list.View() + "\n") 
}

// functions

// Initial model (AKA first View)
func ChooseProject(pr models.GormProjectRepository, er models.GormEntryRepository) {
	input := textinput.New()
	input.Prompt = "$ "
	input.Placeholder = "Project name..."
	input.CharLimit = 250
	input.Width = 50

	
	projects, err := pr.GetAllProjects()
	if err != nil {
		log.Fatal(err)
	}
	items := projectsToItems(projects)
	m := initModel(items, input, &pr, &er)
	
	tea.LogToFile("debug.log", "debug")
	p := tea.NewProgram(m)
	p.EnterAltScreen()

	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func initModel(items []list.Item, input textinput.Model, pr *models.GormProjectRepository, er *models.GormEntryRepository) tea.Model {
	m := model{list: list.NewModel(items, list.NewDefaultDelegate(), 0, 0), input: input, pr: pr, er: er, keymap: 
	keymap{
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
	return m 
}

// TODO: use generics
// convert []model.Project to []list.Item
func projectsToItems(projects []models.Project) []list.Item {
	items := make([]list.Item, len(projects))
	for i, proj := range projects {
		items[i] = list.Item(proj)
	}
	return items
}
