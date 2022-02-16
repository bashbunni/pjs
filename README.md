
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=bashbunni-pjm&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=bashbunni-pjm)

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

### Entry Commands

#### Create Entry
`./project-management <projID> entry -ce`

For example:
`./project-management 2 entry -ce`

#### Delete Entry
`./project-management <projID> entry -de <entryID>`

For example:
`./project-management 2 entry -de 15`

### Output Commands
#### Output to PDF
`./project-management <projID> output -pdf`

For example:
`./project-management 2 output -pdf`

#### Output to Markdown
`./project-management <projID> output -md`

For example:
`./project-management 2 output -md`

### Project Commands
#### List All Projects
`./project-management <projID> project -lp`

For example:
`./project-management 2 project -lp`

#### Delete Project
`./project-management <projID> project -dp`

For example:
`./project-management 2 project -dp`

#### Rename Project (broken)
`./project-management <projID> project -ep "New Project Name"`

For example:
`./project-management 2 project -ep "New Project Name"`


## Collaboration
Because this is a learning project, I'm not currently accepting pull requests. I will be accepting pull requests on this project in the near future, so stay tuned!
