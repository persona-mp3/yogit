package yogit

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"v1/common"

	"github.com/liamg/tml"
)

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

	// creating default master branch, and writing default commit of 00000000
	f, err := os.Create(common.MASTER_BRANCH)
	if err != nil {
		log.Fatalf("error: occured while creating default master branch %s\n", err)
	}
	defer f.Close()
	if _, err := fmt.Fprintf(f, "000000000000000"); err != nil {
		log.Fatalf("error: occured while writing default commit\n %s\n", err)
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
