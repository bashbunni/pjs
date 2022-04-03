package entry

/*
need to implement func Error() string for it to implement Error interface
*/

const (
	cannotDeleteEntry          = "cannot delete entry"
	cannotCreateEntry          = "cannot create entry"
	cannotFindProject          = "cannot find project"
	errCannotCreateFile        = "cannot create file"
	errCannotSaveFile          = "cannot save file"
	errPandoc                  = "cannot pipe stdin to pandoc"
	errCannotWriteToFilePandoc = "cannot write to output pdf file"
	errCannotRunPandoc         = "cannot run pandoc"
)

// generic errors and wrap with more specific info
// return errors.Wrap(err, ErrEmptyTable)
