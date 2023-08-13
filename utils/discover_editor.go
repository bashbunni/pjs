package utils

import (
	"errors"
	"os/exec"
)

func GetAvailableEditor() (string, error) {
	editors := []string{"vim", "nvim", "emacs", "hx", "nano", "vi"}

	for _, editor := range editors {
		if _, err := exec.LookPath(editor); err == nil {
			return editor, nil
		}
	}

	return "", errors.New("no editor found. Please install a suitable editor, or specify an existing one using the EDITOR environment variable")
}
