package main

import (
	"fmt"
	"os"
)

/*
TODO:
- render all entries to markdown
- render specific date frame
- render specific date?
*/

func OutputMarkdown(entries []Entry) error {
	file, err := os.OpenFile("./output.md", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close() // want defer as close to acquisition of resources as possible
	var output string
	for _, entry := range entries {
		output += entry.Message + "\n"
	}
	fmt.Println(output)
	_, err = file.WriteString(output)
	if err != nil {
		return err
	}
	return nil
}
