package data

import (
	"io"

	"github.com/spf13/afero"
	"github.com/bashbunni/project-management/database/dbconn"
)

func NewStdinPlainReader(readFrom io.Reader) stdinPlainReader {
	return stdinPlainReader{
		readFrom: readFrom,
	}
}

func OpenDBConnection(path string) (dbconn.GormWrapper, error) {
	return constOpenDBConnection(path)
}

func OverloadUC(overload func() (string, error)) func() {
	ucRef := uc
	uc = overload
	return func() { uc = ucRef }
}

func OverloadFS(overload afero.Fs) func() {
	fsRef := fs
	fs = overload
	return func() { fs = fsRef }
}

func OverloadOpenDBConnection(overload func(string) (dbconn.GormWrapper, error)) func() {
	openDBConnectionRef := openDBConnection
	openDBConnection = overload
	return func() { openDBConnection = openDBConnectionRef }
}

func OverloadPlainPromptReader(overload plainReader) func() {
	plainPromptReaderRef := plainPromptReader
	plainPromptReader = overload
	return func() { plainPromptReader = plainPromptReaderRef }
}

func OverloadPasswordPromptReader(overload passwordReader) func() {
	passwordPromptReaderRef := passwordPromptReader
	passwordPromptReader = overload
	return func() { passwordPromptReader = passwordPromptReaderRef }
}
