package utils

import "errors"

/*
need to implement func Error() string for it to implement Error interface
*/

var (
	ErrEmptyTable = errors.New("empty table")
	ErrProjectNotFound  = errors.New("project not found")
	ErrCannotCreateProject  = errors.New("cannot create project")
	ErrCannotDeleteProject = errors.New("cannot delete project")
	ErrEntryNotFound    = errors.New("entry not found")
	ErrCannotCreateFile = errors.New("cannot create file")
	ErrCannotOpenEditor = errors.New("cannot open editor")
	ErrCannotSaveFile   = errors.New("cannot save file")
	ErrPandoc           = errors.New("cannot pipe stdin to pandoc")
	// creation errors
)

const (
	CannotWriteToFilePandoc        = "cannot write to output pdf file"
	CannotRunPandoc                = "cannot run pandoc"
	CannotCreateProjectWithEntries = "cannot create ProjectWithEntries"
	CannotUpdateEntries            = "cannot update entries"
	CannotDeleteEntry              = "cannot delete entry"
)

// TODO: generic errors and wrap with more specific info
// return errors.Wrap(err, ErrProjectNotFound)
