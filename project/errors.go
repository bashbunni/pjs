package project

/*
need to implement func Error() string for it to implement Error interface
*/

const (
	emptyTable          = "empty table"
	cannotCreateProject = "cannot create project"
	cannotDeleteProject = "cannot delete project"
	cannotFindProject   = "cannot find project"
)

// generic errors and wrap with more specific info
// return errors.Wrap(err, ErrEmptyTable)
