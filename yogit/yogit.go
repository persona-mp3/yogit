package yogit

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	// "encoding/json"
	"crypto/sha1"
)

type Commit struct {
	Parent *HashId `json:"parent"`
	Id HashId `json:"hashId"`
	Author string `json:"author"`
	Contact string `json:"contact"`
	CommitMsg string `json:"commitMsg"`
	CommittedAt time.Time `json:"committedAt"`
}

type HashId struct {
	Id string
}

func YoGit() {
	fmt.Println("WELCOME TO YOGIT")
}

func Init() {
	fmt.Println("standback, making folders")
	// use an input package to set the global variables
	err := os.MkdirAll(".yogit", 0777)
	LogErr(err, "Error in making yogit folders")

	subFolders := []string{"object", "log", "refs"}
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

func Add() {
	fmt.Println("adding all files onto the staging area")
	// write all current files to staging area : index, but for now, lets just add one
}

func CommitFunc(message string) {
	// parent := ParentCommit()
	// -> we first need to check if all the folders have been properly initialised and theres a staging area
	// for now lets just take the hash of commitMsg
	hasher := sha1.New()
	hasher.Write([]byte(message))
	hashContent := hasher.Sum(nil)
	hashId := hex.EncodeToString(hashContent)

	hash := HashId {
		Id: hashId,
	}

	commit := Commit {
		// Author -> read from globals
	}

	commit.CommitMsg = message
	commit.CommittedAt = time.Now()
	commit.Id = hash
	commit.Author = "persona-mp3@github.com"
	commit.Contact = "persona-mp3@github.com"

	fmt.Println(" yo! your new commit has been saved")

	fmt.Printf(" Commit Message: %-20s\n Commit Id: %s\n CommittedAt: %s\n", commit.CommitMsg, commit.Id.Id, commit.CommittedAt.Format("Jan 2, 2006, 3:04 PM"))

	logFile, err := os.OpenFile(".yogit/log/logs", os.O_CREATE | os.O_RDWR | os.O_APPEND, 0777)
	LogErr(err, "Error opening log file to record new commit")
	defer logFile.Close()
	
	logger := log.New(logFile, "", 0)
	logger.Printf("%-30s  %-30s  %-20s  %-40s %s\n",commit.Author, commit.Contact, hash.Id, message, time.Now().Format("Jan 2, 2006, 3:04 PM") ) 

	updateBranch(hash)

	fmt.Printf("\njust recorded your commit to ./.yogit/log/logs\n")
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
