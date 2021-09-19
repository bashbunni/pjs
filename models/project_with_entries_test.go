package models

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCreateProjectWithEntries(t *testing.T) {
	tc := []struct {
		name string
		got  ProjectWithEntries
		want ProjectWithEntries
	}{
		{
			"test1: happy path",
			*CreateProjectWithEntries(
				Project{ID: 1, Name: "Ikiris"},
				MockEntryRepository{map[uint]*Entry{
					1: {ProjectID: 1, Message: "ikiris is cool"},
					2: {ProjectID: 1, Message: ""},
					3: {ProjectID: 1, Message: "ikiris has a diamond"},
					4: {ProjectID: 1, Message: "*FEJS-()' I like special characters"},
				}}),
			ProjectWithEntries{
				Project{ID: 1, Name: "Ikiris"},
				[]Entry{
					{ProjectID: 1, Message: "ikiris is cool"},
					{ProjectID: 1, Message: ""},
					{ProjectID: 1, Message: "ikiris has a diamond"},
					{ProjectID: 1, Message: "*FEJS-()' I like special characters"},
				}},
		},
		{
			"test2: not happy path",
			*CreateProjectWithEntries(
				Project{ID: 1, Name: "Ikiris"},
				MockEntryRepository{map[uint]*Entry{
					1: {ProjectID: 1, Message: "ikiris is cool"},
					2: {ProjectID: 3, Message: ""},
					3: {ProjectID: 1, Message: "ikiris has a diamond"},
					4: {ProjectID: 1, Message: "*FEJS-()' I like special characters"},
				}}),
			ProjectWithEntries{
				Project{ID: 1, Name: "Ikiris"},
				[]Entry{
					{ProjectID: 1, Message: "ikiris is cool"},
					{ProjectID: 1, Message: "ikiris has a diamond"},
					{ProjectID: 1, Message: "*FEJS-()' I like special characters"},
				}},
		},
		{
			"test3: negative projectID",
			*CreateProjectWithEntries(
				Project{ID: 1000, Name: "Ikiris"},
				MockEntryRepository{map[uint]*Entry{
					1: {ProjectID: 250, Message: "ikiris is cool"},
					2: {ProjectID: 1000032, Message: ""},
					3: {ProjectID: 1000000000000000, Message: "ikiris has a diamond"},
					4: {ProjectID: 1, Message: "*FEJS-()' I like special characters"},
				}}),
			ProjectWithEntries{
				Project{ID: 1000, Name: "Ikiris"},
				nil},
		},
	}
	for _, c := range tc {
		if diff := cmp.Diff(c.want, c.got); diff != "" {
			t.Errorf("CreateProjectWithEntries() %s:\n%v", c.name, diff)
		}
	}
}
