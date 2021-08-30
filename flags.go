package main

import (
	"flag"
	"time"
)

var (
	// stringvar := flag.String("optionname", "defaultvalue", "description of the flag")
	cEntry      = flag.Int("ce", -1, "create a new entry for a project; default projID is -1")
	deleteEntry = flag.Int("de", -1, "delete an existing entry; default is -1")
	listProj    = flag.Bool("lp", false, "display all projects")
	deleteProj  = flag.Int("dp", -1, "delete an existing project; default is -1")
	editProj    = flag.Int("ep", -1, "rename an existing project; default is empty string")
	markdown    = flag.Bool("md", false, "output all entries to markdown file")
	pdf         = flag.Bool("pdf", false, "output all entries to pdf file")
	start       = flag.String("s", "", "start date for date range")
	end         = flag.String("e", time.Now().Format("2006-01-02"), "end date for date range")
)
