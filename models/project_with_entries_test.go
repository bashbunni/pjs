package models

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCreateProjectWithEntries(t *testing.T) {
	got := *CreateProjectWithEntries(
		Project{ID: 1, Name: "Ikiris"},
		MockEntryRepository{map[uint]*Entry{
			1: {ProjectID: 1, Message: "ikiris is cool"},
			2: {ProjectID: 1, Message: ""},
			3: {ProjectID: 1, Message: "ikiris has a diamond"},
			4: {ProjectID: 1, Message: "*FEJS-()' I like special characters"},
		}})
	want := ProjectWithEntries{
		Project{ID: 1, Name: "Ikiris"},
		[]Entry{
			{ProjectID: 1, Message: "ikiris is cool"},
			{ProjectID: 1, Message: ""},
			{ProjectID: 1, Message: "ikiris has a diamond"},
			{ProjectID: 1, Message: "*FEJS-()' I like special characters"},
		},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("CreateProjectWithEntries() mismatch:\n%v", diff)
	}
}
