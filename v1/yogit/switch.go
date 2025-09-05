package yogit

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"v1/common"

	"github.com/liamg/tml"
)

type commitState struct {
	// full path to parent blob
	parent string
	// full path to tree blob
	tree string
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

	// create path to the tree location
	treePath := filepath.Join(common.ROOT_DIR_OBJECTS, commitTree[:2], commitTree[2:])
	parentPath := filepath.Join(common.ROOT_DIR_OBJECTS, commitParent[:2], commitParent[2:])

	return commitState{
		parent: parentPath,
		tree:   treePath,
	}
}

type fileSnapshot struct {
	// perm file permissions in octal notation ie 644
	perm     string
	blobPath string
	fileName string
}

// findState takes the commitState type and uses the tree property to find each blob
func (state commitState) findState() {

	f, err := os.OpenFile(state.tree, os.O_RDONLY, 0555)
	if err != nil {
		log.Fatalf("error: occured in opening tree:\n %s\n", err)
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var s fileSnapshot
		line := scanner.Text()
		line = strings.ReplaceAll(line, " ", "")
		s.perm = line[:3]
		id := parseId(line)
		s.blobPath = filepath.Join(common.ROOT_DIR_OBJECTS, id[:2], id[2:])
		s.fileName = parseName(line, id)
		s.buildState()
	}
	out := tml.Sprintf("<green>Full state restored</green>")
	fmt.Println(out)

}

// parseId extracts the hashId for each line in STAGE file saved at blob level
func parseId(line string) string {
	var id string
	var revId string
	for i := len(line) - 1; i > 0 && len(id) < 40; i-- {
		id += string(line[i])
	}

	for i := len(id) - 1; i >= 0; i-- {
		revId += string(id[i])
	}
	return revId
}

// parseName extracts the fileName for each line in STAGE file saved at blob level
func parseName(line string, id string) string {
	nameHashId := line[3:] // offseting from the octal permission
	fileName, _, _ := strings.Cut(nameHashId, id)
	return fileName
}

// buildState extracts the content from a blob and copies it to the file in the current state.
//
// # When the blob exists, but the file doesnt exist in the root repository, it creates it and copies
//
// the previous content into the newly created file, otherwise, truncates it
func (s fileSnapshot) buildState() {
	src, err := os.Open(s.blobPath)
	if err != nil {
		log.Fatal(err)
	}
	defer src.Close()

	octalPerm, err := strconv.ParseInt(s.perm, 8, 32)
	if err != nil {
		fmt.Printf("could not convert %s to octal\n", s.perm)
		panic(err)
	}
	dst, err := os.OpenFile(s.fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(octalPerm))
	if err != nil {
		log.Fatal(err)
	}
	defer dst.Close()

	io.Copy(dst, src)

	out := tml.Sprintf("<green> Snapshot done for %s </green>", s.fileName)
	fmt.Println(out)
}

func Switch(branchName string) {
	repoState := findBranch(branchName)
	repoState.findState()
}
