package yogit

import (
	"fmt"
	"log"
)


func LogErr(err error, msg string) {
	if err != nil {
		fmt.Println("\n", msg)
		log.Fatal(err)
	}
}
