package main

import (
	"log"
	"os"
)

/*
TODO:
- render all entries to markdown
- render specific date frame
- render specific date?
*/

func OutputMarkdown(entries []Entry) {
	file, err := os.OpenFile("./output.md", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 640) // change to current date
	if err != nil {
		log.Fatal("unable to create file. error: %s", err.Error()) // TODO: is this how I want to handle this error?
	}
	for _, entry := range entries {
		file.WriteString(entry.Message)
		if err != nil {
			log.Fatal("unable to write file. error: %s", err.Error())
		}
	}
	file.Close()
}
