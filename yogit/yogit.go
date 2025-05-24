package yogit

import (
	"fmt"
)

func YoGit() {
	fmt.Println("YOGIT initiated")

}

func Init() {
	fmt.Println("standback, making folders ./.yogit")
}

func Add() {
	fmt.Println("adding all files onto the staging area")
}

func Commit(message string) {
	fmt.Println("saving state message ----> ", message)
}
