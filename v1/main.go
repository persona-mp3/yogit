package main

import (
	"fmt"
	"os"
	"v1/yogit"
)

var man string = `
yogit is a compact and friendly interface to Git, implementing the same lower-level 
linux-based APIs. 

Usage:
	yogit init
	yogit add .
	yogit new -b <branch-name>
	yogit save <message>
	yogit switch <branch-name>
	yogit help


	yogit init 
	Initialises a new .yogit folder in the current directory, synonymous to git init

	yogit add . 
	Saves the current working state of the repository, excluding node_modules, cache_files and dot files.
	Synonymous with git add .

	yogit save <message> 
	Attaches the saved state to a message, synonymous to git commit -m <message>

	yogit new -b <branch-name>
	Creates a new branch. By default a master branch is created, if you'd want to make a new branch 
	called 'playground' simply enter:
		yogit new -b playground
	Which means, yogit, make a new branch playground and take me there

	yogit switch <branch-name>
	To go back to a previous branch, it returns the current state of the branch to the last state of the destination 
	branch. 
	Example
		All branches -> master, playground, shenanigance
		Current-Branch -> playground
		yogit switch master

		Now all files that were in master are now available to you, explicitly. 
`

func cmdArgs(args []string, BASE string) {
	if len(args) < 1 {
		fmt.Println("fatal: invalid arguments, run yogit help for instructions")
		return
	}

	switch cmds := len(args); cmds {
	case 2:
		if args[1] == "help" {
			fmt.Println(man)
		} else if args[1] == "init" {
			yogit.Init(BASE)
		} else {
			fmt.Println("Invalid use", man)
		}
	case 3:
		if args[1] == "add" && args[2] == "." {
			yogit.Add(BASE)
		} else if args[1] == "save" && len(args[2]) > 1 {
			yogit.Save(args[2])
		} else if args[1] == "switch" && len(args[2]) > 1 {
			yogit.Switch(args[2])
		} else {
			fmt.Println("Invalid use", man)
		}
	case 4:
		if args[1] == "new" && args[2] == "-b" && len(args[3]) > 1 {
			yogit.Checkout(args[3])
		} else {
			fmt.Println("Invalid use", man)
		}
	default:
		fmt.Println("Invalid use", man)
	}

}

func main() {
	BASE := "."
	cmdArgs(os.Args, BASE)
}
