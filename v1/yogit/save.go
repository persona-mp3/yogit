package yogit

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
	"v1/common"
	"v1/utils"
)

type HashId string

// Blob is used to represent all files that will be saved in the objects folder
type Blob struct {
	// Id represents the sha1Hash of file converted to HexString
	Id HashId

	// Parent represents the full path where BlobName will be saved in object folder
	Parent HashId

	// BlobName represents the full path to where the blob will be created
	BlobName string
}
type Commit struct {
	Id           common.CommitId
	Tree         HashId
	ParentCommit common.CommitId
	CommitMsg    string
	CommittedAt  time.Time
}

// saveState hashes the stage file and it's content
//
// This is because the staage file represents the current state of
// the repository after the Add command.
//
// # This is the tree level as we can now backtrack a specific commmit with
//
// this stage file as it has mappings for each file and their hashId
func saveState() Blob {
	STATE, err := os.ReadFile(common.ROOT_STAGE_FILE)
	if err != nil {
		log.Fatalf("error: occured in opeing STAGE file for saving \n %s\n", err)
	}

	var blob Blob
	fileHash := string(utils.Hasher(STATE))

	// create parent path
	parent := filepath.Join(common.ROOT_DIR_OBJECTS, fileHash[:2])
	blob.Parent = HashId(parent)

	// create blob's own path in relation to parent
	blobPath := filepath.Join(parent, fileHash[2:])
	blob.BlobName = blobPath
	blob.Id = HashId(fileHash)

	// first create the parent directory for the blob
	if err := os.Mkdir(parent, 0o755); err != nil && !os.IsExist(err) {
		log.Fatalf("error occured trying to create blob parent\n %s\n", err)
	}

	f, err := os.Create(blobPath)
	if err != nil {
		log.Fatalf("error occured trying to create tree file\n %s\n", err)
	}
	defer f.Close()

	if _, err := fmt.Fprintf(f, "%s", STATE); err != nil {
		log.Fatal(err)
	}

	return blob
}

// Create a new commit struct including its meta data to be stored
func commitMeta(commitMsg string, treeBlob Blob) Commit {
	c := Commit{}

	// so we need to get the previous parent commit

	c.ParentCommit = GetParentCommit()
	c.CommitMsg = commitMsg
	c.CommittedAt = time.Now()
	c.Tree = treeBlob.Id

	// save the whole commmit meta data by hashing the struct as is
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	if err := enc.Encode(c); err != nil {
		log.Fatalf("error occured trying to encode commit meta-data\n %s\n", err)
	}

	hashId := utils.Hasher(buffer.Bytes())
	c.Id = common.CommitId(hashId)
	return c

	// save commit meta data to objects at blob level then
	// update the current branch and log file with this commit
}

// Save command takes in the commit message and saves it to the blob level alongside updating
//
// logs, and branches. It starts by hashing the index file, calling saveState(), and then
//
// creates a commit type with its metadata, commitMeta().
//
// The commit is then saved at the blob level, updating the branch and log with the latest commit.
func Save(commitMsg string) {
	treeBlob := saveState()
	c := commitMeta(commitMsg, treeBlob)
	c.SaveCommitBlob()
	c.UpdateBranch()
	c.UpdateLog()

}
