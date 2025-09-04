package yogit

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"v1/common"
)

// GetParentCommit gets the latest commit by calling GetActiveBranch and returns
//
// the previous commit made in CommitId struct
func GetParentCommit() common.CommitId {
	activeBranch := GetActiveBranch()
	PARENT_COMMIT, err := os.ReadFile(activeBranch)
	if err != nil {
		log.Fatalf("error: occcurred trying to read latest commit\n %s\n", err)
	}

	return common.CommitId(string(PARENT_COMMIT))
}

// GetActiveBranch gets the current branch the HEADER is pointing to and returns
//
// the resolved path to as a string using filepath.Join ready for file operations.
func GetActiveBranch() string {

	ACTIVE_BRANCH, err := os.ReadFile(common.ROOT_HEADER_FILE)
	if err != nil {
		log.Fatalf("error: occured trying to get parent commit\n %s\n\n", err)

	}

	_, branchName, isFound := strings.Cut(string(ACTIVE_BRANCH)+"/", common.BRANCH_REFS)
	if !isFound {
		fmt.Println(branchName, isFound)
		panic("Something happend")
	}

	// get path to branch
	branchPath := filepath.Join(common.BRANCH_PATH, branchName)

	return branchPath
}

// SaveCommitBlob just like saveBlob saves a commit metadata at the blob level in the objects folder
func (c Commit) SaveCommitBlob() {
	// to save a commit, we'll need to create its path in object folder
	parentFolder := string(c.Id[:2])
	blobName := string(c.Id[2:])

	parentPath := filepath.Join(common.ROOT_DIR_OBJECTS, parentFolder)
	blobPath := filepath.Join(parentPath, blobName)

	if err := os.Mkdir(parentPath, 0o755); err != nil && !os.IsExist(err) {
		log.Fatal("error from method: saveCommit", err)
	}

	f, err := os.Create(blobPath)
	if err != nil {
		log.Fatal("error: occured in creating blob", blobPath)
	}
	defer f.Close()

	logFormat := fmt.Sprintf(
		`id: %s    tree: %s    parent: %s    msg: %s  time: %s `,
		string(c.Id), c.Tree, c.ParentCommit, c.CommitMsg, c.CommittedAt.Format("Jan 2, 1990 3:04 PM"),
	)

	if _, err := fmt.Fprintf(f, "%s", logFormat); err != nil {
		panic(err)
	}

	fmt.Printf("Commit Detetails -> \n %s\n", logFormat)
}

// UpdateBranch calls GetActiveBranch, truncates the branch and writes the new commit Id to the branch
func (c Commit) UpdateBranch() {
	branchPath := GetActiveBranch()
	f, err := os.OpenFile(branchPath, os.O_TRUNC|os.O_WRONLY, 0o755)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	latestCommit := string(c.Id)

	if _, err := fmt.Fprintf(f, "%s", latestCommit); err != nil {
		log.Fatal(err)
	}
}

func (c Commit) UpdateLog() {
	f, err := os.OpenFile(common.LOG_LOGS_FILE, os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	logFormat := fmt.Sprintf(
		`id: %s    Tree: %s    parent: %s    msg: %s  at: %s `,
		string(c.Id), c.Tree, c.ParentCommit, c.CommitMsg, c.CommittedAt.Format("Jan 2, 1990 3:04 PM"),
	)

	if _, err := fmt.Fprintf(f, "%s\n", logFormat); err != nil {
		panic(err)
	}

	fmt.Println("log files updated")
}
