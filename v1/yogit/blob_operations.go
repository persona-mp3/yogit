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
// It then locates that branch and retrieves it's current commit, and returns it,
//
// stored in common.CommitId string struct
func GetParentCommit() common.CommitId {
	activeBranch := GetActiveBranch()
	PARENT_COMMIT, err := os.ReadFile(activeBranch)
	if err != nil {
		log.Fatalf("error: occurred trying to read latest commit\n %s\n", err)
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

// UpdateBranch calls GetActiveBranch, truncates the branch and writes the new commit Id to the branch
func UpdateBranch(commitId string) {
	branchPath := GetActiveBranch()
	f, err := os.OpenFile(branchPath, os.O_TRUNC|os.O_WRONLY, 0o755)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// latestCommit := string(c.Id)
	latestCommit := commitId

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

}

// CreateBranch creates a new branch in .yogit/refs/heads/name path
//
// The file returned is for os.O_WRONLY and os.O_CREATE. The caller is responsible
// for closing the file
func CreateBranch(name string) (*os.File, error) {
	// first check if the  branch exists
	branchPath := filepath.Join(common.BRANCH_PATH, name)
	if _, err := os.Stat(branchPath); err == nil {
		return nil, fmt.Errorf("branch already exists in %s ", branchPath)

	}

	f, err := os.OpenFile(branchPath, os.O_WRONLY|os.O_CREATE, 0o755)
	if err != nil {
		fmt.Printf("error: occured in creating new branch: \n %s\n", err)
		return nil, err
	}

	return f, nil
}

// updateHeader takes in the name of a branch provided and writes it to the HEADER file.
//
// Formatted as refs:/refs/heads/master.
func UpdateHeader(branchName string) {

	HEADER, err := os.OpenFile(common.ROOT_HEADER_FILE, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0o777)
	if err != nil {
		fmt.Println("error in updatingHeader")
		log.Fatal(err)
	}
	defer HEADER.Close()
	if _, err := fmt.Fprintf(HEADER, "%s/%s", common.BRANCH_REFS, branchName); err != nil {
		log.Fatal(err)
	}

}
