package tests

import (
	"fmt"
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

func TestHasProjects(t *testing.T) {
	var tests = []struct {
		testname string
		data models.ProjectRepository
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
			got := tt.data.HasProjects()

			if got != tt.want {
				t.Errorf("got %t want %t", got, tt.want)
			}	
		})
	}
}

