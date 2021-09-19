package utils

import "errors"

/*
need to implement func Error() string for it to implement Error interface
*/

var (
	ErrProjectNotFound  = errors.New("project not found")
	ErrEntryNotFound    = errors.New("entry not found")
	ErrCannotCreateFile = errors.New("cannot create file")
	ErrCannotOpenEditor = errors.New("cannot open editor")
	ErrCannotSaveFile   = errors.New("cannot save file")
	ErrPandoc           = errors.New("cannot pipe stdin to pandoc")
)

const (
	CannotWriteToFilePandoc = "cannot write to output pdf file"
	CannotRunPandoc         = "cannot run pandoc"
)

// TODO: generic errors and wrap with more specific info
// return errors.Wrap(err, ErrProjectNotFound)
