package yogit

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"v1/common"

	"github.com/liamg/tml"
)

// updateHeader takes in the name of a branch provided and writes it to the HEADER file.
//
// Formatted as refs:/refs/heads/master.
func UpdateHeader(branchName string) {
	// fullBranch := fmt.Sprintf("%s/%s", common.BRANCH_REFS, branchName)
	// fmt.Printf("Updating header to %s \n", fullBranch)

	HEADER, err := os.OpenFile(common.ROOT_HEADER_FILE, os.O_CREATE|os.O_RDWR, 0o777)
	if err != nil {
		fmt.Println("error in updatingHeader")
		log.Fatal(err)
	}
	defer HEADER.Close()
	fmt.Println(HEADER.Name())

	if _, err := fmt.Fprintf(HEADER, "%s/%s", common.BRANCH_REFS, branchName); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("written to HEADER, %s%s\n", common.BRANCH_REFS, branchName)

}

// initFolders create all neccessary base folders and file structure
func initFolders(base string) {
	root := filepath.Join(base, common.ROOT)
	if err := os.Mkdir(root, 0o777); err != nil {
		panic(err)
	}

	rootFolders := []string{common.ROOT_DIR_OBJECTS, common.ROOT_DIR_LOG, common.ROOT_DIR_REFS}
	rootFiles := []string{common.ROOT_STAGE_FILE, common.ROOT_HEADER_FILE, common.ROOT_GLOBALS_FILE}

	refsChildrenFolders := []string{"/heads", "/tags"}

	// creating root folders
	for _, rootFolder := range rootFolders {
		if err := os.MkdirAll(rootFolder, 0o755); err != nil {
			log.Fatalf("error: %s\n", err)
		}
	}

	// creating root files
	for _, rootFile := range rootFiles {
		if _, err := os.Create(rootFile); err != nil {
			log.Fatalf("error: %s\n", err)
		}
	}

	// creating refs children folders
	for _, refFolders := range refsChildrenFolders {
		fullPath := fmt.Sprintf("%s%s", common.ROOT_DIR_REFS, refFolders)
		if err := os.MkdirAll(fullPath, 0o755); err != nil {
			log.Fatalf("error: %s\n", err)
		}
	}

	// creating log file
	if _, err := os.Create(common.LOG_LOGS_FILE); err != nil {
		log.Fatalf("error: %s\n", err)
	}

	// creating default master branch
	if _, err := os.Create(common.MASTER_BRANCH); err != nil {
		log.Fatalf("error: %s\n", err)
	}

	out := tml.Sprintf("<green>Yogit initialised</green>")
	fmt.Println(out)
}

// Init command for yogit creates file structures and creates the master branch by default.
//
// If the .yogit folder has already been initialsed, all execution is stopped
func Init(base string) {
	_, err := os.Stat(common.ROOT)
	if err != nil && os.IsExist(err) {
		log.Fatal(".yogit already exists")
	}

	initFolders(base)
	UpdateHeader("master")
}
