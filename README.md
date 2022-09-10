# PJ - Project Journal

Staying organized in the easiest way possible.

## Purpose
*A tool for per-project timestamped work logging.*
The point is so that when your boss asks you seven months later why you made a very specific design decision, you can send them the whole list of progress updates on the project throughout its lifecycle. 

It backs up to a sqlite db which could easily be backed up.

## Build Instructions

If you're new to Go, check out ["How to Build and Install Go Programs"](https://www.digitalocean.com/community/tutorials/how-to-build-and-install-go-programs)

## Running the Program
navigate to the root directory of the project and run either `go run .` to run the program or `go build` to build the program binary.
If you built the binary, you can run it with `./project-management` or even add it to your PATH so you can just run `project-management` from anywhere. 

## TODOs

- [ ] print to PDF
- [ ] print to markdown
- [ ] backup to `charm-fs`

### Entry View
- [ ] hold entries in entry Model, convert that to string instead of all entries to a single string
- [ ] look over fields, cut unnecessary ones

## Collaboration
You're welcome to write features and report issues for this project.
