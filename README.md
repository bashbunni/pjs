
## Purpose
*A tool for per-project timestamped work logging.*
The point is so that when your boss asks you seven months later why you made a very specific design decision, you can send them the whole list of progress updates on the project throughout its lifecycle. 

## Status
- [] create entry -> functionality is blocked by [this issue](https://github.com/charmbracelet/bubbletea/issues/171)
- [] add scrolling to pagination provided by slides

## Build Instructions

This project tracker can be built with `go build -o <outdir/executablefilename>`

where the output directory and executable name are provided after the `-o` flag. 

ex. 

`go build` - will create an executable named "project-management"

you can execute the file by running `./project-management "my message"`

## Running the Program
navigate to cmd/project-management and run either `go run .` to run the program or `go build` to build the program binary.
If you built the binary, you can run it with `./project-management` or even add it to your PATH so you can just run `project-management` from anywhere. 

## Collaboration
You're welcome to write features and report issues for this project.
It's still a learning project for me, so I can't confirm that all PRs will be merged, but this project is open to contributions. 
