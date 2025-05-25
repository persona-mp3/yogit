package yogit

import (
	"fmt"
	"log"
	"os"
	"crypto/sha1"
	"encoding/hex"
)

const (
	OBJECT_PATH = ".yogit/objects"
)


func LogErr(err error, msg string) {
	if err != nil {
		fmt.Println("\n", msg)
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

	_, errW := fmt.Fprintf(commitFile, "author:%s  id:%s  message:%s tree:%s at:%s\n", commit.Author, commit.Id, commit.CommitMsg,commit.Tree, commit.CommittedAt.Format("Jan 2, 1990 3:04 PM"))
	LogErr(errW, "Error writing to commit file")

	fmt.Printf("successfully wrote commit to file at %s | %s\n", folderName, fileName)
}
