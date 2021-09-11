package main

import (
	"flag"
)

// TODO: init proj -> then support subcommands

var (
	// entry
	entryCommands = flag.NewFlagSet("entry", flag.ExitOnError)
	createEntry   = entryCommands.Bool("ce", false, "create a new entry for a project; default projID is -1")
	deleteEntry   = entryCommands.Bool("de", false, "delete an existing entry; default is -1")

	// output
	outputCommands = flag.NewFlagSet("output", flag.ExitOnError)
	markdown       = outputCommands.Bool("md", false, "output all entries to markdown file")
	pdf            = outputCommands.Bool("pdf", false, "output all entries to pdf file")

	// project
	projectCommands = flag.NewFlagSet("project", flag.ExitOnError)
	listAllProjects = projectCommands.Bool("lp", false, "display all projects")
	deleteProject   = projectCommands.Bool("dp", false, "delete an existing project; default is -1")
	editProject     = projectCommands.Bool("ep", false, "rename an existing project; default is empty string")
)
