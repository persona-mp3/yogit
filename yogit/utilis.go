package yogit

import (
	"fmt"
	"log"
	"os"
	"crypto/sha1"
	"strings"
	"encoding/hex"
	"github.com/liamg/tml"
)

const (
	OBJECT_PATH = ".yogit/objects"
	BRANCH_PATH = ".yogit/refs/heads"
	HEADER_PATH = ".yogit/HEADER"
)

func LogErr(err error, msg string) {
	if err != nil {
		yoErr := tml.Sprintf("Error: <red>%s</red>", msg)
		fmt.Println(yoErr)
		log.Fatal(err)
	}
}

func HasherFn(content []byte) Sha1Hash {
	hasher := sha1.New()
	hasher.Write(content)
	hashedContent := hasher.Sum(nil)

	sha1Hash := hex.EncodeToString(hashedContent)

	return Sha1Hash{Hash: sha1Hash}
}

func SaveCommitToObj_u(commitHash Sha1Hash, commit Commit) {
	folderName := commitHash.Hash[:2]
	fileName := commitHash.Hash[2:]

	savePath := fmt.Sprintf("%s/%s", OBJECT_PATH, folderName)
	filePath := fmt.Sprintf("%s/%s", savePath, fileName)

	err := os.Mkdir(savePath, 0777)
	LogErr(err, "Error in making path for commit")

	// creating commit file 
	commitFile, err := os.OpenFile(filePath, os.O_CREATE | os.O_RDWR, 0777)
	LogErr(err, "Error making commit file")
	defer commitFile.Close()
	commit.Id = commitHash
	// commit.Parent = Parent

	_, errW := fmt.Fprintf(commitFile, "parent:%s author:%s  commitHash:%s  commitMessage:%s tree:%s committedAt:%s\n", 
													commit.Parent.Hash, commit.Author, commit.Id.Hash, commit.CommitMsg,
													commit.Tree.Hash, commit.CommittedAt.Format("Jan 2, 1990 3:04 PM"))

	LogErr(errW, "Error writing to commit file")

	fmt.Printf("successfully wrote commit to file at %s | %s\n", folderName, fileName)
}


func GetParentCommit_u() Sha1Hash{
	// we basically just want to get the previous commit 
	HEAD, err := os.ReadFile(HEADER_PATH)	
	LogErr(err, "Error finding HEAD in SaveCommit()")
	_, activeBranch, _ := strings.Cut(string(HEAD), ":")

	pathTo := fmt.Sprintf(".yogit/%s", activeBranch)
	prevCommitHash, err := os.ReadFile(pathTo)
	LogErr(err, "Error finding active branch in SaveCommit()")

	return Sha1Hash{Hash: string(prevCommitHash)}
}
