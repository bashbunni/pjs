package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"gorm.io/gorm"
)

const divider = "_______________________________________"

func getOutput(entries []Entry) []byte {
	var output string
	for _, entry := range entries {
		output += fmt.Sprintf("ID: %d\nCreated: %s\nMessage:\n %s\n %s\n", entry.ID, entry.CreatedAt.Format("2006-01-02"), entry.Message, divider)
	}
	return []byte(output)
}

func OutputMarkdownRange(start time.Time, end time.Time, db *gorm.DB) {
	entries := GetEntriesByDate(start, end, db)
	OutputMarkdown(entries)
}

func OutputMarkdown(entries []Entry) error {
	file, err := os.OpenFile("./output.md", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close() // want defer as close to acquisition of resources as possible
	output := getOutput(entries)
	_, err = file.Write(output)

	if err != nil {
		return err
	}
	return nil
}

func OutputPdf(entries []Entry) error {
	output := getOutput(entries)                               // []byte
	pandoc := exec.Command("pandoc", "-s", "-o", "output.pdf") // c is going to run pandoc, so I'm assigning the pipe to c
	wc, wcerr := pandoc.StdinPipe()                            // io.WriteCloser, err
	if wcerr != nil {
		return wcerr
	}
	goerr := make(chan error)
	done := make(chan bool)
	go func() {
		defer wc.Close()
		_, err := wc.Write(output)
		fmt.Println("I wrote things I think")
		goerr <- err
		close(goerr)
		close(done)
	}()
	if err := <-goerr; err != nil {
		return err
	}
	fmt.Println("I'm about to run things")
	err := pandoc.Run()
	if err != nil {
		return err
	}
	return nil
}
