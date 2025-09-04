package yogit

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"v1/common"
	"v1/utils"

	"github.com/liamg/tml"
)

type stagingInfo struct {
	name   string
	hashId string
	perm   fs.FileMode
}

var stageFiles []stagingInfo

// stagingArea reads all the files in the current directory and saves them to the stage or file.
func stagingArea(base string) {
	dirEntries, err := os.ReadDir(base)
	if err != nil {
		log.Fatalf("error: occured in staging area %s\n", err)
	}

	for _, path := range dirEntries {
		fileInfo := stagingInfo{}
		if slices.Contains(common.IGNORE_FILES, path.Name()) {
			continue
		} else if path.IsDir() {

			filepath.WalkDir(path.Name(), func(path string, d fs.DirEntry, err error) error {
				if strings.Contains(path, "/") {
					fileInfo.name = path
					info, _ := d.Info()
					fileInfo.perm = info.Mode().Perm()
					stageFiles = append(stageFiles, fileInfo)
				}
				return nil
			})
			continue
		}

		fileInfo.name = path.Name()
		fsInfo, _ := path.Info()
		fileInfo.perm = fsInfo.Mode().Perm()
		stageFiles = append(stageFiles, fileInfo)
	}

}

// saveBlob converts all file contents into blobs by hashing the content with sha1
func saveToBlob() {

	for idx, path := range stageFiles {
		src, err := os.ReadFile(path.name)
		if err != nil {
			log.Fatalf("error: occured in reading %s\n %s\n", path.name, err)
		}

		// Get hashed content
		hashId := utils.Hasher(src)
		parentFolder, blobName := hashId[:2], hashId[2:]

		stageFiles[idx].hashId = hashId

		// Create parent folder
		objectParent := filepath.Join(common.ROOT_DIR_OBJECTS, parentFolder)
		if err := os.Mkdir(objectParent, 0o755); err != nil && !os.IsExist(err) {
			log.Fatalf("error: occured in creating object parent: %s\n %s\n", objectParent, err)
		}

		// Begin writing src to the blob file
		blobFile := filepath.Join(objectParent, blobName)
		f, err := os.Create(blobFile)
		if err != nil {
			log.Fatalf("error: occured in creating object parent: %s\n %s\n", objectParent, err)
		}
		defer f.Close()

		if _, err := fmt.Fprint(f, string(src)); err != nil {
			log.Fatal(err)
		}

	}

}

// The Add command simply saves the current state of the whole repostiory. When executed, all files are saved
//
// at the blob level except node_modules and more. After that, the STAGE area is used to store the map of all
// these files
//
// All files, their meta-data ie file-permissions, path and folder structure are recorded, along side their hashId
//
// for future mapping and referencing when backtracking
func Add(base string) {
	stagingArea(base)
	saveToBlob()

	// Open stage file
	f, err := os.OpenFile(common.ROOT_STAGE_FILE, os.O_CREATE|os.O_WRONLY, 0o766)
	if err != nil {
		log.Fatal("error: ", err)
	}

	for _, stageInfo := range stageFiles {
		if _, err := fmt.Fprintf(f, "%s %s %s\n", stageInfo.perm, stageInfo.name, stageInfo.hashId); err != nil {
			panic(err)
		}
	}

	out := tml.Sprintf("<green>Staged all files</green>")
	fmt.Println(out)
}
