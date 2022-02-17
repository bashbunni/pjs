package frontend

import (
	"fmt"
	"log"
	"os"

	"github.com/bashbunni/project-management/models"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

// TODO: render list of entries -> need to update m.list with entries instead of projects and change title to proj name
// TODO: stop cursor from moving in edit mode -> override list.Model

const content = `
Today’s Menu
| Name        | Price | Notes                           |
| ---         | ---   | ---                             |
| Tsukemono   | $2    | Just an appetizer               |
| Tomato Soup | $4    | Made with San Marzano tomatoes  |
| Okonomiyaki | $4    | Takes a few minutes to make     |
| Curry       | $3    | We can add squash if you’d like |
Seasonal Dishes
| Name                 | Price | Notes              |
| ---                  | ---   | ---                |
| Steamed bitter melon | $2    | Not so bitter      |
| Takoyaki             | $3    | Fun to eat         |
| Winter squash        | $3    | Today it's pumpkin |
Desserts
| Name         | Price | Notes                 |
| ---          | ---   | ---                   |
| Dorayaki     | $4    | Looks good on rabbits |
| Banana Split | $5    | A classic             |
| Cream Puff   | $3    | Pretty creamy!        |
All our dishes are made in-house by Karen, our chef. Most of our ingredients
are from our garden or the fish market down the street.
Some famous people that have eaten here lately:
* [x] René Redzepi
* [x] David Chang
* [ ] Jiro Ono (maybe some day)
Bon appétit!
`

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

// implements tea.Model (Init, Update, View)
type model struct {
	state    string
	viewport viewport.Model
	list     list.Model
	input    textinput.Model
	pr       *models.GormProjectRepository
	er       *models.GormEntryRepository
	keymap   keymap
	mode     string
	err      error // TODO: does this get used
}

type keymap struct {
	create key.Binding
	enter  key.Binding
	rename key.Binding
	delete key.Binding
	back   key.Binding
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch m.state {
	case "viewProjectList":
		return m.handleProjectList(msg, cmds, cmd)
	case "viewEntries":
		return m.handleEntriesList(msg, cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	switch m.state {
	case "viewEntries":
		return docStyle.Render(m.viewport.View())
	default:
		if m.input.Focused() {
			return docStyle.Render(m.list.View() + "\n" + m.input.View())
		}
		return docStyle.Render(m.list.View() + "\n")
	}
}
func (m model) handleEntriesList(msg tea.Msg, cmds []tea.Cmd, cmd tea.Cmd) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		// check for my cmds being done running
		// check for keypresses?
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		default:
			m.viewport, cmd = m.viewport.Update(msg)
			cmds = append(cmds, cmd)
		}
	}
	return m, tea.Batch(cmds...)
}

func initEntries(m tea.Model) (viewport.Model, error) {
	vp := viewport.New(78, 20)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)
	renderer, err := glamour.NewTermRenderer(glamour.WithStylePath("notty"))
	if err != nil {
		return vp, err
	}
	str, err := renderer.Render(content)
	if err != nil {
		return vp, err
	}
	vp.SetContent(str)
	return vp, nil
}

func StartTea(pr models.GormProjectRepository, er models.GormEntryRepository) {
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
	m := initProjectView(items, input, &pr, &er)
	tea.LogToFile("debug.log", "debug")
	p := tea.NewProgram(m)
	p.EnterAltScreen()
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
