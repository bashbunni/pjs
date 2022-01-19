package tests

import (
	"testing"

	"github.com/bashbunni/project-management/mocks"
	"github.com/bashbunni/project-management/models"
)

// func NewMock() {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		log.Fatalf("unable to open stub database: %v", err)
// 	}
// }

// TODO: figure out Setup() if it would be useful

func TestGetProjectByID(t *testing.T) {
	var tests = []struct {
		testname string
		repository models.ProjectRepository
		id uint
		want models.Project
	}{
		{"doesn't exist: too high", mocks.MockProjectRepository{
			Projects: map[uint]*models.Project{
			1: &models.Project{Name: "project1"}, 
			2: &models.Project{Name: "project2"},
		},}, 3, models.Project{}},
		{"does exist", mocks.MockProjectRepository{
			Projects: map[uint]*models.Project{
			1: &models.Project{Name: "project1"}, 
			2: &models.Project{Name: "project2"},
		},}, 2, models.Project{Name: "project2"}},
		{"doesn't exist: too low", mocks.MockProjectRepository{
			Projects: map[uint]*models.Project{
			1: &models.Project{Name: "project1"}, 
			2: &models.Project{Name: "project2"},
		},}, 0, models.Project{}},
	}
		for _, tt := range tests {
			t.Run(tt.testname, func(t *testing.T) {
				got := tt.repository.GetProjectByID(tt.id)
				if got != tt.want {
					t.Errorf("got %s want %s", got.Name, tt.want.Name)
				}
			})
		}
}

func TestHasProjects(t *testing.T) {
	var tests = []struct {
		testname string
		repository models.ProjectRepository
		want bool
	}{
		{"has projects", mocks.MockProjectRepository{Projects: map[uint]*models.Project{
		1: &models.Project{Name: "project1"}, 
		2: &models.Project{Name: "project2"},
	}}, true},
		{"no projects", mocks.MockProjectRepository{Projects: map[uint]*models.Project{}}, false},
	}

	for _, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			got := tt.repository.HasProjects()

			if got != tt.want {
				t.Errorf("got %t want %t", got, tt.want)
			}	
		})
	}
}

