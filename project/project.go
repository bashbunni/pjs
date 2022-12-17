package project

import (
	"bufio"
	"fmt"
	"os"
)

// NewProjectPrompt create a new project from user input to console
func NewProjectPrompt() string {
	var name string
	fmt.Println("what would you like to name your project?")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	name = scanner.Text()
	return name
}
