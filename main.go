package main 

import (
	"fmt"	
	"os"
	"yogit/yogit"
)

func ReadCommandLine(args []string) {
	cmds := len(args)

	if cmds <= 1 {
		yogit.YoGit()
		// fmt.Println("Print instruction manual or enter interactive mode?")
		return
	}

	switch {

	case cmds == 2: {
		// fmt.Println("command line args passed init")
		if args[1] == "init" {
			yogit.Init()
			return
		}

		if args[1] == "add" {
			yogit.Add()
			return
		}

		return
	}

	case cmds == 3: {
		// fmt.Println("command line args has two arguments passed in")
		if args[1] == "save" && args[2] != "" {
			yogit.CommitFunc(args[2])
			return
		}

		return
	}

	default  :
		fmt.Println("idk what this guy entered -> ", cmds)
		return
	}

}

func main() {
	
	ReadCommandLine(os.Args)
}
