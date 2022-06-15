# PJ - Project Journal

Staying organized in the easiest way possible.

## Purpose
*A tool for per-project timestamped work logging.*
The point is so that when your boss asks you seven months later why you made a very specific design decision, you can send them the whole list of progress updates on the project throughout its lifecycle. 

## Build Instructions

This project tracker can be built with `go build -o <outdir/executablefilename>`

where the output directory and executable name are provided after the `-o` flag. 

ex. 

`go build` - will create an executable named "project-management"

you can execute the file by running `./project-management "my message"`

## Running the Program
navigate to the root directory of the project and run either `go run .` to run the program or `go build` to build the program binary.
If you built the binary, you can run it with `./project-management` or even add it to your PATH so you can just run `project-management` from anywhere. 

## Collaboration
You're welcome to write features and report issues for this project.
