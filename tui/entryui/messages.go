package entryui

import "os"

type errMsg struct{ error } // TODO: have this implement Error()
type updateEntryListMsg struct{input []byte}
type updatedMsg struct{}

type editorFinishedMsg struct {
	err error
	file *os.File
}
