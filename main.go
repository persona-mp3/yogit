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
		if args[1] == "init" {
			yogit.Init()
			return
		}

		// if args[1] == "logs" {
		// 	yogit.Init()
		// 	return
		// }
		
		if args[1] == "clogs" {
			yogit.CheckLogs()
			return
		}
		return
	}

	case cmds == 3: {
		// fmt.Println("command line args has two arguments passed in")
		if args[1] == "save" && args[2] != "" {
			yogit.SaveCommit(args[2])
			return
		}

		if args[1] == "add" && args[2] != "" {
			// yogit.Add(args[2])
			yogit.StagingArea()
			return
		}

		if args[1] == "ntimeline" && args[2] != "" {
			yogit.NewTimeLine(args[2])
		}

		if args[1] == "travelto" && args[2] != "" {
			yogit.TravelTo(args[2])
		}

		if args[1] == "switchto" && args[2] != "" {
			yogit.SwitchTo(args[2])
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
