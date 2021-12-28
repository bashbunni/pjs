package frontend

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// TODO: add lipgloss for prettier frontend
// TODO: have tabs for different projects if multiple are selected (see lipgloss docs)

type projectModel struct {
	projects []string // list of projects
	cursor   int
	selected map[int]struct{} // keep track of which projects are selected
}

func newProjectModel() projectModel {
	return projectModel{
		projects: []string{"project 1", "project 2"}, // TODO: actually get projects
		selected: make(map[int]struct{}),
	}
}

// TODO: look into this
func (p projectModel) Init() tea.Cmd {
	return nil
}

func (p projectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return p, tea.Quit
		case "up", "k":
			if p.cursor > 0 {
				p.cursor--
			}
		case "down", "j":
			if p.cursor < len(p.projects)-1 {
				p.cursor++
			}
		case "enter", " ":
			_, ok := p.selected[p.cursor]
			if ok {
				delete(p.selected, p.cursor)
			} else {
				p.selected[p.cursor] = struct{}{} // TODO: why is this a struct?
			}
		}
	}
	return p, nil
}

func (p projectModel) View() string {
	content := "Which project would you like to choose?\n\n"
	for i, choice := range p.projects {
		cursor := " "
		if p.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := p.selected[i]; ok {
			checked = "x"
		}
		content += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}
	content += "\n(esc) to quit"
	return content
}

func main() {
	ui := tea.NewProgram(newProjectModel())
	if err := ui.Start(); err != nil {
		fmt.Printf("unable to create UI: %v", err)
		os.Exit(1)
	}
}
