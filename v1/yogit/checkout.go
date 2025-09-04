package yogit

import (
	"fmt"
	"github.com/liamg/tml"
	"log"
)

// To checkout we need to change the HEADER to point to a new branch per-say v1
//
// - We'd have to do get the current branch the HEADER is pointing at
//
// - Get the latest commit from the current branch too
//
// - Change the header to point to the new branch
//
// - Update the new branch to point to the latest commit
func Checkout(newBranch string) {
	currCommit := string(GetParentCommit())
	f, err := CreateBranch(newBranch)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// writing current commit to new branch
	if _, err := fmt.Fprintf(f, "%s", currCommit); err != nil {
		log.Fatal(err)
	}

	// and then we can update the HEADER
	UpdateHeader(newBranch)

	tml.Printf("<green> *%s</green>\n", newBranch)
}
