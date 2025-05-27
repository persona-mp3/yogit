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
	"log"
	"bufio"
	
)

type Commit struct {
	Parent Sha1Hash `json:"parent"`
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

const (
	LOG_PATH = ".yogit/log/logs"
)

func Init() {
	fmt.Println("standback, making folders")
	// use an input package to set the global variables
	err := os.MkdirAll(".yogit", 0777)
	LogErr(err, "Error in making yogit folders")

	subFolders := []string{"objects", "log", "refs"}
	subFiles := []string{"stage", "HEADER", "GLOBALS"}
	logSubFiles := []string{"logs"}
	// logSubFolders := []string{"refs"}
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

func updateBranch(hash Sha1Hash) {
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

	n, err := branch.Write([]byte(hash.Hash))
	LogErr(err, "Error in updating branch")

	fmt.Printf("updated your current branch with the your state, now the header knows where you are and your state\nN bytes written %d\n", n)
}

func updateLog(c Commit) {
	logFile, err := os.OpenFile(LOG_PATH, os.O_CREATE | os.O_RDWR | os.O_APPEND, 0777)
	LogErr(err, "Error opening log file: updateLog()")
	defer logFile.Close()

	fmtC:= fmt.Sprintf( "author:%s  id:%s  message:%s tree:%s at:%s\n", c.Author, c.Id.Hash, c.CommitMsg, c.Tree.Hash, c.CommittedAt.Format("Jan 2, 1990 3:04 PM"))
	logger := log.New(logFile, "", 0)
	logger.Println(fmtC)
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

// this function is called  in StaginArea(), it serves as hashing and saving to the object store at blob level
// so for every file, we get their hash and compressed content, and write their hash and name to the STAGE file 
// the file writing is done in StagingArea()
// the tree object is created by TreeObject() which basically does the same thing as Add() but its the index file itself
func SaveToObject(file string) Sha1Hash{
	fmt.Println("adding all files onto the staging area")

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
	// fmt.Println(blobName)
	blob, err := os.Create(blobName)
	LogErr(err, "Error in making blob")

	// gzipWriter := gzip.NewWriter(blob)
	// gzipWriter.Write(content)
	// gzipWriter.Close()
	byteReader := bytes.NewReader(content)
	io.Copy(blob, byteReader)
	
	fmt.Printf("%s saved at %s\n", file,  hashId)

	return Sha1Hash{Hash:hashId}
}

func SaveCommit(msg string) {
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
	// we basically just want to get the previous commit 
	HEAD, err := os.ReadFile(HEADER_PATH)	
	LogErr(err, "Error finding HEAD in SaveCommit()")
	_, activeBranch, _ := strings.Cut(string(HEAD), ":")

	pathTo := fmt.Sprintf(".yogit/%s", activeBranch)
	prevCommitHash, err := os.ReadFile(pathTo)
	LogErr(err, "Error finding active branch in SaveCommit()")
	
	parent := Sha1Hash {
		Hash: string(prevCommitHash),
	}

	commit.Parent = parent
	SaveCommitToObj_u(commitHash, commit)

	updateBranch(commitHash)
	commit.Id = commitHash
	updateLog(commit)
}

func NewTimeLine(timeLine string) {
	parent := GetParentCommit_u()
	HEAD, err := os.OpenFile(HEADER_PATH, os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("HEAD FILE NOT FOUND IN", HEADER_PATH)
		panic(err)
	}
	defer HEAD.Close()
	
	refs := fmt.Sprintf("ref:refs/heads/%s", timeLine)
	HEAD.Truncate(0)
	HEAD.Write([]byte(refs))
	newTimeLine := fmt.Sprintf("%s/%s", BRANCH_PATH, timeLine)

	
	newBranch, err := os.OpenFile(newTimeLine, os.O_CREATE | os.O_RDWR, 0777)
	LogErr(err, "Error creating new timeline in NewTimeLine()")
	defer newBranch.Close()

	newBranch.Write([]byte(parent.Hash))
	fmt.Printf("new timeline made, have fun in %s", timeLine)

}


type State struct {
	Perm string
	Sha1Id string
	File string
}
// we coud basically check for all the files in the current directory firstly
// and then if they exist, we go look for their hash and content 
// overwrite the current state with the content from the object by trunc && copying
// and for files that dont exist, we can just remake them according to Perm

func TravelTo(hash string) {
	s := time.Now()
	// go to object path to look for the first two letters of the hash passed in 
	folderName := hash[:2]
	fileName := hash[2:]
	fmt.Println("locating multiverse for commt....")
	commitLocation := fmt.Sprintf("%s/%s/%s", OBJECT_PATH, folderName, fileName)
	content, err := os.ReadFile(commitLocation)
	LogErr(err, "Error in getting commit in multiverse, TravelTo()")

	fmt.Println("here is the file state from the past")
	_, treeInfo, _ := strings.Cut(string(content), "tree:{")
	fmt.Printf("\nthis is da tree Id -> %s\n", treeInfo[:40])

	treeHash := treeInfo[:40]
	fmt.Println("locating stage area as that then...")
	
	treeFolder := treeHash[:2] 
	treeName := treeHash[2:]
	treePath := fmt.Sprintf("%s/%s/%s", OBJECT_PATH, treeFolder, treeName) 

	dirSnapshot, err := os.OpenFile(treePath, os.O_RDONLY, 0)
	LogErr(err, "Error in getting directorySnapshot, TravelTo()")
	defer dirSnapshot.Close()

	// fmt.Printf("\nDirectory snapshot to be read into struct -> \n%s\n", string(dirSnapshot))
	d := time.Since(s)
	fmt.Printf("function took %s to execute\n", d)

	scanner := bufio.NewScanner(dirSnapshot)

	for scanner.Scan() {
		state := State{}
		line := scanner.Text()
		fmt.Println(line)
		// the file returns content stoed as 0664  {40digithash} %20sfileName
		perm, content, _ := strings.Cut(line, "{")
		fileHash, fileName, _ := strings.Cut(content, "}")

		state.Perm = perm
		state.File = strings.TrimSpace(fileName)
		state.Sha1Id = fileHash
		
		fmt.Printf("add to struct for this line %s | %s | %s->\n", perm,  fileHash, strings.TrimSpace(fileName))

		BuildState(state)
	}
}


func BuildState(state State) {
	// just regular file path build up 
	folderLocation := state.Sha1Id[:2]
	fileHash := state.Sha1Id[2:]

	pathToF := fmt.Sprintf("%s/%s/%s", OBJECT_PATH, folderLocation, fileHash)
	
	snapshot, err := os.Open(pathToF)
	LogErr(err, "Error in opening snapshot file, BuildState")
	defer snapshot.Close()


	// tame to open file in current directory
	dst, err := os.OpenFile(state.File, os.O_CREATE | os.O_TRUNC, 0)
	LogErr(err, "Error in opening file in current directory, BuildState()")
	defer dst.Close()

	io.Copy(dst, snapshot)

	fmt.Printf("done copying prev state -> %s-> to %s\n ", state.Sha1Id, state.File)
}
