package yogit

import (
	// "compress/gzip"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"
	"io"
	"bytes"
	"crypto/sha1"
)

type Commit struct {
	// Parent *HashId `json:"parent"`
	Id Sha1Hash `json:"hashId"`
	Tree Sha1Hash `json:"tree"`
	Author string `json:"author"`
	Contact string `json:"contact"`
	CommitMsg string `json:"commitMsg"`
	CommittedAt time.Time `json:"committedAt"`
}

type HashId struct {
	Id string
}

type Sha1Hash struct {
	Hash string
}

func YoGit() {
	fmt.Println("WELCOME TO YOGIT")
}

func Init() {
	fmt.Println("standback, making folders")
	// use an input package to set the global variables
	err := os.MkdirAll(".yogit", 0777)
	LogErr(err, "Error in making yogit folders")

	subFolders := []string{"objects", "log", "refs"}
	subFiles := []string{"index", "HEADER", "GLOBALS"}
	logSubFiles := []string{"logs"}
	refSubFolders := []string{"heads", "tags"}

	for _, folder := range subFolders {
		path := fmt.Sprintf(".yogit/%s", folder)
		err := os.MkdirAll(path, 0777)
		LogErr(err, "Error in making subfolders")
	}

	for _, file := range subFiles {
		path := fmt.Sprintf(".yogit/%s",file)
		_, err := os.OpenFile(path, os.O_CREATE | os.O_RDWR |os.O_TRUNC,  0777)
		LogErr(err, "Error making sub file")
	}

	for _, logFile := range logSubFiles {
		path := fmt.Sprintf(".yogit/log/%s", logFile)
		_, err := os.Create(path)
		LogErr(err, "Error making logFile")
	}

	for _, ref := range refSubFolders {
		path := fmt.Sprintf(".yogit/refs/%s", ref)
		err := os.Mkdir(path, 0777)
		LogErr(err, "Error in making refSubFolders")
	}

	// make master branch by defualt and write it to the HEADER in .yogit
	master, err := os.Create(".yogit/refs/heads/master")
	LogErr(err, "Error in making master branch in heads")
	master.Close()

	header, err := os.OpenFile(".yogit/HEADER", os.O_WRONLY | os.O_TRUNC, 0777)
	LogErr(err, "Error in finding HEADER")
	defer header.Close()

	n, err := fmt.Fprintf(header, "ref:refs/heads/master")
	LogErr(err, "Error in making default write to header")
	fmt.Printf("number of bytes written to HEADER -> %d", n)

	fmt.Println("DONE")
	fmt.Printf("So the magic that just happend is that key folders in this directory ./.yogit\n")
}

func updateBranch(hash HashId) {
	// Find what the current HEAD is pointing at
	header, err := os.ReadFile(".yogit/HEADER")
	LogErr(err, "Error in reading HEADER")
	pointingTo := string(header)
	// expect return to be refs/heads/master or any other branch
	_, activeBranch, _:= strings.Cut(pointingTo, ":")

	pathTo := fmt.Sprintf(".yogit/%s", activeBranch)
	branch, err := os.OpenFile(pathTo, os.O_CREATE | os.O_TRUNC | os.O_WRONLY, 0777)
	LogErr(err, "Error in finding the active branch")
	defer branch.Close()

	n, err := branch.Write([]byte(hash.Id))
	LogErr(err, "Error in updating branch")

	fmt.Printf("updated your current branch with the your state, now the header knows where you are and your state\nN bytes written %d\n", n)
}

// this function is called  in StaginArea(), it serves as hashing and saving to the object store at blob level
// so for every file, we get their hash and compressed content, and write their hash and name to the STAGE file 
// the file writing is done in StagingArea()
// the tree object is created by TreeObject() which basically does the same thing as Add() but its the index file itself
func SaveToObject(file string) Sha1Hash{
	fmt.Println("adding all files onto the staging area")
	// write all current files to staging area : index, but for now, lets just add one
	// open the file that wants to be saved

	content, err := os.ReadFile(file)
	LogErr(err, "Error occured, could not find file specified")
	// hash the content, get the first 2 names
	hasher := sha1.New() 
	hasher.Write([]byte(content))
	hashedContent := hasher.Sum(nil)
	hashId := hex.EncodeToString(hashedContent)
	objectFolder := hashId[:2]

	saveAt := fmt.Sprintf(".yogit/objects/%s", objectFolder)
	errM := os.Mkdir(saveAt, 0777)
  LogErr(errM, "check Add()")
	
	blobPath := fmt.Sprintf("%s/%s", saveAt, objectFolder[2:])
	blobName := fmt.Sprintf("%s%s", blobPath, hashId[2:])
	fmt.Println(blobName)
	blob, err := os.Create(blobName)
	LogErr(err, "Error in making blob")

	// gzipWriter := gzip.NewWriter(blob)
	// gzipWriter.Write(content)
	// gzipWriter.Close()
	byteReader := bytes.NewReader(content)
	io.Copy(blob, byteReader)
	
	fmt.Println("hashed_content  --- ",hashId)
	fmt.Println("name to store folder --- ",objectFolder)

	// fmt.Printf("\nFile compressed successfully...\n")
	return Sha1Hash{Hash:hashId}
}

func StagingArea() {
	// we need to write all current files into the staging area in .yogit/stage
	dirEntries, err := os.ReadDir(".")
	LogErr(err, "Error in reading StagingArea(), reading directory")

	stage, err := os.OpenFile(".yogit/stage", os.O_CREATE | os.O_APPEND | os.O_RDWR, 0777)
	LogErr(err, "Error in opening stage file")
	defer stage.Close()

	for _, path := range dirEntries {
		if path.Name() == ".git" || path.Name() == ".yogit" || path.IsDir() {
			continue
		} 

		info, err := path.Info()
		LogErr(err, "Error in getting file info")

		hashId := SaveToObject(info.Name())
		_, errw := fmt.Fprintf(stage, "%04o %-20s %20s\n", info.Mode().Perm(), hashId, info.Name() )
		LogErr(errw, "Error in writing to file")
	}
	fmt.Println("check ./yogit/stage")
}

func SaveCommit(msg string) {
	// save commit should create a commit data type reprsented by struct as 
	// {Author string, Tree Sha1Hash}
	// but before that the Sha1Hash is gotten from the SaveToObj() as it hashes this data and returns the hash
	// we can find where it is based on this.
	
	fmt.Println("processing index file")
	const INDEX_PATH = ".yogit/stage"

	content, err := os.ReadFile(INDEX_PATH)
	LogErr(err, "Error occured, could not find file specified")

	treeHash := HasherFn(content)
	hashId := treeHash.Hash
	objectFolder := treeHash.Hash[:2]

	saveAt := fmt.Sprintf(".yogit/objects/%s", objectFolder)
	errM := os.Mkdir(saveAt, 0777)
  LogErr(errM, "check Add()")
	
	treePath := fmt.Sprintf("%s/%s", saveAt, objectFolder[2:])
	treeName := fmt.Sprintf("%s%s", treePath, hashId[2:])
	fmt.Println(treeName)
	tree, err := os.Create(treeName)
	LogErr(err, "Error in making tree")
	defer tree.Close()

	byteReader := bytes.NewReader(content)
	io.Copy(tree, byteReader)
	
	fmt.Println("hashed_content for tree is store at --- ",hashId)

	// create a hash for commit data
	commit := Commit{}
	commit.Author = "archive@persona-mp3jpeg"
	// commit.Id = 
	commit.Contact = "<archive@persona-mp3>"
	commit.Tree = treeHash
	commit.CommitMsg = msg
	commit.CommittedAt = time.Now()

	// save commit as it is 
	s := fmt.Sprintf("%v", commit)
	commitHash := HasherFn([]byte(s))
	// commit.Id = commitHash
	SaveCommitToObj_u(commitHash, commit)
}
