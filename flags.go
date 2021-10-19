package main

import (
	"flag"
)

// TODO: init proj -> then support subcommands

var (
	// entry
	entryCommands = flag.NewFlagSet("entry", flag.ExitOnError)
	createEntry   = entryCommands.Bool("ce", false, "create a new entry for the selected project")
	deleteEntry   = entryCommands.Uint("de", 0, "delete an existing entry; default is 0") // TODO: check this is true

	// output
	outputCommands = flag.NewFlagSet("output", flag.ExitOnError)
	markdown       = outputCommands.Bool("md", false, "output all entries to markdown file")
	pdf            = outputCommands.Bool("pdf", false, "output all entries to pdf file")

	// project
	// TODO: make this not need the project ID for arg1
	projectCommands = flag.NewFlagSet("project", flag.ExitOnError)
	listAllProjects = projectCommands.Bool("lp", false, "display all projects")
	deleteProject   = projectCommands.Bool("dp", false, "delete an existing project; default is -1")
	editProject     = projectCommands.Bool("ep", false, "rename an existing project; default is empty string")
	// TODO: don't prompt for user entry
)
