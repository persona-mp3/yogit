package main

import (
	"v1/yogit"
)

func main() {
	BASE := "."
	yogit.Init(BASE)
	yogit.Add(BASE)
	yogit.Save("feat: saving commit messages")
}
