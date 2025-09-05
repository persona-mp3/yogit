package main

import (
	"v1/yogit"
)

func main() {
	BASE := "."
	yogit.Init(BASE)
	yogit.Add(BASE)
	yogit.Checkout("version1")
	yogit.Save("feat: added switch feature")
	yogit.Checkout("refactor")
	// yogit.Switch("version1")
}
