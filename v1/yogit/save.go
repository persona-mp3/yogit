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

// stagingArea reads all the files in the current directory and saves them to the stage or file.
func stagingArea(base string) []string {
	var files []string
	dirEntries, err := os.ReadDir(base)
	if err != nil {
		log.Fatalf("error: occured in staging area %s\n", err)
	}

	for _, path := range dirEntries {
		if slices.Contains(common.IGNORE_FILES, path.Name()) {
			continue
		} else if path.IsDir() {

			filepath.WalkDir(path.Name(), func(path string, d fs.DirEntry, err error) error {
				if strings.Contains(path, "/") {
					files = append(files, path)
				} else {
					return nil
				}
				return nil
			})
			continue
		}
		files = append(files, path.Name())
	}

	out := tml.Sprintf("<yellow>Done reading all files</yellow>")
	fmt.Println(out)

	return files
}

// saveBlob converts all file contents into blobs by hashing the content with sha1
func saveToBlob(files []string) {

	for _, path := range files {
		src, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("error: occured in reading %s\n %s\n", path, err)
		}

		hashId := utils.Hasher(src)
		parentFolder, blobName := hashId[:2], hashId[2:]

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

		if _, err := fmt.Fprint(f, string(src)); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("written content successfully. HashID: %s, Parent: %s, Blob: %s\n", hashId, parentFolder, blobFile)
	}

}

func Add(base string) {
	files := stagingArea(base)
	saveToBlob(files)
}
