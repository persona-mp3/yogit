package yogit

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"v1/common"
)

type commitState struct {
	parent string
	tree   string
}

func findBranch(name string) commitState {
	currentBranch := GetActiveBranch()
	destPath := filepath.Join(common.BRANCH_PATH, name)

	// check if we're already on the branch
	if currentBranch == destPath {
		fmt.Printf("Already on %s, check-> %s | %s\n", name, currentBranch, destPath)
		return commitState{}
	}

	// check if the branch actually exists
	content, err := os.ReadFile(destPath)
	if err != nil && !os.IsExist(err) {
		fmt.Printf("branch %s does not exist in %s\n", name, destPath)
		fmt.Println(err)
		return commitState{}
	}

	// since the branch exists, lets get the last commit it had
	latestCommit := string(content)
	fmt.Printf("latest commit on: %s is %s\n", name, latestCommit)

	// now search the object folder for the commit details
	commitPath := filepath.Join(common.ROOT_DIR_OBJECTS, latestCommit[:2], latestCommit[2:])
	metadata, err := os.ReadFile(commitPath)
	if err != nil {
		fmt.Println("error: occured trying to read meta data of commit:", commitPath)
		log.Fatal(err)
	}
	commitData := string(metadata)

	_, commitTreeParentMsg, found := strings.Cut(commitData, "tree:")
	if !found {
		fmt.Println(commitTreeParentMsg, found, commitData)
		panic("Could not find the tree and parent in commit data")
	}

	commitTreeParent, _, found := strings.Cut(commitTreeParentMsg, "msg")
	if !found {
		fmt.Println(commitTreeParent, found, commitData)
		panic("Could not find the msg commit data")
	}

	treeParent := strings.ReplaceAll(commitTreeParent, " ", "")  // remove all whitespace
	commitTree, commitParent := treeParent[:40], treeParent[47:] // after first 40hex string, len(parent) so we offset by 7
	fmt.Println("CommitTree:", commitTree)
	fmt.Println("CommitParent", commitParent)

	return commitState{
		parent: commitParent,
		tree:   commitTree,
	}
}

func findState(repoState commitState) {}

func Switch(branchName string) {
	findBranch(branchName)
}
